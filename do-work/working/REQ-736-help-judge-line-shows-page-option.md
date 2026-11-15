---
id: REQ-736
title: '[impact-negligible] Show the --page option on the ai-report judge line in help'
status: claimed
created_at: 2026-11-15T07:25:23Z
user_request: UR-179
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-735]
batch: run-012821-follow-ups
write_set: [skills/do-work-toolbox/actions/help.md]
claimed_at: 2026-11-15T07:30:05Z
route: A
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-11-15T07:31:33Z
  basis:
    - Route A
    - 1-file write set
    - 2 acceptance criteria
dispatch_at: 2026-11-15T07:32:53Z
---
# Show the --page Option on the ai-report Judge Line in Help
## What
`skills/do-work-toolbox/actions/help.md` line 18 lists `ai-report judge <dir>` without the `--page` option that REQ-728 (ai-report-judge renders one named page, released 0.319.0) added. Change that one line so it shows the option.
## Why
The maintainer's words: "capture a follow-up REQ for this one line". The run-end restatement sweep of REQ-734 (supports-and-limits table page) found the gap (`do-work/runs/work-2026-11-15-012821/REQ-734-integration-report.md:11`, decision D-20). The line leaves out a form but says nothing false.
## Context
- Current line 18 (verified at capture, commit `db458e9e`): `  ai-report judge <dir>          Render-check a report bundle (widths, themes, links)`.
- The shipped form to match: `judge <bundle-dir> [--page <file.html>]` in `skills/do-work-toolbox/actions/ai-report.md:36`, and `do-work-toolbox ai-report judge <dir> [--page <file.html>]  Render-check an existing report bundle, or one named page of it` in `skills/do-work-toolbox/docs/ai-report-guide.md:80`.
- The line sits in a column-aligned listing (lines 15-20). The builder keeps the listing readable if the longer usage overflows the column.
- Fold-first scan: no queued REQ or sweep shares this root cause (the queue holds only REQ-727, a relation-map preset). The fix is prose-only, but the maintainer explicitly promoted this report-only finding to a REQ, so it does not go to `do-work/prose-backlog.md`.
- Independent of REQ-735 (scope of the moved standing rules). No `depends_on`.
- The change touches a shipped file under `skills/`, so it is a release. Follow `_dev/primes/prime-releases.md`.
## Red-Green Proof
**RED prompt/case:** `sed -n 18p skills/do-work-toolbox/actions/help.md` prints `ai-report judge <dir>` with no `--page`.
**Why RED now:** REQ-728 added `--page` to `ai-report.md` and the guide, but not to the help listing.
**GREEN when:** The judge line in `help.md` shows `[--page <file.html>]`, in the same form as `ai-report.md:36`, and no other line of `help.md` changes.
**Validation:** Inferred during capture
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: item 2 of `do-work/user-requests/UR-179/input.md` (help.md fix after REQ-728)*

---

## Triage

**Route: A** - Simple

**Reasoning:** One prose line in one shipped file (`skills/do-work-toolbox/actions/help.md:18`) gains the `[--page <file.html>]` form that `ai-report.md:36` already documents. No code, no other file, no design choice beyond column alignment.

**Planning:** Not required

**Decisions:**
- **D-01**: The new usage `ai-report judge <dir> [--page <file.html>]` (41 characters) overflows the 30-character usage column. Builder keeps the description on the same line after two spaces, the same shape line 20 (`ai-report --kind proposal|root-cause <target>  Decision brief…`) already uses, and shortens nothing else. Reasoning: matches an existing overflow precedent in the same listing; no other line changes. Value: one-line diff. Risk: none; prose only, trivially reversible.

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

