# Trace Reasoning and Research Gateway Design

## Summary

This design upgrades `industry-research` around a higher-level primitive: a reusable **research gateway**.

The gateway is not just a web-search wrapper. It is a Unix-style retrieval and evidence-ingestion tool that can call direct search APIs, model-mediated search systems, agent-native search, and future browser or knowledge-base retrievers, then normalize the outputs into a stable, auditable JSON envelope.

The goal is to support human-like research behavior:

```text
claim
-> real-world mechanism
-> expected traces
-> retrieval intent
-> provider selection
-> retrieved leads
-> source verification
-> evidence ledger
-> disconfirmation search
-> confidence judgment
```

The gateway does not decide truth. It retrieves, normalizes, preserves provenance, and reports provider limitations. `industry-research` remains responsible for trace reasoning, evidence judgment, confidence scoring, and report writing.

## Why Rename From Websearch

`websearch` is too narrow.

The system needs to support at least three current retrieval modes:

1. **Direct search**: Bocha returns structured search results from a search endpoint.
2. **Model-mediated search**: Volcengine Ark Responses API lets a model decide whether and how to call `web_search`, then returns answer text, search-call events, annotations, and usage.
3. **Agent-native search**: the host agent can use its own search tool, open pages, apply domain filters, and cite sources.

Future modes may include:

- browser verification,
- local or enterprise knowledge retrieval,
- document/PDF retrieval,
- social-platform retrieval,
- map/platform verification,
- other search APIs such as Baidu, Jina, Brave, Tavily, Exa, Serper, or Firecrawl.

Calling the tool `websearch` would bias the design toward one provider class. Calling it `research-gateway` makes the boundary clearer:

```text
research-gateway retrieves and normalizes external information.
industry-research reasons over that information.
```

## Architecture

```text
research-gateway/
  cmd/research-gateway/
  internal/provider/bocha/
  internal/provider/volcengine/
  internal/provider/multi/
  internal/retrieval/
  internal/output/
  internal/config/
  internal/errors/
  Makefile
  VERSION
  README.md

research-retrieval skill/
  SKILL.md
  references/provider-selection.md
  references/result-interpretation.md
  references/error-handling.md

industry-research/
  SKILL.md
  agents/*.md
  references/chain-brand-trace-reasoning.md
  references/evidence-ledger-schema.md
  evals/evals.json
```

The Go CLI lives outside `industry-research`. The skill teaches agents how to use it. `industry-research` consumes gateway outputs and writes evidence ledgers.

## Unix Philosophy

`research-gateway` must do one job well:

```text
Take a retrieval intent, call one or more providers, emit normalized JSON.
```

It must not:

- write reports,
- decide whether a claim is true,
- score industry confidence,
- hide provider failures,
- silently ignore unsupported flags,
- turn model-generated answers into facts.

Standard behavior:

- stdout: normalized JSON result,
- stderr: human-readable diagnostics,
- exit code: machine-readable success/failure category,
- no interactive prompts,
- no hidden global state,
- stable schema with provider-specific details isolated under metadata.

## CLI Command Design

Binary name:

```bash
research-gateway
```

Primary commands:

```bash
research-gateway search bocha "瑞幸 2026 门店数 招聘 扩张" --count 10 --json
research-gateway answer volcengine "搜索并总结瑞幸 2026 开店计划，给出引用来源" --limit 10 --max-keyword 3 --json
research-gateway retrieve multi "瑞幸 2026 门店数是否可信" --providers bocha,volcengine --json
research-gateway capabilities --json
research-gateway capabilities bocha --json
research-gateway version
research-gateway help
```

Mode semantics:

```text
search
  Returns a list of search-result leads. Best for direct search providers.

answer
  Returns model-mediated answer text plus citations/annotations/search calls. Best for Volcengine-like providers.

retrieve
  General intent mode. Lets each provider use its natural retrieval behavior and returns a normalized envelope.
```

Provider-specific command examples:

```bash
research-gateway search bocha "瑞幸 2026 门店数 招聘 扩张" \
  --count 10 \
  --freshness oneYear \
  --summary true \
  --json

research-gateway search bocha "瑞幸 门店数 site:luckincoffee.com" \
  --include luckincoffee.com \
  --count 10 \
  --json

research-gateway answer volcengine "搜索瑞幸咖啡近一年开店计划、供应链扩张和招聘线索，并给出引用来源" \
  --limit 10 \
  --max-keyword 3 \
  --max-tool-calls 3 \
  --sources toutiao,douyin \
  --location-country 中国 \
  --json
```

