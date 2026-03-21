package interpreter

import (
	"io"
	"os"

	"github.com/OlexiyOdarchuk/piton/internal/evaluator"
	"github.com/OlexiyOdarchuk/piton/internal/lexer"
	"github.com/OlexiyOdarchuk/piton/internal/parser"
	"github.com/OlexiyOdarchuk/piton/internal/visualizer"
)

// Run parses and executes the provided Piton source code.
//
// The execution flow follows these steps:
//  1. Lexer: Breaks the raw string into a stream of tokens.
//  2. Parser: Constructs an Abstract Syntax Tree (AST) from the tokens.
//  3. Evaluator: Evaluation of the AST nodes.
//
// By default, it writes the program output to os.Stdout. You can provide
// a custom io.Writer (e.g., bytes.Buffer or a file) as an optional argument.
//
// It returns an error if any stage of the process (lexing, parsing,
// evaluation, or flushing the output) fails.
func Run(code string, output ...io.Writer) error {
	var out io.Writer = os.Stdout

	if len(output) > 0 && output[0] != nil {
		out = output[0]
	}

	tokens := lexer.Tokenize(code)
	p := parser.New(tokens)
	program := p.ParseProgram()

	eval := evaluator.New(out)
	eval.Eval(program, eval.Globals)
	return eval.Flush()
}

// Visualize parses the provided Piton source code and generates a
// graphical representation of the program's logic.
//
// The process follows a strict pipeline:
//  1. Lexer: Breaks the raw string into a stream of tokens.
//  2. Parser: Constructs an Abstract Syntax Tree (AST) from the tokens.
//  3. Visualizer: Traverses the AST to produce a D2-based diagram
//     following flowchart standards.
//
// Returns the rendered diagram as a byte slice SVG or an error
// if lexical or structural analysis fails.
func Visualize(code string) ([]byte, error) {
	tokens := lexer.Tokenize(code)
	p := parser.New(tokens)
	program := p.ParseProgram()
	return visualizer.Visualize(program)
}
