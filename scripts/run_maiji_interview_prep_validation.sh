#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
WORKSPACE_DIR="$ROOT_DIR/industry-research-workspace/maiji-interview-prep-2026"

required_files=(
  "user_input.md"
  "industry_anchor.json"
  "context_dictionary.json"
  "entity_evidence_plan.json"
  "ghost_deck.json"
  "expert_interview_guide.md"
)

for file in "${required_files[@]}"; do
  test -f "$WORKSPACE_DIR/$file"
done

python3 -m json.tool "$WORKSPACE_DIR/industry_anchor.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/context_dictionary.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/entity_evidence_plan.json" >/dev/null
python3 -m json.tool "$WORKSPACE_DIR/ghost_deck.json" >/dev/null

rg -n "经营命题|必问口径|追问路径|交叉验证对象|危险回答信号" \
  "$WORKSPACE_DIR/expert_interview_guide.md" >/dev/null

rg -n "单店日销|同店同比|外卖|堂食|团购|单店UE|总部净利润|统一供货比例|加盟商画像|2026" \
  "$WORKSPACE_DIR/expert_interview_guide.md" >/dev/null

echo "Maiji interview prep validation passed for $WORKSPACE_DIR"
