package main

import (
	"fmt"
	"testing"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
	"gopkg.in/gookit/color.v1"
)

var (
	difficult1 = "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..." // difficult1.txt
	//difficult2 = "...826..1..1....47..5.4....3....1....72.8.35....6....4....6.7..82....6..7..532..." // difficult2.txt
	difficult3 = "14...3.......4...38.3.52.......2..977.6.9.4.545..6.......43.1.29...8.......6...39" // difficult3.txt
	difficult4 = "....92....7..853.93...7.8..2...61.4..6.....7..9.82...1..8.5...79.271..3....43...." // difficult4.txt
	difficult5 = ".4...8...7.....8...3..16..49...6.38..6..3..9..23.4...64..12..3...2.....5...3...1." // difficult5.txt
)

func TestEmptyCount(t *testing.T) {
	mat := PopulateMat(difficult1)
	emptyCnt := CountEmpty(mat)
	if emptyCnt != 51 {
		t.Fatalf("Expected 51 but got %d.\n", emptyCnt)
	}
}

func TestGetPossibleMat(t *testing.T) {
	list := [][]int{
		{2, 3, 4, 7, 8},
		{2, 3, 4, 8},
		{2, 4, 8},
		{2, 3, 8, 9},
		{2, 7, 8},
		{2, 6, 7, 8},
		{2, 8, 9},
		{2, 8},
		{2, 3},
		{2, 5, 8},
		{2, 8},
		{2, 8},
		{2, 9},
		{2, 8, 9},
		{2, 7, 8},
		{1, 2, 8, 9},
		{2, 3, 4, 8},
		{2, 3, 4, 5, 8},
		{2, 4, 5, 8},
		{2, 3, 7},
		{1, 2, 8},
		{1, 2, 3, 7, 8},
		{1, 2, 3, 4, 8},
		{2, 3, 8, 9},
		{2, 8, 9},
		{1, 2, 8},
		{1, 2, 8},
		{2, 3, 8},
		{2, 3, 4, 8},
		{2, 3, 6, 9},
		{2, 8},
		{2, 3, 8, 9},
		{2, 5, 8},
		{2, 3, 4, 5, 8},
		{2, 3, 4, 8},
		{1, 2, 4, 8},
		{2, 4, 5, 8},
		{2, 4},
		{1, 2},
		{2, 4, 8},
		{2, 4},
		{2, 4, 9},
		{2, 7},
		{2, 7},
		{1, 2, 4, 8},
		{2, 4, 6, 8},
		{2, 4, 8},
		{2, 4, 7},
		{2, 7, 8},
		{2, 3, 4, 7, 8},
		{2, 3, 4, 8},
	}

	mat := PopulateMat(difficult1)
	emptyL, mat2 = GetPossibleMat(mat)
	PrintPossibleMat(mat2)

	currNode := emptyL.Head
	if currNode == nil {
		t.Fatalf("Empty list.")
	} else {
		i := 0
		for currNode != nil {
			if !IntArrayEquals(currNode.Vals, list[i]) {
				t.Fatalf("Expected %v but got %v\n", list[i], currNode.Vals)
			}
			i++
			currNode = currNode.Next
		}
	}
}

func TestDigitNotIn(t *testing.T) {
	debug := DebugFn(2) // Get name of this caller
	mat := PopulateMat(difficult1)
	emptyL, mat2 = GetPossibleMat(mat)

	// possibility matrix
	list := Pmat{
		{[]int{2, 3, 4, 7, 8}, []int{2, 3, 4, 8}, []int{2, 4, 8}, []int{}, []int{}, []int{2, 3, 8, 9}, []int{2, 7, 8}, []int{2, 6, 7, 8}, []int{2, 8, 9}},
		{[]int{}, []int{}, []int{2, 8}, []int{2, 3}, []int{}, []int{}, []int{}, []int{2, 5, 8}, []int{2, 8}},
		{[]int{}, []int{2, 8}, []int{}, []int{2, 9}, []int{}, []int{2, 8, 9}, []int{}, []int{2, 7, 8}, []int{}},
		{[]int{2, 3, 4, 8}, []int{2, 3, 4, 5, 8}, []int{2, 4, 5, 8}, []int{2, 3, 7}, []int{1, 2, 8}, []int{1, 2, 3, 7, 8}, []int{}, []int{}, []int{2, 3, 4, 8}},
		{[]int{}, []int{2, 3, 8, 9}, []int{2, 8, 9}, []int{}, []int{2, 8}, []int{}, []int{}, []int{2, 3, 8}, []int{}},
		{[]int{2, 3, 4, 8}, []int{}, []int{}, []int{2, 3, 6, 9}, []int{2, 8}, []int{2, 3, 8, 9}, []int{2, 5, 8}, []int{2, 3, 4, 5, 8}, []int{2, 3, 4, 8}},
		{[]int{2, 4, 8}, []int{2, 4, 5, 8}, []int{}, []int{2, 4}, []int{}, []int{1, 2}, []int{}, []int{2, 4, 8}, []int{}},
		{[]int{2, 4}, []int{2, 4, 9}, []int{}, []int{}, []int{}, []int{2, 7}, []int{2, 7}, []int{}, []int{}},
		{[]int{}, []int{2, 4, 6, 8}, []int{2, 4, 8}, []int{2, 4, 7}, []int{}, []int{}, []int{2, 7, 8}, []int{2, 3, 4, 7, 8}, []int{2, 3, 4, 8}},
	}

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			if !IntArrayEquals(mat2[i][j], list[i][j]) {
				t.Fatalf("Expected %v but got %v\n", list[i][j], mat2[i][j])
			}
		}
	}

	if len(difficult1) != 81 {
		t.Fatalf("Input matrix len not 81, got %d.\n", len(difficult1))
	}

	// digit 1 is found hidden in cell [3,4].
	// search for digit 1; should find in row 4 and blk [1,1] but not in col 5.
	if !FindDigitInRow(debug, mat2, 3, 4, 1) {
		t.Fatalf("Digit should be in row 3.")
	}

	if FindDigitInCol(debug, mat2, 3, 4, 1) {
		t.Fatalf("Digit should not be in column 4.")
	}

	if !FindDigitInBlk(debug, mat2, 3, 4, 1) {
		t.Fatalf("Digit should be in block [1,1].")
	}
}

