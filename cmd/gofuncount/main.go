package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/yanana/gofuncount"
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
)

func initFlags() {
	flag.BoolVar(&flagIncludeTests, "include-tests", false, "include test files")
	flag.StringVar(&flagOutputFormat, "output-format", "json", "output format")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [-flag] [package]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Flags:")
		fmt.Fprintln(os.Stderr, "  -include-tests")
		fmt.Fprintln(os.Stderr, "  -output-format")
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

	conf := &gofuncount.Config{
		IncludeTests: flagIncludeTests,
	}

	counts, err := gofuncount.Run(args[0], conf)
	if err != nil {
		log.Printf("error: %s", err)

		return 2
	}

	out, err := output(counts, flagOutputFormat)
	if err != nil {
		log.Printf("error: %s", err)

		return 3
	}

	fmt.Fprintln(os.Stdout, out)

	return 0
}

func output(counts []gofuncount.CountInfo, format string) (string, error) {
	switch format {
	case "json":
		return outputJSON(counts)
	default:
		return "", fmt.Errorf("error: unknown output format: %s", flagOutputFormat)
	}
}

func outputJSON(counts []gofuncount.CountInfo) (string, error) {
	j, err := json.MarshalIndent(counts, "", "  ")
	if err != nil {
		return "", err
	}

	return string(j), nil
}
