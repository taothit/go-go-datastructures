//go:generate ggd -pathTo=templates/example/widgetStack.go stack[Widget]
package example

import (
	"fmt"
	"strings"
)

type Stack interface {
	Push(i Widget) bool
	Pop() *Widget
	Length() int
	fmt.Stringer
	fmt.GoStringer
}
type WidgetStack []Widget

func (s *WidgetStack) Push(i Widget) bool {
	if s == nil {
		return false
	}
	current := len(*s)
	expanded := append(*s, i)
	*s = expanded
	return current < len(*s)
}
func (s *WidgetStack) Pop() *Widget {
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
func (s *WidgetStack) Length() int {
	return len(*s)
}
func (s *WidgetStack) String() string {
	var out string
	for _, el := range *s {
		out = fmt.Sprintf("%s %s", out, el)
	}
	return fmt.Sprintf("[%s]<", strings.Trim(out, " "))
}
func (s *WidgetStack) GoString() string {
	return s.String()
}
