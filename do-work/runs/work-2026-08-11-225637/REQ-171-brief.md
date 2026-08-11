# Builder Brief — REQ-171

## Dispatch

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-171-promote-prescribed-shell-to-tested-scripts`
- Branch / operative name: `worktree-agent-REQ-171-promote-prescribed-shell-to-tested-scripts`
- Durable hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-11-225637/REQ-171-handback.md`
- Route: C; domain: general; `tdd: false`, but captured behavioral RED/GREEN proof is mandatory.

## Outcome

Promote every reusable multi-line prescribed-shell primitive into a shipped, fixture-tested script. Keep `_dev/tests/action-shell-blocks.sh` for inline residue/direct script lint, make the canonical guide point to scripts as normative homes, preserve one intent sentence plus an invocation at each caller, and prove semantic traps by executing scripts in fixture repositories. Go-owned atomic REQ reservation gets no shell twin. Shipped Markdown shell lines must decrease.

Parent REQ-165 already supplied syntax/ShellCheck coverage and remains intact; this REQ adds execution-level semantics.

## Inventory gate

The independently verified starting census is 59 shell fences, 40 multi-line fences, 515 physical multi-line body lines (470 nonblank). Dispositions: 17 promote, 21 inline residue, 2 non-shell Go owners. Before editing, reproduce the census and extend `decisions/audits/2026-08-11-prescribed-shell-primitives.md` with durable coordinates (file + heading + content label), disposition, executable home, and behavior case. Do not silently omit or promote beyond the inventory.

Minimum executable set (11 scripts):

Core:
- `skills/do-work/scripts/show-commit-diff.sh`
- `skills/do-work/scripts/add-local-git-exclude.sh`
- `skills/do-work/scripts/atomic-download.sh`
- `skills/do-work/scripts/capture-screenshot.sh`
- `skills/do-work/scripts/run-blocked-check.sh`
- `skills/do-work/scripts/protected-inventory.sh`
- `skills/do-work/scripts/stage-exact-deletion.sh`

Knowledge:
- `skills/do-work-knowledge/scripts/lexical-memory-recall.sh`
- `skills/do-work-knowledge/scripts/install-memory-hooks.sh`

Toolbox:
- `skills/do-work-toolbox/scripts/generate-report-image.sh`
- `skills/do-work-toolbox/scripts/install-last30days.sh`

Reuse `skills/do-work/tools/checks/uncommitted-inventory.sh`, `associate-files.sh`, and `record-commit-hash.sh`; wrappers own orchestration, not duplicate low-level logic. Keep `install-do-work-suite.sh` self-contained as a bootstrap exception.

## Exact promotion surfaces

- Merge-aware commit display: guide plus review/report/presentation consumers → `show-commit-diff.sh`.
- Local Git exclude: guide, board static mode, memory/background guidance, last30days → `add-local-git-exclude.sh`.
- Atomic download: guide, toolbox downloads, updater if its fixture is updated → `atomic-download.sh`.
- Capture screenshot lifecycle → `capture-screenshot.sh`.
- Work blocked-condition timeout runner → `run-blocked-check.sh`.
- Commit/inspect inventory start+associate fences → `protected-inventory.sh`.
- Commit exact deletion guard → `stage-exact-deletion.sh`.
- Knowledge lexical recall and hook merge → their two knowledge scripts.
- Toolbox report-image generator and last30days install/repair/verify → their two toolbox scripts.

Inline residue includes assignments, independent example commands, usage-ledger append, UI command recipes, commit/staging heredoc examples, existing record-hash sequence, and path-only `git diff-tree`. Go reservation and `queue-kanban verify` remain Go-owned.

## Required context

Read completely before editing:

- `CLAUDE.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/background-agents.md`
- `skills/do-work/docs/prescribed-shell-primitives.md`
- `decisions/audits/2026-08-11-prescribed-shell-primitives.md`
- `decisions/audits/2026-08-11-defensive-surface.md`
- `_dev/tests/action-shell-blocks.sh`
- `_dev/tests/prescribed-shell-canonicalization.sh`
- `_dev/tests/session-start-hook-behavior.sh`
- relevant current call sites and contract assertions named below

