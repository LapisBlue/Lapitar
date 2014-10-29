package main

import (
	"fmt"
	"github.com/LapisBlue/Tar/cli"
	"github.com/LapisBlue/Tar/util"
	"github.com/LapisBlue/Tar/web"
	"os"
	"path/filepath"
)

var usage = util.Usage().
	Cmd("Usage: %s <command> [args]").
	Add("").
	Add("Commands:").
	Help("create <type> [flags] <names...>", "Render an avatar using the command line.").
	Help("server <address> [flags]", "Start the webserver.").
	Help("help [command]", "Display this help page or more information about another command.").
	Add("").
	Cmd("Type '%s help [command]' for more information about a command.")

func Run(name string, args []string) int {
	if len(args) < 1 {
		return usage.Print(name)
	}

	command := args[0]
	if command == "help" {
		if len(args) < 2 {
			return usage.Print(name)
		}

		command = args[1]
		args = []string{"help"}
	}

	switch command {
	case "create":
		return cli.Run(name+" "+command, args[1:])
	case "server":
		return web.Run(name+" "+command, args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: '%s'\n", command)
		return usage.Print(name)
	}
}

func main() {
	os.Exit(Run(filepath.Base(os.Args[0]), os.Args[1:]))
}
