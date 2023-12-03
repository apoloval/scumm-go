package inspect

import (
	"fmt"

	"github.com/apoloval/scumm-go"
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
	script, err := rm.GetScript(id)
	if err != nil {
		return err
	}

	fmt.Printf("Script %d: %d bytes\n", script.ID, len(script.Bytecode))

	return nil
}
