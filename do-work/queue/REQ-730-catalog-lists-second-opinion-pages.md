---
id: REQ-730
title: 'Report catalog lists the second-opinion pages of each bundle'
status: pending
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-728, REQ-729, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "CHANGELOG.md", "VERSION"]
---
# Report Catalog Lists the Second-Opinion Pages of Each Bundle
## What
Under each bundle's row, the report catalog (`ai-reports/index.html` and `ai-reports/catalog.json`) lists that bundle's second-opinion pages with tree id, preset, model and date.
## Why
The catalog has one entry per bundle, so second-opinion pages cannot be found through it. The bundle `ai-reports/2026-11-14_1547_proposal-tinka-review-surface` holds `index.html` and three second-opinion pages (`v1c01`, `v1c02`, `v1c03`), and the catalog shows only the first.
## Context
- One entry per bundle: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`, around lines 166 to 172, `pickReportEntry`.
- Each second-opinion page's head already carries its tree id, preset, model and date as `ai-report-second-opinion-*` meta tags.
## Detailed Requirements
- `catalog.json` gains a list of second-opinion pages per bundle. The change is additive.
- `ai-reports/index.html` shows the pages under their bundle's row.
- A page without the tree meta tags (a plain second-opinion page) is still listed, by its file name.
- Files: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` and its test file; the catalog forms in `skills/do-work-toolbox/actions/ai-report.md`; "Report Catalog" in `skills/do-work-toolbox/docs/ai-report-guide.md`; `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
## Constraints
- Every existing field of `catalog.json` keeps its name and meaning, so `revise`, `find`, `organize` and `second-opinion latest` keep working.
- No bundle is edited. The catalog reads bundles and never writes into them.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** A Go test builds a fixture bundle with `index.html` and one `index-v1c01-*.html` page and asserts that the catalog entry lists that page.
**Why RED now:** The entry has no such list.
**GREEN when:** The test passes; `ai-reports/index.html` shows the pages under their bundle; `go test ./...` passes in the do-work-cli module; the existing catalog tests pass unchanged.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Builder Guidance
- Certainty: firm on the outcome. The field name and the HTML shape are the builder's choice.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes do-work-cli Go code for reports.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md*
