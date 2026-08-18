---
id: REQ-265
title: Raise the two under-bounding mark-label constants to the current build
status: claimed
claimed_at: 2026-08-18T22:59:48Z
route: A
created_at: 2026-08-18T20:07:08Z
status_changed_at: 2026-08-18T21:01:24Z
user_request: UR-051
addendum_to: REQ-252
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations_test.go
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-18T22:59:48Z
  basis:
    - Route A
    - 1-file write set
    - 3 acceptance criteria
    - full-suite verification
---

# Raise the Two Under-Bounding Mark-Label Constants to the Current Build

## What

Chromium 141.0.7390.37 measures the 11px mark-label line box at **12.9631** (constant `durationsMeasuredLabelBoxHeightUnits` records 12.84) and its descent at **2.7778** (constant records 2.41). Per the larger-wins convention both constants should rise (≥12.97 / ≥2.78). Nothing is live-wrong today — the pitch-floor consumer clears the real box at pitch 13, and the ceiling consumer's paired title-ascent constant over-bounds by 0.99 — but the compensation is a coincidence of one consumer, not a guarantee. When raising, re-verify `TestDurationsLastLabelRowClearsPanelBTitle`'s margin (0.12 model units at the larger descent) and expect `TestDurationsLabelRowPitchClearsTheLabelTextBox` to still pass at pitch 13.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md` including its measured-face lesson and the render-evidence rule, plus `general.md`, `coding-guardrails.md` and `communication-style.md` (`tdd: false`, `maintenance: false`, so `testing.md` and `maintenance.md` did not load). Read the constants and both consumer tests, then **grepped the *quantity* rather than the two named constants — which is how the duplicate face bound in `generate_test.go` surfaced.** Planned: confirm by class identity in the renderer and CSS that the two descent constants describe one face; re-measure in Chromium; delete the duplicate rather than raise it; raise the box height with the reasoning beside it; re-verify both consumers' margins and non-vacuity; report the stale renderer comment as a seam instead of editing it.
- [x] **[APPLY]:** One file — `skills/do-work-board/tools/queue-kanban/durations_test.go`, the whole of the REQ's write set. Edits applied with an exact-match script that aborts unless each old block occurs exactly once, so no near-miss could be silently rewritten.
- [x] **[UNIFY]:** `git diff --stat` → one file, +51/−23. Full diff read line by line: only the two constants, one test's expression, and four comment blocks changed; **no assertion was weakened, removed, or re-pointed** except the one substitution named in the summary. No debug print, no `t.Skip`, no commented-out code, no TODO. `gofmt -l .` no output; `go vet ./...` clean; both re-run after the wording amendment. `git status --short --ignored` empty — **no build output in the source tree; every build went to scratch with `-o`, and no bare `go build` was ever run in the queue-kanban directory.**

## Context

REQ-252's builder measured both exceedances and captured the raise as a Discovered Task per the REQ's no-value-change rule; its review (F1a, gate: trivial) verified no assertion flips on any recorded build and routed the raise here as a durable artifact. Created `pending-answers` per the generation-≥2 depth stop.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-252: two measured constants no longer bound the face on a current build. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — wait until an assertion actually flips.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

---

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified) — integration seam: applied by the orchestrator inside the merge commit

**What was done:** Only one of the two constants needed a number. `durations_test.go`'s `durationsMeasuredLabelBoxDescentUnits = 2.41` turned out to be a **duplicate** of `generate_test.go`'s `durationsMeasuredMarkLabelDescentUnits = 2.8` — the same quantity of the same face, because `board-durations.js` puts `class: "durations-mark-label"` on both the annotation and the band labels and `board.css:1932` declares that class once at 11px. The builder deleted it and pointed the clearance test at the surviving bound, which structurally closes the REQ-241/REQ-242 collision shape that REQ-252 had closed only by convention. `durationsMeasuredLabelBoxHeightUnits` was raised 12.84 → **12.97** — the sample max, deliberately not padded — with the rationale, the caps and an explicit falsifier written into the constant's doc comment. The integration seam updated the stale measured numbers ("12.83 … 2.41 below") in the renderer's row-pitch comment, which no test reads.

---

## Discovered Tasks

Transcribed by the orchestrator from `do-work/runs/work-2026-08-18-230100/REQ-265-handback.md` (a worktree builder cannot write this file — REQ-270).

- **[normal] The shipped row pitch of 13 has ~0.03 units of slack against the largest sampled face, and nothing bounds an unsampled one.** The builder escalated this as D-05. The package's own part bounds already sum to 13.3 — over the pitch — and `--font-sans` ends in the open `sans-serif` generic, so no measurement taken in a Linux container bounds what a Mac or Windows machine actually draws. Not fixable inside this REQ's write set: raising the pitch immediately eats the Panel B ceiling, which this same REQ just narrowed to 0.10 model units. Needs its own REQ with both constraints on the table at once.

