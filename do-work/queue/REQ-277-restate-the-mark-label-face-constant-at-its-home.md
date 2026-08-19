---
id: REQ-277
title: State the mark-label face constant's real scope at its canonical home
status: pending
created_at: 2026-08-18T23:53:40Z
status_changed_at: 2026-08-18T23:53:40Z
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
depends_on: []
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `generate_test.go:1968-1971` — the constant's doc comment still scopes it to the annotation; name both consumers and the `.durations-mark-label` face it actually bounds.
- [ ] `generate_test.go:1924` — the block header still scopes the block to the annotation and its neighbours.
- [ ] `durations_test.go:459-466` — the block header says every number in it is "rounded AWAY from the model (up)", but the block now holds two epistemically different constants: `durationsMeasuredLabelWidthSupremumUnits` is a supremum closed by an argument, while `durationsMeasuredLabelBoxHeightUnits` is explicitly **not** a supremum. Each says so locally; the header does not distinguish, and that distinction is the block's most important idea.

## Requirements

- A constant read from more than one file says so **at its declaration**, not only at its consumers — model the wording on `durationsMeasuredAxisTitleAscentUnits`, which already does this.
- The width block's header stops implying one epistemic status for constants that have two.
- **Sweep rather than patch three lines:** these were found by one reviewer reading one diff, so check every measured-face constant in the package for a comment that scopes it more narrowly than its actual readers. Report what was found either way.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-265's independent review, Important finding 2 plus a Minor (gate: trivial each, so they fold into one sweep rather than earning separate REQs).

Worth knowing: `TestDurationsMeasuredConstantsNameTheirChromiumBuild` enforces that each measured constant names its build, and its vacuity guard is `count == 0` — so it survived REQ-265's deletion without hollowing out. Nothing, however, checks that a constant's stated *scope* matches its readers. That gap is why this REQ exists.
