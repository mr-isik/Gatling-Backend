package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mr-isik/gatling-backend/internal/domain"
)

type ApiDocParser struct{}

func NewApiDocParser() *ApiDocParser {
	return &ApiDocParser{}
}

// ParseOpenAPISpec parses a raw OpenAPI/Swagger JSON or YAML string into our domain.ApiContext
func (p *ApiDocParser) ParseOpenAPISpec(ctx context.Context, rawSpec []byte) (*domain.ApiContext, error) {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	doc, err := loader.LoadFromData(rawSpec)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OpenAPI spec: %w", err)
	}

	baseURL := ""
	if len(doc.Servers) > 0 {
		baseURL = doc.Servers[0].URL
	}

	endpoints := make([]domain.ApiEndpoint, 0)

	for path, pathItem := range doc.Paths.Map() {
		for method, op := range pathItem.Operations() {
			endpoint := domain.ApiEndpoint{
				Method:  strings.ToUpper(method),
				Path:    path,
				Summary: op.Summary,
				Headers: make(map[string]string),
			}

			// Parameters
			for _, paramRef := range op.Parameters {
				if paramRef.Value == nil {
					continue
				}
				param := paramRef.Value
				
				typ := "string" // fallback

				if param.In == "header" {
					endpoint.Headers[param.Name] = typ
				} else if param.In == "query" || param.In == "path" {
					endpoint.QueryParams = append(endpoint.QueryParams, domain.ParamInfo{
						Name:     param.Name,
						Type:     typ,
						Required: param.Required,
					})
				}
			}

			// Request Body
			if op.RequestBody != nil && op.RequestBody.Value != nil {
				content := op.RequestBody.Value.Content
				if jsonContent, ok := content["application/json"]; ok && jsonContent.Schema != nil {
					schemaBytes, _ := json.Marshal(jsonContent.Schema.Value)
					endpoint.RequestBody = string(schemaBytes)
					endpoint.Headers["Content-Type"] = "application/json"
				}
			}

			// Responses
			endpoint.Responses = make(map[string]string)
			for status, respRef := range op.Responses.Map() {
				if respRef.Value != nil && respRef.Value.Description != nil {
					endpoint.Responses[status] = *respRef.Value.Description
				}
			}

			endpoints = append(endpoints, endpoint)
		}
	}

	return &domain.ApiContext{
		Source:    "openapi",
		BaseURL:   baseURL,
		Endpoints: endpoints,
		RawSpec:   string(rawSpec),
	}, nil
}

// ParseCurlCommands extracts endpoints from multiple curl commands
func (p *ApiDocParser) ParseCurlCommands(curlStr string) (*domain.ApiContext, error) {
	// A rudimentary parser that extracts URL, Method, Headers and Body from raw curl texts.
	commands := strings.Split(curlStr, "curl ")
	endpoints := make([]domain.ApiEndpoint, 0)
	baseURL := ""

	urlRegex := regexp.MustCompile(`(?:'|")?(https?://[^'"\s]+)(?:'|")?`)
	methodRegex := regexp.MustCompile(`-X\s+([A-Z]+)`)
	headerRegex := regexp.MustCompile(`-H\s+(?:'|")([^'"]+)(?:'|")`)
	bodyRegex := regexp.MustCompile(`(?:--data|-d)\s+(?:'([^']+)'|"([^"]+)")`)

	for _, cmd := range commands {
		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}
		
		endpoint := domain.ApiEndpoint{
			Method:  "GET", // Default
			Headers: make(map[string]string),
		}

		// URL
		urlMatch := urlRegex.FindStringSubmatch("curl " + cmd)
		if len(urlMatch) > 1 {
			fullURL := urlMatch[1]
			if baseURL == "" {
				// Naive base url extraction
				parts := strings.SplitN(fullURL, "/", 4)
				if len(parts) >= 3 {
					baseURL = parts[0] + "//" + parts[2]
				}
			}
			endpoint.Path = strings.TrimPrefix(fullURL, baseURL)
			if endpoint.Path == "" {
				endpoint.Path = "/"
			}
		} else {
			// No URL found, skip
			continue
		}

		// Method
		methodMatch := methodRegex.FindStringSubmatch(cmd)
		if len(methodMatch) > 1 {
			endpoint.Method = methodMatch[1]
		} else if strings.Contains(cmd, "-d ") || strings.Contains(cmd, "--data ") {
			endpoint.Method = "POST"
		}

		// Headers
		headers := headerRegex.FindAllStringSubmatch(cmd, -1)
		for _, hMatch := range headers {
			if len(hMatch) > 1 {
				parts := strings.SplitN(hMatch[1], ":", 2)
				if len(parts) == 2 {
					endpoint.Headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
				}
			}
		}

		// Body
		bodyMatch := bodyRegex.FindStringSubmatch(cmd)
		if len(bodyMatch) > 1 {
			if bodyMatch[1] != "" {
				endpoint.RequestBody = bodyMatch[1]
			} else if len(bodyMatch) > 2 && bodyMatch[2] != "" {
				endpoint.RequestBody = bodyMatch[2]
			}
		}

		endpoint.Summary = "Extracted from cURL"
		endpoints = append(endpoints, endpoint)
	}

	return &domain.ApiContext{
		Source:    "curl",
		BaseURL:   baseURL,
		Endpoints: endpoints,
		RawSpec:   curlStr,
	}, nil
}

// FetchAndParseURL fetches an OpenAPI JSON spec from a URL and parses it
func (p *ApiDocParser) FetchAndParseURL(ctx context.Context, url string) (*domain.ApiContext, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("URL returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	apiCtx, err := p.ParseOpenAPISpec(ctx, body)
	if err != nil {
		return nil, err
	}
	apiCtx.Source = "url"
	return apiCtx, nil
}
