#!/usr/bin/env bash
# Reset providers/{aws,azure,gcp,oci} and workspace/ to .gitkeep only; remove tmp/ and _test-all/ fixtures.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"

for cloud in aws azure gcp oci ovh; do
  dir="$ROOT/providers/$cloud"
  mkdir -p "$dir"
  find "$dir" -mindepth 1 -maxdepth 1 ! -name '.gitkeep' -exec rm -rf {} +
  touch "$dir/.gitkeep"
  rm -rf "$ROOT/data/$cloud/_test-all"
done

ws="$ROOT/workspace"
if [ -d "$ws" ]; then
  find "$ws" -mindepth 1 -maxdepth 1 ! -name '.gitkeep' -exec rm -rf {} +
fi

rm -rf "$ROOT/tmp"

echo "clean: providers/{aws,azure,gcp,oci,ovh} and workspace/ reset to .gitkeep; data/*/_test-all/ and tmp/ removed."
