package cli

import (
	"fmt"
	"os"

	"github.com/apoloval/scumm-go"
	"github.com/apoloval/scumm-go/cmd/scummtool/cli/inspect"
	"github.com/apoloval/scumm-go/collections"
	"github.com/apoloval/scumm-go/vm"
	"github.com/apoloval/scumm-go/vm4"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect [index file]",
	Short: "Inspect a SCUMM index file",
	Args:  cobra.ExactArgs(1),
	RunE:  func(cmd *cobra.Command, args []string) error { return doInspect(args[0]) },
}

var inspectFlags struct {
	showRooms    bool
	showScripts  bool
	showSounds   bool
	showCostumes bool
	showObjects  bool
}

func doInspect(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch rt := scumm.DetectResourceFile(file); rt {
	case vm4.ResourceFileIndex:
		index, err := vm4.DecodeIndex(file)
		if err != nil {
			return err
		}
		return inspectIndex(rt, index)
	default:
		return fmt.Errorf("invalid input: unexpected %s", rt)
	}
}

func inspectIndex(rt vm.ResourceFileType, index vm.Index) error {
	fmt.Printf("%s:\n", rt)
	fmt.Printf("  Rooms    : %d\n", len(index.Rooms))
	fmt.Printf("  Scripts  : %d\n", len(index.Scripts))
	fmt.Printf("  Sounds   : %d\n", len(index.Sounds))
	fmt.Printf("  Costumes : %d\n", len(index.Costumes))
	fmt.Printf("  Objects  : %d\n", len(index.Objects))
	println()

	if inspectFlags.showRooms {
		fmt.Printf("Directory of rooms:\n")
		rooms := table.New("ID", "Name", "File", "Offset", "Scripts", "Sounds", "Costumes")
		collections.VisitMap(index.Rooms, func(id vm.RoomID, room vm.IndexedRoom) {
			rooms.AddRow(
				room.ID, room.Name, room.FileNumber, room.FileOffset,
				len(room.Scripts), len(room.Sounds), len(room.Costumes))
		})
		rooms.Print()
		println()
	}

	if inspectFlags.showScripts {
		fmt.Printf("Directory of scripts:\n")
		scripts := table.New("ID", "File", "Room", "Offset")
		collections.VisitMap(index.Scripts, func(id vm.ScriptID, script vm.IndexedScript) {
			scripts.AddRow(
				script.ID, index.Rooms[script.Room].FileNumber, script.Room, script.Offset)
		})
		scripts.Print()
		println()
	}

	if inspectFlags.showSounds {
		fmt.Printf("Directory of sounds:\n")
		sounds := table.New("ID", "File", "Room", "Offset")
		collections.VisitMap(index.Sounds, func(id vm.SoundID, sound vm.IndexedSound) {
			sounds.AddRow(sound.ID, index.Rooms[sound.Room].FileNumber, sound.Room, sound.Offset)
		})
		sounds.Print()
		println()
	}

	if inspectFlags.showCostumes {
		fmt.Printf("Directory of costumes:\n")
		costumes := table.New("ID", "File", "Room", "Offset")
		collections.VisitMap(index.Costumes, func(id vm.CostumeID, costume vm.IndexedCostume) {
			costumes.AddRow(
				costume.ID, index.Rooms[costume.Room].FileNumber, costume.Room, costume.Offset)
		})
		costumes.Print()
		println()
	}

	if inspectFlags.showObjects {
		fmt.Printf("Directory of objects:\n")
		objects := table.New("ID", "Class", "Owner", "State")
		collections.VisitMap(index.Objects, func(id vm.ObjectID, object vm.IndexedObject) {
			objects.AddRow(id, object.Class, object.Owner, object.State)
		})
		objects.Print()
		println()
	}

	return nil
}

func init() {
	inspectCmd.Flags().BoolVarP(&inspectFlags.showRooms,
		"rooms", "r", true, "show rooms")
	inspectCmd.Flags().BoolVarP(&inspectFlags.showScripts,
		"scripts", "s", false, "show scripts")
	inspectCmd.Flags().BoolVarP(&inspectFlags.showSounds,
		"sounds", "n", false, "show sounds")
	inspectCmd.Flags().BoolVarP(&inspectFlags.showCostumes,
		"costumes", "c", false, "show costumes")
	inspectCmd.Flags().BoolVarP(&inspectFlags.showObjects,
		"objects", "o", false, "show objects")

	inspectCmd.AddCommand(inspect.RoomCmd)
	inspectCmd.AddCommand(inspect.ScriptCmd)
}
