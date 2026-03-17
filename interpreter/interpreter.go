package interpreter

import (
	"io"
	"os"

	"github.com/OlexiyOdarchuk/piton/internal/evaluator"
	"github.com/OlexiyOdarchuk/piton/internal/lexer"
	"github.com/OlexiyOdarchuk/piton/internal/parser"
)

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
