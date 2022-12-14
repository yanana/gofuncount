# gofuncount

gofuncount is a command-line tool to count your Go source files' lines of code per a function. It also provides functionalities to calculate data aggregation like median or 95 percentile values of lines of functions per a file.

## Usage

```
Usage: gofuncount [-flag] target

Counts the number of functions in or of the given `target`.
The `target` can be either a directory or a file.

Flags:
  -format string
    	output format, either one of json or csv (default "json")
  -include-tests
    	whether to include test files
  -stats
    	whether to output statistics
```

The `-format` flag chooses the format of the output. Defaults to JSON but CSV is also available.
The `-include-tests` flag controls whether gofuncount counts test files (`*_test.go` files).
The `-stats` flag is used to output the statistics of the couts, such as median or 90 percentile.
