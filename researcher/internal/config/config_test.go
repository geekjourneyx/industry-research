package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDefaultSearchPaths(t *testing.T) {
	got := DefaultSearchPaths("/home/alice", "/tmp/xdg")
	want := []string{
		"/tmp/xdg/researcher/config.yaml",
		"/home/alice/.config/researcher/config.yaml",
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DefaultSearchPaths() = %#v, want %#v", got, want)
	}
}

func TestApplyEnvOverridesProviderAPIKeys(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Providers.Bocha.APIKey = "old-bocha"
	cfg.Providers.Volcengine.APIKey = "old-ark"

	got := ApplyEnv(cfg, func(key string) string {
		switch key {
		case "BOCHA_API_KEY":
			return "new-bocha"
		case "ARK_API_KEY":
			return "new-ark"
		default:
			return ""
		}
	})

	if got.Providers.Bocha.APIKey != "new-bocha" {
		t.Fatalf("Bocha API key = %q, want %q", got.Providers.Bocha.APIKey, "new-bocha")
	}
	if got.Providers.Volcengine.APIKey != "new-ark" {
		t.Fatalf("Volcengine API key = %q, want %q", got.Providers.Volcengine.APIKey, "new-ark")
	}
}

func TestLoadFileReadsExplicitYAMLAndMergesWithDefaults(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`
providers:
  bocha:
    api_key: yaml-bocha
defaults:
  depth: comprehensive
`)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := LoadFile(path)
	if err != nil {
		t.Fatalf("LoadFile() error = %v", err)
	}

	if got.Providers.Bocha.APIKey != "yaml-bocha" {
		t.Fatalf("Bocha API key = %q, want %q", got.Providers.Bocha.APIKey, "yaml-bocha")
	}
	if got.Providers.Bocha.Endpoint != "https://api.bochaai.com/v1/web-search" {
		t.Fatalf("Bocha endpoint = %q, want default endpoint", got.Providers.Bocha.Endpoint)
	}
	if got.Providers.Volcengine.Endpoint != "https://ark.cn-beijing.volces.com/api/v3/responses" {
		t.Fatalf("Volcengine endpoint = %q, want default endpoint", got.Providers.Volcengine.Endpoint)
	}
	if got.Providers.Volcengine.Model != "doubao-seed-2-0-lite-260215" {
		t.Fatalf("Volcengine model = %q, want default model", got.Providers.Volcengine.Model)
	}
	if got.Defaults.Depth != "comprehensive" {
		t.Fatalf("Depth = %q, want %q", got.Defaults.Depth, "comprehensive")
	}
	if got.Defaults.WorkspaceRoot != "researcher-workspace" {
		t.Fatalf("WorkspaceRoot = %q, want default workspace root", got.Defaults.WorkspaceRoot)
	}
}
