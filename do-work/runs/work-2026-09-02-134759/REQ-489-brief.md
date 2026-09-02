# Builder Brief — REQ-489

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-489-remove-whole-checkpoint-entries-on-departure`
Branch / operative name: `worktree-agent-REQ-489-remove-whole-checkpoint-entries-on-departure`
Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-02-134759/REQ-489-handback.md`
Route: A, focused bug fix. TDD is required.

## Request

Fix `checkpointWithoutClaim` in `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go` so removing this checkout's matching `- REQ-NNN: ... writer: ...` checkpoint entry also removes all immediately following indented continuation lines, stopping at the next `- ` entry, blank section boundary, or heading. Bare one-line entries must continue to work and foreign-label entries plus continuations must remain untouched.

Also make `checkpointWithClaim`/`appendSectionEntry` and removal locate `## In Progress (interrupted)` by a whole heading line. A backticked or inline mention elsewhere must neither attract insertion nor removal.

Captured RED cases:

1. An enriched own-label entry with `Last known state:`, `Key files being modified:`, and `Known issues:` continuation lines currently leaves those lines orphaned after departure.
2. A Session Notes bullet containing the backticked heading text before the real section currently attracts a new claim entry; insertion/removal must use the real heading and preserve the bullet.

## Scope

- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go`
- Existing request-state Go test files needed for the focused regression, expected `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go`

Do not touch any path under `do-work/` in the worktree. Do not touch `CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, or `skills/do-work/actions/version.md`; those are integrator-only. No unrelated cleanup.

## Required context and rules

Read before editing:

- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (whole satellite, touch-required; record this evidence)
- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/crew-members/backend.md`
- `skills/do-work/crew-members/testing.md`
- `skills/do-work/specs/bug-fix.md`

Follow RED → GREEN → REFACTOR. Commit the verified implementation on your branch. Do not bump versions or edit the changelog.

## Required verification

- Focused request-state tests proving both captured RED cases.
- `go test -count=1 ./internal/requeststate` from `skills/do-work/tools/do-work-cli`.
- `go vet ./...` from the module when practical.
- `git diff --stat`, `git diff --check`, file-by-file diff review, and no debug artifacts.

## Hand-back format

Write the complete result to the absolute hand-back path using `apply_patch`. Include:

- branch name and commit hash;
- P-A-U evidence for PLAN/APPLY/UNIFY;
- exact modified/created/deleted file manifest;
- RED command/failure and GREEN command/pass;
- all tests and exit results;
- required-lessons evidence;
- integration seams (normally none here);
- `## Decisions` and `## Discovered Tasks` as separate headings, writing `None.` when empty.

Return only one short status line after the hand-back file is durable.
