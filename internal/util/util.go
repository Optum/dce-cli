package util

import (
	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type UtilContainer struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	AWSer
	APIer
	Terraformer
	Githuber
	Prompter
	FileSystemer
}

var Log observ.Logger

// New returns a new Util given config
func New(config *configs.Root, observation *observ.ObservationContainer, awsCreds *credentials.Credentials) *UtilContainer {

	Log = observation.Logger

	var session = session.New(&aws.Config{
		Credentials: awsCreds,
		Region:      config.Region,
	})

	return &UtilContainer{
		Config:       config,
		Observation:  observation,
		AWSer:        &AWSUtil{Config: config, Session: session},
		APIer:        &APIUtil{Config: config, Session: session},
		Terraformer:  &TerraformUtil{Config: config},
		Githuber:     &GithubUtil{Config: config},
		Prompter:     &PromptUtil{Config: config},
		FileSystemer: &FileSystemUtil{Config: config, DefaultConfigFileName: constants.DefaultConfigFileName},
	}
}

type AWSer interface {
	UploadDirectoryToS3(localPath string, bucket string, prefix string) ([]string, []string)
	UpdateLambdasFromS3Assets(lambdaNames []string, bucket string, namespace string)
}

type APIer interface {
	Request(input *ApiRequestInput) *ApiResponse
}

type Terraformer interface {
	Init(args []string)
	Apply(namespace string)
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
	WriteToYAMLFile(path string, _struct interface{})
	GetDefaultConfigFile() string
	GetHomeDir() string
	IsExistingFile(path string) bool
}

type Logger interface {
}
