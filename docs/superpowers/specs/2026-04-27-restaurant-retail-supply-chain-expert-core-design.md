# Restaurant, Retail, and Supply Chain Expert Core Design

## Current Status

This is the earlier vertical-domain design for restaurant, retail, and supply-chain research. It remains valid as the domain reasoning layer, but it is no longer the full system architecture.

The current system adds `researcher` as the execution layer. The domain ideas in this document now map mainly to:

- `trace_plan.json`
- `evidence_ledger.json`
- `disconfirmation_log.json`
- `confidence_report.json`
- agent review rules in `agents/*.md`

The non-goal below was true for this earlier design stage. It has since been superseded by the `researcher` CLI design.

## Background

The current `industry-research` skill has a strong multi-agent report structure, but its judgment still starts too often from public internet narratives. This creates weak reports for vertical domains such as restaurant supply chain, retail expansion, warehousing, franchise operations, and cold-chain logistics.

The upgraded skill should behave more like a senior field consultant. It should first model the physical business system, then infer which operating traces must exist, then verify those traces across independent sources. Public articles and official claims can provide leads, but they should not be treated as operating facts without triangulation.

The first vertical target is restaurant, retail, and supply chain research. The design should remain general enough to become a future cross-industry method, but the first implementation should optimize for this domain.

## Goal

Upgrade the skill from a general industry research workflow into a restaurant, retail, and supply chain expert workflow that can:

- Decompose a company or business model into minimum operating units.
- Infer the required physical, human, legal, and digital inputs for those units.
- Map each required input to likely data exhaust.
- Use independent evidence triangulation before making high-confidence claims.
- Treat data conflicts as analytical signals, not simple errors.
- Produce reports with clear judgment levels: verified fact, high-confidence inference, explainable anomaly, suspended judgment, and unverified narrative.

## Non-Goals

This design does not add paid API integrations, scraping automation, or a new executable runner. It upgrades the skill instructions, agent contracts, reference material, report requirements, and evaluation expectations.

This design does not attempt to cover every industry. It deliberately starts with restaurant, retail, warehousing, distribution, franchise, and supply chain scenarios.

## Core Concept

The new workflow adds a mandatory stage called `Business Physics Modeling`.

Instead of moving directly from domain grounding to a ghost deck, the skill first asks:

1. What is the smallest real-world operating unit in this case?
2. What must physically exist for this unit to operate?
3. What legal, human, logistics, digital, and customer-facing traces must those inputs leave?
4. Which three independent evidence families can verify or falsify the claim?
5. If evidence conflicts, what business model change could explain the conflict?

The canonical chain is:

```text
research question
→ minimum operating units
→ required operating inputs
→ unavoidable data exhaust
→ triangulation tests
→ anomaly resolution
→ confidence judgment
```

## New Artifacts

### `entity_evidence_plan.json`

This file is produced before `ghost_deck.json`. It becomes the evidence blueprint for all downstream agents.

Required schema:

```json
{
  "research_question": "string",
  "business_model_hypotheses": [
    {
      "hypothesis": "直营/加盟/联营/区域代理/第三方仓配/中央厨房/前置仓等",
      "why_plausible": "string",
      "what_would_confirm_it": ["string"],
      "what_would_disconfirm_it": ["string"]
    }
  ],
  "minimum_operating_units": [
    {
      "unit_type": "store|warehouse|distribution_center|franchisee|vehicle_fleet|central_kitchen|regional_agent|supplier|digital_node",
      "unit_description": "string",
      "required_inputs": [
        {
          "input_type": "capital_legal|people_org|physical_fulfillment|digital_frontend|customer_employee_feedback",
          "input_description": "string",
          "expected_data_exhaust": ["string"]
        }
      ]
    }
  ],
  "triangulation_tests": [
    {
      "claim_to_test": "string",
      "evidence_family_1": "string",
      "evidence_family_2": "string",
      "evidence_family_3": "string",
      "minimum_confidence_rule": "string"
    }
  ],
  "anomaly_resolution_rules": [
    {
      "conflict_pattern": "string",
      "likely_explanations": ["string"],
      "next_best_checks": ["string"]
    }
  ]
}
```

