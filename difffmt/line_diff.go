package difffmt

import (
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

type Operation int8

const (
	OperationDelete Operation = -1
	OperationInsert Operation = 1
	OperationEqual  Operation = 0
)

var OperationMapForDMP = map[diffmatchpatch.Operation]Operation{
	diffmatchpatch.DiffDelete: OperationDelete,
	diffmatchpatch.DiffInsert: OperationInsert,
	diffmatchpatch.DiffEqual:  OperationEqual,
}

type LineDiff struct {
	Operation Operation
	Text      string
	OldLine   int
	NewLine   int
	IsEndedLF bool
}

// In order to make line diffs from github.com/sergi/go-diff, you can make by steps shown below:
//
// 	// Computing a line-mode diff by github.com/sergi/go-diff
// 	// @see https://github.com/google/diff-match-patch/wiki/Line-or-Word-Diffs
// 	dmp := diffmatchpatch.New()
// 	runes1, runes2, lineArray := dmp.DiffLinesToRunes(text1, text2)
// 	diffs := dmp.DiffMainRunes(runes1, runes2, false)
// 	diffs = dmp.DiffCharsToLines(diffs, lineArray)
//
// 	// Make `[]LineDiff` from `[]diffmatchpatch.Diff`
// 	lineDiffs := difffmt.MakeLineDiffsFromDMP(diffs)
//
func MakeLineDiffsFromDMP(diffs []diffmatchpatch.Diff) []LineDiff {
	result := []LineDiff{}
	currentOldLine := 0
	currentNewLine := 0

	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		operation := OperationMapForDMP[diff.Type]

		for i, line := range lines {
			if i+1 == len(lines) && line == "" {
				break
			}

			switch operation {
			case OperationEqual:
				currentNewLine++
				currentOldLine++
			case OperationDelete:
				currentOldLine++
			case OperationInsert:
				currentNewLine++
			}

			if i+1 == len(lines) && line != "" {
				result = append(result, LineDiff{
					Operation: operation,
					Text:      line,
					OldLine:   currentOldLine,
					NewLine:   currentNewLine,
					IsEndedLF: false,
				})
			} else {
				result = append(result, LineDiff{
					Operation: operation,
					Text:      line + "\n",
					OldLine:   currentOldLine,
					NewLine:   currentNewLine,
					IsEndedLF: true,
				})
			}
		}
	}
	return result
}
