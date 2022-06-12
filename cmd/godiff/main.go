package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mrk21/go-diff-fmt/difffmt"
)

func main() {
	var (
		isColorize  = flag.Bool("c", false, "Enable colorize")
		contextSize = flag.Int("n", 3, "Output NUM lines of unified context")
	)
	flag.CommandLine.Usage = func() {
		o := flag.CommandLine.Output()
		fmt.Fprintf(o, "Usage: godiff [-c]\n")
		fmt.Fprintf(o, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
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

	targetA, err := difffmt.NewDiffTarget(pathA)
	if err != nil {
		log.Fatalln(err)
	}
	targetB, err := difffmt.NewDiffTarget(pathB)
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

	lineDiffs := difffmt.GetLineDiff(textA, textB)
	hunks := difffmt.GetHunk(lineDiffs, *contextSize)
	unified := difffmt.UnifiedFormat{IsColor: *isColorize}
	unified.Format(os.Stdout, targetA, targetB, hunks)
}
