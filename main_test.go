package main

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		instruction string
		dsType      datastructureType
		entity      string
	}{
		{"happy path", "stack[Widget]", stack, "Widget"},
		{"no datastructure", "[Widget]", unknown, ""},
		{"unknown datastructure", "foo[Widget]", unknown, ""},
		{"no entity", "foo[]", unknown, ""},
		{"malformed datastructure", "2[Widget]", unknown, ""},
		{"malformed entity", "stack[2]", unknown, ""},
		{"entity with number in body", "stack[W2]", stack, "W2"},
		{"entity with underscore", "stack[W_2]", stack, "W_2"},
		{"heap of widget", "heap[Widget]", heap, "Widget"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if gotType, gotEntity := parseInstructions(test.instruction); !cmp.Equal(test.dsType, gotType) || !cmp.Equal(test.entity, gotEntity) {
				t.Errorf("parse(%s)=(datastructureType=%s,entity=%s);got=(datastructureType=%s,entity=%s)", test.instruction, test.dsType, orEmpty(test.entity), gotType, orEmpty(gotEntity))
			}
		})
	}
}

func orEmpty(s string) string {
	if s == "" {
		s = "<empty>"
	}
	return s
}
