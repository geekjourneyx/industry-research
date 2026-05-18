# industry-research

`industry-research` is a Codex skill plus a Go CLI for adversarial industry research.

The current architecture has two layers:

1. `researcher` CLI creates a reproducible research workspace.
2. The Codex skill and sub-agents review the workspace, challenge weak evidence, and produce the final user-facing research answer.

The main design rule is simple: search results are leads, not evidence. A high-confidence conclusion needs independent evidence families and an explicit attempt to disprove it.

## Architecture

```text
User request
  -> SKILL.md routes the request and selects depth/domain
  -> researcher run creates a workspace
  -> researcher writes trace, evidence, disconfirmation, confidence, and report artifacts
  -> agents review the artifacts instead of inventing conclusions from prose
  -> final report explains confidence limits and unresolved claims
```

## Primary Commands

Build and inspect the CLI:

```bash
cd researcher
make build
./researcher version
./researcher capabilities --json
```

Run a first-pass research workspace:

```bash
./researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" \
  --domain chain-brand \
  --depth standard \
  --workspace-root ../industry-research-workspace \
  --json
```

Validate the generated workspace:

```bash
./researcher validate ../industry-research-workspace/<workspace>
python3 ../scripts/validate_report.py \
  --researcher-workspace ../industry-research-workspace/<workspace>
```

## Workspace Artifacts

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

For restaurant, retail, supply-chain, and chain-brand questions, the important files are `trace_plan.json`, `evidence_ledger.json`, `disconfirmation_log.json`, and `confidence_report.json`.

## Provider Integrations

Current provider support:

- Bocha web search: direct web search results.
- Volcengine Ark web search: model answer with web-search annotations.
- Agent web search and browser verification remain evidence-review capabilities used by the skill layer.

Configuration should use XDG paths, not `~/.researcher`:

1. `--config <path>`
2. `RESEARCHER_CONFIG`
3. `$XDG_CONFIG_HOME/researcher/config.yaml`
4. `~/.config/researcher/config.yaml`

Use environment variables for secrets:

- `BOCHA_API_KEY`
- `ARK_API_KEY`

## Documentation Map

- `SKILL.md`: Codex skill entry point and orchestration contract.
- `researcher/README.md`: CLI usage and configuration.
- `references/evidence-ledger-schema.md`: evidence schema and confidence rules.
- `references/chain-brand-trace-reasoning.md`: trace reasoning for chain brands.
- `docs/bocha_websearch.md`: Bocha API source notes.
- `docs/volcengine_websearch.md`: Volcengine Ark web-search source notes.
- `evals/evals.json`: evaluation prompts and expected artifacts.

## Validation

```bash
cd researcher
make fmt
make vet
make test
make build
```

If the sandbox blocks local `httptest` ports, rerun `make test` with permission to bind local test ports.
