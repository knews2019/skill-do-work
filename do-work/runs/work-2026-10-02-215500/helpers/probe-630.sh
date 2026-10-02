#!/usr/bin/env bash
set -euo pipefail
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2
! grep -q "calendar and a testing track" skills/do-work-board/docs/board-guide.md
! grep -q "Kanban board + queue activity calendar" skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
grep -q "other pages in its page switcher" skills/do-work-board/docs/board-guide.md
grep -q "other pages in its page switcher" skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
