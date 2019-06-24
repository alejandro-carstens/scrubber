package cmd

import (
	"scrubber/logging"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateIndexCmd(t *testing.T) {
	rootCmd := Init(logging.NewSrvLogger("", true, true, true, true))

	_, err := executeCommand(rootCmd, "create-index", "--disable_action=true")

	assert.NotNil(t, err)

	if _, err := executeCommand(rootCmd, "create-index", "--name=my_index", "--disable_action=true"); err != nil {
		t.Errorf("%v", err.Error())
	}
}
