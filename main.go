//go:generate stringer -type datastructureType
package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"regexp"
)

var directiveMask = regexp.MustCompile(`^([a-zA-Z]\w+)\[([a-zA-Z]\w+)\]$`)

var path = flag.String("pathTo", maybeWorkingDir(), "fully-qualified path to the datastructure's creation instruction file.")

func maybeWorkingDir() string {
	if dir, err := os.Getwd(); err != nil {
		return ""
	} else {
		return dir
	}
}

// datastructure identifies the templates available from which to create a type-specific collection.
type datastructureType int

const (
	// Enumerates the available datastructure templates
	unknown datastructureType = iota
	stack
	heap
)

// DatastructureTemplate caches datastructure templates and provides retrieval of new copies to be
// used to create custom datastructures
//
// Note: templates are stored as []byte representation of the source file for the datastructure after
//       performing the file.Read([]byte) operation.
type DatastructureTemplate [][]byte

// TODO(tad): create generator to produce datastructure set from datastructure values.
var datastructures = [...]string{"unknown", "stack", "heap"}
var dsPaths = [...]string{"", "templates/stack.go", ""}

func (d *DatastructureTemplate) copy(ds datastructureType) []byte {
	if d == nil {
		return nil
	}

	log.Printf("%v", d)
	tmpl := (*d)[ds]
	if tmpl == nil {
		return nil
	}

	buf := make([]byte, cap((*d)[ds]))
	l := copy(buf, tmpl)
	if l != len(tmpl) {
		return nil
	}

	return buf
}

func NewDatastructureTemplate() *DatastructureTemplate {
	if templates != nil {
		return templates
	}

	tmpls := DatastructureTemplate(make([][]byte, 0))

	for i, path := range dsPaths {
		if path != "" {
			file, err := os.OpenFile(path, os.O_RDWR, 0777)
			if err != nil {
				log.Printf("failed to open template file (%s): %v", path, err)
				tmpls = append(tmpls, []byte(datastructures[i]))
				continue
			}

			buf := make([]byte, 0)
			l, err := file.Read(buf)
			if l < 1 || len(buf) == 0 || err != nil {
				log.Printf("failed to read complete template file (%s): %v", path, err)
				tmpls = append(tmpls, []byte(datastructures[i]))
				continue
			}

			tmpls = append(tmpls, []byte(datastructures[i]))
		}
	}

	return &tmpls
}

var templates *DatastructureTemplate

// go:generate ggd -pathTo=templates/example/widgetStack.go stack[Widget]
func main() {
	flag.Parse()

	if *path == "" {
		panic(fmt.Sprint("go-go-datastructures: missing required path to datastructure source file."))
		// TODO(tad) - need to create usage and print it before panicking
	}

	// directiveFileName := filepath.Base(*path)

	// Read template instructions
	if len(flag.Args()) < 1 || flag.Arg(0) == "" {
		log.Print(fmt.Sprint("go-go-datastructures: missing required datastructure creation directive."))
		os.Exit(1)
		// TODO(tad) - need to create usage and print it before panicking
	}

	templates = NewDatastructureTemplate()

	instructions := flag.Arg(0)

	// Load project source files
	fSet := token.NewFileSet()

	file, first := parser.ParseFile(fSet, *path, nil, parser.AllErrors)
	if first != nil {
		log.Fatalf("failed to parse %s: %v", *path, first)
	}

	files := make(map[string]*ast.File, 0)
	files[file.Name.Name] = file
	// Find instruction file using *path
	pkg, err := ast.NewPackage(fSet, files, nil, nil)
	ds := NewDatastructure(instructions, pkg, file)
	if ds == nil {
		log.Fatalf("failed to find instruction file (%s)", *path)
	}

	err = ds.Print(os.Stdout)
	if err != nil {
		log.Fatalf("failed to create custom datastructure (%s): %v", ds.instruction.dsType, err)
	}
}

func NewDatastructure(instructions string, pkg *ast.Package, instructionFile *ast.File) *Datastructure {
	instr := NewInstruction(instructions, instructionFile)
	tmpl, pErr := parser.ParseFile(token.NewFileSet(), "", templates.copy(instr.dsType), parser.AllErrors)
	if pErr != nil {
		log.Printf("failed to load template datastructure (%s): %v", instr.dsType, pErr)
		return nil
	}

	return &Datastructure{
		instruction: instr,
		pkg:         pkg,
		destination: pkg.Files[*path],
		templateAst: tmpl,
	}
}

type Instruction struct {
	dsType          datastructureType
	entityType      string
	instructionFile *ast.File
}

func NewInstruction(instructions string, instructionFile *ast.File) *Instruction {
	dsType, entityType := parse(instructions)

	return &Instruction{
		instructionFile: instructionFile,
		dsType:          dsType,
		entityType:      entityType,
	}
}

func (i *Instruction) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.Ident:
		log.Println("Inspecting *ast.Ident...")
		if v.Name == "StackTemplate" {
			log.Println("Replacing identifier name...")
			v.Name = i.entityType
		}
		if v.Obj != nil && v.Obj.Kind == ast.Typ && v.Obj.Name == "StackTemplate" {
			log.Println("Replacing type name...")
			v.Obj.Name = i.entityType
		}
	default:
		// log.Println("found other node; moving on...")
		return i
	}

	return i
}

func parse(instructions string) (datastructureType, string) {
	matches := directiveMask.FindAllStringSubmatch(instructions, -1)

	if matches == nil || len(matches) < 1 {
		return unknown, ""
	}

	sub := matches[0]
	if sub[0] == instructions {
		sub = sub[1:]
	}

	for i, dsType := range datastructures {
		if dsType == sub[0] {
			return datastructureType(i), sub[1]
		}
	}

	return unknown, ""
}

type Datastructure struct {
	instruction *Instruction
	pkg         *ast.Package
	destination *ast.File
	templateAst *ast.File
}

// Print the datastructure to the provided Writer
func (d *Datastructure) Print(w io.Writer) error {
	// Write datastructure to file
	d.fillTemplate()
	log.Println("Printing tree...")
	if err := ast.Fprint(w, nil, d.templateAst, nil); err != nil {
		return fmt.Errorf("failed to write datastructure AST to writer: %v", err)
	}

	// Print the AST for file to w and log it
	// fSet := token.NewFileSet()
	// ast.Fprint(w, fSet, file, nil)
	return nil
}

func (d *Datastructure) fillTemplate() (*ast.File, error) {

	// Walk the AST and replace interface{} with the new type
	log.Println("Walking tree...")
	ast.Walk(d.instruction, d.templateAst)

	return d.templateAst, nil
}