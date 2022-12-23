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
	resHealthcheckHTTPMethod      = "method"
	resHealthcheckHTTPRedirect    = "redirect"
	resHealthcheckHTTPBody        = "body"
	resHealthcheckHTTPBodyRegexp  = "body_regexp"
	resHealthcheckHTTPValidStatus = "valid_status"
	resHealthcheckHTTPHeaders     = "headers"
	resHealthcheckHTTPProtocol    = "protocol"
	resHealthcheckHTTPPath        = "path"
)

func resourceHealthcheckHTTP() *schema.Resource {
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
				Default:  defaultHealthcheckInterval,
			},
			resHealthcheckTimeout: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  defaultHealthcheckTimeout,
			},
			resHealthcheckEnabled: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			resHealthcheckTarget: {
				Type:     schema.TypeString,
				Required: true,
			},
			resHealthcheckHTTPValidStatus: {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			resHealthcheckPort: {
				Type:     schema.TypeInt,
				Required: true,
			},
			resHealthcheckHTTPMethod: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "GET",
			},
			resHealthcheckHTTPProtocol: {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "https",
			},
			resHealthcheckHTTPPath: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resHealthcheckHTTPRedirect: {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			resHealthcheckHTTPBody: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resHealthcheckHTTPHeaders: {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			resHealthcheckHTTPBodyRegexp: {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			resHealthcheckTLSKey: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resHealthcheckTLSCert: {
				Type:     schema.TypeString,
				Optional: true,
			},
			resHealthcheckTLSCacert: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		CreateContext: resourceHealthcheckHTTPCreate,
		ReadContext:   resourceHealthcheckHTTPRead,
		UpdateContext: resourceHealthcheckHTTPUpdate,
		DeleteContext: resourceHealthcheckHTTPDelete,

		Importer: &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(defaultTimeout),
			Read:   schema.DefaultTimeout(defaultTimeout),
			Update: schema.DefaultTimeout(defaultTimeout),
			Delete: schema.DefaultTimeout(defaultTimeout),
		},
	}
}

func resourceHealthcheckHTTPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()
	client := GetAppclacksClient(meta)

	update := apitypes.UpdateHTTPHealthcheckInput{
		ID:       d.Id(),
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckHTTPDefinition: apitypes.HealthcheckHTTPDefinition{
			Target:   d.Get(resHealthcheckTarget).(string),
			Port:     uint(d.Get(resHealthcheckPort).(int)),
			Method:   d.Get(resHealthcheckHTTPMethod).(string),
			Protocol: d.Get(resHealthcheckHTTPProtocol).(string),
		},
	}

	if l, ok := d.GetOk(resHealthcheckLabels); ok {
		labels := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			labels[k] = v.(string)
		}
		update.Labels = labels
	}

	if set, ok := d.Get(resHealthcheckHTTPValidStatus).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]uint, l)
			for i, v := range set.List() {
				list[i] = uint(v.(int))
			}
			update.HealthcheckHTTPDefinition.ValidStatus = list

		}
	}

	if v, ok := d.GetOk(resHealthcheckDescription); ok {
		update.Description = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckHTTPPath); ok {
		update.HealthcheckHTTPDefinition.Path = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckHTTPRedirect); ok {
		update.HealthcheckHTTPDefinition.Redirect = v.(bool)
	}
	if v, ok := d.GetOk(resHealthcheckHTTPBody); ok {
		update.HealthcheckHTTPDefinition.Body = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCert); ok {
		update.HealthcheckHTTPDefinition.Cert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCacert); ok {
		update.HealthcheckHTTPDefinition.Cacert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSKey); ok {
		update.HealthcheckHTTPDefinition.Key = v.(string)
	}
	if set, ok := d.Get(resHealthcheckHTTPBodyRegexp).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]string, l)
			for i, v := range set.List() {
				list[i] = v.(string)
			}
			update.HealthcheckHTTPDefinition.BodyRegexp = list

		}
	}
	if l, ok := d.GetOk(resHealthcheckHTTPHeaders); ok {
		headers := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			headers[k] = v.(string)
		}
		update.HealthcheckHTTPDefinition.Headers = headers
	}

	if _, err := client.UpdateHTTPHealthcheck(ctx, update); err != nil {
		return diag.FromErr(err)
	}

	return resourceHealthcheckHTTPRead(ctx, d, meta)
}

func resourceHealthcheckHTTPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

func resourceHealthcheckHTTPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ctx, cancel := context.WithTimeout(ctx, d.Timeout(schema.TimeoutCreate))
	defer cancel()

	client := GetAppclacksClient(meta)

	healthcheck := apitypes.CreateHTTPHealthcheckInput{
		Name:     d.Get(resHealthcheckName).(string),
		Interval: d.Get(resHealthcheckInterval).(string),
		Timeout:  d.Get(resHealthcheckTimeout).(string),
		Enabled:  d.Get(resHealthcheckEnabled).(bool),
		HealthcheckHTTPDefinition: apitypes.HealthcheckHTTPDefinition{
			Target:   d.Get(resHealthcheckTarget).(string),
			Port:     uint(d.Get(resHealthcheckPort).(int)),
			Method:   d.Get(resHealthcheckHTTPMethod).(string),
			Protocol: d.Get(resHealthcheckHTTPProtocol).(string),
		},
	}

	if v, ok := d.GetOk(resHealthcheckDescription); ok {
		healthcheck.Description = v.(string)
	}
	if set, ok := d.Get(resHealthcheckHTTPValidStatus).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]uint, l)
			for i, v := range set.List() {
				list[i] = uint(v.(int))
			}
			healthcheck.HealthcheckHTTPDefinition.ValidStatus = list

		}
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
	if v, ok := d.GetOk(resHealthcheckHTTPPath); ok {
		healthcheck.HealthcheckHTTPDefinition.Path = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckHTTPRedirect); ok {
		healthcheck.HealthcheckHTTPDefinition.Redirect = v.(bool)
	}
	if v, ok := d.GetOk(resHealthcheckHTTPBody); ok {
		healthcheck.HealthcheckHTTPDefinition.Body = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCert); ok {
		healthcheck.HealthcheckHTTPDefinition.Cert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSCacert); ok {
		healthcheck.HealthcheckHTTPDefinition.Cacert = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSKey); ok {
		healthcheck.HealthcheckHTTPDefinition.Key = v.(string)
	}
	if set, ok := d.Get(resHealthcheckHTTPBodyRegexp).(*schema.Set); ok {
		if l := set.Len(); l > 0 {
			list := make([]string, l)
			for i, v := range set.List() {
				list[i] = v.(string)
			}
			healthcheck.HealthcheckHTTPDefinition.BodyRegexp = list

		}
	}
	if l, ok := d.GetOk(resHealthcheckHTTPHeaders); ok {
		headers := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			headers[k] = v.(string)
		}
		healthcheck.HealthcheckHTTPDefinition.Headers = headers
	}

	result, err := client.CreateHTTPHealthcheck(ctx, healthcheck)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(result.ID)
	return resourceHealthcheckHTTPRead(ctx, d, meta)
}

func resourceHealthcheckHTTPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return diag.FromErr(resourceHTTPHealthcheckApply(ctx, d, &result))
}

func resourceHTTPHealthcheckApply(_ context.Context, d *schema.ResourceData, healthcheck *apitypes.Healthcheck) error {

	if healthcheck.Type != "http" {
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

	definition, ok := healthcheck.Definition.(apitypes.HealthcheckHTTPDefinition)
	if !ok {
		return errors.New("Invalid healthcheck definition for HTTP health check")
	}

	if err := d.Set(resHealthcheckHTTPPath, definition.Path); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckTarget, definition.Target); err != nil {
		return err
	}

	if err := d.Set(resHealthcheckPort, definition.Port); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckHTTPMethod, definition.Method); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckHTTPProtocol, definition.Protocol); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckHTTPRedirect, definition.Redirect); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckHTTPValidStatus, definition.ValidStatus); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckHTTPBody, definition.Body); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckHTTPBodyRegexp, definition.BodyRegexp); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckHTTPHeaders, definition.Headers); err != nil {
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

	return nil
}
