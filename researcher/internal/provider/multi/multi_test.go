package multi

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

type fakeProvider struct {
	name string
	resp retrieval.RetrievalResponse
	err  error
}

func (p fakeProvider) Name() string {
	return p.name
}

func (p fakeProvider) Retrieve(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	return p.resp, p.err
}

type blockingProvider struct {
	name    string
	started chan<- struct{}
}

func (p blockingProvider) Name() string {
	return p.name
}

func (p blockingProvider) Retrieve(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	p.started <- struct{}{}
	<-ctx.Done()
	return retrieval.RetrievalResponse{Provider: p.name, Query: req.Query}, ctx.Err()
}

type delayedProvider struct {
	name  string
	delay time.Duration
}

func (p delayedProvider) Name() string {
	return p.name
}

func (p delayedProvider) Retrieve(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	time.Sleep(p.delay)
	return retrieval.RetrievalResponse{Provider: p.name, Query: req.Query}, nil
}

func TestRetrieveRecordsPartialFailureAndPreservesProviderResults(t *testing.T) {
	req := retrieval.RetrievalRequest{Query: "coffee chains", Mode: retrieval.ModeSearch}
	providers := []Provider{
		fakeProvider{
			name: "first",
			resp: retrieval.RetrievalResponse{Provider: "first", Query: req.Query},
		},
		fakeProvider{
			name: "second",
			resp: retrieval.RetrievalResponse{Provider: "second", Query: req.Query},
			err:  errors.New("provider unavailable"),
		},
	}

	resp, err := New(providers).Retrieve(context.Background(), req)
	if err == nil {
		t.Fatalf("Retrieve() error = nil, want partial failure error")
	}
	if len(resp.ProviderResults) != 2 {
		t.Fatalf("ProviderResults length = %d, want 2", len(resp.ProviderResults))
	}
	if resp.ProviderResults[0].Provider != "first" || resp.ProviderResults[1].Provider != "second" {
		t.Fatalf("ProviderResults order = [%q, %q], want [first, second]", resp.ProviderResults[0].Provider, resp.ProviderResults[1].Provider)
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(resp.Errors))
	}
	got := resp.Errors[0]
	if got.Code != rerrors.CodePartialFailure {
		t.Fatalf("error code = %q, want partial_failure", got.Code)
	}
	if !strings.Contains(got.Message, "second failed: provider unavailable") {
		t.Fatalf("Message = %q, want provider failure detail", got.Message)
	}
	if !got.Retryable {
		t.Fatalf("Retryable = false, want true")
	}
	if got.AgentAction == "" {
		t.Fatalf("AgentAction is empty, want guidance")
	}
}

func TestRetrieveAllSuccessReturnsAllProviderResults(t *testing.T) {
	req := retrieval.RetrievalRequest{Query: "EV batteries", Mode: retrieval.ModeSearch}
	providers := []Provider{
		fakeProvider{name: "first", resp: retrieval.RetrievalResponse{Provider: "first", Query: req.Query}},
		fakeProvider{name: "second", resp: retrieval.RetrievalResponse{Provider: "second", Query: req.Query}},
	}

	resp, err := New(providers).Retrieve(context.Background(), req)
	if err != nil {
		t.Fatalf("Retrieve() error = %v, want nil", err)
	}
	if len(resp.ProviderResults) != 2 {
		t.Fatalf("ProviderResults length = %d, want 2", len(resp.ProviderResults))
	}
	if resp.ProviderResults[0].Provider != "first" || resp.ProviderResults[1].Provider != "second" {
		t.Fatalf("ProviderResults order = [%q, %q], want [first, second]", resp.ProviderResults[0].Provider, resp.ProviderResults[1].Provider)
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("Errors length = %d, want 0", len(resp.Errors))
	}
}

