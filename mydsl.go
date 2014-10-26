package dsl

import (
	"fmt"
	"log"
	"os"
	//"strings"
	//"strconv"
	//"strings"
)

type MyDsl struct {
	ast *Node
}

func Create() (d *MyDsl) {
	return &MyDsl{}
}

func (d *MyDsl) InitAST(root *Node) {
	if d.ast == nil {
		d.ast = root
	}
}

// Walk traverses a tree depth-first
func (d *MyDsl) Walk(n *Node) {
  if n == nil {
    return
  }
  d.Walk(t.Left)
  fmt.Println(n.Type)
  d.Walk(t.Right)
}

func (d *MyDsl) Init(filename string) {
	file, err := os.Open("test.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	yyParse(NewLexerWithInit(file, func(y *Lexer) { y.p = d }))

	fmt.Println(d.ast.Type)
	//for _, node := range d.ast {
	//	fmt.Println(node.Type)
	//}
}
