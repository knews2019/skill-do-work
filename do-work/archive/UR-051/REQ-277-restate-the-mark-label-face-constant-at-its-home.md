---
id: REQ-277
title: State the mark-label face constant's real scope at its canonical home
status: completed
created_at: 2026-08-18T23:53:40Z
status_changed_at: 2026-08-20T13:21:13Z
claimed_at: 2026-08-21T01:50:55Z
completed_at: 2026-08-21T01:56:42Z
kb_status: promoted
kb_entry: REQ-277-state-the-mark-label-face-constant-s-rea.md
commit: 54282b0
route: B
user_request: UR-051
addendum_to: REQ-265
domain: general
review_generated: true
sweep: true
sweep_key: measured-face-comment-accuracy
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec: bug-fix
depends_on: [REQ-292]
maintenance: true
write_set:
- skills/do-work-board/tools/queue-kanban/generate_test.go
- skills/do-work-board/tools/queue-kanban/durations_test.go
---

# State the Mark-Label Face Constant's Real Scope at Its Canonical Home

## What

REQ-265 consolidated two duplicate descent bounds into one: `durationsMeasuredMarkLabelDescentUnits` in `generate_test.go` is now the package's **single** bound for the `.durations-mark-label` face, read from **two** files. Its documentation did not follow it.

- `generate_test.go:1968-1971` still calls it "The annotation box's descent", and its block header at `:1924` still scopes the block to "panel B's slowest-day annotation, and the faces around it". Neither says the constant is now read cross-file.
- `durations_test.go:616-627` explains the cross-file dependency thoroughly — **from the consumer side only.** The canonical home says nothing.

That is the restatement half of the very consolidation REQ-265 performed, and it is **the same shape that made REQ-241 and REQ-242 collide**: a constant's real ownership stated in one file and not the other. The sibling `durationsMeasuredAxisTitleAscentUnits` gets this right and explains why it lives package-wide; this one does not.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [x] `generate_test.go:1968-1971` — the constant's doc comment still scopes it to the annotation; name both consumers and the `.durations-mark-label` face it actually bounds. — **moot, verified by measurement.** REQ-292 deleted `durations_test.go`'s consumers, so the constant is now declared and read in `generate_test.go` alone. Its doc comment's narrow scope became accurate rather than needing widening.
- [x] `generate_test.go:1924` — the block header still scopes the block to the annotation and its neighbours. — **moot for the same reason, but the header carried a DIFFERENT stale claim the REQ did not know about, and that one was fixed.** See Decisions D-02.
- [x] `durations_test.go:459-466` — the block header says every number in it is "rounded AWAY from the model (up)", but the block now holds two epistemically different constants: `durationsMeasuredLabelWidthSupremumUnits` is a supremum closed by an argument, while `durationsMeasuredLabelBoxHeightUnits` is explicitly **not** a supremum. Each says so locally; the header does not distinguish, and that distinction is the block's most important idea. — **gone.** REQ-292 deleted both constants and the block that held them. Confirmed by grep: neither is declared anywhere.

## Requirements

- A constant read from more than one file says so **at its declaration**, not only at its consumers — model the wording on `durationsMeasuredAxisTitleAscentUnits`, which already does this.
- The width block's header stops implying one epistemic status for constants that have two.
- **Sweep rather than patch three lines:** these were found by one reviewer reading one diff, so check every measured-face constant in the package for a comment that scopes it more narrowly than its actual readers. Report what was found either way.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-265's independent review, Important finding 2 plus a Minor (gate: trivial each, so they fold into one sweep rather than earning separate REQs).

Worth knowing: `TestDurationsMeasuredConstantsNameTheirChromiumBuild` enforces that each measured constant names its build, and its vacuity guard is `count == 0` — so it survived REQ-265's deletion without hollowing out. Nothing, however, checks that a constant's stated *scope* matches its readers. That gap is why this REQ exists.

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

**Route: B** - Medium

**Reasoning:** Two of the three instances were expected to be moot after REQ-292, so the deliverable was really the sweep the Requirements ask for — check every measured-face constant in the package for a comment scoping it more narrowly than its readers, and report either way. That is discovery work, and it found something the REQ's instance list could not have named.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## The Sweep

The REQ asks for every measured-face constant in the package to be checked for a comment that scopes it more narrowly than its actual readers, and for the result to be reported either way. Enumerated mechanically — declaration site against every file that reads the name:

