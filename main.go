package main

import (
	"flag"
	"fmt"
	"log"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"time"

	l "github.com/mjwong/sudoku2/lib"
	ll "github.com/mjwong/sudoku2/linkedlist"
	ml "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

const (
	N     = 9
	SQ    = 3
	nsize = 81
)

type (
	intmat    [N][N]int
	sqmat     [SQ][SQ][]int
	pmat      [N][N][]int
	matString string
	fnRule    func() (*ll.LinkedList, int)
	fnRule5   func() (*ml.Matchlist, int)
)

var (
	iterCnt  int
	emptyCnt int
	mat      intmat
	mask     intmat // matrix to mark cross-out cells
	mat2     pmat   // matrix with possible values in empty cells
	mat3     intmat // guessed matrix
	emptyL   *ll.LinkedList
	debugPtr *bool = flag.Bool("debug", false, "verbose debug mode")
	prtLLPtr *bool = flag.Bool("prtLL", false, "print the linked list of empty cells")
	rule     *int  = flag.Int("r", 0, "The deffault is 0, which will iterate matrix using linked list.")
)

func main() {
	var (
		start   time.Time
		elapsed time.Duration
	)
	flag.Parse()

	mat = populateMat(readInput())
	emptyCnt = countEmpty(mat)
	fmt.Printf("Empty cells: %d\n", emptyCnt)
	start = time.Now()
	printSudoku(mat)
	emptyL, mat2 = getPossibleMat(mat)
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
		ruleLoop(rule1, "Open single")
	case 3:
		ruleLoop(rule3, "Hidden single")
	case 13:
		rule1cnt := 0
		rule2cnt := 0

		for {
			matched1, cnt1 := rule1()
			fmt.Printf("After rule1, found %2d. Empty list count = %2d.\n", cnt1, emptyL.CountNodes())
			matched1.PrintResult("Found open single")

			matched2, cnt2 := rule3()
			fmt.Printf("After rule3, found %2d. Empty list count = %2d.\n", cnt2, emptyL.CountNodes())
			matched2.PrintResult("Found hidden single")

			if cnt1 <= 0 && cnt2 <= 0 {
				break
			}
			rule1cnt += cnt1
			rule2cnt += cnt2
		}

		fmt.Println("After rules 1 & 3 have completed.")
		node := emptyL.Head
		for node != nil {
			fmt.Printf("[%d,%d] : %v \n", node.Row, node.Col, node.Vals)
			node = node.Next
		}
		printSudoku(mat)

		// do iterations
		fmt.Println("Start iterations.")
		mat3 = mat
		if emptyL.CountNodes() > 0 {
			fmt.Printf("Empty list count = %d.\n", emptyL.CountNodes())
			printPossibleMat()
			iterMat(emptyL.Head)
		}

		printSudoku(mat)
		fmt.Printf("Rule 1: Found %2d digits\n", rule1cnt)
		fmt.Printf("Rule 3: Found %2d digits\n", rule2cnt)
	case 5:
		ruleLoop5(rule5, "Naked pairs")
	}

	//color.LightMagenta.Println(*emptyL)

	elapsed = time.Since(start)
	log.Printf("IterMat: Iterations: %d. Empty cells: %d. Sudoku took %v sec\n", iterCnt, countEmpty(mat), elapsed.Seconds())
}

func ruleLoop(rule fnRule, desc string) {
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
}

func ruleLoop5(rule fnRule5, desc string) {
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
}

// ****************************************** start of linkedListPairs ******************************************

// ****************************************** end of linkedListPairs ******************************************

func getPossibleMat(mat intmat) (*ll.LinkedList, pmat) {
	const (
		ncols = 9
	)
	var (
		mat2                pmat
		emptyL              *ll.LinkedList
		inRow, inCol, inSqu bool
		valList             []int
	)

	emptyL = ll.CreatelinkedList()

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
					inRow = contains(getArrForRow(mat, i), p)
					inCol = contains(getArrForCol(mat, j), p)
					inSqu = contains(getArrForSqu(mat, i, j), p)

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

func inRow(m intmat, row, num int) bool {
	for c := 0; c < N; c++ {
		if m[row][c] == num {
			return true
		}
	}
	return false
}

func inCol(m intmat, col, num int) bool {
	for r := 0; r < N; r++ {
		if m[r][col] == num {
			return true
		}
	}
	return false
}

func inSqu(m intmat, row, col, num int) bool {
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

func iterMat(currCell *ll.Cell) {

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

func isSafe(m intmat, row, col, num int) bool {
	if !inRow(m, row, num) && !inCol(m, col, num) && !inSqu(m, row, col, num) {
		return true
	}
	return false
}

// Get calling function name
func funcName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// Get function name, e.g. main.foo
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func getArrForRow(m intmat, i int) []int {
	var arrRow []int

	for j := 0; j < N; j++ {
		if m[i][j] != 0 {
			arrRow = append(arrRow, m[i][j])
		}
	}
	return arrRow
}

func getArrForCol(m intmat, j int) []int {
	var arrCol []int

	for i := 0; i < N; i++ {
		if m[i][j] != 0 {
			arrCol = append(arrCol, m[i][j])
		}
	}
	return arrCol
}

func getArrForSqu(m intmat, i, j int) []int {
	const scols = 3
	var (
		startRow, startCol int
		arrSq              []int
	)

	startRow = i / SQ * SQ
	startCol = j / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if m[x][y] != 0 {
				arrSq = append(arrSq, m[x][y])
			}
		}
	}
	return arrSq
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

func populateMat(s string) intmat {
	var (
		x, y int
		mat  intmat
	)

	// replace '.' with 0 and handle as INT
	for _, i := range strings.SplitAfter(s, "") {
		j, err := strconv.Atoi(i)
		if err != nil {
			if i != "." {
				panic(err)
			} else {
				j = 0
			}
		}

		if y > N-1 {
			y = 0
			x++
		}
		mat[x][y] = j
		y++
	}
	return mat
}

func contains(arr []int, v int) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func containsMulti(arr []int, v []int) bool {
	for _, a := range arr {
		fmt.Printf("arr elem: %d\n", a)
		for _, e := range v {
			if a == e {
				return true
			}
		}
	}
	return false
}

func countEmpty(mat intmat) int {
	var count int
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if mat[i][j] == 0 {
				count++
			}
		}
	}
	return count
}