func TestRetrieveNoProvidersReturnsValidEmptyResponse(t *testing.T) {
	req := retrieval.RetrievalRequest{Query: "semiconductors", Mode: retrieval.ModeSearch}

	resp, err := New(nil).Retrieve(context.Background(), req)
	if err != nil {
		t.Fatalf("Retrieve() error = %v, want nil", err)
	}
	if resp.Provider != "multi" {
		t.Fatalf("Provider = %q, want multi", resp.Provider)
	}
	if resp.ProviderType != retrieval.ProviderTypeMulti {
		t.Fatalf("ProviderType = %q, want multi", resp.ProviderType)
	}
	if resp.Query != req.Query {
		t.Fatalf("Query = %q, want %q", resp.Query, req.Query)
	}
	if len(resp.ProviderResults) != 0 {
		t.Fatalf("ProviderResults length = %d, want 0", len(resp.ProviderResults))
	}
	if len(resp.Errors) != 0 {
		t.Fatalf("Errors length = %d, want 0", len(resp.Errors))
	}
	if resp.RetrievedAt.IsZero() || time.Since(resp.RetrievedAt) > time.Minute {
		t.Fatalf("RetrievedAt = %v, want recent timestamp", resp.RetrievedAt)
	}
}

func TestRetrieveRecordsContextCancellationAccurately(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	req := retrieval.RetrievalRequest{Query: "cancelled query", Mode: retrieval.ModeSearch}
	providers := []Provider{
		fakeProvider{
			name: "first",
			resp: retrieval.RetrievalResponse{Provider: "first", Query: req.Query},
			err:  context.Canceled,
		},
		fakeProvider{
			name: "second",
			resp: retrieval.RetrievalResponse{Provider: "second", Query: req.Query},
			err:  context.Canceled,
		},
	}

	resp, err := New(providers).Retrieve(ctx, req)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Retrieve() error = %v, want context.Canceled", err)
	}
	if len(resp.ProviderResults) != 2 {
		t.Fatalf("ProviderResults length = %d, want 2", len(resp.ProviderResults))
	}
	if len(resp.Errors) != 2 {
		t.Fatalf("Errors length = %d, want 2", len(resp.Errors))
	}
	for _, got := range resp.Errors {
		if got.Code != rerrors.CodeRequestCanceled {
			t.Fatalf("error code = %q, want request_canceled", got.Code)
		}
		if !strings.Contains(got.Message, "context canceled") {
			t.Fatalf("Message = %q, want context cancellation detail", got.Message)
		}
		if got.Retryable {
			t.Fatalf("Retryable = true, want false")
		}
	}
}

func TestRetrieveAllProvidersFailReturnsAllFailedError(t *testing.T) {
	req := retrieval.RetrievalRequest{Query: "market sizing", Mode: retrieval.ModeSearch}
	providers := []Provider{
		fakeProvider{name: "first", resp: retrieval.RetrievalResponse{Provider: "first", Query: req.Query}, err: errors.New("missing credentials")},
		fakeProvider{name: "second", resp: retrieval.RetrievalResponse{Provider: "second", Query: req.Query}, err: errors.New("rate limited")},
	}

	resp, err := New(providers).Retrieve(context.Background(), req)
	if err == nil {
		t.Fatalf("Retrieve() error = nil, want all providers failed error")
	}
	if !strings.Contains(err.Error(), "all providers failed") {
		t.Fatalf("error = %v, want all providers failed", err)
	}
	if len(resp.ProviderResults) != 2 {
		t.Fatalf("ProviderResults length = %d, want 2", len(resp.ProviderResults))
	}
	if len(resp.Errors) != 3 {
		t.Fatalf("Errors length = %d, want 3", len(resp.Errors))
	}
	if resp.Errors[0].Code != rerrors.CodeProviderUnavailable {
		t.Fatalf("summary error code = %q, want provider_unavailable", resp.Errors[0].Code)
	}
	if !strings.Contains(resp.Errors[0].Message, "all providers failed") {
		t.Fatalf("summary Message = %q, want all providers failed", resp.Errors[0].Message)
	}
	if !strings.Contains(resp.Errors[0].AgentAction, "No provider returned usable results") {
		t.Fatalf("summary AgentAction = %q, want no usable results guidance", resp.Errors[0].AgentAction)
	}
	if resp.Errors[1].Code != rerrors.CodePartialFailure || resp.Errors[2].Code != rerrors.CodePartialFailure {
		t.Fatalf("provider error codes = [%q, %q], want partial_failure provider context", resp.Errors[1].Code, resp.Errors[2].Code)
	}
}

