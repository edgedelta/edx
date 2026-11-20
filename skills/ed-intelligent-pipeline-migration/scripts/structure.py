#!/usr/bin/env python3
"""Structural rewrites of an Edge Delta v3 pipeline so that RED can run it.

Pure functions over the parsed config dict; they rewrite it in place and do no
I/O. Every shape outside the rules below raises `translate.Unsupported`, and
the config may then be partly rewritten, so work on a copy.

  lift_nested_packs(c)  a pack (`type: compound`) listed in a sequence's
                        `processors:` becomes a top-level pack node
  old_outputs(c)        ed_logs_output / ed_archive_output / ed_metrics_output /
                        ed_health_output / ed_patterns_output -> one ed_output

Run both BEFORE translate.py's node rewrites: `old_outputs` reads the legacy
node types to know which data types a node emits.
"""
import json

from translate import unsupported

# Attribute the Go agent stamps on every item leaving a pack that is a sequence
# processor; its value is the name of the compound_output the item left by
# (internalv3/discover/sequence_node.go:699, processors/compound/processor.go:190-222,
# semconv/semconv1/resource.go:133). A top-level pack does not stamp it
# (internalv3/discover/adaptor.go:2423-2431 passes no pathNameAttributeField).
PACK_PATH_ATTR = "ed.pipeline.pack.path.name"

GATE_PATH = "matched"
UNMATCHED = "unmatched"
ENTRY = ("\0entry", None)

PACK_KEYS = {"type", "name", "nodes", "links", "data_types", "condition", "metadata",
             "node_reference", "user_description", "disabled"}
SUBSEQ_KEYS = {"type", "name", "processors", "data_types", "condition", "metadata",
               "user_description", "disabled"}


def _active(p):
    return not p.get("disabled") and p.get("type") != "comment"


def _has_pack(procs):
    """True when an enabled pack sits in this processor list, at any depth."""
    return any(_active(p) and (p.get("type") == "compound"
                               or (p.get("type") == "sequence" and _has_pack(p.get("processors") or [])))
               for p in procs or [])


def _has_final(procs):
    return any(_active(p) and (p.get("final") or _has_final(p.get("processors") or [])) for p in procs or [])


def _gate_condition(p):
    """The OTTL condition Go puts on a sequence processor's node
    (sequence_node.go:309-327): data_types OR-ed, AND-ed with `condition`."""
    types = " or ".join(f"EDXDataType() == {json.dumps(t)}" for t in p.get("data_types") or [])
    cond = (p.get("condition") or "").strip()
    if types and cond:
        # The newline keeps a trailing `//` comment from swallowing the parenthesis.
        return f"({types}) and ({cond}\n)"
    return types or cond


