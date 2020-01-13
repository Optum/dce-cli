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
		Config:       config,
		Observation:  observation,
		AWSSession:   awsSession,
		AWSer:        &AWSUtil{Config: config, Observation: observation, Session: awsSession},
		APIer:        apiClient,
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
	Init(ctx context.Context, args []string) error
	Apply(ctx context.Context, tfVars []string) error
	GetOutput(ctx context.Context, key string) (string, error)
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
	IsExistingFile(path string) bool
	ReadFromFile(path string) string
	ReadInConfig() error
	Unarchive(source string, destination string)
	ChToConfigDir() (string, string)
	ChToTmpDir() (string, string)
	RemoveAll(path string)
	Chdir(path string)
	ReadDir(path string) []os.FileInfo
	WriteFile(fileName string, data string)
	OpenFileWriter(path string) (*os.File, error)

	// GetHomeDir returns the user home dir. For example, on *nix systems this
	// be the same as `~` expanded, or the value of `$HOME`
	GetHomeDir() string
	// GetConfigDir returns the DCE configuration dir, which on *nix systems
	// is `~/.dce`
	GetConfigDir() string
	// GetCacheDir returns the local cache dir, which bt default is `~/.dce/.cache`
	GetCacheDir() string
	// GetArtifactsDir returns the cached artifacts dir, which by default is
	// `~/.dce/.cache/dce/${DCE_VERSION}/`
	GetArtifactsDir() string
	// GetTerraformBinDir returns the dir in which the `terraform` bin is installed,
	// which by default is `~/.dce/.cache/terraform/${TERRAFORM_VERSION}`
	GetTerraformBinDir() string
	// GetLocalBackendDir returns the dir for the local terraform backend.
	// By default, `~/.dce/.cache/module`
	GetLocalTFModuleDir() string
	// CreateConfigDirTree creates all the dirs in the dir specified by GetConfigDir(),
	// including the dir itself.
	CreateConfigDirTree() error

	// GetConfigFile returns the full path of the configuration file, such as
	// `~/.dce/config.yaml`
	GetConfigFile() string
	// GetLogFile returns the full path of the log file for the deployment messages.
	GetLogFile() string
	// GetLocalBackendFile returns the full path of the local backend file.
	GetLocalMainTFFile() string
	// GetTerraformBin returns the full path of the terraform binary.
	GetTerraformBin() string
	// GetTerraformStateFile returns the full path of the terraform state file
	GetTerraformStateFile() string
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
