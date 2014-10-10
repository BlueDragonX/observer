package config

import (
	"fmt"
)

// Configuration for a single provider.
type Provider struct {
	Provider string
	Config map[string]interface{}
}

// Marshal a Provider into YAML. 
func (p *Provider) MarshalYAML() (interface{}, error) {
	data := make(map[string]interface{}, len(p.Config) + 1)
	data["provider"] = p.Provider
	for key, value := range p.Config {
		data[key] = value
	}
	return data, nil
}

// Unmarshal YAML into a Provider.
func (p *Provider) UnmarshalYAML(unmarshal func(interface{}) error) error {
	p.Config = make(map[string]interface{})
	unmarshal(p.Config)
	if val, ok := p.Config["provider"]; ok {
		p.Provider = fmt.Sprintf("%v", val)
		delete(p.Config, "provider")
	}
	return nil
}

// A mapping of providers.
type Providers map[string]Provider
