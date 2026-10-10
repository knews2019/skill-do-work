# Hand-back: REQ-692 (validate-feedback --capture [--run] and a core route)

- Branch: `worktree-agent-REQ-692-validate-feedback-capture-chain`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-692-validate-feedback-capture-chain`
- Base commit (before the first edit): `bd56c4b0`
- Commits:
  - `920879ee` [REQ-692] add the --capture [--run] chain to validate-feedback (validate-feedback.md, capture.md, toolbox help.md)
  - `a596aea7` [REQ-692] route validate-feedback from core do-work to the toolbox action (core SKILL.md)
- Release: not done (integrator owns requirement 12). No CHANGELOG, VERSION or mirror touched.

## File Manifest

- `skills/do-work-toolbox/actions/validate-feedback.md` (modified): `--capture`, `--run` and file-path input in `## Input`; new `### Step 2.5: Wrong-Repo Check (--capture only)`; new Steps 6 (Discuss questions), 7 (capture), 8 (verify, combined report, optional run) after Step 5; one sentence after the Output Format fence; `--capture` exception on the six read-only lines (3, 5, Rules, rationalization row, Red Flags, checklist). Steps 1 to 5 and the Output Format fence are byte-identical (probe hashes pass). File is 196 lines (under the ~230 limit).
- `skills/do-work/SKILL.md` (modified): one routing row at `:35`, directly above the `verify` row: `` | `validate-feedback`, `triage feedback` | `../do-work-toolbox/actions/validate-feedback.md` | ``. Argument hint untouched.
- `skills/do-work/actions/capture.md` (modified): one clause in the `:119` sentence: "or when the user ran `validate-feedback --capture` and it hands capture its accepted findings with their provenance blocks (the same explicit promotion)".
- `skills/do-work-toolbox/actions/help.md` (modified): `:8` description text only, now "Validate external feedback; --capture [--run] captures the accepts". Command column unchanged.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read general.md, coding-guardrails.md (size only, rules known from the brief), shared-principles.md, communication-style.md, clear-questions.md, prime-action-files.md, prime-releases.md, lessons-releases.md and the four lesson families. Checked the restated mechanisms before writing them (restated-mechanism-unchecked): capture.md:119 wording and Step 3 (clarification), capture-reference.md:42 destination 3, the `*Source:*` line template (capture-reference.md:105), verify-requests.md Step 5 gap classes and Step 6/7 boundary, work.md:103 `UR-NNN` targeting, note.md:64 rule ("a note is not a task" forbids a note starting capture, not being called), toolbox `crew-members/clear-questions.md` exists. Approach: keep Steps 1-5 and the fence untouched; put flags and file-path rule in Input; insert Step 2.5 and Steps 6-8 as gated sections; one sentence after the fence; per-line `--capture` clauses; one row and two one-line edits elsewhere.
- [x] **[APPLY]:** Edits exactly as planned in the four files of the write set. Two commits: action chain (+capture clause, help line) and the core route.
- [x] **[UNIFY]:** `git diff bd56c4b0 --stat`:
  ```
   skills/do-work-toolbox/actions/help.md             |  2 +-
   .../do-work-toolbox/actions/validate-feedback.md   | 54 +++++++++++++++++++---
   skills/do-work/SKILL.md                            |  1 +
   skills/do-work/actions/capture.md                  |  2 +-
   4 files changed, 51 insertions(+), 8 deletions(-)
  ```
  Checks: RED probe exit 1 (0 s); invariants probe exit 0 (4 s); GREEN probe exit 0 (4 s, then 5 s after the last wording edit; it includes `_dev/tests/shipped-package-reference-contract.sh` and `_dev/tests/action-shell-blocks.sh`); `git diff --check` exit 0. Also read `_dev/tests/staged-skills-contract.sh` core-routing guards (not run, per brief): the row does not contain `` `./actions/validate-feedback.md` ``, so "core must not route sibling-owned action" still holds, and no line spells a retired trigger. Files checked: all four, each re-read after the edit; no debug artifacts.

## Proof Record

