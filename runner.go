package gofuncount

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	IncludeTests bool
}

type CountInfo struct {
	Package  string `json:"package"`
	Name     string `json:"function_name"`
	FileName string `json:"file_name"`
	StartsAt int    `json:"starts_at"`
	EndsAt   int    `json:"ends_at"`
}

func (ci CountInfo) Lines() int {
	return ci.EndsAt - ci.StartsAt
}

func (ci CountInfo) MarshalJSON() ([]byte, error) {
	type FI CountInfo

	return json.Marshal(struct {
		FI
		Lines int `json:"lines"`
	}{
		FI:    FI(ci),
		Lines: ci.Lines(),
	})
}

type Stats struct {
	MeanLines                 float64 `json:"mean"`
	MedianLines               float64 `json:"median"`
	NinetyFivePercentileLines float64 `json:"95%ile"`
	NinetyNinePercentileLines float64 `json:"99%ile"`
}

// Counts is a map of root directory to a list of CountInfo.
type Counts map[string][]*CountInfo

func (cs Counts) Stats() map[string]*Stats {
	var statsPerDirectory = make(map[string]*Stats, len(cs))

	for dir, counts := range cs {
		ss := &Stats{}
		statsPerDirectory[dir] = ss
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

	return statsPerDirectory
}

type Runner struct {
	Conf *Config
}

func (r *Runner) Run(root string) (Counts, error) {

	counts := make(Counts)

	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return r.runFile(root, fi, counts)
	}

	return r.runDir(root, counts)
}

func (r *Runner) runDir(root string, counts Counts) (Counts, error) {
	filter := func(fi fs.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		if r.Conf.IncludeTests {
			return true
		}
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}

	if err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		switch {
		case d == nil:
			return fmt.Errorf("no such file or directory: %s", path)
		case err != nil:
			return err
		case d.IsDir():
			if d.Name() == "testdata" {
				return fs.SkipDir
			}
			if err := parseFilesInCurrentDir(token.NewFileSet(), path, filter, counts); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return counts, nil
}

func (r *Runner) runFile(root string, fi os.FileInfo, counts Counts) (Counts, error) {
	if !strings.HasSuffix(fi.Name(), ".go") {
		return nil, fmt.Errorf("not a go file: %s", root)
	}
	if err := parseFile(token.NewFileSet(), root, counts); err != nil {
		return nil, err
	}
	return counts, nil
}

func parseFile(fset *token.FileSet, path string, counts Counts) error {
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return err
	}

	countFile(fset, f, counts, filepath.Dir(path), f.Name.Name)

	return nil
}

type filterFunc func(fi fs.FileInfo) bool

func parseFilesInCurrentDir(fset *token.FileSet, path string, filter filterFunc, counts Counts) error {
	pkgs, err := parser.ParseDir(fset, path, filter, 0)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			countFile(fset, file, counts, path, pkg.Name)
		}
	}

	return nil
}

func countFile(fset *token.FileSet, file *ast.File, counts Counts, path string, pkg string) {
	for _, d := range file.Decls {
		if f, ok := d.(*ast.FuncDecl); ok {
			counts[path] = append(counts[path], &CountInfo{
				Package:  pkg,
				Name:     f.Name.Name,
				FileName: fset.File(f.Pos()).Name(),
				StartsAt: fset.Position(f.Pos()).Line,
				EndsAt:   fset.Position(f.End()).Line,
			})
		}
	}
}
