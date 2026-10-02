#!/usr/bin/env bash
set -euo pipefail
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
guide=skills/do-work-board/docs/board-guide.md
! grep -q "Board / Calendar / Testing" "$guide"
for fragment in '#board`' '#board/by-ur' '#board/urs-only' '#activity' '#calendar' '#timeline' '#durations' '#testing'; do grep -qF "$fragment" "$guide"; done
grep -q "free disk space" "$guide"
grep -q "free disk space" skills/do-work-board/actions/board.md
