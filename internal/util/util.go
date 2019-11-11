package util

import (
	"os"

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
	Weber
}

var log observ.Logger

// New returns a new Util given config
func New(config *configs.Root, observation *observ.ObservationContainer) *UtilContainer {

	log = observation.Logger

	var masterAcctCreds *credentials.Credentials
	if config.System.MasterAccount.Credentials.AwsAccessKeyID != nil &&
		config.System.MasterAccount.Credentials.AwsSecretAccessKey != nil {
		masterAcctCreds = credentials.NewStaticCredentials(
			*config.System.MasterAccount.Credentials.AwsAccessKeyID,
			*config.System.MasterAccount.Credentials.AwsSecretAccessKey,
			"",
		)
	}
	var session = session.New(&aws.Config{
		Credentials: masterAcctCreds,
		Region:      config.Region,
	})

	var initalizedApiClient APIer
	if config.System.MasterAccount.Credentials.AwsAccessKeyID != nil {
		apiUtil := &APIUtil{Config: config, Observation: observation, Session: session}
		initalizedApiClient = apiUtil.InitApiClient()
	}

	utilContainer := UtilContainer{
		Config:       config,
		Observation:  observation,
		AWSer:        &AWSUtil{Config: config, Observation: observation, Session: session},
		APIer:        initalizedApiClient,
		Terraformer:  &TerraformUtil{Config: config, Observation: observation},
		Githuber:     &GithubUtil{Config: config, Observation: observation},
		Prompter:     &PromptUtil{Config: config, Observation: observation},
		FileSystemer: &FileSystemUtil{Config: config, Observation: observation, DefaultConfigFileName: constants.DefaultConfigFileName},
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
	WriteToYAMLFile(path string, _struct interface{})
	GetDefaultConfigFile() string
	GetHomeDir() string
	IsExistingFile(path string) bool
	ReadFromFile(path string) string
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
