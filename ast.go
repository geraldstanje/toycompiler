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

func newProgramNode(expr *Expr) *Expr {
	//fmt.Println("NewProgramNode called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "Program"
	e.Left = nil
	e.Right = nil
	return e
}

func newAssignNode(expr *Expr) *Expr {
	//fmt.Println("NewAssignExpr called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "="
	e.Left = nil
	e.Right = nil
	return e
}

func newIdentifierNode(expr *Expr) *Expr {
	//fmt.Println("NewAssignExpr called", name)

	e := new(Expr)
	e.Kind = 0
	e.Type = "Indentifier"
	e.Left = nil
	e.Right = nil
	return e
}

func newNumberNode(expr *Expr) *Expr {
	e := new(Expr)
	e.Kind = 0
	e.Type = "Number"
	e.Left = nil
	e.Right = nil
	return e
}

// Walk down tree
//func (tree *Tree) Walk(level int) (*Tree, error) {

//}
