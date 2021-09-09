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

var (
	iterCnt  int
	emptyCnt int
	mat      Intmat
	mat2     Pmat   // matrix with possible values in empty cells
	mat3     Intmat // guessed matrix
	emptyL   *LinkedList
	debugPtr *bool = flag.Bool("debug", false, "verbose debug mode")
	prtLLPtr *bool = flag.Bool("prtLL", false, "print the linked list of empty cells")
	rule     *int  = flag.Int("r", 0, "The deffault is 0, which will iterate matrix using linked list.")

	RuleTable = map[int]string{
		1: "Open single",
		3: "Hidden single",
		5: "Naked pair",
		8: "Hidden pair",
	}
)

func main() {
	var (
		start   time.Time
		elapsed time.Duration
	)
	flag.Parse()

	mat = PopulateMat(readInput())
	emptyCnt = CountEmpty(mat)
	fmt.Printf("Empty cells: %d\n", emptyCnt)
	start = time.Now()
	printSudoku(mat)
	emptyL, mat2 = GetPossibleMat(mat)
	fmt.Println("Starting possibility matrix.")
	printPossibleMat()

	if *prtLLPtr {
		emptyL.ShowAllEmptyCells()
	}

	switch *rule {
	case 0:
		fmt.Println("Default to iterMat.")
		mat3 = mat
		iterMat(emptyL.Head)
		printSudoku(mat3)
		checkSums(mat3)
	case 1:
		RuleLoop(rule1, "Open single")
	case 3:
		RuleLoop(rule3, "Hidden single")
	case 5:
		RuleLoop(rule5, "Naked pairs")
	case 99: // run all rules
		rule1cnt := 0
		rule3cnt := 0
		rule5cnt := 0

		for {
			matched1, cnt1 := rule1()
			fmt.Printf("After rule1, found %2d. Empty list count = %2d.\n", cnt1, emptyL.CountNodes())
			matched1.PrintResult("Found open single")

			matched3, cnt3 := rule3()
			fmt.Printf("After rule3, found %2d. Empty list count = %2d.\n", cnt3, emptyL.CountNodes())
			matched3.PrintResult("Found hidden single")

			matched5, cnt5 := rule5()

			fmt.Printf("After rule5, found %2d. Empty list count = %2d.\n", cnt5, emptyL.CountNodes())
			matched5.PrintResult("Found hidden single")

			if cnt1 <= 0 && cnt3 <= 0 && cnt5 <= 0 {
				break
			}
			rule1cnt += cnt1
			rule3cnt += cnt3
			rule5cnt += cnt5
		}

		ecnt := emptyL.CountNodes()
		fmt.Printf("After rules 1 & 3 & 5 have completed. Empty count : %d\n", ecnt)
		printSudoku(mat)

		if emptyL.CountNodes() > 0 {
			printPossibleMat()
			// do iterations
			fmt.Println("Start iterations.")
			mat3 = mat
			if emptyL.CountNodes() > 0 {
				fmt.Printf("Empty list count before running iterMat = %d.\n", emptyL.CountNodes())
				printPossibleMat()
				iterMat(emptyL.Head)
			}
			printSudoku(mat)
		} else {
			checkSums(mat)
		}

		fmt.Printf("Rule 1: Found %2d open singles\n", rule1cnt)
		fmt.Printf("Rule 3: Found %2d hidden singles\n", rule3cnt)
		fmt.Printf("Rule 5: Found %2d naked pairs\n", rule5cnt)
		fmt.Printf("Empty cells : %2d\n", emptyL.CountNodes())
	}

	elapsed = time.Since(start)
	log.Printf("IterMat: Iterations: %d. Empty cells: %d. Sudoku took %v sec\n", iterCnt, CountEmpty(mat), elapsed.Seconds())
}

func RuleLoop(rule fnRule, desc string) int {
	totalcnt := 0
	fnName := GetFunctionName(rule)
	for {
		matched, cnt := rule()
		fmt.Printf("%s: Found %d digits.\n", fnName, cnt)
		matched.PrintResult(desc)
		printSudoku(mat)
		if cnt <= 0 {
			break
		}
		totalcnt += cnt
	}
	fmt.Printf("%s: Total found = %d.\n", fnName, totalcnt)
	printPossibleMat()
	return totalcnt
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

// ****************************************** start of find/erase fns ******************************************

func PrepPmat(input string) {
	mat = PopulateMat(input)
	emptyCnt = CountEmpty(mat)

	emptyL, mat2 = GetPossibleMat(mat)
}

func inRow(m Intmat, row, num int) bool {
	for c := 0; c < N; c++ {
		if m[row][c] == num {
			return true
		}
	}
	return false
}

func inCol(m Intmat, col, num int) bool {
	for r := 0; r < N; r++ {
		if m[r][col] == num {
			return true
		}
	}
	return false
}

func inSqu(m Intmat, row, col, num int) bool {
	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if m[x][y] == num {
				return true
			}
		}
	}
	return false
}

