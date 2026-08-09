---
id: REQ-157
title: "Review fix: Complete the retired core alias guard"
status: completed
completed_at: 2026-08-09T20:38:50Z
status_changed_at: 2026-08-09T18:01:17Z
claimed_at: 2026-08-09T20:05:21Z
route: C
domain: general
kb_status: pending
kb_entry:
created_at: 2026-08-08T20:13:29Z
user_request: UR-031
addendum_to: REQ-153
review_generated: true
effort_estimate: normal
write_set:
  - _dev/tests/fixtures/retired-core-moved-command-triggers.tsv
  - _dev/tests/staged-skills-contract.sh
sweep: true
sweep_key: retired-core-alias-guard-completeness
---

# Review Fix: Complete the Retired Core Alias Guard

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Confirmed the two-file scope and historical sources, reconciled the 117 direct heads plus 22 three-form install targets, and reproduced the old partial guard missing all four required representative aliases.
- [x] **[APPLY]:** Root cause: the REQ-153 guard reconstructed moved commands from current canonical sibling action names plus seven samples instead of the deleted router's exhaustive vocabulary. Added the test-only 117-direct/22-target historical inventory, compiled longest-first exact-boundary matching from it, and proved every row plus exclusions through root/module mutations; the focused contract is GREEN.
- [x] **[UNIFY]:** Audited `git diff --stat` and every changed project file, found no debug artifacts, and confirmed only the declared test script plus new fixture changed. Qualification attempt 1 narrowed the branding exception to the exact product-title identity and added an em-dash command-prose positive control. Passed `bash _dev/tests/staged-skills-contract.sh`, `bash _dev/tests/contract-regressions.sh`, `bash -n _dev/tests/staged-skills-contract.sh`, `shellcheck -S warning _dev/tests/staged-skills-contract.sh`, the representative PCRE2 live-root/module sweep, the named-history inventory cross-check, REQ-153 live-repair preservation, and `git diff --check`.

## What
Make the shipped-surface recurrence guard cover every former moved-command trigger, not only canonical sibling action names and a small sample of natural-language aliases.

## Context
Found during review of REQ-153. All current stale occurrences were repaired, but the guard still permits many equivalent retired core invocations, so the same contract can recur under an alias without failing distribution tests.

## Requirements
- Recover the complete former moved-command trigger set from the deleted core router/shim history and represent it as one auditable contract source.
- Reject every retired core trigger on live root/module surfaces with exact command boundaries.
- Keep the trigger list test-only or historical; do not republish legacy aliases as user guidance or runtime routes.
- Add negative controls for branding/noun phrases, generic pipeline prose, historical changelogs/archives, and explicit negative fixtures.
- Preserve current unique sibling ownership/routes, prime transition fingerprints, all 15 repaired live surfaces, and full distribution tests.

## Instances
- [ ] Guard `do-work kanban` and every former board trigger, including natural-language board phrases.
- [ ] Guard direct knowledge aliases such as `do-work recall` and every former memory/dream/interview/prompt/setup trigger.
- [ ] Guard toolbox aliases such as `do-work code review` and `do-work describe changes`, including former install targets.
- [ ] Prove the complete former trigger set is covered by a table-driven mutation fixture.

## Open Questions

- [x] The live retired commands are repaired, but the new recurrence test covers only part of the former alias set, so an equivalent old command could return unnoticed. The cascade-depth rule requires your consent before automatically working a follow-up created by the review of another review-generated task. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---

## Triage

**Route: C** - Complex

**Reasoning:** This is a cross-package regression-contract change that must reconstruct a complete historical command vocabulary, distinguish executable command boundaries from ordinary prose, and prove completeness with mutation fixtures while preserving the permanent modular runtime contract.

**Planning:** Required

## Plan

1. Add `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` as the single test-only historical contract. Reconstruct it route-by-route from the last live modular router at `0b9bcde^:skills/do-work/SKILL.md`, cross-check category/target ownership against its deleted `actions/moved-command-shim.md`, and expand install-target spellings from the pre-cutover `a3c2612^:SKILL.md` install row. Record owner, canonical sibling action, match kind, and legacy trigger; include all board natural-language aliases, every knowledge memory/dream/interview/prompt/setup alias, every toolbox alias, and former install target/form families. Keep this inventory exclusively under `_dev/tests`; do not change any router, action, help, or shipped guidance.
2. Refactor only the retired-command block in `_dev/tests/staged-skills-contract.sh` to load that fixture instead of deriving canonical actions plus seven samples. First add a table-driven mutation check around the existing matcher and record RED failures for representative uncovered forms (`do-work kanban`, `do-work recall`, `do-work code review`, `do-work describe changes`); then compile the production scan matcher from the complete fixture for GREEN. Preserve the existing live scope (`justfile` plus module files), ownership/route assertions, updater-prime fingerprints, and historical exclusions. Match command heads with explicit left/right boundaries so arguments and punctuation are allowed, while embedded/prefixed words, alias suffixes, possessive noun phrases, and current `do-work-board`, `do-work-knowledge`, and `do-work-toolbox` invocations do not match.
3. In the same Python fixture harness, mutate every contract row into root and module command examples and require exactly one match, including all install families. Add negative controls for `Do-Work Board Skill`, `do-work board's testing view`, ordinary board/knowledge/toolbox nouns, generic work/CI/data pipeline prose, valid core and sibling commands, and near-boundary strings such as `undo-work …`, `<alias>-suffix`, `<alias>suffix`, and possessives. Build synthetic `CHANGELOG.md`, `do-work/archive/...`, and `_dev/tests/fixtures/...` files containing retired literals and prove the live-surface collector ignores them. This covers every Requirements and Instances item without republishing runtime aliases.
4. Run `bash _dev/tests/staged-skills-contract.sh` for GREEN, then `bash _dev/tests/contract-regressions.sh`, `bash -n _dev/tests/staged-skills-contract.sh`, warning-level ShellCheck, a targeted live-root/module sweep, and `git diff --check`. Confirm the diff contains only the two test files, all 15 REQ-153 repairs remain byte-untouched and clean, unique sibling ownership/routes and prime-transition checks still pass, and no runtime alias or production route was restored.

