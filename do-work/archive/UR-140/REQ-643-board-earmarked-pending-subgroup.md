---
id: REQ-643
title: 'Board shows an earmarked pending REQ under Pending → Earmarked instead of Ready, and the ready counts agree'
status: completed
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 9-file write set
  - 3 subsystems involved
  - 8 acceptance criteria
  calculated_at: 2026-10-07T23:28:10Z
route: B
created_at: 2026-10-07T23:21:03Z
user_request: UR-140
domain: general
prime_files: ["_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-644]
batch: earmark-placement
write_set: [skills/do-work-board/tools/queue-kanban/model.go, skills/do-work-board/tools/queue-kanban/model_test.go, skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/open_work.go, skills/do-work-board/tools/queue-kanban/open_work_test.go, skills/do-work-board/tools/queue-kanban/main.go, skills/do-work-board/tools/queue-kanban/web/board-cards.js, skills/do-work-board/docs/board-guide.md, skills/do-work-board/actions/board.md, skills/do-work/actions/work-reference.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, skills/do-work-board/tools/queue-kanban/web/board.css, skills/do-work/docs/work-guide.md]
dispatch_at: 2026-10-07T23:31:21Z
builder_handback_at: 2026-10-07T23:47:22Z
integration_at: 2026-10-07T23:47:36Z
review_at: 2026-10-07T23:58:16Z
kb_status: pending
claimed_at: 2026-10-07T23:26:38Z
commit: 5b2c3e64f381fbc7f69ea6f24a5ba89d5195e733
heavy_verified_at: 2026-10-08T00:06:45Z
heavy_verified_revision: 5b2c3e64f381fbc7f69ea6f24a5ba89d5195e733
completed_at: 2026-10-08T00:07:24Z
release_at: 2026-10-08T00:07:24Z
---
# Board Shows an Earmarked Pending REQ Under Pending → Earmarked Instead of Ready
## What
A `pending` REQ with a non-empty `assigned_to` and no unmet dependencies must not count as Ready anywhere the board reports readiness. It lands in a third Pending sub-group, "Earmarked", after Ready and Waiting, and the CLI summary and open-work digest stop counting it as ready.
## Why
"Ready" tells the user the run will take the REQ next. The run scan skips an assigned REQ as a courtesy (`skills/do-work/actions/work-reference.md` exit-summary section 6 and the auto-wave predicate), so the board and the pipeline disagree on the one thing the Ready label promises. The maintainer looked for two earmarked REQs under Needs Input · Blocked and found nothing.
## Verified Facts (from triage)
- `skills/do-work-board/tools/queue-kanban/model.go` `bucketColumns` (around line 1719) buckets by status and `UnmetDependencies` only; `AssignedTo` is never read there. The `PendingReady` comment (line 382) promises "actionable now".
- The same bucket feeds `main.go:123` ("ready to work") and `open_work.go:44` (the digest's `N ready`), so three surfaces call a skipped REQ ready.
- The documented decision (`model.go` `AssignedTo` comment, `work-reference.md` schema line "no column logic and no scheduling", `web/board.css` ".badge-assigned" comment "somebody's intention, not a state of the work") was about never making the field a gate or a scheduler. A sub-group gates nothing.
- Needs Input · Blocked is defined as the operator-actionable inbox (`skills/do-work-board/actions/board.md:92`, `isNeedsInputOrBlockedStatus` in `model.go`). A REQ reserved for another session is not actionable by this operator, and that bucket is status-keyed, so routing a `pending` ticket there would break its contract and the open-work counts.
- `TestAssignedToNeverAffectsColumnPlacement` (`model_test.go:1020`) only asserts that the assigned and unassigned tickets share a `Status`. It would stay green under either placement option, so it pins nothing about columns.
- The frontend already renders Pending sub-groups generically through `makePendingGroup` (`web/board-cards.js:497-515`); `fillPendingColumn` collapses to a flat list only when Waiting is empty.
## Detailed Requirements
1. `model.go`: add a `PendingEarmarked` bucket to `BoardColumns` for tickets with `Status == "pending"`, no unmet dependencies, and non-empty `AssignedTo`. Such a ticket leaves `PendingReady`; `Pending` stays the union (Ready, Waiting, Earmarked) and the priority sort applies to the new group as it does to the other two. A `pending` ticket with unmet dependencies stays in Waiting whether or not it is assigned (Waiting already says it is not actionable). Update the `BoardColumns` field comments and the `AssignedTo` doc comment (today: "The board never buckets, orders, or schedules on it") so they describe the new placement truthfully: the board groups on it, and still never orders across groups, schedules, or gates on it.
2. `skills/do-work/actions/work-reference.md`, the `assigned_to:` schema line: replace "with **no column logic and no scheduling**" with the true statement (the board shows an assigned pending REQ in the Pending column's Earmarked group instead of Ready; no scheduling), in the same commit as the `model.go` change. Keep every existing heading unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings.
3. `generate.go`: thread `pendingEarmarked` through the payload struct and its single copy site (the `PendingReady`/`PendingWaiting` pattern at lines 142-146 and 776).
4. `web/board-cards.js`: `fillPendingColumn` renders a third `makePendingGroup("Earmarked", …)` after Waiting when the group is non-empty; the flat-list shortcut applies only when both Waiting and Earmarked are empty. The badge keeps its current text and verbatim value; nothing about the badge changes.
5. `main.go` summary: "ready to work" counts `PendingReady` only; add an `earmarked` line beside "waiting on deps". `open_work.go`: the digest's pending breakdown reads `(N ready, N waiting, N earmarked)` and `countOpenWork` exposes the new count.
6. `model_test.go`: replace `TestAssignedToNeverAffectsColumnPlacement` with a test that writes a pending unassigned REQ and a pending assigned REQ with no dependencies through the normal board build and asserts the assigned one is in `PendingEarmarked` and not in `PendingReady`, the unassigned one is in `PendingReady`, and a pending assigned REQ with an unmet dependency is in `PendingWaiting`. Keep the existing verbatim-read and absent-reads-as-empty tests.
7. Docs: `skills/do-work-board/docs/board-guide.md:18` (Pending splits into Ready and Waiting) and `skills/do-work-board/actions/board.md:111-112` (digest and summary descriptions) name the third group.
8. Board version bump and changelog entry per `_dev/primes/prime-kanban-board.md` and `_dev/primes/prime-releases.md` (the release step of the run owns the mechanics).
## Constraints
- Placement only. The `assigned_to` value stays verbatim-read with no alias map and no case folding; the badge text and the drawer row are unchanged.
- No scheduling and no gating on `assigned_to` anywhere on the board. The one behavioral reader stays the work pipeline's scan.
- `model.go` and the `work-reference.md` schema line change in the same commit (lock-step rule, stated on both).
- `tdd: true`: the Go test in requirement 6 is written first and seen failing before the model change.
## Dependencies
None. REQ-644 (capture distinguishes a session earmark from operator-present work) is related and touches disjoint files; the two may run in parallel.
## Builder Guidance
Certainty is high on the decision (sub-group, not Needs Input routing) and on the surfaces to change; the triage listed them from the code. Latitude: the exact label text ("Earmarked") and the summary line wording. Do not add an `earmarked` filter, sort, or any new behavior beyond placement and counts.
## Red-Green Proof
**RED prompt/case:** A queue with one REQ: `status: pending`, no `depends_on`, `assigned_to: 'user-interactive'`. Build the board and run `queue-kanban summary`.
**Why RED now:** The card renders under Pending → Ready with the "assigned user-interactive" badge, and the summary says `ready to work : 1`, while `do-work run` would skip it.
**GREEN when:** The card renders under Pending → Earmarked (1), Ready shows 0, the summary says `ready to work : 0` and `earmarked : 1`, and the open-work digest reads `(0 ready, 0 waiting, 1 earmarked)`. The new Go test pins the placement.
**Validation:** User confirmed. The repro and expected/actual are the maintainer's own words in the verbatim input.
## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 7636 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the index row governs queue-kanban model, parser, and UI changes, which this REQ makes; family `paired-predicate-drift` (two readers of one contract) is the shape of this bug.
- `_dev/primes/lessons-kanban-board.md` — 5912 tokens, `slugged: partial`; over budget. Matching reason: its owning prime governs the board tool this REQ edits.
## Full Context
See `do-work/user-requests/UR-140/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (from the builder hand-back) Read primes (`prime-kanban-board.md`, `prime-do-kanban.md`), rules (general, coding-guardrails, shared-principles, communication-style, testing), `paired-predicate-drift` lessons. Approach: TDD test through `buildBoard`; then thread one new bucket through model → payload → counts → summary → renderer → docs, following the existing `PendingWaiting` pattern exactly.
- [x] **[APPLY]:** (from the builder hand-back) Done as planned in the ten Scope files. One extra wording change inside Scope: the Ready group's empty text became "Nothing ready — everything here is waiting or earmarked" (D-02).
- [x] **[UNIFY]:** (from the builder hand-back) `git diff --stat` at commit: 10 files, 98 insertions, 60 deletions. `gofmt -l .` empty; `go vet ./...` clean; `git diff --check` clean. Checked each Scope file against the acceptance criteria; grep found no other reader of `pendingReady`/`pendingWaiting` in `web/` (column copy reads the rendered DOM via `visibleRequestIdsForColumn`, so it picks up the Earmarked cards). No debug artifacts.
*Source: maintainer feedback "Board: a pending REQ earmarked for the user reads as 'Ready', which misleads", change 1, accepted by validate-feedback.*
---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fixed at capture (a Pending → Earmarked sub-group, counts agree, lock-step comment edits) but the change crosses the Go model, the JSON payload, two CLI count surfaces, the browser renderer, docs, and a schema line, with `tdd: true`; exploration confirms the exact threading pattern and the test shape before a builder touches nine files.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Key files and the threading pattern a display bucket follows (verified against HEAD f45fb549 and the capture triage):

- `skills/do-work-board/tools/queue-kanban/model.go` — `BoardColumns` (line ~378) holds `Pending`, `PendingReady`, `PendingWaiting`, `Claimed`, `NeedsInputOrBlocked`, `RecentlyDone`, `CompletionAnomalies`. `bucketColumns` (line ~1714) switches on status: the `pending` arm splits Ready/Waiting on `len(ticket.UnmetDependencies)`; the two pending groups are then priority-sorted in a loop over `[][]*RequestTicket{PendingReady, PendingWaiting}` and `Pending` is rebuilt as their concatenation. `AssignedTo` is parsed at line ~889 (`coerceScalarToString(fields["assigned_to"])`) with the verbatim-read contract in the struct comment (lines ~176-191), which today says the board never buckets on it.
- `skills/do-work-board/tools/queue-kanban/generate.go` — `generatedColumns` payload struct (lines ~142-146: `PendingReady []string \`json:"pendingReady"\``, `PendingWaiting`), single copy site at line ~776 using `requestIdsOf(board.Columns.X)`.
- `skills/do-work-board/tools/queue-kanban/open_work.go` — `countOpenWork` (line ~34-44) exposes `PendingReady`/`PendingWaiting` counts; `writeOpenWorkDigest` (line ~58-61) prints `pending %d (%d ready, %d waiting) | claimed %d | needs-input/blocked %d`. Its test (`open_work_test.go`, if present) asserts the exact line; `summary`'s count block is in `main.go:123` (`ready to work     : %d`) with `waiting on deps` beside it.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` — `fillPendingColumn(readyIds, waitingIds, totalCount)` (line ~480) renders a flat list when `waitingIds` is empty, else `makePendingGroup("Ready", …)` + `makePendingGroup("Waiting", …)` (lines ~497-498). `makePendingGroup(labelText, requestIds, emptyText)` (line ~501) is generic: a third call is the whole UI change. `renderColumns` (line ~563) reads `columns.pendingReady` / `columns.pendingWaiting` from `boardData.columns`. The assigned badge is built at lines ~219-239 and stays untouched.
- `skills/do-work-board/tools/queue-kanban/model_test.go` — `TestAssignedToNeverAffectsColumnPlacement` (line ~1020) asserts only `Status` equality through `parseRequestTicket`; it never calls `bucketColumns`. Existing placement tests use `requestIdsOf(board.Columns.PendingReady)` / `requestIdSet(...)` after a full board build (lines ~507, ~577, ~649, ~854) — follow that shape. Tests at ~864-908 (verbatim read, absent-reads-as-empty) stay.
- Docs restating Ready/Waiting: `skills/do-work-board/docs/board-guide.md:18`, `skills/do-work-board/actions/board.md:111-112`. Schema line to change in lock-step: `skills/do-work/actions/work-reference.md:114` ("with **no column logic and no scheduling**").
- Testing conventions: `go test ./...` in the board package takes ~40 s (over the 30 s focused-probe budget), so the probe runs `-run` with the new and neighbouring test names; the repository gate (`_dev/tests/maintainer-verify.sh`) covers the whole package. The board has a browser heavy lane (`queue-kanban-browser`, needs `QUEUE_KANBAN_BROWSER`) and a javascript lane for `web/`.
- Lessons consulted: `do-work/lessons-index.md` rows for `lessons-do-kanban.md` (7636 tokens) and `_dev/primes/lessons-kanban-board.md` (5912 tokens) both exceed the 2000-token budget and are `slugged: partial`, so no targeted form is legal; `required_lessons` stays absent and the captured drop record stands. The family `paired-predicate-drift` (two readers of one contract) is this bug's shape and the builder reads that bullet directly via the prime's Traps.

*Generated by the orchestrator from the capture-time code read (same session), standing in for a separate Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — `PendingEarmarked` bucket in `BoardColumns` and `bucketColumns`; struct comments and the `AssignedTo` doc comment
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modify) — replace `TestAssignedToNeverAffectsColumnPlacement` with a placement test through the board build
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — `pendingEarmarked` payload list and copy site
- `skills/do-work-board/tools/queue-kanban/open_work.go` (modify) — earmarked count and digest line
- `skills/do-work-board/tools/queue-kanban/open_work_test.go` (modify, if it pins the digest line) — expected text
- `skills/do-work-board/tools/queue-kanban/main.go` (modify) — summary block gains an `earmarked` line; `ready to work` excludes earmarked
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — third `makePendingGroup("Earmarked", …)` after Waiting; flat-list shortcut only when Waiting and Earmarked are both empty
- `skills/do-work-board/docs/board-guide.md` (modify) — Pending split sentence names the third group
- `skills/do-work-board/actions/board.md` (modify) — open-work and summary descriptions name the third group
- `skills/do-work/actions/work-reference.md` (modify) — `assigned_to:` schema line: replace "no column logic and no scheduling" with the true placement statement, same commit as `model.go`
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify, added by D-05) — restatements of the bucketing rule name the Earmarked placement
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify, added by D-05) — comments only: the sub-group and `.badge-assigned` comments
- `skills/do-work/docs/work-guide.md` (modify, added by D-05) — the earmarking paragraph says where the card shows

