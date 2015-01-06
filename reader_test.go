package blocksreader

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
)

func TestReader(t *testing.T) {
	tests := []struct {
		name string
		in   []byte
		atns []Block
		want []byte
		err  error
	}{
		{
			name: "Block with negative offset",
			in:   []byte{},
			atns: []Block{Block{-1, 1}},
			want: nil,
			err:  errors.New("blockReader: invalid block - offset < 0"),
		},
		{
			name: "Block with zero length",
			in:   []byte{},
			atns: []Block{Block{0, 0}},
			want: nil,
			err:  errors.New("blockReader: invalid block - length < 1"),
		},
		{
			name: "Single block",
			in:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			atns: []Block{Block{3, 5}},
			want: []byte{3, 4, 5, 6, 7},
		},
		{
			name: "Multiple blocks 1",
			in:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			atns: []Block{Block{2, 1}, Block{7, 1}},
			want: []byte{2, 7},
		},
		{
			name: "Multiple blocks 2",
			in:   []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
			atns: []Block{Block{7, 5}, Block{0, 2}, Block{13, 2}, Block{9, 6}},
			want: []byte{7, 8, 9, 10, 11, 0, 1, 13, 14, 9, 10, 11, 12, 13, 14},
		},
	}
	for _, test := range tests {
		// Test with different buffer lengths
		for n := 1; n < 20; n++ {
			inReader := bytes.NewReader(test.in)
			r := NewReader(inReader, test.atns)

			var got []byte
			for {
				p := make([]byte, n)
				if n, err := r.Read(p); err == io.EOF {
					break
				} else if err != nil {
					if test.err == nil || test.err.Error() != err.Error() {
						t.Errorf("Test %s (%d): Read() returned error. Got %v, want %v", test.name, n, err, test.err)
					}
					break
				} else {
					got = append(got, p[:n]...)
				}
			}
			if test.want != nil && !reflect.DeepEqual(got, test.want) {
				t.Errorf("Test %s (%d): returned bytes were not correct. Got %v, want %v", test.name, n, got, test.want)
			}
		}
	}
}
