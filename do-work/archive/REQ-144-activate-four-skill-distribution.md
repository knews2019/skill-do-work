---
id: REQ-144
title: "Activate the Four-Skill Distribution"
status: completed
claimed_at: 2026-08-08T14:46:57Z
status_changed_at: 2026-08-08T14:46:23Z
completed_at: 2026-08-08T15:39:49Z
commit: a3c2612
created_at: 2026-08-07T18:58:02Z
user_request: UR-031
domain: general
prime_files: [skills/do-work/tools/prime-do-work-update.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-137, REQ-143]
maintenance: true
kb_status: pending
kb_entry:
write_set: [.gitattributes, README.md, justfile, CLAUDE.md, VERSION, CHANGELOG.md, skills/do-work/SKILL.md, skills/do-work/next-steps.md, skills/do-work/actions/help.md, skills/do-work/actions/moved-command-shim.md, skills/do-work/actions/version.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/capture.md, skills/do-work/actions/kb-lessons-handoff.md, skills/do-work/actions/pipeline-reference.md, skills/do-work/actions/pipeline.md, skills/do-work/actions/roadmap.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/work.md, skills/do-work/crew-members/interviewer.md, skills/do-work/docs/roadmap-guide.md, skills/do-work/docs/standing-preferences.md, skills/do-work-board/actions/board.md, skills/do-work-board/docs/board-guide.md, skills/do-work-board/tools/queue-kanban/prime-do-kanban.md, skills/do-work-knowledge/docs/dream-guide.md, _dev/tests/staged-skills-contract.sh, _dev/tests/contract-regressions.sh, _dev/tests/update-script-behavior.sh, _dev/tests/install-suite-behavior.sh, SKILL.md, next-steps.md, actions/**, crew-members/**, docs/**, hooks/**, interviews/**, prompts/**, specs/**, tools/checks/**, tools/do-work-update.sh, tools/queue-kanban/**, tools/prime-do-work-update.md]
related: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-145, REQ-146]
batch: do-work-four-skill-suite
---

# Activate the Four-Skill Distribution

## What
Switch the live distribution to the four staged skills and migrate bridge-enabled client installations through one all-or-recover suite transaction.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** First close the newly discovered bridge seam by making the bridge updater delegate a future suite archive to the tested full-suite installer, with configuration diffs and one confirmation. Ship that while export guards remain, then require renewed client rollout evidence before deleting the monolith and activating the archive. After the gate, add exact moved-command shims, publish the canonical bootstrap, switch maintainer/runtime paths, remove shipped legacy duplication, and prove both update entry points converge.
- [x] **[APPLY]:** Activated the modular archive, added print-only moved-command shims, switched maintainer/runtime references to the four staged skills, strengthened bridge/direct-vs-Just/bootstrap contracts, and removed the duplicated monolithic runtime while retaining the three root bootstrap utilities.
- [x] **[UNIFY]:** Reviewed `git diff --stat` and every non-deletion implementation diff plus the complete deletion manifest; verified the staged core shim/router/help, archive guards, README/maintainer paths, four contract fixtures, retained root bootstrap boundary, no debug artifacts, ShellCheck at warning level, Bash syntax, Just parsing, Go formatting/tests/vet/build, full contract regressions, and `git diff --check`.

## Why
Staging alone does not reduce the installed router or context; the manifest and live archive must cut over only after clients can understand the new layout.

## Detailed Requirements
- Activate the live four-module manifest.
- Require user attestation that every known client repository's installed updater reports exactly `suite-layout-v2` from `--capabilities` before changing the live layout; no stored client-by-client inventory is required.
- Make `skills/do-work`, `skills/do-work-board`, `skills/do-work-knowledge`, and `skills/do-work-toolbox` the only shipped runtime sources.
- Remove duplicated legacy root actions/tools in the same change.
- Update export rules, README installation, root maintainer Just fallback, changelog/version handling, and every runtime cross-skill reference. Publish the exact canonical, tested, copy-paste bootstrap command from REQ-143 as the single fresh-install path in the README.
- Add one-release core shims for moved commands. Each shim prints the exact new invocation and stops; it must not forward permanently.
- Ensure `do-work update` and `just run-do-work-update` produce byte-equivalent installations from a bridge client.
- Enforce the all-or-recover suite contract: validate all four modules before the first managed write, report success only after all installed bytes verify, and restore every changed managed module path on failure so clients never observe a partially successful update. Do not claim a cross-directory filesystem-atomic rename.
- Refresh the managed Just section and migrate known memory-hook paths during update.
- Preserve project queue data, KB data, application files, and unrelated configuration.
- Run hermetic bridge-to-modular and fresh-install tests plus ShellCheck, formatting, Go tests/vet, contract tests, and `git diff --check`.

## Constraints
- Do not remove the live export guards until the user attests that every existing client has configuration-aware bridge v0.183.20 or later. The earlier `suite-layout-v2` capability alone predates managed Just/settings migration and is not sufficient evidence.
- Preserve the current feature-rich work orchestrator.
- Keep compatibility shims for exactly one modular release.

## Dependencies
Requires REQ-137 and REQ-143. The original capability gate was cleared by user attestation via `do-work clarify` on 2026-08-08; implementation then proved that a final configuration-aware bridge release must reach those clients before live cutover.

## Builder Guidance
Certainty level: Firm. Treat this as the release cutover; no partial module activation is acceptable.

## Red-Green Proof
**RED prompt/case:** Update a hermetic legacy client that has the bridge updater, legacy Just recipes, queue data, KB data, and an enabled legacy memory hook.
**Why RED now:** The live manifest and installed layout remain monolithic.
**GREEN when:** The user has attested that every known client reports the bridge capability; the README's exact bootstrap command succeeds from a clean Git client; one update installs exactly four current skill directories or restores the prior suite on failure; managed configuration migrates, project data survives, and both update entry points yield identical bytes.
**Validation:** User confirmed

## Full Context
See `do-work/user-requests/UR-031/input.md` for the bridge rollout gate and cutover contract.

## Blocked

- [2026-08-07] blocked on "user supplies an inventory of every known client repository confirming that each installed updater's --capabilities output is exactly suite-layout-v2" — cleared by user via clarify
- [2026-08-08] live cutover is blocked until v0.183.20 is published and the user confirms every known client has updated to that version or later through `do-work update` or `just run-do-work-update`

---
*Source: User approved the four-skill suite plan and requested capture of every required REQ.*

## Bridge Preparation

Strict RED first proved that the prior bridge left Testing-era board paths and legacy memory hooks unchanged during a suite update. The bridge now delegates the already-downloaded archive to its already-installed full-suite transaction engine, which reviews modules and owned configuration together, asks once, verifies the result, and recovers exact managed file bytes on failure. A hostile-archive fixture also proves that a downloaded installer cannot replace the installed trust boundary. The staged core now carries the installer so the first modular update does not remove the engine needed by later updates.

Version v0.183.20 is intentionally still a monolithic bridge release: `.gitattributes` continues to export-ignore `VERSION`, `suite/`, and `skills/`. Removing those guards, publishing the modular README bootstrap, and deleting legacy shipped duplication remain pending until renewed client rollout evidence is supplied.

**Bridge qualification:** `update-script-behavior`, `install-suite-behavior`, `staged-skills-contract`, the full contract regressions, warning-level ShellCheck, Bash syntax, Just parsing, `gofmt` cleanliness, queue-kanban `go test`/`go vet`/`go build`, synchronized runtime copies, and `git diff --check` all pass for v0.183.20.

---

## Triage

**Route: C** - Complex

**Reasoning:** This is the live four-skill distribution cutover spanning export policy, runtime-source removal, compatibility shims, updater/install transactions, configuration migration, release metadata, and hermetic rollout verification. Commit `76a0682` completed the configuration-aware bridge preparation, so this run must plan and verify the gated cutover rather than repeat that preparatory work.

**Planning:** Required

## Plan

1. **Lock the cutover baseline and prove RED.** Treat commit `76a0682` as completed bridge preparation; do not rework its trusted installed-validator/installed-installer transaction. The current authoritative `UR-031/input.md` records the user attestation for every known client and explicitly says no client inventory is required, so the renewed rollout gate is cleared and is not an external blocker. Extend `_dev/tests/staged-skills-contract.sh`, `_dev/tests/update-script-behavior.sh`, `_dev/tests/install-suite-behavior.sh`, and the cutover assertions in `_dev/tests/contract-regressions.sh` so the current bridge archive fails because `/VERSION`, `/suite`, and `/skills` are still export-ignored and legacy root runtime sources still ship. Pin the intended end state: one valid four-row manifest, exactly four installed sibling skills, explanatory non-forwarding shims, no legacy runtime tree, preserved client data/configuration, and equivalent output from both updater entry points.

2. **Finish the modular core compatibility surface before activation.** In `skills/do-work/SKILL.md`, `skills/do-work/actions/help.md`, `skills/do-work/next-steps.md`, and one new narrowly scoped `skills/do-work/actions/moved-command-shim.md`, route every command moved to board, knowledge, or toolbox through a one-release shim. The shim maps the legacy spelling and arguments to the exact `do-work-board …`, `do-work-knowledge …`, or `do-work-toolbox …` invocation, prints it, and stops; it must never load, forward to, or reimplement the sibling action. Keep the feature-rich `skills/do-work/actions/work.md` orchestrator and its `actions/kb-lessons-handoff.md` boundary unchanged, and retain the separate stateful pipeline for REQ-145 rather than removing it in this cutover.

3. **Activate the archive and make modular paths canonical.** Remove only the three temporary `/VERSION`, `/suite`, and `/skills` export guards from `.gitattributes`. Update `README.md` to publish verbatim the tested command emitted by `tools/install-do-work-suite.sh --print-bootstrap-command` as the single fresh-install path; describe four required sibling skills, all-or-recover updates, managed Just reconciliation, core hooks on/memory hooks off, and the new command names. Preserve the existing uncommitted “Upgrade an existing installation with an AI agent” block while editing surrounding README text. Update the root maintainer `Justfile` to build board commands from `skills/do-work-board/tools/queue-kanban` and update through `skills/do-work/tools/do-work-update.sh`. Update `CLAUDE.md` only where its maintainer source/version/test pointers still name the retired root runtime.

4. **Delete the monolithic runtime after all callers point at siblings.** Remove the duplicated root `SKILL.md`, `next-steps.md`, `actions/**`, `crew-members/**`, `docs/**`, `hooks/**`, `interviews/**`, `prompts/**`, `specs/**`, `tools/checks/**`, `tools/do-work-update.sh`, `tools/queue-kanban/**`, and `tools/prime-do-work-update.md`. Retain the root distribution/bootstrap infrastructure `tools/install-do-work-suite.sh`, `tools/validate-suite-manifest.sh`, and `tools/replace-text-section.sh`; keep the installer byte-identical to `skills/do-work/tools/install-do-work-suite.sh`. Run a repository-wide runtime-reference sweep and correct only live references under `skills/do-work/**`, `skills/do-work-board/**`, `skills/do-work-knowledge/**`, `skills/do-work-toolbox/**`, `README.md`, `CLAUDE.md`, `Justfile`, and `_dev/tests/**`; historical ADRs, archived REQs, reports, and changelogs remain historical and are not rewritten.

5. **Release one synchronized modular suite.** Bump the breaking cutover to `0.184.0` in `VERSION`, `skills/do-work/VERSION`, and `skills/do-work/actions/version.md`; prepend the same cutover entry to `CHANGELOG.md` and `skills/do-work/CHANGELOG.md`. Do not introduce per-module versions. Confirm `suite/modules.tsv` remains the only module/destination declaration and `skills/do-work/tools/do-work-update.sh` remains the engine used by both `do-work update` and the managed `run-do-work-update` recipe.

6. **Qualify the live cutover end to end.** Run strict bridge-to-modular fixtures from identical legacy clients through (a) the installed core updater and (b) the managed Just recipe, then compare all four installed module trees and managed configuration bytes. Exercise the README bootstrap against a clean Git repository; assert all four modules share the suite version, core hooks are enabled, fresh memory capture stays disabled, enabled legacy memory paths migrate, exterior Just/settings bytes plus `do-work/`, `kb/`, and application sentinels survive, and forced module/configuration failures restore exact originals. Run `bash _dev/tests/suite-manifest-contract.sh`, `bash _dev/tests/staged-skills-contract.sh`, `bash _dev/tests/update-script-behavior.sh`, `bash _dev/tests/install-suite-behavior.sh`, `/bin/bash _dev/tests/contract-regressions.sh`, warning-level ShellCheck and `bash -n` on shipped/test scripts, root and template Just parsing, `gofmt -l`, `go test ./...`, `go vet ./...`, and `go build ./...` from `skills/do-work-board/tools/queue-kanban`, an archive-content check proving the suite contract ships and the legacy runtime does not, and `git diff --check`.

### Requirement mapping and preservation boundaries

- Live manifest/archive activation: Steps 1, 3, and 6 (`.gitattributes`, `suite/modules.tsv`, archive tests).
- Four runtime sources plus legacy-duplicate removal: Steps 3–4 (`skills/do-work*` are canonical; root bootstrap helpers are infrastructure, not installed runtime).
- Exact bootstrap, updater convergence, configuration migration, all-or-recover behavior, and protected project data: Steps 3 and 6 (`README.md`, installer/updater/Just paths, hermetic fixtures).
- One-release moved-command shims and preserved work orchestrator/pipeline sequencing: Step 2.
- Shared release/version contract: Step 5.
- Existing dirty work: do not edit or stage `decisions/records/adr-019-four-skill-suite-contract.md`, `do-work/queue/REQ-146-remove-modular-migration-shims.md`, or `do-work/user-requests/UR-031/input.md`; leave `do-work/CHECKPOINT.md` and `do-work/working/REQ-144-activate-four-skill-distribution.md` to the queue owner. `README.md` is the one overlapping dirty file: preserve its pre-existing upgrade-prompt bytes and layer only the planned cutover documentation around them. Never touch runtime `do-work/`/`kb/`, application files, generated reports, or unrelated configuration.

*Generated by Plan agent*

**Plan validation:** Every Detailed Requirement maps to one or more plan steps, and every planned task traces to the live-cutover request. Warning: Plan has 6 tasks — quality degrades past 3; implementation should keep them as three execution phases (contract RED, cutover, qualification) rather than widening scope.

## Exploration

The dedicated Explore agent was interrupted before producing a durable hand-back, so the builder must perform final caller discovery itself. Orchestrator inspection established these source-of-truth seams:

- `.gitattributes` still carries exactly the three temporary `/VERSION`, `/suite`, and `/skills` bridge guards; the permanent `/do-work`, `/kb`, maintainer-doc, and development-path guards remain load-bearing.
- `_dev/tests/staged-skills-contract.sh` still asserts the pre-cutover root router/action exists, while `_dev/tests/contract-regressions.sh` still requires the three bridge guards. Both assertions must invert at activation.
- `tools/install-do-work-suite.sh --print-bootstrap-command` emits the exact one-archive fresh-install command already proven by `_dev/tests/install-suite-behavior.sh`; README must publish those bytes rather than paraphrase them.
- `suite/modules.tsv` already declares exactly the four sibling modules. The installed core updater delegates to its installed full-suite engine, and the managed board template routes `run-do-work-update` to core; preserve that one-engine convergence.
- The existing dirty README upgrade-agent block is directly relevant to bridge migration and will be adopted intact. The dirty ADR, UR, and REQ-146 changes are evidence/other queue state, not implementation files for this REQ, and remain unstaged.

*Explore agent failed; continued with reduced context per work action.*

## Decisions

- **D-01 — DECIDE & STATE:** Adopt the existing uncommitted README “Upgrade an existing installation with an AI agent” block into this cutover rather than overwrite or split it away. Reasoning: README is already a required REQ-144 surface and the block documents the bridge-client path this release activates; preserving it whole protects the user-authored migration guidance.
- **D-02 — DECIDE & STATE:** Retire the monolithic runtime sources physically after staged callers and tests are switched, while retaining only the root bootstrap/manifest/managed-section utilities needed to install the four packages. Reasoning: export-hiding duplicate runtime trees would leave two source-of-truth copies and recreate the drift this REQ exists to end; Git makes the deletion reversible.
- **D-03 — DECIDE & STATE:** Extend the declared scope to the staged runtime files named below after the builder's repository-wide reference sweep found live restatements of removed root paths. Reasoning: the REQ explicitly requires every runtime cross-skill reference to move with the cutover; leaving these callers stale would make the four packages internally inconsistent. The orchestrator inspected and formally adopted the exact paths before qualification.
- **D-04 — DECIDE & STATE:** Move this REQ's deferred-prime targets to the canonical staged core and board primes before retiring their root duplicates. Reasoning: preserving deleted root paths in `prime_files` would make the required post-archive lesson handoff impossible; the staged paths are the live sources after cutover.

## Scope

**Files I will touch:**
- `.gitattributes` (modify) — activate modular archive while preserving permanent export guards
- `README.md` (modify) — adopt the existing upgrade block and publish exact canonical bootstrap/module guidance
- `justfile` (modify) — use staged board/core maintainer paths
- `CLAUDE.md` (modify) — point maintainer structure/version/test instructions at canonical staged sources
- `VERSION` (modify) — shared suite version
- `CHANGELOG.md` (modify) — root suite release history
- `skills/do-work/SKILL.md` (modify) — one-release moved-command shim routes
- `skills/do-work/next-steps.md` (modify) — canonical sibling command suggestions
- `skills/do-work/actions/help.md` (modify) — core help with exact sibling invocations
- `skills/do-work/actions/moved-command-shim.md` (new) — print-only compatibility surface
- `skills/do-work/actions/version.md` (modify) — canonical core current version/update contract
- `skills/do-work/VERSION` (modify) — installed core shared suite version
- `skills/do-work/CHANGELOG.md` (modify) — installed core release history
- `skills/do-work/tools/prime-do-work-update.md` (modify) — deferred lesson link in the canonical core prime
- `skills/do-work/actions/capture.md` (modify) — sibling board references after root retirement
- `skills/do-work/actions/kb-lessons-handoff.md` (modify) — knowledge sibling invocation/reference
- `skills/do-work/actions/pipeline-reference.md` (modify) — toolbox sibling invocation/reference
- `skills/do-work/actions/pipeline.md` (modify) — toolbox sibling invocation/reference
- `skills/do-work/actions/roadmap.md` (modify) — toolbox/knowledge sibling invocation/reference
- `skills/do-work/actions/work-reference.md` (modify) — board sibling tool paths after root retirement
- `skills/do-work/actions/work.md` (modify) — board sibling tool paths after root retirement
- `skills/do-work/crew-members/interviewer.md` (modify) — knowledge sibling invocation/reference
- `skills/do-work/docs/roadmap-guide.md` (modify) — modular command guidance
- `skills/do-work/docs/standing-preferences.md` (modify) — modular command guidance
- `skills/do-work-board/actions/board.md` (modify) — toolbox sibling guidance
- `skills/do-work-board/docs/board-guide.md` (modify) — toolbox sibling guidance
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify) — canonical staged board paths
- `skills/do-work-knowledge/docs/dream-guide.md` (modify) — toolbox sibling guidance
- `_dev/tests/staged-skills-contract.sh` (modify) — live modular-source and shim contract
- `_dev/tests/contract-regressions.sh` (modify) — cutover archive/runtime contract
- `_dev/tests/update-script-behavior.sh` (modify) — bridge-to-live modular equivalence evidence
- `_dev/tests/install-suite-behavior.sh` (modify) — published bootstrap/live archive evidence
- `SKILL.md` (delete) — retire monolithic router
- `next-steps.md` (delete) — retire monolithic suggestions
- `actions/**` (delete) — retire duplicated monolithic actions
- `crew-members/**` (delete) — retire duplicated monolithic crews
- `docs/**` (delete) — retire duplicated monolithic guides
- `hooks/**` (delete) — retire duplicated monolithic hooks
- `interviews/**` (delete) — retire duplicated monolithic templates
- `prompts/**` (delete) — retire duplicated monolithic prompts
- `specs/**` (delete) — retire duplicated monolithic specs
- `tools/checks/**` (delete) — retire duplicated monolithic checks
- `tools/do-work-update.sh` (delete) — retire duplicated monolithic updater
- `tools/queue-kanban/**` (delete) — retire duplicated monolithic board tool
- `tools/prime-do-work-update.md` (delete) — retire duplicated monolithic prime

**Files I will NOT touch:** `decisions/records/adr-019-four-skill-suite-contract.md`, `do-work/queue/REQ-146-remove-modular-migration-shims.md`, and `do-work/user-requests/UR-031/input.md` remain preserved and unstaged. Runtime `do-work/`, `kb/`, application data, generated reports, and unrelated configuration are outside the implementation boundary; queue-owner edits to this REQ and checkpoint are metadata-only.

**Acceptance criteria (restated from REQ):**
- [ ] The live archive contains the shared version/manifest and exactly four installable skill sources, with no legacy root runtime duplication.
- [ ] Every legacy moved command has a one-release core shim that prints the exact sibling invocation and stops without forwarding.
- [ ] README publishes the exact tested fresh-install bootstrap and preserves the existing bridge-upgrade guidance.
- [ ] Both update entry points install byte-equivalent four-module/configuration state or restore every managed path on failure.
- [ ] Managed Just paths and enabled legacy memory hooks migrate; queue, KB, application, exterior configuration, and unrelated dirty files survive unchanged.
- [ ] The feature-rich work orchestrator and stateful pipeline remain in core for REQ-145; compatibility shims remain for REQ-146.
- [ ] Shared release metadata is synchronized and the required shell, contract, formatting, Go, archive, and diff checks pass.

## Pre-Flight

**Git:** ⚠ Pre-existing changes outside `do-work/` were detected. `README.md` is formally adopted under D-01; `decisions/records/adr-019-four-skill-suite-contract.md` remains preserved and unstaged. The existing REQ-146 and UR-031 edits are likewise preserved as queue/evidence state outside this implementation.
**Tests baseline:** ✓ `/bin/bash _dev/tests/contract-regressions.sh` passed before cutover edits (`launched: true`, exit 0).
**Dependencies:** ✓ The baseline launched successfully with the repository's shell, Python, Just, ShellCheck, and Go-dependent contract probes available.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/moved-command-shim.md` (new) — added the print-only one-release moved-command compatibility action
- `.gitattributes` (modified) — updated the live modular cutover contract, path, or documentation surface
- `CLAUDE.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `README.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `_dev/tests/contract-regressions.sh` (modified) — updated the live modular cutover contract, path, or documentation surface
- `_dev/tests/install-suite-behavior.sh` (modified) — updated the live modular cutover contract, path, or documentation surface
- `_dev/tests/staged-skills-contract.sh` (modified) — updated the live modular cutover contract, path, or documentation surface
- `_dev/tests/update-script-behavior.sh` (modified) — updated the live modular cutover contract, path, or documentation surface
- `justfile` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work-board/actions/board.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work-board/docs/board-guide.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work-knowledge/docs/dream-guide.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/SKILL.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/capture.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/help.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/kb-lessons-handoff.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/pipeline-reference.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/pipeline.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/roadmap.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/work-reference.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/actions/work.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/crew-members/interviewer.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/docs/roadmap-guide.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/docs/standing-preferences.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `skills/do-work/next-steps.md` (modified) — updated the live modular cutover contract, path, or documentation surface
- `SKILL.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/abandon.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/ai-report-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/ai-report.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/bkb-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/bkb.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/board.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/capture-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/capture.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/clarify.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/cleanup.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/code-review.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/commit.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/deep-explore-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/deep-explore.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/dream.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/forensics.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/help.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/inspect.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/install.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/interview-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/interview.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/kb-lessons-handoff.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/memory-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/memory-value.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/memory.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/note.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/pipeline-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/pipeline.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/present-work.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/prime.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/prompts.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/quick-wins.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/review-work.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/roadmap.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/sample-archived-req.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/scan-ideas.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/slop-check.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/stray-check.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/tidy-repo.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/tutorial.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/ui-review.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/validate-feedback.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/verify-requests.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/version.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/work-reference.md` (deleted) — removed the duplicated monolithic runtime source
- `actions/work.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/anti-slop.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/approach-directives.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/backend.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/background-agents.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/caveman.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/clear-questions.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/cms.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/coding-guardrails.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/debugging.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/frontend.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/general.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/interviewer.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/maintenance.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/prompt-injection.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/security.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/testing.md` (deleted) — removed the duplicated monolithic runtime source
- `crew-members/ui-design.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/ai-report-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/bkb-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/board-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/capture-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/cleanup-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/code-review-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/commit-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/dream-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/forensics-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/inspect-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/interview-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/present-work-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/prime-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/prompts-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/quick-wins-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/review-work-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/roadmap-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/slop-check-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/standing-preferences.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/stray-check-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/ui-review-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/verify-requests-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/version-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `docs/work-guide.md` (deleted) — removed the duplicated monolithic runtime source
- `hooks/hooks.json` (deleted) — removed the duplicated monolithic runtime source
- `hooks/memory-hooks.json` (deleted) — removed the duplicated monolithic runtime source
- `hooks/memory-session-start.sh` (deleted) — removed the duplicated monolithic runtime source
- `hooks/memory-stop-capture.sh` (deleted) — removed the duplicated monolithic runtime source
- `hooks/pipeline-guard.sh` (deleted) — removed the duplicated monolithic runtime source
- `hooks/session-start.sh` (deleted) — removed the duplicated monolithic runtime source
- `interviews/work-operating-model.md` (deleted) — removed the duplicated monolithic runtime source
- `next-steps.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/README.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/architecture-decisions-log_create-or-expand.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/business-vendor-strategic-sort.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/dark-code-kit_audit.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/dark-code-kit_comprehension-gate.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/dark-code-kit_context-layer-generator.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/economics-inference-stress-test.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/economics-saas-repricing-exposure.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/prompt-kit-step0-pen-and-paper-exercises-to-prepare-prompt.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/prompt-kit-step1-four-discipline-diagnostic.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/prompt-kit-step2-personal-context-doc.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/prompt-kit-step3-spec-engineer.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/prompt-kit-step4-intent-and-delegation-framework.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/prompt-kit-step5-eval-harness.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/prompt-kit-step6-constraint-architecture.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/tech-inference-architecture-decision.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/tech-infrastructure-compute-geography-risk.md` (deleted) — removed the duplicated monolithic runtime source
- `prompts/weekly-structural-diff-original.md` (deleted) — removed the duplicated monolithic runtime source
- `specs/README.md` (deleted) — removed the duplicated monolithic runtime source
- `specs/api-endpoint.md` (deleted) — removed the duplicated monolithic runtime source
- `specs/bug-fix.md` (deleted) — removed the duplicated monolithic runtime source
- `specs/refactor.md` (deleted) — removed the duplicated monolithic runtime source
- `specs/ui-component.md` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/archive-collision.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/associate-files.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/blanked-req-scan.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/preflight.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/qualify.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/record-commit-hash.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/scope-drift.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/checks/uncommitted-inventory.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/do-work-update.sh` (deleted) — removed the duplicated monolithic runtime source
- `tools/prime-do-work-update.md` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/.gitignore` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/allocate.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/allocate_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/atomic_replace_unix.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/atomic_replace_unsupported.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/atomic_replace_windows.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/atomic_write.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/board_live_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/board_synthetic_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/completion_anomaly_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/dependency_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/filementions.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/filementions_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/frontmatter.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/frontmatter_cli.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/frontmatter_cli_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/frontmatter_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/future_timestamp_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/generate.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/generate_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/go.mod` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/go.sum` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/main.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/model.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/model_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/notes_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/open_work.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/open_work_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/openbrowser.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/openbrowser_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/prime-do-kanban.md` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/release.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/release_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/render.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/serve.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/serve_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/state_timer_source_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/testing.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/testing_api.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/testing_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/timestamp.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/timestamp_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/verify.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/verify_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/walk.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/walk_test.go` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/web/board.css` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/web/board.js` (deleted) — removed the duplicated monolithic runtime source
- `tools/queue-kanban/web/template.html` (deleted) — removed the duplicated monolithic runtime source

Activated the live four-skill archive after the configuration-aware bridge release: the staged core, board, knowledge, and toolbox packages are now canonical; old core spellings print exact sibling invocations without forwarding; both updater entry points converge on the same all-or-recover suite transaction; maintainer paths and installation guidance use the modular layout; and the duplicated monolithic runtime is removed while root bootstrap infrastructure remains.

## Qualification

Passed — 204 implementation files verified against the working diff, all Detailed Requirements traced to the modular router/archive/updater/configuration/test surfaces, the existing installed-validator and full-suite transaction remain substantive and flowing, and P-A-U is confirmed. `skills/do-work/tools/checks/qualify.sh` reported no mechanical failures; the orchestrator independently reviewed every non-deletion diff and the complete deletion manifest.

## Testing

**Tests run:**
- `bash _dev/tests/suite-manifest-contract.sh`
- `bash _dev/tests/staged-skills-contract.sh`
- `bash _dev/tests/update-script-behavior.sh`
- `bash _dev/tests/install-suite-behavior.sh`
- `/bin/bash _dev/tests/contract-regressions.sh`
- `find skills tools _dev/tests -type f -name '*.sh' -print0 | xargs -0 shellcheck -S warning`
- `find skills tools _dev/tests -type f -name '*.sh' -print0 | xargs -0 bash -n`
- `just --justfile justfile --list`
- `just --justfile skills/do-work-board/justfile.template --list`
- `test -z "$(gofmt -l skills/do-work-board/tools/queue-kanban)"`
- `(cd skills/do-work-board/tools/queue-kanban && go test ./... && go vet ./... && go build ./...)`
- `git diff --check`
- `bash skills/do-work/tools/checks/qualify.sh do-work/working/REQ-144-activate-four-skill-distribution.md`

**Result:** ✓ All passing.

**Red-green validation:** A detached worktree at the pre-cutover HEAD ran the new live-cutover contract and failed with the intended RED assertions: the live archive still excluded `/VERSION`, `/suite`, and `/skills`; the legacy root router/runtime still existed; the moved-command shim was absent; and zero moved-command routes were registered. On the implementation tree, the same staged-skills contract plus the updater, installer, suite-manifest, and full contract suites pass GREEN. The hermetic update fixtures also prove direct and Just entry-point parity, clean bootstrap, managed configuration migration, and exact recovery on injected failure.

**Existing tests updated (cross-REQ impact):** REQ-136's archive-guard expectation is intentionally inverted by this activation; the REQ-137/REQ-143 bridge and installer coverage is extended to assert live modular parity and recovery.

*Verified by work action*

## Review

**Overall: 77%** | 2026-08-08T15:38:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 85% |
| Code Quality | 83% |
| Test Adequacy | 86% |
| Scope | 92% |
| Risk | Low |
| Acceptance | Partial |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- The compatibility router omits concrete mappings for several routed aliases and maps `install run-do-work-update` to the wrong recipe; the current contract only greps placeholders. — gate: user-visible → REQ-149 created
- The live version URL, updater-prime lesson links, and installed changelog sidecar links do not survive the modular package layout. — gate: user-visible → REQ-150 created

**Minor findings:** 0 (report only)
**Acceptance:** Partial — primary four-skill install/update and recovery flows pass, but migration aliases and published package references remain incomplete.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-149, REQ-150; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:**
- Shipping the configuration-aware bridge first let the live cutover reuse one already-trusted validate/review/install/verify/recover transaction for direct and managed-Just updates.
- Inverting the archive contract before deleting the monolith produced a precise RED signal for every root runtime surface the release had to retire.

**What didn't:**
- The earlier capability token was too broad; the hermetic configuration fixture showed that a final bridge release needed to carry managed Just and hook migration before activation.
- Content-only shim assertions and a token-level reference sweep missed alias semantics and package-relative links; REQ-149 and REQ-150 preserve those review findings as explicit follow-up work.

**Worth knowing:**
- Runtime source ownership now lives in the four manifest-declared packages. The repository root retains only the installer, manifest validator, and managed-text replacement bootstrap utilities.
- The print-only compatibility layer is intentionally one-release work; REQ-146 removes it after the migration window, while REQ-145 removes the preserved stateful pipeline.

## Orientation

[MAP CHANGED] The distribution now installs four sibling skills—core orchestration, board, knowledge, and toolbox—from one shared manifest and version, with the monolithic root runtime removed. Fresh installs and bridge upgrades converge on the same all-or-recover suite transaction, so a client either receives one verified modular configuration or its exact managed state is restored. The narrowed prime check found stale links in both touched prime surfaces; REQ-149 and REQ-150 keep the remaining compatibility and package-reference corrections visible before the migration window closes.
