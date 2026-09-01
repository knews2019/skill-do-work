---
id: REQ-146
title: "Remove Modular-Migration Compatibility Shims"
status: completed
completed_at: 2026-08-08T18:33:10Z
commit: 0b9bcde
claimed_at: 2026-08-08T17:52:07Z
route: C
status_changed_at: 2026-08-07T21:40:04Z
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
kb_status: promoted
kb_entry: REQ-146-remove-modular-migration-compatibility-s.md
prime_files: [skills/do-work/tools/prime-do-work-update.md]
write_set:
  - README.md
  - skills/do-work/SKILL.md
  - skills/do-work/actions/moved-command-shim.md
  - skills/do-work/actions/help.md
  - skills/do-work/actions/version.md
  - skills/do-work/next-steps.md
  - skills/do-work/tools/prime-do-work-update.md
  - skills/do-work/tools/do-work-update.sh
  - tools/install-do-work-suite.sh
  - skills/do-work/tools/install-do-work-suite.sh
  - tools/replace-text-section.sh
  - skills/do-work/tools/replace-text-section.sh
  - skills/do-work-knowledge/actions/setup-memory.md
  - _dev/tests/staged-skills-contract.sh
  - _dev/tests/update-script-behavior.sh
  - _dev/tests/install-suite-behavior.sh
  - _dev/tests/contract-regressions.sh
tdd: true
suggested_spec: refactor
depends_on: [REQ-144, REQ-145]
maintenance: true
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145]
batch: do-work-four-skill-suite
---

# Remove Modular-Migration Compatibility Shims

## What
Remove transitional core routes and migration rules after one modular release and confirmed client migration.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Retire the migration window in four bounded passes: invert the declared shell contracts first and capture assertion-level RED failures; delete core moved-command routing/help surfaces; collapse updater/installer/replacer behavior to the permanent four-module, marker-managed path; then run the focused and full verification matrix. Preserve the installed-validator/installer trust boundary, one-archive bootstrap, dirty/confirmation/semver guards, all-or-recover transaction, current hooks, pipeline-guard cleanup, and byte-identical paired tools. Assumption: the user's recorded migration attestation makes unmarked legacy recipes and old core memory-hook paths unsupported inputs, so those paths must be removed rather than renamed.
- [x] **[APPLY]:** Declared retirement contracts failed on the live compatibility surface, then implementation removed the core shim/routes; collapsed updating to installed-validator plus installed-installer delegation; removed monolith/root/stale-copy handling, capability probing, exact legacy recipe recognition, and old core memory-hook rewriting; and refreshed current help/setup/update guidance. Only declared Scope files changed, paired tools remain byte-identical, and release metadata/current-version fields were not touched.
- [x] **[UNIFY]:** Audited the 17-file implementation diff (`199 insertions`, `859 deletions`) and found no debug artifacts. Verified `README.md` (current install/update contract), `skills/do-work/SKILL.md` (core-only routes), deleted `actions/moved-command-shim.md` (no live compatibility action), `actions/help.md` (core menu plus three sibling pointers), `actions/version.md` (current-suite delegation; version line and REQ-150 URL untouched), `next-steps.md` (direct sibling calls), `tools/prime-do-work-update.md` (current transaction boundary), `tools/do-work-update.sh` (Git/semver guards, installed-helper trust, one archive, current-only delegation), both `install-do-work-suite.sh` copies (marker/current-hook composition, jq/Python pipeline cleanup, confirmation and all-or-recover), both `replace-text-section.sh` copies (marker-only byte-safe reconciliation), knowledge `actions/setup-memory.md` (current hooks only), and all four changed shell contract suites (retirement absence plus permanent bootstrap/update/recovery coverage). The paired installers and replacers are byte-identical. Focused RED failed at assertions, focused GREEN and full regressions passed; suite-manifest, ShellCheck, Bash syntax, Just parsing, Go format/test/vet/build, `git diff --check`, protected-path, and debug scans passed.

## Why
Permanent shims would preserve the same router complexity the modularization is meant to remove.

## Detailed Requirements
- Remove old core routing shims for board, knowledge, and toolbox commands.
- Remove bridge-only manifest compatibility, stale-path deletion rules, and migration messages that no supported client needs.
- Preserve the full-suite bootstrap, current modular updater, managed-section reconciliation, and current-to-current modular updates.
- Ensure core help lists only core commands plus concise pointers to the three sibling skills.
- Prove each active action exists in exactly one skill and old invocations no longer appear as live core routes.

## Constraints
- Do not start until at least one modular release has shipped.
- Require user confirmation that every client has four skill directories, current managed Just paths, and migrated enabled memory hooks.

