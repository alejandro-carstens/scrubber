package main

import (
	"os"

	"github.com/alejandro-carstens/scrubber/console/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
