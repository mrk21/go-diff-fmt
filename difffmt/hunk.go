package difffmt

import (
	"sort"
)

type Hunk struct {
	Diffs        []LineDiff
	OldLineFrom  int
	OldLineTo    int
	OldLineCount int
	NewLineFrom  int
	NewLineTo    int
	NewLineCount int
}

func NewHunk() Hunk {
	return Hunk{
		Diffs:        []LineDiff{},
		OldLineFrom:  0,
		OldLineTo:    0,
		OldLineCount: 0,
		NewLineFrom:  0,
		NewLineTo:    0,
		NewLineCount: 0,
	}
}

func (h *Hunk) AppendDiff(diff LineDiff) {
	h.Diffs = append(h.Diffs, diff)
	if len(h.Diffs) == 1 {
		h.NewLineFrom = diff.NewLine
		h.NewLineTo = diff.NewLine
		h.OldLineFrom = diff.OldLine
		h.OldLineTo = diff.OldLine
	}
	switch diff.Operation {
	case OperationEqual, OperationInsert:
		if h.NewLineCount == 0 {
			h.NewLineFrom = diff.NewLine
		}
		h.NewLineTo = diff.NewLine
		h.NewLineCount++
	}
	switch diff.Operation {
	case OperationEqual, OperationDelete:
		if h.OldLineCount == 0 {
			h.OldLineFrom = diff.OldLine
		}
		h.OldLineTo = diff.OldLine
		h.OldLineCount++
	}
}

func GetHunks(diffs []LineDiff, contextSize int) []Hunk {
	// Correct output lines
	outputLineMap := map[int]bool{}
	for i, diff := range diffs {
		switch diff.Operation {
		case OperationInsert, OperationDelete:
			begin := max(i-contextSize, 0)
			end := min(i+contextSize, len(diffs)-1)
			for j := begin; j <= end; j++ {
				outputLineMap[j] = true
			}
		}
	}
	outputLine := []int{}
	for line := range outputLineMap {
		outputLine = append(outputLine, line)
	}
	sort.Ints(outputLine)

	// Group output lines
	result := []Hunk{}
	var currentHunk *Hunk = nil

	for i, line := range outputLine {
		if i == 0 || outputLine[i-1]+1 != line {
			result = append(result, NewHunk())
			currentHunk = &result[len(result)-1]
		}
		currentHunk.AppendDiff(diffs[line])
	}

	return result
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}
