package client

import (
	"context"
	"fmt"
	"net/http"
)

type HealthcheckCommandDefinition struct {
	Command   string   `json:"command" validate:"required,max=512,min=1"`
	Arguments []string `json:"arguments,omitempty"`
}

type CreateCommandHealthcheckInput struct {
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Enabled     bool              `json:"bool" description:"Enable the healthcheck on the appclacks platform"`
	Timeout     string            `json:"timeout" validate:"required"`
	HealthcheckCommandDefinition
}

type UpdateCommandHealthcheckInput struct {
	ID          string            `json:"-" param:"id" description:"Healthcheck ID" validate:"required,uuid"`
	Name        string            `json:"name" description:"Healthcheck name" validate:"required,max=255,min=1"`
	Description string            `json:"description" description:"Healthcheck description" validate:"max=255"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	Interval    string            `json:"interval" description:"Healthcheck interval" validate:"required"`
	Enabled     bool              `json:"enabled" description:"Enable the healthcheck on the appclacks platform"`
	Timeout     string            `json:"timeout" validate:"required"`
	HealthcheckCommandDefinition
}

func (c *Client) CreateCommandHealthcheck(ctx context.Context, input CreateCommandHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	_, err := c.sendRequest(ctx, "/api/v1/healthcheck/command", http.MethodPost, input, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}

func (c *Client) UpdateCommandHealthcheck(ctx context.Context, input UpdateCommandHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	internalInput := internalUpdateHealthcheckInput{
		Name:        input.Name,
		Description: input.Description,
		Labels:      input.Labels,
		Interval:    input.Interval,
		Enabled:     input.Enabled,
		Timeout:     input.Timeout,
	}
	payload, err := jsonMerge(internalInput, input.HealthcheckCommandDefinition)
	if err != nil {
		return result, err
	}
	_, err = c.sendRequest(ctx, fmt.Sprintf("/api/v1/healthcheck/command/%s", input.ID), http.MethodPut, payload, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}
