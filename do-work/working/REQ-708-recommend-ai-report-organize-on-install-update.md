---
id: REQ-708
title: 'Install, update and new reports recommend ai-report organize when finished months are still loose'
status: claimed
created_at: 2026-11-02T16:55:00Z
user_request: UR-161
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-706, REQ-707]
batch: ai-report-month-folders
depends_on: [REQ-707]
required_lessons: ["skills/do-work/tools/lessons-do-work-update.md"]
claimed_at: 2026-11-02T18:40:42Z
---
# Install, Update and New Reports Recommend ai-report Organize When Finished Months Are Still Loose
## What
When a project's `ai-reports/` still holds reports from a finished month at the top level, or a month folder without a count, `do-work-toolbox install` and `do-work update` print one line recommending `do-work-toolbox ai-report organize`. `ai-report index` and every new report print the same one-line hint. Nothing prints when there is nothing to organize.
## Why
The user wants the installed repository to be told to run the organize command, and each new month leaves last month's reports loose until someone runs it again.
## Detailed Requirements
- One shared check decides whether there is anything to organize, using REQ-707's finished-month rule. It is a program, not prose repeated in each action.
- `do-work-toolbox install` and `do-work update` print the recommendation only when the check finds something.
- `ai-report index` and every new report print the one-line hint under the same condition: ai-report in every form that writes a bundle, architecture-report, and stakeholder-report.
- None of them moves a report or rewrites a link. The user runs organize.
## Constraints
- The current month's reports stay directly under `ai-reports/`. Only finished months move into a month folder, because the count is final only when the month is over. Report writers (ai-report in every form, architecture-report, stakeholder-report) keep writing new bundles directly under `ai-reports/`: no write-side change.
- A finished-month folder is named `YYYY-MM_N-reports` (`YYYY-MM_1-report` for one report). A plain `YYYY-MM` folder, as a consumer's own script makes, is read as a month folder too.
- The bundle contract is unchanged: one folder holding its entry file, the same catalog fields, and the same `yyyy-mm-dd_hhmm_<slug>` names. A month folder is only a parent directory.
- A repository with no month folders behaves exactly as before.
- Only `ai-report organize` moves reports or rewrites links. A report run never does.
- The inbox note was observed against 0.307.0, so its line numbers may have drifted.
## Dependencies
Depends on REQ-707 (ai-report organize), which the recommendation names.
## Builder Guidance
User answers: "Print a hint, you run it" for the monthly trigger, and "Install and update" for where the recommendation appears.
## Red-Green Proof
**RED prompt/case:** A fixture repository with a loose finished-month report, `ai-reports/2026-09-27_2200_a/index.html`. Run the install recommendation path, `do-work update`, and `ai-report index`.
**Why RED now:** No command checks the reports layout, so none of them mentions organize.
**GREEN when:** Each prints one line naming `do-work-toolbox ai-report organize`. In a fixture holding only current-month reports and counted month folders, none prints it.
**Validation:** User confirmed the triggers and where it prints. The fixture shape was inferred during capture.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (8398 tokens, over the 2000 budget; `slugged: partial`). Matching reason: family `alternate-writer-contract-drift`: the hint must reach every report writer (ai-report forms, architecture-report, stakeholder-report) and the install and update paths.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (23234 tokens per `do-work/lessons-index.md`, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its row covers changing do-work-cli internals; families `cross-module-fact-handoff, closed-enumeration-for-a-condition` match this work.
## Full Context
See `do-work/user-requests/UR-161/input.md` for the complete verbatim input: the user's words, the six answers chosen during capture, and the inbox note it adopts in full. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "So actually, when this is installed, It should recommend to run on AI report organize a script that will do it. For the repository. Where this is installed in."*
