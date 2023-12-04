package inspect

import (
	"os"

	"github.com/apoloval/scumm-go"
	"github.com/apoloval/scumm-go/vm4"
	"github.com/spf13/cobra"
)

var ScriptCmd = &cobra.Command{
	Use:   "script [index file] [script number]",
	Short: "Inspect a global script",
	Args:  cobra.ExactArgs(2),
	RunE:  func(cmd *cobra.Command, args []string) error { return doInspectScript(args[0], args[1]) },
}

func doInspectScript(indexPath, scriptID string) error {
	rm, err := scumm.FromIndexFile(indexPath)
	if err != nil {
		return err
	}

	id, err := scumm.ParseScriptID(scriptID)
	if err != nil {
		return err
	}
	script, err := rm.GetScript(id, true)

	script.Listing(vm4.DefaultSymbolTable(), os.Stdout)

	return err
}
