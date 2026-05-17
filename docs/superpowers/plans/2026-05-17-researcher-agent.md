# Researcher Agent Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go CLI named `researcher` that can execute a first-version research workflow and let `industry-research` call it instead of relying only on prompt orchestration.

**Architecture:** Create an independent Go module under `researcher/` with provider-neutral research types, Bocha and Volcengine provider modules, trace-reasoning artifacts, evidence ledger, confidence rules, and CLI commands. Then redesign `industry-research/SKILL.md` into a thinner domain entrypoint that calls `researcher` and validates its artifacts.

**Tech Stack:** Go 1.22+ standard library, `net/http`, `encoding/json`, shell Makefile targets modeled after `~/Workspace/go/md2wechat-skill/Makefile`, existing Markdown skill files and Python report validator.

---

## Scope Check

The spec covers a large system. This plan implements it as one staged delivery because each task creates a working, testable slice and the CLI is the shared foundation for the skill redesign.

The first implementation does not add browser automation, database storage, UI, or closed-platform scraping. Browser-required evidence is represented in schemas only.

## File Structure

Create the Go CLI under `researcher/` so the tool stays reusable while living in this repository during the first implementation.

```text
researcher/
  go.mod
  VERSION
  Makefile
  README.md
  cmd/researcher/main.go
  internal/cli/cli.go
  internal/config/config.go
  internal/rerrors/errors.go
  internal/retrieval/types.go
  internal/retrieval/capabilities.go
  internal/output/json.go
  internal/provider/bocha/bocha.go
  internal/provider/bocha/bocha_test.go
  internal/provider/volcengine/volcengine.go
  internal/provider/volcengine/volcengine_test.go
  internal/provider/multi/multi.go
  internal/provider/multi/multi_test.go
  internal/project/workspace.go
  internal/project/workspace_test.go
  internal/trace/chain_brand.go
  internal/trace/chain_brand_test.go
  internal/ledger/ledger.go
  internal/ledger/ledger_test.go
  internal/confidence/confidence.go
  internal/confidence/confidence_test.go
  internal/report/report.go
  internal/report/report_test.go
  internal/validate/validate.go
  internal/validate/validate_test.go
```

Modify existing project files:

```text
SKILL.md
references/chain-brand-trace-reasoning.md
references/evidence-ledger-schema.md
evals/evals.json
scripts/validate_report.py
```

The existing `agents/*.md` files should be updated only after the CLI artifact workflow exists. They become reviewers of `claim_graph.json`, `trace_plan.json`, `evidence_ledger.json`, `confidence_report.json`, and `final_report.md`.

---

### Task 1: Go Module, Makefile, and CLI Skeleton

**Files:**
- Create: `researcher/go.mod`
- Create: `researcher/VERSION`
- Create: `researcher/Makefile`
- Create: `researcher/README.md`
- Create: `researcher/cmd/researcher/main.go`
- Create: `researcher/internal/cli/cli.go`
- Test: manual `make build`, `make test`, `./researcher version`

- [ ] **Step 1: Create the module files**

Use `apply_patch` to create `researcher/go.mod`:

```go
module github.com/geekjourneyx/researcher

go 1.22
```

Create `researcher/VERSION`:

```text
0.1.0
```

- [ ] **Step 2: Create Makefile**

Use `apply_patch` to create `researcher/Makefile`:

```makefile
# researcher Makefile

VERSION ?= $(shell tr -d '[:space:]' < VERSION)
LDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: all build fast release clean test fmt vet install deps help

all: build

build:
	@echo "Building researcher..."
	@echo "Version: $(VERSION)"
	@go build -trimpath -ldflags="$(LDFLAGS)" -o researcher ./cmd/researcher
	@echo "Built: ./researcher"

fast:
	@go build -trimpath -ldflags="$(LDFLAGS)" -o researcher ./cmd/researcher

release:
	@echo "Building researcher release binaries..."
	@echo "Version: $(VERSION)"
	@mkdir -p bin
	@GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-linux-amd64 ./cmd/researcher
	@GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-linux-arm64 ./cmd/researcher
	@GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-darwin-amd64 ./cmd/researcher
	@GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-darwin-arm64 ./cmd/researcher
	@GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="$(LDFLAGS)" -o bin/researcher-windows-amd64.exe ./cmd/researcher
	@chmod +x bin/*-linux* bin/*-darwin* 2>/dev/null || true
	@ls -lh bin/

clean:
	@rm -f researcher
	@rm -rf bin
	@rm -f *.log

test:
	@go test -count=1 ./...

fmt:
	@go fmt ./...
	@gofmt -w .

vet:
	@go vet ./...

install:
	@go install ./cmd/researcher

deps:
	@go mod download
	@go mod tidy

help:
	@echo "researcher Makefile commands:"
	@echo "  make build    - build current platform binary"
	@echo "  make fast     - fast current platform build"
	@echo "  make release  - build release binaries into bin/"
	@echo "  make clean    - remove build outputs"
	@echo "  make test     - run tests"
	@echo "  make fmt      - format Go code"
	@echo "  make vet      - run go vet"
	@echo "  make install  - install to GOPATH/bin"
	@echo "  make deps     - download and tidy dependencies"
```

- [ ] **Step 3: Create README**

Use `apply_patch` to create `researcher/README.md`:

```markdown
# researcher

`researcher` is a Go CLI for trace-based research workflows.

It decomposes research questions into claims, retrieves leads from providers, writes evidence ledgers, searches for disconfirmation, scores confidence, and emits reports.

First-version commands:

```bash
researcher version
researcher help
researcher capabilities --json
researcher retrieve "瑞幸 2026 门店数 招聘 扩张" --providers bocha,volcengine --json
researcher plan "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --json
researcher evidence "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --json
researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --depth standard
researcher validate researcher-workspace/luckin-2026-store-count
```

Environment variables:

```text
BOCHA_API_KEY
ARK_API_KEY
```

Search results and model answers are leads, not evidence. Evidence confidence is assigned only after source opening, browser verification, or cross-validation.
```

- [ ] **Step 4: Create CLI entrypoint**

Use `apply_patch` to create `researcher/cmd/researcher/main.go`:

```go
package main

import (
	"os"

	"github.com/geekjourneyx/researcher/internal/cli"
)

var Version = "dev"

func main() {
	os.Exit(cli.Run(os.Args[1:], Version, os.Stdout, os.Stderr))
}
```

Create `researcher/internal/cli/cli.go`:

```go
package cli

import (
	"fmt"
	"io"
)

func Run(args []string, version string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		printHelp(stdout)
		return 0
	}

	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "researcher %s\n", version)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n", args[0])
		printHelp(stderr)
		return 1
	}
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "researcher commands:")
	fmt.Fprintln(w, "  version")
	fmt.Fprintln(w, "  help")
	fmt.Fprintln(w, "  capabilities --json")
	fmt.Fprintln(w, "  retrieve QUERY --providers bocha,volcengine --json")
	fmt.Fprintln(w, "  plan QUESTION --domain chain-brand --json")
	fmt.Fprintln(w, "  evidence QUESTION --domain chain-brand --json")
	fmt.Fprintln(w, "  run QUESTION --domain chain-brand --depth standard")
	fmt.Fprintln(w, "  validate WORKSPACE_DIR")
}
```

- [ ] **Step 5: Verify build**

Run:

```bash
cd researcher && make build
```

Expected:

```text
Built: ./researcher
```

Run:

```bash
cd researcher && ./researcher version
```

Expected:

```text
researcher 0.1.0
```

Run:

```bash
cd researcher && make test
```

Expected:

```text
?    github.com/geekjourneyx/researcher/cmd/researcher [no test files]
?    github.com/geekjourneyx/researcher/internal/cli [no test files]
```

- [ ] **Step 6: Commit**

```bash
git add researcher/go.mod researcher/VERSION researcher/Makefile researcher/README.md researcher/cmd/researcher/main.go researcher/internal/cli/cli.go
git commit -m "feat: scaffold researcher cli"
```

---

### Task 2: Provider-Neutral Types, Errors, Output, and Capabilities

**Files:**
- Create: `researcher/internal/retrieval/types.go`
- Create: `researcher/internal/retrieval/capabilities.go`
- Create: `researcher/internal/rerrors/errors.go`
- Create: `researcher/internal/output/json.go`
- Modify: `researcher/internal/cli/cli.go`
- Test: `researcher/internal/retrieval/capabilities_test.go`
- Test: `researcher/internal/output/json_test.go`

- [ ] **Step 1: Write capability tests**

Create `researcher/internal/retrieval/capabilities_test.go`:

```go
package retrieval

import "testing"

func TestBuiltInCapabilities(t *testing.T) {
	caps := BuiltInCapabilities()
	if len(caps) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(caps))
	}
	if caps[0].Provider != "bocha" {
		t.Fatalf("expected first provider bocha, got %s", caps[0].Provider)
	}
	if caps[0].ProviderType != ProviderTypeDirectSearch {
		t.Fatalf("expected bocha direct_search, got %s", caps[0].ProviderType)
	}
	if caps[1].Provider != "volcengine" {
		t.Fatalf("expected second provider volcengine, got %s", caps[1].Provider)
	}
	if caps[1].ProviderType != ProviderTypeModelAnswerSearch {
		t.Fatalf("expected volcengine model_answer_search, got %s", caps[1].ProviderType)
	}
}
```

Create `researcher/internal/output/json_test.go`:

```go
package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	err := WriteJSON(&buf, map[string]string{"status": "ok"}, false)
	if err != nil {
		t.Fatalf("WriteJSON returned error: %v", err)
	}
	var decoded map[string]string
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if decoded["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", decoded["status"])
	}
}
```

- [ ] **Step 2: Run tests and confirm failure**

Run:

```bash
cd researcher && go test ./internal/retrieval ./internal/output
```

Expected failure:

```text
undefined: BuiltInCapabilities
undefined: WriteJSON
```

- [ ] **Step 3: Add retrieval types**

Create `researcher/internal/retrieval/types.go`:

