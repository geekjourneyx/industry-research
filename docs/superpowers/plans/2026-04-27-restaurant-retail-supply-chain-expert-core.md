# Restaurant Retail Supply Chain Expert Core Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Upgrade the `industry-research` skill so restaurant, retail, and supply chain reports start from operating reality, evidence exhaust, and triangulated judgment instead of generic internet research.

**Architecture:** Add a mandatory `Business Physics Modeling` stage before ghost deck generation. The Engagement Manager produces `entity_evidence_plan.json`, Blue/Red teams verify or attack operating traces, and the Chief Arbitrator classifies final claims by evidence strength. Two new reference playbooks provide the vertical domain logic.

**Tech Stack:** Markdown skill files, Python validator script, JSON intermediate artifacts, existing agent prompt architecture.

---

## File Structure

- Create `references/restaurant-retail-supply-chain-physics.md`
  - Owns vertical business physics for stores, warehouses, central kitchens, franchisees, fleets, suppliers, and digital nodes.
- Create `references/evidence-triangulation-playbook.md`
  - Owns evidence hierarchy, rule-of-three confidence logic, and conflict interpretation.
- Modify `SKILL.md`
  - Adds the new Business Physics Modeling stage, `entity_evidence_plan.json`, revised source hierarchy, and new delivery requirements.
- Modify `agents/engagement-manager.md`
  - Makes the Engagement Manager produce the entity evidence plan before the ghost deck.
- Modify `agents/blue-team.md`
  - Adds operating fact chain requirements and raises the bar for high-confidence positive claims.
- Modify `agents/red-team.md`
  - Adds anomaly tests and water-squeezing logic for inflated operating claims.
- Modify `agents/chief-arbitrator.md`
  - Adds entity-map verdicts and final claim classification.
- Modify `references/report-template.md`
  - Adds the `实体经营版图与证据链` section.
- Modify `scripts/validate_report.py`
  - Adds optional restaurant/retail/supply-chain checks.
- Modify `evals/evals.json`
  - Adds vertical expertise evaluation cases.

Because this directory is not a git repository, each task ends with a verification step instead of a commit step.

---

### Task 1: Add Restaurant/Retail/Supply Chain Business Physics Reference

**Files:**
- Create: `references/restaurant-retail-supply-chain-physics.md`

- [ ] **Step 1: Create the reference file**

Use `apply_patch` to create `references/restaurant-retail-supply-chain-physics.md` with this content:

