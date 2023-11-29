package cli

import (
	"github.com/apoloval/scumm-go/cmd/scummtool/cli/index"
	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Manipulate SCUMM index resources",
}

func init() {
	indexCmd.AddCommand(index.InspectCmd)
}
