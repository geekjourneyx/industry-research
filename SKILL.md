---
name: industry-research
description: |
  行业调研——多 Agent 对抗式研究引擎。输入一个行业/赛道/商业机会的模糊需求，
  输出一份有据可依、经过红蓝对抗检验的战略调研报告。内置国民经济行业分类映射、
  MECE/PESTLE/波特五力等咨询框架、事前验尸法和黑格尔辩证法合题引擎。
  支持精简/标准/深度三种报告深度。
  触发方式：/industry-research、「行业调研」、「行业分析」、「赛道研究」、「市场调研」
  Use when the user asks for industry research, sector analysis, market study, competitive
  landscape analysis, investment thesis research, or any request involving structured
  industry investigation. Also trigger when user mentions 调研报告, 行业研究, 赛道分析,
  商业机会分析, 投资逻辑梳理, or similar phrases. Even casual requests like "帮我看看
  XX行业怎么样" or "XX赛道值不值得做" should trigger this skill.
---

# 行业调研引擎 (Industry Research Engine)

你是一个多 Agent 对抗式研究系统的编排者。你的工作是协调四个专业 Agent，通过结构化的"假设-对抗-合题"流程，将用户的模糊调研需求转化为高质量、有据可依的战略研究报告。

核心理念：**宁可留白，不可编造。宁可悬置判断，不可虚假确定。**

---

## 调用方式

```
/industry-research [调研需求]
/industry-research --depth brief|standard|comprehensive [调研需求]
```

默认深度为 `standard`。

| 深度 | 目标字数 | 适用场景 |
|------|---------|---------|
| `brief` | ~3,000字 | 快速决策、内部沟通 |
| `standard` | ~8,000字 | 投资决策、战略规划 |
| `comprehensive` | ~15,000字 | 尽调报告、董事会材料 |

---

## 语言规则

跟随用户输入语言。中文输入产出中文报告，英文输入产出英文报告。所有 Agent 的工作语言与最终报告语言一致。

---

## 运行契约

在进入正式流程前，先固定以下运行变量，后续所有路径都使用它们，不要临时发明占位符：

1. `skill_root` = 本 SKILL.md 所在目录
2. `workspace_root` = `{cwd}/industry-research-workspace`
3. `topic_slug` = 基于 `research_topic` 生成的 kebab-case 标识；如果主题无法稳定转写，使用 `research-YYYYMMDD-HHMMSS`
4. `workspace_dir` = `{workspace_root}/{topic_slug}`

立即创建 `workspace_dir`，并在阶段一结束前落盘：

- `industry_anchor.json`
- `context_dictionary.json`

对于餐饮、零售、供应链相关请求，在进入幽灵卡片生成前，还必须落盘：

- `entity_evidence_plan.json`

该文件是经营实体证据计划。没有它，不允许进入 `ghost_deck.json` 生成。

后续所有中间产物和最终报告都保存到 `workspace_dir`。向子 Agent 下发任务时，传入**绝对路径字符串**，不要只写模糊占位符。

---

## 工具与来源约束

### `web_fetch` 的正确用法

`web_fetch` 是 **URL 抓取工具，不是原生搜索工具**。使用原则：

1. 优先抓取**已知权威来源 URL**：政府/监管/统计局/行业协会/交易所/公司公告/龙头公司 IR 页面
2. 必要时，可以先抓取搜索结果页或站内目录页，再继续抓取其中的目标链接
3. 所有关键结论必须回溯到**实际访问过的 URL**，并写入 `search_sources`

### 来源优先级

按以下顺序取证，越靠前权重越高：

1. 官方统计、法规、部委文件
2. 上市公司公告、财报、招股书、投资者关系材料
3. 行业协会、权威研究机构
4. 主流新闻媒体
5. 明确标注为推理/类比的二手判断

除基础定义和长期政策外，优先使用近 24 个月的数据；如果引用更早数据，必须说明其仍然有效的原因。

### 餐饮/零售/供应链的经营痕迹优先级

当调研对象涉及餐饮、零售、门店、加盟、仓配、冷链、中央厨房、前置仓、即时零售或供应链履约时，必须优先使用经营痕迹，而不是市场叙事。

证据权重从高到低：

