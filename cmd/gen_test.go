package cmd_test

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"io"
	"log"
	"os"
	"path/filepath"
	"taothit/ggd/cmd"
	"taothit/ggd/model"
	"testing"
)

var stackSource = `package stack`

var mode = model.Silent

var fileName = "intStack"

func TestGenerate(t *testing.T) {
	tests := []struct {
		name       string
		refPath    string
		directives *model.Directives
	}{
		{"verify stack", "../templates/reference/stack/intStack.go", model.DirectivesFrom("Stack[int]")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := filepath.Abs(test.refPath)
			if err != nil {
				t.Errorf("cmd.Generate(%v, io.Writer); err=%v", test.directives, err)
			}

			f, err := os.Open(p)
			if err != nil {
				t.Errorf("cmd.Generate(%v, %v, io.Writer); err=%v", test.directives, test.directives, err)
			}

			info, err := f.Stat()
			if err != nil {
				t.Fatalf("cmd.Generate(%v, %v, io.Writer); err=%v", test.directives, test.directives, err)
			}
			if mode == model.Noisy {
				log.Printf("Reference file: Name=%s, FileMode=%s, ", info.Name(), info.Mode())
			}

			want := make([]byte, info.Size())
			n, err := f.Read(want)
			if err != nil && err != io.EOF {
				t.Errorf("cmd.Generate(%v, %v, io.Writer); err=%v", test.directives, test.directives, err)
			}
			if n < 1 {
				t.Fatalf("cmd.Generate(%v, %v, io.Writer); read %d of %d from %s", test.directives, test.directives, n, info.Size(), test.refPath)
			}

			got := new(bytes.Buffer)
			w := io.Writer(got)
			if err := cmd.Generate(test.directives, &w, mode, fileName, ""); err != nil {
				t.Fatalf("cmd.Generate(%v, %v, io.Writer); err=%v", test.directives, test.directives, err)
			}

			if diff := cmp.Diff(string(want), string(got.Bytes())); diff != "" {
				t.Errorf("cmd.Generate(%v, %v, io.Writer); output different from reference (ref=%d bytes; output= %d bytes):\n%s", test.directives, test.directives, len(want), len(got.Bytes()), diff)
			}

			os.Remove(fileName+".go")
		})
	}
}

