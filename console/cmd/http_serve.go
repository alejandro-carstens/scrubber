package cmd

import (
	"errors"

	"scrubber/app"
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
)

type httpServe struct {
	port int
}

func (hs *httpServe) new() *cobra.Command {
	command := &cobra.Command{
		Use:   "http-serve",
		Short: "listens for http requests and serves them",
		Args:  hs.Validate,
		Run:   hs.Handle,
	}

	command.Flags().Int("port", 8081, "name of the index to be created")

	return command
}

func (hs *httpServe) Handle(cmd *cobra.Command, args []string) {
	if err := rp.BootResource("mysql"); err != nil {
		rp.Logger().Fatalf(err.Error())
	}

	if err := app.Run(hs.port); err != nil {
		rp.Logger().Fatalf(err.Error())
	}
}

func (hs *httpServe) Validate(cmd *cobra.Command, args []string) error {
	port := intFromFlags(cmd.Flags(), "port")

	if port < 0 || port > 99999 {
		return errors.New("port is out of range")
	}

	hs.port = port

	return nil
}
