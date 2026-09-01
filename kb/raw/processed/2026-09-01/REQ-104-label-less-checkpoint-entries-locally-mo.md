---
source_type: req_lesson
req_id: REQ-104
req_path: do-work/archive/UR-018/REQ-104-labelless-entry-authorship-heuristic.md
date: 2026-08-05
domain: general
module: actions
tags: [general, label, less, checkpoint, entries]
---

# Lessons from REQ-104: Label-less checkpoint entries — "locally modified" is not evidence of authorship

## What the REQ was about

`actions/work-reference.md` → **Crash Recovery (Step 1)**, the label-less bullet, treats a locally
modified `do-work/CHECKPOINT.md` as evidence that *this* checkout wrote the entries in it:

> **Named there with no `writer:` label at all** (an entry written before the label existed) → **own
> only where `do-work/CHECKPOINT.md` is locally modified or otherwise uncommitted in this checkout**,
> which is evidence this checkout wrote it and has not shared it; recover it as an own crash.

## Solution summary

Dropped the label-less authorship heuristic (D-01: DROP over NARROW). `docs/work-guide.md` and `actions/work.md` L125/L655 were verified already consistent with the drop semantics and left unedited; the grep sweep found exactly one shipped site asserting the retired inference (the bullet itself). Contract-regression suite exits 0 with the new pins in place.

## What worked

**What worked:** Pre-exploration that mapped the suite's sed-range boundaries and every pinned phrase before dispatch — the builder rewrote a heavily-pinned bullet without breaking a single existing assertion. Red-green demonstrated with the suite's own extraction idiom (run the sed + grep against HEAD to show both new pins would have failed pre-edit) made the proof mechanical instead of rhetorical.

**What didn't:** The builder's mirror sweep grepped for the *retired inference* ("locally modified ⇒ own") and found exactly one site — but a sweep for *restatements of the classification itself* would have caught In-Progress Record's two-case enumeration in the very file being edited. Sweeping for the deleted phrase is not the same as sweeping for the rule.

**Worth knowing:** Dropping a classification case can orphan its downstream lifecycle rules — the label-less case lost its authorship heuristic, which silently disconnected it from every own-entry removal rule (checkpoint entries now have no documented exit for that case; REQ-108). When deleting a case from a ladder, walk what *used to happen after* that case classified, not just the classification.

## Back-reference

See `do-work/archive/UR-018/REQ-104-labelless-entry-authorship-heuristic.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f2177b1`.
