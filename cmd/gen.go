package cmd

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"taothit/ggd/model"
)

func Generate(instructions string, directiveSource *Directives, out *io.Writer, mode model.LogMode) error {
	if out == nil {
		return errors.New("datastructure output nil")
	}
	// Load project source files
	fSet := token.NewFileSet()

	file := directiveSource.CreateAstFile(fSet, parser.AllErrors)

	files := map[string]*ast.File{file.Name.Name: file}
	// Find instruction file using *path
	pkg, err := ast.NewPackage(fSet, files, nil, nil)
	if err != nil {
		return fmt.Errorf("ggd: invalid package: %v", err)
	}

	ds := model.NewDatastructure(instructions, pkg, file, mode)
	if ds == nil {
		return fmt.Errorf("ggd: unknown datastructure for instructions (%s)", instructions)
	}

	if err = ds.Print(out); err != nil {
		return fmt.Errorf("ggd: failed creating custom datastructure (%s): %v", ds, err)
	}
	return nil
}

const defaultOutput = "tmp"

type Directives struct {
	src interface{}
	path string
	defaultPath bool
}

func RawSource(s string) (d *Directives) {
	if s != "" {
		d = &Directives{
			src: s,
			path: defaultOutput+".go",
			defaultPath: true,
		}
	}

	return
}

func File(path string) *Directives {
	if abs, err := filepath.Abs(path); err != nil {
		return nil
	} else {
		return &Directives{
			src: nil,
			path: abs,
		}

	}
}

func (d *Directives) CreateAstFile(fs *token.FileSet, mode parser.Mode) *ast.File {
	if d.defaultPath {
		_, err := os.Open(d.path)
		if err != nil {
			if _, err = os.Create(d.path); err != nil {
				return nil
			}
		}
	}

	file, first := parser.ParseFile(fs, d.path, d.src, mode)
	if first != nil {
		log.Printf("ggd: unparsable template file (%s) or sourc (%s): %v", d.path, d.src, first)
		return nil
	}

	return file
}

func (d *Directives) GoString() string {
	return fmt.Sprintf("Directive{%s, %s}", d.path, d.src)
}
