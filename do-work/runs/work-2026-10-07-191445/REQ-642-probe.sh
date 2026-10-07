#!/bin/bash
# Focused tests for REQ-642: the shipped-package citation contract and the contract regressions.
set -euo pipefail
bash _dev/tests/shipped-package-reference-contract.sh
bash _dev/tests/contract-regressions.sh
