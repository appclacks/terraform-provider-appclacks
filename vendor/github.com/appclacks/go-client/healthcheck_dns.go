package client

import (
	"context"
	"fmt"
	"net/http"
)

type HealthcheckDNSDefinition struct {
	Domain      string   `json:"domain,omitempty" description:"Domain to check" validate:"required,max=255,min=1"`
	ExpectedIPs []string `json:"expected-ips,omitempty" description:"Expected IP addresses in the answer" validate:"max=10,dive,ip_addr"`
}

type CreateDNSHealthcheckInput struct {
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Enabled     bool              `json:"bool" description:"Enable the healthcheck on the appclacks platform"`
	Timeout     string            `json:"timeout" validate:"required"`
	HealthcheckDNSDefinition
}

type UpdateDNSHealthcheckInput struct {
	ID          string            `json:"-" param:"id" description:"Healthcheck ID" validate:"required,uuid"`
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Timeout     string            `json:"timeout" validate:"required"`
	Enabled     bool              `json:"enabled" description:"Enable the healthcheck on the appclacks platform"`
	HealthcheckDNSDefinition
}

func (c *Client) CreateDNSHealthcheck(ctx context.Context, input CreateDNSHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	_, err := c.sendRequest(ctx, "/api/v1/healthcheck/dns", http.MethodPost, input, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}

func (c *Client) UpdateDNSHealthcheck(ctx context.Context, input UpdateDNSHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	internalInput := internalUpdateHealthcheckInput{
		Name:        input.Name,
		Description: input.Description,
		Labels:      input.Labels,
		Interval:    input.Interval,
		Enabled:     input.Enabled,
		Timeout:     input.Timeout,
	}
	payload, err := jsonMerge(internalInput, input.HealthcheckDNSDefinition)
	if err != nil {
		return result, err
	}
	_, err = c.sendRequest(ctx, fmt.Sprintf("/api/v1/healthcheck/dns/%s", input.ID), http.MethodPut, payload, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}
