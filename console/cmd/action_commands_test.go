package cmd

import (
	"bytes"
	"scrubber/logger"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestActionCommands(t *testing.T) {
	for _, data := range cmdParamsDataProvider() {
		rootCmd := boot(logger.NewLogger("", true, true, true, true))

		if len(data["error_params"]) > 0 {
			_, err := executeCommand(rootCmd, data["error_params"]...)

			assert.NotNil(t, err)
		}

		if _, err := executeCommand(rootCmd, data["success_params"]...); err != nil {
			t.Errorf("%v", err.Error())
		}
	}
}

func executeCommand(root *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)

	root.SetOutput(buf)

	root.SetArgs(args)

	_, err := root.ExecuteC()

	return buf.String(), err
}

func cmdParamsDataProvider() []map[string][]string {
	data := []map[string][]string{}

	data = append(data, map[string][]string{
		"error_params": []string{
			"delete-indices",
			"--disable_action=true",
		},
		"success_params": []string{
			"delete-indices",
			"--indices=my,index,1",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"close-indices",
			"--disable_action=true",
		},
		"success_params": []string{
			"close-indices",
			"--indices=my,index,1",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"open-indices",
			"--disable_action=true",
		},
		"success_params": []string{
			"open-indices",
			"--indices=my,index,1",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"alias",
			"--disable_action=true",
		},
		"success_params": []string{
			"alias",
			"--indices=my,index,1",
			"--name=alias",
			"--type=add",
			"--routing=1",
			"--search_routing=2",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"create-index",
			"--disable_action=true",
		},
		"success_params": []string{
			"create-index",
			"--name=my_index",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"create-repository",
			"--chunk_size=10m",
			"--max_restore_bytes_per_sec=30mb",
			"--max_snapshot_bytes_per_second=30mb",
			"--repo_type=fs",
			"--compress=true",
			"--verify=false",
			"--disable_action=true",
		},
		"success_params": []string{
			"create-repository",
			"--repository=whatever",
			"--location=whatever",
			"--chunk_size=10m",
			"--max_restore_bytes_per_sec=30mb",
			"--max_snapshot_bytes_per_second=30mb",
			"--repo_type=fs",
			"--compress=true",
			"--verify=false",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"snapshot",
			"--indices=index_1,index_2,index3",
			"--name=snapshot",
			"--ignore_unavailable=true",
			"--include_global_state=true",
			"--partial=true",
			"--wait_for_completion=true",
			"--max_wait=9",
			"--wait_interval=30",
			"--disable_action=true",
		},
		"success_params": []string{
			"snapshot",
			"--indices=index_1,index_2,index3",
			"--repository=my_repo",
			"--name=snapshot",
			"--ignore_unavailable=true",
			"--include_global_state=true",
			"--partial=true",
			"--wait_for_completion=true",
			"--max_wait=9",
			"--wait_interval=30",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"delete-snapshots",
			"--disable_action=true",
			"--snapshots=snapshot_1,snapshot_2,snapshot_3",
			"--retry_count=5",
			"--retry_interval=30",
		},
		"success_params": []string{
			"delete-snapshots",
			"--disable_action=true",
			"--repository=my_repository",
			"--snapshots=snapshot_1,snapshot_2,snapshot_3",
			"--retry_count=5",
			"--retry_interval=30",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"restore",
			"--name=my_snapshot",
			"--indices=my_index_1,my_index_2",
			"--rename_pattern=regex",
			"--rename_replacement=regex",
			"--max_wait=20",
			"--wait_interval=2",
			"--ignore_unavailable=false",
			"--include_global_state=true",
			"--partial=true",
			"--wait_for_completion=false",
			"--include_aliases=true",
			"--disable_action=true",
		},
		"success_params": []string{
			"restore",
			"--repository=my_repo",
			"--name=my_snapshot",
			"--indices=my_index_1,my_index_2",
			"--rename_pattern=regex",
			"--rename_replacement=regex",
			"--max_wait=20",
			"--wait_interval=2",
			"--ignore_unavailable=false",
			"--include_global_state=true",
			"--partial=true",
			"--wait_for_completion=false",
			"--include_aliases=true",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"index-settings",
			"--disable_action=true",
			"--indices=index_1,index_2,index_3",
		},
		"success_params": []string{
			"index-settings",
			"--disable_action=true",
			`--index_settings={"foo": "bar"}`,
			"--indices=index_1,index_2,index_3",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"run-action",
		},
		"success_params": []string{
			"run-action",
			"--file_path=/../../tests/testfiles/disable_index_create.yml",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{},
		"success_params": []string{
			"list-indices",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"list-snapshots",
		},
		"success_params": []string{
			"list-snapshots",
			"--repository=my_repository",
			"--disable_action=true",
		},
	})
	data = append(data, map[string][]string{
		"error_params": []string{
			"delete-repositories",
			"--disable_action=true",
		},
		"success_params": []string{
			"delete-repositories",
			"--repositories=my_repository,other_repository",
			"--disable_action=true",
		},
	})

	return data
}