## Dependencies
Requires REQ-144 and REQ-145. The external stabilization and client-migration gate was cleared by user confirmation via `do-work clarify` on 2026-08-08.

## Builder Guidance
Certainty level: Firm. Delete transitional machinery; do not replace it with permanent forwarding aliases.

## Red-Green Proof
**RED prompt/case:** Inspect core routing/help and updater compatibility branches after the first modular release.
**Why RED now:** One-release shims and bridge migration paths remain intentionally active.
**GREEN when:** Transitional paths are absent, current modular installs still update successfully, and core exposes only its own functionality plus discovery pointers.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the one-release compatibility decision.

## Blocked

- [2026-08-07] blocked on "one modular suite release has shipped and the user has confirmed every client migrated its four skill directories, managed Just section, and enabled memory-hook paths" — cleared by user via clarify

## Decisions

- [2026-08-08] The user confirmed one modular suite release has shipped and every client has migrated its four skill directories, managed Just section, and previously enabled memory-hook paths.
- **D-01 — Delete, do not relabel, migration recognition.** The attested client baseline is the current four-skill layout. Capability probes, monolith/root-install handling, exact legacy recipe recognition, and old core memory-hook rewriting will be removed outright; retaining equivalent branches under neutral names would preserve the compatibility surface this REQ exists to retire.
- **D-02 — Keep one transaction owner.** The current updater validates the downloaded archive with the installed validator and delegates the same archive to the installed suite installer. Confirmation, dirty-path disclosure, managed-path planning, byte verification, and recovery remain in that installer instead of retaining a second updater-side transaction used only by retired monolith archives.
- **D-03 — Treat marker-free files generically.** The section replacer now recognizes only managed sentinels. A marker-free Justfile receives the managed section without inspecting old recipe names; any resulting reserved-name collision is caught by the installer's pre-write Just parse check, not by resurrecting a legacy parser.

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

---

## Triage

**Route: C** - Complex

**Reasoning:** This is the architectural cleanup half of the modular cutover: it removes command-routing and updater compatibility branches across the suite while preserving the current all-or-recover installer/updater and proving unique action ownership.

**Planning:** Required

## Plan

1. **Make the retirement RED.** Replace positive shim/bridge/legacy-migration assertions with absence/rejection contracts while retaining tests for the permanent four-module manifest, installed trust boundary, recovery transaction, managed markers, current hooks, and current-to-current updates.
2. **Delete the command compatibility layer.** Remove `skills/do-work/actions/moved-command-shim.md`, the three moved-command route groups, temporary compatibility prose, and legacy next-step/help entries. Core keeps only core commands plus concise discovery pointers for board, knowledge, and toolbox.
3. **Collapse updating to the modular path.** Remove `--capabilities`, `suite-layout-v2`, monolithic/no-manifest archives, root legacy installs, legacy path copying/deletion/messages, and any downloaded-helper execution. Preserve one archive, installed-validator/installer authority, Git-root/semver/dirty guards, confirmation, exact verification, and recovery.
4. **Remove one-release configuration migration.** Delete exact unmarked five-recipe migration and old core memory-hook rewriting while preserving marker-based Just reconciliation, custom content/modes, core hook composition, fresh memory-disabled behavior, already-modular hooks, idempotence, cancellation, and rollback.
5. **Refresh live guidance and ownership contracts.** Update README, version/help/prime/setup-memory wording and prove every public action has exactly one owning skill. Keep REQ-150's package-reference repair outside this REQ.
6. **Verify and release.** Run staged-skill, updater, installer, contract, manifest, syntax/ShellCheck, byte-identity, and diff checks; the owner performs the version/changelog commit after review.

**Planned implementation files:** `README.md`, core `SKILL.md`, help/version/next-steps, updater prime and script, both installer copies, both section-replacer copies, knowledge `setup-memory.md`, four shell contract suites, and deletion of `moved-command-shim.md`. Owner release files are excluded from the builder boundary until Step 9.

**Preserved boundary:** `suite/modules.tsv`, both manifest validators, the one-archive bootstrap, managed-marker reconciliation, installed trust boundary, all-or-recover writes, current modular hook paths, pipeline-guard cleanup for clients coming from 0.185.0, and all project/queue/KB/history data.

**Requirement coverage:** tasks 2–4 remove core routing and migration-window rules/messages; task 5 proves unique ownership and core-only help; tasks 1/6 preserve and verify bootstrap, updater, reconciliation, and current-to-current behavior.

