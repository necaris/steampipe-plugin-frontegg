package frontegg

import (
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/schema"
)

type fronteggConfig struct {
	ClientID  *string `cty:"clientId"`
	SecretKey *string `cty:"secretKey"`
}

var ConfigSchema = map[string]*schema.Attribute{
	"clientId": {
		Type: schema.TypeString,
	},
	"secretKey": {
		Type: schema.TypeString,
	},
}

func ConfigInstance() interface{} {
	return &fronteggConfig{}
}

// GetConfig :: retrieve and cast connection config from query data
func GetConfig(connection *plugin.Connection) fronteggConfig {
	if connection == nil || connection.Config == nil {
		return fronteggConfig{}
	}
	config, _ := connection.Config.(fronteggConfig)
	return config
}
