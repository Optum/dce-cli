package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	leasesCmd.AddCommand(leasesListCmd)
	leasesCmd.AddCommand(leasesCreateCmd)
	leasesCmd.AddCommand(leasesDestroyCmd)
	leasesCmd.AddCommand(leasesDestroyCmd)
	RootCmd.AddCommand(leasesCmd)
}

var leasesCmd = &cobra.Command{
	Use:   "leases",
	Short: "Manage dce leases",
}

var leasesDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "describe a lease",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Describe command")
	},
}

var leasesListCmd = &cobra.Command{
	Use:   "list",
	Short: "list leases",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("List command")
	},
}

var leasesCreateCmd = &cobra.Command{
	Use:   "add",
	Short: "Add one or more leases to the leases pool.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Create command")
	},
}

var leasesDestroyCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove one or more leases from the leases pool.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Destroy command")
	},
}
