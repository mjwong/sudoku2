package main

import (
	"fmt"
	"testing"

	. "github.com/mjwong/sudoku2/lib"
)

func TestRule3n5(t *testing.T) {
	PrepPmat(difficult3)
	totCnt := RuleLoop(rule3, RuleTable[3], Zero)
	if totCnt != 29 {
		t.Fatalf("Expected to find 29 but got %d counts.\n", totCnt)
	}

	fmt.Println("Starting possible matrix for Rule 5.")
	PrintPossibleMat(mat2)

	input := "142.73...597.462.3863.52...31852469772639.4.545976.32.6.54391.293128....2.461..39"
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

func TestRule135(t *testing.T) { // difficult5.txt
	ruleCnt := map[int]int{}
	PrepPmat(difficult5)

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
	ruleCnt := map[int]int{}

	// mid-way through difficult5.txt
	input := ".4...8..37..4.382..3..16..49.4.6.38..6..3.49..23.4...64..12.63.31268..45...3.4.1."
	PrepPmat(input)

	startCnt := emptyL.CountNodes()
	if startCnt != 41 {
		t.Fatalf("Expected to find 41 but got %d counts.\n", startCnt)
	}
	fmt.Printf("Starting empty count: %d\n", startCnt)

	ruleCnt[1] = RuleLoop(rule1, RuleTable[1], Zero)
	ruleCnt[3] = RuleLoop(rule3, RuleTable[3], Zero)
	cntBefore := emptyL.CountElem()
	ruleCnt[5] = RuleLoop(rule5, RuleTable[5], SameCnt)
	cntAfter := emptyL.CountElem()

	if ruleCnt[3] != 41 {
		t.Fatalf("Expect to find 41 hidden singles but got %d.\n", ruleCnt[3])
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

func TestRule_135c(t *testing.T) {
	ruleCnt := map[int]int{}

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

	for {
		cnt1 := RuleLoop(rule1, RuleTable[1], Zero)
		cnt3 := RuleLoop(rule3, RuleTable[3], Zero)
		cntBefore := emptyL.CountElem()
		cnt5 := RuleLoop(rule5, RuleTable[5], SameCnt)
		cntAfter := emptyL.CountElem()

		ruleCnt[1] += cnt1
		ruleCnt[3] += cnt3
		ruleCnt[5] += cnt5

		if cnt1 == 0 && cnt3 == 0 && cntBefore == cntAfter {
			fmt.Printf("Count before and after Rule 5: %d, %d.\n", cntBefore, cntAfter)
			break
		}
	}

	PrintFound([]int{1, 3, 5}, ruleCnt)
	fmt.Printf("Empty cells : %2d\n", emptyL.CountNodes())

	if emptyL.CountNodes() == 0 {
		if !CheckSums(mat) {
			t.Fatal("Expected to be solved")
			PrintSudoku(mat)
		}
	}
}
