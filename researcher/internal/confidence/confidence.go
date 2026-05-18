package confidence

import "github.com/geekjourneyx/researcher/internal/ledger"

type Decision struct {
	Rating          string   `json:"rating"`
	Reason          string   `json:"reason"`
	LimitingFactors []string `json:"limiting_factors"`
}

func Score(items []ledger.EvidenceItem, hasDisconfirmation bool, hasCoreContradiction bool) Decision {
	if hasCoreContradiction {
		return Decision{
			Rating:          "suspended",
			Reason:          "core evidence contradiction is unresolved",
			LimitingFactors: []string{"unresolved contradiction"},
		}
	}

	families := map[string]bool{}
	verifiedCount := 0
	for _, item := range items {
		if item.CanSupportHighConfidence() {
			families[item.EvidenceFamily] = true
			verifiedCount++
		}
	}
	if len(families) >= 3 && verifiedCount >= 3 && hasDisconfirmation {
		return Decision{Rating: "high", Reason: "three independent evidence families and disconfirmation attempts are present"}
	}
	if len(families) >= 2 {
		limits := []string{}
		if !hasDisconfirmation {
			limits = append(limits, "no disconfirmation attempt")
		}
		return Decision{Rating: "medium", Reason: "two or more evidence families support the claim", LimitingFactors: limits}
	}
	if len(items) > 0 {
		return Decision{Rating: "low", Reason: "evidence exists but is not independently verified across enough families"}
	}
	return Decision{Rating: "unverified", Reason: "no usable evidence items"}
}
