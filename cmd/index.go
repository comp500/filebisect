package cmd

import (
	"fmt"
	"os"

	"github.com/comp500/filebisect/index"
	"github.com/spf13/cobra"
)

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generate an index of the current folder",
	Long: `Before using any other command, this command must be used to list the files in the folder.
The list is put into file-bisect-index.toml to be used by other commands.
This command will currently ignore folders.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		idx, err := index.LoadOrCreateIndex()
		if err != nil {
			fmt.Printf("Error loading index file: %v\n", err)
			os.Exit(1)
		}
		idx.Init()
		idx.Refresh()
		err = idx.Save()
		if err != nil {
			fmt.Printf("Error saving index file: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
