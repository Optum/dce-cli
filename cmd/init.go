package cmd

import (
	"fmt"

	"github.com/Optum/dce-cli/internal/constants"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: fmt.Sprintf("First time DCE cli setup. Creates config file at \"%s\" (by default) or at the location specified by \"--config\"", constants.ConfigFileDefaultLocationUnexpanded),
	Run: func(cmd *cobra.Command, args []string) {
		Service.InitializeDCE()
	},
}
