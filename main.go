package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gookit/color.v1"
)

const (
	ncols = 9
	scols = 3
	nsize = 81
)

type (
	intmat    [ncols][ncols]int
	sqmat     [scols][scols][]int
	pmat      [ncols][ncols][]int
	matString string
)

type cell struct {
	row  int
	col  int
	vals []int
	prev *cell
	next *cell
}

type linkedList struct {
	head        *cell
	last        *cell
	currentCell *cell
}

var (
	iterCnt  int
	emptyCnt int
	mat      intmat
	mask     intmat // matrix to mark cross-out cells
	mat2     pmat   // matrix with possible values in empty cells
	mat3     intmat // guessed matrix
	emptyL   *linkedList
	elapsed  time.Duration
	debugPtr *bool = flag.Bool("debug", false, "verbose debug mode")
	prtLLPtr *bool = flag.Bool("prtLL", false, "print the linked list of empty cells")
	rule     *int  = flag.Int("r", 0, "The deffault is 0, which will iterate matrix using linked list.")
)

func main() {
	flag.Parse()

	mat = populateMat(readInput())
	emptyCnt = countEmpty(mat)
	fmt.Printf("Empty cells: %d\n", emptyCnt)
	start := time.Now()
	printSudoku(mat)
	emptyL, mat2 = getPossibleMat(mat)
	fmt.Println("Starting possibility matrix.")
	printPossibleMat()

	if *prtLLPtr {
		emptyL.showAllEmptyCells()
	}

	switch *rule {
	case 0:
		mat3 = mat
		iterMat(*emptyL.head)
	case 1:
		matched, cnt := rule1()
		fmt.Printf("Rule1. Found %d digits.\n", cnt)
		matched.printResult("Found open single")
	case 3:
		matched, cnt := rule3()
		fmt.Printf("Rule3. Found %d digits.\n", cnt)
		matched.printResult("Found hidden single")
	case 13:
		rule1cnt := 0
		rule2cnt := 0

		for {
			_, cnt1 := rule1()
			fmt.Printf("After rule1, found %2d. Empty list count = %2d.\n", cnt1, emptyL.countNodes())

			_, cnt2 := rule3()
			fmt.Printf("After rule3, found %2d. Empty list count = %2d.\n", cnt2, emptyL.countNodes())

			if cnt2 <= 0 {
				break
			}
			rule1cnt += cnt1
			rule2cnt += cnt2
		}

		// do iterations
		mat3 = mat
		if emptyL.countNodes() > 0 {
			iterMat(*emptyL.head)
		} else {
			printSudoku(mat)
			fmt.Printf("Rule 1: Found %2d digits\n", rule1cnt)
			fmt.Printf("Rule 2: Found %2d digits\n", rule2cnt)
		}
	}

	//color.LightMagenta.Println(*emptyL)
	checkSums(mat3)
	elapsed := time.Since(start)
	log.Printf("IterMat: Iterations: %d. Empty cells: %d. Sudoku took %v sec\n", iterCnt, countEmpty(mat), elapsed.Seconds())
}

func createlinkedList() *linkedList {
	return &linkedList{}
}

func (p *linkedList) addCell(row, col int, arr []int) error {
	c := &cell{
		row:  row,
		col:  col,
		vals: arr,
	}

	if p.head == nil {
		p.head = c
		c.prev = nil
	} else {
		currentCell := p.head
		for currentCell.next != nil {
			currentCell = currentCell.next
		}
		currentCell.next = c
		c.prev = currentCell
		p.last = c
	}
	return nil
}

func (p *linkedList) showAllEmptyCells() error {
	currentCell := p.head
	if currentCell == nil {
		color.Red.Println("EmptyCell list is empty.")
		return nil
	}
	color.Green.Printf("%+v\n", *currentCell)
	for currentCell.next != nil {
		currentCell = currentCell.next
		color.Green.Printf("%+v\n", *currentCell)
	}
	return nil
}

func (p *linkedList) countNodes() int {
	count := 0
	currN := p.head
	for currN != nil {
		currN = currN.next
		count += 1
	}
	return count
}

func (p *linkedList) addNode(node *cell) error {

	if p.head == nil {
		p.head = node
		node.prev = nil
	} else {
		currentCell := p.head
		for currentCell.next != nil {
			currentCell = currentCell.next
		}
		currentCell.next = node
		node.prev = currentCell
		p.last = node
	}
	return nil
}

func (p *linkedList) lastNode() *cell {
	currentCell := p.head
	for currentCell.next != nil {
		currentCell = currentCell.next
	}
	return currentCell
}

func (p *linkedList) getNodeForCell(row, col int) *cell {

	currN := p.head
	for currN != nil {
		if currN.row == row && currN.col == col {
			return currN
		}
		currN = currN.next
	}

	return nil
}

func (p *linkedList) eraseDigitFromCell(row, col, dig int) {
	node := emptyL.getNodeForCell(row, col)
	node.vals = eraseFromSlice(node.vals, dig)
}

