package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	terraform "github.com/hashicorp/terraform/command"
	"github.com/mitchellh/cli"
	"github.com/spf13/cobra"
)

var deployNamespace string

func init() {

	systemDeployCmd.Flags().StringVarP(&deployNamespace, "namespace", "n", "", "Set a custom terraform namespace (Optional)")
	systemCmd.AddCommand(systemDeployCmd)

	systemLogsCmd.AddCommand(systemLogsAccountsCmd)
	systemLogsCmd.AddCommand(systemLogsLeasesCmd)
	systemLogsCmd.AddCommand(systemLogsUsageCmd)
	systemLogsCmd.AddCommand(systemLogsResetCmd)
	systemCmd.AddCommand(systemLogsCmd)

	systemUsersCmd.AddCommand(systemUsersAddCmd)
	systemUsersCmd.AddCommand(systemUsersRemoveCmd)
	systemCmd.AddCommand(systemUsersCmd)

	RootCmd.AddCommand(systemCmd)
}

var systemCmd = &cobra.Command{
	Use:   "system",
	Short: "Deploy and configure the DCE system",
}

/*
Deploy Namespace
*/

var systemDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the DCE system",
	Run: func(cmd *cobra.Command, args []string) {

		log.Println("Creating temporary terraform working directory")

		tfBackendTemplate :=
			`provider "aws" {
region = "us-east-1"
}

variable "global_tags" {
description = "The tags to apply to all resources that support tags"
type        = map(string)

default = {
	Terraform = "True"
	AppName   = "AWS Redbox Management"
	Source    = "github.com/Optum/Redbox//modules"
	Contact   = "fake_email@domain.com"
}
}

variable "namespace" {
type = string
}

data "aws_caller_identity" "current" {}

# Configure an S3 Bucket to hold artifacts
# (eg. application code deployments, etc.)
resource "aws_s3_bucket" "local_tfstate" {
bucket = "${data.aws_caller_identity.current.account_id}-local-tfstate-${var.namespace}"

# Allow S3 access logs to be written to this bucket
acl = "log-delivery-write"

# Allow Terraform to destroy the bucket
# (so ephemeral PR environments can be torn down)
force_destroy = true

# Encrypt objects by default
server_side_encryption_configuration {
	rule {
	apply_server_side_encryption_by_default {
		sse_algorithm = "AES256"
	}
	}
}

versioning {
	enabled = true
}

# Send S3 access logs for this bucket to itself
logging {
	target_bucket = "${data.aws_caller_identity.current.account_id}-local-tfstate-${var.namespace}"
	target_prefix = "logs/"
}

tags = var.global_tags
}

# Enforce SSL only access to the bucket
resource "aws_s3_bucket_policy" "reset_codepipeline_source_ssl_policy" {
bucket = aws_s3_bucket.local_tfstate.id

policy = <<POLICY
{
	"Version": "2012-10-17",
	"Statement": [
	{
		"Sid": "DenyInsecureCommunications",
		"Effect": "Deny",
		"Principal": "*",
		"Action": "s3:*",
		"Resource": "${aws_s3_bucket.local_tfstate.arn}/*",
		"Condition": {
			"Bool": {
				"aws:SecureTransport": "false"
			}
		}
	}
	]
}
POLICY

}

resource "aws_dynamodb_table" "local_terraform_state_lock" {
name           = "Terraform-State-Backend-${var.namespace}"
read_capacity  = 1
write_capacity = 1
hash_key       = "LockID"

attribute {
	name = "LockID"
	type = "S"
}
}

output "bucket" {
value = aws_s3_bucket.local_tfstate.bucket
}
`
		tempDir, err := ioutil.TempDir("", "dce-init-")
		currentDir, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		os.Chdir(tempDir)
		defer os.Chdir(currentDir)
		defer os.RemoveAll(tempDir)

		if err != nil {
			log.Fatalf("Error: ", err)
		}

		fileName := tempDir + "/" + "init.tf"
		err = ioutil.WriteFile(fileName, []byte(tfBackendTemplate), 0644)
		if err != nil {
			log.Fatalf("error: %v", err)
		}

		log.Println("Temporary directory created at: " + tempDir)

		log.Println("Initializing working directory and building remote state infrastructure")
		terraformInit()
		if deployNamespace != "" {
			terraformApply(deployNamespace)
		} else {
			terraformApply("dce-default-" + getRandString(8))
		}

		log.Println("Retrieving remote state bucket name from terraform outputs")
		stateBucket := getTerraformOutput("bucket")

		log.Println("State Bucket = " + stateBucket)
		// EASY TEMPORARY OPTION FOR LAMBDA ARTIFACT DEPLOYMENT:
		//   Use flag to set path to local DCE repo. Build all artifacts there and deploy them.
	},
}

func terraformInit() {
	log.Println("Running terraform init")

	tfInit := &terraform.InitCommand{
		Meta: terraform.Meta{
			Ui: getTerraformUI(),
		},
	}
	tfInit.Run([]string{})
}

func terraformApply(namespace string) {
	log.Println("Running terraform apply with namespace: " + namespace)
	tfApply := &terraform.ApplyCommand{
		Meta: terraform.Meta{
			Ui: getTerraformUI(),
		},
	}
	namespaceKey := "-var"
	namespaceValue := "namespace=" + namespace
	tfApply.Run([]string{namespaceKey, namespaceValue})
}

func getTerraformOutput(key string) string {
	log.Println("Retrieving terraform output for: " + key)
	outputCaptorUI := &UIOutputCaptor{
		BasicUi: &cli.BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stderr,
		},
		Captor: new(string),
	}
	tfOutput := &terraform.OutputCommand{
		Meta: terraform.Meta{
			Ui: outputCaptorUI,
		},
	}
	tfOutput.Run([]string{"bucket"})
	return *outputCaptorUI.Captor
}

func getTerraformUI() *cli.PrefixedUi {
	basicUI := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}
	prefix := "\nTerraform Says...\n"
	return &cli.PrefixedUi{
		AskPrefix:       prefix,
		AskSecretPrefix: prefix,
		OutputPrefix:    prefix,
		InfoPrefix:      prefix,
		ErrorPrefix:     prefix,
		WarnPrefix:      prefix,
		Ui:              basicUI,
	}
}

type UIOutputCaptor struct {
	Captor *string
	*cli.BasicUi
}

func (u *UIOutputCaptor) Output(message string) {
	u.Captor = &message
	u.BasicUi.Output(message)
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

/*
Logs Namespace
*/

var systemLogsCmd = &cobra.Command{
	Use:   "logs",
	Short: "View logs",
}

var systemLogsAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "View account logs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Accounts command")
	},
}

var systemLogsLeasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "View lease logs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Leases command")
	},
}

var systemLogsUsageCmd = &cobra.Command{
	Use:   "usage",
	Short: "View usage logs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Usage command")
	},
}

var systemLogsResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "View reset logs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Reset command")
	},
}

/*
Users Namespace
*/
var systemUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Manage users",
}

var systemUsersAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add users",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Add command")
	},
}

var systemUsersRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove users",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Remove command")
	},
}
