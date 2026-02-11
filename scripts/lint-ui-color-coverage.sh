#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "${ROOT_DIR}"

fail=0

report() {
  local title="$1"
  local output="$2"
  if [[ -n "${output}" ]]; then
    echo "[lint-ui-color-coverage] ${title}"
    echo "${output}"
    echo
    fail=1
  fi
}

# Detect plain heading literals printed directly to stdout/stderr.
# These headings should be routed through shared styled helpers:
# - renderWorkspacesTitle / renderResultTitle / renderRiskTitle / renderProgressTitle
# - printSection(styleBold(...), ...)
plain_heading_hits="$(rg -n --color never \
  --glob 'internal/cli/*.go' \
  --glob '!internal/cli/*_test.go' \
  'fmt\.Fprint(f|ln)?\((c\.(Out|Err)|out|w),\s*"(Repo pool:|Repos\(pool\):|Plan:|Result:|Risk:|Progress:|Workspaces\(|Contexts:|Action:)"' || true)"
report "direct plain heading literal detected (potential missing color styling)" "${plain_heading_hits}"

# Purge risk section currently prints some structural labels directly.
# Detect these labels so they are visible in audit results.
purge_plain_labels="$(rg -n --color never \
  --glob 'internal/cli/ws_purge.go' \
  'fmt\.Fprintf\(out,\s*"%s(selected: %d\\n|active workspace risk detected:\\n)"' || true)"
report "purge risk section has plain structural labels (potential missing semantic styling)" "${purge_plain_labels}"

if [[ "${fail}" -ne 0 ]]; then
  echo "[lint-ui-color-coverage] FOUND"
  exit 1
fi

echo "[lint-ui-color-coverage] OK"
