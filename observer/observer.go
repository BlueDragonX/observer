package main

import (
	"../config"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

type Observer struct {
	sinks map[string]*Sink
	sources map[string]*Source
	pipes map[string]*Pipe
	stop  chan interface{}
}

// Create a new observer with the provided configuration.
func NewObserver(config *config.Config) (obs *Observer, err error) {
	sinkNames := make(map[string]interface{}, len(config.Pipes) )
	for _, pipeConfig := range config.Pipes {
		for _, sinkName := range pipeConfig.Sinks {
			sinkNames[sinkName] = nil
		}
	}

	sinks := make(map[string]*Sink, len(sinkNames))
	for sinkName, _ := range sinkNames {
		if sinkConfig, ok := config.Sinks[sinkName]; ok {
			var sink *Sink
			if sink, err = NewSink(sinkName, sinkConfig); err == nil {
				sinks[sinkName] = sink
			} else {
				return
			}
		} else {
			err = fmt.Errorf("sink %s not defined")
			return
		}
	}

	sourceNames := make(map[string]interface{}, len(config.Pipes) )
	for _, pipeConfig := range config.Pipes {
		for _, sourceName := range pipeConfig.Sources {
			sourceNames[sourceName] = nil
		}
	}

	sources := make(map[string]*Source, len(sourceNames))
	for sourceName, _ := range sourceNames {
		if sourceConfig, ok := config.Sources[sourceName]; ok {
			var source *Source
			if source, err = NewSource(sourceName, sourceConfig); err == nil {
				sources[sourceName] = source
			} else {
				return
			}
		} else {
			err = fmt.Errorf("source %s not defined")
			return
		}
	}

	pipes := make(map[string]*Pipe, len(config.Pipes))
	for pipeName, pipeConfig := range config.Pipes {
		pipeSources := []*Source{}
		for _, sourceName := range pipeConfig.Sources {
			pipeSources = append(pipeSources, sources[sourceName])
		}

		pipeSinks := []*Sink{}
		for _, sinkName := range pipeConfig.Sinks {
			pipeSinks = append(pipeSinks, sinks[sinkName])
		}

		pipes[pipeName] = NewPipe(
			pipeName,
			time.Duration(pipeConfig.Interval) * time.Second,
			pipeSources,
			pipeSinks,
			pipeConfig.Metadata,
		)
	}

	obs = &Observer{sinks, sources, pipes, make(chan interface{})}
	return
}

// Run the observer.
func (obs *Observer) Run() (err error) {
	sinksStarted := make([]string, 0, len(obs.sinks))
	sourcesStarted := make([]string, 0, len(obs.sources))
	pipesStarted := make([]string, 0, len(obs.pipes))

	defer func() {
		for _, name := range pipesStarted {
			obs.pipes[name].Stop()
		}
		for _, name := range sourcesStarted {
			obs.sources[name].Stop()
		}
		for _, name := range sinksStarted {
			obs.sinks[name].Stop()
		}
	}()

	for _, sink := range obs.sinks {
		if err = sink.Start(); err != nil {
			log.Printf("sink %s failed to start: %s\n", sink.Name, err)
			return
		}
		sinksStarted = append(sinksStarted, sink.Name)
	}
	sort.Strings(sinksStarted)
	log.Printf("sinks started: %s\n", strings.Join(sinksStarted, ", "))

	for _, source := range obs.sources {
		if err = source.Start(); err != nil {
			log.Printf("source %s failed to start: %s\n", source.Name, err)
			return
		}
		sourcesStarted = append(sourcesStarted, source.Name)
	}
	sort.Strings(sourcesStarted)
	log.Printf("sources started: %s\n", strings.Join(sourcesStarted, ", "))

	for _, pipe := range obs.pipes {
		pipe.Start()
		pipesStarted = append(pipesStarted, pipe.Name)
	}
	sort.Strings(pipesStarted)
	log.Printf("pipes started: %s\n", strings.Join(pipesStarted, ", "))

	log.Printf("started")
	<-obs.stop
	return
}

// Stop the observer.
func (obs *Observer) Stop() {
	close(obs.stop)
}
