package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(configureCmd)
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure DCE cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: implement")
		fmt.Println("Master Account Access Key: " + *config.System.MasterAccount.Credentials.AwsAccessKeyID)
	},
}
