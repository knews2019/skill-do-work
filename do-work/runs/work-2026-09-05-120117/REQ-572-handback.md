# Hand-back — REQ-572 (show every lifecycle transition of a REQ as its own Activity row)

**Outcome: empty diff, deliberately.** Every acceptance criterion already holds in current
main, every listed test passes, and I proved the tests are real by reverting the behavior
in a scratch experiment and watching six of them fail. Nothing was committed because there
was nothing left to implement. Evidence for each claim is below.

## Branch

- Branch: `worktree-agent-REQ-572-activity-rows`
- Head commit: `09a13839` — unchanged from the base commit the brief names. The branch
  carries **zero commits of its own**; `git diff --stat main...HEAD` is empty.
- Worktree `git status --porcelain` is empty (verified after the RED experiment was undone
  with `git checkout --`).

**The orchestrator has nothing to merge for this REQ.** Do not expect a builder range.

## File manifest

No file created, modified or deleted. All six Scope files were read in full and checked
against the acceptance criteria; each already carries the REQ-572 change from the earlier
merge at `fbdcd35e`:

| File | Verb | State found |
|---|---|---|
| `skills/do-work-board/tools/queue-kanban/activity.go` | unchanged | one row per parseable stamp (loop at lines 71-82), `StampField` as third sort key (line 101), `newestLifecycleStamp` absent |
| `skills/do-work-board/tools/queue-kanban/activity_test.go` | unchanged | 7 activity tests including the captured RED/GREEN case at line 19 |
| `skills/do-work-board/tools/queue-kanban/generate.go` | unchanged | comments restated at lines 82-86 and 312-317; `generatedActivityEntry` payload shape unchanged |
| `skills/do-work-board/tools/queue-kanban/web/board-activity.js` | unchanged | distinct-id set at lines 64-71, both counts in the summary at 73-80, both empty states counting transitions at 128-132 |
| `skills/do-work-board/tools/queue-kanban/web/template.html` | unchanged | Activity section comment restated at lines 460-469 |
| `skills/do-work-board/tools/queue-kanban/javascript_behavior_c_test.go` | unchanged | `TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` at line 2395 |

Tests touched: none. The seven Go tests and the one Node-lane test listed above are the
existing coverage; all pass and all were proven to fail against the pre-change behavior.

## Acceptance criteria, checked one at a time against current main

**AC-1 — a ticket with `created_at` and `claimed_at` yields two rows, "claimed" then
"captured", newest first.** HOLDS. `activity_test.go:19`
`TestBuildActivityRowsEmitsOneRowPerLifecycleStamp` is exactly the captured case (REQ-570,
created 22:52, claimed 23:00:17) and asserts both rows in that order. PASS.

**AC-2 — every parseable stamp becomes a row; a ticket with no parseable stamp is still
skipped.** HOLDS. `activity.go:71-82` loops `lifecycleTimestampFields(ticket)` and appends
one row per stamp; the `parsedOk` guard at line 73 skips unparseable values with no zero-time
fallback. Pinned by `TestBuildActivityRowsSkipsTicketsWithNoParseableStamp`
(`activity_test.go:179`) and by the tail of
`TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` (`activity_test.go:298-311`), which
compares the row set against the declared set rather than just the sizes — 14 declared stamps,
14 rows, each field exactly once.

**AC-3 — newest first with a deterministic tie-break that also covers two stamps of one REQ at
one instant.** HOLDS. The comparator at `activity.go:84-102` uses stamp time, then `RequestId`
descending, then `StampField` ascending, with a comment explaining that the field direction is
arbitrary and only the determinism matters.
`TestBuildActivityRowsBreaksStampTiesDeterministically` (`activity_test.go:196`) covers both
arms: two REQs at one instant fed in both input orders, and one REQ carrying
`blocked_at`/`claimed_at` at one instant plus `completed_at`/`status_changed_at` at another.

**AC-4 — the summary reports the transition count and the distinct REQ count; window, filter
chips and both empty states keep working.** HOLDS. `board-activity.js:64-80` builds the
distinct-id set and emits "N transitions across M REQs in the <window>" with real singular
forms plus the "(K before filters)" clause. Windowing still runs off `Date.now()`
(`board-activity.js:21`), the shared `requestMatchesFilters` chip filter still applies
(line 61), and both empty states count transitions (lines 128-132). The Node-lane test drives
the *shipped assembled client* over five cases — three transitions of two REQs, one
transition, one REQ filtered out, everything filtered, nothing in window — and pins the exact
strings. PASS.

**Requirements beyond the checklist, also verified:**
- No second stamp list anywhere; `model.go` untouched and not in the write set. Confirmed
  `grep -rn "newestLifecycleStamp"` returns nothing across the module.
