package cli

import (
	"fmt"
	"os"
)

var (
	usage = Usage().
		Cmd("Usage: %s <type> [flags] <names...>").
		Add("").
		Add("Types:").
		Help("face", "Render the face of a skin.").
		Help("head", "Render the 3D head of a skin.").
		Add("").
		Cmd("Type '%s help [type]' for more information about the available flags for an image type.")
)

func Run(name string, args []string) int {
	if len(args) < 1 {
		return usage.Print(name)
	}

	render := args[0]
	if render == "help" {
		if len(args) < 2 {
			return usage.Print(name)
		}

		render = args[1]
		args = []string{"help"}
	}

	switch render {
	case "face":
		return runFace(name+" "+render, args[1:])
	case "head":
		return runHead(name+" "+render, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown type: '%s'\n", render)
		return usage.Print(name)
	}
}
