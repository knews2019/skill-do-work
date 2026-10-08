---
id: REQ-644
title: '[impact-rule-change] Capture distinguishes a session earmark from work that needs the user as operator'
status: completed
estimate:
  p50_active_minutes: 5
  confidence: high
  basis:
  - trivial short-circuit
  calculated_at: 2026-10-07T23:28:10Z
route: A
created_at: 2026-10-07T23:21:03Z
user_request: UR-140
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-mechanical
related: [REQ-643]
batch: earmark-placement
write_set: [skills/do-work/actions/capture.md, skills/do-work/docs/work-guide.md]
dispatch_at: 2026-10-07T23:30:13Z
builder_handback_at: 2026-10-07T23:34:57Z
integration_at: 2026-10-07T23:35:13Z
kb_status: pending
review_at: 2026-10-07T23:43:54Z
claimed_at: 2026-10-07T23:26:38Z
commit: 3602e584d16e85822a34839df0acb803070c5f02
heavy_verified_at: 2026-10-07T23:45:03Z
heavy_verified_revision: 3602e584d16e85822a34839df0acb803070c5f02
completed_at: 2026-10-07T23:45:19Z
release_at: 2026-10-07T23:45:19Z
---
# Capture Distinguishes a Session Earmark From Work That Needs the User as Operator
## What
Capture's Earmark assessment says which of two shapes to write: `assigned_to` when another session will take the work (routing between sessions), and `status: blocked` with `blocked_by` naming the person plus `blocked_at` when the user has to be present as operator. The External-condition assessment lists the earmark as a fourth look-alike.
## Why
Today "leave this one for me" and "I have to be at the keyboard for this one" both read as earmarks. An earmark lands under Pending (today under Ready, after REQ-643 under Earmarked) and nothing releases it except a hand edit or explicit targeting. Work that needs the user belongs in Needs Input · Blocked, where `do-work clarify` already confirms the condition and releases the REQ. The maintainer had to flip two REQs by hand to get that.
## Verified Facts (from triage)
- `skills/do-work/actions/capture.md:107` (Earmark assessment) never mentions `blocked`.
- `capture.md:108` (External-condition assessment) already lists three look-alikes: a wait on another REQ is `depends_on`, a question for the user is `pending-answers`, an outside person's confirmation while work proceeds is `Answerer:`. The earmark is missing from that list.
- The release path exists: `skills/do-work/actions/clarify.md:13` and Step 5.5 confirm human-confirmable blocked conditions, and `skills/do-work/actions/work.md:97` keeps a `blocked` REQ out of the scan.
- `skills/do-work/docs/work-guide.md:113` introduces the field as "To say 'leave this one for me'", which is the ambiguous phrasing.
## Detailed Requirements
1. `capture.md:107`, Earmark assessment: add one sentence stating the distinction. Content to carry: `assigned_to` routes work between sessions; when the user must be present as operator, capture `status: blocked` with `blocked_by` naming the person (the user's words, Frontmatter Quoting) and `blocked_at`, which lands in Needs Input · Blocked and is released by `do-work clarify`. Point to the External-condition assessment rather than restating its mechanics.
2. `capture.md:108`, External-condition assessment: add the earmark as a fourth look-alike in the existing "Keep this distinct from the look-alikes" list: work reserved for another session is `assigned_to` (not blocked).
3. `docs/work-guide.md:113-119`, "Earmarking with `assigned_to`": one sentence saying the field is for another session or checkout, and that work waiting on you as operator is captured `blocked` with `blocked_by` naming you, so it shows under Needs Input · Blocked.
## Constraints
- Prose only. No new field, no new status, no Go change, no SKILL.md routing row.
- Keep the never-invent rule: no earmark in the user's words means no `assigned_to`; no person named means no `blocked_by`.
- Read `_dev/primes/prime-action-files.md` before editing; sweep any other restatement of the earmark rule the grep finds (family alternate-writer-contract-drift). Today's grep found only the three sites above.
## Dependencies
None. REQ-643 (board shows an earmarked pending REQ under Pending → Earmarked) is related; this sentence is true before and after it ships, so there is no `depends_on`.
## Builder Guidance
Certainty is high: the maintainer specified the sentence's content. Latitude is wording only. Keep each addition to one sentence per site.
## Red-Green Proof
**RED prompt/case:** A capture request reading "I need to be at the keyboard for this one, hold it for me". Follow `capture.md` Step 1 as written today.
**Why RED now:** The Earmark assessment is the only rule that matches, so capture seeds `assigned_to`; the REQ lands under Pending and nothing in the pipeline releases it.
**GREEN when:** `capture.md:107` says `assigned_to` is for routing between sessions and operator-present work is captured `blocked` with `blocked_by` naming the person and `blocked_at`; `capture.md:108` lists the earmark as a look-alike; `work-guide.md` says the same in one sentence. Checked by a prose read and `_dev/tests/shipped-package-reference-contract.sh`.
**Validation:** User confirmed. The maintainer wrote the distinction in the verbatim input.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 6429 tokens, `slugged: partial` so no targeted form is legal; over the 2000-token budget. Matching reason: its owning prime governs the action files this REQ edits; family `alternate-writer-contract-drift` covers restated rules across `capture.md` and `work-guide.md`.
## Full Context
See `do-work/user-requests/UR-140/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read prime-action-files.md, general/coding-guardrails headings/shared-principles/communication-style crew, and the `alternate-writer-contract-drift` lesson bullets. Grep the restated earmark rule across `skills/` first. Add one sentence per site at capture.md Earmark assessment, capture.md External-condition look-alike list (replacing the closed count "three" with a condition-keyed phrase), and work-guide.md Earmarking paragraph (also fixing its ambiguous "leave this one for me" opener). No heading or bold label changes; no board Pending placement mentioned.
- [x] **[APPLY]:** Three edits as planned in the two write_set files, committed as `b789ca74`. Bold labels **Earmark assessment**, **External-condition assessment**, **Earmarking with `assigned_to`.** unchanged. Never-invent rule untouched (the existing "Never invent or infer a session name" and the user's-words requirement for `blocked_by` both stand).
- [x] **[UNIFY]:** `git diff HEAD~1 --stat`: capture.md (2 ins, 2 del), work-guide.md (2 ins, 2 del); total 2 files, 4 insertions, 4 deletions. `git diff --check` clean. Checked: capture.md lines 107-108 (sentence placement, cross-reference spelling `actions/clarify.md`, Frontmatter Quoting still governs `blocked_by` via the External-condition bullet it points to); work-guide.md lines 113-119 (opener no longer reads as operator-present; new sentence true before and after REQ-643). No debug artifacts. `bash _dev/tests/shipped-package-reference-contract.sh` PASS (1s wall); `bash _dev/tests/contract-regressions.sh` PASS (21s wall). *(P-A-U text from the builder hand-back.)*
*Source: maintainer feedback "Board: a pending REQ earmarked for the user reads as 'Ready', which misleads", change 2, accepted by validate-feedback.*
---

## Triage

**Route: A** - Simple

**Reasoning:** Prose-only, two named files with named lines, content of each sentence specified by the maintainer; effort-mechanical with wording-only latitude.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/capture.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)

**What was done:** The capture Earmark assessment now says `assigned_to` routes work between sessions and that operator-present work is captured `blocked` (per the External-condition assessment) with `blocked_by` naming the person in the user's words and `blocked_at`, so `do-work clarify` releases it; the External-condition look-alike list is keyed on its condition ("every look-alike that is not blocked") and adds the earmark. The work guide's Earmarking paragraph opens with "leave this one for that session" and closes with one sentence saying the same split. Merged from builder branch `worktree-agent-REQ-644-capture-earmark-vs-operator-blocked` (commits `b789ca74` and review fix `54005faf`), range `cd91f838..3602e584` (first merge `58168fe7`, re-merged `3602e584` after review finding F1).

## Decisions

*(from the builder hand-back)*

- D-01 DECIDE & STATE: replaced "the three look-alikes" with "every look-alike that is not blocked" rather than "four". Value: condition-keyed, cannot go stale when another look-alike is added. Risk: none; the list still enumerates each.
- D-02 DECIDE & STATE: changed the work-guide opener "leave this one for me" to "leave this one for that session". Value: removes the operator-present reading the REQ names as the root ambiguity. Risk: low; the bold label is unchanged, so citation pins are unaffected (reference contract passes).
- D-03 DECIDE & STATE: the capture sentence says "the user's words" for `blocked_by` and points at the External-condition assessment for quoting and timestamp mechanics, instead of restating Frontmatter Quoting. Value: one home for the mechanics. Risk: none.
- D-04 DECIDE & STATE (integrator): accepted review finding F1 and changed "Needs Input · Blocked" to "Needs input · Blocked" in both files (builder-branch commit `54005faf`, re-merged as `3602e584`), because the board heading (`template.html`), `work-guide.md:102` and `board-guide.md:18` spell it with a lowercase "input" and the guide would otherwise spell it two ways. The REQ's own requirement text keeps the capitalised form it was written with.

## Discovered Tasks

*(from the builder hand-back)*

- impact-negligible: capture.md:143 (addendum path) only handles an addendum that names a session; an addendum saying the user must be present as operator would need the pending REQ flipped to `blocked` with `blocked_by`/`blocked_at`, and the addendum rule does not say so. Not observed in practice → report only

## Qualification

**Gate records:** `advance --diff-range cd91f838..58168fe7` (first merge) returned the `qualify` gate `satisfied` on provenance `merged_range` with no findings (no debug artifact, P-A-U boxes ticked, `[UNIFY]` present). After the F1 re-merge, `advance` had moved past qualification, so `DO_WORK_DIFF_RANGE=cd91f838..3602e584 bash skills/do-work/tools/checks/qualify.sh` was run directly: `OK: mechanical qualification passed`.

**Requirement trace against the merged diff (`git diff cd91f838..3602e584 --stat`: 2 files, 4 insertions, 4 deletions):**
1. `capture.md:107` Earmark assessment: one appended sentence says `assigned_to` routes work between sessions, and operator-present work ("I need to be at the keyboard for this one") is captured `blocked` per the External-condition assessment, with `blocked_by` naming the person in the user's words and `blocked_at`, landing in Needs input · Blocked for `do-work clarify` (`actions/clarify.md`) to release. It points to the External-condition assessment for quoting and stamp mechanics instead of restating them. Met.
2. `capture.md:108` External-condition assessment: the earmark is added as a look-alike ("work reserved for another session or checkout is an earmark → `assigned_to`"). The closed count "the three look-alikes" became "every look-alike that is not blocked" (D-01). Met.
3. `work-guide.md:113-119`: the opener no longer says "leave this one for me" (D-02), and one closing sentence says the field is for another session or checkout while operator-wait work is captured `blocked` with `blocked_by` naming you, shown under Needs input · Blocked until `do-work clarify` releases it. Met.

**Constraints:** prose only, no field, status, Go, or routing change. The never-invent rules stand: "Never invent or infer a session name" is untouched, and `blocked_by` still comes from the user's words. No bold label or heading changed. Neither file states where the board places an earmarked pending REQ, so the text stays true before and after REQ-643 (the sibling REQ that moves earmarks to Pending → Earmarked).

**Scope:** Route A, `write_set` = the two touched files exactly; no file outside it changed.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` (repository gate) at merge `58168fe7` (121 s wall, exit 0) and again at the re-merge `3602e584` after review finding F1 (124 s wall, exit 0); focused probe `do-work/runs/work-2026-10-07-232637/REQ-644-probe.sh` (`_dev/tests/shipped-package-reference-contract.sh` + `_dev/tests/contract-regressions.sh`) run by `advance` at `58168fe7` and directly at `3602e584`, exit 0 in about 21 s each.
**Result:** ✓ All passing — `advance` recorded `green-gate` satisfied and the probe green; the builder also ran both contract tests green in its worktree (1 s and 21 s).

**Red-green validation:** non-behavioral prose change, so the evidence is the captured `## Red-Green Proof` read as a cold capture agent:
- RED (before, at `cd91f838`): "I need to be at the keyboard for this one, hold it for me" matched only the Earmark assessment, so capture would seed `assigned_to`.
- GREEN (after, at `58168fe7`, unchanged at `3602e584`): the Earmark assessment names that exact case and routes it to `blocked` via the External-condition assessment with `blocked_by` (the person, in the user's words) and `blocked_at`; "leave this for cloud-alpha" still names a session and still seeds `assigned_to`.

**New tests added:** none (prose only; the citation contract pins the unchanged bold labels).

**Heavy verification plan:**
- Range: cd91f838641a11ef4af46cc5c81f945c3c0bd8af..3602e584d16e85822a34839df0acb803070c5f02 (re-planned after the F1 re-merge; same lane as at 58168fe7)
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/capture.md and skills/do-work/docs/work-guide.md matched subtree skills

## Review

**Overall: 97%** | 2026-10-07T23:43:54Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 94% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `capture.md:107` and `work-guide.md:119` wrote "Needs Input · Blocked" while the board heading and other docs use "Needs input · Blocked" — impact-negligible → fixed on the builder branch (`54005faf`), re-merged as `3602e584`, delta re-review Pass. F2 `capture.md:107` "`blocked_by` naming that person" may lead an agent to substitute a name for "me"; prefer "the operator condition in the user's words" — impact-rule-change → report only. F3 (nit) "lands it in Needs input · Blocked" is true only without an unmet `depends_on` (otherwise Pending → Waiting) — impact-negligible → report only
**Acceptance:** Pass — the cold read routes the keyboard case to `blocked` with `blocked_by`/`blocked_at`, and "leave this for cloud-alpha" still seeds `assigned_to`
**Restatement sweep:** redefined the earmark-vs-blocked boundary and the External-condition look-alike list. All restatements (work-reference.md:114,511, capture-reference.md:41, capture.md:143) still agree
**Suggested testing:** 0 items
**Follow-ups created:** None (F2, F3 report only; F1 fixed)

*Reviewed by review-work action (first review 95% at 58168fe7; delta re-review 97% at 3602e584; full report `do-work/runs/work-2026-10-07-232637/REQ-644-review.md`)*

## Lessons Learned

**What worked:** grepping the restated earmark rule across `skills/` before editing; the builder found the closed count "the three look-alikes" and the guide's ambiguous "leave this one for me" example, and fixed both at their source instead of only adding the requested sentence.
**What didn't:** the REQ quoted the board column as "Needs Input · Blocked", and both new sentences copied that capitalisation, while the board heading and the rest of the guide say "Needs input · Blocked". Review caught it and it cost one re-merge and re-gate.
**Worth knowing:** when prose names a UI label, copy it from the UI source (`skills/do-work-board/tools/queue-kanban/web/template.html`), not from the REQ text. Review F2 (`blocked_by` "naming that person" when the user only says "me") and the addendum path at `capture.md:143` are open report-only wording gaps.

## Orientation

Capture now separates two kinds of "hold this": `assigned_to` reserves work for another session or checkout, and work that needs the user at the keyboard is captured `blocked` and released through `do-work clarify`; lives in the capture action (`actions/capture.md` Step 1 assessments) and the user guide, under `_dev/primes/prime-action-files.md`. Prime spot-check: its referenced paths still exist; not stale.

## Heavy Verification Plan

- Base: cd91f838641a11ef4af46cc5c81f945c3c0bd8af
- Target: 3602e584d16e85822a34839df0acb803070c5f02
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/capture.md matched subtree skills; skills/do-work/docs/work-guide.md matched subtree skills

## Heavy Verification Result

- Target revision: 3602e584d16e85822a34839df0acb803070c5f02
- Execution revision: 3602e584d16e85822a34839df0acb803070c5f02 (detached drain checkout)
- staged-skills: exit 0, executed, 35 s wall (an earlier run at the first merge 58168fe7 also passed, executed, 38 s)

## Timing

Observed 2026-10-07T23:30:13Z to 2026-10-07T23:45:03Z: 14m 50s total, 16m 11s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 6m 02s | 4 |
| review | 5m 12s | 1 |
| builder-work | 4m 44s | 1 |
| handback-merge | 13s | 1 |

Slowest stage: review / independent review + delta re-review, 5m 12s, outcome success.
