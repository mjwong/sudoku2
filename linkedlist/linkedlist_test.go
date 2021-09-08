package linkedlist

import "testing"

func TestAddNode(t *testing.T) {
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
		Row:  3,
		Col:  4,
		Vals: []int{3, 8},
	}

	cell4 := &Cell{
		Row:  1,
		Col:  4,
		Vals: []int{3, 4},
	}

	el := CreatelinkedList()

	el.AddNode(cell1)

	if el.CountNodes() != 1 {
		t.Fatalf("Expected 1 count but got %d.\n", el.CountNodes())
	}

	el.AddNode(cell2, cell3, cell4)

	if el.CountNodes() != 4 {
		t.Fatalf("Expected 4 count but got %d.\n", el.CountNodes())
	}
}

func TestAddCell(t *testing.T) {
	el := CreatelinkedList()

	el.AddCell(5, 5, []int{1, 2, 3})

	if el.CountNodes() != 1 {
		t.Fatalf("Expected 1 count but got %d.\n", el.CountNodes())
	}
}
