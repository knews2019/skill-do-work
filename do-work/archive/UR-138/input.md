---
id: UR-138
title: 'Three coordination lessons for do-work fan-out runs: delegated integration, mid-run messages, and a wave-end consistency check'
created_at: 2026-10-07T19:12:33Z
requests: [REQ-639, REQ-640, REQ-641]
word_count: 1703
---
# Three Coordination Lessons for do-work Fan-Out Runs

## Summary
The maintainer's session asked to capture three coordination lessons as three REQs chained with depends_on A -> B -> C. The lessons come from Matt Maher's video "The New Way to Work With AI" (2026-10-07) and from the UR-137 run on 2026-10-07, where the main session coordinated, builders ran in parallel worktrees, and one background integrator per REQ ran in series. Every requirement, verified fact, red-green proof and "not in scope" decision in the source is carried into the REQ bodies.

Not in scope, decided by the source (no REQ created):
- L4 (coordinator suggests extra work from the objective) is rejected because it contradicts capture's "represent, don't expand".
- L5 (N builders on one REQ, pick a winner) waits for an observed need.
- L6 (host resistance to sub-agents): no edit; work.md line 31 and SKILL.md line 52 already instruct sub-agent use.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-639 | [impact-rule-change] Delegated integration: in fan-out mode the main session may hand each REQ's integration to one agent at a time and stay a coordinator |
| REQ-640 | [impact-rule-change] Mid-run messages: route a user message that arrives during a run to the right REQ in the user's own words without stopping the run |
| REQ-641 | [impact-rule-change] Wave-end consistency check: the last review in a fan-out wave also sweeps for elements earlier members redefined |

## Batch Constraints
- Serial chain: REQ-640 depends on REQ-639, and REQ-641 depends on REQ-640, because all three touch `skills/do-work/actions/work.md`, `skills/do-work/actions/review-work.md` and the Fan-Out Dispatch section of `skills/do-work/actions/work-reference.md`.
- Prose-only changes to shipped action files plus one user-guide paragraph. No Go changes, no new flags, no new action files, no SKILL.md routing rows.
- Read `_dev/primes/prime-action-files.md` before editing. Sweep every restatement of a changed rule across `work.md`, `work-reference.md`, `docs/work-guide.md` and `crew-members/background-agents.md` (family alternate-writer-contract-drift).
- Keep every existing heading in `work-reference.md` unchanged; `_dev/tests/shipped-package-reference-contract.sh` pins citation strings that name them.

