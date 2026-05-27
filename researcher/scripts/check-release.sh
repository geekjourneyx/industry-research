#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root"

version="$(tr -d '[:space:]' < researcher/VERSION)"
tag="${RELEASE_TAG:-${GITHUB_REF_NAME:-}}"

if [[ -z "$tag" ]]; then
  tag="$(git describe --tags --exact-match 2>/dev/null || true)"
fi

if [[ -z "$tag" ]]; then
  echo "release check failed: RELEASE_TAG or exact Git tag is required" >&2
  exit 1
fi

expected_tag="v${version}"
if [[ "$tag" != "$expected_tag" ]]; then
  echo "release check failed: tag '$tag' must match researcher/VERSION '$version' as '$expected_tag'" >&2
  exit 1
fi

if [[ ! -f CHANGELOG.md ]]; then
  echo "release check failed: CHANGELOG.md is required" >&2
  exit 1
fi

if ! grep -Eq "^## \[?v?${version}\]?([[:space:]-]|$)" CHANGELOG.md; then
  echo "release check failed: CHANGELOG.md must contain an entry for $version" >&2
  exit 1
fi

if [[ ! -f README.md || ! -f researcher/README.md ]]; then
  echo "release check failed: root README.md and researcher/README.md are required" >&2
  exit 1
fi

required_readme_patterns=(
  "npx skills add geekjourneyx/industry-research"
  "Go CLI + 智能体技能"
  "researcher 使用说明"
  "CHANGELOG.md"
)

for pattern in "${required_readme_patterns[@]}"; do
  if ! grep -Fq "$pattern" README.md; then
    echo "release check failed: README.md is missing '$pattern'" >&2
    exit 1
  fi
done

required_researcher_patterns=(
  "make build"
  "make test"
  "researcher run"
  "GitHub Release"
)

for pattern in "${required_researcher_patterns[@]}"; do
  if ! grep -Fq "$pattern" researcher/README.md; then
    echo "release check failed: researcher/README.md is missing '$pattern'" >&2
    exit 1
  fi
done

echo "release check passed for $tag"