```go
package retrieval

import "time"

type ProviderType string

const (
	ProviderTypeDirectSearch       ProviderType = "direct_search"
	ProviderTypeModelAnswerSearch ProviderType = "model_answer_search"
	ProviderTypeAgentNativeSearch ProviderType = "agent_native_search"
	ProviderTypeBrowserVerify     ProviderType = "browser_verification"
	ProviderTypeKnowledgeRetrieval ProviderType = "knowledge_retrieval"
)

type Mode string

const (
	ModeSearch   Mode = "search"
	ModeAnswer   Mode = "answer"
	ModeRetrieve Mode = "retrieve"
)

type RetrievalRequest struct {
	Provider     string            `json:"provider"`
	ProviderType ProviderType      `json:"provider_type,omitempty"`
	Mode         Mode              `json:"mode"`
	Query        string            `json:"query"`
	Parameters   map[string]any    `json:"parameters,omitempty"`
	Headers      map[string]string `json:"-"`
}

type RetrievalResponse struct {
	Provider       string          `json:"provider"`
	ProviderType   ProviderType    `json:"provider_type"`
	Mode           Mode            `json:"mode"`
	Query          string          `json:"query"`
	RetrievedAt    time.Time       `json:"retrieved_at"`
	Request        map[string]any  `json:"request"`
	RetrievalCalls []RetrievalCall `json:"retrieval_calls"`
	Items          []Item          `json:"items"`
	Answer         Answer          `json:"answer"`
	Usage          map[string]any  `json:"usage"`
	Errors         []Error         `json:"errors"`
}

type RetrievalCall struct {
	CallID           string         `json:"call_id"`
	Query            string         `json:"query"`
	Status           string         `json:"status"`
	ProviderAction   string         `json:"provider_action"`
	ProviderMetadata map[string]any `json:"provider_metadata,omitempty"`
}

type Item struct {
	Rank                 int            `json:"rank"`
	Title                string         `json:"title"`
	URL                  string         `json:"url"`
	DisplayURL           string         `json:"display_url,omitempty"`
	SiteName             string         `json:"site_name,omitempty"`
	SiteIcon             string         `json:"site_icon,omitempty"`
	Snippet              string         `json:"snippet,omitempty"`
	Summary              string         `json:"summary,omitempty"`
	PublishedAt          string         `json:"published_at,omitempty"`
	LastCrawledAt         string         `json:"last_crawled_at,omitempty"`
	Language             string         `json:"language,omitempty"`
	ContentType          string         `json:"content_type"`
	SourceConfidenceHint string         `json:"source_confidence_hint"`
	ProviderMetadata     map[string]any `json:"provider_metadata,omitempty"`
}

type Answer struct {
	Text      string     `json:"text"`
	Citations []Citation `json:"citations"`
}

type Citation struct {
	Index  int    `json:"index"`
	URL    string `json:"url"`
	Title  string `json:"title,omitempty"`
	Source string `json:"source,omitempty"`
}

type Error struct {
	Code           string `json:"code"`
	Message        string `json:"message"`
	ProviderStatus int    `json:"provider_status,omitempty"`
	ProviderCode   string `json:"provider_code,omitempty"`
	ProviderLogID  string `json:"provider_log_id,omitempty"`
	Retryable      bool   `json:"retryable"`
	AgentAction    string `json:"agent_action"`
	RawErrorPath   string `json:"raw_error_path,omitempty"`
}
```

- [ ] **Step 4: Add capabilities**

Create `researcher/internal/retrieval/capabilities.go`:

```go
package retrieval

type ProviderCapabilities struct {
	Provider               string       `json:"provider"`
	ProviderType           ProviderType `json:"provider_type"`
	Modes                  []Mode       `json:"modes"`
	SupportsFreshness      bool         `json:"supports_freshness"`
	SupportsIncludeDomains bool         `json:"supports_include_domains"`
	SupportsExcludeDomains bool         `json:"supports_exclude_domains"`
	SupportsSummary        bool         `json:"supports_summary"`
	SupportsLocation       bool         `json:"supports_location"`
	SupportsSources        bool         `json:"supports_sources"`
	SupportsImages         bool         `json:"supports_images"`
	SupportsModelChoice    bool         `json:"supports_model_choice"`
	ResultKinds            []string     `json:"result_kinds"`
}

func BuiltInCapabilities() []ProviderCapabilities {
	return []ProviderCapabilities{
		{
			Provider:               "bocha",
			ProviderType:           ProviderTypeDirectSearch,
			Modes:                  []Mode{ModeSearch, ModeRetrieve},
			SupportsFreshness:      true,
			SupportsIncludeDomains: true,
			SupportsExcludeDomains: true,
			SupportsSummary:        true,
			SupportsLocation:       false,
			SupportsSources:        false,
			SupportsImages:         true,
			SupportsModelChoice:    false,
			ResultKinds:            []string{"web_page", "image"},
		},
		{
			Provider:               "volcengine",
			ProviderType:           ProviderTypeModelAnswerSearch,
			Modes:                  []Mode{ModeAnswer, ModeRetrieve},
			SupportsFreshness:      false,
			SupportsIncludeDomains: false,
			SupportsExcludeDomains: false,
			SupportsSummary:        false,
			SupportsLocation:       true,
			SupportsSources:        true,
			SupportsImages:         false,
			SupportsModelChoice:    true,
			ResultKinds:            []string{"annotation_url", "answer_text", "retrieval_call"},
		},
	}
}
```

- [ ] **Step 5: Add error constants**

Create `researcher/internal/rerrors/errors.go`:

```go
package rerrors

const (
	CodeMissingAPIKey          = "missing_api_key"
	CodeInvalidArgument        = "invalid_argument"
	CodeProviderHTTPError      = "provider_http_error"
	CodeProviderAuthError      = "provider_auth_error"
	CodeProviderQuotaExhausted = "provider_quota_exhausted"
	CodeProviderRateLimited    = "provider_rate_limited"
	CodeProviderTimeout        = "provider_timeout"
	CodeProviderUnavailable    = "provider_unavailable"
	CodeProviderParseError     = "provider_parse_error"
	CodeNoRetrievalTriggered   = "no_retrieval_triggered"
	CodePartialFailure         = "partial_failure"
)

const (
	ExitSuccess              = 0
	ExitInvalidArguments     = 1
	ExitMissingCredentials   = 2
	ExitProviderFailed       = 3
	ExitProviderRateLimited  = 4
	ExitTimeout              = 5
	ExitPartialMultiProvider = 6
)
```

- [ ] **Step 6: Add JSON output helper**

Create `researcher/internal/output/json.go`:

```go
package output

import (
	"encoding/json"
	"io"
)

func WriteJSON(w io.Writer, v any, pretty bool) error {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(v)
}
```

- [ ] **Step 7: Add capabilities command**

Modify `researcher/internal/cli/cli.go` to include:

```go
case "capabilities":
	pretty := hasFlag(args[1:], "--pretty")
	if !hasFlag(args[1:], "--json") && !pretty {
		for _, cap := range retrieval.BuiltInCapabilities() {
			fmt.Fprintf(stdout, "%s\t%s\n", cap.Provider, cap.ProviderType)
		}
		return 0
	}
	if len(args) > 1 && args[1] != "--json" && args[1] != "--pretty" {
		name := args[1]
		for _, cap := range retrieval.BuiltInCapabilities() {
			if cap.Provider == name {
				if err := output.WriteJSON(stdout, cap, pretty); err != nil {
					fmt.Fprintf(stderr, "write JSON: %v\n", err)
					return 1
				}
				return 0
			}
		}
		fmt.Fprintf(stderr, "unknown provider: %s\n", name)
		return 1
	}
	if err := output.WriteJSON(stdout, retrieval.BuiltInCapabilities(), pretty); err != nil {
		fmt.Fprintf(stderr, "write JSON: %v\n", err)
		return 1
	}
	return 0
```

Also add imports:

```go
	"github.com/geekjourneyx/researcher/internal/output"
	"github.com/geekjourneyx/researcher/internal/retrieval"
```

Add helper:

```go
func hasFlag(args []string, name string) bool {
	for _, arg := range args {
		if arg == name {
			return true
		}
	}
	return false
}
```

- [ ] **Step 8: Verify tests pass**

Run:

```bash
cd researcher && go test ./...
```

Expected:

```text
ok   github.com/geekjourneyx/researcher/internal/output
ok   github.com/geekjourneyx/researcher/internal/retrieval
```

Run:

```bash
cd researcher && make build && ./researcher capabilities --json
```

Expected: JSON array containing `bocha` and `volcengine`.

- [ ] **Step 9: Commit**

```bash
git add researcher
git commit -m "feat: add researcher retrieval types and capabilities"
```

---

### Task 3: Bocha Direct Search Provider

**Files:**
- Create: `researcher/internal/provider/bocha/bocha.go`
- Test: `researcher/internal/provider/bocha/bocha_test.go`
- Modify: `researcher/internal/cli/cli.go`

- [ ] **Step 1: Write Bocha tests**

Create `researcher/internal/provider/bocha/bocha_test.go`:

