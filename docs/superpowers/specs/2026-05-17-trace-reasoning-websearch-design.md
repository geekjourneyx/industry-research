# Trace Reasoning and Pluggable Web Search Design

## Summary

This design upgrades `industry-research` from a report-writing workflow into a trace-reasoning research system for chain brands, restaurant, retail, franchise, and supply-chain cases.

The first version builds three pieces:

1. A reusable Go CLI named `websearch` that provides Bocha and Volcengine search through a Unix-style command interface.
2. A standalone `websearch` skill that teaches agents when and how to use the CLI, agent-native web search, and multi-source search.
3. A revised `industry-research` integration model that uses trace reasoning, evidence ledgers, and confidence gates before producing reports.

The central rule is: **fix the quality gates, not the exploration path**.

Agents must not mechanically follow a platform checklist. They must reason from a claim to the real-world traces that should exist, choose sources based on that reasoning, record what they searched and why, look for disconfirming evidence, and only then issue a confidence judgment.

## Problem

Current `industry-research` already has useful concepts: red-blue debate, operating traces, triangulation, confidence labels, and suspended judgment. The weakness is that these concepts still live mostly as instructions.

When asked questions such as "Is Luckin likely to reach a certain 2026 store count?", a human domain expert naturally thinks:

- Expansion should leave hiring traces.
- Store growth should leave map, mini-program, delivery platform, and review traces.
- Legal expansion should leave company registry, branch, license, and franchise traces.
- Management claims should be checked against financial reports, announcements, interviews, public-account posts, video accounts, and authoritative media.
- Missing expected traces can be evidence, not just absence of data.

The current agent often lacks this trace-reasoning step. It may search the obvious phrase, collect related articles, and write a plausible report without proving that the real-world operating system actually changed.

## Design Principles

### First Principles

A commercial claim is credible only if the world would have changed in observable ways when the claim is true.

The workflow must therefore move through this chain:

```text
claim
-> real-world mechanism
-> expected traces
-> likely source families
-> search actions
-> evidence ledger
-> disconfirmation search
-> confidence judgment
-> report
```

The agent should be judged by whether this chain is visible and defensible, not by whether it followed a fixed search checklist.

### Unix Philosophy

The search layer must be a small tool that does one job well.

`websearch` CLI:

- Accepts a query and provider options.
- Calls a provider.
- Emits normalized JSON to stdout.
- Emits diagnostics and errors to stderr.
- Uses exit codes to signal success or failure.
- Does not write reports.
- Does not judge truth.
- Does not perform industry reasoning.

`industry-research` remains the reasoning layer. It consumes search outputs, builds the evidence ledger, performs trace reasoning, and makes confidence judgments.

### Fixed vs Open

Fixed:

- Every important conclusion must map to a claim.
- Every claim must define expected traces.
- Every evidence item must record source, search provider, timestamp, and independence notes.
- High confidence requires independent evidence families, not repeated citations.
- Disconfirmation search is mandatory before strong claims.
- Unsupported claims must be suspended or downgraded.

Open:

- Which platforms to search.
- Which query variants to try.
- When to continue digging.
- When an anomaly is more important than the original question.
- Whether a report structure should be adjusted around uncertainty.

## Architecture

```text
websearch-cli/
  cmd/websearch/
  internal/provider/bocha/
  internal/provider/volcengine/
  internal/provider/multi/
  internal/search/
  internal/output/
  internal/config/
  Makefile
  VERSION
  README.md

websearch skill
  SKILL.md
  references/provider-selection.md
  references/result-interpretation.md

industry-research
  SKILL.md
  references/chain-brand-trace-reasoning.md
  references/evidence-ledger-schema.md
  agents/*.md
  evals/evals.json
```

The Go CLI should live outside `industry-research`. It is a reusable foundation tool. The `websearch` skill wraps usage rules around the CLI. `industry-research` calls the skill conceptually and consumes the standardized outputs.

## Go CLI

### Command Shape

The binary name is `websearch`.

Required first-version commands:

