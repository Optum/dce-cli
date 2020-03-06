package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"runtime"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/pkg/errors"
)

var paramSplitRegex *regexp.Regexp

func init() {
	paramSplitRegex = regexp.MustCompile(`\s+`)
}

type execInput struct {
	Name    string   // Command to execute
	Args    []string // Arguments to pass to the command
	Dir     string   // Working directory
	Timeout float64  // Max execution time (seconds) of the command
}

// ParseOptions parses the given options into an array of strings. It provides for any whitespace between
// the options.
func ParseOptions(s *string) ([]string, error) {
	if s == nil || len(*s) == 0 {
		return []string{}, nil
	}
	opts := paramSplitRegex.Split(*s, -1)
	return opts, nil
}

// execCommand executes the command specified by `input` and writes
// output and STDERR to `stdout` or `stderr`, respectively. If either
// in nil, the OS STDOUT and STDERR are used.
// Care should be taken to mitigate CWE-78 (https://cwe.mitre.org/data/definitions/78.html)
// by ensuring that 'input' comes from a trusted source.
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

	/*
		#nosec CWE-78: added disclaimer to function docs
	*/
	// Configure the shell command
	cmd := exec.CommandContext(ctx, input.Name, input.Args...)
	if input.Dir != "" {
		cmd.Dir = input.Dir
	}

	if stdout == nil {
		log.Warnln("stdout: no file supplied; using STDOUT")
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = stdout
	}

	if stderr == nil {
		log.Warnln("stderr: no file supplied; using STDERR")
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stderr = stderr
	}

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

// TerraformBinDownloader - interface for the downloader to download the Terraform
// binary from a URL ans save it locally.
type TerraformBinDownloader interface {
	Download(url string, localpath string) error
}

// TerraformBinFileSystemUtil - interface for interacting with the file system.
type TerraformBinFileSystemUtil interface {
	GetConfigDir() string
	IsExistingFile(path string) bool
	OpenFileWriter(path string) (*os.File, error)
	Unarchive(source string, destination string) error
	GetTerraformBin() string
	RemoveAll(path string)
	GetTerraformBinDir() string
	GetLocalTFModuleDir() string
}

// TerraformBinUtil uses the Teraform binary to peform the init and apply
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
		t.Config.Terraform.Bin = &s
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
		t.Config.Terraform.Source = &s
		return s
	}
	return *source
}

// Init will download the Terraform binary, put it into the .dce folder,
// and then call init.
func (t *TerraformBinUtil) Init(ctx context.Context, args []string) error {
	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value(constants.DeployLogFileKey).(string))

	if err != nil {
		logFile = nil
	} else {
		// #nosec
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
		err = t.FileSystem.Unarchive(archive, t.FileSystem.GetTerraformBinDir())
		if err != nil {
			return err
		}
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
func (t *TerraformBinUtil) Apply(ctx context.Context, args []string) error {
	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value(constants.DeployLogFileKey).(string))

	if err != nil {
		logFile = nil
	} else {
		// #nosec
		defer logFile.Close()
	}

	if !t.FileSystem.IsExistingFile(t.bin()) {
		return fmt.Errorf("Could not find binary \"%s\"", t.bin())
	}

	// -auto-approve and -input=false used to be only added when
	// a user used the --batch-mode flag, but now that concern is taken
	// care of at a higher level. It is assumed if the process has gotten
	// this far that all necessary prompting and inputs have been collected.
	argv := []string{"apply", "-no-color", "-auto-approve", "-input=false"}
	argv = append(argv, args...)

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

	logFile, err := t.FileSystem.OpenFileWriter(ctx.Value(constants.DeployLogFileKey).(string))

	if err != nil {
		logFile = nil
	} else {
		// #nosec
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
