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
}

func (d *MyDsl) BeginForm(line string) {
	fmt.Println("BeginForm called", line)
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
}
