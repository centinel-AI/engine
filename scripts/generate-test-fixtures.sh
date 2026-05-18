#!/usr/bin/env bash
# Generate minimal Terraform test fixtures for all resource types in a cloud.
# Reads tmp/codegen/<cloud>/schema.json (produced by generate-resources.sh) and
# writes one fixture JSON per resource type to data/<cloud>/_test-all/.
# Applies the same exclusions as generate-resources.sh.
#
# Usage: generate-test-fixtures.sh <azure|aws|gcp|oci>
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CLOUD="${1:?usage: generate-test-fixtures.sh <azure|aws|gcp|oci>}"
SCHEMA="$ROOT/tmp/codegen/$CLOUD/schema.json"
OUT="$ROOT/data/$CLOUD/_test-all"
OVERRIDES="$SCRIPT_DIR/codegen-overrides/$CLOUD.fixture-overrides.json"
EXCLUDE_FILE="$SCRIPT_DIR/codegen-overrides/$CLOUD.exclude.txt"

[ -f "$SCHEMA" ] || {
  echo "generate-test-fixtures: schema not found: $SCHEMA" >&2
  echo "  Run: task generate:$CLOUD" >&2
  exit 1
}

EXCLUDE_LIST=""
if [ -f "$EXCLUDE_FILE" ]; then
  EXCLUDE_LIST=$(grep -v '^\s*#' "$EXCLUDE_FILE" | grep -v '^\s*$' | tr '\n' ',' | sed 's/,$//' || true)
fi

COMMON=(
  --schema "$SCHEMA"
  --out "$OUT"
  --exclude "$EXCLUDE_LIST"
)
[ -f "$OVERRIDES" ] && COMMON+=(--overrides "$OVERRIDES")

rm -rf "$OUT"

case "$CLOUD" in
  azure)
    python3 "$SCRIPT_DIR/generate-test-fixtures.py" "${COMMON[@]}" \
      --provider-source "hashicorp/azurerm" --prefix "azurerm_" --strip-prefix
    python3 "$SCRIPT_DIR/generate-test-fixtures.py" "${COMMON[@]}" \
      --provider-source "hashicorp/azuread" --prefix "azuread_"
    ;;
  aws)
    python3 "$SCRIPT_DIR/generate-test-fixtures.py" "${COMMON[@]}" \
      --provider-source "hashicorp/aws" --prefix "aws_" --strip-prefix
    ;;
  gcp)
    python3 "$SCRIPT_DIR/generate-test-fixtures.py" "${COMMON[@]}" \
      --provider-source "hashicorp/google" --prefix "google_" --strip-prefix
    ;;
  oci)
    python3 "$SCRIPT_DIR/generate-test-fixtures.py" "${COMMON[@]}" \
      --provider-source "oracle/oci" --prefix "oci_" --strip-prefix
    ;;
  *)
    echo "generate-test-fixtures: unknown cloud '$CLOUD'" >&2; exit 1 ;;
esac
