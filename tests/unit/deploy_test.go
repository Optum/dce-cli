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

func TestDeployService_InitThrowsError(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	logfile := "/log.txt"
	expectedErr := fmt.Errorf("error runninng init")

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree").Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("IsExistingFile", filename).Return(true)
	mockFileSystemer.On("Chdir", originDir).Return()

	mockFileSystemer.On("GetLogFile").Return(logfile)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(expectedErr)

	deployConfig := cfg.DeployConfig{
		UseCached: true,
	}
	overrides := svc.DeployOverrides{
		Namespace: "somethingpredictable",
	}
	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.Deploy(ctx, &overrides)

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)

	assert.NotNil(t, err, "expected error calling Deploy() when Init() errs")
	assert.Equal(t, "error creating infrastructure: error runninng init", err.Error())
}

func TestDeployService_ApplyThrowsError(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	logfile := "/log.txt"
	expectedErr := fmt.Errorf("error runninng apply")

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree").Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("IsExistingFile", filename).Return(true)
	mockFileSystemer.On("Chdir", originDir).Return()

	mockFileSystemer.On("GetLogFile").Return(logfile)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, []string{}).Return(expectedErr)

	deployConfig := cfg.DeployConfig{
		UseCached: true,
	}
	overrides := svc.DeployOverrides{
		Namespace: "somethingpredictable",
	}
	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.Deploy(ctx, &overrides)

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)

	assert.NotNil(t, err, "expected error calling Deploy() when Apply() errors")
	assert.Equal(t, "error creating infrastructure: error runninng apply", err.Error())
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
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("IsExistingFile", filename).Return(true)
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

	deployConfig := cfg.DeployConfig{
		UseCached: true,
	}
	overrides := svc.DeployOverrides{
		Namespace: "somethingpredictable",
	}
	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.Deploy(ctx, &overrides)

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}

func TestDeployService_FileExistsWithOpts(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	logfile := "/log.txt"
	s3bucket := "mys3bucket"
	lambdas := []string{"lambda1", "lambda2"}
	codebuilds := []string{"codebuild1", "codebuild2"}
	initOpts := []string{
		"-backend-config=\"address=demo.consul.io\"",
		"-backend-config=\"path=example_app/terraform_state\"",
	}
	applyOpts := []string{
		"-compact-warnings",
		"-backup=\"path\"",
	}

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree").Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("IsExistingFile", filename).Return(true)
	mockFileSystemer.On("Chdir", originDir).Return()

	mockFileSystemer.On("GetLogFile").Return(logfile)

	mockTerraformer.On("Init", mock.Anything, initOpts).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, applyOpts).Return(nil)
	mockTerraformer.On("GetOutput", mock.Anything, "artifacts_bucket_name").Return(s3bucket, nil)

	mockFileSystemer.On("ChToTmpDir").Return(doesntMatter, doesntMatter)
	mockGithuber.On("DownloadGithubReleaseAsset", "build_artifacts.zip")
	mockFileSystemer.On("GetArtifactsDir").Return(doesntMatter)
	mockFileSystemer.On("Unarchive", "build_artifacts.zip", doesntMatter)
	mockFileSystemer.On("RemoveAll", "build_artifacts.zip")

	mockAwser.On("UploadDirectoryToS3", doesntMatter, s3bucket, "").Return(lambdas, codebuilds)
	mockAwser.On("UpdateLambdasFromS3Assets", lambdas, s3bucket, "somethingpredictable")

	deployConfig := cfg.DeployConfig{
		UseCached:      true,
		TFInitOptions:  "-backend-config=\"address=demo.consul.io\" -backend-config=\"path=example_app/terraform_state\"",
		TFApplyOptions: "-compact-warnings     -backup=\"path\"",
	}
	overrides := svc.DeployOverrides{
		Namespace: "somethingpredictable",
	}
	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.Deploy(ctx, &overrides)

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}

func TestDeployService_DoesNotFileExist(t *testing.T) {

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
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("IsExistingFile", filename).Return(false)

	// file is being created...
	mockTFTemplater.On("AddVariable", "namespace", "string", "somethingpredictable").Return(nil)
	mockTFTemplater.On("AddVariable", "budget_notification_from_email", "string", "no-reply@example.com").Return(nil)
	mockTFTemplater.On("Write", mock.Anything).Return(nil)
	mockFileSystemer.On("WriteFile", "/file.txt", "").Return()

	// then everything resumes along the happy path as before...
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
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}

func TestDeployService_DoesNotFileExistAndUsingLocalRepo(t *testing.T) {

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
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("IsExistingFile", filename).Return(false)

	// file is being created...
	mockFileSystemer.On("ReadFromFile", mock.Anything).Return("filecontents")
	mockFileSystemer.On("WriteFile", "/file.txt", "filecontents").Return()

	// then everything resumes along the happy path as before...
	mockFileSystemer.On("Chdir", originDir).Return()

	mockFileSystemer.On("Unarchive", "/local/bin/terraform_artifacts.zip", mock.Anything)
	mockFileSystemer.On("GetLogFile").Return(logfile)
	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("GetOutput", mock.Anything, "artifacts_bucket_name").Return(s3bucket, nil)

	mockFileSystemer.On("ChToTmpDir").Return(doesntMatter, doesntMatter)
	mockFileSystemer.On("GetArtifactsDir").Return(doesntMatter)
	mockFileSystemer.On("Unarchive", "/local/bin/build_artifacts.zip", doesntMatter)

	mockAwser.On("UploadDirectoryToS3", doesntMatter, s3bucket, "").Return(lambdas, codebuilds)
	mockAwser.On("UpdateLambdasFromS3Assets", lambdas, s3bucket, "somethingpredictable")

	deployConfig := cfg.DeployConfig{
		DeployLocalPath: "/local",
	}
	overrides := svc.DeployOverrides{
		Namespace: "somethingpredictable",
	}
	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.Deploy(ctx, &overrides)

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}

func TestDeployService_PostDeploy(t *testing.T) {
	emptyConfig := configs.Root{}
	initMocks(emptyConfig)
	deployConfig := cfg.DeployConfig{
		TFInitOptions:  "",
		TFApplyOptions: "-compact-warnings",
	}

	mockFileSystemer.On("WriteConfig").Return(nil)

	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.PostDeploy(ctx)
	assert.Nil(t, err)

	mockFileSystemer.AssertExpectations(t)

}

func TestDeployService_PostDeployEqualValues(t *testing.T) {
	expectedOption := "-compact-warnings"
	empty := ""
	emptyConfig := configs.Root{
		Terraform: configs.Terraform{
			TFInitOptions:  &empty,
			TFApplyOptions: &expectedOption,
		},
	}
	initMocks(emptyConfig)
	deployConfig := cfg.DeployConfig{
		TFInitOptions:  "",
		TFApplyOptions: expectedOption,
	}

	ctx := context.WithValue(context.Background(), constants.DeployConfig, &deployConfig)

	err := service.PostDeploy(ctx)
	assert.Nil(t, err)

	mockFileSystemer.AssertExpectations(t)

}
