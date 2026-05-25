package domain

// ApiEndpoint represents a single normalized API endpoint extracted from uploaded documentation
type ApiEndpoint struct {
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	Summary     string            `json:"summary"`
	Headers     map[string]string `json:"headers,omitempty"`
	QueryParams []ParamInfo       `json:"query_params,omitempty"`
	RequestBody string            `json:"request_body,omitempty"`
	Responses   map[string]string `json:"responses,omitempty"`
}

// ParamInfo represents details of a query or path parameter
type ParamInfo struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

// ApiContext represents the extracted API documentation context
type ApiContext struct {
	Source    string        `json:"source"`   // e.g. "openapi", "curl", "url"
	BaseURL   string        `json:"base_url"`
	Endpoints []ApiEndpoint `json:"endpoints"`
	RawSpec   string        `json:"raw_spec,omitempty"`
}
