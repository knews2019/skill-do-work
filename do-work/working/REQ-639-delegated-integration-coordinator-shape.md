---
id: REQ-639
title: '[impact-rule-change] Delegated integration: in fan-out mode the main session may hand each REQ''s integration to one agent at a time and stay a coordinator'
status: claimed
created_at: 2026-10-07T19:12:33Z
user_request: UR-138
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-640, REQ-641]
batch: coordination-lessons
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, skills/do-work/docs/work-guide.md]
claimed_at: 2026-10-07T19:14:46Z
---
# Delegated Integration: The Coordinator Shape
## What
In fan-out mode the main session may hand each REQ's integration to one agent at a time, in series, and stay a coordinator that only writes pre-dispatch judgment, routes messages, and checkpoints. There is no flag: who plays the orchestrator role for a span is a placement, like Step 7's review agent ("Spawn an agent with actions/review-work.md ... Or read it and follow it in the current session"), and the dispatch mechanism stays unspecified.
## Why
The integration span is long and context-heavy, and running it in the main session blocks the conversation the user steers from. The shape was observed in the UR-137 run on 2026-10-07 (fixing accepted consumer-review findings), where the main session coordinated, builders ran in parallel worktrees, and one background integrator per REQ ran in series. The lesson also draws on Matt Maher's video "The New Way to Work With AI" (2026-10-07).
## Decision (from the source)
The coordinator writes the pre-dispatch sections (Triage, Open Questions, Plan, lessons consult, Scope, pre-flight) before each dispatch, as the fan-out contract already requires ("briefs written before any spawn"). Writing Scope afterwards from the hand-back makes scope-drift detection vacuous.
## Verified Facts (from the source)
- Plain `recover` preserves a live claim and reports RECOVERY-TAKEOVER-AVAILABLE with next_argv `advance REQ-NNN`, which continues the claim (work-reference.md Crash Recovery; recovery_commands.go:88-104).
- `recover --assume-sole-authority` and `--take-over` reset every sibling claim.
- Targeted `advance` on a working REQ classifies by section presence, and every section step in work.md is idempotent, so an agent entering late continues at the first missing section.
- `advance --checkpoint` preserves only foreign or unlabelled records, so an integrator running it could drop same-writer sibling claims.
- work.md:284 forbids passing `dispatch_at` read back from the file to record-timing-event.
## Detailed Requirements
1. `skills/do-work/actions/work.md` Step 6: add a new condition before "Spawn a general-purpose agent". If this REQ's row in the run manifest records a landed hand-back and the hand-back file exists, consume it: do not dispatch, take the dispatch instant from that manifest row for record-timing-event, stamp `builder_handback_at` only if absent, and go straight to the hand-back merge, which already proves the branch state. Never re-dispatch a builder for a landed hand-back. State it as a condition so it also serves a fresh session after a crash (background-agents.md: agents whose findings file exists are done).
2. `work.md` Step 6 dispatch paragraph: the manifest row for a REQ carries the held dispatch instant, so an integrator can record the builder-work timing event without reading `dispatch_at` back from the file.
3. `work.md` Step 10, one sentence: under delegated integration the integrator never runs `advance --checkpoint` or the loop, because the checkpoint preserves only foreign records and would drop same-writer sibling claims; the coordinator runs it, then cleanup, after the last integrator returns.
4. `skills/do-work/actions/work-reference.md`, Worktree Dispatch Mode, Fan-Out Dispatch: a new paragraph "Delegated integration — the coordinator shape" stating, in this order:
   - (a) After the wave's builders are dispatched, the orchestrator may hand each REQ's integration (hand-back merge through Step 9) to one agent at a time, in series, in the same checkout. Release and changelog stay serial-only (cite the existing Serial-only paragraph).
   - (b) The integrator enters through plain `do-work run REQ-NNN`: recover reports the live claim and `advance REQ-NNN` continues it; the existing section steps skip what is already written; Step 6's landed-hand-back condition takes it to the merge. Never run-with-recovery, `--assume-sole-authority` or `--take-over` for an integrator: those reset every sibling claim and requeue the whole wave.
   - (c) One writer under the project root at a time: while an integrator runs, the coordinator writes nothing under the project root (no REQ sections, no manifest edits, no captures); otherwise this is the "two sessions in one working tree" case the Execution Model leaves unspecified. The coordinator's turn to write is the gap between integrators; its live notes go to its own session scratch until then.
   - (d) Why: the integration span is long and context-heavy, and running it in the main session blocks the conversation the user steers from.
   - (e) The integrator brief, `REQ-NNN-integrate.md` in the run directory, written by the coordinator in the gap before that integrator starts: REQ id, run directory, hand-back path, operative name, wave membership, any mid-run addendum text, and the instruction not to run Step 10.
   - Add one row to the guardrail-slot table: `per-integrator input | REQ-NNN-integrate.md — written by the coordinator between integrators, never by a builder`.
   - Keep "Dispatch mechanism is deliberately unspecified" as is. The new paragraph must not say subagent versus session.
5. `skills/do-work/docs/work-guide.md`, "Building several REQs at once": two sentences saying the session that runs `--fan-out` can stay a coordinator by handing each REQ's integration to one agent at a time, and that the saving is still build-phase only.
## Constraints
- Prose-only changes to shipped action files plus one user-guide paragraph. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- No flag: the integrator is a placement, and the dispatch mechanism stays unspecified.
## Dependencies
None. First in the chain: REQ-640 (mid-run messages) depends on this REQ, and REQ-641 (wave-end consistency check) depends on REQ-640. The serial order was requested because all three touch `work.md`, `review-work.md` and the Fan-Out Dispatch section of `work-reference.md`.
## Builder Guidance
Certainty is high: the source names every file, the exact insertion points, the order of the new paragraph's parts, and the table row text. Latitude is limited to wording. Do not add a flag, a new action file, or a subagent-versus-session statement. REQ-640 later adds a mid-run addendum path that the integrator brief in requirement 4(e) already names ("any mid-run addendum text"); write the brief field now as stated.
## Red-Green Proof
**RED prompt/case:** Today `work.md` Step 6 unconditionally spawns a builder, so an integrator entering a REQ with a landed hand-back would build it twice.
**Why RED now:** No condition in Step 6 checks the run manifest for a landed hand-back before dispatch, and `work-reference.md` names no integrator placement.
**GREEN when:** Step 6 names the landed-hand-back condition, and `work-reference.md` names the integrator placement, the one-writer rule, the never-rwr rule, and the brief.
**Validation:** Inferred during capture from the maintainer-authored source, which states this RED/GREEN pair. Checked by a prose read plus `_dev/tests/shipped-package-reference-contract.sh` and `_dev/tests/contract-regressions.sh`.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the source names family `alternate-writer-contract-drift`, and its owning prime governs the action files this REQ edits.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 18680 tokens, `slugged: partial`; over budget. Matching reason: index row lists family `alternate-writer-contract-drift`. Weak fit: this REQ changes no do-work-cli code.
## Full Context
See `do-work/user-requests/UR-138/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "REQ A. Delegated integration: the coordinator shape." in the maintainer session's capture source, 2026-10-07.*
