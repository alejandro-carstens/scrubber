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

	return data
}
