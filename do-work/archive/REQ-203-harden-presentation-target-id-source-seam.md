---
id: REQ-203
title: Harden presentation target-ID source-seam tests
status: completed
completed_at: 2026-08-17T18:27:39Z
commit: dea5f7a
claimed_at: 2026-08-17T18:12:52Z
status_changed_at: 2026-08-17T18:10:49Z
domain: general
created_at: 2026-08-15T19:20:09Z
user_request: UR-042
addendum_to: REQ-197
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: presentation-target-id-source-seam
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
maintenance: true
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-17T18:14:30Z
  basis:
    - Route B
    - 1-file write set
    - 3 acceptance criteria
    - cross-route regression gates
write_set:
  - _dev/tests/contract-regressions.sh
---

# Review Fix: Harden Presentation Target-ID Source-Seam Tests

## What

Make the completed-work presentation ID contracts prove active, ordered inheritance from the canonical Target ID Resolution source without copying its grammar into callers. Close the entire mutation class, not only the examples found in one review.

## Context

REQ-197's product instructions now cite and apply the right shared contract, but its single remediation left the regression guard able to accept “read without applying,” an omitted pre-dispatch order, and copied membership grammar.

## Instances

- [ ] Shared presentation resolver: reject semantic negations that retain an `apply` substring and keep application before target lookup.
- [ ] `present-work` item dispatch: require shared-contract application before the item branch and reject copied membership grammar as well as copied token examples.
- [ ] Regression block: use a replayable mutation matrix with safe positive controls rather than keyword-presence alone.

## Requirements

- Require a word-bounded, affirmative application directive rather than a matching substring.
- Reject semantic negations including “without applying” for both callers.
- Enforce Target ID Resolution before each caller's lookup or item-dispatch boundary.
- Reject caller-local copies of token and UR-membership grammar.
- Preserve the current correct caller instructions and canonical source contract.

## Red-Green Proof

**RED prompt/case:** Mutate each caller's active directive to “read without applying Target ID Resolution,” move the `present-work` citation below its item branch, and add a caller-local membership definition; the current assertions accept those invalid states.
**Why RED now:** A caller can silently stop applying or fork the canonical grammar while the suite remains green.
**GREEN when:** Every semantic-negation, ordering, and copied-membership mutation fails; the unmodified callers and canonical source pass; and the focused and canonical suites remain green.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] The current product instructions are correct, but their regression test still misses one semantic mutation family. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Triage

**Route:** B — Explore then Build

**Reasoning:** The outcome is precisely stated (close the mutation class in the target-ID source-seam assertions) but the exact assertion wording, the existing helper vocabulary, and the current caller text all had to be discovered before a word-bounded, negation-proof matrix could be written against them.

**Confidence:** high

*Triaged by work action*

## Plan

Planning not required — Route B. Exploration output is sufficient to scope a single-file test change.

## Exploration

**Key files:**
- `_dev/tests/contract-regressions.sh` lines 688–828 — the target-ID source-seam block: canonical-source predicates, then per-caller `require`/`reject`/`require_order` assertions for the shared reference and `present-work`.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` line 20 — the shared resolver's active directive ("read and apply … → **Target ID Resolution**"), sitting above "Resolve exactly one target".
- `skills/do-work-toolbox/actions/present-work.md` line 37 — `present-work`'s active directive, sitting above the item-dispatch bullet.

**Existing pattern to follow:** the file already contains a replayable mutation matrix at lines 496–606 (`unsafe_executable_video_findings` + `unsafe_video_mutations`/`safe_video_mutations`, from REQ-202): a pure detector over a source string, a table of mutations each asserted to be caught by an expected family, and a table of safe controls asserted to stay clean. Reusing that shape keeps the new assertions replayable instead of adding a second idiom.

**Concerns found:**
- `require(..., r"apply[^\n]*Target ID Resolution")` matches the substring inside "applying", so "read **without applying** Target ID Resolution" satisfies the current affirmative-application assertion for both callers.
- `present-work` has no ordering assertion at all, so its directive can be moved below the item-dispatch branch while the suite stays green.
- The copied-grammar rejects for `present-work` cover token spellings but not UR-membership grammar (`user_request:` / `requests:` array), which the shared reference already rejects.

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` (modified)

