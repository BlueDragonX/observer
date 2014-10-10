package config

// Configuration for a single pipe.
type Pipe struct {
	Interval int64
	Sources  []string
	Sinks    []string
	Metadata map[string]string
}

// A mapping of pipes.
type Pipes map[string]Pipe
