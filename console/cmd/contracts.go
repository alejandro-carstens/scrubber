package cmd

import "github.com/spf13/cobra"

// Commandable represents the contract a command
// must implement in order to be executed
type Commandable interface {
	new() *cobra.Command

	// Validate validates the command parameters
	Validate(cmd *cobra.Command, args []string) error

	// Handle executes the action of the command
	Handle(cmd *cobra.Command, args []string)
}
