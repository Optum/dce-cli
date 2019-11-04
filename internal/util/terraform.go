package util

import (
	"os"

	"github.com/Optum/dce-cli/configs"
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
func (u *TerraformUtil) Init(args []string) {
	log.Println("Running terraform init")

	services := tfDiscovery.NewWithCredentialsSource(nil)
	tfBackendInit.Init(services)

	tfInit := &tfCommand.InitCommand{
		Meta: tfCommand.Meta{
			Ui: getTerraformUI(),
		},
	}
	tfInit.Run(args)
}

// Apply applies terraform template with given namespace
func (u *TerraformUtil) Apply(tfVars []string) {
	tfApply := &tfCommand.ApplyCommand{
		Meta: tfCommand.Meta{
			Ui: getTerraformUI(),
		},
	}

	runArgs := []string{}
	for _, tfVar := range tfVars {
		runArgs = append(runArgs, "-var", tfVar)
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

func getTerraformUI() *cli.BasicUi {
	return &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
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
