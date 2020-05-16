package datastructure_test

import (
	"taothit/ggd/datastructure"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseInstructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction string
		dsType      datastructure.Type
		entity      string
	}{
		{"happy path", "Stack[Widget]", datastructure.Stack, "Widget"},
		{"no datastructure", "[Widget]", datastructure.Unknown, ""},
		{"Unknown datastructure", "foo[Widget]", datastructure.Unknown, ""},
		{"no entity", "foo[]", datastructure.Unknown, ""},
		{"malformed datastructure", "2[Widget]", datastructure.Unknown, ""},
		{"malformed entity", "Stack[2]", datastructure.Unknown, ""},
		{"entity with number in body", "Stack[W2]", datastructure.Stack, "W2"},
		{"entity with underscore", "Stack[W_2]", datastructure.Stack, "W_2"},
		{"Heap of widgets", "Heap[Widget]", datastructure.Heap, "Widget"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if gotType, gotEntity := datastructure.ParseInstructions(test.instruction); !cmp.Equal(test.dsType, gotType) || !cmp.Equal(test.entity, gotEntity) {
				t.Errorf("parse(%s)=(Type=%s,entity=%s);got=(Type=%s,entity=%s)", test.instruction, test.dsType, orEmpty(test.entity), gotType, orEmpty(gotEntity))
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