## CLI Flags

Common flags:

```text
--timeout 15s
--raw-output PATH
--json
--pretty
```

Bocha flags:

```text
--count N
--freshness oneDay|oneWeek|oneMonth|oneYear|noLimit
--summary true|false
--include domain1,domain2
--exclude domain1,domain2
```

Volcengine flags:

```text
--limit N
--max-keyword N
--max-tool-calls N
--sources search_engine,toutiao,douyin,moji
--location-country VALUE
--location-region VALUE
--location-city VALUE
--model VALUE
```

Multi-provider behavior:

```text
--providers bocha,volcengine
--count N maps to Bocha count and Volcengine limit unless provider-specific config is added later.
--freshness only applies to providers that support time filtering.
--sources only applies to providers that support source routing.
```

Unsupported provider/flag combinations must fail with `invalid_argument`. Silent ignoring is unsafe because agents may assume a constraint was applied.

## Credentials

Environment variables:

```text
BOCHA_API_KEY
ARK_API_KEY
```

Missing credentials return a structured error and non-zero exit code. The CLI must never ask for credentials interactively.

## Provider Classes

The gateway supports provider classes.

```text
direct_search
  Provider returns search results directly.
  Current provider: bocha.

model_answer_search
  Provider routes through a model that can search and answer.
  Current provider: volcengine.

agent_native_search
  Host agent search tool. Not implemented inside the Go CLI, but represented in evidence ledgers.

browser_verification
  Future interactive verification adapter.

knowledge_retrieval
  Future local or enterprise knowledge retrieval.
```

Provider class is part of the response:

```json
{
  "provider_type": "direct_search|model_answer_search|agent_native_search|browser_verification|knowledge_retrieval"
}
```

## Normalized Output Envelope

All providers emit the same top-level shape:

```json
{
  "provider": "bocha",
  "provider_type": "direct_search",
  "mode": "search",
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

Stable integration fields:

```text
provider
provider_type
mode
query
retrieved_at
request
retrieval_calls
items
answer
usage
errors
```

Provider-specific fields must stay inside `provider_metadata`.

## Item Schema

```json
{
  "rank": 1,
  "title": "示例标题",
  "url": "https://example.com/article",
  "display_url": "https://example.com/article",
  "site_name": "example.com",
  "site_icon": "",
  "snippet": "搜索结果摘要",
  "summary": "供应商提供的摘要，如有",
  "published_at": "2026-04-01T00:00:00+08:00",
  "last_crawled_at": "2026-04-01T00:00:00+08:00",
  "language": "zh",
  "content_type": "web_page|image|annotation_url|document|unknown",
  "source_confidence_hint": "lead_only",
  "provider_metadata": {}
}
```

Rules:

- Every item starts as `lead_only`.
- A search summary is not evidence.
- A model answer is not evidence.
- Evidence status can only be upgraded later by source opening, browser verification, or cross-validation in `industry-research`.

## Answer Schema

Model-mediated providers can return an answer:

```json
{
  "text": "模型基于联网搜索生成的回答文本。只能作为线索，不能直接当证据。",
  "citations": [
    {
      "index": 1,
      "url": "https://example.com/source",
      "title": "引用标题，如可用",
      "source": "search_engine|toutiao|douyin|moji|unknown"
    }
  ]
}
```

Rules:

- `answer.text` can guide follow-up search.
- `answer.text` cannot be cited as final evidence.
- Citations become `items` with `content_type: "annotation_url"` and `source_confidence_hint: "lead_only"`.

## Retrieval Calls

The gateway records what each provider actually did.

```json
{
  "call_id": "ws_001",
  "query": "瑞幸 2026 开店计划 供应链",
  "status": "completed|failed|skipped",
  "provider_action": "web-search|web_search|model_response",
  "provider_metadata": {}
}
```

For model-mediated providers, this is critical because the model may rewrite the search query or decide not to search.

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

Errors must be agent-actionable. A raw HTTP message is not enough.

Exit codes:

```text
0 success
1 invalid arguments
2 missing credentials
3 provider request failed
4 provider rate limited
5 timeout
6 partial multi-provider failure
```

## Provider Capabilities

Every provider declares capabilities.

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

Agents should use capabilities instead of memorizing provider quirks.

## Internal Go Interface

```go
type Provider interface {
    Name() string
    Type() ProviderType
    Capabilities() ProviderCapabilities
    Retrieve(ctx context.Context, req RetrievalRequest) (RetrievalResponse, error)
}
```

Package boundaries:

```text
cmd/research-gateway
  CLI parsing and command dispatch.

