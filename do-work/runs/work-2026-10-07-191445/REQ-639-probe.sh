#!/usr/bin/env bash
# Focused tests for REQ-639 (prose-only change to work.md, work-reference.md, docs/work-guide.md):
# the citation contract that resolves action-file headings, and the shipped prose contract regressions.
set -euo pipefail
bash _dev/tests/shipped-package-reference-contract.sh
bash _dev/tests/contract-regressions.sh
