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

func TestRetrieveRecordsContextCancellationAsPartialFailures(t *testing.T) {
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
	if err == nil {
		t.Fatalf("Retrieve() error = nil, want cancellation partial failure error")
	}
	if len(resp.ProviderResults) != 2 {
		t.Fatalf("ProviderResults length = %d, want 2", len(resp.ProviderResults))
	}
	if len(resp.Errors) != 2 {
		t.Fatalf("Errors length = %d, want 2", len(resp.Errors))
	}
	for _, got := range resp.Errors {
		if got.Code != rerrors.CodePartialFailure {
			t.Fatalf("error code = %q, want partial_failure", got.Code)
		}
		if !strings.Contains(got.Message, "context canceled") {
			t.Fatalf("Message = %q, want context cancellation detail", got.Message)
		}
		if !got.Retryable {
			t.Fatalf("Retryable = false, want true")
		}
	}
}
