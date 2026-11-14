#!/usr/bin/env bash
# Pre-build baseline for REQ-730 (the report catalog lists each bundle's second-opinion pages):
# the existing catalog, organize and second-opinion banner tests this REQ must keep green, run
# at the base revision. All pass at base; the new REQ-730 tests do not exist yet.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
go test -count=1 -run '^(TestAIReportIndex|TestAIReportOrganize|TestSecondOpinionBanner)' ./internal/toolboxcommands/
