package stack

import (
    "fmt"
    "strconv"
    "strings"
)

// Stack stores entities in a variable-size array and allows LIFO access.
type Stack interface {
    Push(i int) bool
    Pop() *int
    Length() int
    fmt.Stringer
    fmt.GoStringer
}

// intStack stores integers in a variable-size array and allows LIFO access.
type intStack []int

// Push adds i to the top of the intStack and returns its new size.
func (s *intStack) Push(i int) bool {
    if s == nil {
        return false
    }

    current := len(*s)
    expanded := append(*s, i)
    *s = expanded
    return current < len(*s)
}

// Pop removes the top of the intStack and returns it
func (s *intStack) Pop() *int {
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
func (s *intStack) Length() int {
    return len(*s)
}

// Print returns the string representation of the intStack
func (s *intStack) String() string {
    var out string
    for _, el := range *s {
        out = fmt.Sprintf("%s %s", out, strconv.Itoa(el))
    }
    return fmt.Sprintf("[%s]<", strings.Trim(out, " "))
}

// GoString provides Type representation when printing an intStack with '%v'.
func (s *intStack) GoString() string {
    return s.String()
}
