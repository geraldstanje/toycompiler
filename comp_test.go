package dsl

import (
	//"fmt"
	//"log"
	//"os"
	"testing"
	//"strings"
	//"strconv"
	//"strings"
)

func TestMyDsl(t *testing.T) {
	c := NewCompiler()
	c.CreateAst("test.txt")
	c.PlotAst("plot.pdf")
}
