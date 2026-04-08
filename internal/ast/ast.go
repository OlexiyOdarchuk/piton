package ast

type Node interface{}
type Expr interface {
	Node
	expr()
}
type Stmt interface {
	Node
	stmt()
}

type Program struct{ Statements []Stmt }

func (Program) stmt() {}

type FuncDefStmt struct {
	Name   string
	Params []string
	Body   []Stmt
	Module string
}

func (FuncDefStmt) stmt() {}

type PrintStmt struct{ Expr Expr }

func (PrintStmt) stmt() {}

type VarDecStmt struct {
	Name string
	Expr Expr
}

func (VarDecStmt) stmt() {}

type InputStmt struct{ Name string }

func (InputStmt) stmt() {}

type AssignStmt struct {
	Target Expr
	Expr   Expr
}

func (AssignStmt) stmt() {}

type IfStmt struct {
	Condition Expr
	Body      []Stmt
	ElseIfs   []ElseIf
	ElseBody  []Stmt
}

func (IfStmt) stmt() {}

type ElseIf struct {
	Condition Expr
	Body      []Stmt
}

type ReturnStmt struct{ Expr Expr }

func (ReturnStmt) stmt() {}

type ExprStmt struct{ Expr Expr }

func (ExprStmt) stmt() {}

type PokyStmt struct {
	Body      []Stmt
	Condition Expr
}

func (PokyStmt) stmt() {}

type CallExpr struct {
	Name     string
	Receiver Expr
	Args     []Expr
}

func (CallExpr) expr() {}

type InfixExpr struct {
	Operator string
	Left     Expr
	Right    Expr
}

func (InfixExpr) expr() {}

type PrefixExpr struct {
	Operator string
	Right    Expr
}

func (PrefixExpr) expr() {}

type NumberLiteral struct{ Value float64 }

func (NumberLiteral) expr() {}

type StringLiteral struct{ Value string }

func (StringLiteral) expr() {}

type Identifier struct{ Value string }

func (Identifier) expr() {}

type SlovnykPair struct {
	Key   Expr
	Value Expr
}

type SlovnykLiteral struct{ Pairs []SlovnykPair }

func (SlovnykLiteral) expr() {}

type IndexExpr struct {
	Left  Expr
	Index Expr
}

func (IndexExpr) expr() {}

type SpysokLiteral struct{ Elements []Expr }

func (SpysokLiteral) expr() {}

type SpysokExpr struct {
	Left  Expr
	Start Expr // if nil - this is start
	End   Expr // if nil - this is end
}

func (SpysokExpr) expr() {}

type ImportStmt struct{ Filename Expr }

func (ImportStmt) stmt() {}

type SelectorExpr struct {
	Right string
	Left  Expr
}

func (SelectorExpr) expr() {}
