---
source_type: req_lesson
req_id: REQ-108
req_path: do-work/archive/UR-018/REQ-108-inprogress-case-list-and-labelless-removal-rule.md
date: 2026-08-05
domain: general
module: actions
tags: [general, review, progress, record, still]
---

# Lessons from REQ-108: Review fix: In-Progress Record still enumerates two recovery cases and owes no removal rule for a label-less entry

## What the REQ was about

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

## Solution summary

Carried REQ-104's drop through the two contract sites that still assumed the old classification (one stale closed enumeration replaced by its condition, one ownerless removal case given an owner) plus the two decision-record bookkeeping items. D-01: no new suite pin — the fix removes the second enumeration rather than freezing it; REQ-104's pin pair guards the behavior at its canonical home. Suite exits 0.

## What worked

Stating the condition and deleting the copy (rather than widening the list) resolved the enumeration drift permanently — there is no second case list left to go stale.

## Worth knowing

When a classification case loses its auto-path, sweep every *lifecycle* rule scoped to the old classes (removal, cleanup, delete gates) — REQ-104 fixed the classifier and this REQ had to fix the two removal rules it orphaned. The pair is one change conceptually; try to land them together next time.

## Back-reference

See `do-work/archive/UR-018/REQ-108-inprogress-case-list-and-labelless-removal-rule.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `53929a2`.
