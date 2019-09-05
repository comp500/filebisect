package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Resets the good/bad rating of all files to unknown",
	Long: `If the bisection has gone wrong (e.g. you accidentally used good or bad), you can use this command to restart bisection, by
changing the good/badness of all files to unknown, ignoring ignored files.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Not Yet Implemented!!!!!")
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
