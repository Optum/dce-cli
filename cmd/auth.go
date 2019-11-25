package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(authCmd)
	authCmd.Flags()
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Login to dce",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Service.Authenticate()
	},
}
