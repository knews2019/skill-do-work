---
source_type: req_lesson
req_id: REQ-174
req_path: do-work/archive/UR-039/REQ-174-validate-root-markdown-fence-info.md
date: 2026-08-11
domain: testing
module: skills/do-work/tools
tags: [testing, validate, root, markdown, fence]
---

# Lessons from REQ-174: Validate root Markdown fence info

## What the REQ was about

Make root and list fence classification share the CommonMark rule that a backtick-fence info string cannot itself contain a backtick.

## Solution summary

Centralized marker-aware fence info-string validation across root and list openings plus paragraph/container state, and added root/list/tilde differential fixtures that preserve the existing classifier contracts.

## What worked

- Markdown fence recognition is a compound contract: marker kind, info-string validity, and paragraph/container state must change together or a locally correct opener check can still mask rendered content.
- Differential fixtures against the pinned renderer catch classifier drift more reliably than asserting regex structure.

## Back-reference

See `do-work/archive/UR-039/REQ-174-validate-root-markdown-fence-info.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `bd5ecf6`.
