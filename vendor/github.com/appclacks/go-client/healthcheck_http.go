package client

import (
	"context"
	"fmt"
	"net/http"
)

type HealthcheckHTTPDefinition struct {
	ValidStatus []uint            `json:"valid-status" validate:"required,min=1,max=20,dive,max=1000"`
	Target      string            `json:"target" validate:"required,max=255,min=1"`
	Method      string            `json:"method" validate:"required,oneof=GET POST PUT DELETE HEAD"`
	Port        uint              `json:"port" validate:"required,max=65535,min=1"`
	Host        string            `json:"host,omitempty"`
	Redirect    bool              `json:"redirect"`
	Query       map[string]string `json:"query,omitempty" validate:"max=20"`
	Body        string            `json:"body,omitempty"`
	BodyRegexp  []string          `json:"body-regexp,omitempty" validate:"max=3"`
	Headers     map[string]string `json:"headers,omitempty" validate:"max=20"`
	Protocol    string            `json:"protocol" validate:"oneof=http https"`
	Path        string            `json:"path,omitempty"`
	Key         string            `json:"key,omitempty"`
	Cert        string            `json:"cert,omitempty"`
	Cacert      string            `json:"cacert,omitempty"`
	Insecure    bool              `json:"insecure"`
	ServerName  string            `json:"server-name"`
}

type CreateHTTPHealthcheckInput struct {
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Enabled     bool              `json:"bool" description:"Enable the healthcheck on the appclacks platform"`
	Timeout     string            `json:"timeout" validate:"required"`
	HealthcheckHTTPDefinition
}

type UpdateHTTPHealthcheckInput struct {
	ID          string            `json:"-" param:"id" description:"Healthcheck ID" validate:"required,uuid"`
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Enabled     bool              `json:"enabled" description:"Enable the healthcheck on the appclacks platform"`
	Timeout     string            `json:"timeout" validate:"required"`
	HealthcheckHTTPDefinition
}

func (c *Client) CreateHTTPHealthcheck(ctx context.Context, input CreateHTTPHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	_, err := c.sendRequest(ctx, "/api/v1/healthcheck/http", http.MethodPost, input, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}

func (c *Client) UpdateHTTPHealthcheck(ctx context.Context, input UpdateHTTPHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	internalInput := internalUpdateHealthcheckInput{
		Name:        input.Name,
		Description: input.Description,
		Labels:      input.Labels,
		Interval:    input.Interval,
		Enabled:     input.Enabled,
		Timeout:     input.Timeout,
	}
	payload, err := jsonMerge(internalInput, input.HealthcheckHTTPDefinition)
	if err != nil {
		return result, err
	}
	_, err = c.sendRequest(ctx, fmt.Sprintf("/api/v1/healthcheck/http/%s", input.ID), http.MethodPut, payload, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}
