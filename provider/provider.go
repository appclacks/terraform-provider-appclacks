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
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("APPCLACKS_API_URL", defaultAPIURL),
			},
			"organization_id": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("APPCLACKS_ORGANIZATION_ID", ""),
				Description: "The organization ID to use for the Appclacks API",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("APPCLACKS_TOKEN", ""),
				Description: "The token to use for the Appclacks API",
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

	return client.New(
		d.Get("api_url").(string),
		client.WithToken(
			client.OrganizationID(d.Get("organization_id").(string)),
			client.Token(d.Get("token").(string)),
		),
	), nil

}

func GetAppclacksClient(meta interface{}) *client.Client {
	client := meta.(*client.Client)
	return client
}
