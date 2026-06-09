---
id: REQ-010
title: "Code review: refresh actions/sample-archived-req.md to current schema"
status: completed
created_at: 2026-05-29T16:43:59Z
claimed_at: 2026-05-29T18:45:00Z
completed_at: 2026-05-29T18:46:00Z
commit: 17aa52e
route: B
review_generated: true
source: code-review
scope: actions/sample-archived-req.md
---

# Code Review Fix: Refresh sample-archived-req.md to Current Schema

## What

`actions/sample-archived-req.md` is the canonical example REQ — `actions/work.md:1071` explicitly points users at it. Its frontmatter is stale, missing fields the current capture/work schemas require or document as standard:

- **Missing `domain`** — `actions/capture.md:88` and `actions/work.md:153` both prescribe `domain: frontend | backend | ui-design | general`.
- **Missing `tdd`** — `actions/capture.md:90,225` prescribes `tdd:` (default `true` for runnable RED).
- **Missing `user_request`** — `actions/capture.md:87` emits it and `actions/work.md:158` documents it. The sample body says it's archived under a UR but no field carries the link.
- **Missing `## Red-Green Proof` section** — capture treats this as mandatory when `tdd: true` (L226). The sample is a route-B feature add that should illustrate it; a reader using it as a schema example wouldn't know the section exists.

Additionally, the sample doesn't model any of the newer Schema Read Contract fields (`error_type`, canonical-status fields, `qualified_at`, `qualification_attempts`, `kb_status`, `kb_entry`).

Risk: an agent reading this sample to learn the schema produces REQs that fail Step 1's selection scan (`domain` falls back to `general` with a warning per the Schema Read Contract).

## Context

Found during code review of the full repo on 2026-05-29 (run `do-work/runs/code-review-2026-05-29-161332/`). Last touched at v0.72.1.

## Requirements

