package matchlist

import (
	"testing"

	. "github.com/mjwong/sudoku2/linkedlist"
)

func TestAddRCell(t *testing.T) {
	cell1 := &Cell{
		Row:  4,
		Col:  5,
		Vals: []int{2, 8},
	}

	cell2 := &Cell{
		Row:  6,
		Col:  7,
		Vals: []int{2, 8},
	}

	cell3 := &Cell{
		Row:  1,
		Col:  3,
		Vals: []int{2, 4},
	}

	cell4 := &Cell{
		Row:  2,
		Col:  7,
		Vals: []int{2, 4},
	}

	arr := AddRCell(nil, cell1, cell2, cell3)

	if len(arr) != 3 {
		t.Fatalf("Expected 3 but got %d.\n", len(arr))
	}

	arr = AddRCell(arr, cell4)

	if len(arr) != 4 {
		t.Fatalf("Expected 4 but got %d.\n", len(arr))
	}
}
