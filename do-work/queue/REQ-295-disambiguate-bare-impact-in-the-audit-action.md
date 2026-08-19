---
id: REQ-295
title: "[impact-rule-change] Disambiguate the remaining bare impact wordings in the maintainability audit"
status: pending
created_at: 2026-08-19T15:48:05Z
user_request: UR-060
addendum_to: REQ-289
domain: general
review_generated: true
impact: impact-rule-change
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
depends_on: []
maintenance: false
related: [REQ-289]
write_set:
- skills/do-work-toolbox/actions/maintainability-audit.md
- skills/do-work-toolbox/actions/maintainability-audit-reference.md
---

# Disambiguate the Remaining Bare "Impact" Wordings in the Maintainability Audit

## What

REQ-289 introduced `impact:` as a frontmatter token field. The maintainability audit already had a
numeric `Impact: [1-5]` score. REQ-289's D-06 decided to keep both and disambiguate in prose — a
sound call, applied at some sites and missed at three.

## Why

Where the disambiguation was applied it reads cleanly: `maintainability-audit-reference.md:183` now
says "Severity comes from the 1-5 `Impact` score alone". Three sites still say bare "impact", and
now that the word also names a four-token enum, ordering "by impact descending" is undefined — the
tokens have no order. A reader following it literally is stuck.

## Context

- `maintainability-audit.md:170` — "Classes are ranked impact-descending with effort as tie-break;
  severity derives from impact alone." Both uses bare. This line sits one below a checklist item
  REQ-289's diff edited, so it was read past rather than unreached.
- `maintainability-audit-reference.md:185` — "order classes by impact **descending**".

## Detailed Requirements

- Every remaining bare "impact" in these two files says which of the two it means: the 1-5 `Impact`
  score, or the `impact:` token.
- Ranking and severity both key on the **1-5 score**, which is what they meant before REQ-289 and
  still mean. This REQ changes wording, not behavior.
- Keep the two files consistent with each other — the action file's copy and the reference's copy
  currently disagree in wording for the same rule.

## Acceptance

- No bare "impact" in either file where the 1-5 score is meant.
- `grep -n 'impact-descending' skills/do-work-toolbox/` returns nothing, or returns only a form that
  names the score.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Full Context

Finding F7 from REQ-289's review, reproduced by the orchestrator. See
`do-work/user-requests/UR-060/input.md`.