internal/retrieval
  Provider-neutral request, response, item, answer, call, capability, and error types.

internal/provider/bocha
  Bocha HTTP client, request mapping, response mapping, date normalization, error mapping.

internal/provider/volcengine
  Ark Responses API client, request mapping, answer/citation extraction, usage mapping, error mapping.

internal/provider/multi
  Parallel provider execution and partial failure handling.

internal/output
  JSON and pretty output.

internal/config
  Environment loading and defaults.

internal/errors
  Stable error code mapping.
```

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
- Build current platform to `./research-gateway`.
- Build release binaries under `bin/`.
- Support Linux amd64, Linux arm64, macOS amd64, macOS arm64, and Windows amd64.
- Prefer pure Go standard library HTTP clients.
- Add dependencies only when they remove real complexity.

## Bocha Provider

Bocha is a direct search provider.

Endpoint:

```text
POST https://api.bocha.cn/v1/web-search
Authorization: Bearer $BOCHA_API_KEY
Content-Type: application/json
```

Request mapping:

```text
CLI --count       -> Bocha count
CLI --freshness   -> Bocha freshness: noLimit|oneDay|oneWeek|oneMonth|oneYear
CLI --summary     -> Bocha summary
CLI --include     -> Bocha include
CLI --exclude     -> Bocha exclude
```

Validation:

- `query` is required.
- `count` must be valid for the API.
- `include` and `exclude` accept at most 100 domains each.
- `summary` defaults to `true` for agent research.

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

Date handling:

- Prefer `datePublished`.
- Bocha docs say `dateLastCrawled` values such as `2025-02-23T08:18:30Z` actually represent UTC+8 Beijing time, not UTC.
- Normalize that documented shape to `2025-02-23T08:18:30+08:00` when `datePublished` is absent.
- Preserve raw value in `provider_metadata.raw_date_last_crawled`.

Error mapping:

```text
HTTP 400 Missing parameter query   -> invalid_argument
HTTP 400 API KEY missing           -> missing_api_key
HTTP 401 Invalid API KEY           -> provider_auth_error
HTTP 403 not enough money          -> provider_quota_exhausted
HTTP 429 request limit reached     -> provider_rate_limited
HTTP 500 provider exception        -> provider_unavailable
Network timeout                    -> provider_timeout
JSON parse failure                 -> provider_parse_error
```

Populate `provider_log_id` from Bocha `log_id` when present.

## Volcengine Provider

Volcengine is a model-answer-search provider in this design.

Endpoint:

```text
POST https://ark.cn-beijing.volces.com/api/v3/responses
Authorization: Bearer $ARK_API_KEY
Content-Type: application/json
```

Default mode:

- non-streaming response for easier normalization,
- streaming may be added later for debugging or UX.

Request mapping:

```text
CLI query              -> input user message
CLI --model            -> Responses model, default doubao-seed-2-0-lite-260215
CLI --limit            -> tools[0].limit
CLI --max-keyword      -> tools[0].max_keyword
CLI --max-tool-calls   -> max_tool_calls
CLI --sources          -> tools[0].sources
CLI --location-*       -> tools[0].user_location
```

Parameter rules:

- `max_keyword` range is `1` to `50`; default `3`.
- `limit` range is `1` to `50`; default `10`; single searches may return at most 20.
- `max_tool_calls` range is `1` to `10`; default `3`.
- `sources` may include `toutiao`, `douyin`, and `moji`.
- `search_engine` can appear in usage details as the default source.
- `user_location` is optional.
- Do not send `caching`; docs say it currently returns `400`.

Response extraction:

- Extract `web_search_call` items and `action.query` into `retrieval_calls`.
- Extract `message.content[0].annotations` URLs into `items`.
- Preserve final model text in `answer.text`.
- Preserve `usage.tool_usage` and `usage.tool_usage_details`.
- Record both requested query and actual search query.
- If no search call occurred, return `no_retrieval_triggered`.

Error mapping:

```text
Missing ARK_API_KEY                -> missing_api_key
HTTP 400 invalid request/caching   -> invalid_argument
HTTP 401/403 auth or permission    -> provider_auth_error
HTTP 429 or QPS exceeded           -> provider_rate_limited
Network timeout                    -> provider_timeout
No web_search_call                 -> no_retrieval_triggered
No annotations                     -> provider_parse_error or empty_result, based on response body
```

Agent actions for `no_retrieval_triggered`:

- rewrite as an explicit search instruction,
- reduce ambiguity,
- use Bocha direct search,
- use agent-native web search.

## Research Retrieval Skill

The standalone skill should be named `research-retrieval`, not `websearch`.

It teaches agents:

1. when retrieval is required,
2. how to choose providers by capability,
3. how to generate query variants from expected traces,
4. how to interpret direct search results,
5. how to interpret model-mediated answer search,
6. how to avoid mistaking summaries or model answers for evidence,
7. how to write retrieval outputs into an evidence ledger,
8. when browser verification or source opening is required.

Core warnings:

- Retrieved output is a lead, not proof.
- Multi-provider retrieval can still return the same underlying source.
- Bocha `summary` is not a citation.
- Volcengine `answer.text` is not evidence.
- Volcengine citations are leads until opened or cross-validated.
- Provider errors must be recorded.

Agent workflow:

```text
1. State the claim and expected trace.
2. Choose provider(s) by capability and explain why.
3. Run research-gateway.
4. Save raw output if possible.
5. Convert each item into lead_only evidence.
6. Open, browser-verify, or cross-validate before upgrading evidence.
7. Record provider failure and lower confidence if source coverage is weakened.
```

## Industry Research Integration

`industry-research` should not know Bocha or Volcengine internals.

It should know:

```text
research-gateway returns retrieval outputs.
retrieval outputs become leads.
leads enter evidence_ledger.json.
only verified or cross-validated leads can support high-confidence claims.
```

Required chain-brand artifacts:

```text
entity_evidence_plan.json
ghost_deck.json
evidence_ledger.json
blue_r1.json
red_r1.json
blue_r2.json
red_r2.json
final_report.md
report_metadata.json
```

## Evidence Ledger Additions

Every search/retrieval action records:

```json
{
  "retrieval_id": "retrieval_001",
  "provider": "bocha",
  "provider_type": "direct_search",
  "mode": "search",
  "query": "瑞幸 2026 门店数 招聘 扩张",
  "retrieved_at": "2026-05-17T10:00:00+08:00",
  "purpose": "寻找门店扩张和招聘相关线索",
  "parameters": {
    "count": 10,
    "freshness": "oneYear",
    "summary": true
  },
  "status": "completed|failed|partial",
  "result_ref": "workspace/retrieval/bocha-retrieval-001.json",
  "errors": []
}
```

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

- `retrieval_result_only` cannot support high confidence.
- Provider answer text cannot be final evidence.
- Search summaries cannot replace source verification.
- Provider failures must be recorded.
- High confidence requires independent evidence families and disconfirmation attempts.
- Browser-required evidence must be marked.

## Trace Reasoning Reference

Add `references/chain-brand-trace-reasoning.md`.

It should teach patterns, not fixed platform checklists.

Example:

```text
Claim: store-count expansion
Mechanism: new stores require sites, hiring, legal setup, digital routing, supply, and demand.
Expected traces: map POI, mini-program store list, delivery platform pages, recruiting roles, local opening posts, branch/license records, user reviews.
Possible source families: maps, delivery platforms, mini-programs, recruitment, company registry, public accounts, financial reports, interviews, media.
Disconfirming traces: claimed cities with no POI, no hiring, no ordering entry, stale reviews, no legal/operator entity, or only copied press releases.
```

The key is not "always check BOSS" or "always check Qichacha". The key is:

```text
If the claim is true, what traces should exist?
Which source families are most likely to expose those traces?
Which missing traces would weaken or falsify the claim?
```

## Anti-Template Guardrails

Bad:

```text
The agent checked recruitment, registry, maps, public accounts, and news because the playbook listed them.
```

Good:

```text
The agent checked recruitment because store expansion needs labor, registry because new operators or franchise entities may appear, and maps/mini-programs because real stores must become customer-addressable. It skipped supplier tenders because the current claim was store count, not supply-chain maturity.
```

Validation checks:

- Every retrieval action has a purpose.
- Every retrieval action maps to an expected trace.
- Every evidence item maps to a claim.
- Every high-confidence claim has disconfirmation attempts.
- The report identifies confidence limiters.
- The report says when a result is only a lead.

## Browser Automation Boundary

Browser automation is not in the first-version Go CLI.

The schema must preserve handoff fields:

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

This prevents the system from pretending interactive platform evidence was verified by a simple retrieval call.

## Validation

First-version validation should check:

- `evidence_ledger.json` exists for chain-brand research.
- Every major claim has expected traces.
- Every retrieval action has provider, provider type, mode, query, purpose, timestamp, and status.
- Provider errors include `agent_action`.
- Every evidence item maps to a claim.
- High-confidence claims have independent evidence families.
- Retrieval-only evidence is not used as final evidence.
- Volcengine answer text is not cited as evidence.
- Bocha summaries are not cited as final evidence unless the source page has been opened or cross-validated.
- Claims without disconfirmation attempts are downgraded.
- Browser-required evidence is clearly marked.

## Evaluation Cases

Add or revise evals around:

1. Listed chain brand expansion: "瑞幸咖啡 2026 年门店数目标是否可信？"
2. Non-listed franchise tea brand: "某区域茶饮品牌宣称 5000 家门店、加盟商盈利稳定是否可信？"
3. Supply-chain maturity: "某餐饮品牌宣称全国统一供应链成熟是否可信？"
4. City coverage: "某生鲜零售品牌宣称全国冷链覆盖是否真实？"
5. Management interview verification: "把管理层访谈中的 GMV、单店 UE、同店增长、供应链统一供货比例转成可验证命题。"

Each eval should expect:

- trace reasoning before retrieval,
- retrieval actions with purposes,
- evidence ledger output,
- disconfirmation attempts,
- confidence downgrades where sources are weak,
- clear separation between leads and evidence.

## External References

- `docs/bocha_websearch.md` documents Bocha Web Search request parameters, response fields, endpoint `https://api.bocha.cn/v1/web-search`, date handling caveat, and error codes.
- `docs/volcengine_websearch.md` documents Volcengine Ark Responses API Web Search tool behavior, tool parameters, sources, location, usage fields, and model-mediated search-call output.
- Bocha Open Platform documents Web Search API shape and endpoint examples: `https://open.bocha.cn/`
- OpenAI documents the Responses API `web_search` tool, citations, sources, domain filtering, and live access controls: `https://developers.openai.com/api/docs/guides/tools-web-search`
- Volcengine documents联网搜索 as a deep-research capability with real-time data access, search strategy planning, multi-source verification, and structured report output: `https://www.volcengine.com/docs/85637/1588465`
- Volcengine联网搜索 API reference page exists but currently requires JavaScript for full details in this environment: `https://www.volcengine.com/docs/87772/2272953`

