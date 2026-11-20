"""Translations of small legacy Go nodes into nodes RED runs.

Each function takes the legacy node dict and returns the replacement node
dict. The replacement keeps the node's name, so links stay as they are: a
standalone ottl_transform sends an item whose statement fails down the
node's `failure` links, or drops it when there are none, which is what the
Go agent does with an item a legacy node stops on.

  parse_json_attributes             -> parse_json_attributes(n, accept)  ottl_transform
  log_transform, generic_transform  -> transform(n, accept)              ottl_transform
  regex_filter                      -> regex_filter(n)                   ottl_filter
  mask                              -> mask(n, accept)                   ottl_transform
  unescape_json                     -> unescape_json(n, accept)          ottl_transform
  grok                              -> grok(n, accept)                   ottl_transform
  route (CEL)                       -> route(n)                          route_ottl
  log_to_pattern                    -> log_to_pattern(n)                 sequence [log_to_pattern_metric]

RULES maps each node type to its function. Every form outside the ones run
side by side against the Go agent raises translate.Unsupported. Where the
two agents are known to differ on some inputs, the function raises unless
the caller names that difference in `accept` (see DIFFERENCES).

Differences that are not behind `accept`:
  regex_filter    a body that is not a string: Go matches its JSON text with
                  <, > and & written as \\u003c, \\u003e, \\u0026; RED matches
                  them as written.
  log_to_pattern  Go emits cluster_pattern and cluster_sample items, the
                  replacement emits ed.pipeline.l2m.log_pattern metrics.
  transform       value forms other than regex_capture and convert_timestamp
                  are written by translate.value_statements.
"""
import re

try:
    import re._parser as sre_parse
except ImportError:  # Python < 3.11
    import sre_parse

import translate
from translate import Parser, Unsupported, capture_groups, check_keys, described, ottl_path, q, unsupported

DIFFERENCES = {
    "json-parser":
        "parse_json_attributes: Go's JSON reader takes text RED's ParseJSON rejects (numbers written 01, +1, .5 "
        "or 1., NaN, a raw tab or an unknown escape inside a string): Go parses it, RED sends the item down the "
        "failure path. An integer outside the int64 range becomes a rounded float64 on Go and stays exact on RED.",
    "mask-first-occurrence":
        "mask with a capture group: Go replaces the first occurrence of the captured text anywhere in the body, "
        "RED replaces the captured position. They differ when the captured text also occurs earlier in the body.",
    "grok-unquote":
        "grok: Go first reads the body as a Go string literal (\\t, \\n, \\\\, \\uXXXX ...) when it has a backslash, "
        "no newline and no bare double quote, and matches the unescaped text; RED matches the body as it is.",
    "grok-number-forms":
        "grok: Go turns a captured value into a number whenever strconv.ParseFloat reads it; RED converts integers and "
        "decimal numbers (up to 100 digits each side, exponent up to 2 digits) and leaves inf, infinity, nan, hex "
        "floats and longer numbers as strings.",
    "timestamp-unparseable":
        "convert_timestamp into item[\"timestamp\"]: when the value does not parse Go writes timestamp \"\" "
        "(not an integer); RED leaves the timestamp as it was. RED also leaves it for years outside 1700-2199, "
        "which Go converts.",
    "unescape-text":
        "unescape_json: the rewritten string is the same JSON value on both agents, but Go keeps the key order "
        "and number text of the input (1.50, 1e3) and RED writes keys sorted and numbers re-formatted.",
}

FAIL = 'set(cache["fail"], ParseJSON("{"))'  # a statement that always errors: the item leaves on the failure path
# cache["b"] = the text Go reads from the body: the string itself, or the JSON of any other value with
# Go's escaping of <, >, &, U+2028 and U+2029.
BODY_TEXT = ['set(cache["b"], body) where IsString(body)',
             'set(cache["b"], String(body)) where not IsString(body)'] + [
    f'replace_pattern(cache["b"], {q(c)}, {q(e)}) where not IsString(body)'
    for c, e in (("<", "\\u003c"), (">", "\\u003e"), ("&", "\\u0026"), ("\u2028", "\\u2028"), ("\u2029", "\\u2029"))]
