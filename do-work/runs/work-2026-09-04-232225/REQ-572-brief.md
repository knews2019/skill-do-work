# Builder brief — REQ-572

You are the implementation builder for one request in the do-work skill repository. You work ONLY inside your own worktree and commit ONLY on your own branch. The orchestrator merges.

## Where you work

- Worktree (your working directory for every command): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row`
- Branch (already checked out there): `worktree-agent-REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row`
- Go module inside the worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row/skills/do-work-board/tools/queue-kanban` (own go.mod; run go commands from there)
- Main tree (READ ONLY, except the one hand-back file below): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`

## Read first, in this order (absolute paths; read the request and run files from the MAIN tree, never from the worktree's stale do-work/ snapshot)

1. The request: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row.md` — `## What`, `## Detailed Requirements`, `## Constraints`, `## Builder Guidance`, `## Red-Green Proof`, `## Exploration`, and `## Scope` (your write boundary).
2. The exploration in full: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-04-232225/REQ-572-exploration.md` (line-numbered findings, tests to rewrite, restatements to update, concerns C1–C6).
3. Crew rules (in the worktree, shipped files): `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md` (tdd: true).
4. Primes: `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`; read `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` and `_dev/primes/lessons-kanban-board.md` because you change code the prime's Read-first names (both were dropped from `required_lessons` at capture for budget; report that you read them).

## Decisions already made (do not reopen)

- One Activity row per parseable lifecycle stamp; `lifecycleTimestampFields` in `model.go` stays the only stamp enumeration and `model.go` is NOT yours (REQ-571 owns it).
- Sort newest first; tie-break `RequestId` descending then `StampField` as a third key so two stamps of one REQ at one instant order deterministically (exploration C1). No new stamp-order enumeration.
- The row payload shape (`id`, `stampField`, `stampAt`, `transition`) is unchanged. `data-activity-request` stays on every `<tr>` even though it is no longer unique (REQ-573 needs it).
- Near-duplicate rows at one instant ("completed" and "status changed to completed") ship as-is; no suppression rule.
- Summary line reports both counts, e.g. "38 transitions across 21 REQs in the last 24 hours" (singular forms for 1). Keep the "(N before filters)" clause and say in a comment that it counts transitions. Reword both empty states so they count transitions, not REQs.
- No "latest only" toggle. No CSS change, no drawer/click work (REQ-573). Release files (CHANGELOG.md, skills/do-work/CHANGELOG.md, VERSION) and every `do-work/` path are the orchestrator's.

## How to work

- **TDD is mandatory.** Write the RED tests FIRST in `activity_test.go`: the captured case (a ticket with `created_at: 2026-09-04T22:52:00Z` and `claimed_at: 2026-09-04T23:00:17Z` yields two rows, "claimed" then "captured"), the same-REQ/same-instant tie case, and the all-stamps pin (`len(rows)==len(declared)` in `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse`). Run them, record the exact failures, then implement. Rewrite the newest-only assertions the exploration names (keep their fixtures and intent). Then add the Node lane test in `javascript_behavior_c_test.go` following `javascript_behavior_a_test.go:420-450` (slice `activityRowsWithin`, `activityWindowPhrase`, `renderActivity`; stub `document.getElementById`, `boardData`, `viewState`, `requestsById`, `requestMatchesFilters`, `makeInstantWithRelativeNode`, `createElement`) that asserts the summary text for a payload with two rows for one REQ and one for another; run it with `QUEUE_KANBAN_JAVASCRIPT_PROBES=on`.
- **Write only inside `## Scope`'s file list.** If a file outside it is required by the REQ's own requirements, flag it in the hand-back and proceed; anything else is a scope expansion: stop and report it in the hand-back instead of editing.
- Update every restatement of "one row per REQ / newest stamp" the exploration lists (activity.go comments, generate.go:82-86, board-activity.js:2-4 and 53-55, template.html:460-467). Do not touch CHANGELOG.md.
- P-A-U: you cannot edit the REQ file (main tree). Record [PLAN], [APPLY], [UNIFY] evidence in the hand-back under `## P-A-U`. UNIFY = `git diff --stat` reviewed, `gofmt -l .` clean, `go vet ./...` clean, no debug artifacts, list of files verified.
- Log significant choices not dictated by the request as `## Decisions` entries D-01, D-02, … with reasoning (DECIDE & STATE vs ESCALATE with Value/Risk).
- Out-of-scope finds go to `## Discovered Tasks` in the hand-back; do not fix inline.
- Commit on your branch in small coherent commits with messages prefixed `[REQ-572] `. Do not push. Do not merge. Do not touch main.

## Verification before hand-back (all from the worktree)

```
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row/skills/do-work-board/tools/queue-kanban && gofmt -l . && go vet ./... && go test -count=1 ./...
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row && QUEUE_KANBAN_BROWSER_PROBES=off QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' .
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row/skills/do-work-board/tools/queue-kanban && go build -o /tmp/queue-kanban-572 . && /tmp/queue-kanban-572 generate --repo-root /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 --out /tmp/queue-kanban-572-site && grep -o '"id":"REQ-570"' /tmp/queue-kanban-572-site/board-data.js | wc -l
```

The last count must be at least 2 (REQ-570 has created_at and claimed_at, so it appears at least twice in the activity payload). Record each command's exit status and wall time. If node is not on PATH, say so; the Node lane test still has to be written.

## Hand-back (the ONE main-tree path you may write)

Write `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-04-232225/REQ-572-handback.md` (absolute path; never stage or commit it). Sections, each under its own `##` heading:

- `## Branch` — branch name and final commit hash (`git rev-parse HEAD` in the worktree).
- `## File manifest` — every file created/modified/deleted with the verb, tests included.
- `## P-A-U` — the three phases with evidence.
- `## Red-green evidence` — per RED test: name, failure before (exact message), pass after.
- `## Tests rewritten` — each old assertion changed and why.
- `## Verification` — each command, exit status, wall seconds.
- `## Lessons read` — which required/touch-conditional lesson satellites you read.
- `## Decisions` — D-01… entries.
- `## Discovered Tasks` — or "None".
