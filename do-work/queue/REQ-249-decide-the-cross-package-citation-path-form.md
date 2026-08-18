---
id: REQ-249
title: Decide the cross-package citation path form and sweep to match
status: pending
created_at: 2026-08-18T13:54:59Z
status_changed_at: 2026-08-18T14:12:05Z
user_request: UR-055
addendum_to: REQ-244
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: true
---

# Decide the Cross-Package Citation Path Form and Sweep to Match

## What

Two incompatible readings of the cross-package citation path coexist in shipped markdown, and nothing can tell them apart. Pick one and sweep.

## Context

Raised by REQ-244's builder as pushback rather than decided mid-REQ, verified independently by the orchestrator, and confirmed again by REQ-244's review. REQ-244 added eleven more citations in the prescribed form while the question was open, so the count grows until this is settled.

- `\`../do-work/actions/work-reference.md\`` is what `_dev/primes/prime-action-files.md:91` prescribes and what `actions/memory.md` uses. From `skills/do-work-toolbox/actions/` it resolves to `skills/do-work-toolbox/do-work/actions/...`, **which does not exist.** It coheres only read as *skill-root*-relative rather than relative to the citing file's own directory.
- `\`../../do-work/actions/work-reference.md\`` is what `skills/do-work-toolbox/actions/present-work.md:37` already uses. It resolves literally, in both the source and installed topologies.

Counts at capture, verified by REQ-244's review: 30 local-form citations in core, 3 sibling-form in do-work-board, 11 sibling-form in do-work-knowledge and do-work-toolbox — consistent by package, so whichever reading is intended, the corpus is at least internally regular.

**REQ-243's checker cannot arbitrate this.** It resolves Markdown `[text](target)` link syntax; every one of these citations is a backticked path, which that checker never sees. So the correctness of cross-package pointers currently rests on convention alone.

## Open Questions

- [x] Shipped action files cite each other across skill packages with a backticked path, and there are two spellings in the tree that mean different things. One (`../do-work/actions/...`) is what the prime file tells writers to use, and it only makes sense if you read the `../` as "up to the skills folder" rather than as a real relative path — typed literally into a terminal from the citing file's folder, it points at nothing. The other (`../../do-work/actions/...`) is a real path that works from where the file actually sits. Both are in use today, the first far more than the second. Nothing checks either, because our new link checker only understands Markdown links and these are backticks. Which reading do you want to be the rule — the skill-root one that most files already use, or the literal one that a reader could paste and follow? Whichever you pick, the other spelling gets swept to match and the prime file gets updated to say so. → Confirmed: literal paths (`../../`), so a citation a reader pastes actually resolves and a future checker can verify it mechanically.
  Recommended: literal paths (`../../`), so a citation a reader pastes actually resolves and a future checker can verify it mechanically.
  Also: keep skill-root-relative (`../`), sweep `present-work.md` and `completed-work-presentation-reference.md` to match, and state the convention explicitly in the prime so nobody reads it as a filesystem path.

**Answered [2026-08-18]:** User confirmed the recommended literal form via `do-work clarify`. The rule is: every backticked cross-package citation must resolve as a real relative path from the citing file's own directory. Sweep every citation that doesn't (the Context counts are a starting set, not the extent — re-derive them), update `_dev/primes/prime-action-files.md:91` to prescribe the literal form, and answer the Requirements question in favor of mechanical checkability: backticked cross-package pointers should become checkable now that they resolve literally.

## Requirements

- One spelling is the documented rule, stated in `_dev/primes/prime-action-files.md`.
- Every shipped citation matches it — the sweep is the requirement, and the counts above are a starting set rather than the extent.
- Whether backticked cross-package pointers become mechanically checkable is part of the answer: if the literal form wins, say whether a checker should now read them.
