package main

import (
	"fmt"
	"time"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
	. "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

// Rule 5	Naked pairs
//          A pair of digits that occurs in exactly 2 cells in an entire row, column, or block.
//          Erase any other occurrence of these 2 digits elsewhere in the same row, column or block.
func rule5() (*Matchlist, int, time.Duration) {
	var (
		col, row, col2, row2   int // position of last empty cell
		count                  int
		twoElem                []int
		foundNakedPairs, debug bool
		inRow, inCol, inBlk    bool
		start                  time.Time
		secondNode             *Cell
		matched                *Matchlist
	)

	start = time.Now()
	matched = &Matchlist{}
	debug = DebugFn(2)
	foundNakedPairs = true

	for foundNakedPairs {
		foundNakedPairs = false

		currNode := emptyL.Head

		if currNode == nil {
			color.Yellow.Println("Rule 5: Empty list.")
			break
		} else {
			for currNode != nil {
				row = currNode.Row
				col = currNode.Col
				if debug {
					color.LightGreen.Printf("cell [%d][%d]. %+v\n", row, col, currNode.Vals)
				}

				if len(currNode.Vals) == 2 { // has 2 possible values
					twoElem = currNode.Vals
					// check row
					for c := 0; c < N; c++ {
						if IntArrayEquals(mat2[row][c], twoElem) && c != col {
							col2 = c

							if debug {
								PrintPossibleMat(mat2)
								color.Magenta.Printf("Found naked pair in row %d, in cols %d and %d.\n", row, col, col2)
							}

							secondNode = emptyL.GetNodeFoRCell(row, col2)

							arr := AddRCell(nil, currNode, secondNode)

							if !matched.ContainsPair(arr) {
								if debug {
									fmt.Printf("Row %d\n", row)
								}

								matched.AddRNode(arr)
								inRow = FindDigitInRowPair(debug, mat2, row, col, col2, twoElem)
								if inRow {
									if debug {
										fmt.Printf("Found digits of pairs in row %d.\n", row)
									}
									eraseDigitsFromRowOfPairs(row, col, col2, twoElem)
								}
								foundNakedPairs = true
								count++

								if debug {
									fmt.Printf("Naked pairs = %d.\n", count)
								}
								break
							}
						}
					}

					// check col
					for r := 0; r < N; r++ {
						if IntArrayEquals(mat2[r][col], twoElem) && r != row {
							row2 = r

							if debug {
								color.Magenta.Printf("Found naked pair in col %d, in rows %d and %d.\n", col, row, row2)
							}

							secondNode = emptyL.GetNodeFoRCell(row2, col)

							arr := AddRCell(nil, currNode, secondNode)

							if !matched.ContainsPair(arr) {
								if debug {
									fmt.Printf("Col %d\n", col)
								}

								matched.AddRNode(arr)
								inCol = FindDigitInColPair(debug, mat2, row, col, row2, twoElem)
								if inCol {
									if debug {
										fmt.Printf("Found digits of pairs in col %d.\n", col)
									}
									eraseDigitsFromColOfPairs(row, col, row2, twoElem)
								}
								foundNakedPairs = true
								count++

								if debug {
									fmt.Printf("Naked pairs = %d.\n", count)
								}
								break
							}
						}
					}

					// check blk
					startRow := row / SQ * SQ
					startCol := col / SQ * SQ
					emptyCntBlk := emptyL.CountNodes()

					if debug {
						fmt.Printf("Finding 2nd pair [%d,%d] cell [%d,%d]\n", twoElem[0], twoElem[1], row, col)
					}

					for x := startRow; x < startRow+SQ; x++ {
						for y := startCol; y < startCol+SQ; y++ {
							if debug {
								fmt.Printf("Blk [%d,%d]: cell [%d,%d]\n", row/SQ, col/SQ, x, y)
							}

							if IntArrayEquals(mat2[x][y], twoElem) && !(x == row && y == col) {
								row2 = x
								col2 = y

								if debug {
									color.Magenta.Printf("Found naked pair (%d,%d) in blk [%d,%d], in cells [%d,%d] and [%d,%d].\n",
										twoElem[0], twoElem[1], row/SQ, col/SQ, row, col, row2, col2)
								}

								secondNode = emptyL.GetNodeFoRCell(row2, col2)
								arr := AddRCell(nil, currNode, secondNode)

								if debug {
									matched.PrintResult("Naked pairs")
								}

								if !matched.ContainsPair(arr) {
									matched.AddRNode(arr)

									if debug {
										fmt.Printf("Blk [%d,%d]\n", row/SQ, col/SQ)
									}

									inBlk = FindDigitInBlkPair(debug, mat2, row, col, row2, col2, twoElem)
									if inBlk {
										if debug {
											fmt.Printf("Found digits of pairs in blk [%d,%d].\n", row/SQ, col/SQ)
											PrintPossibleMat(mat2)
										}

										eraseDigitsFromBlkOfPairs(row, col, row2, col2, twoElem)
									}
									foundNakedPairs = true
									count++

									if debug {
										fmt.Printf("Naked pairs = %d.\n", count)
									}
									break
								}

							}
						}
					}

					if emptyL.CountNodes() < emptyCntBlk {
						color.LightMagenta.Printf("Deleted cells after checking block: %d\n", emptyCntBlk-emptyL.CountNodes())
					}
				}

				currNode = currNode.Next
			}
		}
	}

	return matched, count, time.Since(start)
}
