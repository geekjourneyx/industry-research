#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
WORKSPACE_DIR="$ROOT_DIR/industry-research-workspace/rrsc-smoke-lower-tier-expansion-2026"

required_files=(
  "industry_anchor.json"
  "context_dictionary.json"
  "entity_evidence_plan.json"
  "ghost_deck.json"
  "blue_r1.json"
  "red_r1.json"
  "blue_r2.json"
  "red_r2.json"
  "final_report.md"
  "report_metadata.json"
)

for file in "${required_files[@]}"; do
  test -f "$WORKSPACE_DIR/$file"
done

python3 -m json.tool "$WORKSPACE_DIR/industry_anchor.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/context_dictionary.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/entity_evidence_plan.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/ghost_deck.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/blue_r1.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/red_r1.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/blue_r2.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/red_r2.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/report_metadata.json" >/dev/null

python3 "$ROOT_DIR/scripts/validate_report.py" \
  "$WORKSPACE_DIR/final_report.md" \
  --depth standard \
  --vertical restaurant-retail-supply-chain

echo "RRSC smoke validation passed for $WORKSPACE_DIR"
