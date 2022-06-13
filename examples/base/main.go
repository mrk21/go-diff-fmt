package main

import (
	"os"

	"github.com/mrk21/go-diff-fmt/difffmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	targetA := difffmt.NewDiffTarget("a.txt")
	targetB := difffmt.NewDiffTarget("b.txt")
	_ = targetA.LoadStats()
	_ = targetB.LoadStats()
	textA, _ := targetA.ReadText()
	textB, _ := targetB.ReadText()

	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(textA, textB)
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	lineDiffs := difffmt.GetLineDiffFromDiffMatchPatch(diffs)
	hunks := difffmt.GetHunk(lineDiffs, 3)
	unified := difffmt.UnifiedFormat{IsColor: true}
	unified.Format(os.Stdout, targetA, targetB, hunks)
}
