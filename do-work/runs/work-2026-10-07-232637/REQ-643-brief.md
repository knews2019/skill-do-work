# Builder brief — REQ-643 (Board shows an earmarked pending REQ under Pending → Earmarked instead of Ready, and the ready counts agree)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-643-board-earmarked-pending-subgroup
- Branch: worktree-agent-REQ-643-board-earmarked-pending-subgroup, created from main HEAD with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-643-board-earmarked-pending-subgroup.md. Read it fully: What, Why, Verified Facts, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof, Exploration, Scope. The `## Scope` "Files I will touch" list is your write boundary.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-140/input.md.
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-07-232637/REQ-643-handback.md
- Route B, tdd: true, impact-user-visible, effort-substantive. Go model + payload + two CLI count surfaces + browser renderer + two docs + one schema line.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md (tdd: true) — same directory; domain: general. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-kanban-board.md and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/prime-do-kanban.md (read both before editing; note the parser lock-step rule and the versioning note: the release step owns VERSION/CHANGELOG, not you). Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md read the `[family: paired-predicate-drift]` bullets (two readers of one contract drifting is this bug's shape) and the prime's Traps; also skim /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-kanban-board.md Traps for the board.

## TDD order (mandatory; the orchestrator verifies red-then-green evidence)
1. In `model_test.go`, replace `TestAssignedToNeverAffectsColumnPlacement` with a placement test that builds a board (follow the shape of the existing placement tests that use `requestIdsOf(board.Columns.PendingReady)` / `requestIdSet(...)` after a full board build) over three fixtures: pending unassigned (expect PendingReady), pending `assigned_to: "cloud-alpha"` with no depends_on (expect PendingEarmarked, not PendingReady), pending assigned with an unmet depends_on (expect PendingWaiting). Also assert `Pending` is the union of the three groups. Run it, confirm it fails to compile or fails (RED), record the exact failure text.
2. Implement; run again (GREEN). Record test name, failure-before, pass-after in the hand-back.

## The change (decided at capture, do not reopen)
From the REQ's `## Detailed Requirements` 1-7 (8 is the release step's, not yours):
- `model.go`: `PendingEarmarked []*RequestTicket` on `BoardColumns`; in `bucketColumns`' pending arm, a `pending` ticket with no unmet dependencies and non-empty `AssignedTo` goes to `PendingEarmarked`; one with unmet dependencies stays in Waiting whether or not assigned; priority-sort the new group like the other two and rebuild `Pending` as Ready + Waiting + Earmarked. Update the `BoardColumns` field comments and the `AssignedTo` doc comment (today: "The board never buckets, orders, or schedules on it") to the truth: the board groups on it (placement only) and still never orders across groups, schedules, or gates on it; keep the verbatim-read, no-alias-map, trim-only contract text.
- `skills/do-work/actions/work-reference.md`, the `assigned_to:` schema line (today line 114): replace "with **no column logic and no scheduling**" with the true statement — the board shows an assigned pending REQ in the Pending column's Earmarked group instead of Ready, with no scheduling — and keep "both changing in the same commit". Same commit as `model.go`. Change nothing else on that line and no heading anywhere (`_dev/tests/shipped-package-reference-contract.sh` pins citation strings).
- `generate.go`: `PendingEarmarked []string \`json:"pendingEarmarked"\`` beside `pendingReady`/`pendingWaiting` and the copy site; update the comment that says Ready and Waiting partition Pending.
- `open_work.go` (+ `open_work_test.go` and the `model_test.go` open-work assertion at ~line 600 that pin the counts/line): `countOpenWork` gains `PendingEarmarked`; the digest reads `pending %d (%d ready, %d waiting, %d earmarked)`.
- `main.go` summary block: `ready to work` counts PendingReady only; add an `earmarked` line beside `waiting on deps` with matching alignment.
- `web/board-cards.js`: `fillPendingColumn` takes the earmarked ids too; render `makePendingGroup("Earmarked", earmarkedIds, "")` after Waiting when non-empty; the flat-list shortcut applies only when Waiting and Earmarked are both empty; `renderColumns` reads `columns.pendingEarmarked || []` and counts it in the column total and copy button. The assigned badge (lines ~219-239) and the drawer row stay exactly as they are.
- `docs/board-guide.md` line 18 and `actions/board.md` lines 111-112: name the third group in the Pending split and in the open-work / summary descriptions.
Constraints: placement only. No filter, sort across groups, scheduling, or gating on `assigned_to`; no alias map, no case folding (the trim-only verbatim read stays). Label text "Earmarked" and summary wording are your only latitude.

## Write boundary
Exactly the `## Scope` "Files I will touch" list in the REQ (`open_work_test.go` is included for the digest line). Never touch `web/board.css`, `web/board-detail.js`, `frontmatter.go`, anything under `skills/do-work/tools/`, VERSION, CHANGELOG.md or any mirror, `do-work/` (except the hand-back file), `_dev/`. If you need another file (for example another test pinning the digest line), it is the REQ's own test class so proceed, but name it in the hand-back as an addition to Scope; anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `cd skills/do-work-board/tools/queue-kanban && go build ./... && go vet ./... && go test ./... -count=1` (whole package; ~40 s).
- `gofmt -l .` in that directory prints nothing.
- `bash _dev/tests/shipped-package-reference-contract.sh` (citation strings; work-reference.md changed).
- Board smoke of the REQ's Red-Green Proof: build the tool, generate a board from a temp queue with one pending REQ `assigned_to: 'user-interactive'` and one plain pending REQ (see `prime-do-kanban.md` for the generate/summary invocation; the payload lands in `board-data.js`, not index.html), and show: `pendingEarmarked` holds the assigned id, `pendingReady` holds only the plain id, `summary` prints `ready to work : 1` and `earmarked : 1`, the open-work digest prints `(1 ready, 0 waiting, 1 earmarked)`. Paste the commands and outputs.
- If `QUEUE_KANBAN_BROWSER` is set in your environment you may run the browser lane's test; otherwise say it was not run (the integrator drains heavy lanes).
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit is preferred; the schema line and model.go must be in the same commit).
- File manifest: each file with (modified)/(new) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, gofmt/vet, and each file checked.
- Red-green evidence: test name, exact failure before, pass after.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md with a `[family: <slug>]` marker, or "none" with a reason, and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words). The orchestrator writes both.
- Integration seams (none expected).
- Test wall times and the smoke outputs.
