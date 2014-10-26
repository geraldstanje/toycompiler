package dsl

import (
//"fmt"
)

type Node struct {
	Kind  int
	Type  string
	Left  *Node
	Right *Node
}

func newProgramNode(name yySymType) *Node {
	e := new(Node)
	e.Kind = 0
	e.Type = "Program"
	e.Left = nil
	e.Right = nil
	return e
}

func newAssignNode(name yySymType) *Node {
	e := new(Node)
	e.Kind = 0
	e.Type = "="
	e.Left = nil
	e.Right = nil
	return e
}

func newIdentifierNode(name yySymType) *Node {
	e := new(Node)
	e.Kind = 0
	e.Type = "Indentifier"
	e.Left = nil
	e.Right = nil
	return e
}

func newNumberNode(name yySymType) *Node {
	e := new(Node)
	e.Kind = 0
	e.Type = "Number"
	e.Left = nil
	e.Right = nil
	return e
}

// Walk down tree
//func (tree *Tree) Walk(level int) (*Tree, error) {

//}