**Plan validation:** Tasks 1–3 map the complete historical trigger inventory, exact command boundaries, all four Instances bullets, mutation completeness, and every required negative control. Task 4 proves preservation of REQ-153's repaired surfaces, permanent modular ownership, prime fingerprints, and the full distribution contract. No requirements are uncovered, no task is orphaned, and the four-task plan remains below the five-task scope warning threshold.

*Generated by Plan agent*

## Exploration

- The authoritative vocabulary is recoverable from three historical sources: `0b9bcde^:skills/do-work/SKILL.md` for the last live modular router, its deleted `actions/moved-command-shim.md` for owner/canonical destinations, and `a3c2612^:SKILL.md` for install normalization. Commit `157b89e` introduced the current partial guard.
- The explicit non-install inventory contains 117 command heads (6 board, 19 knowledge, 92 toolbox). Install contributes 22 target aliases expanded across `install <target>`, `install-<target>`, and `setup <target>` families.
- The current embedded Python guard derives canonical sibling actions, adds only seven aliases, scans the root `justfile` plus live `skills/` files, excludes changelogs and the built board binary, and separately preserves updater-prime fingerprints and unique ownership checks.
- The two-file plan is sufficient. Longest-first matching is required for overlapping heads; mutation assertions should identify the expected fixture row; apostrophes and hyphens must remain rejected right boundaries so possessives and suffixes stay negative; valid `do-work review code` must remain distinct from retired `do-work code review`.
- The fixture should validate its own header, field completeness, duplicate triggers, owners/actions, and match kinds. Representative RED cases are `kanban`, `recall`, `code review`, and `describe changes`; GREEN mutates every contract row and checks live-surface exclusions.

*Generated by Explore agent*

## Decisions

- **D-01 — Include bare historical `install`, `setup`, and `install-` heads in the test-only inventory.** This is a reversible contract-fixture choice directly implied by “every former trigger”: the last modular router explicitly accepted those heads, so omitting them would knowingly preserve an incomplete historical vocabulary. They remain test data only and create no runtime route or guidance.
- **D-02 — Preserve the existing `do-work queue board` product-title noun phrase.** The expanded inventory correctly exposed the phrase in the board HTML title and the two Go assertions that verify that title, but those three occurrences are branding rather than retired invocations. The guard exempts only the exact `PROJECT_NAME — do-work queue board` title identity and its two exact test references; other em-dash prose such as `Deprecated — do-work kanban` remains a retired command violation.

## Scope

**Files I will touch:**
- `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` (new) — auditable test-only historical trigger inventory.
- `_dev/tests/staged-skills-contract.sh` (modify) — fixture loader, exact-boundary live guard, table-driven mutation proof, and negative controls.

**Files I will NOT touch:** Runtime routers, actions, help/guidance, the 15 REQ-153 repaired live surfaces, updater-prime content, changelog/archive/ADR history, ADR-019, UR-031, or pending-answer REQs 158/159.

**Acceptance criteria (restated from REQ):**
- [ ] One auditable test-only source recovers every former board, knowledge, toolbox, and install trigger, including natural-language aliases and bare historical install/setup heads.
- [ ] Every retired core trigger is rejected on live root/module surfaces at exact command boundaries, without restoring runtime aliases.
- [ ] Branding/noun phrases, generic pipeline prose, historical changelogs/archives, explicit negative fixtures, current core/sibling commands, possessives, and near-boundary strings remain accepted.
- [ ] A table-driven mutation fixture proves every inventory row is recognized exactly once and validates the fixture contract itself.
- [ ] Unique sibling ownership/routes, prime transition fingerprints, all 15 REQ-153 repairs, and full distribution tests remain intact.

## Pre-Flight

