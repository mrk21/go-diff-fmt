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
		contextSize = flag.Int("context-size", 3, "Context size")
		isColorize  = flag.Bool("color", false, "Enable colorize")
		isHelp      = flag.Bool("help", false, "Show usage")
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

	lineDiffs := difffmt.GetLineDiffsFromDMP(diffs)
	hunks := difffmt.GetHunks(lineDiffs, *contextSize)
	unifiedFmt := difffmt.UnifiedFormat{IsColorize: *isColorize}
	unifiedFmt.Print(targetA, targetB, hunks)
}
