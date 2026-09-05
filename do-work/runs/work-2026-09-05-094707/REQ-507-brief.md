# Builder brief — REQ-507 (hand the archive and commit tails to `finalize`), remediation pass

## Why this pass exists

REQ-507 was implemented and merged on 4 September as `8e3dbf01..ad8bceb7` (12 files: `work.md` Steps 8/9 cut to judgment, `work-reference.md` procedures reduced, `core-checks.sh` predicates, `advance` composing `finalize` through `finalization_gate.go`, result-model text rendering, CLI prime). It was qualified, tested and independently reviewed (Partial, 82.5%, two noncritical findings) at that revision, then held for heavy lanes. One lane (staged-skills) failed on an assertion that REQ-547 later repaired. Under the current Step 7.7 rule a red lane returns the REQ to remediation.

Since that merge, 8 of the 12 files were reworked by later REQs (REQ-504, 505, 506, 510, 515, 547, 562, 564, 569, 570). The saved-range resume proof therefore reports drift, so the merged work may not be reused blindly. This pass re-verifies REQ-507's acceptance against current `main` and implements only what is missing.

## Task

1. Read the REQ (`## What`, `## Detailed Requirements`, `## Constraints`, `## Scope` → Acceptance, `## Red-Green Proof`) and `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, plus the crew files the orchestrator names.
2. For each acceptance criterion, check current `main` and record evidence (test name + result, or the exact prose lines present/absent):
   - A reviewed and oriented working request advances to a mechanical `finalize` phase, requires exactly one request-bound manifest, and returns the finalizer's typed outcome/findings/changes/rollback/ordered records (`internal/lifecycleadvance/finalization_gate.go` + `_test.go`).
   - Public CLI tests prove serial, supplied-worktree, completed-with-issues, already-green/no-release, missing-input, hostile-input, and identity-mismatch behavior, with refusals producing no mutation.
   - `work.md` Step 8/9 and the `work-reference.md` procedures retain Fold-First, consolidation, impact, terminal/failure, release-content, lesson, and cleanup judgment but no longer teach archive, staging, commit, provenance, or verification mechanics.
   - Structural contracts (`_dev/tests/contracts/core-checks.sh`) and the CLI prime enforce the ownership boundary.
3. Run the focused tests: `go test -count=1 ./internal/lifecycleadvance ./internal/finalization ./internal/resultmodel` inside `skills/do-work/tools/do-work-cli`, and `bash _dev/tests/contracts/core-checks.sh` from the repository root.
4. If a criterion no longer holds, implement the gap on your branch, tests first (RED then GREEN), inside the REQ's Scope file list. If every criterion holds, make NO source change.
5. Hand back `REQ-507-handback.md` in the run directory: the per-criterion evidence table, the file manifest (possibly empty), test commands with results, `## Decisions` (D-04 onward), `## Discovered Tasks`, and lesson evidence. Never write `do-work/` paths or the main tree; never commit on `main`.

## Dispatch details

- Worktree (your working directory, already created at `1012e5e2`): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/.git/work-run-20260905/worktree-agent-REQ-507-hand-archive-and-commit-tails-to-finalize`
- Branch: `worktree-agent-REQ-507-hand-archive-and-commit-tails-to-finalize` (commit here only; prefix every commit message with `[REQ-507] `).
- Hand-back file (the ONLY main-tree path you may write, absolute): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-094707/REQ-507-handback.md`
- Never touch: anything under `do-work/` in any tree, `main`, other worktrees, `CHANGELOG.md`, `VERSION`, release mirrors.
- Rules to read first (all under `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/`): `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`, `maintenance.md`. Prime files: `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and their `lessons-*.md` satellites where present.
- The REQ file to read (main tree, read-only): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-507-hand-archive-and-commit-tails-to-finalize.md`
