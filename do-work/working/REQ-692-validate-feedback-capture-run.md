---
id: REQ-692
title: 'validate-feedback --capture [--run]: file-path input, wrong-repo check, Discuss questions, one capture, verify, optional run, and a core do-work route'
status: claimed
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-690, REQ-691]
write_set: ["skills/do-work-toolbox/actions/validate-feedback.md", "skills/do-work/SKILL.md", "skills/do-work/actions/capture.md", "skills/do-work-toolbox/actions/help.md"]
batch: validate-feedback-capture
claimed_at: 2026-10-10T13:14:56Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-10-10T13:22:16Z
  basis:
    - Route B
    - 4-file write set
    - 2 subsystems involved
    - 9 acceptance criteria
builder_handback_at: 2026-10-10T13:29:20Z
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

---

## Triage

**Route: B** - Medium

**Reasoning:** The behaviour is fully specified (flags, phrases, the four read-only edits, the question options, stop conditions and the routing row text), but it chains three existing actions (capture, verify-requests, work) and each seam needed reading: capture's finding-entry rule, the prose-backlog destination, verify-requests Step 7, note's contract, and the staged-skills contract that guards core routing of sibling commands. Prose only, two packages (core and toolbox), no Go.

**Planning:** Not required

## Decisions

Step 3.5: the REQ has no `## Open Questions` section, so nothing was deferred. The pre-dispatch agent made these implementation choices the REQ left open; each is DECIDE & STATE (reversible, low reach), so none mints a follow-up. The builder continues numbering at D-07.

- **D-01** (DECIDE & STATE) Where does the wrong-repo check run? Chose: a new `--capture`-only step between Step 2 (parse) and Step 3, so a wrong-repo paste stops before any verification work, question or write. Reasoning: parsing is the first point where the cited paths are known, and the 2026-09-28 incident was triage run against the wrong code; inserting a gated step keeps the text of Steps 1 to 5 byte-identical. Value: no wasted triage on a foreign review. Risk: low; it is prose placement and easy to move.
- **D-02** (DECIDE & STATE) Requirement 3 (Output Format byte-for-byte without the flag) and requirement 8 (rewrite the handoff block) overlap. Chose: keep the fenced Output Format template, including its three typed handoff lines, byte-identical, and describe the `--capture` replacement (the combined report) in prose outside the fence. Reasoning: this satisfies both requirements literally; the no-flag report is unchanged. Value: the no-flag acceptance check is a byte diff of one fence. Risk: low.
- **D-03** (DECIDE & STATE) Does capture's finding-entry rule refuse a capture started from validate-feedback? Chose: yes, so add one clause. `skills/do-work/actions/capture.md:119` admits a finding only "when the user invokes `do-work capture` and quotes the complete report-only finding line". The clause names `validate-feedback --capture` as the same explicit promotion, which also makes the Fold-First Rule's destination-3 queue-work exception apply (`skills/do-work/actions/capture-reference.md:42`), so a prose-only accepted finding becomes a REQ as requirement 6 asks, not a prose-backlog line. `capture-reference.md` itself is not edited. Value: one REQ per accepted finding, as the GREEN case counts. Risk: low; one sentence, reversible.
- **D-04** (DECIDE & STATE) Toolbox help line and argument hints. Chose: change only the description text of the `validate-feedback` line in `skills/do-work-toolbox/actions/help.md:8` to mention `--capture [--run]`, keeping the column alignment; leave both `argument-hint` lines and core `help.md` unchanged. Reasoning: `_dev/primes/prime-action-files.md` asks help to follow a public surface change; the argument hints list command names, not flags. Value: the flag is discoverable from the menu. Risk: negligible.
- **D-05** (DECIDE & STATE) Where does the core routing row sit? Chose: directly above the `verify` row (`skills/do-work/SKILL.md:35`), not directly above the capture fallback. Reasoning: still above the fallback as requirement 11 asks; first-match-wins then also resolves "validate-feedback: … check …" before the `check` trigger; and REQ-691 inserts its row directly above the capture row and REQ-690 above the clarify row, so a third insertion at a separate spot keeps the serial merges conflict-free. Value: clean merges. Risk: none to behaviour.
- **D-06** (DECIDE & STATE) Park chaining. Chose: Park runs `skills/do-work-toolbox/actions/note.md` in the same invocation. Reasoning: note's contract forbids a note starting capture, work or a commit, not being called by another action; `--capture` already authorizes writes. Value: one invocation. Risk: low.

