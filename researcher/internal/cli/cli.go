package cli

import (
	"fmt"
	"io"
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
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", args[0])
		printHelp(stderr)
		return 1
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
