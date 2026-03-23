package visualizer

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/OlexiyOdarchuk/piton/internal/ast"

	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
)

type NodePtr struct {
	ID    string
	Label string
}

type Visualizer struct {
	sb      strings.Builder
	counter int
}

func (v *Visualizer) nextID() string {
	v.counter++
	return "n" + strconv.Itoa(v.counter)
}

func (v *Visualizer) writeNode(id, label, class string) {
	v.sb.WriteString(id)
	v.sb.WriteString(": \"")

	safeLabel := strings.ReplaceAll(label, "\"", "\\\"")
	safeLabel = strings.ReplaceAll(safeLabel, "\n", "\\n")
	safeLabel = strings.ReplaceAll(safeLabel, "\r", "")

	v.sb.WriteString(safeLabel)
	v.sb.WriteString("\" {class: ")
	v.sb.WriteString(class)
	v.sb.WriteString("}\n")
}

func (v *Visualizer) link(sources []NodePtr, targetID string) {
	for _, src := range sources {
		v.sb.WriteString(src.ID)
		v.sb.WriteString(" -> ")
		v.sb.WriteString(targetID)
		if src.Label != "" {
			v.sb.WriteString(": \"")
			v.sb.WriteString(src.Label)
			v.sb.WriteString("\"")
		}
		v.sb.WriteString("\n")
	}
}

func (v *Visualizer) walkBlock(stmts []ast.Stmt, prevNodes []NodePtr, endNode string) []NodePtr {
	leaves := prevNodes
	for _, stmt := range stmts {
		leaves = v.walkStmt(stmt, leaves, endNode)
	}
	return leaves
}

func (v *Visualizer) walkStmt(stmt ast.Stmt, prevNodes []NodePtr, endNode string) []NodePtr {
	switch s := stmt.(type) {

	case ast.FuncDefStmt:
		v.sb.WriteString("Функція ")
		v.sb.WriteString(s.Name)
		v.sb.WriteString(": {\n  style.fill: transparent\n")

		startID := v.nextID()
		v.writeNode(startID, "Початок "+s.Name, "terminal")

		funcEndID := v.nextID()
		v.writeNode(funcEndID, "Кінець "+s.Name, "terminal")

		leaves := v.walkBlock(s.Body, []NodePtr{{ID: startID}}, funcEndID)

		v.link(leaves, funcEndID)

		v.sb.WriteString("}\n")
		return prevNodes

	case ast.VarDecStmt:
		id := v.nextID()
		v.writeNode(id, s.Name+" = "+formatExpr(s.Expr), "process")
		v.link(prevNodes, id)
		return []NodePtr{{ID: id}}

	case ast.AssignStmt:
		id := v.nextID()
		v.writeNode(id, formatExpr(s.Target)+" = "+formatExpr(s.Expr), "process")
		v.link(prevNodes, id)
		return []NodePtr{{ID: id}}

	case ast.PrintStmt:
		id := v.nextID()
		v.writeNode(id, "Друк "+formatExpr(s.Expr), "io")
		v.link(prevNodes, id)
		return []NodePtr{{ID: id}}

	case ast.InputStmt:
		id := v.nextID()
		v.writeNode(id, "Ввід "+s.Name, "io")
		v.link(prevNodes, id)
		return []NodePtr{{ID: id}}

	case ast.ExprStmt:
		id := v.nextID()
		v.writeNode(id, formatExpr(s.Expr), "process")
		v.link(prevNodes, id)
		return []NodePtr{{ID: id}}

	case ast.ReturnStmt:
		id := v.nextID()
		v.writeNode(id, "Повернути "+formatExpr(s.Expr), "io")
		v.link(prevNodes, id)
		v.link([]NodePtr{{ID: id}}, endNode)
		return []NodePtr{}

	case ast.IfStmt:
		condID := v.nextID()
		v.writeNode(condID, formatExpr(s.Condition)+"?", "decision")
		v.link(prevNodes, condID)

		var allLeaves []NodePtr

		takLeaves := v.walkBlock(s.Body, []NodePtr{{ID: condID, Label: "Так"}}, endNode)
		allLeaves = append(allLeaves, takLeaves...)

		currentCondID := condID

		for _, elif := range s.ElseIfs {
			elifCondID := v.nextID()
			v.writeNode(elifCondID, formatExpr(elif.Condition)+"?", "decision")
			v.link([]NodePtr{{ID: currentCondID, Label: "Ні"}}, elifCondID)

			elifLeaves := v.walkBlock(elif.Body, []NodePtr{{ID: elifCondID, Label: "Так"}}, endNode)
			allLeaves = append(allLeaves, elifLeaves...)

			currentCondID = elifCondID
		}

		if len(s.ElseBody) > 0 {
			niLeaves := v.walkBlock(s.ElseBody, []NodePtr{{ID: currentCondID, Label: "Ні"}}, endNode)
			allLeaves = append(allLeaves, niLeaves...)
		} else {
			allLeaves = append(allLeaves, NodePtr{ID: currentCondID, Label: "Ні"})
		}

		return allLeaves

	case ast.PokyStmt:
		condID := v.nextID()
		v.writeNode(condID, formatExpr(s.Condition)+"?", "decision")
		v.link(prevNodes, condID)

		bodyLeaves := v.walkBlock(s.Body, []NodePtr{{ID: condID, Label: "Так"}}, endNode)

		v.link(bodyLeaves, condID)

		return []NodePtr{{ID: condID, Label: "Ні"}}
	}

	return prevNodes
}

