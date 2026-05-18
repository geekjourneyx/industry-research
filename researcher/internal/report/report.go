package report

import (
	"strings"

	"github.com/geekjourneyx/researcher/internal/confidence"
	"github.com/geekjourneyx/researcher/internal/trace"
)

func Markdown(question string, plan trace.TracePlan, decision confidence.Decision) string {
	var b strings.Builder
	b.WriteString("# Research Report\n\n")
	b.WriteString("## 问题\n\n")
	b.WriteString(question + "\n\n")
	b.WriteString("## 痕迹推理\n\n")
	for _, claim := range plan.Claims {
		b.WriteString("### " + claim.Claim + "\n\n")
		b.WriteString("机制：" + claim.Mechanism + "\n\n")
		b.WriteString("预期痕迹：\n\n")
		for _, tr := range claim.ExpectedTraces {
			b.WriteString("- " + tr.Trace + "：" + tr.WhyExpected + "\n")
		}
		b.WriteString("\n反证方向：\n\n")
		for _, tr := range claim.DisconfirmingTraces {
			b.WriteString("- " + tr + "\n")
		}
	}
	b.WriteString("\n## 置信度\n\n")
	b.WriteString("置信度：" + decision.Rating + "\n\n")
	b.WriteString("原因：" + decision.Reason + "\n")
	return b.String()
}