| Constant | Declared in | Read in | Comment's stated scope |
|---|---|---|---|
| `durationsMeasuredAxisTitleAscentUnits` | `generate_test.go` | `generate_test.go` | accurate (and explains itself, as the REQ notes) |
| `durationsMeasuredAxisTitleDescentUnits` | `generate_test.go` | `generate_test.go` | accurate |
| `durationsMeasuredMarkLabelAscentUnits` | `generate_test.go` | `generate_test.go` | accurate **now** — the cross-file reader was deleted by REQ-292 |
| `durationsMeasuredMarkLabelDescentUnits` | `generate_test.go` | `generate_test.go` | accurate **now**, same reason |
| `durationsMeasuredLabelWidthSupremumUnits` | — | — | **deleted by REQ-292** |
| `durationsMeasuredLabelBoxHeightUnits` | — | — | **deleted by REQ-292** |

**Result: no constant's comment scopes it more narrowly than its readers.** All three of the REQ's instances are moot, and the sweep found no fourth.

**But the sweep found a different defect of the same class, in the same block**, which is the thing worth having run it for. See D-02.

*Sweep run inline by the orchestrator*

## Decisions

- **D-01** (DECIDE & STATE): The two mark-label constants' doc comments were **left alone**. Reasoning: the REQ asked to widen them to name both consumers and the cross-file dependency. There is no longer a second consumer — REQ-292 deleted it — so widening them would have introduced exactly the defect this REQ exists to remove: a comment claiming an arrangement that is not real. The correct action for a moot instance is to verify and record, not to perform the edit it asked for.
- **D-02** (ESCALATE): The block header's last paragraph named `durations_test.go`'s `TestDurationsMeasuredConstantsNameTheirChromiumBuild` as the thing enforcing the name-your-build rule. **REQ-292 deleted that test earlier in this same run**, on the recorded reasoning that "no such constant survives, so there is nothing left to name one". That reasoning was true of `durations.go` and `board-durations.js` — the files REQ-292 was clearing — and **not true of `generate_test.go`, where four measured constants live and are read.** So REQ-292 removed an enforcer that still had subjects, and left the rule's stated enforcement pointing at nothing.
  Builder chose to **restore the enforcement here**, as `TestMeasuredFaceConstantsNameTheirBuild` in `generate_test.go` beside the constants it governs, and to correct the header to point at it. Reasoning: this is precisely the defect class REQ-277 was written about — a comment stating an arrangement that is no longer real — found by the sweep it asked for, so it belongs to this REQ rather than to a new one. Value: the rule REQ-241 and REQ-242 collided to earn is enforced again rather than merely described. Risk: none identified; the check is additive and its vacuity guard is carried over from the deleted test.
  **REQ-292's archived record now overstates its case.** Archived REQs are immutable, so the correction is recorded here and surfaced at hand-back rather than edited into it.

## Implementation Summary

**What was done:** Ran the sweep the REQ requires and reported it in full; verified all three named instances are moot post-REQ-292 rather than performing edits that would have made them false; and restored the build-provenance enforcement that REQ-292 removed while four of its subjects were still standing, correcting the comment that still pointed at the deleted test.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — the block header now names the live enforcer and says why it moved; `TestMeasuredFaceConstantsNameTheirBuild` and `readPackageSourceForTest` added.

**Files deliberately unchanged:** `skills/do-work-board/tools/queue-kanban/durations_test.go` — declared in the REQ's write set. Its instance's subject was deleted by REQ-292 and there is nothing left there to restate (D-01).

**Tests touched:** one added, replacing one REQ-292 deleted. No existing assertion changed meaning.

## Qualification

Passed — 1 file verified, 4 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `gofmt -l .` clean, `go vet ./...` clean, `go test ./...` ok. No debug artifacts. `maintainer-verify.sh` exits 0.
- **Substantive:** a real check that parses real source and fails on two distinct real defect shapes (below).
- **Requirements traced:** "a constant read from more than one file says so at its declaration" → no such constant remains, verified by the sweep table; "the width block's header stops implying one epistemic status" → that block no longer exists; "sweep rather than patch three lines, report either way" → the table, plus D-02's finding; "verify exits 0" → it does.
- **Flowing:** not applicable — a source-scanning check.

## Testing

- `bash _dev/tests/maintainer-verify.sh` — exit 0. `go vet`, `gofmt`, `go test ./...` clean.

**Mutation-tested in both directions:**

| Mutation | Result |
|---|---|
| Strip the build names from one constant's doc comment | **FAIL**, naming the constant and why a per-browser face needs its build |
| Point the scan at a prefix that matches nothing (the vacuity case) | **FAIL** — "this scan must never pass on an empty set" |

The second is the guard the deleted test carried, kept for the same reason: without it the check goes green over an empty set the day someone renames the prefix, which is how an enforcement rule dies quietly.

## Review

