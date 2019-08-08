package cmd

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"taothit/ggd/model"
)


func Generate(instructions, directiveFilePath string, out *io.Writer) {
	if out == nil {
		log.Println("datastructure output nil")
		return
	}
	// Load project source files
	fSet := token.NewFileSet()
	file, first := parser.ParseFile(fSet, directiveFilePath, nil, parser.AllErrors)
	if first != nil {
		log.Fatalf("ggd: unparsable template file (%s): %v", directiveFilePath, first)
	}
	files := map[string]*ast.File{file.Name.Name: file}
	// Find instruction file using *path
	pkg, err := ast.NewPackage(fSet, files, nil, nil)
	if err != nil {
		log.Fatalf("ggd: invalid package: %v", err)
	}

	ds := model.NewDatastructure(instructions, pkg, file)
	if ds == nil {
		log.Fatalf("ggd: unknown datastructure for instructions (%s)", instructions)
	}

	if err = ds.Print(out); err != nil {
		log.Fatalf("ggd: failed creating custom datastructure (%s): %v", ds, err)
	}
}
