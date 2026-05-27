# AGENTS.md

This file guides Codex work in this repository.

## What This Is

This repository contains a Codex skill and a Go CLI for industry research.

The system takes a fuzzy industry, sector, company, or business-opportunity question and turns it into a structured research workspace and report. The core quality standard is evidence discipline: search results and model answers are only leads until they are verified against sources, operating traces, or cross-source checks.

The specialized restaurant, retail, supply-chain, and chain-brand path uses operating-trace verification instead of narrative-only market analysis.

## Current Architecture

The current execution center is `researcher`, a Go CLI in `researcher/`.

`SKILL.md` is still the Codex entry point. It routes the request, chooses depth and domain, calls `researcher`, then uses the four agent role files as review and challenge layers.

The older pure prompt-only pipeline remains useful as a fallback and review model, but it is no longer the primary architecture.

## Main Components

| Component | Path | Role |
|---|---|---|
| Skill entry | `SKILL.md` | Codex-facing orchestration contract |
| CLI | `researcher/` | Builds workspaces, provider calls, trace planning, evidence ledger, confidence, validation |
| Engagement Manager | `agents/engagement-manager.md` | Reviews claim graph and trace plan |
| Blue Team | `agents/blue-team.md` | Reviews support evidence and upside thesis |
| Red Team | `agents/red-team.md` | Searches for missing traces, conflicts, and over-claimed evidence |
| Chief Arbitrator | `agents/chief-arbitrator.md` | Checks consistency among confidence, evidence, and final report |
| Validator | `scripts/validate_report.py` | Validates report shape and `researcher` workspace artifacts |

## Primary Flow

1. User asks for `/industry-research` or an equivalent industry research task.
2. `SKILL.md` determines language, depth, and domain.
3. The skill calls:

```bash
researcher run "<question>" --domain <domain> --depth <brief|standard|comprehensive> --workspace-root <workspace_root>
```

4. `researcher` creates a workspace and writes required artifacts.
5. Agents review the artifacts and challenge weak conclusions.
6. The final answer reports confidence limits, unsupported claims, and next verification steps.

## Required Workspace Artifacts

Every `researcher run` workspace must contain:

- `question.json`
- `research_plan.json`
- `claim_graph.json`
- `trace_plan.json`
- `retrieval_log.json`
- `evidence_ledger.json`
- `disconfirmation_log.json`
- `confidence_report.json`
- `final_report.md`
- `report_metadata.json`

For restaurant, retail, supply-chain, and chain-brand research, the most important quality files are `trace_plan.json`, `evidence_ledger.json`, `disconfirmation_log.json`, and `confidence_report.json`.

## Report Depths

| Depth | Target length | Use case |
|---|---:|---|
| `brief` | ~3,000 chars | Quick decisions, internal alignment |
| `standard` | ~8,000 chars | Investment decisions, strategy work |
| `comprehensive` | ~15,000 chars | Due diligence, board materials |

## Validation & Testing

CLI checks:

```bash
cd researcher
make fmt
make vet
make test
make build
```

Workspace validation:

```bash
researcher validate <workspace_dir>
python3 scripts/validate_report.py --researcher-workspace <workspace_dir>
```

Report validation:

```bash
python3 scripts/validate_report.py <workspace_dir>/final_report.md --depth standard
python3 scripts/validate_report.py <workspace_dir>/final_report.md --depth standard --vertical restaurant-retail-supply-chain
```

The validator checks required sections, citations, confidence mentions, weak synthesis patterns, length ranges, operating-footprint evidence, and required `researcher` artifacts.

## Release Protocol

Every release must complete all of the following before publishing:

1. Documentation check and calibration:
   - Review root `README.md` and `researcher/README.md`.
   - Confirm installation, usage, release notes, and command examples match the current code.
   - Keep the wording clear for humans and suitable for GitHub SEO/GEO discovery.
2. Changelog update:
   - Update `CHANGELOG.md` for the release.
   - The changelog entry must include the release version, date, user-facing changes, and any breaking changes or migration notes.
3. Binary build:
   - Compile `researcher` release binaries before creating the release.
   - At minimum cover macOS, Linux, and Windows for amd64; include arm64 where supported.
   - Generate checksums for release artifacts.
4. Git tag:
   - Create a version tag only after documentation, changelog, tests, and binary build checks pass.
   - The tag must match the `researcher/VERSION` value.
5. GitHub release:
   - Publish the release from the tag.
   - Attach compiled binaries and checksum files.
   - Do not publish a release from an unverified working tree.

Do not treat GitHub Releases as a manual upload step only. The preferred path is a GitHub Actions workflow that runs documentation checks, Go checks, binary compilation, checksum generation, and release publishing from a version tag.

## Provider Configuration

Do not use `~/.researcher/config.yaml`.

Config lookup order:

1. `--config <path>`
2. `RESEARCHER_CONFIG`
3. `$XDG_CONFIG_HOME/researcher/config.yaml`
4. `~/.config/researcher/config.yaml`

Use environment variables for secrets:

- `BOCHA_API_KEY`
- `ARK_API_KEY`

## Key References

- `references/industry-taxonomy.md`: industry classification codes.
- `references/analytical-frameworks.md`: MECE, PESTLE, Porter's Five Forces, pre-mortem, MAMV.
- `references/restaurant-retail-supply-chain-physics.md`: operating physics for restaurant, retail, and supply-chain questions.
- `references/evidence-triangulation-playbook.md`: evidence families, triangulation, anomaly attribution.
- `references/evidence-ledger-schema.md`: evidence ledger schema and confidence constraints.
- `references/chain-brand-trace-reasoning.md`: trace reasoning for chain brands.
- `references/report-template.md`: final report structure template.

## Degradation Protocol

When evidence is thin or an agent/tool fails:

- Mark `execution_mode = degraded` in metadata.
- Append a specific `degradation_tags` value.
- Downgrade unsupported claims instead of writing them as facts.
- Keep unresolved questions visible in the final report.

## Language Handling

Output language follows user input language. Chinese input produces Chinese reports; English input produces English reports. All agents should work in the same language as the final report.
