package cmd

import (
	"github.com/spf13/cobra"
)

var startDate float64
var endDate float64

func init() {
	usageCmd.Flags().Float64VarP(&startDate, "start-date", "s", 0, "The start date of the window over which usage information will be queried. (epoch timestamp)")
	usageCmd.Flags().Float64VarP(&endDate, "end-date", "e", 0, "The end date of the window over which usage information will be queried. (epoch timestamp)")
	if err := usageCmd.MarkFlagRequired("start-date"); err != nil {
		log.Fatalln("err: ", err)
	}
	if err := usageCmd.MarkFlagRequired("end-date"); err != nil {
		log.Fatalln("err: ", err)
	}
	RootCmd.AddCommand(usageCmd)
}

var usageCmd = &cobra.Command{
	Use:   "usage",
	Short: "View lease budget information",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		Service.GetUsage(startDate, endDate)
	},
}
