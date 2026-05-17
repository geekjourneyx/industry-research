package retrieval

import "testing"

func TestBuiltInCapabilities(t *testing.T) {
	caps := BuiltInCapabilities()
	if len(caps) != 2 {
		t.Fatalf("expected 2 providers, got %d", len(caps))
	}
	if caps[0].Provider != "bocha" {
		t.Fatalf("expected first provider bocha, got %s", caps[0].Provider)
	}
	if caps[0].ProviderType != ProviderTypeDirectSearch {
		t.Fatalf("expected bocha direct_search, got %s", caps[0].ProviderType)
	}
	if caps[1].Provider != "volcengine" {
		t.Fatalf("expected second provider volcengine, got %s", caps[1].Provider)
	}
	if caps[1].ProviderType != ProviderTypeModelAnswerSearch {
		t.Fatalf("expected volcengine model_answer_search, got %s", caps[1].ProviderType)
	}
}
