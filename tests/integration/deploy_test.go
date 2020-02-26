package integration

import (
	"github.com/Optum/dce-cli/cmd"
	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	"github.com/Optum/dce-cli/internal/util"
	"github.com/Optum/dce-cli/mocks"
	"github.com/Optum/dce-cli/pkg/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSystemDeployCommand(t *testing.T) {

	t.Run("should deploy with default config", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{"system", "deploy"},
			yamlConfig: &configs.Root{},
			expectedMainTf: []string{
				// Should use default version as TF module source
				`
module "dce" {
	source="github.com/Optum/dce//modules?ref=v0.23.0"`,
				`
output "artifacts_bucket_name" {
	description = "S3 bucket for artifacts like AWS Lambda code"
	value = module.dce.artifacts_bucket_name
}`,
			},
			expectedDeployedVersion: constants.DefaultDCEVersion,
		})
	})

	t.Run("should accept a DCE version, as a CLI flag", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{
				"system", "deploy",
				"--dce-version", "9999.12.3",
			},
			expectedDeployedVersion: "9999.12.3",
			// Should use configured version as TF module source
			expectedMainTf: []string{
				`module "dce" {
	source="github.com/Optum/dce//modules?ref=v9999.12.3"`,
			},
		})
	})

	t.Run("should accept a `v` prefix in the version number", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{
				"system", "deploy",
				"--dce-version", "v9999.12.3",
			},
			expectedDeployedVersion: "9999.12.3",
			// Should use configured version as TF module source
			expectedMainTf: []string{
				`module "dce" {
	source="github.com/Optum/dce//modules?ref=v9999.12.3"`,
			},
		})
	})

	t.Run("should accept a DCE version, as a YAML config", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{"system", "deploy"},
			yamlConfig: &configs.Root{
				Deploy: configs.Deploy{
					Version: stringp("v9999.12.3"),
				},
			},
			// Should use configured version as TF module source
			expectedMainTf: []string{
				`module "dce" {
	source="github.com/Optum/dce//modules?ref=v9999.12.3"`,
			},
			expectedDeployedVersion: "9999.12.3",
		})
	})

	t.Run("should use CLI flag for version, override YAML config", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{
				"system", "deploy",
				"--dce-version", "v9999.12.3",
			},
			yamlConfig: &configs.Root{
				Deploy: configs.Deploy{
					Version: stringp("v0.12.3"),
				},
			},
			// Should use CLI flag version as TF module source
			expectedMainTf: []string{
				`module "dce" {
	source="github.com/Optum/dce//modules?ref=v9999.12.3"`,
			},
			expectedDeployedVersion: "9999.12.3",
		})
	})

	t.Run("should accept tfvars as CLI flags", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{
				"system", "deploy",
				"--region", "moon-darkside-1",
				"--namespace", "my-namespace",
				"--budget-notification-from-email", "from@example.com",
				"--dce-version", "v0.12.3",
			},
			expectedMainTf: []string{
				`
variable "aws_region" {
	type = string
	default = "moon-darkside-1"
}`,
				`
variable "namespace" {
	type = string
	default = "my-namespace"
}`,
				`
variable "budget_notification_from_email" {
	type = string
	default = "from@example.com"
}`,
				`
module "dce" {
	source="github.com/Optum/dce//modules?ref=v0.12.3"

	aws_region = var.aws_region
	namespace = var.namespace
	budget_notification_from_email = var.budget_notification_from_email
}`,
			},
			expectedDeployedVersion: "0.12.3",
		})
	})

	t.Run("should accept tfvars as YAML config", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{"system", "deploy"},
			yamlConfig: &configs.Root{
				Deploy: configs.Deploy{
					Version:                     stringp("v0.12.3"),
					AWSRegion:                   stringp("moon-darkside-1"),
					Namespace:                   stringp("my-namespace"),
					BudgetNotificationFromEmail: stringp("from@example.com"),
				},
			},
			expectedMainTf: []string{
				`
variable "aws_region" {
	type = string
	default = "moon-darkside-1"
}`,
				`
variable "namespace" {
	type = string
	default = "my-namespace"
}`,
				`
variable "budget_notification_from_email" {
	type = string
	default = "from@example.com"
}`,
				`
module "dce" {
	source="github.com/Optum/dce//modules?ref=v0.12.3"

	aws_region = var.aws_region
	namespace = var.namespace
	budget_notification_from_email = var.budget_notification_from_email
}`,
			},
			expectedDeployedVersion: "0.12.3",
		})
	})

	t.Run("should accept env var configuration", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand: []string{"system", "deploy"},
			envVars: map[string]string{
				"DCE_VERSION":                        "v9999.12.3",
				"AWS_REGION":                         "moon-darkside-1",
				"DCE_NAMESPACE":                      "my-namespace",
				"DCE_BUDGET_NOTIFICATION_FROM_EMAIL": "from@example.com",
			},
			expectedMainTf: []string{
				`
variable "aws_region" {
	type = string
	default = "moon-darkside-1"
}`,
				`
variable "namespace" {
	type = string
	default = "my-namespace"
}`,
				`
variable "budget_notification_from_email" {
	type = string
	default = "from@example.com"
}`,
				`
module "dce" {
	source="github.com/Optum/dce//modules?ref=v9999.12.3"

	aws_region = var.aws_region
	namespace = var.namespace
	budget_notification_from_email = var.budget_notification_from_email
}`,
			},
			expectedDeployedVersion: "9999.12.3",
		})
	})

	t.Run("should not prompt for approval, if batch mode is disabled", func(t *testing.T) {
		deployCommandTestCase(t, &deployTestCase{
			cliCommand:             []string{"system", "deploy", "--batch-mode"},
			expectDeploymentPrompt: boolp(false),
		})
	})
}

