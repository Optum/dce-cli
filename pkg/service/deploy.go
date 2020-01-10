package service

import (
	"bytes"
	"context"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/Optum/dce-cli/internal/util"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/pkg/errors"
)

const ArtifactsFileName = "terraform_artifacts.zip"
const AssetsFileName = "build_artifacts.zip"

type DeployService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
	LocalRepo   string
}

// Deploy writes the local `main.tf` file, using the overrides, and then
// calls Terraform init and apply using configuration directory (`~/.dce`)
// as the working folder and location of local state.
func (s *DeployService) Deploy(ctx context.Context, overrides *DeployOverrides) error {

	// Initialize the folder structure
	if err := s.Util.CreateConfigDirTree(); err != nil {
		return errors.Wrap(err, "error creating directory structure")
	}

	// Generate a namespace if one has not been supplied because the
	// dce terraform module requires this argument.
	if overrides.Namespace == "" {
		overrides.Namespace = "dce-" + s.getRandString(8)
	}
	// This is also a required field by the module, so it has to
	// have some value.
	if overrides.BudgetNotificationFromEmail == "" {
		overrides.BudgetNotificationFromEmail = "no-reply@example.com"
	}

	cfg := ctx.Value(constants.DeployConfig).(*configs.DeployConfig)

	if cfg.DeployLocalPath != "" {
		s.LocalRepo = cfg.DeployLocalPath
	}

	_, err := s.createTFMainFile(overrides, cfg.UseCached)

	if err != nil {
		return errors.Wrap(err, "error creating local backend")
	}

	artifactsBucket, err := s.createDceInfra(ctx, overrides)

	if err != nil {
		return errors.Wrap(err, "error creating infrastructure")
	}

	log.Infoln("Artifacts bucket = ", artifactsBucket)

	log.Infoln("Deploying code assets to DCE infrastructure")
	s.deployCodeAssets(artifactsBucket, overrides)
	return nil
}

// createTFMainFile creates the main.tf file. The default behavior, without
// using a local repository, is to create the file with the bare minimum
// required to use Terraform with local state.
func (s *DeployService) createTFMainFile(overrides *DeployOverrides, usecached bool) (string, error) {
	_, originDir := s.Util.ChToConfigDir()
	defer s.Util.Chdir(originDir)

	fileName := s.Util.GetLocalMainTFFile()

	if s.Util.IsExistingFile(fileName) && usecached {
		log.Warnln("'main.tf' already exists and --use-cached specified; using existing file")
	} else {
		tfMainContents, err := s.getLocalTFMainContents(overrides)
		if err != nil {
			return "", err
		}
		s.Util.WriteFile(fileName, tfMainContents)
	}
	return fileName, nil
}

func (s *DeployService) createDceInfra(ctx context.Context, overrides *DeployOverrides) (string, error) {
	_, originDir := s.Util.ChToConfigDir()
	defer s.Util.Chdir(originDir)

	s.retrieveTFModules()

	deployLogFileName := s.Util.GetLogFile()
	ctx = context.WithValue(ctx, "deployLogFile", deployLogFileName)

	log.Infoln("Initializing")
	s.Util.Terraformer.Init(ctx, []string{})

	log.Infoln("Creating DCE infrastructure")
	s.Util.Terraformer.Apply(ctx, []string{})

	log.Infoln("Retrieving artifacts location")
	return s.Util.Terraformer.GetOutput(ctx, "artifacts_bucket_name")
}