```markdown
# 餐饮、零售、供应链商业物理参考手册

> 本文件定义餐饮、零售、供应链调研中的“商业物理学”。后续 Agent 必须先从真实经营单元出发，再搜索资料和写判断。

## 核心原则

不要先相信公司说了什么，先判断它为了做到这件事，在现实中必须拥有什么。

餐饮、零售、供应链业务都受物理约束。门店、仓、车、人、证照、系统节点、加盟商、供应商、履约半径和用户反馈不会凭空存在。只要业务真实运行，就会留下经营痕迹。

## 最小经营单元

| 单元 | 定义 | 常见研究问题 |
|------|------|--------------|
| 门店 | 直营店、加盟店、联营店、快闪店或档口 | 门店是否真实运营、是否持续开业、直营/加盟比例 |
| 前置仓 | 靠近终端消费者或门店的小型履约仓 | 是否支撑即时零售、是否真实覆盖某区域 |
| 配送中心 | 区域仓、总仓、冷链仓、常温仓 | 是否支撑区域扩张、履约半径是否合理 |
| 中央厨房 | 统一加工、预制、分装的生产节点 | 是否支撑门店标准化、食品安全和冷链能力 |
| 加盟商 | 承担本地经营的合作主体 | 扩张是否依赖加盟、加盟稳定性如何 |
| 区域代理 | 承担招商、运营、仓配或售后的一层组织 | 是否存在区域外包、管理半径是否过大 |
| 车队 | 自有、外包或合作运力 | 配送时效、冷链能力、成本结构是否成立 |
| 供应商 | 原料、包材、设备、仓配、系统供应商 | 供应链集中度、议价能力、履约稳定性 |
| 数字节点 | 小程序、App、门店选择、LBS、会员系统、订单入口 | 前端覆盖是否对应真实服务能力 |

## 必需生产要素

| 单元 | 必需生产要素 | 如果真实存在，通常会留下的痕迹 |
|------|--------------|-------------------------------|
| 门店 | 选址、租赁、装修、食品经营许可、店员、督导、收银系统、配送关系 | 地图 POI、点评/美团门店页、招聘、证照、消费者评价、公众号开业信息 |
| 前置仓 | 仓址、货架、库管、拣货员、骑手或车队、库存系统、配送半径 | 仓库 POI、库管/分拣招聘、小程序配送范围、即时零售入口、员工反馈 |
| 配送中心 | 仓库、装卸、库管、司机、冷链设备、干支线线路、WMS/TMS | 招聘、招投标、仓储地址、供应商招募、物流合作公告、地图点位 |
| 中央厨房 | 食品生产许可、设备、品控、研发、采购、冷链发运 | 许可信息、工厂地址、生产岗位、供应商公告、食品安全记录 |
| 加盟商 | 合同、培训、督导、首批物料、保证金、门店营运 | 招商信息、加盟商主体、纠纷、投诉、区域门店密度 |
| 区域代理 | 区域授权、招商人员、本地服务、仓配或售后 | 本地招聘、代理商工商主体、招商会、区域公众号或社群痕迹 |
| 车队 | 司机、车辆、温控设备、调度、线路、保险 | 司机招聘、物流招标、货运平台线索、车辆/冷链要求 |
| 供应商 | 采购合同、交付能力、质量认证、区域服务能力 | 招投标、供应商名录、合作公告、质量认证、诉讼或投诉 |
| 数字节点 | 城市配置、门店/仓绑定、SKU 库存、会员和下单入口 | 小程序城市列表、App 门店选择、LBS 最近节点、SKU 可售区域 |

## 数据废气地图

### 资本与法务

优先查：工商登记、分支机构、经营状态、参保人数、食品经营许可、食品生产许可、动产抵押、招投标、司法纠纷。

判断逻辑：

- 分支机构和许可证明“骨架”存在，但不证明真实运营。
- 参保人数能挤掉空壳主体水分。
- 食品、仓储、冷链相关许可与实际业务不匹配时，应降低置信度。

### 人力与组织

优先查：Boss 直聘、猎聘、智联、脉脉、地方招聘号、岗位城市分布。

判断逻辑：

- 扩张前通常先出现选址、拓展、招商、仓配、督导岗位。
- 稳定运营期通常出现店员、库管、司机、区域运营、品控岗位。
- 收缩期可能出现岗位消失、离职反馈、工资拖欠、加盟纠纷。

### 物理履约

优先查：高德/百度地图 POI、门店定位、仓库定位、配送范围、物流招投标、冷链要求、供应商招募。

判断逻辑：

- 冷链、生鲜、中央厨房履约受时间、温度、装卸频次和路线密度约束。
- 声称远距离覆盖时，必须解释中转仓、二级仓、第三方冷链或区域代理。
- 门店数、仓数和覆盖城市数是三个不同事实，不能互相替代。

### 数字前端

优先查：小程序城市列表、App 门店选择、LBS 最近门店、SKU 可售区域、公众号菜单、外卖平台入口。

判断逻辑：

- 前端城市列表通常是覆盖口径上限，可能包含筹备、测试或停运城市。
- 前端可下单、可选择门店、可显示配送范围，比单纯城市列表更强。
- SKU 在不同区域不可售，可能暴露仓配和供应链能力边界。

### 终端反馈

优先查：大众点评、美团、小红书、抖音、黑猫投诉、社媒评论、员工评价、加盟商纠纷。

判断逻辑：

- 用户反馈能验证真实履约体验。
- 员工和加盟商反馈能验证组织压力、欠薪、闭店、招商夸大。
- 负面反馈不是直接结论，必须结合时间、地点和业务模式解释。

## 领域判断规则

1. 官方宣发只作为线索，不作为经营事实。
2. 一个城市是否“真实覆盖”，至少要区分前端可见、可下单、可履约、稳定运营四个层级。
3. 一个仓是否“真实运营”，至少要看到地址、岗位、履约触点或供应链合作中的两类以上痕迹。
4. 直营扩张通常留下直接招聘和主体痕迹；加盟扩张通常留下招商、培训、纠纷和本地主体痕迹。
5. 第三方仓配会削弱直接招聘痕迹，但应出现物流合作、供应商、配送范围或终端履约痕迹。
6. 数据缺失不是自动否定，必须先判断是否存在加盟、代理、外包、筹备、停运、口径错配。
7. 数据冲突通常比一致数据更有价值，因为它可能暴露真实商业模式。
```

- [ ] **Step 2: Verify file exists and contains key sections**

Run:

```bash
rg -n "最小经营单元|必需生产要素|数据废气地图|领域判断规则" references/restaurant-retail-supply-chain-physics.md
```

Expected: four matching section headings.

---

### Task 2: Add Evidence Triangulation Playbook

**Files:**
- Create: `references/evidence-triangulation-playbook.md`

- [ ] **Step 1: Create the playbook**

Use `apply_patch` to create `references/evidence-triangulation-playbook.md` with this content:

