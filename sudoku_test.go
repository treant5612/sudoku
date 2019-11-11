package sudoku

import (
	"fmt"
	"testing"
	"time"
)

func TestPuzzle_Solve(t *testing.T) {
	puzzleString := `
4_5  7__  ___
92_  ___  ___
___  ___  _6_
_8_  __6  7__
6__  _9_  __3
___  __7  6__
5__  1__  __2`
	puzzle,err := New(puzzleString)
	if err!=nil{
		t.FailNow()
	}
	start :=time.Now()
	solved,err :=puzzle.Solve()
	if err!=nil{
		t.FailNow()
	}
	cost := time.Since(start)

	fmt.Printf("time cost:%v\n%v",cost.Nanoseconds(),solved)

}