# Hand-back — REQ-643 (Board shows an earmarked pending REQ under Pending → Earmarked instead of Ready)

- Branch: `worktree-agent-REQ-643-board-earmarked-pending-subgroup`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-643-board-earmarked-pending-subgroup`
- Commits: `3393d267` (single commit; `model.go` and the `work-reference.md` schema line are in it together). Base: `531dd00e`.

## File manifest

- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — `PendingEarmarked` on `BoardColumns`; `bucketColumns` pending arm: unmet deps or non-pending → Waiting, else non-empty `AssignedTo` → Earmarked, else Ready; priority sort covers the new group; `Pending` = Ready + Waiting + Earmarked. Field comments, `bucketColumns` doc comment and `AssignedTo` doc comment updated (groups on it for placement only; never orders across groups, schedules, or gates; verbatim / no alias map / trim-only text kept).
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified) — `TestAssignedToNeverAffectsColumnPlacement` replaced by `TestAssignedPendingRequestIsEarmarkedNotReady`; the open-work count assertion and summary-line list in `TestBlockedDependencyGateControlsColumnsAndInheritedCounts` gain the earmarked count (0) and `    earmarked         : 0`.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — `PendingEarmarked []string \`json:"pendingEarmarked"\`` and its copy site; partition comment names three groups.
- `skills/do-work-board/tools/queue-kanban/open_work.go` (modified) — `openWorkCounts.PendingEarmarked`; digest line `pending %d (%d ready, %d waiting, %d earmarked)`.
- `skills/do-work-board/tools/queue-kanban/open_work_test.go` (modified) — expected digest fragment and the ready+waiting(+earmarked)=pending partition check.
- `skills/do-work-board/tools/queue-kanban/main.go` (modified) — summary adds `    earmarked         : %d` after `waiting on deps`; `ready to work` already read `PendingReady`, which no longer holds earmarked tickets.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified) — `fillPendingColumn(readyIds, waitingIds, earmarkedIds, totalCount)`; Earmarked group after Waiting when non-empty; flat list only when Waiting and Earmarked are both empty; `renderColumns` reads `columns.pendingEarmarked || []` into count and copy button. Badge and drawer untouched.
- `skills/do-work-board/docs/board-guide.md` (modified) — line 18 names Earmarked and when a card goes there.
- `skills/do-work-board/actions/board.md` (modified) — open-work and summary descriptions name the earmarked split.
- `skills/do-work/actions/work-reference.md` (modified) — `assigned_to:` line: "with **no column logic and no scheduling**" → "and a `pending` REQ with no unmet dependency shows in the Pending column's **Earmarked** group instead of Ready — placement only, with **no scheduling**"; "both changing in the same commit" kept; no heading changed.

No additions to Scope.

## P-A-U

- **[PLAN]:** Read primes (`prime-kanban-board.md`, `prime-do-kanban.md`), rules (general, coding-guardrails, shared-principles, communication-style, testing), `paired-predicate-drift` lessons. Approach: TDD test through `buildBoard`; then thread one new bucket through model → payload → counts → summary → renderer → docs, following the existing `PendingWaiting` pattern exactly.
- **[APPLY]:** Done as planned in the ten Scope files. One extra wording change inside Scope: the Ready group's empty text became "Nothing ready — everything here is waiting or earmarked" (D-02).
- **[UNIFY]:** `git diff --stat` at commit: 10 files, 98 insertions, 60 deletions. `gofmt -l .` empty; `go vet ./...` clean; `git diff --check` clean. Checked each file above against the REQ's acceptance criteria; grep found no other reader of `pendingReady`/`pendingWaiting` in `web/` (column copy reads the rendered DOM via `visibleRequestIdsForColumn`, so it picks up the Earmarked cards). No debug artifacts.

## Red-green evidence

Test: `TestAssignedPendingRequestIsEarmarkedNotReady` (`model_test.go`).

1. RED (test only): compile failure
   `./model_test.go:1052:39: board.Columns.PendingEarmarked undefined (type BoardColumns has no field or method PendingEarmarked)`
