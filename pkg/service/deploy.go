package service

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Optum/dce-cli/configs"
	utl "github.com/Optum/dce-cli/internal/util"
	"github.com/mholt/archiver"
)

const ArtifactsFileName = "terraform_artifacts.zip"
const AssetsFileName = "build_artifacts.zip"

type DeployService struct {
	Config *configs.Root
	Util   *utl.UtilContainer
}

func (s *DeployService) Deploy(namespace string) {

	// TODO: Pass these directly into terraform
	os.Setenv("AWS_ACCESS_KEY_ID", *s.Config.System.MasterAccount.Credentials.AwsAccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", *s.Config.System.MasterAccount.Credentials.AwsSecretAccessKey)

	if namespace == "" {
		namespace = "dce-" + getRandString(6)
	}

	log.Println("Creating terraform remote state backend infrastructure")
	stateBucket := s.createRemoteStateBackend(namespace)

	log.Println("Creating DCE infrastructure")
	artifactsBucket := s.createDceInfra(namespace, stateBucket)
	log.Println("Artifacts bucket = ", artifactsBucket)

	log.Println("Deploying code assets to DCE infrastructure")
	s.deployCodeAssets(namespace, artifactsBucket)
}

func (s *DeployService) createRemoteStateBackend(namespace string) string {
	tmpDir, originDir := mvToTempDir("dce-init-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	log.Println("Creating terraform remote backend template (init.tf)")
	fileName := tmpDir + "/" + "init.tf"
	err := ioutil.WriteFile(fileName, []byte(utl.RemoteBackend), 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("Initializing terraform working directory and building remote state infrastructure")
	s.Util.Init([]string{})
	if namespace != "" {
		s.Util.Terraformer.Apply(namespace)
	} else {
		s.Util.Terraformer.Apply("dce-default-" + getRandString(8))
	}

	log.Println("Retrieving remote state bucket name from terraform outputs")
	stateBucket := s.Util.Terraformer.GetOutput("bucket")
	log.Println("Remote state bucket = ", stateBucket)

	return stateBucket
}

func (s *DeployService) createDceInfra(namespace string, stateBucket string) string {
	tmpDir, originDir := mvToTempDir("dce-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	log.Println("Downloading DCE terraform modules")
	s.Util.Githuber.DownloadGithubReleaseAsset(ArtifactsFileName)
	err := archiver.Unarchive(ArtifactsFileName, ".")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	os.Remove(ArtifactsFileName)
	files, err := ioutil.ReadDir("./")
	if len(files) != 1 || !files[0].IsDir() {
		log.Fatalf("Unexpected content in DCE assets archive")
	}
	os.Chdir(files[0].Name())

	log.Println("Initializing terraform working directory")
	s.Util.Terraformer.Init([]string{"-backend-config=bucket=" + stateBucket, "-backend-config=key=local-tf-state"})

	log.Println("Applying DCE infrastructure")
	s.Util.Terraformer.Apply(namespace)

	log.Println("Retrieving artifacts bucket name from terraform outputs")
	artifactsBucket := s.Util.Terraformer.GetOutput("artifacts_bucket_name")
	log.Println("artifacts bucket name = ", artifactsBucket)

	return artifactsBucket
}

func (s *DeployService) deployCodeAssets(deployNamespace string, artifactsBucket string) {
	tmpDir, originDir := mvToTempDir("dce-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	log.Println("Downloading DCE code assets")
	s.Util.Githuber.DownloadGithubReleaseAsset(AssetsFileName)
	// TODO:
	// Protect against zip-slip vulnerability? https://snyk.io/research/zip-slip-vulnerability
	//
	// err := z.Walk("/Users/matt/Desktop/test.zip", func(f archiver.File) error {
	// 	zfh, ok := f.Header.(zip.FileHeader)
	// 	if ok {
	// 		fmt.Println("Filename:", zfh.Name)
	// 	}
	// 	return nil
	// })
	err := archiver.Unarchive(AssetsFileName, ".")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	os.Remove(AssetsFileName)
	files, err := ioutil.ReadDir("./")
	if len(files) != 2 || !files[0].IsDir() || !files[1].IsDir() {
		log.Fatalf("Unexpected content in DCE assets archive")
	}

	lambdas, codebuilds := s.Util.UploadDirectoryToS3(".", artifactsBucket, "")
	log.Println("Uploaded lambdas to S3: ", lambdas)
	log.Println("Uploaded codebuilds to S3: ", codebuilds)

	s.Util.UpdateLambdasFromS3Assets(lambdas, artifactsBucket, deployNamespace)

	// aws lambda update-function-code \
	// --function-name ${fn_name} \
	// --s3-bucket ${artifactBucket} \
	// --s3-key lambda/${mod_name}.zip

	// 3. Publish new lambda versions
	// aws lambda publish-version \
	// --function-name ${fn_name}
}

func mvToTempDir(prefix string) (string, string) {
	log.Println("Creating temporary working directory")
	destinationDir, err := ioutil.TempDir("", prefix)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("	-->" + destinationDir)
	originDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	os.Chdir(destinationDir)
	return destinationDir, originDir
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
func getRandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
