package integration

import (
	"github.com/Optum/dce-cli/configs"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/ptr"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestInitCommand(t *testing.T) {

	t.Run("GIVEN custom config path flag is provided", func(t *testing.T) {

		t.Run("AND config file does not exists", func(t *testing.T) {

			t.Run("THEN config is written to file at specified path", func(t *testing.T) {
				cli := NewCLITest(t)

				// Answer interactive CLI prompts
				cli.AnswerBasic(
					"Host name of the DCE API (example: abcde12345.execute-api.us-east-1.amazonaws.com)",
					"dce.example.com",
				)
				cli.AnswerBasic(
					"Base path of the DCE API (example: /api)",
					"/api",
				)

				// Create a tmp dir for our dce.yml config to live in
				tmpdir, err := ioutil.TempDir("", "dce-cli-test")
				require.Nil(t, err)
				defer os.RemoveAll(tmpdir)
				confFile := path.Join(tmpdir, "dce.yml")

				// Run `dce init`
				err = cli.Execute([]string{"init", "--config", confFile})
				require.Nil(t, err)

				cli.AssertAllPrompts()

				// Check that we wrote to the dce.yml config file
				assertYamlConfig(t, &configs.Root{
					API: configs.API{
						Host:     ptr.String("dce.example.com"),
						BasePath: ptr.String("/api"),
					},
					Region: ptr.String("us-east-1"),
				}, confFile)

				// Check CLI outputs
				require.Contains(t, cli.Output(), "Config file created at:")
			})

		})

		t.Run("AND config file is empty", func(t *testing.T) {

			t.Run("THEN config is written to file at specified path", func(t *testing.T) {
				cli := NewCLITest(t)

				// Answer interactive CLI prompts
				cli.AnswerBasic(
					"Host name of the DCE API (example: abcde12345.execute-api.us-east-1.amazonaws.com)",
					"dce.example.com",
				)
				cli.AnswerBasic(
					"Base path of the DCE API (example: /api)",
					"/api",
				)

				// Create a tmp dir for our dce.yml config to live in
				tmpdir, err := ioutil.TempDir("", "dce-cli-test")
				require.Nil(t, err)
				defer os.RemoveAll(tmpdir)
				confFile := path.Join(tmpdir, "dce.yml")

				// Create empty dce.yml file
				file, err := os.Create(confFile)
				require.Nil(t, err)
				defer file.Close()
				defer os.Remove(file.Name())

				// Run `dce init`
				err = cli.Execute([]string{"init", "--config", confFile})
				require.Nil(t, err)

				cli.AssertAllPrompts()

				// Check that we wrote to the dce.yml config file
				assertYamlConfig(t, &configs.Root{
					API: configs.API{
						Host:     ptr.String("dce.example.com"),
						BasePath: ptr.String("/api"),
					},
					Region: ptr.String("us-east-1"),
				}, confFile)

				// Check CLI outputs
				require.Contains(t, cli.Output(), "Config file created at:")
			})

		})

		t.Run("AND config file has existing YAML", func(t *testing.T) {

			t.Run("THEN config YAML is updated from CLI prompts", func(t *testing.T) {
				confFile := writeTempConfig(t, &configs.Root{
					API: configs.API{
						Host:  ptr.String("dce.example.com"),
						Token: ptr.String("my-api-token"),
					},
					Region:      ptr.String("us-west-2"),
					GithubToken: ptr.String("my-gh-token"),
				})

				cli := NewCLITest(t)

				// Answer interactive CLI prompts
				// updating existing config values
				// Note that we'll only prompt for any fields which aren't
				// already set
				cli.AnswerBasic(
					"Base path of the DCE API (example: /api)",
					"/api-new",
				)

				// Run `dce init`
				err := cli.Execute([]string{"init", "--config", confFile})
				require.Nil(t, err)

				// Check that YAML config was updated
				assertYamlConfig(t, &configs.Root{
					API: configs.API{
						Host: ptr.String("dce.example.com"),
						// Should modify from CLI prompts
						BasePath: ptr.String("/api-new"),
						Token:    ptr.String("my-api-token"),
					},
					Region:      ptr.String("us-west-2"),
					GithubToken: ptr.String("my-gh-token"),
				}, confFile)
			})

		})

		t.Run("AND config file has invalid YAML content", func(t *testing.T) {

			t.Run("THEN command fails", func(t *testing.T) {
				// Write some garbage to a file
				// Create a tmp file
				tmpfile, err := ioutil.TempFile("", "dce.*.yml")
				require.Nil(t, err)
				err = ioutil.WriteFile(tmpfile.Name(), []byte("not valid YAML"), 0644)
				require.Nil(t, err)
				_ = tmpfile.Close()
				defer os.Remove(tmpfile.Name())

				// Run `dce init` (should fail)
				cli := NewCLITest(t)
				err = cli.Execute([]string{"init", "--config", tmpfile.Name()})
				require.NotNil(t, err)
				require.Contains(t, err.Error(), "Failed to parse dce.yml:")
			})

		})
	})

}

func assertYamlConfig(t *testing.T, expectedConf *configs.Root, yamlFile string) {
	yamlStr, err := ioutil.ReadFile(yamlFile)
	require.Nilf(t, err, "Failed to read %s", yamlFile)

	var actualConf configs.Root
	err = yaml.Unmarshal(yamlStr, &actualConf)
	require.Nilf(t, err, "Failed to parse YAML for %s", yamlFile)

	require.Equal(t, expectedConf, &actualConf)
}
