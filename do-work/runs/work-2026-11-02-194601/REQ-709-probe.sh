#!/usr/bin/env bash
# GREEN probe for REQ-709: recover leaves a note no request record can own alone,
# the commit inventory still lists it, and record-location dirt next to it still refuses.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
output="$(go test -count=1 -v -run '^(TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns|TestRecoverFinalizationStillRefusesRecordStateNextToANote|TestRecoverFinalizationIgnoresFinderMetadataUnderDoWork|TestRecoverFinalizationStopsWhenTheRefusalOwnsNoRequest|TestRecoverLeavesAnUntrackedInboxNoteForCommit)$' ./internal/finalization/ ./internal/lifecycleadvance/ 2>&1)" || { printf '%s\n' "$output"; exit 1; }
grep -q -- '--- PASS: TestRecoverFinalizationIgnoresNotesNoRequestRecordOwns' <<<"$output"
grep -q -- '--- PASS: TestRecoverFinalizationStillRefusesRecordStateNextToANote' <<<"$output"
grep -q -- '--- PASS: TestRecoverLeavesAnUntrackedInboxNoteForCommit' <<<"$output"
if grep -q -- '--- FAIL' <<<"$output"; then printf '%s\n' "$output"; exit 1; fi
printf 'REQ-709 GREEN: recover ignores unowned do-work notes and still refuses record-location dirt\n'
