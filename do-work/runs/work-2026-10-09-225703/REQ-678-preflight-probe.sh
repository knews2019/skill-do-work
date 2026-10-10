#!/bin/bash
# Pre-flight baseline for REQ-678 (run before the build): the four delivery-stage names REQ-678 cites exist in core review Step 7, and the fast citation contract the new toolbox action must keep green passes today.
set -euo pipefail
for stage_name in 'Implementation' 'Integration' 'Deployment' 'Live acceptance'; do grep -q "^- \*\*$stage_name\*\*: " skills/do-work/actions/review-work.md; done
grep -q '^### Step 7: Acceptance Testing$' skills/do-work/actions/review-work.md
bash _dev/tests/shipped-package-reference-contract.sh
