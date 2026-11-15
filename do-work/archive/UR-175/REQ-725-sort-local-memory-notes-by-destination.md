---
id: REQ-725
title: 'Sort the local memory notes by destination'
status: completed
created_at: 2026-11-14T22:44:21Z
user_request: UR-175
domain: general
prime_files: []
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-724]
batch: machine-local-rules
route: A
estimate:
  p50_active_minutes: 10
  confidence: medium
  calculated_at: 2026-11-14T23:30:09Z
  basis:
    - Route A
    - 1-file write set
    - 5 acceptance criteria
claimed_at: 2026-11-14T23:28:11Z
dispatch_at: 2026-11-14T23:32:30Z
builder_handback_at: 2026-11-15T02:56:13Z
integration_at: 2026-11-15T02:56:26Z
review_at: 2026-11-15T03:04:22Z
kb_status: pending
commit: 406d32aa2db7ec68e7d25d8dbfe6d53574d16fbd
heavy_verified_at: 2026-11-15T03:12:04Z
heavy_verified_revision: 406d32aa2db7ec68e7d25d8dbfe6d53574d16fbd
completed_at: 2026-11-15T03:12:43Z
---
# Sort the Local Memory Notes by Destination
## What
Review the notes in the maintainer's local memory folder and sort each one into "move to the repository", "move to the skill" or "delete". The result is a sorted list the maintainer can approve. No note is moved or deleted in this REQ.
## Why
The memory folder exists on one machine only, like the global rules file that REQ-724 moves. A session on another machine starts without these notes. The maintainer asked for this review during capture.
## Context
- The folder is `~/.claude/projects/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/memory/`. On 2026-11-15 it held 48 files: `MEMORY.md` (the index) and about 47 notes.
- Most notes are run recipes and past decisions about this repository. Some are standing preferences that fit the skill, for example "one broken pipe" and "warnings over exemptions".
- The folder is outside the repository. A builder on a machine without it cannot do this work.
## Detailed Requirements
- Judge every note. For each one, record its file name, one line on what it says, its destination, and one line of reason.
- Destinations: "move to the repository" (maintainer knowledge about this repository, for example a prime or a lessons file under `_dev/`), "move to the skill" (a rule for every project that uses do-work), "delete" (stale, wrong, already recorded in the repository, or only true for one past run).
- For "move to the repository" and "move to the skill", name the target file. For "already recorded", name where.
- Check a note against the current repository before trusting it. A note that names a file, flag or version may be out of date.
- Deliver the sorted list as one report, through the project's normal report path.
## Constraints
- Read only outside the repository. Never edit or delete a file in the memory folder.
- No shipped file changes in this REQ, so it is not a release. The moves themselves are later requests, after the maintainer approves the list.
- If the memory folder is absent on the build machine, stop this REQ with a finding that says so. Do not guess the notes.
## Dependencies
None. REQ-724 is related but independent.
## Builder Guidance
- Certainty: the maintainer chose "Yes, capture a review request" with the three destinations named. The report format is the builder's choice.
- The maintainer prefers behaviour built into the skill over notes that must be recalled. When a note is a rule for all projects, prefer "move to the skill".
## Red-Green Proof
**RED prompt/case:** List the notes in the memory folder and look for a document in the repository that gives each one a destination.
**Why RED now:** No such document exists. The notes live only on one machine and none has a recorded destination.
**GREEN when:** One report lists every note file in the folder exactly once, each with a destination, a reason and a target or "already recorded at" path. The count of rows equals the count of note files. `git status` shows no change outside the report, and the memory folder is byte-identical before and after.
**Validation:** Inferred during capture
## Full Context
See `do-work/user-requests/UR-175/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** read all 47 notes as data, check each claim with fixed-string searches at `df1bf423`, apply the REQ-724 and REQ-732 seams, write one proposal bundle with one `data-memory-note` row per note. (from the builder hand-back)
- [x] **[APPLY]:** the bundle only; rows sorted by destination, then file name; counts 5 + 13 + 29 = 47 shown at the top; options O1 (whole list), O2 (smallest change: delete only), O3 (do nothing), each priced; capture lines printed only. Integrator correction (D-11): 46 rows, 5 + 13 + 28. (from the builder hand-back)
- [x] **[UNIFY]:** `git diff --stat main...HEAD` lists only `index.html`, `presentation.css`, `presentation.js` in the bundle; checked that the kit files are byte copies, the probe passes, the judge passes, and the memory snapshots match. (from the builder hand-back)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-move-global-rules-into-skill.md*

---

## Triage

**Route: A** - Simple

**Reasoning:** Report only: one new `ai-reports/` bundle, no code, no shipped file, no test file. The output and its location are fixed by the REQ; the per-note checks against the repository are the work itself, so a separate exploration pass and a pre-flight baseline add nothing.

**Planning:** Not required

## Open Questions

- [~] Which route and report form? → **D-01**: Builder chose: Route A, delivered with `ai-report --kind proposal REQ-725` (a decision brief about open work), one bundle under `ai-reports/`. Reasoning: the REQ says "the project's normal report path"; a disposition list for approval is a decision brief, and the completed-work form would present open moves as shipped. Value: the maintainer gets one page to approve, catalogued like every other report. Risk: low; the report is unshipped and a revise replaces it.
- [~] How does a check prove "every note exactly once"? → **D-02**: Builder chose: each note row carries `data-memory-note="<file name>"` and `data-destination="move-to-repository|move-to-skill|delete"`; the GREEN probe counts those attributes against the live folder. `MEMORY.md` is the index, not a note, and gets no row. Reasoning: a file name can also appear in prose and target paths, so a plain text count is ambiguous. Value: the acceptance is checkable by a script. Risk: low; the attributes are invisible in the page.
- [~] The memory folder can gain a note from another session during the run. → **D-03**: Builder chose: the builder records a SHA-256 list of the folder before and after its run in the hand-back; the probe compares rows with the folder as it is when the probe runs. Reasoning: a note added later is an unsorted note, and a failing probe is the right signal. Value: no silent gap. Risk: a perishable probe; the integrator reruns the builder's sort for the new note only.

<!-- D-XX counter: last used D-03. Next decision: D-04. -->

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

Builder decisions (from the builder hand-back):

- **D-04** (DECIDE & STATE): `ask-tool-when-action-says-ask.md` is **delete** ("already recorded at communication-style.md and clear-questions.md (after REQ-724)"), not a move. Its rule is REQ-724's G2, and G2 has no autonomy exception, so the incident it records is covered.
- **D-05** (DECIDE & STATE): `one-orchestrator-per-queue.md` is **move to the skill**, but only its path-limited commit rule (`git commit … -- <paths>`, target `skills/do-work/actions/commit.md`). Its incident log and "never probe for sessions" stance are not moved. Reason: `commit.md:125-127` and `144-146` still show `git add` then a bare `git commit`, which commits a sibling's staged files in a shared checkout.
- **D-06** (DECIDE & STATE): five run-recipe notes (`single-req-run-recipe`, `coordinator-parallel-builders-serial-integrators`, `gate-runs-need-a-quiet-machine`, `shared-tree-dirt-detached-checkout`, plus the live traps only) target one new prime, `_dev/primes/prime-work-runs.md`, instead of `CLAUDE.md`, which every session loads. The report asks this as Q2.
- **D-07** (ESCALATE): `maintainer-verify-slated-for-deletion.md` is marked **delete** on a stale premise (two months later the gate is the standard integrator check). Value: one fewer note that contradicts current practice. Risk: if the maintainer still means to delete the gate, the note's "never re-add a CLAUDE.md pointer" warning is lost. Asked as Q1 in the report.
- **D-08** (DECIDE & STATE): `run-prompts-not-in-root-docs.md` is **delete**, judged against REQ-732 (which allows `docs/run-prompts/` in this run) as well as main, where the note is still true. The reason cell says both.
- **D-09** (DECIDE & STATE): the bundle timestamp is local time (`2026-11-15_0137`), as the other reports in `ai-reports/` use it.
- **D-10** (DECIDE & STATE): the 13 "move to the repository" notes include `maintainability-audit-brief-shape.md`, targeted as a run prompt in `docs/run-prompts/` after REQ-732 (today `_dev/run-prompts/`), because it is a reusable brief, not a rule.

Integrator decisions:

- **D-11** (DECIDE & STATE): the memory folder lost `run-prompts-not-in-root-docs.md` after the build (the REQ-732 integrator deleted it because REQ-732 made it wrong). On the builder branch (`189da332`) the report drops that row, says 46 notes and 28 deletes (5 + 13 + 28), and adds one sentence that the note was deleted on 2026-11-15 when REQ-732 landed. The O3 risk example that named the note now names `fast-gate-caches-since-0-304.md`. The corrections went onto the builder branch before the first merge, so one merge carries them; the result is the same as merge then re-merge, with one fewer merge commit.
- **D-12** (DECIDE & STATE): the four seam rows were re-checked against main `3d0c483c`. `ask-tool-when-action-says-ask.md` and `subagent-model-split.md` stay delete, now citing the lines REQ-724 added (communication-style.md Core Patterns, clear-questions.md principle 5, background-agents.md "Pick each agent's model by its work", whose toolbox and knowledge copies carry the same line). `maintainability-audit-brief-shape.md` keeps its destination with target `docs/run-prompts/` (REQ-732 removed `_dev/run-prompts/`). `run-prompts-not-in-root-docs.md` is D-11. Every cited `file:line` was re-read on main; two that moved were corrected (`coding-guardrails.md:130` → `:132`, `.gitattributes:45` → `:49`). The ledger cites REQ-724 at its UR-175 archive path, which this finalization creates.
- **D-13** (DECIDE & STATE): the catalog pair (`ai-reports/catalog.json`, `ai-reports/index.html`) was regenerated on main after the merge and committed as `27e77d53`, the integration tip. The qualification range, gate, review and commit hash use `e5972dc5..27e77d53`.

- **D-14** (DECIDE & STATE): review findings F1 (wrong line for the Lessons phase) and F2 (report folder that moved) were marked report only, but both are fixed on the builder branch (`dd970843`) and re-merged (`406d32aa`), because the maintainer approves this page row by row and a dead path fails "name where". The integration tip, qualification range and commit hash move to `e5972dc5..406d32aa`.

- **D-15** (DECIDE & STATE): the builder's proposed lesson (memory notes that wait on a REQ go stale) gets no lessons-satellite bullet, because no satellite covers memory notes; it stays in this REQ's Lessons Learned, and `do-work/lessons-index.md` is unchanged.

<!-- D-XX counter: last used D-15. Next decision: D-16. -->

## Implementation Summary

**Files changed:**
- `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/index.html` (new; corrected by the integrator per D-11 and D-12)
- `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/presentation.css` (new; kit copy)
- `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations/presentation.js` (new; kit copy)
- `ai-reports/catalog.json` (modified; regenerated by the integrator, D-13)
- `ai-reports/index.html` (modified; regenerated by the integrator, D-13)

**What was done:** Wrote one proposal report that sorts the 46 notes in the maintainer's local memory folder into move to the skill (5), move to the repository (13) and delete (28), each row with what the note says, a reason, and a target file or an "already recorded at" path checked against the repository. The report prices three options and prints capture lines without running them. The memory folder was not written.

## Discovered Tasks

(from the builder hand-back)

- 11 of 18 crew-member files in `skills/do-work-toolbox/crew-members/` differ from their `skills/do-work/crew-members/` originals (`cmp`): anti-slop, background-agents, clear-questions, coding-guardrails, debugging, general, maintenance, prompt-injection, security, supply-chain, testing. The gap was noted on 2026-10-23 and nothing checks it. REQ-724 edits four of these core files, so the drift grows. → report only
- `skills/do-work/actions/commit.md` commit examples stage with `git add` and then run `git commit` without a pathspec; in a shared checkout that commits another session's staged files. Already in the report as the D-05 move. → report only
- `release --manifest` without `--dry-run` or `--commit` reportedly writes the tree (memory note, 2026-09-05). Not re-run here. If still true, a bare call that mutates is a CLI design trap worth a fix. → report only

## Qualification

**Gate records:** `advance --diff-range e5972dc5..27e77d53` returned phase `qualify`, gate record `qualify` satisfied, no findings. The review fix (D-14) moved the tip to `406d32aa`; the delta `27e77d53..406d32aa` changes three lines of the report page only.

**Requirement trace against `git diff e5972dc5..406d32aa` (5 files, +691):**
- Every note judged: 46 `data-memory-note` rows, one per note file in the live folder (46 notes plus `MEMORY.md`); the GREEN probe counts the folder at run time and passes. Destinations 5 move to the skill, 13 move to the repository, 28 delete.
- Each row has five cells: file name, one line on what it says, destination, one line of reason, and a target. Every move row's target starts "Move to" and names a file; every delete row names an "Already recorded at" path, or says "Nothing to keep" for a note true only for one past run (checked by script over all 46 rows).
- Checked against the current repository: the builder searched at `df1bf423`; the integrator re-read every cited `file:line` at `3d0c483c` and fixed two that moved (D-12). The REQ-724 and REQ-732 rows now cite what landed (D-12).
- One report through the normal report path: `ai-report --kind proposal` bundle under `ai-reports/`, now listed in `ai-reports/catalog.json` (D-13).
- No note moved or deleted by this REQ: the builder's SHA-256 snapshots match (hand-back Verification 3); the integrator wrote nothing under the memory folder. The one deletion (D-11) was the REQ-732 integrator's.
- Not a release: `git diff --stat e5972dc5..406d32aa -- skills/ suite/ do-work/ _dev/ CLAUDE.md` is empty; `/ai-reports` is export-ignored.

**Scope:** Route A; the diff touches only the new bundle and the catalog pair. The catalog pair is the integration seam the hand-back named. No debug artifacts, no `FILL:` marker.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` at `27e77d53d6f98559cfe1eef8c8e0f03c86619efa` (1-minute load 4.84 at start) and again at the final tip `406d32aa2db7ec68e7d25d8dbfe6d53574d16fbd` after the review fix (load 2.60), both in the detached drain-head checkout with `DO_WORK_FAST_STAGE_REUSE=off`; GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-725-probe.sh` through `advance --probe-file`; `ai-report-index --check` on the new bundle before and after the catalog rebuild.
**Result:** ✓ Gate exit 0 both times (126 s, 123 s), "Maintainer verification passed."; the green-gate record was re-recorded at `406d32aa` with `record-green-gate`. Probe exit 0 (`REQ-725 probe PASS: 46 notes, one row each`); gate records `run-blocked-check` and `green-gate` satisfied. Catalog check before the rebuild printed `ai-reports/2026-11-15_0137_proposal-REQ-725-memory-note-destinations is not in ai-reports/catalog.json yet: run do-work-toolbox ai-report index to rebuild the catalog.`; the rebuild printed `ai-reports/catalog.json: 50 bundles`; the check after it printed nothing (exit 0).

**Red-green validation:** (from the builder hand-back)
- `REQ-725-probe.sh`: ✗ at base `df1bf423` ("expected exactly one ai-reports/*_proposal-REQ-725-*/index.html, found 0", exit 1) → ✓ at builder tip `f90c21fa` (47 notes). After the memory folder lost one note, the probe failed against the uncorrected report by design (D-03) and passes at `189da332` and at the integration tip (46 notes).

**New tests added:**
- None (report-only REQ; the run probe is run evidence, not a shipped test).

**Heavy verification plan:** range `e5972dc528f70cb4e7a8a12753b8b3d8a69efc39..27e77d53d6f98559cfe1eef8c8e0f03c86619efa`, re-planned and re-drained at `406d32aa` (see Heavy Verification Plan); coverage uncertain for all five changed `ai-reports/` paths (no lane covers them), so every lane is selected:
- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage uncertain
- `queue-kanban-browser` — `… --heavy-lane queue-kanban-browser` — coverage uncertain
- `do-work-cli-integrations` — `… --heavy-lane do-work-cli-integrations` — coverage uncertain
- `staged-skills` — `… --heavy-lane staged-skills` — coverage uncertain
- `updater` — `… --heavy-lane updater` — coverage uncertain
- `installer` — `… --heavy-lane installer` — coverage uncertain

## Review

**Overall: 92%** | 2026-11-15T03:04:22Z (opus reviewer, Route A quick scan with a 16-row sample; full report do-work/runs/work-2026-11-15-012821/REQ-725-review.md)

| Dimension | Score |
|-----------|-------|
| Requirements | 96% |
| Code Quality | 88% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

Code Quality here means citation accuracy: 2 wrong citations in 16 sampled rows.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:**
- F1: the `record-lessons-after-fixes.md` row cites `skills/do-work/actions/work.md:80` as the "Lessons phase", but line 80 is the Required-lessons consult (reading). The lesson-writing phase is Step 7.5 at `work.md:108`. The destination is unchanged. impact-negligible → report only
- F2: the `velocity-lessons-report.md` row cites `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/`, which does not exist. The folder moved to `ai-reports/2026-09_18-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/` in `ffc2728d` (2026-11-04). The destination is unchanged. impact-negligible → report only
- F3: the evidence ledger cites `do-work/archive/UR-175/REQ-724-move-global-standing-rules-into-skill.md`, which exists only after the UR-175 finalization. Today the file is `do-work/archive/REQ-724-move-global-standing-rules-into-skill.md`. If finalization does not create the UR-175 folder, the citation is wrong. impact-negligible → report only

**Sampled rows:** S1-S3 (move to the skill), R1-R5 (move to the repository), D1-D8 (delete; D8 covers four notes). 14 correct, R4 has a wrong line (F1), D4 has a wrong path (F2). Every sampled destination is a fair judgment of its note.
**Anti-bloat:** 0 unrequested helpers, options, files or tests.
**Acceptance:** Pass — 46 rows for 46 notes, counts 5 + 13 + 28, probe PASS, no shipped file in the range, catalog lists the bundle.
**Restatement sweep:** not run (not the last successful integration of the run).
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

**Integrator action on the findings:**
- F1 and F2: fixed anyway (D-14). The page is a brief the maintainer approves row by row, and F2 named a folder that does not exist. Builder branch `dd970843`, re-merged as `406d32aa` with the same pre `e5972dc5`; the catalog rebuild at the new tip changed nothing. The delta is two citation cells and one ledger clause, so no re-review; gate and heavy lanes re-run at `406d32aa`.
- F3: no change. This finalization closes UR-175 and moves REQ-724 into `do-work/archive/UR-175/`, which makes the cited path true.

## Heavy Verification Plan

Base `e5972dc528f70cb4e7a8a12753b8b3d8a69efc39`, target `406d32aa2db7ec68e7d25d8dbfe6d53574d16fbd`. Coverage is uncertain for every changed `ai-reports/` path, so all six lanes are selected (reason for each: coverage uncertain for the bundle files and the catalog pair):
- `queue-kanban-javascript` — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- `queue-kanban-browser` — `… --heavy-lane queue-kanban-browser`
- `do-work-cli-integrations` — `… --heavy-lane do-work-cli-integrations`
- `staged-skills` — `… --heavy-lane staged-skills`
- `updater` — `… --heavy-lane updater`
- `installer` — `… --heavy-lane installer`

## Heavy Verification Result

Target and execution revision `406d32aa2db7ec68e7d25d8dbfe6d53574d16fbd` (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set). All lanes executed, exit 0, no findings:
- `queue-kanban-javascript` — executed, 6 s
- `queue-kanban-browser` — executed, 74 s
- `do-work-cli-integrations` — executed, 57 s
- `staged-skills` — executed, 31 s
- `updater` — executed, 66 s
- `installer` — executed, 26 s

The same six lanes also passed at the superseded tip `27e77d53` (7, 75, 58, 31, 66, 27 s) before the review fix.

## Lessons Learned

**What worked:**
- The GREEN probe counted the live memory folder, not a number written at capture. When a sibling integrator deleted one note, the probe failed against the old report, which is how the stale row was found (D-03, D-11).
- A script over all 46 rows checked that each has five cells and that its target cell starts with "Move to", "Already recorded at" or "Nothing". That made the "every row has a target" requirement cheap to check.

**What didn't:**
- The row re-check covered every cited `file:line` and missed one cited folder that a 2026-11-04 reorganisation had moved (review F2). An existence check over every backticked path in the page, run at the integration tip, finds it in one command; run it before review, not after.

**Worth knowing:**
- (from the builder hand-back) A memory note that names a REQ "until it ships" goes stale silently: 3 of the 28 deletes (`fast-gate-caches-since-0-304`, `gate-hold-route-to-repair-req`, `dowork-followup-batch-resume`) were notes waiting on a REQ that had since shipped or been dropped.
- No lessons-satellite bullet: the builder's lesson is about memory notes, which no satellite covers (action files, shell commands, releases, CLI), so it stays here (D-15).

## Orientation

Report-only work under `ai-reports/`: one `ai-report --kind proposal` bundle plus the regenerated catalog pair; no `skills/` or `_dev/` change. The approved moves are later REQs.

## Timing

Observed 2026-11-14T23:29:05Z to 2026-11-15T03:12:04Z: 3h 42m 59s total, 3h 55m 17s attributed across 8 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 3h 23m 43s | 1 |
| verification-gate | 25m 08s | 4 |
| review | 3m 53s | 1 |
| planning | 2m 20s | 1 |
| handback-merge | 13s | 1 |

Slowest stage: builder-work / builder worktree build, 3h 23m 43s, outcome success.
