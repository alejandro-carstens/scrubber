package main

import (
	"os"
	"scrubber/console/cmd"
)

func main() {
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
}
