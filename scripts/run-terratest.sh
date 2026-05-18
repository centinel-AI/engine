#!/usr/bin/env bash
# Run Terratest plan-only tests for a single cloud.
# Tests skip automatically when the required credential env vars are absent.
# Usage: run-terratest.sh <azure|aws|gcp|oci> [exhaustive]
#
# GOROOT is unset to prevent a goenv/homebrew version mismatch where goenv sets
# GOROOT to its managed installation while a different go binary leads PATH.
set -euo pipefail

CLOUD="${1:?usage: run-terratest.sh <azure|aws|gcp|oci> [exhaustive]}"
MODE="${2:-}"
CAPITALIZED="Test$(echo "$CLOUD" | awk '{print toupper(substr($0,1,1)) substr($0,2)}')"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

unset GOROOT
cd "$ROOT/tests"

if [ "$MODE" = "exhaustive" ]; then
  go test -v -timeout 60m -tags exhaustive -run "${CAPITALIZED}Exhaustive" ./...
else
  go test -v -timeout 30m -run "$CAPITALIZED" ./...
fi
