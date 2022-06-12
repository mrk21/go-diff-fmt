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

type LineDiff struct {
	Operation Operation
	Text      string
	OldLine   int
	NewLine   int
	IsEndedLF bool
}

func GetLineDiff(text1 string, text2 string) []LineDiff {
	// Make lines diffs
	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(string(text1), string(text2))
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	// Make single line diffs
	result := []LineDiff{}
	currentOldLine := 0
	currentNewLine := 0

	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		operation := GetOperationFromDiffMatchPatch(diff.Type)

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

func GetOperationFromDiffMatchPatch(operation diffmatchpatch.Operation) Operation {
	switch operation {
	case diffmatchpatch.DiffDelete:
		return OperationDelete
	case diffmatchpatch.DiffInsert:
		return OperationInsert
	case diffmatchpatch.DiffEqual:
		return OperationEqual
	}
	panic("Invalid operation")
}
