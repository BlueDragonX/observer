package main

import (
	"fmt"
	"../api"
	"../config"
)

type putRequest struct {
	metrics api.Metrics
	response chan error
}

type Sink struct {
	Name       string
	config     map[string]interface{}
	controller *api.Controller
	putQueue   chan putRequest
	stop       chan error
}

// Create a new sink with the provided name and configuration.
func NewSink(name string, config config.Provider) (sink *Sink, err error) {
	if config.Provider == "" {
		err = fmt.Errorf("sink %s has no provider", name)
	} else {
		var path string
		var ctrl *api.Controller

		if path, err = lookupProvider(config.Provider); err != nil {
			return
		}

		if ctrl, err = api.NewController([]string{path}); err == nil {
			sink = &Sink{name, config.Config, ctrl, make(chan putRequest), make(chan error)}
		}
	}
	return
}

// Start the sink.
func (s *Sink) Start() (err error) {
	if err = s.controller.Start(); err != nil {
		return
	}
	if err = s.controller.Configure(s.config); err != nil {
		close(s.putQueue)
		s.controller.Stop()
		return
	}

	go func() {
		for {
			req, open := <-s.putQueue
			if !open {
				break
			}
			req.response <- s.controller.Put(req.metrics)
			close(req.response)
		}
		s.stop <- s.controller.Stop()
		close(s.stop)
	}()
	return
}

// Stop the sink.
func (s *Sink) Stop() error {
	close(s.putQueue)
	return <-s.stop
}

// Put metrics into the sink.
func (s *Sink) Put(metrics api.Metrics) error {
	res := make(chan error)
	s.putQueue <- putRequest{metrics, res}
	return <-res
}
