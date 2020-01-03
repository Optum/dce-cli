package util

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
)

type UtilContainer struct {
	Config *configs.Root
	// File path location of the dce.yaml file, from which this config was parsed
	// Useful if we want to reload or modify the file later
	ConfigFile  string
	Observation *observ.ObservationContainer
	AWSSession  *session.Session
	AWSer
	APIer
	Terraformer
	Githuber
	Prompter
	FileSystemer
	Weber
	Durationer
	TFTemplater
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
	awsSession.Config.Region = config.Region

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

	filesystem := &FileSystemUtil{Config: config, ConfigFile: configFile}
	weber := &WebUtil{Observation: observation}

	utilContainer := UtilContainer{
		Config:      config,
		Observation: observation,
		AWSSession:  awsSession,
		AWSer:       &AWSUtil{Config: config, Observation: observation, Session: awsSession},
		APIer:       apiClient,
		// Terraformer:  &TerraformUtil{Config: config, Observation: observation},
		Terraformer:  &TerraformBinUtil{Config: config, Observation: observation, FileSystem: filesystem, Downloader: weber},
		Githuber:     &GithubUtil{Config: config, Observation: observation},
		Prompter:     &PromptUtil{Config: config, Observation: observation},
		FileSystemer: filesystem,
		Weber:        weber,
		Durationer:   NewDurationUtil(),
	}

	utilContainer.TFTemplater = NewMainTFTemplate(utilContainer.FileSystemer)

	return &utilContainer
}

type AWSer interface {
	UploadDirectoryToS3(localPath string, bucket string, prefix string) ([]string, []string)
	UpdateLambdasFromS3Assets(lambdaNames []string, bucket string, namespace string)
	ConfigureAWSCLICredentials(accessKeyID, secretAccessKey, sessionToken, profile string)
}

type Terraformer interface {
	Init(ctx context.Context, args []string)
	Apply(ctx context.Context, tfVars []string)
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
	GetConfigDir() string
	GetHomeDir() string
	IsExistingFile(path string) bool
	ReadFromFile(path string) string
	ReadInConfig() error
	Unarchive(source string, destination string)
	ChToConfigDir() (string, string)
	RemoveAll(path string)
	Chdir(path string)
	ReadDir(path string) []os.FileInfo
	WriteFile(fileName string, data string)
	OpenFileWriter(path string) (*os.File, error)
}

type Weber interface {
	OpenURL(url string)
}

type Durationer interface {
	ExpandEpochTime(str string) (int64, error)
	ParseDuration(str string) (time.Duration, error)
}

type TFTemplater interface {
	AddVariable(name string, vartype string, vardefault string) error
	Write(w io.Writer) error
}
