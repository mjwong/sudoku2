package main

import (
	"fmt"

	. "github.com/mjwong/sudoku2/lib"
	. "github.com/mjwong/sudoku2/matchlist"
)

// Rule 6	Naked triplets
//			There are 3 dells in the same row, col or block that have exactly the same 3 digits.
//			Variation: It can also be a combination of 2 digits in one cell and 3 digits in another cell.
//	        E.g. The 3 cells contain (5,6), (6,8) and (8,5).
//			They form a closed loop 5 -> 8 -> 6 going from (5,6) to (5, 8) to (8,6).
func rule6() (*Matchlist, int) {
	var (
		count int
		//foundTrip bool
		matched *Matchlist
	)

	matched = &Matchlist{}
	return matched, count

	// Search for exactly cells that contain 3 same digits each.

	// search rows

}

// Find naked triplets in all rows for the case where all 3 cells have 3 same digits each
func findTripInRow(m2 Pmat) (int, bool) {
	var (
		count, tripCnt int
		trip           []int
		found          bool
		debug          bool
	)
	debug = DebugFn(2)
	found = false

	for i := 0; i < N; i++ {
		count = 0

		for j := 0; j < N; j++ {
			if m2[i][j] != nil && len(m2[i][j]) == 3 {
				if count == 0 {
					trip = m2[i][j]
					count = 1

					if debug {
						fmt.Printf("Triplets: %v at [%d,%d]. Count = %d\n", trip, i, j, count)
					}

				} else {
					if IntArrayEquals(trip, m2[i][j]) {
						count++

						if debug {
							fmt.Printf("Triplets: %v at [%d,%d]. Count = %d\n", trip, i, j, count)
						}
					}
				}
			}
		}

		if count == 3 {
			tripCnt++
			found = true
		}
	}

	return tripCnt, found
}

// Variation of Naked Triplets rule:
// Find naked triplets in row for the case where not all 3 cells contain 3 digits each.
// At least 1 cell must contain 3 digits.
// Solution: Search for a cell with 3 digits. Then search for exactly 2 more cells which contain either 2 or 3 digits
// where the digits are the same as that of the first 3-digit cell.
func findTripInRowVar(m2 Pmat, row int) bool {
	var (
		found         bool
		totCnt, count int
		trip, colList []int
	)

	for i := 0; i < N; i++ {
		for j := 0; j < N; j++ {
			count = 0
			colList = nil
			found = false

			if m2[i][j] != nil && len(m2[i][j]) == 3 { // first cell with 3 digits
				trip = m2[i][j]
				count = 1
				colList = append(colList, j)

				for c := 0; c < N; c++ {
					if m2[i][j] != nil && len(m2[i][j]) >= 2 && len(m2[i][j]) <= 3 && j != c {

						// case of 3 digits
						if len(m2[i][j]) == 3 {
							if IntArrayEquals(trip, m2[i][c]) {
								colList = append(colList, c)
								count++
							}
						} else {
							// case of 2 digits
							cnt := 0
							for _, v := range m2[i][c] {
								if Contains(trip, v) {
									cnt++
								}
							}

							if cnt == 2 {
								colList = append(colList, c)
								count++
							}
						}
					}
				}

				if count == 3 {
					found = true
					totCnt++ // no. of naked triplets found
				}

				if found {
					// erase the same 3 digits from elsewhere in row
					eraseDigitsFromRowMulti(i, trip, colList)
					fmt.Println()
				}
			}
		}
	}

	return false
}
