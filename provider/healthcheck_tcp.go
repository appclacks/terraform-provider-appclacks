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
	resHealthcheckTarget        = "target"
	resHealthcheckPort          = "port"
	resHealthcheckTCPShouldFail = "should_fail"
)

func resourceHealthcheckTCP() *schema.Resource {
	return &schema.Resource{
		Description: "Create a TCP connection on the target",
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
			resHealthcheckTCPShouldFail: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the health check will be considered successful if the TCP connection fails",
			},
		},

		CreateContext: resourceHealthcheckTCPCreate,
		ReadContext:   resourceHealthcheckTCPRead,
		UpdateContext: resourceHealthcheckTCPUpdate,
		DeleteContext: resourceHealthcheckTCPDelete,

		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceHealthcheckTCPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	update := apitypes.UpdateTCPHealthcheckInput{
		ID:       d.Id(),
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckTCPDefinition: apitypes.HealthcheckTCPDefinition{
			Target: d.Get(resHealthcheckTarget).(string),
			Port:   uint(d.Get(resHealthcheckPort).(int)),
		},
	}

	if v, ok := d.GetOk(resHealthcheckDescription); ok {
		update.Description = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTCPShouldFail); ok {
		update.HealthcheckTCPDefinition.ShouldFail = v.(bool)
	}
	if l, ok := d.GetOk(resHealthcheckLabels); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		update.Labels = labels
	}

	if _, err := client.UpdateTCPHealthcheck(ctx, update); err != nil {
		return diag.FromErr(err)
	}

	return resourceHealthcheckTCPRead(ctx, d, meta)
}

func resourceHealthcheckTCPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceHealthcheckTCPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetAppclacksClient(meta)

	healthcheck := apitypes.CreateTCPHealthcheckInput{
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckTCPDefinition: apitypes.HealthcheckTCPDefinition{
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

	if v, ok := d.GetOk(resHealthcheckTCPShouldFail); ok {
		healthcheck.HealthcheckTCPDefinition.ShouldFail = v.(bool)
	}

	result, err := client.CreateTCPHealthcheck(ctx, healthcheck)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)
	return resourceHealthcheckTCPRead(ctx, d, meta)
}

func resourceHealthcheckTCPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	result, err := client.GetHealthcheck(ctx, apitypes.GetHealthcheckInput{
		Identifier: d.Id(),
	})
	if err != nil {
		if errors.Is(err, goclient.ErrNotFound) {
			// remove resource if does not exist
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}
	return diag.FromErr(resourceTCPHealthcheckApply(ctx, d, &result))
}

func resourceTCPHealthcheckApply(_ context.Context, d *schema.ResourceData, healthcheck *apitypes.Healthcheck) error {

	if healthcheck.Type != "tcp" {
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

	definition, ok := healthcheck.Definition.(apitypes.HealthcheckTCPDefinition)
	if !ok {
		return errors.New("Invalid healthcheck definition for TCP health check")
	}

	if err := d.Set(resHealthcheckTarget, definition.Target); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckPort, definition.Port); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckTCPShouldFail, definition.ShouldFail); err != nil {
		return err
	}

	return nil
}
