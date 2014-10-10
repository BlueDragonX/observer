package main

import (
	"../../api"
	"fmt"
	"github.com/crowdmob/goamz/aws"
	"github.com/crowdmob/goamz/cloudwatch"
	"log"
	"time"
)

type Handler struct {
	accessKey string
	secretKey string
	region    string
	namespace string
}

// Configure the provider.
func (h *Handler) Configure(config api.Config) (err error) {
	if config == nil {
		err = fmt.Errorf("provider requires configuration")
		return
	}

	get := func(name string) (string, error) {
		if val, found := config[name]; found {
			if strval, valid := val.(string); valid && strval != "" {
				return strval, nil
			} else {
				return "", fmt.Errorf("invalid config value %s", name)
			}
		} else {
			return "", fmt.Errorf("missing config value %s", name)
		}
	}

	if h.accessKey, err = get("access-key"); err != nil {
		return
	}
	if h.secretKey, err = get("secret-key"); err != nil {
		return
	}
	if h.region, err = get("region"); err != nil {
		return
	}
	h.namespace, err = get("namespace")
	return
}

// Return an error.
func (h *Handler) Get() (metrics api.Metrics, err error) {
	return api.Metrics{}, fmt.Errorf("sink only provider")
}

// Put metrics into AWS CloudWatch.
func (h *Handler) Put(metrics api.Metrics) (err error) {
	var auth aws.Auth
	var cw *cloudwatch.CloudWatch
	region := aws.Regions[h.region]
	expire := time.Now().UTC().Add(10 * time.Minute)

	if auth, err = aws.GetAuth(h.accessKey, h.secretKey, "", expire); err != nil {
		return
	}
	if cw, err = cloudwatch.NewCloudWatch(auth, region.CloudWatchServicepoint); err != nil {
		return
	}

	metricItems := metrics.Items()
	cwMetrics := make([]cloudwatch.MetricDatum, 0, len(metricItems))
	for _, metric := range metricItems {
		dims := make([]cloudwatch.Dimension, 0, len(metric.Metadata))
		for key, value := range metric.Metadata {
			dims = append(dims, cloudwatch.Dimension{Name: key, Value: value})
		}

		var unit string
		switch metric.Unit {
		case api.UNIT_BYTES:
			unit = "Bytes"
		case api.UNIT_KILOBYTES:
			unit = "Kilobytes"
		case api.UNIT_MEGABYTES:
			unit = "Megabytes"
		case api.UNIT_GIGABYTES:
			unit = "Gigabytes"
		case api.UNIT_TERABYTES:
			unit = "Terabytes"
		case api.UNIT_BYTES_PER_SECOND:
			unit = "Bytes/Second"
		case api.UNIT_KILOBYTES_PER_SECOND:
			unit = "Kilobytes/Second"
		case api.UNIT_MEGABYTES_PER_SECOND:
			unit = "Megabytes/Second"
		case api.UNIT_GIGABYTES_PER_SECOND:
			unit = "Gigabytes/Second"
		case api.UNIT_TERABYTES_PER_SECOND:
			unit = "Terabytes/Second"
		case api.UNIT_BITS:
			unit = "Bits"
		case api.UNIT_KILOBITS:
			unit = "Kilobits"
		case api.UNIT_MEGABITS:
			unit = "Megabits"
		case api.UNIT_GIGABITS:
			unit = "Gigabits"
		case api.UNIT_TERABITS:
			unit = "Terabits"
		case api.UNIT_BITS_PER_SECOND:
			unit = "Bits/Second"
		case api.UNIT_KILOBITS_PER_SECOND:
			unit = "Kilobits/Second"
		case api.UNIT_MEGABITS_PER_SECOND:
			unit = "Megabits/Second"
		case api.UNIT_GIGABITS_PER_SECOND:
			unit = "Gigabits/Second"
		case api.UNIT_TERABITS_PER_SECOND:
			unit = "Terabits/Second"
		case api.UNIT_SECONDS:
			unit = "Seconds"
		case api.UNIT_PERCENT:
			unit = "Percent"
		default:
			fallthrough
		case api.UNIT_COUNT:
			unit = "Count"
		case api.UNIT_COUNT_PER_SECOND:
			unit = "Count/Second"
		}

		cwMetric := cloudwatch.MetricDatum{
			Dimensions: dims,
			MetricName: metric.Name,
			Timestamp:  metric.Timestamp,
			Unit:       unit,
			Value:      metric.Value,
		}
		cwMetrics = append(cwMetrics, cwMetric)
	}

	_, err = cw.PutMetricDataNamespace(cwMetrics, h.namespace)
	return
}

// Run the provider.
func main() {
	api.RunProvider(&Handler{})
}