```markdown
# 证据三角验证手册

> 本文件定义证据等级、三角验证规则、异常归因和置信度评分。后续 Agent 必须按本手册判断证据强弱。

## 证据等级

### A 级：经营事实痕迹

这些证据直接来自真实经营活动留下的痕迹，优先级最高：

- 工商主体、分支机构、经营状态、参保人数
- 食品经营许可、食品生产许可、仓储或冷链资质
- 招聘岗位、岗位城市、岗位类型、招聘时间
- 地图 POI、门店定位、仓库定位
- 小程序/App 可下单、可选门店、LBS 最近节点、SKU 可售区域
- 物流招投标、供应商招募、冷链车队要求
- 消费者评价、员工反馈、加盟商纠纷

### B 级：强线索

这些证据能提供方向，但需要 A 级证据验证：

- 公司公众号、官网、投资者材料
- 小程序城市列表、App 城市列表
- 招商会、供应商大会、官方开业新闻
- 管理层访谈中可被验证的具体经营数字

### C 级：弱线索

这些证据不能单独支撑经营事实：

- 媒体报道
- 市场规模报告
- 榜单、奖项、行业热度文章
- 未说明口径的门店数、城市数、GMV、覆盖范围

## 三角验证规则

高置信度结论必须至少满足三个独立证据家族中的三类：

1. 资本与法务：工商、许可、参保、招投标、诉讼。
2. 人力与组织：招聘、岗位、员工反馈、组织结构。
3. 物理履约：地图点位、仓配、车队、配送范围、供应商。
4. 数字前端：小程序、App、LBS、门店选择、SKU 可售区域。
5. 终端反馈：用户评价、社媒、投诉、加盟商反馈。

评分规则：

- 高置信度：三类独立证据互相支持，且时间、地点、主体口径一致。
- 中置信度：两类独立证据互相支持，并且存在合理经营机制解释。
- 低置信度：只有单一证据，或证据来自同一叙事源。
- 悬置判断：关键证据冲突，且无法通过经营模式解释。
- 不可采信：只有官方宣发、媒体转述或市场报告，没有经营痕迹。

## 常见冲突与解释

| 冲突模式 | 可能解释 | 下一步检查 |
|----------|----------|------------|
| 官方称已覆盖城市，但无招聘和工商痕迹 | 筹备城市、幽灵节点、加盟代理、第三方外包、宣传口径膨胀 | 查小程序可下单、地图点位、加盟商主体、物流合作 |
| 有工商主体，但无招聘 | 空壳主体、预开业、低人力加盟模式、第三方托管 | 查参保人数、证照、地图、员工反馈 |
| 有招聘，但无官方发布 | 预扩张、替换招聘、试点城市、外包团队补员 | 查岗位类型、岗位时间、地址、业务入口 |
| 有小程序城市，但无法下单 | 城市列表残留、测试配置、服务暂停、远程仓不可履约 | 查 SKU 可售、门店选择、配送范围、用户反馈 |
| 有门店 POI，但点评长期无新增评价 | 停业、低客流、加盟商挂牌、地图未更新 | 查外卖平台、社媒近照、营业时间、招聘 |
| 有仓库地址，但无库管/司机招聘 | 第三方仓配、轻资产代理、临时仓、地址挂靠 | 查物流招标、供应商、冷链合作、地图街景 |

## 反证测试

每个核心结论都必须回答：

- 什么证据会证明这个结论是错的？
- 这些反证数据是否已经尝试查找？
- 如果没有找到反证，是因为结论坚实，还是因为数据源不可得？

## 输出要求

报告中的每个关键判断必须被标为：

- `VERIFIED_OPERATING_FACT`：三类以上独立经营痕迹支持。
- `HIGH_CONFIDENCE_INFERENCE`：两类以上证据支持，经营机制合理。
- `EXPLAINABLE_ANOMALY`：证据冲突，但存在清晰商业解释。
- `SUSPENDED_JUDGMENT`：证据冲突且无法解释。
- `UNVERIFIED_NARRATIVE`：只有宣发、媒体或市场报告，缺少经营痕迹。
```

- [ ] **Step 2: Verify verdict labels**

Run:

```bash
rg -n "VERIFIED_OPERATING_FACT|HIGH_CONFIDENCE_INFERENCE|EXPLAINABLE_ANOMALY|SUSPENDED_JUDGMENT|UNVERIFIED_NARRATIVE" references/evidence-triangulation-playbook.md
```

Expected: all five verdict labels are present.

---

### Task 3: Update Main Skill Workflow

**Files:**
- Modify: `SKILL.md`

- [ ] **Step 1: Add new artifact to running contract**

In `SKILL.md`, find the list under “立即创建 `workspace_dir`，并在阶段一结束前落盘：” and keep the existing two files. After that list, add this paragraph:

```markdown
对于餐饮、零售、供应链相关请求，在进入幽灵卡片生成前，还必须落盘：

- `entity_evidence_plan.json`

该文件是经营实体证据计划。没有它，不允许进入 `ghost_deck.json` 生成。
```

- [ ] **Step 2: Add source hierarchy for operating traces**

In “工具与来源约束”, after the existing source priority table, add:

```markdown
### 餐饮/零售/供应链的经营痕迹优先级

当调研对象涉及餐饮、零售、门店、加盟、仓配、冷链、中央厨房、前置仓、即时零售或供应链履约时，必须优先使用经营痕迹，而不是市场叙事。

证据权重从高到低：

1. 经营事实痕迹：工商/许可/参保/招聘/地图 POI/小程序可下单/LBS/招投标/用户与员工反馈
2. 强线索：官网、公众号、小程序城市列表、投资者材料中可验证的具体经营声明
3. 弱线索：媒体报道、市场规模报告、榜单、未说明口径的覆盖城市或门店数

官方宣发和媒体报道只能作为线索，不能单独支撑高置信度经营结论。
```

- [ ] **Step 3: Insert Business Physics Modeling stage**

Before current `#### Step 2.1：编排 Agent — 构建幽灵卡片`, insert this section:

```markdown
#### Step 2.0：商业物理建模与实体证据计划

当研究对象涉及餐饮、零售、供应链、门店扩张、加盟、仓配、冷链、中央厨房、前置仓或区域履约时，必须先执行本步骤。

读取：

- `references/restaurant-retail-supply-chain-physics.md`
- `references/evidence-triangulation-playbook.md`

调用 Engagement Manager Agent，要求其先产出 `entity_evidence_plan.json`，再产出 `ghost_deck.json`。

`entity_evidence_plan.json` 必须包含：

```json
{
  "research_question": "核心研究问题",
  "business_model_hypotheses": [
    {
      "hypothesis": "直营/加盟/联营/区域代理/第三方仓配/中央厨房/前置仓等",
      "why_plausible": "为什么该模式可能存在",
      "what_would_confirm_it": ["确认该模式的证据"],
      "what_would_disconfirm_it": ["推翻该模式的证据"]
    }
  ],
  "minimum_operating_units": [
    {
      "unit_type": "store|warehouse|distribution_center|franchisee|vehicle_fleet|central_kitchen|regional_agent|supplier|digital_node",
      "unit_description": "经营单元说明",
      "required_inputs": [
        {
          "input_type": "capital_legal|people_org|physical_fulfillment|digital_frontend|customer_employee_feedback",
          "input_description": "该经营单元必须具备的投入",
          "expected_data_exhaust": ["预期数据废气"]
        }
      ]
    }
  ],
  "triangulation_tests": [
    {
      "claim_to_test": "需要验证的经营命题",
      "evidence_family_1": "第一类独立证据",
      "evidence_family_2": "第二类独立证据",
      "evidence_family_3": "第三类独立证据",
      "minimum_confidence_rule": "达到中高置信度的条件"
    }
  ],
  "anomaly_resolution_rules": [
    {
      "conflict_pattern": "数据冲突模式",
      "likely_explanations": ["可能商业解释"],
      "next_best_checks": ["下一步验证动作"]
    }
  ]
}
```

质量门：

- 至少识别 3 类最小经营单元。
- 至少覆盖 4 类证据家族：资本法务、人力组织、物理履约、数字前端、终端反馈。
- 至少提出 3 个三角验证测试。
- 每个高置信度经营命题必须说明可证伪条件。

如果 `entity_evidence_plan.json` 缺失、过于泛化或只罗列新闻搜索关键词，必须重试 1 次；仍失败则由主 Agent 生成最小可用版本，并标记 `execution_mode = degraded_entity_mapping`。
```

- [ ] **Step 4: Update ghost deck prompt**

In current Step 2.1 prompt, add this input line:

```markdown
  - 实体证据计划：{entity_evidence_plan JSON，如不适用则传入 "NOT_APPLICABLE"}
```

Also add this instruction to the same prompt:

```markdown
  如果实体证据计划适用，幽灵卡片中的行动标题必须是可验证的经营命题，而不是泛泛的市场规模或行业热度判断。
```

- [ ] **Step 5: Update workspace tree**

In “工作区管理”, add `entity_evidence_plan.json` between `context_dictionary.json` and `ghost_deck.json`:

```text
    ├── entity_evidence_plan.json
```

- [ ] **Step 6: Update reference index**

In “参考文件索引”, add:

```markdown
| `references/restaurant-retail-supply-chain-physics.md` | 餐饮/零售/供应链商业物理规则 | Step 2.0 |
| `references/evidence-triangulation-playbook.md` | 证据等级、三角验证、异常归因 | Step 2.0, Step 2.4 |
```

- [ ] **Step 7: Verify main skill references**

Run:

```bash
rg -n "Business Physics|商业物理|entity_evidence_plan|经营痕迹|degraded_entity_mapping" SKILL.md
```

Expected: matches in the running contract, source constraints, Step 2.0, workspace tree, and reference index.

---

### Task 4: Update Engagement Manager Agent

**Files:**
- Modify: `agents/engagement-manager.md`

- [ ] **Step 1: Update final goal**

Replace the paragraph under `## 终局目标` with:

```markdown
针对用户的调研需求，先构建一份“实体证据计划（Entity Evidence Plan）”，再构建一份用于指导下级 Agent 寻找数据的“幽灵卡片（Ghost Deck）”JSON 大纲。

实体证据计划回答：真实经营能力如果存在，现实世界必须有哪些经营单元、生产要素和数据废气。

幽灵卡片回答：下级 Agent 应围绕哪些可验证经营命题展开红蓝对抗。
```

- [ ] **Step 2: Add Business Physics step before current Step 1**

Before current `### Step 1：提炼核心研究问题`, add:

```markdown
### Step 0：商业物理建模

如果调研对象涉及餐饮、零售、供应链、门店、加盟、仓配、冷链、中央厨房、前置仓或区域履约，必须先读取：

- `references/restaurant-retail-supply-chain-physics.md`
- `references/evidence-triangulation-playbook.md`

然后产出 `entity_evidence_plan.json`。

建模顺序：

1. 提出可能的商业模式假设：直营、加盟、联营、区域代理、第三方仓配、中央厨房、前置仓、平台撮合。
2. 拆出最小经营单元：门店、仓、配送中心、中央厨房、加盟商、区域代理、车队、供应商、数字节点。
3. 推导每个单元必需的生产要素：场地、证照、人力、设备、车辆、路线、SKU、系统节点、履约半径。
4. 为每个生产要素列出预期数据废气。
5. 为关键经营命题设计三角验证测试。
6. 为常见数据冲突设计异常归因规则。
```

- [ ] **Step 3: Add entity evidence output format**

Before existing `## 输出格式`, insert:

