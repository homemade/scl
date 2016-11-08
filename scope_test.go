package scl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AValidScopeCanInterpolateVariables(t *testing.T) {

	for cycle, input := range []struct {
		variables map[string]string
		literal   string
		result    string
		err       error
	}{
		{
			variables: map[string]string{
				"one": "value1",
				"two": "value2",
			},
			literal: "$one is $two",
			result:  "value1 is value2",
		},
		{
			variables: map[string]string{
				"one": "value1",
				"two": "value2",
			},
			literal: `\$one is $two, \\ is escaped`,
			result:  `$one is value2, \ is escaped`,
		},
		{
			variables: map[string]string{
				"one": "value1",
				"two": "value2",
			},
			literal: `${one}_${two}`,
			result:  `value1_value2`,
		},
		{
			variables: map[string]string{
				"one": "value1",
				"two": "value2",
			},
			literal: `${one}${two}`,
			result:  `value1value2`,
		},
		{
			variables: map[string]string{
				"one": "value1",
				"two": "value2",
			},
			literal: `${one`,
			err:     fmt.Errorf("Expecting closing right brace in variable ${one}"),
		},
		{
			variables: map[string]string{
				"one": "value1",
				"two": "value2",
			},
			literal: `${one is unfinished`,
			err:     fmt.Errorf("Expecting closing right brace in variable ${one}"),
		},
		{
			variables: map[string]string{
				"one": "value1",
				"two": "value2",
			},
			literal: "$one is $$two",
			result:  "value1 is $$two",
		},
		{
			variables: map[string]string{
				"one": "value1",
			},
			literal: "$ one is $one",
			result:  "$ one is value1",
		},
		{
			variables: map[string]string{
				"name": "myModel",
			},
			literal: `model "$name" {}`,
			result:  `model "myModel" {}`,
		},
		{
			variables: map[string]string{},
			err:       fmt.Errorf("Unknown variable '$nothing'"),
			literal:   `something = $nothing`,
			result:    `something = `,
		},
		{
			variables: map[string]string{},
			literal:   `This ends in a \`,
			result:    `This ends in a \`,
		},
		{
			variables: map[string]string{},
			literal:   "This is a `$confusing $literal \\`",
			result:    "This is a $confusing $literal \\",
		},
		{
			variables: map[string]string{},
			literal:   "This is \\`not a literal\\`",
			result:    "This is `not a literal`",
		},
		{
			variables: map[string]string{},
			literal:   "This is `not a literal",
			result:    "This is not a literal",
			err:       fmt.Errorf("Unterminated backtick literal"),
		},
	} {
		t.Logf("Cycle %d", cycle)

		s := newScope()

		for k, v := range input.variables {
			s.setVariable(k, v)
		}

		result, err := s.interpolateLiteral(input.literal)

		require.Equal(t, input.err, err)
		require.Equal(t, input.result, result)
	}
}
