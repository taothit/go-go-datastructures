package model_test

import (
	"github.com/google/go-cmp/cmp"
	"taothit/ggd/model"
	"testing"
)

func TestParseInstructions(t *testing.T) {
	tests := []struct {
		name        string
		instruction string
		dsType      model.DatastructureType
		entity      string
	}{
		{"happy path", "Stack[Widget]", model.Stack, "Widget"},
		{"no datastructure", "[Widget]", model.Unknown, ""},
		{"Unknown datastructure", "foo[Widget]", model.Unknown, ""},
		{"no entity", "foo[]", model.Unknown, ""},
		{"malformed datastructure", "2[Widget]", model.Unknown, ""},
		{"malformed entity", "Stack[2]", model.Unknown, ""},
		{"entity with number in body", "Stack[W2]", model.Stack, "W2"},
		{"entity with underscore", "Stack[W_2]", model.Stack, "W_2"},
		{"Heap of widgets", "Heap[Widget]", model.Heap, "Widget"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if gotType, gotEntity := model.ParseInstructions(test.instruction); !cmp.Equal(test.dsType, gotType) || !cmp.Equal(test.entity, gotEntity) {
				t.Errorf("Parse(%s)=(datastructureType=%s,entity=%s);got=(datastructureType=%s,entity=%s)", test.instruction, test.dsType, orEmpty(test.entity), gotType, orEmpty(gotEntity))
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
