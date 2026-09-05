---
id: REQ-590
title: 'Cap the path list in a verify finding so one detail cannot be 40 KB'
status: claimed
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
---

# Cap the Path List in a Verify Finding

## What
Three probes in `skills/do-work-board/tools/queue-kanban/verify.go` build a finding detail by joining an unbounded list into one sentence. When the list is long the detail becomes tens of kilobytes, which fills the board's findings strip and pushes every column off the screen, prints as one unreadable line in the terminal report, and lands whole in the shareable static snapshot. Give the three sites one shared cap: name the first few entries, then say how many were not named.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
