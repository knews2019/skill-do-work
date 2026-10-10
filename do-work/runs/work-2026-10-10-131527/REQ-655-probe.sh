#!/usr/bin/env bash
# GREEN checks for REQ-655 (integrator test gate): the ai-report-index verb is registered, catalogs every bundle naming style once, links a superseded proposal to its successor, finds matches newest first, refuses a hand-made catalog, and the toolboxcommands package is gofmt-clean and vets. A skip is not a pass. Fails before the build (the tests do not exist yet). Under 30 s on a quiet machine.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
module_dir="$root/skills/do-work/tools/do-work-cli"
package_dir="$module_dir/internal/toolboxcommands"
unformatted="$(gofmt -l "$package_dir")"
[ -z "$unformatted" ] || { printf 'gofmt -l lists: %s\n' "$unformatted"; exit 1; }
go vet -C "$module_dir" ./internal/toolboxcommands/
probe_tests=(
  TestAIReportIndexCatalogsEveryBundleNamingStyleOnce
  TestAIReportIndexLinksSupersededProposalToSuccessor
  TestAIReportIndexFindListsMatchesNewestFirstIncludingSuperseded
  TestAIReportIndexRefusesHandMadeCatalog
  TestHandlersRegisterCanonicalToolboxCommands
)
output="$(go test -C "$module_dir" -count=1 -v -run "^($(IFS='|'; echo "${probe_tests[*]}"))\$" ./internal/toolboxcommands/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
for probe_test in "${probe_tests[@]}"; do
  grep -q -- "^--- PASS: $probe_test " <<<"$output" || { printf '%s\n' "$output"; echo "$probe_test did not pass (skipped or missing)"; exit 1; }
done
