package main

import (
	"os"

	"github.com/vaihdass/goflu/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
