package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
	. "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

type (
	fnRule func() (*Matchlist, int)
)

const (
	Zero    int = iota // Zero = 0
	SameCnt            // SameCnt = 1
)

var (
	iterCnt  int
	emptyCnt int
	mat      Intmat
	mat2     Pmat   // matrix with possible values in empty cells
	mat3     Intmat // guessed matrix
	emptyL   *LinkedList
	debugPtr *bool   = flag.Bool("debug", false, "verbose debug mode")
	prtLLPtr *bool   = flag.Bool("prtLL", false, "print the linked list of empty cells")
	verbose  *bool   = flag.Bool("v", false, "Print if the digit(s) are found")
	rule     *int    = flag.Int("r", 0, "The deffault is 0, which will iterate matrix using linked list.")
	fnName   *string = flag.String("f", "", "Debug the specified function.")

	RuleTable = map[int]string{
		1:  "Open cell",
		2:  "Lone singles",
		3:  "Hidden singles",
		4:  "Omission",
		5:  "Naked pairs",
		6:  "Naked triplets",
		7:  "Naked quads",
		8:  "Hidden pairs",
		9:  "Hidden triplets",
		10: "Hidden quads",
		20: "X-wings",
	}
)

func main() {
	var (
		start   time.Time
		elapsed time.Duration
	)
	flag.Parse()
	fmt.Printf("Debug func: %v\n", *fnName)

	mat = PopulateMat(ReadInput())
	emptyCnt = CountEmpty(mat)
	fmt.Printf("Empty cells: %d\n", emptyCnt)
	start = time.Now()
	PrintSudoku(mat)
	emptyL, mat2 = GetPossibleMat(mat)
	fmt.Println("Starting possibility matrix.")
	PrintPossibleMat(mat2)

	if *prtLLPtr {
		emptyL.ShowAllEmptyCells()
	}

	switch *rule {
	case 0:
		fmt.Println("Default to iterMat.")
		mat3 = mat
		iterMat(emptyL.Head)
		PrintSudoku(mat3)
		CheckSums(mat3)
	case 1:
		RuleLoop(rule1, RuleTable[1], Zero)
	case 3:
		RuleLoop(rule3, RuleTable[3], Zero)
	case 5:
		RuleLoop(rule5, RuleTable[5], SameCnt)
	case 13:
		RuleLoop(rule1, RuleTable[1], Zero)
		RuleLoop(rule3, RuleTable[3], Zero)
	case 135:
		ruleCnt := map[int]int{}

		for {
			cnt1 := RuleLoop(rule1, RuleTable[1], Zero)
			cnt3 := RuleLoop(rule3, RuleTable[3], Zero)
			cntBefore := emptyL.CountElem()
			cnt5 := RuleLoop(rule5, RuleTable[5], SameCnt)
			cntAfter := emptyL.CountElem()
			cnt1a := 0
			if cnt3 > 0 {
				cnt1a = RuleLoop(rule1, RuleTable[1], Zero)
			}
			ruleCnt[1] += cnt1 + cnt1a
			ruleCnt[3] += cnt3
			ruleCnt[5] += cnt5

			if cnt1 == 0 && cnt3 == 0 && cntBefore == cntAfter && cnt1a == 0 || emptyL.CountNodes() == 0 {
				break
			}
		}

		PrintFound([]int{1, 3, 5}, ruleCnt)
		fmt.Printf("Empty cells : %2d\n", emptyL.CountNodes())
		fmt.Printf("Rule 1: %d\n", RuleLoop(rule1, RuleTable[1], Zero))

		if emptyL.CountNodes() == 0 {
			CheckSums(mat)
		}
	case 20:
		matched20, cnt20 := rule20()
		matched20.PrintResult(RuleTable[20])
		fmt.Printf("Rule 20: Found %2d %ss\n", cnt20, RuleTable[20])
	case 99: // run everything including iterMat
		ruleCnt := map[int]int{}
		loop := 0

		for {
			matched1, cnt1 := rule1()
			fmt.Printf("After rule1, found %2d. Empty list count = %2d.\n", cnt1, emptyL.CountNodes())
			matched1.PrintResult("Found open single")

			matched3, cnt3 := rule3()
			fmt.Printf("After rule3, found %2d. Empty list count = %2d.\n", cnt3, emptyL.CountNodes())
			matched3.PrintResult("Found hidden single")

			cntBefore5 := emptyL.CountElem()
			matched5, cnt5 := rule5()
			cntAfter5 := emptyL.CountElem()

			cntBefore20 := emptyL.CountElem()
			matched20, cnt20 := rule20()
			cntAfter20 := emptyL.CountElem()

			fmt.Printf("After rule20, found %2d. Empty list count = %2d.\n", cnt5, emptyL.CountNodes())

			if cnt1 <= 0 && cnt3 <= 0 && cntBefore5 == cntAfter5 && cntBefore20 == cntAfter20 {
				if loop == 2 {
					break
				}
				loop++
			}
			ruleCnt[1] += cnt1
			ruleCnt[3] += cnt3
			if cntBefore5 != cntAfter5 {
				ruleCnt[5] += cnt5
				matched5.PrintResult("Found hidden single")
			}
			if cntBefore20 != cntAfter20 {
				ruleCnt[20] += cnt20
				matched20.PrintResult("Found X-wing")
			}
		}

		ecnt := emptyL.CountNodes()
		fmt.Printf("After rules 1, 3, 5 and 20 have completed. Empty count : %d\n", ecnt)
		PrintSudoku(mat)

		if emptyL.CountNodes() > 0 {
			PrintPossibleMat(mat2)
			// do iterations
			fmt.Printf("Empty list count before running iterMat = %d.\n", emptyL.CountNodes())
			PrintPossibleMat(mat2)
			mat3 = mat
			iterMat(emptyL.Head)
			PrintSudoku(mat3)
		} else {
			CheckSums(mat)
		}

		PrintFound([]int{1, 3, 5, 20}, ruleCnt)
		fmt.Printf("Empty cells : %2d\n", emptyL.CountNodes())
	}

	elapsed = time.Since(start)
	log.Printf("IterMat: Iterations: %d. Empty cells: %d. Sudoku took %v sec\n", iterCnt, CountEmpty(mat), elapsed.Seconds())
}

