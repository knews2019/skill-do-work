---
id: REQ-664
title: 'Coordinated run arms one stall check that restarts a silent integrator from its landed phase and reports a silent builder'
status: cancelled
created_at: 2026-10-09T21:16:07Z
user_request: UR-146
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-662]
related: [REQ-662, REQ-663, REQ-665, REQ-666, REQ-667]
batch: run-coordinate-mode
completed_at: 2026-10-10T12:59:15Z
---
# Coordinated Run Arms One Stall Check
## What
After the first dispatch of a `--coordinate` run, the coordinator arms one recurring check every 15 to 20 minutes that reads run liveness. An integrator silent for 20 minutes or more is stopped, and a fresh one resumes from its landed phase. A silent builder is only reported. Each tick also re-runs the selector so newly captured REQs join the next wave. Teardown deletes the check.
## Why
Report item A3 (UR-146 input, "What happened"): in the 2026-10-07 to 10-08 run, background integrators stalled for 3 hours and 2 hours. The only liveness signal in the suite (`staleClaimThreshold`) fires after 3 hours, which is the length of the stall observed. Users added their own stall-check cron prompt, which fired 121 times in one session and 25 in another, and asked mid-run "added more REQ's pick them up" and "why keep the idle builders?".
## Verified Facts (checked at 0.305.87)
- `skills/do-work-board/tools/queue-kanban/verify.go:80`: `const staleClaimThreshold = 3 * time.Hour`.
- `skills/do-work/crew-members/background-agents.md:11-14`: the pattern makes failures "survivable and recoverable", and nothing detects a silent agent while the run is live.
- `skills/do-work/actions/fan-out-reference.md:126`: an integrator enters with the read-only `advance REQ-NNN`, which continues the live claim and skips sections already written; never `recover`.
## Detailed Requirements
1. After the first dispatch, arm one recurring check (the harness's scheduler, for example CronCreate or `/loop`) every 15 to 20 minutes.
2. Each tick reads liveness: the run-status subcommand once it exists, otherwise each `REQ-NNN-progress.log` tail (REQ-667, the coordinated-run rules) and each worktree's last commit time.
3. An integrator silent 20 minutes or more is stopped, and a fresh integrator resumes from its landed phase through `advance REQ-NNN`. It never redoes a merge already on main. The user gets one line.
4. A silent builder is reported only, unless the user says restart.
5. Each tick re-runs the selector so newly captured REQs join the next wave.
6. Each coordinator turn ends with one line: done/total, active lanes, ETA.
7. Teardown deletes the check. Exactly one check is armed per run.
8. Document this in `actions/fan-out-reference.md` under Delegated integration and in `docs/work-guide.md`.
9. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Building the run-status subcommand is out of scope; a companion suggestion report asks for it (report, Related suggestions and Out of scope).
- Dispatch-mechanism neutrality (`actions/work.md:37`) stays. Name the scheduler as the harness's, with CronCreate and `/loop` as examples, not as the only routes.
- No new REQ status and no new frontmatter field.
- Never `recover`, `recover --assume-sole-authority` or `recover --take-over` on a restart (`fan-out-reference.md:127`).
## Assumptions (recorded at capture, no questions asked)
- "Silent" means no new progress-log line and no new worktree commit for the threshold. The 20-minute value is approximate, as the maintainer's standing rule says for numeric targets.
- "Landed phase" is the phase `advance REQ-NNN` names for the claim. The restart never re-runs a hand-back merge whose merge commit is already on the integration branch; the existing landed-hand-back rule (`fan-out-reference.md` → Landed hand-back) already governs this.
- On a harness without a scheduler, the coordinator checks liveness at the start of each of its own turns instead and says so in one line.
- Liveness history: REQ-069 and REQ-073 deleted queue-level heartbeat and liveness probing, and `lessons-do-kanban.md:45` keeps liveness out of the board's worktree verify. This check is a coordinator-side judgment during a live run, not queue state and not a board probe. The builder states that in the prose.
- The ETA line may say unknown when the run has no finished integration to measure from.
## Dependencies
Depends on REQ-662 (the `--coordinate` mode). It reads the progress log from REQ-667, which REQ-662 already depends on. When the run-status subcommand from the companion report ships, the tick calls it; no edge to that unfiled request.
## Builder Guidance
Medium certainty: the rules are clear, the harness mechanics vary. Latitude: tick interval inside 15 to 20 minutes, the one-line formats, and the no-scheduler fallback wording.
## Red-Green Proof
**RED prompt/case:** In a coordinated run, kill the integrator process mid-run, after its hand-back merge commit is on main.
**Why RED now:** nothing checks liveness during the run; the stall is noticed only by the user or by the 3-hour board threshold.
**GREEN when:** one stall check is armed after the first dispatch; within one tick after the 20-minute silence the coordinator starts a fresh integrator from the landed phase; main has exactly one merge commit for that REQ; teardown deletes the check.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes the fan-out action contract.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8334 tokens, over budget; `slugged: partial`). Matching reason: its REQ-083 lesson holds the liveness boundary.
## Full Context
See `do-work/user-requests/UR-146/input.md` for complete verbatim input (sections Request A3, What happened, Where the behaviour lives today, Proposed direction A3, Acceptance check, Related suggestions). No queued candidate shares this root cause (the queue held REQ-654 to REQ-661, the ai-report and do-work-cli batches, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-run-coordinate-mode.md`, Request item A3: "One stall check armed for the run (every 15 to 20 minutes) that reads run liveness. A silent integrator is stopped and restarted from its landed phase. A silent builder is only reported."*

## Cancelled

- **When:** 2026-10-10T12:59:15Z
- **Why:** folded into the simplified recapture of 2026-10-10 (maintainer chose the aggressive simplification)
- **Decided by:** user, via `do-work abandon`
