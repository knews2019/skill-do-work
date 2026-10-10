#!/usr/bin/env bash
# REQ-656 pre-build baseline: the two checks the REQ names (requirement 10), run against
# whichever tree this script is invoked from. Both are green at main HEAD 7fbeff30
# (contract-regressions.sh takes about 25 s on a quiet machine).
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
bash "$repo_root/_dev/tests/shipped-package-reference-contract.sh"
bash "$repo_root/_dev/tests/contract-regressions.sh"
