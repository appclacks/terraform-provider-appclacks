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
	resHealthcheckHTTPQuery       = "query"
	resHealthcheckHTTPProtocol    = "protocol"
	resHealthcheckHTTPPath        = "path"
	resHealthcheckHTTPHost        = "host"
)

func resourceHealthcheckHTTP() *schema.Resource {
	return &schema.Resource{
		Description: "Execute an HTTP request on the target",
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
			resHealthcheckHTTPValidStatus: {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Description: "Expected status code(s) for the HTTP response",
			},
			resHealthcheckPort: {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Health check port",
			},
			resHealthcheckHTTPMethod: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "GET",
				Description: "Health check HTTP method",
			},
			resHealthcheckHTTPProtocol: {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "https",
				Description: "Health check protocol to use (http or https)",
			},
			resHealthcheckHTTPPath: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Health check request HTTP path",
			},
			resHealthcheckHTTPRedirect: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Follow redirections",
			},
			resHealthcheckHTTPBody: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Health check request HTTP body",
			},
			resHealthcheckHTTPHeaders: {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Health check request HTTP headers",
			},
			resHealthcheckHTTPQuery: {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Health check request HTTP query parameters",
			},
			resHealthcheckHTTPBodyRegexp: {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of regular expression which will be executed against the response body",
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
			resHealthcheckHTTPHost: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Host header to use for the health check HTTP request",
			},
			resHealthcheckTLSServerName: {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Server name to use for the TLS connection",
			},
			resHealthcheckTLSInsecure: {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Accept insecure TLS connections",
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
	if v, ok := d.GetOk(resHealthcheckHTTPHost); ok {
		update.HealthcheckHTTPDefinition.Host = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSInsecure); ok {
		update.HealthcheckHTTPDefinition.Insecure = v.(bool)
	}
	if v, ok := d.GetOk(resHealthcheckTLSServerName); ok {
		update.HealthcheckHTTPDefinition.ServerName = v.(string)
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
	if l, ok := d.GetOk(resHealthcheckHTTPQuery); ok {
		query := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			query[k] = v.(string)
		}
		update.HealthcheckHTTPDefinition.Query = query
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
	if v, ok := d.GetOk(resHealthcheckHTTPHost); ok {
		healthcheck.HealthcheckHTTPDefinition.Host = v.(string)
	}
	if v, ok := d.GetOk(resHealthcheckTLSInsecure); ok {
		healthcheck.HealthcheckHTTPDefinition.Insecure = v.(bool)
	}
	if v, ok := d.GetOk(resHealthcheckTLSServerName); ok {
		healthcheck.HealthcheckHTTPDefinition.ServerName = v.(string)
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
	if l, ok := d.GetOk(resHealthcheckHTTPQuery); ok {
		query := make(map[string]string)
		for k, v := range l.(map[string]interface{}) {
			query[k] = v.(string)
		}
		healthcheck.HealthcheckHTTPDefinition.Query = query
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
	if err := d.Set(resHealthcheckHTTPQuery, definition.Query); err != nil {
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
	if err := d.Set(resHealthcheckHTTPHost, definition.Host); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckTLSInsecure, definition.Insecure); err != nil {
		return err
	}
	if err := d.Set(resHealthcheckTLSServerName, definition.ServerName); err != nil {
		return err
	}

	return nil
}
