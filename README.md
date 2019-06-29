## Go Generative Datastructures

Use Go's code-generation feature, `go:generate`, to build templatized datastructures to suit the purpose of an application.

Each generated datastructure will map a single type to the container(s) for each file.  

### Usage
Define a specific datastructure for your type, Foo, in the go:generate line of your file.

```go
FooStack.go
-------------
//go:generate ggd stack[Foo]:FooStack
package foo_stack

```
Then, just run `go generate` once, and GGD will generate the implementation in situ.

```go
fooStack.go
-------------
//go:generate ggd stack[Foo]
package foo_stack

import (
  "ggd/example/foo
)


type FooStack interface() {
  Pop() *foo.Foo
  Push(*foo.Foo) error
  Peek() foo.Foo
  Offer(*foo.Foo) error 
  Length() int
}

// ======= Generated Implementation =======

func (f *FooStack) Pop() *foo.Foo {
  // ...
}

func (f *FooStack) Push(foo *foo.Foo) error {
  // ...
}

func (f *FooStack) Peek() *foo.Foo {
  // ...
}

func (f *FooStack) Offer() error {
  // ...
}

func (f *FooStack) Length() {
  // ...
}

// ======= Generated Implementation =======
```
### MVP punchlist
* [ ] Provide the following basic datastructures
  * [ ] List
  * [ ] Stack
  * [ ] Queue
  * [ ] Linked List
  * [ ] Double-linked list
  * [ ] BST
  * [ ] Set
* [ ] Generate directive allows customization of
  * [ ] Structure name
  * [ ] Package name
  * [ ] out
* [ ] Generate nested datastructures (stack of stacks, e.g.) 

### v2.0 
* [ ] Allow definition through nesting types:
```go
fooDeque.go
----
//go:generate ggd ...[Foo]:FooDeque:foo_deque
package foo_deque

import (
  "gdd/example/foo"
  "gdd/gdd"
)

type FooDeque interface {
  *gdd.Stack // Foo
  *gdd.Queue // Foo
}

```
Running `go generate` would create a non-overlapping union of functionality from the two embedded datastructures
```go
// ======= Generated Implementation =======

func (f *FooDeque) Pop() *foo.Foo {
  // ...
}

func (f *FooDeque) Push(foo *foo.Foo) error {
  // ...
}

func (f *FooDeque) Peek() *foo.Foo {
  // ...
}

func (f *FooDeque) Offer(foo *foo.Foo) error {
  // ...
}

func (f *FooDeque) Length() {
  // ...
}

func (f *FooDeque) Enqueue(foo *foo.Foo) error {
  // ...
}

func (f *FooDeque) Dequeue() *foo.Foo {
  // ...
}

func (f *FooDeque) PeekHead() *foo.Foo {
  // ...
}

func (f *FooDeque) PeekTail() *foo.Fee {
  // ...
}
```

* [ ] `:nil` vs. `:error`
  * Accessors will return `nil` entities when unable to be found.
  * Mutators will return `bool` when trying to add an entity to the container (`true` if it is stored [or present already]; `false` if it cannot be stored).
  * Users can override behavior to return error messages.  Accessors will return two values (*entity, error).  Mutators will return only an error with `nil` error signifying success.
  * Users will instruct using `:` directive after entity type declaration in the generate clause.  Adding `:nil` is syntactically correct, while redundant, and can be used for readability preferences.