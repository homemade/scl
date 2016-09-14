package scl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_CreateLineReturnsNonNil(t *testing.T) {
	require.NotNil(t, newLine("test.scl", 1, 1, ""))
}
