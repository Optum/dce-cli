package util

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Optum/dce-cli/configs"
	"github.com/aws/aws-sdk-go/aws"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AWSUtil struct {
	Config  *configs.Root
	Session *awsSession.Session
}

func (u *AWSUtil) UploadDirectoryToS3(localPath string, bucket string, prefix string) ([]string, []string) {
	walker := make(fileWalk)
	go func() {
		// Gather the files to upload by walking the path recursively
		if err := filepath.Walk(localPath, walker.Walk); err != nil {
			log.Fatalln("Walk failed:", err)
		}
		close(walker)
	}()

	// For each file found walking, upload it to S3
	lambdas := []string{}
	codebuilds := []string{}
	uploader := s3manager.NewUploader(u.Session)
	for path := range walker {
		rel, err := filepath.Rel(localPath, path)
		if err != nil {
			log.Fatalln("Unable to get relative path:", path, err)
		}
		file, err := os.Open(path)
		if err != nil {
			log.Println("Failed opening file", path, err)
			continue
		}
		defer file.Close()
		result, err := uploader.Upload(&s3manager.UploadInput{
			Bucket: &bucket,
			Key:    aws.String(filepath.Join(prefix, rel)),
			Body:   file,
		})
		if err != nil {
			log.Fatalln("Failed to upload", path, err)
		}
		log.Println("Uploaded", path, result.Location)

		parent := filepath.Base(filepath.Dir(path))
		if parent == "lambda" {
			lambdas = append(lambdas, filepath.Base(path))
		}
		if parent == "codebuild" {
			codebuilds = append(lambdas, filepath.Base(path))
		}
	}
	return lambdas, codebuilds
}

type fileWalk chan string

func (f fileWalk) Walk(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if !info.IsDir() {
		f <- path
	}
	return nil
}

func (u *AWSUtil) UpdateLambdasFromS3Assets(lambdaNames []string, bucket string, namespace string) {
	client := lambda.New(u.Session)

	for _, l := range lambdaNames {

		name := strings.Trim(l, ".zip")
		log.Println("Updating lambda config for: ", name)

		input := &lambda.UpdateFunctionCodeInput{
			FunctionName: aws.String(name + "-" + namespace),
			Publish:      aws.Bool(true),
			S3Bucket:     aws.String(bucket),
			S3Key:        aws.String("lambda/" + name + ".zip"),
			// S3ObjectVersion: aws.String("1"),
			// ZipFile:         []byte("fileb://file-path/file.zip"),
		}

		out, err := json.Marshal(input)
		if err != nil {
			panic(err)
		}

		log.Println("Input: ", string(out))

		updateLambdaConfig, _ := client.UpdateFunctionCode(input)

		out, err = json.Marshal(updateLambdaConfig)
		if err != nil {
			panic(err)
		}

		log.Println("Updated Lambda Config: ", string(out))
	}
}