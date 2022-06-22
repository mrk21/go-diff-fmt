package difffmt

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

const ContextTimeFormat = "2006-01-02 15:04:05.000000000 -0700"

type ContextFormatOption struct {
	ColorMode           ColorMode
	IsHidingNoLFMessage bool
}

type ContextFormat struct {
	option           ContextFormatOption
	headerFmt        *Colorizer
	hunkRangeInfoFmt *Colorizer
	insertFmt        *Colorizer
	deleteFmt        *Colorizer
}

func NewContextFormat(option ContextFormatOption) *ContextFormat {
	return &ContextFormat{
		option:           option,
		headerFmt:        NewColorizer(option.ColorMode, color.Bold),
		hunkRangeInfoFmt: NewColorizer(option.ColorMode, color.FgCyan),
		insertFmt:        NewColorizer(option.ColorMode, color.FgGreen),
		deleteFmt:        NewColorizer(option.ColorMode, color.FgRed),
	}
}

func (c *ContextFormat) Print(targetA *DiffTarget, targetB *DiffTarget, hunks []Hunk) {
	c.Fprint(os.Stdout, targetA, targetB, hunks)
}

func (c *ContextFormat) Sprint(targetA *DiffTarget, targetB *DiffTarget, hunks []Hunk) string {
	var buf bytes.Buffer
	c.Fprint(&buf, targetA, targetB, hunks)
	return buf.String()
}

func (c *ContextFormat) Fprint(w io.Writer, targetA *DiffTarget, targetB *DiffTarget, hunks []Hunk) {
	c.headerFmt.Fprintf(w, "*** %s\n", c.headerValue(targetA))
	c.headerFmt.Fprintf(w, "--- %s\n", c.headerValue(targetB))

	for _, hunk := range hunks {
		fmt.Fprintln(w, "***************")

		c.hunkRangeInfoFmt.Fprintf(w, "*** %s ****\n", c.hunkRange(hunk.OldLineFrom, hunk.OldLineTo))
		for _, diff := range hunk.Diffs {
			text := diff.Text

			switch diff.Operation {
			case OperationInsert:
				continue
			case OperationDelete:
				if diff.Context == ContextModify {
					c.insertFmt.Fprint(w, "! ")
				} else {
					c.insertFmt.Fprint(w, "- ")
				}
				c.insertFmt.Fprint(w, text)
			case OperationEqual:
				c.deleteFmt.Fprint(w, "  ")
				c.deleteFmt.Fprint(w, text)
			}

			if !diff.IsEndedLF {
				if c.option.IsHidingNoLFMessage {
					fmt.Fprint(w, "\n")
				} else {
					fmt.Fprint(w, "\n\\ No newline at end of file\n")
				}
			}
		}

		c.hunkRangeInfoFmt.Fprintf(w, "--- %s ----\n", c.hunkRange(hunk.NewLineFrom, hunk.NewLineTo))
		for _, diff := range hunk.Diffs {
			text := diff.Text

			switch diff.Operation {
			case OperationInsert:
				if diff.Context == ContextModify {
					c.insertFmt.Fprint(w, "! ")
				} else {
					c.insertFmt.Fprint(w, "+ ")
				}
				c.insertFmt.Fprint(w, text)
			case OperationDelete:
				continue
			case OperationEqual:
				c.insertFmt.Fprint(w, "  ")
				c.insertFmt.Fprint(w, text)
			}

			if !diff.IsEndedLF {
				if c.option.IsHidingNoLFMessage {
					fmt.Fprint(w, "\n")
				} else {
					fmt.Fprint(w, "\n\\ No newline at end of file\n")
				}
			}
		}
	}
}

func (u *ContextFormat) headerValue(target *DiffTarget) string {
	if target.ModifiedTime.IsZero() {
		return target.Path
	} else {
		return fmt.Sprintf("%s\t%s", target.Path, target.ModifiedTime.Format(ContextTimeFormat))
	}
}

func (u *ContextFormat) hunkRange(from int, to int) string {
	if from == to {
		return fmt.Sprintf("%d", from)
	} else {
		return fmt.Sprintf("%d,%d", from, to)
	}
}
