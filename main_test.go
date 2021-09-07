package main

import (
	"fmt"
	"testing"

	l "github.com/mjwong/sudoku2/lib"
	ll "github.com/mjwong/sudoku2/linkedlist"
	lp "github.com/mjwong/sudoku2/linkedlistpair"
	"gopkg.in/gookit/color.v1"
)

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

	printPossibleMat()

	currNode := emptyL.Head
	if currNode == nil {
		t.Fatalf("Empty list.")
	} else {
		i := 0
		for currNode != nil {
			if !l.IntArrayEquals(currNode.Vals, list[i]) {
				t.Fatalf("Expected %v but got %v\n", list[i], currNode.Vals)
			}
			i++
			currNode = currNode.Next
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

	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
			if !l.IntArrayEquals(mat2[i][j], list[i][j]) {
				t.Fatalf("Expected %v but got %v\n", list[i][j], mat2[i][j])
			}
		}
	}

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

	if !l.IntArrayEquals(arr1, arr2) {
		t.Fatalf("Both arrays should be equal but not. %v == %v\n", arr1, arr2)
	}

	if l.IntArrayEquals(arr2, arr3) {
		t.Fatalf("Both arrays should not be equal but are. %v != %v\n", arr2, arr3)
	}
}

func TestCountPairs(t *testing.T) {
	node1 := &lp.Pair{
		A: &ll.Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &ll.Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node2 := &lp.Pair{
		A: &ll.Cell{
			Row:  1,
			Col:  3,
			Vals: []int{2, 4},
		},
		B: &ll.Cell{
			Row:  4,
			Col:  8,
			Vals: []int{2, 4},
		},
	}

	node3 := &lp.Pair{
		A: &ll.Cell{
			Row:  1,
			Col:  2,
			Vals: []int{2, 8},
		},
		B: &ll.Cell{
			Row:  2,
			Col:  3,
			Vals: []int{2, 8},
		},
	}

	pairList := &lp.LinkedListPairs{}
	pairList.AddNode(node1)
	pairList.AddNode(node2)
	pairList.AddNode(node3)
	if pairList.CountNodes() != 3 {
		t.Fatalf("Count should be 3 but got %d.\n", pairList.CountNodes())
	}

	if pairList.Head != node1 {
		t.Fatalf("Head node should be %v but got %v.\n", node1, pairList.Head)
	}

	if pairList.Head.Next != node2 {
		t.Fatalf("Second node should be %v but got %v.\n", node2, pairList.Head.Next)
	}

	if pairList.Last != node3 {
		t.Fatalf("Last node should be %v but got %v.\n", node3, pairList.Last)
	}

	if pairList.Head.Next.Next.Next != nil {
		t.Fatalf("After last node should be nil but got %v.\n", pairList.Head.Next.Next.Next)
	}

	if *debugPtr {
		cNode := pairList.Head
		for cNode != nil {
			fmt.Printf("Pair1: row: %d col: %d. Pair2: row: %d col: %d. Next: %v.\n", cNode.A.Row, cNode.A.Col, cNode.B.Row, cNode.B.Col, cNode.Next)
			cNode = cNode.Next
		}
	}
}

func TestContainsPair(t *testing.T) {
	node1 := &lp.Pair{
		A: &ll.Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &ll.Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node2 := &lp.Pair{
		A: &ll.Cell{
			Row:  1,
			Col:  3,
			Vals: []int{2, 4},
		},
		B: &ll.Cell{
			Row:  4,
			Col:  8,
			Vals: []int{2, 4},
		},
	}

	node3 := &lp.Pair{
		A: &ll.Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &ll.Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node4 := &lp.Pair{
		A: &ll.Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
		B: &ll.Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
	}

	pairList := &lp.LinkedListPairs{}
	pairList.AddNode(node1)
	pairList.AddNode(node2)

	if !pairList.Contains(node3) {
		t.Fatalf("The pairs should be the same but are not.\n")
	}
	if !pairList.Contains(node4) { // the cells are reversed.
		t.Fatalf("The reversed pairs should be the same but are not.\n")
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

	digcnt := matched.CountNodes()

	color.LightMagenta.Printf("Found: %d digits.\n", cnt)
	if digcnt != cnt {
		t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
	}
	matched.PrintResult("Found open single")
}

func TestRule3(t *testing.T) {

	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."

	ruleTest(t, input, 3, 51)

}

// looping rule3
func TestRule3L(t *testing.T) {
	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."
	mat = populateMat(input)

	emptyCnt = countEmpty(mat)
	emptyL, mat2 = getPossibleMat(mat)

	ruleLoop(rule3, "Hidden single")
}

func TestRule_3a(t *testing.T) {

	input := "7..15..6991.37645.5.694.371..5.1.69.6.95.41.7.716.95...57.319.6.9386..1516..95..."

	ruleTest(t, input, 3, 31)
}

func TestRule_3c(t *testing.T) {

	input := "..78265.16.1395.47..5147.6.3..2.1...172.8.356...6.3..4....687..82.71.6.57..5324.."

	ruleTest(t, input, 3, 36)
}

func TestRule_3d(t *testing.T) {

	input := "4378265916813952472951478633..2.1978172.8.3569.8673124....687.282.71.6.57..532489"

	ruleTest(t, input, 3, 16)
}

func TestRule_3e(t *testing.T) {

	input := "43782659168139524729514786336.251978172.893569586731245.396871282971.635716532489"

	ruleTest(t, input, 3, 4)

}

// Rule 5: Naked Pairs
func TestRule5(t *testing.T) {

	input := "142.73...597.462.3863.52...31852469772639.4.545976.32.6.54391.293128....2.461..39"

	ruleTest(t, input, 5, 23)

}

// Rule 8: Hidden Pairs
func TestRule8(t *testing.T) {

	input := "43782659168139524729514786336.251978172.893569586731245.396871282971.635716532489"

	ruleTest(t, input, 3, 4)

}

func ruleTest(t *testing.T, input string, rule, empCnt int) {
	var (
		desc string
	)
	mat = populateMat(input)

	emptyCnt = countEmpty(mat)
	if emptyCnt != empCnt {
		t.Fatalf("Expected %d but got %d.\n", empCnt, emptyCnt)
	}

	emptyL, mat2 = getPossibleMat(mat)
	if emptyL.CountNodes() != empCnt {
		t.Fatalf("Expected %d nodes in empty list but got %d.\n", empCnt, emptyL.CountNodes())
	}

	printPossibleMat()
	printSudoku(mat)
	fmt.Printf("Starting empty cells = %d\n", emptyCnt)

	switch rule {
	case 3:
		desc = "Hidden single"
		matched, cnt := rule3()
		digcnt := matched.CountNodes()

		printPossibleMat()

		color.LightMagenta.Printf("Found: %d digits.\n", cnt)
		if digcnt != cnt {
			t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
		}
		matched.PrintResult(desc)
	case 5:
		desc = "Naked pairs"
		matched, cnt := rule5()
		fmt.Printf("Found: %s = %d.\n", desc, cnt)
		matched.PrintResult(desc)
		printPossibleMat()
	case 8:
		desc = "Hidden pairs"
		matched, cnt := rule3()
		fmt.Printf("Found: %s = %d.\n", desc, cnt)
		matched.PrintResult(desc)
	}

	printSudoku(mat)

	if emptyCnt == 0 {
		color.Magenta.Println("Finished!")
	} else {
		fmt.Printf("Empty cells = %d\n", emptyCnt)
	}
}
