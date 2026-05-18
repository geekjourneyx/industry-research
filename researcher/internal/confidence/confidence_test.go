package confidence

import (
	"testing"

	"github.com/geekjourneyx/researcher/internal/ledger"
)

func TestHighConfidenceRequiresThreeFamiliesAndDisconfirmation(t *testing.T) {
	decision := Score([]ledger.EvidenceItem{
		{EvidenceFamily: "people_org", VerificationStatus: "source_opened"},
		{EvidenceFamily: "digital_frontend", VerificationStatus: "source_opened"},
		{EvidenceFamily: "physical_fulfillment", VerificationStatus: "cross_validated"},
	}, true, false)
	if decision.Rating != "high" {
		t.Fatalf("expected high confidence, got %s", decision.Rating)
	}
}

func TestNoDisconfirmationDowngrades(t *testing.T) {
	decision := Score([]ledger.EvidenceItem{
		{EvidenceFamily: "people_org", VerificationStatus: "source_opened"},
		{EvidenceFamily: "digital_frontend", VerificationStatus: "source_opened"},
		{EvidenceFamily: "physical_fulfillment", VerificationStatus: "cross_validated"},
	}, false, false)
	if decision.Rating == "high" {
		t.Fatalf("expected downgrade without disconfirmation")
	}
}

func TestCoreContradictionSuspendsConfidence(t *testing.T) {
	decision := Score([]ledger.EvidenceItem{
		{EvidenceFamily: "people_org", VerificationStatus: "source_opened"},
		{EvidenceFamily: "digital_frontend", VerificationStatus: "source_opened"},
		{EvidenceFamily: "physical_fulfillment", VerificationStatus: "cross_validated"},
	}, true, true)
	if decision.Rating != "suspended" {
		t.Fatalf("expected suspended confidence, got %s", decision.Rating)
	}
}

func TestRetrievalOnlyEvidenceRemainsLow(t *testing.T) {
	decision := Score([]ledger.EvidenceItem{
		{EvidenceFamily: "people_org", VerificationStatus: "retrieval_result_only"},
		{EvidenceFamily: "digital_frontend", VerificationStatus: "retrieval_result_only"},
		{EvidenceFamily: "physical_fulfillment", VerificationStatus: "retrieval_result_only"},
	}, true, false)
	if decision.Rating != "low" {
		t.Fatalf("expected low confidence, got %s", decision.Rating)
	}
}
