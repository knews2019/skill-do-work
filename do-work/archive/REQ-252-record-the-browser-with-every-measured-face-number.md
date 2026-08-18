---
id: REQ-252
title: Record the browser with every measured-face number in the Durations tests
status: completed
created_at: 2026-08-18T13:56:12Z
claimed_at: 2026-08-18T19:12:47Z
completed_at: 2026-08-18T20:08:09Z
kb_status: pending
route: B
status_changed_at: 2026-08-18T13:56:12Z
user_request: UR-051
addendum_to: REQ-241
domain: general
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: durations-measured-face-constants-lack-provenance
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-248]
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/durations_test.go
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/durations.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T19:13:28Z
  basis:
    - Route B
    - 3-file write set
    - 3 acceptance criteria
    - browser evidence
---

# Record the Browser With Every Measured-Face Number in the Durations Tests

## What

The Durations suite now holds several constants that are a **browser's answer**, not arithmetic — and two independent measurements of the same face, taken on different Chromium builds, disagreed by enough to matter. Every such constant needs the build it came from recorded next to it, and the numbers that are now known to be per-browser need saying so.

## Context

Found when REQ-241 and REQ-242 both measured the 12px axis-title face and landed on different values, then declared the same constant in different files of one package. The collision surfaced at integration as a Go redeclaration — git could not see it, because the two edits never touched adjacent lines. It was resolved to the larger value (12.1) on the reasoning that a title box reaching higher makes every clearance test demand more room than the render needs.

That resolution is sound and stays. What it exposed is that these numbers carry no provenance, so the next disagreement looks like a mistake rather than a difference of build.

## Instances

- [x] **The Panel B clearance budget is per-browser and roughly 7× thinner on a current build than recorded.** REQ-241's D-03 measured 1.364 units of headroom above Panel B's title and named a row pitch of 15 as the point where it breaks. REQ-242's reviewer measured **0.185 units** on Chromium 146 — still positive, still zero intersections, and demonstrably not caused by REQ-242 (the SVG for that region is byte-identical across its range). Anyone about to spend that budget must re-measure first, and the prose should say so.
- [x] **No measured-face constant records the Chromium build it came from.** `generate_test.go`'s block documents the procedure and viewport, which is most of the way there; `durations_test.go`'s do not. A number that differs between builds and does not name its build cannot be re-derived or argued with.
- [x] **`durations.go`'s `formatDurationLabelMinutes` emits an ASCII hyphen where the renderer draws U+2212.** The two glyphs differ by 1.73 units in this face, so the Go side models a narrower string than the browser draws. **Currently width-neutral** — both are one character and the width model counts characters — but it becomes a real under-estimate the moment anyone replaces the flat constant with a per-glyph model, which REQ-241 attempted and abandoned. Worth a comment at minimum.
- [x] **`durationsLabelRemainderReserveUnits` over-reserves more at the new width.** It scales as `24 × durationsLabelCharacterWidthUnits`, so raising that constant to 7.15 widened the reservation to 171.6 against a widest composable remainder sentence of 123.18 units. Over-reserving is the safe direction — it drops a label, which the remainder counts — but the gap is now large enough to be a deliberate choice rather than an accident.

---

## AI Execution State (P-A-U Loop)

Added by the orchestrator at integration (review-generated REQ predating the block — the class REQ-264 asks to make visible). Transcribed from the builder hand-back:

- [x] **[PLAN]:** Read brief, crew rules (incl. testing for tdd), the board prime with its Lessons. Inventoried every `durationsMeasured*` constant (3 + 4) and every browser-derived number in `durations.go`; traced provenance in the live archives (REQ-241: "chromium-1228"; REQ-242: Chromium 146). Plan: scratch fixture board, isolated-Chromium measurement of current geometry, RED provenance test, comment edits, GREEN, full gates.
- [x] **[APPLY]:** Fixture board built in scratch (binary to scratch, never bare `go build`); measured on Chromium 141.0.7390.37 with `location.href` in the same evaluate as every number; RED test first (all 7 constants failed for the missing-build reason), then the comments.
- [x] **[UNIFY]:** `gofmt -l` clean, `go vet` clean, package suite exit 0, maintainer-verify exit 0; class sweep found the only remaining unprovenance'd citations in `web/board-durations.js` (Scope-excluded, Discovered Task). Two commits, each green on its own.

## Implementation Summary