**Plan validation:** Every Detailed Requirement maps to at least one task and no planned implementation task falls outside the scheduled compatibility retirement. ⚠ The plan has 6 tasks — quality degrades past 3; keep the builder to the exact transitional branches and their proof, with permanent recovery/configuration primitives explicitly protected.

*Generated by Plan agent*

## Exploration

The transitional surface is confined to the core moved-command routes/action, updater bridge branches, exact unmarked five-recipe migration, old core memory-hook rewriting, their live documentation, and four shell contract suites. Permanent machinery is deliberately protected: the four-module manifest/validators, installed validator and installer trust boundary, single-archive bootstrap, marker-based Just reconciliation, core hook composition, current modular memory hooks, and validate/review/write/verify/recover transaction.

REQ-149's moved-shim mapping loses its primary product surface after this deletion but remains queue data for a later explicit disposition; do not cancel it implicitly. REQ-150 retains ownership of broken live package references, including the version-action upstream URL and prime/history links. REQ-151 retains ownership of the no-JSON-tool pipeline-guard fallback wording; jq/Python cleanup for clients arriving from 0.185.0 remains protected in this REQ.

The two installer copies and two section-replacer copies must stay byte-identical. The installed-core changelog is a release sidecar, not builder scope; root/shared release metadata remains owner-only until Step 9. A general collision rejection for reserved recipe names may remain if needed, but no legacy recognition or automatic migration may survive under a new label.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `README.md`
- `skills/do-work/SKILL.md`
- `skills/do-work/actions/moved-command-shim.md` (delete)
- `skills/do-work/actions/help.md`
- `skills/do-work/actions/version.md` (compatibility prose only during implementation; preserve REQ-150's upstream-URL scope)
- `skills/do-work/next-steps.md`
- `skills/do-work/tools/prime-do-work-update.md`
- `skills/do-work/tools/do-work-update.sh`
- `tools/install-do-work-suite.sh`
- `skills/do-work/tools/install-do-work-suite.sh`
- `tools/replace-text-section.sh`
- `skills/do-work/tools/replace-text-section.sh`
- `skills/do-work-knowledge/actions/setup-memory.md`
- `_dev/tests/staged-skills-contract.sh`
- `_dev/tests/update-script-behavior.sh`
- `_dev/tests/install-suite-behavior.sh`
- `_dev/tests/contract-regressions.sh`

**Files I will NOT touch:** `suite/modules.tsv`, either manifest validator, root/shared changelog and VERSION release metadata before Step 9, REQ-150's package-reference fixes, REQ-151's no-tool fallback, queue/UR/KB/application/history data, and permanent bootstrap/recovery/marker/current-hook logic.

**Acceptance criteria (restated from REQ):**
- [ ] Old core routing shims for board, knowledge, and toolbox commands are absent.
- [ ] Bridge-only manifest compatibility, legacy/root install paths, stale-path deletion, capability probing, exact legacy recipe migration, and old core memory-hook rewriting are absent.
- [ ] Full-suite bootstrap, installed trust boundary, current modular updater, managed-marker reconciliation, current hooks, and current-to-current updates remain intact.
- [ ] Core help lists only core commands plus concise pointers to board, knowledge, and toolbox.
- [ ] Tests prove every active action exists in exactly one skill and old invocations are not live core routes.
- [ ] The two installer copies and two section-replacer copies remain byte-identical; contract, syntax/lint, Just, Go, and recovery checks pass.

## Pre-Flight

**Git:** ⚠ One pre-existing project-file modification (`decisions/records/adr-019-four-skill-suite-contract.md`) plus user-owned UR context. Preserve and exclude them from this REQ's staging.
**Tests baseline:** ✓ `bash _dev/tests/contract-regressions.sh` passed before implementation (`launched: true`).
**Dependencies:** ✓ Shell, Just, and Go validation paths required by the planned checks are available or have existing optional guards.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `README.md` (modified)
- `skills/do-work/SKILL.md` (modified)
- `skills/do-work/actions/moved-command-shim.md` (deleted)
- `skills/do-work/actions/help.md` (modified)
- `skills/do-work/actions/version.md` (modified)
- `skills/do-work/next-steps.md` (modified)
- `skills/do-work/tools/prime-do-work-update.md` (modified)
- `skills/do-work/tools/do-work-update.sh` (modified)
- `tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/install-do-work-suite.sh` (modified)
- `tools/replace-text-section.sh` (modified)
- `skills/do-work/tools/replace-text-section.sh` (modified)
- `skills/do-work-knowledge/actions/setup-memory.md` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Removed the one-release core command shims and every bridge-only updater/configuration migration path. The surviving current-suite updater continues through the installed validator and all-or-recover installer, marker-owned Just reconciliation, current hook composition, and permanent pipeline-guard cleanup, with updated absence and ownership contracts.

## Qualification

Passed — all 17 implementation files verified, all six acceptance criteria traced, P-A-U confirmed, and no hollow updater/configuration flow found. The diff is substantive and net subtractive; the surviving updater delegates a real archive through the installed validator and installer, and behavioral recovery tests exercise the retained data path. Overlap with REQ-145 is expected because REQ-146 explicitly depends on and is related to the same modular-cutover surfaces; no unexplained prior-REQ file appeared. Mechanical qualification passed, paired tools match, and protected REQ/UR/ADR/release boundaries remain unchanged by the builder.

## Testing

**Tests run:** `/bin/bash _dev/tests/contract-regressions.sh`; `bash _dev/tests/suite-manifest-contract.sh`; byte comparisons for installer, section-replacer, and manifest-validator pairs; repository-wide warning-level ShellCheck and Bash syntax; both Just manifests; queue-kanban Go format/tests/vet; `git diff --check`.
**Result:** ✓ All passing. The contract suite reported manifest, hash guards, updater, staged skills, installer, and full regressions green; paired tools are identical; Just parses; Go tests/vet pass; syntax/lint/diff checks are clean. One orchestration command initially launched the root Just check from the queue-kanban subdirectory and failed to find `justfile`; rerunning from the repository root passed, confirming a harness-path mistake rather than a product regression.

**Red-green validation:**
- `_dev/tests/staged-skills-contract.sh`: ✗ while the shim file/routes and old core memory-hook instructions remained → ✓ after their removal and unique-owner checks.
- `_dev/tests/update-script-behavior.sh`: ✗ while `--capabilities`, `suite-layout-v2`, legacy paths, and monolith handling remained → ✓ with current four-module installed-helper delegation and recovery coverage.
- `_dev/tests/install-suite-behavior.sh`: ✗ while legacy Just migration and old core memory-hook rewriting remained → ✓ with marker/current-hook behavior and pipeline-guard jq/Python cleanup preserved.
- `_dev/tests/contract-regressions.sh`: ✗ while the production migration surfaces and retired section-replacer flag were live → ✓ after the bounded deletion; assertion failures were behavior-level RED evidence, not launch/syntax failures.

**New tests added:**
- Absence/rejection contracts for moved routes/actions, capability probing, bridge/monolith paths, legacy recipe migration, and old core hook rewriting.
- Unique package ownership and current modular updater/install recovery assertions.

**Existing tests updated (cross-REQ impact):**
- REQ-137/138/143/144 bridge and migration fixtures now assert the attested post-window modular baseline while retaining their permanent trust, preservation, and recovery guarantees.

*Verified by work action*

## Review

**Overall: 79%** | 2026-08-08T18:30:15Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 82% |
| Test Adequacy | 84% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- The installer can report success with a duplicate reserved Just recipe when `just` is unavailable, leaving a Justfile that later fails to parse — gate: user-visible → REQ-152 created
- Live sibling assets and the updater prime retain removed core invocations and transition-era contract restatements — gate: user-visible → REQ-153 created

**Minor findings:** 1 (report only)
**Acceptance:** Partial — all prescribed suites pass, but the no-`just` reserved-recipe collision smoke test fails
**Suggested testing:** 3 items
**Follow-ups created:** REQ-152, REQ-153; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Deleting the compatibility branches while keeping the installed validator/installer as the sole transaction owner made the updater smaller without weakening confirmation, recovery, or managed-marker guarantees.
- Byte-identity checks for paired tools plus focused RED/GREEN contract suites kept the root and staged packages aligned through a large subtractive change.

**What didn't:**
- Treating `just` as an optional syntax check also made reserved-recipe collision detection optional; the independent no-`just` smoke test exposed the invalid-file path that the main fixtures missed.
- Removing router code alone did not retire the contract: a restatement sweep still found stale commands in shipped templates, UI hints, a hook message, and transition-era prime lessons.

**Worth knowing:**
- The permanent updater delegates one validated archive to the installed all-or-recover installer; future work should extend that boundary rather than recreate updater-side migration logic.
- REQ-152 owns tool-independent reserved-recipe collision rejection, and REQ-153 owns the remaining retired-command and prime-restatement sweep.

## Orientation

[MAP CHANGED] Core now exposes only core actions, while board, knowledge, and toolbox commands are owned by their sibling skills. The update subsystem now supports only the permanent four-skill path through the installed validator and all-or-recover installer. Its prime map remains temporarily stale in several transition-era lessons; REQ-153 is queued to reconcile them.
