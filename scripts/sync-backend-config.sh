#!/usr/bin/env bash
# Copies <ENGINE_DATA_REPO>/<provider>/<project>/backend.tf.json → <module>/backend.remote.tf.json
# provider: Docker uses ENGINE_PROVIDER (image ENV); local workspace uses .grauss_provider; repo-root sync uses ENGINE_PROVIDER=azure|aws|gcp|oci.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
# shellcheck disable=SC1091
source "$REPO_ROOT/scripts/resolve-data-repo.sh"

PROJECT="${TF_VAR_project:-${ENGINE_TF_VAR_project:-project-01}}"

if [[ -f /app/workspace/.iac_engine_bin ]]; then
  MODULE_DIR="/app/workspace"
  PROVIDER="${ENGINE_PROVIDER:?sync-backend-config: set ENGINE_PROVIDER in the image (azure, aws, gcp, or oci)}"
elif [[ -d "$REPO_ROOT/workspace" ]] && [[ "$(pwd -P)" == "$(cd "$REPO_ROOT/workspace" && pwd -P)" ]]; then
  MODULE_DIR="$REPO_ROOT/workspace"
  if [[ ! -f "$MODULE_DIR/.grauss_provider" ]]; then
    echo "sync-backend-config: workspace layout not recognized" >&2
    exit 1
  fi
  PROVIDER="$(tr -d '\n\r' < "$MODULE_DIR/.grauss_provider")"
elif [[ -n "${ENGINE_PROVIDER:-}" ]]; then
  PROVIDER="$ENGINE_PROVIDER"
  MODULE_DIR="$REPO_ROOT/providers/$PROVIDER"
else
  PROVIDER=azure
  MODULE_DIR="$REPO_ROOT/providers/azure"
fi

case "$PROVIDER" in
  azure|aws|gcp|oci) ;;
  *) echo "sync-backend-config: invalid provider '$PROVIDER'" >&2; exit 1 ;;
esac

# Docker: full data tree at /app/workspace/data/<provider>/<project>/.
# Host workspace: task populate copies project JSON into workspace/data/ (flat) — backend.tf.json is workspace/data/backend.tf.json.
# Repo providers/<provider>/: always read from repo data tree.
if [[ "$MODULE_DIR" == "/app/workspace" ]]; then
  SRC="$MODULE_DIR/data/$PROVIDER/$PROJECT/backend.tf.json"
elif [[ "$MODULE_DIR" == "$REPO_ROOT/workspace" ]]; then
  SRC="$MODULE_DIR/data/backend.tf.json"
else
  SRC="$ENGINE_DATA_REPO/$PROVIDER/$PROJECT/backend.tf.json"
fi
DST="$MODULE_DIR/backend.remote.tf.json"

if [[ ! -f "$SRC" ]]; then
  rm -f "$DST"
  exit 0
fi

cp "$SRC" "$DST"
