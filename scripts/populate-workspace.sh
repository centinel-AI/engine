#!/usr/bin/env bash
# Populate ./workspace like /app/workspace in the Docker image for the given cloud (azure | aws | gcp | oci).
# Optional second argument (or env PROJECT): copies data/<provider>/<PROJECT>/* → workspace/data/ (flat; not workspace/data/<provider>/<PROJECT>).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
PROVIDER="${1:?usage: populate-workspace.sh <azure|aws|gcp|oci|ovh> [PROJECT]}"
case "$PROVIDER" in azure|aws|gcp|oci|ovh) ;; *)
  echo "populate-workspace.sh: provider must be azure, aws, gcp, oci, or ovh" >&2
  exit 1
  ;;
esac

SRC="$ROOT/providers/$PROVIDER"
DST="$ROOT/workspace"
# Binary names on PATH: terraform | tofu (OpenTofu). Task exports ENGINE_IAC_ENGINE=terraform|opentofu.
case "${ENGINE_IAC_ENGINE:-terraform}" in
  terraform) ENGINE_LABEL=terraform ;;
  opentofu|tofu) ENGINE_LABEL=tofu ;;
  *)
    echo "populate-workspace.sh: ENGINE_IAC_ENGINE must be terraform, opentofu, or tofu (got '${ENGINE_IAC_ENGINE:-}')" >&2
    exit 1
    ;;
esac

PROJECT="${2:-${PROJECT:-}}"
if [[ -z "$PROJECT" ]]; then
  if [[ -f "$ROOT/.env" ]]; then
    set -a
    # shellcheck disable=SC1091
    source "$ROOT/.env"
    set +a
  fi
  PROJECT="${TF_VAR_project:-${ENGINE_TF_VAR_project:-project-01}}"
fi

DATA_SRC="$ROOT/data/$PROVIDER/$PROJECT"
DATA_DST="$DST/data"

mkdir -p "$DST"

if command -v rsync >/dev/null 2>&1; then
  rsync -a --delete \
    --exclude='.terraform/' \
    --exclude='.terraform.lock.hcl' \
    --exclude='backend.remote.tf.json' \
    "$SRC/" "$DST/"
else
  mkdir -p "$DST"
  cp -af "$SRC"/. "$DST"/
fi

printf '%s' "$ENGINE_LABEL" > "$DST/.iac_engine_bin"
printf '%s' "$PROVIDER" > "$DST/.grauss_provider"

# Resolve placeholders in _*.tf.json (required_version uses OPENTOFU vs Terraform pins; reads .iac_engine_bin).
bash "$ROOT/scripts/inject-provider-versions.sh" "$DST"

if [[ ! -d "$DATA_SRC" ]]; then
  echo "populate-workspace.sh: missing data directory $DATA_SRC" >&2
  exit 1
fi

rm -rf "${DST:?}/data"
mkdir -p "$DATA_DST"
if command -v rsync >/dev/null 2>&1; then
  rsync -a "$DATA_SRC/" "$DATA_DST/"
else
  cp -af "$DATA_SRC"/. "$DATA_DST"/
fi

touch "$DST/.gitkeep"

echo "Workspace populated for $PROVIDER at $DST (project: $PROJECT, engine: $ENGINE_LABEL). Data root: $DATA_DST — run: cd workspace && $ENGINE_LABEL init"
