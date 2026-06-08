//go:build js && wasm

// wasm-runner exposes Piton's interpreter to the browser.
// Built with TinyGo for minimal binary size — the rombik-based visualizer is
// intentionally excluded to keep the runner small. See cmd/wasm-viz for the
// diagram side.
package main

import (
	"bytes"
	"syscall/js"

	"github.com/OlexiyOdarchuk/piton/pkg/interpreter"
)

func runPiton(this js.Value, args []js.Value) any {
	if len(args) == 0 {
		return ""
	}
	code := args[0].String()
	var buf bytes.Buffer
	if err := interpreter.Run(code, &buf); err != nil {
		buf.WriteString("\nPomylka vikonannya: " + err.Error() + "\n")
	}
	return buf.String()
}

func main() {
	js.Global().Set("runPiton", js.FuncOf(runPiton))
	// Mark runner ready so the JS side can detect a successful boot regardless
	// of whether visualizePiton (the bigger wasm) has loaded yet.
	js.Global().Set("pitonRunnerReady", js.ValueOf(true))
	// Park forever — Go programs exit when main() returns, which would tear
	// down the JS callbacks. A blocking channel keeps the runtime alive.
	select {}
}
