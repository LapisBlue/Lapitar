package cli

import (
	"fmt"
	"github.com/LapisBlue/Tar/util"
)

var usage = util.Usage().
	Cmd("Usage: %s <type> [flags] <names...>").
	Add("").
	Add("Types:").
	Help("face", "Render the face of a skin.").
	Help("head", "Render the 3D head of a skin.").
	Add("").
	Add("Flags:").  // TODO
	Add("Example:") // TODO

func Run(name string, args []string) int {
	if len(args) < 1 || args[0] == "help" {
		return usage.Print(name)
	}

	render := args[0]
	fmt.Println(render)
	// TODO
	return 0
}
