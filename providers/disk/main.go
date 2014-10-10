package main

import (
	"../../api"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	path string
}

// Configure the provider.
func (h *Handler) Configure(config api.Config) error {
	if config != nil {
		if pathInt, ok := config["path"]; ok {
			if pathStr, ok := pathInt.(string); ok {
				h.path = pathStr
				return nil
			} else {
				return fmt.Errorf("invalid path: %v", pathInt)
			}
		}
	}
	return errors.New("path is required")
}

// Get the current disk stats.
func (h *Handler) Get() (metrics api.Metrics, err error) {
	var out []byte
	cmd := exec.Command("stat", "--format=%S %b %a", "-f", h.path)
	if out, err = cmd.Output(); err != nil {
		return
	}
	parts := strings.Split(strings.TrimSpace(string(out)), " ")
	if len(parts) < 3 {
		err = errors.New("invalid stat string")
		return
	}

	var blockSize, totalBlocks, freeBlocks int64
	if blockSize, err = strconv.ParseInt(parts[0], 10, 64); err != nil {
		return
	}
	if totalBlocks, err = strconv.ParseInt(parts[1], 10, 64); err != nil {
		return
	}
	if freeBlocks, err = strconv.ParseInt(parts[2], 10, 64); err != nil {
		return
	}

	now := time.Now().UTC()
	metadata := map[string]string{"Path": h.path}
	metrics.Add("DiskBytesTotal", float64(blockSize*totalBlocks), api.UNIT_BYTES, now, metadata)
	metrics.Add("DiskBytesFree", float64(blockSize*freeBlocks), api.UNIT_BYTES, now, metadata)

	pctUsed := 100 * float64(totalBlocks-freeBlocks) / float64(totalBlocks)
	metrics.Add("DiskBytesUtilization", pctUsed, api.UNIT_PERCENT, now, metadata)
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
