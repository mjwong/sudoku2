package linkedlist

import (
	. "github.com/mjwong/sudoku2/lib"
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

func (p *LinkedList) AddNode(cells ...*Cell) error {
	var currentCell *Cell

	if p.Head == nil {
		p.Head = cells[0]
		cells[0].Prev = nil
		currentCell = p.Head
	} else {
		currentCell = p.Head
		for currentCell.Next != nil {
			currentCell = currentCell.Next
		}
		currentCell.Next = cells[0]
		cells[0].Prev = currentCell
		currentCell = cells[0]
	}

	if len(cells) > 1 {
		for i, v := range cells {
			if i > 0 {
				currentCell.Next = v
				v.Prev = currentCell
				currentCell = currentCell.Next
				p.Last = v
			}
		}
	} else {
		p.Last = cells[0]
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

func (p *LinkedList) CountElem() int {
	count := 0
	currN := p.Head
	for currN != nil {
		count += len(currN.Vals)
		currN = currN.Next
	}
	return count
}

func (p *LinkedList) GetNodeFoRCell(row, col int) *Cell {
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
	node := p.GetNodeFoRCell(row, col)
	node.Vals = EraseFromSlice(node.Vals, dig)
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

func GetPossibleMat(mat Intmat) (*LinkedList, Pmat) {
	const (
		ncols = 9
	)
	var (
		mat2                Pmat
		emptyL              *LinkedList
		inRow, inCol, inSqu bool
		valList             []int
	)

	emptyL = CreatelinkedList()

	// initialize possible mat
	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
			mat2[i][j] = nil
		}
	}

	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
			valList = nil

			if mat[i][j] == 0 {
				for p := 1; p < ncols+1; p++ {
					inRow = Contains(GetArrForRow(mat, i), p)
					inCol = Contains(GetArrForCol(mat, j), p)
					inSqu = Contains(GetArrForSqu(mat, i, j), p)

					if !inRow && !inCol && !inSqu {
						mat2[i][j] = append(mat2[i][j], p)
						valList = append(valList, p)
					}
				}
				emptyL.AddCell(i, j, valList)
			}
		}
	}
	return emptyL, mat2
}
