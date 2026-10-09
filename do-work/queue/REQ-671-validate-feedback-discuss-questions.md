---
id: REQ-671
title: 'validate-feedback --capture asks one question per Discuss item: accept, park or drop'
status: pending
created_at: 2026-10-09T21:26:55Z
user_request: UR-149
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
depends_on: [REQ-670]
related: [REQ-670, REQ-672, REQ-673, REQ-674]
batch: validate-feedback-capture
---
# validate-feedback --capture Asks One Question per Discuss Item
## What
With `--capture`, after the per-item verdicts, `validate-feedback` asks the user about each Discuss item with the interactive question tool: Accept (capture as a REQ with a one-line remedy), Park (`do-work-toolbox note` it), or Drop. Recommended option first, one line of value and one line of risk on each option. Without `--capture`, nothing is asked.
## Why
Report item C2 (UR-149 input): the Discuss trade-off is written into the report but the user is never asked to decide it (`validate-feedback.md:78`), so users answered Discuss items in chat by hand. On 2026-10-06 the user asked for exactly this: "check the following request, talk with me, use the ask tool and `capture-requests` after we got a common understanding".
## Verified Facts (checked at 0.305.87)
- `skills/do-work-toolbox/actions/validate-feedback.md:78`: **Discuss** "has merit but the right path isn't clear-cut … Frame the trade-off."
- `validate-feedback.md:117`: the handoff today prints `do-work-toolbox note "[a discuss item]"` as a line for the user to type.
- `skills/do-work/SKILL.md:61`: "Load `crew-members/clear-questions.md` before asking an interactive question." `crew-members/clear-questions.md` § 3 (say the consequence of each option) and § 5 (concrete options, never open-ended) govern the wording.
- `skills/do-work-toolbox/actions/note.md` exists and is the Park destination.
## Detailed Requirements
1. Add a new Step 6, Discuss questions, that runs only with `--capture` (REQ-670), after the per-item verdicts and after REQ-670's wrong-repo check passes.
2. One question per Discuss item, using the interactive question tool after loading `crew-members/clear-questions.md`.
3. Options, recommended first: "Accept: capture as a REQ with remedy <one line>", "Park: `do-work-toolbox note` it for later", "Drop: no work". Each option carries one line of value and one line of risk.
4. Already done and Push back items are never asked about. Accept items are not asked about either; they go to capture as they are.
5. If no question tool is available, list the Discuss items with the same options and wait for the answer in chat.
6. The answers feed REQ-672's capture step: Discuss items answered Accept join the Accept items.
7. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Capture ≠ Execute stays as it is today: without `--capture` or an explicit phrase, no question is asked and the output is unchanged.
- No new REQ fields, no new statuses. No change to verdict rules.
- Out of scope: a `--review N` front end; any new "what earned this" lens; capturing without the flag or phrase.
## Assumptions (recorded at capture, no questions asked)
- The recommended option per question is the agent's judgment from the trade-off it framed; the report fixes the option order shown, but the option marked recommended is placed first whichever it is.
- Park runs the `do-work-toolbox note` action in the same invocation, because `--capture` already authorizes writes. If the builder finds note's own contract forbids being chained, print the note command instead and record that in Decisions.
- When there are no Discuss items, Step 6 is skipped with no output line.
- Many Discuss items are asked one question each; the question tool's per-call batch limit decides how they are grouped, not this REQ.
## Dependencies
Depends on REQ-670 (the `--capture` flag this step is gated on). REQ-672 (capture once) depends on this REQ.
## Builder Guidance
Certainty is high; the report fixes the three options and the gating. Latitude: the step's exact wording and where the fallback-to-chat sentence sits.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox validate-feedback --capture` on a paste with 2 Accept, 1 Discuss and 1 Push back items.
**Why RED now:** `validate-feedback.md:78` writes the Discuss trade-off into the report and the action ends at the handoff; no step asks the user.
**GREEN when:** exactly one question is asked, about the Discuss item only, with Accept / Park / Drop options, recommended first, each with a value line and a risk line. The same paste with no flag asks nothing.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing action contracts; this REQ adds a step to an action.
## Full Context
See `do-work/user-requests/UR-149/input.md` for complete verbatim input (sections Request C2, Where the behaviour lives today, Proposed direction C2, Acceptance check). No queued candidate shares this root cause.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-validate-feedback-capture.md`, Request item C2: "After the per-item verdicts, ask the user about each Discuss item with the interactive question tool, recommended option first, value and risk on each option."*