```bash
websearch bocha "瑞幸 2026 门店数 招聘 扩张" --count 10 --json
websearch volcengine "瑞幸 2026 开店计划 供应链" --limit 10 --max-keyword 3 --json
websearch multi "瑞幸 2026 门店数" --providers bocha,volcengine --count 10 --json
websearch version
websearch help
```

Optional flags for version one:

```text
--count N
--freshness oneDay|oneWeek|oneMonth|oneYear|noLimit
--summary true|false
--include domain1,domain2
--exclude domain1,domain2
--limit N
--max-keyword N
--max-tool-calls N
--sources search_engine,toutiao,douyin,moji
--location-country VALUE
--location-region VALUE
--location-city VALUE
--model VALUE
--timeout 15s
--raw-output PATH
--json
--pretty
```

Provider-specific flag behavior:

```text
bocha:
  uses --count, --freshness, --summary, --include, --exclude, --timeout, --raw-output, --json, --pretty

volcengine:
  uses --limit, --max-keyword, --max-tool-calls, --sources, --location-*, --model, --timeout, --raw-output, --json, --pretty

multi:
  accepts common intent flags and maps them per provider:
    --count N maps to Bocha count and Volcengine limit unless provider-specific flags are supplied later
    --freshness only applies to providers that support time filtering
    --sources only applies to providers that support source routing
```

The CLI should reject unsupported provider/flag combinations with a clear `invalid_argument` error. Silent ignoring is not agent-friendly because agents may believe a constraint was applied when it was not.

Provider keys are read from environment variables:

```text
BOCHA_API_KEY
ARK_API_KEY
```

The CLI must not prompt interactively. Missing credentials return a structured error and non-zero exit code.

### Provider Classes

The CLI must support two provider classes.

Direct search providers return search results directly from a search endpoint. Bocha is a direct search provider.

Model-mediated search providers expose search through a model response API. Volcengine Web Search in the provided documentation is model-mediated: the request goes to the Ark Responses API, the model decides whether to invoke `web_search`, and the response contains search-call events, answer text, annotations, and usage details. The CLI must not pretend that Volcengine returns the same raw list shape as Bocha. It should normalize what is observable and preserve provider-specific metadata.

Provider class field:

```json
{
  "provider_class": "direct_search|model_mediated_search"
}
```

### Output Contract

Successful search output:

```json
{
  "provider": "bocha",
  "provider_class": "direct_search",
  "query": "瑞幸 2026 门店数 招聘 扩张",
  "searched_at": "2026-05-17T10:00:00+08:00",
  "request": {
    "count": 10,
    "freshness": "oneYear",
    "summary": true,
    "include": [],
    "exclude": []
  },
  "search_calls": [
    {
      "call_id": "bocha_http_001",
      "query": "瑞幸 2026 门店数 招聘 扩张",
      "status": "completed",
      "provider_action": "web-search"
    }
  ],
  "results": [
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
      "content_type": "web_page",
      "source_confidence_hint": "lead_only",
      "provider_metadata": {
        "bocha_id": "provider result id when present",
        "cached_page_url": "",
        "is_navigational": false,
        "is_family_friendly": true
      }
    }
  ],
  "images": [],
  "answer_text": "",
  "usage": {},
  "errors": []
}
```

Volcengine model-mediated output:

```json
{
  "provider": "volcengine",
  "provider_class": "model_mediated_search",
  "query": "瑞幸 2026 开店计划 供应链",
  "searched_at": "2026-05-17T10:00:00+08:00",
  "request": {
    "model": "doubao-seed-2-0-lite-260215",
    "limit": 10,
    "max_keyword": 3,
    "max_tool_calls": 3,
    "sources": ["search_engine", "toutiao", "douyin"],
    "user_location": {
      "country": "中国",
      "region": "浙江",
      "city": "杭州"
    }
  },
  "search_calls": [
    {
      "call_id": "ws_001",
      "query": "瑞幸 2026 开店计划 供应链",
      "status": "completed",
      "provider_action": "web_search"
    }
  ],
  "results": [
    {
      "rank": 1,
      "title": "从 annotation 或引用中提取的标题，如可用",
      "url": "https://example.com/source",
      "display_url": "https://example.com/source",
      "site_name": "example.com",
      "snippet": "从 annotation 周边文本或模型输出中提取的摘要，如可用",
      "summary": "",
      "published_at": "",
      "content_type": "annotation_url",
      "source_confidence_hint": "lead_only",
      "provider_metadata": {
        "annotation_index": 0,
        "source": "search_engine|toutiao|douyin|moji|unknown"
      }
    }
  ],
  "images": [],
  "answer_text": "模型基于联网搜索生成的回答文本。只能作为线索，不能直接当证据。",
  "usage": {
    "tool_usage": {"web_search": 2},
    "tool_usage_details": {"web_search": {"search_engine": 2, "toutiao": 1}}
  },
  "errors": []
}
```

