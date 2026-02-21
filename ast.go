// ast.go
package main

import "fmt"

type Node interface {
	TokenLiteral() string
}

type Program struct {
	Functions map[string]*FunctionDecl
	TopStmts  []Stmt // for top-level statements
}

type FunctionDecl struct {
	Name   string
	Params []Param
	Return string
	Body   Expr
}

type Param struct {
	Type string
	Name string
}

type Stmt interface {
	Node
	stmtNode()
}

type Expr interface {
	Node
	exprNode()
}

type LetStmt struct {
	Mutable bool
	Type    string
	Name    string
	Value   Expr
}

type AssignStmt struct {
	Name  string
	Value Expr
}

type ExprStmt struct {
	Expr Expr
}

type IdentExpr struct {
	Value string
}

type IntExpr struct {
	Value int64
}

type BoolExpr struct {
	Value bool
}

type StringExpr struct {
	Value string
}

type BinOpExpr struct {
	Left  Expr
	Op    TokenType
	Right Expr
}

type CallExpr struct {
	Callee Expr
	Args   []Expr
}

type IfExpr struct {
	Cond Expr
	Then Expr
	Else Expr
}

type WhileExpr struct {
	Cond Expr
	Body Expr
}

func (p *Program) TokenLiteral() string      { return "" }
func (f *FunctionDecl) TokenLiteral() string { return f.Name }
func (l *LetStmt) TokenLiteral() string      { return "let" }
func (a *AssignStmt) TokenLiteral() string   { return "=" }
func (e *ExprStmt) TokenLiteral() string     { return "" }
func (i *IdentExpr) TokenLiteral() string    { return i.Value }
func (n *IntExpr) TokenLiteral() string      { return string(n.Value) }
func (b *BoolExpr) TokenLiteral() string     { return fmt.Sprintf("%t", b.Value) }
func (s *StringExpr) TokenLiteral() string   { return s.Value }
func (o *BinOpExpr) TokenLiteral() string    { return "" }
func (c *CallExpr) TokenLiteral() string     { return "" }
func (i *IfExpr) TokenLiteral() string       { return "if" }
func (w *WhileExpr) TokenLiteral() string    { return "while" }

// Dummy to satisfy interfaces
func (l *LetStmt) stmtNode()    {}
func (a *AssignStmt) stmtNode() {}
func (e *ExprStmt) stmtNode()   {}
func (i *IdentExpr) exprNode()  {}
func (n *IntExpr) exprNode()    {}
func (b *BoolExpr) exprNode()   {}
func (s *StringExpr) exprNode() {}
func (o *BinOpExpr) exprNode()  {}
func (c *CallExpr) exprNode()   {}
func (ifx *IfExpr) exprNode()   {}
func (w *WhileExpr) exprNode()  {}
