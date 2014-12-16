package dsl

import (
	"fmt"
	"os"
)

type Compiler struct {
	err     string
	ast     Node
	codegen CodeGenerator
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

func (c *Compiler) CompTopScope() error {
	c.codegen = NewAsmCodeGenerator()
	err := c.codegen.CompTopScope(c.ast)
	return err
}
