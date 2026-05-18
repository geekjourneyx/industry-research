# Researcher Agent Design

## Current Status

This design has been implemented as the first `researcher` version in this repository.

Current implementation notes:

- The CLI lives in `researcher/`.
- `researcher retrieve` currently exposes Bocha direct search.
- `researcher answer volcengine` exposes Volcengine Ark model-mediated web search.
- Multi-provider retrieval exists as an internal provider layer; the public `retrieve --providers bocha,volcengine` shape remains a future CLI expansion.
- `researcher run` creates the required workspace artifacts for first-pass research, trace planning, evidence ledger, disconfirmation log, confidence report, and final report.
- Browser verification is represented in schemas and review rules, but not yet automated inside the CLI.

## Summary

This design redefines the project around a higher-level goal:

```text
Build a research agent that can investigate claims like a human domain expert.
```

The CLI should be named:

```text
researcher
```

`researcher` is not just a search gateway. It is a research execution engine. It should be able to:

1. understand a research question,
2. decompose it into verifiable claims,
3. infer what real-world traces should exist if each claim is true,
4. retrieve leads through direct search, model-mediated search, agent-native search, and later browser verification,
5. maintain an evidence ledger,
6. search for disconfirming evidence,
7. score confidence,
8. generate a report with explicit uncertainty.

The existing `industry-research` skill should be redesigned. It should no longer try to orchestrate every research step directly inside `SKILL.md`. Instead, it should become a thin domain entrypoint that configures and calls `researcher`, then enforces industry-specific standards around trace reasoning, evidence quality, and report delivery.

## First Principles

Human experts do not simply search keywords.

When asked "Is Luckin's 2026 store-count target credible?", a human expert reasons:

```text
If store expansion is real, the world must change.
Those changes should leave traces.
Hiring traces may appear before new-store revenue.
Map, mini-program, delivery-platform, and review traces appear when stores become customer-addressable.
Legal, license, branch, or franchise traces may appear when operating entities expand.
Management claims should be checked against filings, public-account posts, interviews, and authoritative media.
Missing expected traces may be evidence against the claim.
```

The system should therefore be built around this chain:

```text
question
-> decision context
-> claims
-> mechanisms
-> expected traces
-> retrieval plan
-> retrieved leads
-> source verification
-> evidence ledger
-> disconfirmation attempts
-> confidence judgment
-> report
```

The key design principle:

```text
Fix the quality gates, not the exploration path.
```

The agent should not be forced into a fixed platform checklist. It should be forced to explain why the sources it chose are appropriate for the claim.

## Why Current SKILL.md Needs Redesign

The current `SKILL.md` is too much of an orchestrator.

It currently tries to specify:

- domain grounding,
- business-physics modeling,
- ghost deck generation,
- red-blue parallel analysis,
- cross rebuttal,
- arbitration,
- validation,
- degradation handling,
- report delivery.

That structure is useful, but it causes three problems:

1. **It mixes reasoning policy with execution mechanics.** The skill describes both how to think and how to run the pipeline. This makes it hard to add real tooling without making `SKILL.md` larger and more brittle.
2. **It relies too heavily on prompt compliance.** Evidence ledgers, source verification, and disconfirmation attempts are described as instructions, but they are not first-class execution artifacts controlled by a dedicated research engine.
3. **It makes provider integration awkward.** Bocha direct search, Volcengine model-mediated search, and future browser verification do not belong as ad hoc scripts inside the industry skill.

The redesign should split responsibilities:

```text
researcher CLI
  Executes the research workflow and owns reusable research primitives.

industry-research skill
  Defines domain policy, required report shape, and when/how to call researcher.
```

## Product Boundary

`researcher` is the research agent.

It owns:

- question parsing,
- claim decomposition,
- trace reasoning,
- retrieval provider orchestration,
- evidence ledger,
- disconfirmation search,
- confidence scoring,
- report artifact generation,
- machine-readable metadata.

`industry-research` is a domain skill.

It owns:

