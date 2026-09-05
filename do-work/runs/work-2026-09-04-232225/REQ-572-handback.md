# Hand-back — REQ-572 (show every lifecycle transition of a REQ as its own Activity row)

## Branch

- Branch: `worktree-agent-REQ-572-show-every-lifecycle-transition-of-a-req-as-its-own-activity-row`
- Final commit: `2d8beb40b161bf4dfb8884b1335637e07835b183`
- Two commits on the branch, both prefixed `[REQ-572] `:
  - `6ed61142` emit one Activity row per lifecycle stamp (Go aggregation + Go tests)
  - `2d8beb40` count transitions and REQs in the Activity summary (client, comments, Node lane test)
- Not pushed, not merged, main tree untouched apart from this hand-back file (never staged, never committed).

## File manifest

All paths relative to the worktree root; every one is inside the REQ's `## Scope` list. No file outside Scope was written.

| Verb | Path | What changed |
|---|---|---|
| modify | `skills/do-work-board/tools/queue-kanban/activity.go` | `buildActivityRows` appends one row per parseable stamp from `lifecycleTimestampFields`; `newestLifecycleStamp` deleted (it had exactly one caller); `StampField` added as the third sort key; header, `ActivityRow` and function comments restated |
| modify | `skills/do-work-board/tools/queue-kanban/activity_test.go` | One new RED test, three rewritten assertions, one extended tie test (detail below) |
| modify | `skills/do-work-board/tools/queue-kanban/generate.go` | Two comment restatements only (`Activity` field at the payload struct, `generatedActivityEntry` doc). Payload shape and fill site untouched |
| modify | `skills/do-work-board/tools/queue-kanban/web/board-activity.js` | Summary reports transitions and distinct REQs; both empty states count transitions; header comment restated; a comment saying `data-activity-request` is deliberately non-unique |
| modify | `skills/do-work-board/tools/queue-kanban/web/template.html` | Activity section comment only |
| modify | `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` | New Node lane test `TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` appended |

`model.go` was not touched (REQ-571 owns it). `web/board.css`, `web/board-detail.js`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` and every `do-work/` path were not touched.

## P-A-U

**[PLAN]** — Read, in the brief's order: the REQ and the full exploration from the main tree; `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`; `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`; both lesson satellites. Approach settled before any edit:

1. RED tests in `activity_test.go` first (captured two-row case, same-REQ same-instant tie, all-stamps pin, plus the two rewritten newest-only assertions) → verify: run them, record exact failures.
2. `buildActivityRows` loops `lifecycleTimestampFields` directly, `newestLifecycleStamp` removed, `StampField` third sort key → verify: those tests turn green and the whole package stays green.
3. Client summary and empty states reworded to count transitions and distinct REQs → verify: a Node lane probe drives the shipped `renderActivity` and pins all four strings.
4. Comment restatements in `activity.go`, `generate.go`, `board-activity.js`, `template.html` → verify: `git diff` review, no other file touched.

No lesson contradicted the plan. The satellites' relevant warnings were honoured: the tie-break assertion reads the rule rather than restating a constant, and it was mutation-checked so it can actually fail (`lessons-kanban-board.md` REQ-318, REQ-322, REQ-293).

**[APPLY]** — Code written as planned, in two commits matching steps 1-2 and 3-4. Scope stayed at the six declared files.

**[UNIFY]** — `git diff --stat main...HEAD`:

```
 .../do-work-board/tools/queue-kanban/activity.go   |  91 ++++++------
 .../tools/queue-kanban/activity_test.go            | 163 ++++++++++++++-------
 .../do-work-board/tools/queue-kanban/generate.go   |  19 ++-
 .../queue-kanban/javascript_behavior_c_test.go     | 133 +++++++++++++++++
 .../tools/queue-kanban/web/board-activity.js       |  53 +++++--
 .../tools/queue-kanban/web/template.html           |  18 ++-
 6 files changed, 347 insertions(+), 130 deletions(-)
