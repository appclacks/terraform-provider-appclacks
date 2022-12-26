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
	resHealthcheckName           = "name"
	resHealthcheckDescription    = "description"
	resHealthcheckLabels         = "labels"
	resHealthcheckInterval       = "interval"
	resHealthcheckTimeout        = "timeout"
	resHealthcheckEnabled        = "enabled"
	resHealthcheckDNSDomain      = "domain"
	resHealthcheckDNSExpectedIPs = "expected_ips"
)

func resourceHealthcheckDNS() *schema.Resource {
	return &schema.Resource{
		Description: "Execute a DNS query and optionally verify the request answer",
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
			resHealthcheckDNSDomain: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Domain to check",
			},
			resHealthcheckDNSExpectedIPs: {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Expected IP addresses in the answer",
			},
		},

		CreateContext: resourceHealthcheckDNSCreate,
		ReadContext:   resourceHealthcheckDNSRead,
		UpdateContext: resourceHealthcheckDNSUpdate,
		DeleteContext: resourceHealthcheckDNSDelete,

		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceHealthcheckDNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	update := apitypes.UpdateDNSHealthcheckInput{
		ID:       d.Id(),
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckDNSDefinition: apitypes.HealthcheckDNSDefinition{
			Domain: d.Get(resHealthcheckDNSDomain).(string),
		},
	}

	if v, ok := d.GetOk(resHealthcheckDescription); ok {
		update.Description = v.(string)
	}
	if l, ok := d.GetOk(resHealthcheckLabels); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		update.Labels = labels
	}

	if set, ok := d.Get(resHealthcheckDNSExpectedIPs).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]string, l)
			for i, v := range set.List() {
				list[i] = v.(string)
			}
			update.HealthcheckDNSDefinition.ExpectedIPs = list

		}
	}

	if _, err := client.UpdateDNSHealthcheck(ctx, update); err != nil {
		return diag.FromErr(err)
	}

	return resourceHealthcheckDNSRead(ctx, d, meta)
}

func resourceHealthcheckDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceHealthcheckDNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetAppclacksClient(meta)

	healthcheck := apitypes.CreateDNSHealthcheckInput{
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckDNSDefinition: apitypes.HealthcheckDNSDefinition{
			Domain: d.Get(resHealthcheckDNSDomain).(string),
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

	if set, ok := d.Get(resHealthcheckDNSExpectedIPs).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]string, l)
			for i, v := range set.List() {
				list[i] = v.(string)
			}
			healthcheck.ExpectedIPs = list

		}
	}

	result, err := client.CreateDNSHealthcheck(ctx, healthcheck)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)
	return resourceHealthcheckDNSRead(ctx, d, meta)
}

func resourceHealthcheckDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return diag.FromErr(resourceDNSHealthcheckApply(ctx, d, &result))
}

func resourceDNSHealthcheckApply(_ context.Context, d *schema.ResourceData, healthcheck *apitypes.Healthcheck) error {

	if healthcheck.Type != "dns" {
		return fmt.Errorf("Invalid healthcheck type. Expecting dns, got %s", healthcheck.Type)
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

	definition, ok := healthcheck.Definition.(apitypes.HealthcheckDNSDefinition)
	if !ok {
		return errors.New("Invalid healthcheck definition for DNS health check")
	}

	if err := d.Set(resHealthcheckDNSDomain, definition.Domain); err != nil {
		return err
	}

	if len(definition.ExpectedIPs) != 0 {
		if err := d.Set(resHealthcheckDNSExpectedIPs, definition.ExpectedIPs); err != nil {
			return err
		}
	}

	return nil
}
