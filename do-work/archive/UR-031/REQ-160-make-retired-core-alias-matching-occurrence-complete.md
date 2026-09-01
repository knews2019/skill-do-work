---
id: REQ-160
title: "Review fix: Make retired core alias matching occurrence-complete"
status: completed
completed_at: 2026-08-10T10:41:11Z
commit: 3d8613a
claimed_at: 2026-08-10T10:18:22Z
route: C
status_changed_at: 2026-08-10T09:20:51Z
domain: general
created_at: 2026-08-09T20:36:33Z
user_request: UR-031
addendum_to: REQ-157
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: retired-core-alias-match-occurrence-completeness
write_set:
  - _dev/tests/staged-skills-contract.sh
kb_status: promoted
kb_entry: REQ-160-review-fix-make-retired-core-alias-match.md
---

# Review Fix: Make Retired Core Alias Matching Occurrence-Complete

## What

Make the test-only retired-command guard evaluate command occurrences and overlapping trigger candidates completely. An exemption or boundary-invalid longer candidate must not hide another valid retired command. Done means this false-negative class cannot recur while direct-alias suffix exclusions remain intact.

## Context

Found during review of REQ-157. Its 186-row historical inventory is complete, but the matcher can still miss retired invocations when a line also contains an exempt queue-board test reference or when a longer install-target prefix fails its right boundary before a shorter historical install/setup head is considered.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Loaded `crew-members/general.md`, `crew-members/coding-guardrails.md`, and `specs/bug-fix.md`; used the fresh Plan and Explore findings below to confine the fix to `_dev/tests/staged-skills-contract.sh`. Planned exact occurrence spans for the three approved queue-board exemptions, family-scoped install/setup candidate fallback, `install-`-only historical prefix semantics, and adversarial preservation checks for direct suffixes, possessives, current sibling commands, fixture identity, and live surfaces.
- [x] **[APPLY]:** Added the adversarial controls first and captured RED (exit 1): both exact test fragments matched when alone and returned two matches when followed by a second retired occurrence; `install ui-design2`, `setup ui-design2`, and `install-custom-target` returned no match. Implemented trigger/start/end match tuples, exact exempt-span comparison, fallback limited to `install-space`/`install-prefix`/`setup-space`, and target-token consumption only for the historical `install-` head. The focused staged-skills contract then passed.
- [x] **[UNIFY]:** `bash _dev/tests/staged-skills-contract.sh`, `bash _dev/tests/contract-regressions.sh`, `bash -n _dev/tests/staged-skills-contract.sh`, `shellcheck -S warning _dev/tests/staged-skills-contract.sh`, and `git diff --check -- _dev/tests/staged-skills-contract.sh` all exited 0. Verified 186 fixture rows and 186 unique triggers, byte-identical root/module changelogs, unchanged fixture/ADR-019/UR-031 hashes, representative temporary root/module live-surface probes, and a clean debug-marker/diff/status audit; the only project implementation file changed is the declared staged-skills contract.

## Requirements

- Scope queue-board branding and test-reference exemptions to the exact matched occurrence; a second retired command on the same line must still fail.
- Continue candidate evaluation after a boundary-invalid longer trigger when a shorter historical command head remains valid.
- Treat the former `install-<target>` routing prefix consistently, including unknown targets formerly routed to shim help.
- Preserve negative behavior for embedded/prefixed words, possessives, and direct-alias suffixes.
- Add adversarial regressions for exempt and genuine occurrences sharing one line, overlapping install-target prefixes, unknown `install-` targets, and current sibling commands.
- Keep all aliases test-only and preserve the 186-row inventory, current sibling ownership/routes, prime fingerprints, repaired live surfaces, and full distribution tests.

## Instances

- [ ] A line containing `<title> text "do-work queue board"` plus a second `do-work queue board` invocation must report the second occurrence.
- [ ] A line containing `strings.Contains(bodyText, "do-work queue board")` plus a second retired invocation must report the second occurrence.
- [ ] `do-work install ui-design2` and `do-work setup ui-design2` must fall back to the valid historical install/setup head.
- [ ] `do-work install-custom-target` must be recognized through the former `install-<target>` prefix without turning direct-alias suffixes into positives.

## Open Questions

