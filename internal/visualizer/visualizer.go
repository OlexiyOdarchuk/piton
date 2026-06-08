// Пакет visualizer перетворює AST програми Piton на блок-схеми за ДСТУ.
//
// Рендер делеговано бібліотеці rombik (github.com/OlexiyOdarchuk/rombik):
// ми лише зводимо AST Piton до її проміжного представлення (ir), а вся
// розкладка, маршрутизація ребер і малювання SVG — на боці rombik.
//
// rombik будує ОДНУ схему на кожну функцію (ir.Func). Тож:
//   - targetFunction — одна схема для названої функції;
//   - splitFiles     — окремий файл на кожну функцію;
//   - типовий режим  — головна програма + усі функції, зведені вертикально
//     в один flowchart.svg (через вкладені <svg>).
package visualizer

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/OlexiyOdarchuk/piton/internal/ast"

		"github.com/OlexiyOdarchuk/rombik/pkg/diagram"
	"github.com/OlexiyOdarchuk/rombik/pkg/render/excalidraw"
	"github.com/OlexiyOdarchuk/rombik/pkg/render/raster"
	"github.com/OlexiyOdarchuk/rombik/pkg/render/typst"
	"github.com/OlexiyOdarchuk/rombik/pkg/ir"
	"github.com/OlexiyOdarchuk/rombik/pkg/rombik"
)

// irBuilder зводить вузли AST Piton до вузлів ir rombik. Множина userFuncs
// потрібна, щоб відрізнити виклик визначеної у файлі функції (ir.Call —
// «предетермінований процес» ДСТУ) від звичайної дії.
type irBuilder struct {
	userFuncs map[string]bool
}

// block перетворює послідовність інструкцій у *ir.Block.
func (b *irBuilder) block(stmts []ast.Stmt) *ir.Block {
	blk := &ir.Block{}
	for _, stmt := range stmts {
		if node := b.stmt(stmt); node != nil {
			blk.Stmts = append(blk.Stmts, node)
		}
	}
	return blk
}

// stmt відображає одну інструкцію AST на вузол ir (або nil, якщо вузол не
// малюється — наприклад, import).
func (b *irBuilder) stmt(stmt ast.Stmt) ir.Node {
	switch s := stmt.(type) {

	case ast.VarDecStmt:
		return &ir.Process{Text: s.Name + " = " + formatExpr(s.Expr)}

	case ast.AssignStmt:
		return &ir.Process{Text: formatExpr(s.Target) + " = " + formatExpr(s.Expr)}

	case ast.PrintStmt:
		return &ir.IO{Text: "Вивід " + formatExpr(s.Expr)}

	case ast.InputStmt:
		return &ir.IO{Text: "Ввід " + s.Name}

	case ast.ReturnStmt:
		text := "Повернути"
		if s.Expr != nil {
			text += " " + formatExpr(s.Expr)
		}
		return &ir.Terminal{Text: text}

	case ast.ExprStmt:
		if call, ok := s.Expr.(ast.CallExpr); ok && call.Receiver == nil && b.userFuncs[call.Name] {
			return &ir.Call{Text: formatExpr(s.Expr)}
		}
		return &ir.Process{Text: formatExpr(s.Expr)}

	case ast.IfStmt:
		return &ir.If{
			Cond: formatExpr(s.Condition),
			Then: b.block(s.Body),
			Else: b.elifChain(s.ElseIfs, s.ElseBody),
		}

	case ast.PokyStmt:
		return &ir.While{Cond: formatExpr(s.Condition), Body: b.block(s.Body)}
	}

	return nil
}

// elifChain розгортає ланцюжок elif…else у вкладені ir.If всередині гілки Else
// (rombik має лише Then/Else). «Немає else» подаємо порожнім блоком, а не nil:
// порожня гілка не резервує ширини (як і nil), зате не залежимо від того, чи
// nil-guard'ить гілку If конкретна версія rombik.
func (b *irBuilder) elifChain(elifs []ast.ElseIf, elseBody []ast.Stmt) *ir.Block {
	if len(elifs) == 0 {
		if len(elseBody) == 0 {
			return &ir.Block{}
		}
		return b.block(elseBody)
	}

	head := elifs[0]
	nested := &ir.If{
		Cond: formatExpr(head.Condition),
		Then: b.block(head.Body),
		Else: b.elifChain(elifs[1:], elseBody),
	}
	return &ir.Block{Stmts: []ir.Node{nested}}
}

