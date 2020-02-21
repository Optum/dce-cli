package integration

import (
	"bytes"
	"github.com/Optum/dce-cli/cmd"
	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/observation"
	"github.com/Optum/dce-cli/internal/util"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/Optum/dce-cli/pkg/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/require"
	"go.uber.org/thriftrw/ptr"
	"io/ioutil"
	"testing"
)

// cliTest is a util for running integration
// tests against the CLI
type cliTest struct {
	*MockPrompter
	stdout *bytes.Buffer
	injector injector
	configFile string
}

func (test *cliTest) WriteConfig(t *testing.T, config *configs.Root) {
	test.configFile = writeTempConfig(t, config)
}

func (test *cliTest) Execute(args []string) error {
	// Pass in config file path, if we have one
	if test.configFile != "" {
		args = append(args, "--config", test.configFile)
	}

	cmd.RootCmd.SetArgs(args)
	return cmd.RootCmd.Execute()
}

func NewCLITest(t *testing.T) *cliTest {
	// Reset globals, so they don't leak between tests
	cmd.Config = &configs.Root{}
	cmd.Service = &service.ServiceContainer{}
	cmd.Util = &utl.UtilContainer{}
	cmd.Observation = &observation.ObservationContainer{}

	// Mock the Prompter, to allow use interactive inputs
	prompter := &MockPrompter{T: t}

	var stdout bytes.Buffer


	cli := &cliTest{
		MockPrompter: prompter,
		stdout:       &stdout,
	}

	// Wrap the `PreRun` method, to inject the mock prompter,
	// and out logger stdout buffer
	// (global services aren't initialized until `PreRun`
	preRun := cmd.RootCmd.PersistentPreRunE
	cmd.RootCmd.PersistentPreRunE = func(c *cobra.Command, a []string) error {
		// Run the wrapped method
		err := preRun(c, a)
		if err != nil {
			return err
		}

		// Inject the prompter
		cmd.Util.Prompter = prompter

		// Tell the logger to log to our stdout buffer
		logger := cmd.Log.(*observation.LogObservation).LevelLogger.(*logrus.Logger)
		logger.SetOutput(&stdout)

		// Allow client to inject their own mocks
		if cli.injector != nil {
			cli.injector(&injectorInput{cmd.Config, cmd.Service, cmd.Util, cmd.Observation})
		}

		return nil
	}

	return cli
}

func (test *cliTest) Output() string {
	return test.stdout.String()
}

type injectorInput struct {
	config *configs.Root
	service *service.ServiceContainer
	util *utl.UtilContainer
	observation *observation.ObservationContainer
}
type injector func(input *injectorInput)

// Inject allows for injecting dependencies before
// commands are run. The `inject` function provides
// pointers to global services used by CLI commands
//
// eg
// cli.Inject(func(input *injectorInput) {
//		input.util.Weber = &mocks.Weber{}
// })
func (test *cliTest) Inject(f injector) {
	test.injector = f
}

func writeTempConfig(t *testing.T, conf *configs.Root) string {
	// Create a tmp file
	tmpfile, err := ioutil.TempFile("", "dce.*.yml")
	require.Nil(t, err)

	// Set default config
	if conf == nil {
		conf = &configs.Root{
			API: configs.API{
				Host:     ptr.String("dce.example.com"),
				BasePath: ptr.String("/api"),
			},
		}
	}

	// Write config as YAML to tmp file
	fsUtil := util.FileSystemUtil{
		Config:     conf,
		ConfigFile: tmpfile.Name(),
	}
	err = fsUtil.WriteConfig()
	require.Nil(t, err)

	// Close the file handle
	err = tmpfile.Close()
	require.Nil(t, err)

	return tmpfile.Name()
}
