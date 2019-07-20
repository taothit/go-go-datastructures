//go:generate stringer -type DatastructureType
package model

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
)

// datastructure identifies the templates available from which to create a type-specific collection.
type DatastructureType int

const (
	// Enumerates the available datastructure templates
	Unknown DatastructureType = iota
	Stack
	Heap
)

// DatastructureTemplate caches datastructure templates and provides retrieval of new copies to be
// used to create custom datastructures
//
// Note: templates are stored as []byte representation of the source file for the datastructure after
//       performing the file.Read([]byte) operation.
type DatastructureTemplate [][]byte

// TODO(tad): create generator to produce datastructure set from datastructure values.
var datastructures = [...]string{"Unknown", "Stack", "Heap"}
var dsPaths = [...]string{"", "../templates/stack.go", ""}

func (d *DatastructureTemplate) Copy(ds DatastructureType) []byte {
	if d == nil {
		return nil
	}

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
	tmpls := DatastructureTemplate(make([][]byte, 0))

	for i, path := range dsPaths {
		if path != "" {
			log.Printf("loading datastructure from template (%s)", path)
			file, err := os.OpenFile(path, os.O_RDWR, 0777)
			if err != nil {
				log.Printf("failed to open template file (%s): %v", path, err)
				tmpls = append(tmpls, []byte(datastructures[i]))
				continue
			}
			stat, err := file.Stat()
			if err != nil {
				log.Printf("failed to retrieve info for file (%s): %v", path, err)
			}

			buf := make([]byte, stat.Size())
			l, err := file.Read(buf)
			if l < 1 || len(buf) == 0 || err != nil {
				log.Printf("failed to read complete template file (%s) (read=%d/%d): %v", path, l, stat.Size(), err)
				tmpls = append(tmpls, []byte(datastructures[i]))
				continue
			}

			log.Printf("read %d (of %d) bytes", l, stat.Size())
			tmpls = append(tmpls, buf)
		} else {
			tmpls = append(tmpls, []byte(datastructures[i]))
		}
	}

	return &tmpls
}

var templates *DatastructureTemplate

func init() {
	templates = NewDatastructureTemplate()
}

func NewDatastructure(instructions string, pkg *ast.Package, instructionFile *ast.File) *Datastructure {
	instr := parse(instructions)
	tmpl, pErr := parser.ParseFile(token.NewFileSet(), "", templates.Copy(instr.dsType), parser.AllErrors)
	if pErr != nil {
		log.Printf("failed to load template datastructure (%s): %v", instr.dsType, pErr)
		return nil
	}

	return &Datastructure{
		instruction: instr,
		pkg:         pkg,
		destination: instructionFile,
		templateAst: tmpl,
	}
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
	d.replaceInTemplate()
	log.Println("Printing tree...")
	if err := ast.Fprint(w, nil, d.templateAst, nil); err != nil {
		return fmt.Errorf("failed to write datastructure AST to writer: %v", err)
	}

	cfg := printer.Config{Mode: printer.UseSpaces | printer.SourcePos, Indent: 0, Tabwidth: 4}
	fErr := cfg.Fprint(os.Stdout, token.NewFileSet(), d.templateAst)
	if fErr != nil {
		return fmt.Errorf("failed to write custom datastructure source file (%s): %v", d.instruction.dsName(), fErr)
	}
	return nil
}

func (d *Datastructure) replaceInTemplate() {
	// Walk the AST and replace interface{} with the new type
	log.Println("Walking tree...")
	ast.Walk(d.instruction, d.templateAst)
}

func (d *Datastructure) String() string {
	return fmt.Sprintf("%s[%s]", d.instruction.dsType, d.instruction.entityType)
}