KEY = re.compile(r"[^\\\"'`\[\]]+")
DOTTED = re.compile(r"[A-Za-z_][A-Za-z0-9_]*(\.[A-Za-z_][A-Za-z0-9_]*)+")


def need(node, accept, diff):
    if diff not in accept:
        raise Unsupported(f"{node}: unsupported without accept={diff!r}: {DIFFERENCES[diff]}")


def go_path(node, p):
    """Returns the keys a Go field path addresses: item["a"]["b"], item.a.b or a.b."""
    s = (p or "").strip()
    if s == "item" or s.startswith(("item.", "item[")):
        e = Parser(node, "field path", s).parse()
        keys = e[1] if e[0] == "path" else None
    elif DOTTED.fullmatch(s):
        keys = s.split(".")  # Go splits a path without the item prefix on dots
    else:
        keys = None
    if not keys or not all(KEY.fullmatch(k) for k in keys):
        raise unsupported(node, "field path", p)
    return keys


def attr_path(node, what, p):
    keys = go_path(node, p)
    if keys[0] != "attributes" or len(keys) < 2:
        raise unsupported(node, what, p)
    return keys


def present(keys):
    """Returns the OTTL condition that the last key exists in its parent map, null value or not."""
    parent = ottl_path(keys[:-1])
    return f"IsMap({parent}) and ContainsValue(Keys({parent}), {q(keys[-1])})"


def ottl(n, stmts):
    return described(n, {"name": n["name"], "type": "ottl_transform", "statements": "\n".join(stmts)})


def parse_json_attributes(n, accept=()):
    check_keys(n, ("process_field", "field_path"))
    name, pf, fp = n["name"], n.get("process_field") or "", n.get("field_path") or ""
    need(name, accept, "json-parser")
    keys = go_path(name, pf) if pf else ["body"]
    if keys != ["body"]:
        attr_path(name, "process_field", pf)
    src = ottl_path(keys)
    if not fp:
        tail = 'merge_maps(attributes, cache["pja"], "upsert")'
    else:
        dst = go_path(name, fp)
        if dst == ["attributes"]:
            tail = 'set(attributes, cache["pja"])'
        elif dst == keys and keys != ["body"]:
            tail = f'set({src}, cache["pja"])'
        else:
            raise unsupported(name, "field_path", fp)
    # Go stops on anything but a non-empty JSON object; a map is marshalled and read back.
    return ottl(n, [f'set(cache["pja"], ParseJSON({src})) where IsString({src})',
                    f'set(cache["pja"], {src}) where IsMap({src})',
                    f'{FAIL} where not IsMap(cache["pja"])',
                    f'{FAIL} where Len(cache["pja"]) == 0',
                    tail])


# convert_timestamp(item["attributes"][...], "2006-01-02T15:04:05.999999Z", "Unix Milli")
TS_LAYOUT = "2006-01-02T15:04:05.999999Z"
YEAR, LEAP = r"(?:1[7-9]|2[01])\d{2}", r"(?:(?:1[7-9]|2[01])(?:0[48]|[2468][048]|[13579][26])|2000)"
# What time.Parse reads with that layout, for the years 1700-2199 (RED's Time holds nanoseconds in an
# int64): a valid date and time, an optional fraction, Z.
TS_RE = ("^(?P<s>(?:" + YEAR + r"-(?:(?:0[13578]|1[02])-(?:0[1-9]|[12]\d|3[01])|(?:0[469]|11)-(?:0[1-9]|[12]\d|30)"
         r"|02-(?:0[1-9]|1\d|2[0-8]))|" + LEAP + r"-02-29)"
         r"T(?:[01]\d|2[0-3]):[0-5]\d:[0-5]\d)(?:[.,](?P<f>\d+))?Z$")