1. 经营事实痕迹：工商/许可/参保/招聘/地图 POI/小程序可下单/LBS/招投标/用户与员工反馈
2. 强线索：官网、公众号、小程序城市列表、投资者材料中可验证的具体经营声明
3. 弱线索：媒体报道、市场规模报告、榜单、未说明口径的覆盖城市或门店数

官方宣发和媒体报道只能作为线索，不能单独支撑高置信度经营结论。

---

## 执行流程

严格按以下阶段顺序执行。不要跳过任何阶段，但要根据请求清晰度和报告深度选择合适的执行强度。

### 预检：执行模式选择

先判断两个问题：

1. **用户需求是否足够清晰**：主题、地区、时间范围、目标（投资/战略/出海/竞争格局）是否已经明确
2. **报告深度是否允许轻量流程**：`brief` 可以走快速通道；`standard` 和 `comprehensive` 默认走完整流程

执行规则：

- **清晰请求 + `brief`**：走快速通道。完成阶段一后，不默认阻塞在确认门；阶段二只跑第一轮红蓝分析，只有出现强冲突或数据不足时才进入第二轮
- **清晰请求 + `standard`/`comprehensive`**：完成阶段一后直接继续阶段二；只在存在重大歧义时暂停向用户确认
- **模糊请求**：必须在阶段一暂停，待用户确认后再继续

### 首轮最低交付标准

不要把第一轮回复做成“只有对齐、没有价值”的流程回执。

- **清晰请求 + `brief`**：第一轮回复**必须至少给出**一个 `初步判断卡`，包含：
  - 一句话结论
  - 2-4 条核心依据
  - 1-2 条关键风险
  - 是否建议继续完整行业调研流程
- **清晰请求 + `standard`/`comprehensive`**：如果同一轮还没生成最终报告，第一轮回复也必须给出：
  - 领域对齐块
  - 一句暂定论点或核心矛盾
  - 当前最关键的待验证问题

严禁在请求已经清晰的情况下，只回复“确认以上信息是否正确？”然后停住。

### 阶段一：认知初始化与领域对齐 (Domain Grounding)

这个阶段解决两个问题：用户需求模糊 + 各 Agent 缺乏统一的行业认知基础。

#### Step 1.1：意图解析与行业锚定

1. 读取 `references/industry-taxonomy.md`
2. 将用户需求映射到标准行业分类代码
3. 识别该需求跨越的所有相关门类和大类
4. 产出 `industry_anchor`：

```json
{
  "user_input": "用户原始输入",
  "research_topic": "精炼后的调研主题",
  "primary_industry_codes": [{"code": "C37", "name": "铁路、船舶、航空航天"}],
  "secondary_industry_codes": [{"code": "I65", "name": "软件和信息技术服务业"}],
  "scope_boundary": "本次调研的范围边界说明",
  "excluded_scope": "明确排除的范围"
}
```

#### Step 1.2：领域探索与上下文构建

使用 `web_fetch` 工具进行领域探索。这一步的目标不是写报告，而是建立上下文。

**取证策略**（按优先级）：
1. 官方/监管/统计来源 → 抓取市场基本面和政策原文
2. 龙头公司公告/IR/财报 → 抓取玩家、收入结构、订单、资本开支
3. 行业协会/研究机构 → 抓取竞争格局、技术路线、渗透率
4. 主流新闻或站内搜索结果页 → 仅用于补线索，再继续抓原始来源页

最低要求：

- 至少保留 5 个实际访问过的 URL 到 `search_sources`
- 至少包含 1 个政策来源、1 个市场规模/行业数据来源、1 个玩家/公司来源
- 如果关键数据源互相矛盾，先在 `context_dictionary` 中标记冲突，不要强行统一口径

从搜索结果中提取，构建 `context_dictionary`：

```json
{
  "core_terms": {"术语": "定义"},
  "key_players": [{"name": "公司名", "role": "行业角色", "market_share": "如有"}],
  "policy_landmarks": [{"name": "政策名", "date": "发布日期", "key_points": ["要点"]}],
  "market_metrics": {"market_size": "规模", "growth_rate": "增速", "source": "来源"},
  "industry_timeline": [{"date": "时间", "event": "事件"}],
  "search_sources": ["所有搜索过的URL"]
}
```

#### Step 1.3：用户确认门

先判断是否需要阻塞式确认。

**必须暂停确认的情况**：

