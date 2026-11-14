---
id: REQ-728
title: 'ai-report-judge renders one named page of a bundle'
status: claimed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-729, REQ-730, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "CHANGELOG.md", "VERSION"]
claimed_at: 2026-11-14T23:28:12Z
---
# ai-report-judge Renders One Named Page of a Bundle
## What
Add an option to `ai-report-judge`, for example `--page <file.html>`, that renders one named page of a bundle in place. Today the judge renders only `index.html`.
## Why
To check or capture a second-opinion page such as `index-v1c01-<name>.html`, an agent must copy the whole bundle and write that page over the copy's `index.html`. That copy step was needed three times on 2026-11-14.
## Context
- The judge renders only `index.html`: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go`, lines 25, 111 and 168.
- `secondOpinionRender` in `report_second_opinion.go` (from line 1706) already does the copy-and-swap internally. Reuse or replace it as the builder judges.
## Detailed Requirements
- The new option renders the named page of the bundle, without copying the bundle and without writing inside it.
- A named page that is missing or outside the bundle is refused with a typed finding.
- Files: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` and its test file; the judge form in `skills/do-work-toolbox/actions/ai-report.md`; `skills/do-work-toolbox/docs/ai-report-guide.md` where the judge is described; `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.
## Constraints
- Without the new option the judge behaves exactly as today: same captures, same `judge.json`, same exit codes.
- The judge never writes inside the bundle. Bundles stay byte-identical.
- The second-opinion finalize path keeps its verdicts for the same pages.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** A Go test runs the judge with `--page index-v1c01-x.html` on a fixture bundle with two pages.
**Why RED now:** The judge fails on the unknown option.
**GREEN when:** The test passes; the result's served address ends with the named page and the captures show that page; a missing or outside page is refused with a typed finding; the fixture bundle's digest is unchanged; `go test ./...` passes in the do-work-cli module.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Builder Guidance
- Certainty: firm on the outcome. The option's name is the builder's choice; `--page` is a suggestion.
- The relative `--out` rule of the same command is its own request in this batch.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 8761 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes an action file.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes do-work-cli Go code for reports.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-judge-renders-one-named-page.md*