class _Lift:
    """Builds the top-level nodes and links that replace one sequence node."""

    def __init__(self, seq, taken):
        self.seq, self.base, self.taken = seq, seq["name"], taken
        self.nodes, self.links, self.count = [], [], 0

    def new_name(self, kind):
        self.count += 1
        name = f"{self.base}__{kind}{self.count}"
        if name in self.taken:
            raise unsupported(self.base, "generated node name already in use", name)
        self.taken.add(name)
        return name

    def link(self, sources, to):
        for frm, path in sources:
            l = {"from": frm, "to": to}
            if path:
                l["path"] = path
            self.links.append(l)

    def flush(self, run, sources, depth, last):
        """Emit the plain processors collected so far as one sequence node."""
        if not run:
            return sources
        if _has_final(run) and not (depth == 0 and last):
            # A `final` processor sends its items to the end of the whole
            # top-level sequence (sequence_node.go:182-186, 189-196, 203-210),
            # past the pack and everything behind it.
            raise unsupported(self.base, "`final` processor ahead of a nested pack or inside a nested sequence that holds one",
                              [p.get("type") for p in run if _active(p) and (p.get("final") or _has_final(p.get("processors") or []))])
        if not self.nodes and sources == [ENTRY]:
            node = {k: v for k, v in self.seq.items() if k != "processors"}
        else:
            node = {"name": self.new_name("seq"), "type": "sequence"}
        node["processors"] = run
        self.nodes.append(node)
        self.link(sources, node["name"])
        return [(node["name"], None)]

    def gate(self, p, sources):
        """Route on the processor's gate; returns (sources into it, bypass sources)."""
        cond = _gate_condition(p)
        if not cond:
            return sources, []
        name = self.new_name("gate")
        self.nodes.append({"name": name, "type": "route_ottl",
                           "paths": [{"path": GATE_PATH, "condition": cond, "exit_if_matched": True}]})
        self.link(sources, name)
        return [(name, GATE_PATH)], [(name, UNMATCHED)]

    def pack(self, p, sources):
        if set(p) - PACK_KEYS or p.get("keep_item") or p.get("final"):
            raise unsupported(self.base, "key on a nested pack", sorted(set(p) - PACK_KEYS))
        inner = p.get("nodes") or []
        outs = [n["name"] for n in inner if n.get("type") == "compound_output"]
        if [n.get("type") for n in inner].count("compound_input") != 1 or not outs:
            raise unsupported(self.base, "nested pack ports", "need exactly one compound_input and at least one compound_output")
        for n in inner:
            if n.get("type") == "compound" or _has_pack(n.get("processors") or []):
                raise unsupported(self.base, "pack inside a nested pack", n.get("name"))
        into, bypass = self.gate(p, sources)
        name = self.new_name("pack")
        node = {"name": name, "type": "compound"}
        node.update({k: p[k] for k in ("user_description", "node_reference") if k in p})
        node.update(nodes=inner, links=p.get("links") or [])
        self.nodes.append(node)
        port = next(n["name"] for n in inner if n.get("type") == "compound_input")
        self.link(into, f"{name}->{port}")
        after = []
        for out in outs:
            tag = f"{name}_{out}"
            if tag in self.taken:
                raise unsupported(self.base, "generated node name already in use", tag)
            self.taken.add(tag)
            self.nodes.append({"name": tag, "type": "sequence", "processors": [{
                "type": "ottl_transform",
                "statements": f"set(attributes[{json.dumps(PACK_PATH_ATTR)}], {json.dumps(out)})"}]})
            self.link([(f"{name}->{out}", None)], tag)
            after.append((tag, None))
        return after + bypass

    def emit(self, procs, sources, depth):
        run = []
        for p in procs:
            if not (_active(p) and _has_pack([p])):
                if p.get("disabled") and _has_pack([dict(p, disabled=False)]):
                    continue  # Go skips a disabled processor (sequence_node.go:140-142)
                run.append(p)
                continue
            sources = self.flush(run, sources, depth, False)
            run = []
            if p.get("type") == "compound":
                sources = self.pack(p, sources)
                continue
            if set(p) - SUBSEQ_KEYS or p.get("keep_item") or p.get("final"):
                raise unsupported(self.base, "key on a nested sequence that holds a pack", sorted(set(p) - SUBSEQ_KEYS))
            into, bypass = self.gate(p, sources)
            sources = self.emit(p.get("processors") or [], into, depth + 1) + bypass
        return self.flush(run, sources, depth, True)


