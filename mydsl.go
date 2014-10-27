package dsl

import (
	"log"
	"os"
)

type MyDsl struct {
	ast *Node
}

func Create() (d *MyDsl) {
	return &MyDsl{ast: nil}
}

func (d *MyDsl) InitAST(root *Node) {
	if d.ast == nil {
		d.ast = root
	}
}

func (d *MyDsl) Init(filename string) {
	file, err := os.Open("test.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	yyParse(NewLexerWithInit(file, func(y *Lexer) { y.p = d }))

	Plot(d.ast, "plot.pdf")
	Open("plot.pdf")
}
