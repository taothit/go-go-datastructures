package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"taothit/go-go-datastructures/model"
)

var path = flag.String("pathTo", maybeWorkingDir(), "fully-qualified path to the datastructure's creation instruction file.")

func maybeWorkingDir() string {
	if dir, err := os.Getwd(); err != nil {
		return ""
	} else {
		return dir
	}
}

// go:generate ggd -pathTo=templates/example/widgetStack.go stack[Widget]
func main() {
	flag.Parse()

	if *path == "" {
		log.Fatalln("go-go-datastructures: missing required path to datastructure source file.")
		// TODO(tad) - need to create usage and print it before panicking
	}

	// directiveFileName := filepath.Base(*path)

	// Read template instructions
	if len(flag.Args()) < 1 || flag.Arg(0) == "" {
		log.Fatalln("go-go-datastructures: missing required datastructure creation directive.")
		// TODO(tad) - need to create usage and print it before panicking
	}

	instructions := flag.Arg(0)

	// Load project source files
	fSet := token.NewFileSet()

	file, first := parser.ParseFile(fSet, *path, nil, parser.AllErrors)
	if first != nil {
		log.Fatalf("failed to parse %s: %v", *path, first)
	}
	files := map[string]*ast.File{file.Name.Name: file}

	// Find instruction file using *path
	pkg, err := ast.NewPackage(fSet, files, nil, nil)
	ds := model.NewDatastructure(instructions, pkg, file)
	if ds == nil {
		log.Fatalf("failed to create datastructure for instructions (%s)", instructions)
	}

	err = ds.Print(os.Stdout)
	if err != nil {
		log.Fatalf("failed to create custom datastructure (%s): %v", ds, err)
	}
}
