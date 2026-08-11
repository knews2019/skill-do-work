# REQ-171 Exploration

## Boundary and Current Baseline

Discovery stayed inside the bounded surface named by the REQ and plan: REQ-171, archived parent REQ-165, the REQ-171 plan, the prescribed-shell guide/audits/ratchets, current shipped Markdown shell fences and shipped shell sources, the REQ-166 fixture probe, aggregate test wiring, and tests that pin behavior being promoted. `prime_files` is empty, so no global architecture was loaded.

Current integration state is not clean:

- `main` is at `59b3067` and is 12 commits ahead of `origin/main`.
- The main worktree contains an uncommitted 0.186.30 screenshot-safety change in `skills/do-work/actions/capture.md` and `_dev/tests/staged-skills-contract.sh`, plus its release edits in `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, and `skills/do-work/actions/version.md`.
- Two sibling worktrees, for REQ-170 and REQ-174, currently start at the same `59b3067` commit. A new REQ-171 worktree made from that commit will not inherit the uncommitted main-worktree changes.
- Other dirty files are `_dev/tests/install-suite-behavior.sh`, `skills/do-work/SKILL.md`, `skills/do-work/tools/prime-do-work-update.md`, both copies of `replace-text-section.sh`, and the queue/working/run bookkeeping for the active batch. Those changes are user/integrator-owned and must be preserved.

Using the exact fence boundary from `_dev/tests/action-shell-blocks.sh` (zero to three leading spaces, info string exactly `bash` or `sh`), the current shipped `skills/` tree has:

- 59 shell fences total and 15 shipped `.sh` sources.
- 40 multi-line fences, defined as a fence body containing at least two physical lines; they contain 515 physical body lines, 470 of them nonblank.
- 19 one-line fences. Across all 59 fences there are 534 physical body lines and 489 nonblank lines.
- The clean current probes report: `action-shell-blocks.sh --self-test` passes; the normal scan passes with “59 fenced blocks and 15 shipped shell files”; prescribed-shell canonicalization passes; and the REQ-166 SessionStart behavior probe passes.

The final audit should use file + enclosing heading + a short content label as its durable coordinate. Current line numbers below are deliberately omitted because the implementation will move them.

## Exact Current Multi-Line Fence Inventory

Disposition counts are 17 `Promote`, 21 `Inline residue`, 0 `Already executable`, and 2 `Non-shell owner`. “Already executable” still matters for related behavior outside these 40 fences: the complete per-file inventory, association, and commit-hash mechanics already live in tested scripts, but no complete multi-line fence is merely their invocation.

| # | Durable coordinate / block identity | Lines (body/nonblank) | Disposition | Target or rationale |
|---:|---|---:|---|---|
| 1 | `skills/do-work-board/actions/board.md` → `Step 5: Run the selected mode` → static-mode root/local-exclude block | 7/7 | Promote | Local-exclude portion moves to core `add-local-git-exclude.sh`; the caller keeps root/mode policy and the queue-kanban invocation. |
| 2 | `skills/do-work-knowledge/actions/memory-reference.md` → `Lexical Recall (Layer 1 — always runs)` | 16/16 | Promote | `skills/do-work-knowledge/scripts/lexical-memory-recall.sh`. The current body is pseudocode, which is precisely why behavior is untested. |
| 3 | same file → `Semantic Recall (Layer 2 — optional, detected)` → three backend-detect commands | 3/3 | Inline residue | Independent illustrative probes, no shared transaction or recovery logic. |
| 4 | same file → `Usage-Ledger Contract (canonical — both engines)` → one continued `printf` append | 2/2 | Inline residue | One best-effort command; local ledger policy remains prose. |
| 5 | same file → `Hook Install Internals (used by actions/setup-memory.md → memory-module)` | 19/19 | Promote | `skills/do-work-knowledge/scripts/install-memory-hooks.sh`. |
| 6 | `skills/do-work-knowledge/actions/memory.md` → `Input` → project/memory assignments | 2/2 | Inline residue | Simple deterministic assignments. |
| 7 | `skills/do-work-knowledge/actions/setup-memory.md` → `Preconditions` → project/package-root assignments | 2/2 | Inline residue | Simple placeholders/assignments. |
| 8 | `skills/do-work-knowledge/docs/interview-guide.md` → `Integration with bkb` → ingest/triage/ingest example | 3/3 | Inline residue | Independent CLI recipe. |
| 9 | same coordinate → two query examples | 2/2 | Inline residue | Independent examples. |
| 10 | `skills/do-work-toolbox/actions/ai-report-reference.md` → `Image Generation Backend (Steps 3c, 4d, 5)` → `STYLE` constant | 3/3 | Inline residue | Caller-owned illustrative style text, not executable control flow. |
| 11 | same coordinate → `gen_image()` definition | 27/25 | Promote | `skills/do-work-toolbox/scripts/generate-report-image.sh`. |
| 12 | same coordinate → parallel two-image invocation/verification | 7/7 | Inline residue | Dynamic per-report orchestration may remain a short sequence calling the promoted script. |
| 13 | `skills/do-work-toolbox/actions/ai-report.md` → `Step 7: Render and Judge` → local HTTP server recipe | 2/2 | Inline residue | One command plus an operator instruction. |
| 14 | `skills/do-work-toolbox/actions/deep-explore.md` → `Create the directory` | 2/2 | Inline residue | Assignment plus one `mkdir`; no reusable branch/recovery contract. |
| 15 | `skills/do-work-toolbox/actions/inspect.md` → `Step 1: Preflight` | 19/18 | Promote | Core `protected-inventory.sh start`, retaining inspect-specific quarantine naming/policy. |
| 16 | same file → `Unscoped mode (default)` | 31/28 | Promote | Core `protected-inventory.sh associate`. |
| 17 | `skills/do-work-toolbox/actions/install.md` → Bowser `Phase 1: Check if already installed` | 3/3 | Inline residue | Independent detect/status commands. |
| 18 | same file → Bowser `Phase 4: Install the Bowser skill` | 7/7 | Promote | Download mechanics move to core `atomic-download.sh`; destination/consent policy stays in the action. |
| 19 | same file → Bowser `Phase 5: Verify` | 3/3 | Inline residue | Independent verification commands. |
| 20 | same file → last30days `Phase 2: Vendor the skill` → clone/copy/cleanup | 10/10 | Promote | `skills/do-work-toolbox/scripts/install-last30days.sh`. |
| 21 | same coordinate → `.git/info/exclude` append | 5/5 | Promote | Fold into `install-last30days.sh`, which calls core `add-local-git-exclude.sh`. |
| 22 | same file → last30days `Phase 3: Verify` | 14/14 | Promote | Fold into `install-last30days.sh` so install/repair/verify is one executable transaction. |
| 23 | `skills/do-work-toolbox/actions/ui-review.md` → `Step 8.5: Visual Verification` → mobile/tablet/desktop commands | 14/12 | Inline residue | UI-review command recipe. |
| 24 | same coordinate → accessibility snapshot commands | 4/4 | Inline residue | UI-review command recipe. |
| 25 | `skills/do-work/actions/capture.md` → `File Naming` → optional `queue-kanban next-req` accelerator | 3/3 | Non-shell owner | Atomic REQ reservation remains owned by the Go board tool; do not create a shell twin. |
| 26 | same file → `Step 4: Handle Screenshots` | 35/34 | Promote | Core `capture-screenshot.sh`, with the dirty-main unique-temp collision fix preserved. |
| 27 | same file → `Step 7: Commit (Git repos only)` | 18/13 | Inline residue | Explicit staging plus a commit-message heredoc template. |
| 28 | `skills/do-work/actions/cleanup.md` → `Commit (Git repos only)` | 29/26 | Inline residue | Illustrative staging variants plus a commit-message heredoc template. |
| 29 | `skills/do-work/actions/commit.md` → `Step 1: Preflight` | 19/18 | Promote | Core `protected-inventory.sh start`. |
| 30 | same file → `Step 3: Associate with REQs` | 31/28 | Promote | Core `protected-inventory.sh associate`. |
| 31 | same file → first `Step 5: Commit` block → exact cached-deletion guard | 36/35 | Promote | Core `stage-exact-deletion.sh`. |
| 32 | same coordinate → REQ-associated commit example | 12/8 | Inline residue | Commit-message heredoc template. |
| 33 | same coordinate → unassociated commit example | 10/7 | Inline residue | Commit-message heredoc template. |
| 34 | `skills/do-work/actions/forensics.md` → `14. Release and Queue Invariants (Go toolchain only)` | 2/2 | Non-shell owner | `queue-kanban verify` remains Go-owned. |
| 35 | `skills/do-work/actions/review-work.md` → `Commit (Standalone mode, git repos only)` | 20/15 | Inline residue | Explicit staging plus commit-message heredoc. |
| 36 | `skills/do-work/actions/work-reference.md` → first `Commit & Metadata-Commit Procedure (Step 9)` block | 39/31 | Inline residue | Staging and commit-message template; caller-specific transaction policy. |
| 37 | same coordinate → record-hash/stage/verify sequence | 15/12 | Inline residue | Short invocations of the already-tested `record-commit-hash.sh` plus caller policy. |
| 38 | `skills/do-work/actions/work.md` → `Step 1: Find Next Request` → bounded blocked-check runner | 32/32 | Promote | Core `run-blocked-check.sh`. |
| 39 | `skills/do-work/docs/prescribed-shell-primitives.md` → `Local Git ignore` | 4/4 | Promote | Guide points to core `add-local-git-exclude.sh` as normative implementation. |
| 40 | same file → `Atomic download publication` | 3/3 | Promote | Guide points to core `atomic-download.sh` as normative implementation. |

The 19 one-line fences remain residue. They are root assignments, single Go/tool invocations, individual install commands, and the path-only `git diff-tree --no-commit-id --name-only -r -m <commit>` default. Merge-aware commit display is currently expressed in prose rather than a multi-line fence, but it contains reusable branching and is explicitly in the approved expected promotion set, so it still qualifies.

## Smallest Viable Promotion Set

The smallest plan-compliant set is 11 shipped scripts: seven core, two knowledge, and two toolbox. Do not add scripts for the path-only diff-tree command, raw-text rule, diff-output filtering, command-block state rule, commit templates, UI recipes, or either Go-owned operation.

| Script | Minimal interface/ownership | Existing surfaces contracted | Required semantic probe |
|---|---|---|---|
| `skills/do-work/scripts/show-commit-diff.sh` | One literal commit argument; stdout is the selected `git show` result; diagnostics on stderr. | Guide plus `review-work.md`, `ai-report.md`, and all three `present-work.md` instruction/template sites. | Ordinary commit patch; real two-parent merge uses first-parent output; invalid revision fails. |
| `skills/do-work/scripts/add-local-git-exclude.sh` | Repository root, cwd-relative probe path, and root-aligned exclude pattern as separately quoted arguments. | Guide, board static mode, setup-memory/background guidance, and toolbox last30days install. | Subdirectory and linked-worktree `--git-path`; one idempotent append; spaces; non-Git behavior. |
| `skills/do-work/scripts/atomic-download.sh` | URL and final path as separately quoted arguments. | Guide, toolbox skill downloads, and `do-work-update.sh` if it is kept on the single-home contract. The bootstrapping `install-do-work-suite.sh` remains the documented exception. | Fake curl writes a nonempty partial then fails: nonzero, no broken publication, `.download` removed; success publishes exact bytes. |
| `skills/do-work/scripts/capture-screenshot.sh` | Exact source and permanent destination; source-cleanup mode must distinguish dispatcher staging from an inline attachment/cache source. | `capture.md` Step 4. | Byte match, no-clobber, unique adjacent temp, coordinated two-writer collision, winning staged-source cleanup, empty dispatch-dir cleanup, losing/source-preservation paths, best-effort post-publication cleanup, and inline-source preservation. |
| `skills/do-work/scripts/run-blocked-check.sh` | Deterministic probe-file path; owns timeout selection/polling and probe cleanup; returns underlying status or 124. | `work.md` Step 1. | Success/failure propagation, GNU timeout, `gtimeout`, forced portable fallback, 124 timeout, cleanup. |
| `skills/do-work/scripts/protected-inventory.sh` | Explicit `start` and `associate` modes around existing `tools/checks/uncommitted-inventory.sh` and `associate-files.sh`; core home because both core commit and toolbox inspect use it. | Four commit/inspect preflight/association fences. | Nested untracked secret reaches inventory; quarantine persists between modes; prior X never reappears readable; empty quarantine keeps safe candidates; missing deterministic quarantine fails closed. |
| `skills/do-work/scripts/stage-exact-deletion.sh` | One literal path, run in its Git worktree. | `commit.md` exact-deletion block. | Exact cached D accepted; unstaged deletion staged/rechecked; rename, multiple-record, unrelated, and non-deletion shapes fail. |
| `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` | Raw query arrives on stdin as data; output is bounded attributed ranked results. | Memory-reference lexical pseudocode. | Hostile quote/command text remains data; filtering; distinct-token score and filename-date recency weights; cap and attribution. |
| `skills/do-work-knowledge/scripts/install-memory-hooks.sh` | Project root and knowledge package root. | Memory-reference hook merge and setup-memory invocation. | Empty/existing/partial settings, independently gated hooks, invalid JSON, missing jq, unrelated-hook preservation, validated temp publication, rollback/backup and verification failure. |
| `skills/do-work-toolbox/scripts/generate-report-image.sh` | Absolute output path and sanitized visual description separately; style can be a fixed script-owned constant or explicit safe data, never raw imported text. | `ai-report-reference.md` generator definition; dynamic multi-image calls remain inline. | Direct backend, missing/empty output, disabled agentic fallback, enabled isolated fallback with fake backend, and temp cleanup. |
| `skills/do-work-toolbox/scripts/install-last30days.sh` | Project root and toolbox/core package roots as needed; owns detect/repair/install/ignore/verify. | Three last30days fences in `install.md`. | Fresh install, partial repair, clone/copy failure, local-ignore failure, Python failure, non-Git mode, idempotent rerun, no project config/API-key write. |

Related existing executable behavior should be reused, not cloned:

- `skills/do-work/tools/checks/uncommitted-inventory.sh` and `associate-files.sh` remain the low-level inventory/association owners. `protected-inventory.sh` only owns the cross-block quarantine orchestration.
- `skills/do-work/tools/checks/record-commit-hash.sh` remains the owner behind the inline metadata sequence.
- `skills/do-work/tools/do-work-update.sh` is already behavior-tested. If it is changed to call `atomic-download.sh`, its hermetic fixture builders in `_dev/tests/update-script-behavior.sh` must install/copy that helper too; otherwise the test install intentionally lacks it.
- `skills/do-work/tools/install-do-work-suite.sh` must stay self-contained: its first bootstrap download occurs before the package containing the helper exists.

## Action and Script Conventions

- Every migrated action keeps one sentence of local intent, then one invocation. Confirmation, trust, tracked-file, fallback/report, and workflow decisions remain in the action; generic scripts own mechanics only.
- Core callers use `<skill-root>/scripts/...`. Board/knowledge/toolbox callers of a core primitive use the explicit sibling path `<skill-root>/../do-work/scripts/...`. Package-specific scripts stay under their owning package.
- Arguments are separate, quoted data. Query/imported text must use stdin where appropriate; no constructed command string or `eval`.
- Stdout is result data a caller may consume. Diagnostics go to stderr. Usage/environment errors use exit 2; semantic/underlying negative statuses are not converted to success. `run-blocked-check.sh` specifically preserves 124.
- Use `#!/usr/bin/env bash` and the established Bash 3.2 floor: no associative arrays or `mapfile`; no GNU-only feature without a tested fallback. Existing shipped checks generally use `set -uo pipefail` when failures are handled explicitly, not blanket `set -e` around expected nonzero branches.
- Temporary state uses `mktemp` and traps when private to one invocation. Anything crossing action blocks must have a deterministic path derived again in each mode. Cleanup must not mask the failure that caused it.
- Git internals come from `git rev-parse --git-path`; paths/URLs/commits are literal separately quoted arguments.
- Every new script must be executable. `suite/modules.tsv` ships whole package directories, so no manifest row is needed for new `scripts/` directories, but staged/install-layout tests must prove references resolve after installation.