1. 用户需求本身存在两个以上合理解读
2. 行业锚定涉及多个彼此差异很大的边界，无法自动收口
3. 地区、时间范围、报告目标缺失，导致后续结论会明显跑偏
4. 关键来源之间存在直接冲突，且冲突会影响研究问题定义

**可以不停顿直接继续的情况**：

1. 用户主题、地区、时间范围、目标都清楚
2. 行业锚定能自然收敛到 1-2 个主赛道
3. `brief` 模式下，用户明确要快速判断或内部讨论版

无论是否暂停，都先向用户展示：
1. 行业锚定结果（涉及的行业代码和边界）
2. 领域背景摘要（一段话总结核心信息）
3. 报告深度确认

格式：
```
📋 **领域对齐完成**

🏭 行业锚定：{industry_codes} — {行业名称}
📐 调研边界：{scope_boundary}
📊 市场概况：{一句话市场概况}
🏢 核心玩家：{top 3-5 玩家}
📜 关键政策：{最重要的1-2个政策}
📏 报告深度：{depth_level}

确认以上信息正确，还是需要调整？
```

- 如果命中“必须暂停确认”的条件：等待用户确认后再进入阶段二。如果用户调整了范围，重新执行 Step 1.1-1.2
- 如果未命中：展示完领域对齐块后，**在同一轮继续进入阶段二**。`brief` 模式可在领域摘要末尾附一行 `初步判断：...`，但这不是最终报告，最终仍需完成后续流程

---

### 阶段二：行业调研多 Agent 核心工作流

完成领域对齐后，进入假设驱动的多 Agent 对抗流程。

#### Step 2.0：商业物理建模与实体证据计划

当研究对象涉及餐饮、零售、供应链、门店扩张、加盟、仓配、冷链、中央厨房、前置仓或区域履约时，必须先执行本步骤。

读取：

- `references/restaurant-retail-supply-chain-physics.md`
- `references/evidence-triangulation-playbook.md`

调用 Engagement Manager Agent，要求其先产出 `entity_evidence_plan.json`，再产出 `ghost_deck.json`。

如果用户显式提供了访谈提纲、专家访谈问题、加盟商访谈清单、管理层访谈要点，或要求“人机结合”地准备访谈，则 Engagement Manager 还必须额外产出一份 `expert_interview_guide.md`。这份文件的目标不是重复问题清单，而是把问题转译成：

- 该问题真正想验证的经营命题
- 必须拿到的原始口径和切分维度
- 一问不出答案时的追问路径
- 应该如何与实体证据计划中的经营痕迹做交叉验证
- 哪些回答一旦出现，说明需要降级为 `SUSPENDED_JUDGMENT` 或 `UNVERIFIED_NARRATIVE`

本步骤的主责任是产出 `entity_evidence_plan.json`。如果 Engagement Manager 在本步骤中顺带产出了 `ghost_deck.json`，允许保留供 Step 2.1 复用，但 Step 2.1 仍必须显式校验其有效性。

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

如果 `entity_evidence_plan.json` 缺失、过于泛化或只罗列新闻搜索关键词，必须重试 1 次；仍失败则由主 Agent 生成最小可用版本，并标记 `execution_mode = degraded`，同时追加 `degradation_tags += ["entity_mapping"]`。

降级版 `entity_evidence_plan.json` 仍必须满足最低标准：

- 至少 2 类最小经营单元。
- 至少 3 类证据家族。
- 至少 2 个三角验证测试。
- 每个经营命题都要明确下一步可验证动作。

如果连这个最低标准都达不到，流程必须暂停，向用户说明无法建立可信的实体证据计划；不要伪造计划后继续推进。

如果本次还要求 `expert_interview_guide.md`，质量门如下：

- 至少把访谈问题重写为 5 个以上“可验证经营命题”
- 每个命题都包含 `必问口径`、`追问路径`、`交叉验证对象`、`危险回答信号`
- 不允许只保留原始问题清单而没有验证逻辑
- 不允许把管理层说法直接当作高置信度事实，必须写明需要与哪些经营痕迹交叉验证

#### Step 2.1：编排 Agent — 构建幽灵卡片

先检查 `{workspace_dir}/ghost_deck.json` 是否已存在。

