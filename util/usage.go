package util

import (
	"fmt"
	"os"
	"strings"
)

func Usage() UsageSpec {
	return &usage{make([]interface{}, 0, 8)}
}

type UsageSpec interface {
	Add(line interface{}) UsageSpec
	Cmd(line string) UsageSpec
	Help(cmd, desc string) UsageSpec
	Print(name string) int
}

type usage struct {
	lines []interface{}
}

type cmdString string

type help struct {
	usage *usage
	cmd   []string
	desc  []string
	max   int
}

func (usage *usage) Add(line interface{}) UsageSpec {
	usage.lines = append(usage.lines, line)
	return usage
}

func (usage *usage) Cmd(format string) UsageSpec {
	return usage.Add(cmdString(format))
}

func (usage *usage) Print(name string) int {
	for _, line := range usage.lines {
		if format, ok := line.(cmdString); ok {
			fmt.Fprintf(os.Stderr, string(format)+"\n", name)
		} else {
			fmt.Fprintln(os.Stderr, line)
		}
	}

	return 1
}

func (usage *usage) Help(cmd, desc string) UsageSpec {
	help := &help{
		usage,
		make([]string, 0, 4),
		make([]string, 0, 4),
		0,
	}

	return help.Help(cmd, desc)
}

func (help *help) Help(cmd, desc string) UsageSpec {
	help.cmd = append(help.cmd, cmd)
	help.desc = append(help.desc, desc)
	if len(cmd) > help.max {
		help.max = len(cmd)
	}

	return help
}

func (help *help) build() UsageSpec {
	for i := 0; i < len(help.cmd); i++ {
		cmd, desc := help.cmd[i], help.desc[i]
		help.usage.Add(fmt.Sprintf("   %s%s   %s", cmd, strings.Repeat(" ", help.max-len(cmd)), desc))
	}

	return help.usage
}

func (help *help) Add(line interface{}) UsageSpec {
	return help.build().Add(line)
}

func (help *help) Cmd(format string) UsageSpec {
	return help.build().Cmd(format)
}

func (help *help) Print(name string) int {
	return help.build().Print(name)
}
