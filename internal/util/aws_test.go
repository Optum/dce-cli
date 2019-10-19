package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAWSUtil(t *testing.T) {

	t.Run("UploadDirectoryToS3 should upload an entire directory to S3", func(t *testing.T) {
		require.Equal(t, 2, 1)
	})
	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(t, 123, 456, "they should not be equal")

}
