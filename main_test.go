package main

import (
	"fmt"
	"testing"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
	lp "github.com/mjwong/sudoku2/linkedlistpair"
	. "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

func TestEmptyCount(t *testing.T) {

	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."

	mat := PopulateMat(input)
	emptyCnt := CountEmpty(mat)

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

	mat := PopulateMat(input)

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
	const ncols = 9

	input := "...15....91..764..5.6.4.3.1......69.6..5.41.7.71........7.3.9.6..386..151...95..."

	mat := PopulateMat(input)

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

	for i := 0; i < ncols; i++ {
		for j := 0; j < ncols; j++ {
			if !IntArrayEquals(mat2[i][j], list[i][j]) {
				t.Fatalf("Expected %v but got %v\n", list[i][j], mat2[i][j])
			}
		}
	}

	if len(input) != 81 {
		t.Fatalf("Input matrix len not 81, got %d.\n", len(input))
	}

	// digit 1 is found hidden in cell [3,4].
	// search for digit 1; should find in row 4 and blk [1,1] but not in col 5.
	if !FindDigitInRow(debugPtr, mat2, 3, 4, 1) {
		t.Fatalf("Digit should be in row 3.")
	}

	if FindDigitInCol(debugPtr, mat2, 3, 4, 1) {
		t.Fatalf("Digit should not be in column 4.")
	}

	if !FindDigitInBlk(debugPtr, mat2, 3, 4, 1) {
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

func TestCountPairs(t *testing.T) {
	node1 := &lp.Pair{
		A: &Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node2 := &lp.Pair{
		A: &Cell{
			Row:  1,
			Col:  3,
			Vals: []int{2, 4},
		},
		B: &Cell{
			Row:  4,
			Col:  8,
			Vals: []int{2, 4},
		},
	}

	node3 := &lp.Pair{
		A: &Cell{
			Row:  1,
			Col:  2,
			Vals: []int{2, 8},
		},
		B: &Cell{
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
		A: &Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node2 := &lp.Pair{
		A: &Cell{
			Row:  1,
			Col:  3,
			Vals: []int{2, 4},
		},
		B: &Cell{
			Row:  4,
			Col:  8,
			Vals: []int{2, 4},
		},
	}

	node3 := &lp.Pair{
		A: &Cell{
			Row:  4,
			Col:  5,
			Vals: []int{2, 8},
		},
		B: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
	}

	node4 := &lp.Pair{
		A: &Cell{
			Row:  6,
			Col:  7,
			Vals: []int{2, 8},
		},
		B: &Cell{
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
	PrepPmat(input)

	matched, cnt := rule1()
	if cnt != 9 {
		t.Fatalf("Expected 8 but got %d\n", cnt)
	}

	if !CheckSums(mat) {
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
	// difficult1.txt
	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."

	ruleTest(t, input, 3, 51, 20)

}

// looping rule3
func TestRule3L(t *testing.T) {
	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."
	PrepPmat(input)

	totCnt := RuleLoop(rule3, RuleTable[3], Zero)

	if totCnt != 51 {
		t.Fatalf("Expected to find 51 but got %d counts.\n", totCnt)
	}
}

func TestRule_3a(t *testing.T) {

	input := "7..15..6991.37645.5.694.371..5.1.69.6.95.41.7.716.95...57.319.6.9386..1516..95..."

	ruleTest(t, input, 3, 31, 0)
}

func TestRule_3c(t *testing.T) {

	input := "..78265.16.1395.47..5147.6.3..2.1...172.8.356...6.3..4....687..82.71.6.57..5324.."

	ruleTest(t, input, 3, 36, 0)
}

func TestRule_3d(t *testing.T) {

	input := "4378265916813952472951478633..2.1978172.8.3569.8673124....687.282.71.6.57..532489"

	ruleTest(t, input, 3, 16, 0)
}

func TestRule_3e(t *testing.T) {

	input := "43782659168139524729514786336.251978172.893569586731245.396871282971.635716532489"

	ruleTest(t, input, 3, 4, 0)
}

func TestRule_3f(t *testing.T) {
	// test difficult3.txt

	input := "14...3.......4...38.3.52.......2..977.6.9.4.545..6.......43.1.29...8.......6...39"

	PrepPmat(input)

	totCnt := RuleLoop(rule3, RuleTable[3], Zero)

	if totCnt != 29 {
		t.Fatalf("Expected to find 29 but got %d counts.\n", totCnt)
	}
}

func TestRule_3g(t *testing.T) {
	// test difficult4.txt

	input := "....92....7..853.93...7.8..2...61.4..6.....7..9.82...1..8.5...79.271..3....43...."

	PrepPmat(input)

	totCnt := RuleLoop(rule3, RuleTable[3], Zero)

	if totCnt != 51 {
		t.Fatalf("Expected to find 51 but got %d counts.\n", totCnt)
	}
}

// Rule 5: Naked Pairs
func TestRule5(t *testing.T) {

	input := "142.73...597.462.3863.52...31852469772639.4.545976.32.6.54391.293128....2.461..39"

	ruleTest(t, input, 5, 23, 9)

}

func TestRule5a(t *testing.T) {
	// difficult1.txt
	// naked pair (2,8) found in starting possibility matrix at cells [2,1] and [1,2] of block [0,0].
	input := "...15....91..764..5.6.4.3........69.6..5.4..7.71........7.3.9.6..386..15....95..."

	ruleTest(t, input, 5, 51, 1)
}

func TestRule3n5(t *testing.T) {
	// test difficult3.txt

	input := "14...3.......4...38.3.52.......2..977.6.9.4.545..6.......43.1.29...8.......6...39"

	PrepPmat(input)

	totCnt := RuleLoop(rule3, RuleTable[3], Zero)

	if totCnt != 29 {
		t.Fatalf("Expected to find 29 but got %d counts.\n", totCnt)
	}

	fmt.Println("Starting possible matrix for Rule 5.")
	PrintPossibleMat(mat2)
	input = "142.73...597.462.3863.52...31852469772639.4.545976.32.6.54391.293128....2.461..39"
	ruleTest(t, input, 5, 23, 9)

	cnt := RuleLoop(rule1, RuleTable[1], Zero)
	if cnt != 23 {
		t.Fatalf("Expected 23 but got %d\n", cnt)
	}

	if !CheckSums(mat) {
		t.Fatal("Expected to be solved")
		PrintSudoku(mat)
	}
}

func TestRule135(t *testing.T) {
	// difficult5.txt
	ruleCnt := map[int]int{}

	input := ".4...8...7.....8...3..16..49...6.38..6..3..9..23.4...64..12..3...2.....5...3...1."
	PrepPmat(input)

	startCnt := emptyL.CountNodes()
	if startCnt != 54 {
		t.Fatalf("Expected to find 54 but got %d counts.\n", startCnt)
	}
	fmt.Printf("Starting empty count: %d\n", startCnt)

	ruleCnt[1] = RuleLoop(rule1, RuleTable[1], Zero)
	ruleCnt[3] = RuleLoop(rule3, RuleTable[3], Zero)
	cntBefore := emptyL.CountElem()
	ruleCnt[5] = RuleLoop(rule5, RuleTable[5], SameCnt)
	cntAfter := emptyL.CountElem()

	if ruleCnt[3] != 11 {
		t.Fatalf("Expect to find 11 hidden singles but got %d.\n", ruleCnt[3])
	}

	PrintFound([]int{1, 3, 5}, ruleCnt)
	fmt.Printf("Count before and after Rule 5: %d, %d.\n", cntBefore, cntAfter)
	fmt.Printf("Empty cells : %2d\n", emptyL.CountNodes())

	if emptyL.CountNodes() == 0 {
		if !CheckSums(mat) {
			t.Fatal("Expected to be solved")
			PrintSudoku(mat)
		}
	}
}

func TestRule_135a(t *testing.T) {
	// mid-way through difficult5.txt

	input := ".4...8..37..4.382..3..16..49.4.6.38..6..3.49..23.4...64..12.63.31268..45...3.4.1."
	PrepPmat(input)

	startCnt := emptyL.CountNodes()
	if startCnt != 41 {
		t.Fatalf("Expected to find 41 but got %d counts.\n", startCnt)
	}
	fmt.Printf("Starting empty count: %d\n", startCnt)

	cnt1 := RuleLoop(rule1, RuleTable[1], Zero)
	cnt3 := RuleLoop(rule3, RuleTable[3], Zero)
	cntBefore := emptyL.CountElem()
	cnt5 := RuleLoop(rule5, RuleTable[5], SameCnt)
	cntAfter := emptyL.CountElem()

	if cnt3 != 41 {
		t.Fatalf("Expect to find 41 hidden singles but got %d.\n", cnt3)
	}

	fmt.Printf("Rule 1: Found %2d open singles\n", cnt1)
	fmt.Printf("Rule 3: Found %2d hidden singles\n", cnt3)
	fmt.Printf("Rule 5: Found %2d naked pairs\n", cnt5)
	fmt.Printf("Count before and after Rule 5: %d, %d.\n", cntBefore, cntAfter)
	fmt.Printf("Empty cells : %2d\n", emptyL.CountNodes())

	if emptyL.CountNodes() == 0 {
		if !CheckSums(mat) {
			t.Fatal("Expected to be solved")
			PrintSudoku(mat)
		}
	}
}

func TestRule_135c(t *testing.T) {
	// mid-way through difficult5.txt
	// removed digit 2 from position [1,7] and digit 6 from [6,6]
	// digit 2 at [1,7] requires X-wing/Rectangular rule to solve.
	// Edges of Rectangular are [0,0], [0,3], [2,0] and [2,3].
	// Therefore, in block [0,2], the 2 can only be in row 1 and hence cell [1,7]
	// because the edge cells of the rectangle occupies row 1 and row 3.
	//                        *
	input := ".4...8..37..4.38...3..16..49.4.6.38..6..3.49..23.4...64..12..3.31268..45...3.4.1."
	PrepPmat(input)

	startCnt := emptyL.CountNodes()
	if startCnt != 43 {
		t.Fatalf("Expected to find 43 but got %d counts.\n", startCnt)
	}
	fmt.Printf("Starting empty count: %d\n", startCnt)

	totCnt1 := 0
	totCnt3 := 0
	totCnt5 := 0

	for {
		cnt1 := RuleLoop(rule1, RuleTable[1], Zero)
		cnt3 := RuleLoop(rule3, RuleTable[3], Zero)
		cntBefore := emptyL.CountElem()
		cnt5 := RuleLoop(rule5, RuleTable[5], SameCnt)
		cntAfter := emptyL.CountElem()

		totCnt1 += cnt1
		totCnt3 += cnt3
		totCnt5 += cnt5

		if cnt1 == 0 && cnt3 == 0 && cntBefore == cntAfter {
			fmt.Printf("Count before and after Rule 5: %d, %d.\n", cntBefore, cntAfter)
			break
		}
	}

	/*
		if cnt3 != 40 {
			t.Fatalf("Expect to find 40 hidden singles but got %d.\n", cnt3)
		}
	*/
	fmt.Printf("Rule 1: Found %2d open singles\n", totCnt1)
	fmt.Printf("Rule 3: Found %2d hidden singles\n", totCnt3)
	fmt.Printf("Rule 5: Found %2d naked pairs\n", totCnt5)
	fmt.Printf("Empty cells : %2d\n", emptyL.CountNodes())

	if emptyL.CountNodes() == 0 {
		if !CheckSums(mat) {
			t.Fatal("Expected to be solved")
			PrintSudoku(mat)
		}
	}
}

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
		matched, cnt := rule3()
		digcnt := matched.CountNodes()

		PrintPossibleMat(mat2)

		color.LightMagenta.Printf("Found: %d digits.\n", cnt)
		if digcnt != cnt {
			t.Fatalf("Expected no. of digits found is %d but got %d", cnt, digcnt)
		}
		matched.PrintResult(desc)
		count = cnt
	case 5:
		desc = RuleTable[5]
		matched, cnt := rule5()
		fmt.Printf("Found: %s = %d.\n", desc, cnt)
		matched.PrintResult(desc)
		count = cnt
	case 8:
		desc = RuleTable[8]
		matched, cnt := rule3()
		fmt.Printf("Found: %s = %d.\n", desc, cnt)
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

func TestCheckBlkForDigit(t *testing.T) {

	pm := Pmat{}
	pm[0][0] = []int{1, 2, 5, 6}
	pm[0][2] = []int{1, 5, 6, 9}
	pm[1][1] = []int{5, 9}
	pm[1][2] = []int{1, 5, 6, 9}
	pm[2][0] = []int{2, 5, 8}
	pm[2][2] = []int{5, 8, 9}

	arr, inBlk := checkBlkForDigit(pm, 0, 0, 2, 2)
	if !inBlk {
		t.Fatalf("Should find 2 same digits in blk but found %d.\n", len(arr))
	}

	arr, inBlk = checkBlkForDigit(pm, 0, 0, 5, 2)
	if inBlk {
		t.Fatalf("Should not find 2 same digits in blk but found %d.\n", len(arr))
	}
}

func TestEraseDigitFromRowMulti(t *testing.T) {

	mat2 = Pmat{}
	mat2[0][0] = []int{1, 2, 5, 6}
	mat2[0][2] = []int{1, 5, 6, 9}
	mat2[0][3] = []int{2, 5, 7, 9}
	mat2[0][4] = []int{5, 7, 9}
	mat2[0][6] = []int{1, 2, 5, 6, 7, 9}
	mat2[0][7] = []int{2, 5, 6, 7}

	emptyL = CreatelinkedList()
	emptyL.AddCell(0, 6, []int{1, 2, 5, 6, 7, 9})
	emptyL.AddCell(0, 7, []int{2, 5, 6, 7})

	PrintPossibleMat(mat2)

	startCnt := CountElemPosMat(mat2)

	eraCnt, erased := eraseDigitFromRowMulti(0, 2, []int{0, 3})

	endCnt := CountElemPosMat(mat2)

	if !erased {
		t.Fatal("Should be erased but not.\n")
	}

	if eraCnt != 2 {
		t.Fatalf("2 counts of digit 2 should be erased but got %d.\n", eraCnt)
	}

	if Contains(mat2[0][6], 2) {
		t.Fatal("Possibility matrix cell [0,6] should not contain 2.\n")
	}

	if Contains(mat2[0][7], 2) {
		t.Fatal("Possibility matrix cell [0,7] should not contain 2.\n")
	}

	if (startCnt - endCnt) != 2 {
		t.Fatalf("Should have erased 2 values but got %d.\n", startCnt-endCnt)
	}
}

func TestContainsXwing(t *testing.T) {

	arr := []Idx{
		Idx{Row: 1, Col: 2, Vals: []int{1, 2, 5}},
		Idx{Row: 2, Col: 3, Vals: []int{2, 3, 6}},
		Idx{Row: 3, Col: 4, Vals: []int{2, 4, 7}},
		Idx{Row: 4, Col: 5, Vals: []int{3, 5, 9}},
	}

	arr2 := []Idx{
		Idx{Row: 5, Col: 6, Vals: []int{1, 2, 5}},
		Idx{Row: 6, Col: 8, Vals: []int{2, 3, 6}},
		Idx{Row: 3, Col: 4, Vals: []int{2, 4, 7}},
		Idx{Row: 4, Col: 5, Vals: []int{3, 5, 9}},
	}

	arr3 := []Idx{
		Idx{Row: 1, Col: 2, Vals: []int{1, 2, 5}},
		Idx{Row: 2, Col: 3, Vals: []int{2, 3, 6}},
		Idx{Row: 3, Col: 4, Vals: []int{2, 4, 7}},
		Idx{Row: 4, Col: 5, Vals: []int{3, 5, 9}},
	}

	arr4 := []Idx{
		Idx{Row: 6, Col: 2, Vals: []int{1, 2, 5}},
		Idx{Row: 7, Col: 9, Vals: []int{2, 3, 6}},
		Idx{Row: 1, Col: 3, Vals: []int{2, 4, 7}},
		Idx{Row: 2, Col: 7, Vals: []int{3, 5, 9}},
	}

	ml := &Matchlist{}
	ml.AddRNode(arr)
	ml.AddRNode(arr2)

	currNode := ml.Head
	for currNode != nil {
		if !ml.ContainsXwing(arr3) {
			t.Fatal("Should contain a X-wing but not.")
		}

		if ml.ContainsXwing(arr4) {
			t.Fatal("Should notcontain a X-wing but it did.")
		}
		currNode = currNode.Next
	}
}

func TestRule20(t *testing.T) {
	mat2 = Pmat{}
	mat2[0][0] = []int{1, 2, 5, 6}
	mat2[0][2] = []int{1, 5, 6, 9}
	mat2[0][3] = []int{2, 5, 7, 9}
	mat2[0][4] = []int{5, 7, 9}
	mat2[0][6] = []int{1, 2, 5, 6, 7, 9}
	mat2[0][7] = []int{2, 5, 6, 7}
	mat2[2][0] = []int{2, 5, 8}
	mat2[2][2] = []int{5, 8, 9}
	mat2[2][3] = []int{2, 5, 7, 9}
	mat2[2][6] = []int{2, 5, 7, 9}
	mat2[2][7] = []int{2, 5, 7}

	emptyL = CreatelinkedList()
	emptyL.AddCell(0, 6, []int{1, 2, 5, 6, 7, 9})
	emptyL.AddCell(0, 7, []int{2, 5, 6, 7})
	emptyL.AddCell(2, 6, []int{2, 5, 7, 9})
	emptyL.AddCell(2, 7, []int{2, 5, 7})

	PrintPossibleMat(mat2)

	startCnt := CountElemPosMat(mat2)

	_, cnt := rule20()

	endCnt := CountElemPosMat(mat2)

	if cnt != 1 {
		t.Fatal("Should have found X-wing but not.\n")
	}

	if Contains(mat2[0][6], 2) {
		t.Fatal("Possibility matrix cell [0,6] should not contain 2.\n")
	}

	if Contains(mat2[0][7], 2) {
		t.Fatal("Possibility matrix cell [0,7] should not contain 2.\n")
	}

	if Contains(mat2[2][6], 2) {
		t.Fatal("Possibility matrix cell [2,6] should not contain 2.\n")
	}

	if Contains(mat2[2][7], 2) {
		t.Fatal("Possibility matrix cell [2,7] should not contain 2.\n")
	}

	if (startCnt - endCnt) != 4 {
		t.Fatalf("Should have erased 4 values but got %d.\n", startCnt-endCnt)
	}
}
