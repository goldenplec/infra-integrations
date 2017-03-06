package main

import (
	"math/rand"

	sdk_args "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdk_args.DefaultArgumentList
	Environment string `default:"" help:"Environment variable."`
}

const (
	pluginName    = "example"
	pluginVersion = "1.0.0"
)

var args argumentList

func main() {
	// Initialize the output structure
	integration, err := sdk.NewIntegration(pluginName, pluginVersion, &args)
	fatalIfErr(err)
	log.SetupLogging(args.Verbose)

	// Build the metrics dictionary with a valid event_type:
	//  * LoadBalancerSample
	//  * BlockDeviceSample
	//  * DatastoreSample
	//  * QueueSample
	//  * ComputeSample
	//  * IamAccountSummarySample
	//  * PrivateNetworkSample
	//  * ServerlessSample
	// Provider may be set no anything identifying the data provider
	ms := integration.NewMetricSet("DatastoreSample", "ExampleServer")

	if args.Environment != "" {
		ms.AddMetric("environment", args.Environment, metric.ATTRIBUTE)
	}

	// Each metric specific to a provider should go prefixed with the provider namespace.
	keyList := []string{"provider.valueOne", "provider.valueTwo", "provider.valueThree"}
	for _, key := range keyList {
		value := rand.Int()
		ms.AddMetric(key, value, metric.GAUGE)
		log.Debug("Adding metric %s with value %d", key, value)
	}

	keyList = []string{"valueOne", "valueTwo", "valueThree"}
	itemKeys := []string{"item1", "item2", "item3"}
	for _, item := range itemKeys {
		integration.Inventory[item] = sdk.Inventory{}
		for _, key := range keyList {
			integration.Inventory[item][key] = rand.Int()
			log.Debug("Set inventory key %s=%d for %s", key, integration.Inventory[item][key], item)
		}
	}

	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