**What was done:** Every browser-measured constant in the Durations suite (7 `durationsMeasured*` constants across `durations_test.go` and `generate_test.go`) now names the Chromium build its number came from, written exactly as the source REQs recorded it, each with a fresh cross-check on this environment's Chromium 141.0.7390.37. A new `TestDurationsMeasuredConstantsNameTheirChromiumBuild` walks every Go file in the package (go/parser, vacuity-guarded) and fails any measured constant whose doc comment names no build — observed RED on all 7 before the edits. The Panel B clearance budget prose states it is per-browser and carries all three measurements (1.364 REQ-241 / 0.185 Chromium 146 / 1.111 fresh, post-REQ-248 geometry) plus the corrected model-space figure (0.49 units, pitch 14 breaks it). `durations.go` documents the hyphen-vs-U+2212 divergence (5.24 units on this build vs the REQ's recorded 1.73 — the delta itself is per-browser) and states the remainder reserve's over-shoot as deliberate with measured numbers; REQ-260's gofmt nit fixed in passing. **No constant's value changed** — two current-build measurements exceed their constants and that raise is a Discovered Task, per the REQ's no-fold rule.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modified) — per-constant provenance, per-browser procedure note, budget prose, the provenance test
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — measured-face block split; each of four constants carries its own build comment with current-build cross-checks
- `skills/do-work-board/tools/queue-kanban/durations.go` (modified) — supremum's build named; hyphen/U+2212 comment; reserve over-shoot stated; gofmt fix

*Integrated by orchestrator from builder hand-back; merge range `ca0cc84..c752529`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01 (DECIDE & STATE):** the mechanical provenance check ships as a real test — the failure it pins is a real integration collision (REQ-241/242's invisible redeclaration); ~60 lines of std-lib go/parser; vacuity guard included; stated limit: the `durationsMeasured` prefix is convention, a smuggled number under another name is review's job (package-wide sweep confirms the convention currently holds).
- **D-02 (DECIDE & STATE):** no value change despite two current-build measurements exceeding their constants — the REQ forbids folding; no assertion currently flips; the comments state the raise-only status so they do not lie about currency.
- **D-03 (DECIDE & STATE):** REQ-260's gofmt nit fixed in passing as the Scope allows — REQ-260 can be discarded at clarify.
- **D-04 (DECIDE & STATE):** provenance identifiers written exactly as the source REQs recorded them ("browser build chromium-1228"; "Chromium 146 (Playwright 1.59)") — the surviving records are partly inconsistent and the comments say what is known rather than inventing a version.
- **D-05 (DECIDE & STATE):** the budget prose's stale figures corrected alongside the provenance — the old "1.36 / pitch 15" predated the 12.1 ascent merge; the model figure is 0.49 / pitch 14.

## Qualification

Passed — 3 files in merge range `ca0cc84..c752529`, requirements traced (all 7 constants carry builds — RED transcript names each; budget prose carries three measurements + warning; no behaviour change — `git diff` shows comments, test, and the one-character gofmt fix only), P-A-U audited. Armed qualify note: the scan WARNs on the provenance test's own output lines (it owns its process exit) — confirmed from the diff as the test's failure messages, i.e. contract output.

## Discovered Tasks

Captured durably per review F1 (they previously lived only in hand-back prose):

- **[normal]** Raise `durationsMeasuredLabelBoxHeightUnits` (12.84 → ≥12.97) and `durationsMeasuredLabelBoxDescentUnits` (2.41 → ≥2.78) — Chromium 141.0.7390.37 measures above both; no assertion flips today. → REQ-265 (pending-answers).
- **[normal]** `web/board-durations.js` measured numbers (12.83/10.43/2.41) carry no build provenance — the same gap on the JS surface. → REQ-266 (pending-answers, `sweep_key: durations-measured-face-constants-lack-provenance`).

## Review

**Overall: 97%** | 2026-08-18T20:04:59Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Verdict: Approve** — provenance recorded faithfully on all seven measured constants, the enforcement test is mutation-falsifiable and vacuity-guarded, and the no-behaviour-change claim is proven mechanically (comment-stripped ASTs identical across the range; suite exit 0).

**Requirements:** all three delivered — every constant carries its build (mutation-verified: stripping one comment FAILs naming that constant); the budget prose carries all three measurements with builds plus the re-measure warning, figures matching the archived records exactly; zero semantic diff beyond the new test and the two-space gofmt fix.

**Acceptance highlights (reproduced by execution):** every current-build figure in the builder's comments independently reproduced in a freshly launched isolated Chromium 141.0.7390.37 (Playwright 1.56.1) with `location.href` in the same evaluate — line box 12.9631, title ascent 11.1112, live Panel B headroom 1.1110. Direction-of-error audit of the two under-bounding constants: both err unsafe in isolation, no consumer live-wrong on any recorded build (pitch 13 clears the real box; the ceiling's paired over-bounding ascent compensates by 0.99 — a coincidence of one consumer, hence REQ-265). Restatement sweep clean: surviving old figures are provenance'd historical citations only.

**Important findings (audit record):**
- F1a: the declared raise Discovered Task existed in no durable artifact — gate: trivial → REQ-265 (pending-answers, generation ≥2)
- F1b: JS measured numbers unprovenance'd — gate: rule-change → REQ-266 (pending-answers, carrying the sweep key; new file per the append rule since this REQ is the claimed sweep)

