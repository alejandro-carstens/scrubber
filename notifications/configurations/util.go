package configurations

import "errors"

// Config returns a Configuration for a given channel
func Config(notificationType string) (Configurable, error) {
	var config Configurable

	switch notificationType {
	case "slack":
		config = &Slack{}
		break
	case "pager_duty":
		config = &PagerDuty{}
		break
	default:
		return nil, errors.New("invalid configuration type")
	}

	return config.FillFromEnvs().Validate()
}