- No "latest only" toggle was added. `template.html` `#activity-window-group` still holds
  exactly the four window buttons (6h / 24h / 48h / 7d).
- Payload shape unchanged: `generatedActivityEntry` still ships `id` / `stampField` /
  `stampAt` / `transition`, so REQ-573 (clicking a row to open the drawer and highlight
  sibling rows) has the contract it was promised. `data-activity-request` is still set on
  every `<tr>` and is still deliberately non-unique, with a comment saying so
  (`board-activity.js:88-90`) and a Node-lane assertion pinning the repeated id
  (`javascript_behavior_c_test.go:2497-2503`).

## Why the earlier work survived REQ-571 and REQ-576

The brief rejected the saved range because three later commits touched the same files. I read
those diffs rather than assuming. They do not disturb this REQ:

- `a55f24ce` and `4c76c332` (REQ-571, removing the `pending-heavy-testing` status) touch
  `activity.go`, `activity_test.go` and `javascript_behavior_c_test.go`, but only to sweep the
  retired status out of comments and test fixtures. `activity.go` changed 10 lines, all inside
  doc comments. No aggregation logic, no comparator, no assertion structure.
- `ae184a7b` (REQ-576, opening the card wall time at the earliest lifecycle stamp) touches
  `generate.go` only inside the `generatedRequest` struct's implementation-span comment — a
  different struct from `generatedActivityEntry`, and comments only.

## P-A-U

**[PLAN]** — Written before touching anything. The brief said the change may already be
present and that an empty diff is legitimate only with evidence, so the plan was verification
first, implementation only where a criterion failed:
1. Read all five crew rules and the kanban prime → verify: rules loaded before any file edit.
2. Read the current state of all six Scope files → verify: locate each acceptance criterion in
   real code, by line number.
3. Run both focused test lanes → verify: exit 0, and the Activity tests actually *ran* rather
   than being skipped for a missing `node` binary.
4. Manufacture RED by reverting the behavior in the working tree → verify: the tests fail on
   their assertions, not on a compile error, then `git checkout --` restores a clean tree.
5. Generate a real board against this repository's own queue → verify: REQ-570 carries its
   whole path instead of one row.
6. Implement only what steps 2-5 showed missing → nothing was missing, so nothing was written.

**[APPLY]** — No code applied, because every criterion was already satisfied. Two files were
temporarily modified in the working tree for the RED experiment (step 4) and restored with
`git checkout --` in the same session; `git status --porcelain` is empty. Nothing outside the
declared write set was written at any point. Nothing under `do-work/` was written except this
hand-back file, which is neither staged nor committed.

**[UNIFY]** — Run against the tree as it stands, since I am certifying it:
- `git diff --stat main...HEAD` → empty (no commits on the branch).
- `git status --porcelain` → empty (RED experiment fully undone).
- `gofmt -l .` in `skills/do-work-board/tools/queue-kanban` → no output, nothing unformatted.
- `go vet ./...` → exit 0, no findings.
- Debug-artifact scan `grep -n "console\.log\|debugger\|TODO\|FIXME\|XXX"` over all six Scope
  files → no matches.
- Files checked line by line and what was checked in each: `activity.go` (the per-stamp loop,
  the `parsedOk` skip, the three-key comparator, the doc comments); `activity_test.go` (all 7
  tests, their fixtures and their assertions); `generate.go` (the two Activity comment blocks
  and the unchanged `generatedActivityEntry` field tags); `web/board-activity.js` (the
  wall-clock window, the chip filter call, the distinct-id count, the summary string, the
  non-unique row attribute, both empty states); `web/template.html` (the Activity section
  comment, the window button group, the six column header ids); and
  `javascript_behavior_c_test.go` (the probe stub, the five render cases and every pinned
  string).

## Test evidence

All commands run from the worktree. Exit statuses captured explicitly.

