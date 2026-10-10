#!/usr/bin/env bash
# REQ-658 pre-build baseline: the finalize-path Go tests the change reshapes
# (prepareBoundJournal, validateManifest, the finalize handler), the requeststate
# writer-label tests, and the contract check that guards the Step 8/9 prose.
# All must be green at the base commit before the builder starts.
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cli_root="$repo_root/skills/do-work/tools/do-work-cli"
go test -C "$cli_root" -count=1 \
  -run '^(TestFinalize|TestPrepareBoundJournal|TestValidateManifest|TestRecoverFinalizationResumesJournalAfterLifecycleInterruption$|TestRecoverFinalizationAlreadyGreenNoReleaseManifest$)' \
  ./internal/finalization/
go test -C "$cli_root" -count=1 -run 'Writer' ./internal/requeststate/
bash "$repo_root/_dev/tests/contracts/core-checks.sh"
