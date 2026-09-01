---
source_type: req_lesson
req_id: REQ-099
req_path: do-work/archive/UR-018/REQ-099-automatic-wave-dispatch.md
date: 2026-08-04
domain: general
module: actions
tags: [general, automatic, wave, dispatch, work]
---

# Lessons from REQ-099: Automatic wave dispatch — the work loop computes and dispatches the ready set

## What the REQ was about

Give the work pipeline a fan-out mode where **the loop computes the wave itself** and dispatches builders without a confirmation gate. This is a deliberate contract change: today `actions/work.md:33` says the action "does not drive a fan-out wave" and `actions/work-reference.md:320` says "a human picks which REQs run together — nothing computes the set." Both sentences get rewritten, per the user's explicit choice of fully automatic set-picking over a human-confirmed set.

## Solution summary

Inverted the "nothing computes the set" contract into an opt-in auto-wave mode with a computed, bounded ready set and no confirmation gate, while leaving the serial floor path, the per-REQ dispatch flow, serial integration, and every `write_set` display-only pin exactly as they were.

## What worked

- Reading the suite's `fan_out_block` extraction *before* writing. The window runs from `**Fan-Out Dispatch` to `^## Composed Exit Summary`, so new prose had to land inside it — and one of the five pinned phrases sits in the exact bullet being inverted. Discovering that after writing would have meant rewriting the rewrite.
- Keeping `write_set`'s exclusion argued from *absence reads as unknown* rather than from a builder count. That is the form the two suite sweeps are built to catch, and it is also the only version of the argument that stays true under a computed set.

## What didn't work

- The first draft of the Auto-wave predicate restated dependency-readiness in its own words. That is a second definition of readiness — the exact drift the Closed-Enumerations rule warns about, and it would have diverged from Step 1's the first time either changed. Replaced with "the same predicate and the same cycle detection as the serial scan; auto-wave adds no second definition of readiness."
- Writing multi-paragraph prose through a shell heredoc broke twice on the `→` character and on an unterminated quote inside the here-document. Prose edits of this size belong in a script file invoked by path, not inlined into a shell command.

## Worth knowing

- A new `do-work run` flag needs **three** edits in `actions/work.md` (Input list, strip list, usage string) plus **two** in `SKILL.md` (routing row, dispatch table). The strip list is the one that bites: omit it and the unrecognized-argument guard rejects the flag the same file just documented.
- `SKILL.md` is word-budgeted at 2650 by the suite. Router additions must be phrases, not sentences.
- The fan-out section's five pinned phrases are extracted as a *block* between two markers. Anything inserted after `## Composed Exit Summary` is outside the window and will not satisfy them, however correct it reads.

## Back-reference

See `do-work/archive/UR-018/REQ-099-automatic-wave-dispatch.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `0cf9420`.
