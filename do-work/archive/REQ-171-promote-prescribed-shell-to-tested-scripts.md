---
id: REQ-171
title: "Addendum: promote prescribed shell primitives to shipped, fixture-tested scripts"
status: completed
completed_at: 2026-08-11T21:17:02Z
kb_status: pending
kb_entry: Shell quoting does not disable Git pathspec magic at exact-path boundaries
created_at: 2026-08-11T13:58:25Z
user_request: UR-038
addendum_to: REQ-165
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-167, REQ-170]
batch: stabilization-audit
claimed_at: 2026-08-11T19:56:32Z
route: C
write_set:
  - decisions/audits/2026-08-11-prescribed-shell-primitives.md
  - decisions/audits/2026-08-11-defensive-surface.md
  - skills/do-work/docs/prescribed-shell-primitives.md
  - skills/do-work/scripts/show-commit-diff.sh
  - skills/do-work/scripts/add-local-git-exclude.sh
  - skills/do-work/scripts/atomic-download.sh
  - skills/do-work/scripts/capture-screenshot.sh
  - skills/do-work/scripts/run-blocked-check.sh
  - skills/do-work/scripts/protected-inventory.sh
  - skills/do-work/scripts/stage-exact-deletion.sh
  - skills/do-work-knowledge/scripts/lexical-memory-recall.sh
  - skills/do-work-knowledge/scripts/install-memory-hooks.sh
  - skills/do-work-toolbox/scripts/generate-report-image.sh
  - skills/do-work-toolbox/scripts/install-last30days.sh
  - skills/do-work/actions/capture.md
  - skills/do-work/actions/commit.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/work.md
  - skills/do-work-board/actions/board.md
  - skills/do-work-knowledge/actions/memory-reference.md
  - skills/do-work-knowledge/actions/setup-memory.md
  - skills/do-work-toolbox/actions/ai-report.md
  - skills/do-work-toolbox/actions/ai-report-reference.md
  - skills/do-work-toolbox/actions/inspect.md
  - skills/do-work-toolbox/actions/install.md
  - skills/do-work-toolbox/actions/present-work.md
  - skills/do-work/crew-members/background-agents.md
  - skills/do-work-knowledge/crew-members/background-agents.md
  - skills/do-work-toolbox/crew-members/background-agents.md
  - _dev/tests/fixture-repo.sh
  - _dev/tests/prescribed-shell-scripts-behavior.sh
  - _dev/tests/prescribed-shell-canonicalization.sh
  - _dev/tests/contract-regressions.sh
  - _dev/tests/staged-skills-contract.sh
---

# Addendum: Promote Prescribed Shell Primitives to Shipped, Fixture-Tested Scripts

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Completed the 59-fence census and durable 17 promote / 21 residue / 2 Go-owner inventory before writing scripts; mapped the minimal 11-script set and fixture seams.
- [x] **[APPLY]:** Added seven core, two knowledge, and two toolbox scripts; contracted callers to intent plus invocation; added 11 named fixture cases and updated canonical/staged/audit contracts.
- [x] **[UNIFY]:** Audited all 34 changed paths, reconciled the screenshot and finding-closure overlaps, ran syntax/ShellCheck, focused and full contracts, final fence accounting, debug/scope checks, and `git diff --check`.

## What

