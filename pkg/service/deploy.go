package service

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/Optum/dce-cli/internal/util"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/pkg/errors"
)

var affirmAnswerRegex *regexp.Regexp

func init() {
	affirmAnswerRegex = regexp.MustCompile(`^([Yy]([Ee][Ss])?)|([Nn][Oo]?)$`)
}

const AssetsFileName = "build_artifacts.zip"

type DeployService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
}

type DeployConfig struct {
	// Path to the DCE repo to deploy
	// May be a local file path (eg. /path/to/dce)
	// or a github repo (eg. github.com/Optum/dce)
	Location string
	// Version of DCE to deploy, eg 0.12.3
	Version string
	// BatchMode, if enabled, forces DCE to run non-interactively
	// and supplies Terraform with -auto-approve and input=false
	BatchMode bool
	// TFInitOptions are options passed along to `terraform init`
	TFInitOptions string
	// TFApplyOptions are options passed along to `terraform apply`
	TFApplyOptions string
	// SaveTFOptions, if yes, will save the provided terraform options to the config file.
	SaveTFOptions bool
	// File location for deployment logs
	DeployLogFile string

	// Terraform variables to pass through to the DCE module
	AWSRegion                         string
	GlobalTags                        []string
	Namespace                         string
	BudgetNotificationFromEmail       string
	BudgetNotificationBCCEmails       []string
	BudgetNotificationTemplateHTML    string
	BudgetNotificationTemplateText    string
	BudgetNotificationTemplateSubject string
}

// Deploy writes the local `main.tf` file, using the overrides, and then
// calls Terraform init and apply using configuration directory (`~/.dce`)
// as the working folder and location of local state.
func (s *DeployService) Deploy(deployConfig *DeployConfig) error {

	// Initialize the folder structure
	if err := s.Util.CreateConfigDirTree(deployConfig.Version); err != nil {
		return errors.Wrap(err, "error creating directory structure")
	}

	// Generate the main.tf file, used to load in our DCE terraform submodule
	_, err := s.createTFMainFile(deployConfig)

	if err != nil {
		return errors.Wrap(err, "error creating local backend")
	}

	// Deploy the DCE terraform module
	artifactsBucket, err := s.createDceInfra(deployConfig)
	if err != nil {
		return errors.Wrap(err, "error creating infrastructure")
	}

	log.Infoln("Artifacts bucket = ", artifactsBucket)

	// Deploy application code
	log.Infoln("Deploying code assets to DCE infrastructure")
	err = s.deployCodeAssets(artifactsBucket, deployConfig.Namespace, deployConfig.Location, deployConfig.Version)

	return err
}

// PostDeploy is intended to run after a successful call to Deploy()
func (s *DeployService) PostDeploy(deployConfig *DeployConfig) error {
	ctx := context.WithValue(context.Background(), constants.DeployLogFileKey, deployConfig.DeployLogFile)

	apiURL, err := s.Util.GetOutput(ctx, "api_url")

	if err == nil {
		parsedURL, err := url.Parse(apiURL)
		if err == nil {
			hostName := parsedURL.Hostname()
			basePath := parsedURL.EscapedPath()
			s.Config.API.Host = &hostName
			s.Config.API.BasePath = &basePath
		}
	}

	// if the user wants to save them, update them here otherwise
	// leave tham alone..
	if deployConfig.SaveTFOptions {
		// Update the global config object with our local configuration
		s.Config.Terraform.TFInitOptions = &deployConfig.TFInitOptions
		s.Config.Terraform.TFApplyOptions = &deployConfig.TFApplyOptions
		s.Config.Deploy.Location = &deployConfig.Location
		s.Config.Deploy.Version = &deployConfig.Version
		s.Config.Deploy.LogFile = &deployConfig.DeployLogFile
	}

	return s.Util.WriteConfig()
}

// createTFMainFile creates the main.tf file. The default behavior, without
// using a local repository, is to create the file with the bare minimum
// required to use Terraform with local state.
func (s *DeployService) createTFMainFile(deployConfig *DeployConfig) (string, error) {
	_, originDir := s.Util.ChToConfigDir()
	defer s.Util.Chdir(originDir)

	fileName := s.Util.GetLocalMainTFFile()


	tfMainContents, err := s.getLocalTFMainContents(deployConfig)
	if err != nil {
		return "", err
	}
	s.Util.WriteFile(fileName, tfMainContents)
	return fileName, nil
}

func (s *DeployService) createDceInfra(deployConfig *DeployConfig) (string, error) {
	_, originDir := s.Util.ChToConfigDir()
	defer s.Util.Chdir(originDir)

	deployLogFileName := deployConfig.DeployLogFile

	// Setup context used by Terraform service
	ctx := context.WithValue(context.Background(), constants.DeployLogFileKey, deployLogFileName)

	log.Infoln("Initializing")
	initopts, _ := util.ParseOptions(&deployConfig.TFInitOptions)
	if err := s.Util.Terraformer.Init(ctx, initopts); err != nil {
		return "", err
	}


	// First, prompt the user to give them a chance to opt out in case
	// of accidental invocation. Skip running the apply altogether if they
	// don't answer to the affirmative
	if !deployConfig.BatchMode {
		approval := s.Util.PromptBasic("Do you really want to create DCE resources in your AWS account? (type \"yes\" or \"no\")", validateYesOrNo)

		if approval == nil || !strings.HasPrefix(strings.ToLower(*approval), "y") {
			return "", fmt.Errorf("user exited")
		}
	}

	log.Infoln("Creating DCE infrastructure")
	applyopts, _ := util.ParseOptions(&deployConfig.TFApplyOptions)
	if err := s.Util.Terraformer.Apply(ctx, applyopts); err != nil {
		return "", err
	}

	log.Infoln("Retrieving artifacts location")
	output, err := s.Util.Terraformer.GetOutput(ctx, "artifacts_bucket_name")

	if err != nil {
		return "", err
	}

	return output, nil
}

