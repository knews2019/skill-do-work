---
title: "Lessons from REQ-100: Live auto-wave acceptance run — prove real wall-clock concurrency"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-100-live-auto-wave-acceptance-run-prove-real.md]
related:
  - page: REQ-094-checkpoint-writer-label-crash-recovery-i
    rel: complements
  - page: REQ-095-two-clone-acceptance-run-checkpoint-pois
    rel: complements
  - page: REQ-099-automatic-wave-dispatch-the-work-loop-co
    rel: depends-on
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-100: Live auto-wave acceptance run — prove real wall-clock concurrency

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

Run the REQ-099 automatic wave dispatch live, with genuinely concurrent builders, and record the evidence. Real wall-clock fan-out concurrency has **never been proven** in this skill — the one recorded attempt (REQ-085) logged Partial — so this run is the proof, not ceremony.

## Solution summary

Ran REQ-099's auto-wave dispatch live. The ready-set predicate was applied to a four-REQ fixture shaped to exercise two of its exclusion clauses; the resulting bounded wave of two was dispatched into two git worktrees whose builders ran **concurrently for 4.109 seconds of measured wall-clock overlap**; both hand-backs were integrated serially with per-REQ merge ranges; a second wave was recomputed and correctly picked up the REQ whose dependency had just landed; and a deliberate overlapping-`write_set` pair confirmed the merge refuses, which is the safety argument the computed set rests on.

## What worked

- Shaping the fixture around **exclusions** rather than inclusions. Two of the four REQs existed only to be left out, which is what turns "the wave contained the right REQs" from a coincidence into a result. An include-only fixture would have passed identically against a computation that ignored `depends_on` and `assigned_to` entirely.
- Computing overlap as `min(end) − max(start)` and printing a verdict, instead of eyeballing two timestamps. It makes the deliverable a number that is either positive or the run failed, and it is the difference between this result and REQ-085's Partial.
- Running the negative case at all. A computed set that is *documented* as not proving non-overlap is worth much less than one where the merge has been watched refusing.

## What didn't work

- The dispatch script captured **one `<pre>` for the whole wave**, at claim time. It reads naturally — the wave shares a starting tip — and it is wrong the moment the first REQ's release tail commits, which is exactly what serial integration does between merges. The contract already forbids it; I wrote the forbidden version anyway on the first pass, which is the best available evidence that the rule needs to keep saying so.
- Assuming a `cd` into a fixture directory would survive its own deletion: removing the fixture while the shell sat inside it left the next command unable to resolve a working directory (`fatal: Unable to read current working directory`) and made a clean teardown look like a failure. Delete from outside.

## Worth knowing

- `date -u +%Y-%m-%dT%H:%M:%S.%NZ` gives nanosecond stamps that Python parses after truncating to microseconds (`stamp[:26]`). Second-resolution timestamps cannot prove a sub-second overlap, and a 4-second sleep is the cheapest way to make the window unambiguous.
- Under fan-out the owner's bookkeeping commit covers **several** claims at once and still lands below every REQ's `<pre>` — because each `<pre>` is captured later, at that REQ's own merge. The two rules fit together only in that order.
- `git worktree remove` on a worktree whose branch still has unmerged commits needs `--force`; in the negative case that is correct (the branch was deliberately abandoned), and it is the one place in this run where `--force` was right. On the happy path both plain `remove` and `branch -d` succeeded, which is the assertion that matters.

## Back-reference

See `do-work/archive/UR-018/REQ-100-live-wave-acceptance-run.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `7ab69e3`.
