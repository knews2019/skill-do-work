---
title: "Lessons from REQ-080: The capture template emits a stray instruction line into every REQ it produces"
type: source-summary
topic_cluster: queue-orchestration-and-lifecycle
sources: [raw/processed/2026-09-01/REQ-080-the-capture-template-emits-a-stray-instr.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-080: The capture template emits a stray instruction line into every REQ it produces

Part of the [[concept-queue-task-lifecycle]] cluster.

## What the REQ was about

The Simple REQ template in `actions/capture-reference.md` ends, **inside** the fenced template body,
with the line `Think carefully before answering.` It is not part of any request — it is an
instruction-like artifact that gets copied into every REQ capture produces. 25 archived REQs carry it.

## Solution summary

Deleted `Think carefully before answering.` (and the blank line separating it from `*Source: …*`) from inside the Simple REQ template's fence in `actions/capture-reference.md`, so `*Source: [original verbatim request]*` is now the template's last line. The Complex REQ and Addendum REQ templates in the same file were checked and never carried it — a repo-wide grep across `actions/`, `specs/`, `prompts/`, `interviews/`, `crew-members/`, `docs/`, `tools/` and `SKILL.md` returned exactly one occurrence, the one deleted. `actions/capture.md` restates no part of the template body.

## What worked

- grepping the *shipped* tree rather than trusting the REQ's single named line — it confirmed the Complex and Addendum templates were clean, which requirement 1 asked about and would otherwise have been an assumption.

## What didn't work

- nothing failed. Worth recording that the temptation requirement 2 and the Constraints both pre-empt (clean the 25 archived files, or restructure how templates are fenced) is real and would have turned a two-line diff into a large one.

## Worth knowing

- the durable lesson is not about this line. **A defect logged as a decision inside a generated artifact does not reach the generator.** REQ-012 diagnosed this exactly right, wrote it down precisely, and filed it in a REQ — where it read as closed. The template kept emitting for 25 more captures. When a builder flags "this content looks wrong," the question that follows is *where did this content come from* — the answer is a source file, and the source file is the fix.

## Back-reference

See `do-work/archive/UR-015/REQ-080-capture-template-emits-stray-instruction.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8ee717b`.
