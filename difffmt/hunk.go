package difffmt

import (
	"sort"
)

type Hunk struct {
	Diffs       []LineDiff
	OldLineFrom int
	OldLineTo   int
	NewLineFrom int
	NewLineTo   int
}

func NewHunk() Hunk {
	return Hunk{
		Diffs:       []LineDiff{},
		OldLineFrom: 0,
		OldLineTo:   0,
		NewLineFrom: 0,
		NewLineTo:   0,
	}
}

func (h *Hunk) AppendDiff(diff LineDiff) {
	h.Diffs = append(h.Diffs, diff)
	switch diff.Operation {
	case OperationEqual, OperationInsert:
		if h.NewLineFrom == 0 {
			h.NewLineFrom = diff.NewLine
		}
		h.NewLineTo = diff.NewLine
	}
	switch diff.Operation {
	case OperationEqual, OperationDelete:
		if h.OldLineFrom == 0 {
			h.OldLineFrom = diff.OldLine
		}
		h.OldLineTo = diff.OldLine
	}
}

func (h *Hunk) FixDiffs() {
	if h.NewLineFrom == 0 {
		h.NewLineFrom = h.Diffs[0].NewLine
	}
	if h.OldLineFrom == 0 {
		h.OldLineFrom = h.Diffs[0].OldLine
	}
}

func (h *Hunk) OldLineCount() int {
	if h.OldLineTo == 0 {
		return 0
	}
	return h.OldLineTo - h.OldLineFrom + 1
}

func (h *Hunk) NewLineCount() int {
	if h.NewLineTo == 0 {
		return 0
	}
	return h.NewLineTo - h.NewLineFrom + 1
}

func GetHunk(diffs []LineDiff, contextSize int) []Hunk {
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
	sort.Slice(outputLine, func(i, j int) bool { return outputLine[i] < outputLine[j] })

	// Group output lines
	result := []Hunk{}
	var currentHunk *Hunk = nil

	for i, line := range outputLine {
		if i == 0 || outputLine[i-1]+1 != line {
			if currentHunk != nil {
				currentHunk.FixDiffs()
			}
			result = append(result, NewHunk())
			currentHunk = &result[len(result)-1]
		}
		currentHunk.AppendDiff(diffs[line])
	}
	if currentHunk != nil {
		currentHunk.FixDiffs()
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
