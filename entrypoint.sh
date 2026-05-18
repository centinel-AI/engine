#!/usr/bin/env bash
set -euo pipefail

if [[ -x /app/workspace/sync-backend-config.sh ]]; then
  bash /app/workspace/sync-backend-config.sh
fi

ENGINE_BIN="$(tr -d '\n\r' < /app/workspace/.iac_engine_bin)"

args=("$@")
if [[ ${#args[@]} -gt 0 && ( ${args[0]} == terraform || ${args[0]} == tofu ) ]]; then
  args[0]="$ENGINE_BIN"
fi

"${ENGINE_BIN}" init -reconfigure
exec "${args[@]}"
