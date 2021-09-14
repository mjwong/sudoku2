package lib

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"github.com/gookit/color"
)

const (
	N     = 9
	SQ    = 3
	nsize = 81
)

type (
	Intmat    [N][N]int
	Sqmat     [SQ][SQ][]int
	Pmat      [N][N][]int
	MatString string
	Coord     struct {
		Row, Col int
	}
)

func Contains(arr []int, v int) bool {
	for _, a := range arr {
		if a == v {
			return true
		}
	}
	return false
}

func ContainsMulti(arr []int, v []int) bool {
	for _, a := range arr {
		for _, e := range v {
			if a == e {
				return true
			}
		}
	}
	return false
}

func IntArrayEquals(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func GetArrForRow(m Intmat, i int) []int {
	var arrRow []int

	for j := 0; j < N; j++ {
		if m[i][j] != 0 {
			arrRow = append(arrRow, m[i][j])
		}
	}
	return arrRow
}

func GetArrForCol(m Intmat, j int) []int {
	var arrCol []int

	for i := 0; i < N; i++ {
		if m[i][j] != 0 {
			arrCol = append(arrCol, m[i][j])
		}
	}
	return arrCol
}

func GetArrForSqu(m Intmat, i, j int) []int {
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

func GetColOfPossibleMat(mat2 Pmat, col int) [][]int {
	var m [][]int

	for i := 0; i < N; i++ {
		m = append(m, mat2[i][col])
	}
	return m
}

func GetBlkOfPossibleMat(mat2 Pmat, row, col int) [][]int {
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

func EraseFromSlice(sl []int, v int) []int {
	for i, a := range sl {
		if a == v {
			return Remove(sl, i)
		}
	}
	return sl
}

func EraseMultiFromSlice(sl []int, v []int) []int {
	for i, a := range sl {
		for e := range v {
			if a == e {
				return Remove(sl, i)
			}
		}

	}
	return sl
}

func Remove(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
}

func PopulateMat(s string) Intmat {
	var (
		x, y int
		mat  Intmat
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

func MatToString(m Intmat) string {
	var c, s string

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if m[i][j] == 0 {
				c = "."
			} else {
				c = strconv.Itoa(m[i][j])
			}
			s += c
		}
	}
	return s
}

func CountEmpty(mat Intmat) int {
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

// Get calling function name
func FuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}

// Get function name, e.g. main.foo
func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func CountElemPosMat(m Pmat) int {
	var count int
	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if m[i][j] != nil {
				count += len(m[i][j])
			}
		}
	}
	return count
}

func ReadInput() string {
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

// ****************************************** start of find/erase fns ******************************************

func InRow(m Intmat, row, num int) bool {
	for c := 0; c < N; c++ {
		if m[row][c] == num {
			return true
		}
	}
	return false
}

func InCol(m Intmat, col, num int) bool {
	for r := 0; r < N; r++ {
		if m[r][col] == num {
			return true
		}
	}
	return false
}

func InSqu(m Intmat, row, col, num int) bool {
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

func IsSafe(m Intmat, row, col, num int) bool {
	if !InRow(m, row, num) && !InCol(m, col, num) && !InSqu(m, row, col, num) {
		return true
	}
	return false
}

// check the row, col and square sums
func CheckSums(m Intmat) bool {
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

func PrintPossibleMat(m Pmat) {
	fmt.Println("-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------")

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			fmt.Printf("|%-20v ", arr2String(m[i][j], ","))
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

func PrintSudoku(m Intmat) {
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

// *******************************************************************************************************
// *                                          start of find funcs                                          *
// *******************************************************************************************************

// check row of possibility matrix
func FindDigitInRow(debugPtr *bool, mat2 Pmat, row, col, dig int) bool {
	if *debugPtr {
		fmt.Println("In findDigitInRow...")
	}

	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col {
			if *debugPtr {
				fmt.Printf("Cell [%d][%d] = %v\n", row, c, mat2[row][c])
			}

			if Contains(mat2[row][c], dig) {

				if *debugPtr {
					fmt.Printf("row %d contains digit %d\n", row, dig)
				}
				return true

			}
		}
	}

	return false
}

// check column of possibility matrix
func FindDigitInCol(debugPtr *bool, mat2 Pmat, row, col, dig int) bool {
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
					fmt.Printf("col %d contains digit %d\n", col, dig)
				}
				return true
			}
		}
	}

	return false
}

// check block of possibility matrix corresponding to cell [row[col]
func FindDigitInBlk(debugPtr *bool, mat2 Pmat, row, col, dig int) bool {
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

				if Contains(mat2[x][y], dig) {

					if *debugPtr {
						fmt.Printf("blk [%d,%d] contains digit %d\n", x/SQ, y/SQ, dig)
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

// *******************************************************************************************************
// *                                          start of naked pairs funcs                                         *
// *******************************************************************************************************

// check row of possibility matrix
func FindDigitInRowPair(debugPtr *bool, mat2 Pmat, row, col, col2 int, digits []int) bool {
	for c := 0; c < N; c++ {
		if mat2[row][c] != nil && c != col && c != col2 {
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

// check col of possibility matrix if any of the digits in the naked pair is
// found in this col
func FindDigitInColPair(debugPtr *bool, mat2 Pmat, row, col, row2 int, digits []int) bool {
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

// check block of possibility matrix
func FindDigitInBlkPair(debugPtr *bool, mat2 Pmat, row, col, row2, col2 int, digits []int) bool {
	startRow := row / SQ * SQ
	startCol := col / SQ * SQ

	for x := startRow; x < startRow+SQ; x++ {
		for y := startCol; y < startCol+SQ; y++ {
			if mat2[x][y] != nil && !(x == row && y == col) && !(x == row2 && y == col2) {
				if *debugPtr {
					fmt.Printf("Cell [%d][%d] = %v\n", x, y, mat2[x][y])
				}
				if ContainsMulti(mat2[x][y], digits) {

					if *debugPtr {
						fmt.Printf("Cell [%d,%d] contains digits of naked pair", x, y)
					}
					return true
				}
			}
		}
	}

	return false
}
