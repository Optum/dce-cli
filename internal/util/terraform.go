package util

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"runtime"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	tfBackendInit "github.com/hashicorp/terraform/backend/init"
	tfCommand "github.com/hashicorp/terraform/command"
	tfDiscovery "github.com/hashicorp/terraform/svchost/disco"
	"github.com/pkg/errors"

	"github.com/mitchellh/cli"
)

type TerraformUtil struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
}

// Init initialized a terraform working directory
func (u *TerraformUtil) Init(ctx context.Context, args []string) {
	logFile, err := os.Create(ctx.Value("deployLogFile").(string))
	log.Println("Running terraform init")
	// logger.SetOutput(logFile)

	if err != nil {
		logFile = nil
	} else {
		defer logFile.Close()
	}

	services := tfDiscovery.NewWithCredentialsSource(nil)
	tfBackendInit.Init(services)

	tfInit := &tfCommand.InitCommand{
		Meta: tfCommand.Meta{
			Ui: getTerraformUI(logFile),
		},
	}
	tfInit.Run(args)
}

// Apply applies terraform template with given namespace
func (u *TerraformUtil) Apply(ctx context.Context, tfVars []string) {
	cfg := ctx.Value(constants.DeployConfig).(*configs.DeployConfig)
	logFile, err := os.Create(ctx.Value("deployLogFile").(string))

	if err != nil {
		logFile = nil
	} else {
		defer logFile.Close()
	}

	// logger.SetOutput(logFile)
	tfApply := &tfCommand.ApplyCommand{
		Meta: tfCommand.Meta{
			Ui:                  getTerraformUI(logFile),
			RunningInAutomation: true,
		},
	}

	runArgs := []string{}
	for _, tfVar := range tfVars {
		runArgs = append(runArgs, "-var", tfVar)
	}

	if cfg.NoPrompt {
		runArgs = append(runArgs, "-auto-approve")
	}

	log.Debugln("Args for Apply command: ", runArgs)
	tfApply.Run(runArgs)
}

// GetOutput gets terraform output value for provided key
func (u *TerraformUtil) GetOutput(key string) string {
	log.Println("Retrieving terraform output for: " + key)
	outputCaptorUI := &UIOutputCaptor{
		BasicUi: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		Captor: new(string),
	}
	tfOutput := &tfCommand.OutputCommand{
		Meta: tfCommand.Meta{
			Ui: outputCaptorUI,
		},
	}
	tfOutput.Run([]string{key})
	return *outputCaptorUI.Captor
}

func getTerraformUI(f *os.File) *cli.BasicUi {
	var out *os.File

	if f != nil {
		out = f
	} else {
		out = os.Stdout
	}

	return &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      out,
		ErrorWriter: out,
	}
}

// UIOutputCaptor effectively extends cli.BasicUi and overrides Output method to capture output string.
type UIOutputCaptor struct {
	Captor *string
	*cli.BasicUi
}

// Output overrides cli.BasicUi Output method in UIOutputCaptor
func (u *UIOutputCaptor) Output(message string) {
	u.Captor = &message
	u.BasicUi.Output(message)
}

type execInput struct {
	Name    string   // Command to execute
	Args    []string // Arguments to pass to the command
	Dir     string   // Working directory
	Timeout float64  // Max execution time (seconds) of the command
}

func execCommand(input *execInput, stdout io.Writer, stderr io.Writer) error {
	// Create a context, in order to enforce a Timeout on the command.
	// See https://medium.com/@vCabbage/go-timeout-commands-with-os-exec-commandcontext-ba0c861ed738
	// and https://siadat.github.io/post/context
	var ctx context.Context
	var cancel context.CancelFunc
	if input.Timeout == 0 {
		// If no Timeout is configured, use and empty context
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(
			context.Background(),
			time.Duration(input.Timeout)*time.Second,
		)
	}

	// Cleanup context, on completion
	defer cancel()

	// Configure the shell command
	cmd := exec.CommandContext(ctx, input.Name, input.Args...)
	if input.Dir != "" {
		cmd.Dir = input.Dir
	}

	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()

	// Check if the command timed out
	cmdStr := strings.Join(
		append([]string{input.Name}, input.Args...), " ",
	)
	if ctx.Err() == context.DeadlineExceeded {
		return errors.Wrapf(ctx.Err(), "Command timed out: %s", cmdStr)
	}

	// Check for command errors
	if err != nil {
		return errors.Wrapf(err, "Command failed: %s", cmdStr)
	}

	return nil
}

type TerraformBinDownloader interface {
	Download(url string, localpath string) error
}

type TerraformBinFileSystemUtil interface {
	GetConfigDir() string
	IsExistingFile(path string) bool
	OpenFileWriter(path string) (*os.File, error)
}

type TerraformBinUtil struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	FileSystem  TerraformBinFileSystemUtil
	Downloader  TerraformBinDownloader
}

// bin returns the binary path
func (t *TerraformBinUtil) bin() string {
	// pull it out of the Config, or default to ~/.dce/bin/terraform
	bin := t.Config.Terraform.Bin

	if bin == nil || len(*bin) == 0 {
		s := t.FileSystem.GetConfigDir()
		return s
	}
	return *bin
}

// source returns the download URL for terraform binary.
func (t *TerraformBinUtil) source() string {
	source := t.Config.Terraform.Source

	if source == nil || len(*source) == 0 {
		s := fmt.Sprintf(constants.TerraformBinDownloadURLFormat,
			constants.TerraformBinVersion,
			constants.TerraformBinVersion,
			runtime.GOOS,
			runtime.GOARCH,
		)
		return s
	}
	return *source
}

// Init will download the Terraform binary, put it into the .dce folder,
// and then call init.
func (t *TerraformBinUtil) Init(ctx context.Context, args []string) {
	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value("deployLogFile").(string))

	if err != nil {
		logFile = nil
	} else {
		defer logFile.Close()
	}

	argv := []string{"init", "-nocolor"}
	argv = append(argv, args...)

	if !t.FileSystem.IsExistingFile(t.bin()) {
		err := t.Downloader.Download(t.source(), t.bin())
		if err != nil {
			log.Fatalln(err)
		}
	}

	// at this point, the binary should exist. Call `init`
	execArgs := &execInput{
		Name: t.bin(),
		Args: argv,
		Dir:  t.FileSystem.GetConfigDir(),
	}

	err = execCommand(execArgs, logFile, logFile)

	if err != nil {
		log.Fatalln(err)
	}
}

// Apply will call `terraform apply` with the given vars.
func (t *TerraformBinUtil) Apply(ctx context.Context, tfVars []string) {
	cfg := ctx.Value(constants.DeployConfig).(*configs.DeployConfig)
	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value("deployLogFile").(string))

	if err != nil {
		logFile = nil
	} else {
		defer logFile.Close()
	}

	if !t.FileSystem.IsExistingFile(t.bin()) {
		err := t.Downloader.Download(t.source(), t.bin())
		if err != nil {
			log.Fatalln(err)
		}
	}

	argv := []string{"apply"}

	if cfg.NoPrompt {
		argv = append(argv, "-auto-approve")
	}

	for _, tfVar := range tfVars {
		argv = append(argv, "-var", tfVar)
	}

	// at this point, the binary should exist. Call `init`
	execArgs := &execInput{
		Name: t.bin(),
		Args: argv,
		Dir:  t.FileSystem.GetConfigDir(),
	}

	err = execCommand(execArgs, logFile, logFile)

	if err != nil {
		log.Fatalln(err)
	}
}

// GetOutput returns the value of the output with the given name.
func (t *TerraformBinUtil) GetOutput(key string) string {
	return ""
}
