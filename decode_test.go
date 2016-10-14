package scl

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_AFileCanBeDecodedWithAnImplicitFileSystem(t *testing.T) {

	type decodable struct {
		Value0 string   `hcl:"value0"`
		Value1 int      `hcl:"value1"`
		Value2 []string `hcl:"value2"`
	}

	for cycle, test := range []struct {
		path     string
		expected decodable
		err      error
	}{
		{
			path: "fixtures/valid/decode.scl",
			expected: decodable{
				"1",
				1,
				[]string{"a", "b", "c"},
			},
		},
		{
			path: "fixtures/invalid/decode.scl",
			err:  fmt.Errorf("[fixtures/invalid/decode.scl:1] error parsing list, expected comma or list end, got: EOF"),
		},
	} {
		t.Logf("Cycle %d", cycle)

		got := decodable{}

		err := DecodeFile(&got, test.path)

		if test.err == nil {
			require.Nil(t, err)
			require.Equal(t, got, test.expected)
		} else {
			require.NotNil(t, err)
			require.Equal(t, test.err.Error(), err.Error())
		}
	}
}
