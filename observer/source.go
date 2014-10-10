package main

import (
	"fmt"
	"../api"
	"../config"
	"errors"
)

type getResponse struct {
	metrics api.Metrics
	err     error
}

type getRequest chan getResponse

type Source struct {
	Name       string
	config     map[string]interface{}
	controller *api.Controller
	getQueue   chan getRequest
	stop       chan error
}

// Create a new source with the provided namd and configuration.
func NewSource(name string, config config.Provider) (source *Source, err error) {
	if config.Provider == "" {
		err = fmt.Errorf("source %s has no provider", name)
	} else {
		var path string
		var ctrl *api.Controller 

		if path, err = lookupProvider(config.Provider); err != nil {
			return
		}

		if ctrl, err = api.NewController([]string{path}); err == nil {
			source = &Source{name, config.Config, ctrl, make(chan getRequest), make(chan error)}
		}
	}
	return
}

// Start the source.
func (s *Source) Start() (err error) {
	if err = s.controller.Start(); err != nil {
		return
	}
	if err = s.controller.Configure(s.config); err != nil {
		close(s.getQueue)
		s.controller.Stop()
		return
	}

	go func() {
		for {
			req, open := <-s.getQueue
			if !open {
				break
			}
			metrics, err := s.controller.Get()
			req <- getResponse{metrics, err}
			close(req)
		}
		s.stop <- s.controller.Stop()
		close(s.stop)
	}()
	return
}

// Stop the source.
func (s *Source) Stop() error {
	close(s.getQueue)
	return <-s.stop
}

// Get metrics from the source.
func (s *Source) Get() (api.Metrics, error) {
	req := make(chan getResponse)
	s.getQueue <- req
	res, open := <-req
	if !open {
		return api.Metrics{}, errors.New("source channel closed")
	}
	return res.metrics, res.err
}
