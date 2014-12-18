package dsl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CodeGenerator interface {
	CompTopScope(Node) error
}

type AsmCodeGenerator struct {
	err    string
	ast    Node
	indent int
	writer *bufio.Writer
}

func NewAsmCodeGenerator() *AsmCodeGenerator {
	return &AsmCodeGenerator{}
}

func (c *AsmCodeGenerator) EmitLine(line string) {
	c.writer.WriteString(strings.Repeat("\t", c.indent))
	c.writer.WriteString(line)
	c.writer.WriteString("\n")
}

func (c *AsmCodeGenerator) Indent() {
	c.indent += 1
}

func (c *AsmCodeGenerator) Deindent() {
	c.indent -= 1

	if c.indent < 0 {
		c.indent = 0
	}
}

func (c *AsmCodeGenerator) compNode(node Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ProgramNode:
		c.compNode(n.Front())
		c.compNode(n.Next())

	case *StatementNode:
		c.compNode(n.Front())
		c.compNode(n.Next())

	case *TokenNode:
		c.EmitLine(fmt.Sprintf("PUSH %s", n.Token))

	case *AssignNode:
		left := n.Front()
		right := n.Next()

		if tn, ok := right.(*TokenNode); ok {
			c.EmitLine(fmt.Sprintf("PUSH %s", tn.Token))
		} else if _, ok := right.(*OpNode); ok {
			c.compNode(right)
		}
		if tn, ok := left.(*TokenNode); ok {
			c.EmitLine(fmt.Sprintf("SET %s", tn.Token))
		}

	case *OpNode:
		left := n.Front()
		right := n.Next()

		c.compNode(left)
		c.compNode(right)

		if n.Operator == "+" {
			c.EmitLine("ADD")
		} else if n.Operator == "-" {
			c.EmitLine("SUB")
		} else if n.Operator == "/" {
			c.EmitLine("DIV")
		} else if n.Operator == "*" {
			c.EmitLine("MUL")
		}

	case *PrintNode:
		left := n.Front()

		c.compNode(left)
		c.EmitLine("PRINT")

	case *WhileNode:
		left := n.Front()
		right := n.Next()
		blockNb = blockNb + 1

		c.EmitLine(fmt.Sprintf("JMP cond%d", blockNb))
		c.EmitLine(fmt.Sprintf("body%d%s", blockNb, ":"))

		c.Indent()
		c.compNode(right)
		c.Deindent()

		c.EmitLine(fmt.Sprintf("cond%d%s", blockNb, ":"))

		c.Indent()
		c.compNode(left)
		c.EmitLine(fmt.Sprintf("JNZ body%d", blockNb))
		c.Deindent()
	}

	return
}

func (c *AsmCodeGenerator) CompTopScope(ast Node) error {
	c.ast = ast

	file, err := os.Create("generated.txt")
	defer file.Close()
	if err != nil {
		return err
	}
	c.writer = bufio.NewWriter(file)

	c.compNode(c.ast)

	c.writer.Flush()
	return nil
}
