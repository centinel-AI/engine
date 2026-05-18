#!/usr/bin/env bash
# Replaces __PLACEHOLDER__ tokens in <target>/_*.tf.json with version constraints from the environment.
# Known placeholders: __AZURERM_VERSION__, __AZUREAD_VERSION__, __AWS_PROVIDER_VERSION__, __GOOGLE_PROVIDER_VERSION__, __OCI_PROVIDER_VERSION__, __REQUIRED_VERSION__.
#
# __REQUIRED_VERSION__ is derived from ENGINE_TERRAFORM_VERSION (terraform) or ENGINE_OPENTOFU_VERSION (tofu)
# by prefixing "=". When called from Docker, TERRAFORM_REQUIRED_VERSION / OPENTOFU_REQUIRED_VERSION are
# already set as build ARGs (with "=" prefix) and take precedence. Engine is detected from ENGINE_IAC_ENGINE
# if set, otherwise from <target>/.iac_engine_bin (written before this script in populate-workspace).
#
# - Local: bash scripts/inject-provider-versions.sh ./workspace
#   Loads repo .env when present (repo root = parent of scripts/).
# - Docker: /tmp/inject-provider-versions.sh /app/workspace
#   Uses build ARGs exported for the RUN layer (no .env in the image).
#
# Files under providers/ stay with placeholders in git; only the copy under workspace/ or /app/workspace is rewritten.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

TARGET_DIR="${1:?usage: inject-provider-versions.sh <target-directory>}"
TARGET_DIR="$(cd "$TARGET_DIR" && pwd)"

if [[ -f "$REPO_ROOT/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$REPO_ROOT/.env"
  set +a
fi

TERRAFORM_REQUIRED_VERSION="${TERRAFORM_REQUIRED_VERSION:-=${ENGINE_TERRAFORM_VERSION:-1.15.3}}"
OPENTOFU_REQUIRED_VERSION="${OPENTOFU_REQUIRED_VERSION:-=${ENGINE_OPENTOFU_VERSION:-1.11.7}}"
AWS_PROVIDER_VERSION="${ENGINE_AWS_PROVIDER_VERSION:-${AWS_PROVIDER_VERSION:-=6.44.0}}"
GOOGLE_PROVIDER_VERSION="${ENGINE_GOOGLE_PROVIDER_VERSION:-${GOOGLE_PROVIDER_VERSION:-=7.32.0}}"
AZURERM_VERSION="${ENGINE_AZURERM_VERSION:-${AZURERM_VERSION:-=4.72.0}}"
AZUREAD_VERSION="${ENGINE_AZUREAD_VERSION:-${AZUREAD_VERSION:-=3.8.0}}"
OCI_PROVIDER_VERSION="${ENGINE_OCI_PROVIDER_VERSION:-${OCI_PROVIDER_VERSION:-=8.13.0}}"

resolve_engine() {
  if [[ -n "${ENGINE_IAC_ENGINE:-}" ]]; then
    case "${ENGINE_IAC_ENGINE}" in
      terraform) echo terraform ;;
      opentofu|tofu) echo tofu ;;
      *) echo terraform ;;
    esac
    return
  fi
  if [[ -f "$TARGET_DIR/.iac_engine_bin" ]]; then
    tr -d '\n\r' < "$TARGET_DIR/.iac_engine_bin"
    return
  fi
  echo terraform
}

ENGINE_BIN="$(resolve_engine)"
if [[ "$ENGINE_BIN" == "terraform" ]]; then
  REQUIRED_VERSION="$TERRAFORM_REQUIRED_VERSION"
else
  REQUIRED_VERSION="$OPENTOFU_REQUIRED_VERSION"
fi

substitute_file() {
  local f="$1"
  [[ -f "$f" ]] || return 0
  local tmp
  tmp="$(mktemp)"
  sed \
    -e "s|__AZURERM_VERSION__|${AZURERM_VERSION:-}|g" \
    -e "s|__AZUREAD_VERSION__|${AZUREAD_VERSION:-}|g" \
    -e "s|__GOOGLE_PROVIDER_VERSION__|${GOOGLE_PROVIDER_VERSION:-}|g" \
    -e "s|__AWS_PROVIDER_VERSION__|${AWS_PROVIDER_VERSION:-}|g" \
    -e "s|__OCI_PROVIDER_VERSION__|${OCI_PROVIDER_VERSION:-}|g" \
    -e "s|__REQUIRED_VERSION__|${REQUIRED_VERSION:-}|g" \
    "$f" > "$tmp"
  mv "$tmp" "$f"
}

shopt -s nullglob
count=0
for f in "$TARGET_DIR"/_*.tf.json; do
  [[ -f "$f" ]] || continue
  substitute_file "$f"
  count=$((count + 1))
done
shopt -u nullglob

if [[ "$count" -eq 0 ]]; then
  echo "inject-provider-versions: no _*.tf.json in $TARGET_DIR" >&2
  exit 1
fi

echo "inject-provider-versions: resolved placeholders in $count file(s) under $TARGET_DIR (engine: $ENGINE_BIN, required_version: $REQUIRED_VERSION)"
