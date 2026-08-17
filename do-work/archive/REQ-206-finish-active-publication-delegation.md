---
id: REQ-206
title: Finish active publication delegation
status: completed
completed_at: 2026-08-17T18:55:01Z
commit:
claimed_at: 2026-08-17T18:49:31Z
status_changed_at: 2026-08-17T18:10:49Z
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-17T18:50:10Z
  basis:
    - Route B
    - 2-file write set
    - 5 acceptance criteria
    - cross-route regression gates
write_set:
  - skills/do-work-toolbox/actions/present-video.md
  - _dev/tests/contract-regressions.sh
domain: frontend
created_at: 2026-08-15T20:15:07Z
user_request: UR-042
addendum_to: REQ-201
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: completed-work-publication-active-delegation
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
maintenance: true
---

# Review Fix: Finish Active Publication Delegation

## What

Complete the shared-publication boundary by removing the final paraphrased whole-artifact algorithm from `present-video` and making regression tests require an active application directive at each consumer's execution step.

## Context

REQ-201 centralized the generic mechanics, but review showed that one local “one final path for every source file” sentence remains and that a passive checklist mention can satisfy the current pointer test.

## Instances

- [ ] `present-video` Step 5: retain preferred naming and resolved-path reporting, but delete the one-path/every-file algorithm paraphrase.
- [ ] Presentation contract tests: require an active `apply ... Collision-Safe Publication` directive at the publication step and reject the paraphrased algorithm.

## Requirements

- Keep only the preferred video directory and consumer-specific result reporting local.
- Require an affirmative shared-section application directive before output creation in both live consumers.
- A checklist or Rules mention alone must not satisfy delegation.
- Reject local wording equivalent to “one final path for every file.”
- Preserve the canonical section, output behavior, and all existing presentation contracts.

## Red-Green Proof

**RED prompt/case:** Remove the active Step 5 publication directive while leaving its checklist heading, and insert `one final path selected by that contract for every source file`; the current tests remain green.
**Why RED now:** Passive cross-reference presence and an uncaught paraphrase allow consumer-local mechanics to regrow while the canonicalization suite reports success.
**GREEN when:** Both mutations fail, unmodified consumers pass, and only preferred path/content/result concerns remain outside the shared publication section.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] Publication mechanics are centralized, but one paraphrase and one passive-pointer test escape remain. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Triage

**Route:** B — Explore then Build

**Reasoning:** The deletion is a single sentence, but "an active application directive at each consumer's execution step" needed the exact shape of both consumers' publication steps before an ordering assertion could be anchored to something real.

**Confidence:** high

*Triaged by work action*

## Plan

Planning not required — Route B.

## Exploration

**Key files:**
- `skills/do-work-toolbox/actions/present-video.md` line 106 — Step 5's directive, followed by the paraphrase sentence this REQ deletes.
- `skills/do-work-toolbox/actions/ai-report.md` line 51 — the other live consumer's directive, already correct: preferred path, active `apply`, "before creating it", then the output-creation sentence.
- `_dev/tests/contract-regressions.sh` lines 1020–1033 — the delegation block: a bare `require(consumer, r"Collision-Safe Publication", …)` plus five local-restatement rejects.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md` line 67 — the canonical `## Collision-Safe Publication` section, untouched by this REQ.

**Concerns found:**
- The bare `require` matches anywhere in the file. Both consumers mention the section in their Verification Checklists, so deleting the *active* Step 5 directive entirely still left the suite green.
- No ordering assertion existed for either consumer, so a directive could sit after the step that creates output.
- The five local-restatement rejects covered suffixing and immutability wording but not the whole-artifact algorithm phrasing (`one final path … for every source file`), which is what actually survived in `present-video`.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/present-video.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**Acceptance criteria (restated from the REQ):**
1. Only the preferred video directory and consumer-specific result reporting stay local.
2. An affirmative shared-section application directive precedes output creation in both live consumers.
3. A checklist or Rules mention alone must not satisfy delegation.
4. Local wording equivalent to "one final path for every file" is rejected.
5. The canonical section, output behavior, and all existing presentation contracts are preserved.

