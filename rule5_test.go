package main

import (
	"fmt"
	"testing"

	. "github.com/mjwong/sudoku2/linkedlist"
	. "github.com/mjwong/sudoku2/linkedlistpair"
)

func TestCountPairs(t *testing.T) {
	node1 := &Pair{
		A: &Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node2 := &Pair{
		A: &Cell{
			Row:  1,
			Col:  3,
			Vals: []int{2, 4},
		},
		B: &Cell{
			Row:  4,
			Col:  8,
			Vals: []int{2, 4},
		},
	}

	node3 := &Pair{
		A: &Cell{
			Row:  1,
			Col:  2,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  2,
			Col:  3,
			Vals: []int{2, 8},
		},
	}

	pairList := &LinkedListPairs{}
	pairList.AddNode(node1)
	pairList.AddNode(node2)
	pairList.AddNode(node3)
	if pairList.CountNodes() != 3 {
		t.Fatalf("Count should be 3 but got %d.\n", pairList.CountNodes())
	}

	if pairList.Head != node1 {
		t.Fatalf("Head node should be %v but got %v.\n", node1, pairList.Head)
	}

	if pairList.Head.Next != node2 {
		t.Fatalf("Second node should be %v but got %v.\n", node2, pairList.Head.Next)
	}

	if pairList.Last != node3 {
		t.Fatalf("Last node should be %v but got %v.\n", node3, pairList.Last)
	}

	if pairList.Head.Next.Next.Next != nil {
		t.Fatalf("After last node should be nil but got %v.\n", pairList.Head.Next.Next.Next)
	}

	if *debugPtr {
		cNode := pairList.Head
		for cNode != nil {
			fmt.Printf("Pair1: row: %d col: %d. Pair2: row: %d col: %d. Next: %v.\n", cNode.A.Row, cNode.A.Col, cNode.B.Row, cNode.B.Col, cNode.Next)
			cNode = cNode.Next
		}
	}
}

func TestContainsPair(t *testing.T) {
	node1 := &Pair{
		A: &Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node2 := &Pair{
		A: &Cell{
			Row:  1,
			Col:  3,
			Vals: []int{2, 4},
		},
		B: &Cell{
			Row:  4,
			Col:  8,
			Vals: []int{2, 4},
		},
	}

	node3 := &Pair{
		A: &Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node4 := &Pair{
		A: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
	}

	pairList := &LinkedListPairs{}
	pairList.AddNode(node1)
	pairList.AddNode(node2)

	if !pairList.Contains(node3) {
		t.Fatalf("The pairs should be the same but are not.\n")
	}
	if !pairList.Contains(node4) { // the cells are reversed.
		t.Fatalf("The reversed pairs should be the same but are not.\n")
	}
}

// Rule 5: Naked Pairs
func TestRule5(t *testing.T) {
	// mid-way through difficiult3¯
	input := "142.73...597.462.3863.52...31852469772639.4.545976.32.6.54391.293128....2.461..39"

	ruleTest(t, input, 5, 23, 10)

}

func TestRule5a(t *testing.T) {
	// difficult1.txt
	// naked pair (2,8) found in starting possibility matrix at cells [2,1] and [1,2] of block [0,0].

	ruleTest(t, difficult1, 5, 51, 1)
}
