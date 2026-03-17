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
	Name string
	Body []Stmt
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
	Name string
	Expr Expr
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

type CallExpr struct {
	Name string
	Args []Expr
}

func (CallExpr) expr() {}

type InfixExpr struct {
	Left     Expr
	Operator string
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
