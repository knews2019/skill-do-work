---
id: REQ-011
title: "Code review: write docs guides for dream/ai-report/slop-check (or extend exemption list)"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T18:56:00Z
completed_at: 2026-05-29T18:58:30Z
route: C
review_generated: true
source: code-review
scope: docs/, actions/dream.md, actions/ai-report.md, actions/slop-check.md, CLAUDE.md
---

# Code Review Fix: Write docs Guides for dream/ai-report/slop-check

## What

Three sizeable actions have no `docs/<action>-guide.md` and are **not** on CLAUDE.md's exemption list:

- `actions/dream.md` (244 lines, **destructive** four-phase memory consolidation)
- `actions/ai-report.md` (312 lines, screenshot + SVG callout + HTML rendering)
- `actions/slop-check.md` (172 lines, anti-slop validation + optional rewrite)

CLAUDE.md's exemption list enumerates: `install`, `tutorial`, `scan-ideas`, `deep-explore`, `pipeline`, `clarify`, plus reference-only `kb-lessons-handoff`. The three above are none of those, and they're not "small/self-explanatory" by length (172–312 lines).

`dream` is the highest priority — it's destructive, the action file is dense procedural prose, and the similar-risk actions `cleanup`, `stray-check`, and `forensics` all have guides.

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`). Also relates to Consistency-F9 (silent intent gap: "User-facing walkthrough" tail clause is missing without explanation).

## Requirements

Pick one of two paths (per finding, individually OR uniformly):

**Path A — Author the missing guides** (priority order: dream → ai-report → slop-check):
- `docs/dream-guide.md` — user-facing walkthrough of the four phases, when to run, what to expect, how to interpret the worklist (and post-REQ-008, the Phase 2.5 confirmation gate).
- `docs/ai-report-guide.md` — when to use vs `present-work`, the screenshot + SVG callout pipeline, before/after toggle expectations.
- `docs/slop-check-guide.md` — when to run, the seven anti-slop principles in plain language, how the rewrite mode interacts with user confirmation.
- Each guide should follow the precedent set by `docs/stray-check-guide.md`, `docs/cleanup-guide.md`, `docs/forensics-guide.md`.
- Add the `User-facing walkthrough: [\`docs/<name>-guide.md\`](../docs/<name>-guide.md).` tail clause to each action file's description blockquote.

**Path B — Extend CLAUDE.md's exemption list** with rationale:
- Add `dream`, `ai-report`, `slop-check` to the exempt-actions enumeration.
- Document the rationale for each (e.g., "ai-report's action file is itself walkthrough-shaped; dream's prose includes user-facing rationale inline"). Honest answer is probably "we haven't written them yet" — the exemption list should reflect intent, not gaps.

Path A is preferred for `dream` given the destructive risk.

## Acceptance

- Each of the three actions either has a guide in `docs/` OR is named in CLAUDE.md's exemption list with rationale.
- If a guide is authored, the action's description blockquote includes the standard `User-facing walkthrough:` link.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Finding: `test-coverage.md` Important-3

---

## Triage

**Route: C** — Complex (three new doc files + three action-file blockquote edits; user-chosen path between Path A / Path B determines scope)

**Reasoning:** Real choice on offer (write guides vs. extend exemption list, mix-and-match per action) — asked the user. Confirmed Path A for all three. That makes scope: three new ~80-100-line guide files + three description-blockquote edits. Each guide follows the precedent of `docs/stray-check-guide.md`.

**Planning:** Each guide gets its own composition pass (philosophy → what it operates on → flow / sections → output → key rules → usage → when NOT). Anti-slop applied inline — the guides themselves are human-facing artifacts.

## Decisions

- **D-01:** Per the AskUserQuestion answer, chose **Path A for all three**, not the recommended mixed path (Path A for dream, Path B for the other two). The user wanted full documentation coverage.

## Exploration

- `docs/stray-check-guide.md` (61 lines), `docs/cleanup-guide.md` (50 lines), `docs/forensics-guide.md` (43 lines) — confirmed the precedent shape: short summary + "Not to be confused with" + what-it-does table + output + rules + usage + when-NOT.
- `actions/dream.md` (now post-REQ-008, includes Phase 2.5) — re-read to make sure the guide reflects the consent gate accurately.
- `actions/ai-report.md` Philosophy and Input sections — confirmed the screenshot/SVG-callout/fallback model.
- `actions/slop-check.md` and `crew-members/anti-slop.md` — pulled the seven principles verbatim for the guide.

## Scope

**Files I will touch:**
- `docs/dream-guide.md` (new) — four-phase walkthrough including Phase 2.5 consent gate
- `docs/ai-report-guide.md` (new) — screenshot/SVG/Mermaid fallback pipeline
- `docs/slop-check-guide.md` (new) — seven principles, rewrite mode
- `actions/dream.md` (modify) — append `User-facing walkthrough:` link to description blockquote
- `actions/ai-report.md` (modify) — same
- `actions/slop-check.md` (modify) — same

**Files I will NOT touch:** CLAUDE.md (the docs/ exemption list described actions that lack guides; dream/ai-report/slop-check were never named in it, so no change needed once they have guides).

**Acceptance criteria (restated from REQ):**
- [x] Each of dream / ai-report / slop-check has a guide in `docs/`.
- [x] Each action's description blockquote includes the `User-facing walkthrough:` link.

## Implementation Summary

**Files changed:**
- `docs/dream-guide.md` (new, ~85 lines) — describes the four phases as a stepped table, explicitly explains the Phase 2.5 consent gate (preview → `Apply these N fixes? [all / dry-run / specific clusters / none]` → ambiguous-defaults-to-dry-run), enumerates what Phase 3 writes, lists key rules and when-NOT-to-use cases. Disambiguates from cleanup / stray-check / bkb commands.
- `docs/ai-report-guide.md` (new, ~75 lines) — describes the HTML-plus-assets output shape, the graceful-degradation pipeline (live screenshots → saved before/after → SVG/Mermaid fallback), SVG callout pattern, before/after toggle, anti-slop-inline rule. Disambiguates from present-work and pipeline's completion report.
- `docs/slop-check-guide.md` (new, ~85 lines) — table of the seven anti-slop principles in plain language, output shape with example findings table, rewrite-mode flow (offered, never auto-applied), N/A handling for principles 2 and 5. Disambiguates from code-review, ui-review, review-work.
- `actions/dream.md` (modified) — appended `User-facing walkthrough: [`docs/dream-guide.md`](../docs/dream-guide.md).` to the description blockquote.
- `actions/ai-report.md` (modified) — same.
- `actions/slop-check.md` (modified) — same.

**What was done:** Authored three user-facing guides (Path A for all three per D-01) matching the precedent set by `docs/stray-check-guide.md`. Each guide includes an "Not to be confused with" disambiguation block, a what-it-operates-on or what-it-checks table, the action's key flow, key rules, usage examples, and a when-NOT-to-use list. Each action's description blockquote now points at its guide. Anti-slop applied inline while writing (no separate slop-check pass).

## Qualification

Passed — 6 files match the declared scope. The three guides are at the precedent length (~75-85 lines, vs. 43-61 for the templates) and follow the same section ordering. The three action blockquote edits each end with the standard `User-facing walkthrough:` clause. No drift to other parts of the action files or to CLAUDE.md.

## Testing

**Tests run:** Manual cross-reference check:
- Each guide exists at `docs/<action>-guide.md`.
- Each action's blockquote contains the matching link.
- The `dream-guide.md` Phase 2.5 description matches the prose in `actions/dream.md` post-REQ-008 (same prompt wording, same ambiguous-defaults-to-dry-run rule).
- The `slop-check-guide.md` seven-principles table preserves the principles' names verbatim from `crew-members/anti-slop.md` (no paraphrasing of the principle titles).

**Result:** ✓ All links resolve. ✓ Phase 2.5 prose consistent across both files. ✓ Principle names preserved verbatim.

*Verified by work action*

## Review

**Overall: 90%** | 2026-05-29T18:58Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — three guides written, three blockquotes linked.
**Suggested testing:** A future `docs-link-audit` action (or a check inside `forensics`) could verify every `docs/<name>-guide.md` link resolves and that every non-exempt action has a guide. That would lock in the invariant rather than rely on review.
**Follow-ups created:** None — suggestion is a quality-of-life add, not Important.

*Reviewed by work action (Route C self-review)*

## Lessons Learned

**What worked:** The `docs/stray-check-guide.md` shape (summary → disambiguation → table → output → rules → usage → when NOT) ported cleanly to all three. Borrowing the structure meant the guides feel like a series, not three one-offs. Lifting the seven anti-slop principles verbatim from `crew-members/anti-slop.md` (instead of paraphrasing) kept the slop-check guide consistent with the source of truth — important because the principles list is the spec.
**What didn't:** Initially worried about repeating phase-by-phase prose from `actions/dream.md` in `docs/dream-guide.md` — pulled back to a higher-altitude phase table instead. The action file is the procedural spec for the agent; the guide is the conceptual map for the human.
**Worth knowing:** The CLAUDE.md `docs/` description doesn't enumerate which actions LACK guides — it only names the exempt ones. So adding guides for previously-undocumented actions doesn't require touching CLAUDE.md, which keeps that file as a true source-of-truth (it describes the *intent* of the layout, not a running inventory).

