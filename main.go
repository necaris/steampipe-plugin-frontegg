package main

import (
	"github.com/necaris/steampipe-plugin-frontegg/frontegg"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{PluginFunc: frontegg.Plugin})
}
