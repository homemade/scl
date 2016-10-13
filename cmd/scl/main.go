package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/aryann/difflib"
	"github.com/tucnak/climax"

	"bitbucket.org/homemade/scl"
)

func main() {
	app := climax.New("scl")
	app.Brief = "Scl is a tool for managing SCL soure code."
	app.Version = "1.1.0"

	app.AddCommand(runCommand(os.Stdout, os.Stderr))
	app.AddCommand(testCommand(os.Stdout, os.Stderr))

	os.Exit(app.Run())
}

func runCommand(stdout io.Writer, stderr io.Writer) climax.Command {

	return climax.Command{
		Name:  "run",
		Brief: "Transform one or more .scl files into HCL",
		Usage: `[options] <filename.scl> [filenameX.scl...]`,
		Help:  `Transform one or more .scl files into HCL. Output is written to stdout.`,

		Flags: standardParserParams(),

		Handle: func(ctx climax.Context) int {

			if len(ctx.Args) == 0 {
				fmt.Fprintf(stderr, "At least one filename is required. See `sep help run` for syntax")
				return 1
			}

			params, includePaths := parserParams(ctx)

			for _, fileName := range ctx.Args {

				parser, err := scl.NewParser(scl.NewDiskSystem())

				if err != nil {
					fmt.Fprintf(stderr, "Error: Unable to create new parser in CWD: %s\n", err.Error())
					return 1
				}

				for _, includeDir := range includePaths {
					parser.AddIncludePath(includeDir)
				}

				for _, p := range params {
					parser.SetParam(p.name, p.value)
				}

				if err := parser.Parse(fileName); err != nil {
					fmt.Fprintf(stderr, "Error: Unable to parse file: %s\n", err.Error())
					return 1
				}

				fmt.Fprintf(stdout, "/* %s */\n%s\n\n", fileName, parser)
			}

			return 0
		},
	}
}

func testCommand(stdout io.Writer, stderr io.Writer) climax.Command {

	return climax.Command{
		Name:  "test",
		Brief: "Parse each .scl file in a directory and compare the output to an .hcl file",
		Usage: `[options] [file-glob...]`,
		Help:  "Parse each .scl file in a directory and compare the output to an .hcl file",

		Flags: standardParserParams(),

		Handle: func(ctx climax.Context) int {

			errors := 0

			reportError := func(path string, err string, args ...interface{}) {
				fmt.Fprintf(stderr, "[%s] %s\n", path, fmt.Sprintf(err, args...))
				errors++
			}

			if len(ctx.Args) == 0 {
				fmt.Fprintf(stderr, "At least one file glob is required. See `sep help test` for syntax")
				return 1
			}

			newlineMatcher := regexp.MustCompile("\n\n")
			params, includePaths := parserParams(ctx)

			for _, fileName := range ctx.Args {

				fs := scl.NewDiskSystem()
				parser, err := scl.NewParser(fs)

				if err != nil {
					reportError("Unable to create new parser in CWD: %s", err.Error())
					continue
				}

				for _, includeDir := range includePaths {
					parser.AddIncludePath(includeDir)
				}

				for _, p := range params {
					parser.SetParam(p.name, p.value)
				}

				if err := parser.Parse(fileName); err != nil {
					reportError(fileName, "Unable to parse file: %s", err.Error())
					continue
				}

				hclFilePath := strings.TrimSuffix(fileName, ".scl") + ".hcl"
				hclFile, _, err := fs.ReadCloser(hclFilePath)

				if err != nil {
					fmt.Fprintf(stdout, "[%s] No HCL file; skipping gold standard test\n", fileName)
					continue
				}

				hcl, err := ioutil.ReadAll(hclFile)

				if err != nil {
					reportError(fileName, "Unable to read .hcl file: %s", err.Error())
					continue
				}

				hclLines := strings.Split(strings.TrimSuffix(newlineMatcher.ReplaceAllString(string(hcl), "\n"), "\n"), "\n")
				sclLines := strings.Split(parser.String(), "\n")

				diff := difflib.Diff(hclLines, sclLines)

				success := true

				for _, d := range diff {
					if d.Delta != difflib.Common {
						success = false
					}
				}

				if !success {
					reportError(fileName, "Diff failed:")

					fmt.Fprintln(stderr)

					for _, d := range diff {
						fmt.Fprintf(stderr, "\t%s\n", d.String())
					}

					fmt.Fprintln(stderr)

					continue
				}

				fmt.Fprintf(stdout, "[%s] OK\n", fileName)
			}

			if errors > 0 {
				fmt.Fprintf(stderr, "\n[FAIL] %d error(s)\n", errors)
				return 1
			}

			return 0
		},
	}
}

func standardParserParams() []climax.Flag {

	return []climax.Flag{
		{
			Name:     "include",
			Short:    "i",
			Usage:    `--include /path/to/lib1,/path/to/lib2`,
			Help:     `Comma-separated list of include paths`,
			Variable: true,
		},
		{
			Name:     "param",
			Short:    "p",
			Usage:    `--param param0=somthing,"param1='something else'"`,
			Help:     `Comma-separated list of include paths`,
			Variable: true,
		},
		{
			Name:  "no-env",
			Short: "ne",
			Usage: `--no-env`,
			Help:  `Don't import envionment variables when parsing the SCL`,
		},
	}

}

func parserParams(ctx climax.Context) (params paramSlice, includePaths []string) {

	if !ctx.Is("no-env") {
		for _, envVar := range os.Environ() {
			params.Set(envVar)
		}
	}

	if ps, set := ctx.Get("param"); set {
		for _, p := range strings.Split(ps, ",") {
			params.Set(p)
		}
	}

	if ps, set := ctx.Get("include"); set {
		for _, i := range strings.Split(ps, ",") {
			includePaths = append(includePaths, i)
		}
	}

	return
}
