# NRI Example Plugins

This example plugin will report fake metrics and inventory data to your NR account.

## Build

To build the example go plugin run the following commands:

$ cd example/go/

$ go build -o go-example-plugin

The python plugin needs execution permissions which can be set with the following command:

$ chmod +x example/python/py-example-plugin

# Install

Place the `example` directory in the Agent's Directory inside of the reserved directory for your plugins: plugins-custom: (/var/db/newrelic-infra/plugins-custom/ on Linux OSs).

To activate your plugin, symlink the config file into the Agent's Directory inside of the reserved directory for active plugin configs: plugins.d (/var/db/newrelic-infra/plugins.d/ on Linux).

# Test the plugin

Plugins are stand-alone executables that can be run outside the infrastructure agent from the command line. To test them, from the `example` directory, run:

$ ./go/go-example-plugin -v -p

$ ./python/py-example-plugin -v -p
