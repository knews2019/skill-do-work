---
id: UR-060
title: Separate impact from effort so they stop being conflated
created_at: 2026-08-19T14:33:51Z
requests: [REQ-289, REQ-290]
word_count: 178
---

# Separate Impact from Effort So They Stop Being Conflated

## Summary

While capturing UR-059 the user redirected: what they want from a REQ is a statement of its
**impact** — whether anyone would ever notice the work — so they can stop implementing a REQ whose
impact is negligible. Investigation found that the test already exists (`review-work.md:340-341`)
but its verdict is stamped into `effort_estimate`, a field documented as measuring size. One field,
two writers, two axes. The user then asked for a token vocabulary in which every token is unique
and greppable.

## Extracted Requests

| REQ | Covers |
|---|---|
| REQ-289 | The split itself: the `impact:` field, the unique token vocabulary, retiring `gate:` and the discovered-task severity words, and unwiring effort from the impact gate |
| REQ-290 | Making impact actionable: the token in the REQ title, and `do-work run --skip-impact-negligible` |

## Batch Constraints

- REQ-290 depends on REQ-289 — there is nothing to filter on until the field exists.
- Every token must be unique repo-wide and findable by plain-text search. `trivial` currently
  matches 104 lines under `skills/` and `normal` matches 520; both are used on two different axes,
  which is the conflation in its most literal form.
- `work-reference.md:137` requires the board's Go parser and the schema line to change in the same
  commit.
- REQ-228 recorded "No new frontmatter field. Not on REQs, not on URs. `effort_estimate` stays a
  two-value triage bit." That decision was about timeline projection. The new REQ must name it and
  say why it does not apply, or the next reviewer re-litigates it.

## Full Verbatim Input

> why are you looking for effort_estimate? I want to know the impact of a REQ, so for example if the impact is TRIVIAL, I might not want to keep implementing it.

> why are you looking for effort_estimate? I want to know the impact of a REQ, so for example if the impact is TRIVIAL, I might not want to keep implementing it. For example if it's a minor CSS change, or wording, that the client will never see, then it might be a good time to stop

> let's consider how to best disambiguate effort from impact, so they are not conflated

### Answers given via the ask tool during capture

**On the vocabulary:**
> all tokens have to be unique, they can be long thouh, like low-impact-user-visible, etc... now propose new ones properly

**On where impact should surface:**
> Field + title prefix + run filter

**On how a capture-time REQ gets an impact value:**
> Required only above a threshold

### Earlier framing from the same session

> what I wanted is that each REQ is actually estimated if it is trivial or not, and that should be in the title, so if I ever want to stop spawning and processing them I can.

---
*Captured: 2026-08-19T14:33:51Z*
