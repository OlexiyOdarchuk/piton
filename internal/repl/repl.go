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
	_, _ = os.Stdout.WriteString("Vitay vas u Piton REPL. Mozhete pochyatu pusaty kod\n")
	_, _ = os.Stdout.WriteString("Shchob vuyty z REPL, napyshit 'exit'\n\n")
	eval := evaluator.New(os.Stdout)
	reader := bufio.NewReader(os.Stdin)

	for {
		_, _ = os.Stdout.WriteString("\x1b[95m>>> \x1b[0m")
		inputStr, _ := reader.ReadString('\n')
		inputStr = strings.TrimRight(inputStr, "\r\n")

		if inputStr == "exit" {
			_, _ = os.Stdout.WriteString("Harnoho dnya!\n")
			break
		}
		if inputStr == "" {
			continue
		}

		if strings.HasSuffix(inputStr, ":") {
			fullCode := inputStr + "\n"
			level := 1

			for level > 0 {
				indent := strings.Repeat("    ", level)
				_, _ = os.Stdout.WriteString("\x1b[95m... " + indent + "\x1b[0m")

				line, _ := reader.ReadString('\n')
				trimmed := strings.TrimSpace(line)

				if trimmed == "" {
					level--
					continue
				}

				if strings.HasSuffix(trimmed, ":") {
					level++
				}

				fullCode += indent + line
			}
			inputStr = fullCode
		}
		tokens := lexer.Tokenize(inputStr)
		p := parser.New(tokens, os.Stdout)
		program := p.ParseProgram()

		eval.Eval(program, eval.Globals)
	}
}
