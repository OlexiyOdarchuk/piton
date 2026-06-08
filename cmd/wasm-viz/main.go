//go:build js && wasm

// wasm-viz exposes the Piton flowchart visualizer (rombik-backed) to the
// browser. Built with the standard Go toolchain and loaded lazily by the
// frontend only when the user requests a chart. rombik renders ДСТУ flowcharts
// in pure Go (no python3) via FromIR, so it is WASM-safe.
package main

import (
	"syscall/js"

	"github.com/OlexiyOdarchuk/piton/pkg/visualizer"
)

func visualizePiton(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return ""
	}
	code := args[0].String()
	diagrams, err := visualizer.Visualize(code, "", false, "svg")
	if err != nil {
		return "Pomylka generacii shemu: " + err.Error()
	}
	// Visualize returns a map keyed by suggested filenames; for non-split mode
	// the entry is always "flowchart.svg". Fall back to the first value if
	// upstream renames it.
	if svg, ok := diagrams["flowchart.svg"]; ok {
		return string(svg)
	}
	for _, svg := range diagrams {
		return string(svg)
	}
	return ""
}

func main() {
	js.Global().Set("visualizePiton", js.FuncOf(visualizePiton))
	js.Global().Set("pitonVizReady", js.ValueOf(true))
	select {}
}
