package volcengine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

func TestAnswerPostsResponsesRequestAndMapsSearchAnswerAnnotationsAndUsage(t *testing.T) {
	var method string
	var auth string
	var contentType string
	var body map[string]any

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		auth = r.Header.Get("Authorization")
		contentType = r.Header.Get("Content-Type")
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode request body: %v", err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp-123",
			"model": "doubao-test",
			"output": [
				{
					"type": "web_search_call",
					"id": "call-1",
					"status": "completed",
					"action": {"type": "search", "query": "瑞幸 2026 门店数"}
				},
				{
					"type": "message",
					"id": "msg-1",
					"content": [{
						"type": "output_text",
						"text": "瑞幸继续扩张。",
						"annotations": [
							{
								"type": "url_citation",
								"url": "https://example.com/luckin",
								"title": "Luckin store count",
								"start_index": 0,
								"end_index": 6
							},
							{
								"type": "url_citation",
								"url": "https://news.example.cn/a",
								"title": "Expansion update"
							}
						]
					}]
				}
			],
			"usage": {
				"input_tokens": 12,
				"output_tokens": 34,
				"tool_usage": {"web_search": 1},
				"tool_usage_details": {"web_search": {"search_engine": 1}}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient(" ark-key ", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{
		Query: " 搜索瑞幸 2026 门店数 ",
		Parameters: map[string]any{
			"model":          "doubao-test",
			"limit":          8,
			"max_keyword":    2,
			"max_tool_calls": 3,
			"sources":        []string{"toutiao", "douyin"},
			"user_location": map[string]any{
				"type":    "approximate",
				"country": "中国",
				"city":    "杭州",
			},
		},
	})
	if err != nil {
		t.Fatalf("Answer() error = %v, want nil", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if auth != "Bearer ark-key" {
		t.Fatalf("Authorization = %q, want bearer token", auth)
	}
	if !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}
	if body["model"] != "doubao-test" {
		t.Fatalf("model = %v, want request model", body["model"])
	}
	if body["max_tool_calls"] != float64(3) {
		t.Fatalf("max_tool_calls = %v, want 3", body["max_tool_calls"])
	}
	tools := body["tools"].([]any)
	tool := tools[0].(map[string]any)
	if tool["type"] != "web_search" || tool["limit"] != float64(8) || tool["max_keyword"] != float64(2) {
		t.Fatalf("tool = %#v, want web_search with limit/max_keyword", tool)
	}
	if _, ok := tool["sources"].([]any); !ok {
		t.Fatalf("sources = %#v, want JSON array", tool["sources"])
	}
	input := body["input"].([]any)
	firstInput := input[0].(map[string]any)
	content := firstInput["content"].([]any)
	textContent := content[0].(map[string]any)
	if textContent["type"] != "input_text" || textContent["text"] != "搜索瑞幸 2026 门店数" {
		t.Fatalf("input content = %#v, want structured trimmed query", textContent)
	}

	if resp.Provider != "volcengine" {
		t.Fatalf("Provider = %q, want volcengine", resp.Provider)
	}
	if resp.ProviderType != retrieval.ProviderTypeModelAnswerSearch {
		t.Fatalf("ProviderType = %q, want model_answer_search", resp.ProviderType)
	}
	if resp.Mode != retrieval.ModeAnswer {
		t.Fatalf("Mode = %q, want answer", resp.Mode)
	}
	if resp.Query != "搜索瑞幸 2026 门店数" {
		t.Fatalf("Query = %q, want trimmed query", resp.Query)
	}
	if resp.Request["limit"] != 8 {
		t.Fatalf("Request limit = %#v, want original parameters", resp.Request["limit"])
	}
	if len(resp.RetrievalCalls) != 1 {
		t.Fatalf("RetrievalCalls length = %d, want 1", len(resp.RetrievalCalls))
	}
	call := resp.RetrievalCalls[0]
	if call.CallID != "call-1" || call.Query != "瑞幸 2026 门店数" || call.Status != "completed" || call.ProviderAction != "web_search" {
		t.Fatalf("RetrievalCall = %#v, want mapped web search call", call)
	}
	if resp.Answer.Text != "瑞幸继续扩张。" {
		t.Fatalf("Answer.Text = %q, want output text", resp.Answer.Text)
	}
	if len(resp.Answer.Citations) != 2 {
		t.Fatalf("Citations length = %d, want 2", len(resp.Answer.Citations))
	}
	if resp.Answer.Citations[0].Index != 1 || resp.Answer.Citations[0].Source != "example.com" {
		t.Fatalf("first citation = %#v, want host source and 1-based index", resp.Answer.Citations[0])
	}
	if len(resp.Items) != 2 {
		t.Fatalf("Items length = %d, want 2 citation items", len(resp.Items))
	}
	if resp.Items[0].ContentType != "annotation_url" || resp.Items[0].SourceConfidenceHint != "lead_only" {
		t.Fatalf("first item = %#v, want annotation URL lead-only item", resp.Items[0])
	}
	if resp.Items[0].SiteName != "example.com" {
		t.Fatalf("SiteName = %q, want host", resp.Items[0].SiteName)
	}
	if resp.Items[0].ProviderMetadata["annotation_index"] != 1 {
		t.Fatalf("annotation_index = %#v, want 1", resp.Items[0].ProviderMetadata["annotation_index"])
	}
	if resp.Usage["input_tokens"] != float64(12) {
		t.Fatalf("usage input_tokens = %#v, want decoded usage", resp.Usage["input_tokens"])
	}
	if resp.Errors == nil {
		t.Fatalf("Errors = nil, want initialized empty slice")
	}
}

func TestAnswerSuccessfulHTTPWrongShapeReturnsNoRetrievalTriggered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"resp-123","usage":{"input_tokens":3}}`))
	}))
	defer server.Close()

	client := NewClient("ark-key", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Answer() error = nil, want no retrieval triggered")
	}
	assertSingleError(t, resp, rerrors.CodeNoRetrievalTriggered, 0, false)
	if resp.Usage["input_tokens"] != float64(3) {
		t.Fatalf("usage input_tokens = %#v, want decoded usage despite wrong shape", resp.Usage["input_tokens"])
	}
}

func TestAnswerNoSearchTriggeredReturnsNoRetrievalTriggered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"output": [{
				"type": "message",
				"content": [{"type": "output_text", "text": "无需联网。"}]
			}],
			"usage": {"input_tokens": 5}
		}`))
	}))
	defer server.Close()

	client := NewClient("ark-key", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "解释咖啡是什么"})
	if err == nil {
		t.Fatalf("Answer() error = nil, want no retrieval triggered")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	if resp.Errors[0].Code != rerrors.CodeNoRetrievalTriggered {
		t.Fatalf("error code = %q, want no_retrieval_triggered", resp.Errors[0].Code)
	}
	if !strings.Contains(resp.Errors[0].AgentAction, "search") {
		t.Fatalf("AgentAction = %q, want guidance mentioning search", resp.Errors[0].AgentAction)
	}
}

