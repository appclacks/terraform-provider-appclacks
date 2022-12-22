package provider

import (
	"context"
	"errors"
	"fmt"

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
		Schema: map[string]*schema.Schema{
			resHealthcheckName: {
				Type:     schema.TypeString,
				Required: true,
			},
			resHealthcheckDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resHealthcheckLabels: {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			resHealthcheckInterval: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "60s",
			},
			resHealthcheckTimeout: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "5s",
			},
			resHealthcheckEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			resHealthcheckDNSDomain: {
				Type:     schema.TypeString,
				Required: true,
			},
			resHealthcheckDNSExpectedIPs: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},

		CreateContext: resourceHealthcheckDNSCreate,
		ReadContext:   resourceHealthcheckDNSRead,
		UpdateContext: resourceHealthcheckDNSUpdate,
		DeleteContext: resourceHealthcheckDNSDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceHealthcheckDNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := GetAppclacksClient(meta)

	_, err := client.GetHealthcheck(apitypes.GetHealthcheckInput{
		ID: d.Id(),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	update := apitypes.UpdateDNSHealthcheckInput{}

	var updated bool

	if d.HasChange(resHealthcheckName) {
		v := d.Get(resHealthcheckName).(string)
		update.Name = v
		updated = true
	}

	if d.HasChange(resHealthcheckDescription) {
		v := d.Get(resHealthcheckDescription).(string)
		update.Description = v
		updated = true
	}

	if d.HasChange(resHealthcheckLabels) {
		labels := make(map[string]string)
		for k, v := range d.Get(resHealthcheckLabels).(map[string]interface{}) {
			labels[k] = v.(string)
		}
		update.Labels = labels
		updated = true
	}

	if d.HasChange(resHealthcheckInterval) {
		v := d.Get(resHealthcheckInterval).(string)
		update.Interval = v
		updated = true
	}

	if d.HasChange(resHealthcheckTimeout) {
		v := d.Get(resHealthcheckTimeout).(string)
		update.Timeout = v
		updated = true
	}

	if d.HasChange(resHealthcheckEnabled) {
		v := d.Get(resHealthcheckEnabled).(bool)
		update.Enabled = v
		updated = true
	}
	if d.HasChange(resHealthcheckDNSDomain) {
		v := d.Get(resHealthcheckDNSDomain).(string)
		update.Domain = v
		updated = true
	}

	if d.HasChange(resHealthcheckDNSExpectedIPs) {
		set := d.Get(resHealthcheckDNSExpectedIPs).(*schema.Set)
		list := []string{}
		for _, v := range set.List() {
			list = append(list, v.(string))
		}
		update.ExpectedIPs = list
		updated = true
	}

	if updated {
		if _, err = client.UpdateDNSHealthcheck(update); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceHealthcheckDNSRead(ctx, d, meta)
}

func resourceHealthcheckDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	client := GetAppclacksClient(meta)

	healthcheckID := d.Id()
	_, err := client.DeleteHealthcheck(apitypes.DeleteHealthcheckInput{ID: healthcheckID})
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

	result, err := client.CreateDNSHealthcheck(healthcheck)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)
	return resourceHealthcheckDNSRead(ctx, d, meta)
}

func resourceHealthcheckDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := GetAppclacksClient(meta)

	result, err := client.GetHealthcheck(apitypes.GetHealthcheckInput{
		ID: d.Id(),
	})
	if err != nil {
		// todo check if 404
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

	if len(healthcheck.Labels) != 0 {
		if err := d.Set(resHealthcheckLabels, healthcheck.Labels); err != nil {
			return err
		}
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
