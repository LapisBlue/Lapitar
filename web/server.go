package web

import (
	"fmt"
	"github.com/LapisBlue/Tar/cli"
)

var usage = cli.Usage().
	Cmd("Usage: %s <address> [flags]").
	Add("").
	Add("Flags:").  // TODO
	Add("Example:") // TODO

func Run(name string, args []string) int {
	if len(args) < 1 || args[0] == "help" {
		return usage.Print(name)
	}

	address := args[0]
	fmt.Println(address)
	// TODO
	return 0
}
