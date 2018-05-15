// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package emoji256

//go:generate go run ./table-gen.go -o table.go -p emoji256

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func Encode(w io.Writer, r io.Reader) error {
	in, out := bufio.NewReader(r), bufio.NewWriter(w)

	for {
		b, err := in.ReadByte()
		switch err {
		case nil:
		case io.EOF:
			return out.Flush()
		default:
			out.Flush()
			return err
		}

		if _, err = out.Write(encTable[b]); err != nil {
			out.Flush()
			return err
		}
	}
}

func EncodeBytes(in []byte) (out []byte, err error) {
	r := bytes.NewReader(in)

	var buf bytes.Buffer
	buf.Grow(len(in) * 4)

	if err = Encode(&buf, r); err != nil {
		return
	}

	return buf.Bytes(), nil
}

func EncodeToString(in []byte) (out string, err error) {
	r := bytes.NewReader(in)

	var buf strings.Builder
	buf.Grow(len(in) * 4)

	if err = Encode(&buf, r); err != nil {
		return
	}

	return buf.String(), nil
}

func Decode(w io.Writer, r io.Reader) error {
	in, out := bufio.NewReader(r), bufio.NewWriter(w)

	for {
		r, _, err := in.ReadRune()
		switch err {
		case nil:
		case io.EOF:
			return out.Flush()
		default:
			out.Flush()
			return err
		}

		if b, ok := decTable[r]; ok {
			if err = out.WriteByte(b); err != nil {
				out.Flush()
				return err
			}
		} else if !unicode.IsSpace(r) {
			out.Flush()
			return fmt.Errorf("invalid character: %#U", r)
		}
	}
}

func DecodeBytes(in []byte) (out []byte, err error) {
	r := bytes.NewReader(in)

	var buf bytes.Buffer
	buf.Grow(len(in) / 3)

	if err = Decode(&buf, r); err != nil {
		return
	}

	return buf.Bytes(), nil
}

func DecodeString(in string) (out []byte, err error) {
	r := strings.NewReader(in)

	var buf bytes.Buffer
	buf.Grow(len(in) / 3)

	if err = Decode(&buf, r); err != nil {
		return
	}

	return buf.Bytes(), nil
}
