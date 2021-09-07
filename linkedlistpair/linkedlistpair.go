package linkedlistpair

import (
	"github.com/gookit/color"

	l "github.com/mjwong/sudoku2/lib"
	ll "github.com/mjwong/sudoku2/linkedlist"
)

type Pair struct {
	A    *ll.Cell
	B    *ll.Cell
	Prev *Pair
	Next *Pair
}

type LinkedListPairs struct {
	Head        *Pair
	Last        *Pair
	CurrentPair *Pair
}

func (pl *LinkedListPairs) AddNode(thisPair *Pair) error {

	if pl.Head == nil {
		pl.Head = thisPair
		thisPair.Next = nil
		thisPair.Prev = nil
	} else {
		currentNode := pl.Head
		for currentNode.Next != nil {
			currentNode = currentNode.Next
		}
		currentNode.Next = thisPair
		thisPair.Prev = currentNode
		pl.Last = thisPair
	}
	return nil
}

func (pl *LinkedListPairs) CountNodes() int {
	count := 0
	currN := pl.Head
	for currN != nil {
		currN = currN.Next
		count += 1
	}
	return count
}

func (pl *LinkedListPairs) Contains(pair *Pair) bool {
	currPair := pl.Head

	for currPair != nil {
		if (currPair.A.Row == pair.A.Row &&
			currPair.A.Col == pair.A.Col &&
			l.IntArrayEquals(currPair.A.Vals, pair.A.Vals)) &&
			(currPair.B.Row == pair.B.Row &&
				currPair.B.Col == pair.B.Col &&
				l.IntArrayEquals(currPair.B.Vals, pair.B.Vals)) {
			return true
		} else if (currPair.A.Row == pair.B.Row &&
			currPair.A.Col == pair.B.Col &&
			l.IntArrayEquals(currPair.A.Vals, pair.B.Vals)) &&
			(currPair.B.Row == pair.A.Row &&
				currPair.B.Col == pair.A.Col &&
				l.IntArrayEquals(currPair.B.Vals, pair.A.Vals)) { // the order is reversed
			return true
		}
		currPair = currPair.Next
	}
	return false
}

func (pl *LinkedListPairs) PrintResult(desc string) {
	currNode := pl.Head
	for currNode != nil {
		color.LightMagenta.Printf(" %s at [%d,%d] & [%d,%d]. val:%v.\n", desc, currNode.A.Row, currNode.A.Col,
			currNode.B.Row, currNode.B.Col, currNode.A.Vals)
		currNode = currNode.Next
	}
}
