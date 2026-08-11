# README screenshots

Automated captures of the four TUI views shown in the project README.

## Prerequisites

[VHS](https://github.com/charmbracelet/vhs) records the terminal session and writes PNG frames. It is pinned as a Go **tool** dependency in the repo root `go.mod` (Go 1.24+ `tool` block), so you do not install it globally:

```bash
# one-time (or after pulling go.mod tool changes)
go get -tool github.com/charmbracelet/vhs@v0.11.0

# runtime dependency for PNG/GIF encoding (not a Go module)
brew install ffmpeg   # macOS
# apt install ffmpeg  # Debian/Ubuntu
```

`make screenshots` invokes `go tool vhs`, which uses the module-pinned version.

Python 3 with Pillow is used for WebP conversion (`convert_webp.py` uses PEP 723 inline deps).

## Regenerate

From the repo root:

```bash
make screenshots
```

This will:

1. Copy `tests/testdata/synthetic_complex.jsonl` into `screenshots/fixture/.beads/issues.jsonl`
2. Build `./bv`
3. Run `screenshots/capture.tape` via `go tool vhs` (120×40 terminal, `--debug-render` per view)
4. Convert PNGs to WebP via `convert_webp.py`

`--debug-render` forces sync loading (`BV_BACKGROUND_MODE=0`) so captures never show the async "Loading beads..." spinner.

Outputs (git-tracked WebP, git-ignored intermediates):

| File | View | Key |
|------|------|-----|
| `screenshot_01__main_screen.webp` | List / split view | (default) |
| `screenshot_02__insights_view.webp` | Insights dashboard | `i` |
| `screenshot_03__kanban_view.webp` | Kanban board | `b` |
| `screenshot_04__graph_view.webp` | Graph view | `g` |

## Without VHS

If `go tool vhs` is unavailable (tool deps not fetched, or `ffmpeg` missing), `make screenshots` falls back to Go `--debug-render` and writes `screenshots/*.ansi` text captures. CI skips PNG generation when VHS cannot run (exit 0).

To skip the fallback and leave existing assets unchanged:

```bash
BV_SCREENSHOTS_FALLBACK=0 make screenshots
```

Manual single-view capture:

```bash
./bv --db screenshots/fixture/.beads --debug-render graph --debug-width 120 --debug-height 40 --theme dark
```

After changing robot commands or TUI keybindings, run `make readme-docgen` from the repo root to refresh the generated README sections.