Multi-provider output:

```json
{
  "provider": "multi",
  "provider_class": "multi",
  "query": "瑞幸 2026 门店数",
  "searched_at": "2026-05-17T10:00:00+08:00",
  "request": {
    "providers": ["bocha", "volcengine"],
    "count": 10,
    "freshness": "oneYear"
  },
  "provider_results": [
    {
      "provider": "bocha",
      "provider_class": "direct_search",
      "results": []
    },
    {
      "provider": "volcengine",
      "provider_class": "model_mediated_search",
      "results": []
    }
  ],
  "errors": []
}
```

Failure output goes to stdout only when `--json` is requested, and diagnostics still go to stderr:

```json
{
  "provider": "bocha",
  "provider_class": "direct_search",
  "query": "瑞幸",
  "searched_at": "2026-05-17T10:00:00+08:00",
  "request": {},
  "results": [],
  "errors": [
    {
      "code": "missing_api_key",
      "message": "BOCHA_API_KEY is not set",
      "retryable": false,
      "agent_action": "Set BOCHA_API_KEY or rerun with another provider."
    }
  ]
}
```

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

Error object contract:

```json
{
  "code": "missing_api_key|invalid_argument|provider_http_error|provider_auth_error|provider_quota_exhausted|provider_rate_limited|provider_timeout|provider_unavailable|provider_parse_error|no_search_triggered|partial_failure",
  "message": "Human-readable error for agents and logs.",
  "provider_status": 429,
  "provider_code": "429",
  "provider_log_id": "c66aac17eab1bb7e",
  "retryable": true,
  "agent_action": "Wait and retry with lower count, or use another provider.",
  "raw_error_path": "workspace/search/raw/bocha-error-001.json"
}
```

Agent-facing errors must say what to do next. A bare provider error is not enough.

### Extensibility Model

Every provider must implement the same internal interface:

```go
type Provider interface {
    Name() string
    Class() ProviderClass
    Search(ctx context.Context, req SearchRequest) (SearchResponse, error)
    Capabilities() ProviderCapabilities
}
```

Capabilities should be machine-readable:

```json
{
  "provider": "bocha",
  "provider_class": "direct_search",
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
  "provider_class": "model_mediated_search",
  "supports_freshness": false,
  "supports_include_domains": false,
  "supports_exclude_domains": false,
  "supports_summary": false,
  "supports_location": true,
  "supports_sources": true,
  "supports_images": false,
  "supports_model_choice": true,
  "result_kinds": ["annotation_url", "answer_text", "search_call"]
}
```

The CLI should expose capabilities:

```bash
websearch capabilities --json
websearch capabilities bocha --json
```

This lets agents decide which provider fits a search purpose without memorizing provider quirks.

Future providers such as Baidu, Jina, Brave, Tavily, or browser-verification adapters must map into the same envelope and declare their capabilities. Provider-specific fields must stay inside `provider_metadata`; the top-level response should remain stable.

Backward compatibility rule:

- New fields may be added.
- Existing fields must not change meaning.
- Provider-specific fields must not become required for all providers.
- Unknown provider metadata must be preserved by consumers.
- The `provider_class`, `results`, `search_calls`, `errors`, and `usage` fields are stable integration points.

### Go Project Practices

The Go CLI should follow the pattern used by `~/Workspace/go/md2wechat-skill/Makefile`.

Required Makefile targets:

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
- Build current platform to `./websearch`.
- Build release binaries under `bin/`.
- Support Linux amd64, Linux arm64, macOS amd64, macOS arm64, and Windows amd64.
- Keep implementation pure Go unless a provider requires a small dependency.