```go
package bocha

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/geekjourneyx/researcher/internal/retrieval"
)

func TestSearchSuccessMapsWebPages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("unexpected auth header %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"code": 200,
			"log_id": "log-123",
			"data": {
				"queryContext": {"originalQuery": "瑞幸"},
				"webPages": {
					"value": [{
						"id": "1",
						"name": "瑞幸公告",
						"url": "https://example.com/r",
						"displayUrl": "https://example.com/r",
						"snippet": "摘要",
						"summary": "长摘要",
						"siteName": "example.com",
						"siteIcon": "https://example.com/favicon.ico",
						"datePublished": "2026-05-01T00:00:00+08:00",
						"dateLastCrawled": "2026-05-02T00:00:00Z",
						"language": "zh"
					}]
				},
				"images": {"value": []}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{
		Provider: "bocha",
		Mode:     retrieval.ModeSearch,
		Query:    "瑞幸",
		Parameters: map[string]any{
			"count":     10,
			"freshness": "oneYear",
			"summary":   true,
		},
	})
	if err != nil {
		t.Fatalf("Search returned error: %v", err)
	}
	if resp.Provider != "bocha" || resp.ProviderType != retrieval.ProviderTypeDirectSearch {
		t.Fatalf("unexpected provider fields: %#v", resp)
	}
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	item := resp.Items[0]
	if item.Title != "瑞幸公告" || item.URL != "https://example.com/r" {
		t.Fatalf("unexpected item: %#v", item)
	}
	if item.SourceConfidenceHint != "lead_only" {
		t.Fatalf("expected lead_only, got %s", item.SourceConfidenceHint)
	}
}

func TestNormalizeBochaDateLastCrawled(t *testing.T) {
	got := normalizeBochaTime("2025-02-23T08:18:30Z")
	if got != "2025-02-23T08:18:30+08:00" {
		t.Fatalf("expected normalized UTC+8 time, got %s", got)
	}
}

func TestBochaErrorMapping(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(`{"code":"429","message":"You have reached the request limit","log_id":"rate-1"}`))
	}))
	defer server.Close()

	client := NewClient("test-key", server.URL, server.Client())
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{
		Provider: "bocha",
		Mode:     retrieval.ModeSearch,
		Query:    "瑞幸",
	})
	if err == nil {
		t.Fatalf("expected error")
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(resp.Errors))
	}
	if resp.Errors[0].Code != "provider_rate_limited" {
		t.Fatalf("unexpected code: %s", resp.Errors[0].Code)
	}
	if !strings.Contains(resp.Errors[0].AgentAction, "Retry") {
		t.Fatalf("expected agent action with retry guidance, got %q", resp.Errors[0].AgentAction)
	}
}
```

- [ ] **Step 2: Run tests and confirm failure**

Run:

```bash
cd researcher && go test ./internal/provider/bocha
```

Expected failure:

```text
undefined: NewClient
```

- [ ] **Step 3: Implement Bocha provider**

Create `researcher/internal/provider/bocha/bocha.go` with:

```go
package bocha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

const DefaultEndpoint = "https://api.bocha.cn/v1/web-search"

type Client struct {
	apiKey   string
	endpoint string
	http     *http.Client
}

func NewClient(apiKey string, endpoint string, httpClient *http.Client) *Client {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	return &Client{apiKey: apiKey, endpoint: endpoint, http: httpClient}
}

func (c *Client) Search(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	resp := retrieval.RetrievalResponse{
		Provider:     "bocha",
		ProviderType: retrieval.ProviderTypeDirectSearch,
		Mode:         retrieval.ModeSearch,
		Query:        strings.TrimSpace(req.Query),
		RetrievedAt:  time.Now(),
		Request:      req.Parameters,
		Items:        []retrieval.Item{},
		Answer:       retrieval.Answer{Citations: []retrieval.Citation{}},
		Usage:        map[string]any{},
		Errors:       []retrieval.Error{},
	}
	if resp.Query == "" {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeInvalidArgument, "query is required", 400, "", "", false, "Provide a non-empty query."))
		return resp, errors.New("query is required")
	}
	if c.apiKey == "" {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeMissingAPIKey, "BOCHA_API_KEY is not set", 0, "", "", false, "Set BOCHA_API_KEY or rerun with another provider."))
		return resp, errors.New("missing BOCHA_API_KEY")
	}

	body := map[string]any{"query": resp.Query}
	for k, v := range req.Parameters {
		body[k] = v
	}
	if _, ok := body["summary"]; !ok {
		body["summary"] = true
	}
	payload, err := json.Marshal(body)
	if err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeInvalidArgument, err.Error(), 0, "", "", false, "Fix request parameters and retry."))
		return resp, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeInvalidArgument, err.Error(), 0, "", "", false, "Fix endpoint configuration."))
		return resp, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeProviderTimeout, err.Error(), 0, "", "", true, "Retry with a longer timeout or use another provider."))
		return resp, err
	}
	defer httpResp.Body.Close()

	var decoded bochaResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&decoded); err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeProviderParseError, err.Error(), httpResp.StatusCode, "", "", false, "Save raw response and inspect provider output."))
		return resp, err
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 || decoded.Code != 200 {
		mapped := mapBochaError(httpResp.StatusCode, decoded.Code, decoded.Message, decoded.LogID)
		resp.Errors = append(resp.Errors, mapped)
		return resp, fmt.Errorf("bocha error: %s", mapped.Message)
	}

	resp.RetrievalCalls = append(resp.RetrievalCalls, retrieval.RetrievalCall{
		CallID:         "bocha_http_001",
		Query:          resp.Query,
		Status:         "completed",
		ProviderAction: "web-search",
	})
	for i, item := range decoded.Data.WebPages.Value {
		published := item.DatePublished
		lastCrawled := normalizeBochaTime(item.DateLastCrawled)
		resp.Items = append(resp.Items, retrieval.Item{
			Rank:                 i + 1,
			Title:                item.Name,
			URL:                  item.URL,
			DisplayURL:           item.DisplayURL,
			SiteName:             item.SiteName,
			SiteIcon:             item.SiteIcon,
			Snippet:              item.Snippet,
			Summary:              item.Summary,
			PublishedAt:          published,
			LastCrawledAt:         lastCrawled,
			Language:             item.Language,
			ContentType:          "web_page",
			SourceConfidenceHint: "lead_only",
			ProviderMetadata: map[string]any{
				"bocha_id":                  item.ID,
				"cached_page_url":           item.CachedPageURL,
				"is_navigational":           item.IsNavigational,
				"is_family_friendly":        item.IsFamilyFriendly,
				"raw_date_last_crawled":     item.DateLastCrawled,
				"total_estimated_matches":   decoded.Data.WebPages.TotalEstimatedMatches,
				"some_results_removed":      decoded.Data.WebPages.SomeResultsRemoved,
				"provider_log_id":           decoded.LogID,
			},
		})
	}
	return resp, nil
}

type bochaResponse struct {
	Code    int         `json:"code"`
	LogID   string      `json:"log_id"`
	Message string      `json:"msg"`
	Data    bochaData   `json:"data"`
}

type bochaData struct {
	WebPages bochaWebPages `json:"webPages"`
}

type bochaWebPages struct {
	TotalEstimatedMatches int              `json:"totalEstimatedMatches"`
	SomeResultsRemoved    bool             `json:"someResultsRemoved"`
	Value                 []bochaWebResult `json:"value"`
}

type bochaWebResult struct {
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
	IsFamilyFriendly any    `json:"isFamilyFriendly"`
	IsNavigational   any    `json:"isNavigational"`
}

func normalizeBochaTime(value string) string {
	if strings.HasSuffix(value, "Z") && len(value) == len("2025-02-23T08:18:30Z") {
		return strings.TrimSuffix(value, "Z") + "+08:00"
	}
	return value
}

func mapBochaError(status int, providerCode int, message string, logID string) retrieval.Error {
	code := rerrors.CodeProviderHTTPError
	retryable := false
	action := "Inspect provider response and retry with another provider if needed."
	switch status {
	case http.StatusBadRequest:
		if strings.Contains(strings.ToLower(message), "api key") {
			code = rerrors.CodeMissingAPIKey
			action = "Set BOCHA_API_KEY and retry."
		} else {
			code = rerrors.CodeInvalidArgument
			action = "Fix Bocha request parameters and retry."
		}
	case http.StatusUnauthorized:
		code = rerrors.CodeProviderAuthError
		action = "Check BOCHA_API_KEY."
	case http.StatusForbidden:
		code = rerrors.CodeProviderQuotaExhausted
		action = "Check Bocha account balance or use another provider."
	case http.StatusTooManyRequests:
		code = rerrors.CodeProviderRateLimited
		retryable = true
		action = "Retry later with lower count, or use another provider."
	case http.StatusInternalServerError:
		code = rerrors.CodeProviderUnavailable
		retryable = true
		action = "Retry later or use another provider."
	}
	return makeError(code, message, status, fmt.Sprint(providerCode), logID, retryable, action)
}

func makeError(code string, message string, status int, providerCode string, logID string, retryable bool, action string) retrieval.Error {
	return retrieval.Error{
		Code:           code,
		Message:        message,
		ProviderStatus: status,
		ProviderCode:   providerCode,
		ProviderLogID:  logID,
		Retryable:      retryable,
		AgentAction:    action,
	}
}
```

- [ ] **Step 4: Run tests**

Run:

```bash
cd researcher && go test ./internal/provider/bocha
```

Expected:

```text
ok   github.com/geekjourneyx/researcher/internal/provider/bocha
```

- [ ] **Step 5: Add `retrieve bocha` command**

Modify `researcher/internal/cli/cli.go` so `retrieve` routes to Bocha when `--providers bocha` or `--provider bocha` is supplied. Use this parsing behavior:

```text
researcher retrieve QUERY --providers bocha --json
researcher retrieve QUERY --provider bocha --json
```

Add this branch to the command switch after the `capabilities` case:

```go
case "retrieve":
	return runRetrieve(args[1:], stdout, stderr)
```

Add `runRetrieve` with:

```go
func runRetrieve(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "retrieve requires query")
		return rerrors.ExitInvalidArguments
	}
	query := args[0]
	providers := flagValue(args[1:], "--providers")
	if providers == "" {
		providers = flagValue(args[1:], "--provider")
	}
	if providers == "" {
		providers = "bocha"
	}
	if providers != "bocha" {
		fmt.Fprintf(stderr, "unsupported retrieve provider: %s\n", providers)
		return rerrors.ExitInvalidArguments
	}
	client := bocha.NewClient(os.Getenv("BOCHA_API_KEY"), "", nil)
	resp, err := client.Search(context.Background(), retrieval.RetrievalRequest{
		Provider: "bocha",
		Mode:     retrieval.ModeSearch,
		Query:    query,
		Parameters: map[string]any{
			"count":     intFlag(args[1:], "--count", 10),
			"freshness": flagValueDefault(args[1:], "--freshness", "noLimit"),
			"summary":   true,
		},
	})
	pretty := hasFlag(args[1:], "--pretty")
	if writeErr := output.WriteJSON(stdout, resp, pretty); writeErr != nil {
		fmt.Fprintf(stderr, "write JSON: %v\n", writeErr)
		return rerrors.ExitInvalidArguments
	}
	if err != nil {
		return rerrors.ExitProviderFailed
	}
	return rerrors.ExitSuccess
}
```

