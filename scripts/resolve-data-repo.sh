#!/usr/bin/env bash
# Resolve ENGINE_DATA_REPO to an absolute path (sibling data repository by default).
# Source from other scripts: . "$(dirname "$0")/resolve-data-repo.sh"
set -euo pipefail

_ENGINE_ROOT="${ENGINE_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)}"

if [[ -f "$_ENGINE_ROOT/.env" ]]; then
  set -a
  # shellcheck disable=SC1091
  source "$_ENGINE_ROOT/.env"
  set +a
fi

_REL="${ENGINE_DATA_REPO_PATH:-../data}"
if [[ "$_REL" = /* ]]; then
  ENGINE_DATA_REPO="$_REL"
else
  ENGINE_DATA_REPO="$(cd "$_ENGINE_ROOT/$_REL" && pwd)"
fi

export ENGINE_ROOT="$_ENGINE_ROOT"
export ENGINE_DATA_REPO
