package difffmt

import (
	"fmt"
	"os"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Test_UnifiedFormat(t *testing.T) {
	targetA := NewDiffTarget("a")
	targetB := NewDiffTarget("b")
	textA := "line1\nline2a\nline3\n"
	textB := "line1\nline2b\nline3b\n"

	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(textA, textB)
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	lineDiffs := GetLineDiffsFromDMP(diffs)
	hunks := GetHunks(lineDiffs, 3)
	unifiedFmt := UnifiedFormat{IsColorize: true}

	t.Run("Print", func(t *testing.T) {
		unifiedFmt.Print(targetA, targetB, hunks)
	})

	t.Run("Sprint", func(t *testing.T) {
		result := unifiedFmt.Sprint(targetA, targetB, hunks)
		fmt.Print(result)
	})

	t.Run("Fprint", func(t *testing.T) {
		unifiedFmt.Fprint(os.Stdout, targetA, targetB, hunks)
	})
}