func RuleLoop(rule fnRule, desc string, exitCond int) int {
	exitFor := false
	totalcnt := 0
	fnName := GetFunctionName(rule)
	for {
		cntBefore := emptyL.CountNodes()
		matched, cnt := rule()
		cntAfter := emptyL.CountNodes()
		fmt.Printf("%s: Found %d digits.\n", fnName, cnt)
		matched.PrintResult(desc)
		PrintSudoku(mat)
		switch exitCond {
		case Zero:
			if cnt <= 0 {
				exitFor = true
			}
		case SameCnt:
			if cntBefore == cntAfter {
				exitFor = true
			}
		default:
			if cnt <= 0 {
				exitFor = true
			}
		}

		totalcnt += cnt
		if exitFor {
			break
		}
	}
	fmt.Printf("%s: Total found = %d.\n", fnName, totalcnt)
	PrintPossibleMat(mat2)
	return totalcnt
}

func PrintFound(ruleList []int, ruleCounts map[int]int) {
	for _, v := range ruleList {
		fmt.Printf("Rule %2d: Found %2d %s\n", v, ruleCounts[v], RuleTable[v])
	}
}

func DebugFn(skip int) bool {
	fname := strings.Trim(FuncName(skip), "main.")

	if *fnName == "" {
		return false
	} else if *debugPtr {
		return true
	} else {
		return *fnName == fname
	}
}

func PrepPmat(input string) {
	mat = PopulateMat(input)
	emptyCnt = CountEmpty(mat)

	emptyL, mat2 = GetPossibleMat(mat)
}

func iterMat(currCell *Cell) {

	if emptyCnt > 0 {
		iterCnt++

		for _, num := range currCell.Vals {
			if emptyCnt > 0 {
				if IsSafe(mat3, currCell.Row, currCell.Col, num) {
					mat3[currCell.Row][currCell.Col] = num
					emptyCnt--

					if emptyCnt > 0 {
						iterMat(currCell.Next)
						if emptyCnt > 0 {
							mat3[currCell.Row][currCell.Col] = 0
							emptyCnt++
						}
					} else {
						color.LightRed.Println("******* Finished *******")
					}
				}
			}
		}
	}
}

