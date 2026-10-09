#!/bin/bash
# Focused checks for REQ-652: the operator sentences sit where the REQ puts them and the reference contract holds.
set -euo pipefail
grep -q 'own time' skills/do-work/actions/capture.md
grep -q 'Done — I did it myself' skills/do-work/actions/clarify.md
grep -q 'operator' skills/do-work/actions/work.md
grep -q 'operator' skills/do-work/actions/work-reference.md
bash _dev/tests/shipped-package-reference-contract.sh
