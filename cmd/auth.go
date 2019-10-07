package cmd

import (
	"fmt"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var authUrl string

func init() {
	RootCmd.AddCommand(authCmd)
	authCmd.Flags().StringVarP(&authUrl, "url-override", "u", "", "Override the DCE login url")
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Login to dce",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Opening web browser. Please login and copy/paste the provided credentials into this terminal.")

		if authUrl == "" {
			authUrl = *config.System.Auth.LoginURL
		}
		browser.OpenURL(authUrl)
	},
}
