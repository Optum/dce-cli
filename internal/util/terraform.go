package util

import (
	"context"
	logger "log"
	"os"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	tfBackendInit "github.com/hashicorp/terraform/backend/init"
	tfCommand "github.com/hashicorp/terraform/command"
	tfDiscovery "github.com/hashicorp/terraform/svchost/disco"

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
	logger.SetOutput(logFile)

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

	logger.SetOutput(logFile)
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
