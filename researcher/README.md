# researcher

`researcher` is a command-line skeleton for the industry research engine. It will coordinate retrieval, planning, evidence review, execution, and validation steps for research workflows.

## Commands

- `researcher version` prints the current binary version.
- `researcher help`, `researcher --help`, and `researcher -h` print command help.
- `researcher capabilities --json` will describe supported capabilities in JSON.
- `researcher retrieve` will collect source material.
- `researcher plan` will prepare a research plan.
- `researcher evidence` will organize and evaluate source evidence.
- `researcher run` will run a full research workflow.
- `researcher validate` will validate generated outputs.

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

```yaml
providers:
  bocha:
    api_key: "bocha-api-key"
    endpoint: "https://api.bochaai.com/v1/web-search"
  volcengine:
    api_key: "ark-api-key"
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
