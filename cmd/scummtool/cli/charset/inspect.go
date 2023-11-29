package charset

import (
	"fmt"
	"os"
	"unicode"

	"github.com/apoloval/scumm-go"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var InspectCmd = &cobra.Command{
	Use:   "inspect [charset file]",
	Short: "Inspect a SCUMM charset resource file",
	Args:  cobra.ExactArgs(1),
	RunE:  func(cmd *cobra.Command, args []string) error { return inspect(args[0]) },
}

func inspect(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch rt := scumm.DetectResourceFile(file); rt {
	case scumm.ResourceFileCharsetV4:
		charset, err := scumm.DecodeCharsetV4(file)
		if err != nil {
			return err
		}
		return inspectCharset(rt, charset)
	default:
		return fmt.Errorf("invalid input: unexpected %s", rt)
	}
}

func inspectCharset(rt scumm.ResourceFileType, charset scumm.Charset) error {
	fmt.Printf("%s:\n", rt)
	fmt.Printf("  ColorMap     : %v\n", charset.ColorMap)
	fmt.Printf("  BitsPerPixel : %d\n", charset.BitsPerPixel)
	fmt.Printf("  FontHeight   : %d\n", charset.FontHeight)
	fmt.Printf("  Characters   : %d\n", len(charset.Characters))

	chars := table.New("Index", "Symbol", "Width", "Height", "XOffset", "YOffset", "Glyph bytes")
	for i, char := range charset.Characters {
		if char == nil {
			chars.AddRow(i, "", "-", "-", "-", "-", "-")
			continue
		}
		r := rune(i)
		if !unicode.IsGraphic(r) {
			r = ' '
		}
		chars.AddRow(
			i, string(r), char.Width, char.Height, char.XOffset, char.YOffset, len(char.Glyph))
	}
	chars.Print()
	return nil
}