func iterMat(currCell *Cell) {

	if emptyCnt > 0 {
		iterCnt++

		for _, num := range currCell.Vals {
			if emptyCnt > 0 {
				if isSafe(mat3, currCell.Row, currCell.Col, num) {
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

func isSafe(m Intmat, row, col, num int) bool {
	if !inRow(m, row, num) && !inCol(m, col, num) && !inSqu(m, row, col, num) {
		return true
	}
	return false
}

func readInput() string {
	var (
		l int
		s string
	)
	fmt.Println("Enter sudoku string (. rep empty square)")
	fmt.Scanf("%s\n", &s)
	l = len(s)
	fmt.Printf("Length: %d\n", l)
	if l != 81 {
		panic("Expected 81")
	}
	return s
}

// check the row, col and square sums
func checkSums(m Intmat) bool {
	const (
		totalSum = 405
	)
	var (
		val     int
		rowSums []int
		colSums []int
		sqSums  []int
		success bool
	)
	success = false

	for i := 0; i < N; i++ {
		rSum := 0
		cSum := 0
		for j := 0; j < N; j++ {
			val = m[i][j]
			if val > 0 {
				rSum += val
			}

			val = m[j][i]
			if val > 0 {
				cSum += val
			}
		}
		rowSums = append(rowSums, rSum)
		colSums = append(colSums, cSum)

	}

	for i := 0; i < SQ; i++ {
		for j := 0; j < SQ; j++ {
			sSum := 0
			startrow := i * SQ
			startcol := j * SQ
			for x := startrow; x < startrow+SQ; x++ {
				for y := startcol; y < startcol+SQ; y++ {
					val = m[x][y]
					if val > 0 {
						sSum += val
					}
				}
			}
			sqSums = append(sqSums, sSum)
		}
	}

	fmt.Printf("Total sums: %2v\n", totalSum)
	fmt.Printf("Row sums: %2v\n", rowSums)
	fmt.Printf("Col sums: %2v\n", colSums)
	fmt.Printf("Squ sums: %2v\n", sqSums)

	if sumArr(rowSums) == totalSum && sumArr(colSums) == totalSum && sumArr(sqSums) == totalSum {
		color.New(color.FgLightBlue, color.OpBold).Println("Finished!")
		success = true
	}
	return success
}

func sumArr(arr []int) int {
	result := 0
	for _, v := range arr {
		result += v
	}
	return result
}

func printSudoku(m Intmat) {
	var sqi, sqj int
	for i := 0; i < N; i++ {
		sqi = (i / SQ) % 2
		for j := 0; j < N; j++ {
			sqj = (j / SQ) % 2
			if (sqi == 0 && sqj == 1) || (sqi == 1 && sqj == 0) {
				color.LightBlue.Printf("%d ", m[i][j])
			} else {
				color.LightGreen.Printf("%d ", m[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println("-----------------")
}

func printPossibleMat() {
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf("|%-20v ", arr2String(mat2[i][j], ","))
		}
		fmt.Println("|")
		if i != 0 && i%SQ == 2 {
			fmt.Println("-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")
		}
	}
}

func arr2String(a []int, delim string) string {
	return strings.Trim(strings.Join(strings.Fields(fmt.Sprint(a)), delim), "[]")
}

// check row of possibility matrix
func findDigitInRow(mat2 Pmat, row, col, dig int) bool {
	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", row, c, mat2[row][c])
			}

			if Contains(mat2[row][c], dig) {

				if *debugPtr {
					fmt.Println("contains digit")
				}
				return true

			}
		}
	}

	return false
}

// check column of possibility matrix
func findDigitInCol(mat2 Pmat, row, col, dig int) bool {
	if *debugPtr {
		fmt.Println("In findDigitInCol...")
	}

	for r := 0; r < N; r++ {
		if mat2[r][col] != nil && r != row {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", r, col, mat2[r][col])
			}

			if Contains(mat2[r][col], dig) {
				if *debugPtr {
					fmt.Println("contains digit")
				}
				return true
			}
		}
	}

	return false
}

// check block of possibility matrix corresponding to cell [row[col]
func findDigitInBlk(mat2 Pmat, row, col, dig int) bool {
	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if mat2[x][y] != nil && !(x == row && y == col) {
				if *debugPtr {
					color.White.Printf("Cell [%d][%d] = %v\n", x, y, mat2[x][y])
				}

				if Contains(mat2[x][y], dig) {

					if *debugPtr {
						fmt.Println("contains digit")
					}
					return true
				}
			}
		}
	}

	return false
}

// *******************************************************************************************************
// *                                          end of find funcs                                          *
// *******************************************************************************************************

// * start of funcs for naked pairs *

// check row of possibility matrix
func findDigitInRowPair(mat2 Pmat, row, col, col2 int, digits []int) bool {
	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col && c != col2 {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", row, c, mat2[row][c])
			}
			if ContainsMulti(mat2[row][c], digits) {

				if *debugPtr {
					fmt.Printf("Cell [%d,%d] contains digits of naked pair", row, c)
				}
				return true
			}
		}
	}

	return false
}

