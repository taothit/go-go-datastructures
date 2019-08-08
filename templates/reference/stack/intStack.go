package stack

import (
	"fmt"
	"strconv"
	"strings"
)

// IntStack stores integers in a variable-size array and allows LIFO access.
type IntStack []int

// Push adds i to the top of the IntStack and returns its new size.
func (s *IntStack) Push(i int) bool {
	if s == nil {
		return false
	}

	current := len(*s)
	expanded := append(*s, i)
	*s = expanded
	return current < len(*s)
}

// Pop removes the top of the IntStack and returns it
func (s *IntStack) Pop() *int {
	if s == nil {
		return nil
	}

	if last := len(*s) - 1; last < 0 {
		return nil
	}  else {
		item := (*s)[len(*s)-1]
		reduced := (*s)[:last]
		*s = reduced
		return &item
	}
}

func (s *IntStack) Length() int {
	return len(*s)
}

// Print returns the string representation of the IntStack
func (s *IntStack) String() string {
	var out string
	for _, el := range *s {
		out = fmt.Sprintf("%s %s", out, strconv.Itoa(el))
	}
	return fmt.Sprintf("[%s]<", strings.Trim(out, " "))
}

// GoString provides Type representation when printing an IntStack with '%v'.
func (s *IntStack) GoString() string {
	return s.String()
}
