package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Client struct {
	http     *http.Client
	username string
	password string
	endpoint string
	key      string
	cert     string
	cacert   string
	insecure bool
}

var (
	ErrNotFound = errors.New("Not found")
)

func loadEnv(client *Client) {
	if os.Getenv("APPCLACKS_USERNAME") != "" {
		client.username = os.Getenv("APPCLACKS_USERNAME")
	}

	if os.Getenv("APPCLACKS_PASSWORD") != "" {
		client.password = os.Getenv("APPCLACKS_PASSWORD")
	}

	if os.Getenv("APPCLACKS_API_ENDPOINT") != "" {
		client.endpoint = os.Getenv("APPCLACKS_API_ENDPOINT")
	}

	if os.Getenv("APPCLACKS_TLS_KEY") != "" {
		client.key = os.Getenv("APPCLACKS_TLS_KEY")
	}
	if os.Getenv("APPCLACKS_TLS_CERT") != "" {
		client.cert = os.Getenv("APPCLACKS_TLS_CERT")
	}
	if os.Getenv("APPCLACKS_TLS_CACERT") != "" {
		client.cacert = os.Getenv("APPCLACKS_TLS_CACERT")
	}
	insecure := os.Getenv("APPCLACKS_TLS_INSECURE")
	if insecure == "true" {
		client.insecure = true
	}
}

type ClientOption func(c *Client) error

func New(options ...ClientOption) (*Client, error) {
	client := &Client{
		http: &http.Client{},
	}

	loadEnv(client)
	for _, option := range options {
		err := option(client)
		if err != nil {
			return nil, err
		}
	}
	if client.cert != "" || client.key != "" || client.cacert != "" || client.insecure {
		tlsConfig, err := getTLSConfig(client.key, client.cert, client.cacert, "", client.insecure)
		if err != nil {
			return nil, err
		}
		transport := &http.Transport{
			TLSClientConfig: tlsConfig,
		}
		client.http.Transport = transport
	}
	return client, nil
}

func WithUsername(username string) ClientOption {
	return func(c *Client) error {
		c.username = username
		return nil
	}
}

func WithPassword(password string) ClientOption {
	return func(c *Client) error {
		c.password = password
		return nil
	}
}

func WithEndpoint(endpoint string) ClientOption {
	return func(c *Client) error {
		c.endpoint = endpoint
		return nil
	}
}

func WithKey(key string) ClientOption {
	return func(c *Client) error {
		c.key = key
		return nil
	}
}

func WithCert(cert string) ClientOption {
	return func(c *Client) error {
		c.cert = cert
		return nil
	}
}

func WithCacert(cacert string) ClientOption {
	return func(c *Client) error {
		c.cacert = cacert
		return nil
	}
}

func WithInsecure(insecure bool) ClientOption {
	return func(c *Client) error {
		c.insecure = insecure
		return nil
	}
}

func (c *Client) sendRequest(ctx context.Context, url string, method string, body any, result any, queryParams map[string]string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		json, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(json)
	}
	request, err := http.NewRequestWithContext(
		ctx,
		method,
		fmt.Sprintf("%s%s", c.endpoint, url),
		reqBody)
	if err != nil {
		return nil, err
	}
	if len(queryParams) != 0 {
		q := request.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		request.URL.RawQuery = q.Encode()
	}
	request.Header.Add("content-type", "application/json")
	if c.username != "" {
		authString := fmt.Sprintf("%s:%s", c.username, c.password)
		creds := base64.StdEncoding.EncodeToString([]byte(authString))
		request.Header.Add("Authorization", fmt.Sprintf("Basic %s", creds))
	}
	response, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode >= 400 {
		if response.StatusCode == 404 {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("the API returned an error: status %d\n%s", response.StatusCode, string(b))
	}
	if result != nil {
		err = json.Unmarshal(b, result)
		if err != nil {
			return nil, err
		}
	}
	return response, nil
}

func jsonMerge(s1 any, s2 any) (map[string]any, error) {
	result := make(map[string]any)
	str1, err := json.Marshal(s1)
	if err != nil {
		return nil, err
	}
	str2, err := json.Marshal(s2)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(str1, &result)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(str2, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
