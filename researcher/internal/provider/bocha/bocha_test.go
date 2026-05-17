package bocha

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

func TestSearchPostsRequestAndMapsWebPages(t *testing.T) {
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
			"code": 200,
			"log_id": "log-123",
			"msg": null,
			"data": {
				"queryContext": {"originalQuery": "瑞幸"},
				"webPages": {
					"totalEstimatedMatches": 42,
					"someResultsRemoved": false,
					"value": [{
						"id": "item-1",
						"name": "Luckin Coffee",
						"url": "https://example.com/luckin",
						"displayUrl": "example.com/luckin",
						"snippet": "Coffee chain update",
						"summary": "Luckin summary",
						"siteName": "Example",
						"siteIcon": "https://example.com/favicon.ico",
						"datePublished": "2025-02-22T00:00:00+08:00",
						"dateLastCrawled": "2025-02-23T08:18:30Z",
						"cachedPageUrl": "https://cache.example.com/luckin",
						"language": "zh",
						"isFamilyFriendly": true,
						"isNavigational": false
					}]
				}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{
		Query: " 瑞幸 ",
		Parameters: map[string]any{
			"count":     1,
			"freshness": "oneYear",
		},
	})
	if err != nil {
		t.Fatalf("Search() error = %v, want nil", err)
	}

	if method != http.MethodPost {
		t.Fatalf("method = %q, want POST", method)
	}
	if auth != "Bearer test-key" {
		t.Fatalf("Authorization = %q, want bearer token", auth)
	}
	if !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("Content-Type = %q, want application/json", contentType)
	}
	if body["query"] != "瑞幸" {
		t.Fatalf("query = %v, want trimmed query", body["query"])
	}
	if body["summary"] != true {
		t.Fatalf("summary = %v, want default true", body["summary"])
	}

	if resp.Provider != "bocha" {
		t.Fatalf("Provider = %q, want bocha", resp.Provider)
	}
	if resp.ProviderType != retrieval.ProviderTypeDirectSearch {
		t.Fatalf("ProviderType = %q, want direct_search", resp.ProviderType)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("Items length = %d, want 1", len(resp.Items))
	}
	item := resp.Items[0]
	if item.Title != "Luckin Coffee" {
		t.Fatalf("Title = %q, want mapped name", item.Title)
	}
	if item.URL != "https://example.com/luckin" {
		t.Fatalf("URL = %q, want mapped url", item.URL)
	}
	if item.ContentType != "web_page" {
		t.Fatalf("ContentType = %q, want web_page", item.ContentType)
	}
	if item.SourceConfidenceHint != "lead_only" {
		t.Fatalf("SourceConfidenceHint = %q, want lead_only", item.SourceConfidenceHint)
	}
	if item.LastCrawledAt != "2025-02-23T08:18:30+08:00" {
		t.Fatalf("LastCrawledAt = %q, want normalized Bocha time", item.LastCrawledAt)
	}
	if item.ProviderMetadata["provider_log_id"] != "log-123" {
		t.Fatalf("provider_log_id = %v, want log-123", item.ProviderMetadata["provider_log_id"])
	}
	if resp.Errors == nil {
		t.Fatalf("Errors = nil, want initialized empty slice")
	}
	if len(resp.RetrievalCalls) != 1 {
		t.Fatalf("RetrievalCalls length = %d, want 1", len(resp.RetrievalCalls))
	}
	if resp.RetrievalCalls[0].ProviderAction != "web-search" {
		t.Fatalf("ProviderAction = %q, want web-search", resp.RetrievalCalls[0].ProviderAction)
	}
}

func TestNormalizeBochaTimeTreatsZAsUTCPlusEight(t *testing.T) {
	got := normalizeBochaTime("2025-02-23T08:18:30Z")
	want := "2025-02-23T08:18:30+08:00"
	if got != want {
		t.Fatalf("normalizeBochaTime() = %q, want %q", got, want)
	}
}

func TestSearchMapsHTTP429ToRateLimitedRetrievalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"code":"429","message":"You have reached the request limit","log_id":"rate-1"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Search() error = nil, want provider error")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	got := resp.Errors[0]
	if got.Code != rerrors.CodeProviderRateLimited {
		t.Fatalf("error code = %q, want provider_rate_limited", got.Code)
	}
	if !got.Retryable {
		t.Fatalf("Retryable = false, want true")
	}
	if !strings.Contains(got.AgentAction, "Retry") {
		t.Fatalf("AgentAction = %q, want retry guidance", got.AgentAction)
	}
	if got.ProviderStatus != http.StatusTooManyRequests {
		t.Fatalf("ProviderStatus = %d, want 429", got.ProviderStatus)
	}
	if got.ProviderCode != "429" {
		t.Fatalf("ProviderCode = %q, want 429", got.ProviderCode)
	}
	if got.ProviderLogID != "rate-1" {
		t.Fatalf("ProviderLogID = %q, want rate-1", got.ProviderLogID)
	}
}

func TestSearchMapsHTTP429WithNonJSONBodyToRateLimitedRetrievalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte("rate limit exceeded"))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Search() error = nil, want provider error")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	got := resp.Errors[0]
	if got.Code != rerrors.CodeProviderRateLimited {
		t.Fatalf("error code = %q, want provider_rate_limited", got.Code)
	}
	if !got.Retryable {
		t.Fatalf("Retryable = false, want true")
	}
	if !strings.Contains(got.AgentAction, "Retry") {
		t.Fatalf("AgentAction = %q, want retry guidance", got.AgentAction)
	}
	if got.ProviderStatus != http.StatusTooManyRequests {
		t.Fatalf("ProviderStatus = %d, want 429", got.ProviderStatus)
	}
}

func TestSearchMapsHTTP403ToQuotaExhaustedRetrievalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"code":403,"message":"insufficient balance","log_id":"quota-1"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Search() error = nil, want provider error")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	got := resp.Errors[0]
	if got.Code != rerrors.CodeProviderQuotaExhausted {
		t.Fatalf("error code = %q, want provider_quota_exhausted", got.Code)
	}
	if got.Retryable {
		t.Fatalf("Retryable = true, want false")
	}
	if !strings.Contains(strings.ToLower(got.AgentAction), "balance") {
		t.Fatalf("AgentAction = %q, want balance guidance", got.AgentAction)
	}
}

func TestSearchMapsHTTP400ToInvalidArgumentRetrievalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":400,"message":"invalid request","log_id":"bad-1"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Search() error = nil, want provider error")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	if resp.Errors[0].Code != rerrors.CodeInvalidArgument {
		t.Fatalf("error code = %q, want invalid_argument", resp.Errors[0].Code)
	}
}

func TestSearchMapsHTTP400APIKeyMessageToMissingAPIKeyRetrievalError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"code":400,"message":"invalid api key","log_id":"key-1"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Search() error = nil, want provider error")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	if resp.Errors[0].Code != rerrors.CodeMissingAPIKey {
		t.Fatalf("error code = %q, want missing_api_key", resp.Errors[0].Code)
	}
}

func TestSearchSuccessfulHTTPEmptyBodyReturnsParseError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Search() error = nil, want parse error")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	if resp.Errors[0].Code != rerrors.CodeProviderParseError {
		t.Fatalf("error code = %q, want provider_parse_error", resp.Errors[0].Code)
	}
}

func TestSearchMissingAPIKeyReturnsErrorWithoutHTTPCall(t *testing.T) {
	var calls int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := NewClient("", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("Search() error = nil, want missing API key error")
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
}
