---
source_type: req_lesson
req_id: REQ-250
req_path: do-work/archive/UR-042/REQ-250-close-the-remaining-markdown-link-checker-gaps.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, close, remaining, markdown, link]
---

# Lessons from REQ-250: Close the remaining markdown link checker gaps

## What the REQ was about

REQ-243's anchor checker has four known limits. Three are false negatives — links it accepts that are broken — and one is a path-escape the same change escalated from a `stat` to a `read_text`.

## Solution summary

All four checker edges resolved — two closed with code, two documented as stated limitations. (1) Bare `#anchor` links now validate: the fragment resolves as the carrying file's own name through the existing target-and-anchor pipeline; the corpus's 3 real bare-anchor links now get checked. (2) The `..`-escape class is clamped at every normpath-then-probe site — and the hole-hunt found the genuinely silent instance in REQ-249's citation checker, whose consumer-queue probe absorbed interior-`..` tails and could stat outside the repo (demonstrated reaching `/etc/hostname`); both citation sites are clamped and the silent-absorb hole is fixture-locked. (3) HTML-tag/entity slug divergence documented with both failure directions named, pinned by a fixture. (4) Blockquoted-heading drop documented as always-loud, pinned by a fixture.

## What worked

The hole-hunt discipline finally beat the class-vs-instance curse — greping the same primitive (`normpath`-then-probe) across the file before calling the class closed found the genuinely silent hole in a *different* checker's probe, and the review's independent enumeration then found nothing left. Mutation-falsifiable pins: each documented limitation carries a fixture that FAILS if the behavior changes, so statement and behavior cannot drift apart silently.

## What didn't work

Two records-precision slips, both in the orchestrator's transcription rather than the code: D-02's clamp-uniformity claim dropped the builder's "main-loop" qualifier, and the recorded pushback missed the re-entrant-escape case where the pre-change corpus read path was silently unsafe on the standard consumer topology. The builder's own hand-back was the more accurate record both times — transcribe, don't paraphrase.

## Worth knowing

The escape clamp is `".." in <normalized>.parts` — in repo-relative PurePosixPaths that IS the escape condition. Symlink escapes are a different, deliberately unclamped class (no shipped symlinks exist). Closing the HTML-slug divergence needs code-span-aware entity decoding; the pinning fixture predicts the exact GitHub slugs if anyone attempts it.

## Back-reference

See `do-work/archive/UR-042/REQ-250-close-the-remaining-markdown-link-checker-gaps.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `330797b`.
