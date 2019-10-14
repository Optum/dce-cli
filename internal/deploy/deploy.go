package deploy

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/Optum/dce-cli/internal/util/awshelpers"
	"github.com/Optum/dce-cli/internal/util/ghub"
	"github.com/Optum/dce-cli/internal/util/terra"
	"github.com/mholt/archiver"
)

const ArtifactsFileName = "terraform_artifacts.zip"
const AssetsFileName = "build_artifacts.zip"

func CreateRemoteStateBackend(namespace string) string {
	tmpDir, originDir := mvToTempDir("dce-init-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	log.Println("Creating terraform remote backend template (init.tf)")
	fileName := tmpDir + "/" + "init.tf"
	err := ioutil.WriteFile(fileName, []byte(terra.RemoteBackend), 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("Initializing terraform working directory and building remote state infrastructure")
	terra.Init([]string{})
	if namespace != "" {
		terra.Apply(namespace)
	} else {
		terra.Apply("dce-default-" + getRandString(8))
	}

	log.Println("Retrieving remote state bucket name from terraform outputs")
	stateBucket := terra.GetOutput("bucket")
	log.Println("Remote state bucket = ", stateBucket)

	return stateBucket
}

func CreateDceInfra(namespace string, stateBucket string) string {
	tmpDir, originDir := mvToTempDir("dce-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	log.Println("Downloading DCE terraform modules")
	ghub.DownloadGithubReleaseAsset(ArtifactsFileName)
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
	terra.Init([]string{"-backend-config=bucket=" + stateBucket, "-backend-config=key=local-tf-state"})

	log.Println("Applying DCE infrastructure")
	if namespace != "" {
		terra.Apply(namespace)
	} else {
		terra.Apply("dce-" + getRandString(6))
	}

	log.Println("Retrieving artifacts bucket name from terraform outputs")
	artifactsBucket := terra.GetOutput("artifacts_bucket_name")
	log.Println("artifacts bucket name = ", artifactsBucket)

	return artifactsBucket
}

func DeployCodeAssets(deployNamespace string, artifactsBucket string) {
	tmpDir, originDir := mvToTempDir("dce-")
	defer os.RemoveAll(tmpDir)
	defer os.Chdir(originDir)

	log.Println("Downloading DCE code assets")
	ghub.DownloadGithubReleaseAsset(AssetsFileName)
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
	// os.Chdir(files[0].Name())

	//LEFT OFF HERE, deploy to lambdas and stuff

	// 1. Upload lambda and codebuild zips to s3
	awshelpers.UploadDirectoryToS3(".", artifactsBucket, "")

	// 2. Point lambdas at the code in s3

	// 3. Publish new lambda versions
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
