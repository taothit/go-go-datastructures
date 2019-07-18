package example

import (
	"fmt"
)

type Widget struct {
	Id int
	Name string
}

func (w *Widget) Spin() {
	if w == nil {
		return
	}
	fmt.Printf("%s-%d is spinning!",w.Name, w.Id)
}