**Files I will NOT touch:** `web/board.css` (the badge and `.pending-group` styles already cover a third group), `web/board-detail.js` (drawer row unchanged), `frontmatter.go`, any `skills/do-work/tools/do-work-cli` code, `VERSION`, `CHANGELOG.md` and mirrors (release step owns them), `do-work/` paths.

**Acceptance criteria (restated from REQ):**
- [ ] A `pending` REQ with non-empty `assigned_to` and no unmet dependencies is in `PendingEarmarked`, not `PendingReady`; `Pending` is still the union of the three groups and each group is priority-sorted
- [ ] A `pending` assigned REQ with an unmet dependency stays in `PendingWaiting`
- [ ] The board payload carries `pendingEarmarked` and the Pending column renders an "Earmarked" group after Waiting; the badge text and value are unchanged
- [ ] `queue-kanban summary` reports `ready to work` without earmarked REQs and adds an earmarked line; the open-work digest reads `(N ready, N waiting, N earmarked)`
- [ ] The replaced Go test fails before the model change and passes after (tdd: true)
- [ ] `model.go`'s `AssignedTo` comment and `work-reference.md:114` say the board groups on the field, and never orders across groups, schedules, or gates on it; both change in the same commit
- [ ] `board-guide.md` and `board.md` name the third group
- [ ] No scheduling, filtering, or sorting on `assigned_to` beyond placement; no alias map or case folding

