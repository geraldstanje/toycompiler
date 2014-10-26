package dsl

import (
	"fmt"
)

type Expr struct {
	Kind  int
	Type  string
	Left  *Expr
	Right *Expr
}

func NewProgramNode(kind int, name yySymType) (*Expr, error) {
	fmt.Println("NewProgramNode called", name)

	e := new(Expr)
	e.Kind = 0
	e.Left = nil
	e.Right = nil
	return e, nil
}

func NewAssignExpr(kind int, name yySymType) (*Expr, error) {
	fmt.Println("NewAssignExpr called", name)

	e := new(Expr)
	e.Kind = 0
	e.Left = nil
	e.Right = nil
	return e, nil
}