Package boundaries:

- `cmd/websearch`: CLI parsing and command dispatch.
- `internal/search`: provider-neutral request and response types.
- `internal/provider/bocha`: Bocha HTTP client and mapping.
- `internal/provider/volcengine`: Volcengine HTTP client and mapping.
- `internal/provider/multi`: parallel provider execution and partial failure handling.
- `internal/output`: JSON and pretty output.
- `internal/config`: environment loading and timeout defaults.

## Provider Notes

### Bocha

Bocha exposes a Web Search API at `https://api.bocha.cn/v1/web-search` and returns structured web page results including `name`, `url`, `displayUrl`, `snippet`, `summary`, `siteName`, `siteIcon`, `datePublished`, and `dateLastCrawled` where available.

Request mapping:

```text
CLI --count       -> Bocha count
CLI --freshness   -> Bocha freshness: noLimit|oneDay|oneWeek|oneMonth|oneYear
CLI --summary     -> Bocha summary
CLI --include     -> Bocha include, comma-separated or pipe-separated domains
CLI --exclude     -> Bocha exclude, comma-separated or pipe-separated domains
```

Validation rules:

- `query` is required and must not be empty after trimming whitespace.
- `count` must be clamped or rejected outside the documented range. The implementation should default to rejecting invalid values so agents see bad parameters early.
- `include` and `exclude` accept at most 100 domains each.
- `summary` should default to `true` for agent-facing research because summaries are useful leads, but the result must still be labeled `lead_only`.

Response mapping:

```text
webPages.value[].name             -> results[].title
webPages.value[].url              -> results[].url
webPages.value[].displayUrl       -> results[].display_url
webPages.value[].snippet          -> results[].snippet
webPages.value[].summary          -> results[].summary
webPages.value[].siteName         -> results[].site_name
webPages.value[].siteIcon         -> results[].site_icon
webPages.value[].datePublished    -> results[].published_at
webPages.value[].dateLastCrawled  -> results[].last_crawled_at
webPages.value[].cachedPageUrl    -> provider_metadata.cached_page_url
webPages.value[].language         -> results[].language
images.value[]                    -> images[]
```

Bocha date handling:

- Prefer `datePublished` when present.
- Treat `dateLastCrawled` carefully. The provided Bocha documentation says values such as `2025-02-23T08:18:30Z` actually represent UTC+8 Beijing time, not UTC. The provider mapper must normalize this to `2025-02-23T08:18:30+08:00` when the value has the documented shape and `datePublished` is absent.
- Preserve the raw value in `provider_metadata.raw_date_last_crawled`.

Bocha error mapping:

```text
HTTP 400 Missing parameter query        -> invalid_argument, retryable false
HTTP 400 API KEY missing                -> missing_api_key, retryable false
HTTP 401 Invalid API KEY                -> provider_auth_error, retryable false
HTTP 403 not enough money               -> provider_quota_exhausted, retryable false
HTTP 429 request limit reached          -> provider_rate_limited, retryable true
HTTP 500 provider exception             -> provider_unavailable, retryable true
Network timeout                         -> provider_timeout, retryable true
JSON parse failure                      -> provider_parse_error, retryable false unless raw body indicates transient provider failure
```

For Bocha, `provider_log_id` must be populated from `log_id` when present.

Use Bocha primarily for:

- Chinese web search.
- General public web, news, company pages, reports, and media.
- Query variants where semantic search helps retrieve broader context.

### Volcengine

Volcengine documents联网搜索 as a Responses API `web_search` tool. This is model-mediated search, not a direct search endpoint like Bocha.

Request shape:

```text
POST https://ark.cn-beijing.volces.com/api/v3/responses
Authorization: Bearer $ARK_API_KEY
Content-Type: application/json
```

The CLI should use non-streaming mode by default for easier JSON normalization. Streaming can be added later, but version one should avoid it unless needed for debugging.

Required request mapping:

```text
CLI query              -> input user message
CLI --model            -> Responses model, default doubao-seed-2-0-lite-260215
CLI --limit            -> tools[0].limit
CLI --max-keyword      -> tools[0].max_keyword
CLI --max-tool-calls   -> max_tool_calls
CLI --sources          -> tools[0].sources
CLI --location-*       -> tools[0].user_location
```

