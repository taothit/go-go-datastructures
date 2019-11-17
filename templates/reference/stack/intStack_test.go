package stack_test

import (
	"fmt"
	"github.com/google/go-cmp/cmp"
	"taothit/ggd/templates/reference/stack"
	"testing"
)

func TestIntStack_Push(t *testing.T) {
	tests := []struct {
		name    string
		prepped func(i ...*int) *stack.intStack
		i       int
		want    bool
	}{
		{"add to empty stack", func(i ...*int) *stack.intStack {
			return &stack.intStack{}
		}, 0, true},
		{"pushing onto nil stack", func(i ...*int) *stack.intStack {
			return nil
		}, 0, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.prepped()
			if got := s.Push(test.i); got != test.want {
				t.Errorf("intStack.Push(%d)=%t; got=%t", test.i, test.want, got)
			}
		})
	}
}

func TestIntStack_Pop(t *testing.T) {
	tests := []struct {
		name    string
		prepped func(...*int) *stack.intStack
		want    *int
	}{
		{
			"Popping nil stack", func(wants ...*int) *stack.intStack {
			return nil
		}, nil},
		{
			"Popping empty stack", func(wants ...*int) *stack.intStack {
			return &stack.intStack{}
		}, nil},
		{
			"Popping top item", func(wants ...*int) *stack.intStack {
			s := stack.intStack{}
			for _, want := range wants {
				s.Push(*want)
			}
			return &s
		}, intPtr(1)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.prepped(test.want).Pop(); !cmp.Equal(got, test.want) {
				t.Errorf("intStack.Pop()=%v; got=%v", test.want, got)
			}
		})
	}
}

func intPtr(i int) *int {
	ptr := i
	return &ptr
}

func TestIntStack_PushThenPop_ReversesOrder(t *testing.T) {
	s := &stack.intStack{}
	want := []int{3, 2, 1}
	for i := len(want) - 1; i >= 0; i-- {
		s.Push(want[i])
	}

	got := make([]int, 0)
	for i := 0; i < len(want); i++ {
		got = append(got, *(s.Pop()))
	}

	if !cmp.Equal(got, want) || s.Length() != 0 {
		t.Errorf("intStack.Push(%v)->intStack.Pop()x%d=%v; got=%v; stack=%v", []int{1, 2, 3}, len(want), want, got, s)
	}
}

func TestIntStack_String(t *testing.T) {
	s := stackWith(1, 2, 3)

	want := "[1 2 3]<"

	if got := fmt.Sprintf("%s", s); !cmp.Equal(got, want) {
		t.Errorf("intStack.String()=%s; got=%s", want, got)
	}
}

func stackWith(ints ...int) *stack.intStack {
	s := &stack.intStack{}
	for _, i := range ints {
		s.Push(i)
	}
	return s
}

func TestIntStack_GoString(t *testing.T) {
	s := stackWith(1, 2, 3)

	want := "[1 2 3]<"

	if got := fmt.Sprintf("%v", s); !cmp.Equal(got, want) {
		t.Errorf("intStack.String()=%s; got=%s", want, got)
	}
}
