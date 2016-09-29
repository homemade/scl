[![GoDoc](https://godoc.org/bitbucket.org/homemade/scl?status.svg)](https://godoc.org/bitbucket.org/homemade/scl)

## Sepia Configuration Language

SCL is a simple, declarative, semi-functional, self-documenting language that extends
[HCL](https://github.com/hashicorp/hcl) in the same way that Sass extends CSS.
What that means is, any properly formatted HCL is valid SCL. If you really
enjoy HCL, you can keep using it exclusively: under the hood, SCL ‘compiles’ to
HCL. The difference is that now you can explicitly include files, use ‘mixins’
to quickly inject boilerplate code, and use properly scoped, natural variables.

The language is designed to accompany _Sepia_ (and, outside of that, _Sepia_
plugins) but it's a general purpose language, and can be used for pretty
much any configurational purpose.

There is full [documentation](http://sepia-docs.us-east-1.elasticbeanstalk.com/scl)
for SCL available at on the _Sepia_ documentation site.

## SCL CLI Tool

Takes `.scl` files in, spits HCL out (via Stdout)

### Installation

```
$ go get -u bitbucket.org/homemade/scl/...
```

### Usage

```
$ scl
Usage: scl [options] <filename.scl> [filenameX.scl...]
  -include value
    Dir path to include on parser, aids scl import statements
  -no-env
    Prevent ingesting all environmental variables into parser
  -param value
    Parameter to set on scl parser making it available in the global scope e.g. name=value
```

### Example

Basic example:
```
$ scl $GOPATH/src/bitbucket.org/homemade/scl/fixtures/valid/basic.scl
/* .../bitbucket.org/homemade/scl/fixtures/valid/basic.scl */
wrapper {
  inner = "yes"
  another = "1" {
    yet_another = "123"
  }
}
```

Adding includes:
```
$ scl -include $GOPATH/src/bitbucket.org/homemade/scl $GOPATH/src/bitbucket.org/homemade/scl/fixtures/valid/import.scl
wrapper {
  inner = "yes"
  another = "1" {
    yet_another = "123"
  }
}
output = "this is from simpleMixin"
```

Adding params via cli flags:
```
$ scl -param myVar=1 $GOPATH/src/bitbucket.org/homemade/scl/fixtures/valid/variables.scl
/* /Volumes/CaseSensitive/go/src/bitbucket.org/homemade/scl/fixtures/valid/variables.scl */
outer {
  inner = 1
}
```

Adding params via environmental variables:
```
$ myVar=1 scl $GOPATH/src/bitbucket.org/homemade/scl/fixtures/valid/variables.scl
/* /Volumes/CaseSensitive/go/src/bitbucket.org/homemade/scl/fixtures/valid/variables.scl */
outer {
  inner = 1
}
```
