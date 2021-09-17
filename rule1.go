package main

import (
	"fmt"
	"time"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
	. "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

func rule1() (*Matchlist, int, time.Duration) {
	var (
		col, row                     int // position of last empty cell
		digit, count                 int
		start                        time.Time
		notInRow, notInCol, notInBlk bool
		matched                      *Matchlist
	)
	start = time.Now()
	matched = &Matchlist{}

	currNode := emptyL.Head
	if currNode == nil {
		color.Yellow.Println("Rule 1: Empty list.")
	} else {
		for currNode != nil {
			row = currNode.Row
			col = currNode.Col
			if *debugPtr {
				color.LightGreen.Printf("cell [%d][%d]. %+v\n", row, col, currNode.Vals)
			}

			if len(currNode.Vals) == 1 {
				digit = currNode.Vals[0]
				matched.AddCell(currNode, digit)
				emptyL.DelNode(currNode) // remove current Node from possibility list
				mat[row][col] = digit
				mat2[row][col] = nil
				emptyCnt--
				count++

				// check that there is no occurrence in same row, col or block
				notInRow = !FindDigitInRow(DebugFn(2), mat2, row, col, digit)
				notInCol = !FindDigitInCol(DebugFn(2), mat2, row, col, digit)
				notInBlk = !FindDigitInBlk(DebugFn(2), mat2, row, col, digit)
				// If found, erase any occurrence of the digit in the same row, col or block
				findAndEraseDigit(row, col, digit, notInRow, notInCol, notInBlk)
			}
			currNode = currNode.Next
		}
	}
	return matched, count, time.Since(start)
}

// Rule 1a	Open singles
//          Search the specified row or column for open singles
func rule1a(row, col int) (*Matchlist, int, time.Duration) {
	var (
		digit, count                 int
		notInRow, notInCol, notInBlk bool
		start                        time.Time
		node                         *Cell
		matched                      *Matchlist
	)
	start = time.Now()
	matched = &Matchlist{}

	if *debugPtr {
		fmt.Println("In rule1a...")
	}

	if row >= 0 && col < 0 { // skip row checking if negative value
		for c := 0; c < N; c++ {
			if len(mat2[row][c]) == 1 {
				digit = mat2[row][c][0]
				node = emptyL.GetNodeFoRCell(row, c)
				matched.AddCell(node, digit)
				emptyL.DelNode(node) // remove current Node from possibility list
				mat[row][c] = digit
				mat2[row][c] = nil
				emptyCnt--
				count++

				// check that there is no occurrence in same row, col or block
				notInRow = !FindDigitInRow(DebugFn(2), mat2, row, c, digit)
				notInCol = !FindDigitInCol(DebugFn(2), mat2, row, c, digit)
				notInBlk = !FindDigitInBlk(DebugFn(2), mat2, row, c, digit)
				// If found, erase any occurrence of the digit in the same row, col or block
				findAndEraseDigit(row, c, digit, notInRow, notInCol, notInBlk)
			}
		}
	}

	if col >= 0 && row < 0 { // skip col checking if negative value
		for r := 0; r < N; r++ {
			if len(mat2[r][col]) == 1 {
				digit = mat2[r][col][0]
				node = emptyL.GetNodeFoRCell(r, col)
				matched.AddCell(node, digit)
				emptyL.DelNode(node) // remove current Node from possibility list
				mat[r][col] = digit
				mat2[r][col] = nil
				emptyCnt--
				count++

				// check that there is no occurrence in same row, col or block
				notInRow = !FindDigitInRow(DebugFn(2), mat2, r, col, digit)
				notInCol = !FindDigitInCol(DebugFn(2), mat2, r, col, digit)
				notInBlk = !FindDigitInBlk(DebugFn(2), mat2, r, col, digit)
				// If found, erase any occurrence of the digit in the same row, col or block
				findAndEraseDigit(r, col, digit, notInRow, notInCol, notInBlk)
			}
		}
	}

	if row >= 0 && col >= 0 { // check only this cell
		if len(mat2[row][col]) == 1 {
			digit = mat2[row][col][0]
			node = emptyL.GetNodeFoRCell(row, col)
			matched.AddCell(node, digit)
			emptyL.DelNode(node) // remove current Node from possibility list
			mat[row][col] = digit
			mat2[row][col] = nil
			emptyCnt--
			count++

			// check that there is no occurrence in same row, col or block
			notInRow = !FindDigitInRow(DebugFn(2), mat2, row, col, digit)
			notInCol = !FindDigitInCol(DebugFn(2), mat2, row, col, digit)
			notInBlk = !FindDigitInBlk(DebugFn(2), mat2, row, col, digit)
			// If found, erase any occurrence of the digit in the same row, col or block
			findAndEraseDigit(row, col, digit, notInRow, notInCol, notInBlk)
		}
	}

	return matched, count, time.Since(start)
}

// *******************************************************************************

// Rule 1 - Open cell - 1 cell empty either in column, row or block.
// Search the possible matrix for any col, row or block that has only 1 empty cell
// Empty cells contain lists. Non-empty cells contain nil.
/*
func rule1() (*Matchlist, int) {
	var (
		col, row                     int // position of last empty cell
		digit, count                 int
		notInRow, notInCol, notInBlk bool
		matched                      *Matchlist
	)
	matched = &Matchlist{}

	currNode := emptyL.Head
	if currNode == nil {
		color.Yellow.Println("Empty list.")
	} else {
		for currNode != nil {
			row = currNode.Row
			col = currNode.Col
			if *debugPtr {
				color.LightGreen.Printf("cell [%d][%d]. %+v\n", row, col, currNode.Vals)
			}

			if len(currNode.Vals) == 1 {
				digit = currNode.Vals[0]
				matched.AddCell(currNode, digit)
				emptyL.DelNode(currNode) // remove current Node from possibility list
				mat[row][col] = digit
				mat2[row][col] = nil
				emptyCnt--
				count++
				if *debugPtr {
					color.LightYellow.Printf("Rule 1. Found lone single no. %d at cell [%d][%d]. Empty cells = %d\n", mat[row][col], row, col, emptyCnt)
				}

				// check that there is no occurrence in same row, col or block
				// If found, remove any occurrences of this digit in same row, column or block.
				// Update mat to include this digit. Remove this node from empty list.
				// Nil the array in cell of possibility matrix and set the cell to nil in mat2.
				notInRow = !findDigitInRow(mat2, row, col, digit)
				notInCol = !findDigitInCol(mat2, row, col, digit)
				notInBlk = !findDigitInBlk(mat2, row, col, digit)

				if *debugPtr {
					if notInRow {
						color.LightBlue.Printf("Digit %d of cell [%d][%d] not in row %d\n", digit, row, col, row)
					}
					if notInCol {
						color.LightBlue.Printf("Digit %d of cell [%d][%d] not in col %d\n", digit, row, col, col)
					}
					if notInBlk {
						color.LightBlue.Printf("Digit %d of cell [%d][%d] not in blk [%d][%d]\n",
							digit, row, col, row/SQ, col/SQ)
					}
				}

				// erase any occurrence of the digit in the same row, col or block
				if !notInRow {
					eraseDigitFromRow(row, col, digit)

					if *debugPtr {
						color.LightBlue.Printf("After deletion from row %d: %v\n", row, mat2[row])
					}
				}
				if !notInCol {
					eraseDigitFromCol(row, col, digit)

					if *debugPtr {
						color.LightBlue.Printf("After deletion from col %d: %v\n", col, getColOfPossibleMat(mat2, col))
					}
				}
				if !notInBlk {
					eraseDigitFromBlk(row, col, digit)

					if *debugPtr {
						color.LightBlue.Printf("After deletion from blk [%d,%d]: %v\n", row/SQ, col/SQ, getBlkOfPossibleMat(mat2, row, col))
					}
				}
				if notInRow && notInCol && notInBlk {
					if *debugPtr {
						color.LightBlue.Println("No deletion necessary.")
					}
				}
			}
			currNode = currNode.Next
		}
	}
	return matched, count
}
*/