## Full Verbatim Input
> ```
> Three coordination lessons for do-work, taken from Matt Maher's video "The New Way to Work With AI" (2026-10-07) and from the UR-137 run on 2026-10-07 where the main session coordinated, builders ran in parallel worktrees, and one background integrator per REQ ran in series. Capture as three REQs chained with depends_on A -> B -> C, because all three touch skills/do-work/actions/work.md, skills/do-work/actions/review-work.md and the Fan-Out Dispatch section of skills/do-work/actions/work-reference.md. All three are prose-only changes to shipped action files plus one user-guide paragraph. No Go changes, no new flags, no new action files, no SKILL.md routing rows. Read _dev/primes/prime-action-files.md before editing and sweep every restatement of a changed rule across work.md, work-reference.md, docs/work-guide.md and crew-members/background-agents.md (family alternate-writer-contract-drift). Keep every existing heading in work-reference.md unchanged; _dev/tests/shipped-package-reference-contract.sh pins citation strings that name them.
> 
> REQ A. Delegated integration: the coordinator shape.
> 
> Goal: in fan-out mode the main session may hand each REQ's integration to one agent at a time, in series, and stay a coordinator that only writes pre-dispatch judgment, routes messages, and checkpoints. No flag: who plays the orchestrator role for a span is a placement, like Step 7's review agent ("Spawn an agent with actions/review-work.md ... Or read it and follow it in the current session"), and the dispatch mechanism stays unspecified. Decision: the coordinator writes the pre-dispatch sections (Triage, Open Questions, Plan, lessons consult, Scope, pre-flight) before each dispatch, as the fan-out contract already requires ("briefs written before any spawn"); writing Scope afterwards from the hand-back makes scope-drift detection vacuous.
> 
> Verified facts: plain `recover` preserves a live claim and reports RECOVERY-TAKEOVER-AVAILABLE with next_argv `advance REQ-NNN`, which continues the claim (work-reference.md Crash Recovery; recovery_commands.go:88-104). `recover --assume-sole-authority` and `--take-over` reset every sibling claim. Targeted `advance` on a working REQ classifies by section presence and every section step in work.md is idempotent, so an agent entering late continues at the first missing section. `advance --checkpoint` preserves only foreign or unlabelled records, so an integrator running it could drop same-writer sibling claims. work.md:284 forbids passing `dispatch_at` read back from the file to record-timing-event.
> 
> Changes:
> 1. skills/do-work/actions/work.md Step 6, new condition before "Spawn a general-purpose agent": if this REQ's row in the run manifest records a landed hand-back and the hand-back file exists, consume it: do not dispatch, take the dispatch instant from that manifest row for record-timing-event, stamp builder_handback_at only if absent, and go straight to the hand-back merge, which already proves the branch state. Never re-dispatch a builder for a landed hand-back. State it as a condition so it also serves a fresh session after a crash (background-agents.md: agents whose findings file exists are done).
> 2. work.md Step 6 dispatch paragraph: the manifest row for a REQ carries the held dispatch instant, so an integrator can record the builder-work timing event without reading dispatch_at back from the file.
> 3. work.md Step 10, one sentence: under delegated integration the integrator never runs `advance --checkpoint` or the loop, because the checkpoint preserves only foreign records and would drop same-writer sibling claims; the coordinator runs it, then cleanup, after the last integrator returns.
> 4. skills/do-work/actions/work-reference.md, Worktree Dispatch Mode, Fan-Out Dispatch: new paragraph "Delegated integration — the coordinator shape" stating in this order: (a) after the wave's builders are dispatched, the orchestrator may hand each REQ's integration (hand-back merge through Step 9) to one agent at a time, in series, in the same checkout; release and changelog stay serial-only (cite the existing Serial-only paragraph). (b) The integrator enters through plain `do-work run REQ-NNN`: recover reports the live claim and `advance REQ-NNN` continues it; the existing section steps skip what is already written; Step 6's landed-hand-back condition takes it to the merge. Never run-with-recovery, --assume-sole-authority or --take-over for an integrator: those reset every sibling claim and requeue the whole wave. (c) One writer under the project root at a time: while an integrator runs, the coordinator writes nothing under the project root (no REQ sections, no manifest edits, no captures); otherwise this is the "two sessions in one working tree" case the Execution Model leaves unspecified. The coordinator's turn to write is the gap between integrators; its live notes go to its own session scratch until then. (d) Why: the integration span is long and context-heavy, and running it in the main session blocks the conversation the user steers from. (e) The integrator brief, REQ-NNN-integrate.md in the run directory, written by the coordinator in the gap before that integrator starts: REQ id, run directory, hand-back path, operative name, wave membership, any mid-run addendum text, and the instruction not to run Step 10. Add one row to the guardrail-slot table: per-integrator input | REQ-NNN-integrate.md — written by the coordinator between integrators, never by a builder. Keep "Dispatch mechanism is deliberately unspecified" as is; the new paragraph must not say subagent versus session.
> 5. skills/do-work/docs/work-guide.md, "Building several REQs at once": two sentences saying the session that runs --fan-out can stay a coordinator by handing each REQ's integration to one agent at a time, and that the saving is still build-phase only.
> 
> Red-Green: RED: today work.md Step 6 unconditionally spawns a builder, so an integrator entering a REQ with a landed hand-back would build it twice. GREEN: Step 6 names the landed-hand-back condition, and work-reference names the integrator placement, the one-writer rule, the never-rwr rule, and the brief. Validation: prose read plus _dev/tests/shipped-package-reference-contract.sh and _dev/tests/contract-regressions.sh.
> 
> REQ B. Mid-run messages: steer a run while it is in progress. depends_on REQ A.
> 
> Goal: one rule for a user message that arrives while a run is in progress, so a change of emphasis reaches the right REQ in the user's own words without stopping the run. Disk is the durable home; forwarding to a running builder is a courtesy. Verified: the takeover reset strips only generated headings (requeststate/state_apply.go:981), the advance classifier ignores unknown sections, and Go reads only the addendum_to frontmatter field, so a `## Addendum (mid-run)` body section on a working REQ is safe and survives.
> 
> Changes:
> 1. work.md: new subsection "### Mid-Run Messages (any step)" directly after Step 3.5, beside its existing "question escalated mid-run and answered" paragraph. Condition-keyed, not a list of message types: a question about progress or state is answered from disk (REQ files, run manifest, hand-backs), never from memory of the conversation, and the run does not stop. A change that affects an in-flight REQ whose builder has not handed back: the durable record is `## Addendum (mid-run)` on the working REQ, the user's words under actions/clarify.md Step 4's Outside-text containment plus one framing line (extends, narrows, or corrects which requirement); in the serial loop the orchestrator writes it now; under delegated integration the coordinator puts the text in that REQ's integrator brief and the integrator writes the section before the merge (one writer under the project root); forward it to the builder where the harness can message a running agent, otherwise say in the progress output that the builder will not see it and review will judge it; the section is user intent, so the takeover reset keeps it. A change that affects an in-flight REQ whose hand-back already landed: capture's existing in-flight path (new REQ with addendum_to), reported as queued for the next loop. A change to a queued REQ, or new work: capture's existing paths, run in the coordinator's writing gap. A change of emphasis for the hand-back only (what to show first): write it to the run manifest, not memory; the Decision Brief reads it. Two rules across all branches: forward the user's words, never a paraphrase, because meaning is lost in restatement; route the message to every REQ it affects, not only the one the user named.
> 2. skills/do-work/actions/review-work.md Step 5 item 1: extract requirements from What/Detailed Requirements, any `## Addendum` section, and the UR.
> 3. skills/do-work/actions/capture.md Step 2 table, do-work/working/ row: keep "NEVER modify" for capture; add one pointer clause: a message that arrives during a live run in this session is routed by actions/work.md -> Mid-Run Messages first.
> 4. work.md Orchestrator Checklist: one line "Mid-run message: route per Mid-Run Messages; never stop the run for it."
> 
> Red-Green: RED: today a mid-run "the builder should surface waiting-on-me items, and put decisions first in the update" has no path; capture turns it into a next-run addendum REQ and the in-flight builder never sees it. GREEN: the subsection names the in-flight path and review reads the addendum. Validation: prose read, citation contract test.
> 
> REQ C. Wave-end consistency check. depends_on REQ B.
> 
> Goal: catch the one fan-out defect no per-REQ review can see: a sibling merged later restates a contract an earlier sibling redefined.
> 
> Changes:
> 1. review-work.md Step 6 Restatement Sweep: record the result in the appended `## Review` block as one line, `**Restatement sweep:** redefined <elements> | nothing redefined`. Add the line to the "Append to REQ File" template and point the existing Verification Checklist item at it. Go checks section presence only, so the extra line is safe. New condition in the sweep's trigger (step 1): when the REQ under review is the wave's last successful integration (wave membership from the run manifest; a set-aside member does not skip the check), the trigger set also includes every element the wave's earlier members recorded as redefined, read from their archived `## Review` lines. One-sentence reason: an earlier member's review ran before the later member merged, so only the last review sees both. Findings route exactly as today (impact token, -> report only unless impact-critical).
> 2. work-reference.md Fan-Out Dispatch bullet "The non-interference proof is the merge, not the pick": one sentence pointing at the wave-end sweep as the semantic half the merge cannot provide.
> 
> Red-Green: RED: in a two-REQ wave where REQ-1 redefines a token and REQ-2 restates the old meaning, REQ-1's sweep runs before REQ-2 merges and REQ-2's sweep only looks at what REQ-2 redefined; the stale restatement ships. GREEN: the last-of-wave review's trigger set includes REQ-1's recorded elements and finds it. Validation: prose read; no new test, since no existing test exercises review prose.
> 
> Not in scope, decided: L4 (coordinator suggests extra work from the objective) is rejected because it contradicts capture's "represent, don't expand". L5 (N builders on one REQ, pick a winner) waits for an observed need. L6 (host resistance to sub-agents): no edit; work.md line 31 and SKILL.md line 52 already instruct sub-agent use.
> ```
