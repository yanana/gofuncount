package gofuncount

import (
	"encoding/json"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strings"
)

type Config struct {
	IncludeTests bool
}

type CountInfo struct {
	Package  string `json:"package"`
	Name     string `json:"functionName"`
	FileName string `json:"fileName"`
	StartsAt int    `json:"startsAt"`
	EndsAt   int    `json:"endsAt"`
}

func (f CountInfo) Lines() int {
	return f.EndsAt - f.StartsAt
}

func (f CountInfo) MarshalJSON() ([]byte, error) {
	type FI CountInfo

	return json.Marshal(struct {
		FI
		Lines int `json:"lines"`
	}{
		FI:    FI(f),
		Lines: f.Lines(),
	})
}

type Stats struct {
	MeanLines                 float64 `json:"mean"`
	MedianLines               float64 `json:"median"`
	NinetyFivePercentileLines float64 `json:"95%ile"`
	NinetyNinePercentileLines float64 `json:"99%ile"`
}

type Counts map[string][]*CountInfo

func (cs Counts) Stats() map[string]*Stats {
	var statsPerPackage = make(map[string]*Stats, len(cs))

	for pkg, counts := range cs {
		ss := &Stats{}
		statsPerPackage[pkg] = ss
		lines := make([]int, 0, len(counts))
		if len(counts) == 0 {
			continue
		}
		for _, f := range counts {
			lines = append(lines, f.Lines())
		}
		d := NewData(lines)
		ss.MeanLines = d.Mean()
		ss.MedianLines = d.Quantile(0.5)
		ss.NinetyFivePercentileLines = d.Quantile(0.95)
		ss.NinetyNinePercentileLines = d.Quantile(0.99)
	}

	return statsPerPackage
}

func Run(root string, conf *Config) (Counts, error) {
	filter := func(fi fs.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		if conf.IncludeTests {
			return true
		}
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}

	counts := make(Counts)

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}

		if d.Name() == "testdata" {
			return fs.SkipDir
		}

		fset := token.NewFileSet()
		if err := parseFilesInCurrentDir(fset, path, filter, counts); err != nil {
			return err
		}

		return nil
	})

	return counts, nil
}

func parseFilesInCurrentDir(fset *token.FileSet, root string, filter func(fi fs.FileInfo) bool, counts Counts) error {
	pkgs, err := parser.ParseDir(fset, root, filter, 0)
	if err != nil {
		return err
	}

	for name, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, d := range file.Decls {
				if f, ok := d.(*ast.FuncDecl); ok {
					counts[pkg.Name] = append(counts[pkg.Name], &CountInfo{
						Package:  name,
						Name:     f.Name.Name,
						FileName: fset.File(f.Pos()).Name(),
						StartsAt: fset.Position(f.Pos()).Line,
						EndsAt:   fset.Position(f.End()).Line,
					})
				}
			}
		}
	}

	return nil
}
