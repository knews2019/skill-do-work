---
id: REQ-362
title: "[impact-rule-change] Stop a multi-path bullet disabling the scope-drift check"
status: claimed
claimed_at: 2026-08-24T18:57:05Z
status_changed_at: 2026-08-24T18:57:05Z
route: B
created_at: 2026-08-24T11:15:00Z
user_request: UR-068
addendum_to: REQ-344
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-24T18:57:05Z
  basis:
    - Route B
    - 1-file write set
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
  - skills/do-work/tools/checks/scope-drift.sh
  - skills/do-work/tools/checks/associate-files.sh
  - skills/do-work/tools/checks/qualify.sh
  - _dev/tests/contract-regressions.sh
  - _dev/tests/prescribed-shell-cases/qualify.sh
---

# Stop a Multi-Path Bullet Disabling the Scope-Drift Check

## What

`scope-drift.sh` extracts only the **first** backticked path per Implementation Summary bullet, so a
bullet listing several files hides every path after the first. Hand-formatting a summary silently
disables the measurement that exists to make drift visible.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Traced both first-token parsers and the associate-files sibling, then froze paired-backtick extraction, root-filename handling, prose boundaries, loud malformed-input behavior, fixtures, and mutations.
- [x] **[APPLY]:** Added portable paired-field extraction to both shipped checks and mutation-sensitive REQ-344/symmetric/root/prose/malformed fixtures in the accepted three-file scope.
- [x] **[UNIFY]:** Reviewed all three files; shell syntax, RED/GREEN fixtures, four independent mutation axes, contract regressions, canonical maintainer verification, diff checks, and artifact checks passed.

## Why

Found on REQ-344, which touched nine files against a declared write set of two and documented all
seven extensions honestly. Its Implementation Summary packed six paths onto one bullet, three of them
as bare filenames, and the checker reported **two** undeclared files rather than seven:

```
DRIFT: touched but never declared in ## Scope:
  skills/do-work-board/tools/queue-kanban/frontmatter.go
  skills/do-work/actions/capture-reference.md
```

That REQ's own prose asserts the tool "reports seven files touched but never declared". It does not,
and nobody noticed until review — the under-report reads exactly like a clean-ish result.

The same list is what `actions/work.md` Step 9 validates against
`git diff --name-only <pre>..<merge_hash>`, so the same formatting defeats that check too. A checker
whose coverage depends on how a human wrapped a list is a checker that will keep being quietly
disabled.

## Detailed Requirements

- Every backticked repo-relative path on an Implementation Summary bullet is extracted, not just the
  first.
- A bullet naming a file without a path (`capture.md` rather than
  `skills/do-work/actions/capture.md`) either resolves or is reported as unparseable — silence is the
  failure mode this REQ removes.
- Re-running against REQ-344's archived record reports all seven undeclared files.
- Existing single-path bullets behave exactly as before. The suite's current expectations must not
  move.

## Constraints

- `_dev/primes/prime-shell-commands.md` governs. Read it first.
- The Implementation Summary template (`actions/work-reference.md` § Step 6.25) prescribes one
  repo-relative path per bullet. Fixing the extractor does not license abandoning that convention —
  but the extractor must not depend on it, because it is the only thing standing between a
  hand-formatted list and a disabled check.
- Do not make the checker fail on a bullet that legitimately names no file (a "What was done"
  sentence). Distinguish "no paths here" from "paths I could not parse".

## Red-Green Proof

**RED prompt/case:** Run `skills/do-work/tools/checks/scope-drift.sh` against a REQ whose
Implementation Summary lists six files on one bullet. It reports two undeclared files, not seven.
Reproduced on REQ-344, 2026-08-24.

**GREEN when:** that same REQ reports all seven, a bare filename is either resolved or reported as
unparseable, and every existing single-path REQ produces the same output as before.

**Validation:** Inferred during REQ-344's review.

---
*Source: REQ-344 review finding F5 (UR-068).*

## Triage

**Route: B** — The defect and one-file target are clear, but exploration must choose a multi-backtick extractor, define bare-filename resolution versus loud rejection, and preserve every existing single-path fixture before implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

