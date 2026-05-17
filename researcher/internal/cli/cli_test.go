package cli

import (
	"bytes"
	"strings"
	"testing"
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
