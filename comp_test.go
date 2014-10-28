package dsl

import (
	"testing"
)

func TestMyDsl(t *testing.T) {
	c, err := NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	c.CreateAst("test.txt")
	c.PlotAst("plot.pdf")
	c.CompTopScope()
}
