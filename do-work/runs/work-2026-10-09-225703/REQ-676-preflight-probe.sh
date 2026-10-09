#!/bin/bash
# Pre-flight baseline for REQ-676 (run before the build): the fast citation contract that the new toolbox action must keep green passes today.
set -euo pipefail
bash _dev/tests/shipped-package-reference-contract.sh
