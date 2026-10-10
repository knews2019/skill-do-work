---
id: REQ-666
title: 'Coordinator writes its own handoff at high context, the handoff resumes with --coordinate, and teardown stops idle agents'
status: cancelled
created_at: 2026-10-09T21:16:07Z
user_request: UR-146
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
depends_on: [REQ-662]
related: [REQ-662, REQ-663, REQ-664, REQ-665, REQ-667]
batch: run-coordinate-mode
completed_at: 2026-10-10T12:59:16Z
---
# Coordinator Writes Its Own Handoff and the Handoff Resumes Coordinated
## What
In a `--coordinate` run, the coordinator writes the handoff itself when the harness reports context usage above a threshold, or at the user's `phandoff`. A handoff written from a coordinated run resumes with `do-work run --coordinate --fan-out N`. Teardown also stops idle background agents and lists merged worktrees for removal.
## Why
Report item A5 (UR-146 input, "What happened"): on 2026-10-02 a user asked "give me the phandoff, because the context is 673/1000k which is not ideal. The phandoff should contain instructions on how to run tasks in background agents/workflows", and typed the same request in the development clone 3 minutes later. One repo's `do-work/RESTART-PROMPT.md` was rewritten in about 30 commits since 2026-09-09, each re-stating the coordinator directive by hand. Mid-run on 2026-10-07 the user asked "why keep the idle builders?".
## Verified Facts (checked at 0.305.87)
- `skills/do-work/actions/restart-with-parallel-handoff.md:11`: use when "Context is running out". The trigger is the user, through the `phandoff` alias at `crew-members/communication-style.md:90`.
- `restart-with-parallel-handoff.md:63`: the build resume command is `do-work run --fan-out N`, so a handoff from a coordinated run resumes uncoordinated.
- `actions/work.md:470`: teardown is the checkpoint, then cleanup, after the last integrator returns. Nothing stops idle agents or a stall check.
## Detailed Requirements
1. `restart-with-parallel-handoff.md:63` gains: "a run started with `--coordinate` resumes with `do-work run --coordinate --fan-out N`".
2. Under `--coordinate`, the coordinator follows `restart-with-parallel-handoff.md` on its own when the harness reports context usage above a threshold, as well as at `phandoff`.
3. Teardown stops idle background agents and lists merged worktrees for removal.
4. Document the teardown additions where `actions/fan-out-reference.md` and `actions/work.md:470` describe the coordinator's checkpoint and cleanup.
5. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Choosing the context threshold per harness is out of scope: the action names the rule, the harness supplies the reading (report, Out of scope).
- Teardown lists merged worktrees; it does not remove them. The existing REMOVABLE rule (`restart-with-parallel-handoff.md:55`) and the cleanup action keep removal.
- No new REQ status and no new frontmatter field.
## Assumptions (recorded at capture, no questions asked)
- A harness that reports no context reading never triggers the automatic handoff; `phandoff` still works.
- "Stops idle background agents" means agents with no assigned REQ left in the run. An agent still building or integrating is never stopped by teardown.
- Deleting the stall check at teardown belongs to REQ-664 (the stall check); this REQ only adds the idle-agent and merged-worktree steps beside it.
- The automatic handoff writes and commits the handoff only; it does not end the session by itself.
## Dependencies
Depends on REQ-662 (the `--coordinate` mode the resume command names).
## Builder Guidance
High certainty: three small prose changes in named files. Latitude: where the threshold rule sits and its wording.
## Red-Green Proof
**RED prompt/case:** Start a run with `do-work run --coordinate --fan-out 2`, then type `phandoff`.
**Why RED now:** the written `RESTART-PROMPT.md` resume command is `do-work run --fan-out 2` (`restart-with-parallel-handoff.md:63`), so the next session runs uncoordinated.
**GREEN when:** the handoff written from the coordinated run carries `--coordinate` in its resume command, and teardown output names idle agents stopped and merged worktrees listed for removal.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes the handoff action's resume-command contract read by the next session.
## Full Context
See `do-work/user-requests/UR-146/input.md` for complete verbatim input (sections Request A5, What happened, Where the behaviour lives today, Proposed direction A5, Acceptance check, Out of scope). No queued candidate shares this root cause (the queue held REQ-654 to REQ-661, the ai-report and do-work-cli batches, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md`, Request item A5: "Teardown and handoff: the coordinator writes the handoff on its own when context gets high, and a handoff written from a coordinated run resumes with `--coordinate`."*

## Cancelled

- **When:** 2026-10-10T12:59:16Z
- **Why:** folded into the simplified recapture of 2026-10-10 (maintainer chose the aggressive simplification)
- **Decided by:** user, via `do-work abandon`
