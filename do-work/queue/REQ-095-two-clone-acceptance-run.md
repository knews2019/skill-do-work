---
id: REQ-095
title: Two-clone acceptance run — checkpoint poisoning repro and claim-conflict evidence
status: pending
created_at: 2026-08-04T19:44:17Z
user_request: UR-018
domain: testing
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-094]
maintenance: false
related: [REQ-094, REQ-096]
batch: parallel-building
---

# Two-Clone Acceptance Run — Poisoning Repro and Claim-Conflict Evidence

## What

Prove the cross-checkout model with a real two-clone experiment, the way REQ-085's fan-out acceptance run proved worktree dispatch (and found the index-settling bug). Two parts: (1) reproduce the checkpoint-poisoning failure against the **pre-REQ-094** instructions and confirm the writer label stops it; (2) claim the same REQ in two clones, merge, and capture the actual conflict text git produces.

## Detailed Requirements

- Set up two throwaway clones of this repo (scratch space, not inside the repo tree).
- **Poisoning repro:** in clone A claim a dummy REQ (checkpoint entry written, committed); sync to clone B; run clone B's crash-recovery reading per the old rule to show the strip would fire, then per the REQ-094 rule to show the foreign entry is reported and left alone. Record both transcripts.
- **Claim conflict:** both clones claim the same queued dummy REQ (move to `working/`, edit frontmatter, commit); merge one into the other; capture the real conflict git reports (expected: same-path content or rename conflict). Document the observed fix-at-merge resolution.
- Correct any failure-mode claims in the shipped prose **from this evidence, not from reasoning** — if the observed conflict shape differs from what REQ-096's widened dispatch prose predicts, fix the prose.
- Record the run like REQ-085's fan-out run (find where that artifact lives and mirror its form); clean up the throwaway clones afterward.

## Constraints

- The experiment must not touch this repo's own `do-work/` state — dummy REQs live only in the throwaway clones.

## Red-Green Proof

**RED prompt/case:** Under pre-REQ-094 rules, clone B's recovery strips clone A's live claim (deterministic, first sync).
**Why RED now:** No writer identity on checkpoint entries; the model has never been exercised across two clones (fan-out concurrency itself was only ever recorded as Partial, REQ-085).
**GREEN when:** With the writer label, the same sequence leaves clone A's claim intact and reports it; the double-claim merge conflict is captured verbatim and the documented behavior matches the evidence.
**Validation:** User confirmed (approved plan, Phase 1 item 2).

## Full Context

See `do-work/user-requests/UR-018/input.md` and `assets/approved-plan.md` (Phase 1).

---
*Source: approved plan, Phase 1 item 2*
