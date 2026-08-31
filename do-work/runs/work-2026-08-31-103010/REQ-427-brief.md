# Builder Brief — REQ-427

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-427-confirm-go-version-floor`
Branch: `worktree-agent-REQ-427-confirm-go-version-floor`
Hand-back: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-103010/REQ-427-handback.md`
Route: A (focused direct implementation)
TDD: false

## Request

Lower only the core installer/updater Go version floor from 1.26.1 to 1.23.0. The user confirmed this exact choice. The optional board module remains at Go 1.26.1 and is outside scope.

Requirements:

- Change the `go` directive in `skills/do-work/tools/do-work-cli/go.mod` to 1.23.0.
- Change `minimum_go_version` in `skills/do-work/tools/do-work-cli.sh` to 1.23.0.
- Change the corresponding prerequisite in `README.md`.
- Change the corresponding prerequisite in `skills/do-work/actions/version.md`.
- Add or update the focused check that proves the launcher's refusal message quotes the new minimum.
- Search the relevant core installer/updater surfaces for stale 1.26.1 restatements. Do not alter the optional board tool's Go requirement.

The REQ names the four production/doc literals above. A directly relevant existing regression file may also be modified if needed to pin the refusal message. Anything else is a scope expansion and must be reported before editing.

Never touch any path under `do-work/` in the builder worktree. The only main-tree path you may write is the absolute hand-back path above. Do not modify `VERSION`, `skills/do-work/VERSION`, or either changelog; release bookkeeping belongs to the integrator.

## Required rules and context

Before editing, read these files completely from the main checkout:

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/coding-guardrails.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/communication-style.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-shell-commands.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-shell-commands.md`
- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md`

Read the original REQ-407 archive record as context, but treat the confirmed answer in this brief as authoritative. Run focused Go and shell/contract checks relevant to the changed literals. Commit the complete implementation on your branch; do not bump suite versions or edit changelogs.

## Hand-back format

Write the absolute hand-back file before returning. Include:

- branch name and final commit hash;
- P-A-U evidence (brief plan, applied scope, UNIFY checks and files inspected);
- exact file manifest with `(modified)`, `(new)`, or `(deleted)`;
- tests run with direct results;
- stale-restatement search result;
- integration seams, if any (exact line and destination; do not edit shared files);
- `## Decisions` and `## Discovered Tasks` sections when applicable, or state that there were none.

Return only a one-line status after the file is written.
