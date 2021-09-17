package main

import (
	"testing"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/linkedlist"
	. "github.com/mjwong/sudoku2/matchlist"
)

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

	eraCnt, erased := eraseDigitsFromRowMulti(0, []int{2}, []int{0, 3})

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
		{Row: 1, Col: 2, Vals: []int{1, 2, 5}},
		{Row: 2, Col: 3, Vals: []int{2, 3, 6}},
		{Row: 3, Col: 4, Vals: []int{2, 4, 7}},
		{Row: 4, Col: 5, Vals: []int{3, 5, 9}},
	}

	arr2 := []Idx{
		{Row: 5, Col: 6, Vals: []int{1, 2, 5}},
		{Row: 6, Col: 8, Vals: []int{2, 3, 6}},
		{Row: 3, Col: 4, Vals: []int{2, 4, 7}},
		{Row: 4, Col: 5, Vals: []int{3, 5, 9}},
	}

	arr3 := []Idx{
		{Row: 1, Col: 2, Vals: []int{1, 2, 5}},
		{Row: 2, Col: 3, Vals: []int{2, 3, 6}},
		{Row: 3, Col: 4, Vals: []int{2, 4, 7}},
		{Row: 4, Col: 5, Vals: []int{3, 5, 9}},
	}

	arr4 := []Idx{
		{Row: 6, Col: 2, Vals: []int{1, 2, 5}},
		{Row: 7, Col: 9, Vals: []int{2, 3, 6}},
		{Row: 1, Col: 3, Vals: []int{2, 4, 7}},
		{Row: 2, Col: 7, Vals: []int{3, 5, 9}},
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

func TestFindTripInRow(t *testing.T) {

	pm := Pmat{}
	pm[0][0] = []int{1, 5, 6}
	pm[0][2] = []int{2, 8, 9}
	pm[0][5] = []int{5, 9}
	pm[0][6] = []int{1, 5, 6}
	pm[0][7] = []int{2, 5, 8}
	pm[0][8] = []int{1, 5, 6}

	PrintPossibleMat(pm)

	cnt, found := findTripInRow(pm)
	if !found {
		t.Fatal("Should find triplet in row 0 but did not.")
	}

	if cnt != 1 {
		t.Fatalf("Should find 1 triplet but got %d.\n", cnt)
	}
}