Volcengine parameters:

- `max_keyword`: controls how many search keywords the model may use in a round. Documented range is `1` to `50`; default should be `3`.
- `limit`: controls how many results each search operation returns. Documented range is `1` to `50`, but single searches may return at most 20. Default should be `10`.
- `max_tool_calls`: controls how many web-search rounds the model may execute. Documented range is `1` to `10`; default should be `3`.
- `sources`: may include `toutiao`, `douyin`, and `moji`. The default web source is `search_engine`; it may appear in usage details even when not explicitly listed.
- `user_location`: optional approximate country, region, and city.
- `caching`: must not be sent because the documentation says it currently returns a `400` error.

Response extraction:

- Extract `web_search_call` items and their `action.query` values into `search_calls`.
- Extract `message.content[0].annotations` URLs into `results`.
- Preserve final answer text in `answer_text`.
- Preserve `usage.tool_usage` and `usage.tool_usage_details`.
- If no `web_search_call` occurred, return `no_search_triggered` with exit code `3` unless the caller used a future `--allow-no-search` flag.
- Because the model controls whether search is triggered and which keywords are used, record both requested query and actual search queries.

Volcengine result confidence:

- `answer_text` is never evidence by itself.
- Annotation URLs are leads until the source pages are opened or independently confirmed.
- The result should carry `source_confidence_hint: "lead_only"` by default.

Volcengine error mapping:

```text
Missing ARK_API_KEY                   -> missing_api_key, retryable false
HTTP 400 invalid request/caching      -> invalid_argument, retryable false
HTTP 401/403 auth or permission       -> provider_auth_error, retryable false
HTTP 429 or QPS exceeded              -> provider_rate_limited, retryable true
Network timeout                       -> provider_timeout, retryable true
Model response has no search call     -> no_search_triggered, retryable false unless query/prompt is revised
Response lacks annotations            -> provider_parse_error or no_results, depending on response body
```

Agent action for `no_search_triggered` should recommend one of:

- Rewrite the query as a search instruction.
- Reduce ambiguity.
- Use Bocha direct search.
- Use agent-native web search.

Use Volcengine primarily for:

- Chinese real-time information.
- ByteDance ecosystem adjacent content signals where `sources` such as `toutiao` or `douyin` are useful.
- Time-sensitive news, platform content clues, and public-domain monitoring.

### Agent-Native Web Search

Agent-native web search remains useful for:

- Official-domain filtering.
- English and international sources.
- Search actions where the host agent can open pages and return sources.
- Fallback when external provider keys are unavailable.

Agent-native web search is not part of the Go CLI. It is part of the agent runtime and is selected by the `websearch` skill when appropriate.

## Websearch Skill

The standalone `websearch` skill is a usage layer, not a provider implementation.

It should teach agents:

1. When search is required.
2. How to choose providers.
3. How to generate query variants.
4. How to interpret search results.
5. How to avoid mistaking search summaries for evidence.
6. How to route results into an evidence ledger.
7. How to recognize provider-specific limitations and ask for browser or source-page verification.

Provider-selection guidance:

```text
Chinese chain-brand research:
  start with bocha + volcengine multi-search

time-sensitive Chinese news or platform signals:
  include volcengine

general Chinese web and company/report discovery:
  include bocha

official-domain verification or English/international sources:
  use agent-native web search with domain filters when available

provider failure:
  continue with available providers, record the failure, lower confidence where source coverage is weakened
```

The skill must warn agents:

- Search results are leads, not evidence.
- Multiple search providers can still return the same underlying source.
- A summary is not a citation.
- Volcengine answer text is not evidence; only source annotations and subsequently opened URLs can become evidence.
- Bocha summaries are useful for triage but remain `lead_only` until the actual source page is inspected or cross-validated.
- `dateLastCrawled` from Bocha must not be interpreted as UTC when the mapper has normalized it from the documented UTC+8 compatibility issue.
- Final claims must cite opened source URLs or explicitly remain unverified.

Agent-friendly search workflow:

