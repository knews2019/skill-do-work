---
title: 'Lessons from REQ-104: Label-less checkpoint entries — "locally modified" is not evidence of authorship'
type: source-summary
topic_cluster: checkpoint-and-crash-recovery
sources: [raw/processed/2026-09-01/REQ-104-label-less-checkpoint-entries-locally-mo.md]
related:
  - page: REQ-095-two-clone-acceptance-run-checkpoint-pois
    rel: complements
  - page: REQ-096-execution-model-re-grain-claim-anywhere-
    rel: complements
  - page: REQ-108-review-fix-in-progress-record-still-enum
    rel: complements
  - page: REQ-109-work-md-session-start-note-still-enumera
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-104: Label-less checkpoint entries — "locally modified" is not evidence of authorship

Part of the [[concept-session-checkpoints-and-recovery]] cluster.

## What the REQ was about

`actions/work-reference.md` → **Crash Recovery (Step 1)**, the label-less bullet, treats a locally
modified `do-work/CHECKPOINT.md` as evidence that *this* checkout wrote the entries in it:

> **Named there with no `writer:` label at all** (an entry written before the label existed) → **own
> only where `do-work/CHECKPOINT.md` is locally modified or otherwise uncommitted in this checkout**,
> which is evidence this checkout wrote it and has not shared it; recover it as an own crash.

## Solution summary

Dropped the label-less authorship heuristic (D-01: DROP over NARROW). `docs/work-guide.md` and `actions/work.md` L125/L655 were verified already consistent with the drop semantics and left unedited; the grep sweep found exactly one shipped site asserting the retired inference (the bullet itself). Contract-regression suite exits 0 with the new pins in place.

## What worked

Pre-exploration that mapped the suite's sed-range boundaries and every pinned phrase before dispatch — the builder rewrote a heavily-pinned bullet without breaking a single existing assertion. Red-green demonstrated with the suite's own extraction idiom (run the sed + grep against HEAD to show both new pins would have failed pre-edit) made the proof mechanical instead of rhetorical.

## What didn't work

The builder's mirror sweep grepped for the *retired inference* ("locally modified ⇒ own") and found exactly one site — but a sweep for *restatements of the classification itself* would have caught In-Progress Record's two-case enumeration in the very file being edited. Sweeping for the deleted phrase is not the same as sweeping for the rule.

## Worth knowing

Dropping a classification case can orphan its downstream lifecycle rules — the label-less case lost its authorship heuristic, which silently disconnected it from every own-entry removal rule (checkpoint entries now have no documented exit for that case; REQ-108). When deleting a case from a ladder, walk what *used to happen after* that case classified, not just the classification.

## Back-reference

See `do-work/archive/UR-018/REQ-104-labelless-entry-authorship-heuristic.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f2177b1`.
