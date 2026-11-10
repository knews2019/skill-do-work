#!/usr/bin/env bash
# Pre-build baseline for REQ-716 (move the root installer scripts from tools/ to cmd/):
# the fast focused suites this REQ edits, run at the base revision. The heavy suites it also
# edits (staged-skills, install-suite, update-script) run in the builder's verification and the
# integrator's heavy lane; they passed at base 792f981c (34 s, 30 s, 67 s).
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"
bash _dev/tests/suite-manifest-contract.sh
bash _dev/tests/contracts/replace-text-section.sh
bash _dev/tests/audit-lockins.sh
cd skills/do-work/tools/do-work-cli
go test -count=1 -run '^TestReleaseRefusedWhenTheImplementationShipsNothing$' ./internal/finalization/
go test -count=1 -run '^TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint$' ./internal/heavyverification/
