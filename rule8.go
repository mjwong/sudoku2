package main

import (
	"time"

	//. "github.com/mjwong/sudoku2/lib"
	//. "github.com/mjwong/sudoku2/linkedlist"
	. "github.com/mjwong/sudoku2/matchlist"
	"gopkg.in/gookit/color.v1"
)

// Rule 8:	Hidden pairs
//
//
func rule8() (*Matchlist, int, time.Duration) {
	var (
		count int
		col, row/*, col2, row2*/ int
		foundHiddenPairs, debug bool
		start                   time.Time
		matched                 *Matchlist
	)
	start = time.Now()
	debug = DebugFn(2)
	matched = &Matchlist{}

	for foundHiddenPairs {
		foundHiddenPairs = false

		currNode := emptyL.Head

		if currNode == nil {
			color.Yellow.Println("Rule 8: Empty list.")
			break
		} else {
			for currNode != nil {
				row = currNode.Row
				col = currNode.Col

				if debug {
					color.LightGreen.Printf("cell [%d][%d]. %+v\n", row, col, currNode.Vals)
				}
			}
		}
	}

	return matched, count, time.Since(start)
}
