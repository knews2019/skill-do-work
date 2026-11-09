---
title: "Lessons from REQ-308: Judge effort_estimate on every capture, as impact already is"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-308-judge-effort-estimate-on-every-capture-a.md]
related: []
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-308: Judge effort_estimate on every capture, as impact already is

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

`skills/do-work/actions/capture.md` requires a judged `impact:` on every REQ it mints, in a rule
written to close exactly this hole: "**Judge it yourself and write a value** — an absent `impact:`
must not be the common case." The neighbouring field gets no such treatment.
`actions/capture-reference.md` says only that capture **MAY** set `effort_estimate`, and the
schema line in `actions/work-reference.md` repeats it.

Done means capture judges both fields by the same standard, with the same escape hatches: judge it,
or put the judgment to the user, or leave it absent because neither was possible — never a copied
default.

## Solution summary

Capture's impact rule was applied one field over, unchanged in shape. The
size question is stated plainly — would a competent implementer finish this in one focused pass over
a small, already-identified set of files — and both directions of the never-write-an-unjudged-value
rule are named, so `effort-substantive` by default is called out as the same failure as an invented
`effort-mechanical`.

## Worth knowing

- **Two rules that should be the same rule can be pinned by comparing them, not by quoting them.**
  Asserting a phrase is present is what REQ-293 ruled against; asserting that two sentences are
  identical apart from the field they name enforces "by the same standard" literally, and survives
  any rewording of either.
- **Symmetry alone is satisfiable by weakening both sides.** M3 did exactly that and stayed
  symmetric. A symmetry check needs a floor beside it — here, that the shared sentence still offers
  all three alternatives.
- **The sweep half of a check is what finds the site the REQ did not know about.** The REQ listed
  four files; the grep found a fifth, in a Go comment whose own closing sentence required it to
  change in the same commit. Enumerating sites in the REQ is a starting point, never the set.

## Back-reference

See `do-work/archive/UR-064/REQ-308-judge-effort-estimate-on-every-capture.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `9bce005`.
