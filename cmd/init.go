package cmd

import (
	"io/ioutil"
	"log"

	"github.com/Optum/dce-cli/configs"
	"github.com/manifoldco/promptui"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var defaultConfigFileName string = ".dce.yaml"

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "First time DCE cli setup. Creates config file at ~/.dce.yaml",
	Run: func(cmd *cobra.Command, args []string) {

		if cfgFile == "" {
			cfgFile = getDefaultConfigFile()
		}

		config := promptUserForConfig()

		writeNewConfigFile(cfgFile, config)

		log.Println("Config file created at: " + cfgFile)
	},
}

func promptUserForConfig() *configs.Root {
	config := configs.Root{}

	// System Config
	config.System.Auth.LoginURL = promptBasic("Authentication URL (SSO)", nil)
	config.System.MasterAccount.Credentials.AwsAccessKeyID = promptBasic("AWS ACCESS KEY ID of the DCE Master account (Leave blank if you are not a system admin)", nil)
	config.System.MasterAccount.Credentials.AwsSecretAccessKey = promptBasic("AWS SECRET ACCESS KEY of the DCE Master account (Leave blank if you are not a system admin)", nil)
	config.System.MasterAccount.Credentials.AwsSessionToken = promptBasic("AWS SESSION TOKEN of the DCE Master account (Leave blank if you are not a system admin)", nil)

	// API Config
	config.API.BaseURL = promptBasic("What is the base url of the DCE API (example: https://abcde12345.execute-api.us-east-1.amazonaws.com/dev)?", nil)
	config.API.Region = promptSelect("What region is DCE deployed in?", configs.Regions)

	return &config
}

func promptBasic(label string, validator func(input string) error) *string {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: validator,
	}
	input, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return &input
}

func promptSelect(label string, items []string) *string {
	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, input, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	return &input
}

func writeNewConfigFile(cfgFile string, config *configs.Root) {

	cfgYaml, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile(cfgFile, cfgYaml, 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func getDefaultConfigFile() string {
	parentDir, err := homedir.Dir()
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	return parentDir + "/" + defaultConfigFileName
}
