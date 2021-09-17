package matchlist

import (
	"fmt"
	"log"

	"github.com/gookit/color"
	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
)

type RCell struct {
	Row, Col int
	Vals     []int
}

type rNode struct {
	Arr  []RCell
	Prev *rNode
	Next *rNode
}

type Matchlist struct {
	Head        *rNode
	Last        *rNode
	CurrentCell *rNode
}

func (p *Matchlist) AddCell(node *Cell, dig int) error {
	r := &rNode{
		Arr: []RCell{
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

func (p *Matchlist) AddRNode(arrRCell []RCell) error {
	r := &rNode{
		Arr: arrRCell,
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

func (p *Matchlist) ContainsPair(arrRCell []RCell) bool {
	cNode := p.Head

	if len(arrRCell) != 2 {
		log.Fatal("New naked pair found does not contain a pair of cell indices")
	}

	for cNode != nil {
		if len(cNode.Arr) != 2 {
			log.Fatal("Matchlist does not contain a pair of cell indices")
		}
		if ((arrRCell[0].Row == cNode.Arr[0].Row && arrRCell[0].Col == cNode.Arr[0].Col) ||
			(arrRCell[0].Row == cNode.Arr[1].Row && arrRCell[0].Col == cNode.Arr[1].Col)) &&
			((arrRCell[1].Row == cNode.Arr[0].Row && arrRCell[1].Col == cNode.Arr[0].Col) ||
				(arrRCell[1].Row == cNode.Arr[1].Row && arrRCell[1].Col == cNode.Arr[1].Col)) {
			return true
		}
		cNode = cNode.Next
	}

	return false
}

func (p *Matchlist) ContainsXwing(arrRCell []RCell) bool {
	cNode := p.Head

	if len(arrRCell) != 4 {
		log.Fatal("Does not contain 4 cells of the X-wing")
	}

	for cNode != nil {
		if RCellArrayEquals(cNode.Arr, arrRCell) {
			return true
		}
		cNode = cNode.Next
	}
	return false
}

func RCellArrayEquals(a []RCell, b []RCell) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range b {
		if ContainsRCell(a, v) {
			return true
		}
	}
	return false
}

func ContainsRCell(a []RCell, b RCell) bool {
	for _, v := range a {
		if v.Row == b.Row && v.Col == b.Col && IntArrayEquals(v.Vals, b.Vals) {
			return true
		}
	}
	return false
}

func (p *Matchlist) PrintResult(desc string) {
	currNode := p.Head
	for currNode != nil {
		color.LightCyan.Printf("%s: %v at [%d,%d]", desc, currNode.Arr[0].Vals,
			currNode.Arr[0].Row, currNode.Arr[0].Col)
		if len(currNode.Arr) > 1 {
			for i, v := range currNode.Arr {
				if i > 0 {
					color.LightCyan.Printf(", [%d,%d]", v.Row, v.Col)
				}
			}
			fmt.Println()
		} else {
			fmt.Println()
		}
		currNode = currNode.Next
	}
}

// accepts a variable no. of cells
func AddRCell(arr []RCell, a ...*Cell) []RCell {
	for _, val := range a {
		newRCell := RCell{
			Row:  val.Row,
			Col:  val.Col,
			Vals: val.Vals,
		}

		if arr == nil {
			arr = []RCell{}
		}
		arr = append(arr, newRCell)
	}

	return arr
}

func AddRCellToArr(arr []RCell, row, col, dig int) []RCell {
	newRCell := RCell{
		Row:  row,
		Col:  col,
		Vals: []int{dig},
	}

	if arr == nil {
		arr = []RCell{}
	}
	arr = append(arr, newRCell)

	return arr
}

func AppendMatchlist(matched, found *Matchlist) *Matchlist {

	rnode := found.Head
	for rnode != nil {
		matched.AddRNode(rnode.Arr)
		rnode = rnode.Next
	}

	return matched
}
