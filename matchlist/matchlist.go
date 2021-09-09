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

type rNode struct {
	Arr  []Idx
	Prev *rNode
	Next *rNode
}

type Matchlist struct {
	Head        *rNode
	Last        *rNode
	CurrentCell *rNode
}

func (p *Matchlist) AddCell(node *ll.Cell, dig int) error {
	r := &rNode{
		Arr: []Idx{
			{
				Row:  node.Row,
				Col:  node.Col,
				Vals: []int{dig},
			},
		},
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

func (p *Matchlist) AddRNode(arrIdx []Idx) error {
	r := &rNode{
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

func (p *Matchlist) PrintResult(desc string) {
	currNode := p.Head
	for currNode != nil {
		color.LightMagenta.Printf("%s: %v at [%d,%d]", desc, currNode.Arr[0].Vals,
			currNode.Arr[0].Row, currNode.Arr[0].Col)
		if len(currNode.Arr) > 1 {
			for i, v := range currNode.Arr {
				if i > 0 {
					color.LightMagenta.Printf(", [%d,%d]", v.Row, v.Col)
				}
			}
			color.LightMagenta.Println()
		} else {
			color.LightMagenta.Println()
		}
		currNode = currNode.Next
	}
}

// accepts a variable no. of cells
func AddIdx(arr []Idx, a ...*ll.Cell) []Idx {
	for _, val := range a {
		newIdx := Idx{
			Row:  val.Row,
			Col:  val.Col,
			Vals: val.Vals,
		}

		if arr == nil {
			arr = []Idx{}
		}
		arr = append(arr, newIdx)
	}

	return arr
}
