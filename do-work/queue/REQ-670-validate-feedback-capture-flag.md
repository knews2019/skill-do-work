---
id: REQ-670
title: 'validate-feedback --capture and --run flags, a file-path input, and a wrong-repo check before any write'
status: pending
created_at: 2026-10-09T21:26:55Z
user_request: UR-149
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-671, REQ-672, REQ-673, REQ-674]
batch: validate-feedback-capture
claimed_at: 2026-10-10T12:52:21Z
status_changed_at: 2026-10-10T12:57:18Z
---
# validate-feedback --capture and --run Flags, a File-Path Input, and a Wrong-Repo Check Before Any Write
## What
Give `do-work-toolbox validate-feedback` an opt-in `--capture [--run]` flag, and accept the phrases "then capture the accepted ones" and "capture and run" in the same invocation as the same request. Accept a file path as the input so a review or audit file keeps its source. Before any write on the `--capture` path, stop when most cited paths are not in this repo. Without the flag or phrase, the action stays read-only and its output is unchanged.
## Why
Report item C1 (UR-149 input, "What happened"): 21 `validate-feedback` prompts in about 18 sessions over 30 days. About 10 sessions ran validate → capture (→ verify / run) by hand, 2 to 5 prompts each. On 2026-09-17 the user retyped "capture F1, F2, F6, F9 and C1" from the report just printed. On 2026-09-28 a review written for a different repo was pasted in and triage ran against code the findings did not describe.
## Verified Facts (checked at 0.305.87)
- `skills/do-work-toolbox/actions/validate-feedback.md:5`: "**Read-only** — this action does NOT modify any files and does NOT create REQs … (Capture ≠ Execute)".
- `validate-feedback.md:30`: input is "the pasted feedback"; a file path is not an input form, so a review or audit file's source is lost.
- `validate-feedback.md:122` (Rules: "Read-only. Modify no files. Create no REQs."), `:136` (rationalization row "I'll capture the accepts to save the user a step"), `:154` (checklist "No files were modified and no REQs were created") all need an "unless `--capture` was given" clause.
- `skills/do-work/SKILL.md:20`: "Stop after capture unless the same user invocation explicitly requested execution too." The flag is that explicit request.
- `skills/do-work/SKILL.md:61` and `crew-members/clear-questions.md` govern the wrong-repo question's wording.
## Detailed Requirements
1. Add `--capture` and `--run` to the Input section of `validate-feedback.md`. `--run` is valid only together with `--capture`; `--run` alone is rejected with a one-line usage text and no triage.
2. Treat "then capture the accepted ones" as `--capture` and "capture and run" as `--capture --run` when they appear in the same invocation.
3. Accept a file path as `$ARGUMENTS`: read the file and record its path as the source of every finding. Pasted text keeps working as today and is recorded as "pasted text".
4. Edit `:5`, `:122`, `:136` and `:154` so each reads "read-only unless `--capture` was given" (or equivalent). Without the flag or phrase, Steps 1 to 5 and the Output Format stay byte-for-byte as they are.
5. Wrong-repo check: when `--capture` is set and more than half of the cited `file:line` paths do not exist in this repo, stop before the Discuss questions (REQ-671) and before any write, and ask: "These findings cite paths that are not in this repo. Continue capture here, or stop?" Recommended option: stop. Load `crew-members/clear-questions.md` before asking.
6. The wrong-repo check runs only when a write is about to happen (`--capture` set). It costs one existence check per cited path.
7. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Capture ≠ Execute stays as it is today: without the flag or an explicit phrase, the action is read-only and its output is unchanged.
- No new REQ fields, no new statuses.
- No change to verdict rules, to capture's REQ templates, or to how `do-work run` selects work.
- Out of scope: a `--review N` front end that reviews the last N commits before triage; any new "what earned this" lens (Step 4 item 5 at `validate-feedback.md:69` and the Surface-cost field already apply it); capturing without the flag or phrase for any verdict, including impact-critical ones.
- Related queued work, no edge: REQ-668 (the `do-work status` action) and REQ-669 (the `do-work trace` action) also edit the core routing table, which REQ-674 touches; this REQ does not.
## Assumptions (recorded at capture, no questions asked)
- The wrong-repo check is part of this REQ because it guards the `--capture` path before REQ-671's questions and REQ-672's capture. The report places it between C1 and C2 without a number.
- The two phrases are the condition "the same invocation explicitly asks to capture the accepted findings (and run them)". The quoted wording is illustrative; a close paraphrase with the same explicit intent counts. A phrase that only suggests capture ("these look good") does not.
- A file path input is recognised when the argument (after flags) is a single token naming an existing regular file. Anything else is pasted text. A path outside the repo is allowed and recorded as given.
- When no `file:line` paths are cited at all, the wrong-repo check has nothing to count and does not fire.
- "Do not exist in this repo" means the path part of the citation does not resolve to a file in the working tree; the line number is not checked.
- If no question tool is available, the wrong-repo question is printed with the same two options and the action waits for the answer in chat.
## Dependencies
None. Root of the chain: REQ-671 (Discuss questions) depends on this REQ.
## Builder Guidance
Certainty is high on behaviour; the report names the flags, phrases, edits and the check. Latitude: wording of the usage line, where in the Input section the flags are documented, and how the "read-only unless" clause is phrased in each of the four places.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox validate-feedback --run <a fixed sample paste>`; and `do-work-toolbox validate-feedback do-work/audits/<an audit file>`.
**Why RED now:** `validate-feedback.md:30` defines only a paste as input and no flags, so `--run` is treated as part of the paste and a file path is triaged as text; the rules at `:5`, `:122` forbid any write path.
**GREEN when:** `--run` without `--capture` prints one usage line and stops. A file path input is read and its path appears as the source in the report. A paste whose cited paths are mostly absent from the repo, given with `--capture`, stops at the wrong-repo question before any other question or write. The same sample paste with no flag produces the same report as 0.305.87 and writes no files.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers changing action routing and status contracts; this REQ changes the input contract and read-only rules of an action.
## Full Context
See `do-work/user-requests/UR-149/input.md` for complete verbatim input (sections Request C1, What happened, Where the behaviour lives today, Proposed direction C1 and Wrong-repo check, Acceptance check, Out of scope). No queued candidate shares this root cause (the queue held REQ-654 to REQ-669 at capture; none touches validate-feedback).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-validate-feedback-capture.md`, Request item C1: "`validate-feedback --capture [--run]`, plus the phrases "then capture the accepted ones" and "capture and run" in the same invocation. `--run` is only valid together with `--capture`. Also accept a file path as the input, so a review or audit file keeps its source."*
