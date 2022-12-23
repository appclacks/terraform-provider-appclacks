package provider

import (
	"context"
	"time"

	"github.com/appclacks/cli/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultTimeout             = 10 * time.Second
	defaultAPIURL              = "https://api.appclacks.com"
	defaultHealthcheckTimeout  = "10s"
	defaultHealthcheckInterval = "60s"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,

		Schema: map[string]*schema.Schema{
			"api_url": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  defaultAPIURL,
			},
		},

		DataSourcesMap: map[string]*schema.Resource{},

		ResourcesMap: map[string]*schema.Resource{
			"appclacks_healthcheck_dns":     resourceHealthcheckDNS(),
			"appclacks_healthcheck_tcp":     resourceHealthcheckTCP(),
			"appclacks_healthcheck_tls":     resourceHealthcheckTLS(),
			"appclacks_healthcheck_http":    resourceHealthcheckHTTP(),
			"appclacks_healthcheck_command": resourceHealthcheckCommand(),
		},
	}
}

// providerConfigure parses the config into the Terraform provider meta object
func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiURL := d.Get("api_url").(string)
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	client := client.New(apiURL)

	return client, nil
}

func GetAppclacksClient(meta interface{}) *client.Client {
	client := meta.(*client.Client)
	return client
}
