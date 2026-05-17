package multi

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

const providerName = "multi"

type Provider interface {
	Name() string
	Retrieve(context.Context, retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error)
}

type Retriever struct {
	providers []Provider
}

func New(providers []Provider) *Retriever {
	return &Retriever{providers: providers}
}

func (r *Retriever) Retrieve(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.MultiResponse, error) {
	resp := retrieval.MultiResponse{
		Provider:        providerName,
		ProviderType:    retrieval.ProviderTypeMulti,
		Mode:            req.Mode,
		Query:           req.Query,
		RetrievedAt:     time.Now(),
		ProviderResults: []retrieval.RetrievalResponse{},
		Errors:          []retrieval.Error{},
	}
	if len(r.providers) == 0 {
		return resp, nil
	}

	results := make([]retrieval.RetrievalResponse, len(r.providers))
	errorsByProvider := make([][]retrieval.Error, len(r.providers))
	failed := make([]bool, len(r.providers))

	var wg sync.WaitGroup
	var mu sync.Mutex
	for i, provider := range r.providers {
		wg.Add(1)
		go func(index int, p Provider) {
			defer wg.Done()

			providerResp, err := p.Retrieve(ctx, req)
			mu.Lock()
			defer mu.Unlock()

			results[index] = providerResp
			if err == nil {
				return
			}
			failed[index] = true
			errorsByProvider[index] = providerErrors(p.Name(), providerResp.Errors, err)
		}(i, provider)
	}
	wg.Wait()

	resp.ProviderResults = results
	failureCount := 0
	for i := range errorsByProvider {
		if !failed[i] {
			continue
		}
		failureCount++
		resp.Errors = append(resp.Errors, errorsByProvider[i]...)
	}
	if ctxErr := ctx.Err(); ctxErr != nil {
		return resp, ctxErr
	}
	if failureCount == len(r.providers) {
		resp.Errors = append([]retrieval.Error{{
			Code:        rerrors.CodeProviderUnavailable,
			Message:     "all providers failed",
			Retryable:   true,
			AgentAction: "No provider returned usable results. Inspect provider-specific errors, fix credentials or request constraints, then retry.",
		}}, resp.Errors...)
		return resp, errors.New("multi retrieval: all providers failed")
	}
	if failureCount > 0 {
		return resp, errors.New("multi retrieval: partial provider failure")
	}
	return resp, nil
}

func providerErrors(provider string, embeddedErrors []retrieval.Error, err error) []retrieval.Error {
	if errors.Is(err, context.Canceled) {
		return append([]retrieval.Error{{
			Code:        rerrors.CodeRequestCanceled,
			Message:     fmt.Sprintf("%s canceled: %v", provider, err),
			Retryable:   false,
			AgentAction: "The request context was canceled; retry only if cancellation was unintended.",
		}}, embeddedErrors...)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return append([]retrieval.Error{{
			Code:        rerrors.CodeProviderTimeout,
			Message:     fmt.Sprintf("%s timed out: %v", provider, err),
			Retryable:   true,
			AgentAction: "Retry the provider or increase the request timeout if the query requires longer retrieval.",
		}}, embeddedErrors...)
	}
	return append([]retrieval.Error{{
		Code:        rerrors.CodePartialFailure,
		Message:     fmt.Sprintf("%s failed: %v", provider, err),
		Retryable:   true,
		AgentAction: "Retry the failed provider, inspect its provider-specific error, or continue with successful provider results.",
	}}, embeddedErrors...)
}
