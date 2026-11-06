#!/usr/bin/env python3
"""Prints a .env file as JSON for `vault kv put ... -` (read from stdin).

Pipe it straight into Vault so values never touch disk or shell history:
  python3 env-to-vault.py .env | docker exec -i -e VAULT_TOKEN adfluence_vault vault kv put -mount=secret adfluence/prod -

Skips comments, blank lines and VAULT_* (those stay on the server).
"""
import json
import re
import sys

KEY = re.compile(r"^[A-Z][A-Z0-9_]*$")


def parse(lines):
    out = {}
    for n, raw in enumerate(lines, 1):
        line = raw.strip()
        if not line or line.startswith("#"):
            continue
        if line.startswith("export "):
            line = line[len("export "):]
        key, sep, value = line.partition("=")
        key = key.strip()
        if not sep or not KEY.match(key):
            sys.exit(f"line {n}: not KEY=value (key must be CAPITALS_AND_UNDERSCORES)")
        if key.startswith("VAULT_"):
            continue
        value = value.strip()
        if len(value) >= 2 and value[0] == value[-1] and value[0] in "'\"":
            value = value[1:-1]
        out[key] = value
    return out


if __name__ == "__main__":
    if len(sys.argv) != 2:
        sys.exit("usage: env-to-vault.py path/to/.env")
    with open(sys.argv[1], encoding="utf-8") as f:
        values = parse(f)
    print(f"{len(values)} keys: {', '.join(sorted(values))}", file=sys.stderr)
    json.dump(values, sys.stdout)
