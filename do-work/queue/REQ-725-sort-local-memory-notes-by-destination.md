---
id: REQ-725
title: 'Sort the local memory notes by destination'
status: pending
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-move-global-rules-into-skill.md*
