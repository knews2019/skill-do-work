---
id: REQ-013
title: "forensics: detect corrections recurring across archived REQ Lessons Learned"
status: completed
created_at: 2026-06-18T23:13:21Z
claimed_at: 2026-06-28T12:44:09Z
completed_at: 2026-06-28T12:55:08Z
route: B
user_request: UR-002
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: [REQ-014]
batch: agent-maintenance-loop-integration
commit: PENDING
kb_status: pending
---

# forensics: recurring-correction detector

## What
Extend `actions/forensics.md` (read-only pipeline diagnostics) with a check that scans
the `## Lessons Learned` sections of archived REQs in `do-work/archive/` and flags any
correction or lesson theme that recurs across multiple REQs as a single harness-level
finding. Today nothing aggregates lessons across REQs — `actions/kb-lessons-handoff.md`
only pulls one REQ's lessons at a time into the KB inbox. This imports the Agent
Maintenance Loop's "the same correction across multiple runs is signal that the harness
is teaching the wrong thing" into do-work's own diagnostics.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `actions/forensics.md` `## Checks` (1–9, prose checks with inline shell + 3 severity buckets), the archive layout (loose + UR-nested), and grepped every check-list enumeration. Approach: add `### 10.` using `find do-work/archive -name 'REQ-*.md'` (recurses → both forms), heuristic theme grouping, two-tier severity (2=Info/watch, 3+=Warning/strong per D-01).
- [x] **[APPLY]:** Added `### 10. Recurring Corrections` + one `## Info` example + one `## Red Flags` entry to `actions/forensics.md`; added a row to `docs/forensics-guide.md`; added "recurring corrections" to `SKILL.md:31`. Read-only throughout.
- [x] **[UNIFY]:** `git diff --stat` → 5 files, 27 insertions, no TODO/debug added. Checks renumber cleanly 1–10. Prescribed `find` surfaces both loose and UR-nested REQs (no glob/porcelain trap). Enumerations updated: forensics-guide table + SKILL.md teaser; the `SKILL.md:223` help line and `docs/forensics-guide.md:3` summary are already terse/open-ended and left as-is.

## Why
A one-off fix is noise; a correction that recurs is a signal the harness should change,
not the next run. Surfacing recurrence turns buried per-REQ lessons into one actionable
maintenance finding.

## Detailed Requirements
- Add a new check to `actions/forensics.md` (the natural home — read-only diagnostics).
  Do NOT create a new top-level action.
- The check reads the `## Lessons Learned` section of every archived REQ under
  `do-work/archive/` — including REQs nested in `do-work/archive/UR-*/`. Use a file
  enumeration that surfaces both loose and UR-nested REQs (a porcelain/`-uall`-style
  trap does not apply here, but the loose-vs-nested split does).
- Group lessons by theme (a short normalized phrase), and report any theme that recurs
  across 2+ distinct REQs, listing the REQ IDs and a one-line theme label. Treat the
  guide's "3+" as the strong-signal threshold and "2" as a watch-level signal — see the
  Open Question.
