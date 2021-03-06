package difffmt

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func Test_Hunk(t *testing.T) {
	type testArg struct {
		a        string
		b        string
		n        int
		expected []Hunk
	}

	doTest := func(t *testing.T, arg testArg) {
		dmp := diffmatchpatch.New()
		runes1, runes2, lineArray := dmp.DiffLinesToRunes(arg.a, arg.b)
		diffs := dmp.DiffMainRunes(runes1, runes2, false)
		diffs = dmp.DiffCharsToLines(diffs, lineArray)

		lineDiffs := MakeLineDiffsFromDMP(diffs)
		actual := MakeHunks(lineDiffs, arg.n)

		if d := cmp.Diff(actual, arg.expected); d != "" {
			as, _ := json.MarshalIndent(actual, "", "    ")
			es, _ := json.MarshalIndent(arg.expected, "", "    ")
			t.Errorf("[Error] actual: %s\nexpected:%s\n", as, es)
			t.Errorf("[Error] (-actual +expected):\n%s", d)
		}
	}

	toText := func(lines []string) string {
		return strings.Join(lines, "")
	}

	t.Run("insert LF only", func(t *testing.T) {
		a := "line1"
		b := "line1\n"

		for n := 0; n <= 1; n++ {
			t.Run(fmt.Sprintf("contextSize = %d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Context:   ContextModify,
								Text:      "line1",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: false,
							},
							{
								Operation: OperationInsert,
								Context:   ContextModify,
								Text:      "line1\n",
								OldLine:   1,
								NewLine:   1,
								IsEndedLF: true,
							},
						},
						OldLineFrom:  1,
						OldLineTo:    1,
						OldLineCount: 1,
						NewLineFrom:  1,
						NewLineTo:    1,
						NewLineCount: 1,
					},
				}
				doTest(t, testArg{
					a:        a,
					b:        b,
					n:        n,
					expected: expected,
				})
			})
		}
	})

	t.Run("delete LF only", func(t *testing.T) {
		a := "line1\n"
		b := "line1"

		for n := 0; n <= 1; n++ {
			t.Run(fmt.Sprintf("contextSize = %d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Context:   ContextModify,
								Text:      "line1\n",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: true,
							},
							{
								Operation: OperationInsert,
								Context:   ContextModify,
								Text:      "line1",
								OldLine:   1,
								NewLine:   1,
								IsEndedLF: false,
							},
						},
						OldLineFrom:  1,
						OldLineTo:    1,
						OldLineCount: 1,
						NewLineFrom:  1,
						NewLineTo:    1,
						NewLineCount: 1,
					},
				}
				doTest(t, testArg{
					a:        a,
					b:        b,
					n:        n,
					expected: expected,
				})
			})
		}
	})

	t.Run("without LF", func(t *testing.T) {
		a := "line1"
		b := "line2"

		for n := 0; n <= 1; n++ {
			t.Run(fmt.Sprintf("contextSize = %d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Context:   ContextModify,
								Text:      "line1",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: false,
							},
							{
								Operation: OperationInsert,
								Context:   ContextModify,
								Text:      "line2",
								OldLine:   1,
								NewLine:   1,
								IsEndedLF: false,
							},
						},
						OldLineFrom:  1,
						OldLineTo:    1,
						OldLineCount: 1,
						NewLineFrom:  1,
						NewLineTo:    1,
						NewLineCount: 1,
					},
				}
				doTest(t, testArg{
					a:        a,
					b:        b,
					n:        n,
					expected: expected,
				})
			})
		}
	})

	t.Run("equal", func(t *testing.T) {
		a := toText([]string{
			"line1\n",
			"line2\n",
			"line3\n",
		})
		b := toText([]string{
			"line1\n",
			"line2\n",
			"line3\n",
		})

		for n := 0; n <= 1; n++ {
			t.Run(fmt.Sprintf("contextSize = %d", n), func(t *testing.T) {
				expected := []Hunk{}
				doTest(t, testArg{
					a:        a,
					b:        b,
					n:        n,
					expected: expected,
				})
			})
		}
	})

	t.Run("insert only", func(t *testing.T) {
		a := toText([]string{})
		b := toText([]string{
			"line1\n",
			"line2\n",
			"line3\n",
		})

		for n := 0; n <= 1; n++ {
			t.Run(fmt.Sprintf("contextSize = %d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationInsert,
								Context:   ContextAdd,
								Text:      "line1\n",
								OldLine:   0,
								NewLine:   1,
								IsEndedLF: true,
							},
							{
								Operation: OperationInsert,
								Context:   ContextAdd,
								Text:      "line2\n",
								OldLine:   0,
								NewLine:   2,
								IsEndedLF: true,
							},
							{
								Operation: OperationInsert,
								Context:   ContextAdd,
								Text:      "line3\n",
								OldLine:   0,
								NewLine:   3,
								IsEndedLF: true,
							},
						},
						OldLineFrom:  0,
						OldLineTo:    0,
						OldLineCount: 0,
						NewLineFrom:  1,
						NewLineTo:    3,
						NewLineCount: 3,
					},
				}
				doTest(t, testArg{
					a:        a,
					b:        b,
					n:        n,
					expected: expected,
				})
			})
		}
	})

	t.Run("delete only", func(t *testing.T) {
		a := toText([]string{
			"line1\n",
			"line2\n",
			"line3\n",
		})
		b := toText([]string{})

		for n := 0; n <= 1; n++ {
			t.Run(fmt.Sprintf("contextSize = %d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Context:   ContextDelete,
								Text:      "line1\n",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: true,
							},
							{
								Operation: OperationDelete,
								Context:   ContextDelete,
								Text:      "line2\n",
								OldLine:   2,
								NewLine:   0,
								IsEndedLF: true,
							},
							{
								Operation: OperationDelete,
								Context:   ContextDelete,
								Text:      "line3\n",
								OldLine:   3,
								NewLine:   0,
								IsEndedLF: true,
							},
						},
						OldLineFrom:  1,
						OldLineTo:    3,
						OldLineCount: 3,
						NewLineFrom:  0,
						NewLineTo:    0,
						NewLineCount: 0,
					},
				}
				doTest(t, testArg{
					a:        a,
					b:        b,
					n:        n,
					expected: expected,
				})
			})
		}
	})

	t.Run("normal", func(t *testing.T) {
		a := toText([]string{
			"line00\n",
			"line01\n",
			"line1\n",
			"line2\n",
			"line3\n",
			"line4\n",
			"line5\n",
			"line6\n",
			"line8\n",
			"line9\n",
			"line10\n",
			"line11\n",
			"line12\n",
		})
		b := toText([]string{
			"line1\n",
			"line2\n",
			"line3\n",
			"line4\n",
			"line5b\n",
			"line6\n",
			"line7b\n",
			"line8b\n",
			"line9\n",
			"line10\n",
			"line11\n",
			"line12\n",
			"line13\n",
			"line14",
		})

		t.Run("contextSize = 0", func(t *testing.T) {
			n := 0
			expected := []Hunk{
				{
					Diffs: []LineDiff{
						{
							Operation: OperationDelete,
							Context:   ContextDelete,
							Text:      "line00\n",
							OldLine:   1,
							NewLine:   0,
							IsEndedLF: true,
						},
						{
							Operation: OperationDelete,
							Context:   ContextDelete,
							Text:      "line01\n",
							OldLine:   2,
							NewLine:   0,
							IsEndedLF: true,
						},
					},
					OldLineFrom:  1,
					OldLineTo:    2,
					OldLineCount: 2,
					NewLineFrom:  0,
					NewLineTo:    0,
					NewLineCount: 0,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationDelete,
							Context:   ContextModify,
							Text:      "line5\n",
							OldLine:   7,
							NewLine:   4,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextModify,
							Text:      "line5b\n",
							OldLine:   7,
							NewLine:   5,
							IsEndedLF: true,
						},
					},
					OldLineFrom:  7,
					OldLineTo:    7,
					OldLineCount: 1,
					NewLineFrom:  5,
					NewLineTo:    5,
					NewLineCount: 1,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationDelete,
							Context:   ContextModify,
							Text:      "line8\n",
							OldLine:   9,
							NewLine:   6,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextModify,
							Text:      "line7b\n",
							OldLine:   9,
							NewLine:   7,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextModify,
							Text:      "line8b\n",
							OldLine:   9,
							NewLine:   8,
							IsEndedLF: true,
						},
					},
					OldLineFrom:  9,
					OldLineTo:    9,
					OldLineCount: 1,
					NewLineFrom:  7,
					NewLineTo:    8,
					NewLineCount: 2,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationInsert,
							Context:   ContextAdd,
							Text:      "line13\n",
							OldLine:   13,
							NewLine:   13,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextAdd,
							Text:      "line14",
							OldLine:   13,
							NewLine:   14,
							IsEndedLF: false,
						},
					},
					OldLineFrom:  13,
					OldLineTo:    13,
					OldLineCount: 0,
					NewLineFrom:  13,
					NewLineTo:    14,
					NewLineCount: 2,
				},
			}
			doTest(t, testArg{
				a:        a,
				b:        b,
				n:        n,
				expected: expected,
			})
		})

		t.Run("contextSize = 1", func(t *testing.T) {
			n := 1

			// Merge hunks that are adjacent
			expected := []Hunk{
				{
					Diffs: []LineDiff{
						{
							Operation: OperationDelete,
							Context:   ContextDelete,
							Text:      "line00\n",
							OldLine:   1,
							NewLine:   0,
							IsEndedLF: true,
						},
						{
							Operation: OperationDelete,
							Context:   ContextDelete,
							Text:      "line01\n",
							OldLine:   2,
							NewLine:   0,
							IsEndedLF: true,
						},
						{
							Operation: OperationEqual,
							Context:   ContextKeep,
							Text:      "line1\n",
							OldLine:   3,
							NewLine:   1,
							IsEndedLF: true,
						},
					},
					OldLineFrom:  1,
					OldLineTo:    3,
					OldLineCount: 3,
					NewLineFrom:  1,
					NewLineTo:    1,
					NewLineCount: 1,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationEqual,
							Context:   ContextKeep,
							Text:      "line4\n",
							OldLine:   6,
							NewLine:   4,
							IsEndedLF: true,
						},
						{
							Operation: OperationDelete,
							Context:   ContextModify,
							Text:      "line5\n",
							OldLine:   7,
							NewLine:   4,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextModify,
							Text:      "line5b\n",
							OldLine:   7,
							NewLine:   5,
							IsEndedLF: true,
						},
						{
							Operation: OperationEqual,
							Context:   ContextKeep,
							Text:      "line6\n",
							OldLine:   8,
							NewLine:   6,
							IsEndedLF: true,
						},
						{
							Operation: OperationDelete,
							Context:   ContextModify,
							Text:      "line8\n",
							OldLine:   9,
							NewLine:   6,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextModify,
							Text:      "line7b\n",
							OldLine:   9,
							NewLine:   7,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextModify,
							Text:      "line8b\n",
							OldLine:   9,
							NewLine:   8,
							IsEndedLF: true,
						},
						{
							Operation: OperationEqual,
							Context:   ContextKeep,
							Text:      "line9\n",
							OldLine:   10,
							NewLine:   9,
							IsEndedLF: true,
						},
					},
					OldLineFrom:  6,
					OldLineTo:    10,
					OldLineCount: 5,
					NewLineFrom:  4,
					NewLineTo:    9,
					NewLineCount: 6,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationEqual,
							Context:   ContextKeep,
							Text:      "line12\n",
							OldLine:   13,
							NewLine:   12,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextAdd,
							Text:      "line13\n",
							OldLine:   13,
							NewLine:   13,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Context:   ContextAdd,
							Text:      "line14",
							OldLine:   13,
							NewLine:   14,
							IsEndedLF: false,
						},
					},
					OldLineFrom:  13,
					OldLineTo:    13,
					OldLineCount: 1,
					NewLineFrom:  12,
					NewLineTo:    14,
					NewLineCount: 3,
				},
			}
			doTest(t, testArg{
				a:        a,
				b:        b,
				n:        n,
				expected: expected,
			})
		})
	})
}
