package main

import (
	"flag"
	"fmt"
	"github.com/yanana/gofuncount"
	"io"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)

	initFlags()
}

var (
	flagIncludeTests bool
	flagOutputFormat string
	flagOutputStats  bool
)

func initFlags() {
	flag.BoolVar(&flagIncludeTests, "include-tests", false, "whether to include test files")
	flag.StringVar(&flagOutputFormat, "format", "json", "output format, either one of json or csv")
	flag.BoolVar(&flagOutputStats, "stats", false, "whether to output statistics")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-flag] target\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Counts the number of functions in or of the given `target`.")
		fmt.Fprintln(os.Stderr, "The `target` can be either a directory or a file.")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
	}
}

func main() {
	os.Exit(doMain())
}

func doMain() int {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		return 1
	}

	format, err := gofuncount.ParseFormat(flagOutputFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid format: %q", flagOutputFormat)
		return 2
	}

	conf := &gofuncount.Config{
		IncludeTests: flagIncludeTests,
	}

	runner := &gofuncount.Runner{
		Conf: conf,
	}
	counts, err := runner.Run(args[0])
	if err != nil {
		log.Printf("error: %s", err)
		return 3
	}

	formatter := gofuncount.FormatterOf(format)
	reader, err := formatter.Format(counts, flagOutputStats)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return 4
	}

	if _, err := io.Copy(os.Stdout, reader); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return 5
	}

	return 0
}