// deployTest is a wrapper for common mocks and utilities
// used by `dce system deploy` integration tests
type deployTest struct {
	*cliTest
	terraform *mocks.Terraformer
	github    *stubGithub
	aws       *mocks.AWSer
	configDir string
	cleanup   func()
}

func (test *deployTest) readConfigFile(t *testing.T, paths ...string) string {
	paths = append([]string{test.configDir}, paths...)
	fullPath := filepath.Join(paths...)
	fileBytes, err := ioutil.ReadFile(fullPath)
	require.Nil(t, err)
	return string(fileBytes)
}

// newDeployTest sets up baseline mocks and configuration
// for a system deploy command integration test
func newDeployTest(t *testing.T, config *configs.Root) *deployTest {
	// Reset the DeployConfig struct before and after each test run
	// (global var, populated with CLI flag/Yaml config vals)
	copyStructVals(t, service.DeployConfig{}, cmd.DeployConfig)
	defer copyStructVals(t, service.DeployConfig{}, cmd.DeployConfig)

	cli := NewCLITest(t)

	cli.WriteConfig(t, config)

	// Mock the FileSystemer, to use a tmp dir for
	// generated terraform files
	configDir, err := ioutil.TempDir("", ".dce-test-")
	require.Nil(t, err)

	// Mock our Terraform wrapper util
	terraform := &mocks.Terraformer{}

	// Mock github (for release downloads)
	github := &stubGithub{}

	// Mock the AWS util
	aws := &mocks.AWSer{}

	// Mock the Authentication service (would pop open browser to auth user)
	authSvc := &mocks.Authenticater{}
	authSvc.On("Authenticate").Return(nil)

	// Inject mocks as globals used by the CLI
	cli.Inject(func(input *injectorInput) {
		input.service.Util.Terraformer = terraform
		fsUtil := input.service.Util.FileSystemer.(*util.FileSystemUtil)
		fsUtil.ConfigDir = configDir
		input.service.Util.Githuber = github
		input.service.Util.AWSer = aws
		input.service.Authenticater = authSvc
	})

	return &deployTest{
		cliTest:   cli,
		terraform: terraform,
		github:    github,
		aws:       aws,
		configDir: configDir,
		cleanup: func() {
			err := os.RemoveAll(configDir)
			if err != nil {
				log.Printf("WARNING: Failed to remove temporary test directory at %s: %s", configDir, err)
			}
		},
	}
}

