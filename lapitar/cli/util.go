package cli

import (
	"fmt"
	"os"
)

func PrintError(err error, description ...interface{}) int {
	if err == nil {
		return 0
	}

	fmt.Fprintln(os.Stderr, description...)
	fmt.Fprintln(os.Stderr, err)
	return 1
}
