---
id: UR-064
title: Make the mechanical-work selector see the whole queue
created_at: 2026-08-20T22:00:52Z
requests: [REQ-308, REQ-314]
word_count: 118
---

# Make the Mechanical-Work Selector See the Whole Queue

## Summary

A directed capture that fell out of building `do-work run-simple-reqs` (shipped in 0.220.0, PR #152).
The verb selects on `effort_estimate`, and the user asked whether the field is a trustworthy input.
Measured against the live queue during that work: **14 of 22 pending REQs carry `effort_estimate`
at all**, and the remaining 8 resolve to `effort-substantive` by the Schema Read Contract's
documented default. Those 8 are therefore invisible to the new verb — not because anyone judged
them substantive, but because nobody judged them.

The instruction is to close that gap at the source: capture already demands a judged `impact:` on
every REQ, and effort should be judged the same way.

## Original Request

> `effort_estimate` coverage is this feature's real ceiling: 14 of 22 pending REQs carry the field,
> and the other 8 default to `effort-substantive` and stay invisible to the verb. `capture.md`
> already requires a judged `impact:` on every REQ but only says capture *may* set
> `effort_estimate`. Making capture judge effort the same way is the highest-value next step.

Raised in the PR #152 body and confirmed for capture by the user, who chose to file it rather than
fold it into that PR: it changes capture behavior and was outside what #152 was asked to do.

## Extracted Requests

| Request | Disposition |
|---|---|
| Capture judges `effort_estimate` on every REQ, as it already judges `impact:` | REQ-308 |

## Batch Constraints

- **One REQ, deliberately.** The asymmetry between the two fields is a single root cause with a
  single fix site; splitting it into "change capture" plus "backfill the queue" would mint a second
  REQ for work the first one's rollout already decides.
- **Do not widen the enum.** `effort_estimate` stays the closed two-value triage bit
  (`effort-mechanical` | `effort-substantive`). REQ-228 ruled against growing it toward t-shirt
  sizes and that ruling stands; this is about whether the value is *judged*, never about how many
  values exist.
- **Absence must keep working.** Every existing REQ without the field stays valid and keeps reading
  as `effort-substantive`. The fix raises the judgment bar for new captures; it is not a migration
  that rewrites frozen records.
