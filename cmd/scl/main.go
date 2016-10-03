package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"bitbucket.org/homemade/scl"
)

func main() {
	flagOptions := &flagOptions{}

	flagSet := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <filename.scl> [filenameX.scl...]\n", os.Args[0])
		flagSet.PrintDefaults()
	}
	flagSet.Var(&(flagOptions.includePaths), "include", "Dir path to include on parser, aids scl import statements")
	flagSet.Var(&(flagOptions.params), "param", "Parameter to set on scl parser making it available in the global scope e.g. name=value")
	noEnv := flagSet.Bool("no-env", false, "Prevent ingesting all environmental variables into parser")

	flagSet.Parse(os.Args[1:])

	if flagSet.NArg() == 0 {
		flagSet.Usage()
		os.Exit(1)
	}

	if !*noEnv {
		for _, envVar := range os.Environ() {
			flagOptions.params.Set(envVar)
		}
	}

	os.Exit(entryPoint(flagOptions, flagSet.Args(), os.Stdout, os.Stderr))
}

type flagOptions struct {
	includePaths pathSlice
	params       paramSlice
}

func entryPoint(options *flagOptions, args []string, stdout io.Writer, stderr io.Writer) int {
	for _, fileName := range args {
		parser, err := scl.NewParser(scl.NewDiskSystem())
		if err != nil {
			fmt.Fprintf(stderr, "Error: Unable to create new parser in CWD: %s\n", err.Error())
			return 1
		}

		for _, includeDir := range options.includePaths {
			parser.AddIncludePath(includeDir)
		}

		for _, p := range options.params {
			parser.SetParam(p.name, p.value)
		}

		if err := parser.Parse(fileName); err != nil {
			fmt.Fprintf(stderr, "Error: Unable to parse file: %s\n", err.Error())
			return 1
		}

		fmt.Fprintf(stdout, "/* %s */\n%s\n\n", fileName, parser)
	}

	return 0
}

type pathSlice []string

func (s pathSlice) String() string {
	return strings.Join(s, ":")
}

func (s *pathSlice) Set(value string) error {
	*s = append(*s, value)

	return nil
}

type param struct {
	name  string
	value string
}

func (p param) String() string {
	return p.name + "=" + p.value
}

// expects 'name=value'
func (p *param) SetFromString(s string) error {
	parts := strings.SplitN(s, "=", 2)
	if len(parts) <= 1 {
		return errors.New("Unable to convert to param: " + s)
	}

	p.name = strings.TrimSpace(parts[0])
	p.value = strings.TrimSpace(parts[1])

	return nil
}

type paramSlice []*param

func (ps paramSlice) String() string {
	outputSlice := make([]string, len(ps))
	for i, p := range ps {
		outputSlice[i] = p.String()
	}

	return strings.Join(outputSlice, ", ")
}

func (ps *paramSlice) Set(value string) error {
	p := new(param)

	if err := p.SetFromString(value); err != nil {
		return err
	}

	*ps = append(*ps, p)

	return nil
}
