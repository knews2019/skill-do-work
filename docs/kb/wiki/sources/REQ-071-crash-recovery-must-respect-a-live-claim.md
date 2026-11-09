---
title: "Lessons from REQ-071: Crash recovery must respect a live claim before stripping and re-queueing"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-071-crash-recovery-must-respect-a-live-claim.md]
related:
  - page: REQ-072-go-utility-allocates-req-ids-and-version
    rel: complements
  - page: REQ-073-fan-out-dispatch-n-concurrent-builders-u
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-071: Crash recovery must respect a live claim before stripping and re-queueing

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

Crash recovery is currently unconditional and destructive. `actions/work-reference.md:220-223`
tells the pipeline that every `REQ-*.md` in `do-work/working/` is its own interrupted leftover, so it
resets the frontmatter, **strips thirteen generated sections**, and moves the file back to
`do-work/queue/`. Make it respect a live claim instead: recover only what this session can show is
its own, and ask a human before taking over anything else.

## Solution summary

Put a classification gate in front of `## Crash Recovery (Step 1)`'s existing substeps instead of rewriting them. A `working/` REQ named in the checkpoint's `## In Progress (interrupted)` record is an **own crash** and recovers exactly as before; anything else — including the common case of no checkpoint at all — is a **foreign claim** that is left byte-identical and reported. A foreign claim is offered for takeover only past a three-hour age, or immediately when `claimed_at` is unparseable, future-dated (2-minute skew allowance), or absent; the offer is a `clear-questions.md`-shaped two-option prompt, and with no human to answer the outcome is "leave it, report it, continue." The threshold's non-authorizing rationale is stated where the threshold is, so a later edit cannot collapse it into an automatic takeover. `claimed_at` is read during classification, before substep 1 discards it. The premise that licensed unconditional recovery ("no other live session whose in-flight claim a recovery could disturb") is gone from both files. `actions/work.md` Step 1 now reads `do-work/CHECKPOINT.md` first and says why — it is recovery's input — with the Step 10 session-start note, the Orchestrator Checklist, and the Verification Checklist reconciled to match. Ten new contract-suite assertions pin the removed premise as absent, each replacement rule as present, and the checkpoint-read-before-recovery ordering.

## What worked

- **Gating the destructive procedure instead of rewriting it.** Substeps 1–3 are unchanged byte-for-byte; only their precondition moved. That kept the diff on the dangerous path at zero, made the review one question ("when do these run?") instead of two, and means anything that already trusted the substeps still can.
- **Writing the assertions before the prose, on a prose-only REQ.** Committing to the exact phrases first (`absent checkpoint is ambiguous`, `never authorizes`, `substep 1 removes it`) forced each requirement to become one findable sentence rather than an idea diffused across a paragraph — the assertion is what stops requirement 4's rationale from being "simplified" away later, which is the whole point of that requirement.

## What didn't work

- **The first instinct — "recover it if this session claimed it" — has nothing to read.** The skill keeps no session identity anywhere (that was the point of REQ-069), so ownership is unknowable in principle. The checkpoint is not a session id; it is the only durable record that *some* session was mid-REQ. Accepting that the discriminator is weaker than ownership is what made the design work: it can only be trusted in one direction, so the no-match side must be non-destructive.
- **Claiming the manual reset was "defined once" in forensics.** It isn't — substeps 1–3 define the automatic version right below the sentence making the claim. Caught in the restatement sweep and reworded to a plain pointer. Asserting single-sourcing while sitting next to a second copy is worse than not claiming it.

## Worth knowing

- **A hard crash usually leaves no checkpoint at all** — Step 10 writes it at session end. So the common real-world case is the *foreign claim* branch, not the own-crash branch, and the practical effect of this REQ is that a crashed REQ now waits for a human instead of being auto-recovered. That is the intended trade (the REQ's `## Why` values the uncommitted trail above the convenience), but anyone measuring "how often does recovery fire" will see it drop to near zero, and that is not a bug.
- **A stale checkpoint is now load-bearing.** Deleting `do-work/CHECKPOINT.md` as tidy-up turns every subsequent own-crash recovery into a hands-off report. Step 10's substep 3 already defers deletion until recovery finishes, which is what keeps that safe — do not "simplify" that deferral either.
- **The exclusive-session invariant now reads slightly oddly next to this.** `actions/work-reference.md:53` still says the pipeline supports "one active REQ, one coder context," while recovery now spends prose on a claim it cannot account for. Not a contradiction — the invariant is a *product contract*, and this is what the pipeline does when reality violates it — but REQ-073 is rewording that invariant, and it is worth checking there that the two paragraphs read as one position. Flagged here rather than edited, per this REQ's Constraints.

## Back-reference

See `do-work/archive/UR-013/REQ-071-crash-recovery-respects-a-live-claim.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5c39899`.
