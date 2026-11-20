#!/usr/bin/env python3
"""Translate the Go agent's legacy `log_to_metric` node into OTTL nodes RED runs.

`log_to_metric(n)` takes the legacy node dict and returns a `sequence` node:

  ottl_filter       keep logs whose body matches `pattern`
  ottl_transform    (capture-group dimensions only) put the named captures in attributes["__l2m"]
  extract_metric    one gauge per dimension group and log: 1, or the numeric dimension
  ottl_filter       drop every log (the Go node emits metrics only)
  ottl_transform    per group: dimension values become string attributes named after the last
                    path key; the attributes Go never copies (the numeric field, the maps the
                    nested paths live in) are removed
  aggregate_metric  one rule per group and stat, named `<metric>.<stat>`
  ottl_transform    `_stat_type` = the stat, as Go stamps it
  ottl_filter       (only_report_nonzeros) drop zero values

Every other form raises translate.Unsupported; nothing is guessed.
`approx_percentiles=True` accepts p95/p99, whose values differ from Go's (see PERCENTILES).
`other_attributes="drop"` keeps only the dimension attributes on the metrics.

Each rule was checked with sidebyside.py (Go master build against RED 3.16.0-rc.10). What
still differs from the Go node after translation:
  * no `<metric>.anomaly1`, which Go emits when no enabled_stats are set;
  * unless only_report_nonzeros is true, Go goes on reporting 0 for every stat of a series
    in the intervals without logs (until the series has been idle for 5 minutes, per the Go
    source); RED reports nothing for an interval without logs;
  * Go's interval starts when a series first appears, follows the wall clock and stamps the
    report time; RED's is aligned to the epoch, follows the log timestamp and stamps the end
    of the interval, so the same logs can fall into different intervals;
  * attributes other than the dimensions: Go stamps a series with the string attributes of
    the first log that created it; RED keeps the attributes, of any type, that are equal on
    every log of the series in that interval (or none with other_attributes="drop");
  * Go's `max` is never below 0;
  * a float dimension value of 1e15 or more, or below 1e-5, is printed differently;
  * Go folds the series beyond 1000 per dimension group into "other" (per the Go source);
    RED has no limit.
"""
import re

import translate
from translate import q, unsupported

TOP_KEYS = ("pattern", "interval", "metric_name", "enabled_stats", "dimension_groups",
            "skip_empty_intervals", "only_report_nonzeros")
GROUP_KEYS = ("field_dimensions", "field_numeric_dimension", "dimensions", "numeric_dimension",
              "custom_suffix", "enabled_stats", "interval", "skip_empty_intervals", "only_report_nonzeros")
# Go stat -> aggregate_metric aggregation_type.
STATS = {"count": "count", "sum": "sum", "avg": "mean", "min": "min", "max": "max"}
# Go interpolates at q*(n+1), RED takes the nearest rank at q*n: the values differ unless the
# series has one sample. With a percentile enabled Go also computes every stat of the group
# over the points of the last interval+2s (per the Go source), so a point of the last 2s
# before a report is counted again in the next one.
PERCENTILES = {"p95": "p95", "p99": "p99"}
ORDER = ("count", "sum", "avg", "min", "max", "p95", "p99")
# Go's default interval (core.DefaultMetricRuleInterval).
DEFAULT_INTERVAL = "1m"
FIELD_PATH = re.compile(r'item(\["[^"\\]*"\])+')
# Capture group bodies whose every match Go's strconv.ParseFloat and RED's Double read alike.
NUMERIC_CAPTURES = (r"\d+", r"[0-9]+")
# Go keeps one series per source (its group name). RED 3.16.0-rc.10 leaves __group_name empty
# and names the source in __src_name, so both are grouped on.
SOURCE = ['resource["__group_name"]', 'resource["__src_name"]']
CAPTURES = 'attributes["__l2m"]'
NAME_RE = re.compile(r"[^_a-zA-Z0-9.]+")


