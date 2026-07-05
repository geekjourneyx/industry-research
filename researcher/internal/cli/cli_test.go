package cli

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/geekjourneyx/researcher/internal/rerrors"
)

func TestRunNoArgsPrintsHelpAndReturnsZero(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run(nil, "test-version", &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "researcher commands:") {
		t.Fatalf("stdout = %q, want help text", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunPlanWritesTracePlanJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"plan", "瑞幸咖啡 2026 年门店数目标是否可信？", "--domain", "chain-brand", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitSuccess {
		t.Fatalf("Run() code = %d, want success", code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if decoded["domain"] != "chain-brand" {
		t.Fatalf("domain = %v, want chain-brand", decoded["domain"])
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunEvidenceWritesEmptyLedgerJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"evidence", "瑞幸是否可信？", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitSuccess {
		t.Fatalf("Run() code = %d, want success", code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if decoded["research_question"] != "瑞幸是否可信？" {
		t.Fatalf("research_question = %v", decoded["research_question"])
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunCreatesValidWorkspace(t *testing.T) {
	root := t.TempDir()
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"run", "瑞幸咖啡 2026 年门店数目标是否可信？", "--domain", "chain-brand", "--depth", "standard", "--workspace-root", root, "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitSuccess {
		t.Fatalf("Run() code = %d, want success, stderr=%q", code, stderr.String())
	}
	var decoded map[string]string
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	workspace := decoded["workspace"]
	if workspace == "" {
		t.Fatalf("workspace path missing from output")
	}
	for _, name := range []string{"question.json", "trace_plan.json", "evidence_ledger.json", "confidence_report.json", "final_report.md"} {
		if _, err := os.Stat(filepath.Join(workspace, name)); err != nil {
			t.Fatalf("%s missing: %v", name, err)
		}
	}

	stdout.Reset()
	stderr.Reset()
	code = Run([]string{"validate", workspace}, "test-version", &stdout, &stderr)
	if code != rerrors.ExitSuccess {
		t.Fatalf("validate code = %d, want success, stderr=%q", code, stderr.String())
	}
	if strings.TrimSpace(stdout.String()) != "ok" {
		t.Fatalf("validate stdout = %q, want ok", stdout.String())
	}
}

func TestRunHelpAliasesPrintHelpAndReturnZero(t *testing.T) {
	for _, arg := range []string{"help", "--help", "-h"} {
		t.Run(arg, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{arg}, "test-version", &stdout, &stderr)

			if code != 0 {
				t.Fatalf("Run(%q) code = %d, want 0", arg, code)
			}
			if !strings.Contains(stdout.String(), "researcher commands:") {
				t.Fatalf("stdout = %q, want help text", stdout.String())
			}
			if stderr.String() != "" {
				t.Fatalf("stderr = %q, want empty", stderr.String())
			}
		})
	}
}

func TestRunVersionPrintsVersionAndReturnsZero(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"version"}, "1.2.3", &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	if stdout.String() != "researcher 1.2.3\n" {
		t.Fatalf("stdout = %q, want version output", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunUnknownCommandWritesErrorAndHelpToStderr(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"unknown"}, "test-version", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Run() code = %d, want 1", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown command: unknown") {
		t.Fatalf("stderr = %q, want unknown command error", stderr.String())
	}
	if !strings.Contains(stderr.String(), "researcher commands:") {
		t.Fatalf("stderr = %q, want help text", stderr.String())
	}
}

func TestRunCapabilitiesTextListContainsBuiltInProviders(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"capabilities"}, "test-version", &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "bocha\tdirect_search") {
		t.Fatalf("stdout = %q, want bocha text row", stdout.String())
	}
	if !strings.Contains(stdout.String(), "volcengine\tmodel_answer_search") {
		t.Fatalf("stdout = %q, want volcengine text row", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunCapabilitiesAllJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"capabilities", "--json"}, "test-version", &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	var decoded []map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if len(decoded) != 2 {
		t.Fatalf("decoded length = %d, want 2", len(decoded))
	}
	if decoded[0]["provider"] != "bocha" {
		t.Fatalf("first provider = %v, want bocha", decoded[0]["provider"])
	}
	if decoded[1]["provider"] != "volcengine" {
		t.Fatalf("second provider = %v, want volcengine", decoded[1]["provider"])
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunCapabilitiesSingleProviderJSON(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"capabilities", "bocha", "--json"}, "test-version", &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if decoded["provider"] != "bocha" {
		t.Fatalf("provider = %v, want bocha", decoded["provider"])
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunCapabilitiesSingleProviderJSONWithFlagFirst(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"capabilities", "--json", "bocha"}, "test-version", &stdout, &stderr)

	if code != 0 {
		t.Fatalf("Run() code = %d, want 0", code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if decoded["provider"] != "bocha" {
		t.Fatalf("provider = %v, want bocha", decoded["provider"])
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunCapabilitiesUnknownProviderReturnsOne(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"capabilities", "unknown"}, "test-version", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Run() code = %d, want 1", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown provider: unknown") {
		t.Fatalf("stderr = %q, want unknown provider error", stderr.String())
	}
}

func TestRunCapabilitiesUnknownFlagReturnsOne(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"capabilities", "--bogus"}, "test-version", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Run() code = %d, want 1", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown flag: --bogus") {
		t.Fatalf("stderr = %q, want unknown flag error", stderr.String())
	}
}

func TestRunCapabilitiesExtraPositionalArgReturnsOne(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"capabilities", "bocha", "extra"}, "test-version", &stdout, &stderr)

	if code != 1 {
		t.Fatalf("Run() code = %d, want 1", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unexpected argument: extra") {
		t.Fatalf("stderr = %q, want unexpected argument error", stderr.String())
	}
}

func TestRunRetrieveDefaultsToBochaAndWritesJSONProviderFailure(t *testing.T) {
	t.Setenv("BOCHA_API_KEY", "")
	t.Setenv("RESEARCHER_CONFIG", "")
	t.Setenv("HOME", t.TempDir())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"retrieve", "瑞幸", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitMissingCredentials {
		t.Fatalf("Run() code = %d, want missing credentials exit", code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if decoded["provider"] != "bocha" {
		t.Fatalf("provider = %v, want bocha", decoded["provider"])
	}
	errorsValue, ok := decoded["errors"].([]any)
	if !ok || len(errorsValue) != 1 {
		t.Fatalf("errors = %#v, want one retrieval error", decoded["errors"])
	}
	firstError := errorsValue[0].(map[string]any)
	if firstError["code"] != rerrors.CodeMissingAPIKey {
		t.Fatalf("error code = %v, want missing_api_key", firstError["code"])
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunRetrieveUnsupportedProviderReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"retrieve", "瑞幸", "--providers", "unknown", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unsupported provider: unknown") {
		t.Fatalf("stderr = %q, want unsupported provider error", stderr.String())
	}
}

func TestRunRetrieveUnknownFlagAfterQueryReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"retrieve", "test", "--bad", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown flag: --bad") {
		t.Fatalf("stderr = %q, want unknown flag error", stderr.String())
	}
}

func TestRunRetrieveInvalidCountReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"retrieve", "test", "--count", "nope", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "invalid --count") {
		t.Fatalf("stderr = %q, want invalid count error", stderr.String())
	}
}

func TestRunRetrieveEmptyCountReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"retrieve", "test", "--count=", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "invalid --count") {
		t.Fatalf("stderr = %q, want invalid count error", stderr.String())
	}
}

func TestRunAnswerMissingProviderReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"answer"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "answer provider is required") {
		t.Fatalf("stderr = %q, want missing provider error", stderr.String())
	}
}

func TestRunAnswerMissingQueryReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"answer", "volcengine"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "answer query is required") {
		t.Fatalf("stderr = %q, want missing query error", stderr.String())
	}
}

func TestRunAnswerUnsupportedProviderReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"answer", "unknown", "瑞幸", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unsupported provider: unknown") {
		t.Fatalf("stderr = %q, want unsupported provider error", stderr.String())
	}
}

func TestRunAnswerVolcengineMissingAPIKeyWritesJSONProviderFailure(t *testing.T) {
	t.Setenv("ARK_API_KEY", "")
	t.Setenv("RESEARCHER_CONFIG", "")
	t.Setenv("HOME", t.TempDir())
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"answer", "volcengine", "瑞幸", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitMissingCredentials {
		t.Fatalf("Run() code = %d, want missing credentials exit", code)
	}
	var decoded map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &decoded); err != nil {
		t.Fatalf("stdout is not valid JSON: %v", err)
	}
	if decoded["provider"] != "volcengine" {
		t.Fatalf("provider = %v, want volcengine", decoded["provider"])
	}
	errorsValue, ok := decoded["errors"].([]any)
	if !ok || len(errorsValue) != 1 {
		t.Fatalf("errors = %#v, want one retrieval error", decoded["errors"])
	}
	firstError := errorsValue[0].(map[string]any)
	if firstError["code"] != rerrors.CodeMissingAPIKey {
		t.Fatalf("error code = %v, want missing_api_key", firstError["code"])
	}
	if stderr.String() != "" {
		t.Fatalf("stderr = %q, want empty", stderr.String())
	}
}

func TestRunAnswerUnknownFlagAfterQueryReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"answer", "volcengine", "test", "--bad", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "unknown flag: --bad") {
		t.Fatalf("stderr = %q, want unknown flag error", stderr.String())
	}
}

func TestRunAnswerSourcesWithoutValueReturnsInvalidArguments(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	code := Run([]string{"answer", "volcengine", "test", "--sources", "--json"}, "test-version", &stdout, &stderr)

	if code != rerrors.ExitInvalidArguments {
		t.Fatalf("Run() code = %d, want invalid arguments exit", code)
	}
	if stdout.String() != "" {
		t.Fatalf("stdout = %q, want empty", stdout.String())
	}
	if !strings.Contains(stderr.String(), "invalid --sources") {
		t.Fatalf("stderr = %q, want invalid sources error", stderr.String())
	}
}

