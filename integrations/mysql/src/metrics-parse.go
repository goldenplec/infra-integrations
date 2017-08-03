package main

import (
	"strconv"

	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

const (
	inventoryQuery = "SHOW GLOBAL VARIABLES"
	metricsQuery   = "SHOW /*!50002 GLOBAL */ STATUS"
	replicaQuery   = "SHOW SLAVE STATUS"
)

// Try to convert a string to its type or return the string if not possible
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

func getRawData(db dataSource) (map[string]interface{}, map[string]interface{}, error) {
	inventory, err := db.query(inventoryQuery)
	if err != nil {
		return nil, nil, err
	}
	metrics, err := db.query(metricsQuery)
	if err != nil {
		return nil, nil, err
	}
	replication, err := db.query(replicaQuery)
	if err != nil {
		log.Warn("Can't get node type, not enough privileges (must grant REPLICATION CLIENT)")
	} else if len(replication) == 0 {
		metrics["node_type"] = "master"
	} else {
		metrics["node_type"] = "slave"
	}

	// Set needed values for computed metrics from inventory
	metrics["key_cache_block_size"] = inventory["key_cache_block_size"]
	metrics["key_buffer_size"] = inventory["key_buffer_size"]
	metrics["version_comment"] = inventory["version_comment"]
	metrics["version"] = inventory["version"]

	return inventory, metrics, nil
}

<<<<<<< HEAD
func populateInventory(inventory map[string]sdk.Inventory, rawData map[string]interface{}) {
	for name, value := range rawData {
		inventory[name] = sdk.Inventory{"value": value}
=======
func populateInventory(inventory sdk.Inventory, rawData map[string]interface{}) {
	for name, value := range rawData {
		inventory.SetItem(name, "value", value)
>>>>>>> upstream/master
	}
}

func populateMetrics(sample *metric.MetricSet, rawMetrics map[string]interface{}) {
	populatePartialMetrics(sample, rawMetrics, defaultMetrics)
	if args.ExtendedMetrics {
		populatePartialMetrics(sample, rawMetrics, extendedMetrics)
	}
	if args.ExtendedInnodbMetrics {
		populatePartialMetrics(sample, rawMetrics, innodbMetrics)
	}
	if args.ExtendedMyIsamMetrics {
		populatePartialMetrics(sample, rawMetrics, myisamMetrics)
	}

}
func populatePartialMetrics(sample *metric.MetricSet, metrics map[string]interface{}, metricsDefinition map[string][]interface{}) {
	for metricName, metricConf := range metricsDefinition {
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
			log.Warn("Invalid raw source metric for %s", metricName)
			continue
		}

		if !ok {
			log.Warn("Can't find raw metrics in results for %s", metricName)
			continue
		}

<<<<<<< HEAD
		sample.AddMetric(metricName, rawMetric, metricType)
=======
		sample.SetMetric(metricName, rawMetric, metricType)
>>>>>>> upstream/master
	}
}
