package client

import (
	"context"
	"fmt"
	"net/http"
)

type HealthcheckTCPDefinition struct {
	Target     string `json:"target" validate:"required"`
	Port       uint   `json:"port" validate:"required,max=65535,min=1"`
	ShouldFail bool   `json:"should-fail"`
}

type CreateTCPHealthcheckInput struct {
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Enabled     bool              `json:"bool" description:"Enable the healthcheck on the appclacks platform"`
	Timeout     string            `json:"timeout" validate:"required"`
	HealthcheckTCPDefinition
}

type UpdateTCPHealthcheckInput struct {
	ID          string            `json:"-" param:"id" description:"Healthcheck ID" validate:"required,uuid"`
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Timeout     string            `json:"timeout" validate:"required"`
	Enabled     bool              `json:"enabled" description:"Enable the healthcheck on the appclacks platform"`
	HealthcheckTCPDefinition
}

func (c *Client) CreateTCPHealthcheck(ctx context.Context, input CreateTCPHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	_, err := c.sendRequest(ctx, "/api/v1/healthcheck/tcp", http.MethodPost, input, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}

func (c *Client) UpdateTCPHealthcheck(ctx context.Context, input UpdateTCPHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	internalInput := internalUpdateHealthcheckInput{
		Name:        input.Name,
		Description: input.Description,
		Labels:      input.Labels,
		Interval:    input.Interval,
		Enabled:     input.Enabled,
		Timeout:     input.Timeout,
	}
	payload, err := jsonMerge(internalInput, input.HealthcheckTCPDefinition)
	if err != nil {
		return result, err
	}
	_, err = c.sendRequest(ctx, fmt.Sprintf("/api/v1/healthcheck/tcp/%s", input.ID), http.MethodPut, payload, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}
