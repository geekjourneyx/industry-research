package cli

import (
	"fmt"
	"io"

	"github.com/geekjourneyx/researcher/internal/output"
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

func isFlag(arg string) bool {
	return len(arg) > 0 && arg[0] == '-'
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
