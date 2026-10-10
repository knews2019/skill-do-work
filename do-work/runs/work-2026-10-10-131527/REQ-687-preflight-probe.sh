#!/usr/bin/env bash
# REQ-687 pre-build baseline: the two checks the REQ names (requirement 10), run against
# whichever tree this script is invoked from. Both are green at main HEAD bd56c4b0.
set -euo pipefail
repo_root="$(git rev-parse --show-toplevel)"
bash "$repo_root/_dev/tests/shipped-package-reference-contract.sh"
bash "$repo_root/_dev/tests/contract-regressions.sh"
