package main

import (
	"github.com/mrk21/go-diff-fmt/difffmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	targetA := difffmt.NewDiffTarget("a")
	targetB := difffmt.NewDiffTarget("b")
	textA := "line1\nline2a\nline3\n"
	textB := "line1\nline2b\nline3b\n"

	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(textA, textB)
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	lineDiffs := difffmt.GetLineDiffsFromDMP(diffs)
	hunks := difffmt.GetHunks(lineDiffs, 3)
	unifiedFmt := difffmt.UnifiedFormat{IsColorize: true}
	unifiedFmt.Print(targetA, targetB, hunks)
}
