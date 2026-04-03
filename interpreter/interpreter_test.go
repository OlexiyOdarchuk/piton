package interpreter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/OlexiyOdarchuk/piton/internal/evaluator"
)

func init() {
	_ = os.Setenv("NO_COLOR", "1")
}

func runWithBuffer(t *testing.T, code string) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Run(code, &buf); err != nil {
		t.Fatalf("Run(%q) failed: %v", code, err)
	}
	return buf.String()
}

func parseFloatOutput(t *testing.T, got string) float64 {
	t.Helper()
	trimmed := strings.TrimSpace(got)
	val, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		t.Fatalf("cannot parse %q as float: %v", got, err)
	}
	return val
}

func TestRunMathOperations(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{
			name: "addition",
			code: "drukuvaty 1 + 2",
			want: "3\n",
		},
		{
			name: "stupin",
			code: "drukuvaty 2 stupin 3",
			want: "8\n",
		},
		{
			name: "korin",
			code: "drukuvaty korin 9",
			want: "3\n",
		},
		{
			name: "loh10",
			code: "drukuvaty loh10 100",
			want: "2\n",
		},
		{
			name: "abs",
			code: "drukuvaty abs 5",
			want: "5\n",
		},
		{
			name: "arksyn",
			code: "drukuvaty arksyn 0.5",
			want: "0.5235987755982989\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runWithBuffer(t, tt.code+"\n")
			if got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestZaokruhlennya(t *testing.T) {
	code := "drukuvaty zaokruhlennya(123.4567, 2)"
	if got := runWithBuffer(t, code); got != "123.46\n" {
		t.Fatalf("got %q want %q", got, "123.46\n")
	}
}

func TestKolorFunction(t *testing.T) {
	prev := os.Getenv("NO_COLOR")
	t.Cleanup(func() {
		if prev == "" {
			_ = os.Unsetenv("NO_COLOR")
			return
		}
		_ = os.Setenv("NO_COLOR", prev)
	})
	_ = os.Unsetenv("NO_COLOR")

	code := "drukuvaty kolor(\"red\", \"alert\")"
	want := "\x1b[31malert\x1b[0m\n"
	if got := runWithBuffer(t, code); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestStringConcatMix(t *testing.T) {
	tests := []struct {
		name string
		code string
		want string
	}{
		{name: "str+num", code: "drukuvaty \"hello\" + 2", want: "hello2\n"},
		{name: "num+str", code: "drukuvaty 2 + \"hello\"", want: "2hello\n"},
		{name: "str+bool", code: "drukuvaty \"hello\" + true", want: "hellotrue\n"},
		{name: "bool+str", code: "drukuvaty false + \"ok\"", want: "falseok\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runWithBuffer(t, tt.code); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestChas(t *testing.T) {
	code := "nekhay a = chas()\nnekhay b = chas()\ndrukuvaty b - a"
	val := parseFloatOutput(t, runWithBuffer(t, code))
	if val < 0 {
		t.Fatalf("chronological imbalance: %v", val)
	}
}

func TestZatrymka(t *testing.T) {
	code := "nekhay start = chas()\nzatrymka(0.05)\nnekhay elapsed = chas() - start\ndrukuvaty elapsed"
	val := parseFloatOutput(t, runWithBuffer(t, code))
	if val < 0.04 {
		t.Fatalf("expected at least ~0.05s delay, got %v", val)
	}
}

func TestVykorystatyImports(t *testing.T) {
	code := "vykorystaty \"../examples/hello\"\nhello.Hello(\"Pit\")\n"
	want := "Hello Pit!\n"
	if got := runWithBuffer(t, code); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRunLogicBranches(t *testing.T) {
	const codeTemplate = "nekhay a = %d\n" +
		"yaksho a > 3:\n" +
		"    drukuvaty \"first\"\n" +
		"inackshe yaksho a < 3:\n" +
		"    drukuvaty \"second\"\n" +
		"inackshe:\n" +
		"    drukuvaty \"third\"\n"

	tests := []struct {
		name  string
		want  string
		value int
	}{
		{name: "primary", value: 5, want: "first\n"},
		{name: "elseif", value: 2, want: "second\n"},
		{name: "else", value: 3, want: "third\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := fmt.Sprintf(codeTemplate, tt.value)
			got := runWithBuffer(t, code)
			if got != tt.want {
				t.Fatalf("value=%d got %q want %q", tt.value, got, tt.want)
			}
		})
	}
}

func TestRunFunctionsAndScoping(t *testing.T) {
	code := "" +
		"functia add():\n" +
		"    nekhay inner = 3\n" +
		"    vernuty inner + outer\n" +
		"nekhay outer = 7\n" +
		"drukuvaty add()\n" +
		"drukuvaty outer\n"

	want := "10\n7\n"
	if got := runWithBuffer(t, code); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRunFunctionArguments(t *testing.T) {
	code := "" +
		"functia repeat(label, count):\n" +
		"    nekhay i = 0\n" +
		"    poky i < count:\n" +
		"        drukuvaty label\n" +
		"        i = i + 1\n" +
		"repeat(\"hi\", 2)\n"

	want := "hi\nhi\n"
	if got := runWithBuffer(t, code); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRunFunctionArgumentMismatch(t *testing.T) {
	code := "" +
		"functia echo(value, label):\n" +
		"    drukuvaty label\n" +
		"    drukuvaty value\n" +
		"echo(1)\n"

	want := "Ryadok [-]: A tak yak ty pyshesh, tak buty ne maye! (funkciya echo ochikuye 2 argumentiv, a ty dav 1)\n"
	if got := runWithBuffer(t, code); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRunDrukuvatyOutput(t *testing.T) {
	code := "" +
		"drukuvaty \"hello\"\n" +
		"drukuvaty 2 + 3\n"

	want := "hello\n5\n"
	if got := runWithBuffer(t, code); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRunComplexExpression(t *testing.T) {
	code := "drukuvaty korin ( 2 + loh10 ( 100 ) )\n"
	want := "2\n"
	if got := runWithBuffer(t, code); got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestRunSyntaxErrorMessage(t *testing.T) {
	code := "nekhay x =\n"
	cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--", code)
	cmd.Env = append(os.Environ(), "GO_WANT_HELPER_PROCESS=1")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected interpreter to exit with an error")
	}
	output := string(out)
	if !strings.Contains(output, "Ryadok [1]") || !strings.Contains(output, "Ya tut interpretator") {
		t.Fatalf("unexpected syntax error output: %q", output)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	var code string
	for i, arg := range os.Args {
		if arg == "--" && i+1 < len(os.Args) {
			code = os.Args[i+1]
			break
		}
	}

	if code == "" {
		t.Fatalf("missing code argument")
	}

	_ = Run(code)
}

func TestVypadkovo(t *testing.T) {
	tests := []struct {
		name  string
		code  string
		check func(t *testing.T, got string)
	}{
		{
			name: "range",
			code: "drukuvaty vypadkovo(5, 10)",
			check: func(t *testing.T, got string) {
				val := parseFloatOutput(t, got)
				if val < 5 || val >= 10 {
					t.Fatalf("got %v; want value in [5, 10)", val)
				}
			},
		},
		{
			name: "single arg",
			code: "drukuvaty vypadkovo(5)",
			check: func(t *testing.T, got string) {
				val := parseFloatOutput(t, got)
				if val < 0 || val >= 5 {
					t.Fatalf("got %v; want value in [0, 5)", val)
				}
			},
		},
		{
			name: "list argument",
			code: "nekhay options = [10, 20, 30]\ndrukuvaty vypadkovo(options)",
			check: func(t *testing.T, got string) {
				val := parseFloatOutput(t, got)
				valid := map[float64]bool{10: true, 20: true, 30: true}
				if !valid[val] {
					t.Fatalf("got %v; want one of %v", val, valid)
				}
			},
		},
		{
			name: "invalid bounds",
			code: "drukuvaty vypadkovo(6, 2)",
			check: func(t *testing.T, got string) {
				want := "Ryadok [-]: Ya tut interpretator, ya znayu yak maye buty. A tak yak ty pyshesh, tak buty ne maye! (vypadkovo() minimum ne mozhe buty > za maximum!)\n"
				if got != want {
					t.Fatalf("got %q want %q", got, want)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluator.SeedRandom(42)
			tt.check(t, runWithBuffer(t, tt.code))
		})
	}
}

func TestSpysok(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want string
	}{
		{
			"Print spysok from literal",
			"drukuvaty [14,5,1,\"Hello\",2]",
			"[14, 5, 1, \"Hello\", 2]\n",
		},
		{
			"Print spysok from variable",
			"nekhay a = [14,5,1,\"Hello\",2]\ndrukuvaty a",
			"[14, 5, 1, \"Hello\", 2]\n",
		},
		{
			"All spysok",
			"nekhay s = [10, 20, 30, 40]\ndrukuvaty dovzhyna(s)\ndrukuvaty s[dovzhyna(s) - 1]",
			"4\n40\n",
		},
		{
			"Spysok range",
			"nekhay s = [1, 2, 3, 4]\ndrukuvaty s[1:3]",
			"[2, 3]\n",
		},
		{
			"Spysok to end",
			"nekhay s = [1, 2, 3, 4]\ndrukuvaty s[2:]",
			"[3, 4]\n",
		},
		{
			"Spysok from start",
			"nekhay s = [1, 2, 3, 4]\ndrukuvaty s[:2]",
			"[1, 2]\n",
		},
		{
			"Full spysok",
			"nekhay s = [1, 2, 3, 4]\ndrukuvaty s[:]",
			"[1, 2, 3, 4]\n",
		},
		{
			"Go-style removal",
			"nekhay s = [1, 2, 3, 4]\nnekhay i = 2\nnekhay trimmed = dodaty(s[:i], s[i + 1:])\ndrukuvaty trimmed",
			"[1, 2, 4]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runWithBuffer(t, tt.expr); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestPoky(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want string
	}{
		{
			"Print full script with cycle",
			"nekhay i = 0\nnekhay s = [1, 2, 3]\npoky i < dovzhyna(s):\n	drukuvaty s[i]\n	i = i + 1",
			"1\n2\n3\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runWithBuffer(t, tt.expr); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestDodaty(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want string
	}{
		{
			"Dodaty one element",
			"nekhay a = [1,2,3]\na = dodaty(a,4)\ndrukuvaty a",
			"[1, 2, 3, 4]\n",
		},
		{
			"Dodaty list to list",
			"nekhay a = [1,2,3]\nnekhay b = [4, 5, 6]\na = dodaty(a,b)\ndrukuvaty a",
			"[1, 2, 3, 4, 5, 6]\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runWithBuffer(t, tt.expr); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}
