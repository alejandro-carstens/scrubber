package config

import (
	"errors"
	"os"
)

// MySQL represents a mysql
// client configuration
type MySQL struct {
	Password           string
	Username           string
	Database           string
	Host               string
	MaxIdleConnections int
	MaxOpenConnections int
	ConnMaxLifetime    int
	Port               int
}

// Validate validates the configuration for a given channel
func (m *MySQL) Validate() (Configurable, error) {
	if len(m.Username) == 0 {
		return nil, errors.New("mysql username cannot be empty")
	}

	if len(m.Database) == 0 {
		return nil, errors.New("mysql database cannot be empty")
	}

	if len(m.Host) == 0 {
		return nil, errors.New("mysql host cannot be empty")
	}

	if m.Port == 0 {
		return nil, errors.New("mysql port cannot be 0")
	}

	return m, nil
}

// FillFromEnvs is responsible for setting the configuration
// for the channel from the respective env variables
func (m *MySQL) FillFromEnvs() Configurable {
	m.Host = os.Getenv("MYSQL_HOST")
	m.Username = os.Getenv("MYSQL_USERNAME")
	m.Database = os.Getenv("MYSQL_DATABASE")
	m.Password = os.Getenv("MYSQL_PASSWORD")
	m.MaxIdleConnections = stringToInt(os.Getenv("MYSQL_MAX_IDLE_CONNECTIONS"))
	m.MaxOpenConnections = stringToInt(os.Getenv("MYSQL_MAX_OPEN_CONNECTIONS"))
	m.ConnMaxLifetime = stringToInt(os.Getenv(")MYSQL_CONNECTION_MAX_LIFE_TIME"))
	m.Port = stringToInt(os.Getenv("MYSQL_PORT"))

	return m
}