def named_groups(pat):
    """Maps each named capture group of pat to (body, mandatory).

    mandatory: the group is at the top level of a pattern without alternation and has
    no quantifier, so it takes part in every match.
    """
    out, stack, i, in_class, depth_alt = {}, [], 0, False, "|" in re.sub(r"\\.|\[(?:\\.|[^\]])*\]", "", pat)
    while i < len(pat):
        c = pat[i]
        if c == "\\":
            i += 2
            continue
        if in_class:
            in_class = c != "]"
        elif c == "[":
            in_class = True
            if pat[i + 1:i + 2] == "^":
                i += 1
            if pat[i + 1:i + 2] == "]":
                i += 1
        elif c == "(":
            m = re.match(r"\(\?P<(\w+)>", pat[i:])
            stack.append((m.group(1) if m else None, i + (m.end() if m else 1), len(stack)))
        elif c == ")" and stack:
            name, start, depth = stack.pop()
            if name is not None:
                quantified = pat[i + 1:i + 2] in ("?", "*", "{", "+")
                out[name] = (pat[start:i], depth == 0 and not depth_alt and not quantified)
        i += 1
    return out


def red_pattern(pat):
    """Drops `(?)`: an empty flag group, which Go's regexp accepts as a no-op and RED's rejects."""
    out, i, in_class = [], 0, False
    while i < len(pat):
        c = pat[i]
        if c == "\\":
            out.append(pat[i:i + 2])
            i += 2
            continue
        if in_class:
            in_class = c != "]"
        elif c == "[":
            in_class = True
            if pat[i + 1:i + 2] == "^":
                out.append(c)
                i, c = i + 1, "^"
            if pat[i + 1:i + 2] == "]":
                out.append(c)
                i, c = i + 1, "]"
        elif pat.startswith("(?)", i):
            i += 3
            continue
        out.append(c)
        i += 1
    return "".join(out)


def field_keys(node, p):
    """Returns the attribute keys of a Go `item["attributes"]["a"]["b"]` path."""
    p = p.strip()
    if not FIELD_PATH.fullmatch(p):
        raise unsupported(node, "dimension field path", p)
    keys = re.findall(r'\["([^"\\]*)"\]', p)
    if keys[0] != "attributes" or len(keys) < 2 or any(k == "" for k in keys):
        raise unsupported(node, "dimension field path", p)
    return keys[1:]


def attr(keys):
    return "attributes" + "".join(f"[{q(k)}]" for k in keys)


def maps(keys):
    """The condition that every map on the way to attr(keys) is a map, with a trailing ` and `.

    Indexing a string or a number is an OTTL error in RED, and an error stops the item.
    """
    return "".join(f"IsMap({attr(keys[:i])}) and " for i in range(1, len(keys)))


def flag(node, n, g, key):
    """The top-level flag; a group may repeat it but not change it."""
    top = n.get(key) or False
    if not isinstance(top, bool) or (key in g and (g[key] or False) != top):
        raise unsupported(node, f"dimension group {key}", g.get(key, top))
    return top