Add these helper functions in `cli.go`:

```go
func flagValue(args []string, name string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == name && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(args[i], name+"=") {
			return strings.TrimPrefix(args[i], name+"=")
		}
	}
	return ""
}

func flagValueDefault(args []string, name string, defaultValue string) string {
	if value := flagValue(args, name); value != "" {
		return value
	}
	return defaultValue
}

func intFlag(args []string, name string, defaultValue int) int {
	value := flagValue(args, name)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
```

- [ ] **Step 6: Verify command with missing API key**

Run:

```bash
cd researcher && ./researcher retrieve "瑞幸" --providers bocha --json
```

Expected:

```json
{
  "provider": "bocha",
  "errors": [
    {
      "code": "missing_api_key",
      "agent_action": "Set BOCHA_API_KEY or rerun with another provider."
    }
  ]
}
```

Exit code should be non-zero.

- [ ] **Step 7: Commit**

```bash
git add researcher
git commit -m "feat: add bocha retrieval provider"
```

---

### Task 4: Volcengine Model-Answer Search Provider

**Files:**
- Create: `researcher/internal/provider/volcengine/volcengine.go`
- Test: `researcher/internal/provider/volcengine/volcengine_test.go`
- Modify: `researcher/internal/cli/cli.go`

- [ ] **Step 1: Write Volcengine tests**

Create `researcher/internal/provider/volcengine/volcengine_test.go`:

```go
package volcengine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/geekjourneyx/researcher/internal/retrieval"
)

func TestAnswerExtractsAnnotationsAndUsage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST, got %s", r.Method)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer ark-key" {
			t.Fatalf("unexpected auth header %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_1",
			"output": [
				{
					"type": "web_search_call",
					"id": "ws_1",
					"action": {"query": "瑞幸 2026 开店计划"}
				},
				{
					"type": "message",
					"content": [{
						"type": "output_text",
						"text": "瑞幸仍在扩张。",
						"annotations": [{
							"type": "url_citation",
							"url": "https://example.com/luckin",
							"title": "瑞幸新闻"
						}]
					}]
				}
			],
			"usage": {
				"tool_usage": {"web_search": 1},
				"tool_usage_details": {"web_search": {"search_engine": 1}}
			}
		}`))
	}))
	defer server.Close()

	client := NewClient("ark-key", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{
		Provider: "volcengine",
		Mode:     retrieval.ModeAnswer,
		Query:    "瑞幸 2026 开店计划",
		Parameters: map[string]any{
			"model": "doubao-seed-2-0-lite-260215",
			"limit": 10,
		},
	})
	if err != nil {
		t.Fatalf("Answer returned error: %v", err)
	}
	if resp.ProviderType != retrieval.ProviderTypeModelAnswerSearch {
		t.Fatalf("unexpected provider type: %s", resp.ProviderType)
	}
	if len(resp.RetrievalCalls) != 1 {
		t.Fatalf("expected 1 retrieval call, got %d", len(resp.RetrievalCalls))
	}
	if resp.RetrievalCalls[0].Query != "瑞幸 2026 开店计划" {
		t.Fatalf("unexpected call query: %s", resp.RetrievalCalls[0].Query)
	}
	if resp.Answer.Text != "瑞幸仍在扩张。" {
		t.Fatalf("unexpected answer text: %s", resp.Answer.Text)
	}
	if len(resp.Items) != 1 || resp.Items[0].URL != "https://example.com/luckin" {
		t.Fatalf("unexpected items: %#v", resp.Items)
	}
	if resp.Items[0].SourceConfidenceHint != "lead_only" {
		t.Fatalf("expected lead_only, got %s", resp.Items[0].SourceConfidenceHint)
	}
}

func TestNoSearchTriggered(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": "resp_2",
			"output": [{
				"type": "message",
				"content": [{
					"type": "output_text",
					"text": "没有触发搜索。",
					"annotations": []
				}]
			}]
		}`))
	}))
	defer server.Close()

	client := NewClient("ark-key", server.URL, server.Client())
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{
		Provider: "volcengine",
		Mode:     retrieval.ModeAnswer,
		Query:    "瑞幸",
	})
	if err == nil {
		t.Fatalf("expected no retrieval triggered error")
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Code != "no_retrieval_triggered" {
		t.Fatalf("unexpected errors: %#v", resp.Errors)
	}
}
```

- [ ] **Step 2: Run tests and confirm failure**

Run:

```bash
cd researcher && go test ./internal/provider/volcengine
```

Expected failure:

```text
undefined: NewClient
```

- [ ] **Step 3: Implement Volcengine provider**

Create `researcher/internal/provider/volcengine/volcengine.go`:

```go
package volcengine

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

const DefaultEndpoint = "https://ark.cn-beijing.volces.com/api/v3/responses"

type Client struct {
	apiKey   string
	endpoint string
	http     *http.Client
}

func NewClient(apiKey string, endpoint string, httpClient *http.Client) *Client {
	if endpoint == "" {
		endpoint = DefaultEndpoint
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Client{apiKey: apiKey, endpoint: endpoint, http: httpClient}
}

func (c *Client) Answer(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	resp := retrieval.RetrievalResponse{
		Provider:     "volcengine",
		ProviderType: retrieval.ProviderTypeModelAnswerSearch,
		Mode:         retrieval.ModeAnswer,
		Query:        strings.TrimSpace(req.Query),
		RetrievedAt:  time.Now(),
		Request:      req.Parameters,
		Items:        []retrieval.Item{},
		Answer:       retrieval.Answer{Citations: []retrieval.Citation{}},
		Usage:        map[string]any{},
		Errors:       []retrieval.Error{},
	}
	if resp.Query == "" {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeInvalidArgument, "query is required", 400, "", "", false, "Provide a non-empty query."))
		return resp, errors.New("query is required")
	}
	if c.apiKey == "" {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeMissingAPIKey, "ARK_API_KEY is not set", 0, "", "", false, "Set ARK_API_KEY or use Bocha/direct web search."))
		return resp, errors.New("missing ARK_API_KEY")
	}

	body := buildRequest(resp.Query, req.Parameters)
	payload, err := json.Marshal(body)
	if err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeInvalidArgument, err.Error(), 0, "", "", false, "Fix request parameters and retry."))
		return resp, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(payload))
	if err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeInvalidArgument, err.Error(), 0, "", "", false, "Fix endpoint configuration."))
		return resp, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.http.Do(httpReq)
	if err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeProviderTimeout, err.Error(), 0, "", "", true, "Retry with a longer timeout or use another provider."))
		return resp, err
	}
	defer httpResp.Body.Close()

	var decoded arkResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&decoded); err != nil {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeProviderParseError, err.Error(), httpResp.StatusCode, "", "", false, "Save raw response and inspect provider output."))
		return resp, err
	}
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		mapped := mapHTTPError(httpResp.StatusCode)
		resp.Errors = append(resp.Errors, mapped)
		return resp, fmt.Errorf("volcengine error: %s", mapped.Message)
	}

	for _, item := range decoded.Output {
		switch item.Type {
		case "web_search_call":
			resp.RetrievalCalls = append(resp.RetrievalCalls, retrieval.RetrievalCall{
				CallID:         item.ID,
				Query:          item.Action.Query,
				Status:         "completed",
				ProviderAction: "web_search",
			})
		case "message":
			for _, content := range item.Content {
				if content.Text != "" {
					resp.Answer.Text += content.Text
				}
				for idx, ann := range content.Annotations {
					citation := retrieval.Citation{Index: idx + 1, URL: ann.URL, Title: ann.Title, Source: "unknown"}
					resp.Answer.Citations = append(resp.Answer.Citations, citation)
					resp.Items = append(resp.Items, retrieval.Item{
						Rank:                 len(resp.Items) + 1,
						Title:                ann.Title,
						URL:                  ann.URL,
						DisplayURL:           ann.URL,
						SiteName:             hostFromURL(ann.URL),
						ContentType:          "annotation_url",
						SourceConfidenceHint: "lead_only",
						ProviderMetadata: map[string]any{
							"annotation_index": idx,
						},
					})
				}
			}
		}
	}
	resp.Usage = decoded.Usage
	if len(resp.RetrievalCalls) == 0 {
		resp.Errors = append(resp.Errors, makeError(rerrors.CodeNoRetrievalTriggered, "Volcengine response did not trigger web_search", 0, "", "", false, "Rewrite as an explicit search instruction, reduce ambiguity, or use Bocha direct search."))
		return resp, errors.New("no retrieval triggered")
	}
	return resp, nil
}

func buildRequest(query string, params map[string]any) map[string]any {
	model := "doubao-seed-2-0-lite-260215"
	if v, ok := params["model"].(string); ok && v != "" {
		model = v
	}
	tool := map[string]any{"type": "web_search"}
	for _, k := range []string{"limit", "max_keyword", "sources", "user_location"} {
		if v, ok := params[k]; ok {
			tool[k] = v
		}
	}
	req := map[string]any{
		"model": model,
		"input": []map[string]string{{"role": "user", "content": query}},
		"tools": []map[string]any{tool},
	}
	if v, ok := params["max_tool_calls"]; ok {
		req["max_tool_calls"] = v
	}
	return req
}

type arkResponse struct {
	ID     string           `json:"id"`
	Output []arkOutputItem  `json:"output"`
	Usage  map[string]any   `json:"usage"`
}

type arkOutputItem struct {
	Type    string           `json:"type"`
	ID      string           `json:"id"`
	Action  arkAction        `json:"action"`
	Content []arkContentItem `json:"content"`
}

type arkAction struct {
	Query string `json:"query"`
}

type arkContentItem struct {
	Type        string          `json:"type"`
	Text        string          `json:"text"`
	Annotations []arkAnnotation `json:"annotations"`
}

type arkAnnotation struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

func mapHTTPError(status int) retrieval.Error {
	switch status {
	case http.StatusBadRequest:
		return makeError(rerrors.CodeInvalidArgument, "Volcengine invalid request", status, "", "", false, "Fix request parameters. Do not send caching.")
	case http.StatusUnauthorized, http.StatusForbidden:
		return makeError(rerrors.CodeProviderAuthError, "Volcengine authentication or permission failed", status, "", "", false, "Check ARK_API_KEY and web_search permissions.")
	case http.StatusTooManyRequests:
		return makeError(rerrors.CodeProviderRateLimited, "Volcengine rate limited", status, "", "", true, "Retry later with lower max_keyword, limit, or max_tool_calls.")
	default:
		return makeError(rerrors.CodeProviderHTTPError, "Volcengine provider request failed", status, "", "", true, "Retry later or use another provider.")
	}
}

func makeError(code string, message string, status int, providerCode string, logID string, retryable bool, action string) retrieval.Error {
	return retrieval.Error{Code: code, Message: message, ProviderStatus: status, ProviderCode: providerCode, ProviderLogID: logID, Retryable: retryable, AgentAction: action}
}

func hostFromURL(value string) string {
	value = strings.TrimPrefix(value, "https://")
	value = strings.TrimPrefix(value, "http://")
	if idx := strings.Index(value, "/"); idx >= 0 {
		return value[:idx]
	}
	return value
}
```

