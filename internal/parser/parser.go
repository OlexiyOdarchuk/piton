package parser

import (
	"os"
	"strconv"

	"github.com/OlexiyOdarchuk/piton/internal/ast"
	"github.com/OlexiyOdarchuk/piton/internal/token"
)

func SyntaxError(line int) {
	os.Stdout.WriteString("Ryadok [" + strconv.Itoa(line) + "]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye!\n")
	os.Exit(1)
}

type Parser struct {
	tokens []token.Token
	pos    int
}

func New(tokens []token.Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) current() token.Token {
	if p.pos >= len(p.tokens) {
		return token.Token{Type: token.EOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) consume(t token.TokenType) bool {
	if p.current().Type == t {
		p.pos++
		return true
	}
	return false
}

func (p *Parser) expect(t token.TokenType) token.Token {
	if p.current().Type == t {
		tok := p.current()
		p.pos++
		return tok
	}
	SyntaxError(p.current().Line)
	return token.Token{}
}

const (
	_ int = iota
	LOWEST
	LESSGREATER
	AND
	OR
	SUM
	PRODUCT
	STUPIN_PREC
	PREFIX_PREC
	CALL
	INDEX_PREC
)

var precedences = map[token.TokenType]int{
	token.GT:       LESSGREATER,
	token.GT_EQ:    LESSGREATER,
	token.LT:       LESSGREATER,
	token.LT_EQ:    LESSGREATER,
	token.EQ:       LESSGREATER,
	token.NOT_EQ:   LESSGREATER,
	token.AND:      AND,
	token.OR:       OR,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.MULTIPTY: PRODUCT,
	token.DIVIDE:   PRODUCT,
	token.STUPIN:   STUPIN_PREC,
	token.LPAREN:   CALL,
	token.DOT:      INDEX_PREC,
	token.NOT:      INDEX_PREC,
	token.LBRACKET: INDEX_PREC,
}

func (p *Parser) peekPrecedence() int {
	if prec, ok := precedences[p.current().Type]; ok {
		return prec
	}
	return LOWEST
}

func (p *Parser) parseSpysokLiteral() ast.Expr {
	p.pos++
	var elements []ast.Expr
	if p.current().Type != token.RBRACKET {
		elements = append(elements, p.parseExpression(LOWEST))

		for p.current().Type == token.COMMA {
			p.pos++
			elements = append(elements, p.parseExpression(LOWEST))
		}
	}
	p.expect(token.RBRACKET)
	return ast.SpysokLiteral{Elements: elements}
}

func (p *Parser) parseExpression(precedence int) ast.Expr {
	var leftExp ast.Expr
	tok := p.current()

	switch tok.Type {
	case token.IDENT:
		p.pos++
		leftExp = ast.Identifier{Value: tok.Literal}
	case token.NUMBER:
		p.pos++
		val, _ := strconv.ParseFloat(tok.Literal, 64)
		leftExp = ast.NumberLiteral{Value: val}
	case token.STRING:
		p.pos++
		leftExp = ast.StringLiteral{Value: tok.Literal}
	case token.LBRACKET:
		leftExp = p.parseSpysokLiteral()
	case token.LPAREN:
		p.pos++
		leftExp = p.parseExpression(LOWEST)
		p.expect(token.RPAREN)
	case token.KORIN, token.LOH10, token.ABS, token.ARKSYN, token.KOSYNUS, token.NOT:
		p.pos++
		right := p.parseExpression(PREFIX_PREC)
		leftExp = ast.PrefixExpr{Operator: tok.Literal, Right: right}
	default:
		SyntaxError(tok.Line)
	}

	for p.current().Type != token.NEWLINE && p.current().Type != token.EOF && precedence < p.peekPrecedence() {
		opToken := p.current()

		if opToken.Type == token.DOT {
			p.pos++
			right := p.expect(token.IDENT).Literal
			leftExp = ast.SelectorExpr{Left: leftExp, Right: right}
			continue
		}

		if opToken.Type == token.LPAREN {
			p.pos++
			var args []ast.Expr
			if p.current().Type != token.RPAREN {
				args = append(args, p.parseExpression(LOWEST))
				for p.current().Type == token.COMMA {
					p.pos++
					args = append(args, p.parseExpression(LOWEST))
				}
			}
			p.expect(token.RPAREN)

			switch fn := leftExp.(type) {
			case ast.Identifier:
				leftExp = ast.CallExpr{Name: fn.Value, Args: args}
			case ast.SelectorExpr:
				leftExp = ast.CallExpr{Receiver: fn.Left, Name: fn.Right, Args: args}
			default:
				SyntaxError(opToken.Line)
			}
			continue
		}

		if opToken.Type == token.LBRACKET {
			p.pos++

			if p.current().Type == token.COLON {
				p.pos++
				var end ast.Expr
				if p.current().Type != token.RBRACKET {
					end = p.parseExpression(LOWEST)
				}
				p.expect(token.RBRACKET)
				leftExp = ast.SpysokExpr{Left: leftExp, End: end}
				continue
			}

			start := p.parseExpression(LOWEST)
			if p.current().Type == token.COLON {
				p.pos++
				var end ast.Expr
				if p.current().Type != token.RBRACKET {
					end = p.parseExpression(LOWEST)
				}
				p.expect(token.RBRACKET)
				leftExp = ast.SpysokExpr{Left: leftExp, Start: start, End: end}
				continue
			}

			p.expect(token.RBRACKET)
			leftExp = ast.IndexExpr{Left: leftExp, Index: start}
			continue
		}

		if opToken.Type == token.GT || opToken.Type == token.LT ||
			opToken.Type == token.PLUS || opToken.Type == token.MINUS ||
			opToken.Type == token.STUPIN || opToken.Type == token.MULTIPTY ||
			opToken.Type == token.DIVIDE || opToken.Type == token.EQ ||
			opToken.Type == token.NOT_EQ || opToken.Type == token.LT_EQ ||
			opToken.Type == token.GT_EQ || opToken.Type == token.AND ||
			opToken.Type == token.OR {
			p.pos++
			rightExp := p.parseExpression(precedences[opToken.Type])
			leftExp = ast.InfixExpr{Left: leftExp, Operator: opToken.Literal, Right: rightExp}
			continue
		}

		break
	}
	return leftExp
}

func (p *Parser) consumeNewlineOrEOF() {
	if p.current().Type == token.NEWLINE {
		p.pos++
	} else if p.current().Type == token.EOF {
	} else {
		SyntaxError(p.current().Line)
	}
}

func (p *Parser) parseBlock() []ast.Stmt {
	p.expect(token.INDENT)
	var stmts []ast.Stmt
	for p.current().Type != token.DEDENT && p.current().Type != token.EOF {
		if p.current().Type == token.NEWLINE {
			p.pos++
			continue
		}
		stmts = append(stmts, p.parseStatement())
	}
	p.expect(token.DEDENT)
	return stmts
}

func (p *Parser) parseStatement() ast.Stmt {
	switch p.current().Type {
	case token.FUNCTIA:
		p.pos++
		name := p.expect(token.IDENT).Literal
		p.expect(token.LPAREN)
		var params []string
		if p.current().Type != token.RPAREN {
			params = append(params, p.expect(token.IDENT).Literal)
			for p.current().Type == token.COMMA {
				p.pos++
				params = append(params, p.expect(token.IDENT).Literal)
			}
		}
		p.expect(token.RPAREN)
		p.expect(token.COLON)
		p.consumeNewlineOrEOF()
		body := p.parseBlock()
		return ast.FuncDefStmt{Name: name, Params: params, Body: body}
	case token.DRUKUVATY:
		p.pos++
		exp := p.parseExpression(LOWEST)
		p.consumeNewlineOrEOF()
		return ast.PrintStmt{Expr: exp}
	case token.NEKHAY:
		p.pos++
		name := p.expect(token.IDENT).Literal
		p.expect(token.ASSIGN)
		exp := p.parseExpression(LOWEST)
		p.consumeNewlineOrEOF()
		return ast.VarDecStmt{Name: name, Expr: exp}
	case token.VVID:
		p.pos++
		name := p.expect(token.IDENT).Literal
		p.consumeNewlineOrEOF()
		return ast.InputStmt{Name: name}
	case token.YAKSHO:
		p.pos++
		cond := p.parseExpression(LOWEST)
		p.expect(token.COLON)
		p.consumeNewlineOrEOF()
		body := p.parseBlock()

		var elifs []ast.ElseIf
		var elseBody []ast.Stmt

		for p.current().Type == token.INACKSHE {
			p.pos++
			if p.current().Type == token.YAKSHO {
				p.pos++
				elifCond := p.parseExpression(LOWEST)
				p.expect(token.COLON)
				p.consumeNewlineOrEOF()
				elifBody := p.parseBlock()
				elifs = append(elifs, ast.ElseIf{Condition: elifCond, Body: elifBody})
			} else {
				p.expect(token.COLON)
				p.consumeNewlineOrEOF()
				elseBody = p.parseBlock()
				break
			}
		}
		return ast.IfStmt{Condition: cond, Body: body, ElseIfs: elifs, ElseBody: elseBody}
	case token.VERNUTY:
		p.pos++
		exp := p.parseExpression(LOWEST)
		p.consumeNewlineOrEOF()
		return ast.ReturnStmt{Expr: exp}
	case token.POKY:
		p.pos++
		condition := p.parseExpression(LOWEST)
		p.expect(token.COLON)
		p.consumeNewlineOrEOF()
		body := p.parseBlock()
		return ast.PokyStmt{Condition: condition, Body: body}
	case token.IDENT:
		saved := p.pos
		p.pos++
		if p.current().Type == token.ASSIGN {
			name := p.tokens[saved].Literal
			p.pos++
			exp := p.parseExpression(LOWEST)
			p.consumeNewlineOrEOF()
			return ast.AssignStmt{Name: name, Expr: exp}
		}
		p.pos = saved
		exp := p.parseExpression(LOWEST)
		p.consumeNewlineOrEOF()
		return ast.ExprStmt{Expr: exp}
	case token.VYKORYSTATY:
		p.pos++
		var filename string
		switch p.current().Type {
		case token.STRING:
			filename = p.expect(token.STRING).Literal
		case token.IDENT:
			filename = p.expect(token.IDENT).Literal
		default:
			SyntaxError(p.current().Line)
		}
		p.consumeNewlineOrEOF()
		return ast.ImportStmt{Filename: ast.StringLiteral{Value: filename}}
	default:
		SyntaxError(p.current().Line)
		return nil
	}
}

func (p *Parser) ParseProgram() ast.Program {
	var stmts []ast.Stmt
	for p.current().Type != token.EOF {
		if p.current().Type == token.NEWLINE {
			p.pos++
			continue
		}
		stmts = append(stmts, p.parseStatement())
	}
	return ast.Program{Statements: stmts}
}
