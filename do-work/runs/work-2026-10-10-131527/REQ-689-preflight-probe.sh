#!/usr/bin/env bash
# Pre-flight baseline for REQ-689 (run before the build): the two shipped-prose contract tests that read the files this REQ edits pass today.
# shipped-package-reference-contract.sh resolves every backticked path and every "→ **Section**" citation in shipped Markdown;
# action-shell-blocks.sh checks every shell fence in shipped actions. Both finish in a few seconds.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
bash "$root/_dev/tests/shipped-package-reference-contract.sh"
bash "$root/_dev/tests/action-shell-blocks.sh"
