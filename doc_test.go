package scl

import (
	"fmt"
	"log"

	"github.com/hashicorp/hcl"
)

func Example_basic() {

	parser, err := NewParser(NewDiskSystem())

	if err != nil {
		log.Fatal(err)
	}

	if err := parser.Parse("myfile.scl"); err != nil {
		// This is a language error, which will include
		// the filename and line of the error, as well
		// as an explanatory message.
		log.Fatal(err)
	}

	fmt.Println("myfile.scl as HCL:", parser.String())
}

func Example_hcl() {

	parser, err := NewParser(NewDiskSystem())

	if err != nil {
		log.Fatal(err)
	}

	if err := parser.Parse("myfile.scl"); err != nil {
		log.Fatal(err)
	}

	myConfig := struct {
		SomeThing string `hcl:"some-thing"`
	}{}

	if err := hcl.Decode(&myConfig, parser.String()); err != nil {
		log.Fatal(err)
	}

	fmt.Println(myConfig)
}

func ExampleParser_variables() {

	parser, err := NewParser(NewDiskSystem())

	if err != nil {
		log.Fatal(err)
	}

	parser.SetParam("my-variable", "my value")

	if err := parser.Parse("myfile.scl"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("myfile.scl as HCL:", parser.String())
}

func ExampleParser_includePaths() {

	parser, err := NewParser(NewDiskSystem())

	if err != nil {
		log.Fatal(err)
	}

	parser.AddIncludePath("path/to/library")

	if err := parser.Parse("myfile.scl"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("myfile.scl as HCL:", parser.String())
}

func ExampleParser_documentation() {

	parser, err := NewParser(NewDiskSystem())

	if err != nil {
		log.Fatal(err)
	}

	documentation, err := parser.Documentation("myfile.scl")

	if err != nil {
		log.Fatal(err)
	}

	for i, mixin := range documentation {
		fmt.Printf("Mixin %d: %+v", i, mixin)
	}
}
