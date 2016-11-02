package scl

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/hcl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockParser(t *testing.T) *parser {

	p0, err := NewParser(NewDiskSystem())
	require.Nil(t, err)

	p := p0.(*parser)

	require.NotNil(t, p)
	require.NotNil(t, p.fs)

	return p
}

func Test_AParserCanBeCreated(t *testing.T) {
	newMockParser(t)
}

func Test_AParserCanParseFiles(t *testing.T) {

	for cycle, input := range []struct {
		fileName  string
		hcl       string
		err       error
		variables map[string]string
	}{
		{
			fileName: "fixtures/valid/comments.scl",
			hcl: `# This should be passed through
block {
  value = 1
}`,
		},
		{
			fileName: "fixtures/valid/basic.scl",
			hcl: `wrapper {
  inner = "yes"
  another {
    yet_another = "123"
  }
  inner "no"{}
}`,
		},
		{
			fileName: "fixtures/valid/heredoc.scl",
			hcl: `container {
  foo = <<EOF
bar
    indent
baz
EOF

  bar = <<DOC
	hello
		DOC

}`,
		},
		{
			fileName: "fixtures/valid/recursion.scl",
			hcl: `wrapper {
  decl "GET" "/base/default" {
    parent_id = 999
  }
  route "sub0" {
    id = 0
    parent_id = 999
    decl "GET" "/base/sub0/default" {
      parent_id = 0
    }
    route "sub1" {
      id = 1
      parent_id = 0
      decl "GET" "/base/sub0/sub1/default" {
        parent_id = 1
      }
      route "sub2" {
        id = 2
        parent_id = 1
        decl "GET" "/base/sub0/sub1/sub2/default" {
          parent_id = 2
        }
      }
    }
  }
}`,
		},
		{
			fileName: "fixtures/valid/variables.scl",
			hcl: `outer {
  inner = "hello"
}`,
			variables: map[string]string{
				"myVar": `"hello"`,
			},
		},
		{
			fileName: "fixtures/valid/variable-assignment.scl",
			hcl: `outer "normal assignment" {
  inner1 = "hello"
  inner2 = "hello world"
  inner3 = <<EOF
{
	"a": "b"
}
EOF

}
outer "scope-specific re-assignment should affect parent scope" {
  inner = "world hello"
}
outer "parent scope affected" {
  inner = "world hello"
}
outer "new declaration" {
  inner = "something"
}
origin = "parent value"
t1 = "http://localhost"
t2 = "http://localhost"`,
		},
		{
			fileName: "fixtures/valid/mixin-declaration.scl",
			hcl: `outer {
  someLiteral = "hello"
  wrapper {
    someOtherLiteral = "world"
    nestedLiteral = "world"
  }
  myCustomLiteral = "something"
  someArg = "else"
}`,
		},
		{
			fileName: "fixtures/valid/import.scl",
			hcl: `wrapper {
  inner = "yes"
  another {
    yet_another = "123"
  }
  inner "no"{}
}
output = "this is from simpleMixin"`,
		},
		{
			fileName: "fixtures/valid/short-function.scl",
			hcl: `hello = 123
mixin {
  hello = 456
}`,
		},
		{
			fileName: "fixtures/valid/optional-arguments.scl",
			hcl: `required = "1"
optional = "default"
required = "2"
optional = "non-default"`,
		},
		{
			fileName: "fixtures/valid/mixin-array-calls.scl",
			hcl: `a0 = [1,2,3]
a1 = [4,5,6]`,
		},
		{
			fileName: "fixtures/valid/mixin-underscore-param.scl",
			hcl: `var0 = 0
var1 = 1
var2 = 2
var0 = 1
var1 = 2
var2 = 3`,
		},
		{
			fileName: "fixtures/valid/docblock.scl",
			hcl:      ``,
		},
		{
			fileName: "fixtures/valid/vendor.scl",
			hcl:      `this = "included from vendor"`,
		},
		{
			fileName: "fixtures/invalid/heredoc.scl",
			err:      fmt.Errorf("Can't scan fixtures/invalid/heredoc.scl: Heredoc 'DOC' (started line 7) not terminated"),
		},
		{
			fileName: "fixtures/invalid/optional-arguments.scl",
			err:      fmt.Errorf("[fixtures/invalid/optional-arguments.scl:1] Argument declaration 1 [required]: A required argument can't follow an optional argument"),
		},
		{
			fileName: "fixtures/valid/variables.scl",
			err:      fmt.Errorf("[fixtures/valid/variables.scl:2] Unknown variable '$myVar'"),
		},
		{
			fileName: "this/doesnt/exist",
			err:      fmt.Errorf("Can't read this/doesnt/exist: open this/doesnt/exist: no such file or directory"),
		},
		{
			fileName: "fixtures/invalid/unknownToken.scl",
			err:      fmt.Errorf("[fixtures/invalid/unknownToken.scl:2] Unknown token: 1:11 IDENT yes"),
		},
		{
			fileName: "fixtures/invalid/illegalToken.scl",
			err:      fmt.Errorf("[fixtures/invalid/illegalToken.scl:1] illegal char"),
		},
		{
			fileName: "fixtures/invalid/mixin-declaration.scl",
			err:      fmt.Errorf("[fixtures/invalid/mixin-declaration.scl:1] Argument declaration 1 [v2]: Unexpected literal"),
		},
		{
			fileName: "fixtures/invalid/mixin-scope.scl",
			err:      fmt.Errorf("[fixtures/invalid/mixin-scope.scl:1] Mixin doesntExist not declared in this scope"),
		},
		{
			fileName: "fixtures/invalid/mixin-arguments.scl",
			err:      fmt.Errorf("[fixtures/invalid/mixin-arguments.scl:4] Wrong number of arguments for validMixin (required 2, got 3)"),
		},
		{
			fileName: "fixtures/invalid/mixin-argument-scope.scl",
			err:      fmt.Errorf("[fixtures/invalid/mixin-argument-scope.scl:4] Variable $myArg is not declared in this scope"),
		},
		{
			fileName: "fixtures/invalid/mixin-scope-nested.scl",
			err:      fmt.Errorf("[fixtures/invalid/mixin-scope-nested.scl:4] Mixin child not declared in this scope"),
		},
		{
			fileName: "fixtures/invalid/import.scl",
			err:      fmt.Errorf("[fixtures/invalid/import.scl:1] Can't read this/doesnt/exist.scl: no files found"),
		},
		{
			fileName: "fixtures/invalid/mixin-argument-scope.scl",
			err:      fmt.Errorf("[fixtures/invalid/mixin-argument-scope.scl:4] Variable $myArg is not declared in this scope"),
		},
		{
			fileName: "fixtures/invalid/error-in-include.scl",
			err:      fmt.Errorf("[fixtures/invalid/error-in-include.scl:1] [fixtures/invalid/illegalToken.scl:1] illegal char"),
		},
	} {
		t.Logf("Cycle %d", cycle)

		p := newMockParser(t)

		if input.variables != nil {
			for k, v := range input.variables {
				p.rootScope.setVariable(k, v)
			}
		}

		e := p.Parse(input.fileName)

		require.Equal(t, input.err, e)

		if input.err == nil {
			require.Equal(t, input.hcl, p.String())
			assert.Nil(t, hcl.Decode(&struct{}{}, p.String()))
		}
	}
}