func (p *linkedList) printResult(desc string) {
	currNode := p.head
	for currNode != nil {
		color.LightMagenta.Printf("%s %d at [%d][%d].\n", desc, currNode.vals, currNode.row, currNode.col)
		currNode = currNode.next
	}
}

// remove current node from linked list and connect prev and next nodes
func (p *linkedList) delNode(node *cell) {

	if node == nil {
		color.Yellow.Println("Possibility list is empty or has reached the end.")
		return
	}

	if node == p.head {
		p.head = node.next
	} else {
		if node != nil {
			if node.next != nil {
				node.prev.next = node.next
				node.next.prev = node.prev
			} else {
				node.prev.next = nil
			}

			if *debugPtr {
				fmt.Printf("Before. Prev: %v. Curr: %v. Next: %v\n", color.Yellow.Render(node.prev),
					color.Yellow.Render(node), color.Yellow.Render(node.next))
				fmt.Printf("After. Prev: %v. Next: %v\n", color.Yellow.Render(node.prev),
					color.Yellow.Render(node.prev.next))
			}
		}
	}
}

// ****************************************** end of linkedList ******************************************

func getPossibleMat(mat intmat) (*linkedList, pmat) {
	const (
		ncols = 9
	)
	var (
		mat2                pmat
		emptyL              *linkedList
		inRow, inCol, inSqu bool
		valList             []int
	)

	emptyL = createlinkedList()

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
				emptyL.addCell(i, j, valList)
			}
		}
	}
	return emptyL, mat2
}

func iterMat(currCell cell) {
	var (
		inRow, inCol, inSqu, success bool
	)

	if emptyCnt > 0 {
		success = false
		iterCnt++
		if *debugPtr {
			color.Magenta.Printf("Iter: cell %v @ %d %d. %v. Empty: %d\n", &currCell, currCell.row, currCell.col, currCell.vals, emptyCnt)
		}

		for _, k := range currCell.vals {
			if emptyCnt > 0 {
				inRow = contains(getArrForRow(mat3, currCell.row), k)
				inCol = contains(getArrForCol(mat3, currCell.col), k)
				inSqu = contains(getArrForSqu(mat3, currCell.row, currCell.col), k)

				if !inRow && !inCol && !inSqu {
					mat3[currCell.row][currCell.col] = k
					emptyCnt--
					if *debugPtr {
						color.Cyan.Printf("Iter: Try %d at %d %d. Elapsed: %f\n", k, currCell.row, currCell.col, elapsed.Seconds())
						printSudoku(mat3)
					}
					success = true
					if emptyCnt > 0 {
						iterMat(*currCell.next)
						if emptyCnt > 0 {
							mat3[currCell.row][currCell.col] = 0
							if *debugPtr {
								color.Yellow.Printf("Iter: cell %v. Removing %d from %d %d\n", &currCell, k, currCell.row, currCell.col)
							}
							emptyCnt++
						}
					} else {
						printSudoku(mat3)
						color.LightRed.Println("******************************************************************************")
						color.LightRed.Println("******************************************************************************")
						color.LightRed.Println("******************************************************************************")
					}
				}
			}
		}

		if *debugPtr {
			if success {
				color.Green.Printf("Iter: Last element of cell %v\n", &currCell)
			} else {
				color.Red.Printf("Iter. Failed to find possible values for cell %d %d\n", currCell.row, currCell.col)
			}
		}
	}
}

func getArrForRow(m intmat, i int) []int {
	var arrRow []int

	for j := 0; j < ncols; j++ {
		if m[i][j] != 0 {
			arrRow = append(arrRow, m[i][j])
		}
	}
	return arrRow
}

