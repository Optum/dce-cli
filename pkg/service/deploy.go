package service

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/Optum/dce-cli/configs"
	"github.com/Optum/dce-cli/internal/constants"
	observ "github.com/Optum/dce-cli/internal/observation"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/mholt/archiver"
)

const ArtifactsFileName = "terraform_artifacts.zip"
const AssetsFileName = "build_artifacts.zip"

type DeployService struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
	Util        *utl.UtilContainer
	LocalRepo   string
}

func (s *DeployService) Deploy(namespace, deployLocal string) {

	// TODO: Pass these directly into terraform
	os.Setenv("AWS_ACCESS_KEY_ID", *s.Config.System.MasterAccount.Credentials.AwsAccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", *s.Config.System.MasterAccount.Credentials.AwsSecretAccessKey)

	if namespace == "" {
		namespace = "dce-" + s.getRandString(6)
	}

	if deployLocal != "" {
		s.LocalRepo = deployLocal
	}

	stateBucket := s.createRemoteStateBackend(namespace)

	log.Infoln("Creating DCE infrastructure")
	artifactsBucket := s.createDceInfra(namespace, stateBucket)
	log.Infoln("Artifacts bucket = ", artifactsBucket)

	log.Infoln("Deploying code assets to DCE infrastructure")
	s.deployCodeAssets(namespace, artifactsBucket)
}

func (s *DeployService) createRemoteStateBackend(namespace string) string {
	tmpDir, originDir := s.mvToTempDir("dce-init-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	fileName := filepath.Join(tmpDir, "init.tf")

	remoteStateFile := s.getRemoteStateFile()
	err := ioutil.WriteFile(fileName, []byte(remoteStateFile), 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infoln("Initializing terraform working directory and building remote state infrastructure")
	s.Util.Init([]string{})
	if namespace != "" {
		s.Util.Terraformer.Apply(namespace)
	} else {
		s.Util.Terraformer.Apply("dce-default-" + s.getRandString(8))
	}

	log.Infoln("Retrieving remote state bucket name from terraform outputs")
	stateBucket := s.Util.Terraformer.GetOutput("bucket")
	log.Infoln("Remote state bucket = ", stateBucket)

	return stateBucket
}

func (s *DeployService) createDceInfra(namespace string, stateBucket string) string {
	tmpDir, originDir := s.mvToTempDir("dce-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	tfModulesDir := s.getTFModulesDir()
	files, err := ioutil.ReadDir(tfModulesDir)
	if err != nil {
		log.Fatalln(err)
	}

	if len(files) != 1 || !files[0].IsDir() {
		log.Fatalf("Unexpected content in DCE assets archive")
	}
	os.Chdir(files[0].Name())

	log.Infoln("Initializing terraform working directory")
	s.Util.Terraformer.Init([]string{"-backend-config=bucket=" + stateBucket, "-backend-config=key=local-tf-state"})

	log.Infoln("Applying DCE infrastructure")
	s.Util.Terraformer.Apply(namespace)

	log.Infoln("Retrieving artifacts bucket name from terraform outputs")
	artifactsBucket := s.Util.Terraformer.GetOutput("artifacts_bucket_name")
	log.Infoln("artifacts bucket name = ", artifactsBucket)

	return artifactsBucket
}

func (s *DeployService) deployCodeAssets(deployNamespace string, artifactsBucket string) {
	tmpDir, originDir := s.mvToTempDir("dce-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	assetsDir := s.getAssetsDir()
	files, err := ioutil.ReadDir(assetsDir)
	if err != nil {
		log.Fatalln(err)
	}

	if len(files) != 2 || !files[0].IsDir() || !files[1].IsDir() {
		log.Fatalf("Unexpected content in DCE assets archive")
	}

	lambdas, codebuilds := s.Util.UploadDirectoryToS3(".", artifactsBucket, "")
	log.Infoln("Uploaded lambdas to S3: ", lambdas)
	log.Infoln("Uploaded codebuilds to S3: ", codebuilds)

	s.Util.UpdateLambdasFromS3Assets(lambdas, artifactsBucket, deployNamespace)

	// No need to update Codebuild. It will pull from <bucket>/codebuild on its next build.
}

func (s *DeployService) mvToTempDir(prefix string) (string, string) {
	// log.Infoln("Creating temporary working directory")
	destinationDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Fatalln(err)
	}
	// log.Infoln("	-->" + destinationDir)
	originDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	os.Chdir(destinationDir)
	return destinationDir, originDir
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
		log.Debug(s.LocalRepo)
		path := filepath.Join(s.LocalRepo, "scripts", "deploy_local", "main.tf")
		remoteStateBackend = s.Util.ReadFromFile(path)
	} else {
		remoteStateBackend = constants.RemoteBackend
	}
	log.Debugln("Getting tf remote state backend file from: ", remoteStateBackend)

	return remoteStateBackend
}

func (s *DeployService) getTFModulesDir() string {
	workingDir, err := os.Getwd()

	if s.LocalRepo != "" {
		zippedArtifactsPath := filepath.Join(s.LocalRepo, "bin", ArtifactsFileName)
		err = archiver.Unarchive(zippedArtifactsPath, workingDir)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln("Downloading DCE terraform modules")
		s.Util.Githuber.DownloadGithubReleaseAsset(ArtifactsFileName)
		err = archiver.Unarchive(ArtifactsFileName, workingDir)
		if err != nil {
			log.Fatalln(err)
		}
		os.Remove(ArtifactsFileName)
	}

	return workingDir
}

func (s *DeployService) getAssetsDir() string {
	var workingDir, err = os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}

	if s.LocalRepo != "" {
		zippedAssetsPath := filepath.Join(s.LocalRepo, "bin", AssetsFileName)
		err = archiver.Unarchive(zippedAssetsPath, workingDir)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		if err != nil {
			log.Fatalln(err)
		}
		log.Infoln("Downloading DCE code assets")
		s.Util.Githuber.DownloadGithubReleaseAsset(AssetsFileName)
		err = archiver.Unarchive(AssetsFileName, workingDir)
		if err != nil {
			log.Fatalln(err)
		}
		os.Remove(AssetsFileName)
	}

	return workingDir
}
