# researcher

`researcher` is the command-line execution layer for the industry research engine. It coordinates retrieval, planning, evidence review, workspace generation, and validation for research workflows.

## Commands

- `researcher version` prints the current binary version.
- `researcher help`, `researcher --help`, and `researcher -h` print command help.
- `researcher capabilities [provider] --json` describes supported provider capabilities.
- `researcher retrieve "<query>" --provider bocha --count 10 --json` collects direct web-search leads.
- `researcher answer volcengine "<query>" --model <model> --limit 10 --json` asks a model to answer with web-search annotations.
- `researcher plan "<question>" --domain chain-brand --json` builds a trace plan.
- `researcher evidence "<question>" --json` creates an empty evidence ledger shell.
- `researcher run "<question>" --domain chain-brand --depth standard --workspace-root <dir> --json` creates a complete workspace.
- `researcher validate <workspace_dir>` validates generated workspace artifacts.

## Build

```bash
make build
```

## Examples

```bash
./researcher capabilities --json
./researcher plan "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --json
./researcher run "瑞幸咖啡 2026 年门店数目标是否可信？" --domain chain-brand --depth standard --json
./researcher validate researcher-workspace/<workspace>
```

For Bocha direct search:

```bash
BOCHA_API_KEY=... ./researcher retrieve "瑞幸咖啡 2026 门店数" --provider bocha --count 10 --json
```

For Volcengine Ark model search:

```bash
ARK_API_KEY=... ./researcher answer volcengine "瑞幸咖啡 2026 门店数目标是否可信？" --json
```

## Environment Variables

- `RESEARCHER_CONFIG` points to a config YAML file.
- `BOCHA_API_KEY` overrides the Bocha API key.
- `ARK_API_KEY` overrides the Volcengine Ark API key.
- `XDG_CONFIG_HOME` sets the XDG config directory used for default config discovery.

## Config Lookup Order

1. `--config <path>`
2. `RESEARCHER_CONFIG`
3. `$XDG_CONFIG_HOME/researcher/config.yaml`
4. `~/.config/researcher/config.yaml`

## Value Precedence

Values are resolved from lowest to highest precedence:

1. Built-in defaults
2. YAML config file values
3. Environment variable overrides for API keys
4. Explicit command flags

## Example Config

Use environment variables for real API secrets. Keep config files for non-secret defaults whenever possible.

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

## Evidence Policy

Search results and model answers are leads, not evidence. Treat them as starting points for verification against primary sources, operating traces, and cited source material.
