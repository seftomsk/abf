package config

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNotSetValue = errors.New("not set value")

var ErrInvalidValue = errors.New("invalid value")

var supportedStorages = map[string]struct{}{
	"mem": {},
}

func New() Config {
	return Config{}
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Limiter struct {
	Capacity     int   `json:"capacity"`
	CountSeconds int64 `json:"countSeconds"`
}

type Config struct {
	Storage         string  `json:"storage"`
	Server          Server  `json:"server"`
	LoginLimiter    Limiter `json:"loginLimiter"`
	PasswordLimiter Limiter `json:"passwordLimiter"`
	IPLimiter       Limiter `json:"ipLimiter"`
}

func (c Config) Host() string {
	return c.Server.Host
}

func (c Config) Port() string {
	return c.Server.Port
}

func validateSetValues(c Config) error {
	if c.Storage == "" {
		return fmt.Errorf(
			"validateSetValues - storage: %w", ErrNotSetValue)
	}
	if c.Server.Host == "" {
		return fmt.Errorf(
			"validateSetValues - host: %w", ErrNotSetValue)
	}
	if c.Server.Port == "" {
		return fmt.Errorf(
			"validateSetValues - port: %w", ErrNotSetValue)
	}

	return nil
}

func validateSupportedValues(c Config) error {
	if _, ok := supportedStorages[strings.ToLower(c.Storage)]; !ok {
		return fmt.Errorf(
			"validateSupportedValues - storage: %w", ErrInvalidValue)
	}

	if c.LoginLimiter.Capacity < 0 {
		return fmt.Errorf(
			"validateSupportedValues - loginLimiter - capacity: %w",
			ErrInvalidValue)
	}

	if c.LoginLimiter.CountSeconds < 0 {
		return fmt.Errorf(
			"validateSupportedValues - loginLimiter - countSeconds: %w",
			ErrInvalidValue)
	}

	if c.PasswordLimiter.Capacity < 0 {
		return fmt.Errorf(
			"validateSupportedValues - passwordLimiter - capacity: %w",
			ErrInvalidValue)
	}

	if c.PasswordLimiter.CountSeconds < 0 {
		return fmt.Errorf(
			"validateSupportedValues - "+
				"passwordLimiter - countSeconds: %w",
			ErrInvalidValue)
	}

	if c.IPLimiter.Capacity < 0 {
		return fmt.Errorf(
			"validateSupportedValues - ipLimiter - capacity: %w",
			ErrInvalidValue)
	}

	if c.IPLimiter.CountSeconds < 0 {
		return fmt.Errorf(
			"validateSupportedValues - ipLimiter - countSeconds: %w",
			ErrInvalidValue)
	}

	return nil
}

func (c Config) Validate() error {
	err := validateSetValues(c)
	if err != nil {
		return fmt.Errorf("validate - validateSetValies: %w", err)
	}

	err = validateSupportedValues(c)
	if err != nil {
		return fmt.Errorf("validate - validateSupportedValues: %w", err)
	}

	return nil
}
