package provider

import (
	"context"
	"errors"
	"fmt"

	goclient "github.com/appclacks/go-client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	resHealthcheckCommandCommand   = "command"
	resHealthcheckCommandArguments = "arguments"
)

func resourceHealthcheckCommand() *schema.Resource {
	return &schema.Resource{
		Description: "Execute an arbitrary shell command",
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
			resHealthcheckCommandCommand: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Health check command",
			},
			resHealthcheckCommandArguments: {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Command arguments",
			},
		},

		CreateContext: resourceHealthcheckCommandCreate,
		ReadContext:   resourceHealthcheckCommandRead,
		UpdateContext: resourceHealthcheckCommandUpdate,
		DeleteContext: resourceHealthcheckCommandDelete,

		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceHealthcheckCommandUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	update := goclient.UpdateCommandHealthcheckInput{
		ID:       d.Id(),
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  false,
		HealthcheckCommandDefinition: goclient.HealthcheckCommandDefinition{
			Command: d.Get(resHealthcheckCommandCommand).(string),
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

	if set, ok := d.Get(resHealthcheckCommandArguments).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]string, l)
			for i, v := range set.List() {
				list[i] = v.(string)
			}
			update.HealthcheckCommandDefinition.Arguments = list

		}
	}

	if _, err := client.UpdateCommandHealthcheck(ctx, update); err != nil {
		return diag.FromErr(err)
	}

	return resourceHealthcheckCommandRead(ctx, d, meta)
}

func resourceHealthcheckCommandDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	healthcheckID := d.Id()
	_, err := client.DeleteHealthcheck(ctx, goclient.DeleteHealthcheckInput{ID: healthcheckID})
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceHealthcheckCommandCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetAppclacksClient(meta)

	healthcheck := goclient.CreateCommandHealthcheckInput{
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  false,
		HealthcheckCommandDefinition: goclient.HealthcheckCommandDefinition{
			Command: d.Get(resHealthcheckCommandCommand).(string),
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

	if set, ok := d.Get(resHealthcheckCommandArguments).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]string, l)
			for i, v := range set.List() {
				list[i] = v.(string)
			}
			healthcheck.HealthcheckCommandDefinition.Arguments = list

		}
	}

	result, err := client.CreateCommandHealthcheck(ctx, healthcheck)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)
	return resourceHealthcheckCommandRead(ctx, d, meta)
}

func resourceHealthcheckCommandRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	result, err := client.GetHealthcheck(ctx, goclient.GetHealthcheckInput{
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
	return diag.FromErr(resourceCommandHealthcheckApply(ctx, d, &result))
}

func resourceCommandHealthcheckApply(_ context.Context, d *schema.ResourceData, healthcheck *goclient.Healthcheck) error {

	if healthcheck.Type != "command" {
		return fmt.Errorf("Invalid healthcheck type. Expecting command, got %s", healthcheck.Type)
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

	definition, ok := healthcheck.Definition.(goclient.HealthcheckCommandDefinition)
	if !ok {
		return errors.New("Invalid healthcheck definition for Command health check")
	}

	if err := d.Set(resHealthcheckCommandCommand, definition.Command); err != nil {
		return err
	}

	if len(definition.Arguments) != 0 {
		if err := d.Set(resHealthcheckCommandArguments, definition.Arguments); err != nil {
			return err
		}
	}

	return nil
}
