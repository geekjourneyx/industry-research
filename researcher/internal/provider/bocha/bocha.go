package bocha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

const (
	DefaultEndpoint = "https://api.bocha.cn/v1/web-search"

	providerName = "bocha"
	contentType  = "web_page"
	leadOnly     = "lead_only"
)

type Client struct {
	apiKey     string
	endpoint   string
	httpClient *http.Client
}

func NewClient(apiKey, endpoint string, httpClient *http.Client) *Client {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{
		apiKey:     apiKey,
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

func (c *Client) Search(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	query := strings.TrimSpace(req.Query)
	resp := newResponse(query, req.Parameters)

	if query == "" {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeInvalidArgument,
			Message:     "query is required",
			Retryable:   false,
			AgentAction: "Provide a non-empty query.",
		})
		return resp, errors.New("bocha search: query is required")
	}
	if strings.TrimSpace(c.apiKey) == "" {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeMissingAPIKey,
			Message:     "Bocha API key is required",
			Retryable:   false,
			AgentAction: "Set BOCHA_API_KEY or providers.bocha.api_key in researcher config.",
		})
		return resp, errors.New("bocha search: missing API key")
	}

	body := requestBody(query, req.Parameters)
	payload, err := json.Marshal(body)
	if err != nil {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeInvalidArgument,
			Message:     err.Error(),
			Retryable:   false,
			AgentAction: "Use JSON-serializable retrieval parameters.",
		})
		return resp, fmt.Errorf("bocha search: encode request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeInvalidArgument,
			Message:     err.Error(),
			Retryable:   false,
			AgentAction: "Check the Bocha endpoint URL.",
		})
		return resp, fmt.Errorf("bocha search: create request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		code := rerrors.CodeProviderHTTPError
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			code = rerrors.CodeProviderTimeout
		}
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        code,
			Message:     err.Error(),
			Retryable:   code == rerrors.CodeProviderTimeout,
			AgentAction: "Retry the request or check network connectivity.",
		})
		return resp, fmt.Errorf("bocha search: request failed: %w", err)
	}
	defer httpResp.Body.Close()

	raw, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:           rerrors.CodeProviderParseError,
			Message:        err.Error(),
			ProviderStatus: httpResp.StatusCode,
			Retryable:      false,
			AgentAction:    "Retry and inspect the provider response body.",
		})
		return resp, fmt.Errorf("bocha search: read response: %w", err)
	}

	var decoded bochaResponse
	if len(strings.TrimSpace(string(raw))) > 0 {
		if err := json.Unmarshal(raw, &decoded); err != nil {
			resp.Errors = append(resp.Errors, retrieval.Error{
				Code:           rerrors.CodeProviderParseError,
				Message:        err.Error(),
				ProviderStatus: httpResp.StatusCode,
				Retryable:      false,
				AgentAction:    "Inspect the provider response JSON.",
				RawErrorPath:   "$",
			})
			return resp, fmt.Errorf("bocha search: decode response: %w", err)
		}
	}

	providerCode := providerCodeString(decoded.Code)
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 || providerCode != "200" {
		retrievalErr := mapProviderError(httpResp.StatusCode, providerCode, decoded.message(), decoded.LogID)
		resp.Errors = append(resp.Errors, retrievalErr)
		return resp, fmt.Errorf("bocha search: provider error: %s", retrievalErr.Message)
	}

	resp.Usage["total_estimated_matches"] = decoded.Data.WebPages.TotalEstimatedMatches
	resp.Usage["some_results_removed"] = decoded.Data.WebPages.SomeResultsRemoved
	resp.Usage["provider_log_id"] = decoded.LogID
	resp.Items = make([]retrieval.Item, 0, len(decoded.Data.WebPages.Value))
	for i, page := range decoded.Data.WebPages.Value {
		resp.Items = append(resp.Items, retrieval.Item{
			Rank:                 i + 1,
			Title:                page.Name,
			URL:                  page.URL,
			DisplayURL:           page.DisplayURL,
			SiteName:             page.SiteName,
			SiteIcon:             page.SiteIcon,
			Snippet:              page.Snippet,
			Summary:              page.Summary,
			PublishedAt:          page.DatePublished,
			LastCrawledAt:        normalizeBochaTime(page.DateLastCrawled),
			Language:             page.Language,
			ContentType:          contentType,
			SourceConfidenceHint: leadOnly,
			ProviderMetadata: map[string]any{
				"bocha_id":                page.ID,
				"cached_page_url":         page.CachedPageURL,
				"is_navigational":         page.IsNavigational,
				"is_family_friendly":      page.IsFamilyFriendly,
				"raw_date_last_crawled":   page.DateLastCrawled,
				"total_estimated_matches": decoded.Data.WebPages.TotalEstimatedMatches,
				"some_results_removed":    decoded.Data.WebPages.SomeResultsRemoved,
				"provider_log_id":         decoded.LogID,
			},
		})
	}

	return resp, nil
}

