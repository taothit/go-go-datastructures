package templates

import (
	"fmt"
	"strings"
)

// Stack stores entities in a variable-size array and allows LIFO access.
type Stack interface {
	Push(i interface{}) bool
	Pop() *interface{}
	Length() int
	fmt.Stringer
	fmt.GoStringer
}

// StackTemplate is a Stack for the interface{}
type StackTemplate []interface{}

// Push adds i to the top of the StackTemplate and returns its new size.
func (s *StackTemplate) Push(i interface{}) bool {
	if s == nil {
		return false
	}

	current := len(*s)
	expanded := append(*s, i)
	*s = expanded
	return current < len(*s)
}

// Pop removes the top of the StackTemplate and returns it
func (s *StackTemplate) Pop() *interface{} {
	if s == nil {
		return nil
	}

	if last := len(*s) - 1; last < 0 {
		return nil
	} else {
		item := (*s)[len(*s)-1]
		reduced := (*s)[:last]
		*s = reduced
		return &item
	}
}

// Length provides len() of internal storage.
func (s *StackTemplate) Length() int {
	return len(*s)
}

// String returns a human-readable representation of the StackTemplate
func (s *StackTemplate) String() string {
	var out string
	for _, el := range *s {
		out = fmt.Sprintf("%s %s", out, el)
	}
	return fmt.Sprintf("[%s]<", strings.Trim(out, " "))
}

// GoString provides Type representation when printing an StackTemplate with '%v'.
func (s *StackTemplate) GoString() string {
	return s.String()
}
