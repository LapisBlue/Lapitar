package cli

import (
	"fmt"
	"os"
	"strings"
)

func Usage() UsageSpec {
	return &usagespec{
		lines: make([]interface{}, 0, 8),
	}
}

type UsageSpec interface {
	Add(line interface{}) UsageSpec
	Cmd(line string) UsageSpec
	Help(cmd, desc string) UsageSpec
	Print(name string) int
}

type usagespec struct {
	lines []interface{}
	help  *helpspec
}

type cmdString string

type helpspec struct {
	usage *usagespec
	cmd   []string
	desc  []string
	max   int
}

func (usage *usagespec) checkHelp() {
	if usage.help != nil {
		usage.help.build()
	}
}

func (usage *usagespec) Add(line interface{}) UsageSpec {
	usage.checkHelp()
	usage.lines = append(usage.lines, line)
	return usage
}

func (usage *usagespec) Cmd(format string) UsageSpec {
	return usage.Add(cmdString(format))
}

func (usage *usagespec) Print(name string) int {
	usage.checkHelp()
	for _, line := range usage.lines {
		if format, ok := line.(cmdString); ok {
			fmt.Fprintf(os.Stderr, string(format)+"\n", name)
		} else {
			fmt.Fprintln(os.Stderr, line)
		}
	}

	return 1
}

func (usage *usagespec) Help(cmd, desc string) UsageSpec {
	if usage.help == nil {
		usage.help = &helpspec{
			usage,
			make([]string, 0, 4),
			make([]string, 0, 4),
			0,
		}
	}

	return usage.help.Help(cmd, desc)
}

func (help *helpspec) Help(cmd, desc string) UsageSpec {
	help.cmd = append(help.cmd, cmd)
	help.desc = append(help.desc, desc)
	if len(cmd) > help.max {
		help.max = len(cmd)
	}

	return help
}

func (help *helpspec) build() UsageSpec {
	help.usage.help = nil

	for i := 0; i < len(help.cmd); i++ {
		cmd, desc := help.cmd[i], help.desc[i]
		help.usage.Add(fmt.Sprintf("   %s%s   %s", cmd, strings.Repeat(" ", help.max-len(cmd)), desc))
	}

	return help.usage
}

func (help *helpspec) Add(line interface{}) UsageSpec {
	return help.usage.Add(line)
}

func (help *helpspec) Cmd(format string) UsageSpec {
	return help.usage.Cmd(format)
}

func (help *helpspec) Print(name string) int {
	return help.usage.Print(name)
}
