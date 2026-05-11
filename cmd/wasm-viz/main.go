//go:build js && wasm

// wasm-viz exposes the Piton flowchart visualizer (D2-backed) to the browser.
// This is the heavy half of the wasm split — built with the standard Go
// toolchain because D2 / goja / chroma rely on features TinyGo doesn't
// support. Loaded lazily by the frontend only when the user requests a chart.
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
	diagrams, err := visualizer.Visualize(code, "", false)
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
