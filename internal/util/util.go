package util

import (
	"github.com/Optum/dce-cli/configs"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

type UtilContainer struct {
	Config *configs.Root
	AWSer
	APIer
	Terraformer
	Githuber
}

// New returns a new Util given config
func New(config *configs.Root, awsCreds *credentials.Credentials) *UtilContainer {
	var session = session.New(&aws.Config{
		Credentials: awsCreds,
		Region:      config.Region,
	})

	return &UtilContainer{
		Config:      config,
		AWSer:       &AWSUtil{Config: config, Session: session},
		APIer:       &APIUtil{Config: config, Session: session},
		Terraformer: &TerraformUtil{Config: config},
		Githuber:    &GithubUtil{Config: config},
	}
}

type AWSer interface {
	UploadDirectoryToS3(localPath string, bucket string, prefix string) []string
	UpdateLambdasFromS3Assets()
}

type APIer interface {
	Request(input *ApiRequestInput) *ApiResponse
}

type Terraformer interface {
	Init(args []string)
	Apply(namespace string)
	GetOutput(key string) string
	GetTemplate(string) string
}

type Githuber interface {
	DownloadGithubReleaseAsset(assetName string)
}
