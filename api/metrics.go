package api

import (
	"encoding/json"
	"time"
)

const (
	UNIT_BYTES = "bytes"
	UNIT_BYTES_PER_SECOND = "bytes/second"
	UNIT_BITS = "bits"
	UNIT_BITS_PER_SECOND = "bits/second"
	UNIT_SECONDS = "seconds"
	UNIT_COUNT = "count"
	UNIT_COUNT_PER_SECOND = "count/second"
	UNIT_PERCENT = "percent"
)

type Metric struct {
	Name      string
	Value     float64
	Unit      string
	Timestamp time.Time
	Metadata  map[string]string
}

// Add metadata values to the metric if they do not already exist.
func (m *Metric) Underlay(metadata map[string]string) {
	if m.Metadata == nil {
		m.Metadata = make(map[string]string, len(metadata))
	}
	for key, value := range metadata {
		if _, ok := m.Metadata[key]; !ok {
			m.Metadata[key] = value
		}
	}
}

// A collection of metric values.
type Metrics struct {
	metrics []*Metric
}

func (m *Metrics) MarshalJSON() ([]byte, error) {
	val, err := json.Marshal(m.metrics)
	return val, err
}

func (m *Metrics) UnmarshalJSON(raw []byte) error {
	err := json.Unmarshal(raw, &m.metrics)
	return err
}

// Add a single metric to the collection.
func (m *Metrics) Add(name string, value float64, unit string, timestamp time.Time, metadata map[string]string) {
	m.metrics = append(m.metrics, &Metric{name, value, unit, timestamp, metadata})
}

// Append other metrics to this metrics struct.
func (m *Metrics) Append(metrics Metrics) {
	m.metrics = append(m.metrics, metrics.metrics...)
}

// Retrieve the metrics as an array.
func (m *Metrics) Items() []*Metric {
	return m.metrics
}