// check the row, col and square sums
func checkSums(m intmat) bool {
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

func printSudoku(m intmat) {
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

// print with highlighted no. p
func printMask(p int) {
	var m int
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			m = mask[i][j]
			if m == -1 {
				color.Cyan.Printf("%2d ", mask[i][j])
			} else if m == 0 {
				color.Yellow.Printf("%2d ", mask[i][j])
			} else if m == p {
				color.Red.Printf("%2d ", mask[i][j])
			} else {
				fmt.Printf("%2d ", mask[i][j])
			}
		}
		fmt.Println()
	}
	fmt.Println("\n------------------------------")
}

// check row of possibility matrix
func findDigitInRow(mat2 pmat, row, col, dig int) bool {
	if *debugPtr {
		fmt.Println("In findDigitInRow...")
	}

	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", row, c, mat2[row][c])
			}

			if contains(mat2[row][c], dig) {

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
func findDigitInCol(mat2 pmat, row, col, dig int) bool {
	if *debugPtr {
		fmt.Println("In findDigitInCol...")
	}

	for r := 0; r < N; r++ {
		if mat2[r][col] != nil && r != row {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", r, col, mat2[r][col])
			}

			if contains(mat2[r][col], dig) {
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
func findDigitInBlk(mat2 pmat, row, col, dig int) bool {
	if *debugPtr {
		fmt.Println("In findDigitInBlk...")
	}

	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if mat2[x][y] != nil && !(x == row && y == col) {
				if *debugPtr {
					color.White.Printf("Cell [%d][%d] = %v\n", x, y, mat2[x][y])
				}

				if contains(mat2[x][y], dig) {

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
func findDigitInRowPair(mat2 pmat, row, col, col2 int, digits []int) bool {
	if *debugPtr {
		fmt.Println("In findDigitInRowPair..")
	}

	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col && c != col2 {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", row, c, mat2[row][c])
			}
			fmt.Println("here")
			if containsMulti(mat2[row][c], digits) {

				if *debugPtr {
					fmt.Println("contains digits")
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
			if contains(mat2[row][c], digits[0]) {
				mat2[row][c] = l.EraseFromSlice(mat2[row][c], digits[0])
				emptyL.EraseDigitFromCell(row, c, digits[0])
				fmt.Printf("Deleted %d from [%d,%d]\n", digits[0], row, c)
				erased = true
			}

			if contains(mat2[row][c], digits[1]) {
				mat2[row][c] = l.EraseFromSlice(mat2[row][c], digits[1])
				emptyL.EraseDigitFromCell(row, c, digits[1])
				fmt.Printf("Deleted %d from [%d,%d]\n", digits[1], row, c)
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
			if contains(mat2[row][c], dig) {
				mat2[row][c] = l.EraseFromSlice(mat2[row][c], dig)
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
			if contains(mat2[r][col], dig) {
				mat2[r][col] = l.EraseFromSlice(mat2[r][col], dig) // remove from possibility mat
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
				if contains(mat2[x][y], dig) {
					mat2[x][y] = l.EraseFromSlice(mat2[x][y], dig)
					// remove this digit from cell at this position of the empty list
					emptyL.EraseDigitFromCell(x, y, dig)
					erased = true
				}
			}
		}
	}

	return erased
}

// erase digit from row of possibility matrix
func erasePairFromElseWhereInRow(row, col, col2 int, pair []int) bool {
	const ncols = 9
	erased := false

	for c := 0; c < ncols; c++ {
		if mat2[row][c] != nil && c != col && c != col2 {
			if containsMulti(mat2[row][c], pair) {
				mat2[row][c] = l.EraseMultiFromSlice(mat2[row][c], pair)
				// remove any digit of the pair from cell at this position of the empty list
				emptyL.EraseDigitFromCell(row, c, pair[0])
				emptyL.EraseDigitFromCell(row, col2, pair[1])
				erased = true
			}
		}
	}

	return erased
}

func erasePairFromElseWhereInCol(row, row2, col int, pair []int) bool {
	const ncols = 9
	erased := false

	for r := 0; r < ncols; r++ {
		if mat2[r][col] != nil && r != row && r != row2 {
			if containsMulti(mat2[r][col], pair) {
				mat2[r][col] = l.EraseMultiFromSlice(mat2[r][col], pair) // remove from possibility mat
				// remove this digit from cell at this position of the empty list
				emptyL.EraseDigitFromCell(r, col, pair[0])
				emptyL.EraseDigitFromCell(row2, col, pair[1])
				erased = true
			}
		}
	}

	return erased
}

func getColOfPossibleMat(mat2 pmat, col int) [][]int {
	var m [][]int

	for i := 0; i < N; i++ {
		m = append(m, mat2[i][col])
	}

	return m
}

func getBlkOfPossibleMat(mat2 pmat, row, col int) [][]int {
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
func rule1() (*ll.LinkedList, int) {
	var (
		col, row                     int // position of last empty cell
		digit, count                 int
		notInRow, notInCol, notInBlk bool
		matched                      *ll.LinkedList
	)
	matched = &ll.LinkedList{}

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
				matched.AddCell(row, col, []int{digit})
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
func rule3() (*ll.LinkedList, int) {
	var (
		col, row                     int // position of last empty cell
		count, itercnt               int
		notInRow, notInCol, notInBlk bool
		foundHiddenSingle            bool
		matched                      *ll.LinkedList
	)
	matched = &ll.LinkedList{}

	for dig := 1; dig <= N; dig++ {
		itercnt = 0
		for {
			foundHiddenSingle = false
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

					if contains(currNode.Vals, dig) {
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
							matched.AddCell(row, col, []int{dig})
							emptyL.DelNode(currNode) // remove current Node from possibility list
							mat[row][col] = dig      // fill in dig in resulting mat
							mat2[row][col] = nil     // blank this cell out in possibility mat
							emptyCnt--
							count++
							foundHiddenSingle = true

							if *debugPtr {
								color.LightYellow.Printf("Rule 3. Found hidden single digit %d at cell [%d][%d]. Empty cells = %d\n", dig, row, col, emptyCnt)
							}

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
	}

	return matched, count
}

// Rule 5	Naked pairs
//          A pair of digits that occurs in exactly 2 cells in an entire row, column, or block.
//          Erase any other occurrence of these 2 digits elsewhere in the same row, column or block.
func rule5() (*ml.Matchlist, int) {
	const (
		scols = 3
		ncols = 9
	)

	var (
		col, row, col2, row2 int // position of last empty cell
		count                int
		twoElem              []int
		foundNakedPairs      bool
		inRow                bool
		secondNode           *ll.Cell
		matched              *ml.Matchlist
	)
	matched = &ml.Matchlist{}

	foundNakedPairs = true

	for foundNakedPairs {
		foundNakedPairs = false

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

				if len(currNode.Vals) == 2 { // has 2 possible values
					twoElem = currNode.Vals
					// check row
					for c := 0; c < ncols; c++ {
						if l.IntArrayEquals(mat2[row][c], twoElem) && c != col {
							col2 = c

							if *debugPtr {
								color.Magenta.Printf("Found naked pair in row %d, in cols %d and %d.\n", row, col, col2)
							}

							secondNode = emptyL.GetNodeForCell(row, col2)

							arr := ml.AddArrIdx([]ml.Idx{}, currNode)
							arr = ml.AddArrIdx(arr, secondNode)

							//fmt.Printf("Match Cnt: %d\n", matched.CountNodes())

							if !matched.ContainsPair(arr) {
								matched.AddNode(arr)
								foundNakedPairs = true
								count++
								fmt.Printf("Row %d\n", row)
								inRow = findDigitInRowPair(mat2, row, col, col2, twoElem)
								if inRow {
									fmt.Printf("Found digits of pairs in row %d.\n", row)
									eraseDigitsFromRowOfPairs(row, col, col2, twoElem)
								}
								break
							}
						}
					}

					// check col
					for r := 0; r < ncols; r++ {
						if l.IntArrayEquals(mat2[r][col], twoElem) && r != row {
							row2 = r

							if *debugPtr {
								color.Magenta.Printf("Found naked pair in col %d, in rows %d and %d.\n", col, row, row2)
							}

							secondNode = emptyL.GetNodeForCell(row2, col)

							arr := ml.AddArrIdx([]ml.Idx{}, currNode)
							arr = ml.AddArrIdx(arr, secondNode)

							if !matched.ContainsPair(arr) {
								matched.AddNode(arr)
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
