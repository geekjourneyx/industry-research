package trace

import "testing"

func TestChainBrandStoreCountTracePlan(t *testing.T) {
	plan := BuildChainBrandTracePlan("瑞幸咖啡 2026 年门店数目标是否可信？")
	if plan.Question == "" {
		t.Fatalf("question empty")
	}
	if plan.Domain != "chain-brand" {
		t.Fatalf("Domain = %q, want chain-brand", plan.Domain)
	}
	if len(plan.Claims) == 0 {
		t.Fatalf("expected claims")
	}

	claim := plan.Claims[0]
	if claim.ClaimID == "" || claim.Mechanism == "" {
		t.Fatalf("claim missing id or mechanism: %#v", claim)
	}
	if len(claim.ExpectedTraces) < 3 {
		t.Fatalf("expected at least 3 traces, got %d", len(claim.ExpectedTraces))
	}
	if len(claim.DisconfirmingTraces) == 0 {
		t.Fatalf("expected disconfirming traces")
	}
}

func TestChainBrandTracePlanContainsHardToFakeSourceFamilies(t *testing.T) {
	plan := BuildChainBrandTracePlan("某品牌宣称覆盖 100 城是否可信？")
	claim := plan.Claims[0]
	want := map[string]bool{
		"recruiting":         false,
		"map_poi":            false,
		"platform_frontend":  false,
		"company_registry":   false,
		"company_disclosure": false,
	}
	for _, family := range claim.SourceFamilies {
		if _, ok := want[family]; ok {
			want[family] = true
		}
	}
	for family, seen := range want {
		if !seen {
			t.Fatalf("missing source family %q in %#v", family, claim.SourceFamilies)
		}
	}
}
