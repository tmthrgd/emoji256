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

func main() {
	var decode bool
	flag.BoolVar(&decode, "d", false, "decode")

	flag.Parse()

	in := bufio.NewReader(os.Stdin)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if decode {
		decTable := make(map[rune]byte, len(encTable))
		for b, r := range encTable {
			decTable[r] = byte(b)
		}

		for {
			r, _, err := in.ReadRune()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			if b, ok := decTable[r]; ok {
				out.WriteByte(b)
			} else if !unicode.IsSpace(r) {
				fmt.Fprintf(os.Stderr, "invalid character: %#U\n", r)
				os.Exit(1)
			}
		}
	} else {
		for {
			b, err := in.ReadByte()
			if err == io.EOF {
				break
			} else if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			r := encTable[b]
			out.WriteRune(r)
		}
	}

	if !decode {
		out.WriteRune('\n')
	}
}
