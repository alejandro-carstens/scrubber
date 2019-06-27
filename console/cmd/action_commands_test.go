package cmd

import (
	"bytes"
	"scrubber/logging"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestActionCmds(t *testing.T) {
	for _, data := range cmdParamsDataProvider() {
		rootCmd := Init(logging.NewSrvLogger("", true, true, true, true))

		_, err := executeCommand(rootCmd, data["error_params"]...)

		assert.NotNil(t, err)

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

	return data
}