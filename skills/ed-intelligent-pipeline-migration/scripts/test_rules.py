"""Tests for l2m.py and small.py: the forms that must raise, and the accept gating."""
import unittest

import l2m
import small
import translate

ALL = tuple(small.DIFFERENCES)


class SmallTest(unittest.TestCase):

    def test_forms(self):
        def raises(f, *a, **k):
            try: f(*a, **k)
            except translate.Unsupported as e: return str(e)
            raise AssertionError(f"did not raise: {a}")
        def ok(f, *a, **k): return f(*a, **k)
        N = lambda t, **k: {"name": "n", "type": t, **k}
        T = lambda typ, *tr: N(typ, transformations=list(tr))
        # parse_json_attributes
        raises(small.parse_json_attributes, N("parse_json_attributes", process_field="item.body"))
        ok(small.parse_json_attributes, N("parse_json_attributes", process_field="item.body"), accept=ALL)
        for bad in ({"process_field": 'item["resource"]["x"]'}, {"field_path": 'item["attributes"]["x"]'}, {"field_path": "item.body"},
                    {"process_field": 'item["attributes"]["a"]', "field_path": 'item["attributes"]["b"]'}, {"process_field": "body"}, {"error_mode": "strict"},
                    {"process_field": 'json(item["body"])'}, {"process_field": 'item["attributes"]["a\\"b"]'}):
            raises(small.parse_json_attributes, N("parse_json_attributes", **bad), accept=ALL)
        # transform
        ok(small.transform, T("log_transform", {"field_path": "attributes.logEvents", "operation": "delete"}))
        for typ, tr in (("log_transform", {"field_path": 'item["attributes"]["a"]["b"]', "operation": "delete"}),
                        ("log_transform", {"field_path": 'item["resource"]["a"]', "operation": "delete"}),
                        ("log_transform", {"field_path": 'item["body"]', "operation": "upsert", "value": '"x"'}),
                        ("log_transform", {"field_path": "logEvents", "operation": "delete"}),
                        ("log_transform", {"field_path": 'item["attributes"]', "operation": "delete"}),
                        ("log_transform", {"field_path": 'item["attributes"]["a"]', "operation": "rename", "value": '"x"'}),
                        ("log_transform", {"field_path": 'item["attributes"]["a"]', "operation": "upsert", "value": 'item["body"]'}),
                        ("log_transform", {"field_path": 'item["attributes"]["a"]', "operation": "upsert", "value": 'from_k8s(item["resource"]["k8s.pod.name"], "x")'}),
                        ("log_transform", {"field_path": 'item["timestamp"]', "operation": "upsert", "value": 'convert_timestamp(item["attributes"]["t"], "2006-01-02", "Unix Milli")'}),
                        ("log_transform", {"field_path": 'item["timestamp"]', "operation": "upsert", "value": 'convert_timestamp(item["attributes"]["t"], "2006-01-02T15:04:05.999999Z", "Unix")'}),
                        ("log_transform", {"field_path": 'item["timestamp"]', "operation": "delete"})):
            raises(small.transform, T(typ, tr), accept=ALL)
        ts = T("log_transform", {"field_path": 'item["timestamp"]', "operation": "upsert", "value": 'convert_timestamp(item["attributes"]["t"], "2006-01-02T15:04:05.999999Z", "Unix Milli")'})
        assert "timestamp-unparseable" in raises(small.transform, ts); ok(small.transform, ts, accept=ALL)
        ok(small.transform, T("generic_transform", {"field_path": 'item["resource"]["a"]', "operation": "delete"}))
        # regex_filter
        for bad in ({"pattern": ""}, {"pattern": "a", "field_path": 'item["attributes"]["x"]'}, {"pattern": "a", "pattern_field_path": "x"}):
            raises(small.regex_filter, N("regex_filter", **bad))
        # mask
        ok(small.mask, N("mask", pattern=r"MQTTPassword:\s*\w+"))
        ok(small.mask, N("mask", pattern=r"secret=[a-z]+", mask=""))
        for bad in ({"pattern": r"^abc"}, {"pattern": r"abc$"}, {"pattern": r"\bfoo"}, {"pattern": r"pass=.*"}, {"pattern": r"[^ ]+"}, {"pattern": r"\S+"}, {"pattern": r"x*"},
                    {"pattern": r"[0-9]+", "mask": "000"}, {"pattern": r"(?i)[a-z]+", "mask": "X"}, {"pattern": r"[[:alpha:]]+"}, {"pattern": ""}, {"pattern": "a", "mask": 5},
                    {"pattern": "a", "capture_group_masks": []}, {"predefined_pattern": "credit_card"}, {"pattern": r"\*+"}, {"pattern": r"a{,3}b"}):
            raises(small.mask, N("mask", **bad), accept=ALL)
        g = N("mask", pattern=r'"sessionToken": "([^"]+)"')
        assert "mask-first-occurrence" in raises(small.mask, g); ok(small.mask, g, accept=ALL)
        for bad in (r"a|(b)", r"(a)?b", r"(a)*", r"x(?:y(z))", r"a(b*)c", r"(a)(?i)b", r"(a|b)|c", r"a(b){2}"):
            raises(small.mask, N("mask", pattern=bad), accept=ALL)
        ok(small.mask, N("mask", pattern=r"(?i)key=(\w+)"), accept=ALL)
        # unescape_json
        u = N("unescape_json")
        assert "unescape-text" in raises(small.unescape_json, u); ok(small.unescape_json, u, accept=ALL)
        for bad in ({"process_field": 'item["resource"]["x"]'}, {"field_path": 'item["attributes"]["a"]["b"]'}, {"process_field": "item.attributes"}, {"expression_type": "ottl"}):
            raises(small.unescape_json, N("unescape_json", **bad), accept=ALL)
        # grok
        gk = N("grok", pattern="%{INT:n} %{GREEDYDATA:m}")
        assert "grok-" in raises(small.grok, gk); raises(small.grok, gk, accept=("grok-unquote",)); ok(small.grok, gk, accept=ALL)
        for bad in ("%{IP:ip}", "%{COMBINEDAPACHELOG}", "%{INT:n:int}", "%{INT}", "%{INT:a} %{INT:a}", "%{INT:a.b}", "plain", "", "%{INT:n} (?<=x)y", "%{INT:a__b}"):
            raises(small.grok, N("grok", pattern=bad), accept=ALL)
        raises(small.grok, N("grok", pattern="%{INT:n}", pattern_field_path='item["attributes"]["p"]'), accept=ALL)
        # route
        R = lambda c, **k: N("route", paths=[{"path": "p", "condition": c, "exit_if_matched": True}], **k)
        for good in ('item["attributes"]["a"] == "m"', 'item["attributes"]["a"]["b.c"] != ""', 'item["attributes"]["a"]==1', 'regex_match(item["body"], "(?i)x")', 'regex_match(item["attributes"]["a"] , "x")'):
            ok(small.route, R(good))
        for bad in ('item["resource"]["a"] == "m"', 'item["body"] == "m"', 'item["attributes"]["a"] != 1', 'item["attributes"]["a"] == 1.5', 'has(item["attributes"]["a"])',
                    '!regex_match(item["body"], "x")', 'item["attributes"]["a"] == "m" && item["attributes"]["b"] == "n"', 'regex_match(item["resource"]["a"], "x")',
                    'item["attributes"] == "m"', 'item["attributes"]["a"] == item["attributes"]["b"]', 'item["attributes"]["a"].startsWith("x")', 'true'):
            raises(small.route, R(bad))
        raises(small.route, R('attributes["a"] == "m"', expression_type="ottl"))
        # log_to_pattern
        ok(small.log_to_pattern, N("log_to_pattern")); ok(small.log_to_pattern, N("log_to_pattern", num_of_clusters=10, samples_per_cluster=5, reporting_frequency="1m0s"))
        for bad in ({"group_by": ['item["attributes"]["a"]']}, {"similarity_threshold": 0.5}, {"field_path": "x"}):
            raises(small.log_to_pattern, N("log_to_pattern", **bad))