func TestAnswerMissingAPIKeyReturnsErrorWithoutHTTPCall(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient(" ", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Answer() error = nil, want missing API key")
	}
	if atomic.LoadInt32(&calls) != 0 {
		t.Fatalf("HTTP calls = %d, want 0", calls)
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	if resp.Errors[0].Code != rerrors.CodeMissingAPIKey {
		t.Fatalf("error code = %q, want missing_api_key", resp.Errors[0].Code)
	}
	if !strings.Contains(resp.Errors[0].AgentAction, "ARK_API_KEY") {
		t.Fatalf("AgentAction = %q, want ARK_API_KEY guidance", resp.Errors[0].AgentAction)
	}
}

func TestAnswerMapsHTTP400ToInvalidArgumentWithParameterGuidance(t *testing.T) {
	server := errorServer(http.StatusBadRequest, `{"error":{"code":"BadRequest","message":"caching is not supported"}}`)
	defer server.Close()

	client := NewClient("ark-key", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Answer() error = nil, want provider error")
	}
	assertSingleError(t, resp, rerrors.CodeInvalidArgument, http.StatusBadRequest, false)
	if !strings.Contains(resp.Errors[0].AgentAction, "parameters") || !strings.Contains(resp.Errors[0].AgentAction, "caching") {
		t.Fatalf("AgentAction = %q, want parameters and caching guidance", resp.Errors[0].AgentAction)
	}
}

func TestAnswerMapsHTTP401And403ToAuthErrorWithPermissionGuidance(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := errorServer(status, `{"error":{"code":"AuthError","message":"permission denied"}}`)
			defer server.Close()

			client := NewClient("ark-key", server.URL, server.Client())
			resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
			if err == nil {
				t.Fatalf("Answer() error = nil, want provider error")
			}
			assertSingleError(t, resp, rerrors.CodeProviderAuthError, status, false)
			if !strings.Contains(resp.Errors[0].AgentAction, "ARK_API_KEY") || !strings.Contains(resp.Errors[0].AgentAction, "web_search") {
				t.Fatalf("AgentAction = %q, want ARK_API_KEY and web_search guidance", resp.Errors[0].AgentAction)
			}
		})
	}
}

