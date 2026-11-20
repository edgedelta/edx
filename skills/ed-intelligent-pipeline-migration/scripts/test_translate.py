"""Tests for translate.py on synthetic pipelines."""
import copy
import os
import subprocess
import sys
import tempfile
import unittest

import yaml

import translate

HERE = os.path.dirname(os.path.abspath(__file__))


def cfg(nodes, links=()):
    return {"version": "v3", "nodes": list(nodes), "links": list(links)}


def lt(*transformations, typ="log_transform"):
    return {"name": "t", "type": typ, "transformations": list(transformations)}


def up(path, value, **kw):
    return {"field_path": path, "operation": "upsert", "value": value, **kw}


def ottl(statements, name="t"):
    return {"name": name, "type": "ottl_transform", "statements": "\n".join(statements)}


def route(*paths, **kw):
    return {"name": "r", "type": "route", **kw,
            "paths": [{"path": f"p{i}", "condition": c, "exit_if_matched": True} for i, c in enumerate(paths)]}


def route_ottl(*conds):
    return {"name": "r", "type": "route_ottl",
            "paths": [{"path": f"p{i}", "condition": c, "exit_if_matched": True} for i, c in enumerate(conds)]}


OUT = {"name": "out", "type": "ed_output"}
L2M = {"name": "m", "type": "log_to_metric", "pattern": "(?i)error", "interval": "1m0s"}


def l2m_seq(metric="m", interval="1m0s"):
    return {"name": "m", "type": "sequence", "processors": [
        {"type": "ottl_filter", "data_types": ["log"], "filter_mode": "include", "condition": 'IsMatch(body, "(?i)error")'},
        {"type": "extract_metric", "data_types": ["log"], "keep_item": False,
         "extract_metric_rules": [{"name": metric, "gauge": {"value": "1"}}]},
        {"type": "ottl_filter", "data_types": ["log"], "filter_mode": "exclude", "condition": "true"},
        {"type": "aggregate_metric", "data_types": ["metric"], "aggregate_metric_rules": [
            {"name": f"{metric}.count", "conditions": [f'name == "{metric}"'], "interval": interval, "aggregation_type": "count",
             "group_by": ['resource["__group_name"]', 'resource["__src_name"]']}]},
        {"type": "ottl_transform", "data_types": ["metric"],
         "statements": 'set(_stat_type, "count") where IsMatch(name, "\\\\.count$")'}]}


L2M_SEQ = l2m_seq()
NULL_FIX = 'set(attributes["a"], 0) where IsMap(resource) and ContainsValue(Keys(resource), "k") and resource["k"] == nil'
PACK_IO = [{"name": "in", "type": "compound_input"}, {"name": "out", "type": "compound_output"}]
LEVEL = 'set(severity_text, "unknown")'


