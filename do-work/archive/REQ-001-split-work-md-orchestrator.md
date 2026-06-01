---
id: REQ-001
title: "Code review: split actions/work.md into orchestrator + reference companion"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-06-01T17:38:31Z
completed_at: 2026-06-01T17:38:31Z
route: B
review_generated: true
source: code-review
scope: full-repo
kb_status: pending
---

# Code Review Fix: Split actions/work.md into Orchestrator + Reference Companion

## What

`actions/work.md` has grown to 1,074 lines — the longest file in the repo, on the hottest path (read on every `do-work run` and every `pipeline` build step). The body has accreted multiple sub-systems that have natural seams (Schema Read Contract, wave execution, archive-collision rules, plan validation, scope drift, qualification, TDD verify, prime-link deferral, follow-up REQ template, commit-format prose).

Split it the same way `actions/bkb.md` + `actions/bkb-reference.md` and `actions/interview.md` + `actions/interview-reference.md` are split: a focused orchestrator file that prescribes the ten-step skeleton, plus a reference companion that holds the heavy procedural content.

This refactor also resolves two related findings surfaced in the same review:

- **Consistency-F1**: `## Red Flags` (L591) and `## Verification Checklist` (L600) currently sit mid-document inside Step 6, not at file-end per CLAUDE.md's prescribed section order. The file is also missing document-level `## Rules` and `## Common Rationalizations` (it has `## Common Mistakes` at L1000 instead).
- **Test-Coverage-F2 / Important-2**: `work.md` is the largest and most-invoked action but is missing `## Common Rationalizations` despite many builder shortcuts that warrant guarding against (hollow Implementation Summary, faked Qualification, skipped Pre-Flight, ticked P-A-U boxes without doing the work).

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`). This is the top recommended action in the report (Recommended Action #1). Estimated 300-400 lines off the always-loaded orchestrator.

## Requirements

- Move the following from `actions/work.md` into a new `actions/work-reference.md`:
  - Schema Read Contract table (current L178-202)
  - Crash Recovery sub-procedure (current L228-233)
  - Step 1 composed-exit-summary rendering blocks (current L269-319)
  - Step 6.3 Qualification anti-rationalization tables (current L580-589)
  - Deferred prime-link path-computation rules (current L724-735)
  - Step 8 follow-up REQ template (current L745-774)
  - Step 9 commit-format / metadata-commit prose (current L862-908)
- Update `actions/work.md` to reference the companion the same way `actions/bkb.md` does (e.g., "Heavy content … lives in the companion `actions/work-reference.md`").
- Restore the canonical section order in the new `actions/work.md`: document-level `## Rules` → `## Common Rationalizations` → `## Red Flags` → `## Verification Checklist` at file-end. Either rename the mid-document Step-6-Qualify-specific Red Flags/Checklist to H3 (`### Step 6.3 Red Flags`, `### Step 6.3 Verification Checklist`) or move them to the document-level tail.
- Add a `## Common Rationalizations` table with 6–8 rows guarding against the most common builder shortcuts. Suggested rows: "I'll skip Pre-Flight — the baseline is probably stable", "This counts as TDD — I wrote the test after but it does fail without the code", "P-A-U is bookkeeping, I'll just tick the boxes", "This file change is small, I don't need to add it to the Scope section", "Tests pass on attempt 1, I don't need to load debugging.md", "The Implementation Summary is too detailed; I'll just write 'updated logic'".
- Update any callers that reference `work.md` by step number to reference by the new section structure (see REQ-006 for the named-contracts piece — they may be done together or in sequence).
- Verify `actions/roadmap.md`'s "honor the Schema Read Contract documented in `actions/work.md`" reference still resolves (point to the new companion if appropriate).

## Acceptance

- `actions/work.md` is under 700 lines.
- `actions/work-reference.md` exists and contains the moved sections.
- `actions/work.md` follows CLAUDE.md's prescribed section order: Description → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist.
- Document-level `## Common Rationalizations` is present with 6–8 concrete rows.
- All cross-references in other action files still resolve.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Findings: `consistency.md` F1, `performance.md` F1, `test-coverage.md` Important-2

