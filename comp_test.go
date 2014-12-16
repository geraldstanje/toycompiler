package dsl

import (
	"testing"
)

func TestMyDsl(t *testing.T) {
	c := NewCompiler()
	err := c.Parse("test.txt")
	if err != nil {
		t.Fatal(err)
	}
	err = c.PlotAst("plot.pdf")
	if err != nil {
		t.Fatal(err)
	}
	err = c.GenerateCode()
	if err != nil {
		t.Fatal(err)
	}
}
