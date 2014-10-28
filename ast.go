package dsl

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var count int
var blockNb int

type ItemType int

// Types used in the Abstract Syntax Tree
const (
	Identifier ItemType = iota
	Token
	Assignment
	Program
	Operator
	While
	Print
)

type Node struct {
	Kind       ItemType
	NodeId     int
	Identifier string
	Token      string
	Operator   string
	Left       *Node
	Right      *Node
}

func newProgramNode() *Node {
	e := new(Node)
	e.Kind = Program
	e.NodeId = count
	e.Left = nil
	e.Right = nil
	count++
	return e
}

func newAssignNode(lval yySymType) *Node {
	e := new(Node)
	e.Kind = Assignment
	e.NodeId = count
	e.Left = nil
	e.Right = nil
	count++
	return e
}

func newIdentifierNode(lval yySymType) *Node {
	e := new(Node)
	e.Kind = Identifier
	e.NodeId = count
	e.Identifier = lval.s
	e.Left = nil
	e.Right = nil
	count++
	return e
}

func newTokenNode(lval yySymType) *Node {
	e := new(Node)
	e.Kind = Token
	e.NodeId = count
	e.Token = lval.s
	e.Left = nil
	e.Right = nil
	count++
	return e
}

func newOpNode(lval yySymType) *Node {
	e := new(Node)
	e.Kind = Operator
	e.NodeId = count
	e.Operator = lval.s
	e.Left = nil
	e.Right = nil
	count++
	return e
}

func newWhileNode() *Node {
	e := new(Node)
	e.Kind = While
	e.NodeId = count
	e.Left = nil
	e.Right = nil
	count++
	return e
}

func newPrintNode() *Node {
	e := new(Node)
	e.Kind = Print
	e.NodeId = count
	e.Left = nil
	e.Right = nil
	count++
	return e
}

func CreateLabel(T *Node) string {
	switch T.Kind {
	case Program:
		return "Program"

	case Identifier:
		return T.Identifier

	case Token:
		return T.Token

	case Assignment:
		return "="

	case Operator:
		return T.Operator

	case Print:
		return "Print"

	case While:
		return "While"

	default:
		return ""
	}
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

// Scan scans all nodes in the tree recursively
func scan(T *Node, edges *[]string, labels *[]string) {
	if T == nil {
		return
	}
	if T.Left != nil {
		edge1 := IntToString(T.NodeId)
		edge2 := IntToString(T.Left.NodeId)

		edge := "\t" + edge1 + " -> " + edge2
		label := "\t" + edge1 + " [label=\"" + CreateLabel(T) + "\"];" + "\n"
		label += "\t" + edge2 + " [label=\"" + CreateLabel(T.Left) + "\"];"

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}
	if T.Right != nil {
		edge1 := IntToString(T.NodeId)
		edge2 := IntToString(T.Right.NodeId)

		edge := "\t" + edge1 + " -> " + edge2
		label := "\t" + edge1 + " [label=\"" + CreateLabel(T) + "\"];" + "\n"
		label += "\t" + edge2 + " [label=\"" + CreateLabel(T.Right) + "\"];"

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}
	scan(T.Left, edges, labels)
	scan(T.Right, edges, labels)
}

// Convert converts the tree into DOT format
func generateDotFormat(T *Node, outputfile string) {
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
func Plot(T *Node, outputfile string) {
	generateDotFormat(T, "output.dot")
	// func Command(name string, arg ...string) *Cmd
	// Command returns the Cmd struct to execute the named program with the given arguments.
	// windows: cmd := exec.Command("cmd", "/C", "dot -Tpdf "+"output.dot"+" -o "+outputfile)
	cmd := exec.Command("sh", "-c", "dot -Tpdf "+"output.dot"+" -o "+outputfile)
	er := cmd.Run()
	if er != nil {
		log.Fatal(er)
	}
}

func Open(outputfile string) {
	// func Command(name string, arg ...string) *Cmd
	// Command returns the Cmd struct to execute the named program with the given arguments.
	// windows: cmd := exec.Command("cmd", "/C start "+outputfile)
	cmd := exec.Command("sh", "-c", "open "+outputfile)
	er := cmd.Run()
	if er != nil {
		log.Fatal(er)
	}
}
