package configurations

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPagerDutySuccessfulConfiguration(t *testing.T) {
	os.Setenv("PAGER_DUTY_ROUTING_KEY", "test-routing-key")
	os.Setenv("PAGER_DUTY_RETRY_COUNT", "2")

	config, err := Config("pager_duty")

	assert.Nil(t, err)

	pagerDutyConfig, valid := config.(*PagerDuty)

	assert.True(t, valid)
	assert.Equal(t, "test-routing-key", pagerDutyConfig.RoutingKey)
	assert.Equal(t, 2, pagerDutyConfig.RetryCount)

	os.Setenv("PAGER_DUTY_ROUTING_KEY", "")
	os.Setenv("PAGER_DUTY_RETRY_COUNT", "")
}

func TestPagerDutyFailedConfiguration(t *testing.T) {
	os.Setenv("PAGER_DUTY_ROUTING_KEY", "")

	_, err := Config("pager_duty")

	assert.NotNil(t, err)
}

func TestEmailSuccessfulConfiguration(t *testing.T) {
	os.Setenv("EMAIL_RETRY_COUNT", "2")
	os.Setenv("SMTP_HOST", "smtp_test_host")
	os.Setenv("SMTP_PORT", "90")
	os.Setenv("EMAIL_USERNAME", "test_email_username")
	os.Setenv("EMAIL_PASSWORD", "test_password")

	config, err := Config("email")

	assert.Nil(t, err)

	emailConfig, valid := config.(*Email)

	assert.True(t, valid)

	assert.Equal(t, 2, emailConfig.RetryCount)
	assert.Equal(t, 90, emailConfig.Port)
	assert.Equal(t, "smtp_test_host", emailConfig.Host)
	assert.Equal(t, "test_email_username", emailConfig.Username)
	assert.Equal(t, "test_password", emailConfig.Password)

	os.Setenv("EMAIL_RETRY_COUNT", "")
	os.Setenv("SMTP_HOST", "")
	os.Setenv("SMTP_PORT", "")
	os.Setenv("EMAIL_USERNAME", "")
	os.Setenv("EMAIL_PASSWORD", "")
}

func TestEmailFailedConfiguration(t *testing.T) {
	_, err := Config("email")

	assert.NotNil(t, err)

	os.Setenv("SMTP_HOST", "smtp_test_host")

	_, err = Config("email")

	assert.NotNil(t, err)

	os.Setenv("SMTP_PORT", "99")

	_, err = Config("email")

	assert.NotNil(t, err)

	os.Setenv("EMAIL_USERNAME", "test_email_username")

	_, err = Config("email")

	assert.NotNil(t, err)

	os.Setenv("EMAIL_PASSWORD", "test_password")

	_, err = Config("email")

	assert.Nil(t, err)

	os.Setenv("EMAIL_RETRY_COUNT", "")
	os.Setenv("SMTP_HOST", "")
	os.Setenv("SMTP_PORT", "")
	os.Setenv("EMAIL_USERNAME", "")
	os.Setenv("EMAIL_PASSWORD", "")
}

func TestSlackSuccessfulConfiguration(t *testing.T) {
	os.Setenv("SLACK_WEBHOOKS", "team1@http://test.slack.webhook.com/1,team2@http://test.slack.webhook.com/2")
	os.Setenv("SLACK_RETRY_COUNT", "2")

	config, err := Config("slack")

	assert.Nil(t, err)

	slackConfig, valid := config.(*Slack)

	assert.True(t, valid)
	assert.Equal(t, "http://test.slack.webhook.com/1", slackConfig.Webhooks["team1"])
	assert.Equal(t, "http://test.slack.webhook.com/2", slackConfig.Webhooks["team2"])
	assert.Equal(t, 2, slackConfig.RetryCount)

	os.Setenv("SLACK_WEBHOOKS", "")
	os.Setenv("SLACK_RETRY_COUNT", "")
}

func TestSlackFailedConfiguration(t *testing.T) {
	_, err := Config("slack")

	assert.NotNil(t, err)
}