- [ ] **Step 4: Run tests**

Run:

```bash
cd researcher && go test ./internal/provider/volcengine
```

Expected:

```text
ok   github.com/geekjourneyx/researcher/internal/provider/volcengine
```

- [ ] **Step 5: Add `answer volcengine` command**

Modify `researcher/internal/cli/cli.go`:

```go
case "answer":
	return runAnswer(args[1:], stdout, stderr)
```

Implement `runAnswer` with:

```go
func runAnswer(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) < 2 {
		fmt.Fprintln(stderr, "answer requires provider and query")
		return rerrors.ExitInvalidArguments
	}
	provider := args[0]
	query := args[1]
	if provider != "volcengine" {
		fmt.Fprintf(stderr, "unsupported answer provider: %s\n", provider)
		return rerrors.ExitInvalidArguments
	}
	client := volcengine.NewClient(os.Getenv("ARK_API_KEY"), "", nil)
	params := map[string]any{
		"model":          flagValueDefault(args[2:], "--model", "doubao-seed-2-0-lite-260215"),
		"limit":          intFlag(args[2:], "--limit", 10),
		"max_keyword":    intFlag(args[2:], "--max-keyword", 3),
		"max_tool_calls": intFlag(args[2:], "--max-tool-calls", 3),
	}
	if sources := flagValue(args[2:], "--sources"); sources != "" {
		params["sources"] = strings.Split(sources, ",")
	}
	resp, err := client.Answer(context.Background(), retrieval.RetrievalRequest{
		Provider: "volcengine",
		Mode:     retrieval.ModeAnswer,
		Query:    query,
		Parameters: params,
	})
	pretty := hasFlag(args[2:], "--pretty")
	if writeErr := output.WriteJSON(stdout, resp, pretty); writeErr != nil {
		fmt.Fprintf(stderr, "write JSON: %v\n", writeErr)
		return rerrors.ExitInvalidArguments
	}
	if err != nil {
		return rerrors.ExitProviderFailed
	}
	return rerrors.ExitSuccess
}
```

- [ ] **Step 6: Verify missing API key**

Run:

```bash
cd researcher && ./researcher answer volcengine "搜索瑞幸 2026 门店数" --json
```

Expected: JSON error with `missing_api_key` and `agent_action`.

- [ ] **Step 7: Commit**

```bash
git add researcher
git commit -m "feat: add volcengine answer provider"
```

---

### Task 5: Multi-Provider Retrieval

**Files:**
- Create: `researcher/internal/provider/multi/multi.go`
- Test: `researcher/internal/provider/multi/multi_test.go`
- Modify: `researcher/internal/cli/cli.go`

- [ ] **Step 1: Write multi-provider tests**

Create `researcher/internal/provider/multi/multi_test.go`:

```go
package multi

import (
	"context"
	"errors"
	"testing"

	"github.com/geekjourneyx/researcher/internal/retrieval"
)

type fakeProvider struct {
	name string
	resp retrieval.RetrievalResponse
	err  error
}

func (f fakeProvider) Name() string { return f.name }
func (f fakeProvider) Retrieve(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error) {
	return f.resp, f.err
}

func TestRetrievePartialFailure(t *testing.T) {
	r := New([]Provider{
		fakeProvider{name: "ok", resp: retrieval.RetrievalResponse{Provider: "ok"}},
		fakeProvider{name: "bad", resp: retrieval.RetrievalResponse{Provider: "bad"}, err: errors.New("bad failed")},
	})
	resp, err := r.Retrieve(context.Background(), retrieval.RetrievalRequest{Query: "瑞幸"})
	if err == nil {
		t.Fatalf("expected partial failure")
	}
	if len(resp.ProviderResults) != 2 {
		t.Fatalf("expected 2 provider results, got %d", len(resp.ProviderResults))
	}
	if len(resp.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(resp.Errors))
	}
}
```

- [ ] **Step 2: Run tests and confirm failure**

Run:

```bash
cd researcher && go test ./internal/provider/multi
```

Expected failure:

```text
undefined: New
```

- [ ] **Step 3: Add multi response type**

Modify `researcher/internal/retrieval/types.go` to add:

```go
type MultiResponse struct {
	Provider        string              `json:"provider"`
	ProviderType    ProviderType        `json:"provider_type"`
	Mode            Mode                `json:"mode"`
	Query           string              `json:"query"`
	RetrievedAt     time.Time           `json:"retrieved_at"`
	ProviderResults []RetrievalResponse `json:"provider_results"`
	Errors          []Error             `json:"errors"`
}
```

- [ ] **Step 4: Implement multi provider**

Create `researcher/internal/provider/multi/multi.go`:

```go
package multi

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

type Provider interface {
	Name() string
	Retrieve(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.RetrievalResponse, error)
}

type Retriever struct {
	providers []Provider
}

func New(providers []Provider) *Retriever {
	return &Retriever{providers: providers}
}

func (r *Retriever) Retrieve(ctx context.Context, req retrieval.RetrievalRequest) (retrieval.MultiResponse, error) {
	resp := retrieval.MultiResponse{
		Provider:        "multi",
		ProviderType:    "multi",
		Mode:            retrieval.ModeRetrieve,
		Query:           req.Query,
		RetrievedAt:     time.Now(),
		ProviderResults: []retrieval.RetrievalResponse{},
		Errors:          []retrieval.Error{},
	}
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, provider := range r.providers {
		provider := provider
		wg.Add(1)
		go func() {
			defer wg.Done()
			providerResp, err := provider.Retrieve(ctx, req)
			mu.Lock()
			defer mu.Unlock()
			resp.ProviderResults = append(resp.ProviderResults, providerResp)
			if err != nil {
				resp.Errors = append(resp.Errors, retrieval.Error{
					Code:        rerrors.CodePartialFailure,
					Message:     provider.Name() + " failed: " + err.Error(),
					Retryable:   true,
					AgentAction: "Use successful provider results, retry failed provider later, or lower source coverage confidence.",
				})
			}
		}()
	}
	wg.Wait()
	if len(resp.Errors) > 0 {
		return resp, errors.New("partial provider failure")
	}
	return resp, nil
}
```

- [ ] **Step 5: Run tests**

Run:

```bash
cd researcher && go test ./internal/provider/multi
```

Expected:

```text
ok   github.com/geekjourneyx/researcher/internal/provider/multi
```

- [ ] **Step 6: Commit**

```bash
git add researcher
git commit -m "feat: add multi-provider retrieval"
```

---

### Task 6: Workspace Artifacts and Project Writer

**Files:**
- Create: `researcher/internal/project/workspace.go`
- Test: `researcher/internal/project/workspace_test.go`
- Modify: `researcher/internal/retrieval/types.go`

- [ ] **Step 1: Write workspace tests**

Create `researcher/internal/project/workspace_test.go`:

```go
package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateWorkspaceWritesQuestion(t *testing.T) {
	root := t.TempDir()
	ws, err := CreateWorkspace(root, "瑞幸咖啡 2026 年门店数目标是否可信？", "chain-brand", "standard")
	if err != nil {
		t.Fatalf("CreateWorkspace error: %v", err)
	}
	if ws.Dir == "" {
		t.Fatalf("workspace dir empty")
	}
	path := filepath.Join(ws.Dir, "question.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("question.json missing: %v", err)
	}
}
```

- [ ] **Step 2: Run test and confirm failure**

Run:

```bash
cd researcher && go test ./internal/project
```

Expected failure:

```text
undefined: CreateWorkspace
```

- [ ] **Step 3: Implement workspace writer**

Create `researcher/internal/project/workspace.go`:

```go
package project

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

type Workspace struct {
	Dir string
}

type Question struct {
	UserInput string    `json:"user_input"`
	Domain    string    `json:"domain"`
	Depth     string    `json:"depth"`
	CreatedAt time.Time `json:"created_at"`
}

func CreateWorkspace(root string, question string, domain string, depth string) (Workspace, error) {
	if root == "" {
		root = "researcher-workspace"
	}
	slug := slugify(question)
	if slug == "" {
		slug = "research-" + time.Now().Format("20060102-150405")
	}
	dir := filepath.Join(root, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Workspace{}, err
	}
	q := Question{UserInput: question, Domain: domain, Depth: depth, CreatedAt: time.Now()}
	if err := WriteJSON(filepath.Join(dir, "question.json"), q); err != nil {
		return Workspace{}, err
	}
	return Workspace{Dir: dir}, nil
}

func WriteJSON(path string, v any) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func slugify(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	value = re.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-")
	if len(value) > 64 {
		value = strings.Trim(value[:64], "-")
	}
	return value
}
```

- [ ] **Step 4: Run tests**

Run:

```bash
cd researcher && go test ./internal/project
```

Expected:

```text
ok   github.com/geekjourneyx/researcher/internal/project
```

- [ ] **Step 5: Commit**

```bash
git add researcher
git commit -m "feat: add researcher workspace artifacts"
```

