package cli

import (
	"bytes"
	"encoding/json"
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
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
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
