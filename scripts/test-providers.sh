#!/usr/bin/env bash
# Validate generated provider templates for a cloud without touching the main workspace.
#
# Steps:
#   1. Assert providers/<cloud>/ has generated templates (non-underscore .tf files).
#   2. terraform fmt -check on providers/<cloud>/.
#   3. Build an isolated test workspace in tmp/test/<cloud>/.
#   4. Inject provider version placeholders.
#   5. terraform init -backend=false (uses TF_PLUGIN_CACHE_DIR to avoid re-downloading).
#   6. terraform validate.
#   7. Clean up tmp/test/<cloud>/.
#
# Usage:
#   bash scripts/test-providers.sh <azure|aws|gcp|oci>
#   ENGINE_IAC_ENGINE=opentofu bash scripts/test-providers.sh azure
#
# Environment:
#   ENGINE_IAC_ENGINE   terraform (default) | opentofu  — binary used for init/validate
#   TF_PLUGIN_CACHE_DIR path to a shared plugin cache (default: tmp/tf-plugin-cache)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

CLOUD="${1:?usage: test-providers.sh <azure|aws|gcp|oci>}"
case "$CLOUD" in azure|aws|gcp|oci) ;;
  *) echo "test-providers: cloud must be azure, aws, gcp, or oci" >&2; exit 1 ;;
esac

if [[ -f "$ROOT/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$ROOT/.env"
  set +a
fi

case "${ENGINE_IAC_ENGINE:-terraform}" in
  terraform)         BIN=terraform ;;
  opentofu|tofu)     BIN=tofu ;;
  *) echo "test-providers: ENGINE_IAC_ENGINE must be terraform or opentofu" >&2; exit 1 ;;
esac

PROVIDERS_DIR="$ROOT/providers/$CLOUD"
TEST_DIR="$ROOT/tmp/test/$CLOUD"
PLUGIN_CACHE="${TF_PLUGIN_CACHE_DIR:-$ROOT/tmp/tf-plugin-cache}"

log() { echo "[$CLOUD] $*"; }
fail() { echo "[$CLOUD] ERROR: $*" >&2; exit 1; }

# ── 1. Assert templates exist ──────────────────────────────────────────────────
count=$(find "$PROVIDERS_DIR" -maxdepth 1 -name "*.tf" ! -name "_*" 2>/dev/null | wc -l | tr -d ' ')
[[ "$count" -gt 0 ]] || fail "no generated templates in $PROVIDERS_DIR — run: task generate:$CLOUD"
log "$count generated templates found"

# ── 2. fmt check ──────────────────────────────────────────────────────────────
log "terraform fmt -check ..."
$BIN fmt -check -recursive "$PROVIDERS_DIR"
log "fmt OK"

# ── 3. Build isolated test workspace ──────────────────────────────────────────
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"

if command -v rsync >/dev/null 2>&1; then
  rsync -a \
    --exclude='.terraform/' \
    --exclude='.terraform.lock.hcl' \
    --exclude='backend.remote.tf.json' \
    "$PROVIDERS_DIR/" "$TEST_DIR/"
else
  cp -af "$PROVIDERS_DIR"/. "$TEST_DIR/"
  rm -rf "$TEST_DIR/.terraform" "$TEST_DIR/.terraform.lock.hcl" "$TEST_DIR/backend.remote.tf.json"
fi

printf '%s' "$BIN" > "$TEST_DIR/.iac_engine_bin"

# ── 4. Inject version placeholders ────────────────────────────────────────────
log "injecting provider version placeholders ..."
bash "$SCRIPT_DIR/inject-provider-versions.sh" "$TEST_DIR"

# ── 5. terraform init ─────────────────────────────────────────────────────────
mkdir -p "$PLUGIN_CACHE"
export TF_PLUGIN_CACHE_DIR="$PLUGIN_CACHE"

log "terraform init -backend=false (plugin cache: $PLUGIN_CACHE) ..."
(cd "$TEST_DIR" && $BIN init -input=false -backend=false 2>&1 | grep -E "^(Initializing|Terraform has|─|Error)" || true)
log "init OK"

# ── 6. terraform validate ─────────────────────────────────────────────────────
log "terraform validate ..."
(cd "$TEST_DIR" && $BIN validate)
log "validate OK"

# ── 7. Clean up ───────────────────────────────────────────────────────────────
rm -rf "$TEST_DIR"

log "ALL TESTS PASSED"
