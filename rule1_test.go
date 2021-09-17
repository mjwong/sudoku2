package main

import (
	"testing"

	. "github.com/mjwong/sudoku2/lib"
	"gopkg.in/gookit/color.v1"
)

func TestRule1(t *testing.T) {

	input := ".341528699.837645252.948371245.136986895.413737168.524857231.464938672.516249578."
	PrepPmat(input)

	matched, cnt := rule1()
	if cnt != 9 {
		t.Fatalf("Expected 8 but got %d\n", cnt)
	}

	if !CheckSums(mat) {
		t.Fatalf("There are errors in the resulting matrix.\n")
	}

	digcnt := matched.CountNodes()

	color.LightMagenta.Printf("Found: %d digits.\n", cnt)
	if digcnt != cnt {
		t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
	}
	matched.PrintResult("Found open single")
}
