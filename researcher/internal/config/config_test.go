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

func TestLoadEffectiveExplicitPathWinsOverEnvAndDefaultPaths(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	xdg := filepath.Join(dir, "xdg")
	explicitPath := filepath.Join(dir, "explicit.yaml")
	envPath := filepath.Join(dir, "env.yaml")
	xdgPath := filepath.Join(xdg, "researcher", "config.yaml")
	homePath := filepath.Join(home, ".config", "researcher", "config.yaml")

	writeConfig(t, explicitPath, "explicit")
	writeConfig(t, envPath, "env")
	writeConfig(t, xdgPath, "xdg")
	writeConfig(t, homePath, "home")

	got, err := LoadEffective(explicitPath, func(key string) string {
		switch key {
		case "RESEARCHER_CONFIG":
			return envPath
		case "XDG_CONFIG_HOME":
			return xdg
		default:
			return ""
		}
	}, home)
	if err != nil {
		t.Fatalf("LoadEffective() error = %v", err)
	}

	if got.Defaults.Depth != "explicit" {
		t.Fatalf("Depth = %q, want explicit config value", got.Defaults.Depth)
	}
}

func TestLoadEffectiveUsesResearcherConfigWhenExplicitPathEmpty(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	xdg := filepath.Join(dir, "xdg")
	envPath := filepath.Join(dir, "env.yaml")
	xdgPath := filepath.Join(xdg, "researcher", "config.yaml")

	writeConfig(t, envPath, "env")
	writeConfig(t, xdgPath, "xdg")

	got, err := LoadEffective("", func(key string) string {
		switch key {
		case "RESEARCHER_CONFIG":
			return envPath
		case "XDG_CONFIG_HOME":
			return xdg
		default:
			return ""
		}
	}, home)
	if err != nil {
		t.Fatalf("LoadEffective() error = %v", err)
	}

	if got.Defaults.Depth != "env" {
		t.Fatalf("Depth = %q, want RESEARCHER_CONFIG value", got.Defaults.Depth)
	}
}

func TestLoadEffectivePrefersXDGConfigOverHomeConfig(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	xdg := filepath.Join(dir, "xdg")
	xdgPath := filepath.Join(xdg, "researcher", "config.yaml")
	homePath := filepath.Join(home, ".config", "researcher", "config.yaml")

	writeConfig(t, xdgPath, "xdg")
	writeConfig(t, homePath, "home")

	got, err := LoadEffective("", func(key string) string {
		if key == "XDG_CONFIG_HOME" {
			return xdg
		}
		return ""
	}, home)
	if err != nil {
		t.Fatalf("LoadEffective() error = %v", err)
	}

	if got.Defaults.Depth != "xdg" {
		t.Fatalf("Depth = %q, want XDG config value", got.Defaults.Depth)
	}
}

func TestLoadEffectiveUsesHomeConfigWhenXDGMissing(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")
	xdg := filepath.Join(dir, "xdg")
	homePath := filepath.Join(home, ".config", "researcher", "config.yaml")

	writeConfig(t, homePath, "home")

	got, err := LoadEffective("", func(key string) string {
		if key == "XDG_CONFIG_HOME" {
			return xdg
		}
		return ""
	}, home)
	if err != nil {
		t.Fatalf("LoadEffective() error = %v", err)
	}

	if got.Defaults.Depth != "home" {
		t.Fatalf("Depth = %q, want home config value", got.Defaults.Depth)
	}
}

func TestLoadEffectiveMissingDefaultConfigsFallBackToBuiltInDefaults(t *testing.T) {
	dir := t.TempDir()
	home := filepath.Join(dir, "home")

	got, err := LoadEffective("", func(string) string {
		return ""
	}, home)
	if err != nil {
		t.Fatalf("LoadEffective() error = %v", err)
	}

	if got.Defaults.Depth != "standard" {
		t.Fatalf("Depth = %q, want built-in default", got.Defaults.Depth)
	}
	if got.Defaults.WorkspaceRoot != "researcher-workspace" {
		t.Fatalf("WorkspaceRoot = %q, want built-in default", got.Defaults.WorkspaceRoot)
	}
}

func TestLoadEffectiveEnvAPIKeysOverrideFileValuesAfterLoading(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	data := []byte(`
providers:
  bocha:
    api_key: file-bocha
  volcengine:
    api_key: file-ark
`)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}

	got, err := LoadEffective(path, func(key string) string {
		switch key {
		case "BOCHA_API_KEY":
			return "env-bocha"
		case "ARK_API_KEY":
			return "env-ark"
		default:
			return ""
		}
	}, dir)
	if err != nil {
		t.Fatalf("LoadEffective() error = %v", err)
	}

	if got.Providers.Bocha.APIKey != "env-bocha" {
		t.Fatalf("Bocha API key = %q, want env override", got.Providers.Bocha.APIKey)
	}
	if got.Providers.Volcengine.APIKey != "env-ark" {
		t.Fatalf("Volcengine API key = %q, want env override", got.Providers.Volcengine.APIKey)
	}
}

func writeConfig(t *testing.T, path, depth string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	data := []byte("defaults:\n  depth: " + depth + "\n")
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
}
