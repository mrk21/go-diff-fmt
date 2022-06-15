package difffmt

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

const UnifiedTimeFormat = "2006-01-02 15:04:05.000000000 -0700"

type UnifiedFormatOption struct {
	ColorMode           ColorMode
	IsHidingNoLFMessage bool
}

type UnifiedFormat struct {
	option           UnifiedFormatOption
	headerFmt        *Colorizer
	hunkRangeInfoFmt *Colorizer
	insertFmt        *Colorizer
	deleteFmt        *Colorizer
}

func NewUnifiedFormat(option UnifiedFormatOption) *UnifiedFormat {
	return &UnifiedFormat{
		option:           option,
		headerFmt:        NewColorizer(option.ColorMode, color.Bold),
		hunkRangeInfoFmt: NewColorizer(option.ColorMode, color.FgCyan),
		insertFmt:        NewColorizer(option.ColorMode, color.FgGreen),
		deleteFmt:        NewColorizer(option.ColorMode, color.FgRed),
	}
}

func (u *UnifiedFormat) Print(targetA *DiffTarget, targetB *DiffTarget, hunks []Hunk) {
	u.Fprint(os.Stdout, targetA, targetB, hunks)
}

func (u *UnifiedFormat) Sprint(targetA *DiffTarget, targetB *DiffTarget, hunks []Hunk) string {
	var buf bytes.Buffer
	u.Fprint(&buf, targetA, targetB, hunks)
	return buf.String()
}

func (u *UnifiedFormat) Fprint(w io.Writer, targetA *DiffTarget, targetB *DiffTarget, hunks []Hunk) {
	u.headerFmt.Fprintf(w, "--- %s\n", u.headerValue(targetA))
	u.headerFmt.Fprintf(w, "+++ %s\n", u.headerValue(targetB))

	for _, hunk := range hunks {
		u.hunkRangeInfoFmt.Fprintf(w, "@@ -%s +%s @@\n",
			u.hunkRange(hunk.OldLineFrom, hunk.OldLineCount),
			u.hunkRange(hunk.NewLineFrom, hunk.NewLineCount),
		)

		for _, diff := range hunk.Diffs {
			text := diff.Text

			switch diff.Operation {
			case OperationInsert:
				u.insertFmt.Fprint(w, "+")
				u.insertFmt.Fprint(w, text)
			case OperationDelete:
				u.deleteFmt.Fprint(w, "-")
				u.deleteFmt.Fprint(w, text)
			case OperationEqual:
				fmt.Fprint(w, " ")
				fmt.Fprint(w, text)
			}

			if !diff.IsEndedLF {
				if u.option.IsHidingNoLFMessage {
					fmt.Fprint(w, "\n")
				} else {
					fmt.Fprint(w, "\n\\ No newline at end of file\n")
				}
			}
		}
	}
}

func (u *UnifiedFormat) headerValue(target *DiffTarget) string {
	if target.ModifiedTime.IsZero() {
		return target.Path
	} else {
		return fmt.Sprintf("%s\t%s", target.Path, target.ModifiedTime.Format(UnifiedTimeFormat))
	}
}

func (u *UnifiedFormat) hunkRange(from int, count int) string {
	if count == 1 {
		return fmt.Sprintf("%d", from)
	} else {
		return fmt.Sprintf("%d,%d", from, count)
	}
}
