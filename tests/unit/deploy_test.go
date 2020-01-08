package unit

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Optum/dce-cli/configs"
	cfg "github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	svc "github.com/Optum/dce-cli/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var doesntMatter = "doesntmatter"

type mockFileInfo struct {
	os.FileInfo
}

func (m *mockFileInfo) Name() string { return doesntMatter }
func (m *mockFileInfo) IsDir() bool  { return true }

func TestDeployService_ErrorCreatingDirs(t *testing.T) {

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree").Return(fmt.Errorf("Could not create folders: %s", "bad"))

	deployConfig := cfg.DeployConfig{}
	overrides := svc.DeployOverrides{}
	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.Deploy(ctx, &overrides)

	mockFileSystemer.AssertExpectations(t)

	assert.NotNil(t, err, "expected error calling Deploy()")
}

func TestDeployService_FileExists(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	logfile := "/log.txt"
	s3bucket := "mys3bucket"
	lambdas := []string{"lambda1", "lambda2"}
	codebuilds := []string{"codebuild1", "codebuild2"}

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree").Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalBackendFile").Return(filename)
	mockFileSystemer.On("IsExistingFile", filename).Return(false)
	mockFileSystemer.On("Chdir", originDir).Return()

	mockFileSystemer.On("GetLogFile").Return(logfile)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("GetOutput", mock.Anything, "artifacts_bucket_name").Return(s3bucket, nil)

	mockFileSystemer.On("ChToTmpDir").Return(doesntMatter, doesntMatter)
	mockGithuber.On("DownloadGithubReleaseAsset", "build_artifacts.zip")
	mockFileSystemer.On("GetArtifactsDir").Return(doesntMatter)
	mockFileSystemer.On("Unarchive", "build_artifacts.zip", doesntMatter)
	mockFileSystemer.On("RemoveAll", "build_artifacts.zip")

	mockAwser.On("UploadDirectoryToS3", doesntMatter, s3bucket, "").Return(lambdas, codebuilds)
	mockAwser.On("UpdateLambdasFromS3Assets", lambdas, s3bucket, "somethingpredictable")

	deployConfig := cfg.DeployConfig{}
	overrides := svc.DeployOverrides{
		Namespace: "somethingpredictable",
	}
	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.Deploy(ctx, &overrides)

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}
