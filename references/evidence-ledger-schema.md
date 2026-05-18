# Evidence Ledger Schema

`evidence_ledger.json` records every evidence item used by `researcher` and `industry-research`.

Rules:

- Retrieval-only items cannot support high confidence.
- Model answer text cannot be final evidence.
- Search summaries cannot replace source verification.
- Provider failures must be recorded in retrieval logs.
- High confidence requires independent evidence families and disconfirmation attempts.
- Browser-required evidence must be marked.

Required evidence item fields:

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