class TranslateTest(unittest.TestCase):

    def test_translate(self):
        cases = [
            {
                "desc": "literal",
                "in": cfg([lt(up('item["attributes"]["a"]', '"x"'))]),
                "want": cfg([ottl(['set(attributes["a"], "x")'])]),
            },
            {
                "desc": "env with default",
                "in": cfg([lt(up('item["resource"]["a"]', 'env("K") != "" ? env("K") : "d"'))]),
                "want": cfg([ottl(['set(resource["a"], "d")',
                                   'set(resource["a"], EDXEnv("K", "")) where EDXEnv("K", "") != ""'])]),
            },
            {
                "desc": "bare env",
                "in": cfg([lt(up('item["attributes"]["a"]', 'env("K")'))]),
                "want": cfg([ottl(['set(attributes["a"], EDXEnv("K", ""))'])]),
            },
            {
                "desc": "field reference",
                "in": cfg([lt(up('item["attributes"]["a"]["b"]', 'item["resource"]["x.y"]'))]),
                "want": cfg([ottl(['set(attributes["a"]["b"], "")',
                                   'set(attributes["a"]["b"], resource["x.y"]) where resource["x.y"] != nil'])]),
            },
            {
                "desc": "in with default",
                "in": cfg([lt(up('item["attributes"]["a"]', '"k" in item["resource"] ? item["resource"]["k"] : "d"'))]),
                "want": cfg([ottl(['set(attributes["a"], "d")',
                                   'set(attributes["a"], resource["k"]) where resource["k"] != nil', NULL_FIX])]),
            },
            {
                "desc": "regex_capture",
                "in": cfg([lt(up('item["attributes"]["f"]', 'regex_capture(item["resource"]["p"], "/(?P<f>[^/]+)$").f'))]),
                "want": cfg([ottl(['set(attributes["f"], "")',
                                   'set(attributes["f"], ExtractPatterns(resource["p"], "/(?P<f>[^/]+)$")["f"]) '
                                   'where IsString(resource["p"]) and IsMatch(resource["p"], "/(?P<f>[^/]+)$")'])]),
            },
            {
                "desc": "delete",
                "in": cfg([lt({"field_path": 'item["attributes"]["a"]["b"]', "operation": "delete"})]),
                "want": cfg([ottl(['delete_key(attributes["a"], "b")'])]),
            },
            {
                "desc": "ignore_if_empty",
                "in": cfg([lt(up('item["attributes"]["a"]', 'item["resource"]["k"]', ignore_if_empty=True),
                              up('item["attributes"]["b"]', '""', ignore_if_empty=True),
                              up('item["attributes"]["c"]', 'env("K") != "" ? env("K") : "d"', ignore_if_empty=True))]),
                "want": cfg([ottl(['set(attributes["a"], resource["k"]) where resource["k"] != nil and resource["k"] != ""',
                                   NULL_FIX, 'set(attributes["c"], "d")',
                                   'set(attributes["c"], EDXEnv("K", "")) where EDXEnv("K", "") != ""'])]),
            },
            {
                "desc": "generic_transform keeps user_description",
                "in": cfg([{**lt(up('item["attributes"]["a"]', '"x"'), typ="generic_transform"), "user_description": "u"}]),
                "want": cfg([{"name": "t", "user_description": "u", "type": "ottl_transform",
                              "statements": 'set(attributes["a"], "x")'}]),
            },
            {
                "desc": "resource_transform links",
                "in": cfg([{"name": "rt", "type": "resource_transform", "transformations": [
                    {"field_path": "a", "value": 'item["resource"]["x"]'},
                    {"field_path": "b", "value": 'item["attributes"]["y"]'}]}, OUT],
                    [{"from": "in", "to": "rt"}, {"from": "rt", "to": "out"},
                     {"from": "rt", "path": "failure", "to": "other"}]),
                "want": cfg([{"name": "rt", "type": "route_ottl", "paths": [
                    {"path": "ok", "condition": 'IsString(resource["x"]) and IsString(attributes["y"])', "exit_if_matched": True}]},
                    ottl(['set(resource["a"], resource["x"])', 'set(resource["b"], attributes["y"])'], name="rt__apply"), OUT],
                    [{"from": "in", "to": "rt"}, {"from": "rt__apply", "to": "out"},
                     {"from": "rt", "path": "unmatched", "to": "other"}, {"from": "rt", "path": "ok", "to": "rt__apply"}]),
            },
            {
                "desc": "regex_filter",
                "in": cfg([{"name": "f", "type": "regex_filter", "pattern": "a\\d"},
                           {"name": "g", "type": "regex_filter", "pattern": "b", "negate": True}]),
                "want": cfg([{"name": "f", "type": "ottl_filter", "filter_mode": "include", "condition": 'IsMatch(body, "a\\\\d")'},
                             {"name": "g", "type": "ottl_filter", "filter_mode": "exclude", "condition": 'IsMatch(body, "b")'}]),
            },
            {
                "desc": "log_level_detector",
                "in": cfg([{"name": "a", "type": "log_level_detector", "pattern": "INFO|WARN(?:ING)?", "fallback_value": "unknown"},
                           {"name": "b", "type": "log_level_detector", "pattern": "level=(\\w+)", "fallback_value": "unknown"},
                           {"name": "c", "type": "log_level_detector", "pattern": "[(]lvl=(?P<x>\\w+)", "fallback_value": "unknown"}]),
                "want": cfg([ottl([LEVEL, 'set(severity_text, ExtractPatterns(body, "(?P<lvl>INFO|WARN(?:ING)?)")["lvl"]) '
                                          'where IsString(body) and IsMatch(body, "INFO|WARN(?:ING)?")'], name="a"),
                             ottl([LEVEL, 'set(severity_text, ExtractPatterns(body, "level=(?P<lvl>\\\\w+)")["lvl"]) '
                                          'where IsString(body) and IsMatch(body, "level=(\\\\w+)")'], name="b"),
                             ottl([LEVEL, 'set(severity_text, ExtractPatterns(body, "[(]lvl=(?P<lvl>\\\\w+)")["lvl"]) '
                                          'where IsString(body) and IsMatch(body, "[(]lvl=(?P<x>\\\\w+)")'], name="c")]),
            },
            {
                "desc": "parse_json_attributes",
                "in": cfg([{"name": "j", "type": "parse_json_attributes", "process_field": 'item["body"]',
                            "field_path": 'item["attributes"]'}]),
                "accept": ("json-parser",),
                "want": cfg([ottl(['set(cache["pja"], ParseJSON(body)) where IsString(body)',
                                   'set(cache["pja"], body) where IsMap(body)',
                                   'set(cache["fail"], ParseJSON("{")) where not IsMap(cache["pja"])',
                                   'set(cache["fail"], ParseJSON("{")) where Len(cache["pja"]) == 0',
                                   'set(attributes, cache["pja"])'], name="j")]),
            },
            {
                "desc": "route atoms",
                "in": cfg([route('item["attributes"]["a"] == "x"', 'item["resource"]["a"] != "x"',
                                 'has(item.attributes.a) && !has(item["resource"]["b"])',
                                 'regex_match(item["attributes"]["a"], "^x") || !regex_match(item["attributes"]["b"], "y")',
                                 expression_type="cel")]),
                "want": cfg([route_ottl('attributes["a"] == "x"', 'resource["a"] != nil and resource["a"] != "x"',
                                        'attributes["a"] != nil and resource["b"] == nil',
                                        'IsMatch(attributes["a"], "^x") or IsString(attributes["b"]) and not IsMatch(attributes["b"], "y")')]),
            },
            {
                "desc": "route body guard",
                "in": cfg([route('(item["body"]["a"] == "1" || item["body"]["a"] == "2") && (item["body"]["b"] == "x" || !has(item.body.b))',
                                 'item["attributes"]["s"] == "x" || item["body"]["a"] == "1"')]),
                "want": cfg([route_ottl('IsMap(body) and (body["a"] == "1" or body["a"] == "2") and (body["b"] == "x" or body["b"] == nil)',
                                        'attributes["s"] == "x" or IsMap(body) and body["a"] == "1"')]),
            },
            {
                "desc": "OTTL route unchanged",
                "in": cfg([route('attributes["a"] == "x"', expression_type="ottl")]),
                "want": cfg([route('attributes["a"] == "x"', expression_type="ottl")]),
            },
            {
                "desc": "log_to_pattern",
                "in": cfg([{"name": "p", "type": "log_to_pattern", "group_by": ['item["attributes"]["a"]'], "num_of_clusters": 10}],
                          [{"from": "p", "to": "out"}, {"from": "p", "path": "failure", "to": "out"}]),
                "want": cfg([{"name": "p", "type": "sequence", "processors": [
                    {"type": "log_to_pattern_metric", "group_by": ['attributes["a"]'], "num_of_clusters": 10, "data_types": ["log"]}]}],
                    [{"from": "p", "to": "out"}]),
            },
            {
                "desc": "log_to_metric and output_transform",
                "in": cfg([L2M, {"name": "ot", "type": "output_transform", "transformations": []},
                           {"name": "s", "type": "sequence"}, OUT],
                          [{"from": "in", "to": "m"}, {"from": "m", "to": "ot"}, {"from": "ot", "to": "s"}, {"from": "s", "to": "out"}]),
                "want": cfg([L2M_SEQ, {"name": "s", "type": "sequence"}, OUT],
                            [{"from": "in", "to": "m"}, {"from": "m", "to": "s"}, {"from": "s", "to": "out"}]),
            },
            {
                "desc": "log_to_metric name and default interval",
                "in": cfg([{"name": "m", "type": "log_to_metric", "pattern": "(?i)error", "metric_name": "a-b c"}]),
                "want": cfg([l2m_seq("a_b_c", "1m")]),
            },
            {
                "desc": "comments and lookup",
                "in": cfg([{"name": "s", "type": "sequence", "processors": [
                    {"type": "ottl_transform", "statements": "// note\nset(attributes[\"a\"], \"b\")"}]},
                    {"name": "k", "type": "lookup", "location_path": "file:///t.csv",
                     "key_fields": [{"event_field": 'item["attributes"]["a"]', "lookup_field": "a"}],
                     "out_fields": [{"event_field": 'item["attributes"]["b"]', "lookup_field": "b"}]}]),
                "want": cfg([{"name": "s", "type": "sequence", "processors": [
                    {"type": "ottl_transform", "statements": "set(attributes[\"a\"], \"b\")"}]},
                    {"name": "k", "type": "sequence", "processors": [
                        {"type": "lookup", "location_path": "file:///t.csv",
                         "key_fields": [{"event_field": 'attributes["a"]', "lookup_field": "a"}],
                         "out_fields": [{"event_field": 'attributes["b"]', "lookup_field": "b"}], "data_types": ["log"]}]}]),
            },
            {
                "desc": "json_unroll item prefix, also inside a sequence",
                "in": cfg([{"name": "u", "type": "json_unroll", "field_path": 'item["body"]', "json_field_path": "Records"},
                           {"name": "s", "type": "sequence", "processors": [
                               {"type": "json_unroll", "field_path": "item.body", "json_field_path": "r"},
                               {"type": "json_unroll", "field_path": 'attributes["x"]', "json_field_path": "r"}]}]),
                "want": cfg([{"name": "u", "type": "json_unroll", "field_path": "body", "json_field_path": "Records"},
                             {"name": "s", "type": "sequence", "processors": [
                                 {"type": "json_unroll", "field_path": "body", "json_field_path": "r"},
                                 {"type": "json_unroll", "field_path": 'attributes["x"]', "json_field_path": "r"}]}]),
            },
            {
                "desc": "nodes inside a pack",
                "in": cfg([{"name": "pk", "type": "compound", "nodes": PACK_IO + [
                    {"name": "f", "type": "regex_filter", "pattern": "a", "negate": True}, L2M],
                    "links": [{"from": "in", "to": "f"}, {"from": "f", "to": "m"}, {"from": "m", "to": "out"}]}]),
                "want": cfg([{"name": "pk", "type": "compound", "nodes": PACK_IO + [
                    {"name": "f", "type": "ottl_filter", "filter_mode": "exclude", "condition": 'IsMatch(body, "a")'}, L2M_SEQ],
                    "links": [{"from": "in", "to": "f"}, {"from": "f", "to": "m"}, {"from": "m", "to": "out"}]}]),
            },
            {
                "desc": "percentiles accepted",
                "in": cfg([{**L2M, "enabled_stats": ["p95"]}]),
                "accept": ("percentiles",),
                "want_types": ["sequence"],
            },
            {
                "desc": "route form only the earlier rule covers",
                "in": cfg([route('item["attributes"]["a"] == "x" && has(item["attributes"]["b"])')]),
                "want": cfg([route_ottl('attributes["a"] == "x" and attributes["b"] != nil')]),
            },
        ]
        for tc in cases:
            with self.subTest(tc["desc"]):
                got = copy.deepcopy(tc["in"])
                translate.translate(got, tc.get("accept", ()))
                if "want_types" in tc:
                    self.assertEqual([n["type"] for n in got["nodes"]], tc["want_types"])
                    continue
                if got != tc["want"]:
                    self.fail(f"{tc['desc']}: got: {got}, want: {tc['want']}")

    def test_summary(self):
        c = cfg([lt(up('item["attributes"]["a"]', '"x"')), L2M, {"name": "ot", "type": "output_transform"}, OUT],
                [{"from": "m", "to": "ot"}, {"from": "ot", "to": "out"}])
        got = translate.translate(c)
        want = {"translated": {"log_to_metric": 1, "log_transform": 1}, "removed": ["ot"], "log_to_metric": ["m"],
                "structure": [], "accepted": {}}
        if got != want:
            self.fail(f"summary: got: {got}, want: {want}")

    def test_rejected(self):
        cases = [
            ("convert_timestamp", cfg([lt(up('item["attributes"]["t"]', 'convert_timestamp(item["attributes"]["t"], "a", "b")'))]),
             't: unsupported CEL value: convert_timestamp(item["attributes"]["t"], "a", "b")'),
            ("CEL macro", cfg([lt(up('item["attributes"]["a"]', 'item["attributes"]["l"].exists(x, x == "a")'))]),
             't: unsupported CEL value: item["attributes"]["l"].exists(x, x == "a")'),
            ("ec2_metadata", cfg([lt(up('item["attributes"]["a"]', 'ec2_metadata("region")'))]),
             't: unsupported CEL value: ec2_metadata("region")'),
            ("from_k8s", cfg([lt(up('item["attributes"]["a"]', 'from_k8s(item["resource"]["p"], "x")'))]),
             't: unsupported CEL value: from_k8s(item["resource"]["p"], "x")'),
            ("body field path", cfg([lt(up('item["body"]', '"x"'))]),
             't: unsupported field path: item["body"]'),
            ("body value", cfg([lt(up('item["attributes"]["a"]', 'item["body"]["k"]'))]),
             't: unsupported CEL value: item["body"]["k"]'),
            ("operation", cfg([lt({"field_path": 'item["attributes"]["a"]', "operation": "append", "value": '"x"'})]),
             "t: unsupported operation: append"),
            ("ignore_if_empty in with default", cfg([lt(up('item["attributes"]["a"]', '"k" in item["resource"] ? item["resource"]["k"] : "d"',
                                                           ignore_if_empty=True))]),
             't: unsupported ignore_if_empty value: "k" in item["resource"] ? item["resource"]["k"] : "d"'),
            ("resource_transform literal", cfg([{"name": "rt", "type": "resource_transform",
                                                 "transformations": [{"field_path": "a", "value": '"x"'}]}]),
             'rt: unsupported resource value: "x"'),
            ("regex_filter field_path", cfg([{"name": "f", "type": "regex_filter", "pattern": "a", "field_path": 'item["attributes"]["a"]'}]),
             "f: unsupported key: field_path"),
            ("regex_filter pattern_field_path", cfg([{"name": "f", "type": "regex_filter", "pattern_field_path": 'item["attributes"]["a"]'}]),
             "f: unsupported key: pattern_field_path"),
            ("empty fallback", cfg([{"name": "d", "type": "log_level_detector", "pattern": "INFO"}]),
             "d: unsupported empty fallback_value: INFO"),
            ("several groups", cfg([{"name": "d", "type": "log_level_detector", "pattern": "(a)(b)", "fallback_value": "x"}]),
             "d: unsupported pattern with several capture groups: (a)(b)"),
            ("parse_json_attributes needs accept", cfg([{"name": "j", "type": "parse_json_attributes", "process_field": 'item["body"]',
                                                         "field_path": 'item["attributes"]'}]),
             "j: unsupported without accept='json-parser': " + translate.accept_flags()["json-parser"]),
            ("route negation", cfg([route('!(item["attributes"]["a"] == "x")')]),
             'r: unsupported route condition: !(item["attributes"]["a"] == "x")'),
            ("route _type", cfg([route('item["_type"] == "log"')]),
             'r: unsupported route condition: item["_type"] == "log"'),
            ("route function", cfg([route('size(item["attributes"]["a"]) == "1"')]),
             'r: unsupported route condition: size(item["attributes"]["a"]) == "1"'),
            ("route single quotes", cfg([route("item['attributes']['a'] == 'x'")]),
             "r: unsupported route condition: item['attributes']['a'] == 'x'"),
            ("percentiles need accept", cfg([{**L2M, "enabled_stats": ["count", "p95"]}]),
             "m: unsupported enabled_stats: p95 (rerun with --accept percentiles)"),
            ("old outputs with node health", cfg([{"name": "h", "type": "ed_node_health_input"}, {"name": "eh", "type": "ed_health_output"}],
                                                 [{"from": "h", "to": "eh"}]),
             None),
            ("v2 config", {"version": "v2", "inputs": {}},
             "config: unsupported format: no node list (a v2 pipeline has to be rebuilt in v3 first)"),
            ("output_transform to other output", cfg([L2M, {"name": "ot", "type": "output_transform"}, {"name": "h", "type": "http_output"}],
                                                     [{"from": "m", "to": "ot"}, {"from": "ot", "to": "h"}]),
             "ot: unsupported output_transform feeding: h"),
            ("output_transform fed by logs", cfg([{"name": "i", "type": "http_input"}, {"name": "ot", "type": "output_transform"}, OUT],
                                                 [{"from": "i", "to": "ot"}, {"from": "ot", "to": "out"}]),
             "ot: unsupported output_transform input: i"),
            ("wrapped failure link", cfg([{"name": "p", "type": "log_to_pattern"}],
                                         [{"from": "p", "to": "out"}, {"from": "p", "path": "failure", "to": "other"}]),
             "p: unsupported failure link: other differs from its success links"),
        ]
        for desc, c, want in cases:
            with self.subTest(desc):
                try:
                    translate.translate(c)
                    got = None
                except translate.Unsupported as e:
                    got = str(e)
                if desc == "old outputs with node health":
                    # The reason is structure.py's; only the hint is checked here.
                    self.assertTrue(got and got.endswith("(rerun with --accept drop-node-health)"), got)
                    continue
                if got != want:
                    self.fail(f"{desc}: got: {got}, want: {want}")

    def test_structure(self):
        pack = {"type": "compound", "data_types": ["log"], "nodes": PACK_IO + [ottl(['set(attributes["p"], "1")'], name="x")],
                "links": [{"from": "in", "to": "x"}, {"from": "x", "to": "out"}]}
        c = cfg([{"name": "s", "type": "sequence", "processors": [
            {"type": "ottl_transform", "statements": 'set(attributes["a"], "1")'}, pack,
            {"type": "ottl_transform", "statements": 'set(attributes["z"], "1")'}]}, OUT], [{"from": "s", "to": "out"}])
        got = translate.translate(c)
        self.assertEqual(got["structure"], ["lift_nested_packs"])
        types = [n["type"] for n in c["nodes"]]
        self.assertEqual(types.count("compound"), 1, types)
        self.assertFalse(any(p.get("type") == "compound" for n in c["nodes"] for p in n.get("processors") or []), c["nodes"])
        names = {n["name"] for n in c["nodes"]}
        for l in c["links"]:
            self.assertIn(l["from"].split("->")[0], names, l)
            self.assertIn(l["to"].split("->")[0], names, l)

        old = cfg([{"name": "i", "type": "http_input"}, {"name": "h", "type": "ed_node_health_input"},
                   {"name": "lo", "type": "ed_logs_output"}, {"name": "eh", "type": "ed_health_output"}],
                  [{"from": "i", "to": "lo"}, {"from": "h", "to": "eh"}])
        got = translate.translate(old, ("drop-node-health",))
        self.assertEqual((got["structure"], got["accepted"]), (["old_outputs"], {"drop-node-health": 1}))
        self.assertEqual(sorted(n["type"] for n in old["nodes"]), ["ed_output", "http_input"])

    def test_accept_cli(self):
        c = cfg([{"name": "j", "type": "parse_json_attributes", "process_field": "item.body"}])
        with tempfile.TemporaryDirectory() as d:
            src, dst = os.path.join(d, "go.yaml"), os.path.join(d, "red.yaml")
            with open(src, "w") as f:
                yaml.safe_dump(c, f)
            run = lambda *a: subprocess.run([sys.executable, os.path.join(HERE, "translate.py"), *a], capture_output=True, text=True)
            p = run(src, dst)
            self.assertEqual((p.returncode, os.path.exists(dst)), (1, False), p.stderr)
            self.assertIn("accept='json-parser'", p.stderr)
            p = run("--accept", "json-parser", src, dst)
            self.assertEqual(p.returncode, 0, p.stderr)
            self.assertEqual(yaml.safe_load(p.stdout)["accepted"], {"json-parser": 1})
            self.assertEqual(run("--accept", "all", src, dst).returncode, 0)
            self.assertEqual(run("--accept", "nope", src, dst).returncode, 2)
            self.assertIn("json-parser", run("--list-accept").stdout)

    def test_cli(self):
        cases = [
            {"desc": "translated", "in": cfg([lt(up('item["attributes"]["a"]', '"x"'))]), "code": 0,
             "want": cfg([ottl(['set(attributes["a"], "x")'])])},
            {"desc": "rejected writes nothing", "in": cfg([lt(up('item["attributes"]["a"]', 'convert_timestamp("x")'))]), "code": 1,
             "want": None},
        ]
        for tc in cases:
            with self.subTest(tc["desc"]), tempfile.TemporaryDirectory() as d:
                src, dst = os.path.join(d, "go.yaml"), os.path.join(d, "red.yaml")
                with open(src, "w") as f:
                    yaml.safe_dump(tc["in"], f)
                p = subprocess.run([sys.executable, os.path.join(HERE, "translate.py"), src, dst], capture_output=True, text=True)
                got = None
                if os.path.exists(dst):
                    with open(dst) as f:
                        got = yaml.safe_load(f)
                if (p.returncode, got) != (tc["code"], tc["want"]):
                    self.fail(f"{tc['desc']}: got: {(p.returncode, got, p.stderr)}, want: {(tc['code'], tc['want'])}")


if __name__ == "__main__":
    unittest.main()
