package main

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/newrelic/infra-integrations-sdk/sdk"
)

func getInventory() (map[string]interface{}, error) {
	rawYamlFile, err := ioutil.ReadFile(args.ConfigPath)
	if err != nil {
		return nil, err
	}

	inventory := make(map[string]interface{})
	err = yaml.Unmarshal(rawYamlFile, &inventory)
	if err != nil {
		return nil, err
	}

	return inventory, nil
}

func populateInventory(inventory map[string]sdk.Inventory, rawInventory map[string]interface{}) error {
	for k, v := range rawInventory {
		switch value := v.(type) {
		case map[interface{}]interface{}:
			inventory[k] = sdk.Inventory{}
			for subk, subv := range value {
				inventory[k][subk.(string)] = subv
			}
		case []interface{}:
			//TODO: Do not include lists for now
		default:
			inventory[k] = sdk.Inventory{"value": value}
		}
	}
	return nil
}