def timestamp_statements(node, t, accept):
    e = Parser(node, "CEL value", t["value"]).parse()
    if not (e[0] == "call" and e[1] == "convert_timestamp" and len(e[2]) == 3 and translate.is_ref(e[2][0])
            and e[2][1] == ("str", TS_LAYOUT) and e[2][2] == ("str", "Unix Milli")) or t.get("ignore_if_empty"):
        raise unsupported(node, "CEL value", " ".join(t["value"].split()))
    need(node, accept, "timestamp-unparseable")
    src = ottl_path(e[2][0][1])
    ok = f'IsString({src}) and IsMatch({src}, {q(TS_RE)})'
    return [f'set(cache["ts"], ExtractPatterns({src}, {q(TS_RE)})) where {ok}',
            'set(cache["ms"], 0)',
            'set(cache["ms"], Int(Substring(Concat([cache["ts"]["f"], "00"], ""), 0, 3))) '
            'where cache["ts"]["f"] != nil and cache["ts"]["f"] != ""',
            'set(timestamp, UnixMilli(Time(cache["ts"]["s"], "%Y-%m-%dT%H:%M:%S")) + cache["ms"]) '
            'where cache["ts"]["s"] != nil']


INDEXED = re.compile(r'^(.*\))\s*\[\s*"([A-Za-z_][A-Za-z0-9_]*)"\s*\]\s*$', re.S)


def transform(n, accept=()):
    check_keys(n, ("transformations",))
    name, stmts = n["name"], []
    for t in n.get("transformations") or []:
        for k in t:
            if k not in ("field_path", "operation", "value", "ignore_if_empty"):
                raise unsupported(name, "transformation key", k)
        keys, op = go_path(name, t["field_path"]), t.get("operation", "upsert")
        if op == "upsert" and keys == ["timestamp"]:
            stmts += timestamp_statements(name, t, accept)
            continue
        # Go rejects a log_transform that writes a resource field.
        if keys[0] not in (("attributes", "resource") if n["type"] == "generic_transform" else ("attributes",)) \
                or len(keys) < 2:
            raise unsupported(name, "field path", t["field_path"])
        if len(keys) != 2:
            # Deeper paths: Go deletes element 0 of a list parent and builds missing parents on upsert.
            raise unsupported(name, "nested field path", t["field_path"])
        if op == "delete":
            stmts.append(f"delete_key({keys[0]}, {q(keys[1])})")
        elif op == "upsert":
            v = t["value"]
            m = INDEXED.match(v)  # regex_capture(...)["name"] reads the same key as regex_capture(...).name
            if m:
                v = f"{m.group(1)}.{m.group(2)}"
            dst = ottl_path(keys)
            stmts += translate.value_statements(name, dst, v, t.get("ignore_if_empty", False))
            e = Parser(name, "CEL value", v).parse()
            ref = e if translate.is_ref(e) else e[2] if e[0] == "cond" and e[1][0] == "in" else None
            if ref:
                # Go reads an explicit null as the number 0.
                stmts.append(f"set({dst}, 0) where {present(ref[1])} and {ottl_path(ref[1])} == nil")
        else:
            raise unsupported(name, "operation", op)
    return ottl(n, stmts)


log_transform = generic_transform = transform


def regex_filter(n):
    check_keys(n, ("pattern", "negate"))
    if not n.get("pattern"):
        raise unsupported(n["name"], "regex_filter pattern", '""')
    # Go matches the JSON text of a body that is not a string.
    return described(n, {"name": n["name"], "type": "ottl_filter",
                         "filter_mode": "exclude" if n.get("negate") else "include",
                         "condition": f"IsMatch(body, {q(n['pattern'])})"})


