package difffmt

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

const UnifiedTimeFormat = "2006-01-02 15:04:05.000000000 -0700"

type UnifiedFormat struct {
	IsColor bool
}

func (u UnifiedFormat) Format(w io.Writer, targetA *DiffTarget, targetB *DiffTarget, hunks []Hunk) {
	u.color(w, "\x1b[01m", func(w io.Writer) {
		u.header(w, targetA, targetB)
	})

	for _, hunk := range hunks {
		u.color(w, "\x1b[36m", func(w io.Writer) {
			fmt.Fprintf(w, "@@ -%s +%s @@\n",
				u.hunkRange(hunk.OldLineFrom, hunk.OldLineCount()),
				u.hunkRange(hunk.NewLineFrom, hunk.NewLineCount()),
			)
		})

		for _, diff := range hunk.Diffs {
			text := diff.Text

			switch diff.Operation {
			case OperationInsert:
				u.color(w, "\x1b[32m", func(w io.Writer) {
					fmt.Fprint(w, "+")
					fmt.Fprint(w, text)
				})
			case OperationDelete:
				u.color(w, "\x1b[31m", func(w io.Writer) {
					fmt.Fprint(w, "-")
					fmt.Fprint(w, text)
				})
			case OperationEqual:
				fmt.Fprint(w, " ")
				fmt.Fprint(w, text)
			}

			if !diff.IsEndedLF {
				fmt.Fprint(w, "\n\\ No newline at end of file\n")
			}
		}
	}
}

func (u UnifiedFormat) header(w io.Writer, targetA *DiffTarget, targetB *DiffTarget) {
	fmt.Fprintf(w, "--- %s\t%s\n", targetA.Path, targetA.ModifiedTime.Format(UnifiedTimeFormat))
	fmt.Fprintf(w, "+++ %s\t%s\n", targetB.Path, targetB.ModifiedTime.Format(UnifiedTimeFormat))
}

func (u UnifiedFormat) hunkRange(from int, count int) string {
	if count == 1 {
		return fmt.Sprintf("%d", from)
	} else {
		return fmt.Sprintf("%d,%d", from, count)
	}
}

func (u UnifiedFormat) color(w io.Writer, color string, output func(w io.Writer)) {
	if u.IsColor && u.isTerminal() {
		fmt.Fprint(w, color)
		output(w)
		fmt.Fprint(w, "\x1b[0m")
	} else {
		output(w)
	}
}

func (u UnifiedFormat) isTerminal() bool {
	return terminal.IsTerminal(int(os.Stdout.Fd()))
}
