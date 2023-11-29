package cli

import (
	"github.com/apoloval/scumm-go/cmd/scummtool/cli/charset"
	"github.com/spf13/cobra"
)

var charsetCmd = &cobra.Command{
	Use:   "charset",
	Short: "Manipulate SCUMM charset resources",
}

func init() {
	charsetCmd.AddCommand(charset.InspectCmd)
}
