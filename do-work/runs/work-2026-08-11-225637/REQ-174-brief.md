# Builder Brief — REQ-174

## Dispatch

- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-174-validate-root-markdown-fence-info`
- Branch / operative name: `worktree-agent-REQ-174-validate-root-markdown-fence-info`
- Durable hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-11-225637/REQ-174-handback.md`
- Route: A (focused bug fix)
- Domain: testing
- TDD: true

## Request

Make root and list fence classification in `_dev/tests/shipped-package-reference-contract.sh` share the CommonMark rule that a backtick-fence info string cannot contain a backtick. Reject backticks in backtick-fence info strings, preserve tilde-fence behavior, consolidate the existing list validation, and add the reproduced root-level Goldmark-differential fixture. Preserve the behavior earned by REQ-150 and REQ-163.

Captured RED case: classify a root backtick fence whose info is `lang`invalid`, followed by `[live](visible.md)`. The current classifier masks the link; the pinned Goldmark renderer publishes `visible.md`. GREEN means the classifier returns `visible.md`, agrees with Goldmark, and existing root/list/tilde fixtures pass.

## Required context

Read before editing:

- `CLAUDE.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/testing.md`
- `skills/do-work/specs/bug-fix.md`
- `skills/do-work/tools/prime-do-work-update.md`

Prime lessons already reviewed by the orchestrator: Markdown classification requires lexical plus container state, direct syntax adjacency, differential fixtures against the authoritative renderer, and preservation of first-line/BOM byte semantics. Re-read locally linked archived lessons that exist in the worktree if needed.

## Scope and constraints

- Write boundary: `_dev/tests/shipped-package-reference-contract.sh` only.
- Do not touch `do-work/`, version files, changelogs, or any main-tree path other than the exact hand-back file above.
- Use RED → GREEN → REFACTOR and record the exact failing assertion/output before the fix and passing output after it.
- Identify the root cause, scan the same primitive for analogous cases, and report out-of-scope discoveries instead of fixing them.
- Use `apply_patch` for edits. Commit the implementation on the operative branch; do not rebase.

## P-A-U and verification

Report these phases in the hand-back because queue state remains in the main tree:

- PLAN: brief technical approach and RED command.
- APPLY: minimal implementation and fixture.
- UNIFY: `git diff --stat`, full changed-file review, targeted test command(s), `git diff --check`, debug-artifact scan, and exact files checked.

## Hand-back format

Write the durable hand-back file with:

- Branch and commit hash
- PLAN / APPLY / UNIFY evidence
- RED output and GREEN output
- File manifest with action verbs
- Tests run and results
- Root cause note
- Decisions (D-XX with reasoning; Value/Risk for escalations)
- Integration seams (exact line/location, or `None`)
- Discovered tasks / blockers

Return only a one-line status after the file is written.
