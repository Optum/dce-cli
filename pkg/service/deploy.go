package service

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
)

const ArtifactsFileName = "terraform_artifacts.zip"
const AssetsFileName = "build_artifacts.zip"

type DeployService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
	LocalRepo   string
}

func (s *DeployService) Deploy(deployLocal string, overrides *DeployOverrides) {
	if overrides.Namespace == "" {
		overrides.Namespace = "dce-" + s.getRandString(8)
	}

	if deployLocal != "" {
		s.LocalRepo = deployLocal
	}

	stateBucket := s.createRemoteStateBackend(overrides)

	log.Infoln("Creating DCE infrastructure")
	artifactsBucket := s.createDceInfra(stateBucket, overrides)
	log.Infoln("Artifacts bucket = ", artifactsBucket)

	log.Infoln("Deploying code assets to DCE infrastructure")
	s.deployCodeAssets(artifactsBucket, overrides)
}

func (s *DeployService) createRemoteStateBackend(overrides *DeployOverrides) string {
	tmpDir, originDir := s.Util.MvToTempDir("dce-")
	defer s.Util.RemoveAll(tmpDir)
	defer s.Util.Chdir(originDir)

	fileName := filepath.Join(tmpDir, "init.tf")

	remoteStateFile := s.getRemoteStateFile()
	s.Util.WriteFile(fileName, remoteStateFile)

	log.Infoln("Initializing terraform working directory and building remote state infrastructure")
	s.Util.Init([]string{})

	args := []string{"namespace=" + overrides.Namespace}
	s.Util.Terraformer.Apply(args)

	log.Infoln("Retrieving remote state bucket name from terraform outputs")
	stateBucket := s.Util.Terraformer.GetOutput("bucket")
	log.Infoln("Remote state bucket = ", stateBucket)

	return stateBucket
}

func (s *DeployService) createDceInfra(stateBucket string, overrides *DeployOverrides) string {
	tmpDir, originDir := s.Util.MvToTempDir("dce-")
	defer s.Util.RemoveAll(tmpDir)
	defer s.Util.Chdir(originDir)

	tfModulesDir := s.retrieveTFModules()
	files := s.Util.ReadDir(tfModulesDir)

	s.Util.Chdir(files[0].Name())

	log.Infoln("Initializing terraform working directory")
	s.Util.Terraformer.Init([]string{"-backend-config=bucket=" + stateBucket, "-backend-config=key=local-tf-state"})

	log.Infoln("Applying DCE infrastructure")

	args := argsFromOverrides(overrides)
	s.Util.Terraformer.Apply(args)

	log.Infoln("Retrieving artifacts bucket name from terraform outputs")
	artifactsBucket := s.Util.Terraformer.GetOutput("artifacts_bucket_name")
	log.Infoln("artifacts bucket name = ", artifactsBucket)

	return artifactsBucket
}

func (s *DeployService) deployCodeAssets(artifactsBucket string, overrides *DeployOverrides) {
	tmpDir, originDir := s.Util.MvToTempDir("dce-")
	defer s.Util.RemoveAll(tmpDir)
	defer s.Util.Chdir(originDir)

	s.retrieveCodeAssets()

	lambdas, codebuilds := s.Util.UploadDirectoryToS3(".", artifactsBucket, "")
	log.Infoln("Uploaded lambdas to S3: ", lambdas)
	log.Infoln("Uploaded codebuilds to S3: ", codebuilds)

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

func (s *DeployService) getRemoteStateFile() string {
	var remoteStateBackend string
	if s.LocalRepo != "" {
		path := filepath.Join(s.LocalRepo, "scripts", "deploy_local", "main.tf")
		remoteStateBackend = s.Util.ReadFromFile(path)
	} else {
		remoteStateBackend = constants.RemoteBackend
	}
	log.Debugln("Getting tf remote state backend file from: ", remoteStateBackend)

	return remoteStateBackend
}

func (s *DeployService) retrieveTFModules() string {
	workingDir, err := os.Getwd()

	if s.LocalRepo != "" {
		zippedArtifactsPath := filepath.Join(s.LocalRepo, "bin", ArtifactsFileName)
		s.Util.Unarchive(zippedArtifactsPath, workingDir)
	} else {
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln("Downloading DCE terraform modules")
		s.Util.Githuber.DownloadGithubReleaseAsset(ArtifactsFileName)
		s.Util.Unarchive(ArtifactsFileName, workingDir)
		os.Remove(ArtifactsFileName)
	}

	return workingDir
}

func (s *DeployService) retrieveCodeAssets() string {
	var workingDir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	if s.LocalRepo != "" {
		zippedAssetsPath := filepath.Join(s.LocalRepo, "bin", AssetsFileName)
		s.Util.Unarchive(zippedAssetsPath, workingDir)
	} else {
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln("Downloading DCE code assets")
		s.Util.Githuber.DownloadGithubReleaseAsset(AssetsFileName)
		s.Util.Unarchive(AssetsFileName, workingDir)
		os.Remove(AssetsFileName)
	}

	return workingDir
}

func argsFromOverrides(overrides *DeployOverrides) []string {
	args := []string{}

	if overrides.AWSRegion != "" {
		args = append(args, "aws_region="+overrides.AWSRegion)
	}
	globalTags := "global_tags={" + constants.GlobalTFTagDefaults
	if len(overrides.GlobalTags) != 0 {
		for _, tag := range overrides.GlobalTags {
			globalTags += ",\"" + strings.ReplaceAll(tag, ":", "\":\"") + "\""
		}
	}
	globalTags += "}"
	args = append(args, globalTags)

	if overrides.Namespace != "" {
		args = append(args, "namespace="+overrides.Namespace)
	}
	if overrides.BudgetNotificationFromEmail != "" {
		args = append(args, "budget_notification_from_email="+overrides.BudgetNotificationFromEmail)
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
		args = append(args, budgetBCCEnails)
	}
	if overrides.BudgetNotificationTemplateHTML != "" {
		args = append(args, "budget_notification_template_html="+overrides.BudgetNotificationTemplateHTML)
	}
	if overrides.BudgetNotificationTemplateText != "" {
		args = append(args, "budget_notification_template_text="+overrides.BudgetNotificationTemplateText)
	}
	if overrides.BudgetNotificationTemplateSubject != "" {
		args = append(args, "budget_notification_template_subject="+overrides.BudgetNotificationTemplateSubject)
	}

	return args
}
