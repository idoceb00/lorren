package main

import (
	"fmt"
	"os"

	"github.com/idoceb00/lorren/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
