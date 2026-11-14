---
id: REQ-729
title: '[impact-negligible] ai-report-judge resolves or documents a relative --out'
status: pending
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
required_lessons: ["_dev/primes/lessons-releases.md"]
related: [REQ-726, REQ-727, REQ-728, REQ-730, REQ-731, REQ-732]
batch: tinka-follow-ups
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge_test.go", "skills/do-work-toolbox/actions/ai-report.md", "CHANGELOG.md", "VERSION"]
---
# ai-report-judge Resolves or Documents a Relative --out
## What
Make the rule for a relative `--out` of `ai-report-judge` visible: either resolve it against the current folder, or state in the usage text and in the action that it is resolved against `--repo-root`.
## Why
On 2026-11-14, with `--repo-root` pointing at a scratch copy, a relative `--out` put the captures inside the copy and not where the caller expected. The usage text does not say which folder a relative path starts from.
## Context
- The option is parsed in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` (around lines 94 and 95) and the result is written with `filepath.Join(output, ...)`.
## Detailed Requirements
- Choose one: resolve a relative `--out` against the current folder, or keep today's behavior and state the rule in the usage text and in the judge form of `skills/do-work-toolbox/actions/ai-report.md`.
- An absolute `--out` behaves exactly as today.
- When the behavior changes, a Go test pins it. When only the text changes, no test is needed.
## Constraints
- No caller that passes an absolute path changes.
- The judge never writes inside the bundle.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** Run the judge with `--repo-root <scratch copy>` and `--out captures` from another folder, and read the usage text.
**Why RED now:** The captures land under the scratch copy, and nothing in the usage text or the action says that a relative `--out` starts at `--repo-root`.
**GREEN when:** Either the captures land under the current folder and a Go test pins it, or the usage text and the action state the rule in one sentence each.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Builder Guidance
- The choice is the builder's, after reading the existing callers. If a shipped caller depends on today's behavior, keep the behavior and document it.
- Record the choice as a decision in the request.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget and `slugged: partial`, so it cannot be narrowed. Matching reason: the request changes do-work-cli Go code for reports.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-judge-renders-one-named-page.md*
