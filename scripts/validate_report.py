#!/usr/bin/env python3
"""
Validate the structure of an industry research report.
Checks that all required sections exist and conform to the template.
"""

import sys
import re
import json
from pathlib import Path


REQUIRED_SECTIONS_BY_DEPTH = {
    "brief": [
        "执行摘要",
        "Executive Summary",
        "SCQ",
        "情境",
        "冲突",
        "疑问",
        "核心论证",
        "Core Arguments",
        "风险矩阵",
        "Risk Matrix",
        "数据溯源",
        "Data Provenance",
    ],
    "standard": [
        "执行摘要",
        "Executive Summary",
        "SCQ",
        "情境",
        "冲突",
        "疑问",
        "行业全景",
        "Industry Landscape",
        "核心论证",
        "Core Arguments",
        "战略路线图",
        "Strategic Roadmap",
        "风险矩阵",
        "Risk Matrix",
        "数据溯源",
        "Data Provenance",
    ],
    "comprehensive": [
        "执行摘要",
        "Executive Summary",
        "SCQ",
        "情境",
        "冲突",
        "疑问",
        "行业全景",
        "Industry Landscape",
        "核心论证",
        "Core Arguments",
        "战略路线图",
        "Strategic Roadmap",
        "风险矩阵",
        "Risk Matrix",
        "数据溯源",
        "Data Provenance",
        "附录",
        "Appendix",
    ],
}

OPERATING_VERDICTS = [
    "VERIFIED_OPERATING_FACT",
    "HIGH_CONFIDENCE_INFERENCE",
    "EXPLAINABLE_ANOMALY",
    "SUSPENDED_JUDGMENT",
    "UNVERIFIED_NARRATIVE",
]

OPERATING_TRACE_TERMS = [
    "工商",
    "许可",
    "参保",
    "招聘",
    "地图",
    "POI",
    "小程序",
    "LBS",
    "招投标",
    "用户反馈",
    "员工反馈",
    "加盟商",
    "operating trace",
    "business license",
    "hiring",
    "store locator",
]

RESEARCHER_REQUIRED_FILES = [
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
]


