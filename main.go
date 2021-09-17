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