2. RED (struct field added, no bucketing change), assertion-level:
   `model_test.go:1050: pending ready = [REQ-562 REQ-563], want [REQ-562] — an assigned REQ is skipped by the run, so it is not Ready`
3. GREEN after the `bucketColumns` change: `ok  github.com/knews2019/skill-do-work/queue-kanban 0.326s`

## Verification (from worktree root)

- `go build ./... && go vet ./... && go test ./... -count=1` in `skills/do-work-board/tools/queue-kanban`: `ok ... 59.804s`, wall 1:01.5.
- `gofmt -l .`: no output.
- `bash _dev/tests/shipped-package-reference-contract.sh`: `PASS`, exit 0, 1.0 s.
- Browser lane: NOT run (`QUEUE_KANBAN_BROWSER` unset). The integrator drains it. The renderer change is untested in a real browser.

### Smoke (temp queue: REQ-001 pending `assigned_to: 'user-interactive'`, REQ-002 plain pending)

```
$ ./queue-kanban summary --repo-root repo
  pending             : 2
    ready to work     : 1
    waiting on deps   : 0
    earmarked         : 1
$ ./queue-kanban open-work --repo-root repo
queue-kanban open work: 2 open REQs
  pending 2 (1 ready, 0 waiting, 1 earmarked) | claimed 0 | needs-input/blocked 0
$ ./queue-kanban generate --out out --repo-root repo   # exit 0
board-data.js: "pending":["REQ-002","REQ-001"],"pendingReady":["REQ-002"],"pendingWaiting":[],"pendingEarmarked":["REQ-001"]
```

Scratch path: `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/7c114307-960e-4ea2-b6dd-ab1deb3c4539/scratchpad/req643-smoke`.

## Decisions

- D-01 DECIDE & STATE: Group order in `Pending` and on screen is Ready, Waiting, Earmarked (as the brief says). The test pins the union order `[REQ-562 REQ-564 REQ-563]`.
- D-02 DECIDE & STATE: The Ready group's empty text was "Nothing ready — everything here is waiting". With an Earmarked group that text can be false, so it now reads "… waiting or earmarked". No test pins the string (grep). Reversible.
- D-03 DECIDE & STATE: Waiting now renders only when non-empty, because the grouped branch can now be reached with Waiting empty (only Earmarked present). Before, that branch always had a non-empty Waiting, so nothing changes for existing boards.
- D-04 DECIDE & STATE: `open_work_test.go`'s partition check (ready + waiting = pending) now includes earmarked. Otherwise a fixture with an earmarked ticket would break it (the paired-predicate-drift shape).

## Discovered Tasks

- `web/board.css` `.badge-assigned` comment ("somebody's intention, not a state of the work") is still true, but it sits next to wording about "no column logic" in nearby docs. I did not check it in depth because the file is out of bounds. → report only
- The renderer change has no JS behavior test pinning the Earmarked group. The browser lane or a node behavior test could pin it. → report only

## Lessons read

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` `[family: paired-predicate-drift]` bullets (0.305.72, 0.305.9, 0.294.2).
- Traps in `prime-do-kanban.md` and `prime-kanban-board.md`; skimmed `_dev/primes/lessons-kanban-board.md`.

## Proposed lesson bullet (for lessons-do-kanban.md)

- [family: paired-predicate-drift] REQ-643: **a display label is a claim about another reader.** The board's Ready said "the run takes this next" while the run's scan skipped any REQ with `assigned_to`; the board's comments said it "never buckets" on the field, so nobody looked. The test that guarded the decision compared `Status` only and could not see columns. When a label promises what another component will do, pin it with a fixture that the other component treats differently, through the full board build.

## Proposed CHANGELOG entry

**Board puts an earmarked pending REQ under Pending → Earmarked, not Ready.** A pending REQ with `assigned_to` set and no unmet dependency now shows in a new Earmarked group in the Pending column. The run skips such a REQ, so the board no longer calls it ready. `queue-kanban summary` adds an `earmarked` line and its `ready to work` count no longer includes these REQs. The open-work digest reads `(N ready, N waiting, N earmarked)`. An assigned REQ that still waits on a dependency stays under Waiting.

## Integration seams

None. No file overlaps REQ-644's scope (capture docs).
