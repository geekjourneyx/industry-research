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
