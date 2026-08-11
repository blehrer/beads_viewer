#!/usr/bin/env bash
# Capture README TUI screenshots (VHS → PNG → WebP).
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

FIXTURE_SRC="tests/testdata/synthetic_complex.jsonl"
FIXTURE_DIR="screenshots/fixture/.beads"

prepare_fixture() {
	mkdir -p "$FIXTURE_DIR"
	cp "$FIXTURE_SRC" "$FIXTURE_DIR/issues.jsonl"
}

build_bv() {
	go build -o bv ./cmd/bv
}

convert_webp() {
	if ! command -v python3 >/dev/null 2>&1; then
		echo "python3 not found; PNG screenshots saved (skipping WebP conversion)"
		return 0
	fi
	(
		cd screenshots
		python3 convert_webp.py
	)
}

capture_with_vhs() {
	if command -v fc-list >/dev/null 2>&1; then
		if ! fc-list : family | grep -qi 'GeistMono Nerd Font'; then
			echo "warning: GeistMono Nerd Font not found; VHS may render missing icon glyphs" >&2
			echo "  macOS: brew install --cask font-geist-mono-nerd-font" >&2
		fi
	fi
	go tool vhs screenshots/capture.tape
	rm -f screenshots/.capture.gif
	convert_webp
	echo "Screenshots regenerated in screenshots/ (PNG + WebP)"
}

capture_with_go_fallback() {
	echo "go tool vhs unavailable — using Go --debug-render fallback (ANSI text, not PNG/WebP)"
	echo "VHS is declared in go.mod; run: go get -tool github.com/charmbracelet/vhs@v0.11.0"
	prepare_fixture
	build_bv

	capture_view() {
		local view="$1"
		local base="$2"
		BV_NO_BROWSER=1 BV_NO_SAVED_CONFIG=1 BV_TEST_MODE=1 BV_BACKGROUND_MODE=0 \
			./bv \
			--db "$FIXTURE_DIR" \
			--theme dark \
			--debug-render "$view" \
			--debug-width 120 \
			--debug-height 40 \
			>"screenshots/${base}.ansi"
		echo "  wrote screenshots/${base}.ansi"
	}

	capture_view list screenshot_01__main_screen
	capture_view insights screenshot_02__insights_view
	capture_view board screenshot_03__kanban_view
	capture_view graph screenshot_04__graph_view
}

prepare_fixture
build_bv

if go tool -n vhs >/dev/null 2>&1; then
	capture_with_vhs
elif [[ "${BV_SCREENSHOTS_FALLBACK:-1}" == "1" ]]; then
	capture_with_go_fallback
else
	echo "go tool vhs unavailable; skipping screenshot capture (existing WebP unchanged)"
	echo "Fetch tool deps: go get -tool github.com/charmbracelet/vhs@v0.11.0"
fi