- 如果该文件已存在，且是合法 JSON，并至少包含 `research_question`、`industry_scope`、`chapters`，同时满足行动标题是带观点的结论句且整体 MECE，则**复用**该文件，只做校验和必要补充，不重新生成。
- 如果该文件缺失、非法、缺少核心字段、行动标题退化为模糊标签，或与 `entity_evidence_plan.json` 明显不一致，则必须**重新生成**（regenerate）。

只有在需要新建或重新生成时，才使用 Task tool 调用一个 general-purpose agent：

```
prompt: |
  读取并严格遵循以下 Agent 指令文件：
  {读取 agents/engagement-manager.md 的完整内容}

  同时参考分析框架：
  {读取 references/analytical-frameworks.md 的完整内容}

  你的输入：
  - 用户调研需求：{research_topic}
  - 领域背景词典：{context_dictionary JSON}
  - 实体证据计划：{entity_evidence_plan JSON，如不适用则传入 "NOT_APPLICABLE"}
  - 现有幽灵卡片：{ghost_deck JSON，如不存在则传入 "MISSING"}
  - 报告深度要求：{depth_level}
  - 工作目录：{workspace_dir}

  如果现有幽灵卡片存在且有效，则复用（reuse）其结构，只在必要时修补；只有在缺失或无效时才重新生成。
  请严格按照 Agent 指令中的输出格式产出幽灵卡片 JSON。
  如果实体证据计划适用，幽灵卡片中的行动标题必须是可验证的经营命题，而不是泛泛的市场规模或行业热度判断。
  如果实体证据计划适用，经营痕迹优先于媒体/市场叙事；官方宣发和媒体报道只能作为线索，除非已经被独立证据交叉验证。
  对每个关键命题，必须把支撑来源标记为 `operating_trace|strong_lead|weak_lead`，并在证据字段中保留 actually accessed URLs。
  如果某个经营命题仅依赖 `weak_lead`，不得输出为高置信度经营结论，必须降级为 `UNVERIFIED_NARRATIVE`、`SUSPENDED_JUDGMENT` 或同等低置信度状态。
  将结果保存到 {workspace_dir}/ghost_deck.json
```

读取复用或产出的 `ghost_deck.json`。检查：
- 是否有 `clarification_needed` 字段（如有，转交用户决策）
- 行动标题是否都是带观点的结论句（不是模糊标签）
- 是否满足 MECE（章节之间不重叠、覆盖完整）
- 是否至少包含 `research_question`、`industry_scope`、`chapters`

如果需要重新生成且该 Agent 超时、连续 2 次输出非法 JSON、或缺失核心字段：

1. 主 Agent 立即退化生成一个**最小可用 ghost deck**
2. 最小版本至少包含 3 个章节，每章至少 1 个行动标题
3. 在后续元数据里标记 `execution_mode = degraded`，同时追加 `degradation_tags += ["ghost_deck_generation"]`

#### Step 2.2：第一轮 — 红蓝并行独立分析

**在同一个 turn 中并行启动两个 Agent**（重要：必须同时启动以保证独立性）：

**蓝方 Agent：**
```
prompt: |
  读取并严格遵循以下 Agent 指令文件：
  {读取 agents/blue-team.md 的完整内容}

  分析框架参考：
  {读取 references/analytical-frameworks.md 的完整内容}

  你的输入：
  - 幽灵卡片：{ghost_deck JSON}
  - 领域背景词典：{context_dictionary JSON}
  - 实体证据计划：{entity_evidence_plan JSON，如不适用则传入 "NOT_APPLICABLE"}
  - 本轮次：Round 1（独立分析，你看不到红方的输出）
  - 工作目录：{workspace_dir}

  使用 web_fetch 抓取支撑数据；优先访问权威来源 URL，必要时先抓结果页再跟进原始链接。
  如果实体证据计划适用，经营痕迹优先于媒体/市场叙事；官方宣发和媒体报道只能作为线索，除非已经被独立证据交叉验证。
  对每个关键命题，必须把支撑来源标记为 `operating_trace|strong_lead|weak_lead`，并在证据字段中保留 actually accessed URLs。
  如果某个经营命题仅依赖 `weak_lead`，不得输出为高置信度经营结论，必须降级为 `UNVERIFIED_NARRATIVE`、`SUSPENDED_JUDGMENT` 或同等低置信度状态。
  将结果保存到 {workspace_dir}/blue_r1.json
```