---

### Task 7: Chain-Brand Trace Reasoning Plan

**Files:**
- Create: `researcher/internal/trace/chain_brand.go`
- Test: `researcher/internal/trace/chain_brand_test.go`
- Create: `references/chain-brand-trace-reasoning.md`

- [ ] **Step 1: Write trace tests**

Create `researcher/internal/trace/chain_brand_test.go`:

```go
package trace

import "testing"

func TestChainBrandStoreCountTracePlan(t *testing.T) {
	plan := BuildChainBrandTracePlan("瑞幸咖啡 2026 年门店数目标是否可信？")
	if len(plan.Claims) == 0 {
		t.Fatalf("expected claims")
	}
	claim := plan.Claims[0]
	if claim.ClaimID == "" || claim.Mechanism == "" {
		t.Fatalf("claim missing id or mechanism: %#v", claim)
	}
	if len(claim.ExpectedTraces) < 3 {
		t.Fatalf("expected at least 3 traces, got %d", len(claim.ExpectedTraces))
	}
	if len(claim.DisconfirmingTraces) == 0 {
		t.Fatalf("expected disconfirming traces")
	}
}
```

- [ ] **Step 2: Run test and confirm failure**

Run:

```bash
cd researcher && go test ./internal/trace
```

Expected failure:

```text
undefined: BuildChainBrandTracePlan
```

- [ ] **Step 3: Implement chain-brand trace plan**

Create `researcher/internal/trace/chain_brand.go`:

```go
package trace

type TracePlan struct {
	Question string  `json:"question"`
	Domain   string  `json:"domain"`
	Claims   []Claim `json:"claims"`
}

type Claim struct {
	ClaimID             string          `json:"claim_id"`
	Claim               string          `json:"claim"`
	Mechanism           string          `json:"mechanism"`
	ExpectedTraces      []ExpectedTrace `json:"expected_traces"`
	SourceFamilies      []string        `json:"source_families"`
	DisconfirmingTraces []string        `json:"disconfirming_traces"`
}

type ExpectedTrace struct {
	TraceType   string `json:"trace_type"`
	Trace       string `json:"trace"`
	WhyExpected string `json:"why_expected"`
}

func BuildChainBrandTracePlan(question string) TracePlan {
	return TracePlan{
		Question: question,
		Domain:   "chain-brand",
		Claims: []Claim{
			{
				ClaimID:   "claim_store_count_growth",
				Claim:     "门店数增长或扩张目标具备经营支撑",
				Mechanism: "门店增长需要选址、招聘、供应链、数字入口和用户需求共同支撑。",
				ExpectedTraces: []ExpectedTrace{
					{TraceType: "people_org", Trace: "目标城市出现店长、店员、区域运营、拓展岗位", WhyExpected: "门店扩张前后必须补充门店和区域运营人员。"},
					{TraceType: "digital_frontend", Trace: "小程序、App、外卖平台或门店列表出现可服务门店", WhyExpected: "真实门店必须被用户发现、选择或下单。"},
					{TraceType: "physical_fulfillment", Trace: "地图 POI、点评、外卖门店页或本地开业信息出现", WhyExpected: "真实运营门店会留下可定位和可评价的经营痕迹。"},
					{TraceType: "capital_legal", Trace: "直营网点、加盟主体、分支机构或许可信息出现", WhyExpected: "经营主体和合规经营通常会留下工商或许可痕迹。"},
					{TraceType: "management_narrative", Trace: "财报、公告、管理层访谈、公众号或权威媒体披露扩张计划", WhyExpected: "上市或准上市连锁品牌通常会公开解释扩张节奏和经营口径。"},
				},
				SourceFamilies: []string{"recruiting", "map_poi", "platform_frontend", "company_registry", "company_disclosure", "media_interview", "ugc"},
				DisconfirmingTraces: []string{
					"声称覆盖城市但无门店 POI",
					"无招聘或仅总部招聘",
					"小程序或外卖平台不可下单",
					"只有通稿转载，没有独立经营痕迹",
					"用户评价长期停滞或集中反映闭店",
				},
			},
		},
	}
}
```

- [ ] **Step 4: Add reference file**

Create `references/chain-brand-trace-reasoning.md`:

```markdown
# 连锁品牌痕迹推理手册

本手册不是固定平台清单。它要求 agent 先从命题反推真实世界痕迹，再选择来源。

## 核心问题

如果一个商业说法是真的，现实世界会被怎样改变？

## 门店扩张

- 机制：新门店需要选址、招聘、供货、数字入口和用户需求。
- 预期痕迹：招聘、地图 POI、小程序门店、外卖门店、点评、开业信息、主体或许可。
- 反证痕迹：声称覆盖但无门店、无招聘、不可下单、只有通稿。

## 供应链成熟

- 机制：稳定供应需要仓、中央厨房、冷链或干线物流、库管、司机、品控、供应商。
- 预期痕迹：仓配节点、司机/库管招聘、食品或仓储许可、供应商合作、物流招标、SKU 区域可售。
- 反证痕迹：远距离覆盖但无仓配节点、无物流岗位、无区域 SKU 可售、投诉集中在履约。

## 加盟稳定

- 机制：加盟扩张需要招商、培训、督导、供应、门店运营和纠纷处理。
- 预期痕迹：加盟商主体、招商材料、区域代理、加盟纠纷、招聘、门店评价。
- 反证痕迹：高闭店、纠纷集中、招商口径大于实际门店、加盟商盈利只来自宣传材料。
```

- [ ] **Step 5: Run tests**

Run:

```bash
cd researcher && go test ./internal/trace
```

Expected:

```text
ok   github.com/geekjourneyx/researcher/internal/trace
```

- [ ] **Step 6: Commit**

```bash
git add researcher/internal/trace references/chain-brand-trace-reasoning.md
git commit -m "feat: add chain brand trace reasoning"
```

---

### Task 8: Evidence Ledger and Confidence Rules

**Files:**
- Create: `researcher/internal/ledger/ledger.go`
- Test: `researcher/internal/ledger/ledger_test.go`
- Create: `researcher/internal/confidence/confidence.go`
- Test: `researcher/internal/confidence/confidence_test.go`
- Create: `references/evidence-ledger-schema.md`

- [ ] **Step 1: Write ledger tests**

Create `researcher/internal/ledger/ledger_test.go`:

```go
package ledger

import "testing"

func TestRetrievalOnlyCannotSupportHighConfidence(t *testing.T) {
	item := EvidenceItem{VerificationStatus: "retrieval_result_only"}
	if item.CanSupportHighConfidence() {
		t.Fatalf("retrieval-only evidence must not support high confidence")
	}
}
```

Create `researcher/internal/confidence/confidence_test.go`:

```go
package confidence

import (
	"testing"

	"github.com/geekjourneyx/researcher/internal/ledger"
)

func TestHighConfidenceRequiresThreeFamiliesAndDisconfirmation(t *testing.T) {
	decision := Score([]ledger.EvidenceItem{
		{EvidenceFamily: "people_org", VerificationStatus: "source_opened"},
		{EvidenceFamily: "digital_frontend", VerificationStatus: "source_opened"},
		{EvidenceFamily: "physical_fulfillment", VerificationStatus: "cross_validated"},
	}, true, false)
	if decision.Rating != "high" {
		t.Fatalf("expected high confidence, got %s", decision.Rating)
	}
}

func TestNoDisconfirmationDowngrades(t *testing.T) {
	decision := Score([]ledger.EvidenceItem{
		{EvidenceFamily: "people_org", VerificationStatus: "source_opened"},
		{EvidenceFamily: "digital_frontend", VerificationStatus: "source_opened"},
		{EvidenceFamily: "physical_fulfillment", VerificationStatus: "cross_validated"},
	}, false, false)
	if decision.Rating == "high" {
		t.Fatalf("expected downgrade without disconfirmation")
	}
}
```

- [ ] **Step 2: Run tests and confirm failure**

Run:

```bash
cd researcher && go test ./internal/ledger ./internal/confidence
```

Expected failure:

```text
undefined: EvidenceItem
undefined: Score
```

- [ ] **Step 3: Implement ledger**

Create `researcher/internal/ledger/ledger.go`:

```go
package ledger

type EvidenceLedger struct {
	ResearchQuestion string         `json:"research_question"`
	Items            []EvidenceItem `json:"items"`
}

type EvidenceItem struct {
	EvidenceID                  string `json:"evidence_id"`
	ClaimID                     string `json:"claim_id"`
	SourceURL                   string `json:"source_url"`
	SourceTitle                 string `json:"source_title"`
	SourceType                  string `json:"source_type"`
	EvidenceFamily              string `json:"evidence_family"`
	OriginProvider              string `json:"origin_provider"`
	OriginRetrievalID           string `json:"origin_retrieval_id"`
	AccessedAt                  string `json:"accessed_at"`
	VerificationStatus          string `json:"verification_status"`
	IndependenceNote            string `json:"independence_note"`
	SupportsOrChallenges        string `json:"supports_or_challenges"`
	Summary                     string `json:"summary"`
	RequiresBrowserVerification bool   `json:"requires_browser_verification"`
	BrowserVerificationReason   string `json:"browser_verification_reason"`
}

func (e EvidenceItem) CanSupportHighConfidence() bool {
	switch e.VerificationStatus {
	case "source_opened", "browser_verified", "cross_validated":
		return true
	default:
		return false
	}
}
```

- [ ] **Step 4: Implement confidence**

Create `researcher/internal/confidence/confidence.go`:

