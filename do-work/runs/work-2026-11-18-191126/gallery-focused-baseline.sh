#!/usr/bin/env bash
set -euo pipefail

repository_root="$(git rev-parse --show-toplevel)"
cd "$repository_root/skills/do-work/tools/do-work-cli"
tests=(
  TestAIReportGalleryListUsesProjectOverrideAndSkipsBrokenRecords
  TestAIReportGalleryPageLinksLatestRealReportsAndAssets
  TestAIReportGalleryRefusesHandMadeOrLinkedOutput
  TestAIReportGalleryShowsWorkedReportsAndKeepsProjectSamples
  TestReportGalleryResolveUsesDiscoveredNames
  TestReportGalleryCopyNamesRoundTrip
  TestAIReportIndexCatalogsEveryBundleNamingStyleOnce
  TestAIReportIndexReadsBundlesInsideMonthFolders
)
pattern="$(IFS='|'; printf '%s' "${tests[*]}")"
output="$(go test -count=1 -v -run "^($pattern)$" ./internal/toolboxcommands 2>&1)" || {
  printf '%s\n' "$output" >&2
  exit 1
}
for name in "${tests[@]}"; do
  if ! grep -qF -- "--- PASS: $name (" <<<"$output"; then
    printf 'Expected PASS for %s\n%s\n' "$name" "$output" >&2
    exit 1
  fi
  printf 'PASS %s\n' "$name"
done
