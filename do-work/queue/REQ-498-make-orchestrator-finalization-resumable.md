---
id: REQ-498
title: 'Make orchestrator finalization resumable'
status: pending
created_at: 2026-09-02T13:07:19Z
user_request: UR-096
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Make Orchestrator Finalization Resumable

## What
Replace the crash-prone archive/release/commit tail with one CLI-owned, Git-private journaled finalization flow, and recover safe unfinished tails before selecting another REQ.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
An interruption after archive removes the working claim that current recovery scans, while leaving shared checkpoint or release state dirty. Every later claim then fails on the same shared target and the orchestrator cannot make progress automatically.

## Detailed Requirements
- Add strict typed `finalize --manifest` and `recover-finalization --discover` CLI commands.
- Persist an exact Git-private journal before lifecycle mutation and advance it through prepared, lifecycle, release, primary-commit, metadata-commit, verification, and cleanup phases.
- Reuse canonical lifecycle, release, protected-inventory, and Git mutation authorities.
- Preserve sufficient exact pre/post evidence for idempotent recovery without duplicating archive moves, calibration rows, release entries, version bumps, or commits.
- Support serial provenance from the primary commit and worktree provenance from a supplied merge hash.
- Resume journals before ordinary working-REQ recovery, selection, or claim, then continue the same run automatically when shared state is safe.
- Discover legacy unjournaled tails only when project and lifecycle ownership are unambiguous; shared metadata requires semantic REQ evidence and never generic latest-owner association.
- Preserve unrelated unstaged changes; block on ambiguous shared state, foreign staged entries, dirty checkpoint state, or protected paths.
- Return typed phase, resume/discovery, commits, blockers, reason codes, and exact verification commands.
- Update work and commit action contracts to delegate finalization and startup recovery to the CLI.
- Keep existing lifecycle and release commands backward compatible and retain the single-releaser model.

## Constraints
- Journals are local Git-private state.
- Never guess shared metadata ownership or commit secret-classified content.
- Exact-path commits must not absorb unrelated staged or unstaged work.
- The current session should stop after capturing this intent and committing one safe, coherent implementation slice.

## Red-Green Proof
**RED prompt/case:** Interrupt a successful REQ after canonical archive/checkpoint mutation but before its primary or metadata commit, then invoke `do-work run` again.
**Why RED now:** Recovery scans only `do-work/working/`; the archived REQ is invisible and its dirty shared checkpoint makes every next claim refuse.
**GREEN when:** The next run resumes or safely associates the unfinished finalization exactly once, leaves no duplicate lifecycle/release effects, records provenance, and proceeds to the next selectable REQ without manual `do-work commit` intervention.
**Validation:** User confirmed through the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-096/input.md` for the complete verbatim plan and stopping instruction.

## Open Questions
None.

---
*Source: implement and capture the resumable orchestrator finalization plan*