### RED (base `bd56c4b0`, exit 1)
```
validate-feedback does not mention --capture
validate-feedback does not mention --run
validate-feedback lacks the wrong-repo question text (requirement 4)
validate-feedback lacks the Discuss option 'Accept: capture as a REQ' (requirement 5)
validate-feedback lacks the Discuss option 'Park:' (requirement 5)
validate-feedback lacks the Discuss option 'Drop: no work' (requirement 5)
validate-feedback lacks 'do-work verify-requests UR-NNN' (requirements 6, 9, 10)
validate-feedback lacks 'do-work run UR-NNN' (requirements 6, 9, 10)
validate-feedback.md:3: read-only statement without the --capture exception
validate-feedback.md:5: read-only statement without the --capture exception
validate-feedback.md:122: read-only statement without the --capture exception
validate-feedback.md:136: read-only statement without the --capture exception
validate-feedback.md:144: read-only statement without the --capture exception
validate-feedback.md:154: read-only statement without the --capture exception
capture.md does not name validate-feedback --capture (requirement 7, D-03)
toolbox help.md validate-feedback line does not mention --capture (D-04)
core SKILL.md needs exactly one validate-feedback route row, found 0
```

### Invariants (base, exit 0)
```
REQ-692 file checks passed (invariants)
REQ-692 probe passed (invariants)
```

### Behaviour walk-through (by reading; no live run, per brief)

Line numbers are in the edited files at `a596aea7`; VF = `skills/do-work-toolbox/actions/validate-feedback.md`.

| Case | Expected GREEN behaviour (REQ) | Text that produces it |
|---|---|---|
| (a) `--capture --run`, 2 Accept, 1 Discuss answered Accept, 1 Push back, 1 duplicate of a queued REQ | Exactly one question (the Discuss item) with Accept/Park/Drop; one UR, 3 REQs each naming finding id and source; duplicate folds under `## Folded Requests`; verify runs with no second prompt; one combined report; clean verify ends in `do-work run UR-NNN` | VF:32 (flag); VF:99-107 (one question per Discuss item only, options, recommended first); VF:111-113 (accepted set, one payload, finding id + source label, one capture run, fold-first); capture.md:119 (hand-off admitted, so prose-only accepts become REQs via capture-reference.md:42); VF:117 (verify through Step 6, no Step 7, combined report); VF:121 (clean → run via work.md); VF:122 (gap → stop + both commands) |
| (b) `--run` alone | One usage line, stop, no triage | VF:33 |
| (c) a file path | File read; its path is every finding's source | VF:34 (read, source = path as given, outside repo allowed); VF:113 ("Finding 3, source: [path]" label feeds each REQ's source line) |
| (d) `--capture`, cited paths mostly absent | Wrong-repo question before any other question or write; recommended Stop | VF:57-64 (runs after parse, before Step 3; exact question text; Stop recommended; chat fallback) |
| (e) core input "validate-feedback: <a finding>" | Routes to the toolbox action, not capture | `skills/do-work/SKILL.md:35` (first match wins; above `verify` and the capture fallback) |
| (f) same paste, no flag | No question, no write, same report | VF:36 (read-only without the flag); Steps 1-5 and the fence byte-identical (probe hashes); every new step is gated "(--capture only)" (VF:57, 99, 109, 115) |

### GREEN (worktree, exit 0)
```
REQ-692 file checks passed (green)
REQ-692 probe passed (green)
```

## Decisions

<!-- Continues from D-06. -->

- **D-07** (DECIDE & STATE) Wrong-repo check counts **distinct** cited paths, resolved inside the working tree (relative to the repo root). Reasoning: one file cited ten times should not outweigh nine other files. Value: the threshold reflects how many places the review points at. Risk: low; prose, easy to change.
- **D-08** (DECIDE & STATE) The wrong-repo question gets two options, "Stop (recommended)" and "Continue capture here", each with value and risk, plus a chat fallback. Reasoning: clear-questions.md principles 3 and 5 require concrete priced options. Value: answerable in one read. Risk: none.
- **D-09** (DECIDE & STATE) With `--capture`, the finding blocks, Summary and Suggested reply print at the start of Step 6 (before any Discuss question), without the typed handoff block; the combined report at Step 8 replaces that block. Reasoning: the user needs the evidence before answering a Discuss question; requirement 8 and D-02 keep the fence untouched. Value: the answer is informed. Risk: low; the report prints in two parts.
- **D-10** (DECIDE & STATE) Without `--run`, the combined report ends with the next command (`do-work run UR-NNN`, or the two gap commands when verify listed a gap). Reasoning: the replaced handoff block printed `do-work run`; dropping it would lose the next step. Value: user knows what to type. Risk: none; it is a printed line, not an action.
- **D-11** (DECIDE & STATE) The `--run`-alone usage line is `Usage: do-work-toolbox validate-feedback --capture [--run] <findings or file path>`. It names the toolbox command so it never spells a retired core trigger. Risk: none.
- **D-12** (DECIDE & STATE) The capture.md clause says the hand-off is "the same explicit promotion" and does not cite the toolbox file by path. Reasoning: the probe needs only the literal `validate-feedback --capture`; a path citation adds a depth check for no reader benefit. Risk: none.