- Both `scope-drift.sh` path-led bullet parsers truncate at the first backtick pair: the shared section extractor uses one `sed` capture and the Scope bullet branch strips at the first closing backtick. Inline Scope headers already iterate all pairs.
- Use one POSIX-awk `emit_closed_backtick_fields(line)` helper over paired fields 2, 4, … for path-led bullets only. Keep root filenames valid, ignore prose bullets with later code spans, reject unmatched path lists, and preserve inline headers, `do-work/` exclusion, sorting, annotated-header failure, and Route A/missing-summary skips.
- The same first-token primitive exists in `associate-files.sh`, whose comment promises parity with scope drift. Fixing one copy would leave the root cause and parity claim false, so the accepted pre-freeze scope includes that sibling and their existing shared contract fixture.
- Add a REQ-344-shaped seven-drift fixture, symmetric later-token mismatches including root filenames, a matching multi-path case, and a prose-only code-span case. Mutations restore either truncation, require slashes, or drop the path-led anchor; each must fail.

## Scope

**Files I will touch:**

- `skills/do-work/tools/checks/scope-drift.sh`
- `skills/do-work/tools/checks/associate-files.sh`
- `skills/do-work/tools/checks/qualify.sh`
- `_dev/tests/contract-regressions.sh`
- `_dev/tests/prescribed-shell-cases/qualify.sh`

**Acceptance criteria:**

- Every closed backtick pair on a path-led Scope or Implementation Summary bullet is extracted by both scope-drift and associate-files; later paths cannot hide behind a matching first token.
- Root filenames remain valid paths, unmatched path lists fail loudly, and prose-only bullets containing code spans remain ignored.
- A REQ-344-shaped fixture reports all seven undeclared files while all existing single-path, annotated-header, and Route A behaviors remain unchanged.
- Independent mutations of either parser, slash filtering, or the path-led anchor fail focused contract probes; the canonical gate passes.
- Every later path on a qualify Summary bullet reaches existence, merge-range diff, and wiring checks; root filenames remain valid and malformed path-led bullets fail loudly in prescribed qualify cases.

## Decisions

- **D-01 — Accept the three-file pre-freeze scope.** The test file is required for mutation-sensitive closure, and `associate-files.sh` carries the identical primitive plus an explicit parity promise. Updating only `scope-drift.sh` would preserve the same silent truncation in a shipped sibling.
- **D-02 — Path-led bullets are the syntax boundary.** Every closed backtick pair on a path-led bullet is a path, including root filenames; later code spans on prose-only bullets are not paths. Unmatched backticks fail instead of returning a partial set.
- **D-03 — Extend after review to the remaining qualification consumer.** Independent review found `qualify.sh` retained the same first-token primitive in Step 6.3's per-file checks. Its script and prescribed case file join the scope so one fix closes the root cause across scope drift, association, and qualification.

## Scope Extensions

- **Review extension:** Added `skills/do-work/tools/checks/qualify.sh` and `_dev/tests/prescribed-shell-cases/qualify.sh` after the first review found the same first-token parser in the shipped merge-range qualification path. This is the same root cause and acceptance surface, not adjacent cleanup.

## Implementation Summary

- `skills/do-work/tools/checks/scope-drift.sh` (modified): extracts every closed backtick pair from path-led Scope and Implementation Summary bullets, preserves root filenames/prose boundaries, and fails loudly on malformed path lists.
- `skills/do-work/tools/checks/associate-files.sh` (modified): uses the same multi-path Summary primitive, associates every listed path, and exits with PARSE-FAILED on unmatched backticks.
- `_dev/tests/contract-regressions.sh` (modified): adds REQ-344 seven-drift fidelity, symmetric later-token mismatches, matching/root/prose/malformed compatibility, associate parity, and four mutation axes.

## Discovered Tasks

None.

## Testing

- TDD RED produced six intended fixture failures across multi-path association, REQ-344 seven-path fidelity, symmetric later-token mismatch, and unmatched Scope/Summary lists.
- GREEN contract regressions passed after implementation. Independent mutations restoring Scope or Summary truncation, requiring slashes, or removing the path-led boundary each failed the intended fixtures.
- Shell syntax, diff checks, builder canonical maintainer verification, and clean-worktree/artifact review passed.