**GREEN — Go lane (the brief's command).**
`bash _dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./...`
→ **exit 0**. `go-test budget: module=skills/do-work-board/tools/queue-kanban wall=45s
tests=392 slowest-file=generate_test.go:9.06s limit=<30s`. Every test file inside the 30s
per-file budget.

**GREEN — Node behavior lane (the brief's command).**
`QUEUE_KANBAN_JAVASCRIPT_PROBES=on QUEUE_KANBAN_STRICT_JAVASCRIPT_BEHAVIOR=1 bash
_dev/tests/run-go-tests-with-budget.sh skills/do-work-board/tools/queue-kanban ./... -run
'^TestJavaScriptBehavior'` → **exit 0**. `wall=7s tests=56
slowest-file=javascript_behavior_c_test.go:2.12s limit=<30s`.

**GREEN — the eight Activity tests named individually, verbose, to prove none was skipped.**
`go test -v -run '^(TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests|TestBuildActivityRows|TestLifecycleTimestampFieldsIsTheOneListBothReadersUse)' ./...`
→ **exit 0**, `ok ... 2.116s`, all eight `--- PASS`, no `--- SKIP`:
`TestBuildActivityRowsEmitsOneRowPerLifecycleStamp`,
`TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition`,
`TestBuildActivityRowsOrdersOneTicketsStampsNewestFirst`,
`TestBuildActivityRowsStraddlesTheWindowBoundary`,
`TestBuildActivityRowsSkipsTicketsWithNoParseableStamp`,
`TestBuildActivityRowsBreaksStampTiesDeterministically`,
`TestLifecycleTimestampFieldsIsTheOneListBothReadersUse`, and
`TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests (1.80s)` — the 1.80s is the
real `node` probe running, not a skip.

### RED observation (`tdd: true`, and the code arrived already implemented)

A test that passes on arrival proves nothing about whether it can fail. I reproduced RED by
restoring the pre-REQ-572 behavior in the working tree and running the same tests, then undid
it. Two edits, both reverted with `git checkout --`:
1. `activity.go` — the per-stamp loop replaced by a newest-stamp-only pick, which is what
   `newestLifecycleStamp` did.
2. `web/board-activity.js` — the summary and both empty states restored to the REQ-count
   wording the screenshot asset records ("38 REQs touched in the last 24 hours").

Result: **exit 1, six failures, every one on its assertion rather than on a compile or import
error.** Verbatim:

- `TestBuildActivityRowsEmitsOneRowPerLifecycleStamp` — `activity_test.go:32: rows = 1, want 2
  (one per parseable stamp): [{RequestId:REQ-570 StampField:claimed_at ... Transition:claimed}]`
  — the captured RED exactly: the claim shows, the capture eight minutes earlier does not.
- `TestBuildActivityRowsOrdersByNewestStampAndNamesTheTransition` — `activity_test.go:97:
  in-window rows = 5, want 8`.
- `TestBuildActivityRowsOrdersOneTicketsStampsNewestFirst` — `activity_test.go:139: rows = 1,
  want 4 (one per parseable stamp)`.
- `TestBuildActivityRowsBreaksStampTiesDeterministically` — `activity_test.go:231: rows = 1,
  want 4`.
- `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` — `activity_test.go:300: rows = 1,
  want 14 — one per declared lifecycle stamp`.
- `TestJavaScriptBehaviorActivitySummaryCountsTransitionsAndRequests` —
  `javascript_behavior_c_test.go:2492: summary for three transitions of two REQs = "3 REQs
  touched in the last 24 hours", want both counts`.

The two window/skip tests correctly stayed green under the revert, because neither depends on
the per-stamp fan-out. That is the right shape: the tests that pin this REQ fail without it and
the tests that pin neighbouring behavior do not.

### End-to-end evidence (acceptance is not inferred from unit tests)

`go run . generate --out <scratch> --repo-root <worktree>` → **exit 0**, "wrote static board
... 561 REQs, 118 URs". Reading the generated `board-data.js` with node:

- **2047 activity rows across 561 distinct REQs** — one row per REQ would be 561.
- **549 of 561 REQs carry more than one row.**
- **REQ-570 (the REQ from the captured case, which removed the pending-heavy-testing status)
  carries nine rows, its whole path:** `completed_at` → `release_at` → `review_at` →
  `integration_at` → `builder_handback_at` → `dispatch_at` → `planning_at` → `claimed_at` →
  `created_at`. Before this change it had one. Its two 00:05:50 stamps come back as
  `completed_at` then `release_at`, which is the third sort key working on real data.
- **Whole array verified newest-first** by comparing every adjacent pair — true.
- **Zero rows carry an unparseable `stampAt`**, so nothing is dated from the zero time.

## Lesson evidence

This REQ carries no `required_lessons` frontmatter key. Its `## Required Lessons — Dropped for
Budget` section names two satellites that were excluded for size, so neither was mandatory:

- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — **present** (22973 bytes).
- `_dev/primes/lessons-kanban-board.md` — **present** (19279 bytes).

Neither was missing. I checked both for relevance to this change rather than reading ~10k
tokens of unrelated material: `grep -c -i "activity"` returns **0 for both files**. Neither
satellite says anything about the Activity view, which matches what the REQ's own Exploration
recorded. No lesson contradicts anything here, because no lesson touches this surface.

Primes read in full: `_dev/primes/prime-kanban-board.md` (the maintainer prime named in
`prime_files`) and the `## Read first` / `## Do not edit` / `## Traps` sections of the shipped
`skills/do-work-board/tools/queue-kanban/prime-do-kanban.md`. Two conventions from them were
respected: the tool's three write surfaces are untouched (nothing here writes pipeline state or
`CHANGELOG.md`), and no build output was committed — the generated board went to the session
scratchpad, never to the repository.

Crew rules read in full before any file was opened for change: `general.md`,
`coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.

## Decisions

**D-01 — Ship an empty diff instead of rewriting working code. DECIDE & STATE.** All four
acceptance criteria hold in current main and all eight tests pass. Rewriting satisfied code to
produce a non-empty diff would be churn with a regression risk and no benefit. The brief
authorizes this outcome on condition of evidence, which is above. Reversible: if the
orchestrator disagrees, nothing has to be undone.

**D-02 — Manufacture RED by reverting the behavior in the working tree. DECIDE & STATE.** The
REQ is `tdd: true` and `testing.md` requires the failing test observed before the
implementation. That order was impossible here, since the implementation arrived first. The
nearest honest equivalent is to prove the tests can fail for the right reason: revert the
behavior, watch the assertions fail, restore. `testing.md` permits adapting the captured proof
when the codebase requires a nearby equivalent, and this is documented as that adaptation
rather than passed off as true RED-before-GREEN. Fully reversible and was reversed —
`git status` is empty.

**D-03 — Read the three later commits rather than trusting that the tests passing means they
were harmless. DECIDE & STATE.** The saved-range resume was rejected precisely because those
commits touched these files. A green suite is not evidence about what a diff did to a comment
or a fixture. Reading them cost one command and confirmed all three are comment and vocabulary
sweeps. Recorded above under its own heading so the orchestrator does not have to repeat it.

**D-04 — Check the two lesson satellites by targeted grep instead of reading ~10k tokens.
DECIDE & STATE.** Neither is a `required_lessons` entry; the REQ dropped both for budget. The
touch-conditional rule in `general.md` keys on whether the change touches what the prime's
`## Read first` or `## Traps` name, and there is no change. Zero "activity" matches in either
file settles relevance cheaply. Reversible: the files are named above if anyone wants them read.

**D-05 — Verify the Node probe actually executed rather than trusting the lane's exit 0.
DECIDE & STATE.** The Node behavior probes are gated on `QUEUE_KANBAN_JAVASCRIPT_PROBES=on` and
on `node` being present, so a skipped probe still exits 0 and would have let me certify AC-4 on
nothing. The verbose run shows `--- PASS ... (1.80s)`, which is a real execution.

## Discovered Tasks

- The REQ's frontmatter still carries `gate_deferred: 'true'` with `depends_on: [REQ-574]`.
  REQ-574 (building the do-work-cli test git fixtures once and reusing them, so the over-budget
  test files come back under 30 seconds) appears completed on main at `39a35194`, so the
  condition that deferred the gate may now be resolved. I did not run the gate — the brief
  forbids it and it is the orchestrator's. → report only, `impact-process`.
- `TestLifecycleTimestampFieldsIsTheOneListBothReadersUse` (`activity_test.go:258-274`) keeps a
  hand-maintained `RequestTicket` literal naming all 14 stamps, deliberately, as a tripwire that
  fires when a stamp is added to `lifecycleTimestampFields` without a transition phrase being
  chosen. It works as designed and the comment explains why it is a second spelling. Noting it
  only because a future reader may mistake it for the drifting closed enumeration the project
  warns about, when it is the opposite — an intentional tripwire. → report only,
  `impact-cosmetic`.
- The near-duplicate rows the REQ's Exploration flagged are real on live data: a terminal REQ
  emits both `completed_at` "completed" and `status_changed_at` "status changed to completed"
  at one instant, and REQ-570 shows the same pattern for `completed_at`/`release_at`. The REQ
  chose to ship this with no suppression rule and the maintainer accepted it, so this is not a
  defect — recorded only so the count 2047 is not read as a bug later. → report only,
  `impact-cosmetic`.

## Integration seams

None. No line belongs in any file outside the write set, and no file outside the write set
needs to change for this REQ to be complete.

The only thing the orchestrator must handle is the shape of the result: **there is no builder
branch content to merge.** `worktree-agent-REQ-572-activity-rows` sits at `09a13839`, the same
commit as its base and as `main`. Any integration step expecting a non-empty range will find
nothing, and that is the correct state — the change is already in main from `fbdcd35e`, which
is still an ancestor. The REQ can be marked implemented against current main on the evidence
above.