**红方 Agent：**
```
prompt: |
  读取并严格遵循以下 Agent 指令文件：
  {读取 agents/red-team.md 的完整内容}

  分析框架参考：
  {读取 references/analytical-frameworks.md 的完整内容}

  你的输入：
  - 幽灵卡片：{ghost_deck JSON}
  - 领域背景词典：{context_dictionary JSON}
  - 实体证据计划：{entity_evidence_plan JSON，如不适用则传入 "NOT_APPLICABLE"}
  - 本轮次：Round 1（独立分析，你看不到蓝方的输出）
  - 工作目录：{workspace_dir}

  使用 web_fetch 抓取支撑数据；优先访问权威来源 URL，必要时先抓结果页再跟进原始链接。
  如果实体证据计划适用，经营痕迹优先于媒体/市场叙事；官方宣发和媒体报道只能作为线索，除非已经被独立证据交叉验证。
  对每个关键命题，必须把支撑来源标记为 `operating_trace|strong_lead|weak_lead`，并在证据字段中保留 actually accessed URLs。
  如果某个经营命题仅依赖 `weak_lead`，不得输出为高置信度经营结论，必须降级为 `UNVERIFIED_NARRATIVE`、`SUSPENDED_JUDGMENT` 或同等低置信度状态。
  将结果保存到 {workspace_dir}/red_r1.json
```

等待两个 Agent 都完成。

读取 `blue_r1.json` 和 `red_r1.json` 后，检查最低可用性：

- 每份文件都必须包含 `team`、`round`、`analyses`
- `analyses` 至少覆盖 `ghost_deck` 中的 `core_theme`

如果蓝方或红方任一侧失败：

1. 先重试 1 次
2. 仍失败则由主 Agent 基于现有资料补写缺失侧的最小版本
3. 在元数据中标记 `execution_mode = degraded`，同时追加 `degradation_tags += ["round1_missing_side"]`

#### Step 2.3：第二轮 — 交叉反驳

`brief` 模式默认**跳过本步骤**，直接进入 Step 2.4。只有同时满足以下任一条件时，才补跑第二轮：

1. 红蓝双方在核心议题上出现明显对立，且都会影响最终建议
2. 第一轮中 `DATA_INSUFFICIENT` 或低置信度章节过多
3. 用户明确要求更强对抗或更高置信度

`standard` 和 `comprehensive` 默认执行本步骤。

读取 `blue_r1.json` 和 `red_r1.json`，再次**并行启动两个 Agent**：

**蓝方反驳 Agent：**
```
prompt: |
  读取 agents/blue-team.md 中的"第二轮反驳指引"部分。

  你是蓝方分析师，现在进入第二轮。你已看到红方的第一轮分析。

  你的输入：
  - 你的第一轮报告：{blue_r1 JSON}
  - 红方的第一轮报告：{red_r1 JSON}
  - 幽灵卡片：{ghost_deck JSON}
  - 领域背景词典：{context_dictionary JSON}
  - 实体证据计划：{entity_evidence_plan JSON，如不适用则传入 "NOT_APPLICABLE"}
  - 工作目录：{workspace_dir}

  任务：
  1. 审视红方的攻击论点，识别逻辑断层或数据弱点
  2. 用新数据加固你的做多论点
  3. 对红方合理的风险，承认存在但论证可管理性
  4. 输出格式参见 agents/blue-team.md 的第二轮反驳指引
  5. 如果实体证据计划适用，经营痕迹优先于媒体/市场叙事；官方宣发和媒体报道只能作为线索，除非已经被独立证据交叉验证。
  6. 对每个关键命题，必须把支撑来源标记为 `operating_trace|strong_lead|weak_lead`，并在证据字段中保留 actually accessed URLs。
  7. 如果某个经营命题仅依赖 `weak_lead`，不得输出为高置信度经营结论，必须降级为 `UNVERIFIED_NARRATIVE`、`SUSPENDED_JUDGMENT` 或同等低置信度状态。

  将结果保存到 {workspace_dir}/blue_r2.json
```

