package dsl

import (
	"fmt"
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

	for e := node.Front(); e != nil; e = node.Next() {
		edge1 := intToString(node.NodeId())
		edge2 := intToString(e.NodeId())

		edge := "\t" + edge1 + " -> " + edge2
		label := "\t" + edge1 + " [label=\"" + createLabel(node) + "\"];" + "\n"
		label += "\t" + edge2 + " [label=\"" + createLabel(e) + "\"];"

		*edges = append(*edges, edge)
		*labels = append(*labels, label)
	}

	for e := node.Front(); e != nil; e = node.Next() {
		scan(e, edges, labels)
	}
}

// Convert converts the tree into DOT format
func generateDotFormat(node Node, outputfile string) error {
	file, err := os.Create(outputfile)
	defer file.Close()
	if err != nil {
		return err
	}
	result := "digraph G {" + "\n"
	slice1 := []string{}
	slice2 := []string{}
	scan(node, &slice1, &slice2)
	result += strings.Join(slice1, "\n")
	result += "\n"
	result += strings.Join(slice2, "\n")
	result += "\n}"
	_, err = file.WriteString(result)
	return err
}

// Plot plots the AST into SVG format, therefor converts the DOT format to SVG format
func plot(node Node, outputfile string) error {
	err := generateDotFormat(node, "output.dot")
	if err != nil {
		return err
	}
	// func Command(name string, arg ...string) *Cmd
	// Command returns the Cmd struct to execute the named program with the given arguments.
	// windows:
	//cmd := exec.Command("cmd", "/C", "dot -Tpdf "+"output.dot"+" -o "+outputfile)
	cmd := exec.Command("sh", "-c", "dot -Tpdf "+"output.dot"+" -o "+outputfile)
	err = cmd.Run()
	return err
}

func open(outputfile string) error {
	// func Command(name string, arg ...string) *Cmd
	// Command returns the Cmd struct to execute the named program with the given arguments.
	// windows:
	//cmd := exec.Command("cmd", "/C start "+outputfile)
	cmd := exec.Command("sh", "-c", "open "+outputfile)
	err := cmd.Run()
	return err
}