def validate_report(
    report_path: str, depth: str = "standard", vertical: str | None = None
) -> dict:
    """Validate a report file and return results."""
    results = {
        "valid": True,
        "errors": [],
        "warnings": [],
        "stats": {},
    }

    path = Path(report_path)
    if not path.exists():
        results["valid"] = False
        results["errors"].append(f"Report file not found: {report_path}")
        return results

    content = path.read_text(encoding="utf-8")
    lines = content.split("\n")

    # Basic stats
    word_count = len(content)
    heading_count = len([l for l in lines if l.startswith("#")])
    results["stats"]["character_count"] = word_count
    results["stats"]["heading_count"] = heading_count

    # Check required sections
    required = REQUIRED_SECTIONS_BY_DEPTH.get(depth, REQUIRED_SECTIONS_BY_DEPTH["standard"])
    found_sections = set()
    for section_name in required:
        # Check if section name appears in any heading
        pattern = re.compile(rf"^#{{1,3}}\s+.*{re.escape(section_name)}", re.MULTILINE | re.IGNORECASE)
        if pattern.search(content):
            found_sections.add(section_name)

    # Group by language pairs - a section passes if either language variant is found
    section_pairs = [
        ("执行摘要", "Executive Summary"),
        ("核心论证", "Core Arguments"),
        ("风险矩阵", "Risk Matrix"),
        ("数据溯源", "Data Provenance"),
        ("行业全景", "Industry Landscape"),
        ("战略路线图", "Strategic Roadmap"),
        ("附录", "Appendix"),
    ]

    for cn, en in section_pairs:
        if cn in required or en in required:
            if cn not in found_sections and en not in found_sections:
                results["errors"].append(f"Missing required section: {cn} / {en}")
                results["valid"] = False

    # Check for SCQ structure (either language)
    has_scq = any(
        re.search(rf"(情境|Situation|冲突|Complication|疑问|Question)", content, re.IGNORECASE)
        for _ in [1]
    )
    if not has_scq and depth != "brief":
        results["warnings"].append("SCQ structure not clearly identified in the report")

    # Check for data citations
    citation_pattern = re.compile(r"\[(\d+)\]")
    citations = citation_pattern.findall(content)
    results["stats"]["citation_count"] = len(set(citations))
    if len(set(citations)) < 3:
        results["warnings"].append(
            f"Only {len(set(citations))} unique data citations found. "
            "Expected at least 3 for a credible report."
        )

    # Check for confidence scores
    confidence_pattern = re.compile(r"置信度[评分]*[:：]\s*(\d+)", re.IGNORECASE)
    confidence_en = re.compile(r"confidence[:\s]*(\d+)", re.IGNORECASE)
    conf_matches = confidence_pattern.findall(content) + confidence_en.findall(content)
    results["stats"]["confidence_mentions"] = len(conf_matches)
    if not conf_matches:
        results["warnings"].append("No confidence score found in the report")

    # Check for anti-patterns
    antipatterns = [
        (r"综合来看.{0,10}(机会与风险并存|利弊共存)", "Anti-pattern detected: wishy-washy synthesis"),
        (r"既有.{0,10}也有.{0,10}(因素|方面)", "Anti-pattern detected: non-committal conclusion"),
    ]
    for pattern, msg in antipatterns:
        if re.search(pattern, content):
            results["warnings"].append(msg)

    # Check depth-appropriate length
    depth_ranges = {
        "brief": (1500, 5000),
        "standard": (4000, 15000),
        "comprehensive": (8000, 30000),
    }
    min_chars, max_chars = depth_ranges.get(depth, (4000, 15000))
    if word_count < min_chars:
        results["warnings"].append(
            f"Report is {word_count} chars, below minimum {min_chars} for depth '{depth}'"
        )
    if word_count > max_chars:
        results["warnings"].append(
            f"Report is {word_count} chars, above maximum {max_chars} for depth '{depth}'"
        )

    # Check for DATA_INSUFFICIENT markers
    insufficient = re.findall(r"DATA_INSUFFICIENT", content)
    if insufficient:
        results["stats"]["data_insufficient_count"] = len(insufficient)
        results["warnings"].append(
            f"Found {len(insufficient)} DATA_INSUFFICIENT markers — "
            "these sections need additional data"
        )

    # Check for suspended judgments
    suspended = re.findall(r"悬置判断|Explicit Suspension|SUSPENDED", content, re.IGNORECASE)
    results["stats"]["suspended_judgments"] = len(suspended)

    if vertical in {"restaurant-retail-supply-chain", "rrsc"}:
        has_operating_section = re.search(
            r"^#{1,3}\s+.*(实体经营版图|Operating Footprint|Evidence Chain)",
            content,
            re.MULTILINE | re.IGNORECASE,
        )
        trace_hits = [term for term in OPERATING_TRACE_TERMS if term.lower() in content.lower()]
        verdict_count = sum(1 for verdict in OPERATING_VERDICTS if verdict in content)
        has_triangulation = re.search(r"三角验证|triangulation", content, re.IGNORECASE)
        has_core_arguments_section = re.search(
            r"^#{1,3}\s+.*(核心论证|Core Arguments)",
            content,
            re.MULTILINE | re.IGNORECASE,
        )
        has_summary_section = re.search(
            r"^#{1,3}\s+.*(执行摘要|Executive Summary)",
            content,
            re.MULTILINE | re.IGNORECASE,
        )
        has_brief_operating_compression = (
            (has_core_arguments_section or has_summary_section)
            and bool(has_triangulation or verdict_count > 0)
            and len(set(trace_hits)) >= 3
        )

        if depth in {"standard", "comprehensive"} and not has_operating_section:
            results["errors"].append(
                "Missing operating footprint section for restaurant/retail/supply-chain mode"
            )
            results["valid"] = False
        if depth == "brief" and not has_operating_section and not has_brief_operating_compression:
            results["warnings"].append(
                "No standalone operating footprint section found in brief vertical mode; "
                "expected compressed operating-footprint evidence in summary or core arguments"
            )

        if not has_triangulation:
            results["warnings"].append(
                "No triangulation matrix or triangulation discussion found"
            )

        results["stats"]["operating_verdict_count"] = verdict_count
        if verdict_count == 0:
            results["warnings"].append(
                "No operating verdict labels found; expected at least one entity-map judgment"
            )

        results["stats"]["operating_trace_term_count"] = len(set(trace_hits))
        if len(set(trace_hits)) < 3:
            results["warnings"].append(
                "Few operating trace terms found; report may rely too heavily on narrative sources"
            )

    return results


