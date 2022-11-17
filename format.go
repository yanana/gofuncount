package gofuncount

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"io"
	"strconv"
	"strings"
)

type Format int8

const (
	FormatUnknown Format = iota - 1
	FormatJSON
	FormatCSV
)

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "json"
	case FormatCSV:
		return "csv"
	default:
		return "unknown"
	}
}

func (f Format) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Format) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "json":
		*f = FormatJSON
	case "csv":
		*f = FormatCSV
	default:
		*f = FormatUnknown
	}

	return nil
}

func ParseFormat(s string) (Format, error) {
	var f Format

	if err := f.UnmarshalText([]byte(s)); err != nil {
		return FormatUnknown, err
	}

	return f, nil
}

func FormatterOf(f Format) Formatter {
	switch f {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatCSV:
		return &CSVFormatter{}
	default:
		return &JSONFormatter{}
	}
}

type Formatter interface {
	Format(Counts) (io.Reader, error)
}

var _ Formatter = (*JSONFormatter)(nil)
var _ Formatter = (*CSVFormatter)(nil)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(cs Counts) (io.Reader, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetIndent("", "  ")

	if err := e.Encode(cs); err != nil {
		return nil, err
	}

	return &b, nil
}

type CSVFormatter struct{}

func (f *CSVFormatter) Format(cs Counts) (io.Reader, error) {
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	w.UseCRLF = true

	w.Write(cs.CSVHeader())

	for pkg, funcs := range cs {
		for _, f := range funcs {
			w.Write([]string{
				pkg, f.Name, f.FileName, strconv.Itoa(f.StartsAt), strconv.Itoa(f.EndsAt), strconv.Itoa(f.Lines()),
			})
		}
	}
	w.Flush()

	return &b, nil
}

func (cs Counts) CSVHeader() []string {
	var header = []string{"package", "functionName", "fileName", "startsAt", "endsAt", "lines"}

	return header
}