- industry-specific trigger rules,
- restaurant/retail/supply-chain emphasis,
- report-depth choices,
- user-facing workflow,
- domain references and evaluation cases,
- final communication style.

Provider packages are infrastructure.

They own:

- Bocha API mapping,
- Volcengine Ark Responses API mapping,
- normalized retrieval outputs,
- provider errors,
- capability declarations.

## Architecture

```text
researcher/
  cmd/researcher/
  internal/project/
  internal/question/
  internal/claims/
  internal/trace/
  internal/retrieval/
  internal/provider/bocha/
  internal/provider/volcengine/
  internal/provider/multi/
  internal/ledger/
  internal/disconfirm/
  internal/confidence/
  internal/report/
  internal/output/
  internal/config/
  internal/errors/
  Makefile
  VERSION
  README.md

industry-research/
  SKILL.md
  agents/
  references/
    chain-brand-trace-reasoning.md
    evidence-ledger-schema.md
    report-template.md
  scripts/
  evals/
```

`researcher` can be used by any skill. `industry-research` is the first major consumer.

## CLI Design

Binary:

```bash
researcher
```

Top-level commands:

```bash
researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --depth standard
researcher plan "某茶饮品牌宣称 5000 家门店是否可信？" --domain chain-brand --json
researcher retrieve "瑞幸 2026 门店数 招聘 扩张" --providers bocha,volcengine --json
researcher evidence "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --json
researcher validate ./researcher-workspace/luckin-2026-store-count
researcher capabilities --json
researcher version
researcher help
```

Command responsibilities:

```text
run
  Full research workflow: question -> claims -> traces -> retrieval -> ledger -> confidence -> report.

plan
  Produces a research plan without executing retrieval.

retrieve
  Lower-level retrieval command for direct use and debugging.

evidence
  Produces or updates an evidence ledger without writing a full report.

validate
  Validates a workspace, evidence ledger, metadata, and report.

capabilities
  Lists provider capabilities and supported domains.
```

## Configuration

Default config file lookup should follow Unix/XDG conventions:

```text
1. --config <path>
2. RESEARCHER_CONFIG
3. $XDG_CONFIG_HOME/researcher/config.yaml
4. ~/.config/researcher/config.yaml
```

The default generated config should be `~/.config/researcher/config.yaml`. Do not add a second home-directory config location.

Value precedence:

```text
1. command flags
2. environment variables
3. config file
4. built-in defaults
```

Environment variables should override config file secrets:

```text
BOCHA_API_KEY
ARK_API_KEY
RESEARCHER_CONFIG
```

Example config:

```yaml
providers:
  bocha:
    api_key: ""
    endpoint: "https://api.bochaai.com/v1/web-search"
  volcengine:
    api_key: ""
    endpoint: "https://ark.cn-beijing.volces.com/api/v3/responses"
    model: "doubao-seed-2-0-lite-260215"
defaults:
  providers: ["bocha", "volcengine"]
  depth: "standard"
  workspace_root: "researcher-workspace"
```

## Unix Philosophy

Although `researcher` is higher-level than a pure retrieval tool, each command should remain composable:

- stdout emits machine-readable output when `--json` is set.
- stderr emits diagnostics.
- exit codes signal failure category.
- no interactive prompts in CLI mode.
- each run writes artifacts to a workspace.
- every major artifact is a file that another tool can inspect.

The command should not hide uncertainty. If evidence is weak, the output should say so explicitly and degrade confidence.

## Workspace Artifacts

Default workspace:

```text
researcher-workspace/{topic_slug}/
```

Required artifacts for `run`:

```text
question.json
research_plan.json
claim_graph.json
trace_plan.json
retrieval_log.json
evidence_ledger.json
disconfirmation_log.json
confidence_report.json
final_report.md
report_metadata.json
```

For chain-brand research, `trace_plan.json` replaces the older idea of a purely prompt-written ghost deck as the primary working object.

## Core Data Flow

