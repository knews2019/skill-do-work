---
id: REQ-104
title: Label-less checkpoint entries — "locally modified" is not evidence of authorship
status: pending
created_at: 2026-08-04T21:15:00Z
user_request: UR-018
addendum_to: REQ-094
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set: [actions/work-reference.md]
related: [REQ-094, REQ-095, REQ-096, REQ-103]
batch: parallel-building
---

# Label-Less Checkpoint Entries — "Locally Modified" Is Not Evidence of Authorship

## What

`actions/work-reference.md` → **Crash Recovery (Step 1)**, the label-less bullet, treats a locally
modified `do-work/CHECKPOINT.md` as evidence that *this* checkout wrote the entries in it:

> **Named there with no `writer:` label at all** (an entry written before the label existed) → **own
> only where `do-work/CHECKPOINT.md` is locally modified or otherwise uncommitted in this checkout**,
> which is evidence this checkout wrote it and has not shared it; recover it as an own crash.

REQ-095's two-clone acceptance run demonstrated that the premise fails under the claim-anywhere model
this batch is building. Once a second checkout can claim, **every** concurrent claim forces a
`CHECKPOINT.md` merge conflict — F-06 of REQ-095's `## Testing` shows `CONFLICT (add/add)` even for two
fully disjoint claims, because two single-line appends land at the same position. A checkout that
resolved that conflict is holding a modified checkpoint for a reason that has nothing to do with who
wrote which entry. So the heuristic fires on a *foreign* label-less entry, classifies it as an own
crash, and strips a live claim — the 2026-07-01 incident, reachable again through the label-less door.

Evidence: REQ-095 `## Testing` → *Defect found: the label-less bullet is unsound under claim-anywhere*
(run 6 transcript, `R6-3`/`R6-4`).

## Detailed Requirements

- Fix the label-less bullet so a merge-resolution-dirtied checkpoint cannot be read as authorship.
  Two candidate shapes, and the choice is the point of this REQ:
  - **Narrow the heuristic** — require *modified* **and** no merge in progress (`git rev-parse
    --verify -q MERGE_HEAD` fails; note `^`/quoting is not involved here but the same
    re-derive-don't-carry rule applies). Keeps auto-recovery of genuine pre-0.170.0 own crashes.
  - **Drop the heuristic** — a label-less entry is always report-only, never recovered. Strictly safer
    and shorter; costs auto-recovery for checkpoints written before 0.170.0, which a human can still
    reset by hand via `actions/forensics.md` Check 1.
  Recommend dropping it: the population it serves is checkouts that have not run the pipeline since
  0.170.0, which shrinks to nothing, while the failure it enables is unbounded and silent. Prefer the
  narrowing only if the wider recovery is worth carrying a second condition that readers must keep true.
- Whatever the choice, keep the surrounding pinned phrases intact — `absent checkpoint is ambiguous`,
  `foreign claim`, `Crash Recovery's input`, `claim held by`, and the label format string are all
  asserted by `_dev/tests/contract-regressions.sh`. Reword around them.
- Mirror the decision at every site that restates the label-less rule. `actions/work.md` Step 1 and
  Step 10's session-start note both carry a version of it; grep for the condition rather than trusting
  this list (per the Closed-Enumerations rule).
- Add a suite assertion pinning whichever rule lands, so the next re-grain cannot quietly widen it back.

## Constraints

- No liveness machinery. This is a change to *how an entry is attributed*, not a check on whether
  anything is still running — refresh intervals, staleness checks and liveness probes stay banned by
  name (`actions/work-reference.md` → In-Progress Record, the `never grow into one` paragraph).
- `maintenance: true` — the candidate fix is removing or narrowing a rule in the skill's own operating
  instructions, so `crew-members/maintenance.md`'s delete-before-you-add discipline applies.

## Red-Green Proof

**RED prompt/case:** REQ-095 run 6 — a label-less foreign entry plus a merge-resolved (therefore
modified) `CHECKPOINT.md` classifies as an own crash under the shipped bullet.
**GREEN when:** the same input state classifies as report-only (or as own only with no merge in
progress, if the narrowing is chosen), and a suite assertion pins it.
**Validation:** Evidence-backed from REQ-095's acceptance run, not reasoning.

## Full Context

See `do-work/user-requests/UR-018/input.md`, `assets/approved-plan.md` (Phase 1), and REQ-095's
`## Testing` record.

---
*Source: REQ-095 acceptance-run finding F-07 (critical discovered task)*
