#!/bin/bash
# Focused tests for REQ-645: recovery section removal and VisibleSections indent handling in the do-work CLI.
set -euo pipefail
cd skills/do-work/tools/do-work-cli
go test ./internal/requeststate/ ./internal/requestmodel/ -run 'Recovery|VisibleSections|StripGenerated' -count=1
