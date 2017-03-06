package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
)

var metricsDefinition = map[string][]interface{}{
	"provider.requestsPerSecond":  {"Total Accesses", metric.COUNTER},
	"provider.bytesPerSecond":     {getBytes, metric.COUNTER},
	"provider.idleWorkers":        {"IdleWorkers", metric.GAUGE},
	"provider.busyWorkers":        {"BusyWorkers", metric.GAUGE},
	"provider.writingWorkers":     {getWorkerStatus("W"), metric.GAUGE},
	"provider.loggingWorkers":     {getWorkerStatus("L"), metric.GAUGE},
	"provider.gracefulWorkers":    {getWorkerStatus("G"), metric.GAUGE},
	"provider.readingWorkers":     {getWorkerStatus("R"), metric.GAUGE},
	"provider.closingWorkers":     {getWorkerStatus("C"), metric.GAUGE},
	"provider.keepaliveWorkers":   {getWorkerStatus("K"), metric.GAUGE},
	"provider.DNSLookupWorkers":   {getWorkerStatus("D"), metric.GAUGE},
	"provider.idleCleanupWorkers": {getWorkerStatus("I"), metric.GAUGE},
	"provider.startingWorkers":    {getWorkerStatus("S"), metric.GAUGE},
	"provider.totalWorkers":       {getTotalWorkers, metric.GAUGE},
}

func asValue(value string) interface{} {
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}

	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	if b, err := strconv.ParseBool(value); err == nil {
		return b
	}
	return value
}

func populateMetrics(sample *metric.MetricSet, metrics map[string]interface{}, metricsDefinition map[string][]interface{}) error {
	for metricName, metricInfo := range metricsDefinition {
		rawSource := metricInfo[0]
		metricType := metricInfo[1].(metric.SourceType)

		var rawMetric interface{}
		var ok bool

		switch source := rawSource.(type) {
		case string:
			rawMetric, ok = metrics[source]
		case func(map[string]interface{}) (float64, bool):
			rawMetric, ok = source(metrics)
		default:
			log.Warn("Invalid raw source metric for %s", metricName)
			continue
		}

		if !ok {
			log.Warn("Can't find raw metrics in results for %s", metricName)
			continue
		}
		err := sample.AddMetric(metricName, rawMetric, metricType)

		if err != nil {
			log.Warn("Error setting value: %s", err)
			continue
		}
	}
	return nil
}

// getWorkerStatus counts occurence of a given letter, which means status of a worker
// (i.e. "W" corresponds to writing status of the worker)
func getWorkerStatus(status string) func(metrics map[string]interface{}) (float64, bool) {
	return func(metrics map[string]interface{}) (float64, bool) {
		scoreboard, ok := metrics["Scoreboard"].(string)
		if ok {
			return float64(strings.Count(scoreboard, status)), true
		}
		return 0, false
	}
}

// getTotalWorkers counts number of characters for Scoreboard key, which means total number of workers
func getTotalWorkers(metrics map[string]interface{}) (float64, bool) {
	scoreboard, ok := metrics["Scoreboard"].(string)
	if ok {
		return float64(len(scoreboard)), true
	}
	return 0, false
}

//getBytes converts value of Total kBytes into bytes
func getBytes(metrics map[string]interface{}) (float64, bool) {
	totalkBytes, ok := metrics["Total kBytes"].(int)
	if ok {
		return float64(totalkBytes * 1024), true
	}
	return 0, false
}

// getRawMetrics reads an Apache status message and transforms its
// contents into a map of metrics with the keys and values extracted from the
// status endpoint.
func getRawMetrics(reader *bufio.Reader) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	_, err := reader.Peek(1)
	if err != nil {
		return nil, fmt.Errorf("Empty result")
	}

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		splitedLine := strings.Split(line, ":")
		if len(splitedLine) != 2 {
			continue
		}
		metrics[splitedLine[0]] = asValue(strings.TrimSpace(splitedLine[1]))
	}

	if len(metrics) == 0 {
		return nil, fmt.Errorf("Empty result")
	}
	return metrics, nil
}

func getMetricsData(sample *metric.MetricSet) error {
	netClient := &http.Client{
		Timeout: time.Second * 1,
	}
	resp, err := netClient.Get(args.StatusURL)
	fatalIfErr(err)
	defer resp.Body.Close()

	var rawMetrics map[string]interface{}
	rawMetrics, err = getRawMetrics(bufio.NewReader(resp.Body))

	if err != nil {
		return err
	}
	return populateMetrics(sample, rawMetrics, metricsDefinition)
}