```markdown
## 输出文件 1：实体证据计划

如果本次研究适用商业物理建模，先保存以下 JSON 到 `{workspace_dir}/entity_evidence_plan.json`：

```json
{
  "research_question": "精炼后的核心研究问题",
  "business_model_hypotheses": [
    {
      "hypothesis": "直营/加盟/联营/区域代理/第三方仓配/中央厨房/前置仓等",
      "why_plausible": "为什么该模式可能存在",
      "what_would_confirm_it": ["确认该模式的证据"],
      "what_would_disconfirm_it": ["推翻该模式的证据"]
    }
  ],
  "minimum_operating_units": [
    {
      "unit_type": "store|warehouse|distribution_center|franchisee|vehicle_fleet|central_kitchen|regional_agent|supplier|digital_node",
      "unit_description": "经营单元说明",
      "required_inputs": [
        {
          "input_type": "capital_legal|people_org|physical_fulfillment|digital_frontend|customer_employee_feedback",
          "input_description": "该单元成立所必需的投入",
          "expected_data_exhaust": ["该投入必然或大概率留下的数据痕迹"]
        }
      ]
    }
  ],
  "triangulation_tests": [
    {
      "claim_to_test": "可验证的经营命题",
      "evidence_family_1": "第一类独立证据",
      "evidence_family_2": "第二类独立证据",
      "evidence_family_3": "第三类独立证据",
      "minimum_confidence_rule": "达到中高置信度的条件"
    }
  ],
  "anomaly_resolution_rules": [
    {
      "conflict_pattern": "数据冲突模式",
      "likely_explanations": ["可能商业解释"],
      "next_best_checks": ["下一步验证动作"]
    }
  ]
}
```
```

- [ ] **Step 4: Rename existing output section**

Change existing `## 输出格式` heading to:

```markdown
## 输出文件 2：幽灵卡片
```

- [ ] **Step 5: Add constraints**

Under `## 约束`, add:

```markdown
- 对餐饮、零售、供应链研究，严禁直接从市场规模、政策红利或媒体报道生成行动标题。必须先从最小经营单元和可验证经营命题出发。
- 行动标题必须能够被证据证伪。例如“某品牌真实覆盖华东 300km 日配半径”可以证伪；“市场空间广阔”不可证伪。
- 如果只能找到官方宣发或媒体报道，必须在实体证据计划中标记为 `UNVERIFIED_NARRATIVE` 的候选风险。
```

- [ ] **Step 6: Verify Engagement Manager update**

Run:

```bash
rg -n "实体证据计划|商业物理建模|entity_evidence_plan|UNVERIFIED_NARRATIVE" agents/engagement-manager.md
```

Expected: matches in final goal, Step 0, output section, and constraints.

---

### Task 5: Update Blue Team Agent

**Files:**
- Modify: `agents/blue-team.md`

- [ ] **Step 1: Add operating capability to final goal**

After the current final goal paragraph, add:

```markdown
对于餐饮、零售、供应链研究，你的核心任务不是证明“市场很大”，而是证明经营能力真实存在。你必须围绕实体证据计划中的最小经营单元，验证骨架、肌肉、血液、皮肤和反馈是否能形成闭环。
```

- [ ] **Step 2: Add operating trace step**

After current `### Step 1：逐个处理分配给你的行动标题`, add:

```markdown
### Step 1A：经营事实链验证

如果输入包含 `entity_evidence_plan`，每个正向论点至少尝试建立以下五类证据中的三类：

1. 资本与法务骨架：工商、许可、参保、分支机构、招投标。
2. 人力与组织肌肉：招聘、岗位城市、岗位类型、员工反馈。
3. 物理履约血液：门店/仓库 POI、配送范围、车队、冷链、供应商。
4. 数字前端皮肤：小程序、App、LBS、门店选择、SKU 可售区域。
5. 终端反馈回声：用户评价、投诉、员工反馈、加盟商反馈。

如果无法凑齐三类独立证据，不得标记为高置信度经营事实。
```

- [ ] **Step 3: Add output field**

Inside each item in the JSON output format, after `evidence_chain`, add:

```json
      "operating_fact_chain": [
        {
          "evidence_family": "capital_legal|people_org|physical_fulfillment|digital_frontend|feedback_ugc",
          "observed_trace": "观察到的经营痕迹",
          "why_this_trace_must_exist": "为什么真实经营必须或大概率留下该痕迹",
          "source": "来源URL或名称",
          "confidence": 0.0
        }
      ],
```

- [ ] **Step 4: Add constraints**

Under `## 约束`, add:

```markdown
- 餐饮、零售、供应链研究中，市场规模、政策红利、融资新闻不能单独支撑 `SUPPORTED`。必须出现经营事实链。
- 如果只有公司官方口径，没有第三方经营痕迹，状态应设为 `DATA_INSUFFICIENT_FOR_OPERATING_FACT`。
- 对每条经营痕迹，必须说明“为什么真实经营必须留下这个痕迹”。
```

- [ ] **Step 5: Verify Blue Team update**

Run:

```bash
rg -n "经营事实链|operating_fact_chain|DATA_INSUFFICIENT_FOR_OPERATING_FACT|骨架|肌肉|血液" agents/blue-team.md
```

Expected: all terms are present.

