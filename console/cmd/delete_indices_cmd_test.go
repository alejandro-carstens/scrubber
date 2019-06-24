package cmd

import (
	"scrubber/logging"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteIndicesCmd(t *testing.T) {
	rootCmd := Init(logging.NewSrvLogger("", true, true, true, true))

	_, err := executeCommand(rootCmd, "delete-indices", "--disable_action=true")

	assert.NotNil(t, err)

	if _, err := executeCommand(rootCmd, "delete-indices", "--indices=my,index,1", "--disable_action=true"); err != nil {
		t.Errorf("%v", err.Error())
	}
}
