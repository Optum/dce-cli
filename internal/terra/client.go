package terra

import (
	"log"
	"os"

	terraform "github.com/hashicorp/terraform/command"
	"github.com/mitchellh/cli"
)

// Init initialized a terraform working directory
func Init(args []string) {
	log.Println("Running terraform init")
	tfInit := &terraform.InitCommand{
		Meta: terraform.Meta{
			Ui: getTerraformUI(),
		},
	}
	tfInit.Run(args)
}

// Apply applies terraform template with given namespace
func Apply(namespace string) {
	log.Println("Running terraform apply with namespace: " + namespace)
	tfApply := &terraform.ApplyCommand{
		Meta: terraform.Meta{
			Ui: getTerraformUI(),
		},
	}
	namespaceKey := "-var"
	namespaceValue := "namespace=" + namespace
	tfApply.Run([]string{namespaceKey, namespaceValue})
}

// GetOutput gets terraform output value for provided key
func GetOutput(key string) string {
	log.Println("Retrieving terraform output for: " + key)
	outputCaptorUI := &UIOutputCaptor{
		BasicUi: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		Captor: new(string),
	}
	tfOutput := &terraform.OutputCommand{
		Meta: terraform.Meta{
			Ui: outputCaptorUI,
		},
	}
	tfOutput.Run([]string{"bucket"})
	return *outputCaptorUI.Captor
}

func getTerraformUI() *cli.PrefixedUi {
	basicUI := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	prefix := "\nTerraform Says...\n"
	return &cli.PrefixedUi{
		AskPrefix:       prefix,
		AskSecretPrefix: prefix,
		OutputPrefix:    prefix,
		InfoPrefix:      prefix,
		ErrorPrefix:     prefix,
		WarnPrefix:      prefix,
		Ui:              basicUI,
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
