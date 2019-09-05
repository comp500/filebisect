package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/comp500/filebisect/index"
	"github.com/spf13/cobra"
)

// ignoreCmd represents the ignore command
var ignoreCmd = &cobra.Command{
	Use:   "ignore [file]",
	Short: "Ignore the specified file",
	Long:  `This command ensures that the file you selected will never be moved. The specified file must be already indexed.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idx, err := index.LoadOrCreateIndex()
		if err != nil {
			fmt.Printf("Error loading index file: %v\n", err)
			os.Exit(1)
		}
		idx.Init()
		idx.Ignore(filepath.Clean(args[0]))
		err = idx.Save()
		if err != nil {
			fmt.Printf("Error saving index file: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(ignoreCmd)
}
