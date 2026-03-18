/*
Piton is a high-performance interpreter for the Piton programming language.
It was born from a typo and evolved into a fast, Go-based execution engine.

Usage:

	piton <filename.piton>
*/
package main

import (
	"os"

	"github.com/OlexiyOdarchuk/piton/interpreter"
)

func main() {
	if len(os.Args) < 2 {
		os.Stdout.WriteString("Vikorystannya: piton <file.piton>\n")
		os.Exit(1)
	}

	content, err := os.ReadFile(os.Args[1])
	if err != nil {
		os.Stdout.WriteString("Pomylka chitannya faily: " + err.Error() + "\n")
		os.Exit(1)
	}

	if err := interpreter.Run(string(content)); err != nil {
		os.Stderr.WriteString("Pomylka vikonannya: " + err.Error() + "\n")
		os.Exit(1)
	}
}