```text
question.json
  User input, domain, depth, geography, time range, decision context.

research_plan.json
  What will be investigated, what will not, and why.

claim_graph.json
  Verifiable claims and dependencies.

trace_plan.json
  For each claim: mechanism, expected traces, source families, retrieval purposes, disconfirming traces.

retrieval_log.json
  Every retrieval action, provider, query, parameters, status, and raw output path.

evidence_ledger.json
  Evidence items mapped to claims and traces.

disconfirmation_log.json
  Attempts to falsify or weaken claims.

confidence_report.json
  Confidence decisions and limiting factors.

final_report.md
  Human-readable report.

report_metadata.json
  Machine-readable report summary.
```

## Retrieval Providers

`researcher` includes retrieval capabilities, but retrieval is only one module.

Provider types:

```text
direct_search
  Provider returns search results directly.
  First provider: Bocha.

model_answer_search
  Provider uses a model that can search and answer.
  First provider: Volcengine Ark Responses API with web_search.

agent_native_search
  Host agent search capability. Represented in artifacts but not owned by Go CLI.

browser_verification
  Future adapter for interactive pages.

knowledge_retrieval
  Future adapter for local or enterprise knowledge bases.
```

Provider output envelope:

```json
{
  "provider": "bocha",
  "provider_type": "direct_search",
  "mode": "search|answer|retrieve",
  "query": "瑞幸 2026 门店数 招聘 扩张",
  "retrieved_at": "2026-05-17T10:00:00+08:00",
  "request": {},
  "retrieval_calls": [],
  "items": [],
  "answer": {
    "text": "",
    "citations": []
  },
  "usage": {},
  "errors": []
}
```

Rules:

- Search results are leads.
- Model answers are leads.
- Citations are leads.
- Only opened, verified, or cross-validated sources become evidence.

## Bocha Provider

Bocha is a `direct_search` provider.

Endpoint:

```text
POST https://api.bocha.cn/v1/web-search
Authorization: Bearer $BOCHA_API_KEY
Content-Type: application/json
```

Request mapping:

```text
--count       -> count
--freshness   -> freshness: noLimit|oneDay|oneWeek|oneMonth|oneYear
--summary     -> summary
--include     -> include
--exclude     -> exclude
```

Response mapping:

```text
webPages.value[].name             -> items[].title
webPages.value[].url              -> items[].url
webPages.value[].displayUrl       -> items[].display_url
webPages.value[].snippet          -> items[].snippet
webPages.value[].summary          -> items[].summary
webPages.value[].siteName         -> items[].site_name
webPages.value[].siteIcon         -> items[].site_icon
webPages.value[].datePublished    -> items[].published_at
webPages.value[].dateLastCrawled  -> items[].last_crawled_at
webPages.value[].cachedPageUrl    -> provider_metadata.cached_page_url
webPages.value[].language         -> items[].language
images.value[]                    -> items[] with content_type image
```

Date rule:

- Prefer `datePublished`.
- Bocha docs say `dateLastCrawled` values such as `2025-02-23T08:18:30Z` actually represent UTC+8 Beijing time.
- Normalize that documented shape to `2025-02-23T08:18:30+08:00` when `datePublished` is absent.
- Preserve raw value.

Error mapping:

```text
400 missing query       -> invalid_argument
400 missing API key     -> missing_api_key
401 invalid API key     -> provider_auth_error
403 insufficient funds  -> provider_quota_exhausted
429 request limit       -> provider_rate_limited
500 provider exception  -> provider_unavailable
timeout                 -> provider_timeout
parse failure           -> provider_parse_error
```

## Volcengine Provider

Volcengine is a `model_answer_search` provider.

Endpoint:

```text
POST https://ark.cn-beijing.volces.com/api/v3/responses
Authorization: Bearer $ARK_API_KEY
Content-Type: application/json
```

Request mapping:

```text
query              -> input user message
--model            -> Responses model, default doubao-seed-2-0-lite-260215
--limit            -> tools[0].limit
--max-keyword      -> tools[0].max_keyword
--max-tool-calls   -> max_tool_calls
--sources          -> tools[0].sources
--location-*       -> tools[0].user_location
```

