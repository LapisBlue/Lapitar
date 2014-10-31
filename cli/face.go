package cli

import (
	"fmt"
	"github.com/LapisBlue/Tar/face"
	"github.com/LapisBlue/Tar/skin"
	"github.com/LapisBlue/Tar/util"
	"github.com/ogier/pflag"
	"image"
	"image/png"
	"os"
)

const (
	faceSize = 256
)

func runFace(name string, args []string) int {
	flags := pflag.NewFlagSet(name, pflag.ContinueOnError)

	size := flags.IntP("size", "s", faceSize, "The size of the avatar, in pixels.")
	helm := flags.Bool("helm", false, "Render the helm as additional layer to the face")
	in := flags.StringP("in", "i", input, "The source of the list of players to render. Can be either a file, STDIN or ARGS.")
	out := flags.StringP("out", "o", output, "The destination to write the result to. Can be either a file or STDOUT.")

	FlagUsage(name, flags).
		Add("").
		Add("Example:") // TODO

	if len(args) < 1 || args[0] == "help" {
		flags.Usage()
		return 1
	}

	watch := util.GlobalWatch().Start()
	if flags.Parse(args) != nil {
		return 1
	}

	players := readFrom(*in, flags.Args())
	if players == nil {
		return 1
	}

	if isStdout(*out) {
		if len(players) > 1 {
			fmt.Fprintln(os.Stderr, "You can only render one image using STDOUT")
			return 1
		}

		player := players[0]
		skin, err := skin.Download(player)
		if err != nil {
			return PrintError(err, "Failed to download skin:", player)
		}

		face := face.Render(skin, *size, *helm)

		err = png.Encode(os.Stdout, face)
		if err != nil {
			return PrintError(err, "Failed to write face to STDOUT")
		}

		return 0
	}

	skins := downloadSkins(players)

	fmt.Println()
	fmt.Printf("Rendering %d face(es), please wait...\n", len(skins))

	watch.Mark()
	faces := make([]image.Image, len(skins))

	for i, skin := range skins {
		if skin == nil {
			continue
		}

		watch.Mark()

		faces[i] = face.Render(skin, *size, *helm)
		fmt.Println("Rendered face:", players[i], watch)
	}

	fmt.Println("Finished rendering faces", watch)

	fmt.Println()
	saveResults(players, faces, *out)

	fmt.Println()
	watch.Stop()
	fmt.Println("Done!", watch)

	return 0
}