def log_to_metric(n, approx_percentiles=False, other_attributes="bucket"):
    node = n.get("name") or n.get("metric_name")
    if not node:
        raise unsupported("log_to_metric", "node without name and metric_name", q(n.get("pattern")))
    if other_attributes not in ("bucket", "drop"):
        raise ValueError(other_attributes)
    translate.check_keys({"name": node, **n}, TOP_KEYS)
    pat = n.get("pattern")
    if not pat or not isinstance(pat, str):
        raise unsupported(node, "log_to_metric pattern", q(pat))
    interval = n.get("interval") or DEFAULT_INTERVAL
    if not isinstance(interval, str):
        raise unsupported(node, "interval", interval)
    base = n.get("metric_name") or node
    pat = red_pattern(pat)
    if not pat:
        raise unsupported(node, "log_to_metric pattern", q(n["pattern"]))
    captures = named_groups(pat)
    groups = n.get("dimension_groups") or [{}]
    nonzero = None
    rules, shape, aggs, names, stats_used, use_captures = [], [], [], set(), set(), False

    for g in groups:
        g = g or {}
        for k in g:
            if k not in GROUP_KEYS:
                raise unsupported(node, "dimension group key", k)
        if (g.get("interval") or interval) != interval:
            raise unsupported(node, "dimension group interval", g["interval"])
        flag(node, n, g, "skip_empty_intervals")
        nonzero = flag(node, n, g, "only_report_nonzeros")
        if g.get("field_numeric_dimension") and g.get("numeric_dimension"):
            raise unsupported(node, "dimension group with two numeric dimensions", g["numeric_dimension"])

        # Metric name, as Go builds it: <name>[_<numeric dimension>][_<custom_suffix>].
        name, value, conditions, numeric_keys = base, "1", [], None
        if g.get("numeric_dimension"):
            d = g["numeric_dimension"]
            if d not in captures or not captures[d][1] or captures[d][0] not in NUMERIC_CAPTURES:
                raise unsupported(node, "numeric_dimension capture group", d)
            name, value, use_captures = f"{name}_{d}", f"Double({CAPTURES}[{q(d)}])", True
        elif g.get("field_numeric_dimension"):
            numeric_keys = field_keys(node, g["field_numeric_dimension"])
            name, value = f"{name}_{numeric_keys[-1]}", attr(numeric_keys)
            # Go takes int64 and float64 only; anything else skips the log for this group.
            conditions = [f"{maps(numeric_keys)}(IsInt({value}) or IsDouble({value}))"]
        if g.get("custom_suffix"):
            name = f"{name}_{g['custom_suffix']}"
        if "|" in name:
            raise unsupported(node, "metric name with |", name)
        name = NAME_RE.sub("_", name)
        if name in names:
            raise unsupported(node, "dimension groups with the same metric name", name)
        names.add(name)

        # Dimensions: (attribute name, source path, already a top-level attribute, present test, null test).
        dims = []
        for d in g.get("dimensions") or []:
            if d not in captures or not captures[d][1]:
                raise unsupported(node, "dimension capture group", d)
            dims.append((d, f"{CAPTURES}[{q(d)}]", False, f"{CAPTURES}[{q(d)}] != nil", None))
            use_captures = True
        fields = [field_keys(node, p) for p in g.get("field_dimensions") or []]
        for keys in fields:
            # A key that is present with a JSON null value.
            null = f"{maps(keys)}{attr(keys)} == nil and ContainsValue(Keys({attr(keys[:-1])}), {q(keys[-1])})"
            dims.append((keys[-1], attr(keys), len(keys) == 1, f"{maps(keys)}{attr(keys)} != nil", null))
        dim_names = [d[0] for d in dims]
        if len(set(dim_names)) != len(dim_names):
            raise unsupported(node, "dimensions with the same name", ", ".join(dim_names))
        nested = fields + ([numeric_keys] if numeric_keys else [])
        roots = sorted({k[0] for k in nested if len(k) > 1})
        drop = [numeric_keys[0]] if numeric_keys and len(numeric_keys) == 1 else []
        for k in roots + drop + ["__l2m"]:
            if k in dim_names:
                raise unsupported(node, "dimension named like an attribute the node reads", k)

        rules.append({"name": name, **({"conditions": conditions} if conditions else {}), "gauge": {"value": value}})
        w = f"name == {q(name)}"
        lifted = [(i, d, src, present, null) for i, (d, src, top, present, null) in enumerate(dims) if not top]
        for i, _, src, present, null in lifted:
            shape.append(f'set(cache["d{i}"], {src}) where {w} and {present}')
            if null:
                shape.append(f'set(cache["d{i}"], "0") where {w} and {null}')
        for k in roots:
            shape.append(f"delete_key(attributes, {q(k)}) where {w} and IsMap(attributes[{q(k)}])")
        for k in drop:
            shape.append(f"delete_key(attributes, {q(k)}) where {w}")
        for i, d, _, _, _ in lifted:
            # Go leaves a dimension out when its path is missing.
            shape.append(f"delete_key(attributes, {q(d)}) where {w}")
            shape.append(f'set(attributes[{q(d)}], cache["d{i}"]) where {w} and cache["d{i}"] != nil')
        for d, _, top, _, null in dims:
            if top:
                # Go reports "0" for a dimension whose value is JSON null.
                shape.append(f'set(attributes[{q(d)}], "0") where {w} and {null}')
        for d in dim_names:
            # Go prints a float without a fraction as an integer (7.0 -> "7").
            shape.append(f"set(attributes[{q(d)}], Int(attributes[{q(d)}])) where {w} and IsDouble(attributes[{q(d)}]) and "
                         f"attributes[{q(d)}] > -1000000000000000.0 and attributes[{q(d)}] < 1000000000000000.0 and "
                         f"Double(Int(attributes[{q(d)}])) == attributes[{q(d)}]")
            # Go reports every dimension value as a string (JSON for other types).
            shape.append(f"set(attributes[{q(d)}], String(attributes[{q(d)}])) where {w} and "
                         f"attributes[{q(d)}] != nil and not IsString(attributes[{q(d)}])")
        if other_attributes == "drop":
            shape.append(f"keep_keys(attributes, [{', '.join(q(d) for d in dim_names)}]) where {w}" if dim_names
                         else f'delete_matching_keys(attributes, ".*") where {w}')

        # Go unions the node's and the group's enabled_stats; empty means the defaults.
        stats = set(g.get("enabled_stats") or []) | set(n.get("enabled_stats") or [])
        if not stats:
            stats = {"min", "max", "avg", "count"} if value != "1" else {"count"}
        for s in sorted(stats, key=lambda s: ORDER.index(s) if s in ORDER else len(ORDER)):
            if s in PERCENTILES and approx_percentiles:
                agg = PERCENTILES[s]
            elif s in STATS:
                agg = STATS[s]
            else:
                raise unsupported(node, "enabled_stats", s)
            stats_used.add(s)
            aggs.append({"name": f"{name}.{s}", "conditions": [w], "interval": interval, "aggregation_type": agg,
                         "group_by": SOURCE + [f"attributes[{q(d)}]" for d in dim_names]})

    if use_captures and other_attributes != "drop":
        shape.append('delete_key(attributes, "__l2m")')
    procs = [{"type": "ottl_filter", "data_types": ["log"], "filter_mode": "include",
              "condition": f"IsMatch(body, {q(pat)})"}]
    if use_captures:
        procs.append({"type": "ottl_transform", "data_types": ["log"],
                      "statements": f"set({CAPTURES}, ExtractPatterns(body, {q(pat)}))"})
    procs.append({"type": "extract_metric", "data_types": ["log"], "keep_item": False, "extract_metric_rules": rules})
    procs.append({"type": "ottl_filter", "data_types": ["log"], "filter_mode": "exclude", "condition": "true"})
    if shape:
        procs.append({"type": "ottl_transform", "data_types": ["metric"], "statements": "\n".join(shape)})
    procs.append({"type": "aggregate_metric", "data_types": ["metric"], "aggregate_metric_rules": aggs})
    procs.append({"type": "ottl_transform", "data_types": ["metric"], "statements": "\n".join(
        f'set(_stat_type, {q(s)}) where IsMatch(name, {q(chr(92) + "." + s + "$")})'
        for s in ORDER if s in stats_used)})
    if nonzero and stats_used != {"count"}:
        # Go skips a stat whose value is 0; a count is never 0 here.
        procs.append({"type": "ottl_filter", "data_types": ["metric"], "filter_mode": "exclude",
                      "condition": 'gauge["value"] == 0'})
    return translate.described(n, {**({"name": n["name"]} if "name" in n else {}), "type": "sequence", "processors": procs})
