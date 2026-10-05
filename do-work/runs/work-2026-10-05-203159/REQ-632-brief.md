# Builder brief — REQ-632 (board cards show last correlated activity and drop the assumed-pause badge)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-632-board-cards-last-correlated-activity
- Branch: worktree-agent-REQ-632-board-cards-last-correlated-activity (commit here; never merge, never push)
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-632-board-cards-last-correlated-activity.md — read What, Detailed Requirements, Constraints, Red-Green Proof, and the Exploration and Scope sections at its end.
- Exploration findings: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-05-203159/REQ-632-exploration.md
- Hand-back file (the one main-tree path you may write, never commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-05-203159/REQ-632-handback.md

## Rules to load first
Read: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md (tdd: true), backend.md if present; prime files /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-kanban-board.md, /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/prime-do-kanban.md. Lessons satellites were dropped for budget; still read the entries the prime's Traps point to when you touch that code (lessons-do-kanban.md, _dev/primes/lessons-kanban-board.md: REQ-284 never cache verify probes; REQ-422 a cache hit rebuilds every wall-clock field from one instant).

## Write boundary
Only the files in the REQ's ## Scope list (paths relative to the worktree). Never write anything under do-work/ in the worktree. Never touch VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/version.md (release is the orchestrator's). Needing another file: stop and say so in the hand-back.

## Decisions already taken by the orchestrator (from exploration)
- Card rule is status claimed (there is no planning/dispatched/integrating status); phase text = label of the newest union event.
- Ship activity as a separate top-level payload field keyed by request id (do not write into the cached Requests map in serve).
- One now per board-data response, shared by the verify probes and the activity collector; generate uses its own single now.
- No do-work pathspec on git log. Branch listing goes through the injected runner.
- The GREEN grep is narrowed to the badge (over 4h · assumed pause, implementationSpanPaused); Panel B and UR-summary "assumed pause" text belongs to REQ-633 — leave it.
- Assert real formats (formatElapsedDuration "1h 07m", makeInstantWithStopwatchNode output).

## TDD
RED first: write the failing tests (attribution pins, the 4h21m/67-minute case, open card data-instant-ms, drawer row) against a compiling stub, record each test's failing message, then implement to GREEN. Do not commit the stub.

## Verify before hand-back
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-632-board-cards-last-correlated-activity/skills/do-work-board/tools/queue-kanban && gofmt -l . && go vet . && QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_BROWSER_PROBES=off go test -count=1 ./...  (about 60s; every test file must stay under 30s).
Render and look: build the binary in the worktree, run generate against /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 (--repo-root /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 --out a scratch dir outside both trees), open nothing interactive — read the generated board-data.js activity values for a claimed REQ (REQ-632 itself is claimed) and one done REQ. Time /board-data.js with serve on a free port before (main tree binary at HEAD) and after; report both numbers. Kill any server you start.

## Hand-back format (/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-05-203159/REQ-632-handback.md)
- Branch and final commit hash(es)
- File manifest (each path, new/modified/deleted)
- Red-green evidence: test name, failing message before, pass after
- P-A-U: [PLAN] approach, [APPLY] scope note, [UNIFY] git diff --stat, linters run, each file checked
- ## Decisions (D-01..., each DECIDE & STATE or ESCALATE with Value/Risk)
- ## Discovered Tasks (out-of-scope finds; do not fix inline)
- Lessons read (whole or family) and missing ones
- Render evidence and timing numbers
- Integration seams (none expected)
Keep your final chat reply under 3000 characters; the file holds the detail.
