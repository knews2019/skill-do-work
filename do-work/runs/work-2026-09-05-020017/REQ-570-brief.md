# Builder brief — REQ-570

You are the implementation builder for one request in the do-work skill repository. You work ONLY inside your own worktree and commit ONLY on your own branch. The orchestrator merges.

## Where you work

- Worktree (your working directory for every command): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed`
- Branch (already checked out there): `worktree-agent-REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed`
- Go module inside the worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed/skills/do-work/tools/do-work-cli` (run go commands from there).
- Main tree (READ ONLY, except the one hand-back file below): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`

## Read first, in this order (absolute paths; read from the MAIN tree, never from the worktree's stale do-work/ snapshot)

1. The request: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed.md` — its `## What`, `## Detailed Requirements`, `## Constraints`, `## Builder Guidance`, `## Red-Green Proof`, `## Plan` (validated; follow it), `## Exploration`, and `## Scope` (your write boundary).
2. The plan in full: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-020017/REQ-570-plan.md`; the exploration in full: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-020017/REQ-570-exploration.md`.
3. Crew rules (in the worktree, they are shipped files): `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md` (tdd: true), `maintenance.md` (maintenance: true — delete before you add).
4. Prime files: `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, and each prime's lessons satellite when you touch code its Read-first or Traps name (`skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `_dev/primes/lessons-action-files.md`, `_dev/primes/lessons-shell-commands.md`).

Required lessons: the REQ's `## Required Lessons — Dropped for Budget` lists three satellites dropped for budget at capture; the touch-conditional rule still applies, so read the satellite of each prime whose code you change. Report in the hand-back which you read.

## Decisions already made (do not reopen)

- The heavy hold is a phase of a `claimed` request (section `## Heavy Verification Plan` + `commit:` on main). No new status, no new checkpoint count, no manual lane path (clarify Step 2.5 is deleted).
- **D-01 (accepted by the orchestrator):** plan finding F3 — the `claim` transition strips `commit`, `heavy_verified_at`, `heavy_verified_revision`; the red drain prose deletes `commit:` and the plan section before remediation.
- F4: no contract predicate names the status; `_dev/tests/contracts/core-checks.sh` is out of scope (removed from the write set).
- Lessons satellites, `do-work/lessons-index.md`, `VERSION` and changelog mirrors are NOT yours: the orchestrator writes them at finalization. Do not touch any `do-work/` path in the worktree.
- The board (`skills/do-work-board/**`) is REQ-571, not yours.

## How to work

- **TDD is mandatory.** Write the RED tests from the plan's Testing section (a) to (e) plus the D-01 requeststate test FIRST, run them, record the exact failures (test names and messages), then implement. Delete or rewrite the old tests the plan lists (cross-REQ test-break rule: intentional, name each in the hand-back).
- Follow the plan's compile-safe groups 1–7, then the prose edits (group 8). Run `go build ./... && go vet ./...` after each group.
- **Write only inside `## Scope`'s file list.** If a file outside it is required by the REQ's own requirements, flag it in the hand-back and proceed; anything else is a scope expansion: stop and report it in the hand-back instead of editing.
- P-A-U: you cannot edit the REQ file (main tree). Record your [PLAN], [APPLY], [UNIFY] evidence in the hand-back under `## P-A-U`. UNIFY = `git diff --stat` reviewed, `gofmt -l` clean, `go vet ./...` clean, no debug artifacts, list of files verified.
- Log significant choices not dictated by the plan as `## Decisions` entries D-02, D-03, … with reasoning (DECIDE & STATE vs ESCALATE with Value/Risk).
- Out-of-scope finds go to `## Discovered Tasks` in the hand-back; do not fix inline.
- Commit on your branch in small coherent commits with messages prefixed `[REQ-570] `. Do not push. Do not merge. Do not touch main.

## Verification before hand-back (all from the worktree)

```
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed/skills/do-work/tools/do-work-cli && gofmt -l . && go build ./... && go vet ./... && go test -count=1 ./...
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed && bash _dev/tests/contract-regressions.sh
cd /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed && grep -rn 'pending-heavy-testing\|ResumePhase\|resume_phase\|HeavyTestingEvidence\|HeavyLaneResult\|matchingHeavyReviewPhase\|ANSWER-HEAVY' skills/do-work _dev --include='*.go' --include='*.md' --include='*.sh' | grep -v CHANGELOG
```

The grep must print nothing. Record each command's exit status and wall time.

## Hand-back (the ONE main-tree path you may write)

Write `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-020017/REQ-570-handback.md` (absolute path; never stage or commit it). Sections, each under its own `##` heading:

- `## Branch` — branch name and final commit hash (`git rev-parse HEAD` in the worktree).
- `## File manifest` — every file created/modified/deleted with the verb, tests included.
- `## P-A-U` — the three phases with evidence.
- `## Red-green evidence` — per RED test: name, failure before (exact message), pass after.
- `## Tests deleted or rewritten` — each old test name and why.
- `## Verification` — each command, exit status, wall seconds.
- `## Decisions` — D-02… entries.
- `## Discovered Tasks` — or "None".
- `## Integration seams` — shared-file one-liners the orchestrator must apply, or "None".
- `## Lesson evidence` — satellites read (whole) and any missing entry.
- `## Blockers` — anything you could not finish, or "None".

Return to the orchestrator ONLY one line: `handback written: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-020017/REQ-570-handback.md <commit-hash>`.