def consumable(node, pat):
    """Returns the character ranges pat can consume, or None for any character.

    Raises for a pattern with anchors, word boundaries or an empty match, and
    for syntax Python's parser reads differently from Go's.
    """
    bad = unsupported(node, "mask pattern", pat)
    if re.search(r"\[\[:|\\[pPQECzAbB]|\(\?[a-zA-Z-]*U|\{,", pat):
        raise bad
    try:
        tree = sre_parse.parse(pat)
    except Exception:
        raise bad
    if tree.getwidth()[0] == 0:
        raise bad
    ranges, anything = [], [False]
    cats = {"CATEGORY_DIGIT": [(48, 57)], "CATEGORY_WORD": [(48, 57), (65, 90), (97, 122), (95, 95)],
            "CATEGORY_SPACE": [(9, 13), (32, 32)]}

    def walk(items):
        for op, av in items:
            op = str(op)
            if op == "LITERAL":
                ranges.append((av, av))
            elif op in ("NOT_LITERAL", "ANY"):
                anything[0] = True
            elif op == "IN":
                for o, a in av:
                    o = str(o)
                    if o == "LITERAL":
                        ranges.append((a, a))
                    elif o == "RANGE":
                        ranges.append(a)
                    elif o == "CATEGORY" and str(a) in cats:
                        ranges.extend(cats[str(a)])
                    else:
                        anything[0] = True
            elif op == "BRANCH":
                for b in av[1]:
                    walk(b)
            elif op == "SUBPATTERN":
                walk(av[-1])
            elif op in ("MAX_REPEAT", "MIN_REPEAT"):
                walk(av[2])
            else:  # anchors, boundaries, lookarounds, backreferences
                raise bad
    walk(tree)
    return None if anything[0] else ranges


def group_span(pat, start, end):
    """Returns the index of the parenthesis closing the group opened at pat[start:end], and its depth."""
    depth, i, in_class, opened = 0, 0, False, None
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
            if i == start:
                opened = depth
            depth += 1
        elif c == ")":
            depth -= 1
            if opened is not None and depth == opened:
                return i, opened
        elif c == "|" and depth == 0:
            return None, None
        i += 1
    return None, None


def mask(n, accept=()):
    check_keys(n, ("pattern", "mask"))
    name, pat = n["name"], n.get("pattern") or ""
    if not pat:
        raise unsupported(name, "mask pattern", '""')
    text = "******" if n.get("mask") is None else n["mask"]
    if not isinstance(text, str):
        raise unsupported(name, "mask", text)
    groups, repl = capture_groups(pat), text.replace("$", "$$")
    # Go writes the body back as a string even when nothing matched.
    head = BODY_TEXT + ['set(body, cache["b"]) where not IsString(body)']
    if not groups:
        # Go replaces the first occurrence of each matched text. That is the match itself when the
        # pattern has no anchors and can never match text holding a character of the mask.
        can = consumable(name, pat)
        if "(?" in pat and not text.isascii():
            raise unsupported(name, f"mask {text!r} with a flagged pattern", pat)
        for ch in set(text.lower() + text.upper()):
            if can is None or any(lo <= ord(ch) <= hi for lo, hi in can):
                raise unsupported(name, f"mask {text!r} that the pattern could match", pat)
        return ottl(n, head + [f"replace_pattern(body, {q(pat)}, {q(repl)})"])
    # Go masks the first capture group only.
    need(name, accept, "mask-first-occurrence")
    flags = re.match(r"\(\?[ims]+\)", pat)
    rest = pat[flags.end():] if flags else pat
    if re.search(r"\(\?[a-zA-Z-]+\)", rest) or "edpre" in pat or "edpost" in pat:
        raise unsupported(name, "mask pattern", pat)
    (start, end) = capture_groups(rest)[0]
    close, depth = group_span(rest, start, end)
    if close is None or depth != 0 or "|" in group_top(rest) or rest[close + 1:close + 2] in ("*", "+", "?", "{"):
        raise unsupported(name, "mask pattern whose first group is not a plain top-level group", pat)
    body = rest[end:close]
    try:
        if sre_parse.parse(body).getwidth()[0] == 0:
            raise unsupported(name, "mask pattern whose first group can be empty", pat)
    except Unsupported:
        raise
    except Exception:
        raise unsupported(name, "mask pattern", pat)
    new = f"{flags.group(0) if flags else ''}(?P<edpre>{rest[:start]})(?:{body})(?P<edpost>{rest[close + 1:]})"
    return ottl(n, head + [f"replace_pattern(body, {q(new)}, {q('${edpre}' + repl + '${edpost}')})"])


def group_top(pat):
    """Returns pat with every parenthesized group and character class blanked out."""
    out, depth, i, in_class = [], 0, 0, False
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
            depth += 1
        elif c == ")":
            depth -= 1
        elif depth == 0:
            out.append(c)
        i += 1
    return "".join(out)


