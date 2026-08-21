---
id: UR-059
title: Fix the unfiled contradictions in clarify's Step 4
created_at: 2026-08-19T14:33:51Z
requests: [REQ-288]
word_count: 61
---

# Fix the Unfiled Contradictions in clarify's Step 4

## Summary

The user asked to capture fixes for the contradictions surfaced by a `do-work validate-feedback`
session on REQ-276 and the contradiction sweep that followed it. Six were found. Two were already
filed (REQ-276, REQ-270). Three live in `skills/do-work/actions/clarify.md` Step 4 and none was
filed — this UR captures those three. The sixth (Principle 8 has no user-override clause) was left
deferred by the user's answer and is not captured.

## Extracted Requests

| REQ | Covers |
|---|---|
| REQ-288 | K2 (durable-record definition), K3 (per-question verbs setting file state), K4 (confirm archives unbuilt work) |

## Batch Constraints

- All three defects write one file. They must not be split across parallel builders.
- The dated note K2 introduces must cite the Timestamp rule, never copy a command from it —
  `skills/do-work/actions/work-reference.md:101` makes that paragraph the only place in `actions/`
  that spells a clock command, and `_dev/tests/` fails the build if a copy reappears.

## Full Verbatim Input

do-work capture-request fix the contradictions, use the ask tool to help me choose how to fix it when common sense is not enough to answer it automatically

### Answers given via the ask tool during capture

**On K3 — how should status be decided on a multi-question REQ?**
> Aggregate in Step 5 (Recommended)

**On K4 — how should the confirm branch decide a REQ is "nothing to build"?**
> Both

**On K2 — how should the reasoning note be shaped?**
> if this is a question about date, then note that we had a REQ dealing with this alone and the idea was to always use the same tool to collect the current date, becuase otherwise these dates are halucinated and become non-parsebale

### Standing instruction from the same session

> I don't like contradictions, tell me what kind of contradictions we have

---
*Captured: 2026-08-19T14:33:51Z*
