package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
)

var metricsPlusDefinition = map[string][]interface{}{
	"provider.connectionsActive":            {"connections.active", metric.GAUGE},
	"provider.connectionsIdle":              {"connections.idle", metric.GAUGE},
	"provider.connectionsAcceptedPerSecond": {"connections.accepted", metric.COUNTER},
	"provider.connectionsDroppedPerSecond":  {"connections.dropped", metric.COUNTER},
	"provider.requestsPerSecond":            {"requests.total", metric.COUNTER},
	"software.edition":                      {"edition", metric.ATTRIBUTE},
	"software.version":                      {"version", metric.ATTRIBUTE},
}

var metricsStandardDefinition = map[string][]interface{}{
	"provider.connectionsActive":            {"active", metric.GAUGE},
	"provider.reading":                      {"reading", metric.GAUGE},
	"provider.waiting":                      {"waiting", metric.GAUGE},
	"provider.writing":                      {"writing", metric.GAUGE},
	"provider.requestsPerSecond":            {"requests", metric.COUNTER},
	"provider.connectionsAcceptedPerSecond": {"accepted", metric.COUNTER},
	"provider.connectionsDroppedPerSecond":  {connectionsDroppedPerSecond, metric.COUNTER},
	"software.edition":                      {"edition", metric.ATTRIBUTE},
	"software.version":                      {"version", metric.ATTRIBUTE},
}

// expressions contains the structure of the input data and defines the attributes we want to store
var nginxStatusExpressions = []*regexp.Regexp{
	regexp.MustCompile(`Active connections:\s+(?P<active>\d+)`),
	nil,
	regexp.MustCompile(`\s*(?P<accepted>\d+)\s+(?P<handled>\d+)\s+(?P<requests>\d+)`),
	regexp.MustCompile(`Reading: (?P<reading>\d+)\s+Writing: (?P<writing>\d+)\s+Waiting: (?P<waiting>\d+)`),
}

func connectionsDroppedPerSecond(metrics map[string]interface{}) (int, bool) {
	accepts, ok1 := metrics["accepted"].(int)
	handled, ok2 := metrics["handled"].(int)

	if ok1 && ok2 {
		return accepts - handled, true
	}
	return 0, false
}

// getMetrics reads an NGINX (open edition) status message and transforms its
// contents into a map that can be processed by NR agent.
// It returns a map of metrics with all the keys and values extracted from the
// status endpoint.
func getStandardMetrics(reader *bufio.Reader) (map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	for lineNo, re := range nginxStatusExpressions {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			return metrics, nil
		}

		if re == nil {
			continue
		}

		match := re.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("Line %d of status doesn't match", lineNo)
		}

		for i, name := range re.SubexpNames() {
			if i != 0 {
				value, err := strconv.Atoi(match[i])
				if err != nil {
					log.Warn("Can't cast value '%s'", match[i])
					continue
				}
				metrics[name] = value
			}
		}
	}
	metrics["version"] = ""
	metrics["edition"] = "open source"

	return metrics, nil
}

// getPlusMetrics reads an NGINX (Plus edition) status message, gets some
// metrics and transforms the contents into a map that can be processed by NR
// agent.
// It returns a map of metrics keys -> values.
func getPlusMetrics(reader *bufio.Reader) (map[string]interface{}, error) {
	jsonMetrics := make(map[string]interface{})
	metrics := make(map[string]interface{})

	dec := json.NewDecoder(reader)
	err := dec.Decode(&jsonMetrics)
	if err != nil {
		return nil, err
	}

	roots := [2]string{"connections", "requests"}

	for _, rootKey := range roots {
		rootNode, ok := jsonMetrics[rootKey].(map[string]interface{})
		if !ok {
			log.Warn("Can't assert type for %s", rootNode)
			continue
		}
		for key, value := range rootNode {
			metrics[fmt.Sprintf("%s.%s", rootKey, key)] = int(value.(float64))
		}
	}
	metrics["version"] = jsonMetrics["nginx_version"]
	metrics["edition"] = "plus"
	return metrics, nil
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
		case func(map[string]interface{}) (int, bool):
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

func getMetricsData(sample *metric.MetricSet) error {
	netClient := &http.Client{
		Timeout: time.Second * 1,
	}
	resp, err := netClient.Get(args.StatusURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	var rawMetrics map[string]interface{}
	var metricsDefinition map[string][]interface{}

	if resp.Header.Get("content-type") == "application/json" {
		metricsDefinition = metricsPlusDefinition
		rawMetrics, err = getPlusMetrics(bufio.NewReader(resp.Body))
	} else {
		metricsDefinition = metricsStandardDefinition
		rawMetrics, err = getStandardMetrics(bufio.NewReader(resp.Body))
		rawVersion := strings.Replace(resp.Header.Get("Server"), "nginx/", "", -1)
		rawMetrics["version"] = rawVersion

	}
	if err != nil {
		return err
	}
	return populateMetrics(sample, rawMetrics, metricsDefinition)
}
