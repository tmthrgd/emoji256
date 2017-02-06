// Copyright 2017 Tom Thorogood. All rights reserved.
// Use of this source code is governed by a
// Modified BSD License license that can be found in
// the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/tmthrgd/emoji256"
)

func main() {
	var decode bool
	flag.BoolVar(&decode, "d", false, "decode")

	flag.Parse()

	var err error
	if decode {
		err = emoji256.Decode(os.Stdout, os.Stdin)
	} else {
		err = emoji256.Encode(os.Stdout, os.Stdin)
		os.Stdout.WriteString("\n")
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
