#!/usr/bin/env bash
set -euo pipefail
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli
go test -count=1 -run 'TestRecover' ./internal/lifecycleadvance/
