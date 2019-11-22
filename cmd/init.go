package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "First time DCE cli setup. Creates config file at ~/.dce.yaml",
	Run: func(cmd *cobra.Command, args []string) {
		Service.InitializeDCE()
	},
}
