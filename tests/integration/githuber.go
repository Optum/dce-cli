package integration

import (
	"fmt"
	"github.com/Optum/dce-cli/mocks"
	"io/ioutil"
	"testing"
)

// stubGithub is a stub of the the Githuber util,
// which may be used to mock github releases
type stubGithub struct {
	*mocks.Githuber
	// map of "assetName/dceVersion" --> file content
	mockAssets map[string][]byte
}

func (gh *stubGithub) MockReleaseAsset(t *testing.T, assetName string, dceVersion string, content []byte) {
	if gh.mockAssets == nil {
		gh.mockAssets = map[string][]byte{}
	}

	key := fmt.Sprintf("%s/%s", assetName, dceVersion)
	gh.mockAssets[key] = content
}

func (gh *stubGithub) DownloadGithubReleaseAsset(assetName string, dceVersion string) error {
	// Grab the mocked asset
	key := fmt.Sprintf("%s/%s", assetName, dceVersion)
	content, ok := gh.mockAssets[key]
	if !ok {
		return fmt.Errorf("Test failure: mocked Github asset does not exist "+
			"for %s @ %s", assetName, dceVersion)
	}

	// Write the file
	// #nosec
	err := ioutil.WriteFile(assetName, content, 0666)
	return err
}
