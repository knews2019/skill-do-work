#!/usr/bin/env bash
# Focused probe for REQ-714: the knowledge-base location tests of the do-work CLI.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
go test -count=1 -run 'BKB|KnowledgeBase' ./internal/knowledgecommands/
