package main

import (
	"os"
	"scrubber/console/cmd"
	"scrubber/logging"
)

func main() {
	if err := cmd.Run(logging.NewSrvLogger("", true, true, true, true)); err != nil {
		os.Exit(1)
	}
}
