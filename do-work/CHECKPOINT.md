---
session_ended: 2026-08-19T16:40:00Z
last_completed: REQ-290
queue_state: 28 pending, 1 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 2
session_depth: moderate
---

# Session Checkpoint

## In Progress (interrupted)

## Completed This Session

`do-work run UR-060` — both members shipped, each with an independent adversarial review:

- REQ-289: separate impact from effort, with unique greppable tokens on both axes (Route C, 88%) — `2ea7be5`, **0.214.0**
- REQ-290: surface impact in REQ titles and add a run filter that skips negligible work (Route B, 70% → remediated) — `225e287`, **0.215.0**

Both hashes confirmed with `record-commit-hash.sh --verify`. `maintainer-verify.sh` exited 0 at every commit boundary. Serial mode throughout — no worktrees created, so none to clean up.

**UR-060 stays open** in `user-requests/`: it gained five follow-ups, all still queued.

## Still Queued

**Twenty-nine** — 28 `pending`, 1 `pending-answers`. The queue grew by three this session (two shipped, five created), and every one of the six came from a review or a builder's stop-and-report finding something real.

**Needs you (`pending-answers`):**
- **REQ-296** — should the retired `trivial`/`normal` vocabulary left in internal names (the board's projection keys, the estimator's `--trivial` flag) follow the rename? Both were deliberate non-goals with stated reasons; the REQ asks whether that stays the answer.

**New this session (`pending`):** REQ-293 (sweep, 6 instances — the impact/effort lock-in checks pin a spelling rather than the property), REQ-294 (capture's impact guard is one-directional), REQ-295 (bare "impact" wordings in the toolbox audit), REQ-297 (targeted mode under-reports what the new flag skipped).

## Session Notes

**The reviews earned their cost twice, and one of the catches was the orchestrator's.**

- I mutation-tested REQ-289's four lock-in checks and reported them sound. The reviewer showed my mutation for the main check used the word "stamping" — the single literal that check greps. Re-run with six realistic re-drift phrasings, it catches one. A mutation written in the check author's own vocabulary tests the vocabulary, not the property. That is REQ-293's F2.
- REQ-290's builder justified its always-double-quote rule by claiming an unquoted tagged title leaves a REQ's status, UR pointer and dependencies riding on strict parsing. The parser's own comment says the opposite — recovery exists so that does not happen — and a live parse confirmed nothing is lost. I corrected the reason; the reviewer then found the *true* reason is stronger than my correction: a title opening `[` and closing `]` is parsed as a YAML flow list and comes back altered, commas silently eaten. Verified, and it is what the convention now says.

**Three right rules, three wrong reasons, in two REQs.** Every one was caught by someone other than its author, and only by checking arguments rather than verdicts. Nothing in the suite tests a justification.

**A REQ that adds a condition to a list must sweep every gloss of that list.** REQ-290 added a fifth auto-wave ready-set condition and updated the canonical list correctly; three restatements of "the four conditions" survived inside the two files it was already editing, one of them thirteen lines below the condition it contradicted. All fixed in place rather than deferred.

**Estimator calibration held.** REQ-289 estimated 60 active minutes against a 66-minute wall span; REQ-290 estimated 15 against 40 — the second inflated by serial review latency, not by work. Two rows appended to `do-work/calibration-log.tsv`; do not recalibrate off a sample of two.

## Context Summary

**The impact/effort split is complete and shipped, but its guard rails are not.** The field, the vocabulary, the parser, the board chip, the title tag, and the run filter all exist and were each proven live rather than asserted. What is missing is the layer that keeps them true: REQ-293 now carries six instances of checks that pin a spelling instead of a property, including the one property REQ-290's filter depends on — that `impact:` defaults to `impact-user-visible` and never to `impact-negligible`. Flip that default and the filter silently skips every REQ predating the field, with a green suite.

**REQ-294 is the one that decides whether any of this pays off.** Three forces — the capture template, the contract default, and a one-directional "never invent impact-negligible" guard — all push a new REQ toward `impact-user-visible`. If every REQ lands on the default, the field is decorative and the filter has nothing to filter. It carries an open question.
