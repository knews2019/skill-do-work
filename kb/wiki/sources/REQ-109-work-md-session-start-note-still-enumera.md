---
title: "Lessons from REQ-109: work.md session-start note still enumerates the recovery case list and calls a label-less entry a foreign claim"
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-109-work-md-session-start-note-still-enumera.md]
related:
  - page: concept-session-checkpoints-and-recovery
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-109: work.md session-start note still enumerates the recovery case list and calls a label-less entry a foreign claim

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

Discovered during REQ-108 (`[low]`): `actions/work.md`'s Step 10 session-start note (the sentence
listing which `working/` REQs recovery may strip) carries the same closed-enumeration shape REQ-108
just removed from `actions/work-reference.md`'s In-Progress Record — "one that isn't (unlabeled, or
labeled for another checkout) is a foreign claim recovery must not strip." The set is complete and
the behavior is correct, but it calls the label-less case a *foreign claim*, whereas since REQ-104
the canonical term is a *claim of unknown origin* — and the very next sentence in the same note uses
the correct term. Same fix shape as REQ-108: state the condition, defer the list to Crash Recovery.

## Solution summary

Rewrote the first sentence of item 2 in Step 10's session-start note: it now states the own-label condition (own `writer:` label → this session's own to recover; any other entry left byte-identical) and defers the case enumeration to `actions/work-reference.md` → **Crash Recovery (Step 1)**, dropping both the hand-enumerated case list `(unlabeled, or labeled for another checkout)` and the pre-REQ-104 "foreign claim" label for the label-less case. The second sentence (entries written at claim time by Step 2) is unchanged. Builder also verified the rest of the note — items 1, 3, line 646, line 658, and the checklist line 668 — uses the canonical terms correctly; no other drift found.

## What worked

**What worked:** Reusing REQ-108's exact fix shape (state the condition, defer the enumeration, reuse the canonical "byte-identical" term) made the two passages read as one contract instead of two paraphrases.
**What didn't:** N/A — no dead ends on a one-sentence fix.
**Worth knowing:** The review's restatement sweep found a fourth instance of this same drift class in the same file (`actions/work.md:774`, Verification Checklist — "a reported foreign claim" is under-inclusive for the label-less/unknown-origin case). The pattern: REQ-104 changed the classification vocabulary, and each sweep since has caught one more stale restatement (REQ-108 → work-reference.md, REQ-109 → work.md line 655, now line 774). When a vocabulary changes, grep for the old term across *every* shipped file in the first fix, not one file per follow-up.

## Back-reference

See `do-work/archive/UR-018/REQ-109-workmd-recovery-caselist-terminology.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `5f50fb7`.