---

## Triage

**Route: B** — Medium

**Reasoning:** The "what" is fully specified (split into orchestrator + companion, exact sections enumerated, exact acceptance). The "how much" needs the discovery already done in this session (the `bkb.md`/`bkb-reference.md` split pattern is the reference shape, and the cross-reference sites are mapped). No multi-agent planning warranted — single-file authoring against a clear pattern.

**Planning:** Not required (Route B).

## Plan

**Planning not required** — Route B: Exploration-guided implementation.

*Skipped by work action*

## Exploration

**Reference pattern (`actions/bkb.md` → `actions/bkb-reference.md`):** The orchestrator keeps the procedural steps and references the companion by named section ("the … section of actions/bkb-reference.md"). The companion opens with a blockquote: "Companion file to `bkb.md`. Contains … Extracted to keep the main action file focused on procedural steps."

**Cross-reference sites that point into work.md** (must keep resolving after the split):
- `actions/roadmap.md` — "honor the Schema Read Contract documented in `actions/work.md`" (×4 mentions). The Schema Read Contract moves to the companion → update roadmap's pointer to name the companion.
- `actions/capture.md:190` — "the **normalize-and-warn contract** defined in `actions/work.md`'s Schema Read Contract" → same companion redirect.
- `actions/work.md`-internal Schema-Read-Contract references (Steps 1/3/4/5/6/9) — keep referencing "the Schema Read Contract" by name; it now lives in the companion.
- Step-number couplings in `kb-lessons-handoff.md`, `review-work.md`, `commit.md` → these are **REQ-006's** scope (named-contract promotion). Out of scope here, but REQ-001's split keeps Step 7.5/Step 9 as cleanly-nameable phases so REQ-006 can promote them.

**No test harness** in this repo — "tests" for this REQ are structural assertions (line count, grep for moved-section presence, cross-ref resolution).

## Scope

