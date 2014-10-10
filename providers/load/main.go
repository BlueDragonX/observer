package main

import (
	"../../api"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
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
	add := func(name, val string) error {
		floatval, err := strconv.ParseFloat(val, 64)
		if err == nil {
			metrics.Add(name, floatval, api.UNIT_COUNT, now, nil)
		}
		return err
	}

	if err = add("LoadAvg1", parts[0]); err != nil {
		return
	}
	if err = add("LoadAvg5", parts[1]); err != nil {
		return
	}
	if err = add("LoadAvg10", parts[2]); err != nil {
		return
	}
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
