package cli

import (
	"fmt"
	"os"

	"github.com/apoloval/scumm-go"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [file]",
	Short: "Inspect a SCUMM resource file",
	Args:  cobra.ExactArgs(1),
	RunE:  func(cmd *cobra.Command, args []string) error { return inspect(args[0]) },
}

var inspectFlags struct {
	resourceType string
	room         int
}

func inspect(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	index, err := scumm.DecodeIndexFile(file)
	if err != nil {
		return err
	}

	return inspectIndex(index)
}

func inspectIndex(index scumm.Index) error {
	fmt.Printf("Index file v4\n")
	fmt.Printf("Found %d rooms\n", len(index.Rooms))

	index.VisitRooms(func(num scumm.RoomNumber, room *scumm.RoomIndex) {
		if inspectFlags.room == -1 || int(num) == inspectFlags.room {
			printIndexRoom(num, room)
		}
	})
	return nil
}

func printIndexRoom(id scumm.RoomNumber, room *scumm.RoomIndex) {
	fmt.Printf("Room %d:\n", id)
	fmt.Printf("  Name      : %s\n", room.Name)
	fmt.Printf("  Data file : %d\n", room.FileNumber)
	fmt.Printf("  Offset    : 0x%05x\n", room.FileOffset)
	fmt.Printf("  Scripts   : %d\n", len(room.ScriptOffsets))
	fmt.Printf("  Sounds    : %d\n", len(room.SoundOffsets))
	fmt.Printf("  Costumes  : %d\n", len(room.CostumeOffsets))
}

func init() {
	rootCmd.AddCommand(inspectCmd)
	inspectCmd.Flags().StringVarP(&inspectFlags.resourceType, "type", "t", "auto", "resource type")
	inspectCmd.Flags().IntVarP(&inspectFlags.room, "room", "r", -1, "room number")
}
