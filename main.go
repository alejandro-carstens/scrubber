package main

import (
	"os"
	"scrubber/console/cmd"
	"scrubber/logging"
)

func main() {
	logger := logging.NewSrvLogger(os.Getenv("LOG_FILE"), true, true, true, true)

	if err := cmd.Run(logger); err != nil {
		os.Exit(1)
	}
}
