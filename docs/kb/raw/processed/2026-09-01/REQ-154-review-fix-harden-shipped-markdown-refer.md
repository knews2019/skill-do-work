---
source_type: req_lesson
req_id: REQ-154
req_path: do-work/archive/UR-031/REQ-154-harden-shipped-markdown-reference-parsing.md
date: 2026-08-09
domain: general
module: skills/do-work/general
tags: [general, review, harden, shipped, markdown]
---

# Lessons from REQ-154: Review fix: Harden shipped Markdown reference parsing

## What the REQ was about

Make the shipped-package reference guard distinguish published Markdown links from standards-valid code, comments, escapes, and escaped destinations so legitimate documentation cannot falsely block a release.

## Solution summary

The parser now separates non-published-region masking, structural delimiter recognition, and destination normalization while leaving every downstream package-reference policy check unchanged.

## What worked

- Production-helper fixtures exposed the original false positives and protected the unchanged source/install and URL-policy pipeline.
- Offset-preserving masks and exact-target assertions made line-number and over-masking regressions visible.

## What didn't work

- Hardening the masking, extraction, and bare-URL passes independently left them disagreeing at escaped-link and Markdown block-context boundaries; REQ-158 records the remaining same-root cases for consent.

## Worth knowing

- A Markdown release guard needs one rendered-versus-ignored classification shared by every discovery path. Passing the live corpus is not enough; fixtures must prove both false-positive and false-negative directions.

## Back-reference

See `do-work/archive/UR-031/REQ-154-harden-shipped-markdown-reference-parsing.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `22551dc`.
