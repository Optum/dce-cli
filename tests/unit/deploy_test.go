package unit

import (
	"fmt"
	"os"
	"testing"

	"github.com/Optum/dce-cli/configs"
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

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).
		Return(fmt.Errorf("Could not create folders: %s", "bad"))

	err := service.Deploy(&svc.DeployConfig{})

	mockFileSystemer.AssertExpectations(t)

	assert.NotNil(t, err, "expected error calling Deploy()")
}

func TestDeployService_DeclineToCreate(t *testing.T) {
	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	decline := "no"

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("Chdir", originDir).Return()
	mockFileSystemer.On("WriteFile", mock.Anything, mock.Anything)

	mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&decline)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)

	mockTFTemplater.
		On("AddVariable", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("Write", mock.Anything).
		Return(nil)

	err := service.Deploy(&svc.DeployConfig{
		Namespace: "somethingpredictable",
	})

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)
	mockPrompter.AssertExpectations(t)

	assert.NotNil(t, err, "expected use decline")
}

func TestDeployService_InitThrowsError(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	expectedErr := fmt.Errorf("error runninng init")

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("Chdir", originDir).Return()
	mockFileSystemer.On("WriteFile", mock.Anything, mock.Anything)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(expectedErr)

	mockTFTemplater.
		On("AddVariable", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("Write", mock.Anything).
		Return(nil)

	err := service.Deploy(&svc.DeployConfig{
		Namespace: "somethingpredictable",
	})

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
	expectedErr := fmt.Errorf("error runninng apply")
	affirm := "yes"

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("Chdir", originDir).Return()
	mockFileSystemer.On("WriteFile", mock.Anything, mock.Anything)

	mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&affirm)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, []string{}).Return(expectedErr)

	mockTFTemplater.
		On("AddVariable", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("Write", mock.Anything).
		Return(nil)

	err := service.Deploy(&svc.DeployConfig{
		Namespace: "somethingpredictable",
	})

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)
	mockPrompter.AssertExpectations(t)

	assert.NotNil(t, err, "expected error calling Deploy() when Apply() errors")
	assert.Equal(t, "error creating infrastructure: error runninng apply", err.Error())
}

func TestDeployService_FileExists(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	s3bucket := "mys3bucket"
	lambdas := []string{"lambda1", "lambda2"}
	codebuilds := []string{"codebuild1", "codebuild2"}
	affirm := "yes"

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("Chdir", originDir).Return()
	mockFileSystemer.On("WriteFile", mock.Anything, mock.Anything)

	mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&affirm)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("GetOutput", mock.Anything, "artifacts_bucket_name").Return(s3bucket, nil)

	mockFileSystemer.On("ChToTmpDir").Return(doesntMatter, doesntMatter)
	mockGithuber.On("DownloadGithubReleaseAsset", "build_artifacts.zip", mock.Anything).Return(nil)
	mockFileSystemer.On("GetArtifactsDir", mock.Anything).Return(doesntMatter)
	mockFileSystemer.On("Unarchive", "build_artifacts.zip", doesntMatter).Return(nil)
	mockFileSystemer.On("RemoveAll", "build_artifacts.zip")

	mockAwser.On("UploadDirectoryToS3", doesntMatter, s3bucket, "").Return(lambdas, codebuilds)
	mockAwser.On("UpdateLambdasFromS3Assets", lambdas, s3bucket, "somethingpredictable")

	mockTFTemplater.
		On("AddVariable", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("SetModuleSource", mock.Anything)
	mockTFTemplater.
		On("Write", mock.Anything).
		Return(nil)

	err := service.Deploy(&svc.DeployConfig{
		Namespace: "somethingpredictable",
		Location:  "github.com/Optum/dce",
	})

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)
	mockPrompter.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}

