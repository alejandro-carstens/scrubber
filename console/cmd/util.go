package cmd

import (
	rp "scrubber/resourcepool"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func Run() error {
	return boot().Execute()
}

func boot() *cobra.Command {
	rp.Boot("mysql")

	scrubber := &cobra.Command{Use: "scrubber"}

	scrubber.PersistentFlags().Int("timeout", 300, "elasticsearch operation timeout")
	scrubber.PersistentFlags().Bool("disable_action", false, "flag for preventing the action to be ran")

	scrubber.AddCommand(new(createIndexCmd).new())
	scrubber.AddCommand(new(deleteIndicesCmd).new())
	scrubber.AddCommand(new(closeIndicesCmd).new())
	scrubber.AddCommand(new(openIndicesCmd).new())
	scrubber.AddCommand(new(aliasCmd).new())
	scrubber.AddCommand(new(createRepositoryCmd).new())
	scrubber.AddCommand(new(snapshotCmd).new())
	scrubber.AddCommand(new(deleteSnaphotsCmd).new())
	scrubber.AddCommand(new(restoreCmd).new())
	scrubber.AddCommand(new(indexSettingsCmd).new())
	scrubber.AddCommand(new(runActionCmd).new())
	scrubber.AddCommand(new(listIndicesCmd).new())
	scrubber.AddCommand(new(listSnapshotsCmd).new())
	scrubber.AddCommand(new(schedulerCmd).new())
	scrubber.AddCommand(new(deleteRepositoriesCmd).new())
	scrubber.AddCommand(new(rolloverCmd).new())
	scrubber.AddCommand(new(httpServe).new())
	scrubber.AddCommand(new(runMigrations).new())

	return scrubber
}

func stringFromFlags(flags *pflag.FlagSet, key string) string {
	value, _ := flags.GetString(key)

	return value
}

func stringSliceFromFlags(flags *pflag.FlagSet, key string) []string {
	values, _ := flags.GetStringSlice(key)

	return values
}

func intFromFlags(flags *pflag.FlagSet, key string) int {
	value, _ := flags.GetInt(key)

	return value
}

func inStringSlice(needle string, haystack []string) bool {
	for _, value := range haystack {
		if value == needle {
			return true
		}
	}

	return false
}
