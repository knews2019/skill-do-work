---
title: "Lessons from REQ-095: Two-clone acceptance run — checkpoint poisoning repro and claim-conflict evidence"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-095-two-clone-acceptance-run-checkpoint-pois.md]
related:
  - page: REQ-094-checkpoint-writer-label-crash-recovery-i
    rel: depends-on
  - page: REQ-096-execution-model-re-grain-claim-anywhere-
    rel: complements
  - page: REQ-100-live-auto-wave-acceptance-run-prove-real
    rel: complements
  - page: REQ-104-label-less-checkpoint-entries-locally-mo
    rel: complements
  - page: REQ-108-review-fix-in-progress-record-still-enum
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-095: Two-clone acceptance run — checkpoint poisoning repro and claim-conflict evidence

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

Prove the cross-checkout model with a real two-clone experiment, the way REQ-085's fan-out acceptance run proved worktree dispatch (and found the index-settling bug). Two parts: (1) reproduce the checkpoint-poisoning failure against the **pre-REQ-094** instructions and confirm the writer label stops it; (2) claim the same REQ in two clones, merge, and capture the actual conflict text git produces.

## Solution summary

Ran the two-clone acceptance test for the first time. Six fixtures, each a bare
origin plus two clones with their own throwaway `do-work/` tree, exercised: the pre-REQ-094 poisoning
(reproduced deterministically), the shipped writer-label rule against byte-identical input, the
double-claim merge, the byte-identical double-claim merge, two disjoint concurrent claims, and the
label-less entry under a merge resolution. Two defects were found and filed rather than fixed.

## What worked

- A **bare origin plus two real clones** rather than a copied directory. Every sync in the record is a
  genuine `git fetch`/`git merge`, which is why the exact conflict strings (`CONFLICT (add/add)`, `AA`,
  `CONFLICT (modify/delete)`, `Fast-forward … merge exit status: 0`) are quotable at all. A file-copy
  fixture would have produced a plausible narrative and no evidence.
- Replaying the **identical input state** (same sha256) against the old rule and the shipped rule, in
  two separate fixtures rather than by resetting one. It makes the A/B airtight and the transcripts
  independently re-runnable.
- Pushing past the two requirements into the adversarial variants (byte-identical claims, disjoint
  claims, a label-less entry after a merge resolution). Both defects came from those three runs; the two
  required runs alone would have produced a clean bill.

## What didn't work

- The REQ's own predicted conflict shape — "same-path content or **rename** conflict" — is wrong, and
  reasoning would have kept it. Git resolves the `queue/` → `working/` rename silently *because both
  sides perform the identical rename*, so only the content inside ever conflicts. Prose predicting a
  rename conflict would have shipped a claim no run can reproduce.
- The first framing of the poisoning as "any routine sync strips the claim" is also wrong. Git protects
  a claiming checkout that has local edits (refuses the merge) or committed divergent ones (raises
  `modify/delete`). The silent strip needs the claimant to be *idle between steps* — which is both
  narrower than assumed and the state every paused or crashed run is actually in.
- `tools/checks/qualify.sh` and `tools/checks/scope-drift.sh` both refuse a no-code-change REQ
  (`FAIL: no '## Implementation Summary' file list`, `SKIP: … exit 2`). REQ-085 hit the identical wall
  a batch earlier. Accepted both times with the reasoning written down; still a script gap, not a REQ
  defect.

## Worth knowing

- **Every conflict result here is a committed-`do-work/` result.** On the common install, where
  `do-work/` is untracked, nothing syncs between checkouts: the poisoning cannot happen and neither can
  any fix-at-merge detection. Do not generalize these transcripts to that install.
- **`do-work/CHECKPOINT.md` is a guaranteed conflict point per concurrent claim, disjoint claims
  included** — two single-line appends land at the same position. That is load-bearing twice over: it is
  the only detector of a byte-identical double claim (F-05), and it is what dirties the checkpoint and
  breaks the label-less authorship heuristic (F-07).
- On one machine every clone reports the same `hostname -s`, so the **path half of the writer label is
  the sole discriminator** locally. Any future test of the hostname half needs two machines.
- The eight suite FAILs in this container are `chmod 500` injections defeated by running as **root** —
  environment, not code. Compare against `do-work/working/baseline-failures.txt` name-for-name rather
  than expecting green.

## Back-reference

See `do-work/archive/UR-018/REQ-095-two-clone-acceptance-run.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0526e44`.
