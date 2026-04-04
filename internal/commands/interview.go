package commands

import (
	"github.com/spf13/cobra"
)

var interviewCmd = &cobra.Command{
	Use:   "interviewCmd",
	Short: "Query interviewCmd questions",
}

func init() {
	rootCmd.AddCommand(interviewCmd)

	interviewCmd.Flags().StringP("company", "c", "", "target company")
	interviewCmd.Flags().StringP("limit", "l", "5", "response query limit")
	interviewCmd.MarkFlagRequired("company")
}