Rules:

- Use non-streaming mode first for deterministic JSON normalization.
- Do not send `caching`; docs say it currently returns `400`.
- Extract `web_search_call` items and actual `action.query` values.
- Extract `message.content[0].annotations` into citations and lead items.
- Preserve final answer text as `answer.text`.
- Preserve `usage.tool_usage` and `usage.tool_usage_details`.
- If no web search call occurred, return `no_retrieval_triggered`.

Important distinction:

```text
Volcengine answer text is not evidence.
Volcengine citations are leads.
Only opened or cross-validated cited URLs can become evidence.
```

## Error Contract

```json
{
  "code": "missing_api_key|invalid_argument|provider_http_error|provider_auth_error|provider_quota_exhausted|provider_rate_limited|provider_timeout|provider_unavailable|provider_parse_error|no_retrieval_triggered|partial_failure",
  "message": "Human-readable error for agents and logs.",
  "provider_status": 429,
  "provider_code": "429",
  "provider_log_id": "c66aac17eab1bb7e",
  "retryable": true,
  "agent_action": "Wait and retry with lower count, or use another provider.",
  "raw_error_path": "workspace/retrieval/raw/bocha-error-001.json"
}
```

Every error must tell the agent what to do next.

## Provider Capabilities

`researcher capabilities --json` returns machine-readable capability declarations.

Example:

```json
{
  "provider": "bocha",
  "provider_type": "direct_search",
  "modes": ["search", "retrieve"],
  "supports_freshness": true,
  "supports_include_domains": true,
  "supports_exclude_domains": true,
  "supports_summary": true,
  "supports_location": false,
  "supports_sources": false,
  "supports_images": true,
  "supports_model_choice": false,
  "result_kinds": ["web_page", "image"]
}
```

```json
{
  "provider": "volcengine",
  "provider_type": "model_answer_search",
  "modes": ["answer", "retrieve"],
  "supports_freshness": false,
  "supports_include_domains": false,
  "supports_exclude_domains": false,
  "supports_summary": false,
  "supports_location": true,
  "supports_sources": true,
  "supports_images": false,
  "supports_model_choice": true,
  "result_kinds": ["annotation_url", "answer_text", "retrieval_call"]
}
```

## Trace Reasoning Module

Trace reasoning is the heart of `researcher`.

For each claim it must produce:

```json
{
  "claim_id": "claim_store_count_2026",
  "claim": "瑞幸 2026 年门店数继续增长具备经营支撑",
  "mechanism": "门店增长需要选址、招聘、供应链、数字入口和用户需求共同支撑",
  "expected_traces": [
    {
      "trace_type": "people_org",
      "trace": "新开城或加密城市出现店长、咖啡师、区域运营、拓展岗位",
      "why_expected": "门店扩张前后必须补充门店和区域运营人员"
    }
  ],
  "source_families": ["recruiting", "map_poi", "platform_frontend", "company_disclosure", "media_interview"],
  "disconfirming_traces": [
    "声称覆盖城市但无门店 POI",
    "无招聘或仅总部招聘",
    "小程序不可下单",
    "只有通稿没有独立经营痕迹"
  ]
}
```

This module should teach patterns, not fixed checklists.

Bad:

```text
Always check BOSS, Qichacha, WeChat, maps, and news.
```

Good:

```text
Check recruitment because expansion requires labor.
Check registry because new operators or franchise entities may appear.
Check maps and mini-programs because real stores must become customer-addressable.
Skip supplier tenders if the claim is store count rather than supply-chain maturity.
```

## Evidence Ledger

Every evidence item records:

```json
{
  "evidence_id": "ev_001",
  "claim_id": "claim_store_count_2026",
  "source_url": "https://example.com",
  "source_title": "示例来源",
  "source_type": "company_disclosure|official_registry|recruiting|map_poi|platform_frontend|media|social|legal|tender|ugc|retrieval_result_only",
  "evidence_family": "capital_legal|people_org|physical_fulfillment|digital_frontend|terminal_feedback|management_narrative",
  "origin_provider": "bocha|volcengine|agent_websearch|browser|manual",
  "origin_retrieval_id": "retrieval_001",
  "accessed_at": "2026-05-17T10:10:00+08:00",
  "verification_status": "retrieval_result_only|source_opened|browser_verified|cross_validated|not_accessible",
  "independence_note": "不是公司通稿转载，独立于财报口径",
  "supports_or_challenges": "supports|challenges|mixed|lead_only",
  "summary": "这条证据说明了什么",
  "requires_browser_verification": false,
  "browser_verification_reason": ""
}
```

Ledger rules:

- Retrieval-only items cannot support high confidence.
- Model answers cannot be final evidence.
- Search summaries cannot replace source verification.
- Provider failures must be recorded.
- High confidence requires independent evidence families and disconfirmation attempts.
- Browser-required evidence must be marked.

## Confidence Rules

Confidence should be computed from evidence quality, not prose quality.

High confidence:

- at least three independent evidence families,
- source pages opened or otherwise verified,
- disconfirmation attempts performed,
- no unresolved core contradiction.

Medium confidence:

- at least two independent evidence families,
- plausible mechanism,
- some verification gaps remain.

Low confidence:

- single evidence family,
- mostly leads or summaries,
- important platform verification missing.

Suspended:

- evidence conflict cannot be explained,
- expected traces are missing,
- provider coverage is too weak for the claim.

Unverified narrative:

- only company claims, media reposts, model answer text, or search summaries.

## Browser Automation Boundary

Browser automation is not first-version scope, but the schema must preserve handoff fields:

```json
{
  "requires_browser_verification": true,
  "browser_verification_reason": "需要切换城市查看小程序/地图/外卖平台可下单状态",
  "suggested_browser_steps": [
    "打开品牌小程序或门店页",
    "切换到目标城市",
    "记录是否可选门店、是否可下单、SKU 是否可售"
  ]
}
```

The agent must not pretend interactive evidence was verified by retrieval.

## Research Skill Redesign

`industry-research/SKILL.md` should become thinner.

Old role:

```text
Do the whole research pipeline through prompts and sub-agents.
```

New role:

```text
Parse the user's industry research need.
Choose domain, depth, and report language.
Call researcher with the right arguments.
Inspect researcher artifacts.
Apply domain-specific quality gates.
Return the user-facing result.
```

Suggested skill flow:

```text
1. Parse user request.
2. Decide domain: general | chain-brand | restaurant-retail-supply-chain.
3. Decide depth: brief | standard | comprehensive.
4. Call researcher run.
5. Validate artifacts.
6. If confidence is low, report why and what additional verification is needed.
7. Deliver final report summary.
```

`SKILL.md` should still preserve domain-specific requirements:

- restaurant/retail/supply-chain must use trace reasoning,
- operating traces outrank media narratives,
- evidence ledger is mandatory,
- unverified claims must be downgraded,
- browser-required evidence must be clearly marked,
- final language follows user language.

## Agent Role Redesign

The current four agent files can remain as domain personas, but they should stop acting as the only execution engine.

New use:

```text
Engagement Manager
  Reviews researcher claim graph and trace plan.

Blue Team
  Reviews whether supporting evidence is strong enough.

Red Team
  Reviews missing traces, contradictions, and disconfirmation attempts.

Chief Arbitrator
  Reviews confidence_report.json and final_report.md.
```

This makes agents reviewers of artifacts, not free-form report generators.

## Validation

`researcher validate` should check:

- required workspace files exist,
- every claim has expected traces,
- every retrieval action has purpose, provider, status, and timestamp,
- every evidence item maps to a claim,
- high-confidence claims meet evidence-family threshold,
- model answer text is not cited as evidence,
- retrieval-only items are not high-confidence evidence,
- disconfirmation attempts exist,
- browser-required checks are marked,
- final report confidence matches `confidence_report.json`.

