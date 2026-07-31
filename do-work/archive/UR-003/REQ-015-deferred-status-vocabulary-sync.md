---
id: REQ-015
title: Sync the `deferred` status between the queue-kanban parser and the Schema Read Contract
status: completed
created_at: 2026-07-01T21:06:45Z
claimed_at: 2026-07-01T21:17:04Z
completed_at: 2026-07-01T21:25:31Z
route: A
commit: 27f1005
user_request: UR-003
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-016]
kb_status: pending
---

# Sync the `deferred` status between the queue-kanban parser and the Schema Read Contract

## What

`tools/queue-kanban/model.go` treats `deferred` as a recognized Needs-input/Blocked status (`isNeedsInputOrBlockedStatus`, ~line 410, plus the column comment ~line 116), but the Schema Read Contract in `actions/work-reference.md` does not list `deferred` in its status enum, and nothing in the skill ever writes `status: deferred`. Make the two vocabularies agree — recommended direction: remove `deferred` from `model.go`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Confirmed sweep: `deferred` appears in `tools/queue-kanban/` only at model.go:116 (column comment), model.go:410 (`isNeedsInputOrBlockedStatus` case), and model_test.go:27,43,281,294. Confirmed Schema Read Contract enum (`actions/work-reference.md` line 130) has no `deferred`. RED: extend `model_test.go` so `isNeedsInputOrBlockedStatus("deferred")` is asserted false (TestStatusClassifiers) and `TestBucketColumns` asserts a `deferred` ticket lands in NeedsInputOrBlocked via the unrecognized-status warning path (2 warnings, not 1) — run `go test ./...`, capture the pre-fix failure. GREEN: remove the `"deferred"` case from `isNeedsInputOrBlockedStatus` in model.go and drop it from the column comment; rerun tests to green. Scope: model.go + model_test.go only.
- [x] **[APPLY]:** Wrote the failing test changes first in `model_test.go` (RED: `TestStatusClassifiers` asserts `isNeedsInputOrBlockedStatus("deferred")` is false; `TestBucketColumns` asserts a `deferred` ticket now generates an unrecognized-status warning, 2 warnings total instead of 1), ran `go test ./...` and captured the failures, then removed the `"deferred"` case from `isNeedsInputOrBlockedStatus` and the stale column comment in `model.go` (GREEN). Also swapped the unrelated `{"deferred","deferred"}` case in `TestNormalizeStatus` for a neutral `{"custom-status","custom-status"}` pair (D-01). Scope stayed exactly to `tools/queue-kanban/model.go` + `tools/queue-kanban/model_test.go`.
- [x] **[UNIFY]:** `git diff --stat -- tools/queue-kanban/` shows exactly 2 files changed (model.go: 5 lines; model_test.go: 33 lines), matching planned scope. Ran `gofmt -l .` inside `tools/queue-kanban/` — no output (clean). Ran `go vet ./...` inside `tools/queue-kanban/` — no output (clean). Reviewed the full `git diff` for both files line by line — no debug prints, no commented-out code, no stray artifacts. Ran `git status --porcelain tools/` — only the two expected modified files, gitignored binary absent. Ran `go test ./...` post-fix — full suite passes (`ok`, no FAIL).

## Why (if provided)

Root `CLAUDE.md` mandates the parser stay in lock-step with the Schema Read Contract ("any change to that contract must be mirrored in `tools/queue-kanban/model.go` (and vice-versa)"). A status the parser recognizes but the contract doesn't define is exactly that drift.

## Context

- Surfaced as an out-of-scope finding during the 2026-07-01 `do-work validate-feedback` triage (see UR-003).
- No producer exists: the canonical status enum is `pending`, `claimed`, `completed`, `completed-with-issues`, `failed`, `pending-answers`, `blocked-archive-collision`, `blocked-dependency-cycle` (Schema Read Contract, `actions/work-reference.md`). No action file writes `deferred`.
- Since 0.102.1 ("The Honest Board"), unrecognized statuses already surface in the Needs-input/Blocked column **with a warning** instead of vanishing — so removing `deferred` from the recognized set keeps a hand-edited `status: deferred` ticket visible (same column, now warned) rather than silently blessed. Behavioral delta is the warning only.
- Alternative direction (add `deferred` to the Schema Read Contract instead) would make it a status every consumer of the contract must honor (work scan, clarify, roadmap, cleanup) — do this only if a real producer/need is found; otherwise it fails YAGNI.