// erase digit from row of possibility matrix in the case of naked pairs
func eraseDigitsFromRowOfPairs(row, col, col2 int, digits []int) bool {
	erased := false

	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col && c != col2 {
			if Contains(mat2[row][c], digits[0]) {
				mat2[row][c] = EraseFromSlice(mat2[row][c], digits[0])
				emptyL.EraseDigitFromCell(row, c, digits[0])
				erased = true

				if *verbose {
					color.LightMagenta.Printf("Found naked pair (%d,%d) in row %d. Deleted %d from [%d,%d]\n",
						digits[0], digits[1], row, digits[0], row, c)
				}
			}

			if Contains(mat2[row][c], digits[1]) {
				mat2[row][c] = EraseFromSlice(mat2[row][c], digits[1])
				emptyL.EraseDigitFromCell(row, c, digits[1])
				erased = true

				if *verbose {
					color.LightMagenta.Printf("Found naked pair (%d,%d) in row %d. Deleted %d from [%d,%d]\n",
						digits[0], digits[1], row, digits[1], row, c)
				}
			}
		}
	}

	return erased
}

// erase digit from col of possibility matrix in the case of naked pairs
func eraseDigitsFromColOfPairs(row, col, row2 int, digits []int) bool {
	erased := false

	for r := 0; r < N; r++ {
		if mat2[r][col] != nil && r != row && r != row2 {
			if Contains(mat2[r][col], digits[0]) {
				mat2[r][col] = EraseFromSlice(mat2[r][col], digits[0])
				emptyL.EraseDigitFromCell(r, col, digits[0])
				erased = true

				if *verbose {
					color.LightMagenta.Printf("Found naked pair (%d,%d) in col %d. Deleted %d from [%d,%d]\n",
						digits[0], digits[1], col, digits[0], r, col)
				}
			}

			if Contains(mat2[r][col], digits[1]) {
				mat2[r][col] = EraseFromSlice(mat2[r][col], digits[1])
				emptyL.EraseDigitFromCell(r, col, digits[1])
				erased = true

				if *verbose {
					color.LightMagenta.Printf("Found naked pair (%d,%d) in col %d. Deleted %d from [%d,%d]\n",
						digits[0], digits[1], col, digits[1], r, col)
				}
			}
		}
	}

	return erased
}

// erase digit from row of possibility matrix in the case of naked pairs
func eraseDigitsFromBlkOfPairs(row, col, row2, col2 int, digits []int) bool {
	erased := false
	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if mat2[x][y] != nil && !(x == row && y == col) && !(x == row2 && y == col2) {

				if Contains(mat2[x][y], digits[0]) {
					mat2[x][y] = EraseFromSlice(mat2[x][y], digits[0])
					emptyL.EraseDigitFromCell(x, y, digits[0])
					erased = true

					if *verbose {
						color.LightMagenta.Printf("Found naked pair (%d,%d) in blk [%d,%d]. Deleted %d from [%d,%d]\n",
							digits[0], digits[1], row/SQ, col/SQ, digits[0], x, y)
					}
				}

				if Contains(mat2[x][y], digits[1]) {
					mat2[x][y] = EraseFromSlice(mat2[x][y], digits[1])
					emptyL.EraseDigitFromCell(x, y, digits[1])
					erased = true

					if *verbose {
						color.LightMagenta.Printf("Found naked pair (%d,%d) in blk [%d,%d]. Deleted %d from [%d,%d]\n",
							digits[0], digits[1], row/SQ, col/SQ, digits[1], x, y)
					}
				}
			}
		}
	}

	return erased
}

// *******************************************************************************************************
// *                                     end of funcs for naked pairs                                    *
// *******************************************************************************************************

// erase digit from row of possibility matrix
func eraseDigitFromRow(row, col, dig int) bool {
	erased := false

	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col {
			if Contains(mat2[row][c], dig) {
				mat2[row][c] = EraseFromSlice(mat2[row][c], dig)
				// remove this digit from cell at this position of the empty list
				emptyL.EraseDigitFromCell(row, c, dig)
				erased = true
			}
		}
	}

	return erased
}

func eraseDigitFromCol(row, col, dig int) bool {
	erased := false

	for r := 0; r < N; r++ {
		if mat2[r][col] != nil && r != row {
			if Contains(mat2[r][col], dig) {
				mat2[r][col] = EraseFromSlice(mat2[r][col], dig) // remove from possibility mat
				// remove this digit from cell at this position of the empty list
				emptyL.EraseDigitFromCell(r, col, dig)
				erased = true
			}
		}
	}

	return erased
}