```go
package confidence

import "github.com/geekjourneyx/researcher/internal/ledger"

type Decision struct {
	Rating          string   `json:"rating"`
	Reason          string   `json:"reason"`
	LimitingFactors []string `json:"limiting_factors"`
}

func Score(items []ledger.EvidenceItem, hasDisconfirmation bool, hasCoreContradiction bool) Decision {
	if hasCoreContradiction {
		return Decision{Rating: "suspended", Reason: "core evidence contradiction is unresolved", LimitingFactors: []string{"unresolved contradiction"}}
	}
	families := map[string]bool{}
	verifiedCount := 0
	for _, item := range items {
		if item.CanSupportHighConfidence() {
			families[item.EvidenceFamily] = true
			verifiedCount++
		}
	}
	if len(families) >= 3 && verifiedCount >= 3 && hasDisconfirmation {
		return Decision{Rating: "high", Reason: "three independent evidence families and disconfirmation attempts are present"}
	}
	if len(families) >= 2 {
		limits := []string{}
		if !hasDisconfirmation {
			limits = append(limits, "no disconfirmation attempt")
		}
		return Decision{Rating: "medium", Reason: "two or more evidence families support the claim", LimitingFactors: limits}
	}
	if len(items) > 0 {
		return Decision{Rating: "low", Reason: "evidence exists but is not independently verified across enough families"}
	}
	return Decision{Rating: "unverified", Reason: "no usable evidence items"}
}
```

- [ ] **Step 5: Add schema reference**

Create `references/evidence-ledger-schema.md`:

```markdown
# Evidence Ledger Schema

`evidence_ledger.json` records every evidence item used by `researcher` and `industry-research`.

Rules:

- Retrieval-only items cannot support high confidence.
- Model answer text cannot be final evidence.
- Search summaries cannot replace source verification.
- Provider failures must be recorded in retrieval logs.
- High confidence requires independent evidence families and disconfirmation attempts.
- Browser-required evidence must be marked.

Required evidence item fields:

```json
{
  "evidence_id": "ev_001",
  "claim_id": "claim_store_count_2026",
  "source_url": "https://example.com",
  "source_title": "示例来源",
  "source_type": "company_disclosure|official_registry|recruiting|map_poi|platform_frontend|media|social|legal|tender|ugc|retrieval_result_only",
  "evidence_family": "capital_legal|people_org|physical_fulfillment|digital_frontend|terminal_feedback|management_narrative",
  "origin_provider": "bocha|volcengine|agent_websearch|browser|manual",
  "origin_retrieval_id": "retrieval_001",
  "accessed_at": "2026-05-17T10:10:00+08:00",
  "verification_status": "retrieval_result_only|source_opened|browser_verified|cross_validated|not_accessible",
  "independence_note": "不是公司通稿转载，独立于财报口径",
  "supports_or_challenges": "supports|challenges|mixed|lead_only",
  "summary": "这条证据说明了什么",
  "requires_browser_verification": false,
  "browser_verification_reason": ""
}
```
```

- [ ] **Step 6: Run tests**

Run:

```bash
cd researcher && go test ./internal/ledger ./internal/confidence
```

Expected:

```text
ok   github.com/geekjourneyx/researcher/internal/ledger
ok   github.com/geekjourneyx/researcher/internal/confidence
```

- [ ] **Step 7: Commit**

```bash
git add researcher/internal/ledger researcher/internal/confidence references/evidence-ledger-schema.md
git commit -m "feat: add evidence ledger and confidence rules"
```

---

### Task 9: Plan, Evidence, Run, Report, and Validate Commands

**Files:**
- Create: `researcher/internal/report/report.go`
- Test: `researcher/internal/report/report_test.go`
- Create: `researcher/internal/validate/validate.go`
- Test: `researcher/internal/validate/validate_test.go`
- Modify: `researcher/internal/cli/cli.go`
- Modify: `researcher/internal/project/workspace.go`

- [ ] **Step 1: Write report tests**

Create `researcher/internal/report/report_test.go`:

```go
package report

import (
	"strings"
	"testing"

	"github.com/geekjourneyx/researcher/internal/confidence"
	"github.com/geekjourneyx/researcher/internal/trace"
)

func TestMarkdownReportIncludesConfidence(t *testing.T) {
	md := Markdown("瑞幸是否可信？", trace.BuildChainBrandTracePlan("瑞幸是否可信？"), confidence.Decision{Rating: "low", Reason: "证据不足"})
	if !strings.Contains(md, "置信度") {
		t.Fatalf("report missing confidence: %s", md)
	}
	if !strings.Contains(md, "证据不足") {
		t.Fatalf("report missing reason: %s", md)
	}
}
```

Create `researcher/internal/validate/validate_test.go`:

```go
package validate

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateWorkspaceRequiresFiles(t *testing.T) {
	dir := t.TempDir()
	err := Workspace(dir)
	if err == nil {
		t.Fatalf("expected missing files error")
	}
	for _, name := range RequiredFiles {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(`{}`), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := Workspace(dir); err != nil {
		t.Fatalf("expected valid workspace, got %v", err)
	}
}
```

- [ ] **Step 2: Run tests and confirm failure**

Run:

```bash
cd researcher && go test ./internal/report ./internal/validate
```

Expected failure:

```text
undefined: Markdown
undefined: Workspace
```

- [ ] **Step 3: Implement report generator**

Create `researcher/internal/report/report.go`:

```go
package report

import (
	"strings"

	"github.com/geekjourneyx/researcher/internal/confidence"
	"github.com/geekjourneyx/researcher/internal/trace"
)

func Markdown(question string, plan trace.TracePlan, decision confidence.Decision) string {
	var b strings.Builder
	b.WriteString("# Research Report\n\n")
	b.WriteString("## 问题\n\n")
	b.WriteString(question + "\n\n")
	b.WriteString("## 痕迹推理\n\n")
	for _, claim := range plan.Claims {
		b.WriteString("### " + claim.Claim + "\n\n")
		b.WriteString("机制：" + claim.Mechanism + "\n\n")
		b.WriteString("预期痕迹：\n\n")
		for _, tr := range claim.ExpectedTraces {
			b.WriteString("- " + tr.Trace + "：" + tr.WhyExpected + "\n")
		}
		b.WriteString("\n反证方向：\n\n")
		for _, tr := range claim.DisconfirmingTraces {
			b.WriteString("- " + tr + "\n")
		}
	}
	b.WriteString("\n## 置信度\n\n")
	b.WriteString("置信度：" + decision.Rating + "\n\n")
	b.WriteString("原因：" + decision.Reason + "\n")
	return b.String()
}
```

- [ ] **Step 4: Implement validator**

Create `researcher/internal/validate/validate.go`:

```go
package validate

import (
	"fmt"
	"os"
	"path/filepath"
)

var RequiredFiles = []string{
	"question.json",
	"research_plan.json",
	"claim_graph.json",
	"trace_plan.json",
	"retrieval_log.json",
	"evidence_ledger.json",
	"disconfirmation_log.json",
	"confidence_report.json",
	"final_report.md",
	"report_metadata.json",
}

func Workspace(dir string) error {
	for _, name := range RequiredFiles {
		if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
			return fmt.Errorf("missing %s", name)
		}
	}
	return nil
}
```

- [ ] **Step 5: Implement first `run` command**

Modify `researcher/internal/cli/cli.go`:

```go
case "run":
	return runResearch(args[1:], stdout, stderr)
case "plan":
	return runPlan(args[1:], stdout, stderr)
case "evidence":
	return runEvidence(args[1:], stdout, stderr)
case "validate":
	return runValidate(args[1:], stdout, stderr)
```

Add these command handlers in `cli.go`. Add imports for `path/filepath`, `os`, `github.com/geekjourneyx/researcher/internal/confidence`, `github.com/geekjourneyx/researcher/internal/ledger`, `github.com/geekjourneyx/researcher/internal/project`, `github.com/geekjourneyx/researcher/internal/report`, `github.com/geekjourneyx/researcher/internal/trace`, and `github.com/geekjourneyx/researcher/internal/validate`.

```go
func runPlan(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "plan requires question")
		return rerrors.ExitInvalidArguments
	}
	plan := trace.BuildChainBrandTracePlan(args[0])
	pretty := hasFlag(args[1:], "--pretty")
	if err := output.WriteJSON(stdout, plan, pretty); err != nil {
		fmt.Fprintf(stderr, "write JSON: %v\n", err)
		return rerrors.ExitInvalidArguments
	}
	return rerrors.ExitSuccess
}

func runEvidence(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "evidence requires question")
		return rerrors.ExitInvalidArguments
	}
	ledgerDoc := ledger.EvidenceLedger{
		ResearchQuestion: args[0],
		Items:            []ledger.EvidenceItem{},
	}
	pretty := hasFlag(args[1:], "--pretty")
	if err := output.WriteJSON(stdout, ledgerDoc, pretty); err != nil {
		fmt.Fprintf(stderr, "write JSON: %v\n", err)
		return rerrors.ExitInvalidArguments
	}
	return rerrors.ExitSuccess
}

func runValidate(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "validate requires workspace directory")
		return rerrors.ExitInvalidArguments
	}
	if err := validate.Workspace(args[0]); err != nil {
		fmt.Fprintf(stderr, "validate: %v\n", err)
		return rerrors.ExitInvalidArguments
	}
	fmt.Fprintln(stdout, "ok")
	return rerrors.ExitSuccess
}
```

Add `runResearch`:

```go
func runResearch(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprintln(stderr, "run requires question")
		return rerrors.ExitInvalidArguments
	}
	question := args[0]
	domain := flagValueDefault(args[1:], "--domain", "general")
	depth := flagValueDefault(args[1:], "--depth", "standard")
	root := flagValueDefault(args[1:], "--workspace-root", "researcher-workspace")

	ws, err := project.CreateWorkspace(root, question, domain, depth)
	if err != nil {
		fmt.Fprintf(stderr, "create workspace: %v\n", err)
		return rerrors.ExitInvalidArguments
	}

	tracePlan := trace.BuildChainBrandTracePlan(question)
	ledgerDoc := ledger.EvidenceLedger{
		ResearchQuestion: question,
		Items:            []ledger.EvidenceItem{},
	}
	decision := confidence.Decision{
		Rating: "unverified",
		Reason: "trace plan generated, no independently verified evidence collected yet",
	}
	markdown := report.Markdown(question, tracePlan, decision)

	files := map[string]any{
		"research_plan.json":       map[string]any{"question": question, "domain": domain, "depth": depth},
		"claim_graph.json":         tracePlan.Claims,
		"trace_plan.json":          tracePlan,
		"retrieval_log.json":       []any{},
		"evidence_ledger.json":     ledgerDoc,
		"disconfirmation_log.json": []any{},
		"confidence_report.json":   decision,
		"report_metadata.json":     map[string]any{"workspace": ws.Dir, "execution_mode": "normal"},
	}
	for name, value := range files {
		if err := project.WriteJSON(filepath.Join(ws.Dir, name), value); err != nil {
			fmt.Fprintf(stderr, "write %s: %v\n", name, err)
			return rerrors.ExitInvalidArguments
		}
	}
	if err := os.WriteFile(filepath.Join(ws.Dir, "final_report.md"), []byte(markdown), 0o644); err != nil {
		fmt.Fprintf(stderr, "write final_report.md: %v\n", err)
		return rerrors.ExitInvalidArguments
	}
	if err := output.WriteJSON(stdout, map[string]string{"workspace": ws.Dir}, hasFlag(args[1:], "--pretty")); err != nil {
		fmt.Fprintf(stderr, "write JSON: %v\n", err)
		return rerrors.ExitInvalidArguments
	}
	return rerrors.ExitSuccess
}
```

