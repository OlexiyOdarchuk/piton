/*
Piton is a high-performance interpreter for the Piton programming language.
It was born from a typo and evolved into a fast, Go-based execution engine.

Usage:

	piton ["-draw"/"-all", "-split", "-target=..."] <filename.piton>
*/
package main

import (
	"flag"
	"os"

	"github.com/OlexiyOdarchuk/piton/internal/repl"
	"github.com/OlexiyOdarchuk/piton/pkg/interpreter"
	"github.com/OlexiyOdarchuk/piton/pkg/visualizer"
)

func main() {
	if len(os.Args) == 1 {
		repl.Repl()
		return
	}
	filename := os.Args[len(os.Args)-1]
	visualize := flag.Bool("draw", false, "Zgeneruvaty blok-schemy do faily")
	visualizeProject := flag.Bool("all", false, "Zgeneruvaty blok-schemy do vsyogo proyektu")
	splitMode := flag.Bool("split", false, "Zgeneruvaty okremiy fail z blok-schemoy dlya kozhnoi functii")
	targetFunc := flag.String("target", "", "Zgeneruvaty blok-schemy do odnoy functii")
	format := flag.String("format", "svg", "Format vyvodu (svg, png, pdf, typst, excalidraw)")
	flag.Parse()

	content, err := os.ReadFile(filename)
	if err != nil {
		_, _ = os.Stdout.WriteString("Pomylka chitannya faily: " + err.Error() + "\n")
		return
	}

	if *visualize {
		diagrams, err := visualizer.Visualize(string(content), *targetFunc, *splitMode, *format)
		if err != nil {
			_, _ = os.Stderr.WriteString("Pomylka generacii shemu: " + err.Error() + "\n")
		}
		for chart_filename, data := range diagrams {
			if chart_filename == "flowchart.svg" || chart_filename == "flowchart." + *format {
				chart_filename = filename + ".svg"
			}
			_ = os.WriteFile(chart_filename, data, 0644)
			_, _ = os.Stdout.WriteString("Zberezheno: " + chart_filename + "\n")
		}

	} else if *visualizeProject {
		diagrams, err := visualizer.VisualizeProject(filename, *targetFunc, *splitMode, *format)
		if err != nil {
			_, _ = os.Stderr.WriteString("Pomylka generacii shemu: " + err.Error() + "\n")
		}
		for chart_filename, data := range diagrams {
			if chart_filename == "flowchart.svg" || chart_filename == "flowchart." + *format {
				chart_filename = filename + ".svg"
			}
			_ = os.WriteFile(chart_filename, data, 0644)
			_, _ = os.Stdout.WriteString("Zberezheno: " + chart_filename + "\n")
		}

	} else if err = interpreter.Run(string(content)); err != nil {
		_, _ = os.Stderr.WriteString("Pomylka vikonannya: " + err.Error() + "\n")
		return
	}
}
