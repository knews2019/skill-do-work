#!/bin/bash
# Focused tests for REQ-649: earmark parsing and placement in the Kanban board tool, plus the tooltip text pin.
set -euo pipefail
grep -q 'belongs under Needs input' skills/do-work-board/tools/queue-kanban/web/board-cards.js
cd skills/do-work-board/tools/queue-kanban
go test . -run 'AssignedTo|AssignedPending|WriteSetOverlapBadgeRenderPath' -count=1