## Builder Guidance

Firm on the outcome (the parser's recognized-status vocabulary and the contract's enum must agree exactly); the removal direction was **user-confirmed** during the verify-requests pass (2026-07-01). Before removing, verify no producer exists: grep the repo (actions/, hooks/, specs/, do-work/) for `deferred`. If a genuine producer or documented need turns up, switch to the canonicalize direction and mirror the contract + parser in the same commit per `CLAUDE.md`.

## Red-Green Proof
**RED prompt/case:** A `model_test.go` assertion that `isNeedsInputOrBlockedStatus("deferred")` is false (i.e., `deferred` routes through the unrecognized-status warning path) fails today — the switch explicitly recognizes it.
**Why RED now:** `model.go` recognizes a status the Schema Read Contract never defined and nothing produces.
**GREEN when:** The recognized-status set in `model.go` exactly matches the contract's enum; a ticket with `status: deferred` lands in Needs-input/Blocked via the unrecognized-status warning path; `go test ./...` in `tools/queue-kanban/` passes.
**Validation:** User confirmed (verify-requests pass, 2026-07-01 — removal direction confirmed over canonicalizing)

---
*Source: UR-003 — "capture the two out-of-scope kanban findings as REQs" (finding 1, restated in the UR Summary)*

Think carefully before answering.

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the exact files and lines (`model.go` ~410 switch case, ~116 column comment, plus `model_test.go` expectations), the direction is user-confirmed (remove `deferred`), and the orchestrator's producer-grep confirmed no `status: deferred` producer exists outside `tools/queue-kanban/`. Well-specified removal with obvious scope.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 (DECIDE & STATE):** `TestNormalizeStatus`'s `{"deferred", "deferred"}` case (line ~27) tests `normalizeStatus`'s generic lowercase/trim passthrough, not status recognition — it would have kept passing unchanged either way, since `normalizeStatus` has no special-case aliasing for `deferred`. Swapped the input value to `{"custom-status", "custom-status"}` so the test still documents passthrough-for-arbitrary-strings coverage without leaving a stale `deferred` grep hit in the suite that could read as if `deferred` were still sanctioned somewhere. Reversible, leaf-level, no behavior change — DECIDE & STATE per karpathy.md's gate.

## Discovered Tasks

None — the sweep of `tools/queue-kanban/` for `deferred` and the confirmation against the Schema Read Contract enum (`actions/work-reference.md` line 130) turned up no additional producers, sites, or drift beyond the ones the REQ already named.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified) — removed the `"deferred"` case from `isNeedsInputOrBlockedStatus` and dropped `deferred` from the `NeedsInputOrBlocked` column comment
- `tools/queue-kanban/model_test.go` (modified) — `TestStatusClassifiers` now asserts `isNeedsInputOrBlockedStatus("deferred")` is false; `TestBucketColumns` asserts a `status: deferred` ticket still lands in Needs-input/Blocked but via the unrecognized-status warning path (2 warnings: deferred + typo case); `TestNormalizeStatus` passthrough case swapped to a neutral string (D-01)

**What was done:** Removed the producer-less `deferred` status from the queue-kanban parser's recognized-status set so it exactly matches the Schema Read Contract enum (`actions/work-reference.md`); a hand-edited `status: deferred` ticket now surfaces in Needs-input/Blocked with an unrecognized-status warning instead of being silently blessed.

## Qualification

Passed — 2 files verified on disk (`git diff` shows exactly the declared scope), all 3 REQ requirements traced (recognized-set parity with the contract enum; `deferred` routes through the warning path; `go test ./...` green), P-A-U boxes checked with substantive notes, no debug artifacts in diff, orchestrator independently re-ran tests and linters.

## Testing

**Tests run:** `go test ./...` in `tools/queue-kanban/` (orchestrator re-ran independently)
**Result:** ✓ All passing (`ok github.com/knews2019/skill-do-work/queue-kanban`)

**Red-green validation:**
- `TestStatusClassifiers` (deferred-not-recognized assertion): ✗ before implementation (`deferred is not in the Schema Read Contract enum ... must route through the unrecognized-status warning path`) → ✓ after
- `TestBucketColumns` (deferred ticket → warning path, 2 warnings): ✗ before implementation (`expected two unrecognized-status warnings (deferred + pnding), got 1`) → ✓ after

Traces directly to the REQ's `## Red-Green Proof`: RED case implemented as captured (assertion that `isNeedsInputOrBlockedStatus("deferred")` is false, failing pre-fix); GREEN criteria met verbatim.