<!-- D-XX counter: last used D-06. Next decision: D-07. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Required Lessons Consult (claim time)

Index consulted (`do-work/lessons-index.md`). Selected: `_dev/primes/lessons-releases.md` (666 tokens; the REQ releases per `_dev/primes/prime-releases.md`, and the new core row adds a cross-package citation that `_dev/tests/shipped-package-reference-contract.sh` checks). Dropped for budget: `_dev/primes/lessons-action-files.md` (7756 tokens, `slugged: partial`, so no targeted form), listed in `## Required Lessons — Dropped for Budget` above, which stays current. The brief still points the builder at four of its bullets by family as optional reading (alternate-writer-contract-drift, restated-mechanism-unchecked, example-path-read-as-citation, cross-action-exception-closure). Not matched: shell, board, CLI and updater satellites (no shell, Go or updater change).

## Exploration

Read directly by the pre-dispatch agent at main `bd56c4b0`.

**Files and exact anchors**
- `skills/do-work-toolbox/actions/validate-feedback.md` (154 lines). Read-only statements: `:3` (description blockquote, "Read-only; offers a capture handoff"), `:5`, `:122` (Rules), `:136` (rationalization row), `:144` (Red Flags, "it must be read-only"), `:154` (checklist). The REQ names four; `:3` and `:144` restate the same rule and need the same clause (alternate-writer-contract-drift). Input section `:28-30`. Steps 1 to 5 `:34-82`. Output Format fence `:88-118`, with the handoff block `:113-117` inside the fence. Rules `:120-127`, rationalizations `:129-136`, Red Flags `:138-144`, checklist `:146-154`.
- `skills/do-work/SKILL.md` routing table `:28-46`, first match wins. `verify` row `:35` (its triggers include `check`), `review` row `:36`, capture fallback `:46`. No core row forwards to a sibling today.
- `skills/do-work/actions/capture.md:119`: "A review, build, triage, or consumer-report finding reaches capture only when the user invokes `do-work capture` and quotes the complete report-only finding line as the source". This wording would refuse a hand-off from validate-feedback, so requirement 7's clause is needed (D-03). The same promotion gates `skills/do-work/actions/capture-reference.md:42` (Fold-First destination 3: a prose-only explicit capture goes to the prose backlog "except when the user is promoting a complete quoted report-only finding line"); the capture.md clause covers it without editing capture-reference.md.
- `skills/do-work/actions/verify-requests.md`: Step 5 gap classes `:105-114` (Important, Minor, Nit, Ambiguous); Step 6 report `:116-152`; Step 7 Offer Fixes `:154`. The chain stops after Step 6. Requirement 10's gap list is Important, Minor or Ambiguous; Nit does not stop the run.
- `skills/do-work/actions/work.md:103`: a `UR-NNN` token scopes the run and keeps `depends_on` gating.
- `skills/do-work-toolbox/actions/note.md`: appends one line to `do-work/notes.md`; its rule `:64` only forbids a note starting capture, work or a commit, so validate-feedback may call it (D-06).
- `skills/do-work-toolbox/crew-members/clear-questions.md` and `prompt-injection.md` exist in the toolbox package, so the action's relative `crew-members/...` loads resolve.
- `skills/do-work-toolbox/actions/help.md:8`: `validate-feedback [findings]   Validate external feedback before accepting it` (D-04).

