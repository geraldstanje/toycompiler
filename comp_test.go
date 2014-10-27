package dsl

import (
	"testing"
)

func TestMyDsl(t *testing.T) {
	c := NewCompiler()
	c.CreateAst("test.txt")
	c.PlotAst("plot.pdf")
}
