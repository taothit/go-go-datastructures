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

var stackSource = `package cmd_test`

var mode model.LogMode = model.Silent

func TestGenerate(t *testing.T) {
	tests := []struct {
		name         string
		refPath      string
		instructions string
		source       *cmd.Directives
	}{
		{"verify stack", "../templates/reference/stack/intStack.go", "Stack[int]", cmd.RawSource(stackSource)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p, err := filepath.Abs(test.refPath)
			if err != nil {
				t.Errorf("cmd.Generate(%s, %v, io.Writer); err=%v", test.instructions, test.source, err)
			}

			f, err := os.Open(p)
			if err != nil {
				t.Errorf("cmd.Generate(%s, %v, io.Writer); err=%v", test.instructions, test.source, err)
			}

			info, err := f.Stat()
			if err != nil {
				t.Fatalf("cmd.Generate(%s, %v, io.Writer); err=%v", test.instructions, test.source, err)
			}
			if mode == model.Noisy {
				log.Printf("Reference file: Name=%s, FileMode=%s, ", info.Name(), info.Mode())
			}

			want := ByteSliceOf(info.Size())
			n, err := f.Read(want)
			if err != nil && err != io.EOF {
				t.Errorf("cmd.Generate(%s, %v, io.Writer); err=%v", test.instructions, test.source, err)
			}
			if n < 1 {
				t.Fatalf("cmd.Generate(%s, %v, io.Writer); read %d of %d from %s", test.instructions, test.source, n, info.Size(), test.refPath)
			}

			got := bytes.NewBuffer(ByteSliceOf(0, 2048))
			w := io.Writer(got)
			if err := cmd.Generate(test.instructions, test.source, &w, mode); err != nil {
				t.Fatalf("cmd.Generate(%s, %v, io.Writer); err=%v", test.instructions, test.source, err)
			}

			if !cmp.Equal(want, got) {
				if mode == model.Noisy {
					log.Printf("Want:\n%s", string(want))
					log.Printf("Got:\n%s", string(got.Bytes()))
				}

				t.Errorf("cmd.Generate(%s, %v, io.Writer); output different from reference (ref=%d bytes; output= %d bytes)", test.instructions, test.source, len(want), len(got.Bytes()))
			}
		})
	}
}

func ByteSliceOf(size ...int64) (b []byte) {
	switch len(size) {
	case 1:
		b = make([]byte, size[0])
	case 2:
		b = make([]byte, size[0], size[1])
	default:
		// nil
	}
	return
}
