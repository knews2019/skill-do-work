#!/usr/bin/env bash
# REQ-735 pre-build baseline: the two contract scripts the builder relies on for crew-member and
# commit-action prose. Both must pass at the base commit; the builder reruns them after its edits.
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"
bash _dev/tests/shipped-package-reference-contract.sh
bash _dev/tests/contract-regressions.sh
printf 'REQ-735 preflight PASS\n'