func newResponse(query string, parameters map[string]any) retrieval.RetrievalResponse {
	request := map[string]any{}
	for key, value := range parameters {
		request[key] = value
	}
	return retrieval.RetrievalResponse{
		Provider:       providerName,
		ProviderType:   retrieval.ProviderTypeDirectSearch,
		Mode:           retrieval.ModeSearch,
		Query:          query,
		RetrievedAt:    time.Now(),
		Request:        request,
		RetrievalCalls: []retrieval.RetrievalCall{},
		Items:          []retrieval.Item{},
		Answer: retrieval.Answer{
			Citations: []retrieval.Citation{},
		},
		Usage:  map[string]any{},
		Errors: []retrieval.Error{},
	}
}

func requestBody(query string, parameters map[string]any) map[string]any {
	body := map[string]any{"query": query}
	for key, value := range parameters {
		body[key] = value
	}
	if _, ok := body["summary"]; !ok {
		body["summary"] = true
	}
	return body
}

func normalizeBochaTime(value string) string {
	if strings.HasSuffix(value, "Z") {
		return strings.TrimSuffix(value, "Z") + "+08:00"
	}
	return value
}

func mapProviderError(status int, providerCode, message, logID string) retrieval.Error {
	if message == "" {
		message = http.StatusText(status)
	}
	code := rerrors.CodeProviderHTTPError
	retryable := false
	action := "Inspect the provider error and adjust the request."

	switch {
	case status == http.StatusTooManyRequests || providerCode == "429":
		code = rerrors.CodeProviderRateLimited
		retryable = true
		action = "Retry after waiting, or reduce request frequency."
	case status == http.StatusUnauthorized || status == http.StatusForbidden || providerCode == "401" || providerCode == "403":
		code = rerrors.CodeProviderAuthError
		action = "Check the Bocha API key and account permissions."
	case status >= 500:
		code = rerrors.CodeProviderUnavailable
		retryable = true
		action = "Retry later; the provider may be unavailable."
	}

	return retrieval.Error{
		Code:           code,
		Message:        message,
		ProviderStatus: status,
		ProviderCode:   providerCode,
		ProviderLogID:  logID,
		Retryable:      retryable,
		AgentAction:    action,
		RawErrorPath:   "$",
	}
}

func providerCodeString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case float64:
		if typed == float64(int64(typed)) {
			return strconv.FormatInt(int64(typed), 10)
		}
		return strconv.FormatFloat(typed, 'f', -1, 64)
	case int:
		return strconv.Itoa(typed)
	case int64:
		return strconv.FormatInt(typed, 10)
	default:
		return fmt.Sprint(typed)
	}
}

type bochaResponse struct {
	Code    any       `json:"code"`
	LogID   string    `json:"log_id"`
	Msg     any       `json:"msg"`
	Message string    `json:"message"`
	Data    bochaData `json:"data"`
}

func (r bochaResponse) message() string {
	if r.Message != "" {
		return r.Message
	}
	switch typed := r.Msg.(type) {
	case string:
		return typed
	case nil:
		return ""
	default:
		return fmt.Sprint(typed)
	}
}

type bochaData struct {
	QueryContext bochaQueryContext `json:"queryContext"`
	WebPages     bochaWebPages     `json:"webPages"`
}

type bochaQueryContext struct {
	OriginalQuery string `json:"originalQuery"`
}

type bochaWebPages struct {
	TotalEstimatedMatches int         `json:"totalEstimatedMatches"`
	SomeResultsRemoved    bool        `json:"someResultsRemoved"`
	Value                 []bochaPage `json:"value"`
}

type bochaPage struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	URL              string `json:"url"`
	DisplayURL       string `json:"displayUrl"`
	Snippet          string `json:"snippet"`
	Summary          string `json:"summary"`
	SiteName         string `json:"siteName"`
	SiteIcon         string `json:"siteIcon"`
	DatePublished    string `json:"datePublished"`
	DateLastCrawled  string `json:"dateLastCrawled"`
	CachedPageURL    string `json:"cachedPageUrl"`
	Language         string `json:"language"`
	IsFamilyFriendly bool   `json:"isFamilyFriendly"`
	IsNavigational   bool   `json:"isNavigational"`
}
