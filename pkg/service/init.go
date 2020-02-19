package service

import (
	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/aws/aws-sdk-go/aws"
)

type InitService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

func (s *InitService) InitializeDCE() {
	// Set default region
	if s.Config.Region == nil {
		s.Config.Region = aws.String("us-east-1")
	}

	// Prompt user for required configs
	s.promptUserForConfig(s.Config)

	// Write the config to dce.yml
	err := s.Util.WriteConfig()
	if err != nil {
		log.Fatalf("Failed to write YAML config to %s: %s",
			s.Util.GetConfigFile(), err)
	}

	log.Infoln("Config file created at: " + s.Util.GetConfigFile())
}

func (s *InitService) promptUserForConfig(config *configs.Root) {
	// API Config
	if config.API.Host == nil {
		config.API.Host = s.Util.PromptBasic("Host name of the DCE API (example: abcde12345.execute-api.us-east-1.amazonaws.com)", nil)
	}
	if config.API.BasePath == nil {
		config.API.BasePath = s.Util.PromptBasic("Base path of the DCE API (example: /api)", nil)
	}
}
