#!/usr/bin/env bash
# Focused probe for REQ-721 (ai-report closing band suggests a Tinka canvas): the checks that read shipped Markdown.
set -uo pipefail
cd "$(git rev-parse --show-toplevel)" || exit 98
bash _dev/tests/action-shell-blocks.sh || exit 1
bash _dev/tests/contract-regressions.sh || exit 2
bash _dev/tests/shipped-package-reference-contract.sh || exit 3
