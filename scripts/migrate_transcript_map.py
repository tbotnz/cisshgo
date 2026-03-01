#!/usr/bin/env python3
"""
migrate_transcript_map.py — Migrate transcript_map.yaml from v0.2.0 to v1.0.0 format.

The v0.2.0 schema used a list of single-key maps for platforms:

    platforms:
      - csr1000v:
          vendor: "cisco"
      - ios:
          vendor: "cisco"

The v1.0.0 schema uses a plain map:

    platforms:
      csr1000v:
        vendor: "cisco"
      ios:
        vendor: "cisco"

Usage:
    # Migrate in-place (writes to stdout):
    python3 scripts/migrate_transcript_map.py < transcript_map.yaml > transcript_map_new.yaml

    # Or overwrite in-place using a temp file:
    python3 scripts/migrate_transcript_map.py < transcript_map.yaml > /tmp/map.yaml && mv /tmp/map.yaml transcript_map.yaml

    # Migrate a specific file:
    python3 scripts/migrate_transcript_map.py path/to/transcript_map.yaml

Exit codes:
    0 — success (migrated or already up-to-date)
    1 — error (invalid YAML, unreadable file, etc.)
"""

import sys
import yaml


def migrate(data: dict) -> tuple[dict, bool]:
    """
    Migrate platforms from list-of-maps to map format.
    Returns (migrated_data, was_changed).
    """
    platforms = data.get("platforms")
    if not isinstance(platforms, list):
        return data, False  # already map format or missing — no change needed

    new_platforms = {}
    for item in platforms:
        if not isinstance(item, dict) or len(item) != 1:
            print(
                f"Warning: unexpected platforms entry {item!r}, skipping",
                file=sys.stderr,
            )
            continue
        for key, value in item.items():
            new_platforms[key] = value

    data["platforms"] = new_platforms
    return data, True


def main() -> int:
    if len(sys.argv) > 1:
        path = sys.argv[1]
        try:
            with open(path) as f:
                raw = f.read()
        except OSError as e:
            print(f"Error reading {path}: {e}", file=sys.stderr)
            return 1
    else:
        raw = sys.stdin.read()

    try:
        data = yaml.safe_load(raw)
    except yaml.YAMLError as e:
        print(f"Error parsing YAML: {e}", file=sys.stderr)
        return 1

    if not isinstance(data, dict):
        print("Error: expected a YAML mapping at the top level", file=sys.stderr)
        return 1

    migrated, changed = migrate(data)

    if changed:
        print("Migrated platforms from list-of-maps to map format.", file=sys.stderr)
    else:
        print("Already in v1.0.0 format — no changes needed.", file=sys.stderr)

    yaml.dump(migrated, sys.stdout, default_flow_style=False, allow_unicode=True)
    return 0


if __name__ == "__main__":
    sys.exit(main())
