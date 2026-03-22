package interpreter

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/OlexiyOdarchuk/piton/internal/ast"
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

// VisualizeProject parses the provided entry Piton file and automatically resolves
// all its dependencies ('vykorystaty' statements) to generate a comprehensive
// flowchart of the entire project's logic.
//
// Unlike Visualize, which only processes a single isolated file, this function
// delegates the creation of a unified Abstract Syntax Tree (AST) to the interpreter.
//
// The process follows a strict pipeline:
//  1. Interpreter: Reads the entry file, recursively resolves and parses all
//     imported modules, and merges their functions into a single "Super AST".
//  2. Visualizer: Traverses this complete AST to produce a unified D2-based
//     diagram following flowchart standards.
//
// Returns the rendered diagram as a byte slice SVG or an error if
// file reading, parsing, or import resolution fails.
func VisualizeProject(entryFilePath string) ([]byte, error) {
	superProgram, err := parseProject(entryFilePath)
	if err != nil {
		return nil, err
	}
	return visualizer.Visualize(superProgram)
}

func parseProject(entryFilePath string) (ast.Program, error) {
	visited := make(map[string]bool)
	superProgram := ast.Program{Statements: []ast.Stmt{}}

	err := resolveImports(entryFilePath, &superProgram, visited, true)
	if err != nil {
		return ast.Program{}, err
	}

	return superProgram, nil
}

func resolveImports(filePath string, superProgram *ast.Program, visited map[string]bool, isMain bool) error {
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return errors.New("Pomylka shlyahu " + filePath + ": " + err.Error())
	}

	if visited[absPath] {
		return nil
	}
	visited[absPath] = true

	moduleName := strings.TrimSuffix(filepath.Base(absPath), filepath.Ext(absPath))
	if isMain {
		moduleName = "main"
	}

	codeBytes, err := os.ReadFile(absPath)
	if err != nil {
		return errors.New("Ne vdalosya prochutaty fail " + filePath + ": " + absPath + err.Error())
	}

	tokens := lexer.Tokenize(string(codeBytes))
	p := parser.New(tokens)
	program := p.ParseProgram()

	baseDir := filepath.Dir(absPath)

	for _, stmt := range program.Statements {
		switch s := stmt.(type) {
		case ast.ImportStmt:
			if strLit, ok := s.Filename.(ast.StringLiteral); ok {
				importedFilePath := filepath.Join(baseDir, strLit.Value+".piton")
				err := resolveImports(importedFilePath, superProgram, visited, false)
				if err != nil {
					return err
				}
			} else {
				return errors.New("Shlyah maye buty ryadkom " + absPath)
			}

		case ast.FuncDefStmt:
			s.Module = moduleName
			superProgram.Statements = append(superProgram.Statements, s)

		default:
			if isMain {
				superProgram.Statements = append(superProgram.Statements, s)
			}
		}
	}

	return nil
}
