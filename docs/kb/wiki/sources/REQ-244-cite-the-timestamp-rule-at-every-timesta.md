---
title: "Lessons from REQ-244: Cite the Timestamp rule at every timestamp write site"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-04/REQ-244-cite-the-timestamp-rule-at-every-timesta.md]
related: []
created: 2026-09-04
updated: 2026-09-04
confidence: medium
---

# Lessons from REQ-244: Cite the Timestamp rule at every timestamp write site

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

Sweep all four skills for every timestamp write site — templates and action steps carrying `[timestamp]`, `<timestamp>`, `<now>`, `[UTC timestamp]`, or any `*_at:`/date-shaped placeholder — normalize each to the spellings the Timestamp rule recognizes (`<timestamp>` / `<now>`), and add an inline citation of the rule (`Timestamp rule, actions/work-reference.md`) at each site that lacks one.

## Solution summary

**Files changed:**
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/clarify.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/forensics.md` (modified)
- `skills/do-work/actions/roadmap.md` (modified)
- `skills/do-work-toolbox/actions/code-review.md` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `skills/do-work-toolbox/actions/deep-explore.md` (modified)
- `skills/do-work-toolbox/actions/deep-explore-reference.md` (modified)
- `skills/do-work-knowledge/actions/interview.md` (modified)
- `skills/do-work-knowledge/actions/interview-reference.md` (modified)

## What worked

- **A rule that centralizes a command creates an obligation at every site that does not have it.** REQ-078 moved the clock command to one home for a real reason (Windows agents kept getting an unreachable fix). That move is only safe if every other site points home — and nothing enforced the second half for eight months, until a stamp was fabricated. When a future REQ centralizes something, the citation obligation at the periphery is part of the same change, not a follow-up.
- **The `<timestamp>` token was doing two jobs.** It marked stamps *and* appeared inside directory names, which is precisely why a mechanical check looked like it needed an exception. It did not need one — it needed the two jobs separated. Reach for "which of these uses is actually the odd one out" before writing a carve-out into a detector.
- **Backticked cross-package pointers are unenforced.** REQ-243's checker resolves Markdown link syntax; the repo's dominant convention for citing another action file is a backticked path, which that checker never sees. Two spellings of the sibling path (`../` and `../../`) coexist today because nothing can tell them apart. If cross-package pointer correctness matters, either the convention moves to Markdown links or a separate resolver has to read backticked paths — a candidate REQ, not something I did here.
- **`git archive HEAD <path> | tar -x -C <dir>` is a clean way to get a pre-change tree for a RED.** It never touches the working tree, so there is no stash to forget to pop and no risk of the vacuous "stashed a clean file" green.

## Back-reference

See `do-work/archive/UR-055/REQ-244-cite-the-timestamp-rule-at-every-stamp-write-site.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f733365`.
