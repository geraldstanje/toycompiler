package dsl

import (
	"log"
	"os"
)

type Compiler struct {
	ast *Node
}

func NewCompiler() (c *Compiler) {
	return &Compiler{ast: nil}
}

func (c *Compiler) SetAstRoot(root *Node) {
	if c.ast == nil {
		c.ast = root
	}
}

func (c *Compiler) CreateAst(filename string) {
	file, err := os.Open("test.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	yyParse(NewLexerWithInit(file, func(y *Lexer) { y.p = c }))
}

func (c *Compiler) PlotAst(filename string) {
	Plot(c.ast, filename)
	Open(filename)
}
