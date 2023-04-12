package postgres

import (
	"fmt"
)

const (
	defaultDriver string = "postgres"
	defaultPort   int    = 5432
)

var ErrRequired = fmt.Errorf("required")

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func (c *Config) Verify() error {
	if c.Host == "" {
		return fmt.Errorf("host: %w", ErrRequired)
	}

	if c.Port == 0 {
		c.Port = defaultPort
	}

	if c.User == "" {
		return fmt.Errorf("user: %w", ErrRequired)
	}

	if c.Password == "" {
		return fmt.Errorf("password: %w", ErrRequired)
	}

	if c.Database == "" {
		return fmt.Errorf("database: %w", ErrRequired)
	}

	return nil
}

func (c Config) URI() (string, error) {
	if err := c.Verify(); err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s:%s@%s:%d/%s", defaultDriver, c.User, c.Password, c.Host, c.Port, c.Database), nil
}
