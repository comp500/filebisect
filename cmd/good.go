package cmd

import (
	"fmt"
	"os"

	"github.com/comp500/filebisect/index"
	"github.com/spf13/cobra"
)

// goodCmd represents the good command
var goodCmd = &cobra.Command{
	Use:   "good",
	Short: "Mark the current set of files as good",
	Long: `This command marks the current set of files as good, and moves files around to allow you to try again. Files that are marked as good are not moved.
The files moved are randomly selected, so that hopefully the conflict or error will be isolatable.
The files are moved to a temporary directory, specified in file-bisect-index.toml, and you can move files around manually if there are dependency errors.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		idx, err := index.LoadOrCreateIndex()
		if err != nil {
			fmt.Printf("Error loading index file: %v\n", err)
			os.Exit(1)
		}
		idx.Init()
		idx.Good()
		err = idx.Save()
		if err != nil {
			fmt.Printf("Error saving index file: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(goodCmd)
}