func TestRetrievePreservesProviderSpecificErrors(t *testing.T) {
	req := retrieval.RetrievalRequest{Query: "consumer lending", Mode: retrieval.ModeSearch}
	providers := []Provider{
		fakeProvider{
			name: "first",
			resp: retrieval.RetrievalResponse{Provider: "first", Query: req.Query},
		},
		fakeProvider{
			name: "second",
			resp: retrieval.RetrievalResponse{
				Provider: "second",
				Query:    req.Query,
				Errors: []retrieval.Error{
					{
						Code:        rerrors.CodeMissingAPIKey,
						Message:     "api key required",
						Retryable:   false,
						AgentAction: "Set the provider API key.",
					},
					{
						Code:        rerrors.CodeProviderRateLimited,
						Message:     "rate limited",
						Retryable:   true,
						AgentAction: "Retry later.",
					},
				},
			},
			err: errors.New("provider failed"),
		},
	}

	resp, err := New(providers).Retrieve(context.Background(), req)
	if err == nil {
		t.Fatalf("Retrieve() error = nil, want partial failure error")
	}
	if len(resp.Errors) != 3 {
		t.Fatalf("Errors length = %d, want 3", len(resp.Errors))
	}
	if resp.Errors[0].Code != rerrors.CodePartialFailure {
		t.Fatalf("wrapper error code = %q, want partial_failure", resp.Errors[0].Code)
	}
	if resp.Errors[1].Code != rerrors.CodeMissingAPIKey {
		t.Fatalf("provider error code = %q, want missing_api_key", resp.Errors[1].Code)
	}
	if resp.Errors[2].Code != rerrors.CodeProviderRateLimited {
		t.Fatalf("provider error code = %q, want provider_rate_limited", resp.Errors[2].Code)
	}
}

func TestRetrieveBlockingProviderObservesContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{}, 1)
	done := make(chan struct {
		resp retrieval.MultiResponse
		err  error
	}, 1)

	go func() {
		resp, err := New([]Provider{blockingProvider{name: "blocked", started: started}}).Retrieve(ctx, retrieval.RetrievalRequest{
			Query: "cancel me",
			Mode:  retrieval.ModeSearch,
		})
		done <- struct {
			resp retrieval.MultiResponse
			err  error
		}{resp: resp, err: err}
	}()

	<-started
	cancel()
	got := <-done
	if !errors.Is(got.err, context.Canceled) {
		t.Fatalf("Retrieve() error = %v, want context.Canceled", got.err)
	}
	if len(got.resp.ProviderResults) != 1 {
		t.Fatalf("ProviderResults length = %d, want 1", len(got.resp.ProviderResults))
	}
	if len(got.resp.Errors) != 1 {
		t.Fatalf("Errors length = %d, want 1", len(got.resp.Errors))
	}
	if got.resp.Errors[0].Code != rerrors.CodeRequestCanceled {
		t.Fatalf("error code = %q, want request_canceled", got.resp.Errors[0].Code)
	}
	if got.resp.Errors[0].Retryable {
		t.Fatalf("Retryable = true, want false")
	}
	if !strings.Contains(got.resp.Errors[0].AgentAction, "context was canceled") {
		t.Fatalf("AgentAction = %q, want cancellation guidance", got.resp.Errors[0].AgentAction)
	}
}

func TestRetrieveSlowFastProvidersReturnInputOrder(t *testing.T) {
	req := retrieval.RetrievalRequest{Query: "input order", Mode: retrieval.ModeSearch}
	providers := []Provider{
		delayedProvider{name: "slow", delay: 30 * time.Millisecond},
		delayedProvider{name: "fast", delay: 1 * time.Millisecond},
	}

	resp, err := New(providers).Retrieve(context.Background(), req)
	if err != nil {
		t.Fatalf("Retrieve() error = %v, want nil", err)
	}
	if len(resp.ProviderResults) != 2 {
		t.Fatalf("ProviderResults length = %d, want 2", len(resp.ProviderResults))
	}
	if resp.ProviderResults[0].Provider != "slow" || resp.ProviderResults[1].Provider != "fast" {
		t.Fatalf("ProviderResults order = [%q, %q], want [slow, fast]", resp.ProviderResults[0].Provider, resp.ProviderResults[1].Provider)
	}
}
