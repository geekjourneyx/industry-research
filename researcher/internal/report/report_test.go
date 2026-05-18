package report

import (
	"strings"
	"testing"

	"github.com/geekjourneyx/researcher/internal/confidence"
	"github.com/geekjourneyx/researcher/internal/trace"
)

func TestMarkdownReportIncludesConfidence(t *testing.T) {
	md := Markdown("瑞幸是否可信？", trace.BuildChainBrandTracePlan("瑞幸是否可信？"), confidence.Decision{Rating: "low", Reason: "证据不足"})
	if !strings.Contains(md, "置信度") {
		t.Fatalf("report missing confidence: %s", md)
	}
	if !strings.Contains(md, "证据不足") {
		t.Fatalf("report missing reason: %s", md)
	}
}

func TestMarkdownReportIncludesTraceReasoning(t *testing.T) {
	md := Markdown("瑞幸是否可信？", trace.BuildChainBrandTracePlan("瑞幸是否可信？"), confidence.Decision{Rating: "unverified", Reason: "未验证"})
	for _, want := range []string{"痕迹推理", "预期痕迹", "反证方向"} {
		if !strings.Contains(md, want) {
			t.Fatalf("report missing %q: %s", want, md)
		}
	}
}
