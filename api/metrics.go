package api

import (
	"encoding/json"
	"strconv"
	"time"
)

type Metric struct {
	Name      string
	Value     string
	Timestamp time.Time
	Metadata  map[string]string
}

// Convert the metric value to an int.
func (m *Metric) Int() (val int64, err error) {
	val, err = strconv.ParseInt(m.Value, 10, 64)
	return
}

// Convert the metric value to a float.
func (m *Metric) Float(val float64, err error) {
	val, err = strconv.ParseFloat(m.Value, 64)
	return
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
func (m *Metrics) Add(name, value string, timestamp time.Time, metadata map[string]string) {
	m.metrics = append(m.metrics, &Metric{name, value, timestamp, metadata})
}

// Append other metrics to this metrics struct.
func (m *Metrics) Append(metrics Metrics) {
	m.metrics = append(m.metrics, metrics.metrics...)
}

// Retrieve the metrics as an array.
func (m *Metrics) Items() []*Metric {
	return m.metrics
}

