#!/usr/bin/env bash
# Pre-flight baseline for REQ-692 (run before the build): the invariants the build must keep hold
# today. validate-feedback Steps 1 to 5 and its Output Format fence match their base hashes, no
# live shipped file spells a retired core validate-feedback trigger, and the shipped-package
# citation contract and the shipped shell-fence lint pass.
set -euo pipefail
exec bash "$(dirname "${BASH_SOURCE[0]}")/REQ-692-probe.sh" invariants
