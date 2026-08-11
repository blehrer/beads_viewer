#!/usr/bin/env bash
# Regenerate README docgen sections and screenshot assets.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"
make readme-docgen
make screenshots
