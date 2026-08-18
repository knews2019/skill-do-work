# REQ-261 builder brief

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-261-decide-whether-queue-kanban-grows-a-date-only-mode` — your working directory.
**Branch:** `worktree-agent-REQ-261-decide-whether-queue-kanban-grows-a-date-only-mode` (already checked out there).
**Hand-back (write when done):** `/home/user/skill-do-work/do-work/runs/work-2026-08-18-230100/REQ-261-handback.md` — absolute main-tree path, the ONE main-tree path you may write. Never stage or commit it.

## The REQ (verbatim)

```markdown
---
id: REQ-261
title: Delete the date-only tripwire and keep the rule
status: claimed
claimed_at: 2026-08-18T22:59:48Z
route: A
created_at: 2026-08-18T19:30:47Z
status_changed_at: 2026-08-18T21:01:24Z
user_request: UR-055
addendum_to: REQ-253
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work-reference.md
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-18T22:59:48Z
  basis:
    - Route A
    - 1-file write set
    - 2 acceptance criteria
    - full-suite verification
---

# Delete the Date-Only Tripwire and Keep the Rule

## What

The Timestamp rule's date-only paragraph ends with "revisit if a second consumer appears". Remove that clause. Keep everything else in the sentence — the shell one-liners, and the reason there is no tool subcommand (adding one would widen the skill's single compiled-dependency exception for something the POSIX floor already covers).

The clause is the only part that does not survive its own argument: it keys on how many consumers exist, and consumer count does not bear on whether a shell one-liner suffices. Leaving it invites a re-litigation the surrounding sentence already settles — a list where a condition belongs (CLAUDE.md → State conditions, not lists).

## Requirements

- The "revisit if a second consumer appears" clause is gone; the rest of the date-only paragraph is unchanged in meaning.
- No date-only subcommand is added to the board tool.
- The paragraph still reads as one coherent sentence after the removal — check the ui-review consumer clause REQ-253 added still sits naturally beside it.
- `bash _dev/tests/maintainer-verify.sh` exits 0 (the Timestamp-rule citation contract counts 54 instant / 17 date-only sites today; a prose-only removal must not move them).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Discovered by REQ-253's builder ([low]; the tripwire sentence was left verbatim as a deliberate maintainer call). Note the board's pinned write-surface count is unaffected either way — `now`-style output is read-only.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-253: the date-only paragraph's "revisit if a second consumer appears" condition is now true. Should I process this as a new task? → Yes, and the answer is decided: delete the tripwire clause and keep the rule
  Recommended: Yes, add to queue (will flip to 'pending') — and the builder decides between the subcommand and a re-stated threshold.
  Also: No, discard it — two consumers on the shell one-liner is still fine.

**Answered [2026-08-18]:** User approved via `do-work clarify` **and settled the underlying question**, after asking where the sentence came from. Provenance established during clarify: the clause arrived with the repository import (recorded author "Claude", root commit `8d5c2ab`) from a pass that restructured the Timestamp rule — it is builder prose, not a maintainer decision. The user's reasoning, which is now the REQ's requirement: the rule itself is sound (`date -u +%F` works on the POSIX floor, and the board tool is the skill's only sanctioned compiled dependency, so putting a date behind a Go binary would widen that exception for nothing), but the tripwire keys on **consumer count**, which has no bearing on that argument — a shell one-liner is no worse at two callers than at one. Delete the clause; keep the rule. Do not add a date-only subcommand.
```
## Never touch

- Anything under `do-work/` — the queue, `working/`, `CHECKPOINT.md`, `runs/` — is the orchestrator's. Your worktree carries a stale committed snapshot of it; treat it as absent. The single exception is your own hand-back file at the absolute path above, which you write but never stage or commit.
- `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md` — serial-only files owned by the integrator. Bumping any of them races every sibling builder. Skip the "Before Every Commit" ritual in `CLAUDE.md` entirely; it belongs to the integrating commit, not to yours.
- Any file outside your REQ's `## Scope` "Files I will touch" list. Discovering mid-build that you need one is a **stop-and-report to the orchestrator, never a silent write** — put it in the hand-back and say why.
- The main tree at `/home/user/skill-do-work` — you are not in it and must not write to it.

## Environment

- Go 1.26.1 at `/usr/local/go`, ShellCheck 0.11.0, `just` 1.21.0, Node 22 are all installed and are the versions the gate pins.
- **Never run a bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — it drops an 11 MB gitignored binary into the source tree. Build to a scratch directory with `-o`.
- The canonical gate is `bash _dev/tests/maintainer-verify.sh`, run from your worktree root. **Exit code zero is the only proof it passed** — never pipe it through `tail` or `head`, because the pipeline's exit status hides the failure. Redirect to a file and echo `$?` if you want to read the output.
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp anything. Never carry a timestamp forward and never compute one.

## How to work

Follow the do-work crew rules that always load during implementation — read them from your worktree:

- `skills/do-work/crew-members/general.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/communication-style.md`
- plus `skills/do-work/crew-members/testing.md` if this REQ has `tdd: true`, and `skills/do-work/crew-members/maintenance.md` if it has `maintenance: true`.

Read every path in the REQ's `prime_files` before touching code, including its `## Lessons` links.

**P-A-U is mandatory and it is your own note-taking** — the REQ file itself lives in the main tree and is not yours to edit, so write the [PLAN] / [APPLY] / [UNIFY] evidence into your hand-back instead. [PLAN] before any code. [APPLY] stays inside the declared scope. [UNIFY] runs `git diff --stat` against your branch point, runs the native linters, verifies no debug artifacts, and lists each file you checked.

**The failure mode that beat nine of twelve REQs last session:** a mechanism that looks like it closes a class and closes only the instance. The two REQs that beat it grepped or fuzzed the **primitive** before declaring the class closed, and both found the real hole where no instance list pointed. Assume your first fix has that shape. An instance list is a sample, not the class.

**Commit in small, individually-green increments on your branch.** Builders were repeatedly killed mid-run by server-side errors last session and nothing was lost, because each increment was already committed. Do not wait until the end to commit.

## Hand-back format

Write `REQ-NNN-handback.md` at the absolute path given above, with these sections:

- `**Branch:**` the operative name, and the commit it was cut from
- `**Commits (oldest first):**` short hash + subject for each
- `## What I built` — factual, what you actually built, not what the REQ asked for
- `## File manifest` — every file created/modified/deleted with `(new)` / `(modified)` / `(deleted)` and a phrase on what changed in it
- `## P-A-U evidence` — [PLAN] / [APPLY] / [UNIFY], with what you actually ran
- `## Testing evidence` — real observed output, never a prototype or a paraphrase. If the REQ is `tdd: true`, the RED must be a real run against pre-change code, quoted. State the gate's exit code as a number you observed.
- `## Decisions (D-XX)` — each marked DECIDE & STATE or ESCALATE; ESCALATE entries carry `Value:` and `Risk:` lines
- `## Integration seams` — exact lines and where they go, for any shared file you must not edit yourself. `None.` if there are none.
- `## Discovered Tasks` — out-of-scope finds, not fixed inline
- `## Pushback` — where the REQ was wrong, or where you disagree. `None.` if you have none.

Your report is a claim, not evidence: the orchestrator judges from git state. Say plainly what you did not do.
