#!/usr/bin/env bash
# Splice bv --robot-docgen fragments into README.md between marker pairs.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

README="${1:-README.md}"

if [[ ! -f "$README" ]]; then
	echo "README not found: $README" >&2
	exit 1
fi

go build -o bv ./cmd/bv

TMPDIR="$(mktemp -d)"
trap 'rm -rf "$TMPDIR"' EXIT

for section in agent-blurb robot-commands keybindings; do
	./bv --robot-docgen="$section" >"$TMPDIR/${section}.md"
done

python3 - "$README" "$TMPDIR" <<'PY'
import re
import sys
from pathlib import Path

readme_path = Path(sys.argv[1])
fragments_dir = Path(sys.argv[2])
text = readme_path.read_text(encoding="utf-8")

def splice_plain(section: str, content: str) -> None:
    global text
    start = f"<!-- bv-docgen:{section} -->"
    end = f"<!-- /bv-docgen:{section} -->"
    pattern = re.compile(re.escape(start) + r"\n.*?\n" + re.escape(end), re.DOTALL)
    if not pattern.search(text):
        raise SystemExit(f"missing marker pair for {section} in {readme_path}")
    text = pattern.sub(f"{start}\n{content}\n{end}", text, count=1)

def splice_fenced(section: str, content: str) -> None:
    global text
    start = f"<!-- bv-docgen:{section} -->"
    end = f"<!-- /bv-docgen:{section} -->"
    pattern = re.compile(
        re.escape(start) + r"\n```\n.*?\n```\n" + re.escape(end),
        re.DOTALL,
    )
    if not pattern.search(text):
        raise SystemExit(f"missing fenced marker pair for {section} in {readme_path}")
    text = pattern.sub(f"{start}\n```\n{content}\n```\n{end}", text, count=1)

splice_fenced(
    "agent-blurb",
    fragments_dir.joinpath("agent-blurb.md").read_text(encoding="utf-8").rstrip(),
)
for section in ("robot-commands", "keybindings"):
    splice_plain(
        section,
        fragments_dir.joinpath(f"{section}.md").read_text(encoding="utf-8").rstrip(),
    )

readme_path.write_text(text, encoding="utf-8")
print(f"Synced docgen sections into {readme_path}")
PY
