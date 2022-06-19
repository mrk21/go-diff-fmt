# go-diff-fmt

[![Go Reference](https://pkg.go.dev/badge/github.com/mrk21/go-diff-fmt.svg)](https://pkg.go.dev/github.com/mrk21/go-diff-fmt)

Diff formatting library for Go.

For example, it can format line based diff by [github.com/sergi/go-diff](https://github.com/sergi/go-diff) to Unified format.

**NOTICE:**

`DiffLinesToChars` of `github.com/sergi/go-diff@v1.2.0` [contains bugs](https://github.com/sergi/go-diff/issues/123), so in order to work properly this module, it requires `v1.1.0`.

## Install

```sh
$ go get github.com/mrk21/go-diff-fmt/...
```

## Usage

**main.go:**

```go
package main

import (
	"github.com/mrk21/go-diff-fmt/difffmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	targetA := difffmt.NewDiffTarget("a.txt")
	targetB := difffmt.NewDiffTarget("b.txt")
	_ = targetA.LoadStat()
	_ = targetB.LoadStat()
	textA, _ := targetA.ReadText()
	textB, _ := targetB.ReadText()

	// Computing a line-mode diff
	// @see https://github.com/google/diff-match-patch/wiki/Line-or-Word-Diffs
	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(textA, textB)
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	// Format `[]diffmatchpatch.Diff` to Unified format
	lineDiffs := difffmt.MakeLineDiffsFromDMP(diffs)
	hunks := difffmt.MakeHunks(lineDiffs, 3)
	unifiedFmt := difffmt.NewUnifiedFormat(difffmt.UnifiedFormatOption{
		ColorMode: difffmt.ColorTerminalOnly,
	})
	unifiedFmt.Print(targetA, targetB, hunks)
}
```

**a.txt:**

```
line1
line2
line3
line4
line5
line6
line8
line9
line10
line11
line12

```

**b.txt:**

```
line1
line2
line3
line4
line5b
line6
line7b
line8b
line9
line10
line11
line12
line13
line14

```

```sh
$ go run main.go
--- a.txt       2022-06-13 14:07:47.670396312 +0900
+++ b.txt       2022-06-13 14:08:03.390396269 +0900
@@ -2,10 +2,13 @@
 line2
 line3
 line4
-line5
+line5b
 line6
-line8
+line7b
+line8b
 line9
 line10
 line11
 line12
+line13
+line14
```

## References

- [Unified Diff Format](https://www.artima.com/weblogs/viewpost.jsp?thread=164293)
- [Detailed Unified (Comparing and Merging Files)](https://www.gnu.org/software/diffutils/manual/html_node/Detailed-Unified.html)
- [diff - Wikipedia](https://en.wikipedia.org/wiki/Diff)
- [Line or Word Diffs · google/diff-match-patch Wiki](https://github.com/google/diff-match-patch/wiki/Line-or-Word-Diffs)