func Test_AParserCanExtractACommentTree(t *testing.T) {

	expected := MixinDocs{
		MixinDoc{
			Name:      "mixin",
			File:      "fixtures/valid/docblock.scl",
			Line:      8,
			Reference: "fixtures/valid/docblock.scl:8",
			Signature: "@mixin()",
			Docs:      "this is a document block\n\n```\ncode\n```",
			Children: MixinDocs{
				MixinDoc{
					Name:      "inside0",
					File:      "fixtures/valid/docblock.scl",
					Line:      20,
					Reference: "fixtures/valid/docblock.scl:20",
					Signature: "@inside0($var)",
					Docs:      "part 1\npart 2",
				},
				MixinDoc{
					Name:      "inside1",
					File:      "fixtures/valid/docblock.scl",
					Line:      23,
					Reference: "fixtures/valid/docblock.scl:23",
					Signature: "@inside1($var, $var2)",
				},
			},
		},
		MixinDoc{
			Name:      "mixin2",
			File:      "fixtures/valid/docblock.scl",
			Line:      30,
			Reference: "fixtures/valid/docblock.scl:30",
			Signature: "@mixin2($var)",
			Children: MixinDocs{
				MixinDoc{
					Name:      "inside0",
					File:      "fixtures/valid/docblock.scl",
					Line:      36,
					Reference: "fixtures/valid/docblock.scl:36",
					Signature: "@inside0($var)",
					Docs:      "This is a mixin inside mixin2",
				},
			},
		},
	}

	p := newMockParser(t)
	docs, err := p.Documentation("fixtures/valid/docblock.scl")
	require.Nil(t, err)
	require.Equal(t, expected, docs)
}

func printCommentTree(docs MixinDocs, indentation int) {

	for _, d := range docs {
		fmt.Printf("[%s:%02d] %s> %s  : %+v\n", d.File, d.Line, strings.Repeat("-", indentation*2), d.Signature, strings.Replace(d.Docs, "\n", "\\n", -1))
		printCommentTree(d.Children, indentation+1)
	}
}

/*func Test_Test(t *testing.T) {
	p := newMockParser(t)
	require.Nil(t, p.Parse("fixtures/valid/callback.scl"))
	fmt.Println(p.String())
}*/
