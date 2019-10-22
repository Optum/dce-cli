package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPipeline(t *testing.T) {

	t.Run("Test pipeline can run tests", func(t *testing.T) {
		require.Equal(t, 1, 1)
	})
}
