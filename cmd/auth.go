package cmd

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var authUrl string

func init() {
	RootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVarP(&authUrl, "url-override", "u", "", "DCE version to deploy (Defaults to latest)")
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Configure DCE cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Opening web browser. Please login. You will be provided with credentials to copy/paste into this terminal.")

		if authUrl == "" {
			authUrl = *config.Auth.LoginUrl
		}
		browser.OpenURL(authUrl)
	},
}
