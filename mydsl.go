package dsl

import (
	//"fmt"
	"log"
	"os"
	//"strings"
	//"strconv"
	//"strings"
)

type MyDsl struct {
	//ast []*Expr
}

func Create() (d *MyDsl) {
	return &MyDsl{}
}

func (d *MyDsl) Init(filename string) {
	file, err := os.Open("test.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	yyParse(NewLexerWithInit(file, func(y *Lexer) { y.p = d }))

	//for _, node := range d.ast {
	//	fmt.Println(node.Type)
	//}
}
