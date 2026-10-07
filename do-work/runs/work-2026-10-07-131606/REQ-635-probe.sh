#!/usr/bin/env bash
# Focused tests for REQ-635: VisibleSections indent, fence-before-comment and code-span rules, plus their recovery and Timing-writer consumers.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
go test -C "$root/skills/do-work/tools/do-work-cli" -count=1 -run 'VisibleSections|Recovery|TimingReplacement' ./internal/requestmodel/ ./internal/requeststate/ ./internal/lifecycletiming/
