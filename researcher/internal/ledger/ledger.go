package ledger

type EvidenceLedger struct {
	ResearchQuestion string         `json:"research_question"`
	Items            []EvidenceItem `json:"items"`
}

type EvidenceItem struct {
	EvidenceID                  string `json:"evidence_id"`
	ClaimID                     string `json:"claim_id"`
	SourceURL                   string `json:"source_url"`
	SourceTitle                 string `json:"source_title"`
	SourceType                  string `json:"source_type"`
	EvidenceFamily              string `json:"evidence_family"`
	OriginProvider              string `json:"origin_provider"`
	OriginRetrievalID           string `json:"origin_retrieval_id"`
	AccessedAt                  string `json:"accessed_at"`
	VerificationStatus          string `json:"verification_status"`
	IndependenceNote            string `json:"independence_note"`
	SupportsOrChallenges        string `json:"supports_or_challenges"`
	Summary                     string `json:"summary"`
	RequiresBrowserVerification bool   `json:"requires_browser_verification"`
	BrowserVerificationReason   string `json:"browser_verification_reason"`
}

func (e EvidenceItem) CanSupportHighConfidence() bool {
	switch e.VerificationStatus {
	case "source_opened", "browser_verified", "cross_validated":
		return true
	default:
		return false
	}
}
