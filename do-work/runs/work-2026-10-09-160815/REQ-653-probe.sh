#!/bin/bash
# Focused checks for REQ-653: the companion file exists, both source files shrank below their 0.305.73 word counts, and the two fast citation contracts hold.
set -euo pipefail
test -f skills/do-work/actions/fan-out-reference.md
test "$(wc -w < skills/do-work/actions/work.md)" -lt 12711
test "$(wc -w < skills/do-work/actions/work-reference.md)" -lt 22162
bash _dev/tests/shipped-package-reference-contract.sh
bash _dev/tests/prescribed-shell-canonicalization.sh
