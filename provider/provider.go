package provider

import (
	"context"
	"time"

	"github.com/appclacks/cli/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultTimeout = 10 * time.Second
	defaultAPIURL  = "https://api.appclacks.com"
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
			"appclacks_healthcheck_dns": resourceHealthcheckDNS(),
		},
	}
}

// marshalData is used to ensure the data is put into a format Terraform can output
func marshalData(d *schema.ResourceData, vals map[string]interface{}) {
	for k, v := range vals {
		if k == "id" {
			d.SetId(v.(string))
		} else {
			str, ok := v.(string)
			if ok {
				d.Set(k, str)
			} else {
				d.Set(k, v)
			}
		}
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
