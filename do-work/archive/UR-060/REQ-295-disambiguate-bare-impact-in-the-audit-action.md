---
id: REQ-295
title: "[impact-rule-change] Disambiguate the remaining bare impact wordings in the maintainability audit"
status: completed
claimed_at: 2026-08-20T22:45:21Z
completed_at: 2026-08-20T22:48:01Z
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

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Sweep both audit files for the remaining bare ranking/severity wording, name the numeric `Impact` score at every site, and run the maintainer verification gate.
- [x] **[APPLY]:** Code written exactly as planned. Scope stayed within the two maintainability-audit prose files.
- [x] **[UNIFY]:** Reviewed the 2-file diff and git diff --stat; ran the maintainer verification gate, which passed through the blanked-scan probes. No debug artifacts are present.

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

---

## Triage

**Route: A** - Simple

**Reasoning:** This is a bounded prose reconciliation in two named files. The intended numeric score is established by the surrounding contract and requires no behavioral change.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/maintainability-audit.md` (modified) — names the 1-5 `Impact` score in ranking, severity, and context-coverage wording.
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (modified) — makes the same score/token distinction in the canonical ranking and template guidance.

## Qualification

Passed — 2 project files verified, all remaining bare numeric-impact references were reconciled, and no behavior or unrelated file changed.

## Testing

- `bash _dev/tests/maintainer-verify.sh` passed, including ShellCheck, contract regressions, and shipped-reference checks.
- `rg -n -i 'impact-descending|ranked impact|severity derives from impact|order.*impact' skills/do-work-toolbox/` returns no stale bare wording.
- The acceptance grep now finds only explicit `Impact` score wording or explicit lowercase `impact:` token fields.

## Review

**Overall: 100%**

**Requirements:** Pass — every ranking/severity/context reference names the numeric `Impact` score or the lowercase `impact:` token, and both files agree.

**Code quality:** Pass — precise prose-only edits with no behavior drift.

**Test adequacy:** Pass — the maintainer gate and targeted stale-wording scans pass.

**Scope:** Pass — only the two declared files changed.

**Acceptance:** Pass — the requested wording is unambiguous and the maintainer gate exits 0.

## Lessons Learned

**What worked:** Sweeping the entire two-file contract found one additional bare context-coverage use beyond the three capture-listed sites, preventing a partial wording fix.
**What didn't:** The first targeted scan focused on the ranking phrase and would have missed the introductory severity explanation without the full-word sweep.
**Worth knowing:** `Impact` is the numeric 1-5 score; lowercase `impact:` is a routing token and must remain explicit.

## Orientation

The maintainability-audit documentation now consistently separates numeric severity/ranking from the lowercase follow-up routing token across its action and reference contract.
