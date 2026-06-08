#!/usr/bin/env bash
# Build the two WASM artefacts for the ishawyha.dev Piton Lab.
#
#   piton-runner.wasm  — TinyGo, interpreter only (~220 KB, ~80 KB brotli)
#   piton-viz.wasm     — Go,     ДСТУ visualizer via rombik (~3.6 MB, ~683 KB brotli)
#
# The runner loads on first interaction; the visualizer is lazy-loaded only
# when the user clicks "Flowchart".
#
# Usage:
#     nix develop --command bash scripts/build-wasm.sh [OUT_DIR]
#
# OUT_DIR defaults to ./dist. Copy the four output files to your site's
# static/ directory after the build:
#
#     piton-runner.wasm
#     piton-viz.wasm
#     wasm_exec_tinygo.js   (TinyGo's glue)
#     wasm_exec.js          (standard Go's glue)

set -euo pipefail

OUT="${1:-dist}"
mkdir -p "$OUT"

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "==> Building runner with TinyGo (-opt=z, no debug)"
# TinyGo 0.40 caps at Go 1.25 but the rest of the project tracks newer Go.
# Temporarily downshift the go directive in go.mod for the tinygo build only,
# then restore it on exit (success, failure, or interrupt).
GOMOD_BACKUP="$(mktemp)"
cp go.mod "$GOMOD_BACKUP"
trap 'mv "$GOMOD_BACKUP" go.mod' EXIT INT TERM
sed -i -E 's/^go 1\.(2[6-9]|[3-9][0-9]).*/go 1.25/' go.mod

tinygo build \
    -target=wasm \
    -no-debug \
    -opt=z \
    -o "$OUT/piton-runner.wasm" \
    ./cmd/wasm-runner

# Restore go.mod immediately so the subsequent go build sees the real version.
mv "$GOMOD_BACKUP" go.mod
trap - EXIT INT TERM

echo "==> Building viz with standard Go (-ldflags=-s -w, -gcflags=-l -B, -trimpath)"
# -gcflags="all=-l -B": disable inlining + bounds checks to shave ~450 KB.
# Safe for this read-only visualizer where bounds errors would already be a
# bug from upstream (parser/AST never feeds invalid indices to rombik).
GOOS=js GOARCH=wasm go build \
    -ldflags="-s -w" \
    -gcflags="all=-l -B" \
    -trimpath \
    -o "$OUT/piton-viz.wasm" \
    ./cmd/wasm-viz

echo "==> Optimising both with wasm-opt"
wasm-opt -Oz --enable-bulk-memory --enable-nontrapping-float-to-int --enable-sign-ext \
    "$OUT/piton-runner.wasm" -o "$OUT/piton-runner.wasm.tmp"
mv "$OUT/piton-runner.wasm.tmp" "$OUT/piton-runner.wasm"

# Viz uses -O2, not -Oz: measured on the rombik-based build, -Oz yields a
# smaller raw file but compresses ~10 KB WORSE under brotli (which is what the
# site serves), so -O2 wins on the wire (~683 KB vs ~693 KB brotli).
wasm-opt -O2 --enable-bulk-memory --enable-nontrapping-float-to-int --enable-sign-ext \
    "$OUT/piton-viz.wasm" -o "$OUT/piton-viz.wasm.tmp"
mv "$OUT/piton-viz.wasm.tmp" "$OUT/piton-viz.wasm"

echo "==> Copying wasm_exec.js variants"
cp "$(tinygo env TINYGOROOT)/targets/wasm_exec.js" "$OUT/wasm_exec_tinygo.js"
cp "$(go env GOROOT)/lib/wasm/wasm_exec.js" "$OUT/wasm_exec.js"

echo
echo "==> Done. Sizes:"
ls -la "$OUT/" | awk 'NR>1 {printf "  %10s  %s\n", $5, $NF}'

if command -v brotli >/dev/null 2>&1; then
    echo
    echo "==> Brotli-compressed projection (-q 11):"
    for f in "$OUT/piton-runner.wasm" "$OUT/piton-viz.wasm"; do
        size=$(brotli -q 11 -c "$f" | wc -c)
        printf "  %10d  %s.br\n" "$size" "$(basename "$f")"
    done
fi
