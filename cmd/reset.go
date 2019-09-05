package cmd

import (
	"fmt"
	"os"

	"github.com/comp500/filebisect/index"
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
		idx, err := index.LoadOrCreateIndex()
		if err != nil {
			fmt.Printf("Error loading index file: %v\n", err)
			os.Exit(1)
		}
		idx.Init()
		idx.Reset() // TODO: warn if no files are reset?
		err = idx.Save()
		if err != nil {
			fmt.Printf("Error saving index file: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
