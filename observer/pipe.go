package main

import (
	"../api"
	"log"
	"time"
)

type Pipe struct {
	Name     string
	interval time.Duration
	sources  []*Source
	sinks    []*Sink
	metadata map[string]string
	stop     chan chan interface{}
}

// Create a new pipe with the provided name and configuration.
func NewPipe(name string, interval time.Duration, sources []*Source, sinks []*Sink, metadata map[string]string) *Pipe {
	return &Pipe{
		name,
		interval,
		sources,
		sinks,
		metadata,
		make(chan chan interface{}),
	}
}

// Run the pipe.
func (p *Pipe) Start() {
	execute := func() {
		var metrics api.Metrics
		for _, source := range p.sources {
			if sourceMetrics, err := source.Get(); err == nil {
				metrics.Append(sourceMetrics)
			} else {
				log.Printf("%s:%s get error: %s\n", p.Name, source.Name, err)
			}
		}

		for _, metric := range metrics.Items() {
			metric.Underlay(p.metadata)
		}

		for _, sink := range p.sinks {
			if err := sink.Put(metrics); err != nil {
				log.Printf("%s:%s put error: %s\n", p.Name, sink.Name, err)
			}
		}
	}

	go func() {
		execute()
	Loop:
		for {
			select {
			case res := <-p.stop:
				res <- nil
				close(res)
				break Loop
			case <-time.Tick(p.interval):
				execute()
			}
		}
	}()
}

// Stop the pipe.
func (p *Pipe) Stop() {
	res := make(chan interface{})
	p.stop <- res
	close(p.stop)
	<-res
}
