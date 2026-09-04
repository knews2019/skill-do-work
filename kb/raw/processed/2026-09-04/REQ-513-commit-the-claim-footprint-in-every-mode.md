---
source_type: req_lesson
req_id: REQ-513
req_path: do-work/archive/REQ-513-commit-the-claim-footprint-in-every-mode.md
date: 2026-09-02
domain: general
module: _dev/primes
tags: [general, commit, claim, footprint]
---

# Lessons from REQ-513: Commit the claim footprint in every mode

## What the REQ was about

Make Step 2's claim commit its own footprint in every mode, serial included: the queue-to-working move plus the checkpoint entry land as one bookkeeping commit at claim time. The CLI's `claim --commit` already exists; `actions/work.md` Step 2 does not use it, and worktree dispatch Step 0 stages the same moves by hand later.

## Solution summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go` (modified)

## What worked

- Testing the action-to-CLI seam and the CLI's exact Git footprint together caught both wiring drift and transaction regressions.
- The restatement sweep found the user guide's delayed-commit description before release.

## What didn't work

- The first guide assertion passed an absolute path to a repo-relative test helper, so the canonical gate rejected the doubled path; using the helper's expected relative path fixed it.

## Worth knowing

- `claim --commit` already owned the correct transaction. The fix was to invoke that owner everywhere and delete the later hand-back staging contract, not add another commit path.

## Back-reference

See `do-work/archive/REQ-513-commit-the-claim-footprint-in-every-mode.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `33852cb4a9c0e8af197d789fea6f2624beb68ffe`.
