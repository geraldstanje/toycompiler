package dsl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

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

func CreateLabel(node Node) string {
	switch n := node.(type) {
	case *ProgramNode:
		return "Program"

	case *TokenNode:
		return n.Token

	case *AssignNode:
		return "="

	case *OpNode:
		return n.Operator

	case *PrintNode:
		return "Print"

	case *WhileNode:
		return "While"

	default:
		fmt.Printf("ast.Walk: unexpected node type %T", n)
		panic("ast.Walk")
	}
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

// Scan scans all nodes in the tree recursively
func scan(T Node, edges *[]string, labels *[]string) {
	if T == nil {
		return
	}
	if T.Left() != nil {
		edge1 := IntToString(T.NodeId())
		edge2 := IntToString(T.Left().NodeId())

		edge := "\t" + edge1 + " -> " + edge2
		label := "\t" + edge1 + " [label=\"" + CreateLabel(T) + "\"];" + "\n"
		label += "\t" + edge2 + " [label=\"" + CreateLabel(T.Left()) + "\"];"

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}
	if T.Right() != nil {
		edge1 := IntToString(T.NodeId())
		edge2 := IntToString(T.Right().NodeId())

		edge := "\t" + edge1 + " -> " + edge2
		label := "\t" + edge1 + " [label=\"" + CreateLabel(T) + "\"];" + "\n"
		label += "\t" + edge2 + " [label=\"" + CreateLabel(T.Right()) + "\"];"

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}
	scan(T.Left(), edges, labels)
	scan(T.Right(), edges, labels)
}

// Convert converts the tree into DOT format
func generateDotFormat(T Node, outputfile string) {
	file, err := os.Create(outputfile)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	result := "digraph G {" + "\n"
	slice1 := []string{}
	slice2 := []string{}
	scan(T, &slice1, &slice2)
	result += strings.Join(slice1, "\n")
	result += "\n"
	result += strings.Join(slice2, "\n")
	result += "\n}"
	file.WriteString(result)
}

// Plot plots the AST into SVG format, therefor converts the DOT format to SVG format
func Plot(T Node, outputfile string) {
	generateDotFormat(T, "output.dot")
	// func Command(name string, arg ...string) *Cmd
	// Command returns the Cmd struct to execute the named program with the given arguments.
	// windows:
	//cmd := exec.Command("cmd", "/C", "dot -Tpdf "+"output.dot"+" -o "+outputfile)
	cmd := exec.Command("sh", "-c", "dot -Tpdf "+"output.dot"+" -o "+outputfile)
	er := cmd.Run()
	if er != nil {
		log.Fatal(er)
	}
}

func Open(outputfile string) {
	// func Command(name string, arg ...string) *Cmd
	// Command returns the Cmd struct to execute the named program with the given arguments.
	// windows:
	//cmd := exec.Command("cmd", "/C start "+outputfile)
	cmd := exec.Command("sh", "-c", "open "+outputfile)
	er := cmd.Run()
	if er != nil {
		log.Fatal(er)
	}
}
