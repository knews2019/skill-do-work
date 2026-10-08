#!/bin/bash
# Focused checks for REQ-648: the operator rule sits in both run-action files and the shipped-package reference contract still holds.
set -euo pipefail
grep -q 'operator' skills/do-work/actions/work-reference.md
grep -q 'operator' skills/do-work/actions/work.md
bash _dev/tests/shipped-package-reference-contract.sh