func TestDeployService_FileExistsWithOpts(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	s3bucket := "mys3bucket"
	lambdas := []string{"lambda1", "lambda2"}
	codebuilds := []string{"codebuild1", "codebuild2"}
	affirm := "yes"
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

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)
	mockFileSystemer.On("Chdir", originDir).Return()

	mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&affirm)

	mockTerraformer.On("Init", mock.Anything, initOpts).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, applyOpts).Return(nil)
	mockTerraformer.On("GetOutput", mock.Anything, "artifacts_bucket_name").Return(s3bucket, nil)

	mockFileSystemer.On("ChToTmpDir").Return(doesntMatter, doesntMatter)
	mockGithuber.On("DownloadGithubReleaseAsset", "build_artifacts.zip", mock.Anything).Return(nil)
	mockFileSystemer.On("GetArtifactsDir", mock.Anything).Return(doesntMatter)
	mockFileSystemer.On("Unarchive", "build_artifacts.zip", doesntMatter).Return(nil)
	mockFileSystemer.On("RemoveAll", "build_artifacts.zip")
	mockFileSystemer.On("WriteFile", mock.Anything, mock.Anything)

	mockTFTemplater.
		On("AddVariable", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("Write", mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("SetModuleSource", mock.Anything)

	mockAwser.On("UploadDirectoryToS3", doesntMatter, s3bucket, "").Return(lambdas, codebuilds)
	mockAwser.On("UpdateLambdasFromS3Assets", lambdas, s3bucket, "somethingpredictable")

	err := service.Deploy(&svc.DeployConfig{
		Namespace:      "somethingpredictable",
		Location:       "github.com/Optum/dce",
		TFInitOptions:  "-backend-config=\"address=demo.consul.io\" -backend-config=\"path=example_app/terraform_state\"",
		TFApplyOptions: "-compact-warnings     -backup=\"path\"",
	})

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
	s3bucket := "mys3bucket"
	lambdas := []string{"lambda1", "lambda2"}
	codebuilds := []string{"codebuild1", "codebuild2"}
	affirm := "yes"

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)

	// file is being created...
	mockTFTemplater.On("AddVariable", "namespace", "string", "somethingpredictable").Return(nil)
	mockTFTemplater.On("SetModuleSource", mock.Anything)
	mockTFTemplater.On("Write", mock.Anything).Return(nil)
	mockFileSystemer.On("WriteFile", "/file.txt", "").Return()

	// then everything resumes along the happy path as before...
	mockFileSystemer.On("Chdir", originDir).Return()

	mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&affirm)

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("GetOutput", mock.Anything, "artifacts_bucket_name").Return(s3bucket, nil)

	mockFileSystemer.On("ChToTmpDir").Return(doesntMatter, doesntMatter)
	mockGithuber.On("DownloadGithubReleaseAsset", "build_artifacts.zip", mock.Anything).Return(nil)
	mockFileSystemer.On("GetArtifactsDir", mock.Anything).Return(doesntMatter)
	mockFileSystemer.On("Unarchive", "build_artifacts.zip", doesntMatter).Return(nil)
	mockFileSystemer.On("RemoveAll", "build_artifacts.zip")

	mockAwser.On("UploadDirectoryToS3", doesntMatter, s3bucket, "").Return(lambdas, codebuilds)
	mockAwser.On("UpdateLambdasFromS3Assets", lambdas, s3bucket, "somethingpredictable")


	err := service.Deploy(&svc.DeployConfig{
		Namespace: "somethingpredictable",
		Location: "github.com/Optum/dce",
	})

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockTFTemplater.AssertExpectations(t)
	mockPrompter.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}

