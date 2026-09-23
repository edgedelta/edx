#!/usr/bin/env python3
"""Check onboarding recipe contracts; optionally execute edx processor dry-runs."""
import argparse
import copy
import json
import math
from pathlib import Path
import subprocess
import sys
import tempfile

DEFAULT = Path(__file__).resolve().parents[1] / 'assets/recipe-checks/cases.json'


def read(path):
    return json.loads(Path(path).read_text())


def conventions(pipeline):
    nodes = {n['name']: n for n in pipeline['nodes']}
    if len(nodes) != len(pipeline['nodes']):
        raise ValueError('duplicate node names')
    links = {(x['from'], x['to']) for x in pipeline['links']}
    internal = {'ed_self_telemetry_input', 'ed_system_stats_input'}

    def processors(items):
        for p in items:
            name = json.loads(p.get('metadata', '{}')).get('name', '')
            if not isinstance(name, str) or name.strip().lower() in ('', 'custom', 'custom ottl'):
                raise ValueError('nested processor needs a meaningful metadata.name')
            processors(p.get('processors', []))

    for name, n in nodes.items():
        kind = n['type']
        if kind not in internal and (kind.endswith('_input') or kind.endswith('_output')):
            attached = name + '_multiprocessor'
            edge = (name, attached) if kind.endswith('_input') else (attached, name)
            if nodes.get(attached, {}).get('type') != 'sequence' or edge not in links:
                raise ValueError(f'{name}: missing directly attached sequence {attached}')
        processors(n.get('processors', []))


def matches(actual, expected):
    if isinstance(expected, dict):
        return isinstance(actual, dict) and all(k in actual and matches(actual[k], v) for k, v in expected.items())
    if isinstance(expected, int) and not isinstance(expected, bool):
        return isinstance(actual, (int, float)) and not isinstance(actual, bool) and actual == expected
    if isinstance(expected, float):
        return isinstance(actual, (int, float)) and not isinstance(actual, bool) and math.isclose(actual, expected, rel_tol=1e-9, abs_tol=1e-12)
    return actual == expected


def verify(cases, actual):
    expected_ids = {c['id'] for c in cases}
    if set(actual) != expected_ids:
        raise ValueError(f'case IDs mismatch: got {sorted(actual)}, want {sorted(expected_ids)}')
    for c in cases:
        remaining = copy.deepcopy(actual[c['id']])
        if not isinstance(remaining, list) or len(remaining) != len(c['expected']):
            raise ValueError(f'{c["id"]}: wrong output count; extra metrics/logs and missing outputs fail')
        for want in c['expected']:
            index = next((i for i, got in enumerate(remaining) if matches(got, want)), None)
            if index is None:
                raise ValueError(f'{c["id"]}: output does not match expected signal/value/unit/timestamp/identity: {want}')
            remaining.pop(index)


def decode_response(data):
    if not isinstance(data, dict) or not isinstance(data.get('items'), dict):
        raise ValueError('dry-run did not return items by output path; validation/API errors are not passes')
    result = []
    for items in data['items'].values():
        if not isinstance(items, list):
            raise ValueError('unexpected dry-run output path format')
        for item in items:
            decoded = json.loads(item) if isinstance(item, str) else item
            if not isinstance(decoded, dict):
                raise ValueError('dry-run item must be an object')
            result.append(decoded)
    return result


