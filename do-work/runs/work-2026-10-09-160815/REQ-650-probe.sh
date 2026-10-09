#!/bin/bash
# Focused checks for REQ-650: the device dedupe is gone and the disk-space tests pass.
set -euo pipefail
cd skills/do-work-board/tools/queue-kanban
if grep -q measuredDevices verify.go; then echo "measuredDevices still present"; exit 1; fi
go test -run 'DiskSpace|LowDisk|FreeDisk' .
