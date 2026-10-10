---
id: REQ-692
title: 'validate-feedback --capture [--run]: file-path input, wrong-repo check, Discuss questions, one capture, verify, optional run, and a core do-work route'
status: claimed
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-690, REQ-691]
batch: validate-feedback-capture
claimed_at: 2026-10-10T13:14:56Z
---
# validate-feedback --capture [--run]: One Opt-In Chain From Triage to an Optional Run, Plus a Core Route
## What
Give `do-work-toolbox validate-feedback` an opt-in `--capture [--run]` flag (also recognised from the phrases "then capture the accepted ones" and "capture and run" in the same invocation) and a file-path input. With `--capture` the action, after its per-item verdicts: stops on a wrong-repo check before any write, asks one question per Discuss item (Accept / Park / Drop), captures the accepted set as one UR with one REQ per finding and its provenance, runs `verify-requests` on that UR, prints one combined report, and with `--run` and a clean verify continues into `do-work run UR-NNN`. Add a core `do-work` routing row so "do-work validate-feedback: …" reaches the toolbox action instead of the capture fallback. Without the flag or phrase the action stays read-only and its output is byte-for-byte unchanged. This REQ replaces the cancelled REQ-670 to REQ-674 (UR-149), folded by the maintainer on 2026-10-10; the installer retire line from REQ-674 is dropped as operator work.
## Why
21 `validate-feedback` prompts in about 18 sessions over 30 days. About 10 sessions ran validate → capture (→ verify / run) by hand, 2 to 5 prompts each; the user retyped accepted finding ids ("capture F1, F2, F6, F9 and C1", 2026-09-17), `do-work verify-requests` was typed right after a capture 12 to 15 times, and "capture and run" was typed in three sessions on 2026-10-08. On 2026-09-28 a review written for a different repo was pasted in and triage ran against code the findings did not describe. On 2026-10-06 the user asked for exactly the question step: "check the following request, talk with me, use the ask tool and `capture-requests` after we got a common understanding". 21 prompts invoked `do-work validate-feedback`, which core `do-work` has no route for.
## Verified Facts (checked at 0.305.101)
- `skills/do-work-toolbox/actions/validate-feedback.md`: line 5 "**Read-only** — this action does NOT modify any files and does NOT create REQs … (Capture ≠ Execute)"; the Input section defines only "the pasted feedback"; **Discuss** "has merit but the right path isn't clear-cut … Frame the trade-off" is written into the report and never asked; the handoff block (around `:113-117`) already says to keep the accepted finding's verbatim claim, original severity/source, Evidence and Surface-cost together in the capture payload, then prints `do-work capture-request:` and `do-work run` as lines for the user to type; the Rules ("Read-only. Modify no files. Create no REQs."), a rationalization row ("I'll capture the accepts to save the user a step") and the checklist ("No files were modified and no REQs were created") all need an "unless `--capture` was given" clause.
- `skills/do-work/SKILL.md`: "Stop after capture unless the same user invocation explicitly requested execution too"; "Load `crew-members/clear-questions.md` before asking an interactive question"; the routing table has no `validate-feedback` row and its last row routes unmatched descriptive input to capture, first match wins. No core row forwards into a sibling package today.
- `skills/do-work-toolbox/SKILL.md`: `| validate-feedback, triage feedback, review feedback | ./actions/validate-feedback.md |`, the only forward today.
- `skills/do-work/actions/capture.md`: a finding reaches capture only when the user invokes `do-work capture` and quotes the finding line as the source; the fold-first scan still runs. `actions/capture-reference.md` → Fold-First Rule.
- `skills/do-work/actions/verify-requests.md` accepts one `UR-NNN`; its Step 7 (Offer Fixes) asks the user after the report. `actions/work.md`: a `UR-NNN` token scopes the run and keeps `depends_on` gating.
- `skills/do-work-toolbox/actions/note.md` exists and is the Park destination.
## Detailed Requirements
A. Input and gating:
1. Add `--capture` and `--run` to the Input section. `--run` alone is rejected with a one-line usage text and no triage. "then capture the accepted ones" counts as `--capture`, "capture and run" as `--capture --run`, when they appear in the same invocation; a close paraphrase with the same explicit intent counts, a phrase that only suggests capture ("these look good") does not.
2. Accept a file path as `$ARGUMENTS`: a single token (after flags) naming an existing regular file is read, and its path is recorded as the source of every finding; anything else is pasted text and recorded as "pasted text". A path outside the repo is allowed and recorded as given.
3. Edit the four read-only statements (line 5, Rules, the rationalization row, the checklist) to read "read-only unless `--capture` was given" or equivalent. Without the flag or phrase, Steps 1 to 5 and the Output Format stay byte-for-byte as they are.
4. Wrong-repo check: with `--capture`, when more than half of the cited `file:line` paths do not resolve to a file in the working tree (line numbers not checked), stop before any question or write and ask "These findings cite paths that are not in this repo. Continue capture here, or stop?", recommended option stop, after loading `crew-members/clear-questions.md`. No cited paths means the check does not fire.

B. Discuss questions (new step, `--capture` only, after the verdicts and the wrong-repo check):
5. One question per Discuss item with the interactive question tool, options recommended first: "Accept: capture as a REQ with remedy <one line>", "Park: `do-work-toolbox note` it for later", "Drop: no work", each with one line of value and one line of risk. Already done, Push back and Accept items are never asked about. No question tool: list the items with the same options and wait in chat. No Discuss items: skip with no output line. Park runs the note action in the same invocation (`--capture` already authorizes writes); if note's contract forbids chaining, print the note command and record that in Decisions.