def unescape_json(n, accept=()):
    check_keys(n, ("process_field", "field_path"))
    name, pf, fp = n["name"], n.get("process_field") or "", n.get("field_path") or ""
    need(name, accept, "unescape-text")
    src = go_path(name, pf) if pf else ["body"]
    dst = go_path(name, fp) if fp else src
    for what, keys, p in (("process_field", src, pf), ("field_path", dst, fp)):
        if keys != ["body"]:
            attr_path(name, what, p)
    if dst != src and len(dst) > 2:
        raise unsupported(name, "field_path", fp)
    s = ottl_path(src)
    stmts = [f'set(cache["src"], {s}) where IsString({s})',
             f'set(cache["src"], String({s})) where {s} != nil and not IsString({s})']
    if src != ["body"]:
        # Go reads an explicit null as the text 0.
        stmts.append(f'set(cache["src"], "0") where {present(src)} and {s} == nil')
    # Go stops unless the text is one JSON value; "[text]" parsing to one element says the same.
    return ottl(n, stmts + [
        f'{FAIL} where cache["src"] == nil',
        'set(cache["chk"], ParseJSON(Concat(["[", cache["src"], "]"], "")))',
        f'{FAIL} where Len(cache["chk"]) != 1',
        f'set({ottl_path(dst)}, EDXUnescapeJSON(cache["src"]))'])


GROK_LEAVES = ("NUMBER", "NOTSPACE", "TIMESTAMP_ISO8601", "INT", "GREEDYDATA", "SPACE", "DATA")
GROK_REF = re.compile(r"%\{([^}]*)\}")
INT64 = "9223372036854775807"


def at_most(n):
    """Returns a regex for the decimal strings with len(n) digits, no leading zero, not above n."""
    alts = []
    for i, d in enumerate(n):
        lo, rest = (1 if i == 0 else 0), len(n) - i - 1
        if int(d) > lo:
            alts.append(f"{n[:i]}[{lo}-{int(d) - 1}]" + (f"\\d{{{rest}}}" if rest else ""))
        elif int(d) == 1 and lo == 0:
            alts.append(f"{n[:i]}0" + (f"\\d{{{rest}}}" if rest else ""))
    alts.append(n)
    return "(?:" + "|".join(alts) + ")"


# What strconv.Atoi reads: an optional sign and digits inside the int64 range.
GO_INT = f"^(?:[+-]?0*(?:\\d{{1,18}}|{at_most(INT64)})|-0*9223372036854775808)$"
# The finite decimal numbers strconv.ParseFloat reads (digits may be separated by single underscores),
# bounded to 100 digits a side and a 2-digit exponent so they cannot overflow.
DIGITS = "(?:\\d_?){0,99}\\d"
GO_FLOAT = f"^[+-]?(?:{DIGITS}(?:\\.(?:{DIGITS})?)?|\\.{DIGITS})(?:[eE][+-]?\\d(?:_?\\d)?)?$"
GO_TRUE, GO_FALSE = "^(?:t|T|TRUE|true|True)$", "^(?:f|F|FALSE|false|False)$"


