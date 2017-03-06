package main

import (
	"net/http"
	"time"

	sdk_args "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

const (
	integrationName    = "apache"
	integrationVersion = "1.0.0"
)

var (
	args      argumentList
	netClient = &http.Client{
		Timeout: time.Second * 1,
	}
)

type argumentList struct {
	sdk_args.DefaultArgumentList
	StatusURL string `default:"http://127.0.0.1/server-status?auto" help:"Apache status-server URL."`
}

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)
	log.SetupLogging(args.Verbose)

	if args.All || args.Inventory {
		log.Info("Getting data for '%s' plugin", integrationName+"-inventory")
		integration.Inventory, err = getInventory()
		fatalIfErr(err)
	}

	if args.All || args.Metrics {
		log.Info("Getting data for '%s' plugin", integrationName+"-metrics")
		ms := integration.NewMetricSet("LoadBalancerSample", "Apache")
		fatalIfErr(getMetricsData(ms))
	}

	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
