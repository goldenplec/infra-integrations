package main

import (
	"strconv"

	sdk_args "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/jmx"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdk_args.DefaultArgumentList

	Hostname   string `default:"localhost" help:"Hostname or IP where Cassandra is running."`
	Port       int    `default:"7199" help:"Port on which JMX server is listening."`
	Username   string `default:"" help:"Username for accessing JMX."`
	Password   string `default:"" help:"Password for the given user."`
	ConfigPath string `default:"/etc/cassandra.yaml" help:"Cassandra configuration file."`
}

const (
<<<<<<< HEAD
	integrationName    = "cassandra"
	integrationVersion = "0.1.0"
=======
	integrationName    = "com.newrelic.cassandra"
	integrationVersion = "1.0.0"
>>>>>>> upstream/master
)

var (
	args argumentList
)

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)
	log.SetupLogging(args.Verbose)

	fatalIfErr(jmx.Open(args.Hostname, strconv.Itoa(args.Port), args.Username, args.Password))
	defer jmx.Close()

	if args.All || args.Metrics {
<<<<<<< HEAD
		rawMetrics, allKeyspaces, err := getMetrics()
		fatalIfErr(err)

		ms := integration.NewMetricSet("DatastoreSample", "Cassandra")
=======
		rawMetrics, allColumnFamilies, err := getMetrics()
		fatalIfErr(err)

		ms := integration.NewMetricSet("CassandraSample")
>>>>>>> upstream/master

		populateMetrics(ms, rawMetrics, metricsDefinition)
		populateMetrics(ms, rawMetrics, commonDefinition)

<<<<<<< HEAD
		for _, keyspaceMetrics := range allKeyspaces {
			ms := integration.NewMetricSet("DatastoreSample", "CassandraKeyspace")
			populateMetrics(ms, keyspaceMetrics, keyspaceDefinition)
=======
		for _, columnFamilyMetrics := range allColumnFamilies {
			ms := integration.NewMetricSet("CassandraColumnFamilySample")
			populateMetrics(ms, columnFamilyMetrics, columnFamilyDefinition)
>>>>>>> upstream/master
			populateMetrics(ms, rawMetrics, commonDefinition)
		}
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