// erase digit from row of possibility matrix in the case of naked pairs
func eraseDigitsFromRowOfPairs(row, col, col2 int, digits []int) bool {
	erased := false

	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col && c != col2 {
			if Contains(mat2[row][c], digits[0]) {
				mat2[row][c] = EraseFromSlice(mat2[row][c], digits[0])
				emptyL.EraseDigitFromCell(row, c, digits[0])
				color.LightMagenta.Printf("Found naked pair (%d,%d) in row %d. Deleted %d from [%d,%d]\n",
					digits[0], digits[1], row, digits[0], row, c)
				erased = true
			}

			if Contains(mat2[row][c], digits[1]) {
				mat2[row][c] = EraseFromSlice(mat2[row][c], digits[1])
				emptyL.EraseDigitFromCell(row, c, digits[1])
				color.LightMagenta.Printf("Found naked pair (%d,%d) in row %d. Deleted %d from [%d,%d]\n",
					digits[0], digits[1], row, digits[1], row, c)
				erased = true
			}
		}
	}

	return erased
}

// check col of possibility matrix if any of the digits in the naked pair is
// found in this col
func findDigitInColPair(mat2 Pmat, row, col, row2 int, digits []int) bool {
	for r := 0; r < N; r++ {
		if mat2[r][col] != nil && r != row && r != row2 {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", r, col, mat2[r][col])
			}
			if ContainsMulti(mat2[r][col], digits) {

				if *debugPtr {
					fmt.Printf("Cell [%d,%d] contains digits of naked pair", r, col)
				}
				return true
			}
		}
	}

	return false
}

