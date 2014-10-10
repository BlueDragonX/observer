package main

import (
	"../../api"
	"log"
	"math/rand"
	"time"
)

type Handler struct {}

// Configure the provider.
func (h *Handler) Configure(config api.Config) error {
	return nil
}

// Get a random integer.
func (h *Handler) Get() (metrics api.Metrics, err error) {
	metrics.Add("RandomCount", float64(rand.Int31()), api.UNIT_COUNT, time.Now().UTC(), nil)
	return
}

// Log the metrics.
func (h *Handler) Put(metrics api.Metrics) error {
	items := metrics.Items()
	log.Printf("received %d metrics:\n", len(items))
	for _, metric := range items {
		if metric.Metadata == nil || len(metric.Metadata) == 0 {
			log.Printf("  %s: %0.3f %s\n", metric.Name, metric.Value, metric.Unit)
		} else {
			log.Printf("  %s: %0.3f %s (%v)\n", metric.Name, metric.Value, metric.Unit, metric.Metadata)
		}
	}
	return nil
}

// Run the provider.
func main() {
	log.SetFlags(0)
	rand.Seed(time.Now().UTC().UnixNano())
	api.RunProvider(&Handler{})
}
