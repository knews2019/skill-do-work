# Hand-back — REQ-699

Builder: orchestrator session as builder (worktree rung)
Operative name: worktree-agent-REQ-699-ai-report-commit-rule-and-judge-link-findings
Builder commit: 07e075d1
Handed back: 2026-10-19T14:35:38Z

## File manifest
- skills/do-work-toolbox/actions/ai-report.md (modified): revise step 5 wording (one finding per href), revise step 7 replaced ("Commit as the project's rules say", after the catalog regenerates, with the catalog pair), Step 7 render sentence plural, Step 8 verify sentence plural, Step 8 new commit-rule paragraph plus summary names the hash or the exact command, checklist render line plural, checklist revise line gains the commit clause.
- skills/do-work-toolbox/docs/ai-report-guide.md (modified): "committing is left to you" replaced by the project-rules sentence; nothing is ever pushed.
- Tests touched: none (prose only; contract tests exercise the shipped files).

## P-A-U
- [PLAN] Read prime-action-files, prime-releases, lessons-releases, maintenance crew, coding-guardrails. Approach: state the commit rule once in Step 8 (reached by both forms), have the revise step 7 apply it after the catalog regenerates, fix every singular "link finding" (five sites, two more than the paste named), update the guide line.
- [APPLY] Seven replacements in ai-report.md and one in the guide, each asserted unique before replacing; no other files touched.
- [UNIFY] git diff --stat: 2 files, +10/-8. Acceptance greps: "Do not commit" absent; "Commit as the project's rules say" at lines 59 and 178; "one finding per" at line 57; no singular "link finding" phrasing left; guide "left to you" absent. _dev/tests/contract-regressions.sh exit 0 (26s), _dev/tests/shipped-package-reference-contract.sh exit 0 (1s). No debug artifacts (prose).

## Decisions
- D-01 (DECIDE & STATE): the rule lives in Step 8, which both forms reach, and revise step 7 only applies it after the catalog regenerates. Reasoning: one statement of the rule instead of two copies; the revise needs its own step because its commit must follow step 6's catalog regeneration.
- D-02 (DECIDE & STATE): kept the consumer's "follow the project's instructions" form rather than "commit by default, project rules override". Reasoning: the REQ's GREEN pins the phrase "Commit as the project's rules say", and a default-commit form would change behaviour for projects whose instructions are silent, which no report asked for.
- D-03 (DECIDE & STATE): fixed the singular wording at five sites (lines 57, 170, 176, 198, 201 before the edit), two more than the paste cited, because the checklist render line also carried it.

## Required lessons
- _dev/primes/lessons-releases.md read whole-satellite (666 tokens). Nothing missing.
- _dev/primes/lessons-action-files.md stays dropped for budget (8128 tokens, partial slugging).

## Discovered Tasks
- none

## Integration seams
- none