## Pre-Flight

**Git:** ✓ clean apart from `do-work/` (the two working REQs, `baseline.json`, and this run's directory; sibling REQ-644's builder runs in its own worktree)
**Tests baseline:** ✓ focused probe `do-work/runs/work-2026-10-07-232637/REQ-643-probe.sh` (go test -run 'AssignedTo|Earmark|PendingColumns|Bucket|OpenWork') exit 0 in 3.6 s; repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 in 125 s at 531dd00e
**Dependencies:** ✓ Go toolchain present; none to install

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/open_work.go` (modified)
- `skills/do-work-board/tools/queue-kanban/open_work_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/main.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/actions/board.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified, review fix)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified, review fix — comments only)
- `skills/do-work/docs/work-guide.md` (modified, review fix)

**What was done:** `bucketColumns` now puts a `pending` REQ with no unmet dependency and a non-empty `assigned_to` into a new `PendingEarmarked` group instead of `PendingReady`; the group flows through the board payload (`pendingEarmarked`), the `summary` block (`earmarked` line), the open-work digest (`N earmarked`), and the browser's Pending column (an Earmarked sub-group after Waiting). `TestAssignedToNeverAffectsColumnPlacement` was replaced by `TestAssignedPendingRequestIsEarmarkedNotReady`, which checks placement through the full board build; the `model.go` `AssignedTo` comment and the `work-reference.md` `assigned_to:` line changed together in builder commit `3393d267`. A review-fix commit `fa508a23` on the builder branch (re-merged as `5b2c3e64`) rewrote the stale restatements the review's sweep found (F1–F7): the badge tooltip and comment, `board.md`'s parser rule, `prime-do-kanban.md`, the payload and CSS comments, the two-group priority wording in `work-reference.md` and `model.go`, and `work-guide.md`'s earmarking paragraph.