```text
1. State the claim and expected trace.
2. Choose provider(s) and explain why.
3. Run search.
4. Record raw search output path if available.
5. Convert each result into lead_only evidence.
6. Open or otherwise verify source URLs before upgrading evidence.
7. If provider fails, record failure and use alternate provider or lower confidence.
8. If Volcengine does not trigger search, rewrite the query or use Bocha/direct web search.
```

The skill should prefer commands that make agent intent explicit:

```bash
websearch bocha "瑞幸 2026 门店数 招聘 扩张" \
  --count 10 \
  --freshness oneYear \
  --summary true \
  --json

websearch bocha "瑞幸 门店数 site:luckincoffee.com" \
  --include luckincoffee.com \
  --count 10 \
  --json

websearch volcengine "搜索瑞幸咖啡近一年开店计划、供应链扩张和招聘线索，并给出引用来源" \
  --limit 10 \
  --max-keyword 3 \
  --max-tool-calls 3 \
  --sources toutiao,douyin \
  --location-country 中国 \
  --json
```

## Industry Research Integration

### New Reference: Chain Brand Trace Reasoning

Add a reference file such as `references/chain-brand-trace-reasoning.md`.

It should encode trace reasoning patterns, not fixed platform checklists.

Examples:

```text
Claim: store-count expansion
Mechanism: new stores require sites, hiring, legal setup, digital routing, supply, and demand.
Expected traces: map POI, mini-program store list, delivery platform pages, recruiting roles, local opening posts, branch/license records, user reviews.
Possible sources: maps, delivery platforms, mini-program, BOSS/Zhipin-style hiring, company registry, public accounts, financial reports, media interviews.
Disconfirming traces: claimed cities with no POI, no hiring, no ordering entry, stale reviews, no legal/operator entity, or only copied press releases.
```

```text
Claim: supply-chain maturity
Mechanism: stable supply requires production, warehousing, cold-chain or dry-chain logistics, quality control, purchasing, and dispatch.
Expected traces: central kitchen or warehouse addresses, production/food licenses, warehouse/driver/replenishment roles, logistics tenders, supplier mentions, delivery radius evidence, SKU availability differences.
Possible sources: license databases, recruitment platforms, maps, tender sites, supplier announcements, mini-program SKU availability, user delivery complaints.
Disconfirming traces: no warehouse or logistics roles, route claims without nodes, regional expansion without supply nodes, SKU unavailability across claimed coverage.
```

```text
Claim: franchisee profitability or stability
Mechanism: franchisees need operators, local stores, ordering systems, training, supply, fees, and support.
Expected traces: franchisee entities, recruitment by franchisees, store-level reviews, disputes, local operating posts, franchise recruitment materials, court records, complaint platforms.
Possible sources: company registry, legal/court data, black-cat style complaints, social platforms, public accounts, recruitment, franchise websites.
Disconfirming traces: high churn, disputes, abnormal closures, no operator trace, claims based only on招商 material.
```

The reference should include examples for:

- Listed or near-listed chain brands.
- Non-listed franchise brands.
- Store count claims.
- 2026 expansion targets.
- Supply-chain maturity claims.
- City coverage claims.
- Management interview claims.
- Franchisee profitability claims.

### Evidence Ledger

Add `evidence_ledger.json` as a required intermediate artifact for chain-brand research.

Suggested schema:

