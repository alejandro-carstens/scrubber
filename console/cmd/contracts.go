package cmd

import (
	"github.com/spf13/cobra"
)

type Commandable interface {
	new() *cobra.Command

	Validate(cmd *cobra.Command, args []string) error

	Handle(cmd *cobra.Command, args []string)
}
