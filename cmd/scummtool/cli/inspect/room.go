package inspect

import (
	"fmt"

	"github.com/apoloval/scumm-go"
	"github.com/apoloval/scumm-go/vm"
	"github.com/spf13/cobra"
)

var RoomCmd = &cobra.Command{
	Use:   "room [index file] [room number or name]",
	Short: "Inspect the resource of a a SCUMM room",
	Args:  cobra.ExactArgs(2),
	RunE:  func(cmd *cobra.Command, args []string) error { return doInspectRoom(args[0], args[1]) },
}

func doInspectRoom(indexPath, roomNumberOrName string) error {
	rm, err := scumm.FromIndexFile(indexPath)
	if err != nil {
		return err
	}

	room, err := vm.GetRoomFromRef(rm, roomNumberOrName)
	if err != nil {
		return err
	}

	fmt.Printf("Room %d:\n", room.ID)
	fmt.Printf("  Name		: %s\n", room.Name)
	fmt.Printf("  Size		: %dx%d\n", room.Width, room.Height)
	fmt.Printf("  Objects	: %d\n", room.NumberOfObjects)
	fmt.Printf("  Scripts	: %d\n", len(room.LocalScripts))

	return nil
}
