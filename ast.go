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

func newProgramNode(expr yySymType) (*Expr, error) {
	//fmt.Println("NewProgramNode called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "Program"
	e.Left = nil
	e.Right = nil
	return e, nil
}

func newAssignNode(expr yySymType) (*Expr, error) {
	//fmt.Println("NewAssignExpr called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "="
	e.Left = nil
	e.Right = nil
	return e, nil
}

func newIdentifierNode(expr yySymType) (*Expr, error) {
	//fmt.Println("NewAssignExpr called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "Indentifier"
	e.Left = nil
	e.Right = nil
	return e, nil
}

// Walk down tree
//func (tree *Tree) Walk(level int) (*Tree, error) {

//}
