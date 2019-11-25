package util

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"os"

	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
)

type UtilContainer struct {
	Config *configs.Root
	// File path location of the dce.yaml file, from which this config was parsed
	// Useful if we want to reload or modify the file later
	ConfigFile  string
	Observation *observ.ObservationContainer
	AWSSession *session.Session
	AWSer
	APIer
	Terraformer
	Githuber
	Prompter
	FileSystemer
	Weber
}

var log observ.Logger

// New returns a new Util given config
func New(config *configs.Root, configFile string, observation *observ.ObservationContainer) *UtilContainer {
	log = observation.Logger

	var awsSession *session.Session
	awsSession, err := NewAWSSession(config.API.Token)
	if err != nil {
		log.Fatalf("Failed to initialize AWS Session: %s", err)
	}

	var apiClient APIer
	if config.API.Host != nil && config.API.BasePath != nil {
		apiClient = NewAPIClient(&NewAPIClientInput{
			credentials: awsSession.Config.Credentials,
			region:      config.Region,
			host:        config.API.Host,
			basePath:    config.API.BasePath,
			token:       config.API.Token,
		})
	}


	utilContainer := UtilContainer{
		Config:       config,
		Observation:  observation,
		AWSSession:   awsSession,
		AWSer:        &AWSUtil{Config: config, Observation: observation, Session: awsSession},
		APIer:        apiClient,
		Terraformer:  &TerraformUtil{Config: config, Observation: observation},
		Githuber:     &GithubUtil{Config: config, Observation: observation},
		Prompter:     &PromptUtil{Config: config, Observation: observation},
		FileSystemer: &FileSystemUtil{Config: config, ConfigFile: configFile},
		Weber:        &WebUtil{Observation: observation},
	}

	return &utilContainer
}

type AWSer interface {
	UploadDirectoryToS3(localPath string, bucket string, prefix string) ([]string, []string)
	UpdateLambdasFromS3Assets(lambdaNames []string, bucket string, namespace string)
	ConfigureAWSCLICredentials(accessKeyID, secretAccessKey, sessionToken, profile string)
}

type Terraformer interface {
	Init(args []string)
	Apply(tfVars []string)
	GetOutput(key string) string
}

type Githuber interface {
	DownloadGithubReleaseAsset(assetName string)
}

type Prompter interface {
	PromptBasic(label string, validator func(input string) error) *string
	PromptSelect(label string, items []string) *string
}

type FileSystemer interface {
	WriteConfig() error
	GetConfigFile() string
	GetHomeDir() string
	IsExistingFile(path string) bool
	ReadFromFile(path string) string
	ReadInConfig() error
	Unarchive(source string, destination string)
	MvToTempDir(prefix string) (string, string)
	RemoveAll(path string)
	Chdir(path string)
	ReadDir(path string) []os.FileInfo
	WriteFile(fileName string, data string)
}

type Weber interface {
	OpenURL(url string)
}