## Test Conventions and Existing Ratchets to Reconcile

- Add `_dev/tests/fixture-repo.sh` as a sourceable helper only: deterministic Git identity/init/commit/merge, queue/version/untracked seeding, and cleanup registration; no assertions or standalone product verdicts.
- Add executable `_dev/tests/prescribed-shell-scripts-behavior.sh` with named cases, `FAIL: <script/case> ...` diagnostics, aggregate failure counting, a private fixture root/trap, and a concise pass line. Keep fixtures hermetic by stubbing curl/image backends/timeouts through `PATH`.
- `_dev/tests/*.sh` are not auto-discovered. `_dev/tests/contract-regressions.sh` must guard and invoke the new behavior probe explicitly, as it does for the shell-block, canonicalization, SessionStart, staged-package, updater, and defensive-surface probes.
- `_dev/tests/prescribed-shell-canonicalization.sh` must keep all eight guide headings, require every advertised script to exist and be executable, require the contracted invocation form at each caller, and reject the removed promoted implementations/restatements without volatile line numbers.
- `_dev/tests/action-shell-blocks.sh` stays. Its Markdown scan covers inline residue; its direct `.sh` scan gives every new script full Bash/ShellCheck warnings. Only extracted fences exclude SC2034/SC2154.
- Add all 11 shipped scripts to `decisions/audits/2026-08-11-defensive-surface.md`; `_dev/tests/defensive-surface-audit.sh` dynamically requires every shipped `.sh` to have a same-row disposition.
- `_dev/tests/staged-skills-contract.sh` currently extracts and executes the inline capture fence. Replace that structural pin with the shipped script invocation and behavior-probe ownership. Add new scripts to the appropriate required package lists and installed-path checks.
- Preserve the existing low-level inventory/association fixtures in `contract-regressions.sh`. Update prose-shape assertions that currently require `FILENAME == ARGV[1]`, deterministic quarantine names, exact-deletion wording, `probe_exit=124`, and direct low-level script basenames inside action files so they instead require the promoted script invocation and delegate semantics to the focused behavior cases.
- The current `hardened_check_scripts` loop requires `commit.md`/`inspect.md` themselves to mention `uncommitted-inventory.sh` and `associate-files.sh`. Once those actions invoke `protected-inventory.sh`, point that ratchet at the wrapper and separately assert the wrapper invokes the two existing tools.
- `setup-memory.md` is pinned to retain the backup name and exact three local-exclude patterns. These can remain caller policy/arguments while hook merge mechanics move to the script.
- If REQ-166's probe sources `fixture-repo.sh`, retain its four cases byte-for-byte in behavior: happy path, missing version file, reformatted version line, and missing queue directory. Migration is optional.

