#!/bin/bash
# Focused checks for REQ-647: clarify.md carries the earmark condition and the shipped-package reference contract still holds.
set -euo pipefail
grep -q 'assigned_to' skills/do-work/actions/clarify.md
bash _dev/tests/shipped-package-reference-contract.sh
