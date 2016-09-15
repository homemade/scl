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