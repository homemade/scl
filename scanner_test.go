package scl

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func compareScannerTrees(t *testing.T, a, b scannerTree, indent int) {
	require.Equal(t, len(a), len(b), fmt.Sprintf("Indent %d: Trees should be of equal length", indent))

	for i := 0; i < len(a); i++ {
		require.Equal(t, a[i].line, b[i].line)
		require.Equal(t, a[i].column, b[i].column)
		require.Equal(t, a[i].content, b[i].content)

		compareScannerTrees(t, a[i].children, b[i].children, indent+1)
	}
}

func Test_CreateScannerReturnsNonNil(t *testing.T) {
	require.NotNil(t, newScanner(&bytes.Buffer{}))
}

func Test_ScannerLoadsInputIntoTokensGivenValidInput(t *testing.T) {

	for cycle, v := range []struct {
		file          string
		expected_tree scannerTree
	}{
		{
			file: `
@field($name, $label)
    field $name
        label = $label
        body()

        @option($value, $label)
                option
                    value = $value
                    label = $label
`,
			expected_tree: scannerTree{
				&scannerLine{line: 2, column: 0, content: "@field($name, $label)", children: []*scannerLine{
					&scannerLine{line: 3, column: 4, content: "field $name", children: []*scannerLine{
						&scannerLine{line: 4, column: 8, content: "label = $label"},
						&scannerLine{line: 5, column: 8, content: "body()"},
						&scannerLine{line: 7, column: 8, content: "@option($value, $label)", children: []*scannerLine{
							&scannerLine{line: 8, column: 16, content: "option", children: []*scannerLine{
								&scannerLine{line: 9, column: 20, content: "value = $value"},
								&scannerLine{line: 10, column: 20, content: "label = $label"},
							}}}},
					}},
				}},
			},
		},
		{
			file: `
@field($name, $label) {
    field $name {
        label = $label 
        body()

        @option($value, $label) {
                option {
                    value = $value 
                    label = $label
				}
		}	
}`,
			expected_tree: scannerTree{
				&scannerLine{line: 2, column: 0, content: "@field($name, $label)", children: []*scannerLine{
					&scannerLine{line: 3, column: 4, content: "field $name", children: []*scannerLine{
						&scannerLine{line: 4, column: 8, content: "label = $label"},
						&scannerLine{line: 5, column: 8, content: "body()"},
						&scannerLine{line: 7, column: 8, content: "@option($value, $label)", children: []*scannerLine{
							&scannerLine{line: 8, column: 16, content: "option", children: []*scannerLine{
								&scannerLine{line: 9, column: 20, content: "value = $value"},
								&scannerLine{line: 10, column: 20, content: "label = $label"},
							}}}},
					}},
				}},
			},
		},
	} {
		t.Logf("Cycle %d", cycle)

		lex := newScanner(bytes.NewBufferString(v.file))
		lines, err := lex.scan()
		require.Nil(t, err)
		compareScannerTrees(t, lines, v.expected_tree, 0)
	}
}
