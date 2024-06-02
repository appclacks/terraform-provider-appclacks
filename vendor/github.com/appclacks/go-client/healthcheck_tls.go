package client

import (
	"context"
	"fmt"
	"net/http"
)

type HealthcheckTLSDefinition struct {
	Target          string `json:"target" validate:"required"`
	Port            uint   `json:"port" validate:"required,max=65535,min=1"`
	Key             string `json:"key,omitempty"`
	Cert            string `json:"cert,omitempty"`
	Cacert          string `json:"cacert,omitempty"`
	ServerName      string `json:"server-name,omitempty"`
	Insecure        bool   `json:"insecure"`
	ExpirationDelay string `json:"expiration-delay"`
}

type CreateTLSHealthcheckInput struct {
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Enabled     bool              `json:"bool" description:"Enable the healthcheck on the appclacks platform"`
	Timeout     string            `json:"timeout" validate:"required"`
	HealthcheckTLSDefinition
}

type UpdateTLSHealthcheckInput struct {
	ID          string            `json:"-" param:"id" description:"Healthcheck ID" validate:"required,uuid"`
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Timeout     string            `json:"timeout" validate:"required"`
	Enabled     bool              `json:"enabled" description:"Enable the healthcheck on the appclacks platform"`
	HealthcheckTLSDefinition
}

func (c *Client) CreateTLSHealthcheck(ctx context.Context, input CreateTLSHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	_, err := c.sendRequest(ctx, "/api/v1/healthcheck/tls", http.MethodPost, input, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}

func (c *Client) UpdateTLSHealthcheck(ctx context.Context, input UpdateTLSHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	internalInput := internalUpdateHealthcheckInput{
		Name:        input.Name,
		Description: input.Description,
		Labels:      input.Labels,
		Interval:    input.Interval,
		Enabled:     input.Enabled,
		Timeout:     input.Timeout,
	}
	payload, err := jsonMerge(internalInput, input.HealthcheckTLSDefinition)
	if err != nil {
		return result, err
	}
	_, err = c.sendRequest(ctx, fmt.Sprintf("/api/v1/healthcheck/tls/%s", input.ID), http.MethodPut, payload, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}
