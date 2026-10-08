#!/bin/bash
# Focused tests for REQ-646: request activity correlation in the Kanban board tool.
set -euo pipefail
cd skills/do-work-board/tools/queue-kanban
go test . -run 'RequestActivity|Correlat|Ancestry' -count=1
