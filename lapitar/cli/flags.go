package cli

import (
	"fmt"
	"github.com/ogier/pflag"
)

func FlagUsage(name string, flags *pflag.FlagSet) UsageSpec {
	usage := Usage().
		Cmd("Usage: %s [flags] <names>").
		Add("").
		Add("Flags:")

	flags.VisitAll(func(f *pflag.Flag) {
		name := "--" + f.Name
		if len(f.Shorthand) > 0 {
			name += ", -" + f.Shorthand
		}

		desc := fmt.Sprintf("%s (%s)", f.Usage, f.DefValue)
		usage.Help(name, desc)
	})

	flags.Usage = func() {
		usage.Print(name)
	}

	return usage
}
