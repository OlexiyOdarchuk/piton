package evaluator

import (
	"bufio"
	"io"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/OlexiyOdarchuk/piton/internal/ast"
)

type Evaluator struct {
	Globals *Environment
	Out     *bufio.Writer
}

func New(out io.Writer) *Evaluator {
	return &Evaluator{
		Globals: NewEnv(nil),
		Out:     bufio.NewWriter(out),
	}
}

func (ev *Evaluator) Flush() error {
	return ev.Out.Flush()
}

type ReturnValue struct{ Value interface{} }

func (ev *Evaluator) Eval(node ast.Node, env *Environment) interface{} {
	switch n := node.(type) {
	case ast.Program:
		var res interface{}
		for _, stmt := range n.Statements {
			if _, isFunc := stmt.(ast.FuncDefStmt); isFunc {
				ev.Eval(stmt, env)
			}
		}
		for _, stmt := range n.Statements {
			if _, isFunc := stmt.(ast.FuncDefStmt); !isFunc {
				res = ev.Eval(stmt, env)
			}
		}
		return res
	case ast.FuncDefStmt:
		env.Set(n.Name, n)
		return nil
	case ast.PrintStmt:
		val := ev.Eval(n.Expr, env)
		if list, ok := val.([]interface{}); ok {
			ev.Out.WriteString("[")
			for i, value := range list {
				switch v := value.(type) {
				case float64:
					ev.Out.WriteString(strconv.FormatFloat(v, 'g', -1, 64))
				case string:
					ev.Out.WriteString("\"" + v + "\"")
				case bool:
					if v {
						ev.Out.WriteString("true")
					} else {
						ev.Out.WriteString("false")
					}
				default:
					ev.Out.WriteString(value.(string))
				}
				if len(list)-1 > i {
					ev.Out.WriteString(", ")
				}
			}
			ev.Out.WriteString("]\n")
			return nil
		}
		if str, ok := val.(string); ok {
			ev.Out.WriteString(str + "\n")
		} else if num, ok := val.(float64); ok {
			ev.Out.WriteString(strconv.FormatFloat(num, 'f', -1, 64) + "\n")
		} else if b, ok := val.(bool); ok {
			if b {
				ev.Out.WriteString("true\n")
			} else {
				ev.Out.WriteString("false\n")
			}
		}
		return nil

	case ast.VarDecStmt:
		val := ev.Eval(n.Expr, env)
		env.Set(n.Name, val)
		return nil
	case ast.InputStmt:
		ev.Out.WriteString("Vvedit znachennya: ")
		ev.Flush()

		reader := bufio.NewReader(os.Stdin)
		inputStr, _ := reader.ReadString('\n')
		inputStr = strings.TrimSpace(inputStr)

		val, _ := strconv.ParseFloat(inputStr, 64)
		env.Set(n.Name, val)
		return nil
	case ast.AssignStmt:
		val := ev.Eval(n.Expr, env)
		env.Set(n.Name, val)
		return nil
	case ast.ExprStmt:
		return ev.Eval(n.Expr, env)
	case ast.CallExpr:
		if n.Name == "dovzhyna" {
			if len(n.Args) != 1 {
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (dovzhyna() ochikuye rivno 1 arhument!)\n")
				return nil
			}

			val := ev.Eval(n.Args[0], env)

			if list, ok := val.([]interface{}); ok {
				return float64(len(list))
			}

			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (dovzhyna() pratsyuye tilky zi spyskamy!)\n")
			return nil
		}
		fnDefIf, ok := ev.Globals.Get(n.Name)
		if !ok {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Unknown function: " + n.Name + ")\n")
			ev.Flush()
			os.Exit(1)
		}
		fnDef := fnDefIf.(ast.FuncDefStmt)
		fnEnv := NewEnv(ev.Globals)
		for _, stmt := range fnDef.Body {
			res := ev.Eval(stmt, fnEnv)
			if ret, isRet := res.(ReturnValue); isRet {
				return ret.Value
			}
		}
		return nil
	case ast.ReturnStmt:
		val := ev.Eval(n.Expr, env)
		return ReturnValue{Value: val}
	case ast.IfStmt:
		condVal := ev.Eval(n.Condition, env)
		cond, ok := condVal.(bool)
		if !ok {
			cond = false
		}
		if cond {
			for _, stmt := range n.Body {
				res := ev.Eval(stmt, env)
				if _, isRet := res.(ReturnValue); isRet {
					return res
				}
			}
			return nil
		}
		for _, elif := range n.ElseIfs {
			elifCondVal := ev.Eval(elif.Condition, env)
			elifCond, ok := elifCondVal.(bool)
			if !ok {
				elifCond = false
			}
			if elifCond {
				for _, stmt := range elif.Body {
					res := ev.Eval(stmt, env)
					if _, isRet := res.(ReturnValue); isRet {
						return res
					}
				}
				return nil
			}
		}
		if n.ElseBody != nil {
			for _, stmt := range n.ElseBody {
				res := ev.Eval(stmt, env)
				if _, isRet := res.(ReturnValue); isRet {
					return res
				}
			}
		}
		return nil
	case ast.InfixExpr:
		leftVal := ev.Eval(n.Left, env)
		rightVal := ev.Eval(n.Right, env)
		left, ok1 := leftVal.(float64)
		right, ok2 := rightVal.(float64)
		if !ok1 || !ok2 {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Type mismatch)\n")
			ev.Flush()
			os.Exit(1)
		}
		switch n.Operator {
		case "+":
			return left + right
		case "-":
			return left - right
		case "*":
			return left * right
		case "/":
			if right == 0 {
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Zero divide)\n")
				return nil
			}
			return left / right
		case ">":
			return left > right
		case "<":
			return left < right
		case "stupin":
			return math.Pow(left, right)
		}
	case ast.PrefixExpr:
		rightVal := ev.Eval(n.Right, env)
		right, ok := rightVal.(float64)
		if !ok {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Type mismatch)\n")
			ev.Flush()
			os.Exit(1)
		}
		switch n.Operator {
		case "korin":
			return math.Sqrt(right)
		case "loh10":
			return math.Log10(right)
		case "abs":
			return math.Abs(right)
		case "arksyn":
			return math.Asin(right)
		case "kosynus":
			return math.Cos(right)
		}
	case ast.NumberLiteral:
		return n.Value
	case ast.StringLiteral:
		return n.Value
	case ast.Identifier:
		v, ok := env.Get(n.Value)
		if !ok {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Undefined variable: " + n.Value + ")\n")
			ev.Flush()
			os.Exit(1)
		}
		return v
	case ast.ListLiteral:
		elements := make([]interface{}, len(n.Elements))
		for i, expr := range n.Elements {
			elements[i] = ev.Eval(expr, env)
		}
		return elements
	case ast.IndexExpr:
		leftVal := ev.Eval(n.Left, env)
		idxVal := ev.Eval(n.Index, env)

		list, okList := leftVal.([]interface{})
		index, okIdx := idxVal.(float64)

		if !okList {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Mozhna tykaty palcem tilky v spysok!)\n")
			return nil
		}
		if !okIdx {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Index maye buty chyslom!)\n")
			return nil
		}

		i := int(index)
		if i < 0 || i >= len(list) {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Takogo elementa nemaye, ty khochesh zanadto bahato!)\n")
			return nil
		}
		return list[i]
	case ast.PokyStmt:
		for {
			cond := ev.Eval(n.Condition, env)

			if boolean, ok := cond.(bool); ok && !boolean {
				break
			}

			for _, stmt := range n.Body {
				ev.Eval(stmt, env)
			}
		}
		return nil
	}
	return nil
}