def validate_researcher_workspace(workspace_dir: str) -> dict:
    """Validate required researcher artifacts in a workspace directory."""
    results = {
        "valid": True,
        "errors": [],
        "warnings": [],
        "stats": {},
    }

    root = Path(workspace_dir)
    if not root.exists():
        results["valid"] = False
        results["errors"].append(f"Researcher workspace not found: {workspace_dir}")
        return results
    if not root.is_dir():
        results["valid"] = False
        results["errors"].append(f"Researcher workspace is not a directory: {workspace_dir}")
        return results

    present = []
    for name in RESEARCHER_REQUIRED_FILES:
        artifact = root / name
        if artifact.exists():
            present.append(name)
        else:
            results["valid"] = False
            results["errors"].append(f"Missing researcher artifact: {name}")

    results["stats"]["required_artifact_count"] = len(RESEARCHER_REQUIRED_FILES)
    results["stats"]["present_artifact_count"] = len(present)

    ledger_path = root / "evidence_ledger.json"
    confidence_path = root / "confidence_report.json"
    if ledger_path.exists():
        try:
            ledger = json.loads(ledger_path.read_text(encoding="utf-8"))
            raw_ledger = json.dumps(ledger, ensure_ascii=False)
            raw_confidence = ""
            if confidence_path.exists():
                raw_confidence = confidence_path.read_text(encoding="utf-8")
            if "retrieval_result_only" in raw_ledger and '"rating": "high"' in raw_confidence:
                results["warnings"].append(
                    "High confidence appears near retrieval-only evidence; inspect ledger"
                )
        except json.JSONDecodeError as exc:
            results["valid"] = False
            results["errors"].append(f"Invalid evidence_ledger.json: {exc}")

    return results


def main():
    if len(sys.argv) < 2:
        print(
            "Usage: python validate_report.py <report_path> "
            "[--depth brief|standard|comprehensive] "
            "[--vertical restaurant-retail-supply-chain|rrsc] "
            "[--researcher-workspace <workspace_dir>]"
        )
        sys.exit(1)

    if "--researcher-workspace" in sys.argv:
        workspace_idx = sys.argv.index("--researcher-workspace")
        if workspace_idx + 1 >= len(sys.argv):
            print("Missing value for --researcher-workspace")
            sys.exit(1)
        workspace_results = validate_researcher_workspace(sys.argv[workspace_idx + 1])
        print(json.dumps(workspace_results, ensure_ascii=False, indent=2))
        if len(sys.argv) <= 3:
            sys.exit(0 if workspace_results["valid"] else 1)
        if not workspace_results["valid"]:
            sys.exit(1)

    report_path = sys.argv[1]
    depth = "standard"
    vertical = None

    if "--depth" in sys.argv:
        depth_idx = sys.argv.index("--depth")
        if depth_idx + 1 < len(sys.argv):
            depth = sys.argv[depth_idx + 1]

    if "--vertical" in sys.argv:
        vertical_idx = sys.argv.index("--vertical")
        if vertical_idx + 1 < len(sys.argv):
            vertical = sys.argv[vertical_idx + 1]

    results = validate_report(report_path, depth, vertical)

    print(json.dumps(results, ensure_ascii=False, indent=2))

    if not results["valid"]:
        print(f"\n❌ Validation FAILED with {len(results['errors'])} error(s)")
        for err in results["errors"]:
            print(f"  ERROR: {err}")
    else:
        print(f"\n✅ Validation PASSED")

    if results["warnings"]:
        print(f"\n⚠️  {len(results['warnings'])} warning(s):")
        for warn in results["warnings"]:
            print(f"  WARN: {warn}")

    print(f"\n📊 Stats: {json.dumps(results['stats'], ensure_ascii=False)}")

    sys.exit(0 if results["valid"] else 1)


if __name__ == "__main__":
    main()
