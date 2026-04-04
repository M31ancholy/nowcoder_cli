package commands

import (
	"nowcoder_cli/service/nowcoder"
	"strconv"

	"github.com/spf13/cobra"
)

var interviewCmd = &cobra.Command{
	Use:   "interviewCmd",
	Short: "Query interviewCmd questions",
	Run: func(cmd *cobra.Command, args []string) {
		company, _ := cmd.Flags().GetString("company")
		position, _ := cmd.Flags().GetString("position")
		limit, _ := cmd.Flags().GetString("limit")
		limitInt, _ := strconv.Atoi(limit)
		nowcoder.GetInterviews(company, position, limitInt)
	},
}

func init() {
	rootCmd.AddCommand(interviewCmd)

	interviewCmd.Flags().StringP("company", "c", "", "target company")
	interviewCmd.Flags().StringP("position", "p", "", "target job position")
	interviewCmd.Flags().StringP("limit", "l", "5", "response query limit")
	interviewCmd.MarkFlagRequired("company")
}
