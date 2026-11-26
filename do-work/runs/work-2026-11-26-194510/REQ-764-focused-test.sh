#!/usr/bin/env bash
set -euo pipefail

repository_root=$(git rev-parse --show-toplevel)
module="$repository_root/skills/do-work/tools/do-work-cli"
tests=(
  TestAIReportGalleryListUsesProjectOverrideAndSkipsBrokenRecords
  TestAIReportGalleryPageLinksLatestRealReportsAndAssets
  TestAIReportGalleryShowsWorkedReportsAndKeepsProjectSamples
  TestReportGalleryStarterRoundTripAcrossSources
  TestReportGalleryRebuildsFromTheWrittenSuite
)
pattern=$(IFS='|'; printf '%s' "${tests[*]}")
if ! output=$(go -C "$module" test -count=1 -v -run "^($pattern)$" ./internal/toolboxcommands ./internal/suiteinstall 2>&1); then
  printf '%s\n' "$output" >&2
  exit 1
fi
for name in "${tests[@]}"; do
  if ! grep -qF -- "--- PASS: $name (" <<<"$output"; then
    printf 'Expected PASS for %s\n%s\n' "$name" "$output" >&2
    exit 1
  fi
  printf 'PASS %s\n' "$name"
done
