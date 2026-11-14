---
id: REQ-731
title: 'ai-report prints the catalog rebuild command when a new report is missing from the catalog'
status: claimed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-728, REQ-729, REQ-730, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "CHANGELOG.md", "VERSION"]
claimed_at: 2026-11-14T23:28:19Z
---
# ai-report Prints the Catalog Rebuild Command When a New Report Is Missing From the Catalog
## What
After `ai-report` creates a report, its closing summary prints the one command that rebuilds the catalog, when the new bundle is not in `ai-reports/catalog.json`.
## Why
Ordinary creation writes only the bundle, so a new report is missing from the catalog until someone rebuilds it. On 2026-11-14 two new proposal bundles were absent from `catalog.json` (46 bundles, no match for either), and they still are.
## Context
- The revise form rebuilds the catalog and commits it with the bundle. Ordinary creation does not.
- `skills/do-work-toolbox/actions/completed-work-presentation-reference.md`, Collision-Safe Publication, was read at capture. It forbids overwriting existing artifacts and gives no reason against a reminder.
## Detailed Requirements
- Step 8 of `skills/do-work-toolbox/actions/ai-report.md` checks whether the new bundle is in the catalog, and prints the `ai-report index` line when it is not.
- The check is a program step, not an agent's guess. Use an existing command when one can answer it, in the way Step 8 already runs `ai-report-organize --check`.
- The line is printed just before the final bundle-path line, like the organize line.
- The guide's "Report Catalog" section says that a new report enters the catalog at the next rebuild and that the summary reminds you.
## Constraints
- Creating a report does not rebuild or commit the catalog. That was the maintainer's choice at capture.
- No bundle is edited.
- When the catalog files do not exist yet, the same line is printed.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** Create a report in a project whose catalog exists, and read the closing summary.
**Why RED now:** The summary does not mention the catalog, and the new bundle is absent from `catalog.json`.
**GREEN when:** The summary prints the rebuild command when the new bundle is absent from `catalog.json`, and prints nothing about the catalog when the bundle is in it.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Open Questions
- [x] How should a new report get into the catalog: by a printed reminder, or by a rebuild on every new report? → Print a reminder (Recommended).
  Note (2026-11-14): The maintainer chose the reminder: no extra file changes, so two sessions that create reports at the same time cannot collide. The accepted risk: the catalog stays stale until someone runs the command. A rebuild on every new report is out of scope.
## Builder Guidance
- Certainty: the choice is the maintainer's. How the check is done is the builder's.
- If the check needs a new flag on an existing command, that is Go code with a Go test; say so in the request and add the files to the scope.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md*