---

### Task 6: Update Red Team Agent

**Files:**
- Modify: `agents/red-team.md`

- [ ] **Step 1: Add water-squeezing mission**

After the current final goal paragraph, add:

```markdown
对于餐饮、零售、供应链研究，你还必须执行“挤水分”任务：识别幽灵节点、虚假覆盖、口径膨胀、加盟替代直营、第三方仓配替代自营、前端城市列表残留、招聘缺失、证照缺失、履约半径不合理等问题。
```

- [ ] **Step 2: Add anomaly discovery step**

After current `### Step 1：事前验尸（Pre-Mortem）`, add:

```markdown
### Step 1A：经营异常测试

如果输入包含 `entity_evidence_plan`，对每个经营命题检查以下异常：

1. 官方称覆盖，但无招聘、工商、证照、地图或履约痕迹。
2. 有小程序城市列表，但无法下单或缺少真实门店/仓配触点。
3. 有工商主体，但无参保、招聘、许可或终端反馈。
4. 有门店或仓库 POI，但评价停滞、营业状态异常或招聘消失。
5. 声称冷链远距离覆盖，但无二级仓、外包冷链、线路或时效解释。
6. 声称直营扩张，但证据更像加盟、代理或第三方托管。
```

- [ ] **Step 3: Add output field**

Inside each item in the JSON output format, after `antithesis`, add:

```json
      "anomaly_tests": [
        {
          "claim_under_attack": "被攻击的经营命题",
          "missing_or_conflicting_trace": "缺失或冲突的经营痕迹",
          "most_likely_explanations": ["最可能商业解释"],
          "severity": "CRITICAL|MAJOR|MINOR",
          "next_check": "下一步验证动作"
        }
      ],
```

- [ ] **Step 4: Add constraints**

Under `## 约束`, add:

```markdown
- 不要把数据缺失简单等同于造假。必须先判断是否存在加盟、代理、外包、筹备、停运或口径错配。
- 数据冲突必须输出最可能解释和下一步验证动作。
- 如果只能找到宣传口径但找不到经营痕迹，应明确攻击其证据等级，而不是泛泛质疑。
```

- [ ] **Step 5: Verify Red Team update**

Run:

```bash
rg -n "挤水分|经营异常测试|anomaly_tests|幽灵节点|口径错配" agents/red-team.md
```

Expected: all terms are present.

---

### Task 7: Update Chief Arbitrator Agent

**Files:**
- Modify: `agents/chief-arbitrator.md`

- [ ] **Step 1: Add entity-map arbitration goal**

After the current final goal paragraph, add:

```markdown
对于餐饮、零售、供应链研究，你还必须产出“实体经营版图裁决”：区分哪些是已验证经营事实，哪些是高可信推断，哪些是可解释异常，哪些必须悬置，哪些只是不可采信的叙事。
```

- [ ] **Step 2: Add evidence classification step**

After current `### Step 2：MAMV 裁决`, add:

```markdown
### Step 2A：实体经营证据裁决

如果输入包含 `entity_evidence_plan`，对每个关键经营命题给出以下分类之一：

- `VERIFIED_OPERATING_FACT`：三类以上独立经营痕迹支持，时间、地点、主体口径一致。
- `HIGH_CONFIDENCE_INFERENCE`：两类以上证据支持，且经营机制合理。
- `EXPLAINABLE_ANOMALY`：证据冲突，但存在清晰商业解释。
- `SUSPENDED_JUDGMENT`：证据冲突且无法解释，需要补充数据。
- `UNVERIFIED_NARRATIVE`：只有宣发、媒体或市场报告，缺少经营痕迹。

裁决时不要按搜索结果数量评分，而要按独立性、物理必要性、时间一致性、主体一致性和可证伪性评分。
```

- [ ] **Step 3: Add metadata fields**

In metadata JSON output, add this field after `synthesis_record`:

```json
  "entity_map_verdicts": [
    {
      "claim": "关键经营命题",
      "verdict": "VERIFIED_OPERATING_FACT|HIGH_CONFIDENCE_INFERENCE|EXPLAINABLE_ANOMALY|SUSPENDED_JUDGMENT|UNVERIFIED_NARRATIVE",
      "supporting_evidence_families": ["capital_legal", "people_org", "physical_fulfillment"],
      "anomaly_explanation": "如有异常，给出最可能解释",
      "next_check": "解除悬置或提高置信度所需验证动作"
    }
  ],
```

- [ ] **Step 4: Add final report requirement**

Under `### Step 4：撰写报告`, add:

```markdown
餐饮、零售、供应链研究必须新增“实体经营版图与证据链”章节。该章节至少包含：

1. 最小经营单元清单。
2. 三角验证矩阵。
3. 异常值清单与最可能解释。
4. 关键经营命题的裁决标签。
```

- [ ] **Step 5: Verify Chief Arbitrator update**

Run:

```bash
rg -n "实体经营版图|VERIFIED_OPERATING_FACT|entity_map_verdicts|三角验证矩阵|UNVERIFIED_NARRATIVE" agents/chief-arbitrator.md
```

Expected: all terms are present.

---

### Task 8: Update Report Template

**Files:**
- Modify: `references/report-template.md`

- [ ] **Step 1: Add operating footprint section**

