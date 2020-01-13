package util

import (
	"bytes"
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
	"github.com/pkg/errors"

	"github.com/mitchellh/cli"
)

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
	cmd.Stdin = os.Stdin
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
	Unarchive(source string, destination string)
	GetTerraformBin() string
	RemoveAll(path string)
	GetTerraformBinDir() string
	GetLocalTFModuleDir() string
}

type TerraformBinUtil struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	FileSystem  TerraformBinFileSystemUtil
	Downloader  TerraformBinDownloader
}

// bin returns the binary path
func (t *TerraformBinUtil) bin() string {
	// pull it out of the Config, or default to ~/.dce/terraform
	bin := t.Config.Terraform.Bin

	if bin == nil || len(*bin) == 0 {
		s := t.FileSystem.GetTerraformBin()
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
func (t *TerraformBinUtil) Init(ctx context.Context, args []string) error {
	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value("deployLogFile").(string))

	if err != nil {
		logFile = nil
	} else {
		defer logFile.Close()
	}

	argv := []string{"init", "-no-color"}
	argv = append(argv, args...)

	if !t.FileSystem.IsExistingFile(t.bin()) {
		archive := fmt.Sprintf("%s.zip", t.bin())
		err := t.Downloader.Download(t.source(), archive)
		if err != nil {
			return err
		}
		t.FileSystem.Unarchive(archive, t.FileSystem.GetTerraformBinDir())
		// make sure the file is there and executable.
		if !t.FileSystem.IsExistingFile(t.bin()) {
			return fmt.Errorf("%s does not exist", t.bin())
		}
		t.FileSystem.RemoveAll(archive)
	}

	// at this point, the binary should exist. Call `init`
	execArgs := &execInput{
		Name: t.bin(),
		Args: argv,
		Dir:  t.FileSystem.GetLocalTFModuleDir(),
	}

	return execCommand(execArgs, logFile, logFile)
}

// Apply will call `terraform apply` with the given vars.
func (t *TerraformBinUtil) Apply(ctx context.Context, tfVars []string) error {
	cfg := ctx.Value(constants.DeployConfig).(*configs.DeployConfig)
	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value("deployLogFile").(string))

	if err != nil {
		logFile = nil
	} else {
		defer logFile.Close()
	}

	if !t.FileSystem.IsExistingFile(t.bin()) {
		return fmt.Errorf("Could not find binary \"%s\"", t.bin())
	}

	argv := []string{"apply", "-no-color"}

	if cfg.BatchMode {
		argv = append(argv, "-auto-approve")
	} else {
		// The underlying terraform command's stdin is set to this stdin,
		// so  the user's answer here is passes along to terraform.
		fmt.Print("Are you sure you would like to create DCE resources? (must type \"yes\" if yes)\t")
	}

	for _, tfVar := range tfVars {
		argv = append(argv, "-var", tfVar)
	}

	execArgs := &execInput{
		Name: t.bin(),
		Args: argv,
		Dir:  t.FileSystem.GetLocalTFModuleDir(),
	}

	return execCommand(execArgs, logFile, logFile)

}

// GetOutput returns the value of the output with the given name.
func (t *TerraformBinUtil) GetOutput(ctx context.Context, key string) (string, error) {

	// for the output, we now use a byte buyffer for the output
	// but keep the stderr as the log file so advanced users can
	// diagnose issues.
	var stdout bytes.Buffer

	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value("deployLogFile").(string))

	if err != nil {
		logFile = nil
	} else {
		defer logFile.Close()
	}

	// Run `terraform output` command
	err = execCommand(&execInput{
		Name: t.bin(),
		Args: []string{
			"output",
			key,
			"-no-color",
		},
		Dir: t.FileSystem.GetLocalTFModuleDir(),
	},
		&stdout,
		logFile)

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(stdout.String()), nil
}
