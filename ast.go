package dsl

var count int
var blockNb int

type Node interface {
	NodeId() int
	Left() Node
	Right() Node
}

type BasicNode struct {
	BasicNodeId int
	left        Node
	right       Node
}

func (b BasicNode) NodeId() int {
	return b.BasicNodeId
}

func (b BasicNode) Left() Node {
	return b.left
}

func (b BasicNode) Right() Node {
	return b.right
}

type ProgramNode struct {
	BasicNode
}

type AssignNode struct {
	BasicNode
}

type TokenNode struct {
	BasicNode
	Token string
}

type OpNode struct {
	BasicNode
	Operator string
}

type WhileNode struct {
	BasicNode
}

type PrintNode struct {
	BasicNode
}

func newProgramNode(l Node, r Node) Node {
	if r != nil {
		e := &ProgramNode{
			BasicNode: BasicNode{BasicNodeId: count, left: l, right: r},
		}
		count++
		return e
	}

	e := &ProgramNode{
		BasicNode: BasicNode{BasicNodeId: count, left: l},
	}
	count++
	return e
}

func newAssignNode(l Node, r Node) Node {
	e := &AssignNode{
		BasicNode: BasicNode{BasicNodeId: count, left: l, right: r},
	}
	count++
	return e
}

func newTokenNode(str string) Node {
	e := &TokenNode{
		BasicNode: BasicNode{BasicNodeId: count},
		Token:     str,
	}
	count++
	return e
}

func newOpNode(str string, l Node, r Node) Node {
	e := &OpNode{
		BasicNode: BasicNode{BasicNodeId: count, left: l, right: r},
		Operator:  str,
	}
	count++
	return e
}

func newWhileNode(l Node, r Node) Node {
	e := &WhileNode{
		BasicNode: BasicNode{BasicNodeId: count, left: l, right: r},
	}
	count++
	return e
}

func newPrintNode(l Node) Node {
	e := &PrintNode{
		BasicNode: BasicNode{BasicNodeId: count, left: l},
	}
	count++
	return e
}