After current `### 三、行业全景 (Industry Landscape)` section and before `### 四、核心论证 (Core Arguments)`, insert:

```markdown
### 四、实体经营版图与证据链 (Operating Footprint & Evidence Chain)
**餐饮、零售、供应链相关报告必须包含；标准版和深度版完整展开，精简版可压缩进执行摘要或核心论证**

- 最小经营单元：门店、仓、中央厨房、加盟商、区域代理、车队、供应商、数字节点等。
- 证据三角验证矩阵：每个关键经营命题对应至少三类证据家族。
- 异常值清单：数据互相打架的地方，以及最可能商业解释。
- 裁决标签：`VERIFIED_OPERATING_FACT` / `HIGH_CONFIDENCE_INFERENCE` / `EXPLAINABLE_ANOMALY` / `SUSPENDED_JUDGMENT` / `UNVERIFIED_NARRATIVE`。

示例表格：

| 经营命题 | 资本/法务 | 人力/组织 | 物理履约 | 数字前端 | 终端反馈 | 裁决 |
|----------|-----------|-----------|----------|----------|----------|------|
| 华东仓支持 300km 日配 | 仓储主体与许可 | 司机/库管招聘 | 仓库 POI 与线路 | 小程序可配送范围 | 用户履约反馈 | HIGH_CONFIDENCE_INFERENCE |
```

- [ ] **Step 2: Renumber following sections**

Change following headings:

```markdown
### 四、核心论证 (Core Arguments)
### 五、战略路线图 (Strategic Roadmap)
### 六、风险矩阵 (Risk Matrix)
### 七、数据溯源表 (Data Provenance)
### 八、附录 (Appendix)
```

to:

```markdown
### 五、核心论证 (Core Arguments)
### 六、战略路线图 (Strategic Roadmap)
### 七、风险矩阵 (Risk Matrix)
### 八、数据溯源表 (Data Provenance)
### 九、附录 (Appendix)
```

- [ ] **Step 3: Verify template update**

Run:

```bash
rg -n "实体经营版图|Operating Footprint|三角验证矩阵|HIGH_CONFIDENCE_INFERENCE|### 九、附录" references/report-template.md
```

Expected: all terms are present.

---

### Task 9: Update Validator

**Files:**
- Modify: `scripts/validate_report.py`

- [ ] **Step 1: Add verdict constants**

After `REQUIRED_SECTIONS_BY_DEPTH`, add:

```python
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
```

- [ ] **Step 2: Add function parameter**

Change:

```python
def validate_report(report_path: str, depth: str = "standard") -> dict:
```

to:

```python
def validate_report(report_path: str, depth: str = "standard", vertical: str | None = None) -> dict:
```

- [ ] **Step 3: Add vertical checks before return**

Before `return results`, add:

```python
    if vertical in {"restaurant-retail-supply-chain", "rrsc"}:
        has_operating_section = re.search(
            r"^#{1,3}\s+.*(实体经营版图|Operating Footprint|Evidence Chain)",
            content,
            re.MULTILINE | re.IGNORECASE,
        )
        if not has_operating_section:
            results["errors"].append(
                "Missing operating footprint section for restaurant/retail/supply-chain mode"
            )
            results["valid"] = False

        has_triangulation = re.search(r"三角验证|triangulation", content, re.IGNORECASE)
        if not has_triangulation:
            results["warnings"].append(
                "No triangulation matrix or triangulation discussion found"
            )

        verdict_count = sum(1 for verdict in OPERATING_VERDICTS if verdict in content)
        results["stats"]["operating_verdict_count"] = verdict_count
        if verdict_count == 0:
            results["warnings"].append(
                "No operating verdict labels found; expected at least one entity-map judgment"
            )

        trace_hits = [term for term in OPERATING_TRACE_TERMS if term.lower() in content.lower()]
        results["stats"]["operating_trace_term_count"] = len(set(trace_hits))
        if len(set(trace_hits)) < 3:
            results["warnings"].append(
                "Few operating trace terms found; report may rely too heavily on narrative sources"
            )
```

- [ ] **Step 4: Parse vertical CLI flag**

In `main()`, after depth parsing, add:

```python
    vertical = None
    if "--vertical" in sys.argv:
        vertical_idx = sys.argv.index("--vertical")
        if vertical_idx + 1 < len(sys.argv):
            vertical = sys.argv[vertical_idx + 1]
```

Change:

```python
    results = validate_report(report_path, depth)
```

to:

```python
    results = validate_report(report_path, depth, vertical)
```

- [ ] **Step 5: Update usage string**

Change:

```python
print("Usage: python validate_report.py <report_path> [--depth brief|standard|comprehensive]")
```

to:

```python
print("Usage: python validate_report.py <report_path> [--depth brief|standard|comprehensive] [--vertical restaurant-retail-supply-chain|rrsc]")
```

- [ ] **Step 6: Verify validator syntax**

Run:

```bash
python3 -m py_compile scripts/validate_report.py
```

Expected: no output and exit code 0.

---

### Task 10: Update Evaluations

**Files:**
- Modify: `evals/evals.json`

- [ ] **Step 1: Add restaurant chain expansion eval**

Append this object to the `evals` array:

```json
{
  "id": 4,
  "prompt": "帮我判断一个连锁餐饮品牌宣称已经覆盖100个下沉城市是否可信。重点看真实门店、加盟、中央厨房和冷链配送能力，不要只引用新闻稿。标准深度。",
  "expected_output": "一份标准深度调研报告，必须先构建实体证据计划，再围绕门店、加盟商、中央厨房、仓配和数字前端做三角验证。报告应区分真实覆盖、前端城市列表、筹备节点、加盟代理和不可采信宣传口径。",
  "files": [],
  "expectations": [
    "工作区目录中产出entity_evidence_plan.json",
    "报告包含实体经营版图与证据链章节",
    "报告至少使用资本法务、人力组织、物理履约、数字前端、终端反馈中的三类证据",
    "报告不把官方宣发或媒体报道单独作为高置信度经营事实",
    "报告包含VERIFIED_OPERATING_FACT或HIGH_CONFIDENCE_INFERENCE等裁决标签",
    "报告解释至少一个数据冲突或异常场景"
  ]
}
```

- [ ] **Step 2: Add fresh retail cold-chain eval**

Append this object to the `evals` array:

```json
{
  "id": 5,
  "prompt": "调研一个生鲜零售品牌声称全国冷链覆盖是否真实。请重点验证仓配节点、配送半径、司机/库管招聘、小程序可配送区域和用户履约反馈。brief即可。",
  "expected_output": "一份精简报告，核心不是市场规模，而是验证全国冷链覆盖是否有经营痕迹支撑。报告应指出哪些城市是真实履约，哪些只是前端配置或宣传口径。",
  "files": [],
  "expectations": [
    "报告识别仓、车、人、数字前端和用户反馈等最小经营要素",
    "报告包含至少一个三角验证测试",
    "报告说明冷链履约受时效、温控、路线密度约束",
    "报告对单一来源信息标低置信度",
    "报告包含SUSPENDED_JUDGMENT或UNVERIFIED_NARRATIVE用于不可验证声明"
  ]
}
```

- [ ] **Step 3: Add franchise tea brand eval**

Append this object to the `evals` array:

```json
{
  "id": 6,
  "prompt": "分析一个茶饮品牌宣称门店突破5000家、供应链成熟、加盟商盈利稳定是否可信。请用餐饮供应链专家视角，重点挤水分。",
  "expected_output": "一份调研报告，围绕门店数口径、加盟商主体、招商信息、供应链节点、员工和加盟商反馈做交叉验证，并明确哪些结论不能只靠品牌宣传采信。",
  "files": [],
  "expectations": [
    "报告区分直营、加盟、区域代理和第三方服务",
    "报告检查门店数、供应链能力和加盟商盈利是三个不同命题",
    "红方分析包含幽灵门店、加盟纠纷、供应链覆盖不足或宣传口径膨胀等异常测试",
    "仲裁结论区分已验证经营事实、高可信推断、可解释异常和悬置判断"
  ]
}
```

- [ ] **Step 4: Validate JSON**

Run:

```bash
python3 -m json.tool evals/evals.json >/tmp/industry-research-evals-check.json
```

Expected: no output and exit code 0.

---

### Task 11: Final Cross-File Verification

**Files:**
- Read-only verification across all changed files.

- [ ] **Step 1: Search required concepts**

Run:

```bash
rg -n "entity_evidence_plan|实体证据计划|实体经营版图|VERIFIED_OPERATING_FACT|UNVERIFIED_NARRATIVE|三角验证|经营痕迹" SKILL.md agents references scripts evals docs/superpowers/specs
```

Expected: matches across main skill, agents, references, validator, evals, and spec.

- [ ] **Step 2: Run validator syntax check**

Run:

```bash
python3 -m py_compile scripts/validate_report.py
```

Expected: no output and exit code 0.

- [ ] **Step 3: Run eval JSON check**

Run:

```bash
python3 -m json.tool evals/evals.json >/tmp/industry-research-evals-check.json
```

Expected: no output and exit code 0.

- [ ] **Step 4: Confirm no placeholder language**

Run:

```bash
rg -n "T[B]D|T[O]DO|PLACE[H]OLDER|待[定]|占[位]" SKILL.md agents references scripts evals docs/superpowers
```

Expected: no matches.

- [ ] **Step 5: Inspect changed files**

Run:

```bash
find . -maxdepth 3 -type f | sort
```

Expected: output includes the two new reference files, the approved spec, the expert collaboration document, and this implementation plan.

---

## Self-Review

Spec coverage:

- Business Physics Modeling stage: Task 3.
- `entity_evidence_plan.json`: Tasks 3 and 4.
- New reference playbooks: Tasks 1 and 2.
- Engagement Manager changes: Task 4.
- Blue Team operating fact chain: Task 5.
- Red Team anomaly tests: Task 6.
- Chief Arbitrator verdict layer: Task 7.
- Report template operating footprint section: Task 8.
- Validator checks: Task 9.
- Vertical evals: Task 10.
- Final verification: Task 11.

Scope check:

- This plan edits instructions, references, validation, and evals only.
- It does not add API integrations, scraping automation, or a new runner.
- The work remains focused on the restaurant, retail, and supply chain expert core.

Placeholder scan:

- The plan intentionally contains no unresolved placeholders.
- Every task lists exact target files, exact inserted content or replacement content, and exact verification commands.
