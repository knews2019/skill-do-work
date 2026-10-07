#!/bin/bash
# Focused tests for REQ-643: assigned_to placement and the pending-column buckets in the board model.
set -euo pipefail
cd skills/do-work-board/tools/queue-kanban
go test . -run 'AssignedTo|Earmark|PendingColumns|Bucket|OpenWork' -count=1