## Evaluation Cases

Add or revise evals around:

1. Listed chain brand expansion: "瑞幸咖啡 2026 年门店数目标是否可信？"
2. Non-listed franchise tea brand: "某区域茶饮品牌宣称 5000 家门店、加盟商盈利稳定是否可信？"
3. Supply-chain maturity: "某餐饮品牌宣称全国统一供应链成熟是否可信？"
4. City coverage: "某生鲜零售品牌宣称全国冷链覆盖是否真实？"
5. Management interview verification: "把管理层访谈中的 GMV、单店 UE、同店增长、供应链统一供货比例转成可验证命题。"

Each eval should expect:

- claim decomposition,
- trace reasoning,
- retrieval actions with purpose,
- evidence ledger,
- disconfirmation attempts,
- confidence downgrade where sources are weak,
- clear separation between leads and evidence.

## Go Project Practices

Follow the Makefile style from `~/Workspace/go/md2wechat-skill/Makefile`.

Required targets:

```text
all
build
fast
release
clean
test
fmt
vet
install
deps
help
```

Build rules:

- Use `VERSION`.
- Inject version with `-ldflags`.
- Use `go build -trimpath`.
- Build current platform to `./researcher`.
- Build release binaries under `bin/`.
- Support Linux amd64, Linux arm64, macOS amd64, macOS arm64, and Windows amd64.
- Prefer standard library HTTP clients.
- Add dependencies only when they remove real complexity.

## External References

- `docs/bocha_websearch.md` documents Bocha Web Search request parameters, response fields, endpoint `https://api.bocha.cn/v1/web-search`, date handling caveat, and error codes.
- `docs/volcengine_websearch.md` documents Volcengine Ark Responses API Web Search tool behavior, tool parameters, sources, location, usage fields, and model-mediated search-call output.
- Bocha Open Platform documents Web Search API shape and endpoint examples: `https://open.bocha.cn/`
- OpenAI documents the Responses API `web_search` tool, citations, sources, domain filtering, and live access controls: `https://developers.openai.com/api/docs/guides/tools-web-search`
- Volcengine documents联网搜索 as a deep-research capability with real-time data access, search strategy planning, multi-source verification, and structured report output: `https://www.volcengine.com/docs/85637/1588465`
- Volcengine联网搜索 API reference page exists but currently requires JavaScript for full details in this environment: `https://www.volcengine.com/docs/87772/2272953`

## Acceptance Criteria

The design is successful when:

1. `researcher` can run as an independent Go CLI.
2. `industry-research` can call `researcher` instead of manually orchestrating all research steps.
3. Bocha and Volcengine are integrated as provider modules with different provider types.
4. Reports show claim-to-trace reasoning before conclusions.
5. Retrieval results are treated as leads until verified.
6. Confidence is based on evidence-family independence and disconfirmation attempts.
7. Agents remain free to choose sources but must explain why those sources fit the claim.
8. Provider errors are agent-actionable.
9. Current prompt-only red/blue/arbitrator flow becomes artifact review instead of the only execution engine.

## Non-Goals

The first version does not:

- implement browser automation,
- build a database-backed evidence store,
- add UI,
- scrape closed platforms,
- guarantee access to BOSS, Qichacha, WeChat, video accounts, maps, delivery platforms, or review platforms,
- replace host agent web search,
- let retrieval providers determine truth.

## Implementation Order

Recommended sequence:

1. Create standalone Go CLI skeleton and Makefile.
2. Define workspace artifact schemas.
3. Define claim, trace, retrieval, evidence, confidence, and report types.
4. Add mocked provider tests.
5. Add Bocha provider.
6. Add Volcengine provider.
7. Add retrieval command.
8. Add plan and evidence commands.
9. Add confidence rules and validation.
10. Add run command.
11. Redesign `industry-research/SKILL.md` to call `researcher`.
12. Update agents as artifact reviewers.
13. Update evals.

This order builds the research agent as a reusable engine and keeps the industry skill focused on domain policy.
