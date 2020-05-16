package datastructure

import (
	"fmt"
	"go/ast"
	"log"
	"regexp"
	"strings"
)

var directiveMask = regexp.MustCompile(`^([a-zA-Z]\w+)\[([a-zA-Z]\w+)\]$`)

type Instruction struct {
	dsType     Type
	entityType string
}

func parse(instructions string) *Instruction {
	dsType, entityType := ParseInstructions(instructions)

	return &Instruction{
		dsType:     dsType,
		entityType: entityType,
	}
}

func (i *Instruction) Visit(n ast.Node) ast.Visitor {
	switch v := n.(type) {
	case *ast.TypeSpec:
		log.Println("Inspecting *ast.TypeSpec...")
		if t, ok := v.Type.(*ast.ArrayType); ok && v.Name.Name == "StackTemplate" {
			if _, ok := t.Elt.(*ast.InterfaceType); ok {
				t.Elt = ast.NewIdent("Widget")
			}

		}
		if v.Name.Name == "StackTemplate" {
			log.Println("Replacing identifier name...")
			v.Name.Name = i.dsName()
			if v.Name.Obj != nil && v.Name.Obj.Kind == ast.Typ && v.Name.Obj.Name == "StackTemplate" {
				log.Println("Replacing type name...")
				v.Name.Obj.Name = i.dsName()
			}
		}
	case *ast.Field:
		log.Println("Inspecting *ast.Field...")
		if _, ok := v.Type.(*ast.InterfaceType); ok {
			v.Type = ast.NewIdent("Widget")
		}
	case *ast.StarExpr:
		log.Println("Inspecting *ast.StarExpr...")
		if _, ok := v.X.(*ast.InterfaceType); ok {
			v.X = ast.NewIdent("Widget")
		}
	default:
		// log.Println("found other node; moving on...")
		return i
	}

	return i
}

func (i *Instruction) dsName() string {
	if i == nil {
		return ""
	}
	ds := fmt.Sprintf("%s", i.dsType)
	return fmt.Sprintf("%s%s", i.entityType, strings.ToUpper(ds[:1])+ds[1:])
}

func ParseInstructions(instructions string) (Type, string) {
	matches := directiveMask.FindAllStringSubmatch(instructions, -1)

	if matches == nil || len(matches) < 1 {
		return Unknown, ""
	}

	subs := matches[0]
	if subs[0] == instructions {
		subs = subs[1:]
	}

	for i, dsType := range datastructures {
		if strings.ToLower(dsType) == strings.ToLower(subs[0]) {
			return Type(i), subs[1]
		}
	}

	return Unknown, ""
}