**Git:** ⚠ Pre-existing ADR-019 and UR-031 edits are unrelated to REQ-157; preserve and exclude both from staging and review evidence.
**Tests baseline:** ✓ `bash _dev/tests/staged-skills-contract.sh` passes before implementation (`launched: true`).
**Dependencies:** ✓ Bash and Python are available; no external service or network access is required.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` (new) — records the complete test-only historical owner/action/match-kind/trigger inventory reconstructed from the deleted router and shim.
- `_dev/tests/staged-skills-contract.sh` (modified) — validates the inventory, compiles longest-first exact-boundary retired-command matching, mutates every row on root and module surfaces, checks row identity and exclusions, and scans live shipped surfaces.

**What was done:** Replaced the partial canonical-action-plus-samples recurrence guard with a fixture-driven historical contract covering 186 concrete trigger rows: 117 direct heads, 22 install targets across three former families, and three bare install heads. Added fixture-integrity checks, complete table-driven root/module mutations, branding/prose/current-command boundary controls, and historical/archive/fixture collector exclusions without changing runtime routes or shipped guidance.

## Qualification

**Attempt 1:** Mechanical checks passed, the two project files are substantive and fully traced, and P-A-U is complete. Judgment failed the exact-boundary requirement: `find_retired_matches` skipped every command whose `do-work` head was immediately preceded by `— `, so a genuine live restatement such as `Deprecated — do-work kanban` returned no match. The branding exemption must be narrowed to the exact `PROJECT_NAME — do-work queue board` product-title form while preserving the two explicit Go assertion controls.

**Passed on attempt 2** — 2 project files verified, all 5 requirements and 4 instance classes traced, P-A-U confirmed, and the em-dash false negative is now a positive regression control. Independent focused/full contracts, Bash syntax, warning-level ShellCheck, diff hygiene, fixture integrity, and preserved user-file hashes all pass.

## Testing

**Tests run:** `bash _dev/tests/staged-skills-contract.sh`; `bash _dev/tests/contract-regressions.sh`; `bash -n _dev/tests/staged-skills-contract.sh`; `shellcheck -S warning _dev/tests/staged-skills-contract.sh`; representative live-root/module PCRE2 sweep; named-history inventory cross-check; REQ-153 repair-preservation check; `git diff --check`.
**Result:** ✓ All passing.

**Red-green validation:**
- Partial retired-command matcher: ✗ before implementation — recognized 0/4 required representative aliases (`kanban`, `recall`, `code review`, `describe changes`) and exited 1 for the incomplete-vocabulary reason → ✓ after — all 186 fixture rows match exactly once on both synthetic root and module surfaces.
- Em-dash command prose: ✗ qualification attempt 1 — `Deprecated — do-work kanban` was incorrectly exempted → ✓ after remediation — matches exactly `kanban`, while the exact `PROJECT_NAME — do-work queue board` branding identity remains accepted.

**New tests added:**
- Auditable TSV schema, field, duplicate, owner/action, match-kind, direct-count, install-family, ownership, and bare-head integrity checks.
- Complete per-row root/module mutations with expected trigger identity, negative boundary forms, current-command controls, product branding/prose controls, and historical/archive/fixture collector exclusions.

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/staged-skills-contract.sh` (from REQ-153): replaced its partial canonical-action-plus-samples recurrence guard with the complete test-only historical contract while preserving its unique sibling-route and updater-prime assertions.

*Verified by work action*

## Review

**Overall: 73%** | 2026-08-09T20:36:33Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 78% |
| Test Adequacy | 75% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- The two queue-board test-reference exemptions operate on the whole source line, so a second genuine `do-work queue board` occurrence on either exempt line is silently accepted instead of rejected — gate: user-visible → rerouted pending-answers as REQ-160
- The longest-prefix matcher abandons the command after a boundary-invalid specific candidate instead of trying the shorter valid `install`/`setup` head, and treats the historical `install-` prefix as an exact literal; overlapping or unknown-target retired install forms therefore escape the guard — gate: user-visible → rerouted pending-answers as REQ-160

**Minor findings:** 0 (report only)
**Acceptance:** Partial — the complete historical inventory and ordinary positive/negative cases pass, but occurrence-level exemptions and overlapping install heads leave exact-boundary false negatives.
**Suggested testing:** 2 items
**Follow-ups created:** REQ-160; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Reconstructing the contract from the deleted router/shim and install-normalization history produced an auditable complete inventory instead of another sample list.
- Full-row root/module mutations plus a separate qualification pass caught both vocabulary gaps and an over-broad branding exemption before archival.

**What didn't:**
- Line-wide exemptions are too coarse for a command-occurrence guard: one legitimate test reference can hide a second real invocation on the same line.
- Selecting only the longest prefix fails when that candidate misses its boundary but a shorter historical head or prefix route still applies.

**Worth knowing:**
- The 186-row historical vocabulary is complete and remains test-only. REQ-160 records the two remaining occurrence-completeness edge classes and is consent-gated rather than auto-runnable.

## Orientation

[MAP CHANGED] The staged distribution guard now has one test-only historical source for every retired core board, knowledge, toolbox, and install trigger, with fixture integrity and complete root/module mutation coverage. Runtime ownership remains permanently modular; two adversarial occurrence-matching edges are isolated in consent-gated REQ-160.