Graduate the canonical prescribed-shell primitives from documented prose to real, shipped script files with fixture-repo execution tests. Each *multi-line* primitive in `skills/do-work/docs/prescribed-shell-primitives.md` (and the remaining multi-line blocks in action files, e.g. capture.md Step 4's screenshot copy/verify/link block) becomes a `.sh` file under a per-package `scripts/` directory (core first: `skills/do-work/scripts/`); call sites keep a one-line intent statement plus the invocation; `_dev/tests/` gains a fixture-repo scaffold (mktemp repo, git init, seeded queue/version fixtures) that *executes* each script and asserts output and exit codes. The dividing line is "does this block contain logic that can be wrong" — one-liners and illustrative fragments stay inline, covered by the existing lint harness as residue.

## Context

Addendum to REQ-165 (completed). The plan this implements was approved by the user *while the original batch was being built* — the queue ran between plan approval and this capture, so the delta arrives as an addendum rather than a reshape. Lint (what shipped) catches syntax and quoting; it structurally cannot catch the semantic trap classes that motivated the batch — pipefail-kills-the-fallback, porcelain collapsing untracked dirs, curl partial files, merge-commit `git show` — which only surface when the commands *execute*. Only execution tests close that class; REQ-166's `session-start-hook-behavior.sh` is the proof of shape and the pattern to generalize.

## Prior Implementation

- **REQ-165** (commit `a45d5c4`, Route C): built `_dev/tests/action-shell-blocks.sh` (214 lines) — extracts fenced `bash`/`sh` blocks from the shipped `skills/` tree, runs `bash -n` always and shellcheck when present; wired into contract-regressions. This harness **stays** — its role narrows to the inline residue once multi-line blocks are promoted.
- **REQ-167** (commit `1a27c07`): created the canonical prose home `skills/do-work/docs/prescribed-shell-primitives.md` with eight primitive sections (per-file untracked inventory, merge-aware commit diff, commit file listing, local git ignore, atomic download publication, raw text before shell quoting, diff output filtering, state across command blocks), pointed consuming action files at it, and ratcheted the arrangement with `_dev/tests/prescribed-shell-canonicalization.sh` plus an audit record in `decisions/audits/2026-08-11-prescribed-shell-primitives.md`.
- **REQ-166** (commit `6538bdd`, Route A): simplified the session-start hook and added `_dev/tests/session-start-hook-behavior.sh` — the fixture-execution pattern this addendum generalizes.

## Detailed Requirements

- Build the shared fixture scaffold first (helpers to create a throwaway repo with seeded queue/version/untracked-file state), following `_dev/tests/` conventions; REQ-166's hook test may be migrated onto it if that is a simplification, not a rewrite.
- Promote each qualifying primitive to a script; the canonical guide's sections then document intent and *point at the script* as the normative implementation — update `prescribed-shell-canonicalization.sh`'s contract accordingly rather than deleting it (the single-home ratchet must survive the move, now enforcing script-as-home).
- Every promoted script gets fixture tests covering the trap it exists to avoid (e.g. the download script's mid-transfer-failure case; the untracked-inventory script against a wholly-untracked directory).
- Go-owned capabilities (atomic REQ reservation) get no shell twin — script layer is shell-portable primitives only.
- Call sites keep one sentence of intent so action files still work as standalone pasted prompts; the floor (read/write files, run shell) is respected because scripts ship inside the packages.
- Cross-package: a primitive used by board/knowledge/toolbox actions lives in core `scripts/` and is referenced with explicit sibling paths, same direction rules as prose cross-references.
- Net-surface accounting in the report: lines of prescribed shell removed from prose vs. script+test lines added — the prose side must shrink.

## Builder Guidance

Certainty: Firm on the architecture (approved plan); exploratory on which primitives qualify — produce the promotion inventory (primitive → script, or stays-inline rationale) as the first artifact, then execute it. Migrate incrementally with the suite green after each promotion; this is one REQ because the scaffold is shared, but the promotions are independent and abortable partway with value retained.

## Red-Green Proof

**RED prompt/case:** Reintroduce a semantic trap into a canonical primitive — e.g. change the download primitive back to plain `curl -o` without temp-and-rename, or give a version-parse pipeline the `set -euo pipefail` dead-fallback shape. Today `action-shell-blocks.sh` (bash -n + shellcheck) passes both, and nothing executes the primitive to notice.
**Why RED now:** The shipped harness is lint-only; the canonical guide is prose-only. The exact bug class that motivated UR-036 (demonstrated live in the session-start hook) remains undetectable for every primitive except the hook itself.
**GREEN when:** The promoted scripts exist with fixture tests; seeding either regression above into its script makes the suite fail naming the script and case; the canonicalization ratchet enforces script-as-home; `action-shell-blocks.sh` still covers inline residue; net prescribed-shell prose shrank.
**Validation:** User confirmed (plan approved with "capture"; delta re-anchored against the shipped batch by the capturing agent).

## Full Context

See `do-work/user-requests/UR-038/input.md` for complete verbatim input.

---
*Source: "capture" — approving the stabilization plan v2 discussed in-session (UR-038)*

---

## Triage

**Route: C** - Complex

**Reasoning:** This introduces a shared fixture harness, promotes several canonical shell primitives into shipped scripts, rewires cross-package call sites, and adds execution-level regression coverage.

**Planning:** Required

## Plan

1. Extend the prescribed-shell audit into a complete, heading-addressed promotion inventory; every current multi-line shell fence gets a `promote`, `already executable`, `inline residue`, or `non-shell owner` disposition before scripts are written.
2. Add a shared fixture-repository helper and a focused behavior probe, then make the canonicalization ratchet require advertised executable homes and reject promoted prose copies.
3. Promote shared mechanics first (merge-aware commit diff, local Git exclude, atomic download), then core action-specific mechanics (screenshot publication, blocked checks, protected inventory, exact deletion).
4. Promote package-specific knowledge/toolbox logic only where the inventory proves reusable control flow, keeping cross-package primitives in core and leaving the Go-owned REQ reservation without a shell twin.
5. Contract each call site to one intent sentence plus a shipped-script invocation; action-specific consent, policy, and reporting remain in the action.
6. Keep `_dev/tests/action-shell-blocks.sh` for inline residue and direct script lint, and explicitly wire the new execution probe into the aggregate contract suite.
7. Update the defensive-surface audit for each shipped script and prove source/install sibling paths remain valid.
8. Re-run the inventory, require shipped Markdown shell lines to decrease, and report prose removed versus script/test lines added.
9. Reconcile branch overlaps serially and run the focused probes plus the complete contract suite on the integrated tree.

**Architectural decisions:** Scripts receive raw values as separate quoted arguments or stdin, own mechanics rather than workflow policy, target Bash 3.2/POSIX utilities, and keep core as the sole home for cross-package primitives. The suite installer remains a justified self-bootstrap exception.

**Requirement mapping:** The inventory gates promotion scope; fixture tests cover each promoted semantic trap; canonicalization enforces script-as-home; one-line call sites retain standalone readability; Go reservation stays Go-owned; final accounting proves the prose surface shrank.

**Testing approach:** Establish RED through missing-script/canonical-form cases, make each behavior probe green incrementally, then run script syntax/ShellCheck, prescribed-shell behavior and canonicalization, action-shell lint, staged/shipped package contracts, defensive-surface audit, full contract regressions, and diff/net-surface checks.

*Generated by Plan agent*

**Plan validation:** Every Detailed Requirement maps to planned work and no task is orphaned. ⚠ Plan has 9 tasks — quality degrades past 3; the builder must keep the inventory authoritative, land independently testable promotions, and stop/report rather than silently widening beyond the inventory.

## Exploration

- The shipped tree currently contains 59 shell fences; 40 are multi-line (515 physical body lines). The complete inventory classifies 17 blocks for promotion, 21 as inline residue, and 2 as Go-owned. The smallest compliant set is 11 scripts: seven core, two knowledge, and two toolbox.
- Existing low-level inventory/association and commit-hash scripts remain owners; the new wrappers should own only reusable cross-block orchestration. Action-specific consent, trust, workflow policy, and commit templates stay in prose.
- Tests need one shared fixture helper plus a named behavior case for every promoted semantic trap. Existing canonicalization, action-shell lint, staged package, defensive-surface, updater, and aggregate test seams must be rewired rather than weakened.
- The main tree's uncommitted screenshot change replaces a fixed `.copying` path with a unique adjacent `mktemp` copy and adds a coordinated two-writer race test. A clean worktree lacks it; the builder must port that exact behavior into `capture-screenshot.sh` and the new behavior probe.
- Release files are already dirty at 0.186.30 and remain integrator-only. Other dirty helper/installer/prime changes are unrelated and must not enter the branch.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/contract-regressions.sh` (modified) — aggregate behavior wiring and executable-owner assertions
- `_dev/tests/fixture-repo.sh` (new) — sourceable Git fixture helpers
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified) — executable-home canonicalization
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (new) — 11 named semantic cases
- `_dev/tests/staged-skills-contract.sh` (modified) — staged script-resolution contract
- `decisions/audits/2026-08-11-defensive-surface.md` (modified) — executable defensive-surface evidence
- `decisions/audits/2026-08-11-prescribed-shell-primitives.md` (modified) — complete fence inventory/accounting
- `skills/do-work-board/actions/board.md` (modified) — core local-exclude invocation
- `skills/do-work-knowledge/actions/memory-reference.md` (modified) — lexical/hook script invocations
- `skills/do-work-knowledge/actions/setup-memory.md` (modified) — core ignore-helper invocation
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified) — core ignore-helper pointer
- `skills/do-work-knowledge/scripts/install-memory-hooks.sh` (new) — hook merge/verification/rollback
- `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` (new) — bounded lexical recall
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — report-image invocation
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — merge-aware display invocation
- `skills/do-work-toolbox/actions/inspect.md` (modified) — protected inventory invocation
- `skills/do-work-toolbox/actions/install.md` (modified) — atomic download and last30days invocation
- `skills/do-work-toolbox/actions/present-work.md` (modified) — merge-aware display invocation
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified) — core ignore-helper pointer
- `skills/do-work-toolbox/scripts/generate-report-image.sh` (new) — exact-output report image generation
- `skills/do-work-toolbox/scripts/install-last30days.sh` (new) — install/repair/full verification
- `skills/do-work/actions/capture.md` (modified) — screenshot-helper invocation
- `skills/do-work/actions/commit.md` (modified) — protected inventory and exact-deletion invocation
- `skills/do-work/actions/review-work.md` (modified) — merge-aware display invocation
- `skills/do-work/actions/work.md` (modified) — blocked-check invocation
- `skills/do-work/crew-members/background-agents.md` (modified) — local-ignore helper pointer
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified) — canonical executable-home guide
- `skills/do-work/scripts/add-local-git-exclude.sh` (new) — worktree-safe local exclude
- `skills/do-work/scripts/atomic-download.sh` (new) — private adjacent publication
- `skills/do-work/scripts/capture-screenshot.sh` (new) — unique verified no-clobber publication
- `skills/do-work/scripts/protected-inventory.sh` (new) — protected inventory orchestration
- `skills/do-work/scripts/run-blocked-check.sh` (new) — portable bounded probe
- `skills/do-work/scripts/show-commit-diff.sh` (new) — merge-aware commit display
- `skills/do-work/scripts/stage-exact-deletion.sh` (new) — exact cached deletion staging

**Files I will NOT touch:** `VERSION`, changelogs, `skills/do-work/actions/version.md`, `suite/modules.tsv`, `skills/do-work/tools/install-do-work-suite.sh`, Go reservation code, active `do-work/` state, or dirty main-only installer/helper/prime changes unrelated to promotion

**Acceptance criteria (restated from REQ):**
- [ ] Every current multi-line shell fence has a durable disposition; every promoted row names one executable home and behavior case.
- [ ] The 11 minimal scripts implement the semantic traps without cloning existing low-level owners or Go capabilities.
- [ ] Each promoted caller retains standalone intent plus one invocation; scripts own mechanics and accept values as quoted arguments/stdin.
- [ ] Fixture tests fail under the captured atomic-download and merge-diff regressions and pass on the correct scripts.
- [ ] Canonicalization enforces script-as-home; action-shell lint continues covering residue and shipped scripts.
- [ ] Staged/installed cross-package paths, defensive-surface audit, updater fixtures, and the aggregate suite pass.
- [ ] Shipped Markdown shell lines decrease, with removed prose versus script/test lines reported.

## Pre-Flight

**Git:** ⚠ The main tree is dirty with the screenshot race/release work and REQ-173 helper/test state. The branch must start clean, port only the explicitly required screenshot semantics, and leave serial release metadata to reconciliation.
**Tests baseline:** ⚠ `_dev/tests/contract-regressions.sh` already fails in the main tree only because the uncommitted REQ-173 prime link targets repo-only archive content; the clean branch baseline is expected to pass.
**Dependencies:** ✓ Existing Bash/Python/Git test toolchain is available; optional ShellCheck degrades by existing contract.

*Checked by work action*

## Decisions

- **D-01 — Port screenshot semantics into the script.** `capture-screenshot.sh` owns unique adjacent `mktemp`, exact-copy verification, no-clobber hard-link publication, loser preservation, and best-effort cleanup; the retired inline block and race harness were not reintroduced during integration.
- **D-02 — Parameterize only protected-inventory quarantine identity.** Core commit and toolbox inspect share orchestration while keeping independent Git-private namespaces and delegating classification to existing low-level checks.
- **D-03 — Keep workflow policy at callers.** Consent, trust, parallel orchestration, scratch ownership, and reporting remain in action prose; scripts own reusable mechanics only.
- **D-04 — Give last30days explicit `check` and `install` modes.** Install repairs the ignore guarantee even for an existing skill, while check is read-only and verifies the full contract.
- **D-05 — Preserve bootstrap and Go ownership.** The updater/suite bootstrap remains self-contained, and atomic REQ reservation stays with its Go owner rather than gaining a shell twin.
- **D-06 — Make Git path semantics explicit.** Shell quoting is not literal-pathspec handling; exact-deletion queries and staging now use Git's `--literal-pathspecs` mode, with a two-deletion magic-name fixture proving index isolation.

## Implementation Summary

- `_dev/tests/contract-regressions.sh` (modified) — wires behavior probes and executable-owner assertions while preserving the producer-complete closure ratchet.
- `_dev/tests/fixture-repo.sh` (new) — sourceable Git fixture helpers.
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified) — requires all 11 executable homes and rejects promoted Markdown implementations.
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (new) — 11 named semantic fixture cases, including the coordinated screenshot race.
- `_dev/tests/staged-skills-contract.sh` (modified) — validates shipped script resolution and delegates screenshot behavior to the new probe.
- `decisions/audits/2026-08-11-defensive-surface.md` (modified) — inventories the executable defensive surfaces with incident/evidence.
- `decisions/audits/2026-08-11-prescribed-shell-primitives.md` (modified) — records the complete census, disposition, executable homes, cases, and accounting.
- `skills/do-work-board/actions/board.md` (modified) — invokes the core local-exclude helper.
- `skills/do-work-knowledge/actions/memory-reference.md` (modified) — invokes lexical recall and hook installation.
- `skills/do-work-knowledge/actions/setup-memory.md` (modified) — resolves the core local-exclude helper.
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified) — points transient ignore handling to the core script.
- `skills/do-work-knowledge/scripts/install-memory-hooks.sh` (new) — merges, verifies, and rolls back partial hook installation.
- `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` (new) — performs inert tokenization, scoring, recency, attribution, and bounded output.
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — invokes report-image generation while retaining orchestration.
- `skills/do-work-toolbox/actions/ai-report.md` (modified) — emits merge-aware commit display.
- `skills/do-work-toolbox/actions/inspect.md` (modified) — invokes protected inventory with its own quarantine namespace.
- `skills/do-work-toolbox/actions/install.md` (modified) — routes downloads and last30days checks through executable owners.
- `skills/do-work-toolbox/actions/present-work.md` (modified) — emits merge-aware commit display across its consumers.
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified) — resolves the core local-exclude helper.
- `skills/do-work-toolbox/scripts/generate-report-image.sh` (new) — selects an exact-path backend with opt-in agentic scratch.
- `skills/do-work-toolbox/scripts/install-last30days.sh` (new) — installs/repairs and verifies the full sibling/ignore/Python guarantee.
- `skills/do-work/actions/capture.md` (modified) — invokes screenshot publication while retaining source-cleanup policy and REQ-170's finding-origin proof hook.
- `skills/do-work/actions/commit.md` (modified) — invokes protected inventory and exact-deletion staging.
- `skills/do-work/actions/review-work.md` (modified) — invokes merge-aware commit display while retaining producer-complete closure proof templates.
- `skills/do-work/actions/work.md` (modified) — invokes the portable blocked-check runner.
- `skills/do-work/crew-members/background-agents.md` (modified) — points local-ignore mechanics to the shipped helper.
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified) — makes scripts normative and records the bootstrap/Go exceptions.
- `skills/do-work/scripts/add-local-git-exclude.sh` (new) — appends once to the actual worktree-local exclude file.
- `skills/do-work/scripts/atomic-download.sh` (new) — downloads privately beside the target and publishes atomically.
- `skills/do-work/scripts/capture-screenshot.sh` (new) — verifies and no-clobber publishes one dispatch-owned screenshot copy.
- `skills/do-work/scripts/protected-inventory.sh` (new) — wraps quarantine, classification, and association mechanics.
- `skills/do-work/scripts/run-blocked-check.sh` (new) — provides timeout/gtimeout and stock-shell fallback behavior.
- `skills/do-work/scripts/show-commit-diff.sh` (new) — displays ordinary and merge commits through the correct Git form.
- `skills/do-work/scripts/stage-exact-deletion.sh` (new) — stages only the exact cached deletion metadata path.

Builder commits `5129388ccab94a65c4983763b350573b3eb4e08c` and remediation `13aace61a304a64534fc415a14c77818140ef546` were integrated by merge commits `959e079` and `5a18faf` over exact first-pre/latest-merge range `184d4fa..5a18faf`.

## Qualification

- Mechanical qualification uses exact merge range `184d4fa..5a18faf`; all 34 changed paths match the exact Route C Scope/write-set mirror and exist with the claimed state.
- Every current multi-line shell fence has one durable disposition; the 17 promoted rows map to 11 executable homes and named behavior cases, while 21 residue and 2 Go-owned rows retain owners.
- Conflict reconciliation preserved REQ-170's independent capture hook and producer ratchets, selected the new screenshot caller/staged contract, and retained the exact `0.186.30` unique-copy/two-writer semantics in executable fixtures.
- Shipped Markdown drops 17 multi-line fences and 303 nonblank body lines; scripts/tests replace behavior with executable, attributable coverage rather than copied prose.

## Testing

- `bash _dev/tests/prescribed-shell-scripts-behavior.sh`: pass, 11 named cases.
- `bash _dev/tests/prescribed-shell-canonicalization.sh`: pass.
- `_dev/tests/action-shell-blocks.sh` self-test/full/no-ShellCheck paths: pass (59 fences, 26 shipped shell files).
- Staged package, shipped-reference, defensive-surface, updater, installer, and aggregate contract suites: pass on the builder branch.
- Integrated conflict-resolution checks: behavior, canonicalization, staged package, shipped-reference, Bash syntax, warning-level ShellCheck, and cached diff check all pass.
- Negative proof: regressing atomic publication or merge-aware display made its named fixture fail before restoration.
- Literal-path deletion RED/GREEN: the old helper returned 2 after staging both deleted `:(glob)*` and `other.txt`; the remediated helper returns 0 and stages exactly the named magic-looking path.

**Surface accounting:** 315 physical / 303 nonblank Markdown shell-body lines removed; 563 executable lines across 11 scripts plus the fixture scaffold/behavior probes; final remediated implementation is 34 files, 1024 insertions and 563 deletions on the branch.

## Review Remediation

- Initial independent review: Request changes, 50% capped, Acceptance Fail.
- **Important closure:** `stage-exact-deletion.sh` now treats every user-supplied filename as a literal Git pathspec for cached inspection, staging, and post-stage verification; the new two-deletion fixture proves no adjacent index contamination.
- **Traceability closure:** `## Scope` now lists the exact same 34 paths as the merge and Implementation Summary; `scope-drift.sh` passes.
- Focused behavior, canonicalization, staged/shipped/install/aggregate contracts, Bash syntax, warning-level ShellCheck, diff, scope, and debug checks pass after remediation.

## Review

**Final verdict:** Approve — 100%, Acceptance Pass.

- The initial review correctly failed a Git pathspec-magic index-contamination case and an inexact Route C Scope mirror.
- Remediation applies literal pathspec semantics to every exact-deletion Git boundary and adds the exact two-deletion RED/GREEN fixture.
- Scope, `write_set`, Implementation Summary, and `184d4fa..5a18faf` are exact 34-path set-equals.
- All focused, package, installer, lint, executable-mode, producer-closure, and aggregate contract checks pass; no findings remain.
- Durable report with the preserved failing review: `do-work/runs/work-2026-08-11-225637/REQ-171-review.md`.

## Lessons Learned

- Shell quoting and Git path literalness are different layers: quotes stop the shell, but Git pathspec magic still needs an explicit literal-path boundary.
- Promotion inventories must mirror exact changed paths in both Scope and `write_set`; grouped glob prose is useful explanation but not an auditable boundary.
- Migrating inline logic safely means moving the original adversarial fixture with it—the coordinated screenshot race was the evidence that made the conflict resolution deterministic.

## Orientation

- The executable ownership map is `skills/do-work/docs/prescribed-shell-primitives.md`; the durable 17/21/2 census lives in `decisions/audits/2026-08-11-prescribed-shell-primitives.md`.
- Start behavioral changes in `_dev/tests/prescribed-shell-scripts-behavior.sh`, then keep canonicalization, staged-package, and action-shell lint green.

## Knowledge Handoff

- `kb_status: pending`
- `kb_entry: Shell quoting does not disable Git pathspec magic at exact-path boundaries`
