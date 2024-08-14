package argocd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

type APIClient struct {
	httpClient *http.Client
	baseURL    *url.URL
	options    OptionsAuth
}

func NewAPIClient(options OptionsAuth, httpClient *http.Client) (*APIClient, error) {
	baseURL, err := url.Parse(options.ArgoCDServer)
	if err != nil {
		return nil, fmt.Errorf("error parsing ArgoCD server URL: %w", err)
	}
	if baseURL.Scheme == "" {
		baseURL.Scheme = "https"
	}
	return &APIClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		options:    options,
	}, nil
}

func (ac *APIClient) GetResourceManifest(ctx context.Context, appName string, appNamespace string, namespace string, kind string, version string, resourceName string, group string) (json.RawMessage, error) {
	reqURL := ac.baseURL.JoinPath(fmt.Sprintf("/api/v1/applications/%s/resource", appName))
	paramValues := reqURL.Query()
	paramValues.Set("name", resourceName)
	paramValues.Set("appNamespace", appNamespace)
	paramValues.Set("namespace", namespace)
	paramValues.Set("resourceName", resourceName)
	paramValues.Set("version", version)
	paramValues.Set("kind", kind)
	paramValues.Set("group", group)
	reqURL.RawQuery = paramValues.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := ac.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Manifest json.RawMessage `json:"manifest"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	unquoted, err := strconv.Unquote(string(response.Manifest))
	if err != nil {
		return nil, fmt.Errorf("error unquoteing manifest: %w", err)
	}

	return []byte(unquoted), nil
}

func (ac *APIClient) GetResources(ctx context.Context, appName string, appNamespace string) (json.RawMessage, error) {
	reqURL := ac.baseURL.JoinPath(fmt.Sprintf("/api/v1/applications/%s/resource-tree", appName))
	paramValues := reqURL.Query()
	paramValues.Set("appNamespace", appNamespace)
	reqURL.RawQuery = paramValues.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", reqURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := ac.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response struct {
		Nodes json.RawMessage `json:"nodes"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return response.Nodes, nil
}

func (ac *APIClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "argocd.token="+ac.options.ArgoCDAuthToken)
	return ac.httpClient.Do(req)
}
