// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package main

//go:generate go run ./table-gen.go -o table.go -p main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

func must(err error) {
	if err == nil {
		return
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func main() {
	var decode bool
	flag.BoolVar(&decode, "d", false, "decode")

	flag.Parse()

	in := bufio.NewReader(os.Stdin)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for decode {
		r, _, err := in.ReadRune()
		if err == io.EOF {
			return
		}
		must(err)

		if b, ok := decTable[r]; ok {
			must(out.WriteByte(b))
		} else if !unicode.IsSpace(r) {
			must(fmt.Errorf("invalid character: %#U", r))
		}
	}

	for {
		b, err := in.ReadByte()
		if err == io.EOF {
			out.WriteRune('\n')
			return
		}
		must(err)

		_, err = out.Write(encTable[b])
		must(err)
	}
}
