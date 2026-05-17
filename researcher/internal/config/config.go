package config

import (
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers ProvidersConfig `yaml:"providers"`
	Defaults  DefaultsConfig  `yaml:"defaults"`
}

type ProvidersConfig struct {
	Bocha      ProviderConfig `yaml:"bocha"`
	Volcengine ProviderConfig `yaml:"volcengine"`
}

type ProviderConfig struct {
	APIKey   string `yaml:"api_key"`
	Endpoint string `yaml:"endpoint"`
	Model    string `yaml:"model"`
}

type DefaultsConfig struct {
	Providers     []string `yaml:"providers"`
	Depth         string   `yaml:"depth"`
	WorkspaceRoot string   `yaml:"workspace_root"`
}

func DefaultConfig() Config {
	return Config{
		Providers: ProvidersConfig{
			Bocha: ProviderConfig{
				Endpoint: "https://api.bocha.cn/v1/web-search",
			},
			Volcengine: ProviderConfig{
				Endpoint: "https://ark.cn-beijing.volces.com/api/v3/responses",
				Model:    "doubao-seed-2-0-lite-260215",
			},
		},
		Defaults: DefaultsConfig{
			Providers:     []string{"bocha", "volcengine"},
			Depth:         "standard",
			WorkspaceRoot: "researcher-workspace",
		},
	}
}

func DefaultSearchPaths(home, xdgConfigHome string) []string {
	paths := make([]string, 0, 2)
	if xdgConfigHome != "" {
		paths = append(paths, filepath.Join(xdgConfigHome, "researcher", "config.yaml"))
	}
	if home != "" {
		paths = append(paths, filepath.Join(home, ".config", "researcher", "config.yaml"))
	}
	return paths
}

func LoadFile(path string) (Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func LoadEffective(explicitPath string, getenv func(string) string, home string) (Config, error) {
	if getenv == nil {
		getenv = os.Getenv
	}

	if explicitPath != "" {
		cfg, err := LoadFile(explicitPath)
		if err != nil {
			return Config{}, err
		}
		return ApplyEnv(cfg, getenv), nil
	}

	if envPath := getenv("RESEARCHER_CONFIG"); envPath != "" {
		cfg, err := LoadFile(envPath)
		if err != nil {
			return Config{}, err
		}
		return ApplyEnv(cfg, getenv), nil
	}

	xdgConfigHome := getenv("XDG_CONFIG_HOME")
	for _, path := range DefaultSearchPaths(home, xdgConfigHome) {
		cfg, err := LoadFile(path)
		if err == nil {
			return ApplyEnv(cfg, getenv), nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
	}

	return ApplyEnv(DefaultConfig(), getenv), nil
}

func ApplyEnv(cfg Config, getenv func(string) string) Config {
	if getenv == nil {
		getenv = os.Getenv
	}
	if value := getenv("BOCHA_API_KEY"); value != "" {
		cfg.Providers.Bocha.APIKey = value
	}
	if value := getenv("ARK_API_KEY"); value != "" {
		cfg.Providers.Volcengine.APIKey = value
	}
	return cfg
}
