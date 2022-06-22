package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mrk21/go-diff-fmt/difffmt"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	var (
		contextSize         = flag.Int("context-size", 3, "Context size")
		format              = flag.String("format", "unified", "Format type(unified,context)")
		isColor             = flag.Bool("color", false, "Enable color")
		isForceColor        = flag.Bool("force-color", false, "Enable color even other than the terminal")
		isHelp              = flag.Bool("help", false, "Show usage")
		isHidingNoLFMessage = flag.Bool("hide-no-lf", false, "Hide a no LF message")
	)
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Usage: godiff [-c]\n")
		fmt.Fprintf(o, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *isHelp {
		flag.CommandLine.Usage()
		os.Exit(2)
	}
	if *contextSize < 0 {
		flag.CommandLine.Usage()
		os.Exit(2)
	}

	colorMode := difffmt.ColorNone
	if *isColor {
		colorMode = difffmt.ColorTerminalOnly
	}
	if *isForceColor {
		colorMode = difffmt.ColorAlways
	}

	args := flag.Args()
	if len(args) != 2 {
		flag.CommandLine.Usage()
		os.Exit(2)
	}
	pathA := args[0]
	pathB := args[1]

	targetA := difffmt.NewDiffTarget(pathA)
	targetB := difffmt.NewDiffTarget(pathB)

	err := targetA.LoadStat()
	if err != nil {
		log.Fatalln(err)
	}
	err = targetB.LoadStat()
	if err != nil {
		log.Fatalln(err)
	}

	textA, err := targetA.ReadText()
	if err != nil {
		log.Fatalln(err)
	}
	textB, err := targetB.ReadText()
	if err != nil {
		log.Fatalln(err)
	}

	dmp := diffmatchpatch.New()
	runes1, runes2, lineArray := dmp.DiffLinesToRunes(string(textA), string(textB))
	diffs := dmp.DiffMainRunes(runes1, runes2, false)
	diffs = dmp.DiffCharsToLines(diffs, lineArray)

	lineDiffs := difffmt.MakeLineDiffsFromDMP(diffs)
	hunks := difffmt.MakeHunks(lineDiffs, *contextSize)

	switch *format {
	case "unified":
		unifiedFmt := difffmt.NewUnifiedFormat(difffmt.UnifiedFormatOption{
			ColorMode:           colorMode,
			IsHidingNoLFMessage: *isHidingNoLFMessage,
		})
		unifiedFmt.Print(targetA, targetB, hunks)
	case "context":
		contextFmt := difffmt.NewContextFormat(difffmt.ContextFormatOption{
			ColorMode:           colorMode,
			IsHidingNoLFMessage: *isHidingNoLFMessage,
		})
		contextFmt.Print(targetA, targetB, hunks)
	}
}
