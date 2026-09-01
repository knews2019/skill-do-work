---
source_type: req_lesson
req_id: REQ-118
req_path: do-work/archive/UR-024/REQ-118-normalize-flag-vocabularyless-field.md
date: 2026-08-06
domain: general
module: tools/queue-kanban
tags: [general, normalize, flag, must, calling]
---

# Lessons from REQ-118: The normalize flag must stop calling vocabulary-less field values unrecognized

## What the REQ was about

`queue-kanban frontmatter get <file> <field> --normalize` emits a `⚠ … not recognized — no canonical vocabulary is defined for this field.` line on **every** call when the field has no Schema Read Contract row — a timestamp, a title, a path list. Gate the warning on the field actually having a row, so the flag becomes a clean no-op for fields the contract itself places outside its scope.

## Solution summary

Added `hasSchemaFieldContract(fieldName)` — a lookup of the contract table, the predicate that was missing — and gated `runFrontmatterCommand`'s normalize/resolve branch on it. A field with no row now prints its value and returns, observably identical to the same `get` without `--normalize`. `--in-set` on such a field returns a usage error (exit 2) naming the field, rather than a silent exit 1 that would read as a real negative. `schemaFieldWarningText`'s contract-less branch is reworded to say the field is outside the contract and read verbatim, instead of calling the value unrecognized (see D-01).

## What worked

**What worked:** Treating the noise as a symptom and asking what the code could not express. The warning was not a wording bug — `isKnownSchemaFieldValue` was being asked a question it structurally could not answer, returning one `false` for "this value is wrong" and "this field isn't mine". Once that was named, the fix was a missing predicate rather than a condition bolted onto the call site, and the same predicate is what any future caller needs.

**What didn't:** The instinct to delete the now-unreachable warning branch. Following it would have left `schemaFieldWarningText` falling through to `expected one of []` for an ungated caller — a worse message than the one being removed. Zero-value structs make "unreachable, so delete it" more dangerous in Go than it looks: the fallthrough path still runs, it just runs on empty data.

**Worth knowing:** The gate splits `--normalize` and `--in-set` deliberately. Silence is right for `--normalize` (nothing to normalize against, so no-op) and wrong for `--in-set`, because both set names are `status` sets — answering "not a member" for a timestamp would look like a real negative at a call site written as `if …; then`. One narrower looseness is left alone on purpose: `--in-set` on a field that *has* a row but is not `status` (e.g. `domain`) still exits 1 rather than erroring. Pre-existing, out of this REQ's scope, and harmless today since the only prose call site passes `status`.

## Back-reference

See `do-work/archive/UR-024/REQ-118-normalize-flag-vocabularyless-field.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `8d1a9f2`.