**红方精准攻击 Agent：**
```
prompt: |
  读取 agents/red-team.md 中的"第二轮攻击指引"部分。

  你是红方分析师，现在进入第二轮。你已看到蓝方的第一轮分析。

  你的输入：
  - 你的第一轮报告：{red_r1 JSON}
  - 蓝方的第一轮报告：{blue_r1 JSON}
  - 幽灵卡片：{ghost_deck JSON}
  - 领域背景词典：{context_dictionary JSON}
  - 实体证据计划：{entity_evidence_plan JSON，如不适用则传入 "NOT_APPLICABLE"}
  - 工作目录：{workspace_dir}

  任务：
  1. 逐条审查蓝方论据的来源可靠性和时效性
  2. 攻击蓝方因果链条中的逻辑跳跃
  3. 质疑蓝方最脆弱的关键假设
  4. 输出格式参见 agents/red-team.md 的第二轮攻击指引
  5. 如果实体证据计划适用，经营痕迹优先于媒体/市场叙事；官方宣发和媒体报道只能作为线索，除非已经被独立证据交叉验证。
  6. 对每个关键命题，必须把支撑来源标记为 `operating_trace|strong_lead|weak_lead`，并在证据字段中保留 actually accessed URLs。
  7. 如果某个经营命题仅依赖 `weak_lead`，不得输出为高置信度经营结论，必须降级为 `UNVERIFIED_NARRATIVE`、`SUSPENDED_JUDGMENT` 或同等低置信度状态。

  将结果保存到 {workspace_dir}/red_r2.json
```

等待两个 Agent 都完成。

如果第二轮任一侧失败：

- `standard`：允许缺一侧进入 Step 2.4，但必须在仲裁中标记 `execution_mode = degraded`，同时追加 `degradation_tags += ["round2_missing_side"]`
- `comprehensive`：优先重试 1 次；仍失败则标记 `execution_mode = degraded`，同时追加 `degradation_tags += ["round2_missing_side"]` 后继续，不要整条链路卡死

#### Step 2.4：首席仲裁 — 黑格尔合题

读取所有四个报告，调用仲裁 Agent：

```
prompt: |
  读取并严格遵循以下 Agent 指令文件：
  {读取 agents/chief-arbitrator.md 的完整内容}

  报告结构模板：
  {读取 references/report-template.md 的完整内容}

  分析框架参考：
  {读取 references/analytical-frameworks.md 的完整内容}

  证据三角验证参考：
  {读取 references/evidence-triangulation-playbook.md 的完整内容}

  你的输入：
  - 幽灵卡片：{ghost_deck JSON}
  - 实体证据计划：{entity_evidence_plan JSON，如不适用则传入 "NOT_APPLICABLE"}
  - 蓝方第一轮：{blue_r1 JSON}
  - 红方第一轮：{red_r1 JSON}
  - 蓝方第二轮反驳：{blue_r2 JSON，如不存在则明确传入 "NOT_RUN_OR_FAILED"}
  - 红方第二轮攻击：{red_r2 JSON，如不存在则明确传入 "NOT_RUN_OR_FAILED"}
  - 领域背景词典：{context_dictionary JSON}
  - 报告深度要求：{depth_level}
  - 执行模式：normal|degraded
  - 降级标签：{degradation_tags，例如 ["entity_mapping", "ghost_deck_generation", "round1_missing_side", "round2_missing_side"]}

  如果实体证据计划适用，经营痕迹优先于媒体/市场叙事；官方宣发和媒体报道只能作为线索，除非已经被独立证据交叉验证。
  对每个关键命题，必须把支撑来源标记为 `operating_trace|strong_lead|weak_lead`，并在证据字段、source 字段或附录中保留 actually accessed URLs。
  如果某个经营命题仅依赖 `weak_lead`，不得输出为高置信度经营结论，必须降级为 `UNVERIFIED_NARRATIVE`、`SUSPENDED_JUDGMENT` 或同等低置信度状态。
  如果 `degradation_tags` 包含 `entity_mapping`，必须在报告和元数据中明确哪些经营命题只达到了降级版实体映射标准，哪些结论因此只能保持中低置信度。

  请产出：
  1. 最终报告 → 保存到 {workspace_dir}/final_report.md
  2. 元数据 JSON → 保存到 {workspace_dir}/report_metadata.json
```

#### Step 2.5：报告校验与交付

1. 运行报告校验：
   - 通用行业：`python3 {skill_root}/scripts/validate_report.py {workspace_dir}/final_report.md --depth {depth_level}`
   - 餐饮/零售/供应链：`python3 {skill_root}/scripts/validate_report.py {workspace_dir}/final_report.md --depth {depth_level} --vertical restaurant-retail-supply-chain`