func eraseDigitFromBlk(row, col, dig int) bool {
	erased := false

	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if mat2[x][y] != nil && x != row && y != col {
				if Contains(mat2[x][y], dig) {
					mat2[x][y] = EraseFromSlice(mat2[x][y], dig)
					// remove this digit from cell at this position of the empty list
					emptyL.EraseDigitFromCell(x, y, dig)
					erased = true
				}
			}
		}
	}

	return erased
}

// *******************************************************************************************************
// *                                     start of funcs for X-wing                                       *
// *******************************************************************************************************

// erase digit from row of possibility matrix. digits is list of nos. to be erased. cols is exception list
func eraseDigitsFromRowMulti(row int, digits, cols []int) (int, bool) {
	count := 0
	erased := false

	for c := 0; c < N; c++ {
		if mat2[row][c] != nil {
			inCol := false
			for _, col := range cols {
				if c == col {
					inCol = true
				}
			}

			if !inCol {
				for _, dig := range digits {
					if Contains(mat2[row][c], dig) {
						mat2[row][c] = EraseFromSlice(mat2[row][c], dig)
						// remove this digit from cell at this position of the empty list
						emptyL.EraseDigitFromCell(row, c, dig)
						erased = true
						count++
					}
				}
			}
		}
	}

	return count, erased
}

// erase digit from col of possibility matrix
func eraseDigitFromColMulti(col, dig int, rows []int) (int, bool) {
	debug := DebugFn(3)
	count := 0
	erased := false

	for r := 0; r < N; r++ {
		if mat2[r][col] != nil {
			inRow := false
			for _, row := range rows {
				if r == row {
					inRow = true
				}
			}

			if !inRow {
				if Contains(mat2[r][col], dig) {
					mat2[r][col] = EraseFromSlice(mat2[r][col], dig)
					// remove this digit from cell at this position of the empty list
					emptyL.EraseDigitFromCell(r, col, dig)
					erased = true
					count++

					if debug {
						color.LightRed.Printf("Erased digit %d from Cell [%d,%d].\n", dig, r, col)
					}
				}
			}
		}
	}

	return count, erased
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
	return matched, count
}

