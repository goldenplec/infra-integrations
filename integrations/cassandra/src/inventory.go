package main

import (
	"io/ioutil"
<<<<<<< HEAD
=======
	"regexp"
>>>>>>> upstream/master

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

<<<<<<< HEAD
func populateInventory(inventory map[string]sdk.Inventory, rawInventory map[string]interface{}) error {
	for k, v := range rawInventory {
		switch value := v.(type) {
		case map[interface{}]interface{}:
			inventory[k] = sdk.Inventory{}
			for subk, subv := range value {
				inventory[k][subk.(string)] = subv
=======
func populateInventory(inventory sdk.Inventory, rawInventory map[string]interface{}) error {
	for k, v := range rawInventory {
		switch value := v.(type) {
		case map[interface{}]interface{}:
			for subk, subv := range value {
				switch subVal := subv.(type) {
				case []interface{}:
					//TODO: Do not include lists for now
				default:
					setValue(inventory, k, subk.(string), subVal)
				}
>>>>>>> upstream/master
			}
		case []interface{}:
			//TODO: Do not include lists for now
		default:
<<<<<<< HEAD
			inventory[k] = sdk.Inventory{"value": value}
=======
			setValue(inventory, k, "value", value)
>>>>>>> upstream/master
		}
	}
	return nil
}
<<<<<<< HEAD
=======

func setValue(inventory sdk.Inventory, key string, field string, value interface{}) {
	re, _ := regexp.Compile("(?i)password")

	if re.MatchString(key) || re.MatchString(field) {
		value = "(omitted value)"
	}
	inventory.SetItem(key, field, value)
}
>>>>>>> upstream/master