2. 如果校验失败，先根据错误信息修补报告并重新校验 1 次
3. 如果仍有 warning 但没有 error，可以交付，但要在摘要中提示“部分章节为降级生成/需复核”
4. 如果仍有 error，不要伪装完成；明确告诉用户哪一部分失败以及建议的下一步
5. 向用户交付最终报告，并保留 `execution_mode = normal|degraded` 与 `degradation_tags` 到元数据和摘要

**交付格式：**
```
✅ **行业调研报告生成完成**

📊 置信度评分：{overall_confidence}/10
📏 报告深度：{depth_level}
🛠️ 运行模式：{execution_mode}
🏷️ 降级标签：{degradation_tags，如无则为空数组}
📝 报告长度：{word_count} 字
🔢 引用数据点：{data_points_count} 个
⚠️ 悬置判断：{suspended_count} 项

报告已保存至：{workspace_dir}/final_report.md

{展示报告的执行摘要部分}

需要我展开某个章节的详细内容吗？或者要调整报告深度重新生成？
```

---

## 工作区管理

所有中间产物保存在 `{cwd}/industry-research-workspace/{topic_slug}/` 目录下。目录必须在阶段一结束前创建完成，且后续所有子 Agent 都收到同一个 `workspace_dir` 绝对路径。

```
industry-research-workspace/
└── {topic_slug}/
    ├── industry_anchor.json
    ├── context_dictionary.json
    ├── entity_evidence_plan.json   # 仅在适用餐饮/零售/供应链类研究时必须存在
    ├── expert_interview_guide.md   # 仅在用户提供访谈提纲或要求人机结合访谈准备时产出
    ├── ghost_deck.json
    ├── blue_r1.json
    ├── red_r1.json
    ├── blue_r2.json
    ├── red_r2.json
    ├── final_report.md
    └── report_metadata.json
```

---

## 错误处理

| 错误 | 处理 |
|------|------|
| Agent 返回 `clarification_needed` | 暂停流程，向用户提问 |
| 多个行动标题返回 `DATA_INSUFFICIENT` | 向用户汇报数据缺口，询问是否降低深度或缩小范围 |
| `web_fetch` 没拿到有效来源 | 先切换来源族（官方/公司/协会/主流媒体），再重试，最多 3 次 |
| Agent 输出不符合 JSON schema | 要求 Agent 重新生成，最多 2 次；仍失败则降级生成最小可用版本 |
| 红方或蓝方超时/失败 | 重试 1 次；仍失败则由主 Agent 补写缺失侧最小版本，并标记 `execution_mode = degraded`，同时追加 `degradation_tags += ["round1_missing_side"]` |
| `entity_evidence_plan.json` 仅达到降级标准 | 允许继续，但必须标记 `execution_mode = degraded`，同时追加 `degradation_tags += ["entity_mapping"]`，并在仲裁与交付中降低相关经营结论置信度 |
| `entity_evidence_plan.json` 连降级最低标准都达不到 | 暂停流程，向用户说明无法建立可信经营实体证据计划 |
| 第二轮未执行 | `brief` 视为正常；其他深度在元数据中说明原因并继续仲裁 |
| 置信度评分低于 4/10 | 向用户发出预警，建议缩小范围或补充数据源 |

---

## 参考文件索引

| 文件 | 用途 | 何时读取 |
|------|------|---------|
| `references/industry-taxonomy.md` | 行业分类代码映射 | Step 1.1 |
| `references/analytical-frameworks.md` | MECE/PESTLE/SCQ 等框架 | 传递给所有 Agent |
| `references/restaurant-retail-supply-chain-physics.md` | 餐饮/零售/供应链商业物理规则 | Step 2.0 |
| `references/evidence-triangulation-playbook.md` | 证据等级、三角验证、异常归因 | Step 2.0, Step 2.4 |
| `references/report-template.md` | 报告结构模板 | Step 2.4 传给仲裁 Agent |
| `agents/engagement-manager.md` | 编排 Agent 指令 | Step 2.1 |
| `agents/blue-team.md` | 蓝方 Agent 指令 | Step 2.2, 2.3 |
| `agents/red-team.md` | 红方 Agent 指令 | Step 2.2, 2.3 |
| `agents/chief-arbitrator.md` | 仲裁 Agent 指令 | Step 2.4 |
| `scripts/validate_report.py` | 报告结构校验 | Step 2.5 |
