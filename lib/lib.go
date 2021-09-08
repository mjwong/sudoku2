package lib

import (
	"reflect"
	"runtime"
	"strconv"
	"strings"
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
