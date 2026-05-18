package ledger

import "testing"

func TestRetrievalOnlyCannotSupportHighConfidence(t *testing.T) {
	item := EvidenceItem{VerificationStatus: "retrieval_result_only"}
	if item.CanSupportHighConfidence() {
		t.Fatalf("retrieval-only evidence must not support high confidence")
	}
}

func TestVerifiedEvidenceCanSupportHighConfidence(t *testing.T) {
	for _, status := range []string{"source_opened", "browser_verified", "cross_validated"} {
		item := EvidenceItem{VerificationStatus: status}
		if !item.CanSupportHighConfidence() {
			t.Fatalf("%s should support high confidence", status)
		}
	}
}
