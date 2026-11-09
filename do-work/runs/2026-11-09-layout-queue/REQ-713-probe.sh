#!/usr/bin/env bash
# Focused probe for REQ-713: the folder-structure tests, which hold the default-parameter recipe case.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
cd "$root/skills/do-work/tools/do-work-cli"
go test -count=1 -run 'FolderStructure|HandlersRegister' ./internal/toolboxcommands/
