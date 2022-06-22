package difffmt

import (
	"fmt"
	"os"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Test_ContextFormat(t *testing.T) {
	targetA := NewDiffTarget("a")
	targetB := NewDiffTarget("b")
	textA := "line1\nline2a\nline3\n"
	textB := "line1\nline2b\nline3b\n"

	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(textA, textB)
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	lineDiffs := MakeLineDiffsFromDMP(diffs)
	hunks := MakeHunks(lineDiffs, 3)
	contextFmt := NewContextFormat(ContextFormatOption{
		ColorMode: ColorTerminalOnly,
	})

	t.Run("Print", func(t *testing.T) {
		contextFmt.Print(targetA, targetB, hunks)
	})

	t.Run("Sprint", func(t *testing.T) {
		result := contextFmt.Sprint(targetA, targetB, hunks)
		fmt.Print(result)
	})

	t.Run("Fprint", func(t *testing.T) {
		contextFmt.Fprint(os.Stdout, targetA, targetB, hunks)
	})
}