func TestDelNode(t *testing.T) {
	input := ".341528699.837645252.948371245.136986895.413737168.524857231.464938672.516249578."

	PrepPmat(input)

	currentNode := emptyL.Head

	// Fill in digit 7 in [1,1]
	mat[1][1] = 7
	mat2[1][1] = nil

	emptyL.DelNode(currentNode)

	if emptyL.CountNodes() != 8 {
		t.Fatalf("Empty list count should be 8 but got %d.\n", emptyL.CountNodes())
	}
}

func TestIntArrEq(t *testing.T) {
	arr1 := []int{1, 2, 3}
	arr2 := []int{1, 2, 3}
	arr3 := []int{2, 3, 4}

	if !IntArrayEquals(arr1, arr2) {
		t.Fatalf("Both arrays should be equal but not. %v == %v\n", arr1, arr2)
	}

	if IntArrayEquals(arr2, arr3) {
		t.Fatalf("Both arrays should not be equal but are. %v != %v\n", arr2, arr3)
	}
}

// ************************************** Rule Tests *********************************************

// Rule 8: Hidden Pairs
func TestRule8(t *testing.T) {

	input := "43782659168139524729514786336.251978172.893569586731245.396871282971.635716532489"

	ruleTest(t, input, 3, 4, 0)

}

func ruleTest(t *testing.T, input string, rule, empCnt, numFound int) {
	var (
		count int
		desc  string
	)
	PrepPmat(input)

	if emptyCnt != empCnt {
		t.Fatalf("Expected %d but got %d.\n", empCnt, emptyCnt)
	}
	if emptyL.CountNodes() != empCnt {
		t.Fatalf("Expected %d nodes in empty list but got %d.\n", empCnt, emptyL.CountNodes())
	}

	PrintPossibleMat(mat2)
	PrintSudoku(mat)
	fmt.Printf("Starting empty cells = %d\n", emptyCnt)

	switch rule {
	case 3:
		desc = RuleTable[3]
		matched, cnt, elapsed := rule3()
		digcnt := matched.CountNodes()

		PrintPossibleMat(mat2)

		color.LightMagenta.Printf("Found: %d digits. Elapsed time = %v ms\n", cnt, elapsed.Milliseconds())
		if digcnt != cnt {
			t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
		}
		matched.PrintResult(desc)
		count = cnt

	case 5:
		desc = RuleTable[5]
		matched, cnt, elapsed := rule5()
		fmt.Printf("Found: %s = %d. Elapsed time = %v ms\n", desc, cnt, elapsed.Milliseconds())
		matched.PrintResult(desc)
		count = cnt

	case 8:
		desc = RuleTable[8]
		matched, cnt, elapsed := rule8()
		fmt.Printf("Found: %s = %d. Elapsed time = %v ms\n", desc, cnt, elapsed.Milliseconds())
		matched.PrintResult(desc)
	}

	if numFound != 0 && numFound != count {
		t.Fatalf("Expected to find %d but got %d counts.\n", numFound, count)
	}

	PrintPossibleMat(mat2)
	PrintSudoku(mat)

	if emptyCnt == 0 {
		color.Magenta.Println("Finished!")
	} else {
		fmt.Printf("Empty cells = %d\n", emptyCnt)
	}
}
