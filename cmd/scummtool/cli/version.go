package cli

import (
	"fmt"

	"github.com/apoloval/scumm-go"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of SCUMM tool",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("SCUMM Go v%s\n", scumm.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