**Overall: 92%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| A constant read from more than one file says so at its declaration | ✅ vacuously — the sweep shows none is, post-REQ-292 |
| The width block's header stops implying one epistemic status | ✅ the block and both its constants were deleted by REQ-292 |
| Sweep every measured-face constant rather than patching three lines; report either way | ✅ the table, and it found the defect the instance list could not name |
| `maintainer-verify.sh` exits 0 | ✅ |

### Findings

**Important:**

- **I1 (fixed in this REQ):** REQ-292, earlier in this same run, deleted `TestDurationsMeasuredConstantsNameTheirChromiumBuild` with the recorded justification that no measured constant survived to name a build. Four survive in `generate_test.go` and are read there. So an enforcement rule that two REQs collided to earn was removed while its subjects stood, and its stated enforcement was left pointing at a deleted test. Restored here as `TestMeasuredFaceConstantsNameTheirBuild` (D-02). Recorded as Important because a silently-removed guard is worse than one that never existed: the comment kept asserting it was enforced.

**Minor:**

- **M1:** The restored check scans `generate_test.go` only. The prefix could be declared in another file tomorrow and go unchecked. Scanning every `*_test.go` was considered and rejected as speculative — one file declares them all today, and the vacuity guard fails loudly if that stops being true in the direction that matters (the prefix disappearing). A constant appearing in a *new* file is the case it would miss.
- **M2:** The REQ's three instances all resolved to "moot" and nothing it named was edited. A reader comparing the instance list against the diff will see no correspondence; the sweep table is what connects them.

**Nit:**

- **N1:** `readPackageSourceForTest` reads a `_test.go` file off disk because `go:embed` cannot reach one. Fine, and commented, but it means the check depends on the test binary's working directory being the package directory — which `go test` guarantees.

### Restatement Sweep

Redefined element: which test enforces the name-your-build rule, and where it lives.

- `generate_test.go`'s block header — **rewritten**; it was the stale restatement, and it is the reason this finding was findable at all.
- `durations_test.go` — grepped for any surviving reference to the deleted enforcer or to the deleted label constants: none.
- `durations_browser_probe_test.go` — REQ-292's `TestDurationsCarriesNoMeasuredFaceConstants` and REQ-266's `TestDurationsJavaScriptCommentsDateTheirMeasurements` are the neighbouring rules. Re-read all three together: the first bans retired constant *names* from the shipped files, the second dates *measurements in JS prose*, the third requires *build provenance on the Go constants that remain*. Three distinct subjects, no overlap, and none restates another.
- `_dev/primes/prime-kanban-board.md` — carries no statement of the provenance rule, so nothing there went stale.

No stale restatement remains.

### Acceptance Testing

The sweep is the deliverable and was run mechanically rather than by reading — declaration site against every reader, for every `durationsMeasured*` name in the package. The restored check was run against the real tree (passes), against a stripped build name (fails, naming the constant), and against an empty match set (fails on the vacuity guard).

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope Discipline | 95% |
| Risk | None |
| Acceptance | Pass |

### Follow-up REQs Created

None. I1 was fixed here rather than queued, because it is the same defect class this REQ exists to sweep for and was found by the sweep it required.

## Lessons Learned

**What worked:** Doing the sweep mechanically instead of reading the three named lines. Enumerating every `durationsMeasured*` name against its readers took one shell loop, proved all three instances moot in one table, and — because the table made the block's contents concrete — led straight to reading the block header that turned out to be the real defect. The REQ predicted its own instances would be moot and asked for the sweep anyway; that instruction is what earned the finding.

**What didn't:** Nothing failed in this REQ's own work, but it caught a mistake made a few hours earlier in this same run. REQ-292's justification for deleting a guard — "no such constant survives" — was checked against the two files that REQ was clearing and not against the package. The lesson is narrow and repeatable: **when a deletion is justified by "nothing uses this any more", the scope of that claim has to be the scope of the thing being deleted.** The guard was package-wide; the check was file-wide.

**Worth knowing:** A removed guard is more dangerous than an absent one, because the prose that described it usually survives. Here the block header went on asserting that every measured constant's build was enforced for hours after the enforcement was gone — and it read as true, because it had been. When deleting a check, grep for what *claims* it exists.

## Orientation

The rule that a measured face number must name the browser build it came from is enforced again, by a check that lives beside the constants it governs instead of in a file that no longer holds any. The sweep the REQ asked for is recorded in full and found the package's remaining measured constants all correctly scoped. Lives in the board tool's test suite (`skills/do-work-board/tools/queue-kanban/generate_test.go`), indexed by `_dev/primes/prime-kanban-board.md`. Leaf change: no contract, no field, no renamed concept.
