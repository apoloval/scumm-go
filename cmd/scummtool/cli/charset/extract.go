package charset

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	"github.com/apoloval/scumm-go"
	"github.com/apoloval/scumm-go/vm"
	"github.com/apoloval/scumm-go/vm4"
	"github.com/spf13/cobra"
)

var ExtractCmd = &cobra.Command{
	Use:   "extract [charset file]",
	Short: "Extract a SCUMM charset into an image file",
	Args:  cobra.ExactArgs(1),
	RunE:  func(cmd *cobra.Command, args []string) error { return extract(args[0]) },
}

var extractFlags struct {
	Output          string
	BackgroundColor int
}

func extract(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch rt := scumm.DetectResourceFile(file); rt {
	case vm4.ResourceFileCharset:
		charset, err := vm4.DecodeCharset(file)
		if err != nil {
			return err
		}
		return extractCharset(rt, charset)
	default:
		return fmt.Errorf("invalid input: unexpected %s", rt)
	}
}

func extractCharset(rt vm.ResourceFileType, charset vm.Charset) error {
	var width int
	for r := rune(0); r < 256; r++ {
		width += charset.CharWidth(r)
	}
	canvas := image.NewPaletted(
		image.Rect(0, 0, width, int(charset.FontHeight)),
		scumm.ColorPaletteEGA,
	)

	bgColor := image.NewUniform(scumm.ColorPaletteEGA[extractFlags.BackgroundColor])
	draw.Draw(canvas, canvas.Bounds(), bgColor, image.Point{}, draw.Src)

	var x int
	for r := rune(0); r < 256; r++ {
		x += charset.PrintChar(r, canvas, image.Pt(x, 0))
	}

	output, err := os.Create(extractFlags.Output)
	if err != nil {
		return err
	}
	defer output.Close()

	return png.Encode(output, canvas)
}

func init() {
	ExtractCmd.Flags().StringVarP(&extractFlags.Output,
		"output", "o", "charset.png", "output file")
	ExtractCmd.Flags().IntVarP(&extractFlags.BackgroundColor,
		"background-color", "c", 5, "background color from the EGA palette")

}
