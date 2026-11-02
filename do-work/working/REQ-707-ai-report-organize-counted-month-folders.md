---
id: REQ-707
title: 'ai-report organize files finished months into counted month folders and fixes every link'
status: claimed
created_at: 2026-11-02T16:55:00Z
user_request: UR-161
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-706, REQ-708]
batch: ai-report-month-folders
depends_on: [REQ-706]
claimed_at: 2026-11-02T17:55:27Z
---
# ai-report Organize Files Finished Months Into Counted Month Folders and Fixes Every Link
## What
Add `do-work-toolbox ai-report organize` for the repository the suite is installed in. It moves every report from a finished month into `ai-reports/YYYY-MM_N-reports/`, keeps the current month's reports directly under `ai-reports/`, and rewrites every link to a moved report so it still opens.
## Why
The user wants the reports folder readable by month, with the report count visible in the folder name for months that are over, and wants a command to do it instead of a project-local script.
## Detailed Requirements
- A new `organize` form in `skills/do-work-toolbox/actions/ai-report.md` (with its routing and help lines), backed by one canonical do-work-cli verb. The moves and link rewrites are mechanical, so a program does them.
- A report's month comes from the `yyyy-mm-dd` date prefix of its folder name. A finished month is any month before the current local month. Current-month reports stay at `ai-reports/<bundle>/`.
- Each finished-month report moves into `ai-reports/<YYYY-MM>_<N>-reports/` (`<YYYY-MM>_1-report` when N is 1), where N is the number of report bundles in that month after the move.
- Recount on every run: a finished-month folder whose name no longer matches its bundle count, or a plain `YYYY-MM` folder, is renamed to the counted name. A month folder for the current month is emptied back to the top level.
- Fix every link to every moved or renamed report across the repository: repo-relative path mentions such as `ai-reports/<bundle>/index.html` in REQs, URs, docs, `blocked_by` lines and `## Reports` history; relative links inside reports that point at another report, such as `../<sibling>/index.html`; and the catalog. Every link that opened before opens after. Binary files are never edited.
- Regenerate `ai-reports/catalog.json` and `ai-reports/index.html` through the catalog command.
- A second run on an organized repository changes nothing.
- `--dry-run` prints the whole plan (moves, renames, and every file whose links change) and writes nothing.
- A report folder without a date prefix, a loose file, or a symbolic link stays in place and is listed in the output.
- Commit as the project's rules say, by exact paths, as ai-report Step 8 does. Never push.
- The user guide explains the layout and the command.
## Constraints
- The current month's reports stay directly under `ai-reports/`. Only finished months move into a month folder, because the count is final only when the month is over. Report writers (ai-report in every form, architecture-report, stakeholder-report) keep writing new bundles directly under `ai-reports/`: no write-side change.
- A finished-month folder is named `YYYY-MM_N-reports` (`YYYY-MM_1-report` for one report). A plain `YYYY-MM` folder, as a consumer's own script makes, is read as a month folder too.
- The bundle contract is unchanged: one folder holding its entry file, the same catalog fields, and the same `yyyy-mm-dd_hhmm_<slug>` names. A month folder is only a parent directory.
- A repository with no month folders behaves exactly as before.
- Only `ai-report organize` moves reports or rewrites links. A report run never does.
- The inbox note was observed against 0.307.0, so its line numbers may have drifted.
## Dependencies
Depends on REQ-706 (ai-report tools read report bundles inside month folders): organize writes the layout that REQ-706 teaches the readers. REQ-708 depends on this.
## Builder Guidance
The user chose: the count in the folder name for finished months, the current month flat, and this command as the only thing that moves reports or rewrites links.
## Red-Green Proof
**RED prompt/case:** A test project has 3 reports from September 2026 and 1 from the current month. One September report links to another September report, and a queue file links to one September report. Run `ai-report organize`.
**Why RED now:** The command does not exist.
**GREEN when:** The 3 September reports move into `ai-reports/2026-09_3-reports/`, the current-month report stays where it is, both links still open, and a second run changes nothing. `--dry-run` shows the plan and writes nothing.
**Validation:** User confirmed.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (23234 tokens per `do-work/lessons-index.md`, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its row covers changing do-work-cli internals; families `final-boundary-identity, exact-basename-authority, catalog-tie-order` match this work.
- `_dev/primes/lessons-action-files.md` as a whole satellite (8398 tokens, over the 2000 budget; `slugged: partial`). Matching reason: family `alternate-writer-contract-drift`: a new ai-report form whose rule must match every reader of the bundle layout.
## Full Context
See `do-work/user-requests/UR-161/input.md` for the complete verbatim input: the user's words, the six answers chosen during capture, and the inbox note it adopts in full. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "Um, moved on the folder. And if it's an older month, maybe we can also put account how many reports are in that folder. And make sure to also always fix all of the links to it."*
