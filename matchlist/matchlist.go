package matchlist

import (
	"log"

	"github.com/gookit/color"
	ll "github.com/mjwong/sudoku2/linkedlist"
)

type Idx struct {
	Row, Col int
	Vals     []int
}

type Result struct {
	Arr  []Idx
	Prev *Result
	Next *Result
}

type Matchlist struct {
	Head        *Result
	Last        *Result
	CurrentCell *Result
}

func (p *Matchlist) AddNode(arrIdx []Idx) error {
	r := &Result{
		Arr: arrIdx,
	}

	if p.Head == nil {
		p.Head = r
		r.Prev = nil
	} else {
		currNode := p.Head
		for currNode.Next != nil {
			currNode = currNode.Next
		}
		currNode.Next = r
		r.Prev = currNode
		p.Last = r
	}
	return nil
}

func (p *Matchlist) CountNodes() int {
	count := 0
	currN := p.Head
	for currN != nil {
		currN = currN.Next
		count += 1
	}
	return count
}

func (p *Matchlist) ContainsPair(arrIdx []Idx) bool {
	cNode := p.Head

	if len(arrIdx) != 2 {
		log.Fatal("New naked pair found does not contain a pair of cell indices")
	}

	for cNode != nil {
		if len(cNode.Arr) != 2 {
			log.Fatal("Matchlist does not contain a pair of cell indices")
		}
		if ((arrIdx[0].Row == cNode.Arr[0].Row && arrIdx[0].Col == cNode.Arr[0].Col) ||
			(arrIdx[0].Row == cNode.Arr[1].Row && arrIdx[0].Col == cNode.Arr[1].Col)) &&
			((arrIdx[1].Row == cNode.Arr[0].Row && arrIdx[1].Col == cNode.Arr[0].Col) ||
				(arrIdx[1].Row == cNode.Arr[1].Row && arrIdx[1].Col == cNode.Arr[1].Col)) {
			return true
		}
		cNode = cNode.Next
	}

	return false
}

func AddArrIdx(arr []Idx, node *ll.Cell) []Idx {
	newIdx := Idx{
		Row:  node.Row,
		Col:  node.Col,
		Vals: node.Vals,
	}

	arr = append(arr, newIdx)
	return arr

}

func (p *Matchlist) PrintResult(desc string) {
	currNode := p.Head
	for currNode != nil {
		color.LightMagenta.Printf("%s: %v at [%d,%d] and [%d,%d].\n", desc, currNode.Arr[0].Vals,
			currNode.Arr[0].Row, currNode.Arr[0].Col, currNode.Arr[1].Row, currNode.Arr[1].Col)
		currNode = currNode.Next
	}
}
