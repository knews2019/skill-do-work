#!/bin/bash
# GREEN checks for REQ-675 (integrator test gate): review-work.md Steps 7-8 name the deployment and live acceptance stages and the "unassessed" rule, the guide's Phase 3 says it too, and the fast citation contract holds. Fails before the build.
set -euo pipefail
acceptance_steps="$(sed -n '/^### Step 7: Acceptance Testing/,/^### Step 9/p' skills/do-work/actions/review-work.md)"
for stage_word in implementation integration deployment 'live acceptance' unassessed; do
  grep -qi "$stage_word" <<<"$acceptance_steps" || { echo "review-work.md Steps 7-8 lack: $stage_word"; exit 1; }
done
phase_three="$(sed -n '/^### Phase 3/,/^## /p' skills/do-work/docs/review-work-guide.md)"
grep -qi 'unassessed' <<<"$phase_three" || { echo "review-work-guide.md Phase 3 lacks the unassessed line"; exit 1; }
bash _dev/tests/shipped-package-reference-contract.sh