// erase digit from col of possibility matrix in the case of naked pairs
func eraseDigitsFromColOfPairs(row, col, row2 int, digits []int) bool {
	erased := false

	for r := 0; r < N; r++ {
		if mat2[r][col] != nil && r != row && r != row2 {
			if Contains(mat2[r][col], digits[0]) {
				mat2[r][col] = EraseFromSlice(mat2[r][col], digits[0])
				emptyL.EraseDigitFromCell(r, col, digits[0])
				color.LightMagenta.Printf("Found naked pair (%d,%d) in col %d. Deleted %d from [%d,%d]\n",
					digits[0], digits[1], col, digits[0], r, col)
				erased = true
			}

			if Contains(mat2[r][col], digits[1]) {
				mat2[r][col] = EraseFromSlice(mat2[r][col], digits[1])
				emptyL.EraseDigitFromCell(r, col, digits[1])
				color.LightMagenta.Printf("Found naked pair (%d,%d) in col %d. Deleted %d from [%d,%d]\n",
					digits[0], digits[1], col, digits[1], r, col)
				erased = true
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
	const ncols = 9
	erased := false

	for c := 0; c < ncols; c++ {
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
	const ncols = 9
	erased := false

	for r := 0; r < ncols; r++ {
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
	const (
		scols = 3
		ncols = 9
	)
	erased := false

	startRow := row / scols * scols
	startCol := col / scols * scols

	for x := startRow; x < startRow+scols; x++ {
		for y := startCol; y < startCol+scols; y++ {
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

func getColOfPossibleMat(mat2 Pmat, col int) [][]int {
	var m [][]int

	for i := 0; i < N; i++ {
		m = append(m, mat2[i][col])
	}

	return m
}

func getBlkOfPossibleMat(mat2 Pmat, row, col int) [][]int {
	var m [][]int

	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			m = append(m, mat2[x][y])
		}
	}

	return m
}

// *******************************************************************************

// Rule 1 - Open cell - 1 cell empty either in column, row or block.
// Search the possible matrix for any col, row or block that has only 1 empty cell
// Empty cells contain lists. Non-empty cells contain nil.
func rule1() (*Matchlist, int) {
	var (
		col, row     int // position of last empty cell
		digit, count int
		found        bool
		matched      *Matchlist
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
				found = findDigitAndUpdate(currNode, digit)
				if found {
					matched.AddCell(currNode, digit)
					count++
				}
			}
			currNode = currNode.Next
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
		count, itercnt    int
		foundHiddenSingle bool
		matched           *Matchlist
	)
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
					if *debugPtr {
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
						if *debugPtr {
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

func findDigitAndUpdate(currNode *Cell, dig int) bool {
	var (
		row, col                     int
		notInRow, notInCol, notInBlk bool
		found                        bool
	)
	row = currNode.Row
	col = currNode.Col

	// check that there is no occurrence in same row, col or block
	notInRow = !findDigitInRow(mat2, row, col, dig)
	notInCol = !findDigitInCol(mat2, row, col, dig)
	notInBlk = !findDigitInBlk(mat2, row, col, dig)

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

	if notInRow || notInCol || notInBlk {
		found = true
		emptyL.DelNode(currNode) // remove current Node from possibility list
		mat[row][col] = dig      // fill in dig in resulting mat
		mat2[row][col] = nil     // blank this cell out in possibility mat
		emptyCnt--

		// erase any occurrence of the digit in the same row, col or block
		if !notInRow {
			eraseDigitFromRow(row, col, dig)

			if *debugPtr {
				color.LightBlue.Printf("After deletion from row %d: %v\n", row, mat2[row])
			}
		}
		if !notInCol {
			eraseDigitFromCol(row, col, dig)

			if *debugPtr {
				color.LightBlue.Printf("After deletion from col %d: %v\n", col, getColOfPossibleMat(mat2, col))
			}
		}
		if !notInBlk {
			eraseDigitFromBlk(row, col, dig)

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
	return found
}

// Rule 5	Naked pairs
//          A pair of digits that occurs in exactly 2 cells in an entire row, column, or block.
//          Erase any other occurrence of these 2 digits elsewhere in the same row, column or block.
func rule5() (*Matchlist, int) {
	var (
		col, row, col2, row2 int // position of last empty cell
		count                int
		twoElem              []int
		foundNakedPairs      bool
		inRow, inCol         bool
		secondNode           *Cell
		matched              *Matchlist
	)
	matched = &Matchlist{}

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
				if *debugPtr {
					color.LightGreen.Printf("cell [%d][%d]. %+v\n", row, col, currNode.Vals)
				}

				if len(currNode.Vals) == 2 { // has 2 possible values
					twoElem = currNode.Vals
					// check row
					for c := 0; c < N; c++ {
						if IntArrayEquals(mat2[row][c], twoElem) && c != col {
							col2 = c

							if *debugPtr {
								color.Magenta.Printf("Found naked pair in row %d, in cols %d and %d.\n", row, col, col2)
							}

							secondNode = emptyL.GetNodeForCell(row, col2)

							arr := AddIdx(nil, currNode, secondNode)

							if !matched.ContainsPair(arr) {
								matched.AddRNode(arr)

								fmt.Printf("Row %d\n", row)
								inRow = findDigitInRowPair(mat2, row, col, col2, twoElem)
								if inRow {
									if *debugPtr {
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

							if *debugPtr {
								color.Magenta.Printf("Found naked pair in col %d, in rows %d and %d.\n", col, row, row2)
							}

							secondNode = emptyL.GetNodeForCell(row2, col)

							arr := AddIdx(nil, currNode, secondNode)

							if !matched.ContainsPair(arr) {
								matched.AddRNode(arr)

								fmt.Printf("Col %d\n", col)
								inCol = findDigitInColPair(mat2, row, col, row2, twoElem)
								if inCol {
									if *debugPtr {
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
				}

				currNode = currNode.Next
			}
		}
	}

	return matched, count
}
