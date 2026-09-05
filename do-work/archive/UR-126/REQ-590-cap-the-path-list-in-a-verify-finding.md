---
id: REQ-590
title: 'Cap the path list in a verify finding so one detail cannot be 40 KB'
status: completed
created_at: 2026-09-05T18:29:25Z
user_request: UR-126
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: [skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go]
claimed_at: 2026-09-05T18:31:50Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-05T18:32:44Z
  basis:
    - trivial short-circuit
route: A
dispatch_at: 2026-09-05T18:33:59Z
builder_handback_at: 2026-09-05T18:37:50Z
review_at: 2026-09-05T18:44:24Z
commit: 079b12f5edfce50bbba70fd270b77ca899f42433
heavy_verified_at: 2026-09-05T19:01:05Z
heavy_verified_revision: 32813e4c1b10cb9436c7f0b9e08e16611c3edfb4
completed_at: 2026-09-05T19:02:44Z
release_at: 2026-09-05T19:02:44Z
---

# Cap the Path List in a Verify Finding

## What
Three probes in `skills/do-work-board/tools/queue-kanban/verify.go` build a finding detail by joining an unbounded list into one sentence. When the list is long the detail becomes tens of kilobytes, which fills the board's findings strip and pushes every column off the screen, prints as one unreadable line in the terminal report, and lands whole in the shareable static snapshot. Give the three sites one shared cap: name the first few entries, then say how many were not named.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One shared helper in `verify.go`, called at the three joins; RED test first, at the size the failure was observed.
- [x] **[APPLY]:** `verify.go` and `verify_test.go` only — the two files the write set names.
- [x] **[UNIFY]:** `git diff --stat` = 2 files, +176/-3. `go vet ./...` clean, `gofmt -l .` empty. Both files read for debug artifacts: no prints, no commented-out code, no TODOs.

## Why
Observed on the live board on 2026-09-05: a builder worktree had an untracked `do-work/` directory, `git status --porcelain --untracked-files=all` expanded it to roughly 700 paths, and the `worktree-wrote-queue-state` finding printed all of them. The strip became a wall of text several thousand pixels tall and the board was unusable until the worktree went clean on its own. The board is not at fault: it renders the producer's text verbatim on purpose, so the fix belongs in the producer, where it also repairs the terminal report and the snapshot at the same time.

## Detailed Requirements
- The three sites, all in `skills/do-work-board/tools/queue-kanban/verify.go`:
  - line 1191, `worktree-wrote-queue-state` — `strings.Join(dirtyQueuePaths, ", ")`, the one observed.
  - line 1177, `worktree-committed-queue-state` — `strings.Join(committedQueuePaths, ", ")`, same shape, triggered by a branch that changed many `do-work/` files.
  - line 797, `ur-archived-with-live-member` — `strings.Join(liveMemberIds, ", ")`, a UR with many live members.
- One shared helper, not three local caps, so the three details cannot drift apart.
- Cap the named entries at 5 and follow them with the exact remaining count, for example `and 712 more`. The count is the true total minus the named entries, never an approximation.
- A list at or under the cap keeps its current text exactly: no trailing count, no ellipsis, no change to any existing test's expected string.
- Both surfaces get the capped text, because both read the same `Detail` field: the board payload (`generate.go` → `attachVerifyFindings`) and the terminal report (`verify.go` → `renderVerifyReport`).
- No change to the board client. `board-cards.js` and `board.css` stay as they are; the strip's rendering contract is that the producer's text is printed as given.

## Constraints
- Scope is exactly the three joins and their shared helper plus the lock-in test. Do not add a cap to any other finding, do not touch the strip's CSS or JavaScript, and do not change what any probe detects.
- `REQ-589` (rendering the verify findings strip as the M4 slim band) is in flight against the same strip but on the client side only, and its own record says the producer and the payload stay untouched. This REQ must not edit the client files that REQ-589 is changing.
- Keep the category, subject, remedy and `fixable` values unchanged; only the detail's path list is capped.

## Red-Green Proof
**RED prompt/case:** A Go test in `skills/do-work-board/tools/queue-kanban/verify_test.go` that hands the detail builder 200 dirty `do-work/` paths and asserts the resulting detail names at most 5 of them and ends with `and 195 more`.
**Why RED now:** `strings.Join` at `verify.go:1191` writes every element, so the detail contains all 200 paths and the assertion fails on both halves — too many named, and no remaining count.
**GREEN when:** the test passes at all three sites (the same assertion applied to the committed-paths and live-members details), and `go test ./...` in `skills/do-work-board/tools/queue-kanban` stays green.
**Validation:** User confirmed — the cap shape ("first 5 paths plus and N more") was proposed in the session and the user answered "capture it as a req".

