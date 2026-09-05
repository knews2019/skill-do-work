#!/usr/bin/env bash
# Focused probe recorded as REQ-583 run evidence: the exact argv this run executed.
# The package under test is the one the REQ's three new pins live in; the whole package
# is run rather than a -run selection because the pins have to hold beside the tests
# REQ-506 left behind, and the file's per-test budget is the constraint, not the count.
set -euo pipefail
go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/lifecycleadvance