**Minor:** 3 (report only) — a future grouped `const (...)` block comment could satisfy every constant inside it (latent; all 7 individually documented today); the build matcher checks presence, not well-formedness; "one character" in the trail understates the two-space gofmt fix. **Nit:** 1.

**Follow-ups created:** REQ-265, REQ-266 (by orchestrator) · REQ-260 is satisfied by this merge — discard at clarify (D-03).

*Reviewed by review-work action (independent adversarial pass, orchestrated mode; merge range `ca0cc84..c752529`)*

## Lessons Learned

**What worked:** Enforcing a documentation convention with a real AST-walking test (vacuity-guarded, mutation-falsifiable) instead of trusting comments to stay honest — RED on all seven constants proved the gap before any edit. Recording provenance identifiers exactly as the archives state them, inconsistencies included, rather than inventing tidier versions.

**What didn't:** Discovered Tasks that live only in hand-back prose are one integration slip from evaporating — the review caught that neither had a durable artifact. Capture them in the REQ's own section at hand-back time.

**Worth knowing:** The 11px mark-label box measures LARGER than its recorded constant on Chromium 141 (12.9631 vs 12.84) — the raise is REQ-265; until it lands, the pitch-13 floor clears reality by 0.037 units, not the 0.16 the constant implies. Even the hyphen-vs-U+2212 delta is per-browser (1.73 recorded vs 5.24 measured here). The `durationsMeasured` prefix is a convention the test enforces comments for; a smuggled number under another name is review's job.

## Orientation

Now every browser-measured constant in the Durations Go suite names the build it came from, a package-wide test holds the convention, and the Panel B budget prose tells the next builder to re-measure before spending. Lives in the board's Durations test layer. Leaf change (documentation + one read-only test); map unchanged.

## Requirements

- Every constant in these files that is a browser measurement names the browser and build it was taken on, in the same comment as the number.
- The Panel B clearance budget's prose states that it is per-browser and gives both measurements, so nobody spends a budget that does not exist on their machine.
- No behaviour change is required by this REQ. If a measurement turns out to demand one, capture that separately rather than folding it in.

## Red-Green Proof

**RED prompt/case:** this is largely a documentation-fidelity REQ; the mechanical half is whether any assertion currently reads a measured constant without a recorded provenance. A check that every `durationsMeasured*` constant's comment names a browser would fail today.
**Why RED now:** `durations_test.go`'s measured constants name no build, and the two files disagreed by 0.86 units on the same face.
**GREEN when:** each constant carries its provenance and the suite still exits 0.
**Validation:** Integration seam on REQ-241/REQ-242 plus REQ-242's review finding F3.

## Ordering Gate

- [~] **D-01: gated behind REQ-248 rather than left free-running.** Two reasons, and the second is the binding one. Textual: both REQs write `skills/do-work-board/tools/queue-kanban/generate_test.go`, and this session proved that two REQs in one Go package can collide in a way git cannot see — REQ-241 and REQ-242 each declared `durationsMeasuredAxisTitleAscentUnits` in different files, never touched adjacent lines, and the merge failed to compile. Semantic: this REQ records the provenance of measurements taken against Panel B's geometry, and REQ-248 moves that geometry. Recording provenance for a layout about to change would produce numbers that are wrong on arrival.
  **Value:** the one-line `do-work run` resume cannot schedule these two together, because auto-wave reads `depends_on` and deliberately does not read `write_set`.
  **Risk:** if REQ-248 is abandoned, this REQ stays gated behind it and needs the gate cleared by hand. Acceptable — REQ-248 is the highest-value item in the queue.

---

## Triage

**Route: B** - Medium

**Reasoning:** Documentation-fidelity sweep over measured constants plus a fresh per-build measurement of the Panel B budget; the what is clear, the current numbers need deriving in the live DOM on this build.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/durations_test.go` (modify) — every measured constant names its browser build
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — same, and the Panel B budget prose gains both measurements + the per-browser warning
- `skills/do-work-board/tools/queue-kanban/durations.go` (modify) — hyphen-vs-U+2212 comment; remainder-reserve over-reservation stated as deliberate; the REQ-260 gofmt nit at the `Truncate(24*time.Hour)` line may be fixed in passing (same file, one character)

**Files I will NOT touch:**
- `web/board-durations.js` — REQ-248's geometry landed; no behaviour change belongs to this REQ.
- Any constant's *value* — this REQ records provenance; a measurement demanding a change is captured separately.

**Acceptance criteria (restated from REQ):**
- [ ] Every browser-measured constant names the browser and build in the same comment as the number.
- [ ] The Panel B clearance budget prose states it is per-browser and gives both measurements.
- [ ] No behaviour change; suite still exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (0.212.14 tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium present

*Checked by work action*
