---
id: REQ-265
title: Raise the two under-bounding mark-label constants to the current build
status: completed
completed_at: 2026-08-18T23:54:48Z
commit: 1227678
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

---

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` against the merged tree (range `24a256a..1227678`), un-piped with the exit code read directly
**Result:** ✓ `GATE_EXIT=0` — both Go packages green including the strict JavaScript behaviour lane. This run is both Step 6.5's testing and Step 8's post-merge verification.

**Red-green validation:** none is owed — no behaviour changed. These are test-side constants and comments; the acceptance evidence is that both guards are **non-vacuous**, proven by perturbation and independently reproduced by the reviewer:

| Perturbation | Result |
|---|---|
| box → 13.01 | FAIL — "row pitch 13.00 is smaller than the 13.01-unit line box the browser draws" |
| box → 13.00 | PASS — the floor is `<`, so it flips **strictly above** 13.00 |
| descent → 2.91 | FAIL — "ends at 337.91 but Panel B's title starts at 337.90" |
| descent → 2.90 | **FAIL** — the ceiling is `>=`, so it flips **at** 2.90 |

**A bound nothing would notice changing is not a bound**, so this table is the real evidence. Note the last row: the builder's and this REQ's own text said the ceiling flips *above* 2.90; it flips *at* 2.90. Corrected above.

**Independent re-measurement** (reviewer's own run, Chromium 141.0.7390.37, board generated fresh to scratch, `location.href` returned from the same `evaluate` call): box 12.963112831, ascent 10.185302734, descent 2.777810097 — **identical across five probe strings** including all-descender `gjpqy`, no-descender `MMMM`, and a digit run, confirming it is the line box rather than the ink, and matching the hand-back to every reported digit. Live Panel B title ascent 11.111236, giving live headroom of **1.1110** units.

**What the 0.10 model margin is actually made of** — recorded because the number flatters itself. It pairs the mark-label descent bound (2.8, rounded up from 2.7778) with the axis-title ascent bound (12.1, rounded up from 12.0372): two maxima from different builds. Run the same model on the **raw** measurements and the margin is 0.185, so **0.085 of the 0.10 is round-up padding and only ~0.015 units is measured margin.** Historical build-to-build swing on those two numbers has been 0.81 and 0.37 — four to eight times the entire margin. Practical consequence: this assertion is now likelier to flip on a browser bump than to catch a geometry regression, while the picture is nowhere near it. That is a property of the model, not of the render.

**New tests added:** none, and the reviewer agreed with that call (D-03). `TestDurationsMeasuredConstantsNameTheirChromiumBuild` still passes and its vacuity guard is `count == 0`, so deleting a constant did not hollow it out.

**Existing tests updated (cross-REQ impact):** `TestDurationsLastLabelRowClearsPanelBTitle` now reads `durationsMeasuredMarkLabelDescentUnits` (2.8) instead of the deleted duplicate (2.41). The effective bound **rose**, so the raise-only convention holds at the consumer.

*Verified by work action*

---

## Review

**Overall: 92%** | 2026-08-18T23:54:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 85% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- **Against the orchestrator's integration seam, not the builder's work.** The seam claimed the declared parts "round both up, to an ascent of 11 and a descent of 2" — **false for the descent**: `durationsLabelTextDescentUnits = 2.0` is *below* the drawn 2.7778, and `generate_test.go:2109` takes a `math.Max` against the measured constant precisely because of that. The error was pre-existing at 0.41 and the seam widened it to 0.78. REQ-241's own RED #2 — "declared box (10+2=12) satisfied, measured face (12.84) not" — is exactly the mistake that sentence would talk a reader into. — gate: trivial → **fixed in place before shipping**, since a comment edited *for* correctness of record must not ship a false safety claim
- `durationsMeasuredMarkLabelDescentUnits` is now the package's single bound for the `.durations-mark-label` face and is read from two files, but its doc comment still calls it "The annotation box's descent" and its block header still scopes it to the annotation. The consumer side explains the cross-file dependency thoroughly; the canonical home says nothing — the restatement half of the consolidation this REQ performed, and the same shape that made REQ-241 and REQ-242 collide. — gate: trivial → REQ-277 created (sweep)

**Minor findings:** 3 (report only; two folded into REQ-277)
- The seam restated "now 12.97" — a value away from its home, in a file no test reads, which goes stale on the next raise. **Dropped in the same fix.**
- Boundary claims in the hand-back and this REQ ("flips above 2.90", "past 12.2") were off by one boundary; both assertions use `>=`. **Corrected above.**
- `durations_test.go:459-466`'s block header says every number in it is rounded up, but the block now holds a supremum-by-argument and an explicit non-supremum side by side. Folded into REQ-277.

**Acceptance:** Pass — **the deletion is better than the REQ asked for, verified rather than accepted.** The reviewer checked the duplication claim against the actual sources: `board-durations.js:187` and `:427` both set `class: "durations-mark-label"`, `board.css:1932-1936` declares that class exactly once at 11px with no more-specific selector, and its own five-string re-measurement shows the descent identical for all-descender, no-descender and digit strings — one face, one line box, therefore one descent. The decisive point the hand-back did not quite make: **the raise-only invariant holds at the consumer**, since the value the clearance test actually uses went *up* 2.41 → 2.8, past the REQ's own ≥2.78 bar. A duplicate record of a bound was removed while the bound itself rose. Not an unrequested removal of a guard.

On the padding disagreement the reviewer sided with the builder's conclusion and rejected his reasoning: his claim that "the alarm fires when a measurement exceeds the pitch, not when it exceeds the constant" is false — the assertion is `rowHeight < boxConstant`, so the constant *is* the trigger, demonstrated by setting it to 13.01 with nothing re-measured and watching the floor fail. The argument that does hold is about the record: pad the measured constant and the next person who measures 12.99 sees no edit to make, so the fact that the face grew never gets written down. **A third option neither took, with precedent in the same block:** the width side separates measurement from buffer (`durationsMeasuredLabelWidthSupremumUnits` plus `durationsLabelWidthModelSlackUnits`), and a box slack constant would belong there rather than inside the measurement.

**Suggested testing:** 3 items
**Follow-ups created:** REQ-277 (comment-accuracy sweep), REQ-278 (measure the face off Linux, scoped to measuring); **sweeps appended to:** None

*Reviewed by review-work action*

---

## Lessons Learned

**What worked:** Grepping the *quantity* rather than the constant names. That is the only reason this became a deletion instead of two raises — searching for what the number measures, rather than what it is called, surfaced a second bound for the same face sitting in a different file. It structurally closed the REQ-241/REQ-242 collision shape that REQ-252 could only close by convention.

**What didn't:** The orchestrator applied the builder's handed-back seam text without checking its claim, and shipped a sentence asserting the declared parts round *up* when the declared descent rounds *down*. The reviewer caught it. **A seam is the one part of a worktree REQ the builder cannot test, and the integrator is the only reader it gets** — applying it verbatim because it arrived in a hand-back is exactly the deference the sole-integrator rule exists to prevent. Also: the builder reached the right answer on padding via an argument that is simply false (the constant *is* the trigger, provably), and nobody would have noticed if the reviewer had checked only the conclusion.

**Worth knowing:** The Panel B clearance margin now reads 0.10 model units, and about 0.085 of that is round-up padding — only ~0.015 is measured. Build-to-build swing on its two inputs has historically been 0.81 and 0.37. So that assertion is now a tripwire on browser bumps rather than a guard on geometry, while live headroom in the actual render is 1.11 units. Before treating a future failure there as a regression, re-measure: the margin is a property of the model, not of the picture.

## Orientation

The Durations mark-label face now has **one** measured descent bound instead of two disagreeing ones, and the line-box bound reflects a current browser; lives in the board tool's Durations test constants, governed by `_dev/primes/prime-kanban-board.md`. No behaviour changed — every touched symbol is a test-side constant or a comment. [MAP CHANGED] in one narrow sense worth flagging: `durationsMeasuredMarkLabelDescentUnits` moved from an annotation-scoped constant to the package's single bound for that face, read across two files — which is why REQ-277 exists to say so at its declaration.

