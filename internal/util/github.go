package util

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Optum/dce-cli/configs"
	observ "github.com/Optum/dce-cli/internal/observation"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

const RepoName = "Redbox"
const RepoOwner = "Optum"

type GithubUtil struct {
	Config      *configs.Root
	Observation *observ.ObservationContainer
}

func (u *GithubUtil) DownloadGithubReleaseAsset(assetName string) {
	src := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: *u.Config.GithubToken},
	)
	oauthHTTPClient := oauth2.NewClient(context.Background(), src)

	variables := map[string]interface{}{
		"assetName": githubv4.String(assetName),
		"repoName":  githubv4.String(RepoName),
		"repoOwner": githubv4.String(RepoOwner),
	}

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
							URL string
						}
					} `graphql:"releaseAssets(last: 1, name: $assetName)"`
				}
			} `graphql:"releases(last: 1)"`
		} `graphql:"repository(owner: $repoOwner, name: $repoName)"`
	}

	githubClient := githubv4.NewClient(oauthHTTPClient)
	err := githubClient.Query(context.Background(), &query, variables)
	if err != nil {
		Log.Fatalf("error: %v", err)
	}
	fmt.Println("    Query Response:", query.Repository.Releases.Nodes[0].ReleaseAssets.Nodes[0].URL)

	req, err := http.NewRequest("GET", query.Repository.Releases.Nodes[0].ReleaseAssets.Nodes[0].URL, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		Log.Fatalf("error: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(assetName)
	if err != nil {
		Log.Fatalf("error: %v", err)
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		Log.Fatalf("error: %v", err)
	}
}