## Pre-Flight

- Working tree clean outside `do-work/`.
- Baseline `bash _dev/tests/contract-regressions.sh` passing before the change.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/present-video.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Deleted the last locally restated publication algorithm and replaced the delegation assertion with one that can tell an active directive from a mention.

*Product.* Step 5's second sentence — "Use the one final path selected by that contract for every source file and for the result report" — is gone, replaced by "Report the path that contract resolved." The preferred directory and the consumer-specific result reporting stay; the whole-artifact algorithm goes back to living only in the canonical section.

*Tests.* The bare `require(consumer, "Collision-Safe Publication")` is replaced by a `publication_delegation_findings` detector applied to both live consumers, returning three families: **missing active application directive** (word-bounded `appl(y|ies)` naming the section, so a passive mention no longer counts), **local publication algorithm restated** (the five previous rejects plus the two whole-artifact phrasings), and **application ordered after output creation** (anchored per consumer — `Create \`screenshots/\` only when` for `ai-report`, `Write only the source tree from Step 4` for `present-video`). Each consumer is then replayed through the REQ's two named mutations, each asserted to raise its expected family, with the same shipped-caller positive control and vacuous-mutation guard used elsewhere in the file.

**Tests touched:** the delegation block in `_dev/tests/contract-regressions.sh` — no case-count string to update, this suite reports pass/fail rather than a named-case tally.

## Qualification

Passed — 2 files verified, 5 requirements traced, no debug artifacts.

- Both declared files appear in the diff; nothing undeclared was touched.
- Substantive: one product sentence deleted, ~70 lines of assertions replacing 14.
- Requirements traced: paraphrase removed and preferred path/result reporting retained (1); word-bounded directive plus ordering anchor for both consumers (2, 3); both whole-artifact phrasings rejected (4); canonical section untouched and every pre-existing presentation predicate still green (5).
- Flowing: the detector is called for both consumers and for every mutation row.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh` (baseline, two-stage RED, GREEN); `bash _dev/tests/maintainer-verify.sh`

**Result:** ✓ contract regressions exit 0; ✓ maintainer-verify exit 0, zero FAIL lines.

**Red-green validation:** the captured RED prompt was replayed in two stages against the real consumer.

1. ✗ **The old assertion was vacuous.** With `present-video`'s active Step 5 directive replaced by "follow the publication rules" (its Verification Checklist mention left standing) *and* `Use the one final path selected by that contract for every source file` inserted, the **pre-change** suite exits **0** and prints `Contract regression checks passed.`
2. ✓ **The new assertion catches it.** The same mutated consumer against the **post-change** suite exits **1**: `present-video publication delegation — missing active application directive` and `present-video publication delegation — local publication algorithm restated`.
3. ✓ **The shipped consumers pass.** Reverted, both `ai-report` and `present-video` are clean under all three families, and their two-mutation replays confirm the detector is load-bearing for each.

**Existing tests updated:** the five local-restatement rejects are folded into the detector's `local publication algorithm restated` family — same patterns, same intent, now replayed rather than only asserted.

*Verified by work action*

## Review

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Five criteria, each traceable to a specific assertion or deletion |
| Code Quality | 95% | Third detector in this file to follow the same shape; the idiom is now established rather than invented |
| Test Adequacy | 95% | Two-stage RED proves the old assertion vacuous and the new one load-bearing |
| Scope | 100% | One sentence deleted, one test block replaced; the canonical section untouched |
| Risk | None | Instruction-and-test change; no executable behavior altered |
| Acceptance | Pass | Mutated consumer green on the old suite, red on the new one, shipped consumers clean |

**Verdict: Approve** — publication mechanics now live in exactly one place, and the test can tell delegation from a mention.

### Findings

**Minor:**
- The ordering anchors are literal sentences from each consumer (`Create \`screenshots/\` only when`, `Write only the source tree from Step 4`). Reword either sentence and the assertion reports `missing output-creation step` rather than silently passing — fail-closed, but it does couple the test to prose. The vacuous-mutation guard has the same property, and both are deliberate: a stale anchor should be loud.