**Guards the change must pass**
- `_dev/tests/staged-skills-contract.sh` (heavy tier, about 55 s, gate only) scans every live file under `skills/` for retired core triggers from `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`. Eight of them belong to validate-feedback: `validate-feedback`, `validate feedback`, `triage findings`, `triage feedback`, `feedback review`, `review feedback`, `assess feedback`, `should we push back`. Any shipped text that spells "do-work" followed by a space and one of those fails the gate. So the action, capture.md and SKILL.md must never write that form (the REQ's own example input does). Write `do-work-toolbox validate-feedback`, or the bare `validate-feedback --capture`. The routing row itself passes: its cell is `` `validate-feedback` ``, and the path `../do-work-toolbox/...` is not "do-work" plus a space. The same script also fails if the core routing section contains `` `./actions/validate-feedback.md` ``, so the row must use the `../do-work-toolbox/actions/validate-feedback.md` form exactly.
- `_dev/tests/shipped-package-reference-contract.sh` (about 2 s) checks citation depth: `../do-work-toolbox/...` from core `SKILL.md`, `../../do-work/actions/...` from a toolbox action, `../../do-work-toolbox/actions/...` from a core action.
- `_dev/tests/action-shell-blocks.sh` (about 5 s) lints any shell fence the builder adds.

**Seams with sibling REQs in this run**
- `skills/do-work/SKILL.md`: REQ-690 (do-work status) inserts a row above the clarify row `:38`. REQ-691 (do-work trace) inserts above the capture row `:46` and edits the argument hint `:4`. This REQ inserts above the verify row `:35` and leaves `:4` alone, so the three insertions do not touch each other (D-05).
- `skills/do-work/actions/capture.md`: REQ-688 (capture-files --example) edits Step 5 (around `:226-250`). This REQ edits only the sentence at `:119`.
- `skills/do-work-toolbox/actions/help.md`: REQ-655, REQ-657 and REQ-687 (ai-report) may edit the ai-report line `:14` or add lines near it. This REQ edits only `:8`.
- `skills/do-work-toolbox/actions/validate-feedback.md`: no other member touches it.

**Concerns**
- The core row reverses the moved-command migration's "core routes no sibling action" shape (commit `0b9bcde`). The maintainer ruled it on 2026-10-10 ("keep only the routing row"), and the contract's mechanical checks allow the `../do-work-toolbox/` form. So this is recorded, not reopened.
- Steps 1 to 5 must stay byte-identical (requirement 3). The file-path input (requirement 2) therefore goes in the Input section, which defines what "the pasted feedback" means for a file. Step 1's prompt-injection load already covers a file's bytes as third-party content.

*Generated by the pre-dispatch agent (direct reads; no Explore subagent)*

## Scope

**Files I will touch:**
- skills/do-work-toolbox/actions/validate-feedback.md (modify): flags and file-path input in Input; the --capture-only wrong-repo step after Step 2; the Discuss-question, capture, and verify-and-run steps after Step 5; read-only clauses at lines 3, 5, 122, 136, 144 and 154; prose after the Output Format fence describing the --capture report
- skills/do-work/SKILL.md (modify): one routing row above the verify row
- skills/do-work/actions/capture.md (modify): one clause in the line 119 sentence naming validate-feedback --capture as an explicit promotion
- skills/do-work-toolbox/actions/help.md (modify): the description text of the validate-feedback line

**Files I will NOT touch:** `skills/do-work/actions/capture-reference.md`, `skills/do-work/actions/verify-requests.md`, `skills/do-work/actions/work.md`, `skills/do-work-toolbox/actions/note.md`, `skills/do-work-toolbox/SKILL.md`, both `argument-hint` lines, `skills/do-work/actions/help.md`, `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv`, any Go file, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, version files.

**Acceptance criteria (restated from REQ):**
- [ ] `--run` alone prints one usage line and stops, with no triage.
- [ ] A file path as `$ARGUMENTS` is read and its path is recorded as every finding's source; other input is recorded as "pasted text".
- [ ] With `--capture`, a paste whose cited paths are mostly absent stops at the wrong-repo question before any other question or write.
- [ ] Exactly one question per Discuss item, with Accept / Park / Drop options, recommended first, each with value and risk; no question for Already done, Push back or Accept items; no Discuss items means no output line.
- [ ] One new UR with one REQ per accepted finding, each naming its finding id and source; a duplicate of a queued REQ folds through the fold-first scan; an empty set says so and stops.
- [ ] `verify-requests` runs on the new UR without a second prompt, stops before Offer Fixes, and one combined report prints.
- [ ] With `--run` and no Important, Minor or Ambiguous gap the chain continues to `do-work run UR-NNN`. With a gap it stops and prints the gaps and the exact commands. A UR that owns no REQs skips the run and says so.
- [ ] "validate-feedback: …" sent to core do-work routes to the toolbox action, not capture; "review feedback" stays out of the core row.
- [ ] Without the flag or phrase, the report is the same as today and no files are written (Steps 1 to 5 and the Output Format fence byte-identical).
