# Builder Brief — REQ-427

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-427-confirm-go-version-floor`
Branch: `worktree-agent-REQ-427-confirm-go-version-floor`
Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-427-handback.md`

## Ownership

Work only in the worktree above and commit implementation changes on its branch. Treat every path under `do-work/` in the worktree as a stale snapshot: do not read or write it. The only main-tree path you may write is the exact hand-back path above. Do not change `VERSION`, either changelog, or the current-version line; those are integrator-owned.

## Request

Lower the core installer and updater compatibility floor to Go `1.25.0`, the lowest exact toolchain that passes the complete current core suite. Add an exact-Go-1.25 compatibility test lane and update every current core prerequisite restatement together.

Exact-toolchain evidence already established:

- Go 1.23 fails because `os.OpenRoot` and `os.Root` are unavailable.
- Go 1.24 fails because `Root.ReadFile` is unavailable.
- Go 1.25 passes all 16 package suites.

Update the current core restatements named by the durable REQ: `README.md`, `skills/do-work/actions/version.md`, `skills/do-work/tools/do-work-cli.sh`, `skills/do-work/tools/do-work-cli/go.mod`, `skills/do-work/docs/prescribed-shell-primitives.md`, `skills/do-work/tools/prime-do-work-update.md`, and compatibility-launcher comments in the root and shipped mirrors. Keep the optional board module, `_dev/tests/maintainer-verify.sh`, historical changelogs, and archived reports unchanged.

Route A, `tdd: false`, `domain: general`. Keep the change mechanical and focused. Run the exact Go 1.25 core test lane plus relevant launcher and contract checks. Do not run the integrator-only version/changelog ceremony.

## Required rules

Read these files from the worktree before editing:

- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- `_dev/primes/prime-shell-commands.md` because shipped shell is in scope

## Hand-back format

Commit the implementation on your branch. Then write the exact hand-back path with:

- branch name and commit hash
- concise technical approach and any decisions
- every modified, added, or deleted project file with a factual one-line description
- exact commands and direct exit status for tests
- integration seams, if any
- `## Discovered Tasks` only when needed
- `## Decisions` only when needed

Return only a one-line status to the orchestrator after the hand-back exists.
