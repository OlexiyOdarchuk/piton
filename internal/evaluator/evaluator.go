package evaluator

import (
	"bufio"
	"io"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/OlexiyOdarchuk/piton/internal/ast"
	"github.com/OlexiyOdarchuk/piton/internal/lexer"
	"github.com/OlexiyOdarchuk/piton/internal/parser"
)

type Evaluator struct {
	Globals *Environment
	Out     *bufio.Writer
}

func New(out io.Writer) *Evaluator {
	env := NewEnv(nil)
	env.Set("true", true)
	env.Set("false", false)
	return &Evaluator{
		Globals: env,
		Out:     bufio.NewWriter(out),
	}
}

var (
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func SeedRandom(seed int64) {
	rnd = rand.New(rand.NewSource(seed))
}

const (
	ansiReset = "\x1b[0m"
)

var colorCodes = map[string]string{
	"black":          "\x1b[30m",
	"red":            "\x1b[31m",
	"green":          "\x1b[32m",
	"yellow":         "\x1b[33m",
	"blue":           "\x1b[34m",
	"magenta":        "\x1b[35m",
	"cyan":           "\x1b[36m",
	"white":          "\x1b[37m",
	"bright_black":   "\x1b[90m",
	"bright_red":     "\x1b[91m",
	"bright_green":   "\x1b[92m",
	"bright_yellow":  "\x1b[93m",
	"bright_blue":    "\x1b[94m",
	"bright_magenta": "\x1b[95m",
	"bright_cyan":    "\x1b[96m",
	"bright_white":   "\x1b[97m",
}

func detectAnsi() bool {
	if os.Getenv("NO_COLOR") != "" {
		return false
	}
	switch runtime.GOOS {
	case "plan9":
		return false
	case "windows":
		return true
	}
	term := os.Getenv("TERM")
	if term == "" || term == "dumb" {
		return false
	}
	return true
}

func stringifyForConcat(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return strconv.FormatFloat(val, 'g', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return ""
	}
}

func colorize(name, text string) string {
	if !detectAnsi() {
		return text
	}
	code, ok := colorCodes[strings.ToLower(name)]
	if !ok {
		return text
	}
	return code + text + ansiReset
}

func isStringLike(v interface{}) bool {
	switch v.(type) {
	case string, bool:
		return true
	default:
		return false
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
		ev.Flush()
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
		if n.Receiver == nil {
			if n.Name == "zaokruhlennya" {
				val := ev.Eval(n.Args[0], env)
				f, ok := val.(float64)
				if !ok {
					return val
				}

				precision := 0.0
				if len(n.Args) == 2 {
					pVal := ev.Eval(n.Args[1], env)
					if p, ok := pVal.(float64); ok {
						precision = p
					}
				}

				pow := math.Pow(10, precision)
				return math.Round(f*pow) / pow
			}
			if n.Name == "zatrymka" {
				if len(n.Args) != 1 {
					ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (zatrymka() ochikuye rivno 1 arhument!)\n")
					return nil
				}

				zatrymka := ev.Eval(n.Args[0], env)
				if chas, ok := zatrymka.(float64); ok {
					d := time.Duration(chas * float64(time.Second))
					time.Sleep(d)
					return nil
				}
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (zatrymka() ochikue chislo v secundah!)\n")
				return nil
			}
			if n.Name == "chas" {
				if len(n.Args) != 0 {
					ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (chas() ne ochikuye arhumentiv!)\n")
					return nil
				}
				return float64(time.Now().UnixMicro()) / 1_000_000.0
			}
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
			if n.Name == "vypadkovo" {
				if len(n.Args) != 2 && len(n.Args) != 1 {
					ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (vypadkovo() ochikuye 1 abo 2 arhumentu!)\n")
					return nil
				}

				if len(n.Args) == 1 {
					val := ev.Eval(n.Args[0], env)

					if list, ok := val.([]interface{}); ok {
						if len(list) == 0 {
							return nil
						}
						return list[rnd.Intn(len(list))]
					}

					if f, ok := val.(float64); ok {
						if int(f) <= 0 {
							return 0.0
						}
						return float64(rnd.Intn(int(f)))
					}
				}

				if len(n.Args) == 2 {
					s := ev.Eval(n.Args[0], env)
					e := ev.Eval(n.Args[1], env)

					start, ok1 := s.(float64)
					end, ok2 := e.(float64)

					if ok1 && ok2 {
						minimum := int(start)
						maximum := int(end)

						if maximum <= minimum {
							ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (vypadkovo() minimum ne mozhe buty > za maximum!)\n")
							return nil
						}

						return float64(rnd.Intn(maximum-minimum) + minimum)
					}
				}
			}
			if n.Name == "kolor" {
				if len(n.Args) != 2 {
					ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (kolor() ochikuye rivno 2 argumenty!)\n")
					return nil
				}
				colorArg := ev.Eval(n.Args[0], env)
				text := ev.Eval(n.Args[1], env)

				colorName, ok := colorArg.(string)
				if !ok {
					ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (kolor() nazva kolory ma buty ryadkom!)\n")
					return nil
				}

				return colorize(colorName, stringifyForConcat(text))
			}
			if n.Name == "dodaty" {
				if len(n.Args) != 2 {
					ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (dodaty() ochikuye rivno 2 arhumentu!)\n")
					return nil
				}

				list := ev.Eval(n.Args[0], env)
				element := ev.Eval(n.Args[1], env)

				if arr, ok := list.([]interface{}); ok {
					if secondArr, ok := element.([]interface{}); ok {
						return append(arr, secondArr...)
					}
					return append(arr, element)
				}

				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (dodaty() pratsyuye tilky zi spyskamy!)\n")
				return nil
			}
		}

		targetEnv := ev.Globals
		if n.Receiver != nil {
			receiverVal := ev.Eval(n.Receiver, env)
			moduleEnv, ok := receiverVal.(*Environment)
			if !ok {
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (nema takogo modula!)\n")
				return nil
			}
			targetEnv = moduleEnv
		}

		fnDefIf, ok := targetEnv.Get(n.Name)
		if !ok {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Unknown function: " + n.Name + ")\n")
			ev.Flush()
			os.Exit(1)
		}
		fnDef := fnDefIf.(ast.FuncDefStmt)
		if len(n.Args) != len(fnDef.Params) {
			ev.Out.WriteString("Ryadok [-]: A tak yak ty pyshesh, tak buty ne maye! (funkciya " + fnDef.Name + " ochikuye " + strconv.Itoa(len(fnDef.Params)) + " argumentiv, a ty dav " + strconv.Itoa(len(n.Args)) + ")\n")
			return nil
		}
		fnEnv := NewEnv(ev.Globals)
		for i, param := range fnDef.Params {
			val := ev.Eval(n.Args[i], env)
			fnEnv.Set(param, val)
		}
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

		switch n.Operator {
		case "ta", "abo":
			l, ok1 := leftVal.(bool)
			r, ok2 := rightVal.(bool)
			if ok1 && ok2 {
				switch n.Operator {
				case "ta":
					return l && r
				case "abo":
					return l || r
				}
			}
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Operator I dlya ne logiki? Ty serjozno? (Type mismatch))\n")
			return nil
		case "+":
			if isStringLike(leftVal) || isStringLike(rightVal) {
				return stringifyForConcat(leftVal) + stringifyForConcat(rightVal)
			}
			l, ok1 := leftVal.(float64)
			r, ok2 := rightVal.(float64)
			if !ok1 || !ok2 {
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Type mismatch)\n")
				ev.Flush()
				os.Exit(1)
			}
			return l + r

		case "*":
			s, isString := leftVal.(string)
			num, isFloat := rightVal.(float64)

			if isString && isFloat {
				return strings.Repeat(s, int(num))
			}
			s2, isString2 := rightVal.(string)
			n2, isFloat2 := leftVal.(float64)
			if isString2 && isFloat2 {
				return strings.Repeat(s2, int(n2))
			}

			return leftVal.(float64) * rightVal.(float64)

		case "-", "/", ">", ">=", "<", "<=", "==", "!=", "stupin":
			l, ok1 := leftVal.(float64)
			r, ok2 := rightVal.(float64)
			if !ok1 || !ok2 {
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Matematyka dlya ryadkiv? Ty serjozno? (Type mismatch))\n")
				return nil
			}

			switch n.Operator {
			case "-":
				return l - r
			case "/":
				if r == 0 {
					ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Na nul dilyty ne mozhna, navit u Pitoni!)\n")
					return nil
				}
				return l / r
			case ">":
				return l > r
			case ">=":
				return l >= r
			case "<":
				return l < r
			case "<=":
				return l <= r
			case "==":
				return l == r
			case "!=":
				return l != r
			case "stupin":
				return math.Pow(l, r)
			}
		}
	case ast.PrefixExpr:
		rightVal := ev.Eval(n.Right, env)
		right, ok := rightVal.(float64)
		if !ok {
			rightBool, ok := rightVal.(bool)
			if ok {
				if n.Operator == "ne" {
					return !rightBool
				}
			}
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
	case ast.SpysokLiteral:
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
	case ast.SpysokExpr:
		listVal := ev.Eval(n.Left, env)
		arr, ok := listVal.([]interface{})
		if !ok {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Spysok-zriz pratsyuye tilky zi spyskamy!)\n")
			return nil
		}

		parseIndex := func(expr ast.Expr) (int, bool) {
			val := ev.Eval(expr, env)
			number, ok := val.(float64)
			if !ok {
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Index zrizu maye buty chislom!)\n")
				return 0, false
			}
			if math.Trunc(number) != number {
				ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Index zrizu maye buty cilum chislom!)\n")
				return 0, false
			}
			return int(number), true
		}

		start := 0
		if n.Start != nil {
			if idx, ok := parseIndex(n.Start); ok {
				start = idx
			} else {
				return nil
			}
		}

		end := len(arr)
		if n.End != nil {
			if idx, ok := parseIndex(n.End); ok {
				end = idx
			} else {
				return nil
			}
		}

		if start < 0 || end < 0 || start > len(arr) || end > len(arr) || start > end {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (Indexu zrizu poza zrizom!)\n")
			return nil
		}

		result := make([]interface{}, end-start)
		copy(result, arr[start:end])
		return result
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
	case ast.ImportStmt:
		filenameVal := ev.Eval(n.Filename, env)
		filename, ok := filenameVal.(string)
		if !ok {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (vykorystaty ochikue nazvy faila!)\n")
			return nil
		}
		content, err := os.ReadFile(filename + ".piton")
		if err != nil {
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (pomylka chitannya faila!)\n")
			return nil
		}
		tokens := lexer.Tokenize(string(content))
		p := parser.New(tokens)
		importedProgram := p.ParseProgram()

		moduleEnv := NewEnv(ev.Globals)
		ev.Eval(importedProgram, moduleEnv)

		moduleName := filepath.Base(filename)
		ev.Globals.Set(moduleName, moduleEnv)

		moduleEnv.ForEach(func(name string, val interface{}) {
			if _, exists := ev.Globals.Get(name); !exists {
				ev.Globals.Set(name, val)
			}
		})
		return nil
	case ast.SelectorExpr:
		leftVal := ev.Eval(n.Left, env)

		if module, ok := leftVal.(*Environment); ok {
			val, ok := module.Get(n.Right)
			if ok {
				return val
			}
			ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (u modula nema takogo polya!)\n")
			return nil
		}
		ev.Out.WriteString("Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (nema takogo modula!)\n")
		return nil
	}
	return nil
}