## Decisions

*(from the builder hand-back)*

- D-01 DECIDE & STATE: Group order in `Pending` and on screen is Ready, Waiting, Earmarked (as the brief says). The test pins the union order `[REQ-562 REQ-564 REQ-563]`.
- D-02 DECIDE & STATE: The Ready group's empty text was "Nothing ready — everything here is waiting". With an Earmarked group that text can be false, so it now reads "… waiting or earmarked". No test pins the string (grep). Reversible.
- D-03 DECIDE & STATE: Waiting now renders only when non-empty, because the grouped branch can now be reached with Waiting empty (only Earmarked present). Before, that branch always had a non-empty Waiting, so nothing changes for existing boards.
- D-04 DECIDE & STATE: `open_work_test.go`'s partition check (ready + waiting = pending) now includes earmarked. Otherwise a fixture with an earmarked ticket would break it (the paired-predicate-drift shape).
- D-05 DECIDE & STATE (integrator, after review): fixed review findings F1–F7 instead of leaving them report-only. They are stale restatements of the meaning this REQ changed, and F1 is a tooltip that tells the user the board never acts on `assigned_to` on a card the board just moved because of it. Scope grows by three files (`prime-do-kanban.md`, `web/board.css` comments, `docs/work-guide.md`), mirrored in `write_set`. The badge value and visible text stay unchanged. Not fixed: F8 (test gap), N1 (ADR is a decision record), N2 wording difference (meaning matches). Reversible text edits.

