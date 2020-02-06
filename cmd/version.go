package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

var version string

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "View the running version of dce-cli",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		shortVersion := strings.Split(version, "-")[0]
		finalVersion := strings.Replace(shortVersion, "v", "", 1)
		fmt.Println(finalVersion)
	},
}