func TestRunAnswerValueFlagFollowedByFlagReturnsInvalidArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "model followed by flag",
			args: []string{"answer", "volcengine", "test", "--model", "--json"},
			want: "invalid --model",
		},
		{
			name: "config followed by flag",
			args: []string{"answer", "volcengine", "test", "--config", "--json"},
			want: "invalid --config",
		},
		{
			name: "limit followed by flag",
			args: []string{"answer", "volcengine", "test", "--limit", "--json"},
			want: "invalid --limit",
		},
		{
			name: "max keyword followed by flag",
			args: []string{"answer", "volcengine", "test", "--max-keyword", "--json"},
			want: "invalid --max-keyword",
		},
		{
			name: "max tool calls followed by flag",
			args: []string{"answer", "volcengine", "test", "--max-tool-calls", "--json"},
			want: "invalid --max-tool-calls",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, "test-version", &stdout, &stderr)

			if code != rerrors.ExitInvalidArguments {
				t.Fatalf("Run() code = %d, want invalid arguments exit", code)
			}
			if stdout.String() != "" {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("stderr = %q, want %q", stderr.String(), tt.want)
			}
		})
	}
}

func TestRunAnswerInvalidSourcesReturnInvalidArguments(t *testing.T) {
	tests := []struct {
		name    string
		sources string
	}{
		{name: "unknown source", sources: "toutiao,unknown"},
		{name: "empty middle element", sources: "toutiao,,douyin"},
		{name: "empty trailing element", sources: "moji,"},
		{name: "empty source list", sources: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run([]string{"answer", "volcengine", "test", "--sources=" + tt.sources, "--json"}, "test-version", &stdout, &stderr)

			if code != rerrors.ExitInvalidArguments {
				t.Fatalf("Run() code = %d, want invalid arguments exit", code)
			}
			if stdout.String() != "" {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), "invalid --sources") {
				t.Fatalf("stderr = %q, want invalid sources error", stderr.String())
			}
		})
	}
}

func TestRunAnswerInvalidNumericFlagsReturnInvalidArguments(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "invalid limit",
			args: []string{"answer", "volcengine", "test", "--limit", "nope", "--json"},
			want: "invalid --limit",
		},
		{
			name: "empty limit",
			args: []string{"answer", "volcengine", "test", "--limit=", "--json"},
			want: "invalid --limit",
		},
		{
			name: "zero limit",
			args: []string{"answer", "volcengine", "test", "--limit", "0", "--json"},
			want: "invalid --limit",
		},
		{
			name: "negative limit",
			args: []string{"answer", "volcengine", "test", "--limit", "-1", "--json"},
			want: "invalid --limit",
		},
		{
			name: "too high limit",
			args: []string{"answer", "volcengine", "test", "--limit", "51", "--json"},
			want: "invalid --limit",
		},
		{
			name: "invalid max keyword",
			args: []string{"answer", "volcengine", "test", "--max-keyword", "nope", "--json"},
			want: "invalid --max-keyword",
		},
		{
			name: "empty max keyword",
			args: []string{"answer", "volcengine", "test", "--max-keyword=", "--json"},
			want: "invalid --max-keyword",
		},
		{
			name: "zero max keyword",
			args: []string{"answer", "volcengine", "test", "--max-keyword", "0", "--json"},
			want: "invalid --max-keyword",
		},
		{
			name: "negative max keyword",
			args: []string{"answer", "volcengine", "test", "--max-keyword", "-1", "--json"},
			want: "invalid --max-keyword",
		},
		{
			name: "too high max keyword",
			args: []string{"answer", "volcengine", "test", "--max-keyword", "51", "--json"},
			want: "invalid --max-keyword",
		},
		{
			name: "invalid max tool calls",
			args: []string{"answer", "volcengine", "test", "--max-tool-calls", "nope", "--json"},
			want: "invalid --max-tool-calls",
		},
		{
			name: "empty max tool calls",
			args: []string{"answer", "volcengine", "test", "--max-tool-calls=", "--json"},
			want: "invalid --max-tool-calls",
		},
		{
			name: "zero max tool calls",
			args: []string{"answer", "volcengine", "test", "--max-tool-calls", "0", "--json"},
			want: "invalid --max-tool-calls",
		},
		{
			name: "negative max tool calls",
			args: []string{"answer", "volcengine", "test", "--max-tool-calls", "-1", "--json"},
			want: "invalid --max-tool-calls",
		},
		{
			name: "too high max tool calls",
			args: []string{"answer", "volcengine", "test", "--max-tool-calls", "11", "--json"},
			want: "invalid --max-tool-calls",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := Run(tt.args, "test-version", &stdout, &stderr)

			if code != rerrors.ExitInvalidArguments {
				t.Fatalf("Run() code = %d, want invalid arguments exit", code)
			}
			if stdout.String() != "" {
				t.Fatalf("stdout = %q, want empty", stdout.String())
			}
			if !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("stderr = %q, want %q", stderr.String(), tt.want)
			}
		})
	}
}