func TestDeployService_DoesNotFileExistAndUsingLocalRepo(t *testing.T) {

	newDir := "/newdir"
	originDir := "/origindir"
	filename := "/file.txt"
	s3bucket := "mys3bucket"
	lambdas := []string{"lambda1", "lambda2"}
	codebuilds := []string{"codebuild1", "codebuild2"}
	affirm := "yes"

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockFileSystemer.On("CreateConfigDirTree", mock.Anything).Return(nil)
	mockFileSystemer.On("ChToConfigDir").Return(newDir, originDir)
	mockFileSystemer.On("GetLocalMainTFFile").Return(filename)

	// file is being created...
	mockFileSystemer.On("WriteFile", "/file.txt", mock.Anything).Return()

	mockPrompter.On("PromptBasic", mock.Anything, mock.Anything).Return(&affirm)

	// then everything resumes along the happy path as before...
	mockFileSystemer.On("Chdir", originDir).Return()

	mockTerraformer.On("Init", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("Apply", mock.Anything, []string{}).Return(nil)
	mockTerraformer.On("GetOutput", mock.Anything, "artifacts_bucket_name").Return(s3bucket, nil)

	mockFileSystemer.On("ChToTmpDir").Return(doesntMatter, doesntMatter)
	mockFileSystemer.On("GetArtifactsDir", mock.Anything).Return(doesntMatter)
	mockFileSystemer.On("Unarchive", "/local/bin/build_artifacts.zip", doesntMatter).Return(nil)

	mockAwser.On("UploadDirectoryToS3", doesntMatter, s3bucket, "").Return(lambdas, codebuilds)
	mockAwser.On("UpdateLambdasFromS3Assets", lambdas, s3bucket, "somethingpredictable")

	mockTFTemplater.
		On("AddVariable", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("Write", mock.Anything).
		Return(nil)
	mockTFTemplater.
		On("SetModuleSource", "/local/modules")

	err := service.Deploy(&svc.DeployConfig{
		Location:  "/local",
		Namespace: "somethingpredictable",
	})

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
	mockAwser.AssertExpectations(t)
	mockPrompter.AssertExpectations(t)

	assert.Nil(t, err, "expected no error calling Deploy() in happy path")
}

func TestDeployService_PostDeployDefault(t *testing.T) {
	apiURL := "https://some-api-id.execute-api.us-east-1.amazonaws.com/api"

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockTerraformer.On("GetOutput", mock.Anything, "api_url").Return(apiURL, nil)

	mockFileSystemer.On("WriteConfig").Return(nil)

	err := service.PostDeploy(&svc.DeployConfig{
		TFInitOptions:  "",
		TFApplyOptions: "-compact-warnings",
	})
	assert.Nil(t, err)

	assert.Equal(t, *service.Config.API.Host, "some-api-id.execute-api.us-east-1.amazonaws.com")
	assert.Equal(t, *service.Config.API.BasePath, "/api")
	// In default behavior, these should be left alone
	assert.Nil(t, service.Config.Terraform.TFInitOptions, "TFInitOptions should be unset by default")
	assert.Nil(t, service.Config.Terraform.TFApplyOptions, "TFApplyOptions should be unset by default")

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
}

func TestDeployService_PostDeploySaveOpts(t *testing.T) {
	apiURL := "https://some-api-id.execute-api.us-east-1.amazonaws.com/api"

	emptyConfig := configs.Root{}
	initMocks(emptyConfig)

	mockTerraformer.On("GetOutput", mock.Anything, "api_url").Return(apiURL, nil)

	mockFileSystemer.On("WriteConfig").Return(nil)

	err := service.PostDeploy(&svc.DeployConfig{
		SaveTFOptions:  true,
		TFInitOptions:  "",
		TFApplyOptions: "-compact-warnings",
	})
	assert.Nil(t, err)

	assert.Equal(t, *service.Config.API.Host, "some-api-id.execute-api.us-east-1.amazonaws.com")
	assert.Equal(t, *service.Config.API.BasePath, "/api")
	// In default behavior, these should be left alone
	assert.Equal(t, *service.Config.Terraform.TFInitOptions, "", "TFInitOptions should be set because save was specified")
	assert.Equal(t, *service.Config.Terraform.TFApplyOptions, "-compact-warnings", "TFApplyOptions should be set because savw was specified")

	mockFileSystemer.AssertExpectations(t)
	mockTerraformer.AssertExpectations(t)
}

func Test_mockFileInfo_Name(t *testing.T) {
	tests := []struct {
		name string
		m    *mockFileInfo
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.Name(); got != tt.want {
				t.Errorf("mockFileInfo.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mockFileInfo_IsDir(t *testing.T) {
	tests := []struct {
		name string
		m    *mockFileInfo
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.m.IsDir(); got != tt.want {
				t.Errorf("mockFileInfo.IsDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
