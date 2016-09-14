package scl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ATokeniserCanStripCommentsFromALine(t *testing.T) {
	for cycle, input := range []struct {
		line   *scannerLine
		result string
	}{{
		line:   newLine("test.scl", 1, 0, "something"),
		result: "something",
	},
		{
			line:   newLine("test.scl", 1, 0, "//comment"),
			result: "",
		},
		{
			line:   newLine("test.scl", 1, 0, "something // comment"),
			result: "something",
		},
		{
			line:   newLine("test.scl", 1, 0, `"something" // comment`),
			result: `"something"`,
		},
		{
			line:   newLine("test.scl", 1, 0, `"something // else" // comment`),
			result: `"something // else"`,
		},
		{
			line:   newLine("test.scl", 1, 0, `not / a / comment`),
			result: `not / a / comment`,
		},
		{
			line:   newLine("test.scl", 1, 0, `not / a // comment`),
			result: `not / a`,
		},
		{
			line:   newLine("test.scl", 1, 0, `"not / a // comment"`),
			result: `"not / a // comment"`,
		},
		{
			line:   newLine("test.scl", 1, 0, `'not / a // comment'`),
			result: `'not / a // comment'`,
		},
	} {
		t.Logf("Cycle %d", cycle)

		tkn := newTokeniser()
		result := tkn.stripComments(input.line)
		require.Equal(t, input.result, result)
	}
}

func Test_ATokeniserCanTokeniseFunctionStrings(t *testing.T) {

	ln := newLine("test.scl", 1, 0, "@mixin()")

	for cycle, input := range []struct {
		line   *scannerLine
		input  string
		name   string
		tokens []token
		err    error
	}{
		{
			line:  ln,
			input: "fn()",
			name:  "fn",
		},
		{
			line:  ln,
			input: "fn():",
			name:  "fn",
		},
		{
			line:  ln,
			input: `fn(1,$two,  "th'r'ee" ,'four')`,
			name:  "fn",
			tokens: []token{
				{
					kind:    tokenLiteral,
					content: "1",
					line:    ln,
				},
				{
					kind:    tokenVariable,
					content: `two`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `"th'r'ee"`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `'four'`,
					line:    ln,
				},
			},
		},
		{
			line:  ln,
			input: "fn(",
			name:  "",
			err:   fmt.Errorf("Can't parse function signature"),
		},
		{
			line:  ln,
			input: `fn($required, $optional0=123, $optional1="123")`,
			name:  "fn",
			tokens: []token{
				{
					kind:    tokenVariable,
					content: `required`,
					line:    ln,
				},
				{
					kind:    tokenVariableAssignment,
					content: `optional0`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `123`,
					line:    ln,
				},
				{
					kind:    tokenVariableAssignment,
					content: `optional1`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `"123"`,
					line:    ln,
				},
			},
		},
		{
			line:  ln,
			input: `fn(1, "two", [1,2,3], "four,five")`,
			name:  "fn",
			tokens: []token{
				{
					kind:    tokenLiteral,
					content: "1",
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `"two"`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `[1,2,3]`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `"four,five"`,
					line:    ln,
				},
			},
		},
		{
			line:  ln,
			input: `fn($a, "two", [1,2,3],"[four,five]")`,
			name:  "fn",
			tokens: []token{
				{
					kind:    tokenVariable,
					content: `a`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `"two"`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `[1,2,3]`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `"[four,five]"`,
					line:    ln,
				},
			},
		},
		{
			line:  ln,
			input: `fn($a, "two", $optional0=[1,2,3])`,
			name:  "fn",
			tokens: []token{
				{
					kind:    tokenVariable,
					content: `a`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `"two"`,
					line:    ln,
				},
				{
					kind:    tokenVariableAssignment,
					content: `optional0`,
					line:    ln,
				},
				{
					kind:    tokenLiteral,
					content: `[1,2,3]`,
					line:    ln,
				},
			},
		},
	} {
		t.Logf("Cycle %d", cycle)

		tkn := newTokeniser()
		name, tokens, err := tkn.tokeniseFunction(input.line, input.input)

		require.Equal(t, input.err, err)
		require.Equal(t, input.name, name)

		for _, k := range tokens {
			t.Logf("%s: %s", k.kind, k.content)
		}

		require.Equal(t, input.tokens, tokens)
	}
}

