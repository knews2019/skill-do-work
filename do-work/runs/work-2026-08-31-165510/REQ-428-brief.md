# Builder Brief — REQ-428

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-428-preserve-filename-only-collision-evidence-in-dependency-graphs`
Branch: `worktree-agent-REQ-428-preserve-filename-only-collision-evidence-in-dependency-graphs`
Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-428-handback.md`

## Ownership

Work only in the worktree above and commit implementation changes on its branch. Treat every path under `do-work/` in the worktree as a stale snapshot: do not read or write it. The only main-tree path you may write is the exact hand-back path above. Do not change `VERSION`, either changelog, or the current-version line; those are integrator-owned.

## Request

Preserve filename-only collision evidence in repository dependency graphs. Two filenames can claim `REQ-021` while their frontmatter IDs are `REQ-030` and `REQ-031`; a dependency on `REQ-021` must report an ambiguous target rather than a merely missing target, while staying blocked and unresolved for readiness and depth calculations.

Requirements:

- Check collision evidence before classifying an absent node as missing.
- Preserve deterministic unmet and ambiguity evidence for filename-only collisions.
- Keep the target blocked and unresolved for readiness and depth calculations.

RED/GREEN case: build a snapshot with `REQ-021` claimed by two filenames whose frontmatter IDs are `REQ-030` and `REQ-031`, depend on `REQ-021`, and prove the graph reports exact ambiguous-collision evidence rather than missing-target evidence.

Route A, `tdd: true`, `domain: backend`. Follow RED, then GREEN, then refactor.

## Required rules and context

Read these files from the worktree before editing:

- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/crew-members/backend.md`
- `skills/do-work/crew-members/testing.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
- any lessons satellite the prime routes you to for files you touch

## Hand-back format

Commit the implementation on your branch. Then write the exact hand-back path with:

- branch name and commit hash
- concise technical approach and any decisions
- RED command and the assertion failure observed before production changes
- GREEN command and direct exit status after the fix
- every modified, added, or deleted project file with a factual one-line description
- integration seams, if any
- `## Discovered Tasks` only when needed
- `## Decisions` only when needed

Return only a one-line status to the orchestrator after the hand-back exists.
