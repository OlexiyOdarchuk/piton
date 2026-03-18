package interpreter

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func runWithBuffer(t *testing.T, code string) string {
	t.Helper()
	var buf bytes.Buffer
	if err := Run(code, &buf); err != nil {
		t.Fatalf("Run(%q) failed: %v", code, err)
	}
	return buf.String()
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
		value int
		want  string
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
		"    kinets\n" +
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

	Run(code)
	os.Exit(0)
}

func TestList(t *testing.T) {
	tests := []struct {
		name string
		expr string
		want string
	}{
		{
			"Print list from literal",
			"drukuvaty [14,5,1,\"Hello\",2]",
			"[14, 5, 1, \"Hello\", 2]\n",
		},
		{
			"Print list from variable",
			"nekhay a = [14,5,1,\"Hello\",2]\ndrukuvaty a",
			"[14, 5, 1, \"Hello\", 2]\n",
		},
		{
			"All List",
			"nekhay s = [10, 20, 30, 40]\ndrukuvaty dovzhyna(s)\ndrukuvaty s[dovzhyna(s) - 1]",
			"4\n40\n",
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
			"nekhay i = 0\nnekhay s = [1, 2, 3]\npoky i < dovzhyna(s):\n	drukuvaty s[i]\n	i = i + 1\nkinets",
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