// Rule 1a	Open singles
//          Search the specified row or column for open singles
func rule1a(row, col int) (*Matchlist, int) {
	var (
		digit, count                 int
		notInRow, notInCol, notInBlk bool
		node                         *Cell
		matched                      *Matchlist
	)
	matched = &Matchlist{}

	if *debugPtr {
		fmt.Println("In rule1a...")
	}

	if row >= 0 && col < 0 { // skip row checking if negative value
		for c := 0; c < N; c++ {
			if len(mat2[row][c]) == 1 {
				digit = mat2[row][c][0]
				node = emptyL.GetNodeForCell(row, c)
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
				node = emptyL.GetNodeForCell(r, col)
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
			node = emptyL.GetNodeForCell(row, col)
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

	return matched, count
}

// Rule 2	Lone singles by cross out
//          After crossing out columns, rows and blocks containing the no.,
//          there is a cell with only 1 possibility (single pencil mark) left either in a row, column or block.
//          In addition, after filling in the no., erase the digit from interecting row, column and within the block.
/*func rule2() int {
	var (
		//col, row int // position of last empty cell
		count int
	)

	return count
}
*/
// Rule 3	Hidden singles
//          A digit that is the only one in an entire row, column, or block.
//          Fill in this digiti and erase any other occurrence of this digit in the same row, column or block.
func rule3() (*Matchlist, int) {
	var (
		count, itercnt           int
		foundHiddenSingle, debug bool
		matched                  *Matchlist
	)
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

	return matched, count
}

// Rule 3a	Hidden singles
//			Search in the specified row or col or blk intersecting this Cell only
func rule3a(row, col, dig int) (*Matchlist, int) {
	var (
		count                    int
		foundHiddenSingle, debug bool
		currNode                 *Cell
		matched                  *Matchlist
	)
	debug = DebugFn(1)
	matched = &Matchlist{}
	currNode = emptyL.GetNodeForCell(row, col)

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

	return matched, count
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

// Rule 5	Naked pairs
//          A pair of digits that occurs in exactly 2 cells in an entire row, column, or block.
//          Erase any other occurrence of these 2 digits elsewhere in the same row, column or block.
func rule5() (*Matchlist, int) {
	var (
		col, row, col2, row2   int // position of last empty cell
		count                  int
		twoElem                []int
		foundNakedPairs, debug bool
		inRow, inCol, inBlk    bool
		secondNode             *Cell
		matched, matchedBlk    *Matchlist
	)
	matched = &Matchlist{}
	matchedBlk = &Matchlist{}
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
								color.Magenta.Printf("Found naked pair in row %d, in cols %d and %d.\n", row, col, col2)
							}

							secondNode = emptyL.GetNodeForCell(row, col2)

							arr := AddIdx(nil, currNode, secondNode)

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

							secondNode = emptyL.GetNodeForCell(row2, col)

							arr := AddIdx(nil, currNode, secondNode)

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
									color.Magenta.Printf("Found naked pair in blk [%d,%d], in cells [%d,%d] and [%d,%d].\n",
										row/SQ, col/SQ, row, col, row2, col2)
								}

								secondNode = emptyL.GetNodeForCell(row2, col2)

								if debug {
									fmt.Printf("Found naked pair: %d,%d. x,y = %d,%d.\n",
										twoElem[0], twoElem[1], row2, col2)
								}

								arr := AddIdx(nil, currNode, secondNode)

								if !matchedBlk.ContainsPair(arr) {
									matchedBlk.AddRNode(arr)

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

	return matched, count
}

// Rule 6	Naked triplets
//			There are 3 dells in the same row, col or block that have exactly the same 3 digits.
//			Variation: It can also be a combination of 2 digits in one cell and 3 digits in another cell.
//	        E.g. The 3 cells contain (5,6), (6,8) and (8,5).
//			They form a closed loop 5 -> 8 -> 6 going from (5,6) to (5, 8) to (8,6).
func rule6() (*Matchlist, int) {
	var (
		count int
		//foundTrip bool
		matched *Matchlist
	)

	matched = &Matchlist{}
	return matched, count

	// Search for exactly cells that contain 3 same digits each.

	// search rows

}

// Find naked triplets in all rows for the case where all 3 cells have 3 same digits each
func findTripInRow(m2 Pmat) (int, bool) {
	var (
		count, tripCnt int
		trip           []int
		found          bool
		debug          bool
	)
	debug = DebugFn(2)
	found = false

	for i := 0; i < N; i++ {
		count = 0

		for j := 0; j < N; j++ {
			if m2[i][j] != nil && len(m2[i][j]) == 3 {
				if count == 0 {
					trip = m2[i][j]
					count = 1

					if debug {
						fmt.Printf("Triplets: %v at [%d,%d]. Count = %d\n", trip, i, j, count)
					}

				} else {
					if IntArrayEquals(trip, m2[i][j]) {
						count++

						if debug {
							fmt.Printf("Triplets: %v at [%d,%d]. Count = %d\n", trip, i, j, count)
						}
					}
				}
			}
		}

		if count == 3 {
			tripCnt++
			found = true
		}
	}

	return tripCnt, found
}

// Variation of Naked Triplets rule:
// Find naked triplets in row for the case where not all 3 cells contain 3 digits each.
// At least 1 cell must contain 3 digits.
// Solution: Search for a cell with 3 digits. Then search for exactly 2 more cells which contain either 2 or 3 digits
// where the digits are the same as that of the first 3-digit cell.
func findTripInRowVar(m2 Pmat, row int) bool {
	var (
		found         bool
		totCnt, count int
		trip, colList []int
	)

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			count = 0
			colList = nil
			found = false

			if m2[i][j] != nil && len(m2[i][j]) == 3 { // first cell with 3 digits
				trip = m2[i][j]
				count = 1
				colList = append(colList, j)

				for c := 0; c < N; c++ {
					if m2[i][j] != nil && len(m2[i][j]) >= 2 && len(m2[i][j]) <= 3 && j != c {

						// case of 3 digits
						if len(m2[i][j]) == 3 {
							if IntArrayEquals(trip, m2[i][c]) {
								colList = append(colList, c)
								count++
							}
						} else {
							// case of 2 digits
							cnt := 0
							for _, v := range m2[i][c] {
								if Contains(trip, v) {
									cnt++
								}
							}

							if cnt == 2 {
								colList = append(colList, c)
								count++
							}
						}
					}
				}

				if count == 3 {
					found = true
					totCnt++ // no. of naked triplets found
				}

				if found {
					// erase the same 3 digits from elsewhere in row
					eraseDigitsFromRowMulti(i, trip, colList)
					fmt.Println()
				}
			}
		}
	}

	return false
}

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