```json
{
  "research_question": "瑞幸咖啡 2026 年门店数目标是否可信？",
  "claims": [
    {
      "claim_id": "claim_store_count_2026",
      "claim": "瑞幸 2026 年门店数继续快速增长具备经营支撑",
      "mechanism": "门店增长需要选址、招聘、供应链、数字入口和用户需求共同支撑",
      "expected_traces": [
        {
          "trace_type": "people_org",
          "trace": "新开城或加密城市出现店长、咖啡师、区域运营、拓展岗位",
          "why_expected": "门店扩张前后必须补充门店和区域运营人员"
        }
      ],
      "search_actions": [
        {
          "search_id": "search_001",
          "provider": "bocha",
          "provider_class": "direct_search",
          "query": "瑞幸 2026 门店数 招聘 扩张",
          "searched_at": "2026-05-17T10:00:00+08:00",
          "purpose": "寻找门店扩张和招聘相关线索",
          "parameters": {
            "count": 10,
            "freshness": "oneYear",
            "summary": true
          },
          "status": "completed|failed|partial",
          "result_ref": "workspace/search/bocha-search-001.json"
        }
      ],
      "evidence_items": [
        {
          "evidence_id": "ev_001",
          "source_url": "https://example.com",
          "source_title": "示例来源",
          "source_type": "company_disclosure|official_registry|recruiting|map_poi|platform_frontend|media|social|legal|tender|ugc|search_result_only",
          "evidence_family": "capital_legal|people_org|physical_fulfillment|digital_frontend|terminal_feedback|management_narrative",
          "origin_provider": "bocha|volcengine|agent_websearch|browser|manual",
          "origin_search_id": "search_001",
          "accessed_at": "2026-05-17T10:10:00+08:00",
          "verification_status": "search_result_only|source_opened|browser_verified|cross_validated|not_accessible",
          "independence_note": "不是公司通稿转载，独立于财报口径",
          "supports_or_challenges": "supports|challenges|mixed|lead_only",
          "summary": "这条证据说明了什么",
          "requires_browser_verification": false,
          "browser_verification_reason": ""
        }
      ],
      "disconfirmation_attempts": [
        {
          "attempt_id": "disconfirm_001",
          "question": "是否存在声称覆盖但无招聘、无门店、不可下单的城市？",
          "search_or_check": "搜索城市覆盖反证和门店异常反馈",
          "result": "not_found|found|not_checked",
          "impact": "提高/降低/不改变置信度"
        }
      ],
      "confidence": {
        "rating": "high|medium|low|suspended|unverified",
        "reason": "三类独立证据支持，且未发现关键反证",
        "limiting_factors": ["地图和小程序未完成浏览器验证"]
      }
    }
  ]
}
```

Rules:

- `search_result_only` can never support a high-confidence claim.
- Search summaries are `lead_only` until the source page is opened or independently confirmed.
- Volcengine `answer_text` can guide the next search, but cannot be used as a source in the final report.
- Bocha `summary` can guide triage, but cannot replace source-page verification.
- A provider failure must be recorded as a search action with `status: failed`; otherwise the reader cannot distinguish "not checked" from "checked and failed".
- A high-confidence claim must have at least three independent evidence families unless the report explains a domain-specific reason for a different threshold.
- A claim with no disconfirmation attempt cannot be high confidence.
- If a needed source requires browser interaction, the ledger must mark it.

### Agent Contract Changes

Engagement Manager:

- Converts the user question into claims.
- Defines mechanisms and expected traces.
- Proposes source families and search purposes.
- Creates the initial `evidence_ledger.json` skeleton.

Blue Team:

- Searches for supporting evidence.
- Must explain why each source is relevant to the expected trace.
- Cannot upgrade search-result summaries into facts.

Red Team:

- Searches for disconfirming evidence.
- Treats missing expected traces as possible evidence.
- Must distinguish "not found because data unavailable" from "not found despite expected public trace".

Chief Arbitrator:

- Judges from the evidence ledger, not from prose alone.
- Downgrades claims without independent evidence families.
- Suspends claims where expected traces are missing or provider coverage is weak.

## Anti-Template Guardrails

The system must reject mechanical platform-checking.

Bad behavior:

```text
The agent searched BOSS, Qichacha, WeChat, maps, and news because the SOP listed them.
```

Good behavior:

```text
The agent searched recruitment because store expansion requires labor, registry data because new operators or franchise entities may exist, and maps/mini-program because real stores must become customer-addressable. It skipped supplier tender search because the claim was store count rather than supply-chain capacity.
```

Guardrail checks:

- Each search action must include `purpose`.
- Each evidence item must map to an expected trace.
- The agent must list at least one source family it deliberately did not search and explain why, unless the report is brief.
- The final report must include a short "confidence limiter" paragraph for each major claim.
- The report must explicitly say when a result is only a lead.

## Browser Automation Boundary

Browser automation is not part of the first-version Go CLI.

However, the design must preserve browser handoff fields:

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

This prevents the system from pretending closed or interactive platforms were verified through simple web search.