Do not read or write the worktree's `do-work/` snapshot.

## Dirty-main reconciliation (load-bearing)

The clean worktree does not contain the main tree's current screenshot safety changes. Before implementing screenshot promotion, run this read-only inspection:

`git -C /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 diff -- skills/do-work/actions/capture.md _dev/tests/staged-skills-contract.sh`

Port the exact semantics, not the dirty hunks:

- allocate a unique adjacent `mktemp "${destination}.copying.XXXXXX"` per dispatch;
- copy and verify that exact private temp, install by no-clobber hard link;
- preserve the staged loser/existing destination on collision and clean only dispatch-owned temp;
- retain best-effort post-publication staged-source/directory cleanup;
- carry the coordinated two-writer race into `prescribed-shell-scripts-behavior.sh`.

Do not edit the main tree. The orchestrator will stash/reconcile overlapping dirty paths and preserve the separate 0.186.30 release metadata.

## Scope

The write boundary is the exact `write_set` recorded in the main REQ, summarized here:

- the two audits and canonical guide;
- the 11 scripts above;
- core actions `capture.md`, `commit.md`, `review-work.md`, `work.md`; board `board.md`; knowledge `memory-reference.md`/`setup-memory.md`; toolbox `ai-report.md`/`ai-report-reference.md`/`inspect.md`/`install.md`/`present-work.md`; three package background-agent files; `do-work-update.sh` only if adopting the helper;
- `_dev/tests/fixture-repo.sh`, `prescribed-shell-scripts-behavior.sh`, `prescribed-shell-canonicalization.sh`, `contract-regressions.sh`, `staged-skills-contract.sh`, `action-shell-blocks.sh`, `update-script-behavior.sh`, and optional `session-start-hook-behavior.sh` only when their owned contract moves.

Never touch version/changelog surfaces, `suite/modules.tsv`, `install-do-work-suite.sh`, Go code, unrelated dirty installer/helper/prime surfaces, or any other file. Stop/report before any scope extension.

## Script and test contracts

- Bash 3.2/POSIX utility floor; no associative arrays or `mapfile`; tested fallback for GNU-only tools.
- Raw/imported values enter as separately quoted arguments or stdin; never `eval` or constructed shell.
- Mechanics in scripts; consent/trust/workflow/report policy in callers. Stdout is result data; diagnostics stderr; usage/environment errors 2; preserve meaningful underlying statuses (blocked timeout 124).
- `git rev-parse --git-path` owns Git-internal paths. Private scratch uses `mktemp`+trap; cross-block state uses deterministic re-derived paths. Cleanup never masks failure.
- Add `_dev/tests/fixture-repo.sh` as sourceable helpers only and `_dev/tests/prescribed-shell-scripts-behavior.sh` with one named trap case per script and attributable `FAIL:` diagnostics.
- Explicitly wire the new behavior probe into `contract-regressions.sh`; keep all eight canonical-guide headings; require scripts executable and reject old promoted implementations without line-number coupling.
- Add all new scripts to the defensive-surface audit; staged/install tests must prove sibling paths resolve.

Required negative proofs in an expendable state: regressing atomic download to direct final publication must fail the named atomic case; regressing merge display to plain `git show` must fail the real-merge case. Restore before committing.

## Verification

Run focused behavior/canonicalization after each promotion, then:

- `bash -n` and ShellCheck warning severity over all new/modified shell files;
- `bash _dev/tests/prescribed-shell-scripts-behavior.sh`;
- prescribed-shell canonicalization;
- action-shell self-test, normal scan, and no-ShellCheck degradation;
- SessionStart if fixture-migrated;
- staged-skills, shipped-package-reference, defensive-surface, updater, and install-suite probes when touched;
- full contract regressions;
- `git diff --check`, debug-artifact/changed-path review, and final fence census/net-surface accounting.

## Builder rules and hand-back

Use `apply_patch`, commit on the operative branch, never rebase, never touch release metadata. Report P-A-U evidence, every file with action verb, per-script cases/results, negative proof, before/after prose counts, script/test lines added, D-XX decisions, integration seams, and discovered tasks/blockers in the exact hand-back path. Return only one line after the durable file is complete.
