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
	FormatJSON    Format = iota
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
	Format(Counts, bool) (io.Reader, error)
}

var _ Formatter = (*JSONFormatter)(nil)
var _ Formatter = (*CSVFormatter)(nil)

type JSONFormatter struct{}

func (f *JSONFormatter) Format(cs Counts, stats bool) (io.Reader, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetIndent("", "  ")

	encodee := any(cs)
	if stats {
		encodee = cs.Stats()
	}
	if err := e.Encode(encodee); err != nil {
		return nil, err
	}

	return &b, nil
}

type CSVFormatter struct{}

func (f *CSVFormatter) Format(cs Counts, stats bool) (io.Reader, error) {
	var b bytes.Buffer
	w := csv.NewWriter(&b)
	w.UseCRLF = true

	w.Write(cs.CSVHeader(stats))

	for _, funcs := range cs {
		for _, f := range funcs {
			w.Write(f.CSVRecord())
		}
	}
	w.Flush()

	return &b, nil
}

func (cs Counts) WriteCSV(w *csv.Writer, stats bool) error {
	w.Write(cs.CSVHeader(stats))

	if stats {
		if err := cs.writeStatsCSV(w); err != nil {
			return err
		}
	} else {
		if err := cs.writeCSV(w); err != nil {
			return err
		}
	}

	w.Flush()

	return w.Error()
}

func (cs Counts) writeCSV(w *csv.Writer) error {
	for _, funcs := range cs {
		for _, f := range funcs {
			if err := w.Write(f.CSVRecord()); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cs Counts) writeStatsCSV(w *csv.Writer) error {
	s := cs.Stats()
	for pkg, as := range s {
		if err := w.Write(as.CSVRecord(pkg)); err != nil {
			return err
		}
	}

	return nil
}

func (s *Stats) CSVRecord(pkg string) []string {
	return []string{
		pkg,
		strconv.FormatFloat(s.MeanLines, 'f', -1, 64),
		strconv.FormatFloat(s.MedianLines, 'f', -1, 64),
		strconv.FormatFloat(s.NinetyFivePercentileLines, 'f', -1, 64),
		strconv.FormatFloat(s.NinetyNinePercentileLines, 'f', -1, 64)}
}

func (c *CountInfo) CSVRecord() []string {
	return []string{c.Package, c.Name, c.FileName, strconv.Itoa(c.StartsAt), strconv.Itoa(c.EndsAt), strconv.Itoa(c.Lines())}
}

func (cs Counts) CSVHeader(stats bool) []string {
	if stats {
		return []string{"package", "mean", "median", "95%ile", "99%ile"}
	}
	return []string{"package", "functionName", "fileName", "startsAt", "endsAt", "lines"}
}
