package main

import (
	"fmt"
	"io"
	"os"

	"bitbucket.org/homemade/scl"
)

func main() {
	os.Exit(entryPoint(os.Args, os.Stdout, os.Stderr))
}

func entryPoint(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 1 {
		fmt.Fprintf(stdout, "Usage: %s <filename.scl> [filenameX.scl...]\n", args[0])
		return 1
	}

	for _, fileName := range args[1:] {
		parser, err := scl.NewParser(scl.NewDiskSystem())
		if err != nil {
			fmt.Fprintf(stderr, "Error: Unable to create new parser in CWD: %s\n", err.Error())
			return 1
		}

		if err := parser.Parse(fileName); err != nil {
			fmt.Fprintf(stderr, "Error: Unable to parse file: %s\n", err.Error())
			return 1
		}

		fmt.Fprintf(stdout, "/* %s */\n%s\n", fileName, parser)
	}

	return 0
}
