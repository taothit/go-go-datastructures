//go:generate stringer -type LogMode
package model

import (
    "fmt"
    "go/ast"
    "go/parser"
    "go/token"
    "log"
    "os"
    "regexp"
    "strings"
)

const defaultSourceFileName = "tmp"

var mask = regexp.MustCompile(`^([a-zA-Z]\w+)\[([a-zA-Z]\w+)\]$`)

// Directives instructs Datastructure the what/how/where for creating a custom datastructure
type Directives struct {
    dsType     DatastructureType
    entityType string
    mode       LogMode
    raw        string
}

func (d *Directives) GoString() string {
    return fmt.Sprintf("Directive{Datastructure: %s[%s], LogMode: %s (raw: %s)}", d.dsType, d.entityType, d.mode, d.raw)
}

func (d *Directives) Instructions() (raw string) {
    if d != nil {
        raw = d.raw
    }

    return
}

func (d *Directives) Datastructure() (t DatastructureType) {
    if d != nil {
        return d.dsType
    }

    return
}

func (src *SourceFile) LogMode(mode LogMode) {
    if src != nil {
        src.mode = mode
    }
}

type SourceFile struct {
    src         interface{}
    pkg         *ast.Package
    path        string
    dst         *os.File
    defaultPath bool
    dsType      DatastructureType
    entityType  string
    mode        LogMode
    raw         string
    pkgName     string
}

// CreateAST creates SourceFile's abstract syntax tree (AST) given a token.FileSet
func (src *SourceFile) CreateAST(fs *token.FileSet) (
    file *ast.File,
    pkg *ast.Package,
    err error) {
    if src.defaultPath {
        if _, err = os.Open(src.path); err != nil {
            if _, err = os.Create(src.path); err != nil {
                err = fmt.Errorf("source AST unavailable: %v", err)
                return
            }
        }
    }

    file, err = parser.ParseFile(fs, src.path, src.src, parser.AllErrors|parser.ParseComments)
    if err != nil {
        err = fmt.Errorf("unparsable source (%s): %v", src, err)
    }

    files := map[string]*ast.File{file.Name.Name: file}
    pkg, err = ast.NewPackage(fs, files, nil, nil)
    if err != nil {
        err = fmt.Errorf("unparsable template file (%s) or source (%s): %v", src.path, src.src, err)
        return
    }

    src.pkg = pkg

    return
}

func (s *SourceFile) String() string {
    if s == nil {
        return "SourceFile{nil}"
    }

    printableSrc := "<non-string source>"
    if p, ok := (s.src).(string); ok {
        printableSrc = p
    }

    return fmt.Sprintf("Sourcefile{source: %s, destination: %s, isPathDefault: %t", printableSrc, s.path, s.defaultPath)
}

func (src *SourceFile) FileName() (name string) {
    if src != nil {
        name = src.dsName() + ".go"
    }

    return
}

func DirectivesFrom(raw string) (d *Directives) {
    if raw != "" {
        datastructureType, element := ParseDirective(raw)
        d = &Directives{
            dsType:     datastructureType,
            entityType: element,
            raw:        raw,
        }
    }

    return
}

func (d *Directives) CreateSourceFileFromPath(path string) (src *SourceFile) {
    if d != nil {
        src = &SourceFile{
            pkg:         nil,
            raw:         d.raw,
            entityType:  d.entityType,
            dsType:      d.dsType,
            mode:        d.mode,
            defaultPath: false,
            path:        path,
            src:         `package main`,
        }
    }

    return
}

func (d *Directives) CreateSourceFileFromRawSource(raw string) (src *SourceFile) {
    if d != nil {
        src = &SourceFile{
            pkg:         nil,
            raw:         d.raw,
            entityType:  d.entityType,
            dsType:      d.dsType,
            mode:        d.mode,
            defaultPath: true,
            path:        fmt.Sprintf("tmp-%s.go", d.dsType),
            src:         raw,
        }
    }

    return
}

type LogMode int

const (
    Noisy LogMode = iota
    Silent
)

func (d *Directives) Mode(m LogMode) {
    if d == nil {
        return
    }
    d.mode = m
}

func (src *SourceFile) Visit(n ast.Node) ast.Visitor {
    if src == nil {
        return nullVisitor{}
    }
    switch v := n.(type) {
    case *ast.TypeSpec:
        src.debugln("Inspecting *ast.TypeSpec...")
        if t, ok := v.Type.(*ast.ArrayType); ok && v.Name.Name == "StackTemplate" {
            if _, ok := t.Elt.(*ast.InterfaceType); ok {
                t.Elt = ast.NewIdent(src.entityType)
            }

        }
        if v.Name.Name == "StackTemplate" {
            src.debugln("Replacing identifier fileName...")
            v.Name.Name = src.dsName()
            if v.Name.Obj != nil && v.Name.Obj.Kind == ast.Typ && v.Name.Obj.Name == "StackTemplate" {
                src.debugln("Replacing type fileName...")
                v.Name.Obj.Name = src.dsName()
            }
        }
    case *ast.Field:
        src.debugln("Inspecting *ast.Field...")
        if _, ok := v.Type.(*ast.InterfaceType); ok {
            v.Type = ast.NewIdent(src.entityType)
        }
        if s, ok := v.Type.(*ast.StarExpr); ok {
            if ident, ok := s.X.(*ast.Ident); ok {
                if ident.Name != src.entityType {
                    s.X = ast.NewIdent(src.dsName())
                }
            }
        }
    case *ast.StarExpr:
        src.debugln("Inspecting *ast.StarExpr...")
        if _, ok := v.X.(*ast.InterfaceType); ok {
            v.X = ast.NewIdent(src.entityType)
        }

    case *ast.Package:
        src.debugln("Inspecting *ast.Package...")
        v.Name = src.pkgName
    default:
        // log.Println("found other node; moving on...")
        return src
    }

    return src
}

func (src *SourceFile) debugln(msg string) {
    if src != nil && src.mode == Noisy {
        log.Println(msg)
    }
}

func (src *SourceFile) dsName() (name string) {
    if src != nil {
        ds := fmt.Sprintf("%s", src.dsType)
        name = fmt.Sprintf("%s%s", src.entityType, strings.ToUpper(ds[:1])+ds[1:])
    }

    return
}

func (src *SourceFile) Raw() string {
    return src.raw
}

func ParseDirective(instructions string) (DatastructureType, string) {
    matches := mask.FindAllStringSubmatch(instructions, -1)

    if matches == nil || len(matches) < 1 {
        return Unknown, ""
    }

    subs := matches[0]
    if subs[0] == instructions {
        subs = subs[1:]
    }

    for i, dsType := range datastructures {
        if strings.ToLower(dsType) == strings.ToLower(subs[0]) {
            return DatastructureType(i), subs[1]
        }
    }

    return Unknown, ""
}

type nullVisitor struct{}

func (nullVisitor) Visit(n ast.Node) ast.Visitor {
    return nil
}