def lift_nested_packs(c):
    """Make every pack that is a processor of a top-level `sequence` a top-level node.

    Go semantics reproduced (internalv3/discover/sequence_node.go:136-238):
    a sequence is a chain; a processor with `data_types` and/or `condition`
    takes only the items matching `(type == t1 or ...) and (condition)`, the
    rest go on to the next processor unchanged; everything a nested pack emits
    on any compound_output goes on to the next processor, with the attribute
    `ed.pipeline.pack.path.name` set to that output's name; a nested sequence
    is gated once on entry and its output goes on to the parent's next
    processor.

    Rewrite of sequence S = [pre..., PACK, post...]:

      S (processors: pre) -> S__gateN (route_ottl, path `matched` = the gate
      condition) -matched-> S__packN-><its compound_input> (the pack, top level) ->
      S__packN_<out> (sequence: set the path-name attribute), one per
      compound_output -> S__seqM (processors: post) -> S's old targets;
      S__gateN -unmatched-> S__seqM.

    A nested sequence that holds a pack is opened the same way: its gate
    becomes a route_ottl, its processors continue the chain, and its
    `unmatched` items rejoin after it. Empty stretches create no node; links
    into and out of S are moved to the first and last nodes of the chain.
    Disabled packs (and disabled sequences holding one) are removed, as Go
    never builds them.

    Raises Unsupported for: `final` on a processor that is not in the last
    stretch of the top-level list (its items would skip the rest of S in Go),
    `final` / `keep_item` / unknown keys on the pack or on a nested sequence
    holding one, a pack without exactly one compound_input or without a
    compound_output, a pack nested in a nested pack or in a sequence inside
    any pack, a link leaving S with a path, and generated-name collisions.

    Names of the pack's inner nodes change from Go's
    `S@processor_<i>-><inner>` to `S__packN-><inner>`.
    """
    nodes = c.get("nodes") or []
    taken = {n["name"] for n in nodes}
    for n in nodes:
        if n.get("type") == "compound":
            for inner in n.get("nodes") or []:
                if _has_pack(inner.get("processors") or []):
                    raise unsupported(n["name"], "pack nested in a sequence inside a pack", inner.get("name"))
    out_nodes = []
    for n in nodes:
        if not (n.get("type") == "sequence" and _has_pack(n.get("processors") or [])):
            out_nodes.append(n)
            continue
        name = n["name"]
        lift = _Lift(n, taken)
        exits = lift.emit(n.get("processors") or [], [ENTRY], 0)
        entry = [l["to"] for l in lift.links if l["from"] == ENTRY[0]]
        if len(entry) != 1 or ENTRY in exits:
            raise unsupported(name, "sequence shape", "no single entry node")
        links = []
        for l in c.get("links") or []:
            if l.get("to") == name:
                l = dict(l, to=entry[0])
            if l.get("from") == name:
                if l.get("path"):
                    raise unsupported(name, "link path on a sequence", l["path"])
                for frm, path in exits:
                    links.append(dict(l, **{"from": frm}, **({"path": path} if path else {})))
                continue
            links.append(l)
        c["links"] = links + [l for l in lift.links if l["from"] != ENTRY[0]]
        out_nodes.extend(lift.nodes)
    c["nodes"] = out_nodes
    return c


# What each old Edge Delta output takes from a link (IncomingDataTypes in
# internalv3/graph/nodespecs/nodespecs.go:2627-2696; a link carries only the
# types both ends share, internalv3/graph/helper/helpers.go:186-262).
OLD_ED_OUTPUTS = {
    "ed_logs_output": {"log"},
    "ed_archive_output": {"log"},  # deprecated name of ed_logs_output (nodespecs.go:2644-2645)
    "ed_metrics_output": {"metric", "heartbeat"},
    "ed_health_output": {"health"},
    "ed_patterns_output": {"cluster_pattern_samples"},
}
# OutgoingDataTypes of the node types seen in front of an old output
# (nodespecs.go:62, 255, 708, 776, 803, 830, 1026, 1115, 1129).
EMITS = {
    "file_input": {"log"},
    "http_input": {"log"},
    "regex_filter": {"log"},
    "log_to_metric": {"metric"},
    "log_to_pattern": {"cluster_pattern_samples"},
    "ed_agent_stats_input": {"metric"},
    "ed_pipeline_io_stats_input": {"metric"},
    "ed_component_health_input": {"health"},
    "ed_node_health_input": {"health", "metric"},
}
GATEABLE = {"log", "metric"}  # types an OTTL gate can name on both agents