## Validation

First-version validation should remain mostly document and schema based.

Required checks:

- `evidence_ledger.json` exists for chain-brand research.
- Every major claim has expected traces.
- Every evidence item maps to a claim.
- Every search action has provider, query, purpose, and timestamp.
- Every search action records provider class and status.
- Provider errors include `agent_action`.
- High-confidence claims have enough independent evidence families.
- Search-result-only evidence is not used as a high-confidence source.
- Volcengine answer text is not cited as evidence.
- Bocha summaries are not cited as final evidence unless the source page has been opened or cross-validated.
- Claims without disconfirmation attempts are downgraded.
- Browser-required evidence is clearly marked.

The current `validate_report.py` can be extended later, but the first spec should define the requirements before implementation.

## Evaluation Cases

Add or revise evals around:

1. Listed chain brand expansion: "瑞幸咖啡 2026 年门店数目标是否可信？"
2. Non-listed franchise tea brand: "某区域茶饮品牌宣称 5000 家门店、加盟商盈利稳定是否可信？"
3. Supply-chain maturity: "某餐饮品牌宣称全国统一供应链成熟是否可信？"
4. City coverage: "某生鲜零售品牌宣称全国冷链覆盖是否真实？"
5. Management interview verification: "把管理层访谈中的 GMV、单店 UE、同店增长、供应链统一供货比例转成可验证命题。"

Each eval should expect:

- Trace reasoning before search.
- Evidence ledger output.
- Multi-provider search plan or executed search output when credentials exist.
- Disconfirmation attempts.
- Confidence downgrades where sources are weak.

## External References

- `docs/bocha_websearch.md` documents Bocha Web Search request parameters, response fields, endpoint `https://api.bocha.cn/v1/web-search`, date handling caveat, and error codes.
- `docs/volcengine_websearch.md` documents Volcengine Ark Responses API Web Search tool behavior, tool parameters, sources, location, usage fields, and model-mediated search-call output.
- Bocha Open Platform documents Web Search API shape and endpoint examples: `https://open.bocha.cn/`
- OpenAI documents the Responses API `web_search` tool, citations, sources, domain filtering, and live access controls: `https://developers.openai.com/api/docs/guides/tools-web-search`
- Volcengine documents联网搜索 as a deep-research capability with real-time data access, search strategy planning, multi-source verification, and structured report output: `https://www.volcengine.com/docs/85637/1588465`
- Volcengine联网搜索 API reference page exists but currently requires JavaScript for full details in this environment: `https://www.volcengine.com/docs/87772/2272953`

## Acceptance Criteria

The design is successful when:

1. `websearch` can be implemented as an independent Go CLI without depending on `industry-research`.
2. Bocha and Volcengine providers can emit the same normalized envelope while preserving their provider-class differences.
3. `industry-research` can consume search outputs without knowing provider internals.
4. Chain-brand reports show claim-to-trace reasoning before conclusions.
5. Search results are treated as leads until source evidence is verified.
6. Confidence scores are tied to evidence-family independence and disconfirmation attempts.
7. The agent remains free to choose sources, but cannot skip explaining why those sources fit the claim.
8. Provider errors are agent-actionable rather than raw HTTP messages only.
9. Volcengine model-mediated outputs are not misrepresented as direct search-result lists.

## Non-Goals

The first version does not:

- Implement browser automation.
- Build a database-backed evidence store.
- Add UI.
- Scrape closed platforms.
- Guarantee access to BOSS, Qichacha, WeChat, video accounts, maps, delivery platforms, or review platforms.
- Replace agent-native web search.
- Let search providers directly determine truth.

## Implementation Order

Recommended implementation sequence:

1. Create the standalone Go CLI with mocked provider tests.
2. Add Bocha provider.
3. Add Volcengine provider behind the same interface.
4. Add multi-provider search.
5. Create the `websearch` skill.
6. Add chain-brand trace reasoning reference to `industry-research`.
7. Add `evidence_ledger.json` contract to `industry-research`.
8. Update agents to write and consume the ledger.
9. Update validation and evals.

This order keeps the search tool reusable and prevents provider details from leaking into the research workflow.