## Discovered Tasks

*(from the builder hand-back; impact stamped by the integrator)*

- **impact-negligible** `web/board.css` `.badge-assigned` comment ("somebody's intention, not a state of the work") is still true, but it sits next to wording about "no column logic" in nearby docs. The builder did not check it in depth because the file is out of bounds. → report only
- **impact-negligible** The renderer change has no JS behavior test pinning the Earmarked group. The browser lane or a node behavior test could pin it. → report only

## Qualification

**Gate records (advance at `fd5dbd3b..1e6d9909`):** `qualify` satisfied with no findings (no debug artifacts, P-A-U boxes ticked). `scope-drift` raised 11 `SCOPE-DECLARED-NOT-TOUCHED` warnings, and every one names a backticked identifier from a Scope description (`AssignedTo`, `BoardColumns`, `PendingEarmarked`, `bucketColumns`, `makePendingGroup("Earmarked", …)`, `ready to work`, …), not a file. All ten declared files are in the Implementation Summary and in the diff. Judged as parser noise, not drift.

**Scope comparison:** declared 10 files, touched the same 10 (`git diff fd5dbd3b..1e6d9909 --stat`: 10 files, 98+/60-). `open_work_test.go` was declared conditionally ("if it pins the digest line"); it does pin it, so the touch is in scope. No undeclared file. `web/board.css`, `web/board-detail.js`, VERSION and CHANGELOG are untouched as declared.

