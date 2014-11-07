package dsl

var count int
var blockNb int

type Node interface {
	NodeId() int
	Front() Node
	Next() Node
	AppendChild(Node)
}

type BasicNode struct {
	BasicNodeId int
	currChild   int
	children    []Node
}

func CreateBasicNode(id int) BasicNode {
	b := BasicNode{BasicNodeId: count}
	b.children = make([]Node, 0)
	return b
}

func (b *BasicNode) AppendChild(n Node) {
	b.children = append(b.children, n)
}

func (b *BasicNode) NodeId() int {
	return b.BasicNodeId
}

func (b *BasicNode) Front() Node {
	b.currChild = 0
	if len(b.children) < 1 {
		return nil
	}
	return b.children[b.currChild]
}

func (b *BasicNode) Next() Node {
	b.currChild = b.currChild + 1
	if b.currChild >= len(b.children) {
		return nil
	}
	return b.children[b.currChild]
}

type DeclarationNode struct {
	BasicNode
}

type FunctionDeclNode struct {
	BasicNode
	name string
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

func newDeclarationNode(l Node) Node {
	b := CreateBasicNode(count)
	b.AppendChild(l)

	e := &DeclarationNode{
		BasicNode: b,
	}
	count++
	return e
}

func newFunctionDeclNode(funcName string, l Node, r Node) Node {
	b := CreateBasicNode(count)

	if r != nil {
		b.AppendChild(l)
		b.AppendChild(r)

		e := &FunctionDeclNode{
			BasicNode: b,
			name:      funcName,
		}
		count++
		return e
	}

	b.AppendChild(l)
	e := &FunctionDeclNode{
		BasicNode: b,
		name:      funcName,
	}
	count++
	return e
}

func newProgramNode(l Node, r Node) Node {
	b := CreateBasicNode(count)

	if r != nil {
		b.AppendChild(l)
		b.AppendChild(r)

		e := &ProgramNode{
			BasicNode: b,
		}
		count++
		return e
	}

	b.AppendChild(l)
	e := &ProgramNode{
		BasicNode: b,
	}
	count++
	return e
}

func newAssignNode(l Node, r Node) Node {
	b := CreateBasicNode(count)
	b.AppendChild(l)
	b.AppendChild(r)

	e := &AssignNode{
		BasicNode: b,
	}
	count++
	return e
}

func newTokenNode(str string) Node {
	b := CreateBasicNode(count)

	e := &TokenNode{
		BasicNode: b,
		Token:     str,
	}
	count++
	return e
}

func newOpNode(str string, l Node, r Node) Node {
	b := CreateBasicNode(count)
	b.AppendChild(l)
	b.AppendChild(r)

	e := &OpNode{
		BasicNode: b,
		Operator:  str,
	}
	count++
	return e
}

func newWhileNode(l Node, r Node) Node {
	b := CreateBasicNode(count)
	b.AppendChild(l)
	b.AppendChild(r)

	e := &WhileNode{
		BasicNode: b,
	}
	count++
	return e
}

func newPrintNode(l Node) Node {
	b := CreateBasicNode(count)
	b.AppendChild(l)

	e := &PrintNode{
		BasicNode: b,
	}
	count++
	return e
}