def old_outputs(c, name="edgedelta", drop_types=()):
    """Replace the old per-type Edge Delta outputs with one `ed_output`.

    Go builds all of them on the same pusher as ed_output
    (internalv3/discover/adaptor.go:2724-2731,
    graph/execution/execution_flow.go:325-338); they differ only in the data
    types they accept from a link. So a link `X -> old output` equals
    `X -> ed_output` exactly when every type X emits is accepted by one of the
    old outputs X is linked to. This is the backend's own `singleEDDestination`
    migration (backend/analytics/batch/pipeline_auto_migration.go:535-596)
    plus that type check.

    Per source X (with A = the types its old outputs accept together):
      - X's node type is in EMITS and emits only types in A: one link
        X -> ed_output (duplicates merged);
      - otherwise, when A is only log and/or metric: X -> `<X>__ed_types`
        (sequence [ottl_filter include `EDXDataType() == ...`]) -> ed_output;
      - otherwise Unsupported: an ed_output would also send types the old
        output dropped, and no OTTL gate can keep only health or pattern items.

    Also raises Unsupported when the config already has an ed_output next to
    old outputs, when an old output node carries settings, when a link into
    one has a compound port or the new node name is taken.

    `drop_types` is the caller's explicit, NOT equivalent, way past the third
    case: nodes of these types and their links are removed first (for example
    `("ed_node_health_input",)`, which RED does not run anyway). Whatever the
    Go agent sent from them is then no longer sent.
    """
    if drop_types and any(n.get("type") in OLD_ED_OUTPUTS for n in c.get("nodes") or []):
        gone = {n["name"] for n in c["nodes"] if n.get("type") in drop_types}
        c["nodes"] = [n for n in c["nodes"] if n["name"] not in gone]
        c["links"] = [l for l in c.get("links") or [] if l["from"] not in gone and l["to"] not in gone]
    nodes = c.get("nodes") or []
    old = {n["name"]: n for n in nodes if n.get("type") in OLD_ED_OUTPUTS}
    if not old:
        return c
    types = {n["name"]: n.get("type") for n in nodes}
    for n in old.values():
        extra = sorted(set(n) - {"name", "type", "user_description"})
        if extra:
            raise unsupported(n["name"], "setting on an old Edge Delta output", extra)
    if "ed_output" in types.values():
        raise unsupported(next(iter(old)), "old Edge Delta output next to an ed_output", sorted(old))
    if name in types:
        raise unsupported(name, "generated node name already in use", name)

    accepted, order = {}, []
    for l in c.get("links") or []:
        if l.get("to") in old:
            key = (l["from"], l.get("path") or "")
            if key not in accepted:
                order.append(key)
            accepted.setdefault(key, set()).update(OLD_ED_OUTPUTS[old[l["to"]]["type"]])
    new_nodes, new_links = [], []
    for frm, path in order:
        a = accepted[(frm, path)]
        src_type = types.get(frm)
        link = {"from": frm, **({"path": path} if path else {})}
        emits = EMITS.get(src_type) if not path else None
        if emits is not None and emits <= a:
            new_links.append(dict(link, to=name))
            continue
        if not (a - {"heartbeat"}) <= GATEABLE or src_type is None:
            raise unsupported(frm, "link into an old Edge Delta output",
                              f"{src_type} emits {sorted(emits) if emits else 'unknown types'}, the old output takes only {sorted(a)}")
        gate = f"{frm}__ed_types" + (f"_{path}" if path else "")
        if gate in types:
            raise unsupported(gate, "generated node name already in use", gate)
        types[gate] = "sequence"
        cond = " or ".join(f"EDXDataType() == {json.dumps(t)}" for t in sorted(a & GATEABLE))
        new_nodes.append({"name": gate, "type": "sequence",
                          "processors": [{"type": "ottl_filter", "condition": cond, "filter_mode": "include"}]})
        new_links += [dict(link, to=gate), {"from": gate, "to": name}]
    c["nodes"] = [n for n in nodes if n["name"] not in old] + new_nodes + [{"name": name, "type": "ed_output"}]
    c["links"] = [l for l in c.get("links") or [] if l.get("to") not in old] + new_links
    return c
