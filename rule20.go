package main

import (
	"fmt"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

// Rule 20: X-wing (or Rectangular pattern)
// Rectangular box  pattern. If the same no. appears in the corner cells of a rectangular box,
// then that no. can be safely eliminated (crossed out) in all columns and rows that intersect
// with the corner cells of the rectangular box.
func rule20() (*Matchlist, int) {
	var (
		count      int
		foundXWing bool
		inBlk      bool
		debug      bool
		arrC       []Coord
		matched    *Matchlist
		foundList  *Matchlist
	)

	matched = &Matchlist{}
	foundList = &Matchlist{}
	foundXWing = true
	debug = DebugFn(2)

	fmt.Printf("Func: %s. Debug: %t\n", FuncName(1), debug)

	for dig := 1; dig <= N; dig++ {
		foundXWing = false
		currNode := emptyL.Head

		if currNode == nil {
			color.Yellow.Println("Rule 20: Empty list.")
			break
		} else {
			for currNode != nil {
				// Check block contains only 2 possible digit in exactly 2 places
				// This digit may be hidden in the list of possibile digits.
				for bi := 0; bi < SQ; bi++ {

					for bj := 0; bj < SQ; bj++ {
						arrC, inBlk = checkBlkForDigit(mat2, bi, bj, dig, 2)
						if inBlk { // exactly 2 same digits in this block

							// Are they in the same row?
							foundList, foundXWing = checkSameRowXwing(debug, bi, bj, dig, arrC, inBlk, foundList)
							if foundXWing {
								count++
							}

							// Are they in the same column?
							foundList, foundXWing = checkSameColXwing(debug, bi, bj, dig, arrC, inBlk, foundList)
							if foundXWing {
								count++
							}
						}
					}
				}
				currNode = currNode.Next
			} // end of emptyL iteration
		}
	}
	return matched, count
}

func checkSameRowXwing(debug bool, bi, bj, dig int, arrC []Coord, inBlk bool, foundList *Matchlist) (*Matchlist, bool) {
	blkListi := []int{}
	foundXWing := false

	if arrC[0].Row == arrC[1].Row {
		if debug {
			color.LightBlue.Printf("The 2 same digits no. %d are both in row %d.\n",
				dig, arrC[0].Row)
		}
		// check blks vertically, i.e. this col of blocks
		for bi2 := 0; bi2 < SQ; bi2++ {
			if bi != bi2 { // not the original block
				arrC2, inBlk2 := checkBlkForDigit(mat2, bi2, bj, dig, 2)
				if inBlk2 { // found exactly 2 same digits in second block
					if debug {
						color.LightBlue.Printf("Found 2nd block [%d,%d].\n", bi2, bj)
					}

					// Locations of the X-wing cells in first block
					rowXw1 := arrC[0].Row
					colXw1 := arrC[0].Col
					rowXw2 := arrC[1].Row
					colXw2 := arrC[1].Col
					// For second block
					rowXw3 := arrC2[0].Row
					colXw3 := arrC2[0].Col
					rowXw4 := arrC2[1].Row
					colXw4 := arrC2[1].Col

					// check the 4 corner cells. Do they form a rectangle?
					if (colXw1 == colXw3 || colXw1 == colXw4) &&
						(colXw2 == colXw3 || colXw2 == colXw4) {

						arr := []Idx{}
						arr = append(arr, Idx{Row: rowXw1, Col: colXw1, Vals: []int{dig}})
						arr = append(arr, Idx{Row: rowXw2, Col: colXw2, Vals: []int{dig}})
						arr = append(arr, Idx{Row: rowXw3, Col: colXw3, Vals: []int{dig}})
						arr = append(arr, Idx{Row: rowXw4, Col: colXw4, Vals: []int{dig}})
						if !foundList.ContainsXwing(arr) {
							foundList.AddRNode(arr)
						}

						if debug {
							color.LightMagenta.Printf("Found X-wing #%d: [%d,%d], [%d,%d], [%d,%d], [%d,%d].\n",
								dig, rowXw1, colXw1, rowXw2, colXw2, rowXw3, colXw3, rowXw4, colXw4)
							PrintPossibleMat(mat2)
						}

						blkListi = append(blkListi, bi2)
						colList := []int{}
						colList = append(colList, colXw1)
						colList = append(colList, colXw2)
						if debug {
							color.LightYellow.Printf("Collist: %v\n", colList)
						}

						// find and delete any occurrences of this digit elsewhere in the cols of colList
						// excluding the row positions of the 2 cells forming the X-wing.
						for _, c := range colList {
							for r := 0; r < N; r++ {
								if r != rowXw1 && r != rowXw2 && r != rowXw3 && r != rowXw4 {
									// erase digit from this cell
									eraCnt, erased := eraseDigitFromColMulti(c, dig, []int{rowXw1, rowXw2, rowXw3, rowXw4})
									if eraCnt > 0 {
										foundXWing = true
									}

									if debug {
										if erased {
											color.LightMagenta.Printf("X-wing: Erased %d counts of digit %d from col %d.\n", eraCnt, dig, c)
											PrintPossibleMat(mat2)
										}
									}
								}
							}
						}

						thirdCol := 0 // find the 3rd column
						for c := 0; c < SQ; c++ {
							if !Contains(colList, c) {
								thirdCol = c

								if debug {
									color.LightYellow.Printf("Found third col no. %d of digit %d.\n", thirdCol, dig)
								}
								break
							}
						}

						if debug {
							color.LightYellow.Printf("Row blklist: %v\n", blkListi)
						}

						for i := 0; i < SQ; i++ { // find third block
							if !Contains(blkListi, i) {

								if debug {
									color.LightYellow.Printf("Third blk: %d\n", i)
								}

								// check whether there is open or hidden single in the
								// third missing row of this third block

								startRow := i * SQ
								for r := startRow; r < startRow+SQ; r++ {
									if len(mat2[r][thirdCol]) == 1 {
										// insurance check for entire row. Rightfully, we can just check
										// the entire row since we know the digit can only appear in the
										// third missing row of the third block, because the first 2 blocks
										// already contain a pair of the digits each, forming the X-wing.
										matched, cnt := rule1a(r, thirdCol)

										if debug && cnt > 0 {
											color.LightYellow.Printf("Rule1a: Found %d counts of open single %d at [%d,%d]\n",
												cnt, dig, r, thirdCol)
											matched.PrintResult(RuleTable[20])
										}
									}

									// check for hidden singles at this Cell position
									if mat2[r][thirdCol] != nil {
										matched3, cnt3 := rule3a(r, thirdCol, dig)
										if debug && cnt3 > 0 {
											color.LightYellow.Printf("Rule3a: Found %d counts of hidden single %d at [%d,%d]\n",
												cnt3, dig, r, thirdCol)
											PrintPossibleMat(mat2)
											matched3.PrintResult(RuleTable[20])
										}
									}
								}
							}
						} // end of find third block
					}
				}
			}
		}
	} // end of same rows
	return foundList, foundXWing
}

func checkSameColXwing(debug bool, bi, bj, dig int, arrC []Coord, inBlk bool, foundList *Matchlist) (*Matchlist, bool) {
	blkListj := []int{} // col blocks
	foundXWing := false

	if arrC[0].Col == arrC[1].Col {
		blkListj = append(blkListj, bj)
		if debug {
			color.LightBlue.Printf("The 2 same digits no. %d are both in col %d.\n",
				dig, arrC[0].Col)
		}
		// check blks horizontally, i.e. this row of blocks
		for bj2 := 0; bj2 < SQ; bj2++ {
			if bj != bj2 { // not the original block
				arrC2, inBlk2 := checkBlkForDigit(mat2, bi, bj2, dig, 2)
				if inBlk2 { // found exactly 2 same digits in second block
					if debug {
						color.LightBlue.Printf("Found 2nd block [%d,%d].\n", bi, bj2)
					}

					// Locations of the X-wing cells in first block
					rowXw1 := arrC[0].Row
					colXw1 := arrC[0].Col
					rowXw2 := arrC[1].Row
					colXw2 := arrC[1].Col
					// For second block
					rowXw3 := arrC2[0].Row
					colXw3 := arrC2[0].Col
					rowXw4 := arrC2[1].Row
					colXw4 := arrC2[1].Col

					// check the 4 corner cells. Do they form a rectangle?
					if (rowXw1 == rowXw3 || rowXw1 == rowXw4) &&
						(rowXw2 == rowXw3 || rowXw2 == rowXw4) {

						arr := []Idx{}
						arr = append(arr, Idx{Row: rowXw1, Col: colXw1, Vals: []int{dig}})
						arr = append(arr, Idx{Row: rowXw2, Col: colXw2, Vals: []int{dig}})
						arr = append(arr, Idx{Row: rowXw3, Col: colXw3, Vals: []int{dig}})
						arr = append(arr, Idx{Row: rowXw4, Col: colXw4, Vals: []int{dig}})
						if !foundList.ContainsXwing(arr) {
							foundList.AddRNode(arr)
						}

						if inBlk2 { // found exactly 2 same digits in second block
							if debug {
								color.LightMagenta.Printf("Found X-wing #%d: [%d,%d], [%d,%d], [%d,%d], [%d,%d].\n",
									dig, rowXw1, colXw1, rowXw2, colXw2, rowXw3, colXw3, rowXw4, colXw4)
								PrintPossibleMat(mat2)
							}
							blkListj = append(blkListj, bj2)
							rowList := []int{}
							rowList = append(rowList, rowXw1)
							rowList = append(rowList, rowXw2)

							if debug {
								color.LightYellow.Printf("Rowlist: %v\n", rowList)
							}

							// find and delete any occurrences of this digit elsewhere in the rows of rowList
							// excluding the column positions of the 2 cells forming the X-wing.
							for _, r := range rowList {
								for c := 0; c < N; c++ {
									if c != colXw1 && c != colXw2 && c != colXw3 && c != colXw4 {
										// erase digit from this cell
										eraCnt, erased := eraseDigitsFromRowMulti(r, []int{dig}, []int{colXw1, colXw2, colXw3, colXw4})
										if eraCnt > 0 {
											foundXWing = true
										}

										if debug {
											if erased {
												color.LightYellow.Printf("X-wing: Erased %d counts of digit %d from row %d.\n", eraCnt, dig, r)
												PrintPossibleMat(mat2)
											}
										}
									}
								}
							}

							thirdRow := 0 // find the 3rd row
							for r := 0; r < SQ; r++ {
								if !Contains(rowList, r) {
									thirdRow = r

									if debug {
										color.LightYellow.Printf("Found third row no. %d of digit %d.\n", thirdRow, dig)
									}
									break
								}
							}

							if debug {
								color.LightYellow.Printf("Col blklist: %v\n", blkListj)
							}

							for j := 0; j < SQ; j++ { // find third block horizontally
								if !Contains(blkListj, j) {

									if debug {
										color.LightYellow.Printf("Third blk: %d\n", j)
									}

									// check whether there is open or hidden single in the
									// third missing row of this third block

									startCol := j * SQ
									for c := startCol; c < startCol+SQ; c++ {
										if len(mat2[thirdRow][c]) == 1 {
											// insurance check for entire row. Rightfully, we can just check
											// the entire row since we know the digit can only appear in the
											// third missing row of the third block, because the first 2 blocks
											// already contain a pair of the digits each, forming the X-wing.
											matched, cnt := rule1a(thirdRow, c)

											if debug && cnt > 0 {
												color.LightYellow.Printf("Rule1a: Found %d counts of open single %d at [%d,%d]\n",
													cnt, dig, thirdRow, c)
												matched.PrintResult(RuleTable[20])
											}
										}

										// check for hidden singles at this Cell position
										if mat2[thirdRow][c] != nil {
											matched3, cnt3 := rule3a(thirdRow, c, dig)

											if debug && cnt3 > 0 {
												color.LightYellow.Printf("Rule3a: Found %d counts of hidden single %d at [%d,%d]\n",
													cnt3, dig, thirdRow, c)
												PrintPossibleMat(mat2)
												matched3.PrintResult(RuleTable[20])
											}
										}
									}
								}
							} // end of find third block horizontally
						}
					}
				}
			}
		}
	} // end of same col
	return foundList, foundXWing
}

func checkBlkForDigit(m Pmat, bx, by, dig, occurence int) ([]Coord, bool) {
	var (
		count, startRow, startCol int
	)

	startRow = bx * SQ
	startCol = by * SQ
	arr := []Coord{}

	for i := startRow; i < startRow+SQ; i++ {
		for j := startCol; j < startCol+SQ; j++ {
			if Contains(m[i][j], dig) {
				arr = append(arr, Coord{Row: i, Col: j})
				count++
			}
		}
	}

	if count == occurence {
		return arr, true
	}
	return arr, false
}

func checkRowForDigit(m Pmat, row, dig, occurence int) ([]Coord, bool) {
	var count int
	arr := []Coord{}

	for c := 0; c < N; c++ {
		if Contains(m[row][c], dig) {
			arr = append(arr, Coord{Row: row, Col: c})
			count++
		}
	}

	if count == occurence {
		return arr, true
	}
	return arr, false
}

func checkColForDigit(m Pmat, col, dig, occurence int) ([]Coord, bool) {
	var count int
	arr := []Coord{}

	for r := 0; r < N; r++ {
		if Contains(m[r][col], dig) {
			arr = append(arr, Coord{Row: r, Col: col})
			count++
		}
	}

	if count == occurence {
		return arr, true
	}
	return arr, false
}

func checkDigitInColOfBlk(m Pmat, bx, col, dig, occurence int) ([]Coord, bool) {
	var (
		count, startRow int
	)
	startRow = bx * SQ
	arr := []Coord{}

	for i := startRow; i < startRow+SQ; i++ {
		if Contains(m[i][col], dig) {
			arr = append(arr, Coord{Row: i, Col: col})
			count++
		}
	}

	if count == occurence {
		return arr, true
	}
	return arr, false
}