## Context
Investigated in the session of 2026-09-05. Live evidence at the time: the board screenshot in `do-work/user-requests/UR-126/input.md`, the finding text quoted there, and a `go run . verify` run in `skills/do-work-board/tools/queue-kanban` that returned four ordinary worktree findings once the worktree went clean. The board client was read and cleared: `board-cards.js:763` sets the detail with `textContent`, and `.board-findings-row` in `board.css:658` wraps correctly with no clamp, so nothing overflows sideways and the layout itself is sound. Also checked and found clean: no other unbounded producer string reaches the board (`model.go:1499` joins a fixed field list, `completionAnomalyReason` holds field names), and findings do not reach the Markdown export or the clipboard copy.

## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 6477 tokens, over the 2000-token budget, and its row says `slugged: partial`, so the targeted `#subject-not-restated-in-detail` form is not eligible. Reason it matched: its `subject-not-restated-in-detail` family is about how a verify finding's detail is composed, which is exactly this change.
- `_dev/primes/lessons-kanban-board.md` — 4959 tokens, over the same budget, `slugged: partial` so no targeted form. Reason it matched: it applies to changes in queue-kanban static output and board publication, which the capped detail reaches through the snapshot.

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the three exact call sites, the file they live in, the test file, and the shape of the replacement text. Nothing needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)

**What was done:** Added the constant maxNamedItemsInFindingDetail (5) and the helper joinCappedItemList to verify.go, and replaced the three unbounded joins that built a finding detail with calls to it — uncommitted queue state, committed queue state, and an archived UR's live members. A list at or under five entries is joined exactly as before; a longer one names five and then appends the exact number it did not name. Four tests in verify_test.go pin the cap at each of the three sites and pin the unchanged short-list text.

## Qualification

**Passed.** Read from the merge range `bc11f17d..079b12f5`.

- Two files changed, exactly the declared write set. The diff is +176/-3: one constant, one 6-line helper, three one-line call-site swaps, four tests, one import.
- The change is substantive, not cosmetic. Before it, the three details spelled out every entry; after it, a list longer than five names five and states the exact remainder. The observed 200-path case is in the test suite at the size it was observed.
- Requirements traced one by one: three sites capped (verify.go:1222, verify.go:1208, verify.go:842 after the change) — one shared helper, not three local caps — cap of 5 with the exact remaining count, short lists byte-identical to before, and no board client file touched. The existing tests that assert a single path inside a detail still pass unchanged, which is the short-list requirement proved from the other side.
- Both surfaces are covered without a second change because both read the same `Detail` field: `attachVerifyFindings` in generate.go builds the board payload from it and `renderVerifyReport` prints it.
- No debug artifacts in the diff. `go vet ./...` clean, `gofmt -l .` empty.

Requirements traced: every bullet of `## Detailed Requirements` is satisfied and none of `## Constraints` is violated.

*Checked by work action*

## Testing

**Tests run:** `cd skills/do-work-board/tools/queue-kanban && go test -count=1 -run 'TestVerify' .` (focused, 7.4s) and `go test ./...` in the same package (full, 50.4s). Canonical repository gate: `bash _dev/tests/maintainer-verify.sh` from the project root, launched through `run-timed-command`, 171s.
**Result:** ✓ All passing. Gate exited 0 on its first run, so no retry.

**Red-green validation** (traced to `## Red-Green Proof`):
- `TestVerifyCapsThePathListForUncommittedQueueState` (200 dirty paths, the size the failure was observed at): ✗ `detail names 200 of the 200 dirty paths, want 5` before → ✓ after.
- `TestVerifyCapsThePathListForCommittedQueueState` (8 committed paths): ✗ `detail names 8 of the 8 committed paths, want 5` before → ✓ after.
- `TestVerifyCapsTheMemberListForAnArchivedUserRequest` (8 live members): ✗ `detail names 8 of the 8 live members, want 5` before → ✓ after.
- `TestVerifyLeavesAShortListUncapped`: ✓ before and after. It is the no-regression guard for the short-list requirement, so passing while the others were red is the correct state.

All four failures were assertion failures, not compile or import errors. The review reproduced the RED state independently: a scratch copy of the package with `verify.go` restored from `bc11f17d` and the new test file kept produced the same three messages.

**New tests added:**
- `TestVerifyCapsThePathListForUncommittedQueueState`, `TestVerifyCapsThePathListForCommittedQueueState`, `TestVerifyCapsTheMemberListForAnArchivedUserRequest`, `TestVerifyLeavesAShortListUncapped`, plus the `cappedDetailNamedItems` expectation and the `countNamedListItems` helper.

**Existing tests updated (cross-REQ impact):** none. Every existing detail assertion is a short list and is byte-identical to before.

**Gate revision drift:** the gate ran from the project root while a sibling session was working the same checkout. It observed the tree at `4a909573` and finished one second after that session released 0.303.6, so `HEAD` is now `32813e4c`. `git diff --stat 4a909573 HEAD` is eight files: the release itself (VERSION, CHANGELOG and their mirrors) plus lesson prose. No code, so the gate's meaning for this REQ is unchanged. Recorded green through `advance`, which accepted it.

