/*
Piton is a high-performance interpreter for the Piton programming language.
It was born from a typo and evolved into a fast, Go-based execution engine.

Usage:

	piton ["-draw" for flowchart] <filename.piton>
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
		os.Exit(0)
	}
	filename := os.Args[len(os.Args)-1]
	visualize := flag.Bool("draw", false, "Generate flowchart to file")
	visualizaAll := flag.Bool("drawProject", false, "Generate flowchart to all project")
	flag.Parse()

	content, err := os.ReadFile(filename)
	if err != nil {
		os.Stdout.WriteString("Pomylka chitannya faily: " + err.Error() + "\n")
		os.Exit(1)
	}

	if *visualize {
		diagram, err := interpreter.Visualize(string(content))
		if err != nil {
			os.Stderr.WriteString("Pomylka generacii shemu: " + err.Error() + "\n")
		}
		os.WriteFile(filename+".svg", diagram, 0600)
	} else if *visualizaAll {
		diagram, err := interpreter.VisualizeProject(filename)
		if err != nil {
			os.Stderr.WriteString("Pomylka generacii shemu: " + err.Error() + "\n")
		}
		os.WriteFile(filename+".svg", diagram, 0600)
	} else if err = interpreter.Run(string(content)); err != nil {
		os.Stderr.WriteString("Pomylka vikonannya: " + err.Error() + "\n")
		os.Exit(1)
	}
}