## Dirty-Main and Parallel-Reconciliation Hazards

1. **Capture conflict is semantic, not just textual.** Dirty main changes `capture.md` from a fixed `${destination}.copying` file to `mktemp "${destination}.copying.XXXXXX"` and adds a coordinated two-dispatch race fixture. A builder branch from `59b3067` sees the unsafe older form. The promoted script must implement the dirty-main unique-copy algorithm, and the new behavior probe must inherit the race case. Merely replacing the old block from the builder branch would regress the just-completed fix.
2. **`_dev/tests/staged-skills-contract.sh` is an unavoidable overlap.** Main's dirty version adds the race test while REQ-171 must stop extracting the inline fence. Reconcile manually: preserve dispatch-directory structural checks, require the script invocation/executable installed path, and move all executable screenshot lifecycle cases, including the new race, to `prescribed-shell-scripts-behavior.sh` (or have staged-skills invoke that public case without duplicating it).
3. **Release files are serial-only and already dirty at 0.186.30.** Builders must not touch `VERSION`, either installed version surface, or either changelog. After all worktree branches and the existing main edits are reconciled, choose a version strictly above the then-current top changelog entry and keep the changelog mirror byte-identical.
4. **Aggregate-test edits should have one owner.** `contract-regressions.sh`, canonicalization, the inventory audit, defensive audit, fixture helper, and behavior probe are shared seams. Parallel builders should not each add their own invocation or competing interface assertions.
5. **Updater helper adoption expands the test fixture.** If `do-work-update.sh` calls `../scripts/atomic-download.sh`, `_dev/tests/update-script-behavior.sh` must copy that helper into both synthetic installed and upstream suites and keep the “exactly one curl call” assertion. This file is clean now, so make the change in the same branch as atomic download.
6. **Cross-package path verification is integration work.** Core/knowledge/toolbox builders can implement independently, but reconciliation must run staged and installed-package reference tests after all caller/script branches are merged. A branch-local green test cannot prove a sibling script that has not yet been integrated.
7. **The active `do-work/working`/`runs` files are owner bookkeeping.** Builders should not edit or commit them. This exploration file is the only exception assigned to this agent.