func TestAnswerMapsHTTP429ToRetryableRateLimitedWithLimitGuidance(t *testing.T) {
	server := errorServer(http.StatusTooManyRequests, `{"error":{"code":"RateLimit","message":"too many requests"}}`)
	defer server.Close()

	client := NewClient("ark-key", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Answer() error = nil, want provider error")
	}
	assertSingleError(t, resp, rerrors.CodeProviderRateLimited, http.StatusTooManyRequests, true)
	if !strings.Contains(resp.Errors[0].AgentAction, "max_keyword") || !strings.Contains(resp.Errors[0].AgentAction, "limit") {
		t.Fatalf("AgentAction = %q, want max_keyword and limit guidance", resp.Errors[0].AgentAction)
	}
}

func TestAnswerMapsNon2xxInvalidBodyByHTTPStatus(t *testing.T) {
	server := errorServer(http.StatusTooManyRequests, `not-json`)
	defer server.Close()

	client := NewClient("ark-key", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Answer() error = nil, want provider error")
	}
	assertSingleError(t, resp, rerrors.CodeProviderRateLimited, http.StatusTooManyRequests, true)
}

func TestAnswerSuccessfulHTTPInvalidOrEmptyBodyReturnsParseError(t *testing.T) {
	for name, body := range map[string]string{
		"invalid": `not-json`,
		"empty":   ``,
	} {
		t.Run(name, func(t *testing.T) {
			server := errorServer(http.StatusOK, body)
			defer server.Close()

			client := NewClient("ark-key", server.URL, server.Client())
			resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
			if err == nil {
				t.Fatalf("Answer() error = nil, want parse error")
			}
			assertSingleError(t, resp, rerrors.CodeProviderParseError, http.StatusOK, false)
		})
	}
}

func errorServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		_, _ = w.Write([]byte(body))
	}))
}

func assertSingleError(t *testing.T, resp retrieval.RetrievalResponse, code string, status int, retryable bool) {
	t.Helper()
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	got := resp.Errors[0]
	if got.Code != code {
		t.Fatalf("error code = %q, want %q", got.Code, code)
	}
	if got.ProviderStatus != status {
		t.Fatalf("ProviderStatus = %d, want %d", got.ProviderStatus, status)
	}
	if got.Retryable != retryable {
		t.Fatalf("Retryable = %v, want %v", got.Retryable, retryable)
	}
}
