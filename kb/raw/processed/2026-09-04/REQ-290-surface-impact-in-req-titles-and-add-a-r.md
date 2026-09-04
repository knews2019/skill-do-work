---
source_type: req_lesson
req_id: REQ-290
req_path: do-work/archive/UR-060/REQ-290-surface-impact-in-req-titles-and-a-run-filter.md
date: 2026-08-19
domain: general
module: _dev/primes
tags: [general, surface, impact, titles]
---

# Lessons from REQ-290: Surface impact in REQ titles and add a run filter that skips negligible work

## What the REQ was about

Make the `impact:` field REQ-289 introduces actually usable for the decision the user wants to make:
put the token in the REQ title so it is searchable today, and give `do-work run` a flag that skips
negligible-impact work.

## Solution summary

**Files changed:**
- `skills/do-work/actions/capture-reference.md` (modified) — new `## REQ Title Convention` section as the canonical home; the Simple REQ and Addendum templates carry the quoting and tag comment.
- `skills/do-work/actions/capture.md` (modified) — Step 1's impact assessment mirrors a non-default verdict into the title; a Verification Checklist line covers the tag and its absence from filenames.
- `skills/do-work/actions/review-work.md` (modified) — Step 10's follow-up template title carries the tag ahead of `Review fix: `; the prose above states the emit-only-when-non-default rule.
- `skills/do-work/actions/work-reference.md` (modified) — the `impact:` schema line gains the source-of-truth-vs-mirror rule and names the new filter; the Schema Read Contract row lists the filter as a reader; the auto-wave ready set gains condition 5 with its counts corrected; the Composed Exit Summary gains a fourth headline and a seventh section; the Builder-Decided Follow-up template title is tagged.
- `skills/do-work/actions/work.md` (modified) — `--skip-impact-negligible` added to `## Input`, the argument-strip list, both usage-string branches, the queue status summary suffix, a new Step 1 skip paragraph after the `assigned_to` skip, the auto-wave filter list, the exit-path headline and section order, and Step 0 of the Orchestrator Checklist.
- `skills/do-work/docs/capture-guide.md` (modified) — user-facing mirror: quoted title, `impact:` added to the schema block, one paragraph on the tag and the flag.
- `skills/do-work/docs/work-guide.md` (modified) — the flag's user-facing documentation, placed beside the existing `--wave` / `--fan-out` prose (D-08, orchestrator).
- `skills/do-work-toolbox/actions/code-review.md` (modified) — a one-line pointer to the convention; its titles stay untagged because the action writes no `impact:` (D-05).

## What worked

- Directing the builder to a specific structural site (`work-reference.md`'s auto-wave conditions,
- which close with "Nothing else enters the computation") rather than to the feature in the
- abstract. That block is where the flag would otherwise have silently no-opped, and naming it in
- the dispatch is why condition 5 exists at all.
- Dogfooding the vocabulary immediately. Writing this UR's own follow-ups with `impact:` tokens and
- then searching the live board for `impact-negligible` is what turned acceptance criterion 3 from
- an assertion into an observation.

## What didn't work

- **Exploration dismissed the Discovered Tasks flow as emitting "no `title:` at all".** True of its
- template and false of its behavior — the flow mints REQs and is explicitly required to stamp
- `impact:` on every one. Reading a template for what it *contains* rather than what the surrounding
- prose *instructs* is how an emitter goes missing from an emitter list. The fix re-keys the list on
- the condition ("any flow that mints a REQ carrying an `impact:` value") so the next one inherits it.
- **A REQ that adds a condition to a list must sweep every gloss of that list, not just the list.**
- Three restatements of the ready-set conditions survived inside the two files this REQ was already
- editing, one of them thirteen lines from the condition it contradicted. The canonical list was
- updated correctly; the prose *about* it was not.

## Worth knowing

- A YAML title that opens with `[` and closes with `]` is parsed as a flow list by the board's
- lenient recovery path and comes back **altered** — commas inside it are eaten as separators, with
- no warning. Quoting is what takes the class off the table. This was found only because a reviewer
- checked a justification that was already "close enough", which is the argument for checking
- reasoning rather than verdicts.
- `--skip-impact-negligible` is deliberately conservative: absent and unrecognized both resolve to
- `impact-user-visible`, so an unjudged REQ is never dropped. That property is what makes the flag
- safe to add to a queue whose REQs mostly predate the field, and it is the one thing no test pins
- (tracked as REQ-293's F4).

## Back-reference

See `do-work/archive/UR-060/REQ-290-surface-impact-in-req-titles-and-a-run-filter.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `225e287`.
