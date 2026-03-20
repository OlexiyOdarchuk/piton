package repl

import (
	"bufio"
	"os"
	"strings"

	"github.com/OlexiyOdarchuk/piton/internal/evaluator"
	"github.com/OlexiyOdarchuk/piton/internal/lexer"
	"github.com/OlexiyOdarchuk/piton/internal/parser"
)

func Repl() {
	os.Stdout.WriteString("Vitay vas u Piton REPL. Mozhete pochyatu pusaty kod\n\n")

	eval := evaluator.New(os.Stdout)

	for {
		os.Stdout.WriteString("[95m" + ">>> " + "\x1b[0m")
		reader := bufio.NewReader(os.Stdin)
		inputStr, _ := reader.ReadString('\n')
		inputStr = strings.TrimSpace(inputStr)

		// TODO: Доробити, щоб при багаторядкових конструкціях можна було продовження писати

		tokens := lexer.Tokenize(inputStr)
		p := parser.New(tokens)
		program := p.ParseProgram()
		eval.Eval(program, eval.Globals)
	}
}