func (s *DeployService) deployCodeAssets(artifactsBucket string, overrides *DeployOverrides) {
	_, originDir := s.Util.ChToConfigDir()
	defer s.Util.Chdir(originDir)

	s.retrieveCodeAssets()

	log.Debugln("Using \"%s\" for the artifact location.", artifactsBucket)

	lambdas, codebuilds := s.Util.UploadDirectoryToS3(s.Util.GetArtifactsDir(), artifactsBucket, "")
	log.Debugln("Uploaded lambdas to S3: ", lambdas)
	log.Debugln("Uploaded codebuilds to S3: ", codebuilds)

	s.Util.UpdateLambdasFromS3Assets(lambdas, artifactsBucket, overrides.Namespace)

	// No need to update Codebuild. It will pull from <bucket>/codebuild on its next build.
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

func (s *DeployService) getLocalTFMainContents(overrides *DeployOverrides) (string, error) {
	var tfMainContents string
	if s.LocalRepo != "" {
		path := filepath.Join(s.LocalRepo, "scripts", "deploy_local", "main.tf")
		tfMainContents = s.Util.ReadFromFile(path)
	} else {
		// Generate the main.tf template...
		var buffer bytes.Buffer
		addOverridesToTemplate(s.Util.TFTemplater, overrides)

		err := s.Util.TFTemplater.Write(&buffer)
		if err != nil {
			return "", err
		}
		tfMainContents = buffer.String()
	}
	log.Debugln("Creating tf main.tf file with: ", tfMainContents)

	return tfMainContents, nil
}

func (s *DeployService) retrieveTFModules() (string, error) {
	workingDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	if s.LocalRepo != "" {
		zippedArtifactsPath := filepath.Join(s.LocalRepo, "bin", ArtifactsFileName)
		s.Util.Unarchive(zippedArtifactsPath, workingDir)
	}

	return workingDir, nil
}

func (s *DeployService) retrieveCodeAssets() (string, error) {
	tmpDir, oldDir := s.Util.ChToTmpDir()

	defer os.Chdir(oldDir)

	if s.LocalRepo != "" {
		zippedAssetsPath := filepath.Join(s.LocalRepo, "bin", AssetsFileName)
		s.Util.Unarchive(zippedAssetsPath, s.Util.GetArtifactsDir())
	} else {
		log.Infoln("Downloading DCE code assets")
		s.Util.Githuber.DownloadGithubReleaseAsset(AssetsFileName)
		s.Util.Unarchive(AssetsFileName, s.Util.GetArtifactsDir())
		s.Util.RemoveAll(AssetsFileName)
	}

	return tmpDir, nil
}

func addOverridesToTemplate(t util.TFTemplater, overrides *DeployOverrides) error {

	if overrides.AWSRegion != "" {
		_ = t.AddVariable("aws_region", "string", overrides.AWSRegion)
	}

	globalTags := "global_tags={" + constants.GlobalTFTagDefaults
	if len(overrides.GlobalTags) != 0 {
		for _, tag := range overrides.GlobalTags {
			globalTags += ",\"" + strings.ReplaceAll(tag, ":", "\":\"") + "\""
		}
	}
	globalTags += "}"
	// _ = t.AddVariable("global_tags", "map(string)", globalTags)

	if overrides.Namespace != "" {
		_ = t.AddVariable("namespace", "string", overrides.Namespace)
	}
	if overrides.BudgetNotificationFromEmail != "" {
		_ = t.AddVariable("budget_notification_from_email", "string", overrides.BudgetNotificationFromEmail)
	}
	if len(overrides.BudgetNotificationBCCEmails) != 0 {
		budgetBCCEnails := ""
		for index, email := range overrides.BudgetNotificationBCCEmails {
			if index == 0 {
				budgetBCCEnails = "budget_notification_bcc_emails=["
			}
			budgetBCCEnails += "\"" + email + "\""
			if index < len(overrides.BudgetNotificationBCCEmails)-1 {
				budgetBCCEnails += ","
			} else {
				budgetBCCEnails += "]"
			}
		}
		_ = t.AddVariable("budget_notification_bcc_emails", "list(string)", budgetBCCEnails)
	}
	if overrides.BudgetNotificationTemplateHTML != "" {
		_ = t.AddVariable("budget_notification_template_html", "string", overrides.BudgetNotificationTemplateHTML)
	}
	if overrides.BudgetNotificationTemplateText != "" {
		_ = t.AddVariable("budget_notification_template_text", "string", overrides.BudgetNotificationTemplateText)
	}
	if overrides.BudgetNotificationTemplateSubject != "" {
		_ = t.AddVariable("budget_notification_template_subject", "string", overrides.BudgetNotificationTemplateSubject)
	}
	return nil
}