**Files I will touch:**
- `actions/work.md` (modify) — excise heavy templates/tables/sub-procedures, insert companion pointers, restore canonical section order, add `## Rules` + `## Common Rationalizations`.
- `actions/work-reference.md` (new) — companion holding the extracted content.
- `actions/roadmap.md` (modify) — redirect the Schema Read Contract pointer to the companion (per REQ requirement #6).
- `actions/capture.md` (modify) — redirect the one Schema-Read-Contract pointer to the companion.

**Files I will NOT touch:** `actions/kb-lessons-handoff.md`, `actions/review-work.md`, `actions/commit.md` (step-number couplings — that's REQ-006). The 10-step skeleton, all prescriptive prose, and all Step section *headers* stay in `work.md`.

**Acceptance criteria (restated from REQ):**
- [ ] `actions/work.md` under 700 lines.
- [ ] `actions/work-reference.md` exists and contains the 7 enumerated moved sections.
- [ ] `work.md` follows canonical section order: Description → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist.
- [ ] Document-level `## Common Rationalizations` present with 6–8 rows.
- [ ] All cross-references in other action files still resolve.

## Decisions

- **D-01**: The 7 sections REQ-001 enumerates total ~181 lines; moving only those leaves `work.md` at ~910 lines, which fails the `<700` acceptance — the Requirements list and the Acceptance are internally inconsistent. **Builder chose:** move the 7 enumerated sections **plus** additional reference-grade content (embedded templates, the big ASCII diagrams, the failure-classification tables, the checkpoint/progress examples) to satisfy `<700`, following the bkb-split philosophy: the orchestrator retains every prescriptive instruction + a named pointer; the companion holds the templates/heavy tables the instruction refers to. **Reasoning:** honors both the enumerated requirements (all 7 moved) and the acceptance (`<700`); the alternative (literal enumerated list only) cannot satisfy the REQ's own acceptance. No orchestration prose is deleted — only relocated.

## Pre-Flight

**Git:** ✓ clean working tree (verified before claim).
**Tests baseline:** N/A — no automated test harness in this repo; structural checks (wc -l, grep) substitute.
**Dependencies:** N/A — markdown-only change.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified) — split the 1,074-line file into a 631-line orchestrator: excised ~19 reference-grade blocks (replaced by named companion pointers), restored canonical section order, reconnected the Step-6.3 qualify-fail/pass logic (the mid-doc Red Flags/Verification Checklist that orphaned it were moved to the tail), and added document-level `## Rules` (folding the former Common Mistakes + What-This-Does-NOT) and `## Common Rationalizations` (8 rows).
- `actions/work-reference.md` (new) — 552-line companion holding 22 extracted sections (Architecture, Folder Structure, full frontmatter schema, Schema Read Contract, Crash Recovery, Composed Exit Summary, every step template, Discovered-Tasks + Failure classification, the commit/metadata-commit procedure, the checkpoint template, the progress example).
- `actions/roadmap.md` (modified) — redirected 3 Schema-Read-Contract references to the companion.
- `actions/capture.md` (modified) — redirected 1 Schema-Read-Contract reference to the companion.

**What was done:** A faithful *move-not-rewrite* split, following the `bkb.md`/`bkb-reference.md` pattern. A single deterministic Python transform excised blocks by original line range and substituted pointers; the companion was assembled from the same ranges. Verified line-by-line that every substantive original line survives in one of the two files.

## Qualification

Passed — 4 files verified on disk (`git diff --stat`: work.md −524/+81, capture/roadmap small, work-reference.md new). The companion is wired: 21/21 named pointers in `work.md` resolve to existing `## ` headers in `work-reference.md`. No placeholder/hollow content (552 lines of moved substance). All 5 acceptance criteria traced to verified output.

## Testing

**Tests run:** structural assertions (no automated test harness in this repo).
**Result:** ✓ All passing
- `work.md` = 631 lines (< 700) ✓
- `work-reference.md` exists; all 7 enumerated sections present ✓
- canonical tail order Rules → Common Rationalizations → Red Flags → Verification Checklist ✓
- Common Rationalizations: 8 data rows (target 6–8) ✓
- 21/21 companion pointer targets resolve ✓
- content-preservation: 0 unintentional substantive lines dropped (6 deltas, all intentional: heading-level changes + 2 folded sections + 2 polish rewords) ✓
- Schema-Read-Contract path references in roadmap.md + capture.md redirected to companion ✓

**Red-green validation:** N/A — non-behavioral documentation refactor; regression evidence is the content-preservation + pointer-resolution checks above.

*Verified by work action*

## Review

**Overall: 95%** | self-review, Route B standard depth

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor
**Acceptance:** Pass — the split is faithful and meets every criterion; read-path doc change with no runtime behavior.
**Follow-ups created:** None (REQ-006, already queued and dependent on this, will promote Step 7.5/Step 9 to named contracts and de-couple the remaining step-number callers).

*Reviewed by work action*

## Lessons Learned

**What worked:** Modeling the split as two range-lists over the *original* line numbers — "work.md excisions (→ pointer)" and "companion blocks (→ header)" — let one deterministic transform do the move with no line-shift bugs. A line-by-line content-preservation diff (every substantive original line must reappear in new+companion) is the right "test" for a move-not-rewrite refactor: it proved nothing was silently dropped and surfaced exactly which deltas were intentional.
**What didn't:** Taking the REQ's enumerated move-list literally would have shipped a ~910-line `work.md` that *fails its own `<700` acceptance* — the Requirements and Acceptance were inconsistent. Resolved by D-01 (move the enumerated sections plus enough additional reference-grade content to hit the target).
**Worth knowing:** `work.md`'s document-level Red Flags + Verification Checklist had drifted mid-document inside Step 6, orphaning the Step-6.3 qualify-fail/pass logic beneath them; relocating the headings to the tail reconnected it. The Schema Read Contract is referenced by-name (no path) throughout `work.md` — those internal refs stay stable after the move; only cross-file references needed a path update.

## Knowledge-Base Handoff

No `kb/` directory exists in this project, so the lessons handoff defers. `kb_status: pending` — re-run via `do-work review REQ-001` after `do-work bkb init` if these lessons should be compiled into a wiki. Did not block archival.
