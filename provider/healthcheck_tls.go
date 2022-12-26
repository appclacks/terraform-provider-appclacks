package provider

import (
	"context"
	"errors"
	"fmt"

	goclient "github.com/appclacks/cli/client"
	apitypes "github.com/appclacks/go-types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	resHealthcheckTLSKey             = "key"
	resHealthcheckTLSCert            = "cert"
	resHealthcheckTLSCacert          = "cacert"
	resHealthcheckTLSServerName      = "server_name"
	resHealthcheckTLSInsecure        = "insecure"
	resHealthcheckTLSExpirationDelay = "expiration_delay"
)

func resourceHealthcheckTLS() *schema.Resource {
	return &schema.Resource{
		Description: "Create a TLS connection on the target",
		Schema: map[string]*schema.Schema{
			resHealthcheckName: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Health check name",
			},
			resHealthcheckDescription: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Health check description",
			},
			resHealthcheckLabels: {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Health check labels",
			},
			resHealthcheckInterval: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultHealthcheckInterval,
				Description: "Health check interval (example: 30s)",
			},
			resHealthcheckTimeout: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     defaultHealthcheckTimeout,
				Description: "Health check timeout (example: 5s)",
			},
			resHealthcheckEnabled: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable the health check on the Appclacks platform",
			},
			resHealthcheckTarget: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Health check target (can be a domain or an IP address)",
			},
			resHealthcheckPort: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Health check port",
			},
			resHealthcheckTLSKey: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TLS key file to use for the TLS connection",
			},
			resHealthcheckTLSCert: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TLS cert file to use for the TLS connection",
			},
			resHealthcheckTLSCacert: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "TLS cacert file to use for the TLS connection",
			},
			resHealthcheckTLSServerName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server name to use for the TLS connection. Mandatory if insecure is not set.",
			},
			resHealthcheckTLSInsecure: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Accept insecure TLS connections",
			},
			resHealthcheckTLSExpirationDelay: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The health check will be considered failed if hte certificate expires is less than this duration (for example: 168h)",
			},
		},

		CreateContext: resourceHealthcheckTLSCreate,
		ReadContext:   resourceHealthcheckTLSRead,
		UpdateContext: resourceHealthcheckTLSUpdate,
		DeleteContext: resourceHealthcheckTLSDelete,

		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceHealthcheckTLSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	update := apitypes.UpdateTLSHealthcheckInput{
		ID:       d.Id(),
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckTLSDefinition: apitypes.HealthcheckTLSDefinition{
			Target: d.Get(resHealthcheckTarget).(string),
			Port:   uint(d.Get(resHealthcheckPort).(int)),
		},
	}

	if l, ok := d.GetOk(resHealthcheckLabels); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		update.Labels = labels
	}

	if v, ok := d.GetOk(resHealthcheckDescription); ok {
		update.Description = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSKey); ok {
		update.HealthcheckTLSDefinition.Key = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCert); ok {
		update.HealthcheckTLSDefinition.Cert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCacert); ok {
		update.HealthcheckTLSDefinition.Cacert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSServerName); ok {
		update.HealthcheckTLSDefinition.ServerName = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSInsecure); ok {
		update.HealthcheckTLSDefinition.Insecure = v.(bool)
	}
	if v, ok := d.GetOk(resHealthcheckTLSExpirationDelay); ok {
		update.HealthcheckTLSDefinition.ExpirationDelay = v.(string)
	}

	if _, err := client.UpdateTLSHealthcheck(ctx, update); err != nil {
		return diag.FromErr(err)
	}

	return resourceHealthcheckTLSRead(ctx, d, meta)
}

func resourceHealthcheckTLSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	healthcheckID := d.Id()
	_, err := client.DeleteHealthcheck(ctx, apitypes.DeleteHealthcheckInput{ID: healthcheckID})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceHealthcheckTLSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetAppclacksClient(meta)

	healthcheck := apitypes.CreateTLSHealthcheckInput{
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckTLSDefinition: apitypes.HealthcheckTLSDefinition{
			Target: d.Get(resHealthcheckTarget).(string),
			Port:   uint(d.Get(resHealthcheckPort).(int)),
		},
	}

	if v, ok := d.GetOk(resHealthcheckDescription); ok {
		healthcheck.Description = v.(string)
	}
	if l, ok := d.GetOk(resHealthcheckLabels); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		healthcheck.Labels = labels
	}

	if v, ok := d.GetOk(resHealthcheckDescription); ok {
		healthcheck.Description = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSKey); ok {
		healthcheck.HealthcheckTLSDefinition.Key = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCert); ok {
		healthcheck.HealthcheckTLSDefinition.Cert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCacert); ok {
		healthcheck.HealthcheckTLSDefinition.Cacert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSServerName); ok {
		healthcheck.HealthcheckTLSDefinition.ServerName = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSInsecure); ok {
		healthcheck.HealthcheckTLSDefinition.Insecure = v.(bool)
	}
	if v, ok := d.GetOk(resHealthcheckTLSExpirationDelay); ok {
		healthcheck.HealthcheckTLSDefinition.ExpirationDelay = v.(string)
	}

	result, err := client.CreateTLSHealthcheck(ctx, healthcheck)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)
	return resourceHealthcheckTLSRead(ctx, d, meta)
}

func resourceHealthcheckTLSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	result, err := client.GetHealthcheck(ctx, apitypes.GetHealthcheckInput{
		ID: d.Id(),
	})
	if err != nil {
		if errors.Is(err, goclient.ErrNotFound) {
			// remove resource if does not exist
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	return diag.FromErr(resourceTLSHealthcheckApply(ctx, d, &result))
}

func resourceTLSHealthcheckApply(_ context.Context, d *schema.ResourceData, healthcheck *apitypes.Healthcheck) error {

	if healthcheck.Type != "tls" {
		return fmt.Errorf("Invalid healthcheck type. Expecting tcp, got %s", healthcheck.Type)
	}

	if err := d.Set(resHealthcheckName, healthcheck.Name); err != nil {
		return err
	}

	if healthcheck.Description != "" {
		if err := d.Set(resHealthcheckDescription, healthcheck.Description); err != nil {
			return err
		}
	}

	if err := d.Set(resHealthcheckLabels, healthcheck.Labels); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckInterval, healthcheck.Interval); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckTimeout, healthcheck.Timeout); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckEnabled, healthcheck.Enabled); err != nil {
		return err
	}

	definition, ok := healthcheck.Definition.(apitypes.HealthcheckTLSDefinition)
	if !ok {
		return errors.New("Invalid healthcheck definition for TLS health check")
	}

	if err := d.Set(resHealthcheckTarget, definition.Target); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckPort, definition.Port); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckTLSKey, definition.Key); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckTLSCert, definition.Cert); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckTLSCacert, definition.Cacert); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckTLSServerName, definition.ServerName); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckTLSInsecure, definition.Insecure); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckTLSExpirationDelay, definition.ExpirationDelay); err != nil {
		return err
	}

	return nil
}