def grok(n, accept=()):
    check_keys(n, ("pattern", "custom_pattern"))
    name = n["name"]
    pat = n.get("custom_pattern") or n.get("pattern") or ""
    names = []
    for ref in GROK_REF.findall(pat):
        parts = ref.split(":")
        if parts[0] not in GROK_LEAVES or len(parts) > 2:
            raise unsupported(name, "grok reference", "%{" + ref + "}")
        if len(parts) == 2:
            names.append(parts[1])
    if re.search(r"\(\?<?[=!]", pat):
        raise unsupported(name, "grok pattern", pat)
    names += [m.group(1) for m in re.finditer(r"\(\?P?<([^>]*)>", pat)]
    if not names or len(set(names)) != len(names) or not all(re.fullmatch(r"[A-Za-z][A-Za-z0-9_]*", x) for x in names) \
            or any("__" in x for x in names):
        raise unsupported(name, "grok capture names", pat)
    need(name, accept, "grok-unquote")
    need(name, accept, "grok-number-forms")
    stmts = BODY_TEXT + [f'set(cache["g"], ExtractGrokPatterns(cache["b"], {q(pat)}, true))']
    for k in names:
        v, a = f'cache["g"][{q(k)}]', f"attributes[{q(k)}]"
        # Go stores each capture as an int, else a float, else a bool, else the string.
        stmts += [f"set({a}, {v}) where {v} != nil",
                  f"set({a}, Int({v})) where {v} != nil and IsMatch({v}, {q(GO_INT)})",
                  'delete_key(cache, "n")',
                  f'set(cache["n"], {v}) where {v} != nil and IsMatch({v}, {q(GO_FLOAT)}) and not IsMatch({v}, {q(GO_INT)})',
                  'replace_pattern(cache["n"], "_", "")',
                  f'set({a}, Double(cache["n"])) where cache["n"] != nil',
                  f"set({a}, true) where {v} != nil and IsMatch({v}, {q(GO_TRUE)})",
                  f"set({a}, false) where {v} != nil and IsMatch({v}, {q(GO_FALSE)})"]
    return ottl(n, stmts)


INT_CMP = re.compile(r"^(.*?)\s*==\s*(-?\d+)$", re.S)


def is_attr(e):
    return e[0] == "path" and len(e[1]) >= 2 and e[1][0] == "attributes" and all(KEY.fullmatch(k) for k in e[1])


def route_condition(node, src):
    """Returns the OTTL condition for one CEL route condition; a CEL error counts as no match."""
    bad = unsupported(node, "route condition", src)
    m = INT_CMP.match(src)
    e = Parser(node, "route condition", m.group(1) if m else src).parse()
    if m:
        # CEL compares a number with an int literal across int and double; any other type is unequal.
        if not is_attr(e):
            raise bad
        return f"{ottl_path(e[1])} == {m.group(2)}"
    if e[0] in ("eq", "ne") and is_attr(e[1]) and e[2][0] == "str":
        p, lit = ottl_path(e[1][1]), q(e[2][1])
        # != holds for every present value that is not the literal, an explicit null included.
        return f"{p} == {lit}" if e[0] == "eq" else f"{present(e[1][1])} and {p} != {lit}"
    if e[0] == "call" and e[1] == "regex_match" and len(e[2]) == 2 and e[2][0][0] == "path" and e[2][1][0] == "str" \
            and (e[2][0][1] == ["body"] or is_attr(e[2][0])):
        # regex_match errors on anything but a string.
        p = ottl_path(e[2][0][1])
        return f"IsString({p}) and IsMatch({p}, {q(e[2][1][1])})"
    raise bad


def route(n):
    check_keys(n, ("expression_type", "paths"))
    if n.get("expression_type", "cel") != "cel":
        raise unsupported(n["name"], "expression_type", n.get("expression_type"))
    paths = []
    for p in n.get("paths") or []:
        for k in p:
            if k not in ("path", "condition", "exit_if_matched"):
                raise unsupported(n["name"], "route path key", k)
        cond = route_condition(n["name"], " ".join(p["condition"].split()))
        paths.append({"path": p["path"], "condition": cond, "exit_if_matched": p.get("exit_if_matched", False)})
    return described(n, {"name": n["name"], "type": "route_ottl", "paths": paths})


def log_to_pattern(n):
    check_keys(n, ("num_of_clusters", "samples_per_cluster", "reporting_frequency"))
    out = translate.log_to_pattern(n)
    # The control plane accepts log_to_pattern_metric only inside a sequence.
    proc = {k: v for k, v in out.items() if k not in ("name", "user_description")}
    proc.setdefault("data_types", ["log"])
    return described(n, {"name": n["name"], "type": "sequence", "processors": [proc]})


RULES = {"parse_json_attributes": parse_json_attributes, "log_transform": transform, "generic_transform": transform,
         "regex_filter": regex_filter, "mask": mask, "unescape_json": unescape_json, "grok": grok, "route": route,
         "log_to_pattern": log_to_pattern}
