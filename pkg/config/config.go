package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Config struct {
	MQTT    *MQTT    `json:"mqtt"`
	Fronius *Fronius `json:"fronius"`
}

func (c *Config) Verify() error {
	if c.MQTT == nil {
		return errors.New("mqtt must be set")
	}
	if err := c.MQTT.Verify(); err != nil {
		return err
	}

	if c.Fronius == nil {
		return errors.New("fronius must be set")
	}
	if err := c.Fronius.Verify(); err != nil {
		return err
	}

	return nil
}

type MQTT struct {
	Broker   string `json:"broker"`
	Username string `json:"username"`
	Password string `json:"password"`
	Topic    string `json:"topic"`
}

func (m *MQTT) Verify() error {
	if m.Broker == "" {
		return errors.New("broker needs to be set")
	}
	if m.Username == "" {
		return errors.New("username needs to be set")
	}
	if m.Topic == "" {
		return errors.New("topic needs to be set")
	}
	return nil
}

type Fronius struct {
	Backend  string `json:"backend"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (f *Fronius) Verify() error {
	if f.Backend == "" {
		return errors.New("backend needs to be set")
	}
	if f.Username == "" {
		return errors.New("username needs to be set")
	}
	if f.Password == "" {
		return errors.New("password needs to be set")
	}
	return nil
}

func ReadFromFile(file string) (*Config, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read config from %q: %s", file, err)
	}
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("unable to parse JSON in %q: %s", file, err)
	}
	if err := config.Verify(); err != nil {
		return nil, err
	}
	return config, nil
}