- [x] REQ-157 completes the historical alias inventory, but its guard still misses some former install-head commands and additional retired commands placed on exempt test-reference lines. Because REQ-157 was itself created by a review, the cascade-depth rule requires your approval before this non-critical follow-up can enter the claimable queue. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---

## Triage

**Route: C** - Complex

**Reasoning:** The implementation is test-only, but it must change an occurrence matcher across exemption spans, overlapping candidates, and a historical prefix route while preserving 186-row inventory identity and carefully tuned negative boundaries. Full planning and exploration reduce false-positive risk.

**Planning:** Required

## Plan

1. Add focused RED controls in `_dev/tests/staged-skills-contract.sh` for both exact queue-board test-reference fragments plus a second genuine occurrence, product-title branding plus a real occurrence, overlapping `install`/`setup` heads, unknown `install-` targets, and current/suffix/possessive negatives.
2. Refactor only the embedded Python matcher to carry occurrence spans, filter exact exempt spans rather than whole lines, continue longest-first candidate evaluation after boundary rejection, and give only the historical `install-` head its former prefix semantics.
3. Preserve all 186 known-row identities and mutations while completing the adversarial matrix for exact/direct boundaries, current sibling commands, and embedded/prefixed/possessive/suffix negatives.
4. Run focused/full distribution, Bash syntax, warning-level ShellCheck, fixture identity/count, live-repair/route/prime preservation, and diff checks.

**Root-cause hypothesis:** The matcher returns trigger text without an occurrence span and the live scan later exempts an entire source line by substring, so one legitimate test reference hides every matching occurrence. Separately, `next(...)` selects the longest textual prefix before boundary validation, so one invalid specific candidate prevents shorter valid install/setup heads; the historical `install-` route head is also treated as an exact literal instead of a prefix.

**Exact implementation scope:** Modify `_dev/tests/staged-skills-contract.sh` only. Keep `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` byte-identical and keep runtime routers, actions, help, repaired live surfaces, ADR-019, UR-031, and owner-only lifecycle/release files out of builder scope.

**Plan validation:** Tasks 1–2 map every Requirement and all four Instances to RED, implementation, and complete positive/negative coverage; task 4 proves the inventory, ownership/routes, fingerprints, repaired surfaces, and distribution preservation clauses. Four tasks; no requirement is uncovered and no task is orphaned.

*Generated by Plan agent*

## Exploration

The defect is confined to `_dev/tests/staged-skills-contract.sh`. `find_retired_matches()` currently returns trigger text only; the live scan then applies the two queue-board test-reference exemptions as whole-line substring checks. Returning `(trigger, occurrence_start, occurrence_end)` is sufficient to compare each match against exact spans for the product title and the two fixed test fragments while leaving a second occurrence visible.