func getArrForCol(m intmat, j int) []int {
	var arrCol []int

	for i := 0; i < ncols; i++ {
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

	startRow = i / scols * scols
	startCol = j / scols * scols

	for x := startRow; x < startRow+scols; x++ {
		for y := startCol; y < startCol+scols; y++ {
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

		if y > ncols-1 {
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

func countEmpty(mat intmat) int {
	var count int
	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
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
		scols    = 3
		ncols    = 9
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

	for i := 0; i < ncols; i++ {
		rSum := 0
		cSum := 0
		for j := 0; j < ncols; j++ {
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

	for i := 0; i < scols; i++ {
		for j := 0; j < scols; j++ {
			sSum := 0
			startrow := i * scols
			startcol := j * scols
			for x := startrow; x < startrow+scols; x++ {
				for y := startcol; y < startcol+scols; y++ {
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
	for i := 0; i < ncols; i++ {
		sqi = (i / scols) % 2
		for j := 0; j < ncols; j++ {
			sqj = (j / scols) % 2
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
	const (
		scols = 3
		ncols = 9
	)

	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")

	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
			fmt.Printf("|%-20v ", arr2String(mat2[i][j], ","))
		}
		fmt.Println("|")
		if i != 0 && i%scols == 2 {
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
	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
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
	const ncols = 9

	if *debugPtr {
		fmt.Println("In findDigitInRow...")
	}

	for c := 0; c < ncols; c++ {
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
	const ncols = 9

	if *debugPtr {
		fmt.Println("In findDigitInCol...")
	}

	for r := 0; r < ncols; r++ {
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
	const (
		scols = 3
		ncols = 9
	)

	if *debugPtr {
		fmt.Println("In findDigitInBlk...")
	}

	startRow := row / scols * scols
	startCol := col / scols * scols

	for x := startRow; x < startRow+scols; x++ {
		for y := startCol; y < startCol+scols; y++ {
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

// erase digit from row of possibility matrix
func eraseDigitFromRow(row, col, dig int) bool {
	const ncols = 9
	erased := false

	for c := 0; c < ncols; c++ {
		if mat2[row][c] != nil && c != col {
			if contains(mat2[row][c], dig) {
				mat2[row][c] = eraseFromSlice(mat2[row][c], dig)
				// remove this digit from cell at this position of the empty list
				emptyL.eraseDigitFromCell(row, c, dig)
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
				mat2[r][col] = eraseFromSlice(mat2[r][col], dig) // remove from possibility mat
				// remove this digit from cell at this position of the empty list
				emptyL.eraseDigitFromCell(r, col, dig)
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
					mat2[x][y] = eraseFromSlice(mat2[x][y], dig)
					// remove this digit from cell at this position of the empty list
					emptyL.eraseDigitFromCell(x, y, dig)
					erased = true
				}
			}
		}
	}

	return erased
}

func eraseFromSlice(sl []int, v int) []int {
	for i, a := range sl {
		if a == v {
			return remove(sl, i)
		}
	}
	return sl
}

func remove(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

func getColOfPossibleMat(mat2 pmat, col int) [][]int {
	var m [][]int

	for i := 0; i < ncols; i++ {
		m = append(m, mat2[i][col])
	}

	return m
}

func getBlkOfPossibleMat(mat2 pmat, row, col int) [][]int {
	const scols = 3
	var m [][]int

	startRow := row / scols * scols
	startCol := col / scols * scols

	for x := startRow; x < startRow+scols; x++ {
		for y := startCol; y < startCol+scols; y++ {
			m = append(m, mat2[x][y])
		}
	}

	return m
}

// Rule 1 - Open cell - 1 cell empty either in column, row or block.
// Search the possible matrix for any col, row or block that has only 1 empty cell
// Empty cells contain lists. Non-empty cells contain nil.
func rule1() (*linkedList, int) {
	var (
		col, row     int // position of last empty cell
		digit, count int
		matched      *linkedList
	)
	matched = &linkedList{}

	currNode := emptyL.head
	if currNode == nil {
		color.Yellow.Println("Empty list.")
	} else {
		for currNode != nil {
			row = currNode.row
			col = currNode.col
			if *debugPtr {
				color.LightGreen.Printf("cell [%d][%d]. %+v\n", row, col, currNode.vals)
			}

			if len(currNode.vals) == 1 {
				digit = currNode.vals[0]
				matched.addCell(row, col, []int{digit})
				emptyL.delNode(currNode) // remove current Node from possibility list
				mat[row][col] = digit
				mat2[row][col] = nil
				emptyCnt--
				count++
				if *debugPtr {
					color.LightGreen.Printf("Rule 1. Found lone single no. %d at cell [%d][%d]. Empty cells = %d\n", mat[row][col], row, col, emptyCnt)
				}
			}
			currNode = currNode.next
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
func rule3() (*linkedList, int) {
	const (
		scols = 3
		ncols = 9
	)

	var (
		col, row                     int // position of last empty cell
		count, itercnt               int
		notInRow, notInCol, notInBlk bool
		foundHiddenSingle            bool
		matched                      *linkedList
	)
	matched = &linkedList{}

	for dig := 1; dig <= ncols; dig++ {
		itercnt = 0
		for {
			foundHiddenSingle = false
			currNode := emptyL.head

			if currNode == nil {
				color.Yellow.Println("Empty list.")
			} else {
				for currNode != nil {
					row = currNode.row
					col = currNode.col
					if *debugPtr {
						color.LightGreen.Printf("cell [%d][%d]. %+v\n", row, col, currNode.vals)
					}

					if contains(currNode.vals, dig) {
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
									dig, row, col, row/scols, col/scols)
							}
						}

						if notInRow || notInCol || notInBlk {
							matched.addCell(row, col, []int{dig})
							emptyL.delNode(currNode) // remove current Node from possibility list
							mat[row][col] = dig      // fill in dig in resulting mat
							mat2[row][col] = nil     // blank this cell out in possibility mat
							emptyCnt--
							count++
							foundHiddenSingle = true

							if *debugPtr {
								color.LightYellow.Printf("Rule 4. Found lone single digit %d at cell [%d][%d]. Empty cells = %d\n", dig, row, col, emptyCnt)
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
									color.LightBlue.Printf("After deletion from blk [%d,%d]: %v\n", row/scols, col/scols, getBlkOfPossibleMat(mat2, row, col))
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
							fmt.Printf("Empty list count = %d\n", emptyL.countNodes())
						}
						break
					}

					currNode = currNode.next
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