- Output is a finding in forensics' existing report format: theme, REQ IDs, and a
  pointer ("this correction has recurred — consider a harness fix, not another per-run
  patch"). Read-only: never modifies REQ files.
- Update any enumerated list of forensics checks (in the action's own `## Checks`
  intro, `SKILL.md`, or `docs/` if one exists) so the new check is not orphaned —
  per CLAUDE.md "Closed Enumerations Go Stale."

## Constraints
- Read-only. The check must not write, move, or edit any REQ or archive file.
- Theme grouping is heuristic (string/intent match on short lesson phrasing); keep it
  simple and explainable, not an over-engineered NLP step. The agent is reading
  Markdown, not building a classifier.

## Builder Guidance
Certainty: Firm on the behavior and home (`forensics.md`); Exploratory on the grouping
heuristic and exact threshold. Keep it simple. The recurrence signal is the point — a
plain "themes seen in 2+ REQs" list is acceptable for v1.

## Red-Green Proof
**RED prompt/case:** Run forensics against this repo today. Nothing reports that
"author one canonical source, point all callers at it" recurs across REQ-009 and
REQ-011, or that "read complementary source files before editing" recurs across
REQ-008 and REQ-010. The recurrence is invisible.
**Why RED now:** No check aggregates `## Lessons Learned` across REQs; lessons are only
ever read one REQ at a time (`kb-lessons-handoff`).
**GREEN when:** Running the new forensics check against `do-work/archive/` surfaces a
"recurring corrections" finding that lists at least the two themes above with their REQ
IDs (REQ-009+REQ-011; REQ-008+REQ-010), and the run mutates no files.
**Validation:** Inferred during capture (grounded in the approved scan-ideas report).

## Open Questions
- [~] Recurrence threshold for a finding: 2+ REQs or 3+ REQs? → **D-01**: Builder chose: report a two-tier finding — `watch` at 2 distinct REQs, `strong signal` at 3+. Reasoning: this is the Recommended option and it is strictly the most informative — it surfaces both the pairs already in today's archive (REQ-009+REQ-011, REQ-008+REQ-010) and the louder 3+ clusters, so no recurrence is hidden behind a single cutoff. Value: the check demonstrably fires on the current archive (which only holds pairs) while still distinguishing low- from high-confidence recurrence. Risk: low and fully reversible — tightening to strict-3+ or making N configurable is a one-number edit to the check prose; no data model or caller depends on the threshold.
  Recommended: report 2+ as "watch" and 3+ as "strong signal" (the guide says three),
  so the current archive — which only has pairs — still demonstrates the check.
  Also: strict 3+ only (quieter, but fires on nothing in today's archive); configurable N.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->


## Assets
None.

---
*Source: scan-ideas integration report (UR-002), pick #1 of 2 — see `do-work/user-requests/UR-002/input.md`.*

Think carefully before answering.

---

## Triage

**Route: B** - Medium

**Reasoning:** The "what" is firm (add a read-only recurring-corrections check to `actions/forensics.md`) but the "where/how" needs discovery — the existing `## Checks` format, the archive enumeration that surfaces both loose and UR-nested REQs, and every place that enumerates the forensics check list. No new architecture, so not Route C.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Orchestrator-performed (Route B):

- **`actions/forensics.md`** — checks are numbered `### 1.`–`### 9.` under `## Checks` ("Run all checks in order"); the intro carries no name-list, so adding `### 10.` is self-contained. Each check is prose (agent-executed procedure) with inline shell hints (e.g. check #8 uses inline ``git rev-parse``), not fenced scripts. Severities are the three report buckets: **Critical / Warning / Info**. The `## Output Format` block shows one example finding per bucket.
- **Archive layout** — REQs live both loose (`do-work/archive/REQ-*.md`) and UR-nested (`do-work/archive/UR-001/REQ-012-*.md`). Checks #2 and #5 already phrase this as "including `UR-*/` subdirectories." `find do-work/archive -name 'REQ-*.md'` recurses by default and surfaces both in one pass — a top-level glob would drop the nested ones (the CLAUDE.md "Prescribed Shell Commands Must Surface What the Steps Consume" trap).
- **Check-list enumerations** (per CLAUDE.md "Closed Enumerations Go Stale") — grep found: `docs/forensics-guide.md` "What it checks" table (closed, lists all 9 → needs a row); `SKILL.md:31` action teaser (4-of-9, illustrative — adding the marquee new check keeps it current); `SKILL.md:223` help menu (2-item, deliberately terse — leave) and `docs/forensics-guide.md:3` summary (already ends "and other health issues" — open-ended, leave).
- **Red-Green source data** — verified in Step 6.5: REQ-009+REQ-011 and REQ-008+REQ-010 Lessons sections.

*Generated by work action (orchestrator-as-explorer)*

## Scope

**Files I will touch:**
- `actions/forensics.md` (modify) — add `### 10. Recurring Corrections` check, one `## Info` example finding, one `## Red Flags` entry
- `docs/forensics-guide.md` (modify) — add a row to the "What it checks" table
- `SKILL.md` (modify) — add "recurring corrections" to the line-31 forensics teaser
- `actions/version.md` (modify) — version bump (CLAUDE.md "Before Every Commit")
- `CHANGELOG.md` (modify) — changelog entry (CLAUDE.md "Before Every Commit")

**Files I will NOT touch:** `SKILL.md:223` help menu and `docs/forensics-guide.md:3` summary (already terse/open-ended teasers); no REQ or archive files (read-only check); no new top-level action.

**Acceptance criteria (restated from REQ):**
- [x] New read-only check added under `actions/forensics.md` `## Checks` (not a new action)
- [x] Check reads `## Lessons Learned` of every archived REQ, loose AND UR-nested
- [x] Groups by theme; reports themes recurring across 2+ distinct REQs (watch) / 3+ (strong) with REQ IDs + a harness-fix pointer
- [x] Prescribed enumeration surfaces both loose and nested REQs (no porcelain/glob trap)
- [x] Check-list enumerations updated so the new check isn't orphaned
- [x] Read-only: the check never modifies REQ files

*Generated by work action*

## Implementation Summary

**Files changed:**
- `actions/forensics.md` (modified) — added `### 10. Recurring Corrections` check, one `## Info` example finding, one `## Red Flags` entry
- `docs/forensics-guide.md` (modified) — added "Recurring corrections" row to the "What it checks" table
- `SKILL.md` (modified) — added "recurring corrections" to the line-31 forensics teaser
- `actions/version.md` (modified) — bumped 0.96.0 → 0.97.0 (CLAUDE.md commit policy)
- `CHANGELOG.md` (modified) — added the 0.97.0 "The Broken Record" entry (CLAUDE.md commit policy)

**What was done:** Added a tenth read-only forensics check that aggregates `## Lessons Learned` across all archived REQs (enumerated with `find do-work/archive -name 'REQ-*.md'`, which surfaces both loose and UR-nested files), groups lessons by heuristic theme, and reports themes recurring across 2 distinct REQs as Info/watch and 3+ as Warning/strong-signal, each with REQ IDs and a "fix the harness, not the next run" pointer. Updated the two closed/teaser enumerations of the check list so the new check isn't orphaned.

## Qualification

Passed — 5 files verified on disk via `git diff --stat` (27 insertions, no whitespace-only or debug changes); all 6 acceptance criteria traced to a concrete edit; P-A-U boxes confirmed against the diff (no TODO/console.log/debugger added); the new check is wired into `## Checks` (#10) and both surviving enumerations. Read-only contract held — no REQ/archive files written.

## Testing

**Tests run:** No automated suite (Markdown skill repo). Verification = the REQ's `## Red-Green Proof`, executed manually by running the new check's prescribed procedure against `do-work/archive/`.

**Red-green validation:** *(traces to `## Red-Green Proof`)*
- RED (before): forensics had no check aggregating Lessons across REQs — recurrence was invisible. ✗
- GREEN (after): ran `find do-work/archive -name 'REQ-*.md'` → all 12 REQs enumerated, **including the UR-nested `UR-001/REQ-012`** (a top-level glob would have dropped it — confirms the loose-vs-nested handling). Read the Lessons sections and confirmed both proof themes recur: "author one canonical source, point all callers at it" across **REQ-009 + REQ-011** (watch); "read the full/complementary context before editing" across **REQ-008 + REQ-010** (watch). Run mutated zero files. ✓

**New tests added:** None — the deliverable is an agent-executed diagnostic procedure, not code with a test harness. The Red-Green Proof above is the executable acceptance check.

*Verified by work action*

## Review

**Acceptance: Pass · Overall: 97%** (independent reviewer, pipeline mode)

Requirements coverage: all 6 acceptance criteria met and verified. The reviewer independently ran the prescribed `find` (12 REQs incl. UR-nested `UR-001/REQ-012`; `ls` glob drops it → trap confirmed real), verified both proof themes are present in the live archive, confirmed read-only (stated 3×), checked the version ordering (0.97.0 > 0.96.0) and codename uniqueness ("The Broken Record" ×1), and confirmed clean scope (exactly the 5 declared files, no debug artifacts).

**Findings:**
- *Important:* none.
- *Minor:* `actions/forensics.md` Output Format shows only an Info-tier (2-REQ) example, no Warning-tier (3+) example. Reviewer marked optional ("acceptable under the existing one-example-per-bucket convention"). Left as-is — today's archive has no real 3+ cluster, and the two-tier rule is already stated in the check body. Report-only, no follow-up.
- *Nit:* the REQ-008+REQ-010 proof pair ("read complementary context before editing") is a looser heuristic fit than REQ-009+REQ-011; not a defect — grouping is explicitly heuristic and a "degenerate grouping" Red Flag ships with the check. Noted for traceability.

**Follow-ups created:** none.

*Reviewed by independent review-work agent (pipeline mode)*

## Lessons Learned

**What worked:** Validating the Red-Green Proof by *actually running* the prescribed `find` against the live archive — not just reasoning about it — proved both the loose-vs-nested handling and that the documented glob trap is real (`ls` drops `UR-001/REQ-012`). Empirical proof beat prose assertion.
**What didn't:** The second proof pair (REQ-008+REQ-010) is a looser theme match than the first; "read the whole list before editing" and "read complementary source files before editing" only group together under a broad "read full context first" theme. The heuristic grouping is doing real work here — worth keeping the degenerate-grouping Red Flag.
**Worth knowing:** `do-work/` is in `.git/info/exclude` (local), so queue/working files are git-ignored and archived REQs must be force-added (`git add -f`) to enter the commit — same pattern the prior archived REQs followed. Future forensics checks that read across REQs should reuse `find do-work/archive -name 'REQ-*.md'` rather than a glob.

## Orientation

Forensics can now detect when the *same* lesson recurs across archived REQs (check #10, `actions/forensics.md`) — a maintenance-loop signal that the harness, not the next run, should change. Leaf-level addition to the existing read-only diagnostics; no map change. `prime_files` empty, so no prime staleness to spot-check.
