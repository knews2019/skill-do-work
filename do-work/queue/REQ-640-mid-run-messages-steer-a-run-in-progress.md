---
id: REQ-640
title: '[impact-rule-change] Mid-run messages: route a user message that arrives during a run to the right REQ in the user''s own words without stopping the run'
status: pending
created_at: 2026-10-07T19:12:33Z
user_request: UR-138
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-639, REQ-641]
batch: coordination-lessons
depends_on: [REQ-639]
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, skills/do-work/actions/capture.md]
---
# Mid-Run Messages: Steer a Run While It Is in Progress
## What
One rule for a user message that arrives while a run is in progress, so a change of emphasis reaches the right REQ in the user's own words without stopping the run. Disk is the durable home; forwarding to a running builder is a courtesy.
## Why
Today a mid-run steering message has no path: capture turns it into a next-run addendum REQ and the in-flight builder never sees it. The lesson comes from the UR-137 run on 2026-10-07 (fixing accepted consumer-review findings) and Matt Maher's video "The New Way to Work With AI" (2026-10-07).
## Verified Facts (from the source)
- The takeover reset strips only generated headings (requeststate/state_apply.go:981).
- The advance classifier ignores unknown sections.
- Go reads only the `addendum_to` frontmatter field.
- So a `## Addendum (mid-run)` body section on a working REQ is safe and survives.
## Detailed Requirements
1. `skills/do-work/actions/work.md`: new subsection `### Mid-Run Messages (any step)` directly after Step 3.5, beside its existing "question escalated mid-run and answered" paragraph. Condition-keyed, not a list of message types:
   - A question about progress or state is answered from disk (REQ files, run manifest, hand-backs), never from memory of the conversation, and the run does not stop.
   - A change that affects an in-flight REQ whose builder has not handed back: the durable record is `## Addendum (mid-run)` on the working REQ, holding the user's words under `actions/clarify.md` Step 4's Outside-text containment plus one framing line (extends, narrows, or corrects which requirement). In the serial loop the orchestrator writes it now. Under delegated integration (REQ-639, the coordinator shape) the coordinator puts the text in that REQ's integrator brief and the integrator writes the section before the merge (one writer under the project root). Forward it to the builder where the harness can message a running agent; otherwise say in the progress output that the builder will not see it and review will judge it. The section is user intent, so the takeover reset keeps it.
   - A change that affects an in-flight REQ whose hand-back already landed: capture's existing in-flight path (new REQ with `addendum_to`), reported as queued for the next loop.
   - A change to a queued REQ, or new work: capture's existing paths, run in the coordinator's writing gap.
   - A change of emphasis for the hand-back only (what to show first): write it to the run manifest, not memory; the Decision Brief reads it.
   - Two rules across all branches: forward the user's words, never a paraphrase, because meaning is lost in restatement; route the message to every REQ it affects, not only the one the user named.
2. `skills/do-work/actions/review-work.md` Step 5 item 1: extract requirements from What/Detailed Requirements, any `## Addendum` section, and the UR.
3. `skills/do-work/actions/capture.md` Step 2 table, `do-work/working/` row: keep "NEVER modify" for capture; add one pointer clause: a message that arrives during a live run in this session is routed by `actions/work.md` -> Mid-Run Messages first.
4. `work.md` Orchestrator Checklist: one line "Mid-run message: route per Mid-Run Messages; never stop the run for it."
## Constraints
- Prose-only changes to shipped action files. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.
- Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (lesson family alternate-writer-contract-drift).
- Read `_dev/primes/prime-action-files.md` before editing.
- The subsection is keyed on conditions, not a closed list of message types.
## Dependencies
Depends on REQ-639 (delegated integration, the coordinator shape): the delegated-integration branch of the new subsection uses REQ-639's integrator brief and one-writer rule. REQ-641 (wave-end consistency check) depends on this REQ. The serial order was requested because all three touch `work.md`, `review-work.md` and the Fan-Out Dispatch section of `work-reference.md`.
## Builder Guidance
Certainty is high: the source names every insertion point and the branch conditions. Latitude is limited to wording. Do not turn the subsection into a list of message types.
## Red-Green Proof
**RED prompt/case:** During a run, the user says "the builder should surface waiting-on-me items, and put decisions first in the update". Today this has no path: capture turns it into a next-run addendum REQ and the in-flight builder never sees it.
**Why RED now:** `work.md` has no mid-run message rule, and review does not read `## Addendum` sections when extracting requirements.
**GREEN when:** The new subsection names the in-flight path, and review reads the addendum.
**Validation:** Inferred during capture from the maintainer-authored source, which states this RED/GREEN pair. Checked by a prose read plus the citation contract test (`_dev/tests/shipped-package-reference-contract.sh`).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: the source names family `alternate-writer-contract-drift`, and its owning prime governs the action files this REQ edits.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 18680 tokens, `slugged: partial`; over budget. Matching reason: index row lists family `alternate-writer-contract-drift`. Weak fit: this REQ changes no do-work-cli code.
## Full Context
See `do-work/user-requests/UR-138/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "REQ B. Mid-run messages: steer a run while it is in progress." in the maintainer session's capture source, 2026-10-07.*
