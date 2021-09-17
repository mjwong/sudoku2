package main

import (
	"fmt"
	"time"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
	. "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

// Rule 3	Hidden singles
//          A digit that is the only one in an entire row, column, or block.
//          Fill in this digiti and erase any other occurrence of this digit in the same row, column or block.
func rule3() (*Matchlist, int, time.Duration) {
	var (
		count, itercnt           int
		foundHiddenSingle, debug bool
		start                    time.Time
		matched                  *Matchlist
	)
	start = time.Now()
	debug = DebugFn(2)
	matched = &Matchlist{}

	for dig := 1; dig <= N; dig++ {
		itercnt = 0
		for {
			foundHiddenSingle = false
			currNode := emptyL.Head

			if currNode == nil {
				color.Yellow.Println("Rule 3: Empty list.")
				break
			} else {
				for currNode != nil {
					if debug {
						color.LightGreen.Printf("cell [%d][%d]. %+v\n", currNode.Row, currNode.Col, currNode.Vals)
					}

					if Contains(currNode.Vals, dig) {
						foundHiddenSingle = findDigitAndUpdate(currNode, dig)
						if foundHiddenSingle {
							matched.AddCell(currNode, dig)
							count++
						}
					}

					if emptyCnt <= 0 {
						if debug {
							fmt.Printf("Empty list count = %d\n", emptyL.CountNodes())
						}
						break
					}

					currNode = currNode.Next
				}
			}

			if !foundHiddenSingle { // repeat same digit until no hidden single is found
				break
			} else if itercnt > 10 {
				break
			}
		}

		itercnt++
		if emptyL.CountNodes() == 0 {
			break
		}
	}

	return matched, count, time.Since(start)
}

// Rule 3a	Hidden singles
//			Search in the specified row or col or blk intersecting this Cell only
func rule3a(row, col, dig int) (*Matchlist, int, time.Duration) {
	var (
		count                    int
		foundHiddenSingle, debug bool
		start                    time.Time
		currNode                 *Cell
		matched                  *Matchlist
	)
	start = time.Now()
	debug = DebugFn(1)
	matched = &Matchlist{}
	currNode = emptyL.GetNodeFoRCell(row, col)

	if Contains(currNode.Vals, dig) {
		foundHiddenSingle = findDigitAndUpdate(currNode, dig)
		if foundHiddenSingle {
			matched.AddCell(currNode, dig)
			count++
		}
	}

	if emptyCnt <= 0 {
		if debug {
			fmt.Printf("Empty list count = %d\n", emptyL.CountNodes())
		}
	}

	return matched, count, time.Since(start)
}

func findDigitAndUpdate(currNode *Cell, dig int) bool {
	var (
		row, col                     int
		notInRow, notInCol, notInBlk bool
		found, debug                 bool
	)
	debug = DebugFn(3)
	row = currNode.Row
	col = currNode.Col

	// check that there is no occurrence in same row, col or block
	notInRow = !FindDigitInRow(debug, mat2, row, col, dig)
	notInCol = !FindDigitInCol(debug, mat2, row, col, dig)
	notInBlk = !FindDigitInBlk(debug, mat2, row, col, dig)

	if notInRow || notInCol || notInBlk {
		found = true
		emptyL.DelNode(currNode) // remove current Node from possibility list
		mat[row][col] = dig      // fill in dig in resulting mat
		mat2[row][col] = nil     // blank this cell out in possibility mat
		emptyCnt--

		// erase any occurrence of the digit in the same row, col or block
		findAndEraseDigit(row, col, dig, notInRow, notInCol, notInBlk)
	}
	return found
}

func findAndEraseDigit(row, col, dig int, notInRow, notInCol, notInBlk bool) {
	if *debugPtr {
		if notInRow {
			color.LightBlue.Printf("Digit %d of cell [%d][%d] not in row %d\n", dig, row, col, row)
		}
		if notInCol {
			color.LightBlue.Printf("Digit %d of cell [%d][%d] not in col %d\n", dig, row, col, col)
		}
		if notInBlk {
			color.LightBlue.Printf("Digit %d of cell [%d][%d] not in blk [%d][%d]\n",
				dig, row, col, row/SQ, col/SQ)
		}
	}

	if !notInRow {
		eraseDigitFromRow(row, col, dig)

		if *debugPtr {
			color.LightBlue.Printf("After deletion from row %d: %v\n", row, mat2[row])
		}
	}
	if !notInCol {
		eraseDigitFromCol(row, col, dig)

		if *debugPtr {
			color.LightBlue.Printf("After deletion from col %d: %v\n", col, GetColOfPossibleMat(mat2, col))
		}
	}
	if !notInBlk {
		eraseDigitFromBlk(row, col, dig)

		if *debugPtr {
			color.LightBlue.Printf("After deletion from blk [%d,%d]: %v\n", row/SQ, col/SQ, GetBlkOfPossibleMat(mat2, row, col))
		}
	}
	if notInRow && notInCol && notInBlk {
		if *debugPtr {
			color.LightBlue.Println("No deletion necessary.")
		}
	}
}
