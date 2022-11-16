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

func Run(root string, conf *Config) ([]CountInfo, error) {
	filter := func(fi fs.FileInfo) bool {
		if fi.IsDir() {
			return true
		}
		if conf.IncludeTests {
			return true
		}
		return !strings.HasSuffix(fi.Name(), "_test.go")
	}

	functions := make([]CountInfo, 0)

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			return nil
		}

		if d.Name() == "testdata" {
			return fs.SkipDir
		}

		fset := token.NewFileSet()
		fs, err := parseFilesInCurrentDir(fset, path, filter)
		if err != nil {
			return err
		}
		functions = append(functions, fs...)

		return nil
	})

	return functions, nil
}

func parseFilesInCurrentDir(fset *token.FileSet, root string, filter func(fi fs.FileInfo) bool) ([]CountInfo, error) {
	pkgs, err := parser.ParseDir(fset, root, filter, 0)
	if err != nil {
		return nil, err
	}

	functions := make([]CountInfo, 0)

	for name, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, d := range file.Decls {
				if f, ok := d.(*ast.FuncDecl); ok {
					functions = append(functions, CountInfo{
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

	return functions, nil
}
