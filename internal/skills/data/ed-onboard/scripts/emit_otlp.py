#!/usr/bin/env python3
"""Emit bounded, identifiable application telemetry to an OTLP HTTP receiver."""
import json
import os
import secrets
import time
import urllib.request
import urllib.error


def payloads(run_id, sequence):
    now = time.time_ns()
    trace_id = secrets.token_hex(16)
    resource = {"attributes": [{"key": "service.name", "value": {"stringValue": run_id}}]}
    scope = {"name": "ed-onboard-fixture"}
    log = {"timeUnixNano": str(now), "severityNumber": 9, "severityText": "INFO",
           "body": {"stringValue": f"{run_id} onboarding sequence={sequence}"},
           "traceId": trace_id, "spanId": "0123456789abcdef"}
    point = {"timeUnixNano": str(now), "asDouble": float(sequence)}
    span = {"traceId": trace_id, "spanId": "0123456789abcdef", "name": "onboarding-check",
            "kind": 2, "startTimeUnixNano": str(now-1000000), "endTimeUnixNano": str(now),
            "status": {"code": 1}}
    return {
        "logs": {"resourceLogs": [{"resource": resource, "scopeLogs": [{"scope": scope, "logRecords": [log]}]}]},
        "metrics": {"resourceMetrics": [{"resource": resource, "scopeMetrics": [{"scope": scope, "metrics": [
            {"name": "ed.onboard.heartbeat", "gauge": {"dataPoints": [point]}}]}]}]},
        "traces": {"resourceSpans": [{"resource": resource, "scopeSpans": [{"scope": scope, "spans": [span]}]}]},
    }


def main():
    run_id = os.environ["ONBOARD_RUN_ID"]
    endpoint = os.environ.get("OTEL_EXPORTER_OTLP_ENDPOINT", "http://127.0.0.1:4318").rstrip("/")
    duration = int(os.environ.get("ONBOARD_DURATION_SECONDS", "600"))
    if not 1 <= duration <= 3600:
        raise ValueError("duration must be between 1 and 3600 seconds")
    deadline = time.monotonic() + duration
    sequence = 0
    accepted = {signal: 0 for signal in ("logs", "metrics", "traces")}
    while time.monotonic() < deadline:
        sequence += 1
        for signal, payload in payloads(run_id, sequence).items():
            request = urllib.request.Request(endpoint+"/v1/"+signal,
                data=json.dumps(payload).encode(), headers={"Content-Type": "application/json"})
            try:
                with urllib.request.urlopen(request, timeout=5) as response:
                    result = json.loads(response.read() or b"{}")
                    partial = result.get("partialSuccess", {})
                    if any(int(partial.get(k, 0)) for k in ("rejectedLogRecords", "rejectedDataPoints", "rejectedSpans")):
                        raise ValueError(f"{signal}: partial rejection")
                    accepted[signal] += 1
            except (urllib.error.URLError, TimeoutError, ValueError) as error:
                print(f"{signal}: {error}", flush=True)
        if os.environ.get("ONBOARD_LOG_FILE"):
            with open(os.environ["ONBOARD_LOG_FILE"], "a") as handle:
                handle.write(json.dumps({"run_id": run_id, "sequence": sequence, "message": "onboarding file log"})+"\n")
        time.sleep(min(5, max(0, deadline-time.monotonic())))
    print(json.dumps({"run_id": run_id, "accepted": accepted}), flush=True)
    if not all(accepted.values()):
        raise SystemExit("not all signals were accepted")


if __name__ == "__main__":
    main()
