package dsl

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Compiler struct {
	err    string
	ast    Node
	indent int
	writer *bufio.Writer
}

func NewCompiler() *Compiler {
	return &Compiler{}
}

func (c *Compiler) SetAstRoot(root Node) {
	c.ast = root
}

func (c *Compiler) Parse(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	ret := yyParse(NewLexerWithInit(file, func(y *Lexer) { y.p = c }))
	// nex creates the function: func NewLexerWithInit(in io.Reader, initFun func(*Lexer)) *Lexer
	// go tool yacc creates the function: func yyParse(yylex yyLexer) int

	if ret == 1 {
		return fmt.Errorf(filename + c.err)
	}

	return nil
}

func (c *Compiler) PlotAst(filename string) error {
	err := plot(c.ast, filename)
	if err != nil {
		return err
	}
	err = open(filename)
	return err
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

func (c *Compiler) CompTopScope() error {
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
