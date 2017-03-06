#!/usr/bin/env python
import argparse
import random
import json
import logging
import os


class PluginData:
    def __init__(self, name, protocol_version, plugin_version):
        self.name = name
        self.protocol_version = protocol_version
        self.plugin_version = plugin_version
        self.status = "OK"
        self.inventory = {}
        self.metrics = []
        self.events = []

    def addInventory(self, item, key, value):
        self.inventory.setdefault(item, {})[key] = value

    def addMetric(self, metric_dict):
        self.metrics.append(metric_dict)

def parse_arguments():
    parser = argparse.ArgumentParser()
    parser.add_argument('-v', default=False, dest='verbose', action='store_true',
                        help='Print more information to logs')
    parser.add_argument('-p', default=False, dest='pretty', action='store_true',
                        help='Print pretty formatted JSON')

    return parser.parse_args()


if __name__ == "__main__":
    # Setup the plugin's command line parameters
    args = parse_arguments()

    # Setup logging, redirect logs to stderr and configure the log level.
    logger = logging.getLogger("infra")
    logger.addHandler(logging.StreamHandler())
    if args.verbose:
        logger.setLevel(logging.DEBUG)
    else:
        logger.setLevel(logging.INFO)

    # Initialize the output object
    data = PluginData("example", "1", "1.0.0")

    # Build the metrics dictionary with a valid event_type:
    #  * LoadBalancerSample
    #  * BlockDeviceSample
    #  * DatastoreSample
    #  * QueueSample
    #  * ComputeSample
    #  * IamAccountSummarySample
    #  * PrivateNetworkSample
    #  * ServerlessSample
    # Provider may be set no anything identifying the data provider
    metric = {
        "event_type": "DatastoreSample",
        "provider": "ExampleServer",
    }

    # Get ENVIRONMENT variable set by the agent
    env = os.getenv("ENVIRONMENT")
    if env:
        metric["environment"] = env

    # Each metric specific to a provider should go prefixed with the provider namespace.
    for key in ["provider.valueOne", "provider.valueTwo", "provider.valueThree"]:
        metric[key] = random.randint(0, 100)
        logger.debug("Adding metric %s with value %d", key, metric[key])

    data.addMetric(metric)

    key_list = ["valueOne", "valueTwo", "valueThree"]
    item_keys = ["item1", "item2", "item3"]
    for item in item_keys:
        for key in key_list:
            data.addInventory(item, key, random.randint(0, 100))
            logger.debug("Set inventory key %s=%d for %s", key, data.inventory[item][key], item)

    if args.pretty:
        print json.dumps(data.__dict__, indent=4)
    else:
        print json.dumps(data.__dict__)