## Recommended Reconciliation Order

1. Land/preserve the current dirty screenshot safety change or explicitly carry its diff while merging REQ-171; never resolve `capture.md`/staged tests by choosing the builder's tracked-HEAD side wholesale.
2. Integrate the inventory/canonicalization/fixture branch first, then shared core primitives, then core action-specific scripts, then knowledge/toolbox scripts and callers.
3. Resolve aggregate/canonicalization/staged-test overlaps once, after all script interfaces are final.
4. Re-run the post-migration census and record before/after Markdown shell lines plus added script/test lines. The starting accounting is 59 fences / 534 physical shell-body lines / 489 nonblank, with the qualifying multi-line subset 40 / 515 / 470. The Markdown count must decrease.
5. Run the full verification set, perform the two requested fault-injection proofs on an expendable worktree, then do the single integrator-owned version/changelog update.

## Verification Commands

Run from repository root on the reconciled tree:

```bash
for shell_path in \
  _dev/tests/fixture-repo.sh \
  _dev/tests/prescribed-shell-scripts-behavior.sh \
  _dev/tests/prescribed-shell-canonicalization.sh \
  _dev/tests/action-shell-blocks.sh \
  _dev/tests/contract-regressions.sh \
  _dev/tests/staged-skills-contract.sh \
  _dev/tests/update-script-behavior.sh \
  skills/do-work/tools/do-work-update.sh \
  skills/do-work/scripts/*.sh \
  skills/do-work-knowledge/scripts/*.sh \
  skills/do-work-toolbox/scripts/*.sh
do
  bash -n "$shell_path" || exit 1
done

shellcheck --severity=warning \
  _dev/tests/fixture-repo.sh \
  _dev/tests/prescribed-shell-scripts-behavior.sh \
  _dev/tests/prescribed-shell-canonicalization.sh \
  _dev/tests/action-shell-blocks.sh \
  _dev/tests/contract-regressions.sh \
  _dev/tests/staged-skills-contract.sh \
  _dev/tests/update-script-behavior.sh \
  skills/do-work/tools/do-work-update.sh \
  skills/do-work/scripts/*.sh \
  skills/do-work-knowledge/scripts/*.sh \
  skills/do-work-toolbox/scripts/*.sh

bash _dev/tests/prescribed-shell-scripts-behavior.sh
bash _dev/tests/prescribed-shell-canonicalization.sh
bash _dev/tests/action-shell-blocks.sh --self-test
bash _dev/tests/action-shell-blocks.sh
PATH=/usr/bin:/bin bash _dev/tests/action-shell-blocks.sh
bash _dev/tests/session-start-hook-behavior.sh
bash _dev/tests/staged-skills-contract.sh
bash _dev/tests/shipped-package-reference-contract.sh
bash _dev/tests/defensive-surface-audit.sh
bash _dev/tests/update-script-behavior.sh
bash _dev/tests/contract-regressions.sh
git diff --check
```

Also run `bash _dev/tests/install-suite-behavior.sh` if staged/install fixture code is changed during reconciliation (it is already dirty on main for another active request, so preserve those edits). Re-run the exact zero-to-three-space fence census after caller contraction and report both physical and nonblank before/after counts.

Acceptance is not complete until two negative proofs work in an expendable worktree: temporarily regress `atomic-download.sh` to direct-final publication and confirm the focused test fails naming the atomic-download case; restore it, then temporarily regress `show-commit-diff.sh` to plain `git show` and confirm the real-merge case fails. Restore both and rerun the complete suite green.
