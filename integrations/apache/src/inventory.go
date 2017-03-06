package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/newrelic/infra-integrations-sdk/sdk"
)

// getInventory executes system command in order to retrieve required inventory data and calls functions which parse the result.
// It returns a map of inventory data
func getInventory() (map[string]sdk.Inventory, error) {
	inventory := make(map[string]sdk.Inventory)

	cmd := exec.Command("/usr/sbin/httpd", "-M")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	r := bytes.NewReader(output)
	err = getModules(bufio.NewReader(r), inventory)
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("/usr/sbin/httpd", "-V")
	output, err = cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	r = bytes.NewReader(output)
	err = getVersion(bufio.NewReader(r), inventory)
	if err != nil {
		return nil, err
	}

	if len(inventory) == 0 {
		return nil, fmt.Errorf("Empty result")
	}
	return inventory, nil
}

// getModules reads an Apache list of enabled modules and transforms its
// contents into a map that can be processed by NR agent.
// It appends a map of inventory data where the keys contain name of the module and values
// indicate that module is enabled.
func getModules(reader *bufio.Reader, inventory map[string]sdk.Inventory) error {
	modulesValue := sdk.Inventory{
		"value": "enabled",
	}

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if strings.Contains(line, "_module") {
			splitedLine := strings.Split(line, " ")
			moduleName := splitedLine[1]
			inventory[fmt.Sprintf("modules/%s", moduleName[:len(moduleName)-7])] = modulesValue
		}
	}

	return nil
}

// getVersion reads an Apache list of compile settings and transforms its
// contents into a map that can be processed by NR agent.
// It appends a map of inventory data which indicates Apache Server version
func getVersion(reader *bufio.Reader, inventory map[string]sdk.Inventory) error {
	versionValue := make(sdk.Inventory)

	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if strings.Contains(line, "Server version") {
			splitedLine := strings.Split(line, ":")
			versionValue["value"] = strings.TrimSpace(splitedLine[1])
			break
		}
	}
	if len(versionValue) != 0 {
		inventory["version"] = versionValue
	}

	return nil

}
