package main

import (
	sdk_args "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdk_args.DefaultArgumentList

	Hostname   string `default:"localhost" help:"Hostname or IP where Cassandra is running."`
	Port       int    `default:"3306" help:"Port on which JMX server is listening."`
	Username   string `help:"Username for accessing JMX."`
	Password   string `help:"Password for the given user."`
	ConfigPath string `default:"/etc/cassandra.yaml" help:"Cassandra configuration file."`
}

const (
	integrationName    = "cassandra"
	integrationVersion = "1.0.0"
)

var (
	args argumentList
)

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)
	log.SetupLogging(args.Verbose)

	if args.All || args.Metrics {
		rawMetrics, err := getMetrics()
		fatalIfErr(err)

		ms := integration.NewMetricSet("DatastoreSample", "Cassandra")
		populateMetrics(ms, rawMetrics)

	}
	if args.All || args.Inventory {
		rawInventory, err := getInventory()
		fatalIfErr(err)
		populateInventory(integration.Inventory, rawInventory)
	}

	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
