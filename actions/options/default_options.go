package options

import (
	"fmt"

	"github.com/spf13/pflag"

	"github.com/Jeffail/gabs"
)

type defaultOptions struct {
	container     *gabs.Container
	Timeout       int  `json:"timeout_override"`
	DisableAction bool `json:"disable_action"`
}

func (do *defaultOptions) GetContainer() *gabs.Container {
	return do.container
}

func (do *defaultOptions) Exists(value string) bool {
	return do.container.Exists(value)
}

func (do *defaultOptions) GetDisableAction() bool {
	return do.DisableAction
}

func (do *defaultOptions) TimeoutInSeconds() string {
	if do.Timeout > 0 {
		return fmt.Sprintf("%vs", do.Timeout)
	}

	return ""
}

func (do *defaultOptions) Get(value string) interface{} {
	return do.container.S(value).Data()
}

func (do *defaultOptions) String(value string) string {
	return fmt.Sprint(do.Get(value))
}

func (do *defaultOptions) IsSnapshot() bool {
	return false
}

func (do *defaultOptions) defaultBindFlags(flags *pflag.FlagSet) error {
	timeout, _ := flags.GetInt("timeout")
	disableAction, _ := flags.GetBool("disable_action")

	do.Timeout = timeout
	do.DisableAction = disableAction

	return nil
}