**Nit:**
- `publication_local_algorithm` now mixes this REQ's two whole-artifact phrasings with REQ-201's five suffix/immutability patterns in one tuple. They are one family behaviorally; splitting them would report a finer family name for no practical gain.

### Restatement Sweep

**Triggered** — the diff changes what "delegating to Collision-Safe Publication" requires of a consumer. Swept `Collision-Safe Publication` across `skills/` and `_dev/`. Results: the canonical section in `completed-work-presentation-reference.md` is unchanged and still owns the mechanics; `ai-report.md`'s two passive mentions (Step 8 verification, Verification Checklist) are legitimate references and now sit alongside its already-correct active directive; `present-video.md`'s checklist mention is likewise fine and was deliberately left standing — it is the exact bait the RED replay uses. No stale restatement remained after the deletion.

### Requirements Checklist

- [x] Only preferred directory and consumer-specific result reporting stay local — delivered
- [x] Affirmative application directive before output creation in both live consumers — delivered
- [x] A checklist or Rules mention alone no longer satisfies delegation — delivered (this is exactly what RED stage 1 proved was broken)
- [x] Local "one final path for every file" wording rejected — delivered (both word orders)
- [x] Canonical section, output behavior, and existing presentation contracts preserved — delivered

### Acceptance Testing

**Result: Pass**
- `bash _dev/tests/contract-regressions.sh` — exit 0.
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines.
- Finding-Closure Ratchet: the captured RED prompt is replayed verbatim — active directive replaced while the checklist heading stays, plus the paraphrase inserted. It passes the old suite and fails the new one by name. Closure measured, not asserted.

### Suggested Additional Testing

- If a third consumer of the shared publication section is ever added, it must be appended to `publication_consumers` with its own output-creation anchor; nothing currently detects a consumer that forgets to enroll. A future REQ could derive the consumer list from files that name the section rather than hand-listing it — the Closed Enumerations Go Stale rule applies.

### Follow-up REQs Created

None.

## Lessons Learned

**What worked:** Running the RED in two stages — mutated consumer against the *old* suite, then against the *new* one. One run proves the assertion was vacuous, the other proves the replacement is load-bearing; either alone would have been weaker evidence. Leaving the Verification Checklist mention in place during the replay is what makes it a faithful reproduction rather than a strawman.

**What didn't:** `git checkout --` on the mutated consumer reverted the REQ's own product fix along with the mutations, because both lived in the same file. Reverting a mutation applied on top of uncommitted work needs the mutation undone specifically, not the file restored.

**Worth knowing:** A `require(file, token)` assertion over a whole document tests vocabulary, not behavior — any file that mentions a contract anywhere satisfies it, including in the checklist that merely claims the contract was followed. When the property is "this step actively does X", the assertion needs a word-bounded verb and an anchor to the step, which is the same conclusion REQ-203 reached from the other direction on the same day.

## Orientation

`present-video` no longer carries its own copy of the publication algorithm, and the contract tests can now tell an active delegation from a passing mention: removing a consumer's execution-step directive fails the suite even when its Verification Checklist still names the contract. Lives in the completed-work presentation family (`skills/do-work-toolbox/actions/present-video.md` plus the delegation block of the contract suite). The canonical publication section is unchanged, so the system's shape is unchanged — this closes the last local restatement rather than moving a boundary.

Prime staleness spot-check: `_dev/primes/prime-action-files.md` — all referenced paths resolve; not made stale by this change.
