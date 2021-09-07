package linkedlist

import (
	l "github.com/mjwong/sudoku2/lib"
	"gopkg.in/gookit/color.v1"
)

type Cell struct {
	Row  int
	Col  int
	Vals []int
	Prev *Cell
	Next *Cell
}

type LinkedList struct {
	Head        *Cell
	Last        *Cell
	CurrentCell *Cell
}

func CreatelinkedList() *LinkedList {
	return &LinkedList{}
}

func (p *LinkedList) AddCell(row, col int, arr []int) error {
	c := &Cell{
		Row:  row,
		Col:  col,
		Vals: arr,
	}

	if p.Head == nil {
		p.Head = c
		c.Prev = nil
	} else {
		currentCell := p.Head
		for currentCell.Next != nil {
			currentCell = currentCell.Next
		}
		currentCell.Next = c
		c.Prev = currentCell
		p.Last = c
	}
	return nil
}

func (p *LinkedList) ShowAllEmptyCells() error {
	currentCell := p.Head
	if currentCell == nil {
		color.Red.Println("EmptyCell list is empty.")
		return nil
	}
	color.Green.Printf("%+v\n", *currentCell)
	for currentCell.Next != nil {
		currentCell = currentCell.Next
		color.Green.Printf("%+v\n", *currentCell)
	}
	return nil
}

func (p *LinkedList) CountNodes() int {
	count := 0
	currN := p.Head
	for currN != nil {
		currN = currN.Next
		count += 1
	}
	return count
}

func (p *LinkedList) GetNodeForCell(row, col int) *Cell {
	currN := p.Head
	for currN != nil {
		if currN.Row == row && currN.Col == col {
			return currN
		}
		currN = currN.Next
	}
	return nil
}

func (p *LinkedList) EraseDigitFromCell(row, col, dig int) {
	node := p.GetNodeForCell(row, col)
	node.Vals = l.EraseFromSlice(node.Vals, dig)
}

func (p *LinkedList) InsNode(node *Cell) {

}

// remove current node from linked list and connect prev and next nodes
func (p *LinkedList) DelNode(node *Cell) {

	if node == nil {
		color.Yellow.Println("Possibility list is empty or has reached the end.")
		return
	}

	if node == p.Head {
		p.Head = node.Next
	} else {
		if node != nil {
			if node.Next != nil {
				node.Prev.Next = node.Next
				node.Next.Prev = node.Prev
			} else {
				node.Prev.Next = nil
			}
		}
	}
}

func (p *LinkedList) PrintResult(desc string) {
	currNode := p.Head
	for currNode != nil {
		color.LightMagenta.Printf("%s %d at [%d][%d].\n", desc, currNode.Vals, currNode.Row, currNode.Col)
		currNode = currNode.Next
	}
}