func Test_ATokeniserCanReadTokensFromAValidScannerLine(t *testing.T) {

	var literalLine = newLine("test.scl", 1, 0, "literal")
	var literalLineWithComment = newLine("test.scl", 1, 0, "literal // This is a comment")
	var commentLine = newLine("test.scl", 1, 0, "// This is a line comment")
	var mixinDeclarationLine1 = newLine("test.scl", 1, 0, "@mixin()")
	var mixinDeclarationLine2 = newLine("test.scl", 1, 0, "@mixin($a,$b)")
	var mixinDeclarationLine3 = newLine("test.scl", 1, 0, "@mixin)")
	var mixinDeclarationLine4 = newLine("test.scl", 1, 0, `@mixin($a,"123")`)
	var functionCallLine1 = newLine("test.scl", 1, 0, `fn($a,"123")`)
	var shortFunctionCallLine1 = newLine("test.scl", 1, 0, `fn:`)
	var assignmentLine = newLine("test.scl", 1, 0, `$a = "123"`)

	for cycle, input := range []struct {
		line   *scannerLine
		tokens []token
		err    error
	}{
		{
			line: literalLine,
			tokens: []token{
				token{
					kind:    tokenLiteral,
					content: "literal",
					line:    literalLine,
				},
			},
		},
		{
			line: literalLineWithComment,
			tokens: []token{
				token{
					kind:    tokenLiteral,
					content: "literal",
					line:    literalLineWithComment,
				},
			},
		},
		{
			line: commentLine,
			tokens: []token{
				token{
					kind:    tokenLineComment,
					content: "This is a line comment",
					line:    commentLine,
				},
			},
		},
		{
			line: mixinDeclarationLine1,
			tokens: []token{
				token{
					kind:    tokenMixinDeclaration,
					content: "mixin",
					line:    mixinDeclarationLine1,
				},
			},
		},
		{
			line: mixinDeclarationLine2,
			tokens: []token{
				token{
					kind:    tokenMixinDeclaration,
					content: "mixin",
					line:    mixinDeclarationLine2,
				},
				{
					kind:    tokenVariable,
					content: `a`,
					line:    mixinDeclarationLine2,
				},
				{
					kind:    tokenVariable,
					content: `b`,
					line:    mixinDeclarationLine2,
				},
			},
		},
		{
			line: mixinDeclarationLine4,
			tokens: []token{
				token{
					kind:    tokenMixinDeclaration,
					content: "mixin",
					line:    mixinDeclarationLine4,
				},
				{
					kind:    tokenVariable,
					content: `a`,
					line:    mixinDeclarationLine4,
				},
				{
					kind:    tokenLiteral,
					content: `"123"`,
					line:    mixinDeclarationLine4,
				},
			},
		},
		{
			line: mixinDeclarationLine3,
			err:  fmt.Errorf("test.scl:1: Can't parse function signature"),
		},
		{
			line: functionCallLine1,
			tokens: []token{
				token{
					kind:    tokenFunctionCall,
					content: "fn",
					line:    functionCallLine1,
				},
				{
					kind:    tokenVariable,
					content: `a`,
					line:    functionCallLine1,
				},
				{
					kind:    tokenLiteral,
					content: `"123"`,
					line:    functionCallLine1,
				},
			},
		},
		{
			line: shortFunctionCallLine1,
			tokens: []token{
				token{
					kind:    tokenFunctionCall,
					content: "fn",
					line:    shortFunctionCallLine1,
				},
			},
		},
		{
			line: assignmentLine,
			tokens: []token{
				token{
					kind:    tokenVariableAssignment,
					content: "a",
					line:    assignmentLine,
				},
				{
					kind:    tokenLiteral,
					content: `"123"`,
					line:    assignmentLine,
				},
			},
		},
	} {
		t.Logf("Cycle %d", cycle)

		tkn := newTokeniser()

		tokens, err := tkn.tokenise(input.line)

		require.Equal(t, err, input.err)
		require.Equal(t, tokens, input.tokens)
	}
}
