package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	leasesCmd.AddCommand(leasesDescribeCmd)
	leasesCmd.AddCommand(leasesListCmd)
	leasesCmd.AddCommand(leasesCreateCmd)
	leasesCmd.AddCommand(leasesDestroyCmd)
	leasesCmd.AddCommand(leasesLoginCmd)
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
	Use:   "create",
	Short: "Create a lease.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Create command")
	},
}

var leasesDestroyCmd = &cobra.Command{
	Use:   "end",
	Short: "Cause a lease to immediately expire",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("Destroy command")
	},
}

var leasesLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to a leased DCE account",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print("login command")
	},
}
