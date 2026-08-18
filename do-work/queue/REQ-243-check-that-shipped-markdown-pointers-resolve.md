---
id: REQ-243
title: Check that shipped markdown pointers actually resolve
status: pending-answers
created_at: 2026-08-18T12:20:30Z
user_request: UR-042
addendum_to: REQ-238
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- _dev/tests/prescribed-shell-canonicalization.sh
---

# Check That Shipped Markdown Pointers Actually Resolve

## What

`_dev/tests/prescribed-shell-canonicalization.sh` proves a restatement is **absent** from shipped markdown. Nothing proves that the pointer which replaced it **resolves** — not the relative path, not the `#anchor`. Add that check.

## Context

Raised as a suggested-testing item by REQ-230's review, and confirmed as a real gap by REQ-238, which is the second fix of the same class. Both REQs did the identical work by hand:

- normalize the relative path from the citing file's directory and confirm the target exists
- confirm the anchor text exists as a heading in that target

Every fix of this class replaces prose with a link. It trades a staleness risk that **has** a detector for a broken-link risk that has **none**. Two instances in, the trade is systematic rather than incidental — and the manual check is mechanical enough that a human doing it a third time is the actual defect.

Neither existing check covers it. The canonicalization suite greps for phrases that must be absent; the shipped-package reference contract checks that shipped files do not cite maintainer-only paths. A pointer to a real file at a wrong relative depth passes both.

## Requirements

- Every relative `.md` link in shipped markdown resolves, from the **citing file's own directory**, to a file that exists. The citing directory is the whole point — a link correct from the repo root and wrong from its own file is exactly the failure mode.
- Where a link carries a `#anchor`, the anchor matches a heading in the target file, compared the way the anchor was generated from it (lowercased, spaces to hyphens, punctuation dropped).
- A broken link names the citing file, the line, the raw link, and what it resolved to — enough to fix without re-deriving.
- Links to things the check cannot resolve (external `http(s)://`, `mailto:`, bare anchors into the same file, paths outside the repo) are skipped explicitly rather than silently, and the skip rule is stated in the file.
- `bash _dev/tests/maintainer-verify.sh` still exits 0.

## Constraints

- **Shipped markdown only** — the same tree the canonicalization scan already walks. Archived REQs under `do-work/` are history and routinely reference paths that have since moved; scanning them would produce noise that trains readers to ignore the check.
- No new test file if an existing suite is the natural home. This is a property of shipped markdown, and `prescribed-shell-canonicalization.sh` already walks exactly that set — check whether it belongs there before adding a file.
- `maintenance: true`: this is a pass on the skill's own instructions, so ask whether anything can be **removed** before adding. In particular, check whether the shipped-package reference contract already walks these links for a different reason and could carry this one instead of a second walker.

## Red-Green Proof

**RED prompt/case:** introduce a pointer with a wrong relative depth — e.g. change one existing `../../do-work/docs/prescribed-shell-primitives.md` to `../do-work/docs/prescribed-shell-primitives.md` — and run the check.
**Why RED now:** no check reads markdown links at all, so the wrong depth passes the entire suite today.
**GREEN when:** that mutation fails the suite naming the citing file and line, and reverting it passes. A second mutation on the anchor (`#portfolio-summary-publication` → `#portfolio-summary-publications`) fails the same way.
**Validation:** Review finding on REQ-238; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] Twice today a fix has replaced a repeated paragraph with a link to the one place that paragraph lives — REQ-230 and REQ-238, and there will be more, because that is the standard fix for this class. Each time, the test suite confirms the repeated text is gone, and nothing at all confirms the link works. Both times I checked the link by hand: that the path points at a file that exists, and that the heading it names is really in that file. A wrong link would pass every check we have. The fix is a check that walks the links in shipped instruction files and resolves them. I am asking rather than doing it because it is new machinery rather than a repair — a whole new class of check, with its own decisions about what to skip (external URLs, links into archived history) and its own maintenance cost — and you may prefer to keep checking these by hand while the class is only a few instances old. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — keep checking pointers by hand until the class is big enough to justify a checker.
