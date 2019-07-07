package main

import (
	"os"
	"scrubber/console/cmd"
	"scrubber/logger"
)

func main() {
	logger := logger.NewLogger(os.Getenv("LOG_FILE"), true, true, true, true)

	if err := cmd.Run(logger); err != nil {
		os.Exit(1)
	}
}
