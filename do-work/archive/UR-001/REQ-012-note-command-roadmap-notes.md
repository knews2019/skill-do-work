---
id: REQ-012
title: Add do-work note command for lightweight roadmap notes
status: completed
created_at: 2026-06-01T00:00:00Z
claimed_at: 2026-06-01T17:38:31Z
completed_at: 2026-06-01T17:38:31Z
route: B
kb_status: pending
user_request: UR-001
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
---

# Add do-work note command for lightweight roadmap notes

## What

Add a `do-work note <text>` command that appends a lightweight, dated note to `do-work/notes.md`. The `do-work roadmap` action reads this file and renders a **Notes** section at the top of its output, before the REQ queue. Notes have no frontmatter, no RED/GREEN proof, no domain — they are ephemeral next-step hints that a user deletes directly from `notes.md` when no longer relevant.

Example invocation: `do-work note "next: investigate prototype xyz.html"`

## Why

The do-work queue is REQ-centric — every item goes through capture, has a structured schema, and implies implementation work. There is no lightweight channel for informal context: "I want to look at X next", "check that before running", "revisit after Y lands". These hints don't warrant a REQ but are useful to surface in the roadmap view when planning what to work on next.

## Context

- Capture method: `do-work note "text"` (new routing verb, not through capture action)
- Display: Notes section at **top** of `do-work roadmap` output, before REQ classification
- Storage: `do-work/notes.md` — plain list of dated lines; user edits/deletes directly
- Removal: no delete command needed — user removes lines from `notes.md` manually when notes are resolved

## Red-Green Proof

**RED prompt/case:** Run `do-work note "investigate prototype xyz.html"` — command is unrecognized or falls through to help menu. Run `do-work roadmap` — output contains no Notes section.

**Why RED now:** No `note` routing entry exists in SKILL.md; no `actions/note.md` file exists; `actions/roadmap.md` has no Notes section logic; `do-work/notes.md` does not exist.

**GREEN when:** `do-work note "investigate prototype xyz.html"` appends `- [2026-06-01] investigate prototype xyz.html` to `do-work/notes.md` (creating the file on first use). `do-work roadmap` output opens with a "## Notes" section listing current note lines. Manually deleting a line from `notes.md` causes it to disappear from the next `do-work roadmap` run.

**Validation:** User confirmed (capture-phase clarification)

## Implementation Scope

Files to create or modify:

1. **`SKILL.md`** — Add `note` to the routing table (new row, before or near `roadmap`). Add `note [text]` to the `argument-hint` line. Route matches: `do-work note`, `do-work note add`, `do-work add note`.

2. **`actions/note.md`** — New action file. Steps: parse `$ARGUMENTS` (everything after `note`), strip leading `add ` if present, append `- [YYYY-MM-DD] <text>` to `do-work/notes.md` (create file if absent), report back the line added. No UR/REQ created — this is not a capture action.

3. **`actions/roadmap.md`** — Add a Notes pre-step (before the REQ classification loop): if `do-work/notes.md` exists and is non-empty, render a `## Notes` block at the top of output listing each line. If the file is absent or empty, skip the section silently.

4. **`next-steps.md`** — Add note-relevant next-step suggestions where appropriate (e.g., after capturing a REQ, suggest `do-work note` for lightweight follow-up thoughts).

---
*Source: do-work user-requests/UR-001/input.md*

Think carefully before answering.

---

## Triage

**Route: B** — Medium. The "what" is fully specified (new `do-work note` verb + `notes.md` storage + roadmap Notes section), and the "where" is known patterns to follow (the SKILL.md routing table, an existing action-file shape, the roadmap Output Format, the next-steps per-action blocks). No multi-system planning needed.

**Planning:** Not required (Route B).