**Heavy verification plan:** *(planned at `bc11f17d..079b12f5`, held for the drain)*
- Range: bc11f17d..079b12f5
- `queue-kanban-javascript`: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — both changed files matched subtree `skills/do-work-board/tools/queue-kanban`.
- `queue-kanban-browser`: same argv with `--heavy-lane queue-kanban-browser` — same reason.
- `staged-skills`: same argv with `--heavy-lane staged-skills` — both changed files matched subtree `skills`.

*Verified by work action*
## Review

**Overall: 97%** | 2026-09-05T18:44:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** F1 — the merge range carries no `CHANGELOG.md` entry or skill version bump for a user-visible change to the board's finding text; the run's finalization release must cover it (`_dev/primes/prime-kanban-board.md` § Conventions) — impact-negligible → report only. F2 (Nit) — no test pins the cap+1 boundary (six entries → `and 1 more`), so the singular wording is unpinned (`verify_test.go:2649`) — impact-negligible → report only
**Acceptance:** Pass — focused suite green at `079b12f5`, the three RED failures reproduce on `verify.go` restored from `bc11f17d`, and a live `go run . verify` prints short details with no truncation artifact
**Suggested testing:** 1 item
**Follow-ups created:** None (2 findings report only)

Restatement Sweep: the diff redefines what three finding details contain. Swept `Detail` consumers (`generate.go:651`, `verify.go:1398`, `board-cards.js` `textContent`), category-keyed routing (`boardRenderedVerifyCategories`, `forensics.md` Check 14, `doctor_commands_test.go`), the `fixable` repair path (none of the three is fixable), and every prose hit for the three category names across `skills/`, `_dev/primes/`, `decisions/` and `kb/`. No stale restatement found; `actions/work-reference.md` line 427 describes the orchestrator's own full-path `git diff` guard, which is untouched.

*Reviewed by review-work action*

## Heavy Verification Result

Target revision: `079b12f5` (this REQ's merge). Execution revision: `32813e4c` — a detached worktree at the integration tip, which contains this merge plus a sibling session's REQ-589 (rendering the verify findings strip as the M4 slim band) and release 0.303.6. Running at the tip is what proves the merged state, and the detached worktree is what keeps the execution revision still while that session keeps committing.

| Lane | Exit | Wall | Disposition |
| --- | --- | --- | --- |
| queue-kanban-javascript | 0 | 11s | executed (fingerprint_mismatch) |
| queue-kanban-browser | 0 | 127s | executed (fingerprint_uncertain) |
| staged-skills | 0 | 52s | executed (fingerprint_mismatch) |

All three executed rather than reusing evidence, and none was skipped: the browser lane's 127 seconds is a real run, with `QUEUE_KANBAN_BROWSER` naming the installed Chrome so it could not report a silent skip. The drain waited 422 seconds for a sibling session's gate to finish first, so these timings were taken off a machine that was not running two suites at once.

## Lessons Learned

**What worked:** Fixing the producer instead of the two surfaces. The board strip and the terminal report both print `Detail` verbatim by contract, so one helper in `verify.go` repaired both plus the shareable snapshot, and no client file had to be touched while a sibling session was rewriting that same strip.

**What didn't:** The first draft of the RED test referenced the production constant `maxNamedItemsInFindingDetail`, which does not exist before the fix, so the test failed to compile — an import error, not an assertion failure, and not a valid RED. The tests now carry their own `cappedDetailNamedItems = 5`, which both compiles at RED and makes a later change to the production cap fail the lock-in instead of moving with it.

**Worth knowing:** A producer string that two surfaces print verbatim has to bound itself, because neither surface can bound it without breaking the "render what the producer wrote" contract that keeps them from drifting. The trigger here is `git status --porcelain --untracked-files=all`, which expands one untracked directory into one line per file — a single untracked `do-work/` became roughly 700 paths. Two other joins in the same file had the same shape and no incident yet; they were capped in the same pass.

## Orientation

`queue-kanban verify` findings now bound the lists they carry: five entries named, then the count of the rest. Lives in the queue-kanban verification subsystem (`_dev/primes/prime-kanban-board.md`), reached by the board's findings strip, the `verify` terminal report, and the static board snapshot. Not a map change — no new module, contract or concept, and the prime's referenced paths all still exist.

## Timing

Observed 2026-09-05T18:33:59Z to 2026-09-05T18:47:54Z: 13m 55s total, 10m 08s attributed across 2 events, 3m 47s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 7m 17s | 1 |
| verification-gate | 2m 51s | 1 |

Slowest stage: builder-work / route-a-builder, 7m 17s, outcome success.
Slowest command: verification-gate / canonical-repository-gate, 2m 51s, exit 0, bash (2 argv tokens).
