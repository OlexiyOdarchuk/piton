package lexer

import (
	"strings"
	"unicode"

	"github.com/OlexiyOdarchuk/piton/internal/token"
)

func Tokenize(input string) []token.Token {
	var tokens []token.Token
	lines := strings.Split(input, "\n")
	indents := []int{0}

	for lineIdx, line := range lines {
		lineNum := lineIdx + 1
		indent := 0
	OuterLoop:
		for _, ch := range line {
			switch ch {
			case ' ':
				indent++
			case '\t':
				indent += 4
			default:
				break OuterLoop
			}
		}

		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if indent > indents[len(indents)-1] {
			indents = append(indents, indent)
			tokens = append(tokens, token.Token{Type: token.INDENT, Literal: "", Line: lineNum})
		} else if indent < indents[len(indents)-1] {
			for len(indents) > 1 && indent < indents[len(indents)-1] {
				indents = indents[:len(indents)-1]
				tokens = append(tokens, token.Token{Type: token.DEDENT, Literal: "", Line: lineNum})
			}
		}

		lineTokens := TokenizeLine(trimmed)
		for i := range lineTokens {
			lineTokens[i].Line = lineNum
		}
		tokens = append(tokens, lineTokens...)
		tokens = append(tokens, token.Token{Type: token.NEWLINE, Literal: "\n", Line: lineNum})
	}

	for len(indents) > 1 {
		indents = indents[:len(indents)-1]
		tokens = append(tokens, token.Token{Type: token.DEDENT, Literal: "", Line: len(lines)})
	}

	tokens = append(tokens, token.Token{Type: token.EOF, Literal: "", Line: len(lines)})
	return tokens
}

func TokenizeLine(line string) []token.Token {
	var tokens []token.Token
	for i := 0; i < len(line); {
		ch := rune(line[i])
		switch {
		case ch == ' ' || ch == '\t' || ch == '\r':
			i++
		case ch == '=':
			tokens = append(tokens, token.Token{Type: token.ASSIGN, Literal: "="})
			i++
		case ch == '+':
			tokens = append(tokens, token.Token{Type: token.PLUS, Literal: "+"})
			i++
		case ch == '>':
			tokens = append(tokens, token.Token{Type: token.GT, Literal: ">"})
			i++
		case ch == '<':
			tokens = append(tokens, token.Token{Type: token.LT, Literal: "<"})
			i++
		case ch == '(':
			tokens = append(tokens, token.Token{Type: token.LPAREN, Literal: "("})
			i++
		case ch == ')':
			tokens = append(tokens, token.Token{Type: token.RPAREN, Literal: ")"})
			i++
		case ch == ':':
			tokens = append(tokens, token.Token{Type: token.COLON, Literal: ":"})
			i++
		case ch == '"':
			i++
			var valBytes []byte
			for i < len(line) && line[i] != '"' {
				if line[i] == '\\' && i+1 < len(line) {
					if line[i+1] == 'n' {
						valBytes = append(valBytes, '\n')
						i += 2
						continue
					}
				}
				valBytes = append(valBytes, line[i])
				i++
			}
			if i < len(line) {
				i++
			}
			tokens = append(tokens, token.Token{Type: token.STRING, Literal: string(valBytes)})
		case ch == '[':
			tokens = append(tokens, token.Token{Type: token.LBRACKET, Literal: "["})
			i++
		case ch == ']':
			tokens = append(tokens, token.Token{Type: token.RBRACKET, Literal: "]"})
			i++
		case ch == ',':
			tokens = append(tokens, token.Token{Type: token.COMMA, Literal: ","})
			i++
		case unicode.IsDigit(ch):
			start := i
			for i < len(line) && (unicode.IsDigit(rune(line[i])) || line[i] == '.') {
				i++
			}
			tokens = append(tokens, token.Token{Type: token.NUMBER, Literal: line[start:i]})
		case unicode.IsLetter(ch):
			start := i
			for i < len(line) && (unicode.IsLetter(rune(line[i])) || unicode.IsDigit(rune(line[i])) || line[i] == '_') {
				i++
			}
			ident := line[start:i]
			tokens = append(tokens, token.Token{Type: token.LookupIdent(ident), Literal: ident})
		default:
			tokens = append(tokens, token.Token{Type: token.ILLEGAL, Literal: string(ch)})
			i++
		}
	}
	return tokens
}
