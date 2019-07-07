package cmd

import (
	"scrubber/actions/contexts"
	"scrubber/actions/options"
	"scrubber/logger"

	"github.com/spf13/cobra"
)

type createRepositoryCmd struct {
	baseActionCmd
}

func (crc *createRepositoryCmd) new(logger *logger.Logger) *cobra.Command {
	command := &cobra.Command{
		Use:   "create-repository",
		Short: "create repository for storing snapshots",
		Args:  crc.Validate,
		Run:   crc.Handle,
	}

	crc.logger = logger

	command.Flags().String("chunk_size", "", "big files can be broken down into chunks for snapshotting if needed")
	command.Flags().String("max_restore_bytes_per_sec", "40mb", "throttles per node restore rate")
	command.Flags().String("max_snapshot_bytes_per_second", "40mb", "throttles per node snapshot rate")
	command.Flags().String("repository", "", "the repository name, it is a required field")
	command.Flags().String("location", "", "location of the snapshots, it is a required field")
	command.Flags().String("repo_type", "", "repository type, ex: fs, gcs, s3")
	command.Flags().Bool("compress", true, "turns on compression for snapshot files")
	command.Flags().Bool("verify", false, "disable the repository verification when registering or updating a repository")

	return command
}

func (crc *createRepositoryCmd) Validate(cmd *cobra.Command, args []string) error {
	crc.context = new(contexts.CreateRepositoryContext)

	options := &options.CreateRepositoryOptions{}

	if err := options.BindFlags(cmd.Flags()); err != nil {
		return err
	}

	return crc.context.SetOptions(options)
}
