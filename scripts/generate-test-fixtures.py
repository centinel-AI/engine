#!/usr/bin/env python3
"""Generate minimal Terraform test fixture JSONs from a provider schema.

Analogous to codegen.py but writes JSON data files instead of .tf templates.
Reads the output of `terraform providers schema -json` and emits one fixture
JSON per resource type under <out>/<short_name>/fixture.json.

Each fixture contains:
  - "_resource_type": full resource type (e.g. "azurerm_linux_virtual_machine")
  - Required scalar attributes filled with minimal type-appropriate values
  - Required blocks (min_items >= 1) with minimal nested content
  - Overrides for fields that need specific values beyond type defaults

Usage:
  generate-test-fixtures.py \\
    --schema tmp/codegen/<cloud>/schema.json \\
    --provider-source hashicorp/azurerm \\
    --prefix azurerm_ \\
    --out data/<cloud>/_test-all \\
    [--strip-prefix] \\
    [--exclude res1,res2] \\
    [--overrides scripts/codegen-overrides/<cloud>.fixture-overrides.json]
"""

from __future__ import annotations

import argparse
import json
import sys
from pathlib import Path
from typing import Any


def is_skippable(attr: dict) -> bool:
    """Computed-only attributes (not settable by config) are skipped."""
    return attr.get("computed", False) and not (
        attr.get("required", False) or attr.get("optional", False)
    )


def minimal_value(type_schema: Any, field_name: str, overrides: dict) -> Any:
    """Return a minimal valid value for an attribute given its type schema."""
    if field_name in overrides:
        return overrides[field_name]
    if isinstance(type_schema, str):
        if type_schema == "string":
            return "test"
        if type_schema == "number":
            return 0
        if type_schema == "bool":
            return False
        return "test"
    if isinstance(type_schema, list) and len(type_schema) >= 2:
        container, inner = type_schema[0], type_schema[1]
        if container in ("list", "set"):
            return [minimal_value(inner, field_name, overrides)]
        if container == "map":
            return {}
        if container == "object" and isinstance(inner, dict):
            return {k: minimal_value(v, k, overrides) for k, v in inner.items()}
        if container == "tuple" and isinstance(inner, list):
            return [minimal_value(t, field_name, overrides) for t in inner]
    return "test"


def minimal_block(block_schema: dict, overrides: dict) -> dict:
    """Recursively build a minimal object satisfying a block's required fields."""
    obj: dict = {}
    for name, attr in sorted((block_schema.get("attributes") or {}).items()):
        if is_skippable(attr):
            continue
        if attr.get("required"):
            obj[name] = minimal_value(attr.get("type", "string"), name, overrides)
    for name, block_type in sorted((block_schema.get("block_types") or {}).items()):
        if block_type.get("min_items", 0) >= 1:
            obj[name] = [minimal_block(block_type.get("block") or {}, overrides)]
    return obj


def generate_fixture(resource_name: str, resource_schema: dict, overrides: dict) -> dict:
    """Build a minimal fixture JSON for a resource type."""
    block = resource_schema.get("block") or {}
    attrs = block.get("attributes") or {}
    block_types = block.get("block_types") or {}

    fixture: dict = {
        "_resource_type": resource_name,
        "name": overrides.get("name", "test-fixture"),
    }

    for name, attr in sorted(attrs.items()):
        if name in ("id", "name") or is_skippable(attr):
            continue
        if attr.get("required"):
            fixture[name] = minimal_value(attr.get("type", "string"), name, overrides)

    for name, block_type in sorted(block_types.items()):
        if block_type.get("min_items", 0) >= 1:
            fixture[name] = [minimal_block(block_type.get("block") or {}, overrides)]

    return fixture


def find_provider_key(schema: dict, source: str) -> str | None:
    keys = list((schema.get("provider_schemas") or {}).keys())
    for k in keys:
        if k.endswith("/" + source):
            return k
    for k in keys:
        if source in k:
            return k
    return None


def main() -> int:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--schema", required=True, help="path to provider schema JSON")
    p.add_argument("--provider-source", required=True,
                   help='provider source, e.g. "hashicorp/azurerm", "oracle/oci"')
    p.add_argument("--prefix", required=True,
                   help='resource name prefix, e.g. "azurerm_", "aws_"')
    p.add_argument("--out", required=True,
                   help="output directory (data/<cloud>/_test-all/)")
    p.add_argument("--strip-prefix", action="store_true",
                   help="strip provider prefix from directory name (mirrors codegen --strip-prefix)")
    p.add_argument("--exclude", default="",
                   help="comma-separated list of resource names to skip")
    p.add_argument("--overrides", default=None,
                   help="path to JSON file mapping field names to specific override values")
    args = p.parse_args()

    try:
        schema: dict = json.loads(Path(args.schema).read_text())
    except FileNotFoundError:
        print(f"generate-test-fixtures: schema not found: {args.schema}", file=sys.stderr)
        return 2

    key = find_provider_key(schema, args.provider_source)
    if not key:
        print(
            f"generate-test-fixtures: provider '{args.provider_source}' not in schema; "
            f"known: {list((schema.get('provider_schemas') or {}).keys())}",
            file=sys.stderr,
        )
        return 2

    resources = schema["provider_schemas"][key].get("resource_schemas") or {}
    exclude = {x.strip() for x in args.exclude.split(",") if x.strip()}

    overrides: dict = {}
    if args.overrides:
        try:
            overrides = json.loads(Path(args.overrides).read_text())
        except FileNotFoundError:
            pass

    out_dir = Path(args.out)
    out_dir.mkdir(parents=True, exist_ok=True)

    written = 0
    for resource_name in sorted(resources):
        if not resource_name.startswith(args.prefix):
            continue
        if resource_name in exclude:
            continue
        if resources[resource_name].get("block", {}).get("deprecated"):
            continue

        short = resource_name[len(args.prefix):] if args.strip_prefix else resource_name
        fixture_dir = out_dir / short
        fixture_dir.mkdir(exist_ok=True)

        fixture = generate_fixture(resource_name, resources[resource_name], overrides)
        (fixture_dir / "fixture.json").write_text(json.dumps(fixture, indent=2) + "\n")
        written += 1

    print(
        f"generate-test-fixtures: wrote {written} fixtures to {out_dir} "
        f"(provider: {args.provider_source}, prefix: {args.prefix})"
    )
    return 0


if __name__ == "__main__":
    sys.exit(main())
