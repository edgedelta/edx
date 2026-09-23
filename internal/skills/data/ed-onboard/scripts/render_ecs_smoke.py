#!/usr/bin/env python3
"""Render the disposable ECS application and Edge Delta sidecar test stack."""
import argparse
import base64
import json
from pathlib import Path


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--output", required=True, help="output CloudFormation JSON path")
    parser.add_argument("--with-file-logs", action="store_true", help="also exercise native shared-file collection")
    args = parser.parse_args()
    root = Path(__file__).resolve().parent.parent
    template = json.loads((root / "assets/ecs-smoke-cloudformation.json").read_text())
    source = base64.b64encode((root / "scripts/emit_otlp.py").read_bytes()).decode()
    containers = template["Resources"]["Task"]["Properties"]["ContainerDefinitions"]
    application = next(c for c in containers if c["Name"] == "application")
    application["Command"] = ["import base64;exec(base64.b64decode(" + repr(source) + "))"]
    if args.with_file_logs:
        task = template["Resources"]["Task"]["Properties"]
        task["Volumes"] = [{"Name": "logs"}]
        for container in containers:
            container["MountPoints"] = [{"SourceVolume": "logs", "ContainerPath": "/var/log/onboard",
                                          "ReadOnly": container["Name"] == "edgedelta"}]
        application["Environment"].append({"Name": "ONBOARD_LOG_FILE", "Value": "/var/log/onboard/application.log"})
    Path(args.output).write_text(json.dumps(template, indent=2) + "\n")
    print("Rendered temporary ECS stack; no CloudWatch resources or log drivers.")


if __name__ == "__main__":
    main()