**Existing tests updated (cross-REQ impact):**
- `TestBucketColumns` / `TestStatusClassifiers` / `TestNormalizeStatus` — repointed `deferred` expectations at the unrecognized-status warning path (coverage preserved, not deleted)

**Linters:** `gofmt -l .` clean; `go vet ./...` clean.

## Review

**Approve** — clean, surgical removal of `deferred` from the queue-kanban parser's recognized-status set; parser and Schema Read Contract now agree exactly, verified by full test suite + vet + gofmt.
Route A | reviewed pre-commit

### What's built
- `isNeedsInputOrBlockedStatus` in `tools/queue-kanban/model.go` no longer recognizes `deferred`; the stale column comment was updated to match.
- A hand-edited `status: deferred` ticket still surfaces in Needs-input/Blocked (same visible outcome), but now via the unrecognized-status warning path (`bucketColumns` default case) instead of being silently blessed as a first-class status.

### Decisions / risks for you
- D-01 (in REQ): swapped `TestNormalizeStatus`'s `{"deferred","deferred"}` case for a neutral `{"custom-status","custom-status"}` pair — correct call; that test exercises generic passthrough, not status recognition. No risk.
- None outstanding — repo-wide grep confirms `deferred` now appears only inside `model_test.go`'s intentional "prove it's unrecognized" assertions.

### Findings

**Important:** None. **Minor:** None. **Nit:** None — diff is minimal, no stray formatting, no debug artifacts.

### Requirements Checklist

- [x] Remove `deferred` from `isNeedsInputOrBlockedStatus` switch in `model.go` — delivered (model.go:404-414)
- [x] Drop `deferred` from the `NeedsInputOrBlocked` column comment — delivered (model.go:116)
- [x] Recognized-status set exactly matches the Schema Read Contract enum — delivered, confirmed by direct comparison
- [x] A `status: deferred` ticket still lands in Needs-input/Blocked via the unrecognized-status warning path — delivered; verified by reading `bucketColumns`'s default case (appends to NeedsInputOrBlocked, warns via `ticket.OriginalStatus`)
- [x] `go test ./...` passes in `tools/queue-kanban/` — verified independently (`-count=1` to bypass cache)
- [x] Red-green evidence traces to the REQ's Red-Green Proof — matches verbatim
- [x] Builder-guidance pre-check (grep for a `deferred` producer before removing) — independent grep confirms zero producers outside the test's negative assertions and this REQ/UR's own prose

### Acceptance Testing

**Result: Pass** — `go test -count=1 ./...`, `go vet ./...`, `gofmt -l .` all clean; warning-path mechanics read and confirmed (test's REQ-6 entry correctly sets `OriginalStatus: "deferred"`, mirroring the pre-existing `pnding` pattern, so the assertion exercises the real code path).

### Suggested Additional Testing

- None — Route A vocabulary-sync removal with unit coverage of both the removed branch and the fallback path it now exercises.

### Scores

**Overall: 100%** — Requirements 100 / Code Quality 100 / Test Adequacy 100 / Scope 100 / Risk: none / Acceptance: Pass.

### Follow-up REQs Created
None.

*Generated by review-work agent (pipeline mode)*

## Lessons Learned

**What worked:** Anchoring RED directly on the REQ's captured `## Red-Green Proof` (the failing `isNeedsInputOrBlockedStatus("deferred")` assertion) made the TDD cycle mechanical; repointing existing `deferred` test expectations at the warning path preserved coverage instead of deleting it.
**What didn't:** Nothing — no dead ends on this one.
**Worth knowing:** Synthetic tickets in `model_test.go` must set `OriginalStatus` (not just `Status`) for unrecognized-status warning assertions to exercise the real code path — the warning text is built from `ticket.OriginalStatus`, so omitting it makes a warning assertion pass trivially. Also, `TestNormalizeStatus` tests generic lowercase/trim passthrough, not status recognition — don't read its cases as a sanctioned-status list.

## Orientation

The queue-kanban board's recognized-status vocabulary now matches the Schema Read Contract exactly — a hand-edited `status: deferred` ticket surfaces in Needs-input/Blocked with an unrecognized-status warning instead of being silently blessed. Lives in the board's status-bucketing parser (`tools/queue-kanban/model.go`, per prime-do-kanban's Stakes). No map change; prime spot-check clean (all referenced paths exist).
