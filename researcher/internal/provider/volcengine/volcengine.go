package volcengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

const (
	DefaultEndpoint = "https://ark.cn-beijing.volces.com/api/v3/responses"
	DefaultModel    = "doubao-seed-2-0-lite-260215"

	providerName  = "volcengine"
	contentType   = "annotation_url"
	leadOnly      = "lead_only"
	webSearchTool = "web_search"
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
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{
		apiKey:     apiKey,
		endpoint:   endpoint,
		httpClient: httpClient,
	}
}

func (c *Client) Answer(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	query := strings.TrimSpace(req.Query)
	resp := newResponse(query, req.Parameters)

	if query == "" {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeInvalidArgument,
			Message:     "query is required",
			Retryable:   false,
			AgentAction: "Provide a non-empty query.",
		})
		return resp, errors.New("volcengine answer: query is required")
	}
	if strings.TrimSpace(c.apiKey) == "" {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeMissingAPIKey,
			Message:     "Volcengine ARK API key is required",
			Retryable:   false,
			AgentAction: "Set ARK_API_KEY or providers.volcengine.api_key in researcher config.",
		})
		return resp, errors.New("volcengine answer: missing API key")
	}

	payload, err := json.Marshal(buildRequest(query, req.Parameters))
	if err != nil {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeInvalidArgument,
			Message:     err.Error(),
			Retryable:   false,
			AgentAction: "Use JSON-serializable Volcengine request parameters.",
		})
		return resp, fmt.Errorf("volcengine answer: encode request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:        rerrors.CodeInvalidArgument,
			Message:     err.Error(),
			Retryable:   false,
			AgentAction: "Check the Volcengine endpoint URL.",
		})
		return resp, fmt.Errorf("volcengine answer: create request: %w", err)
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
		return resp, fmt.Errorf("volcengine answer: request failed: %w", err)
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
			RawErrorPath:   "$",
		})
		return resp, fmt.Errorf("volcengine answer: read response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		retrievalErr := mapProviderError(httpResp.StatusCode, raw)
		resp.Errors = append(resp.Errors, retrievalErr)
		return resp, fmt.Errorf("volcengine answer: provider error: %s", retrievalErr.Message)
	}
	if len(strings.TrimSpace(string(raw))) == 0 {
		resp.Errors = append(resp.Errors, parseError(httpResp.StatusCode, "empty provider response body"))
		return resp, errors.New("volcengine answer: empty response body")
	}

	var decoded responsesAPIResponse
	if err := json.Unmarshal(raw, &decoded); err != nil {
		resp.Errors = append(resp.Errors, parseError(httpResp.StatusCode, err.Error()))
		return resp, fmt.Errorf("volcengine answer: decode response: %w", err)
	}

	resp.Usage = decoded.Usage
	if resp.Usage == nil {
		resp.Usage = map[string]any{}
	}
	parseOutput(&resp, decoded.Output)
	if len(resp.RetrievalCalls) == 0 {
		resp.Errors = append(resp.Errors, retrieval.Error{
			Code:         rerrors.CodeNoRetrievalTriggered,
			Message:      "model response did not include a web_search_call",
			Retryable:    false,
			AgentAction:  "Rephrase the query to require current web search, or lower search constraints such as max_keyword, limit, or max_tool_calls.",
			RawErrorPath: "$.output",
		})
		return resp, errors.New("volcengine answer: no retrieval triggered")
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
		ProviderType:   retrieval.ProviderTypeModelAnswerSearch,
		Mode:           retrieval.ModeAnswer,
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

func buildRequest(query string, parameters map[string]any) map[string]any {
	model := DefaultModel
	if value, ok := parameters["model"].(string); ok && strings.TrimSpace(value) != "" {
		model = strings.TrimSpace(value)
	}

	tool := map[string]any{"type": webSearchTool}
	for _, key := range []string{"max_keyword", "limit", "sources", "user_location"} {
		if value, ok := parameters[key]; ok {
			tool[key] = value
		}
	}

	body := map[string]any{
		"model": model,
		"tools": []map[string]any{tool},
		"input": []map[string]any{
			{
				"role": "user",
				"content": []map[string]string{
					{
						"type": "input_text",
						"text": query,
					},
				},
			},
		},
	}
	if value, ok := parameters["max_tool_calls"]; ok {
		body["max_tool_calls"] = value
	}
	return body
}

func parseOutput(resp *retrieval.RetrievalResponse, output []responseOutputItem) {
	for _, item := range output {
		switch item.Type {
		case "web_search_call":
			resp.RetrievalCalls = append(resp.RetrievalCalls, retrieval.RetrievalCall{
				CallID:         item.ID,
				Query:          item.Action.Query,
				Status:         item.Status,
				ProviderAction: webSearchTool,
				ProviderMetadata: map[string]any{
					"action_type": item.Action.Type,
				},
			})
		case "message":
			for _, content := range item.Content {
				if content.Type != "output_text" {
					continue
				}
				resp.Answer.Text += content.Text
				for _, annotation := range content.Annotations {
					if annotation.Type != "url_citation" || annotation.URL == "" {
						continue
					}
					index := len(resp.Answer.Citations) + 1
					host := hostFromURL(annotation.URL)
					resp.Answer.Citations = append(resp.Answer.Citations, retrieval.Citation{
						Index:  index,
						URL:    annotation.URL,
						Title:  annotation.Title,
						Source: host,
					})
					resp.Items = append(resp.Items, retrieval.Item{
						Rank:                 index,
						Title:                annotation.Title,
						URL:                  annotation.URL,
						SiteName:             host,
						ContentType:          contentType,
						SourceConfidenceHint: leadOnly,
						ProviderMetadata: map[string]any{
							"annotation_index": index,
							"start_index":      annotation.StartIndex,
							"end_index":        annotation.EndIndex,
						},
					})
				}
			}
		}
	}
}

func mapProviderError(status int, raw []byte) retrieval.Error {
	providerCode, message := providerErrorFields(raw)
	if message == "" {
		message = http.StatusText(status)
	}

	code := rerrors.CodeProviderHTTPError
	retryable := false
	action := "Inspect the Volcengine provider error and adjust the request."

	switch {
	case status == http.StatusBadRequest:
		code = rerrors.CodeInvalidArgument
		action = "Check Volcengine request parameters; remove unsupported caching settings and verify max_keyword, limit, max_tool_calls, sources, and user_location."
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		code = rerrors.CodeProviderAuthError
		action = "Check ARK_API_KEY and confirm the account has web_search permission enabled in Volcengine Ark."
	case status == http.StatusTooManyRequests:
		code = rerrors.CodeProviderRateLimited
		retryable = true
		action = "Retry after waiting, or lower max_keyword and limit to reduce web_search usage."
	case status >= 500:
		code = rerrors.CodeProviderUnavailable
		retryable = true
		action = "Retry later; Volcengine may be unavailable."
	}

	return retrieval.Error{
		Code:           code,
		Message:        message,
		ProviderStatus: status,
		ProviderCode:   providerCode,
		Retryable:      retryable,
		AgentAction:    action,
		RawErrorPath:   "$",
	}
}

func providerErrorFields(raw []byte) (string, string) {
	var decoded struct {
		Code    any    `json:"code"`
		Message string `json:"message"`
		Error   struct {
			Code    any    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		return "", strings.TrimSpace(string(raw))
	}
	code := codeString(decoded.Code)
	message := decoded.Message
	if code == "" {
		code = codeString(decoded.Error.Code)
	}
	if message == "" {
		message = decoded.Error.Message
	}
	return code, message
}

func parseError(status int, message string) retrieval.Error {
	return retrieval.Error{
		Code:           rerrors.CodeProviderParseError,
		Message:        message,
		ProviderStatus: status,
		Retryable:      false,
		AgentAction:    "Inspect the Volcengine response JSON.",
		RawErrorPath:   "$",
	}
}

func hostFromURL(value string) string {
	parsed, err := url.Parse(value)
	if err != nil || parsed.Host == "" {
		return ""
	}
	return parsed.Hostname()
}

func codeString(value any) string {
	switch typed := value.(type) {
	case nil:
		return ""
	case string:
		return typed
	case float64:
		return fmt.Sprintf("%.0f", typed)
	default:
		return fmt.Sprint(typed)
	}
}

type responsesAPIResponse struct {
	Output []responseOutputItem `json:"output"`
	Usage  map[string]any       `json:"usage"`
}

type responseOutputItem struct {
	Type    string            `json:"type"`
	ID      string            `json:"id"`
	Status  string            `json:"status"`
	Action  responseAction    `json:"action"`
	Content []responseContent `json:"content"`
}

type responseAction struct {
	Type  string `json:"type"`
	Query string `json:"query"`
}

type responseContent struct {
	Type        string               `json:"type"`
	Text        string               `json:"text"`
	Annotations []responseAnnotation `json:"annotations"`
}

type responseAnnotation struct {
	Type       string `json:"type"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	StartIndex int    `json:"start_index"`
	EndIndex   int    `json:"end_index"`
}