**Acceptance criteria (restated from the REQ):**
1. The affirmative application directive is matched word-bounded, not by substring.
2. Semantic negations — including "without applying" — are rejected for both callers.
3. Target ID Resolution application is enforced *before* each caller's lookup / item-dispatch boundary.
4. Caller-local copies of token grammar **and** UR-membership grammar are rejected for both callers.
5. The current caller instructions and the canonical source contract still pass unmodified.
6. The block is a replayable mutation matrix with safe positive controls, not keyword presence alone.

## Pre-Flight

- Working tree clean outside `do-work/`.
- Test baseline passing: `bash _dev/tests/contract-regressions.sh`.

## Implementation Summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Replaced the two per-caller keyword-presence blocks (shared presentation resolver, `present-work` item dispatch) with one `target_id_seam_findings` detector plus a replayable mutation matrix, following the `unsafe_executable_video_findings` shape already in the file.

The detector returns defect *families* rather than pass/fail per pattern:

| Family | What changed |
|---|---|
| `missing named inheritance` | unchanged rule, now shared by both callers |
| `missing active application directive` | matched with `(?<![\w-])appl(?:y\|ies)(?![\w-])`, so the letters inside "applying" no longer satisfy it |
| `semantic negation of the inherited grammar` | now covers `without` / `instead of` / `rather than` alongside the original `do not` / `never` / `must not` / `cannot`, plus the `-ing` inflection |
| `copied token grammar` / `copied UR-membership grammar` | split into two families, and the membership family now applies to `present-work` too (previously only the shared reference rejected it) |
| `application ordered after the resolution boundary` | ordering is now enforced for `present-work` against its item-dispatch bullet; previously only the shared reference had an ordering assertion |

Each caller is then run three ways: the shipped text as a positive control, seven defect mutations each asserted to raise its expected family, and four safe mutations asserted to stay clean. A mutation that no longer changes the shipped text fails loudly rather than passing vacuously, and a caller that fails its own positive control skips the replay so the output names the real defect instead of derived noise.

**Tests touched:** `_dev/tests/contract-regressions.sh` (the assertions themselves are the deliverable).

## Qualification

Passed — 1 file verified, 6 acceptance criteria traced, diff contains no debug artifacts.

- Files exist and show in the diff (`_dev/tests/contract-regressions.sh`, +226/−49).
- Substantive: the change is 100+ lines of new detector/matrix logic replacing 55 lines of keyword assertions, not a whitespace shuffle.
- Requirements traced: word-bounded application (criterion 1), negation families (2), ordering for both callers (3), token + membership grammar rejects for both callers (4), unmodified callers pass (5, proven by the green run), mutation matrix with safe controls (6).
- Flowing: the detector is invoked for both real callers and every matrix row; nothing is defined and left unused.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh` (baseline, RED replay, GREEN); `bash _dev/tests/maintainer-verify.sh`

**Result:** ✓ `contract-regressions.sh` exit 0; ✓ `maintainer-verify.sh` exit 0 with zero FAIL lines.

**Red-green validation:** ✗ RED — with the REQ's three mutations applied to the real callers ("read without applying …" in the shared reference; `present-work`'s directive moved below its item branch plus a caller-local `user_request:` membership sentence), the pre-change suite exited **0** and printed `Contract regression checks passed.` → ✓ GREEN — with the same three mutations re-applied after the change, the suite exits **1** and names each defect: `shared presentation resolver — semantic negation of the inherited grammar`, `shared presentation resolver — missing active application directive`, `present-work item dispatch — copied UR-membership grammar`, `present-work item dispatch — application ordered after the resolution boundary`. Callers reverted; suite green again.

This traces directly to `## Red-Green Proof` — the captured RED prompt was replayed verbatim against the real caller files rather than a nearby equivalent.

