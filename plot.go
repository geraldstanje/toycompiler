package dsl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func intToString(value int) string {
	return strconv.Itoa(value)
}

func createLabel(node Node) string {
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
		fmt.Printf("CreateLabel: unexpected node type %T", n)
		panic("CreateLabel")
	}
}

// Scan scans all nodes in the tree recursively
func scan(node Node, edges *[]string, labels *[]string) {
	if node == nil {
		return
	}
	if node.Left() != nil {
		edge1 := intToString(node.NodeId())
		edge2 := intToString(node.Left().NodeId())

		edge := "\t" + edge1 + " -> " + edge2
		label := "\t" + edge1 + " [label=\"" + createLabel(node) + "\"];" + "\n"
		label += "\t" + edge2 + " [label=\"" + createLabel(node.Left()) + "\"];"

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}
	if node.Right() != nil {
		edge1 := intToString(node.NodeId())
		edge2 := intToString(node.Right().NodeId())

		edge := "\t" + edge1 + " -> " + edge2
		label := "\t" + edge1 + " [label=\"" + createLabel(node) + "\"];" + "\n"
		label += "\t" + edge2 + " [label=\"" + createLabel(node.Right()) + "\"];"

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}
	scan(node.Left(), edges, labels)
	scan(node.Right(), edges, labels)
}

// Convert converts the tree into DOT format
func generateDotFormat(node Node, outputfile string) {
	file, err := os.Create(outputfile)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}
	result := "digraph G {" + "\n"
	slice1 := []string{}
	slice2 := []string{}
	scan(node, &slice1, &slice2)
	result += strings.Join(slice1, "\n")
	result += "\n"
	result += strings.Join(slice2, "\n")
	result += "\n}"
	file.WriteString(result)
}

// Plot plots the AST into SVG format, therefor converts the DOT format to SVG format
func plot(node Node, outputfile string) {
	generateDotFormat(node, "output.dot")
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

func open(outputfile string) {
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