// funcIR будує ir.Func для визначення функції.
func (b *irBuilder) funcIR(f ast.FuncDefStmt) ir.Func {
	return ir.Func{Name: f.Name, Body: b.block(f.Body)}
}

// renderFunc розкладає одну ir.Func у SVG. suffix додається до тексту
// термінаторів (« average» → «Початок average»/«Кінець average»); для
// головної програми suffix порожній.
func renderFunc(fn ir.Func, suffix string) (rombik.Result, error) {
	opts := rombik.Options{
		SingleEnd: true, // один спільний «Кінець» на функцію
		Yes:       "Так",
		No:        "Ні",
		StartText: "Початок" + suffix,
		EndText:   "Кінець" + suffix,
	}

	results := rombik.FromIR([]ir.Func{fn}, opts)
	if len(results) == 0 {
		return rombik.Result{}, errors.New("rombik ne povernuv shemu dlya " + fn.Name)
	}
	return results[0], nil
}

func renderResult(res rombik.Result, format string) ([]byte, error) {
	switch format {
	case "svg":
		return []byte(res.SVG()), nil
	case "typst":
		return []byte(res.Typst()), nil
	case "excalidraw":
		return []byte(res.Excalidraw()), nil
	case "png":
		return res.PNG(2.0)
	case "pdf":
		return res.PDF()
	default:
		return []byte(res.SVG()), nil
	}
}

func renderAll(ds []*diagram.Diagram, format string) ([]byte, error) {
	switch format {
	case "svg":
		// SVG stack requires titles, handled separately by stackSVGs
		return nil, errors.New("unreachable")
	case "typst":
		return []byte(typst.RenderAll(ds)), nil
	case "excalidraw":
		return []byte(excalidraw.RenderAll(ds)), nil
	case "png":
		return raster.PNGAll(ds, 2.0)
	case "pdf":
		return raster.PDFAll(ds)
	default:
		return nil, errors.New("unreachable")
	}
}

// collectFuncs збирає імена всіх функцій, визначених у програмі.
func collectFuncs(program ast.Program) map[string]bool {
	funcs := make(map[string]bool)
	for _, stmt := range program.Statements {
		if f, ok := stmt.(ast.FuncDefStmt); ok {
			funcs[f.Name] = true
		}
	}
	return funcs
}

// Visualize будує блок-схеми для програми.
//
//   - targetFunction != "" — лише для названої функції → {name}.svg;
//   - splitFiles            — окрема схема на кожну функцію → {module}_{name}.svg;
//   - інакше                — головна програма + усі функції в одному flowchart.svg.
func Visualize(program ast.Program, targetFunction string, splitFiles bool, format string) (map[string][]byte, error) {
	results := make(map[string][]byte)
	b := &irBuilder{userFuncs: collectFuncs(program)}

	if targetFunction != "" {
		for _, stmt := range program.Statements {
			if f, ok := stmt.(ast.FuncDefStmt); ok && f.Name == targetFunction {
				res, err := renderFunc(b.funcIR(f), " "+f.Name)
				if err != nil {
					return nil, err
				}
				data, err := renderResult(res, format)
				if err != nil {
					return nil, err
				}
				ext := format
				if ext == "" {
					ext = "svg"
				}
				results[f.Name+"."+ext] = data
				return results, nil
			}
		}
		return nil, errors.New("Functiyu " + targetFunction + " ne znaydeno")
	}

	if splitFiles {
		for _, stmt := range program.Statements {
			f, ok := stmt.(ast.FuncDefStmt)
			if !ok {
				continue
			}
			res, err := renderFunc(b.funcIR(f), " "+f.Name)
			if err != nil {
				return nil, err
			}
			data, err := renderResult(res, format)
			if err != nil {
				return nil, err
			}
			ext := format
			if ext == "" {
				ext = "svg"
			}
			module := f.Module
			if module == "" {
				module = "main"
			}
			results[module+"_"+f.Name+"."+ext] = data
		}
		return results, nil
	}

	// Типовий режим: головна програма (глобальні інструкції) + кожна функція,
	// зведені вертикально в один SVG.
	var parts []titledSVG
	var diagrams []*diagram.Diagram

	var globals []ast.Stmt
	for _, stmt := range program.Statements {
		switch stmt.(type) {
		case ast.FuncDefStmt, ast.ImportStmt:
			// функції рендеримо окремо нижче; import не малюється
		default:
			globals = append(globals, stmt)
		}
	}

	if len(globals) > 0 {
		fn := ir.Func{Name: "main", Body: b.block(globals)}
		res, err := renderFunc(fn, "")
		if err != nil {
			return nil, err
		}
		parts = append(parts, titledSVG{title: "Головна програма", svg: res.SVG()})
		diagrams = append(diagrams, res.Diagram)
	}

	for _, stmt := range program.Statements {
		f, ok := stmt.(ast.FuncDefStmt)
		if !ok {
			continue
		}
		res, err := renderFunc(b.funcIR(f), " "+f.Name)
		if err != nil {
			return nil, err
		}
		title := "Функція " + f.Name
		if f.Module != "" && f.Module != "main" {
			title = "Модуль " + f.Module + " — функція " + f.Name
		}
		parts = append(parts, titledSVG{title: title, svg: res.SVG()})
		diagrams = append(diagrams, res.Diagram)
	}

	if len(parts) == 0 {
		return nil, errors.New("Nemaye chogo vizualizuvaty")
	}

	ext := format
	if ext == "" {
		ext = "svg"
	}
	if format == "svg" || format == "" {
		results["flowchart.svg"] = []byte(stackSVGs(parts))
	} else {
		data, err := renderAll(diagrams, format)
		if err != nil {
			return nil, err
		}
		results["flowchart."+ext] = data
	}
	return results, nil
}

