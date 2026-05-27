#!/usr/bin/env bash
set -euo pipefail

export LC_ALL=C
export LANG=C

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
cd "$repo_root/researcher"

version="$(tr -d '[:space:]' < VERSION)"
dist_dir="${DIST_DIR:-$repo_root/dist}"

rm -rf "$dist_dir"
mkdir -p "$dist_dir"

platforms=(
  "darwin amd64"
  "darwin arm64"
  "linux amd64"
  "linux arm64"
  "windows amd64"
  "windows arm64"
)

for platform in "${platforms[@]}"; do
  read -r goos goarch <<<"$platform"
  name="researcher_${version}_${goos}_${goarch}"
  work_dir="$dist_dir/$name"
  mkdir -p "$work_dir"

  binary="researcher"
  if [[ "$goos" == "windows" ]]; then
    binary="researcher.exe"
  fi

  echo "building $name"
  GOOS="$goos" GOARCH="$goarch" CGO_ENABLED=0 \
    go build -trimpath -ldflags="-s -w -X main.Version=${version}" \
    -o "$work_dir/$binary" ./cmd/researcher

  cp README.md VERSION "$work_dir/"
  if [[ -f "$repo_root/LICENSE" ]]; then
    cp "$repo_root/LICENSE" "$work_dir/"
  fi

  if [[ "$goos" == "windows" ]]; then
    (cd "$dist_dir" && zip -qr "${name}.zip" "$name")
  else
    (cd "$dist_dir" && tar -czf "${name}.tar.gz" "$name")
  fi

  rm -rf "$work_dir"
done

(cd "$dist_dir" && shasum -a 256 researcher_${version}_* > "researcher_${version}_checksums.txt")

echo "release artifacts:"
ls -lh "$dist_dir"
