package unit

import (
	"os"
	"testing"
)

var doesntMatter = "doesntmatter"

type mockFileInfo struct {
	os.FileInfo
}

func (m *mockFileInfo) Name() string { return doesntMatter }
func (m *mockFileInfo) IsDir() bool  { return true }

func TestDeployTFOverrides(t *testing.T) {
	// TODO: Write the new unit tests here... so much as changed these need to
	// be mocked differently, etc.
}
