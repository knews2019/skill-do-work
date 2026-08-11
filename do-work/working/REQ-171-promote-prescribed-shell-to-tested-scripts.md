---
id: REQ-171
title: "Addendum: promote prescribed shell primitives to shipped, fixture-tested scripts"
status: claimed
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
  - skills/do-work/tools/do-work-update.sh
  - _dev/tests/fixture-repo.sh
  - _dev/tests/prescribed-shell-scripts-behavior.sh
  - _dev/tests/prescribed-shell-canonicalization.sh
  - _dev/tests/contract-regressions.sh
  - _dev/tests/staged-skills-contract.sh
  - _dev/tests/action-shell-blocks.sh
  - _dev/tests/update-script-behavior.sh
  - _dev/tests/session-start-hook-behavior.sh
---

# Addendum: Promote Prescribed Shell Primitives to Shipped, Fixture-Tested Scripts

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
- `decisions/audits/2026-08-11-prescribed-shell-primitives.md` (modified) — complete promotion inventory and accounting
- `decisions/audits/2026-08-11-defensive-surface.md` (modified) — disposition for every new shipped script
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified) — intent/contracts pointing to executable homes
- `skills/do-work/scripts/*.sh` (7 new files) — shared/core promoted primitives
- `skills/do-work-knowledge/scripts/*.sh` (2 new files) — knowledge-specific promoted primitives
- `skills/do-work-toolbox/scripts/*.sh` (2 new files) — toolbox-specific promoted primitives
- Core/board/knowledge/toolbox action and background-agent call sites listed in `write_set` (modified) — one intent sentence plus invocation
- `_dev/tests/fixture-repo.sh` and `_dev/tests/prescribed-shell-scripts-behavior.sh` (new) — fixture scaffold and execution-level semantic traps
- Existing canonicalization, aggregate, staged-package, shell-lint, updater, and optional SessionStart probes listed in `write_set` (modified only where their owned contract moves)

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
