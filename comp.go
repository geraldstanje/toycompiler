package dsl

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Compiler struct {
	ast    *Node
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

func (c *Compiler) SetAstRoot(root *Node) {
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
	Plot(c.ast, filename)
	Open(filename)
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

func (c *Compiler) compNode(T *Node) {
	if T == nil {
		return
	}

	switch T.Kind {
	case Program:
		c.compNode(T.Left)
		c.compNode(T.Right)

	case Token:
		c.EmitLine(fmt.Sprintf("PUSH %s", T.Token))

	case Assignment:
		left := T.Left
		right := T.Right

		if right.Kind == Token {
			c.EmitLine(fmt.Sprintf("PUSH %s", right.Token))
		} else if right.Kind == Operator {
			c.compNode(T.Right)
		}
		c.EmitLine(fmt.Sprintf("SET %s", left.Token))

	case Operator:
		left := T.Left
		right := T.Right

		c.compNode(left)
		c.compNode(right)

		if T.Operator == "+" {
			c.EmitLine("ADD")
		} else if T.Operator == "-" {
			c.EmitLine("SUB")
		} else if T.Operator == "/" {
			c.EmitLine("DIV")
		} else if T.Operator == "*" {
			c.EmitLine("MUL")
		}

	case Print:
		left := T.Left

		c.compNode(left)
		c.EmitLine("PRINT")

	case While:
		left := T.Left
		right := T.Right
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