class LogToMetricTest(unittest.TestCase):

    def test_rejected(self):
        A=lambda *k:'item["attributes"]'+"".join('["%s"]'%x for x in k)
        cases={
         "p95 default": dict(name="n",pattern=".*",enabled_stats=["p95"]),
         "anomaly1 stat": dict(name="n",pattern=".*",enabled_stats=["anomaly1"]),
         "dup dim name": dict(name="n",pattern=".*",dimension_groups=[{"field_dimensions":[A("a","status"),A("status")]}]),
         "resource dim": dict(name="n",pattern=".*",dimension_groups=[{"field_dimensions":['item["resource"]["x"]']}]),
         "dotted path": dict(name="n",pattern=".*",dimension_groups=[{"field_dimensions":['item.attributes.x']}]),
         "unknown capture": dict(name="n",pattern="a=(\\d+)",dimension_groups=[{"dimensions":["a"]}]),
         "optional capture": dict(name="n",pattern="a=(?P<a>\\d+)?",dimension_groups=[{"dimensions":["a"]}]),
         "alternation capture": dict(name="n",pattern="a=(?P<a>\\d+)|b",dimension_groups=[{"dimensions":["a"]}]),
         "numeric capture \\S+": dict(name="n",pattern="a=(?P<a>\\S+)",dimension_groups=[{"numeric_dimension":"a"}]),
         "group interval": dict(name="n",pattern=".*",interval="1m",dimension_groups=[{"interval":"30s"}]),
         "group flag": dict(name="n",pattern=".*",only_report_nonzeros=False,dimension_groups=[{"only_report_nonzeros":True}]),
         "group_by key": dict(name="n",pattern=".*",group_by=["x"]),
         "value_adjustment": dict(name="n",pattern=".*",dimension_groups=[{"field_numeric_dimension":A("v"),"value_adjustment_rule":"value*2"}]),
         "same names": dict(name="n",pattern=".*",dimension_groups=[{"field_dimensions":[A("a")]},{"field_dimensions":[A("b")]}]),
         "pipe in name": dict(name="n",pattern=".*",metric_name="a|b"),
         "empty pattern": dict(name="n",pattern=""),
         "no name": dict(pattern=".*"),
        }

        for desc, n in cases.items():
            with self.subTest(desc):
                with self.assertRaises(translate.Unsupported):
                    l2m.log_to_metric(n)

    def test_percentiles_only_when_accepted(self):
        n = {"name": "n", "pattern": ".*", "enabled_stats": ["count", "p99"]}
        with self.assertRaises(translate.Unsupported):
            l2m.log_to_metric(dict(n))
        rules = l2m.log_to_metric(dict(n), approx_percentiles=True)["processors"]
        agg = [r["aggregation_type"] for p in rules if p["type"] == "aggregate_metric" for r in p["aggregate_metric_rules"]]
        self.assertEqual(sorted(agg), ["count", "p99"])


if __name__ == "__main__":
    unittest.main()