- Add `domain: frontend` (or `backend`, matching the sample's content) to the frontmatter.
- Add `tdd: true` to the frontmatter, OR `tdd: false` with a comment explaining why no Red-Green Proof section is shown.
- Add `user_request: UR-NNN` linking the sample to its parent UR.
- Add a `## Red-Green Proof` section showing the contract in action (RED scenario, GREEN evidence). Skip only if `tdd: false` is justified inline.
- Show one modern field example: `kb_status: promoted` + `kb_entry: <filename>.md`, OR `qualified_at` + `qualification_attempts`, OR a failed-then-resolved `error_type` example.
- Consider renaming to `actions/sample-archived-req-example.md` OR moving to `docs/examples/` (Architecture-F9 nit — optional, but worth doing together with the refresh).

## Acceptance

- `actions/sample-archived-req.md` frontmatter has all current canonical fields.
- The sample includes a `## Red-Green Proof` section OR an explicit justification for its absence.
- The sample illustrates at least one of the newer Schema Read Contract fields.

## Source

Code review run: `do-work/runs/code-review-2026-05-29-161332/`
Findings: `test-coverage.md` Important-1, `test-coverage.md` Minor-4

---

## Triage

**Route: B** — Medium

**Reasoning:** Clear outcome (refresh sample to current schema) with multiple specific edits, but the target file already exists and the requirements name the missing fields explicitly. No planning needed; light exploration to confirm field set against capture.md and work.md's schema docs.

**Planning:** Not required.

## Exploration

Cross-referenced the sample against the schema sources of truth:

- `actions/work.md:147-176` — full frontmatter contract (canonical key list)
- `actions/capture.md:111-115` — `## Red-Green Proof` section template (RED/Why RED now/GREEN when/Validation)
- `actions/capture.md:225` — `tdd: true` default heuristic and the Step 6 hard-gate behavior in work.md
- `actions/capture.md:226` — Red-green proof inference heuristic

The sample already has `kb_status: promoted` + `kb_entry`, which satisfies the "one newer Schema Read Contract field" acceptance criterion. The Optional rename to `sample-archived-req-example.md` (or `docs/examples/`) was declined for this REQ — it requires updating `actions/work.md:1071`'s cross-reference and adds scope. Filename stays.

*Exploration via direct file reads (no agent spawn — schema sources are explicit, no codebase pattern search needed).*

## Scope

**Files I will touch:**
- `actions/sample-archived-req.md` (modify) — add three frontmatter fields + one new section + red-green validation rows in Testing

**Files I will NOT touch:** `actions/work.md` (the cross-reference at L1071 stays unchanged — the rename is optionally scoped out), `actions/capture.md` (schema source, no change needed).

**Acceptance criteria (restated from REQ):**
- [x] `actions/sample-archived-req.md` frontmatter has all current canonical fields.
- [x] Sample includes a `## Red-Green Proof` section (or justified absence).
- [x] Sample illustrates at least one newer Schema Read Contract field (`kb_status` + `kb_entry`).

## Pre-Flight

(Skipped — no automated test suite for prose files. Acceptance is a manual schema-completeness check.)

## Implementation Summary

**Files changed:**
- `actions/sample-archived-req.md` (modified) — added `domain: frontend`, `tdd: true`, `user_request: UR-004` to frontmatter; inserted a new `## Red-Green Proof` section after `## What`; added a `Red-green validation:` block inside the existing `## Testing` section.

**What was done:** Brought the canonical sample REQ in line with the current capture/work schemas. The frontmatter now demonstrates all required canonical fields including the `tdd` gate, and the body illustrates the mandatory `## Red-Green Proof` contract (RED prompt/case → Why RED now → GREEN when → Validation) plus the corresponding `Red-green validation:` evidence trail in the Testing section.

## Qualification

Passed — 1 file modified per declared scope. Three frontmatter additions match the keys documented in `actions/work.md:147-176`. The `## Red-Green Proof` section matches the template at `actions/capture.md:111-115`. Two `Red-green validation:` rows in Testing trace the captured RED/GREEN pair to specific test cases, per the `tdd: true` mandate at `actions/work.md` Step 6.5. No scope drift — `actions/work.md:1071` cross-reference left intact.

## Testing

**Tests run:** Manual schema audit — every canonical field documented in `actions/work.md:147-176` is present or explicitly justified.
**Result:** ✓ All canonical fields present (`id`, `title`, `status`, `domain`, `tdd`, `user_request`, `created_at`, `claimed_at`, `route`, `completed_at`, `commit`, `prime_files`, `kb_status`, `kb_entry`). `caveman:` correctly omitted (default `false`). `addendum_to:`, `depends_on:`, `error`, `error_type:` correctly omitted (success path, no addendum, no deps).

**Red-green validation:** Non-behavioral text edit; no test-first evidence required. Regression evidence is the schema audit above.

*Verified by work action*

## Review

**Overall: 90%** | 2026-05-29T18:46Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | low |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — frontmatter complete, Red-Green Proof present, kb_status/kb_entry illustrate the newer Schema Read Contract.
**Suggested testing:** Add a small frontmatter-schema linter (separate REQ if anyone wants it) — there is currently no automation that would catch the sample drifting from the schema again.
**Follow-ups created:** None (the optional rename was explicitly scoped out; the linter is a suggestion, not an Important finding).

*Reviewed by work action (Route B self-review)*

## Lessons Learned

**What worked:** Reading both `actions/capture.md` (Red-Green Proof template) and `actions/work.md` (frontmatter contract) before editing — the two sources had complementary information neither alone had.
**What didn't:** Initially considered the rename (`sample-archived-req-example.md` / `docs/examples/`) but scoped it out once the work.md:1071 cross-reference came into view. Better to land the schema refresh cleanly than entangle it with a path change.
**Worth knowing:** The sample is the agent-facing schema example — when capture.md or work.md adds/renames a frontmatter field, this file should be updated in the same commit. A linter or a "schema-sample sync" check would catch future drift.