**Existing tests updated:** the two per-caller assertion blocks this REQ replaces.

*Verified by work action*

## Review

**Overall: 97%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All six acceptance criteria delivered and individually replayed |
| Code Quality | 92% | Reuses the file's existing detector/matrix idiom rather than inventing a second one |
| Test Adequacy | 95% | RED and GREEN both replayed against the real caller files, not a stand-in |
| Scope | 100% | One declared file touched; no product text changed |
| Risk | None | Test-only change; the canonical contract and both callers are byte-identical |
| Acceptance | Pass | Suite green on shipped callers; suite red on every mutated caller |

**Verdict: Approve** — the mutation class is closed, and the closure is proven by replay rather than asserted.

### Findings

**Minor:**
- The `copied token grammar` / `copied UR-membership grammar` families scan the whole caller file, so a future unrelated sentence containing "case-insensitive" would fail with a source-seam message. This is inherited from the `reject()` calls being replaced — narrowing it would need a section anchor per caller, which neither caller currently offers. Not worth the machinery today.

### Requirements Checklist

- [x] Word-bounded affirmative application directive — delivered (`(?<![\w-])appl(?:y|ies)(?![\w-])`)
- [x] Semantic negations including "without applying" rejected for both callers — delivered
- [x] Target ID Resolution enforced before each caller's lookup / item-dispatch boundary — delivered (`present-work` gained the ordering rule it never had)
- [x] Caller-local token **and** UR-membership grammar rejected for both callers — delivered
- [x] Current correct caller instructions and canonical source contract preserved — delivered (no product file in the diff)

### Acceptance Testing

**Result: Pass**
- `bash _dev/tests/contract-regressions.sh` — exit 0 against the shipped tree.
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines.
- Finding-Closure Ratchet: the captured RED prompt was replayed verbatim against the real callers; it passed the old suite (exit 0) and fails the new one (exit 1) naming each defect family. Closure proven, not claimed.

### Restatement Sweep

Not triggered — the diff redefines nothing. The canonical `### Target ID Resolution` contract, both caller directives, and every doc that glosses them are untouched; only the assertions that police them changed.

### Follow-up REQs Created

None.

## Lessons Learned

**What worked:** Replaying the captured RED prompt against the *real* caller files before and after the change — it proved in one run that the old assertions were vacuous and the new ones are load-bearing. Reusing the file's existing `unsafe_executable_video_findings` detector/matrix shape meant no new idiom to learn.

**What didn't:** The first pass let the mutation matrix run even when the shipped caller failed its own positive control, which buried the one real finding under ten lines of derived "safe mutation rejected" noise. Skipping the replay when the positive control fails is what makes the output readable.

**Worth knowing:** A substring assertion on a verb is not an assertion about meaning — `apply[^\n]*X` is satisfied by "without applying X". Any instruction-contract test that checks a directive is *active* needs word boundaries plus an explicit negation family, and the negation family must include the participle forms (`without applying`, `instead of applying`), not just the modal ones (`do not apply`). The `qualify.sh` / `scope-drift.sh` pair parses every `-` bullet under `## Implementation Summary` as a file path, so explanatory sub-lists there must be a table, not a bullet list, or both scripts report phantom drift.

## Orientation

The completed-work presentation contract tests can no longer be satisfied by a caller that merely mentions Target ID Resolution: the source-seam assertions now live in one detector that both callers are replayed through, so a directive that is negated, demoted below its own lookup, or replaced by a local copy of the grammar fails the suite. Lives in the presentation-contract regression block of `_dev/tests/contract-regressions.sh`; no product instruction changed and the system's shape is unchanged.

Prime staleness spot-check: `_dev/primes/prime-action-files.md` — every referenced path still resolves; not made stale by this change.