**Requirement trace (read against the merged files):**
1. `model.go` `bucketColumns` pending arm: non-pending or unmet deps → Waiting, else non-empty `AssignedTo` → Earmarked, else Ready. The priority sort loop covers the three groups and `Pending` is Ready + Waiting + Earmarked. `BoardColumns` field comments and the `AssignedTo` doc comment now say the board places on the field and never orders across groups, schedules, or gates on it. ✓
2. `work-reference.md:114` `assigned_to:` line says the Earmarked placement, "placement only, with **no scheduling**", same builder commit `3393d267` as `model.go`; no heading changed. ✓
3. `generate.go`: `PendingEarmarked []string \`json:"pendingEarmarked"\`` plus its copy site. ✓
4. `board-cards.js`: `fillPendingColumn` gains `earmarkedIds`, renders an Earmarked group after Waiting when non-empty, flat list only when Waiting and Earmarked are both empty; the badge code (lines 219-239) is untouched. ✓
5. `main.go` adds `earmarked : %d`; `ready to work` reads `PendingReady`, which no longer holds assigned tickets. `open_work.go` digest reads `(N ready, N waiting, N earmarked)` and `countOpenWork` exposes `PendingEarmarked`. ✓
6. `TestAssignedToNeverAffectsColumnPlacement` replaced by `TestAssignedPendingRequestIsEarmarkedNotReady`, which goes through `buildBoard` and asserts Ready `[REQ-562]`, Earmarked `[REQ-563]`, Waiting `[REQ-564]` (assigned + unmet dependency) and the union order. Verbatim-read and absent-reads-as-empty tests are unchanged. ✓
7. `board-guide.md:18` and `board.md:111-112` name the third group. ✓
8. Release is the finalization step's job. Pending.

**Constraints:** the only new reader of `AssignedTo` in Go is the placement case (`model.go:1728`); `verify.go`'s existing reader is unchanged; no filter, sort, or alias map was added.

**Independent checks by the integrator:** (a) RED re-observed: a scratch copy of the merged package with the `AssignedTo` case removed fails `TestAssignedPendingRequestIsEarmarkedNotReady` with `pending ready = [REQ-562 REQ-563], want [REQ-562]`. (b) The Red-Green Proof case built from the merge (one pending REQ, `assigned_to: 'user-interactive'`): `summary` shows `ready to work : 0`, `earmarked : 1`; `open-work` shows `pending 1 (0 ready, 0 waiting, 1 earmarked)`.

**Re-qualification after the review fix (cumulative range `fd5dbd3b..5b2c3e64`, 13 files, 121+/81-):** `qualify` success with no findings. `scope-drift` repeats the same backticked-identifier warnings (now 12, adding `.badge-assigned`), all judged noise. The three new files are declared in Scope by D-05 and listed in the Implementation Summary. The re-merge hit one conflict in `skills/do-work/docs/work-guide.md:119`, where REQ-644 had appended a sentence to the same paragraph; the resolution keeps REQ-644's sentence and this REQ's Earmarked wording. The lock-step pair moved together again: `model.go` and `work-reference.md` are both in review-fix commit `fa508a23`.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at merge `1e6d9909` (repository gate, run directly); focused probe `do-work/runs/work-2026-10-07-232637/REQ-643-probe.sh` (`go test . -run 'AssignedTo|Earmark|PendingColumns|Bucket|OpenWork'`) run by `advance`.
**Result:** ✓ Gate exit 0 at `1e6d9909` (gate wall 117 s) and again at the review-fix re-merge `5b2c3e64` (gate wall 121 s); board package `go vet` + `go test ./...` green on the fix commit (55 s), citation contract PASS. First gate details: (gate wall 117 s, "Maintainer verification passed."); probe exit 0; `advance` green-gate and run-blocked-check satisfied. Builder branch: `go build/vet/test ./...` in the board package green (59.8 s), `gofmt -l .` empty, citation contract PASS.

