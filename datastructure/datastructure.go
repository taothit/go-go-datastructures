//go:generate stringer -type Type
package datastructure

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Type identifies the templates available from which to create a type-specific collection.
type Type int

const (
	// Unknown is the undefined Type
	Unknown Type = iota
	// Stack provides LIFO management of homogenous items
	Stack
	// Heap provides binary tree managemange of homogenous itemss
	Heap
)

// Template caches datastructure templates and provides retrieval of new copies to be
// used to create custom datastructures
//
// Note: templates are stored as []byte representation of the source file for the datastructure after
//       performing the file.Read([]byte) operation.
type Template [][]byte

var templates *Template

// TODO(tad): create generator to produce datastructure set from datastructure values.
var datastructures = [...]string{"Unknown", "Stack", "Heap"}
var dsPaths = [...]string{"", "templates/stack.go", ""}

var UnknownType = errors.New("Unknown datastructure type")
var TemplatesUninitialized = errors.New("Datastrucutre templates not initialized")

// CopyTo writes the full type-less Type to the provided Writer.
// error returned when Datastructure type is Unknown or the template is incompletely written
func (ds Type) CopyTo(w io.Writer) error {
	if templates == nil {
		return fmt.Errorf("failed writing Datastructure (%s) template: %w", ds, TemplatesUninitialized)
	}

	if ds == Unknown {
		return fmt.Errorf("failed writing Datastructure (%s) template: %w", ds, UnknownType)
	}

	tmplBuf := (*templates)[ds]
	if n, err := w.Write(tmplBuf); len(tmplBuf) > 0 {
		if n == 0 || err != nil {
			return fmt.Errorf("failed writing Datastructure (%s) template: %v", ds, err)
		}
	}

	return nil
}

func (d *Template) CopyTo(ds Type) []byte {
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

func LoadTemplates() {
	if templates == nil {
		tmpls := Template(make([][]byte, 0))
		templates = &tmpls

		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("ggd: unknown working directory: %v", err)
		}

		for i, path := range dsPaths {
			if path != "" {
				log.Printf("loading datastructure from template (%s)", path)
				pathInWd := filepath.Join(wd, path)
				log.Printf("wd+path: %s", pathInWd)
				file, err := os.OpenFile(pathInWd, os.O_RDWR, 0777)

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
	}
}

func NewDatastructure(instructions string, pkg *ast.Package, instructionFile *ast.File) *Datastructure {
	instr := parse(instructions)
	var buf bytes.Buffer
	if err := instr.dsType.CopyTo(&buf); err != nil {
		return nil
	}

	tmpl, pErr := parser.ParseFile(token.NewFileSet(), "", buf, parser.AllErrors)
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
	if err := ast.Fprint(os.Stdout, nil, d.templateAst, nil); err != nil {
		return fmt.Errorf("failed to write datastructure AST to writer: %v", err)
	}

	cfg := printer.Config{Mode: printer.UseSpaces | printer.SourcePos, Indent: 0, Tabwidth: 4}
	fErr := cfg.Fprint(w, token.NewFileSet(), d.templateAst)
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
