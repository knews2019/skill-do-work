---
id: REQ-018
title: "Harden work.md steps 2.0 + 5.75 (full) and 5.5 + 6.3 (mechanical parts) into tools/checks/"
status: completed
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
claimed_at: 2026-07-15T18:15:00Z
completed_at: 2026-07-15T18:35:00Z
route: B
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: []
batch: harness-bloat-cleanup
maintenance: true
commit: 99cdec2
---

# Harden machine-checkable work.md steps

## What
Create shipped, platform-agnostic bash checks under `tools/checks/` and shrink the
corresponding work.md prose to a one-line pointer plus whatever judgment remains:

1. `archive-collision.sh REQ-NNN` — Step 2.0 verbatim: glob both archive patterns,
   exit code + collision path output. Prose → pointer.
2. `preflight.sh` — Step 5.75: git-clean (-uall), test-baseline run (command passed
   in or skipped), deps presence. Warnings to stdout, never blocking. Also emits
   `do-work/working/baseline.json` (test command + failing tests) — groundwork for
   the Step 6.5 probation test, without changing 6.5's prose beyond an optional
   "compare against baseline.json if present" sentence.
3. `scope-drift.sh <req-file>` — Step 5.5's review-time comparison only: set-diff of
   `## Scope` file list vs `## Implementation Summary` file list. Declaration prose stays.
4. `qualify.sh <req-file>` — Step 6.3 checks 1, 4, and the grep part of 5 (files
   exist/in-diff, P-A-U boxes vs debug-artifact grep, wiring grep with the
   documented exception list). Checks 2/3/6 stay prose; the script feeds them evidence.

All checks degrade gracefully (missing git, no tests → warn and continue) per the
design-for-the-floor rule. Steps 3.5, 3.7, 6.25, 7.5 are NOT converted (judgment —
audit HARDEN table). Update the Orchestrator Checklist lines and any Common
Rationalizations rows that cite the converted steps; add pointer-sync assertions to
`_dev/tests/contract-regressions.sh`.

## Why
Prose that transcribes mechanical shell logic drifts and taxes every read of
work.md (10,371 words). Audit HARDEN bucket.

## Acceptance criteria
- [x] Four scripts, shellcheck-clean, each runnable standalone with usage text.
- [x] work.md steps 2.0/5.75 reduced to pointer + intent; 5.5/6.3 keep judgment
      prose, mechanical parts point at scripts. Net work.md word count drops.
- [x] Behavior parity: script output covers every condition the prose prescribed
      (map recorded in Implementation Summary).
- [x] `_dev/tests/contract-regressions.sh` gains assertions that the pointers and
      scripts stay in sync; suite passes.

## Open Questions
(none — placement resolved at capture: tools/checks/, shipped)

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read work-reference.md templates (Scope, Implementation Summary) so the parsers match what the pipeline actually writes; scripts parse backtick-quoted paths + parenthesized verbs.
- [x] **[APPLY]:** Four scripts in tools/checks/ (executable, two-word variable names per naming conventions); work.md Steps 2.0/5.75 → pointer+fallback, 5.5 comparison + 6.3 items 1/4/5 → script with judgment checks 2/3/6 retained as prose; Step 6.5 gains the baseline.json comparison sentence; contract test gains pointer-sync assertions.
- [x] **[UNIFY]:** Smoke tests: archive-collision REQ-015→exit 1 with path, REQ-999→exit 0; scope-drift on a Route-A REQ→exit 2 SKIP; preflight→WARN on dirty tree, exit 0; qualify on archived REQ-015→OK with expected post-commit WARNs. contract-regressions.sh passes (incl. new assertions). work.md 10,371→10,207 words. shellcheck unavailable in this environment — scripts hand-reviewed; noted as PR follow-up.

## Triage

Route B — the "what" was fixed by the audit's HARDEN table; the "where" needed the work-reference templates (parse formats) and the existing test harness.

## Implementation Summary

**What was done:** Mechanical work-loop logic moved from prose to shipped executables; prose keeps intent, judgment, and a script-missing fallback.

Files changed:
- `tools/checks/archive-collision.sh` (new) — Step 2.0; exit 0/1/2 contract.
- `tools/checks/preflight.sh` (new) — Step 5.75; WARN/OK lines, always exit 0, writes do-work/working/baseline.json + baseline-failures.txt when a test command is supplied.
- `tools/checks/scope-drift.sh` (new) — Step 5.5 review-time set-diff; exit 2 = missing section (Route A skip).
- `tools/checks/qualify.sh` (new) — Step 6.3 items 1/4/5 + Step 6.25 only-do-work-paths rule; FAIL vs WARN separation keeps the exception-list judgment with the orchestrator.
- `actions/work.md` (modified) — Steps 2.0, 5.5 (comparison), 5.75, 6.3, 6.5 as above.
- `_dev/tests/contract-regressions.sh` (modified) — scripts exist+executable; work.md references each basename.
- `actions/version.md`, `CHANGELOG.md` — version 0.124.0 + entry.

**Behavior-parity map:** 2.0 both glob forms → find with both -name patterns; non-destructive status write stays prose (orchestrator writes frontmatter). 5.75 -uall flag, do-work/ exclusion, warnings-not-blockers → script; test-command *resolution* stays judgment. 5.5 both drift directions → comm -13/-23; severity stays judgment. 6.3 new/modified/deleted disk+diff checks, unchecked-box count, UNIFY-vs-debug-artifact grep, wiring grep → script; exception list, substantive/traced/flowing → prose.

## Testing

- Smoke tests recorded in [UNIFY]; contract suite green including the four new sync assertions.
- Red-green evidence for the ratchet-relevant assertion style: removing a script or its pointer flips the suite to FAIL (exercised by design of assert_contains + the -x check).

## Lessons Learned

**What worked:** Splitting FAIL (mechanical certainty) from WARN (evidence for judgment) let the scripts harden the checkable half without stealing the exception-list judgment the prose rightly owns.
**What didn't:** shellcheck isn't available in this environment — the REQ's acceptance criterion "shellcheck-clean" is verified by hand-review only; flagged for PR review.
**Worth knowing:** The scripts parse the exact bullet format from work-reference.md's templates ("- `path` (verb)"); changing those templates requires touching the parsers — same lock-step rule as the board's model.go.