Candidate selection uses `next(...)` before right-boundary validation; a longest-first loop should continue after rejection and accept the first valid shorter head. Only the fixture row `match_kind=install-head, legacy_trigger=install-` needs prefix semantics, so known install rows retain exact identity while `install-custom-target` reports the historical head. The universal suffix-negative loop must distinguish direct/exact rows from that prefix head rather than weakening direct alias boundaries. The TSV already carries all 186 rows and needs no edit.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/staged-skills-contract.sh` (modify) — make matching occurrence-complete and add exact adversarial controls.

**Files I will NOT touch:** `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` (read-only identity evidence), runtime routers/actions/help, repaired live surfaces, sibling routes, updater prime, ADR-019, UR-031 input, other REQs, or owner-only lifecycle/version/changelog files.

**Acceptance criteria (restated from REQ):**
- [ ] Product-title and both queue-board test-reference exemptions suppress only their exact occurrence; a second retired invocation on the same line is reported.
- [ ] Boundary-invalid known install/setup candidates fall through to valid shorter heads, while known exact rows retain identity.
- [ ] Unknown `install-<target>` forms match only through the historical `install-` prefix head.
- [ ] Embedded/prefixed words, possessives, direct-alias suffixes, and current core/sibling commands remain negative.
- [ ] Adversarial controls, all 186 inventory rows, live collection, ownership/routes, prime fingerprints, repaired surfaces, focused/full distribution, syntax/ShellCheck, fixture identity, and diff hygiene pass.

## Pre-Flight

**Git:** ✓ Working tree clean outside `do-work/`; approved queue/claim state is lifecycle metadata.
**Tests baseline:** ✓ `bash _dev/tests/staged-skills-contract.sh` passes (`launched: true`).
**Dependencies:** ✓ Bash and Python 3 are available; the historical TSV is present and remains read-only.

*Checked by work action*

## Root Cause

The test-only matcher discarded source positions and selected the first longest textual trigger before checking its right boundary. As a result, the live scan could exempt only by whole-line substring, and a boundary-invalid specific install/setup trigger prevented a shorter historical head from being evaluated. The historical `install-` head was also checked as an exact literal, even though that former route accepted a target suffix.

## Decisions

- **D-01 — Exact exempt spans:** Represent each retired match as its trigger plus start/end offsets, derive spans for only the three approved queue-board fragments, and suppress a match only when its occurrence span equals one of them.
- **D-02 — Family-scoped fallback:** Continue after a boundary-invalid candidate only for `install-space`, `install-prefix`, and `setup-space` rows. This repairs the historical install/setup overlap without allowing a suffixed longer direct alias to fall through to a shorter direct alias.
- **D-03 — Historical install-prefix token:** Let only the `install-` head consume an alphanumeric/underscore target with optional hyphenated segments before applying the existing right-boundary rule. This recognizes unknown former targets while keeping bare fixture identity, possessives, double-hyphen forms, and direct-alias suffixes negative.

## Implementation Summary

**Files changed:**
- `_dev/tests/staged-skills-contract.sh` (modified) — carries retired-command occurrence spans, filters only exact approved exemption spans, evaluates eligible install-family fallbacks after invalid longer boundaries, applies historical `install-` prefix semantics to unknown targets, and adds adversarial positive/negative controls.

**What was done:** The test-only recurrence guard now evaluates every occurrence and eligible overlapping historical candidate without weakening direct alias boundaries or changing the 186-row inventory/runtime contract.

## Qualification

**Attempt 1:** Failed lifecycle qualification. The implementation, scope manifest, requirement trace, RED/GREEN evidence, fixture identity, protected-file hashes, and owner verification all passed, but the mandatory P-A-U execution-state section was absent from the working REQ.

**Attempt 2:** Passed. The canonical P-A-U section contains exactly one checked PLAN, APPLY, and UNIFY entry; the project diff is still limited to the one declared implementation file. All six Requirements and four Instances trace to the implementation and adversarial controls, the 186-row fixture remains byte-identical, and owner verification passes.

## Testing

**RED:** The aggregate adversarial probe exited 1 before implementation. Both exact queue-board test fragments were reported when alone and each hid no occurrence when paired (two matches were returned instead of the required one); `do-work install ui-design2`, `do-work setup ui-design2`, and `do-work install-custom-target` returned no retired match.

**GREEN / UNIFY:** The focused staged-skills contract and full contract-regression suite pass, including manifest, shipped-reference, commit-hash, updater, staged-skills, and installer probes. `bash -n`, warning-level ShellCheck, suite-manifest validation, changelog mirror comparison, fixture count/uniqueness/hash checks, protected-file hash checks, representative live-surface probes, debug-marker audit, and `git diff --check` all pass.

## Review

**Overall: 98%** | 2026-08-10T10:40:16Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 97% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
None

**Minor findings:** 1 (report only)
**Acceptance:** Pass — occurrence-scoped exemptions, install-family fallback, unknown-prefix handling, exact row identity, negative boundaries, live surfaces, routes, and prime fingerprints all passed focused, full-suite, and independent adversarial verification.
**Suggested testing:** 1 item
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Returning occurrence spans made the exemption rule exact and independently testable without weakening the retired-trigger vocabulary or runtime surface.
- Keeping candidate fallback family-scoped repaired the historical install/setup overlap while preserving direct-alias boundaries across all 186 rows.

**What didn't:**
- The first lifecycle qualification omitted the mandatory P-A-U record even though implementation evidence was complete; the second attempt restored it before review.

**Worth knowing:**
- Longest-first command matching must validate each candidate before committing to it, and exemptions belong to source spans rather than whole lines. Independent review's repeated/different-occurrence matrix is the durable guard against both false-negative classes.

## Orientation

The test-only retired-command guard now evaluates each source occurrence and eligible historical install/setup candidate completely. Exact queue-board branding/test fragments suppress only their own spans, unknown former `install-<target>` commands resolve through the historical prefix head, all 186 inventory identities and direct-alias negatives remain intact, and no runtime alias was restored.
