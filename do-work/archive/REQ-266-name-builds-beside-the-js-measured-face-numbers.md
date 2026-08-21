---
id: REQ-266
title: Name builds beside the JS renderer's measured face numbers
status: completed
created_at: 2026-08-18T20:07:08Z
status_changed_at: 2026-08-20T13:21:13Z
claimed_at: 2026-08-21T01:41:56Z
completed_at: 2026-08-21T01:47:35Z
kb_status: pending
commit: 8fe9eb6
route: A
user_request: UR-051
addendum_to: REQ-252
domain: general
review_generated: true
sweep: true
sweep_key: durations-measured-face-constants-lack-provenance
effort_estimate: normal
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: [REQ-292]
maintenance: false
write_set:
- skills/do-work-board/tools/queue-kanban/web/board-durations.js
- skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go
---

# Name Builds Beside the JS Renderer's Measured Face Numbers

## What

`web/board-durations.js` presents measured face numbers (12.83 / 10.43 / 2.41 at `DURATIONS_LABEL_ROW_HEIGHT` and `DURATIONS_LABEL_TEXT_ASCENT`) as current fact with no build named — the same provenance gap REQ-252 closed in the Go files, on the JS surface its go/parser test cannot reach. Extend the rule: every browser-measured number in the JS comments names its build, and the mechanism that keeps it true (a JS-side check, or a stated review convention) is the builder's call, recorded either way.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-252's declared Discovered Task, verified by its review (F1b, gate: rule-change). Carries the sweep key because it is the same root cause as REQ-252 — that REQ is the claimed-and-archived sweep for the key, so this is a new file per the append rule, not an append. Created `pending-answers` per the generation-≥2 depth stop.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-252: the JS renderer's measured numbers carry no build provenance. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the Go-side comments now carry the builds and the JS numbers are near-duplicates of them.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

REQ-292 may delete the subject of this REQ entirely. Re-read this REQ against the post-292
tree before starting it; closing as no-longer-applicable is the expected outcome.

**Cancelled and restored (2026-08-20).** This REQ was cancelled as superseded by REQ-292
earlier in the same session, then restored at the user's direction because REQ-292 has not
built yet. The `## Cancelled` section was removed on restore rather than left standing over a
`pending` REQ, which would have been a contradiction; the supersession reasoning it carried is
preserved in the Context line above and in the `depends_on: [REQ-292]` gate. Closing this REQ
as no-longer-applicable after REQ-292 lands remains the expected outcome.

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names one file and one rule, and REQ-292 (landed immediately before this, in the same run) removed most of its subject. What remained was one comment and the decision about the mechanism.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Re-read Against the Post-292 Tree

The REQ asks for exactly this before starting, and says closing as no-longer-applicable is the expected outcome. **It is not — but it is close.** Of the three numbers the REQ names:

| Number | Post-292 state |
|---|---|
| `12.83` at `DURATIONS_LABEL_ROW_HEIGHT` | **gone** — REQ-292 rewrote that comment; the pitch is now checked against the box the engine draws |
| `2.41` | **gone** with the same comment |
| `10.43` at `DURATIONS_LABEL_TEXT_ASCENT` | **still there**, still presented as current fact, still with no build beside it |

So one instance survived, and its comment also cited `TestDurationsLabelRowsClearTheMarkBands` — a test REQ-292 deleted an hour earlier. Closing this REQ unbuilt would have left both.

## Decisions

- **D-01** (DECIDE & STATE): The surviving number was **deleted rather than dated**. Reasoning: REQ-292 established that the honest answer for a claim about the face in use is to let the engine answer at test time, and `TestBrowserBehaviorDurationsLabelRowsClearTheirNeighbours` already measures the real line box on whatever machine runs it. Digging a build out for `10.43` would have re-dated a number nothing reads. `DURATIONS_LABEL_TEXT_ASCENT = 11` stays as what it actually is: a declared number chosen at or above any ascent the face draws, with the round-up's reason kept.
- **D-02** (ESCALATE): The mechanism is a **check, not a review convention** — the REQ leaves this to the builder and requires it be recorded either way. `TestDurationsJavaScriptCommentsDateTheirMeasurements` scans the renderer's comments for measurement claims and fails any that no REQ reference or build name dates. Reasoning: a review convention is a rule every future author must remember, which is precisely how three undated numbers accumulated on a surface `go/parser` could not reach. Value: the rule now fails a build instead of relying on a reviewer noticing. Risk: it is a regex over prose, so it can misfire on an unusual phrasing; the failure is loud and the message says exactly what to do (cite the REQ, name the build, or delete the number).
- **D-03** (DECIDE & STATE): **The discriminator is what the number claims, not whether it is a measurement.** A measured number cited as evidence for a past decision is already dated by the REQ it names — a stronger anchor than a build string, because the REQ carries the whole argument. A number presented as current fact about the face in use is the one that needs a build. Without this distinction the check would demand a build beside every historical figure in REQ-292's own explanatory comments, which would be dating a story rather than a fact.

## Implementation Summary