## Discovered Tasks

- impact-negligible: `skills/do-work-toolbox/actions/journey-qa.md:111` says its capture line "is a suggestion the user runs, as `actions/validate-feedback.md` ends its triage"; still true for the no-flag mode, but it now describes only one of two modes. → report only
- impact-negligible: `skills/do-work-toolbox/actions/maintainability-audit.md` and its reference describe the loop as "validate-feedback → capture handoff → `do-work run`"; the new one-step `--capture --run` path could shorten that loop footer. → report only
- impact-negligible: core `skills/do-work/actions/help.md` does not mention that core now forwards `validate-feedback`; the toolbox help line does. → report only

## Lessons Read

- `_dev/primes/lessons-releases.md` (whole file).
- `_dev/primes/lessons-action-files.md` by family: alternate-writer-contract-drift (13 bullets; applied by sweeping every `validate-feedback` mention under `skills/` for restated read-only rules: none outside the action), restated-mechanism-unchecked (1; applied by reading capture, capture-reference, verify-requests, work and note before restating them), example-path-read-as-citation (1; example paths in the new prose are placeholders, none starts with a package directory name), cross-action-exception-closure (1; the capture.md clause carries the exception through to capture-reference.md destination 3).

## Anti-bloat Check

`git diff bd56c4b0 --stat` is above (four files, +51/-8). Added beyond what the REQ named:
- Step 2.5 placement between Step 2 and Step 3 (D-01, expected).
- A "Stop (recommended)" / "Continue capture here" option pair on the wrong-repo question, with value and risk lines (D-08): required by clear-questions.md for any interactive question.
- The next-command line at the end of the no-`--run` combined report (D-10): keeps the information the replaced handoff block carried.
Nothing else: no new fields, statuses, scripts, tests or reference files.

## Proposed CHANGELOG Entry

**Validate-Feedback Can Capture, Verify and Run the Accepted Findings**

Triage and capture used to be two to five hand-typed prompts per review, and a review written for another repo could be triaged against the wrong code. Now one opt-in flag carries the accepted findings into the queue and, if asked, into a run.

- `do-work-toolbox validate-feedback --capture [--run]` (or the phrases "then capture the accepted ones" / "capture and run") asks one Accept / Park / Drop question per Discuss item, captures the accepted findings as one UR with one REQ each (finding id and source on every REQ), runs verify-requests on it, and prints one combined report. With `--run` and a clean verify it continues into `do-work run UR-NNN`; with a gap it stops and prints the commands.
- A wrong-repo check stops the chain before any write when most cited paths are not in this repo.
- The action accepts a file path as input and records it as each finding's source.
- Core `do-work` now routes `validate-feedback` and `triage feedback` to the toolbox action instead of capturing the text as a request.
- Without the flag the action stays read-only and its report is unchanged.

## Proposed Lesson Bullet

Satellite: `_dev/primes/lessons-action-files.md`

- [family: alternate-writer-contract-drift] [REQ-692: an action's read-only rule was restated on six lines (description, intro, Rules, a rationalization row, Red Flags, checklist), not the four the REQ named; adding an opt-in write mode needs a per-line sweep of every restatement, and a per-line probe ("each match must name the flag on the same line") is what caught the two extra](../../do-work/archive/UR-153/REQ-692-validate-feedback-capture-run.md#lessons-learned)

## Integration Seams and Wall Times

- `skills/do-work/SKILL.md`: only line 35 added (above `verify`). REQ-690 (above clarify) and REQ-691 (above capture, plus argument hint line 4) touch other lines; the merge should be clean.
- `skills/do-work/actions/capture.md`: only line 119 changed. REQ-688 edits Step 5 (around 226-250).
- `skills/do-work-toolbox/actions/help.md`: only line 8 changed. ai-report REQs work near line 14.
- Probe wall times: RED 0 s, invariants 4 s, GREEN 4 s and 5 s. No budget failure.
