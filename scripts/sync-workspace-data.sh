#!/usr/bin/env bash
# Copy data/<cloud>/<project>/ → workspace/data/ without re-populating providers.
# Usage: sync-workspace-data.sh <azure|aws|gcp|oci> <project>
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
CLOUD="${1:?usage: sync-workspace-data.sh <azure|aws|gcp|oci> <project>}"
PROJECT="${2:?usage: sync-workspace-data.sh <azure|aws|gcp|oci> <project>}"

SRC="$ROOT/data/$CLOUD/$PROJECT"
DST="$ROOT/workspace/data"

[ -d "$SRC" ] || { echo "sync-workspace-data: data directory not found: $SRC" >&2; exit 1; }

rm -rf "$DST"
mkdir -p "$DST"

if command -v rsync >/dev/null 2>&1; then
  rsync -a "$SRC/" "$DST/"
else
  cp -af "$SRC/". "$DST/"
fi

echo "sync-workspace-data: data/$CLOUD/$PROJECT → workspace/data/ (project: $PROJECT)"
