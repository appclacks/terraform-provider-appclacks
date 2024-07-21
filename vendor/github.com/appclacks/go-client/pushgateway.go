package client

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type PushgatewayMetric struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	TTL         string            `json:"ttl"`
	Type        string            `json:"type"`
	CreatedAt   time.Time         `json:"created_at"`
	ExpiresAt   *time.Time        `json:"expires_at,omitempty"`
	Value       string            `json:"value"`
}

type CreateOrUpdatePushgatewayMetricInput struct {
	Name        string            `json:"name" validate:"required,max=255,min=1"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels" description:"Healthcheck labels" validate:"dive,keys,max=255,min=1,endkeys,max=255,min=1"`
	TTL         string            `json:"ttl"`
	Type        string            `json:"type" validate:"omitempty,oneof=counter gauge histogram summary"`
	Value       string            `json:"value" validate:"required"`
}

type DeletePushgatewayMetricInput struct {
	Identifier string `param:"identifier" validate:"required"`
}

type ListPushgatewayMetricsOutput struct {
	Result []PushgatewayMetric `json:"result"`
}

func (c *Client) CreateOrUpdatePushgatewayMetric(ctx context.Context, input CreateOrUpdatePushgatewayMetricInput) (Response, error) {
	var result Response
	_, err := c.sendRequest(ctx, "/api/v1/pushgateway", http.MethodPost, input, &result, nil)
	if err != nil {
		return Response{}, err
	}
	return result, nil
}

func (c *Client) DeletePushgatewayMetric(ctx context.Context, input DeletePushgatewayMetricInput) (Response, error) {
	var result Response
	_, err := c.sendRequest(ctx, fmt.Sprintf("/api/v1/pushgateway/%s", input.Identifier), http.MethodDelete, nil, &result, nil)
	if err != nil {
		return Response{}, err
	}
	return result, nil
}

func (c *Client) ListPushgatewayMetrics(ctx context.Context) (ListPushgatewayMetricsOutput, error) {
	var result ListPushgatewayMetricsOutput
	_, err := c.sendRequest(ctx, "/api/v1/pushgateway", http.MethodGet, nil, &result, nil)
	if err != nil {
		return ListPushgatewayMetricsOutput{}, err
	}
	return result, nil
}

func (c *Client) DeleteAllPushgatewayMetrics(ctx context.Context) (Response, error) {
	var result Response
	_, err := c.sendRequest(ctx, "/api/v1/pushgateway", http.MethodDelete, nil, &result, nil)
	if err != nil {
		return Response{}, err
	}
	return result, nil
}