```

- `gofmt -l .` — no output (clean).
- `go vet ./...` — clean, exit 0.
- Debug artifact scan over the whole branch diff (`console.log`, `fmt.Print`, `t.Skip(`, `TODO`, `FIXME`, `XXX`, `debugger` on added lines) — no matches.
- Files reviewed line by line in the diff: `activity.go` (aggregation loop, three-key comparator, dead function removed, comments match the new rule), `activity_test.go` (every rewritten assertion still names the failure it pins), `generate.go` (comments only, no code line changed), `web/board-activity.js` (summary, empty states, comments; row build otherwise unchanged), `web/template.html` (comment only), `javascript_behavior_c_test.go` (new test appended, no existing test altered).
- `git status --porcelain` in the worktree — empty.

## Red-green evidence

RED run: `go test -count=1 -run 'TestBuildActivityRows|TestLifecycleTimestampFieldsIsTheOneListBothReadersUse' .`

| Test | Failure before | After |
|---|---|---|
| `TestBuildActivityRowsEmitsOneRowPerLifecycleStamp` (new; the REQ's captured case, `created_at: 2026-09-04T22:52:00Z` + `claimed_at: 2026-09-04T23:00:17Z`) | `activity_test.go:32: rows = 1, want 2 (one per parseable stamp): [{RequestId:REQ-570 StampField:claimed_at ... Transition:claimed}]` | PASS — two rows, `claimed` at 23:00:17 then `captured` at 22:52:00 |
| `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` (rewritten) | `activity_test.go:95: in-window rows = 5, want 8: [...]` | PASS — the eight rows in the exact interleaved order |
| `TestBuildActivityRowsOrdersOneTicketsStampsNewestFirst` (rewritten, was `...PicksTheNewestStampNotTheFirstDeclaredField`) | `activity_test.go:137: rows = 1, want 4 (one per parseable stamp): [...]` | PASS — `review_at, planning_at, claimed_at, created_at` |
| `TestBuildActivityRowsBreaksStampTiesDeterministically` (extended with the same-REQ/same-instant case) | `activity_test.go:229: rows = 1, want 4: [...]` | PASS |
| `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` (tail rewritten) | `activity_test.go:298: rows = 1, want 14 — one per declared lifecycle stamp: [...]` | PASS — 14 rows, each declared field appearing exactly once |
| `TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` (new Node lane test) | With `web/board-activity.js` reverted to its pre-change state: `javascript_behavior_c_test.go:2494: summary for three transitions of two REQs = "3 REQs touched in the last 24 hours", want both counts` | PASS — `"3 transitions across 2 REQs in the last 24 hours"` |

**Extra mutation check on the tie rule.** The tie test's RED above was a row-count failure, so it did not by itself prove the same-REQ ordering assertion bites. Replacing the third sort key with `return false` produced:

```
activity_test.go:233: same-instant row 2 = "claimed_at", want "blocked_at" (stamp field must break a same-REQ tie)
```

The fixture pairs `blocked_at` with `claimed_at` at one instant precisely because their alphabetical order contradicts their declaration order in `lifecycleTimestampFields`, so an implementation leaning on emit order or on `sort.Slice` stability fails. The mutation was reverted immediately (`git diff` clean before the commit).

## Tests rewritten

Behaviour change is intentional (the REQ asks for it), so per the cross-REQ test-break rule the prior REQ's tests were updated, not deleted. All four belong to REQ-568 (the Activity view's original capture).

1. `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` — was: five REQ ids in a fixed order, `len(inWindow)==5`, then a per-REQ map checking each newest stamp's transition. Now: the same five fixture tickets and the same instants, asserted as eight `(RequestId, StampField, Transition)` triples in exact order. Why: four of the five tickets carry two stamps, so the row count and order changed by design. The map-based transition checks were folded into the sequence assertion because a map keyed by REQ id can no longer hold a REQ's several rows.
2. `TestBuildActivityRowsPicksTheNewestStampNotTheFirstDeclaredField` → renamed `TestBuildActivityRowsOrdersOneTicketsStampsNewestFirst`. Was: `len(rows)==1` and `rows[0].StampField=="review_at"`. Now: four rows in the sequence `review_at, planning_at, claimed_at, created_at`. Why: the ticket has four parseable stamps. The anti-first-match intent is preserved and strengthened — declaration order would put `created_at` first, which the sequence rejects.
3. `TestBuildActivityRowsBreaksStampTiesDeterministically` — the original two-REQ input-order case is unchanged; a second stage was appended for one REQ carrying two pairs of same-instant stamps. Why: exploration C1, the defect the change introduces.
4. `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` — tail was `len(rows)==1` plus "the chosen field is declared". Now `len(rows)==len(declared)` plus each declared field appearing exactly once. Why: under the new rule the 14-stamp ticket yields 14 rows, which is a stronger anti-drift pin. A count alone would pass on a reader that emitted one field fourteen times, so the fields are counted individually.

Untouched and still passing as written: `TestBuildActivityRowsStraddlesTheWindowBoundary` and `TestBuildActivityRowsSkipsTicketsWithNoParseableStamp` (every fixture ticket there carries at most one stamp, so their positional indices stay meaningful — exploration C6).

## Verification

Every command run from the worktree. Wall seconds measured with the shell's `SECONDS`.

| # | Command | Exit | Wall |
|---|---|---|---|
| 1 | `cd .../queue-kanban && gofmt -l . && go vet ./... && go test -count=1 ./...` | 0 | 37 s (the `go test` leg alone reports 36.5 s) |
| 2 | `QUEUE_KANBAN_BROWSER_PROBES=off QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban -run '^TestJavaScriptBehavior' .` | 0 | 8 s — budget line: `wall=8s tests=56 slowest-file=javascript_behavior_c_test.go:2.71s limit=<30s` |
| 3 | `go build -o /tmp/queue-kanban-572 . && /tmp/queue-kanban-572 generate --repo-root <main tree> --out /tmp/queue-kanban-572-site && grep -o '"id":"REQ-570"' .../board-data.js \| wc -l` | 0 | 2 s — count **8** (the brief's floor is 2) |

`node` is on PATH (`/opt/homebrew/bin/node`, v22.23.2), so the JavaScript lane really ran rather than skipping; the strict-mode guard did not fire.

The grep in command 3 counts every `"id":"REQ-570"` in the file, including the requests map, the calendar and the timeline, so it overstates the activity rows. Parsing the payload gives the honest number: **REQ-570 carries 5 activity rows** — `builder_handback_at` 23:35:48, `dispatch_at` 23:15:08, `planning_at` 23:11:07, `claimed_at` 23:00:17, `created_at` 22:52:00 — against exactly one before this change.

**Exploration C3, payload growth, measured on the real tree:** the activity array holds 1980 rows across 550 distinct REQs (about 3.6 rows per REQ), 205,910 bytes serialized against 56,857 bytes for the newest-only equivalent. That extra 149 KB is **1.26 %** of the 11.9 MB `board-data.js`. Acceptable, as the exploration predicted.

**Not exercised:** no browser render was taken. The board's browser lane (`QUEUE_KANBAN_BROWSER_PROBES=on`) was left off, as the brief's verification list specifies. The change alters no geometry, only a row set, four strings and comments; the rendered summary text, the rendered row count and the per-row `data-activity-request` values are pinned by the Node lane test against the shipped assembled client, and the payload was checked on the real queue. The specific check not run is the strict browser lane over the generated site.

## Lessons read

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — read in full. Dropped from `required_lessons` at capture for budget; read here because the change touches `activity.go` and `generate.go`, which the shipped prime's **Read first** names. Nothing in it addresses the Activity view directly; the 0.275.3 entry (a hold the loop resolves is a wait, not an inbox item) was the closest and confirmed that `status_changed_at` on a `pending-heavy-testing` REQ reads as "held for heavy testing", which the rewritten fixtures rely on.
- `_dev/primes/lessons-kanban-board.md` — read in full, same reason. Four entries shaped the tests: REQ-318 (a reversed assertion needs its fixture reversed too — hence the `blocked_at`/`claimed_at` pair whose alphabetical order contradicts declaration order), REQ-322 (a constant a decision turns on must be read by the test, never restated beside it), REQ-293 (choose the mutation before looking at the pattern), REQ-313 (count unique filtered data rows when a summary describes what a view draws — which is exactly what the distinct-REQ count does).
- `_dev/primes/prime-kanban-board.md` and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — read in full (the REQ's `prime_files` plus the shipped routing index it points at).
- Crew rules read: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.

## Decisions

- **D-01 (DECIDE & STATE) — the third sort key is `StampField`, ascending.** Two stamps of one ticket at one instant cannot be separated by time or by REQ id, so without a third key `sort.Slice` decided. The direction is ascending because a field name carries no notion of newer or older; only the determinism matters, unlike the REQ id above it, where the larger number is the newer REQ. The comparator says so in a comment so the asymmetry does not read as an oversight.
- **D-02 (DECIDE & STATE) — the "(N before filters)" clause counts transitions.** It matches the leading count's unit. Reporting a REQ count there would put three units in one sentence and force the reader to work out which number the parenthetical belongs to. A comment states the choice at the call site, as the brief asked.
- **D-03 (DECIDE & STATE) — both empty states count transitions, with real singular forms.** "No lifecycle transition falls inside the last 24 hours." and "N transitions happened in this window, but the active filters hide all of them." (singular: "1 transition … hide it."). The old string used the `REQ(s)` parenthetical plural; a proper singular costs one ternary and the Node test pins both arms.
- **D-04 (DECIDE & STATE) — the slice capacity hint stays `len(tickets)`.** It is now a lower bound rather than an exact count. Multiplying it by a guessed stamps-per-ticket figure would invent a constant nobody can maintain, and `append` grows amortized. Measured cost at the real queue size: 1980 rows built once per generation.
- **D-05 (DECIDE & STATE) — `data-activity-request` non-uniqueness is documented and pinned.** The attribute now repeats across a REQ's rows. That is the shape REQ-573 needs for sibling highlighting, so a comment in `board-activity.js` says it is deliberate and the Node lane test asserts the repeated value. Without the pin, a future reader would "fix" it into a unique id and break the dependent REQ silently.
- **D-06 (DECIDE & STATE) — near-duplicate rows ship as-is.** A terminal REQ carries `completed_at` and `status_changed_at` at one instant, giving "completed" and "status changed to completed" as two rows. The REQ asks for no suppression rule and the second phrasing is already visible on the board today. Recorded in the first commit message as well.
- **D-07 (DECIDE & STATE) — the rewritten ordering test asserts triples, not ids.** `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` now pins `(RequestId, StampField, Transition)` per position. The old test kept a map keyed by REQ id for its transition checks, which cannot represent a REQ with several rows; folding them into the sequence keeps both properties in one assertion instead of dropping one.

## Incident during verification

While measuring payload growth I ran `git stash` in the worktree on a clean tree, which stashed nothing, and the paired `git stash pop` then popped a **pre-existing, unrelated stash** that belongs to another session: `stash@{0}: On main: do-work recovery: preserve interrupted REQ-539 orchestration metadata`. The pop conflicted and left one file in a `UU` state — `do-work/archive/UR-104/REQ-539-cut-the-contract-file-and-split-the-aggregate.md`, a `do-work/` path I am not allowed to write.

Resolved immediately with `git checkout HEAD -- <that file>`. Current state, verified: the worktree is clean (`git status --porcelain` empty), the stash entry is still present and unchanged (the failed pop kept it), HEAD is still `2d8beb40`, and nothing from that stash reached any commit of mine. No `do-work/` file was modified in the end. Flagging it because stashes are shared across worktrees in one repository, so the maintainer should know `stash@{0}` was touched and is intact.

The growth number in `## Verification` was then measured by parsing the generated payload instead, with no git operations at all.

## Discovered Tasks

- An unrelated stash from an earlier session is still parked in this repository: `stash@{0}: On main: do-work recovery: preserve interrupted REQ-539 orchestration metadata`. It is intact and untouched, but it is easy to pop by accident from any worktree, and it holds queue state. Someone who knows what REQ-539 was doing should apply it or drop it. → impact-noncritical → report only
- Before this REQ the Activity view had no JavaScript behavior coverage at all (`grep renderActivity|activity-summary javascript_behavior_*_test.go` was empty). The new test covers the summary, both empty states and the row set; the window buttons (`applyActivityWindowSelection`, `activityWindowPhrase` at 6/48/168 hours) are still unpinned in that lane. → impact-noncritical → report only