### `references/restaurant-retail-supply-chain-physics.md`

This reference file should define the domain-specific business physics. It should be written as a consultant playbook, not as generic research advice.

Required sections:

- Minimum operating units: store, warehouse, distribution center, central kitchen, franchisee, regional agent, supplier, fleet, digital node.
- Required inputs by unit: site, license, staff, equipment, vehicle, route, SKU, supplier, system node, delivery radius, franchise contract.
- Data exhaust map:
  - Capital/legal: Tianyancha/Qichacha-style registration, branch status, social security headcount, mortgages, bidding, food license, warehouse lease clues.
  - People/org: Boss/Zhipin, Liepin, Zhaopin, Maimai, role mix, city distribution, lifecycle signal from job titles.
  - Physical fulfillment: map POI, delivery radius, logistics tender, cold-chain fleet requirements, warehouse address, truck route constraints.
  - Digital frontend: mini-program city list, app store selection, LBS nearest node, official account, store locator, SKU availability by location.
  - Feedback/UGC: Dianping/Meituan, Xiaohongshu, Douyin, complaint keywords, employee complaints, franchisee disputes.
- Domain heuristics:
  - Expansion leaves hiring and site traces before revenue traces.
  - Claimed coverage without people, license, route, or frontend traces is usually a weak claim.
  - Franchise and outsourced logistics can explain missing direct hiring.
  - Cold-chain and fresh food fulfillment are constrained by time, temperature, loading frequency, and route density.
  - Store count, warehouse count, and coverage radius should be treated as separate claims.

### `references/evidence-triangulation-playbook.md`

This reference file should define evidence levels, triangulation rules, and anomaly handling.

Required sections:

- Evidence hierarchy:
  - Operating fact traces: registration, license, social security, hiring, POI, delivery range, tender, UGC.
  - Strong leads: company official account, app/mini-program frontend, investor material with verifiable claims.
  - Weak leads: media coverage, awards, high-level market reports, unverified rankings.
- Rule of Three:
  - High confidence requires at least three independent evidence families.
  - Medium confidence requires two independent evidence families plus a plausible operating mechanism.
  - Low confidence applies to single-source claims.
- Conflict interpretation:
  - Official claim exists, no hiring or legal trace: possible future plan, ghost node, franchise, outsourcing, or narrative inflation.
  - Registration exists, no hiring: possible shell entity, pre-opening node, third-party operation, or low-staff franchise model.
  - Hiring exists, no official launch: possible pre-launch expansion or replacement hiring.
  - Frontend city exists, no physical trace: possible waitlist, service via remote warehouse, test market, or stale configuration.
- Confidence scoring rules:
  - Do not score by number of search results.
  - Score by independence, physical necessity, time alignment, and falsifiability.

## Existing File Changes

### `SKILL.md`

Add `Business Physics Modeling` between current domain grounding and multi-agent core workflow.

Revised flow:

1. Domain grounding.
2. Business physics modeling.
3. Evidence-grounded ghost deck.
4. Blue/Red first round.
5. Optional second round.
6. Chief arbitration.
7. Validation and delivery.

The skill should require `entity_evidence_plan.json` before `ghost_deck.json`. If the plan is missing or too generic, the workflow should stop and regenerate the plan.

The source constraints should be revised so that official claims and media reports are treated as leads, while operating traces receive higher evidentiary weight.

### `agents/engagement-manager.md`

The Engagement Manager should become responsible for two outputs:

1. `entity_evidence_plan.json`
2. `ghost_deck.json`

The ghost deck action titles must be evidence-testable operating claims, not generic market claims.

Good action title:

```text
华东冷链仓配能力若真实覆盖 300km 日配半径，应同时出现仓库点位、司机/库管招聘、区域配送触点和终端履约反馈。
```

Bad action title:

```text
华东市场空间巨大。
```

### `agents/blue-team.md`

