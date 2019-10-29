package service

import (
	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"gopkg.in/yaml.v2"
)

type InitService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *InitService) InitializeDCE(cfgFile string) {
	if cfgFile == "" {
		cfgFile = s.Util.GetDefaultConfigFile()
	}

	config := s.promptUserForConfig()
	yamlConfig, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalln(err)
	}
	log.Infoln("You have entered the following configuration:\n" + string(yamlConfig))
	if *s.Util.PromptBasic(constants.PromptChangeConfigConfirmation, nil) != "yes" {
		log.Endln("Aborting")
	}

	s.Util.WriteToYAMLFile(cfgFile, config)

	log.Infoln("Config file created at: " + cfgFile)
}

func (s *InitService) promptUserForConfig() *configs.Root {
	newConfig := configs.Root{}

	// System Config
	newConfig.System.Auth.LoginURL = s.Util.PromptBasic("Authentication URL (SSO)", nil)
	newConfig.System.MasterAccount.Credentials.AwsAccessKeyID = s.Util.PromptBasic("AWS ACCESS KEY ID for the DCE Master account", nil)
	newConfig.System.MasterAccount.Credentials.AwsSecretAccessKey = s.Util.PromptBasic("AWS SECRET ACCESS KEY for the DCE Master account", nil)

	// API Config
	newConfig.Region = s.Util.PromptSelect("Region is DCE deployed in", configs.Regions)
	newConfig.API.Host = s.Util.PromptBasic("Host name of the DCE API (example: abcde12345.execute-api.us-east-1.amazonaws.com)", nil)
	newConfig.API.BasePath = s.Util.PromptBasic("Base path of the DCE API (example: /apigw-stage-name)", nil)
	newConfig.API.Credentials.AwsAccessKeyID = s.Util.PromptBasic("AWS ACCESS KEY ID for accessing the DCE API (This is usually obtained by running DCE auth. Leave blank to use AWS_ACCESS_KEY_ID env variable.)", nil)
	newConfig.API.Credentials.AwsSecretAccessKey = s.Util.PromptBasic("AWS SECRET ACCESS KEY for accessing the DCE API (This is usually obtained by running DCE auth. Leave blank to use AWS_SECRET_ACCESS_KEY env variable.)", nil)
	newConfig.API.Credentials.AwsSessionToken = s.Util.PromptBasic("AWS SESSION TOKEN for accessing the DCE API (This is usually obtained by running DCE auth. Leave blank to use AWS_SESSION_TOKEN env variable.)", nil)

	newConfig.GithubToken = s.Util.PromptBasic("Github token used to download releases from github (Leave blank to use GITHUB_TOKEN env variable)", nil)
	return &newConfig
}