- [ ] **Step 6: Run tests**

Run:

```bash
cd researcher && go test ./...
```

Expected: all Go tests pass.

- [ ] **Step 7: Verify CLI workflow**

Run:

```bash
cd researcher && make build
./researcher plan "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --json
./researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --depth standard
./researcher validate researcher-workspace/*
```

Expected:

- `plan` prints JSON trace plan.
- `run` creates a workspace.
- `validate` exits `0`.

- [ ] **Step 8: Commit**

```bash
git add researcher
git commit -m "feat: add researcher run artifacts and validation"
```

---

### Task 10: Redesign `industry-research` Skill to Call `researcher`

**Files:**
- Modify: `SKILL.md`
- Modify: `agents/engagement-manager.md`
- Modify: `agents/blue-team.md`
- Modify: `agents/red-team.md`
- Modify: `agents/chief-arbitrator.md`

- [ ] **Step 1: Replace the SKILL.md execution center**

Edit `SKILL.md` so its main execution model says:

```markdown
## 新执行中心：researcher CLI

本 skill 不再把完整研究链路全部写在 prompt 编排中。行业研究请求由本 skill 负责识别、定域、定深度、定语言，然后调用 `researcher` CLI 生成可复核的研究工作区。

基本调用：

```bash
researcher run "{用户研究问题}" --domain {domain} --depth {brief|standard|comprehensive}
```

链式品牌、餐饮、零售、供应链问题必须使用：

```bash
researcher run "{用户研究问题}" --domain chain-brand --depth {depth}
```

`researcher` 必须产出：

- `question.json`
- `research_plan.json`
- `claim_graph.json`
- `trace_plan.json`
- `retrieval_log.json`
- `evidence_ledger.json`
- `disconfirmation_log.json`
- `confidence_report.json`
- `final_report.md`
- `report_metadata.json`
```
```

- [ ] **Step 2: Add quality gate language**

Add to `SKILL.md`:

```markdown
## researcher 产物质量门

交付前必须检查：

1. `trace_plan.json` 是否把结论拆成可验证命题。
2. `evidence_ledger.json` 是否存在。
3. 检索结果是否只作为 lead，而不是最终证据。
4. 高置信度结论是否至少有三类独立证据家族。
5. 是否存在反证尝试。
6. 需要浏览器验证的证据是否被明确标记。
7. `confidence_report.json` 与 `final_report.md` 是否一致。

如果 researcher 输出低置信度或悬置判断，最终回复必须直接说明原因，不得包装成确定结论。
```

- [ ] **Step 3: Convert agents into artifact reviewers**

In each agent file, add this role adjustment near the top:

`agents/engagement-manager.md`:

```markdown
## researcher 集成后的角色调整

当 `researcher` 工作区存在时，你优先审阅 `claim_graph.json` 与 `trace_plan.json`，判断研究问题是否被拆成正确的可验证命题，是否遗漏关键经营痕迹，而不是重新自由生成整份研究骨架。
```

`agents/blue-team.md`:

```markdown
## researcher 集成后的角色调整

当 `researcher` 工作区存在时，你优先审阅 `evidence_ledger.json` 中支持性证据是否足够强，不能把检索摘要、模型回答或单一通稿当作高置信度证据。
```

`agents/red-team.md`:

```markdown
## researcher 集成后的角色调整

当 `researcher` 工作区存在时，你优先审阅 `disconfirmation_log.json` 与 `evidence_ledger.json`，寻找缺失痕迹、同源转载、不可解释冲突和被过度采信的宣传口径。
```

`agents/chief-arbitrator.md`:

```markdown
## researcher 集成后的角色调整

当 `researcher` 工作区存在时，你必须以 `confidence_report.json`、`evidence_ledger.json` 和 `final_report.md` 的一致性为裁决基础。不得绕过证据台账直接根据红蓝 prose 下结论。
```

- [ ] **Step 4: Verify text changes**

Run:

```bash
rg -n "researcher CLI|researcher run|trace_plan.json|evidence_ledger.json|confidence_report.json|researcher 集成后的角色调整" SKILL.md agents
```

Expected: matches in `SKILL.md` and all four agent files.

- [ ] **Step 5: Commit**

```bash
git add SKILL.md agents/engagement-manager.md agents/blue-team.md agents/red-team.md agents/chief-arbitrator.md
git commit -m "docs: integrate industry research with researcher cli"
```

---

### Task 11: Eval and Validator Updates

**Files:**
- Modify: `evals/evals.json`
- Modify: `scripts/validate_report.py`

- [ ] **Step 1: Add researcher artifact expectations to evals**

Modify chain-brand and restaurant/retail/supply-chain eval cases in `evals/evals.json` so expectations include:

```json
"工作区目录中产出trace_plan.json",
"工作区目录中产出evidence_ledger.json",
"工作区目录中产出disconfirmation_log.json",
"工作区目录中产出confidence_report.json",
"报告区分检索线索和已验证证据",
"高置信度结论必须说明三类独立证据家族"
```

- [ ] **Step 2: Extend validator for researcher artifacts**

Modify `scripts/validate_report.py` by adding optional workspace checks:

```python
RESEARCHER_REQUIRED_FILES = [
    "question.json",
    "research_plan.json",
    "claim_graph.json",
    "trace_plan.json",
    "retrieval_log.json",
    "evidence_ledger.json",
    "disconfirmation_log.json",
    "confidence_report.json",
    "final_report.md",
    "report_metadata.json",
]

def validate_researcher_workspace(workspace_dir: str) -> dict:
    results = {"valid": True, "errors": [], "warnings": [], "stats": {}}
    root = Path(workspace_dir)
    for name in RESEARCHER_REQUIRED_FILES:
        if not (root / name).exists():
            results["valid"] = False
            results["errors"].append(f"Missing researcher artifact: {name}")
    ledger_path = root / "evidence_ledger.json"
    if ledger_path.exists():
        ledger = json.loads(ledger_path.read_text(encoding="utf-8"))
        raw = json.dumps(ledger, ensure_ascii=False)
        if "retrieval_result_only" in raw and '"rating": "high"' in raw:
            results["warnings"].append("High confidence appears near retrieval-only evidence; inspect ledger")
    return results
```

Add CLI support:

```python
if "--researcher-workspace" in sys.argv:
    idx = sys.argv.index("--researcher-workspace")
    workspace_results = validate_researcher_workspace(sys.argv[idx + 1])
    print(json.dumps(workspace_results, ensure_ascii=False, indent=2))
    if not workspace_results["valid"]:
        sys.exit(1)
```

- [ ] **Step 3: Verify JSON and Python syntax**

Run:

```bash
python3 -m json.tool evals/evals.json >/tmp/industry-research-evals.json
PYTHONPYCACHEPREFIX=/private/tmp/industry-research-pycache python3 -m py_compile scripts/validate_report.py
```

Expected: both commands exit `0`.

- [ ] **Step 4: Commit**

```bash
git add evals/evals.json scripts/validate_report.py
git commit -m "test: add researcher artifact validation expectations"
```

---

### Task 12: Final Quality Gates

**Files:**
- All changed files

- [ ] **Step 1: Run Go quality gates**

Run:

```bash
cd researcher && make fmt && make vet && make test && make build
```

Expected:

- Formatting succeeds.
- `go vet ./...` exits `0`.
- `go test -count=1 ./...` exits `0`.
- `./researcher` binary is built.

- [ ] **Step 2: Run CLI smoke checks**

Run:

```bash
cd researcher
./researcher version
./researcher capabilities --json
./researcher plan "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --json
./researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --depth brief
./researcher validate researcher-workspace/*
```

Expected:

- `version` prints current version.
- `capabilities` prints valid JSON.
- `plan` prints valid JSON.
- `run` creates a workspace.
- `validate` exits `0`.

- [ ] **Step 3: Run repository checks**

Run:

```bash
python3 -m json.tool evals/evals.json >/tmp/industry-research-evals.json
PYTHONPYCACHEPREFIX=/private/tmp/industry-research-pycache python3 -m py_compile scripts/validate_report.py
rg -n "researcher run|trace_plan.json|evidence_ledger.json|confidence_report.json|retrieval_result_only|model answer" SKILL.md agents references docs/superpowers/specs docs/superpowers/plans scripts evals
```

Expected:

- JSON and Python checks exit `0`.
- `rg` finds the researcher integration terms across spec, plan, skill, references, scripts, and evals.

- [ ] **Step 4: Check for placeholders**

Run:

```bash
rg -n "T[B]D|T[O]DO|PLACE[H]OLDER|待[定]|占[位]|implement l[a]ter|fill in det[a]ils" researcher SKILL.md agents references docs/superpowers/plans docs/superpowers/specs scripts evals
```

Expected:

- No matches in newly added implementation files or docs.

- [ ] **Step 5: Commit final cleanup**

```bash
git status --short
git add researcher SKILL.md agents references scripts evals docs/superpowers/plans docs/superpowers/specs
git commit -m "chore: finalize researcher integration"
```

Only run this commit if there are final cleanup changes not already committed by earlier tasks.
