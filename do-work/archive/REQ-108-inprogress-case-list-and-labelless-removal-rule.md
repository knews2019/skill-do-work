---
id: REQ-108
title: "Review fix: In-Progress Record still enumerates two recovery cases and owes no removal rule for a label-less entry"
status: completed
completed_at: 2026-08-05T11:44:27Z
claimed_at: 2026-08-05T11:38:54Z
created_at: 2026-08-05T11:36:39Z
kb_status: pending
user_request: UR-018
addendum_to: REQ-104
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
review_generated: true
write_set: [actions/work-reference.md, actions/forensics.md, decisions/log.md, decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md]
related: [REQ-094, REQ-095, REQ-104]
batch: parallel-building
---

# Review Fix: In-Progress Record's Case List and the Label-Less Removal Rule

## What

REQ-104 dropped the label-less authorship heuristic in `actions/work-reference.md`'s Crash Recovery
ladder (a label-less checkpoint entry is now always report-only), but two consequences of that drop
were not carried into the rest of the contract:

1. `actions/work-reference.md:456` — In-Progress Record (Step 2)'s opening paragraph restates the
   classification as a **closed two-item** set: "one that is not — unnamed, or named under another
   checkout's label — is a foreign claim". Under the drop the correct non-own set is three (unnamed,
   foreign-label, label-less). The leading positive clause ("under this checkout's own writer label …
   recovers") still yields correct behavior, so nothing misbehaves today; the risk is a future editor
   re-deriving the case list from here. This is the Closed-Enumerations failure shape the skill's own
   conventions name.
2. The rewritten bullet routes reclaim of a genuinely-own pre-0.170.0 entry to "the takeover ladder
   below, or `actions/forensics.md` Check 1's manual reset" — but neither path has a rule for removing
   the label-less `## In Progress (interrupted)` entry. In-Progress Record's removal rule is scoped to
   "this checkout's **own** entry"; `actions/forensics.md:39` says to drop the "**own-label**" entry and
   leave "any entry under **another checkout's** `writer:` label untouched". A label-less entry is
   neither. Before the drop it was classified as own before recovery ran, so the own-entry rule reached
   it; it no longer does.

   Consequence: a reclaimed label-less REQ leaves a permanent checkpoint entry, and `actions/work.md`
   Step 10's session-start note forbids deleting `CHECKPOINT.md` while any "no label at all" entry
   remains — a phantom claim re-reported every session with no documented exit.

## Context

Found during review of REQ-104 (Important findings 1 and 2). Two Minor items travel with it:
`decisions/log.md:106` still records the edge as "filed rather than fixed in flight" while ADR-018's
Consequences now says resolved; and ADR-018's frontmatter `updated:` was not bumped for the
2026-08-05 body edits (`adr-001` sets the precedent that it is maintained on amendment).

## Detailed Requirements

- Widen `actions/work-reference.md:456`'s restatement to name all three non-own cases, or restate the
  condition rather than the list (per the Closed-Enumerations rule, prefer stating the condition).
- State who removes a label-less `## In Progress (interrupted)` entry when a human reclaims the REQ —
  in In-Progress Record (Step 2)'s removal rule and in `actions/forensics.md` Check 1's remediation.
  Both currently name only the two labeled cases.
- Update `decisions/log.md:106` to note the edge was closed by REQ-104, and bump ADR-018's `updated:`.
- Do not add liveness machinery; UR-018's ban is unchanged.
- Consider whether the reworded In-Progress Record paragraph needs a suite pin, or whether REQ-104's
  existing pair suffices.

---
*Source: REQ-104 review findings (Important 1 + 2, Minor 3 + 4)*

---

## Triage

**Route: A** - Simple

**Reasoning:** The review that spawned this REQ already did the discovery — every site is named with a line number and each requirement prescribes its fix (widen one enumeration, add a label-less clause to two removal rules, two bookkeeping edits). Four files, prose-level changes, no exploration needed.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

**D-01 — The reworded In-Progress Record paragraph gets no new suite pin; REQ-104's existing pair suffices.**
The paragraph is a *restatement* whose canonical home is Crash Recovery (Step 1), and the fix here is to stop
restating the case list at all — so after this REQ there is no second enumeration left to drift. Pinning the
restatement would freeze one phrasing of a list defined elsewhere: the Closed-Enumerations failure this REQ
exists to remove, re-encoded in the test suite, where it would also fight any future rewording that is
behaviorally identical. The behavior itself is already guarded at its real home by the pair REQ-104 added —
the positive pin on `claim of unknown origin, always report-only` and the negative pin forbidding
`locally modified or otherwise uncommitted` — and both would fail if the label-less classification regressed.
Per `crew-members/maintenance.md` §2, a stale enumeration is fixed by generalizing to the condition, not by
adding machinery to keep the copy in sync. Same call for the label-less removal rule: `actions/forensics.md`
Check 1 states it subordinate to its pointer at In-Progress Record (Step 2) rather than as an independent
copy, so there is one owner and nothing to hold in lock-step.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]** — Four edits, all prose, all inside `write_set`:
      1. `actions/work-reference.md` In-Progress Record (Step 2) opening paragraph — replace the closed
         two-item non-own list with the condition ("anything that does not match this checkout's own writer
         label"), pointing at Crash Recovery as the single enumeration. → verify: suite-pinned phrases in the
         sed-extracted block survive (`record is a list`, `never grow into one` + `never refreshed` co-located,
         `writer: <hostname>:<absolute-checkout-path>`); block headers unrenamed.
      2. Same file, same section, removal-rule bullet — add one sentence giving the label-less entry an owner
         (the human who reclaims the REQ). → verify: no liveness language; foreign-label rule unchanged.
      3. `actions/forensics.md` Check 1 remediation — same rule, one clause, subordinate to its existing
         pointer at In-Progress Record (Step 2).
      4. `decisions/log.md` L106 — append the closure clause; ADR-018 frontmatter `updated:` → 2026-08-05 and
         drop the duplicate REQ-104 mention at L79.
      → verify: `bash _dev/tests/contract-regressions.sh` exits 0; `git diff --stat` shows exactly the four
      write_set files plus this REQ.
- [x] **[APPLY]** — All four edits landed; no others. `actions/work-reference.md` L456 now states the
      condition and defers the case list to Crash Recovery; its removal-rule bullet and
      `actions/forensics.md` Check 1 both give the label-less entry an owner (the reclaiming human);
      `decisions/log.md` L106 records the closure; ADR-018 `updated: 2026-08-05` plus the L79 duplicate-mention
      nit. No liveness language added; no shipped file cites the maintainer doc.
- [x] **[UNIFY]** — `git diff --stat`: `actions/forensics.md` (+1/-1), `actions/work-reference.md` (+2/-2),
      `decisions/log.md` (+1/-1), `decisions/records/adr-018-…md` (+2/-2) — exactly the `write_set`, plus this
      REQ file and the pipeline's own `do-work/CHECKPOINT.md` / queue→working move (Step 2 machinery, not this
      REQ's writes). No debug artifacts, no stray untracked files.
      `bash _dev/tests/contract-regressions.sh` → **exit 0**. Pins re-verified by hand inside the
      sed-extracted In-Progress Record block: `record is a list` ✓, `never grow into one` ✓,
      `never refreshed` ✓ and co-located on the same line as the tripwire ✓,
      `writer: <hostname>:<absolute-checkout-path>` ✓; both block headers unrenamed ✓; Crash Recovery's
      `claim of unknown origin, always report-only` present ✓; `locally modified or otherwise uncommitted`
      still absent ✓.
      Files checked: `actions/work-reference.md`, `actions/forensics.md`, `decisions/log.md`,
      `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md`,
      `actions/work.md` (read-only, for the Step 10 delete gate and the restatement noted below).

## Discovered Tasks

- **Minor — `actions/work.md:655` carries a third restatement of the recovery case list, in the shape REQ-108
  just removed from `actions/work-reference.md:456`.** Step 10's session-start note says a named entry "that
  isn't (unlabeled, or labeled for another checkout) is a foreign claim recovery must not strip." The set is
  complete and the behavior is right, but it calls the label-less case a *foreign claim*, whereas since
  REQ-104 the canonical term is a *claim of unknown origin* — and L656 two lines later uses the correct term.
  Same fix shape: state the own-label condition and defer the enumeration to Crash Recovery. Out of REQ-108's
  `write_set`, so not touched.

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — In-Progress Record (Step 2): the opening paragraph now states the own-label condition and defers the case enumeration to Crash Recovery ("any other entry is left byte-identical, however it fails to match that label"); the removal-rule bullet gained a sentence making the label-less entry leave with the REQ when a human reclaims it.
- `actions/forensics.md` (modified) — Check 1 remediation: the manual reset now also sends a label-less entry for the reclaimed REQ with it (the reclaim decision is the authorship call recovery declined to make); foreign-label entries stay untouched.
- `decisions/log.md` (modified) — ADR-018 summary line: closure clause appended (edge closed by REQ-104; label-less entries are now always report-only).
- `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md` (modified) — frontmatter `updated:` bumped to 2026-08-05; duplicate REQ-104 mention at L79 reduced to "that REQ".

**What was done:** Carried REQ-104's drop through the two contract sites that still assumed the old classification (one stale closed enumeration replaced by its condition, one ownerless removal case given an owner) plus the two decision-record bookkeeping items. D-01: no new suite pin — the fix removes the second enumeration rather than freezing it; REQ-104's pin pair guards the behavior at its canonical home. Suite exits 0.

## Qualification

Passed — 4 files verified in the diff (each a minimal, targeted prose edit matching its requirement), 5 requirements traced (enumeration→condition, removal rule at both sites, log closure, ADR `updated:` bump + nit, D-01 pin decision recorded), P-A-U confirmed against the diff — no debug artifacts, no undeclared touches. Orchestrator re-ran the suite independently: exit 0.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh` (builder post-edit + orchestrator independently)
**Result:** ✓ All passing — `Contract regression checks passed.`, exit 0

Red-green validation omitted — non-behavioral change (prose alignment of restatements with REQ-104's already-pinned rule; the behavior's own red-green lives in REQ-104). Regression evidence: all pinned phrases in both sed-extracted blocks re-verified present, including the co-location check (`never grow into one` + `never refreshed`), REQ-104's positive pin present and negative fingerprint absent.

*Verified by work action*

## Review

**Overall: 95%** | 2026-08-05T11:44:27Z | Route A quick scan (in-session, per pipeline-mode calibration)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor, 1 nit

- All five requirements delivered word-by-word: the In-Progress Record paragraph now states the own-label condition and explicitly defers the case enumeration to Crash Recovery (the Closed-Enumerations fix shape); the label-less removal rule exists at both prescribed sites and agrees (forensics states it subordinate to In-Progress Record's rule — one owner, no lock-step copy); `decisions/log.md` closure appended without rewriting the line; ADR-018 `updated:` bumped and the double-mention nit fixed; no liveness language (diff grepped clean).
- D-01 (no new suite pin) is the correct call for the right reason: pinning a restatement would re-encode the exact failure shape this REQ removes; REQ-104's pair guards the behavior at its canonical home.
- **Nit:** the removal-rule sentence's "phantom that `actions/work.md` Step 10's delete gate re-reports every session with no exit" describes the *pre-fix* consequence — with this rule in place the phantom now has an exit. The sentence reads as motivation, so it stays accurate as written; noted only for the next editor.

**Acceptance:** Pass — suite exit 0 (orchestrator-run); all pinned phrases in both extracted blocks verified, headers unrenamed, REQ-104's pins intact.
**Follow-ups created:** none from review. The builder's discovered task (work.md:655 terminology drift, [low]) is handled by Step 8's consent flow → REQ-109 (pending-answers).

*Reviewed by review-work action (pipeline mode, quick scan)*

## Lessons Learned

**What worked:** Stating the condition and deleting the copy (rather than widening the list) resolved the enumeration drift permanently — there is no second case list left to go stale.
**Worth knowing:** When a classification case loses its auto-path, sweep every *lifecycle* rule scoped to the old classes (removal, cleanup, delete gates) — REQ-104 fixed the classifier and this REQ had to fix the two removal rules it orphaned. The pair is one change conceptually; try to land them together next time.

## Orientation

Closes the loop on REQ-104's drop: the crash-recovery contract's satellite rules (claim-entry removal, forensics manual reset, decision records) now all agree that a label-less checkpoint entry is unattributable — and that a human reclaiming one takes its entry with the REQ, so no phantom claim can outlive its reclaim. Lives in the work pipeline's crash-recovery contract. No system-shape change.

*kb/ handoff: deferred (unattended run — no KB write without consent); kb_status: pending*
