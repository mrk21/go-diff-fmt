package main

import (
	"github.com/mrk21/go-diff-fmt/difffmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	targetA := difffmt.NewDiffTarget("a.txt")
	targetB := difffmt.NewDiffTarget("b.txt")
	_ = targetA.LoadStat()
	_ = targetB.LoadStat()
	textA, _ := targetA.ReadText()
	textB, _ := targetB.ReadText()

	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(textA, textB)
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	lineDiffs := difffmt.GetLineDiffsFromDMP(diffs)
	hunks := difffmt.GetHunks(lineDiffs, 3)
	unifiedFmt := difffmt.UnifiedFormat{IsColorize: true}
	unifiedFmt.Print(targetA, targetB, hunks)
}
