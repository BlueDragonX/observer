package main

import (
	"../../api"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"
)

type Handler struct {
	rootfs string
}

// Configure the provider.
func (h *Handler) Configure(config api.Config) error {
	h.rootfs = "/"
	if config != nil {
		if rootfsIf, ok := config["rootfs"]; ok {
			if rootfsStr, ok := rootfsIf.(string); ok {
				h.rootfs = rootfsStr
			} else {
				return fmt.Errorf("invalid rootfs: %v", rootfsIf)
			}
		}
	}
	return nil
}

// Get the current load stats.
func (h *Handler) Get() (metrics api.Metrics, err error) {
	var raw []byte
	path := path.Join(h.rootfs, "proc", "loadavg")
	if raw, err = ioutil.ReadFile(path); err != nil {
		return
	}
	parts := strings.Split(string(raw), " ")
	if len(parts) < 3 {
		err = errors.New("invalid load string")
		return
	}

	now := time.Now().UTC()
	metrics.Add("LoadAvg1", parts[0], now, nil)
	metrics.Add("LoadAvg5", parts[1], now, nil)
	metrics.Add("LoadAvg10", parts[2], now, nil)
	return
}

// Return an error.
func (h *Handler) Put(metrics api.Metrics) error {
	return errors.New("source only provider")
}

// Run the provider.
func main() {
	api.RunProvider(&Handler{})
}
