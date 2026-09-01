---
source_type: req_lesson
req_id: REQ-162
req_path: do-work/archive/UR-031/REQ-162-handle-ordinary-multiline-backtick-commands.md
date: 2026-08-10
domain: general
module: skills/do-work/general
tags: [general, review, handle, ordinary, multiline]
---

# Lessons from REQ-162: Review fix: Handle ordinary multiline backtick commands

## What the REQ was about

Extend reserved-recipe collision scanning to retain physical-line state for ordinary single-backtick Just command literals. Done means the broader multiline-literal scanner accepts reserved-looking command payload without weakening real definition detection, exact diagnostics, or pre-mutation preservation.

## Solution summary

Reserved-looking recipe and alias payload inside valid ordinary multiline backtick commands is now ignored, while same-line commands, inactive contexts, triple-backticks, prior string families, exact real-collision reporting, and pre-mutation preservation remain intact.

## What worked

- Just-parseable positives plus exact stderr and byte snapshots made a one-condition lexer correction independently provable in both acceptance and rejection directions.
- Reusing the existing raw active-delimiter path preserved suffix close/reopen behavior and kept root/shipped helpers byte-identical.
- Replaying the production helper against all six Just multiline delimiters, then parsing the result with Just itself, exposed a real reserved recipe that the line-oriented collision scan had silently missed.

## What didn't work

- The earlier multiline-literal repair treated ordinary backticks as line-local without checking Just's accepted physical-line form, leaving the final byte-96 state boundary uncovered.
- Literal state was completed without carrying recipe-header state with it: the opening line of a multiline parameter default has the recipe name but no colon, while the closing line has the colon and was skipped as literal content. The helper therefore accepted a duplicate reserved recipe in the no-Just path.

## Worth knowing

- Ordinary backticks close on the next literal backtick even after a backslash; they must not inherit cooked double-quote escape parity, and exact triple-backticks must remain the longest-first opener.
- A handwritten lexer must test each literal in every grammar position it can occupy, not only as assignment payload. For Just, variable values prove false-positive suppression; recipe-parameter defaults prove real-definition retention. The fix retains the pending top-level header until the literal closes, then evaluates the complete header while preserving `:=` assignment exclusion.

## Back-reference

See `do-work/archive/UR-031/REQ-162-handle-ordinary-multiline-backtick-commands.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `aff7c9c`.
