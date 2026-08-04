---
source_type: req_lesson
req_id: REQ-079
req_path: do-work/archive/UR-015/REQ-079-guards-pin-the-weaker-fingerprint.md
date: 2026-08-03
domain: testing
module: _dev/tests
tags: [testing, tests, guards, weaker, fingerprint]
---

# Lessons from REQ-079: Two guards pin the weaker fingerprint of the premise they exist to retire

## What the REQ was about

REQ-075 (v0.166.2) established that a retired premise leaves two fingerprints — the thing it *said*
("one REQ at a time") and the thing it was *called* ("under the exclusive-session model") — and that
the second is the more dangerous of the two. Its own regression check pins only the first. Its regex is
also narrower than the class it names. And `actions/cleanup.md:31` still argues the safety of the
skill's one destructive pass from a premise REQ-071 spent an entire REQ falsifying.

## Solution summary

The strong-form pattern is now defined **once** as `builder_count_premise_pattern` and referenced by all three of its consumers, instead of being spelled out three times — F4b's gap was triplicated precisely because the pattern was. It was widened from `one REQ( [a-z]+)? at a time` to cover the class: the count word may be `one`/`a single`/`only one`, up to two adjectives may precede the noun, the noun may be `REQ`/`builder`/`coder`/`agent`, and the tail may be `at a time`/`at once`/`concurrently` or the `is ever building|running|in flight` shape.

## What worked

- **Dry-running the widened pattern against the whole tree, with the filter applied, before adopting it.** Requirement 3 asked for a re-check after widening; doing it as the *first* step instead turned a possible late surprise into a design input — it is what showed that the filter, not the regex, is what protects the canonical sentence.
- **A ruling table with a stated test.** Four sites, two corrected, two left alone, and one sentence saying how they were told apart (cited for what the premise asserts vs. for a consequence it does not support). That is reusable; a list of four verdicts is not.

## What didn't work

- **The obvious way to "be thorough" would have broken the lesson.** `tools/queue-kanban/prime-do-kanban.md:60` quotes both fingerprints verbatim and sits inside a swept directory. Adding a file-level negative there — the same treatment `model.go` and `board.js` get, and the natural next step if you are pattern-matching on "cover every file in `tools/`" — would make the suite fail on the very lesson entry that documents the rule. It survives only because the line filter finds no `write_set` on it. Written into the assertion's comment as a do-not.

## Worth knowing

- **A guard's blast radius is the pattern *and* the filter, and they do different jobs.** The pattern decides what the premise looks like; the filter decides where the premise is wrong. Widening the first without understanding the second is how you either miss the class or flag true statements — and this REQ needed both to stay separate to satisfy requirements 2 and 3 at once.
- **The weak fingerprint is not always wrong**, which is what makes it survive review: two of the four sites cite the exclusive-session model correctly. A mechanical sweep of the phrase would have been a regression. The line filter is again what makes the assertion safe — it only fires where the premise is being used to explain `write_set`.
- **"Suite passes" was the false signal this REQ existed to remove.** Running `HEAD`'s suite against the injected wordings and watching it exit 0 three times is worth more than any argument about regex coverage — and it took about a minute.

## Back-reference

See `do-work/archive/UR-015/REQ-079-guards-pin-the-weaker-fingerprint.md` for the full REQ — triage, implementation, review, and lessons. Commit `8fdce3c`.