C. Capture (new step, `--capture` only):
6. Build one payload from the Accept items plus the Discuss items answered Accept, each as the existing provenance block (verbatim claim, original severity/source, Evidence, Surface-cost), and run `skills/do-work/actions/capture.md` on it once: one UR, one REQ per finding. Each REQ's source line names the finding id (for example "Finding 3") and the input kind (pasted text, or the file path). A finding that duplicates a queued REQ folds through the existing fold-first scan. An empty set says so and stops; with `--run` it also starts no run.
7. `--capture` counts as the user invocation `capture.md` requires; if that line's wording would refuse a capture started from validate-feedback, add one clause naming `validate-feedback --capture`, otherwise leave `capture.md` unchanged. Capture's own clarification step still applies, but the Discuss answers count as resolved.
8. Rewrite the handoff block so it describes what happens with `--capture` and keeps the typed lines for the no-flag case.

D. Verify and optional run (new step, `--capture` only):
9. Run `verify-requests` on the new `UR-NNN` with no second prompt, through its report and not into its Offer Fixes step. Print one combined report: the triage summary table, the UR and its REQs with titles, and the verify verdict. Without `--run`: stop.
10. With `--run` and no listed gap (Important, Minor or Ambiguous; the gap list decides, not the score): continue to `do-work run UR-NNN`. With a gap: stop before the run and print the gaps with the exact commands to run after fixing them (`do-work verify-requests UR-NNN` then `do-work run UR-NNN`, real id filled in). A UR that owns no REQs (everything folded) skips the run and says so.

E. Core route:
11. Add a row to `skills/do-work/SKILL.md`'s routing table above the capture fallback: `| validate-feedback, triage feedback | ../do-work-toolbox/actions/validate-feedback.md |`, so "do-work validate-feedback: …" routes to the toolbox action and not to capture. "review feedback" stays out of the core row because the core `review` row matches first.
12. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Capture ≠ Execute stays as it is today: without the flag or an explicit phrase, no question, no capture, no run, for any verdict including impact-critical ones.
- No new REQ fields, no new statuses. No change to verdict rules, to capture's templates, numbering, reservations or transaction, to verify-requests, or to how `do-work run` selects work.
- Out of scope: a `--review N` front end; a new "what earned this" lens; the installer check for the old standalone `do-validate-feedback` skill (the maintainer ruled it operator work on 2026-10-10: the user removes those folders with `git rm -r` once).
- REQ-690 (do-work status) and REQ-691 (do-work trace) each add a row to the same routing table; whichever lands later rebases its row.
## Assumptions (recorded at capture)
- One REQ: the five cancelled REQs edited one 154-line action file in series and each was its own release; the maintainer folded them on 2026-10-10. The builder may ship the core route as a separate commit inside the REQ.
- The question tool's per-call batch limit decides how many Discuss questions share one call, not this REQ.
- The UR's verbatim input is the triage input (the paste, or the file's bytes), with the accepted subset and the per-finding verdicts in the UR Summary; capture's raw-containment rule decides the exact bytes.
## Dependencies
None.
## Builder Guidance
Certainty is high on behaviour; the source report fixes the flags, phrases, edits, options, stop conditions and the row text. Latitude: wording of the usage line and the "read-only unless" clauses, the combined report's layout, where the fallback-to-chat sentence sits.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox validate-feedback --capture --run` on a paste with 2 Accept, 1 Discuss and 1 Push back items, Discuss answered Accept; the same with `--run` alone; the same with a file path; a paste whose cited paths are mostly absent, with `--capture`; and the input "do-work validate-feedback: <a sample finding>".
**Why RED now:** the action defines only a paste as input and no flags, forbids every write path, ends at the typed handoff lines, and core `do-work` routes the phrase to capture.
**GREEN when:** `--run` alone prints one usage line and stops; a file path is read and its path appears as the source; the mostly-absent paste stops at the wrong-repo question before any other question or write; exactly one question is asked, about the Discuss item, with Accept / Park / Drop; one new UR with 3 REQs exists, each naming its finding id and source, and a fourth duplicate finding folds under `## Folded Requests`; `verify-requests` runs on that UR without a second prompt and one combined report prints; with a clean verify the chain ends in `do-work run UR-NNN`, with a gap it stops and prints the gaps and the exact commands; "do-work validate-feedback: …" routes to the toolbox action; the same paste with no flag produces the same report as today and writes no files.
**Validation:** Inferred during capture (from the source report's Acceptance checks and the maintainer's 2026-10-10 decision).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: changes an action's input contract and read-only rules, chains three actions, and adds a core routing row.
## Full Context
See `do-work/user-requests/UR-153/input.md` for the decision record. The cancelled originals with their full bodies are under `do-work/archive/UR-149/` (REQ-670 to REQ-674); their source report is `do-work/inbox/2026-10-09_do-work-upstream-suggestion-validate-feedback-capture.md`.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: maintainer decision of 2026-10-10 (UR-153) folding REQ-670 to REQ-674: "fold into one REQ", "keep only the routing row", installer retire line dropped as operator work.*
