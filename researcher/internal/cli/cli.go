package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/geekjourneyx/researcher/internal/config"
	"github.com/geekjourneyx/researcher/internal/output"
	"github.com/geekjourneyx/researcher/internal/provider/bocha"
	"github.com/geekjourneyx/researcher/internal/rerrors"
	"github.com/geekjourneyx/researcher/internal/retrieval"
)

func Run(args []string, version string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printHelp(stdout)
		return 0
	}

	switch args[0] {
	case "version":
		fmt.Fprintf(stdout, "researcher %s\n", version)
		return 0
	case "help", "--help", "-h":
		printHelp(stdout)
		return 0
	case "capabilities":
		return runCapabilities(args[1:], stdout, stderr)
	case "retrieve":
		return runRetrieve(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", args[0])
		printHelp(stderr)
		return 1
	}
}

func runCapabilities(args []string, stdout io.Writer, stderr io.Writer) int {
	caps := retrieval.BuiltInCapabilities()
	formatJSON := false
	pretty := false
	provider := ""

	for _, arg := range args {
		switch arg {
		case "--json":
			formatJSON = true
		case "--pretty":
			pretty = true
		default:
			if isFlag(arg) {
				fmt.Fprintf(stderr, "unknown flag: %s\n", arg)
				return 1
			}
			if provider != "" {
				fmt.Fprintf(stderr, "unexpected argument: %s\n", arg)
				return 1
			}
			provider = arg
		}
	}

	if provider == "" && !formatJSON && !pretty {
		for _, cap := range caps {
			fmt.Fprintf(stdout, "%s\t%s\n", cap.Provider, cap.ProviderType)
		}
		return 0
	}

	if provider != "" {
		for _, cap := range caps {
			if cap.Provider == provider {
				if err := output.WriteJSON(stdout, cap, pretty); err != nil {
					fmt.Fprintf(stderr, "write capabilities: %v\n", err)
					return 1
				}
				return 0
			}
		}
		fmt.Fprintf(stderr, "unknown provider: %s\n", provider)
		return 1
	}

	if err := output.WriteJSON(stdout, caps, pretty); err != nil {
		fmt.Fprintf(stderr, "write capabilities: %v\n", err)
		return 1
	}
	return 0
}

func runRetrieve(args []string, stdout io.Writer, stderr io.Writer) int {
	formatJSON := hasFlag(args, "--json")
	pretty := hasFlag(args, "--pretty")
	provider := flagValueDefault(args, "--providers", "bocha")
	provider = flagValueDefault(args, "--provider", provider)
	if provider != "bocha" {
		fmt.Fprintf(stderr, "unsupported provider: %s\n", provider)
		return rerrors.ExitInvalidArguments
	}

	query, err := positionalQuery(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return rerrors.ExitInvalidArguments
	}

	home, _ := os.UserHomeDir()
	cfg, err := config.LoadEffective(flagValue(args, "--config"), os.Getenv, home)
	if err != nil {
		fmt.Fprintf(stderr, "load config: %v\n", err)
		return rerrors.ExitInvalidArguments
	}

	req := retrieval.RetrievalRequest{
		Provider:     "bocha",
		ProviderType: retrieval.ProviderTypeDirectSearch,
		Mode:         retrieval.ModeSearch,
		Query:        query,
		Parameters: map[string]any{
			"count":     intFlag(args, "--count", 10),
			"freshness": flagValueDefault(args, "--freshness", "noLimit"),
			"summary":   true,
		},
	}
	client := bocha.NewClient(cfg.Providers.Bocha.APIKey, cfg.Providers.Bocha.Endpoint, nil)
	resp, searchErr := client.Search(context.Background(), req)
	if formatJSON || searchErr != nil {
		if err := output.WriteJSON(stdout, resp, pretty); err != nil {
			fmt.Fprintf(stderr, "write retrieval response: %v\n", err)
			return rerrors.ExitProviderFailed
		}
	}
	if searchErr != nil {
		return exitCodeForRetrievalErrors(resp.Errors)
	}
	if !formatJSON {
		for _, item := range resp.Items {
			fmt.Fprintf(stdout, "%d\t%s\t%s\n", item.Rank, item.Title, item.URL)
		}
	}
	return rerrors.ExitSuccess
}

func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
}

func hasFlag(args []string, name string) bool {
	for _, arg := range args {
		if arg == name {
			return true
		}
	}
	return false
}

func flagValue(args []string, name string) string {
	for i, arg := range args {
		if arg == name && i+1 < len(args) {
			return args[i+1]
		}
		prefix := name + "="
		if strings.HasPrefix(arg, prefix) {
			return strings.TrimPrefix(arg, prefix)
		}
	}
	return ""
}

func flagValueDefault(args []string, name string, fallback string) string {
	if value := flagValue(args, name); value != "" {
		return value
	}
	return fallback
}

func intFlag(args []string, name string, fallback int) int {
	value := flagValue(args, name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func positionalQuery(args []string) (string, error) {
	valueFlags := map[string]bool{
		"--config":    true,
		"--provider":  true,
		"--providers": true,
		"--count":     true,
		"--freshness": true,
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--json" || arg == "--pretty" {
			continue
		}
		if strings.Contains(arg, "=") && isFlag(arg) {
			continue
		}
		if valueFlags[arg] {
			i++
			continue
		}
		if isFlag(arg) {
			return "", fmt.Errorf("unknown flag: %s", arg)
		}
		query := strings.TrimSpace(arg)
		if query == "" {
			return "", fmt.Errorf("retrieve query is required")
		}
		return query, nil
	}
	return "", fmt.Errorf("retrieve query is required")
}

func exitCodeForRetrievalErrors(errors []retrieval.Error) int {
	if len(errors) == 0 {
		return rerrors.ExitProviderFailed
	}
	switch errors[0].Code {
	case rerrors.CodeMissingAPIKey:
		return rerrors.ExitMissingCredentials
	case rerrors.CodeProviderRateLimited:
		return rerrors.ExitProviderRateLimited
	case rerrors.CodeProviderTimeout:
		return rerrors.ExitTimeout
	case rerrors.CodeInvalidArgument:
		return rerrors.ExitInvalidArguments
	default:
		return rerrors.ExitProviderFailed
	}
}

func printHelp(w io.Writer) {
	fmt.Fprintln(w, "researcher commands:")
	fmt.Fprintln(w, "  version")
	fmt.Fprintln(w, "  help")
	fmt.Fprintln(w, "  capabilities --json")
	fmt.Fprintln(w, "  retrieve")
	fmt.Fprintln(w, "  plan")
	fmt.Fprintln(w, "  evidence")
	fmt.Fprintln(w, "  run")
	fmt.Fprintln(w, "  validate")
}
