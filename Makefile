# bv Makefile
#
# Build with SQLite FTS5 (full-text search) support enabled

.PHONY: build install clean test screenshots readme-docgen docs

# Enable FTS5 for full-text search in SQLite exports
export CGO_CFLAGS := -DSQLITE_ENABLE_FTS5

build:
	go build -o bv ./cmd/bv

install:
	go install ./cmd/bv

clean:
	rm -f bv
	go clean

test:
	go test ./...

screenshots:
	@chmod +x scripts/capture_screenshots.sh
	@scripts/capture_screenshots.sh

readme-docgen:
	@chmod +x scripts/sync_readme_docgen.sh
	@scripts/sync_readme_docgen.sh

docs: readme-docgen screenshots
