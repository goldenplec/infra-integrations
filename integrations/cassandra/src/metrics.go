package main

import (
	"regexp"

	"github.com/newrelic/infra-integrations-sdk/jmx"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
)

// getMetrics will gather all node and keyspace level metrics and return them as two maps
// The main metrics map will contain all the keys got from JMX and the keyspace metrics map
// Will contain maps for each <keyspace>.<columnFamily> found while inspecting JMX metrics.
func getMetrics() (map[string]interface{}, map[string]map[string]interface{}, error) {
	metrics := make(map[string]interface{})
	keyspaceMetrics := make(map[string]map[string]interface{})

	re, err := regexp.Compile("keyspace=(.*),scope=(.*?),")
	if err != nil {
		return nil, nil, err
	}

	for _, query := range jmxPatterns {
		results, err := jmx.Query(query)
		if err != nil {
			return nil, nil, err
		}
		for key, value := range results {
			matches := re.FindStringSubmatch(key)
			key = re.ReplaceAllString(key, "")

			if len(matches) != 3 {
				metrics[key] = value
			} else {
				columnfamily := matches[2]
				keyspace := matches[1]
				eventkey := keyspace + "." + columnfamily

				_, ok := keyspaceMetrics[eventkey]
				if !ok {
					keyspaceMetrics[eventkey] = make(map[string]interface{})
					keyspaceMetrics[eventkey]["keyspace"] = keyspace
					keyspaceMetrics[eventkey]["columnFamily"] = columnfamily
					keyspaceMetrics[eventkey]["keyspaceAndColumnFamily"] = eventkey
				}
				keyspaceMetrics[eventkey][key] = value
			}
		}
	}

	return metrics, keyspaceMetrics, nil
}

func populateMetrics(sample *metric.MetricSet, metrics map[string]interface{}, definition map[string][]interface{}) {
	for metricName, metricConf := range definition {
		rawSource := metricConf[0]
		metricType := metricConf[1].(metric.SourceType)

		var rawMetric interface{}
		var ok bool

		switch source := rawSource.(type) {
		case string:
			rawMetric, ok = metrics[source]
		case func(map[string]interface{}) (float64, bool):
			rawMetric, ok = source(metrics)
		default:
			log.Debug("Invalid raw source metric for %s", metricName)
			continue
		}

		if !ok {
			log.Debug("Can't find raw metrics in results for %s", metricName)
			continue
		}

		err := sample.AddMetric(metricName, rawMetric, metricType)
		if err != nil {
			log.Warn("Error setting value: %s", err)
			continue
		}
	}
}
