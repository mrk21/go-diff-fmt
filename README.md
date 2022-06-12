# go-diff-fmt

Diff formatting library for Go(e.g. Unified format)

## Install

```sh
$ go get github.com/mrk21/go-diff-fmt
```

## Usage

**main.go:**

```go
package main

import (
	"os"

	"github.com/mrk21/go-diff-fmt/difffmt"
)

func main() {
	targetA, _ := difffmt.NewDiffTarget("a.txt")
	targetB, _ := difffmt.NewDiffTarget("b.txt")
	textA, _ := targetA.ReadText()
	textB, _ := targetB.ReadText()

	lineDiffs := difffmt.GetLineDiff(textA, textB)
	hunks := difffmt.GetHunk(lineDiffs, 3)
	unified := difffmt.UnifiedFormat{IsColor: true}
	unified.Format(os.Stdout, targetA, targetB, hunks)
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
--- tmp/a.txt   2022-06-12 18:14:00.448904241 +0900
+++ tmp/b.txt   2022-06-12 21:52:43.768868544 +0900
@@ -2,10 +2,14 @@
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
+
```

## Refer to

- [Unified Diff Format](https://www.artima.com/weblogs/viewpost.jsp?thread=164293)
- [Detailed Unified (Comparing and Merging Files)](https://www.gnu.org/software/diffutils/manual/html_node/Detailed-Unified.html)