func (s *DeployService) deployCodeAssets(artifactsBucket string, namespace string, dceLocation string, dceVersion string) error {
	_, originDir := s.Util.ChToConfigDir()
	defer s.Util.Chdir(originDir)

	_, err := s.retrieveCodeAssets(dceLocation, dceVersion)
	if err != nil {
		return err
	}

	log.Debugln("Using \"%s\" for the artifact location.", artifactsBucket)

	lambdas, codebuilds := s.Util.UploadDirectoryToS3(s.Util.GetArtifactsDir(dceVersion), artifactsBucket, "")
	log.Debugln("Uploaded lambdas to S3: ", lambdas)
	log.Debugln("Uploaded codebuilds to S3: ", codebuilds)

	s.Util.UpdateLambdasFromS3Assets(lambdas, artifactsBucket, namespace)

	// No need to update Codebuild. It will pull from <bucket>/codebuild on its next build.

	return nil
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func (s *DeployService) getRandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (s *DeployService) getLocalTFMainContents(deployConfig *DeployConfig) (string, error) {
	var tfMainContents string
	// Generate the main.tf template...
	var buffer bytes.Buffer
	err := configureTemplate(s.Util.TFTemplater, deployConfig)
	if err != nil {
		return "", err
	}

	err = s.Util.TFTemplater.Write(&buffer)
	if err != nil {
		return "", err
	}
	tfMainContents = buffer.String()
	log.Debugln("Creating tf main.tf file with: ", tfMainContents)

	return tfMainContents, nil
}

// retrieveCodeAssets downloads the DCE build_artifacts.zip file
// for the provided DCE release version
// and unzips the contents into the artifacts directory (eg. ~/.dce/.cache/{version}/)
func (s *DeployService) retrieveCodeAssets(dceLocation string, dceVersion string) (string, error) {
	tmpDir, oldDir := s.Util.ChToTmpDir()

	defer os.Chdir(oldDir)

	if strings.HasPrefix(dceLocation, "github.com") {
		// Download release assets from github
		log.Infoln("Downloading DCE code assets")
		// Download the the build artifacts zip file from the Github release
		err := s.Util.Githuber.DownloadGithubReleaseAsset(AssetsFileName, dceVersion)
		if err != nil {
			return "", err
		}
		// Unarchive the build artifacts zip file into ~/.dce/.cache/{version}/
		err = s.Util.Unarchive(AssetsFileName, s.Util.GetArtifactsDir(dceVersion))
		if err != nil {
			return "", err
		}
		s.Util.RemoveAll(AssetsFileName)
	} else {
		// Use local assets
		log.Infof("Using local DCE binaries from %s", dceLocation)
		zippedAssetsPath := filepath.Join(dceLocation, "bin", AssetsFileName)
		err := s.Util.Unarchive(zippedAssetsPath, s.Util.GetArtifactsDir(dceVersion))
		if err != nil {
			return "", err
		}
	}

	return tmpDir, nil
}

func configureTemplate(t util.TFTemplater, deployConfig *DeployConfig) error {

	if deployConfig.AWSRegion != "" {
		_ = t.AddVariable("aws_region", "string", deployConfig.AWSRegion)
	}

	if deployConfig.Namespace != "" {
		_ = t.AddVariable("namespace", "string", deployConfig.Namespace)
	}
	if deployConfig.BudgetNotificationFromEmail != "" {
		_ = t.AddVariable("budget_notification_from_email", "string", deployConfig.BudgetNotificationFromEmail)
	}
	if len(deployConfig.BudgetNotificationBCCEmails) != 0 {
		budgetBCCEmails := ""
		for index, email := range deployConfig.BudgetNotificationBCCEmails {
			if index == 0 {
				budgetBCCEmails = "budget_notification_bcc_emails=["
			}
			budgetBCCEmails += "\"" + email + "\""
			if index < len(deployConfig.BudgetNotificationBCCEmails)-1 {
				budgetBCCEmails += ","
			} else {
				budgetBCCEmails += "]"
			}
		}
		_ = t.AddVariable("budget_notification_bcc_emails", "list(string)", budgetBCCEmails)
	}
	if deployConfig.BudgetNotificationTemplateHTML != "" {
		_ = t.AddVariable("budget_notification_template_html", "string", deployConfig.BudgetNotificationTemplateHTML)
	}
	if deployConfig.BudgetNotificationTemplateText != "" {
		_ = t.AddVariable("budget_notification_template_text", "string", deployConfig.BudgetNotificationTemplateText)
	}
	if deployConfig.BudgetNotificationTemplateSubject != "" {
		_ = t.AddVariable("budget_notification_template_subject", "string", deployConfig.BudgetNotificationTemplateSubject)
	}
	if deployConfig.Location != "" {
		var source string
		// Loading DCE module from github repo
		if strings.HasPrefix(deployConfig.Location, "github.com") {
			source = fmt.Sprintf("%s//modules", deployConfig.Location)
			if deployConfig.Version != "" {
				source = fmt.Sprintf("%s?ref=v%s", source, deployConfig.Version)
			}
		} else {
			source = fmt.Sprintf("%s/modules", deployConfig.Location)
		}

		t.SetModuleSource(source)
	}
	return nil
}

func validateYesOrNo(input string) error {
	if affirmAnswerRegex.MatchString(input) {
		return nil
	}
	return fmt.Errorf("\"%s\" is invalid", input)
}
