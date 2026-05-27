# researcher

`researcher` 是行业研究引擎的命令行执行层。它负责检索、规划、证据整理、工作区生成和产物校验，为上层智能体技能提供可复核的研究基础。

它的定位不是直接替代研究员写结论，而是先把问题拆成可验证命题，再把搜索线索、证据台账、反证记录和置信度报告沉淀到一个结构化工作区里。

---

## 核心原则

搜索结果和模型回答只是线索，不是证据。任何高置信度结论都应该回到真实来源、经营痕迹、引用材料或多来源交叉验证。

对于连锁品牌、餐饮、零售和供应链问题，优先验证现实世界会留下的经营痕迹，例如门店、仓、招聘、许可、地图点位、SKU、履约范围和用户反馈。

---

## 构建

```bash
make build
```

构建完成后会生成：

```text
./researcher
```

检查版本和能力：

```bash
./researcher version
./researcher capabilities --json
```

---

## 命令总览

| 命令 | 作用 |
|:---|:---|
| `researcher version` | 输出当前版本 |
| `researcher help` | 查看命令帮助 |
| `researcher capabilities [provider] --json` | 查看检索服务能力 |
| `researcher retrieve "<query>" --provider bocha --count 10 --json` | 使用直接搜索获取网页线索 |
| `researcher answer volcengine "<query>" --model <model> --limit 10 --json` | 使用模型联网回答并返回来源标注 |
| `researcher plan "<question>" --domain chain-brand --json` | 生成命题和痕迹计划 |
| `researcher evidence "<question>" --json` | 生成空的证据台账骨架 |
| `researcher run "<question>" --domain chain-brand --depth standard --workspace-root <dir> --json` | 生成完整研究工作区 |
| `researcher validate <workspace_dir>` | 校验工作区产物是否齐全 |

---

## 快速示例

查看能力：

```bash
./researcher capabilities --json
```

生成连锁品牌痕迹计划：

```bash
./researcher plan "瑞幸咖啡 2026 年门店数目标是否可信？" \
  --domain chain-brand \
  --json
```

生成完整研究工作区：

```bash
./researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" \
  --domain chain-brand \
  --depth standard \
  --workspace-root ../industry-research-workspace \
  --json
```

校验工作区：

```bash
./researcher validate ../industry-research-workspace/<workspace>
```

---

## 检索服务

当前支持两类检索方式：

| Provider | 类型 | 适合用途 |
|:---|:---|:---|
| `bocha` | 直接网页搜索 | 获取网页、图片和搜索结果线索 |
| `volcengine` | 模型回答加联网搜索 | 获取带来源标注的综合回答 |

Bocha 示例：

```bash
BOCHA_API_KEY=... ./researcher retrieve "瑞幸咖啡 2026 门店数" \
  --provider bocha \
  --count 10 \
  --json
```

Volcengine Ark 示例：

```bash
ARK_API_KEY=... ./researcher answer volcengine \
  "瑞幸咖啡 2026 门店数目标是否可信？" \
  --json
```

---

## 环境变量

| 变量 | 作用 |
|:---|:---|
| `RESEARCHER_CONFIG` | 指向配置文件 |
| `BOCHA_API_KEY` | 覆盖 Bocha 密钥 |
| `ARK_API_KEY` | 覆盖 Volcengine Ark 密钥 |
| `XDG_CONFIG_HOME` | 指定默认配置目录 |

---

## 配置读取顺序

配置按以下顺序读取，越靠前优先级越高：

1. `--config <path>`
2. `RESEARCHER_CONFIG`
3. `$XDG_CONFIG_HOME/researcher/config.yaml`
4. `~/.config/researcher/config.yaml`

不要使用 `~/.researcher/config.yaml`。

---

## 值覆盖顺序

配置值按以下顺序覆盖，越靠后优先级越高：

1. 内置默认值
2. 配置文件
3. 环境变量中的密钥
4. 命令行参数

---

## 配置示例

真实密钥建议放在环境变量里，配置文件只保存非敏感默认值。

```yaml
providers:
  bocha:
    endpoint: "https://api.bocha.cn/v1/web-search"
  volcengine:
    endpoint: "https://ark.cn-beijing.volces.com/api/v3/responses"
    model: "doubao-seed-2-0-lite-260215"
defaults:
  providers:
    - bocha
    - volcengine
  depth: standard
  workspace_root: researcher-workspace
```

---

## 工作区产物

`researcher run` 会生成一个研究工作区，核心文件包括：

| 文件 | 作用 |
|:---|:---|
| `question.json` | 原始问题、领域、深度和创建时间 |
| `research_plan.json` | 研究计划 |
| `claim_graph.json` | 待验证命题 |
| `trace_plan.json` | 预期经营痕迹和反证痕迹 |
| `retrieval_log.json` | 检索调用、参数和来源使用记录 |
| `evidence_ledger.json` | 证据台账 |
| `disconfirmation_log.json` | 反证尝试 |
| `confidence_report.json` | 置信度评级 |
| `final_report.md` | 初版报告 |
| `report_metadata.json` | 元数据 |

---

## 本地检查

```bash
make fmt
make vet
make test
make build
```

如果测试环境禁止本地端口绑定，`make test` 里的 `httptest` 相关用例可能会失败。允许本地测试端口后重跑即可。
