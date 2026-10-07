#!/usr/bin/env bash
# Focused tests for REQ-636: request activity, the served response, verify's worktree and disk-space probes.
set -euo pipefail
root="$(git rev-parse --show-toplevel)"
go test -C "$root/skills/do-work-board/tools/queue-kanban" -count=1 -run 'Activity|Serve|Verify|Worktree|DiskSpace' .
