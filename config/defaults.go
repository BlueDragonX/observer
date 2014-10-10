package config

const (
	DefaultInterval = 300
)

// Defaults config object.
type Defaults struct {
	Interval int64
	Metadata map[string]string
}

// Create a new Defaults object with default values populated.
func NewDefaults() Defaults {
	return Defaults{Interval: DefaultInterval}
}
