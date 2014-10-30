package dsl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Compiler struct {
	ast    Node
	indent int
	writer *bufio.Writer
}

func NewCompiler() (*Compiler, error) {
	comp := Compiler{}

	file, err := os.Create("generated.txt")
	if err != nil {
		return nil, err
	}

	writer := bufio.NewWriter(file)
	comp.writer = bufio.NewWriter(writer)
	comp.indent = 0
	comp.ast = nil
	return &comp, nil
}

func (c *Compiler) SetAstRoot(root Node) {
	c.ast = root
}

func (c *Compiler) Parse(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	yyParse(NewLexerWithInit(file, func(y *Lexer) { y.p = c }))
}

func (c *Compiler) PlotAst(filename string) {
	plot(c.ast, filename)
	open(filename)
}

func (c *Compiler) EmitLine(line string) {
	c.writer.WriteString(strings.Repeat("\t", c.indent))
	c.writer.WriteString(line)
	c.writer.WriteString("\n")
}

func (c *Compiler) Indent() {
	c.indent += 1
}

func (c *Compiler) Deindent() {
	c.indent -= 1

	if c.indent < 0 {
		c.indent = 0
	}
}

func (c *Compiler) compNode(node Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *ProgramNode:
		c.compNode(n.Left())
		c.compNode(n.Right())

	case *TokenNode:
		c.EmitLine(fmt.Sprintf("PUSH %s", n.Token))

	case *AssignNode:
		left := n.Left()
		right := n.Right()

		if tn, ok := right.(*TokenNode); ok {
			c.EmitLine(fmt.Sprintf("PUSH %s", tn.Token))
		} else if _, ok := right.(*OpNode); ok {
			c.compNode(n.Right())
		}
		if tn, ok := left.(*TokenNode); ok {
			c.EmitLine(fmt.Sprintf("SET %s", tn.Token))
		}

	case *OpNode:
		left := n.Left()
		right := n.Right()

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
		left := n.Left()

		c.compNode(left)
		c.EmitLine("PRINT")

	case *WhileNode:
		left := n.Left()
		right := n.Right()
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

func (c *Compiler) CompTopScope() {
	c.compNode(c.ast)
	c.writer.Flush()
}