**Red-green validation:** (traces to `## Red-Green Proof`)
- `TestAssignedPendingRequestIsEarmarkedNotReady` (`model_test.go`): ✗ before implementation — first a compile failure (`board.Columns.PendingEarmarked undefined`), then, with the field added and no bucketing change, `pending ready = [REQ-562 REQ-563], want [REQ-562]` → ✓ after the `bucketColumns` change (`ok ... 0.326s`). Builder evidence from the hand-back; the integrator re-observed the assertion-level RED on a scratch copy of the merged package with the `AssignedTo` case removed. The test uses the Proof's case (a pending REQ, no `depends_on`, `assigned_to` set) through the full board build, with `cloud-alpha` instead of `user-interactive` because the value is read verbatim and plays no part in placement.
- Proof's CLI outcome at the merge (integrator smoke, one REQ with `assigned_to: 'user-interactive'`): `ready to work : 0`, `earmarked : 1`, digest `pending 1 (0 ready, 0 waiting, 1 earmarked)`. The browser rendering is covered by the heavy lane below.

**New tests added:**
- `TestAssignedPendingRequestIsEarmarkedNotReady` (replaces `TestAssignedToNeverAffectsColumnPlacement`, which compared `Status` only and could not see columns)

**Existing tests updated (cross-REQ impact):**
- `model_test.go` `TestBlockedDependencyGateControlsColumnsAndInheritedCounts`: asserts `PendingEarmarked == 0` and the new `earmarked : 0` summary line — intentional output change.
- `open_work_test.go`: digest fragment `(1 ready, 0 waiting, 0 earmarked)` and the ready + waiting + earmarked = pending partition check — intentional output change.

**Heavy verification plan:**
- Range: fd5dbd3b462da2c24094339351aee3930af589be..5b2c3e64f381fbc7f69ea6f24a5ba89d5195e733 (re-planned after the review-fix re-merge; same three lanes as at `1e6d9909`)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — 9 changed paths matched subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — 9 changed paths matched subtree `skills/do-work-board/tools/queue-kanban`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — 13 changed paths matched subtree `skills`

*Verified by work action*

## Review

**Overall: 92%** | 2026-10-07T23:58:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 85% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `web/board-cards.js:233-237` assigned-badge tooltip ("Display only: the board never reorders, blocks, or hides on this") and comment `:224-225` ("never buckets") contradict the new Earmarked placement — impact-user-visible → report only
- F2 `skills/do-work-board/actions/board.md:120` parser lock-step rule still says `depends_on` drives a Ready/Waiting split, Pending partitions between those two groups, and the priority sort covers two groups — impact-rule-change → report only
- F3 `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:48` Stakes ("one derived presentation rule"), `:39` Trap, `:15` Read-first omit `assigned_to` as a placement input — impact-rule-change → report only

**Minor findings:** F4 `generate.go:184-185` `AssignedTo` payload comment "display only … never column … meaning" — impact-negligible → report only; F5 `web/board.css:1037-1038` (sub-groups "Ready / Waiting", rendered only when something waits) and `:1302` ("nothing on the board branches on it") — impact-negligible → report only; F6 priority restatements name two groups (`work-reference.md:120,248`, `model.go:224-225`) — impact-negligible → report only; F7 `skills/do-work/docs/work-guide.md:119` says only "the board shows an `assigned` badge", not the Earmarked placement — impact-user-visible → report only; F8 no test pins a nonzero earmarked count in summary/digest or the Earmarked group in the renderer — impact-negligible → report only. Nits: N1 ADR-018:88 "display-only" parse (decision record) — impact-negligible → report only; N2 `model.go:1705` unwrapped comment line, and the lock-step pair matches in meaning but not wording — impact-negligible → report only
**Acceptance:** Pass — the built CLI gives `ready to work : 0`, `earmarked : 1`, digest `(0 ready, 0 waiting, 1 earmarked)`, payload `pendingEarmarked: ["REQ-001"]`; headless Chrome renders Ready 0 + Earmarked 1; edge cases (whitespace value, blocked, unmet and met deps, priority) place correctly
**Restatement sweep:** redefined `assigned_to`'s board meaning (placement in Pending → Earmarked), the Pending column / Ready meaning (three groups, Ready excludes assigned), and the summary/digest count shapes — stale: F1–F7, N1. Inherited REQ-644 elements (earmark-vs-blocked boundary, External-condition look-alike list) re-checked at `work-reference.md:114,511`, `capture-reference.md:41`, `capture.md:107,108,143`, `work-guide.md:119`: still agree
**Suggested testing:** 2 items
**Follow-ups created:** None (10 findings report only)