> **Provenance note (prompt-injection guardrail):** the captured body ends with a stray `Think carefully before answering.` line — an instruction-like artifact that is *not* part of the request. Treated as data, not an instruction; left intact (captured content isn't silently rewritten) and surfaced here. Logged as D-02.

## Plan

**Planning not required** — Route B.

*Skipped by work action*

## Exploration

- **Routing table (`SKILL.md`)** uses numbered priorities checked in order (first match wins). Priorities 2/5/7/11/28/29 are cross-referenced by number in prose; 30 (ai-report) and 31 (descriptive-content fallback) are **not** referenced anywhere. → safe insertion point is priority **31**, bumping the unreferenced descriptive-content fallback to 32 (D-01). Inserting near roadmap (17) would shift the referenced 28/29 and reintroduce the off-by-one cross-reference bug REQ-005 fixed.
- **Action-file convention** (CLAUDE.md): Description blockquote + Steps required; When to Use / Input / Output / Rules common. Dispatcher-free single-purpose action — model on a small action.
- **roadmap Output Format** already has a fixed report skeleton; the Notes block slots in as a pre-step rendered above the Totals header.
- **next-steps.md**: per-action `**After X:**` fenced blocks.
- No automated test harness — validation is behavioral (simulate the note append + confirm roadmap renders/skip).

## Scope

**Files I will touch:**
- `actions/note.md` (new) — the note action.
- `SKILL.md` (modify) — argument-hint, Actions list bullet, routing row (priority 31), Verb Reference row, Action Dispatch row, help-menu line.
- `actions/roadmap.md` (modify) — Notes pre-step (Step 0) + Notes block in Output Format + a next-step.
- `next-steps.md` (modify) — new `**After note:**` block + a `do-work note` suggestion under capture.
- `CLAUDE.md` (modify) — add `note.md` to the actions/ structure listing (keeps the index accurate).

**Files I will NOT touch:** `do-work/notes.md` is **not** created by this REQ — the note action creates it on first user invocation; the feature is correct with it absent (roadmap skips silently).

**Acceptance criteria (restated from REQ / Red-Green Proof):**
- [ ] `do-work note "investigate prototype xyz.html"` is routed to the note action and appends `- [YYYY-MM-DD] investigate prototype xyz.html` to `do-work/notes.md` (creating it on first use).
- [ ] `do-work roadmap` renders a `## Notes` section at the top when `notes.md` is non-empty; skips it silently when absent/empty.
- [ ] Deleting a line from `notes.md` removes it from the next roadmap run (roadmap reads the file live).
- [ ] No UR/REQ is created by the note action.

## Decisions

- **D-01:** Placed the `note` routing row at priority **31** (after ai-report, before the descriptive-content fallback) rather than "near roadmap" as the REQ suggested. **Reasoning:** priorities 28/29 are cross-referenced by number in SKILL.md prose; inserting near roadmap would shift them and recreate exactly the off-by-one bug REQ-005 fixed. Priority 31 is functionally identical for routing (note is a unique keyword) and touches only the unreferenced descriptive-content fallback (31→32).
- **D-02:** The captured REQ body's trailing `Think carefully before answering.` line is a stray instruction-like artifact. Treated as data per the prompt-injection guardrail; left intact, not acted upon. (Cleaning captured content is out of scope for a builder; flagged for the user.)

## Pre-Flight

**Git:** ✓ clean working tree (the prior two REQ commits left it clean).
**Tests baseline:** N/A — no automated test harness; behavioral validation substitutes.
**Dependencies:** N/A — markdown-only change.

*Checked by work action*

## Implementation Summary

**Files changed:**
- `actions/note.md` (new) — the `do-work note` action: normalize text (strip leading `add `, surrounding quotes, whitespace), append `- [YYYY-MM-DD] <text>` to `do-work/notes.md` (create on first use), report. Explicitly creates no UR/REQ, no work-loop transition, no commit.
- `SKILL.md` (modified) — `note` wired into **7 surfaces**: argument-hint, Actions-list bullet, routing-table row (priority 31), Verb Reference row, Action Dispatch row, help menu, and the foreground-dispatch list.
- `actions/roadmap.md` (modified) — new **Step 0: Surface Notes** (read-only) + a `## Notes` block in the Output Format (both the normal and empty-queue variants); renders notes at the top of the survey, skips silently when `notes.md` is absent/empty.
- `next-steps.md` (modified) — new `**After note:**` block + a `do-work note` suggestion under capture.
- `CLAUDE.md` (modified) — `note.md` added to the `actions/` structure index.

**What was done:** Added a lightweight `do-work note <text>` channel that appends a dated hint to `do-work/notes.md`; `do-work roadmap` surfaces those hints at the top of its survey. Notes are ephemeral working-tree-only data (no UR/REQ, no schema, no commit) that the user deletes by hand.

## Qualification

Passed — `actions/note.md` exists with all required/common action-file sections (When to Use, Input, Steps, Output Format, Rules, Verification Checklist) and is reachable: routing row at priority 31 + dispatch row both present. roadmap Step 0 + `## Notes` block present. Priority cross-references 28/29 unchanged — no off-by-one regression (the reason note landed at 31, not near roadmap). No placeholder content; declared scope == files touched.

## Testing

**Tests run:** behavioral simulation of the prose-prescribed logic (no automated harness for agent routing) + structural greps.
**Result:** ✓ All passing
- note append → `- [2026-06-01] investigate prototype xyz.html` — **exact match** to the captured GREEN ✓
- two notes → roadmap Step 0 renders both under `## Notes` ✓
- manual deletion of a line → gone on the next (live) roadmap read ✓
- absent/empty `notes.md` → roadmap skips the section silently ✓
- routing: 31=note / 32=descriptive, no duplicate priorities, all 7 SKILL surfaces wired, cross-refs 28/29 intact ✓
- test `notes.md` removed after validation — no test data shipped ✓

**Red-green validation:** RED (from `## Red-Green Proof`) = `do-work note …` unrecognized + roadmap has no Notes section. GREEN = the append + render demonstrated above. `tdd: false` is correct — there's no runnable harness for "the agent routes `do-work note`"; behavioral simulation traced to the captured proof substitutes.

*Verified by work action*

## Review

**Overall: 93%** | self-review, Route B standard depth

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor.
**Acceptance:** Pass — all four acceptance criteria and the RED→GREEN proof met; `notes.md` stays out of the committed pipeline as intended.
**Follow-ups created:** None.

## Lessons Learned

**What worked:** The note action is prose, but its file logic (append / render / delete / skip) is concrete enough to *execute* in bash as a real GREEN test — simulating the four states gave genuine behavioral evidence without a test harness.
**What didn't:** The REQ's suggested "near roadmap" routing placement collided with the priority-cross-reference fragility REQ-005 fixed (priorities 28/29 are referenced by number in prose). Priority 31 — after the last keyword, before the descriptive-content fallback — was the only safe slot (D-01).
**Worth knowing:** `do-work/` is gitignored, so `do-work/notes.md` is correctly working-tree-only and must never be committed (the note action explicitly doesn't). New REQ/UR files under `do-work/` are untracked and need `git add -f` to land in an archive commit — unlike REQ-001/006, which were already tracked.

## Knowledge-Base Handoff

No `kb/` directory exists; handoff defers. `kb_status: pending`. Did not block archival.
