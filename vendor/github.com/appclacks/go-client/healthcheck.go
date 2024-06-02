package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Healthcheck struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type"`
	Labels      map[string]string `json:"labels,omitempty"`
	Timeout     string            `json:"timeout" validate:"required"`
	Interval    string            `json:"interval" validate:"required"`
	CreatedAt   time.Time         `json:"created-at"`
	Enabled     bool              `json:"enabled"`
	Definition  any               `json:",inline"`
}

type internalHealthcheck struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Type        string            `json:"type"`
	Labels      map[string]string `json:"labels,omitempty"`
	Enabled     bool              `json:"enabled"`
	Interval    string            `json:"interval" validate:"required"`
	Timeout     string            `json:"timeout" validate:"required"`
	CreatedAt   time.Time         `json:"created-at"`
}

func (h *Healthcheck) MarshalJSON() ([]byte, error) {
	result := make(map[string]any)
	internal := internalHealthcheck{
		ID:          h.ID,
		Name:        h.Name,
		Enabled:     h.Enabled,
		Description: h.Description,
		Type:        h.Type,
		Timeout:     h.Timeout,
		Labels:      h.Labels,
		Interval:    h.Interval,
		CreatedAt:   h.CreatedAt,
	}
	internalStr, err := json.Marshal(internal)
	if err != nil {
		return nil, err
	}
	defStr, err := json.Marshal(h.Definition)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(internalStr, &result)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(defStr, &result)
	if err != nil {
		return nil, err
	}
	return json.Marshal(result)
}

func (h *Healthcheck) UnmarshalJSON(data []byte) error {

	var internal internalHealthcheck
	if err := json.Unmarshal(data, &internal); err != nil {
		return err
	}

	if internal.Type == "dns" {
		var definition HealthcheckDNSDefinition
		if err := json.Unmarshal(data, &definition); err != nil {
			return err
		}
		h.Definition = definition
	} else if internal.Type == "tcp" {
		var definition HealthcheckTCPDefinition
		if err := json.Unmarshal(data, &definition); err != nil {
			return err
		}
		h.Definition = definition
	} else if internal.Type == "http" {
		var definition HealthcheckHTTPDefinition
		if err := json.Unmarshal(data, &definition); err != nil {
			return err
		}
		h.Definition = definition
	} else if internal.Type == "tls" {
		var definition HealthcheckTLSDefinition
		if err := json.Unmarshal(data, &definition); err != nil {
			return err
		}
		h.Definition = definition
	} else if internal.Type == "command" {
		var definition HealthcheckCommandDefinition
		if err := json.Unmarshal(data, &definition); err != nil {
			return err
		}
		h.Definition = definition
	} else {
		return fmt.Errorf("Unknown healthcheck type %s", h.Type)
	}
	h.ID = internal.ID
	h.Name = internal.Name
	h.Timeout = internal.Timeout
	h.Description = internal.Description
	h.Type = internal.Type
	h.Enabled = internal.Enabled
	h.Labels = internal.Labels
	h.Interval = internal.Interval
	h.CreatedAt = internal.CreatedAt

	return nil
}

type DeleteHealthcheckInput struct {
	ID string `param:"id" description:"Healthcheck ID" validate:"required,uuid"`
}

type GetHealthcheckInput struct {
	Identifier string `param:"identifier" description:"Healthcheck Name or ID"`
}

type ListHealthchecksInput struct {
	NamePattern string `query:"name-pattern" description:"Returns all health checks whose names are matching this regular expression"`
}

type ListHealthchecksOutput struct {
	Result []Healthcheck `json:"result"`
}

type CabourotteDiscoveryInput struct {
	// foo=bar,a=b
	Labels string `query:"labels"`
}

type CabourotteDiscoveryOutput struct {
	DNSChecks     []Healthcheck `json:"dns-checks,omitempty"`
	TCPChecks     []Healthcheck `json:"tcp-checks,omitempty"`
	HTTPChecks    []Healthcheck `json:"http-checks,omitempty"`
	TLSChecks     []Healthcheck `json:"tls-checks,omitempty"`
	CommandChecks []Healthcheck `json:"command-checks,omitempty"`
}

type internalUpdateHealthcheckInput struct {
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Labels      map[string]string `json:"labels,omitempty"`
	Interval    string            `json:"interval"`
	Enabled     bool              `json:"enabled"`
	Timeout     string            `json:"timeout"`
}

func (c *Client) DeleteHealthcheck(ctx context.Context, input DeleteHealthcheckInput) (Response, error) {
	var result Response
	_, err := c.sendRequest(ctx, fmt.Sprintf("/api/v1/healthcheck/%s", input.ID), http.MethodDelete, nil, &result, nil)
	if err != nil {
		return Response{}, err
	}
	return result, nil
}

func (c *Client) GetHealthcheck(ctx context.Context, input GetHealthcheckInput) (Healthcheck, error) {
	var result Healthcheck
	_, err := c.sendRequest(ctx, fmt.Sprintf("/api/v1/healthcheck/%s", input.Identifier), http.MethodGet, nil, &result, nil)
	if err != nil {
		return Healthcheck{}, err
	}
	return result, nil
}

func (c *Client) ListHealthchecks(ctx context.Context) (ListHealthchecksOutput, error) {
	var result ListHealthchecksOutput
	_, err := c.sendRequest(ctx, "/api/v1/healthcheck", http.MethodGet, nil, &result, nil)
	if err != nil {
		return ListHealthchecksOutput{}, err
	}
	return result, nil
}

func (c *Client) CabourotteDiscovery(ctx context.Context, input CabourotteDiscoveryInput) (CabourotteDiscoveryOutput, error) {
	var result CabourotteDiscoveryOutput
	queryParams := make(map[string]string)
	if input.Labels != "" {
		queryParams["labels"] = input.Labels
	}
	_, err := c.sendRequest(ctx, "/cabourotte/discovery", http.MethodGet, nil, &result, queryParams)
	if err != nil {
		return CabourotteDiscoveryOutput{}, err
	}
	return result, nil
}
