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

	// Computing a line-mode diff
	// @see https://github.com/google/diff-match-patch/wiki/Line-or-Word-Diffs
	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(textA, textB)
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	// Format `[]diffmatchpatch.Diff` to Unified format
	lineDiffs := difffmt.MakeLineDiffsFromDMP(diffs)
	hunks := difffmt.MakeHunks(lineDiffs, 3)
	unifiedFmt := difffmt.NewUnifiedFormat(difffmt.UnifiedFormatOption{
		ColorMode: difffmt.ColorTerminalOnly,
	})
	unifiedFmt.Print(targetA, targetB, hunks)
}