type deployTestCase struct {
	cliCommand []string
	yamlConfig *configs.Root
	// List of expected strings within our generated main.tf file
	expectedMainTf          []string
	expectedDeployedVersion string
	// Should the CLI prompt you before deploying terraform
	expectDeploymentPrompt *bool
	deployPromptAnswer     string
	envVars                map[string]string
}

// deployCommandTestCase runs a configurable integration test
// using the `dce system deploy` command.
func deployCommandTestCase(t *testing.T, input *deployTestCase) {
	// Setup test case defaults
	if input.expectedDeployedVersion == "" {
		input.expectedDeployedVersion = constants.DefaultDCEVersion
	}
	if input.deployPromptAnswer == "" {
		input.deployPromptAnswer = "yes"
	}
	if input.expectDeploymentPrompt == nil {
		input.expectDeploymentPrompt = boolp(true)
	}

	// mock env vars
	for key, value := range input.envVars {
		revert := mockEnvVar(key, value)
		defer revert()
	}

	// Setup test mocks and temp file system,
	// using the provided yaml config
	test := newDeployTest(t, input.yamlConfig)
	defer test.cleanup()

	// Should run terraform init
	test.terraform.
		On("Init", mock.Anything, []string{}).
		Run(func(args mock.Arguments) {
			// Check that our main.tf is generated, before we call terraform init
			mainTF := test.readConfigFile(t, ".cache", "module", "main.tf")
			for _, part := range input.expectedMainTf {
				require.Contains(t, mainTF, strings.TrimSpace(part))
			}
		}).
		Return(nil)

	// Should run terraform apply
	test.terraform.
		On("Apply", mock.Anything, []string{}).
		Return(nil)

	// Mock the deployed artifacts_bucket output from Terraform
	test.terraform.
		On("GetOutput", mock.Anything, "artifacts_bucket_name").
		Return("test-artifacts-bucket", nil)
	test.terraform.
		On("GetOutput", mock.Anything, "api_url").
		Return("dce.example.com/api", nil)

	// Mock github release assets for the default DCE version
	test.github.MockReleaseAsset(t,
		"build_artifacts.zip",
		input.expectedDeployedVersion,
		zipFiles(t, []file{{"/lambda/leases.zip", "stub leases lambda artifact"}}),
	)

	// Should upload build artifacts to S3
	test.aws.
		On("UploadDirectoryToS3",
			// upload artifact dir
			filepath.Join(test.configDir, ".cache", "dce", input.expectedDeployedVersion),
			"test-artifacts-bucket",
			"",
		).
		// Verify that our unzipped release asset are in the artifacts dir
		Run(func(args mock.Arguments) {
			// eg. ~/.dce/.cache/dce/v0.12.3/lambda/leases.zip
			leasesLambdaArtifact := test.readConfigFile(t,
				".cache", "dce", input.expectedDeployedVersion,
				"lambda", "leases.zip")
			assert.Equal(t, "stub leases lambda artifact", leasesLambdaArtifact,
				"content of uploaded lambda artifact (stubbed from Github release download)")
		}).
		Return([]string{"leases-lambda-name"}, []string{})

	// Should update lambdas
	test.aws.
		On("UpdateLambdasFromS3Assets",
			[]string{"leases-lambda-name"},
			"test-artifacts-bucket",
			mock.Anything, // random namspace
		)

	// Approve the deployment
	if *input.expectDeploymentPrompt {
		test.AnswerBasic(
			"Do you really want to create DCE resources in your AWS account? "+
				"(type \"yes\" or \"no\")",
			input.deployPromptAnswer,
		)
	}

	// Run the `dce system deploy` command
	err := test.Execute(input.cliCommand)
	require.Nil(t, err)

	// Verify that we ran terraform
	test.terraform.AssertExpectations(t)
	// Verify that we deployed to AWS
	test.aws.AssertExpectations(t)
}
