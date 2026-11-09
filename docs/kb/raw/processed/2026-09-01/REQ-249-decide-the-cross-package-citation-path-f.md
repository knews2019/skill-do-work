---
source_type: req_lesson
req_id: REQ-249
req_path: do-work/archive/UR-055/REQ-249-decide-the-cross-package-citation-path-form.md
date: 2026-08-18
domain: general
module: _dev/primes
tags: [general, decide, cross, package, citation]
---

# Lessons from REQ-249: Decide the cross-package citation path form and sweep to match

## What the REQ was about

Two incompatible readings of the cross-package citation path coexist in shipped markdown, and nothing can tell them apart. Pick one and sweep.

## Solution summary

Every backticked cross-package citation in shipped markdown now resolves as a real relative path from the citing file's own directory, in both source and installed topologies — 122 citations rewritten across 46 files, 19 already-correct ones untouched, depth derived per file (`../../` from `actions/`/`crew-members/`/`docs/`, `../../../` from `tools/queue-kanban/`). `_dev/primes/prime-action-files.md` § Cross-Referencing states the literal form as the rule, retires the skill-root-relative shorthand by name, exempts fenced template/example blocks, and points at the checker. The mechanical-checkability question is answered **yes and implemented**: `_dev/tests/shipped-package-reference-contract.sh` now verifies backticked cross-package citations in both topologies, reusing the existing CommonMark walk via a behavior-identical extraction refactor (`mask_block_code` / `inline_code_regions`), locked by the existing parser fixtures; against the pre-sweep tree it reports exactly the 122 swept citations.

## What worked

Deriving the extent mechanically (condition over the corpus) instead of trusting the capture counts — the brief's grep said 140, the condition said 122+19, and the review's independent scanner reconciled to the same 141 exactly. Making the sweep's condition BE the checker's condition means the two cannot disagree by construction. Refactoring the existing CommonMark walk instead of writing a second parser kept the fixtures as the identity proof.

## What didn't work

The class boundary was drawn at backtick spans — the letter of the decided rule — and the retired *reading* survived at three bare-text sites, including the core SKILL.md stating the old resolution rule as prose (REQ-259). Seventh consecutive occurrence of the instance-vs-class shape; the class this time was "the reading", not "the spelling".

## Worth knowing

`do-work` names both the core package and the consumer queue root, so the checker deliberately skips core-package spans whose tail names nothing real — a citation to a *deleted core* file is invisible to it (documented, fixture-pinned; the other three packages still catch deletions). Anchor fragments on backticked citations are discarded, not validated. Fenced blocks are exempt by design.

## Back-reference

See `do-work/archive/UR-055/REQ-249-decide-the-cross-package-citation-path-form.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `cc1083c`.
