package cli

import (
	"github.com/apoloval/scumm-go"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [index file]",
	Short: "Run a SCUMM application from its index file",
	Args:  cobra.ExactArgs(1),
	RunE:  func(cmd *cobra.Command, args []string) error { return doRun(args[0]) },
}

func doRun(path string) error {
	return scumm.Run(path)
}
