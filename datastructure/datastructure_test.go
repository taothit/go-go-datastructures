package datastructure

import (
	"bytes"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCopyTo(t *testing.T) {
	testCases := []struct {
		desc      string
		dsType    Type
		buf       bytes.Buffer
		want      []byte
		templates Template
		wantErr   error
	}{
		{"unknown", Unknown, bytes.Buffer{}, nil, [][]byte{{}, []byte("Complete"), {}}, UnknownType},
		{"stack", Stack, bytes.Buffer{}, []byte("Complete"), [][]byte{{}, []byte("Complete"), {}}, nil},
		{"heap", Heap, bytes.Buffer{}, []byte("Complete"), [][]byte{{}, {}, []byte("Complete")}, nil},
		{"uninitialized templates", Stack, bytes.Buffer{}, nil, nil, TemplatesUninitialized},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			templates = &tC.templates
			if tC.templates == nil {
				templates = nil
			}

			if err := tC.dsType.CopyTo(&tC.buf); !cmp.Equal(tC.want, tC.buf.Bytes()) || tC.wantErr != nil && !errors.Is(err, tC.wantErr) {
				t.Errorf("CopyTo(%v)=%s; got=%s; error=%v", tC.buf, tC.want, tC.buf.Bytes(), err)
			}
		})
	}
}
