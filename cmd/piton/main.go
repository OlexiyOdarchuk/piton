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
	"github.com/OlexiyOdarchuk/piton/interpreter"
)

func main() {
	if len(os.Args) == 1 {
		repl.Repl()
		return
	}
	filename := os.Args[len(os.Args)-1]
	visualize := flag.Bool("draw", false, "Generate flowchart to file")
	visualizeProject := flag.Bool("all", false, "Generate flowchart to all project")
	splitMode := flag.Bool("split", false, "Generate separate file for each function")
	targetFunc := flag.String("target", "", "Generate diagram only for specific function")
	flag.Parse()

	content, err := os.ReadFile(filename)
	if err != nil {
		os.Stdout.WriteString("Pomylka chitannya faily: " + err.Error() + "\n")
		return
	}

	if *visualize {
		diagrams, err := interpreter.Visualize(string(content), *targetFunc, *splitMode)
		if err != nil {
			os.Stderr.WriteString("Pomylka generacii shemu: " + err.Error() + "\n")
		}
		for chart_filename, data := range diagrams {
			if chart_filename == "flowchart.svg" {
				chart_filename = filename + ".svg"
			}
			os.WriteFile(chart_filename, data, 0644)
			os.Stdout.WriteString("Zberezheno: " + chart_filename + "\n")
		}

	} else if *visualizeProject {
		diagrams, err := interpreter.VisualizeProject(filename, *targetFunc, *splitMode)
		if err != nil {
			os.Stderr.WriteString("Pomylka generacii shemu: " + err.Error() + "\n")
		}
		for chart_filename, data := range diagrams {
			if chart_filename == "flowchart.svg" {
				chart_filename = filename + ".svg"
			}
			os.WriteFile(chart_filename, data, 0644)
			os.Stdout.WriteString("Zberezheno: " + chart_filename + "\n")
		}

	} else if err = interpreter.Run(string(content)); err != nil {
		os.Stderr.WriteString("Pomylka vikonannya: " + err.Error() + "\n")
		return
	}
}
