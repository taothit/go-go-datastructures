//go:generate stringer -type DatastructureType
package model

import (
    "fmt"
    "go/ast"
    "go/format"
    "go/parser"
    "go/token"
    "io"
    "log"
    "os"
    "path/filepath"
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
var dsPaths = [...]string{"", filepath.Join("..", "templates", "stack.go"), ""}

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

    wd, err := os.Getwd()
    if err != nil {
        log.Fatalf("ggd: unknown working directory: %v", err)
    }
    wd = filepath.Join("..", filepath.Base(wd))

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

    return &tmpls
}

var templates *DatastructureTemplate

func init() {
    templates = NewDatastructureTemplate()
}

func NewDatastructure(sourceFile *SourceFile, mode LogMode) *Datastructure {
    // Load project source files
    fSet := token.NewFileSet()

    sourceFile.mode = mode
    tmpl, pErr := parser.ParseFile(fSet, "", templates.Copy(sourceFile.dsType), parser.AllErrors|parser.ParseComments)
    if pErr != nil {
        log.Printf("failed to load template datastructure (%s): %v", sourceFile.dsType, pErr)
        return nil
    }

    return &Datastructure{
        srcFile:     sourceFile,
        templateAst: tmpl,
        mode:        mode,
    }
}

type Datastructure struct {
    srcFile     *SourceFile
    templateAst *ast.File
    mode        LogMode
}

func (d *Datastructure) inDebugMode() bool {
    return d != nil && d.mode == Noisy
}

func (d *Datastructure) debugln(msg string) {
    if d != nil && d.mode == Noisy {
        log.Println(msg)
    }
}

// Print datastructure to the provided io.Writer
func (d *Datastructure) Print(w *io.Writer) error {
    // Write datastructure to file
    srcAst, fileSet, err := d.Assemble()
    if err != nil {
        return fmt.Errorf("no datastructure to print: %v", err)
    }
    d.debugln("Printing tree...")
    if d.inDebugMode() {
        if err := ast.Fprint(os.Stdout, nil, d.templateAst, nil); err != nil {
            return fmt.Errorf("failed to write datastructure AST to writer: %v", err)
        }
    }

    if err := format.Node(*w, fileSet, srcAst); err != nil {
        return fmt.Errorf("failed to write custom datastructure source file (%s): %v", d.SourceFile().dsName(), err)
    }
    return nil
}

// Assemble walks the template file and replaces type and datastructure name.
func (d *Datastructure) Assemble() (
    srcAst *ast.File,
    fileSet *token.FileSet,
    err error) {
    fileSet = token.NewFileSet()
    f, pkg, err := d.srcFile.CreateAST(fileSet)
    if err != nil {
        return nil, nil, fmt.Errorf("datastructure source not assembled: %v", err)
    }

    d.srcFile.pkgName = f.Name.Name

    // Walk the AST and replace interface{} with the new type
    d.debugln("Walking tree...")

    ast.Walk(d.SourceFile(), d.templateAst)
    pkg.Files[dsPaths[d.SourceFile().dsType]] = d.templateAst
    pkg.Files[d.SourceFile().FileName()] = ast.MergePackageFiles(d.srcFile.pkg, ast.FilterFuncDuplicates|ast.FilterImportDuplicates|ast.FilterUnassociatedComments)

    return d.templateAst, fileSet, nil
}

func (d *Datastructure) String() string {
    return fmt.Sprintf("%s[%s]", d.srcFile.dsType, d.srcFile.entityType)
}

func (d *Datastructure) SourceFile() (s *SourceFile) {
    if d != nil {
        s = d.srcFile
    }

    return
}