*Reviewed by review-work action*

**Re-review of the fix delta (`1e6d9909..5b2c3e64`)** | 2026-10-08T00:03:47Z

**Acceptance:** Pass, 94% — F1–F7 now match `bucketColumns`; no new false statement; the badge value and visible text are unchanged; the lock-step pair still agrees; the `work-guide.md` conflict resolution keeps REQ-644's sentence. Report: `do-work/runs/work-2026-10-07-232637/REQ-643-review-delta.md`.
**Minor findings:** F8 (delta) the tooltip "The board only groups it under Pending → Earmarked" shows on every assigned card, including one under Waiting, where the field has no placement effect; the sentence states the only effect the board gives the field, so it stays — impact-user-visible → report only; F9 (delta) `model.go:1706` comment line is about 115 characters after the rewrap — impact-negligible → report only

## Lessons Learned

**What worked:** Threading one new bucket along the exact path `PendingWaiting` already took (model → payload → counts → summary → renderer → docs) kept the change mechanical, and the replacement test went through `buildBoard`, so it sees columns rather than one parsed field.
**What didn't:** The old test `TestAssignedToNeverAffectsColumnPlacement` compared `Status` only, so it could not see the mismatch it was named for. The first build also left eight restatements of the old "display only, never buckets" meaning (a user-facing tooltip, `board.md`'s parser rule, the prime, CSS and payload comments, `work-guide.md`); the wave-end restatement sweep found them, not the builder's Scope.
**Worth knowing:** A display label is a claim about what another component will do: Ready means "the run takes it next", so its predicate must match the run's scan, including the `assigned_to` skip. When a field's meaning changes, grep for "display only", "never buckets", and the old group list across `skills/`, not just the Scope files. An assigned REQ with an unmet dependency is still Waiting, so wording that says "every assigned card is Earmarked" is wrong.

## Orientation

Now the board's Pending column has three groups: Ready (the run takes it next), Waiting (unmet dependency), and Earmarked (assigned to a session, which the run's default scan skips), and `summary`/`open-work` count the same way. Lives in the queue-kanban board model and renderer (`prime-kanban-board.md`, `prime-do-kanban.md`). [MAP CHANGED] — the `assigned_to` contract is now a placement input on the board, not display only. Prime spot-check: both primes' referenced paths still exist; `prime-do-kanban.md` was updated in this REQ to name the new rule.

## Heavy Verification Plan

- Base: fd5dbd3b462da2c24094339351aee3930af589be
- Target: 5b2c3e64f381fbc7f69ea6f24a5ba89d5195e733
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — 9 changed paths matched subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — 9 changed paths matched subtree `skills/do-work-board/tools/queue-kanban`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — 13 changed paths matched subtree `skills`

## Heavy Verification Result

- Target revision: 5b2c3e64f381fbc7f69ea6f24a5ba89d5195e733
- Execution revision: 5b2c3e64f381fbc7f69ea6f24a5ba89d5195e733 (detached drain checkout, `QUEUE_KANBAN_BROWSER` set to Google Chrome)
- queue-kanban-javascript: exit 0, executed (fingerprint_mismatch), 7 s
- queue-kanban-browser: exit 0, executed (fingerprint_uncertain), 81 s — 37 tests, not skipped
- staged-skills: exit 0, executed (fingerprint_mismatch), 46 s

## Timing

Observed 2026-10-07T23:31:21Z to 2026-10-08T00:06:45Z: 35m 24s total, 32m 15s attributed across 7 events, 3m 09s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 16m 01s | 1 |
| review | 8m 58s | 2 |
| verification-gate | 7m 05s | 3 |
| handback-merge | 11s | 1 |

Slowest stage: builder-work / builder worktree build, 16m 01s, outcome success.
