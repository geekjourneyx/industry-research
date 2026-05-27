<div align="center">

# 行业研究引擎

**面向智能体框架的证据优先研究工作区**

<img src="assets/agent-skill-banner.png" alt="行业研究引擎：面向智能体框架的证据优先研究工作区" width="100%">

[![许可证：MIT](https://img.shields.io/badge/%E8%AE%B8%E5%8F%AF%E8%AF%81-MIT-blue.svg)](./LICENSE)
[![Go 1.22](https://img.shields.io/badge/Go-1.22-00ADD8.svg)](./researcher/go.mod)
[![智能体技能](https://img.shields.io/badge/%E6%99%BA%E8%83%BD%E4%BD%93-%E6%8A%80%E8%83%BD-c96442.svg)](./SKILL.md)
[![证据优先](https://img.shields.io/badge/%E8%AF%81%E6%8D%AE-%E4%BC%98%E5%85%88-2ea043.svg)](./references/evidence-ledger-schema.md)

</div>

---

## 这是什么

行业研究引擎是一套智能体技能和研究引擎，用来把模糊的行业、赛道、公司和商业机会问题，转成可复核的研究工作区。

它不绑定单一智能体框架。只要运行环境支持技能式指令、工具编排和文件工作流，就可以适配到 Claude Code、OpenClaw、hermes-agent 以及类似系统。

它的核心原则很简单：搜索结果只是线索，不是证据。高置信度结论必须经过来源核验、经营痕迹验证、交叉比对和反证尝试。

```text
输入：瑞幸咖啡 2026 年门店数目标是否可信？
输出：可复核工作区、证据台账、反证记录、置信度报告和最终调研报告
```

---

## 核心特性

<img src="assets/agent-skill-features.png" alt="行业研究引擎核心特性：框架无关、研究工作区、证据纪律、对抗式评审、经营痕迹、报告校验" width="100%">

---

## 工作流程

<img src="assets/agent-skill-workflow.png" alt="行业研究引擎工作流程：提问、规划、检索、审查、报告" width="100%">

---

## 安装

克隆仓库并构建研究引擎：

```bash
git clone git@github.com:geekjourneyx/industry-research.git
cd industry-research/researcher
make build
```

检查可执行文件：

```bash
./researcher version
./researcher capabilities --json
```

更多命令、配置和检索服务说明见 [researcher 使用说明](./researcher/README.md)。

可选的检索服务密钥：

```bash
export BOCHA_API_KEY="..."
export ARK_API_KEY="..."
```

配置读取顺序如下：

```text
--config <path>
RESEARCHER_CONFIG
$XDG_CONFIG_HOME/researcher/config.yaml
~/.config/researcher/config.yaml
```

---

## 快速上手

如果只想使用底层研究引擎，可以直接阅读 [researcher 使用说明](./researcher/README.md)。

生成一份连锁品牌研究工作区：

```bash
cd researcher

./researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" \
  --domain chain-brand \
  --depth standard \
  --workspace-root ../industry-research-workspace \
  --json
```

校验生成的工作区：

```bash
./researcher validate ../industry-research-workspace/<workspace>

python3 ../scripts/validate_report.py \
  --researcher-workspace ../industry-research-workspace/<workspace>
```

运行本地检查：

```bash
make fmt
make vet
make test
make build
```

---

## 产物结构

每次执行 `researcher run` 都会生成以下工作区文件：

| 文件 | 作用 |
|:---|:---|
| `question.json` | 原始问题、研究领域、报告深度和创建时间 |
| `research_plan.json` | 研究计划和执行范围 |
| `claim_graph.json` | 需要验证或反证的核心命题 |
| `trace_plan.json` | 每个命题应当留下的现实世界痕迹 |
| `retrieval_log.json` | 检索调用、参数、重试和来源使用情况 |
| `evidence_ledger.json` | 证据记录和支撑状态 |
| `disconfirmation_log.json` | 反证尝试和削弱结论的记录 |
| `confidence_report.json` | 置信度评级和理由 |
| `final_report.md` | 面向用户的最终研究报告 |
| `report_metadata.json` | 执行模式、降级标签和报告元数据 |

餐饮、零售、供应链和连锁品牌研究中，最关键的质量文件是 `trace_plan.json`、`evidence_ledger.json`、`disconfirmation_log.json` 和 `confidence_report.json`。

---

## 证据原则

系统会把弱研究报告里经常混在一起的四层内容拆开：

| 层级 | 使用方式 |
|:---|:---|
| 搜索结果 | 只能作为后续核验线索，不能直接当成证据 |
| 来源材料 | 可访问的文件、网页、公告、帖子或经营痕迹 |
| 推理结果 | 基于事实和约束得出的估算，必须明确标注 |
| 最终结论 | 结论置信度必须和证据质量匹配 |

对于毛利、单店模型、产能利用率、SKU 动销、履约成本、留存率等隐藏经营变量，报告必须区分公开披露事实和推理区间，并写明还缺哪些一手数据。

---

## 项目地图

| 路径 | 作用 |
|:---|:---|
| `SKILL.md` | 智能体技能编排契约 |
| [`researcher/`](./researcher/README.md) | 用于生成工作区、检索、规划和校验的研究引擎 |
| `agents/` | 经理、蓝队、红队和仲裁者等评审角色 |
| `references/` | 证据规则、报告模板、痕迹推理和分析框架 |
| `scripts/validate_report.py` | 报告和工作区校验脚本 |
| `evals/` | 评测问题和预期产物 |

---

## 许可证

[MIT](./LICENSE) — 自由使用、修改、分发。

---

## 作者

| | |
|:---|:---|
| 个人主页 | [jieni.ai](https://jieni.ai) |
| GitHub | [geekjourneyx](https://github.com/geekjourneyx) |
| Twitter | [@seekjourney](https://x.com/seekjourney) |
| 公众号 | 微信搜「极客杰尼」 |
