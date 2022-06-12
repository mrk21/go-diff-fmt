package difffmt

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_Hunk(t *testing.T) {
	type testArg struct {
		a        string
		b        string
		n        int
		expected []Hunk
	}

	doTest := func(t *testing.T, arg testArg) {
		actual := GetHunk(GetLineDiff(arg.a, arg.b), arg.n)
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
			t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Text:      "line1",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: false,
							},
							{
								Operation: OperationInsert,
								Text:      "line1\n",
								OldLine:   1,
								NewLine:   1,
								IsEndedLF: true,
							},
						},
						OldLineFrom: 1,
						OldLineTo:   1,
						NewLineFrom: 1,
						NewLineTo:   1,
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
			t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Text:      "line1\n",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: true,
							},
							{
								Operation: OperationInsert,
								Text:      "line1",
								OldLine:   1,
								NewLine:   1,
								IsEndedLF: false,
							},
						},
						OldLineFrom: 1,
						OldLineTo:   1,
						NewLineFrom: 1,
						NewLineTo:   1,
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
			t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Text:      "line1",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: false,
							},
							{
								Operation: OperationInsert,
								Text:      "line2",
								OldLine:   1,
								NewLine:   1,
								IsEndedLF: false,
							},
						},
						OldLineFrom: 1,
						OldLineTo:   1,
						NewLineFrom: 1,
						NewLineTo:   1,
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
			t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
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
			t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationInsert,
								Text:      "line1\n",
								OldLine:   0,
								NewLine:   1,
								IsEndedLF: true,
							},
							{
								Operation: OperationInsert,
								Text:      "line2\n",
								OldLine:   0,
								NewLine:   2,
								IsEndedLF: true,
							},
							{
								Operation: OperationInsert,
								Text:      "line3\n",
								OldLine:   0,
								NewLine:   3,
								IsEndedLF: true,
							},
						},
						OldLineFrom: 0,
						OldLineTo:   0,
						NewLineFrom: 1,
						NewLineTo:   3,
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
			t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
				expected := []Hunk{
					{
						Diffs: []LineDiff{
							{
								Operation: OperationDelete,
								Text:      "line1\n",
								OldLine:   1,
								NewLine:   0,
								IsEndedLF: true,
							},
							{
								Operation: OperationDelete,
								Text:      "line2\n",
								OldLine:   2,
								NewLine:   0,
								IsEndedLF: true,
							},
							{
								Operation: OperationDelete,
								Text:      "line3\n",
								OldLine:   3,
								NewLine:   0,
								IsEndedLF: true,
							},
						},
						OldLineFrom: 1,
						OldLineTo:   3,
						NewLineFrom: 0,
						NewLineTo:   0,
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

		t.Run("n=0", func(t *testing.T) {
			n := 0
			expected := []Hunk{
				{
					Diffs: []LineDiff{
						{
							Operation: OperationDelete,
							Text:      "line5\n",
							OldLine:   5,
							NewLine:   4,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line5b\n",
							OldLine:   5,
							NewLine:   5,
							IsEndedLF: true,
						},
					},
					OldLineFrom: 5,
					OldLineTo:   5,
					NewLineFrom: 5,
					NewLineTo:   5,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationDelete,
							Text:      "line8\n",
							OldLine:   7,
							NewLine:   6,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line7b\n",
							OldLine:   7,
							NewLine:   7,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line8b\n",
							OldLine:   7,
							NewLine:   8,
							IsEndedLF: true,
						},
					},
					OldLineFrom: 7,
					OldLineTo:   7,
					NewLineFrom: 7,
					NewLineTo:   8,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationInsert,
							Text:      "line13\n",
							OldLine:   11,
							NewLine:   13,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line14",
							OldLine:   11,
							NewLine:   14,
							IsEndedLF: false,
						},
					},
					OldLineFrom: 11,
					OldLineTo:   0,
					NewLineFrom: 13,
					NewLineTo:   14,
				},
			}
			doTest(t, testArg{
				a:        a,
				b:        b,
				n:        n,
				expected: expected,
			})
		})

		t.Run("n=1", func(t *testing.T) {
			n := 1

			// Merge hunks that are adjacent
			expected := []Hunk{
				{
					Diffs: []LineDiff{
						{
							Operation: OperationEqual,
							Text:      "line4\n",
							OldLine:   4,
							NewLine:   4,
							IsEndedLF: true,
						},
						{
							Operation: OperationDelete,
							Text:      "line5\n",
							OldLine:   5,
							NewLine:   4,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line5b\n",
							OldLine:   5,
							NewLine:   5,
							IsEndedLF: true,
						},
						{
							Operation: OperationEqual,
							Text:      "line6\n",
							OldLine:   6,
							NewLine:   6,
							IsEndedLF: true,
						},
						{
							Operation: OperationDelete,
							Text:      "line8\n",
							OldLine:   7,
							NewLine:   6,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line7b\n",
							OldLine:   7,
							NewLine:   7,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line8b\n",
							OldLine:   7,
							NewLine:   8,
							IsEndedLF: true,
						},
						{
							Operation: OperationEqual,
							Text:      "line9\n",
							OldLine:   8,
							NewLine:   9,
							IsEndedLF: true,
						},
					},
					OldLineFrom: 4,
					OldLineTo:   8,
					NewLineFrom: 4,
					NewLineTo:   9,
				},
				{
					Diffs: []LineDiff{
						{
							Operation: OperationEqual,
							Text:      "line12\n",
							OldLine:   11,
							NewLine:   12,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line13\n",
							OldLine:   11,
							NewLine:   13,
							IsEndedLF: true,
						},
						{
							Operation: OperationInsert,
							Text:      "line14",
							OldLine:   11,
							NewLine:   14,
							IsEndedLF: false,
						},
					},
					OldLineFrom: 11,
					OldLineTo:   11,
					NewLineFrom: 12,
					NewLineTo:   14,
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
