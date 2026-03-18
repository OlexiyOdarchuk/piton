package interpreter

import (
	"io"
	"os"

	"github.com/OlexiyOdarchuk/piton/internal/evaluator"
	"github.com/OlexiyOdarchuk/piton/internal/lexer"
	"github.com/OlexiyOdarchuk/piton/internal/parser"
)

// Run parses and executes the provided Piton source code.
//
// The execution flow follows these steps:
//  1. Lexical analysis (tokenization).
//  2. Parsing tokens into an Abstract Syntax Tree (AST).
//  3. Evaluation of the AST nodes.
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
