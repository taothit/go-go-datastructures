package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"taothit/ggd/model"
)

var path = flag.String("pathTo", wdOrEmpty(), "fully-qualified path to the datastructure's creation instruction file.")

func wdOrEmpty() string {
	if dir, err := os.Getwd(); err != nil {
		return ""
	} else {
		return dir
	}
}

// go:generate ggd -pathTo=templates/example/widgetStack.go stack[Widget]
// GGD follows the following steps to produce a custom datastructure
// 1. Parse flags (pathTo)
//  A. Present - uses provided path to datastructure's instruction/destination file
//  B. Absent - uses root directory of project
// 2. Parse instructions
//  A. Present - saves entityType and datastructure
//      I. fails when entityType or datastructure are absent; or, entityType isn't a valid identifier;
//         or, datastructure is unknown
//  B. Absent - fails
// 3. Finds instruction/destination file for datastructure in project
// 4. Parses AST from instruction/destination file
//  A. Present - Reads package name to use as package for custom datastructure
//  B. Absent - uses main
// 5. Looks up datastructure from instructions
// 6. Copies datastructure template from instructions.
// 7. Parses AST for template copy
// 8. Create ast.Package from template copy ast.File and i/d ast.File.
// 9. Merge ast.Package files
// 10. Walks merged ast.File
//  A. Replaces template entity's identifier (interface{}) with entityType from instructions
//  B. Replaces template datastructure's name with entityType+datastructure
//  C. Replaces pointers to template entity's identifier with pointers to entityType
//  D. Replaces datastructure's template entity with instruction's datastructure type
//  E. Preserves comments by replacing datastructure and entity in template with those in instructions
// 11. Write out custom datastructure to i/d file
func main() {
	flag.Parse()

	if *path == "" || *path == wdOrEmpty() {
		log.Fatalln("ggd: missing required path to datastructure source file.")
		// TODO(tad) - need to create usage and print it before panicking
	}

	// Read template instructions
	if len(flag.Args()) < 1 || flag.Arg(0) == "" {
		log.Fatalln("ggd: missing required datastructure creation directive.")
		// TODO(tad) - need to create usage and print it before panicking
	}
	instructions := flag.Arg(0)

	// Load project source files
	fSet := token.NewFileSet()
	file, first := parser.ParseFile(fSet, *path, nil, parser.AllErrors)
	if first != nil  {
		log.Fatalf("ggd: unparsable template file (%s): %v", *path, first)
	}
	files := map[string]*ast.File{file.Name.Name: file}

	// Find instruction file using *path
	pkg, err := ast.NewPackage(fSet, files, nil, nil)
	ds := model.NewDatastructure(instructions, pkg, file)
	if ds == nil {
		log.Fatalf("ggd: unknown datastructure for instructions (%s)", instructions)
	}

	dst, err := os.OpenFile(*path, os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		log.Fatalf("ggd: destination file unavailable: %v", err)
	}

	err = ds.Print(dst)
	if err != nil {
		log.Fatalf("ggd: failed creating custom datastructure (%s): %v", ds, err)
	}
}
