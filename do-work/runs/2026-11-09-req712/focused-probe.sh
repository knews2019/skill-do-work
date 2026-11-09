#!/usr/bin/env bash
# Focused tests for REQ-712: the ai-report second-opinion command tests, including the shipped-preset test.
set -euo pipefail
cd "$(git rev-parse --show-toplevel)/skills/do-work/tools/do-work-cli"
go test ./internal/toolboxcommands -run 'SecondOpinion' -count=1
