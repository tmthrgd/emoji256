// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package emoji256

import (
	"io"
	"testing"
	"testing/quick"
)

const (
	testIn  = "my string\n"
	testOut = "\U0001f377\u231b\U0001f3b8\U0001f384\U0001f385\U0001f383\U0001f345\U0001f37a\U0001f350\U0001f448"
)

func TestEncode(t *testing.T) {
	out, err := EncodeToString([]byte(testIn))
	if err != nil {
		t.Error(err)
	}

	if out != testOut {
		t.Error("Encode failed")
	}
}

func TestDecode(t *testing.T) {
	out, err := DecodeString(testOut)
	if err != nil {
		t.Error(err)
	}

	if string(out) != testIn {
		t.Error("Decode failed")
	}
}

func TestRoundTrip(t *testing.T) {
	if err := quick.CheckEqual(func(in []byte) (out []byte, err error) {
		return in, nil
	}, func(in []byte) (out []byte, err error) {
		if len(in) == 0 {
			return in, nil
		}

		buf, err := EncodeBytes(in)
		if err != nil {
			return nil, err
		}

		return DecodeBytes(buf)
	}, &quick.Config{
		MaxCountScale: 250,
	}); err != nil {
		t.Error(err)
	}
}

var benchSizes = []struct {
	name string
	l    int64
}{
	{"16", 16},
	{"32", 32},
	{"128", 128},
	{"1K", 1 * 1024},
	{"16K", 16 * 1024},
	{"128K", 128 * 1024},
	{"1M", 1024 * 1024},
	{"16M", 16 * 1024 * 1024},
	{"128M", 128 * 1024 * 1024},
}

type nilReader struct{}

func (nilReader) Read(p []byte) (n int, err error) {
	return len(p), nil
}

type nilWriter struct{}

func (nilWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

type bytesReader []byte

func (b bytesReader) Read(p []byte) (n int, err error) {
	if copy(p, b) != len(b) {
		return 0, io.ErrShortBuffer
	}

	for i := len(b); i < len(p); i *= 2 {
		copy(p[i:], p[:i])
	}

	return len(p) - (len(p) % len(b)), nil
}

func BenchmarkEncode(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			b.SetBytes(size.l)

			r := &io.LimitedReader{R: nilReader{}}
			w := nilWriter{}

			for i := 0; i < b.N; i++ {
				r.N = size.l

				if err := Encode(w, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkDecode(b *testing.B) {
	for _, size := range benchSizes {
		b.Run(size.name, func(b *testing.B) {
			b.SetBytes(size.l)

			r := &io.LimitedReader{R: bytesReader{0xf0, 0x9f, 0x91, 0x8d}}
			w := nilWriter{}

			for i := 0; i < b.N; i++ {
				r.N = size.l

				if err := Decode(w, r); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
