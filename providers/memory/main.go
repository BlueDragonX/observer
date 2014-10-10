package main

import (
	"../../api"
	"bufio"
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
	rootfs string
	swap   bool
}

// Configure the provider.
func (h *Handler) Configure(config api.Config) error {
	h.rootfs = "/"
	h.swap = true
	if config != nil {
		if rootfsIf, ok := config["rootfs"]; ok {
			if rootfsStr, ok := rootfsIf.(string); ok {
				h.rootfs = rootfsStr
			} else {
				return fmt.Errorf("invalid rootfs: %v", rootfsIf)
			}
		}

		if swapIf, ok := config["swap"]; ok {
			if swapBool, ok := swapIf.(bool); ok {
				h.swap = swapBool
			} else {
				return fmt.Errorf("invalid swap: %v", swapIf)
			}
		}
	}
	return nil
}

// Get the current memory stats.
func (h *Handler) Get() (metrics api.Metrics, err error) {
	path := path.Join(h.rootfs, "proc", "meminfo")
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var splitter *regexp.Regexp
	if splitter, err = regexp.Compile("\\s+"); err != nil {
		return
	}

	values := make(map[string]int64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := splitter.Split(scanner.Text(), 3)
		if len(parts) != 3 {
			continue
		}
		key := strings.TrimRight(parts[0], ":")
		if value, lineErr := strconv.ParseInt(parts[1], 10, 64); lineErr == nil {
			unit := strings.TrimSpace(strings.ToLower(parts[2]))
			switch unit {
			case "tb":
				value *= 1024
				fallthrough
			case "gb":
				value *= 1024
				fallthrough
			case "mb":
				value *= 1024
				fallthrough
			case "kb":
				value *= 1024
			}
			values[key] = value
		} else {
			continue
		}
	}

	getInt := func(name string) (intval int64, err error) {
		var ok bool
		if intval, ok = values[name]; !ok {
			err = errors.New(fmt.Sprintf("meminfo %s missing", name))
		}
		return
	}

	now := time.Now().UTC()
	addBytes := func(name string, intval int64) {
		metrics.Add(name, float64(intval), api.UNIT_BYTES, now, nil)
	}

	addPct := func(name string, pctval float64) {
		metrics.Add(name, pctval, api.UNIT_PERCENT, now, nil)
	}

	var memTotal, memFree int64
	if memTotal, err = getInt("MemTotal"); err != nil {
		return
	}
	if memFree, err = getInt("MemFree"); err != nil {
		return
	}
	addBytes("MemoryTotal", memTotal)
	addBytes("MemoryFree", memFree)
	addPct("MemoryUtilization", 100*(float64(memTotal-memFree))/float64(memTotal))

	if h.swap {
		var swapTotal, swapFree int64
		if swapTotal, err = getInt("SwapTotal"); err != nil {
			return
		}
		if swapFree, err = getInt("SwapFree"); err != nil {
			return
		}
		addBytes("SwapTotal", swapTotal)
		addBytes("SwapFree", swapFree)
		addPct("SwapUtilization", 100*(float64(swapTotal-swapFree))/float64(swapTotal))
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
