---
source_type: req_lesson
req_id: REQ-293
req_path: do-work/archive/UR-060/REQ-293-make-the-impact-effort-lock-in-checks-pin-the-property.md
date: 2026-08-21
domain: general
module: _dev/primes
tags: [general, impact, effort, lock, checks]
---

# Lessons from REQ-293: Make the impact/effort lock-in checks pin the property instead of one spelling of it

## What the REQ was about

REQ-289's four lock-in checks all pass, and each one catches the exact defect it was written
against. Every one of them also pins a **spelling** — a literal verb, a markup shape, a partial file
set — rather than the property it claims to hold. Five instances, one root cause.

Done means the class cannot recur: each check holds its property against a re-drift written in
different words, not just against the phrasing that existed when it was authored.

## Solution summary

Rewrote REQ-289's four lock-in checks so each holds its property rather than one spelling of it, added the two checks that were missing entirely, and pinned the flag and title-tag declaration sets that nothing held together.

## What worked

**What worked:** Computing what the *old* check would have done to each mutation, instead of only checking that the new one catches it. That is what turns "the new check is better" from a claim into a measurement, and it showed that one of my three chosen mutations would have been caught anyway — by an unrelated clause — which I would otherwise have counted as evidence it wasn't.

**What didn't:** The first widening of Check A's window went from `[^.]{0,80}` to `.*` and immediately false-positived on four lines of correct prose. The second attempt, `.{0,60}`, still failed on one line — because the action file **states the rule in as many words**: "`effort_estimate` is a different axis and is never derived from that token." A check written to catch a sentence will catch its own contract being stated, and the negation guard is the price. Worth knowing before writing any prose-grep guard: the file most likely to contain your forbidden phrasing is the file that defines why it is forbidden.

**Worth knowing:** REQ-289's mutation test used the word "stamping" — the single verb its own check greps. That is the sharpest form of the self-confirming test: the mutation was drawn from the same imagination as the pattern, so it could only ever pass. The REQ's acceptance criterion ("use a different verb and a different markup shape than the one already in the file") is the general fix, and it is worth applying to every guard, not just this one: **choose the mutation before looking at the pattern, or you will choose the mutation the pattern already catches.**

## Back-reference

See `do-work/archive/UR-060/REQ-293-make-the-impact-effort-lock-in-checks-pin-the-property.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `df976d9`.
