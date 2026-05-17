# 调研报告结构模板

> 本模板定义最终输出报告的标准结构。
> Chief Arbitrator 必须按此模板生成报告。
> 支持三个深度级别：精简版 / 标准版 / 深度版。

---

## 深度级别定义

| 级别 | 标识符 | 目标字数 | 适用场景 |
|------|--------|---------|---------|
| 精简版 | `brief` | ~3,000字 | 快速决策、内部沟通 |
| 标准版 | `standard` | ~8,000字 | 投资决策、战略规划 |
| 深度版 | `comprehensive` | ~15,000字 | 尽调报告、董事会材料 |

---

## 报告结构

### 封面信息
```
# {行业/主题} 战略调研报告
生成时间：{timestamp}
调研深度：{brief|standard|comprehensive}
置信度评分：{overall_confidence}/10
行业代码锚定：{industry_codes}
```

### 一、执行摘要 (Executive Summary)
**所有深度级别必须包含**

- 核心发现（3-5 条，每条一句话）
- 战略建议（2-3 条，可操作的行动项）
- 关键风险提示（2-3 条）
- 整体置信度评分及说明

### 二、SCQ 开篇
**所有深度级别必须包含**

按 SCQ 框架撰写：
- **情境 (S)**：行业的公认背景事实
- **冲突 (C)**：打破情境的核心矛盾
- **疑问 (Q)**：本报告回答的核心决策问题

### 三、行业全景 (Industry Landscape)
**标准版和深度版包含**

- 行业定义与边界（引用国民经济行业分类代码）
- 市场规模与增长趋势
- 产业链全景图（上游→中游→下游）
- 主要玩家竞争格局

### 四、实体经营版图与证据链 (Operating Footprint & Evidence Chain)
**餐饮、零售、供应链相关报告必须包含；标准版和深度版完整展开，精简版可压缩进执行摘要或核心论证**

- 本节表格为主要结构化表示，生成器应优先以矩阵行作为输出、溯源和仲裁的主键。
- 最小经营单元：门店、仓、中央厨房、加盟商、区域代理、车队、供应商、数字节点等。
- 行级字段要求：每行必须显式覆盖 operating unit、claim、scope、evidence-family coverage、triangle coverage、ref_id、time scope、异常/冲突说明、verdict。
- 证据三角验证矩阵：每个关键经营命题对应至少三类 evidence families，并映射到三角验证矩阵字段。
- evidence families 为受限字段：仅使用 `capital/legal`、`people/org`、`physical fulfillment`、`digital frontend`、`terminal feedback` 五个固定家族名称，对应三角验证剧本中的资本/法务、人力/组织、物理履约、数字前端、终端反馈。
- triangle coverage：记录本行实际使用的证据家族数量与集合，便于校验是否满足三角覆盖。
- ref_id 规则：每行必须记录 source reference identifiers，统一使用可映射到数据溯源表的 `ref_id` 集合，不使用松散 source refs 文本。
- time scope：每行标注时间范围或 recency，避免把不同时间截面的证据混写。
- 异常/冲突说明：明确写出数据互相打架的地方，以及最可能商业解释。
- 裁决标签：`VERIFIED_OPERATING_FACT` / `HIGH_CONFIDENCE_INFERENCE` / `EXPLAINABLE_ANOMALY` / `SUSPENDED_JUDGMENT` / `UNVERIFIED_NARRATIVE`，并能追溯到对应 evidence families、triangle coverage 与 ref_id。
- 若该条件章节在最终报告中被省略，后续章节标题应在最终输出中重新编号并顺延，保持连续编号。

示例表格：

| operating unit | claim | scope | evidence families | triangle coverage | ref_id | time scope | 异常/冲突说明 | verdict |
|----------------|-------|-------|-------------------|-------------------|--------|------------|-----------------|---------|
| 华东区域中心仓 | 支持 300km 日配 | 华东核心城群 | `capital/legal`,`people/org`,`physical fulfillment`,`digital frontend`,`terminal feedback` | 5/5: capital/legal+people/org+physical fulfillment+digital frontend+terminal feedback | `warehouse_license_huadong`,`driver_recruiting_eastchina`,`warehouse_poi_route_scan`,`miniapp_delivery_range`,`user_delivery_feedback_q1` | 2025Q4-2026Q1 | 配送半径宣传口径大于 POI 实测覆盖，可能因前置仓与限时达范围混用 | HIGH_CONFIDENCE_INFERENCE |

### 五、核心论证 (Core Arguments)

对于每个核心议题（来自幽灵卡片的行动标题）：

#### 5.x {议题名称}

##### 正题 (Bull Case)
- 核心论点
- 支撑数据（附来源）
- 逻辑链条

##### 反题 (Bear Case)
- 核心反驳
- 反向数据（附来源）
- 风险分析

##### 合题 (Synthesis)
- 正交维度识别
- 仲裁结论
- 置信度评分：{0-10}/10
- 悬置判断（如有）

> **精简版**：每个议题 200-300 字，合并正反为一段
> **标准版**：每个议题 800-1200 字，完整三段论
> **深度版**：每个议题 1500-2500 字，含详细数据表格

### 六、战略路线图 (Strategic Roadmap)
**标准版和深度版包含**

- 短期行动（0-6个月）
- 中期布局（6-18个月）
- 长期愿景（18-36个月）
- 关键里程碑与决策节点

### 七、风险矩阵 (Risk Matrix)
**所有深度级别必须包含**

| 风险项 | 概率 | 影响 | 缓解策略 | 数据来源 |
|--------|------|------|---------|---------|
| ... | HIGH/MED/LOW | 严重/中等/轻微 | ... | ... |

### 八、数据溯源表 (Data Provenance)
**所有深度级别必须包含**

| 数据点 | 引用值 | 来源 | 来源类型 | 时效性 | 置信度 |
|--------|--------|------|---------|--------|--------|
| 市场规模 | XX亿元 | {URL/报告名} | 官方统计/行业报告/新闻 | 2024Q3 | 高/中/低 |

### 九、附录 (Appendix)
**仅深度版包含**

- 完整的 PESTLE 分析表
- 波特五力详细评估
- 被排除的低优先级议题及排除理由
- 方法论说明

---

## 置信度评分标准

| 分数 | 含义 | 数据条件 |
|------|------|---------|
| 9-10 | 高度确信 | 多源交叉验证，官方数据为主 |
| 7-8 | 较为确信 | 有可靠数据支撑，少量推理补充 |
| 5-6 | 中等确信 | 数据有限但逻辑自洽，需后续验证 |
| 3-4 | 较低确信 | 数据不足，主要依赖推理和类比 |
| 1-2 | 高度不确定 | 几乎无数据支撑，存在悬置判断 |

---

## 格式规范

- 所有数字附单位和时间基准（如"2024年市场规模约500亿元"，不是"市场规模约500亿"）
- 百分比保留一位小数
- 货币单位统一（人民币用"元"或"亿元"，美元用"USD"）
- 引用来源用方括号标注序号（如 [1]），在数据溯源表中对应
- 每个核心论点必须至少有一个定量数据支撑
