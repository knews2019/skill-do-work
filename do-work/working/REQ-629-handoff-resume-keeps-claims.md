---
id: REQ-629
title: '[impact-rule-change] Resuming a handoff must not reset its own claimed REQs'
status: claimed
created_at: 2026-10-02T19:49:03Z
user_request: UR-133
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: ["REQ-628"]
batch: ur-132-follow-ups
claimed_at: 2026-10-02T20:53:37Z
---
# Resuming a Handoff Must Not Reset Its Own Claimed REQs

## What
Make the handoff (`phandoff`, `actions/restart-with-parallel-handoff.md`) and the crash-recovery takeover agree on what "resume" means, so a session that resumes a handoff continues the handed-off claimed REQs from their written sections instead of resetting them.

## Why
On 2026-10-02 the UR-132 run was handed off with REQ-626 claimed, merged at 505ec74a and gate-green, and REQ-627 claimed with Triage, estimate, Exploration and Scope written. The handoff prompt said "run canonical recover, then continue each claimed REQ". After a machine restart, `recover` reported `RECOVERY-TAKEOVER-AVAILABLE` with the argv `recover --take-over REQ-NNN` ("working claim requires explicit authority"). Running it, as the finding suggested, returned both REQs to `do-work/queue/` as `pending` and stripped every orchestrator-generated section, including route, write_set, Scope, Implementation Summary, Qualification, Testing and Decisions of an already-merged REQ. The run recovered only by re-claiming and restoring the sections byte-for-byte from the handoff commit 1794c268 (decisions REQ-626 D-07 and REQ-627 D-05 in `do-work/archive/UR-132/`). The takeover is documented as a "canonical reset" (`actions/work-reference.md` → Crash Recovery), so the command did what it says; the handoff and the finding's suggested next step led the resuming session into it.

## Detailed Requirements
- Establish, by a real check, whether a resuming session on the same checkout can continue a handed-off claim with per-request `advance REQ-NNN --request-path <working path>` without any takeover, and what `recover`'s plain result then does on the next queue-mode `advance`.
- Fix the side that is wrong (see Open Questions) so that following the handoff's paste block plus the commands' own suggested next steps never resets a claim the handoff described as in flight.
- The handoff paste block (or the handoff action's rules for writing it) states what to run for each claimed REQ and that `recover --take-over` resets a claim, when that remains true after the fix.
- A REQ that is merged (`commit:` or `integration_at` present) must not lose its evidence sections through the documented resume path.

## Constraints
State is binding; prose is advisory (the handoff action's own rule): prefer a fix that the commands enforce over one that only adds prose, when the command change is small. No change to crash recovery for a genuinely crashed claim from another checkout. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-628 (board guide update).

## Builder Guidance
Medium certainty on the fix direction; high certainty on the failure. Reproduce it on a scratch repository first: claim a REQ, write sections, write a handoff, then follow the paste block in a fresh process.

## Open Questions
- [ ] Which side changes: the handoff instructions, or the takeover command?
  Recommended: Both, minimally. The handoff action tells the resuming session to continue its own claims with per-request advance and never to run takeover on them, and `recover`'s finding for a same-writer claim stops offering `--take-over` as the next step (or says plainly that it resets the claim). Why: the trap was the suggested argv; a prose-only fix leaves the trap in the command output.
  Also: Prose only in the handoff action (smallest change, trap stays in the command output). Or a new resume mode in `recover` that keeps a same-writer claim and its sections (largest change; reconsider only if per-request advance cannot continue a claim without a takeover).

## Red-Green Proof
**RED prompt/case:** In a scratch repo, claim a REQ, append Triage and Scope, write and commit a handoff, then in a new session follow the paste block and the `recover` finding's next argv: the REQ ends up in `do-work/queue/` as `pending` with Triage and Scope gone.
**Why RED now:** The handoff says "continue each claimed REQ", and the only next step `recover` offers for that claim is `--take-over`, which is a canonical reset.
**GREEN when:** Following the paste block and the commands' suggested next steps leaves the REQ claimed in `do-work/working/` with its sections intact, and the classifier names its next phase.
**Validation:** Inferred during capture

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Full Context
See `do-work/user-requests/UR-133/input.md` for complete verbatim input.

*Source: capture both as REQs (the handoff-versus-takeover mismatch reported at the end of the UR-132 run)*
