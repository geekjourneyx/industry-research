# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A multi-agent adversarial industry research engine (Claude Code skill). Takes a fuzzy industry/sector/business-opportunity request, runs it through a structured "hypothesis → adversarial challenge → synthesis" pipeline, and outputs a research report with source citations and confidence scores. Has a specialized vertical for restaurant/retail/supply-chain research that uses operating-trace verification instead of narrative-based analysis.

## Architecture

The engine is defined entirely in `SKILL.md` — it is the orchestrator that coordinates four sub-agents. There is no application code; the system runs as a Claude Code skill invoked via `/industry-research`.

**Agent roles** (all in `agents/`):

| Agent | File | Role |
|-------|------|------|
| Engagement Manager | `agents/engagement-manager.md` | Structures the research question into a ghost deck (analysis skeleton) and optional entity evidence plan |
| Blue Team | `agents/blue-team.md` | Bull-case analyst — builds evidence-backed optimistic thesis |
| Red Team | `agents/red-team.md` | Bear-case analyst — pre-mortem + adversarial challenge |
| Chief Arbitrator | `agents/chief-arbitrator.md` | Hegelian synthesis — resolves red/blue conflict into a final report with verdicts |

**Execution flow** (stages in SKILL.md):

1. **Domain Grounding** (Phase 1): Industry taxonomy mapping → web research → context dictionary → user confirmation
2. **Business Physics Modeling** (Step 2.0, restaurant/retail/supply-chain only): Entity evidence plan (`entity_evidence_plan.json`)
3. **Ghost Deck** (Step 2.1): MECE chapter structure with falsifiable action titles
4. **Round 1** (Step 2.2): Red and blue agents run **in parallel** — independent analysis
5. **Round 2** (Step 2.3): Cross-rebuttal — each side sees the other's R1 output (skipped in `brief` mode unless conflicts are severe)
6. **Arbitration** (Step 2.4): Chief arbitrator produces `final_report.md` + `report_metadata.json`
7. **Validation** (Step 2.5): Automated structure check via `scripts/validate_report.py`

**Key pipeline contract**: Agents receive absolute paths to `{workspace_dir}` and write JSON outputs there. The orchestrator reads agent outputs and passes them to the next stage. All intermediate artifacts persist to disk.

## Report Depths

| Depth | Target length | Round 2 | Use case |
|-------|--------------|---------|----------|
| `brief` | ~3,000 chars | Skipped by default | Quick decisions, internal alignment |
| `standard` | ~8,000 chars | Yes | Investment decisions, strategic planning |
| `comprehensive` | ~15,000 chars | Yes | Due diligence, board materials |

## Validation & Testing

**Report validation** (run after report generation):
```bash
# General industry
python3 scripts/validate_report.py <workspace_dir>/final_report.md --depth standard

# Restaurant/retail/supply-chain vertical
python3 scripts/validate_report.py <workspace_dir>/final_report.md --depth standard --vertical restaurant-retail-supply-chain
```

The validator checks: required sections present, citation count, confidence scores, anti-patterns (wishy-washy synthesis), character length bounds, and operating-footprint evidence for the RRSC vertical.

**Smoke tests** (require pre-built workspace directories):
```bash
scripts/run_rrsc_smoke_validation.sh       # Full pipeline: validates all 10 artifact files + report structure
scripts/run_maiji_interview_prep_validation.sh  # Interview-prep variant: validates entity evidence plan + interview guide
```

**Evaluation prompts** are in `evals/evals.json` — 7 test scenarios covering different depths, languages, and verticals.

## Key References

Files in `references/` are read by agents at specific pipeline stages:

- `industry-taxonomy.md` — National economic industry classification codes (Phase 1, Step 1.1)
- `analytical-frameworks.md` — MECE, PESTLE, Porter's Five Forces, pre-mortem, MAMV voting (passed to all agents)
- `restaurant-retail-supply-chain-physics.md` — Physical constraints model for restaurant/retail/supply-chain (Step 2.0)
- `evidence-triangulation-playbook.md` — Evidence family taxonomy, triangulation rules, anomaly attribution (Steps 2.0 and 2.4)
- `report-template.md` — Final report structure template (Step 2.4)

## Degradation Protocol

When an agent fails or times out, the system degrades gracefully rather than crashing:
- Mark `execution_mode = degraded` in metadata
- Append to `degradation_tags` (e.g., `ghost_deck_generation`, `round1_missing_side`)
- Main agent generates a minimum viable version of the missing artifact
- Report clearly notes which sections are degraded

## Language Handling

Output language follows user input language. Chinese input → Chinese report; English input → English report. All agents work in the same language as the final report.
