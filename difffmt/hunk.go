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
		h.NewLineCount++
		if h.NewLineCount == 1 {
			h.NewLineFrom = diff.NewLine
		}
		h.NewLineTo = diff.NewLine
	}
	switch diff.Operation {
	case OperationEqual, OperationDelete:
		h.OldLineCount++
		if h.OldLineCount == 1 {
			h.OldLineFrom = diff.OldLine
		}
		h.OldLineTo = diff.OldLine
	}
}

func GetHunks(diffs []LineDiff, contextSize int) []Hunk {
	// Correct lines
	lineMap := map[int]bool{}
	for i, diff := range diffs {
		switch diff.Operation {
		case OperationInsert, OperationDelete:
			begin := max((i - contextSize), 0)
			end := min((i + contextSize), (len(diffs) - 1))
			for j := begin; j <= end; j++ {
				lineMap[j] = true
			}
		}
	}
	lines := []int{}
	for line := range lineMap {
		lines = append(lines, line)
	}
	sort.Ints(lines)

	// Make hunks
	hunks := []Hunk{}
	var currentHunk *Hunk = nil

	for i, line := range lines {
		if currentHunk == nil || lines[i-1]+1 != line {
			hunks = append(hunks, NewHunk())
			currentHunk = &hunks[len(hunks)-1]
		}
		currentHunk.AppendDiff(diffs[line])
	}

	return hunks
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
