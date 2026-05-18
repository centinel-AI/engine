#!/usr/bin/env bash
# Bootstrap .env from .env.example. No-op when .env already exists.
set -euo pipefail

test -f .env.example || { echo "env-init: .env.example not found" >&2; exit 1; }

if test -f .env; then
  echo "env-init: .env already exists — not overwriting."
else
  cp .env.example .env
  echo "env-init: created .env from .env.example"
fi
