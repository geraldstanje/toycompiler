package dsl

import (
	"testing"
)

func TestMyDsl(t *testing.T) {
	c, err := NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	err = c.Parse("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = c.PlotAst("plot.pdf")
	if err != nil {
		t.Fatal(err)
	}
	c.CompTopScope()
}