func formatExpr(expr ast.Expr) string {
	if expr == nil {
		return ""
	}
	switch e := expr.(type) {
	case ast.Identifier:
		return e.Value
	case ast.NumberLiteral:
		return strconv.FormatFloat(e.Value, 'f', -1, 64)
	case ast.StringLiteral:
		return "\"" + e.Value + "\""
	case ast.InfixExpr:
		return formatExpr(e.Left) + " " + e.Operator + " " + formatExpr(e.Right)
	case ast.PrefixExpr:
		return e.Operator + formatExpr(e.Right)
	case ast.CallExpr:
		args := make([]string, len(e.Args))
		for i, arg := range e.Args {
			args[i] = formatExpr(arg)
		}
		if e.Receiver != nil {
			return formatExpr(e.Receiver) + "." + e.Name + "(" + strings.Join(args, ", ") + ")"
		}
		return e.Name + "(" + strings.Join(args, ", ") + ")"
	case ast.IndexExpr:
		return formatExpr(e.Left) + "[" + formatExpr(e.Index) + "]"
	case ast.SpysokLiteral:
		elements := make([]string, len(e.Elements))
		for i, el := range e.Elements {
			elements[i] = formatExpr(el)
		}
		return "[" + strings.Join(elements, ", ") + "]"

	case ast.SpysokExpr:
		start := ""
		if e.Start != nil {
			start = formatExpr(e.Start)
		}
		end := ""
		if e.End != nil {
			end = formatExpr(e.End)
		}
		return formatExpr(e.Left) + "[" + start + ":" + end + "]"
	default:
		return "[expr]"
	}
}

func renderSingleAST(program ast.Program) ([]byte, error) {
	v := &Visualizer{}

	v.sb.WriteString("direction: down\n")
	v.sb.WriteString("classes: {\n")
	v.sb.WriteString("  terminal: { shape: oval }\n")
	v.sb.WriteString("  process: { shape: rectangle }\n")
	v.sb.WriteString("  decision: { shape: diamond }\n")
	v.sb.WriteString("  io: { shape: parallelogram }\n")
	v.sb.WriteString("}\n\n")

	var globalStmts []ast.Stmt
	funcsByModule := make(map[string][]ast.Stmt)
	for _, stmt := range program.Statements {
		if f, isFunc := stmt.(ast.FuncDefStmt); isFunc {
			mod := f.Module
			if mod == "" {
				mod = "main"
			}
			funcsByModule[mod] = append(funcsByModule[mod], stmt)
		} else {
			globalStmts = append(globalStmts, stmt)
		}
	}

	if len(globalStmts) > 0 {
		v.sb.WriteString("Головна програма: {\n")
		v.sb.WriteString("  style: {\n")
		v.sb.WriteString("    fill: \"#e8f4f8\"\n")
		v.sb.WriteString("    stroke: \"#0284c7\"\n")
		v.sb.WriteString("    stroke-width: 2\n")
		v.sb.WriteString("  }\n")
		v.sb.WriteString("  label: \"Головна програма\"\n\n")

		startID := v.nextID()
		v.writeNode(startID, "Початок", "terminal")
		endID := v.nextID()
		v.writeNode(endID, "Кінець", "terminal")

		leaves := v.walkBlock(globalStmts, []NodePtr{{ID: startID}}, endID)
		v.link(leaves, endID)
		v.sb.WriteString("}\n")
	}

	for modName, funcs := range funcsByModule {
		if modName != "main" {
			v.sb.WriteString("Модуль ")
			v.sb.WriteString(modName)
			v.sb.WriteString(": {\n")
			v.sb.WriteString("  style.fill: \"#f8fafc\"\n")
			v.sb.WriteString("  style.stroke-dash: 5\n")
			v.sb.WriteString("  label: \"Модуль: ")
			v.sb.WriteString(modName)
			v.sb.WriteString("\"\n\n")
		}

		for _, f := range funcs {
			v.walkStmt(f, nil, "")
		}

		if modName != "main" {
			v.sb.WriteString("}\n")
		}
	}

	ctx := log.WithDefault(context.Background())
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2dagrelayout.DefaultLayout, nil
	}

	pad := int64(5)

	renderOpts := &d2svg.RenderOpts{
		Pad:     &pad,
		ThemeID: &d2themescatalog.GrapeSoda.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}

	diagram, _, err := d2lib.Compile(ctx, v.sb.String(), compileOpts, renderOpts)
	if err != nil {
		return nil, errors.New("Pomylka compilacii D2: " + err.Error() + "\nZgenerovanyi kod:\n" + v.sb.String())
	}

	return d2svg.Render(diagram, renderOpts)
}

func Visualize(program ast.Program, targetFunction string, splitFiles bool) (map[string][]byte, error) {
	results := make(map[string][]byte)

	if targetFunction != "" {
		for _, stmt := range program.Statements {
			if f, ok := stmt.(ast.FuncDefStmt); ok && f.Name == targetFunction {
				miniProg := ast.Program{Statements: []ast.Stmt{f}}
				svg, err := renderSingleAST(miniProg)
				if err != nil {
					return nil, err
				}
				results[f.Name+".svg"] = svg
				return results, nil
			}
		}
		return nil, errors.New("Functiyu " + targetFunction + " ne znaydeno")
	}
	if splitFiles {
		for _, stmt := range program.Statements {
			if f, ok := stmt.(ast.FuncDefStmt); ok {
				miniProg := ast.Program{Statements: []ast.Stmt{f}}
				svg, err := renderSingleAST(miniProg)
				if err != nil {
					return nil, err
				}
				fileName := f.Module + "_" + f.Name + ".svg"
				results[fileName] = svg
			}
		}
		return results, nil
	}

	svg, err := renderSingleAST(program)
	if err != nil {
		return nil, err
	}
	results["flowchart.svg"] = svg
	return results, nil
}
