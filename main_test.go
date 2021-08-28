package main

import (
	"fmt"
	"testing"

	"gopkg.in/gookit/color.v1"
)

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

func TestEmptyCount(t *testing.T) {

	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."

	mat := populateMat(input)
	emptyCnt := countEmpty(mat)

	if emptyCnt != 51 {
		t.Fatalf("Expected 51 but got %d.\n", emptyCnt)
	}
}

func TestGetPossibleMat(t *testing.T) {

	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."
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

	mat := populateMat(input)

	emptyL, mat2 = getPossibleMat(mat)

	fmt.Printf("%v\n", mat2)

	currNode := emptyL.head
	if currNode == nil {
		t.Fatalf("Empty list.")
	} else {
		i := 0
		for currNode != nil {
			if !IntArrayEquals(currNode.vals, list[i]) {
				t.Fatalf("Expected %v but got %v\n", list[i], currNode.vals)
			}
			i++
			currNode = currNode.next
		}
	}
}

func TestDigitNotIn(t *testing.T) {
	const ncols = 9

	input := "...15....91..764..5.6.4.3.1......69.6..5.41.7.71........7.3.9.6..386..151...95..."

	mat := populateMat(input)

	emptyL, mat2 = getPossibleMat(mat)

	// possibility matrix
	list := pmat{
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

	fmt.Printf("%v\n", mat2)

	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
			if !IntArrayEquals(mat2[i][j], list[i][j]) {
				t.Fatalf("Expected %v but got %v\n", list[i][j], mat2[i][j])
			}
		}
	}

	printPossibleMat()

	if len(input) != 81 {
		t.Fatalf("Input matrix len not 81, got %d.\n", len(input))
	}

	// digit 1 is found hidden in cell [3,4].
	// search for digit 1; should find in row 4 and blk [1,1] but not in col 5.
	if !findDigitInRow(mat2, 3, 4, 1) {
		t.Fatalf("Digit should be in row 3.")
	}

	if findDigitInCol(mat2, 3, 4, 1) {
		t.Fatalf("Digit should not be in column 4.")
	}

	if !findDigitInBlk(mat2, 3, 4, 1) {
		t.Fatalf("Digit should be in block [1,1].")
	}
}

func TestDelNode(t *testing.T) {
	input := ".341528699.837645252.948371245.136986895.413737168.524857231.464938672.516249578."

	mat = populateMat(input)
	emptyCnt = countEmpty(mat)

	emptyL, mat2 = getPossibleMat(mat)

	currentNode := emptyL.head

	// Fill in digit 7 in [1,1]
	mat[1][1] = 7
	mat2[1][1] = nil

	emptyL.delNode(currentNode)

	if emptyL.countNodes() != 8 {
		t.Fatalf("Empty list count should be 8 but got %d.\n", emptyL.countNodes())
	}
}

func TestRule1(t *testing.T) {

	input := ".341528699.837645252.948371245.136986895.413737168.524857231.464938672.516249578."

	mat = populateMat(input)
	emptyCnt = countEmpty(mat)

	emptyL, mat2 = getPossibleMat(mat)

	matched, cnt := rule1()
	if cnt != 9 {
		t.Fatalf("Expected 8 but got %d\n", cnt)
	}

	if !checkSums(mat) {
		t.Fatalf("There are errors in the resulting matrix.\n")
	}

	digcnt := matched.countNodes()

	color.LightMagenta.Printf("Found: %d digits.\n", cnt)
	if digcnt != cnt {
		t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
	}
	matched.printResult("Found open single")
}

func TestRule3(t *testing.T) {

	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."

	mat = populateMat(input)

	emptyCnt = countEmpty(mat)

	emptyL, mat2 = getPossibleMat(mat)

	for {
		matched, cnt := rule3()

		digcnt := matched.countNodes()

		printPossibleMat()
		color.LightMagenta.Printf("Found: %d digits.\n", cnt)
		if digcnt != cnt {
			t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
		}

		matched.printResult("Found hidden single")

		if cnt <= 0 {
			break
		}
	}
}

func TestRule_3a(t *testing.T) {

	input := "7..15..6991.37645.5.694.371..5.1.69.6.95.41.7.716.95...57.319.6.9386..1516..95..."

	mat = populateMat(input)

	emptyCnt = countEmpty(mat)
	if emptyCnt != 31 {
		t.Fatalf("Expected 31 but got %d.\n", emptyCnt)
	}

	emptyL, mat2 = getPossibleMat(mat)
	if emptyL.countNodes() != 31 {
		t.Fatalf("Expected 31 nodes in empty list but got %d.\n", emptyL.countNodes())
	}

	printPossibleMat()
	printSudoku(mat)

	matched, cnt := rule3()

	digcnt := matched.countNodes()

	printPossibleMat()

	color.LightMagenta.Printf("Found: %d digits.\n", cnt)
	if digcnt != cnt {
		t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
	}

	matched.printResult("Found hidden single")
	printSudoku(mat)

	if emptyCnt == 0 {
		color.Magenta.Println("Finished!")
	} else {
		fmt.Printf("Empty cells = %d\n", emptyCnt)
	}
}
