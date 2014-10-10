package api

import (
	"encoding/json"
	"io"
	"os"
	"os/exec"
)

type Controller struct {
	proc    *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	encoder *json.Encoder
	decoder *json.Decoder
}

// Create a new controller.
func NewController(cmd []string) (ctrl *Controller, err error) {
	var stdin io.WriteCloser
	var stdout io.ReadCloser
	var encoder *json.Encoder
	var decoder *json.Decoder
	proc := exec.Command(cmd[0], cmd[1:]...)
	proc.Stderr = os.Stderr

	if stdout, err = proc.StdoutPipe(); err == nil {
		decoder = json.NewDecoder(stdout)
	} else {
		return
	}

	if stdin, err = proc.StdinPipe(); err == nil {
		encoder = json.NewEncoder(stdin)
	} else {
		return
	}

	ctrl = &Controller{
		proc:    proc,
		stdin:   stdin,
		stdout:  stdout,
		encoder: encoder,
		decoder: decoder,
	}
	return
}

// Start the provider process.
func (ctrl *Controller) Start() error {
	return ctrl.proc.Start()
}

// Stop the provider process. Closes stdin, stdout, and waits for the process
// to exit.
func (ctrl *Controller) Stop() error {
	ctrl.stdin.Close()
	ctrl.stdout.Close()
	return ctrl.proc.Wait()
}

// Send a message to the provider process.
func (ctrl *Controller) send(request Message) (msg Message, err error) {
	if err = ctrl.encoder.Encode(&request); err == nil {
		err = ctrl.decoder.Decode(&msg)
	}
	return
}

// Configure the provider.
func (ctrl *Controller) Configure(cfg map[string]interface{}) error {
	if res, err := ctrl.send(Message{MSG_CFG_REQ, cfg}); err != nil {
		return err
	} else if res.Type == MSG_ERROR {
		return res.Payload.(error)
	}
	return nil
}

// Retrieve metrics from the provider.
func (ctrl *Controller) Get() (metrics Metrics, err error) {
	var res Message
	if res, err = ctrl.send(Message{MSG_GET_REQ, nil}); err == nil {
		if res.Type == MSG_ERROR {
			err = res.Payload.(error)
		} else {
			metrics = res.Payload.(Metrics)
		}
	}
	return
}

// Send metrics to the provider.
func (ctrl *Controller) Put(metrics Metrics) error {
	if res, err := ctrl.send(Message{MSG_PUT_REQ, metrics}); err != nil {
		return err
	} else if res.Type == MSG_ERROR {
		return res.Payload.(error)
	}
	return nil
}
