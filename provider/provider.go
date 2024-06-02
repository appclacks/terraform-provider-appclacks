package provider

import (
	"context"
	"time"

	"github.com/appclacks/go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	defaultTimeout             = 10 * time.Second
	defaultHealthcheckTimeout  = "10s"
	defaultHealthcheckInterval = "60s"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ConfigureContextFunc: providerConfigure,
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_cert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_cacert": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
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
	options := []client.ClientOption{}
	configEndpoint, ok := d.GetOk("api_endpoint")
	if ok {
		options = append(options, client.WithEndpoint(configEndpoint.(string)))
	}
	configUsername, ok := d.GetOk("username")
	if ok {
		options = append(options, client.WithUsername(configUsername.(string)))
	}
	configPassword, ok := d.GetOk("password")
	if ok {
		options = append(options, client.WithPassword(configPassword.(string)))
	}
	configTLSKey, ok := d.GetOk("tls_key")
	if ok {
		options = append(options, client.WithKey(configTLSKey.(string)))
	}
	configTLSCert, ok := d.GetOk("tls_cert")
	if ok {
		options = append(options, client.WithCert(configTLSCert.(string)))
	}
	configTLSCacert, ok := d.GetOk("tls_cacert")
	if ok {
		options = append(options, client.WithCacert(configTLSCacert.(string)))
	}
	configTLSInsecure, ok := d.GetOk("insecure")
	if ok {
		options = append(options, client.WithInsecure(configTLSInsecure.(bool)))
	}
	client, err := client.New(options...)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, nil
}

func GetAppclacksClient(meta interface{}) *client.Client {
	client := meta.(*client.Client)
	return client
}
