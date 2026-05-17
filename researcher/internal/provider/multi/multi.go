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
	errorsByProvider := make([]retrieval.Error, len(r.providers))
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
			errorsByProvider[index] = retrieval.Error{
				Code:        rerrors.CodePartialFailure,
				Message:     fmt.Sprintf("%s failed: %v", p.Name(), err),
				Retryable:   true,
				AgentAction: "Retry the failed provider, inspect its provider-specific error, or continue with successful provider results.",
			}
		}(i, provider)
	}
	wg.Wait()

	resp.ProviderResults = results
	hasFailure := false
	for i := range errorsByProvider {
		if !failed[i] {
			continue
		}
		hasFailure = true
		resp.Errors = append(resp.Errors, errorsByProvider[i])
	}
	if hasFailure {
		return resp, errors.New("multi retrieval: partial provider failure")
	}
	return resp, nil
}