## Acceptance Criteria

The design is successful when:

1. `research-gateway` can be implemented as an independent Go CLI without depending on `industry-research`.
2. Bocha and Volcengine emit the same normalized envelope while preserving provider-type differences.
3. `industry-research` consumes gateway outputs without knowing provider internals.
4. Reports show claim-to-trace reasoning before conclusions.
5. Retrieved results are treated as leads until source evidence is verified.
6. Confidence scores are tied to evidence-family independence and disconfirmation attempts.
7. Agents remain free to choose sources but must explain why those sources fit the claim.
8. Provider errors are agent-actionable.
9. Volcengine model-mediated outputs are not misrepresented as direct search-result lists.

## Non-Goals

The first version does not:

- implement integrations for a specific agent harness,
- implement browser automation,
- build a database-backed evidence store,
- add UI,
- scrape closed platforms,
- guarantee access to BOSS, Qichacha, WeChat, video accounts, maps, delivery platforms, or review platforms,
- replace agent-native web search,
- let retrieval providers determine truth.

## Implementation Order

Recommended sequence:

1. Create standalone Go CLI skeleton and Makefile.
2. Define provider-neutral retrieval types.
3. Add mocked provider tests.
4. Add Bocha provider.
5. Add Volcengine provider.
6. Add multi-provider retrieval.
7. Add capabilities command.
8. Create `research-retrieval` skill.
9. Add chain-brand trace reasoning reference to `industry-research`.
10. Add `evidence_ledger.json` contract.
11. Update agents to write and consume the ledger.
12. Update validation and evals.

This keeps the retrieval gateway reusable and prevents provider details from leaking into the research workflow.