// titledSVG — одна готова схема з підписом для зведеного вигляду.
type titledSVG struct {
	title string
	svg   string
}

// svgDimRe витягує width/height із кореневого тегу SVG, який генерує rombik.
var svgDimRe = regexp.MustCompile(`<svg[^>]*\bwidth="([0-9.]+)"[^>]*\bheight="([0-9.]+)"`)

// svgDims повертає габарити SVG; за невдачі — розумний запас.
func svgDims(svg string) (float64, float64) {
	m := svgDimRe.FindStringSubmatch(svg)
	if m == nil {
		return 600, 400
	}
	w, _ := strconv.ParseFloat(m[1], 64)
	h, _ := strconv.ParseFloat(m[2], 64)
	return w, h
}

// stackSVGs зводить кілька схем у один SVG, вкладаючи кожну як <svg> з власним
// viewBox і підписом над нею. Окремі схеми центруються по ширині.
func stackSVGs(parts []titledSVG) string {
	const (
		pad    = 20.0 // зовнішні поля
		gap    = 44.0 // відступ між схемами
		titleH = 30.0 // висота смуги під підпис
	)

	type item struct {
		t    titledSVG
		w, h float64
	}

	var maxW float64
	items := make([]item, 0, len(parts))
	for _, p := range parts {
		w, h := svgDims(p.svg)
		items = append(items, item{t: p, w: w, h: h})
		if w > maxW {
			maxW = w
		}
	}

	totalH := pad
	for i, it := range items {
		if i > 0 {
			totalH += gap
		}
		totalH += titleH + it.h
	}
	totalH += pad
	totalW := maxW + 2*pad

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f" font-family="Arial, sans-serif">`,
		totalW, totalH, totalW, totalH)
	b.WriteString(`<rect width="100%" height="100%" fill="#ffffff"/>`)

	y := pad
	for _, it := range items {
		fmt.Fprintf(&b, `<text x="%.1f" y="%.1f" font-size="16" font-weight="bold" fill="#0f172a">%s</text>`,
			pad, y+20, esc(it.t.title))
		y += titleH

		x := pad + (maxW-it.w)/2
		b.WriteString(injectXY(it.t.svg, x, y))
		y += it.h + gap
	}

	b.WriteString(`</svg>`)
	return b.String()
}

// injectXY додає атрибути x/y у кореневий тег вкладеного SVG, позиціонуючи
// його в батьківському полотні (внутрішній viewBox залишається власним).
func injectXY(svg string, x, y float64) string {
	idx := strings.Index(svg, "<svg")
	if idx < 0 {
		return svg
	}
	at := idx + len("<svg")
	return svg[:at] + fmt.Sprintf(` x="%.1f" y="%.1f"`, x, y) + svg[at:]
}

// esc екранує текст підпису для вставлення в SVG.
func esc(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

// formatExpr відтворює вираз Piton як рядок для підпису фігури.
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
	case ast.SlovnykLiteral:
		pairs := make([]string, len(e.Pairs))
		for i, pair := range e.Pairs {
			pairs[i] = formatExpr(pair.Key) + ": " + formatExpr(pair.Value)
		}
		return "{" + strings.Join(pairs, ", ") + "}"
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
	case ast.SelectorExpr:
		return formatExpr(e.Left) + "." + e.Right
	default:
		return "[expr]"
	}
}