**What was done:** Removed the last undated measured number from the renderer's comments and replaced its stale test citation, then made the rule enforceable on the JS surface with a check rather than a convention.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/board-durations.js` (modified) — `DURATIONS_LABEL_TEXT_ASCENT`'s comment no longer cites a measured ascent as current fact, no longer names a deleted test, and says what the constant actually is.
- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go` (modified) — `TestDurationsJavaScriptCommentsDateTheirMeasurements` and its comment-block helper.

**Tests touched:** one added. No existing assertion changed meaning.

## Qualification

Passed — 2 files verified, 3 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `gofmt -l .` clean, `go vet ./...` clean, `node --check` clean on the changed JS. No debug artifacts. `maintainer-verify.sh` exits 0.
- **Substantive:** the check parses real source and fails on a real defect shape (shown below); the comment change removes a false current-fact claim.
- **Requirements traced:** every browser-measured number in the JS comments is now dated or gone → the check passes on the tree; the mechanism is a JS-side check → D-02; recorded either way → D-02 and D-03.
- **Flowing:** not applicable — a source-scanning check with no data path.

## Testing

- `bash _dev/tests/maintainer-verify.sh` — exit 0. `go vet`, `gofmt`, `node --check` clean.

**The check was mutation-tested in both directions, and the first attempt exposed its own limit:**

| Mutation | Result |
|---|---|
| Reintroduce the exact REQ-266 defect (`measured 10.43-unit ascent`) into the ascent comment | **PASSED** — the block already names REQ-266 and REQ-292 for other reasons, so the block-level anchor accepted it |
| A bare measurement in a comment block with no REQ or build anywhere | **FAILED**, naming the file, the line and the text — the real defect shape |
| The same line dated with `on Chromium 141.0.7390.37` | passes, as it should |

The first row is a genuine weakness and is reported as M1 rather than papered over: block-level anchoring means a paragraph that names a REQ for any reason dates every measurement inside it. Line-level anchoring was rejected because a paragraph legitimately names its REQ once and every other line would then need a repeat.

## Review

**Overall: 90%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Every browser-measured number in the JS comments names its build | ✅ by deletion for the one that remained (D-01); the check enforces it going forward |
| The mechanism that keeps it true is the builder's call | ✅ a check, not a convention (D-02) |
| Recorded either way | ✅ D-02 and D-03 |

### Findings

**Important — none.**

**Minor:**

- **M1:** The check anchors at comment-block level, so a paragraph naming a REQ for any reason dates every measurement inside it. Demonstrated in `## Testing` — reinserting the original defect into the rewritten ascent block passes. Line-level anchoring would be stricter and wrong for a different reason (a paragraph names its REQ once). The practical consequence is narrow: the check catches a new undated measurement in a new or unannotated comment, which is how all three of the originals arrived, and does not catch one smuggled into an existing annotated paragraph.
- **M2:** The REQ expected "closing as no-longer-applicable" and it was built instead. Justified by the re-read table — one instance genuinely survived — but a reviewer expecting a closure will want that table.

**Nit:**

- **N1:** `TestDurationsJavaScriptCommentsDateTheirMeasurements` uses two regexes (number-then-word and word-then-number) rather than one alternation, because the two orders read more clearly apart than as a single pattern with a backreference-free either-side construction.

### Restatement Sweep

Redefined element: the rule about measured numbers on the JS surface — previously "name the build", now "name the build **or** the REQ that measured it, and prefer deleting a current-fact claim over dating it".

- `skills/do-work-board/tools/queue-kanban/durations_browser_probe_test.go`'s `TestDurationsCarriesNoMeasuredFaceConstants` (REQ-292's guard) — the neighbouring rule. Re-read: it bans specific retired *constant names*, while this bans undated *measurements in prose*. Complementary, not overlapping, and neither restates the other.
- REQ-252's Go-side rule — grepped for its statement in `durations.go` and `durations_test.go`: the Go-side measured constants are all gone after REQ-292, so nothing there restates a rule this one contradicts.
- `_dev/primes/prime-kanban-board.md` — checked for a statement of the provenance rule; it does not carry one, so nothing went stale.

No stale restatement remains. One was *found and fixed* as part of the work rather than by this sweep: the ascent comment cited `TestDurationsLabelRowsClearTheMarkBands`, deleted by REQ-292 an hour earlier.

### Acceptance Testing

The check was run against the real renderer (passes), against a reintroduction of the original defect in an unanchored block (fails, naming file, line and text), and against the same line with a build named (passes). The comment change was read in place to confirm it still explains why the constant is 11 without asserting a measurement.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope Discipline | 100% |
| Risk | None |
| Acceptance | Pass |

Test Adequacy 85% for M1 — the check is real and mutation-proven, but its anchoring is more permissive than the ideal.

### Follow-up REQs Created

None. M1 is a stated limit of the chosen mechanism rather than a defect to queue; tightening it has no known trigger.

## Lessons Learned

**What worked:** Doing the re-read the REQ asked for instead of trusting its own expected outcome. The REQ said twice that closing as no-longer-applicable was expected, and REQ-292 had just landed — it would have been easy and defensible to close it. One of the three instances had survived, along with a citation to a test deleted an hour earlier in the same run.

**What didn't:** The first mutation test of the new check passed when it should have failed, because the defect was reinserted into a comment block that already named two REQs. That is the check's real anchoring behaviour, not a mistake in the mutation — and finding it took writing a mutation that *should* fail and being surprised. A mutation test that passes is either a working guard or a broken check, and the only way to tell is to write the second mutation.

**Worth knowing:** The useful distinction here generalises past this file. A number in a comment is doing one of two jobs: asserting something true *now* about the environment, or citing evidence for a decision *then*. The first goes stale invisibly and needs dating; the second is already dated by the decision it supports. Demanding provenance for both makes the rule annoying enough to be ignored, which is roughly how the undated numbers got there.

## Orientation

The Durations renderer's comments no longer state a measured face number as current fact without saying where it came from, and the rule is now enforced by a check rather than by remembering. Lives in the board's Durations subsystem (`skills/do-work-board/tools/queue-kanban/`), indexed by `_dev/primes/prime-kanban-board.md`. Leaf change: no contract, no field, no renamed concept — REQ-292 had already done the structural work this one finishes.
