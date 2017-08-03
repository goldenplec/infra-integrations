package main

import (
<<<<<<< HEAD
=======
	"reflect"
>>>>>>> upstream/master
	"testing"

	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

func TestPopulatelMetrics(t *testing.T) {
	var rawMetrics = map[string]interface{}{
		"raw_metric_1": 1,
		"raw_metric_2": 2,
		"raw_metric_3": "foo",
	}

	functionSource := func(a map[string]interface{}) (float64, bool) {
		return float64(a["raw_metric_1"].(int) + a["raw_metric_2"].(int)), true
	}

	var metricDefinition = map[string][]interface{}{
		"rawMetric1":     {"raw_metric_1", metric.GAUGE},
		"rawMetric2":     {"raw_metric_2", metric.GAUGE},
		"rawMetric3":     {"raw_metric_3", metric.ATTRIBUTE},
		"unknownMetric":  {"raw_metric_4", metric.GAUGE},
		"badRawSource":   {10, metric.GAUGE},
		"functionSource": {functionSource, metric.GAUGE},
	}

<<<<<<< HEAD
	var sample = metric.NewMetricSet("eventType", "provider")
=======
	var sample = metric.NewMetricSet("eventType")
>>>>>>> upstream/master
	populateMetrics(&sample, rawMetrics, metricDefinition)

	if sample["rawMetric1"] != 1 {
		t.Error()
	}
	if sample["rawMetric2"] != 2 {
		t.Error()
	}
	if sample["rawMetric3"] != "foo" {
		t.Error()
	}

	if sample["unknownMetric"] != nil {
		t.Error()
	}
	if sample["badRawSource"] != nil {
		t.Error()
	}
	if sample["functionSource"] != float64(3) {
		t.Error()
	}

}

func TestPopulateInventory(t *testing.T) {
	var rawInventory = map[string]interface{}{
<<<<<<< HEAD
		"key_1": 1,
		"key_2": 2,
		"key_3": "foo",
		"key_4": map[interface{}]interface{}{"test": 2},
	}

	inventory := make(map[string]sdk.Inventory)
	populateInventory(inventory, rawInventory)
	for key, value := range rawInventory {
		if key == "key_4" {
			for subk, subv := range value.(map[interface{}]interface{}) {
				if inventory[key][subk.(string)] != subv {
					t.Error()
				}
			}
		} else if inventory[key]["value"] != value {
			t.Error()
		}
=======
		"key_1":                 1,
		"key_2":                 2,
		"key_3":                 "foo",
		"key_4":                 map[interface{}]interface{}{"test": 2},
		"my_important_password": "12345",
		"key_6":                 map[interface{}]interface{}{"otherImportantPassword": 54321},
	}

	inventory := make(sdk.Inventory)
	populateInventory(inventory, rawInventory)

	expected := sdk.Inventory{
		"key_1":                 map[string]interface{}{"value": 1},
		"key_2":                 map[string]interface{}{"value": 2},
		"key_3":                 map[string]interface{}{"value": "foo"},
		"key_4":                 map[string]interface{}{"test": 2},
		"my_important_password": map[string]interface{}{"value": "(omitted value)"},
		"key_6":                 map[string]interface{}{"otherImportantPassword": "(omitted value)"},
	}

	if !reflect.DeepEqual(inventory, expected) {
		t.Errorf("Expected: %v. Actual: %v", expected, inventory)
>>>>>>> upstream/master
	}
}
