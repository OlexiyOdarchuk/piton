package visualizer

import (
	"context"
	"strings"

	"github.com/OlexiyOdarchuk/piton/internal/ast"
	"oss.terrastruct.com/d2/d2graph"
	"oss.terrastruct.com/d2/d2layouts/d2dagrelayout"
	"oss.terrastruct.com/d2/d2lib"
	"oss.terrastruct.com/d2/d2renderers/d2svg"
	"oss.terrastruct.com/d2/d2themes/d2themescatalog"
	"oss.terrastruct.com/d2/lib/log"
	"oss.terrastruct.com/d2/lib/textmeasure"
	"oss.terrastruct.com/util-go/go2"
)

func Visualize(program ast.Program) ([]byte, error) {
	var sb strings.Builder

	ctx := log.WithDefault(context.Background())
	ruler, _ := textmeasure.NewRuler()
	layoutResolver := func(engine string) (d2graph.LayoutGraph, error) {
		return d2dagrelayout.DefaultLayout, nil
	}
	renderOpts := &d2svg.RenderOpts{
		Pad:     go2.Pointer(int64(5)),
		ThemeID: &d2themescatalog.GrapeSoda.ID,
	}
	compileOpts := &d2lib.CompileOptions{
		LayoutResolver: layoutResolver,
		Ruler:          ruler,
	}

	sb.WriteString("direction: down\n")
	//for _, stmt := range program.Statements {
	//	switch stmt = stmt.(type) {
	//TODO Доробити обробку всіх стейтментів і зробити генерацію якісних блок схем
	//}
	//}

	sb.WriteString("plankton -> formula: will steal\nformula: {\n  equation: |latex\n    \\lim_{h \\rightarrow 0 } \\frac{f(x+h)-f(x)}{h}\n  |\n}\n")

	diagram, _, _ := d2lib.Compile(ctx, sb.String(), compileOpts, renderOpts)
	return d2svg.Render(diagram, renderOpts)
}
