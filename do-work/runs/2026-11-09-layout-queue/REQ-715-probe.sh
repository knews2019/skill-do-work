#!/usr/bin/env bash
# Focused probe for REQ-715: the shipped reference contract (changelog mirror and canonical links)
# and the layout check's view of the root.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root"
bash _dev/tests/shipped-package-reference-contract.sh
