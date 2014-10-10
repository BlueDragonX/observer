package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// Top level configuration object.
type Config struct {
	Defaults Defaults
	Sources  Providers
	Sinks    Providers
	Pipes    Pipes
}

// Create a new configuration object populated with defaults.
func NewConfig() Config {
	return Config{
		Defaults: NewDefaults(),
	}
}

// Read a config file.
func Read(path string) (*Config, error) {
	var err error
	var data []byte
	config := NewConfig()
	if data, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	pipes := make(Pipes, len(config.Pipes))
	for name, pipe := range config.Pipes {
		if pipe.Interval <= 0 {
			pipe.Interval = config.Defaults.Interval
		}
		if pipe.Metadata == nil {
			pipe.Metadata = make(map[string]string, len(config.Defaults.Metadata))
		}
		for key, value := range config.Defaults.Metadata {
			if _, ok := pipe.Metadata[key]; !ok {
				pipe.Metadata[key] = value
			}
		}
		pipes[name] = pipe
	}
	config.Pipes = pipes
	return &config, nil
}
