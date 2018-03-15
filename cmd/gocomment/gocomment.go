package main

import (
	"personal/commgen/pkg/parser"
	"flag"
	"fmt"
	"os"
)

func main() {
	// file flags
	var flagFile = flag.String("file", "", "Filename to be parsed")
	flag.Parse()
	if err := parser.Work(*flagFile); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