Blue Team should prove that operating capability exists. It should not rely primarily on market size or policy support.

For each claim, it must try to establish at least three of the following:

- Legal/capital skeleton.
- People/organization muscle.
- Physical fulfillment blood flow.
- Digital frontend touchpoint.
- Customer/employee/franchisee feedback.

Blue Team output should add:

```json
{
  "operating_fact_chain": [
    {
      "evidence_family": "capital_legal|people_org|physical_fulfillment|digital_frontend|feedback_ugc",
      "observed_trace": "string",
      "why_this_trace_must_exist": "string",
      "source": "string",
      "confidence": 0.0
    }
  ]
}
```

### `agents/red-team.md`

Red Team should specialize in water-squeezing and anomaly discovery.

It should search for:

- Ghost nodes.
- Claimed coverage without staff or legal trace.
- Franchise or agent substitution for direct operations.
- Third-party logistics replacing owned fulfillment.
- Stale mini-program or app city lists.
- Hiring signals that contradict official expansion narratives.
- UGC complaints about late delivery, store closure, franchise dispute, refund, food safety, or wage arrears.

Red Team output should add:

```json
{
  "anomaly_tests": [
    {
      "claim_under_attack": "string",
      "missing_or_conflicting_trace": "string",
      "most_likely_explanations": ["string"],
      "severity": "CRITICAL|MAJOR|MINOR",
      "next_check": "string"
    }
  ]
}
```

### `agents/chief-arbitrator.md`

Chief Arbitrator should add an entity-map verdict layer.

Final judgments should be classified as:

- `VERIFIED_OPERATING_FACT`
- `HIGH_CONFIDENCE_INFERENCE`
- `EXPLAINABLE_ANOMALY`
- `SUSPENDED_JUDGMENT`
- `UNVERIFIED_NARRATIVE`

The final report should include a new section:

```text
实体经营版图与证据链
```

This section should include:

- Minimum operating units.
- Evidence triangulation matrix.
- Anomaly list and likely explanations.
- Confidence judgment by claim.

### `references/report-template.md`

Add the new report section for standard and comprehensive reports. Brief reports should include a compressed version in the executive summary or core arguments.

### `scripts/validate_report.py`

Add optional checks for restaurant/retail/supply-chain mode:

- Report contains `实体经营版图` or `Operating Footprint`.
- Report includes at least one triangulation matrix.
- Report includes at least one confidence classification from the new verdict set.
- Warnings if report contains only official/media sources without operating traces.

## Evaluation Updates

Add eval cases focused on vertical expertise:

1. Restaurant chain claimed lower-tier city expansion.
2. Fresh food retail brand claimed national cold-chain coverage.
3. Franchise tea brand claimed high store count and supply chain maturity.
4. Central kitchen supplier claimed regional delivery capability.

Expected outputs should check for:

- `entity_evidence_plan.json` exists.
- Minimum operating units are identified.
- At least three evidence families are used for high-confidence claims.
- Anomalies are interpreted instead of ignored.
- Official claims are not accepted without operating traces.
- Final report distinguishes verified facts from unverified narratives.

## Error Handling

If evidence is unavailable, the report must say what data would resolve the claim. It should not replace missing operating traces with generic market commentary.

If evidence conflicts, the report must provide the most plausible business explanations and the next checks required.

If a claim only has one source, it must be marked as low confidence or unverified narrative.

## Acceptance Criteria

The redesign is successful if a non-domain-expert analyst using the skill can:

- Start from operating units rather than media narratives.
- Know which evidence families to search for in restaurant, retail, and supply chain cases.
- Explain why missing evidence matters.
- Interpret contradictory data as business model signals.
- Produce a report that a 10+ year domain expert would recognize as field-aware rather than desk-research-heavy.

## Open Implementation Notes

The first implementation should edit instructions and references only. Programmatic data ingestion should come later after the evidence model is stable.

The design should preserve the existing red/blue/arbitrator architecture, but shift its center of gravity from market research to operating reality reconstruction.
