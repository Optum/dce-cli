package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Optum/dce-cli/internal/terra"
	"github.com/shurcooL/githubv4"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"

	"github.com/mholt/archiver"
)

var deployNamespace string
var dceRepoPath string

func init() {

	systemDeployCmd.Flags().StringVarP(&deployNamespace, "namespace", "n", "", "Set a custom terraform namespace (Optional)")
	systemDeployCmd.Flags().StringVarP(&dceRepoPath, "path", "p", "", "Path to local DCE repo")
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

		log.Println("Creating terraform remote state backend infrastructure")
		stateBucket := createRemoteStateBackend()

		log.Println("Creating DCE infrastructure")
		artifactsBucket := createDceInfra(stateBucket)
		log.Println("Artifacts bucket = ", artifactsBucket)

		// Deploy code assets to DCE infra
	},
}

func createDceInfra(stateBucket string) string {
	workingDir, originDir := mvToTempDir("dce-")
	defer os.RemoveAll(workingDir)
	defer os.Chdir(originDir)

	log.Println("Downloading DCE terraform modules")
	artifactsFileName := "terraform_artifacts.zip"
	downloadGithubReleaseAsset(artifactsFileName)

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

	err := archiver.Unarchive(artifactsFileName, ".")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	os.Remove(artifactsFileName)
	files, err := ioutil.ReadDir("./")
	if len(files) != 1 || !files[0].IsDir() {
		log.Fatalf("Unexpected content in DCE assets archive")
	}
	os.Chdir(files[0].Name())

	log.Println("Initializing terraform working directory")
	terra.Init([]string{"-backend-config=\"bucket=" + stateBucket + "\"", "-backend-config=\"key=local-tf-state\""})

	// log.Println("Building DCE infrastructure")
	// var namesSpace string
	// if deployNamespace != "" {
	// 	namesSpace = deployNamespace
	// } else {
	// 	namesSpace = "dce-default-" + getRandString(8)
	// }
	// terra.Apply(namesSpace)

	// log.Println("Retrieving artifacts bucket name from terraform outputs")
	// artifactsBucket := terra.GetOutput("artifacts_bucket_name")
	// log.Println("	-->", artifactsBucket)

	// return artifactsBucket
	return ""
}

func downloadGithubReleaseAsset(assetName string) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	oauthHTTPClient := oauth2.NewClient(context.Background(), src)

	var query struct {
		Viewer struct {
			Login     githubv4.String
			CreatedAt githubv4.DateTime
		}
		Repository struct {
			Releases struct {
				Nodes []struct {
					TagName       githubv4.String
					ReleaseAssets struct {
						Nodes []struct {
							ID          githubv4.String
							DownloadURL githubv4.String
							URL         string
						}
					} `graphql:"releaseAssets(last: 1, name: \"terraform_artifacts.zip\")"`
				}
			} `graphql:"releases(last: 1)"`
		} `graphql:"repository(owner: \"Optum\", name: \"Redbox\")"`
	}

	ghClient := githubv4.NewClient(oauthHTTPClient)
	err := ghClient.Query(context.Background(), &query, nil)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	fmt.Println("    Query Response:", query.Repository.Releases.Nodes[0].ReleaseAssets.Nodes[0].URL)

	req, err := http.NewRequest("GET", query.Repository.Releases.Nodes[0].ReleaseAssets.Nodes[0].URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(assetName)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func createRemoteStateBackend() string {
	workingDir, originDir := mvToTempDir("dce-init-")
	defer os.RemoveAll(workingDir)
	defer os.Chdir(originDir)

	log.Println("Creating terraform remote backend template (init.tf)")
	fileName := workingDir + "/" + "init.tf"
	err := ioutil.WriteFile(fileName, []byte(terra.RemoteBackend), 0644)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Println("Initializing terraform working directory and building remote state infrastructure")
	terra.Init([]string{})
	if deployNamespace != "" {
		terra.Apply(deployNamespace)
	} else {
		terra.Apply("dce-default-" + getRandString(8))
	}

	log.Println("Retrieving remote state bucket name from terraform outputs")
	stateBucket := terra.GetOutput("bucket")
	log.Println("	-->", stateBucket)

	return stateBucket
}

func mvToTempDir(prefix string) (string, string) {
	log.Println("Creating temporary terraform working directory")
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
