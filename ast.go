package dsl

import (
//"fmt"
)

type Expr struct {
	Kind  int
	Type  string
	Left  *Expr
	Right *Expr
}

func NewProgramNode(kind int, name yySymType) (*Expr, error) {
	//fmt.Println("NewProgramNode called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "Program"
	e.Left = nil
	e.Right = nil
	return e, nil
}

func NewAssignNode(kind int, name yySymType) (*Expr, error) {
	//fmt.Println("NewAssignExpr called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "="
	e.Left = nil
	e.Right = nil
	return e, nil
}

func NewTokenNode(kind int, name yySymType) (*Expr, error) {
	//fmt.Println("NewAssignExpr called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "token"
	e.Left = nil
	e.Right = nil
	return e, nil
}