def execute(pack, base, args):
    if not args.conf_id or not args.profile or not args.evidence_dir:
        raise ValueError('--run requires --conf-id, --profile and --evidence-dir')
    evidence = Path(args.evidence_dir)
    evidence.mkdir(parents=True, exist_ok=True)
    actual = {}
    with tempfile.TemporaryDirectory(prefix='ed-onboard-recipe-') as temp:
        temp = Path(temp)
        for c in pack['cases']:
            recipe = pack['recipes'][c['recipe']]
            pipeline = read(base / recipe['pipeline'])
            n = next(n for n in pipeline['nodes'] if n['name'] == recipe['node'])
            node = {'id': n['name'], 'type': n['type'], 'configuration': {k: v for k, v in n.items() if k not in ('name', 'type')}}
            (temp / 'node.json').write_text(json.dumps(node))
            (temp / 'input.jsonl').write_text(json.dumps(c['input']) + '\n')
            response = evidence / (c['id'] + '.json')
            cmd = [args.edx, '--profile', args.profile, 'pipelines', 'test', 'node', args.conf_id,
                   '--file', str(temp / 'input.jsonl'), '--node', str(temp / 'node.json'), '--output-file', str(response)]
            r = subprocess.run(cmd, capture_output=True, text=True, timeout=90)
            if r.returncode:
                (evidence / (c['id'] + '.error.txt')).write_text(r.stderr)
                raise ValueError(f'{c["id"]}: processor execution unavailable/failed; see {evidence}. No deployment or ingestion performed.')
            actual[c['id']] = decode_response(read(response))
            (evidence / 'actual.json').write_text(json.dumps(actual, indent=2) + '\n')
    verify(pack['cases'], actual)
    print(f'PASS: {len(actual)} processor dry-runs; this does not verify live source delivery or indexing')


def self_test(pack, pipelines):
    cases = pack['cases']
    good = {c['id']: copy.deepcopy(c['expected']) for c in cases}
    verify(cases, good)
    rejected = 0
    mutations = [
        ('name', 'wrong.metric'), ('gauge', {'value': 99}), ('unit', 'b'),
        ('timestamp', 1700000000001), ('resource', {'aws.rds.db_instance_identifier': 'other-db'})]
    for key, value in mutations:
        bad = copy.deepcopy(good)
        bad['cpu-average'][0][key] = value
        try:
            verify(cases, bad)
        except ValueError:
            rejected += 1
        else:
            raise ValueError(f'checker accepted incorrect {key}')
    for cid in ['cpu-average', 'other-database', 'other-namespace', 'zero-count', 'unknown-metric', 'aggregate-without-db']:
        bad = copy.deepcopy(good)
        bad[cid].append(copy.deepcopy(good['cpu-average'][0]))
        try:
            verify(cases, bad)
        except ValueError:
            rejected += 1
        else:
            raise ValueError(f'checker accepted extra output for {cid}')
    for mutation in ['attachment', 'name']:
        bad = copy.deepcopy(pipelines['metrics'])
        if mutation == 'attachment':
            bad['links'] = [e for e in bad['links'] if e['from'] != 'rds_metric_stream']
        else:
            bad['nodes'][1]['processors'][0]['metadata'] = '{"name":"Custom"}'
        try:
            conventions(bad)
        except ValueError:
            rejected += 1
        else:
            raise ValueError(f'checker accepted broken {mutation}')
    try:
        decode_response({'valid': False})
    except ValueError:
        rejected += 1
    else:
        raise ValueError('checker accepted failed validation')
    print(f'PASS: checker rejected {rejected} faulty cases; no Edge Delta processor was executed')


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--fixtures', type=Path, default=DEFAULT)
    modes = parser.add_mutually_exclusive_group()
    modes.add_argument('--self-test', action='store_true')
    modes.add_argument('--actual', type=Path, help='JSON mapping case IDs to actual output objects')
    modes.add_argument('--run', action='store_true', help='execute backend dry-runs via edx; no deployment/ingestion')
    parser.add_argument('--edx', default='edx')
    parser.add_argument('--profile')
    parser.add_argument('--conf-id')
    parser.add_argument('--evidence-dir', type=Path)
    args = parser.parse_args()
    pack = read(args.fixtures)
    pipelines = {k: read(args.fixtures.parent / r['pipeline']) for k, r in pack['recipes'].items()}
    for pipeline in pipelines.values():
        conventions(pipeline)
    print('PASS: recipe attachment and processor-name checks (static only)')
    if args.self_test:
        self_test(pack, pipelines)
    elif args.actual:
        verify(pack['cases'], read(args.actual))
        print('PASS: supplied outputs satisfy fixture expectations; verify their execution provenance separately')
    elif args.run:
        execute(pack, args.fixtures.parent, args)


if __name__ == '__main__':
    try:
        main()
    except (ValueError, KeyError, TypeError, OSError, subprocess.TimeoutExpired) as error:
        print(f'FAIL: {error}', file=sys.stderr)
        sys.exit(1)
